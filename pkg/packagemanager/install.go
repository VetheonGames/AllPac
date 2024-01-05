package packagemanager

// This package is responsible for handling our actual install logic. We could have probably gotten away with
// implementing this into the packagemanager package, but this seems like a better way
// because this provides a single interface for all our install functions

import (
    "fmt"
    "os"
    "os/exec"
    "os/user"
    "path/filepath"
    "pixelridgesoftworks.com/AllPac/pkg/logger"
)

// InstallPackagePacman installs a package using Pacman and logs the installation
func InstallPackagePacman(packageName string) error {
    cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        logger.Errorf("error installing package with Pacman: %s, %v", output, err)
        return fmt.Errorf("error installing package with Pacman: %s, %v", output, err)
    }

    version, err := GetPacmanPackageVersion(packageName)
    if err != nil {
        logger.Errorf("An error has occured:", err)
        return err
    }

    if err := LogInstallation(packageName, "pacman", version); err != nil {
        logger.Errorf("error logging installation: %v", err)
        return fmt.Errorf("error logging installation: %v", err)
    }
    return nil
}

// InstallPackageSnap installs a package using Snap and logs the installation
func InstallPackageSnap(packageName string) error {
    cmd := exec.Command("sudo", "snap", "install", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        logger.Errorf("error installing package with Snap: %s, %v", output, err)
        return fmt.Errorf("error installing package with Snap: %s, %v", output, err)
    }

    version, err := GetVersionFromSnap(packageName)
    if err != nil {
        logger.Errorf("An error has occured:", err)
        return err
    }

    if err := LogInstallation(packageName, "snap", version); err != nil {
        logger.Errorf("error logging installation: %v", err)
        return fmt.Errorf("error logging installation: %v", err)
    }
    return nil
}

// InstallPackageFlatpak installs a package using Flatpak and logs the installation
func InstallPackageFlatpak(packageName string) error {
    cmd := exec.Command("flatpak", "install", "-y", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        logger.Errorf("error installing package with Flatpak: %s, %v", output, err)
        return fmt.Errorf("error installing package with Flatpak: %s, %v", output, err)
    }

    version, err := GetVersionFromFlatpak(packageName)
    if err != nil {
        logger.Errorf("An error has occured:", err)
        return err
    }

    if err := LogInstallation(packageName, "flatpak", version); err != nil {
        logger.Errorf("error logging installation: %v", err)
        return fmt.Errorf("error logging installation: %v", err)
    }
    return nil
}

// cloneAndInstallFromAUR clones the given AUR repository and installs it
func CloneAndInstallFromAUR(repoURL string, skipConfirmation bool) (string, error) {
    // Request root permissions
    if !skipConfirmation && !requestRootPermissions() {
        logger.Warnf("root permissions denied")
        return "", fmt.Errorf("root permissions denied")
    }

    // Confirm before proceeding with each step
    if !skipConfirmation && !confirmAction("Do you want to download and build package from " + repoURL + "?") {
        logger.Warnf("user aborted the action")
        return "", fmt.Errorf("user aborted the action")
    }
    // Get the current user's home directory
    usr, err := user.Current()
    if err != nil {
        logger.Errorf("error getting current user: %v", err)
        return "", fmt.Errorf("error getting current user: %v", err)
    }

    // Define the base directory for AllPac cache
    baseDir := filepath.Join(usr.HomeDir, ".allpac", "cache")

    // Ensure the base directory exists
    if err := os.MkdirAll(baseDir, 0755); err != nil {
        logger.Errorf("error creating base directory: %v", err)
        return "", fmt.Errorf("error creating base directory: %v", err)
    }

    // Clone the repository
    cmdGitClone := exec.Command("git", "clone", repoURL, baseDir)
    if output, err := cmdGitClone.CombinedOutput(); err != nil {
        logger.Errorf("error cloning AUR repo: %s, %v", output, err)
        return "", fmt.Errorf("error cloning AUR repo: %s, %v", output, err)
    }

    // Determine the name of the created directory (and the package name)
    repoName := filepath.Base(repoURL)
    repoDir := filepath.Join(baseDir, repoName)

    // Change directory to the cloned repository
    if err := os.Chdir(repoDir); err != nil {
        logger.Errorf("error changing directory: %v", err)
        return "", fmt.Errorf("error changing directory: %v", err)
    }

    // Build the package using makepkg
    cmdMakePkg := exec.Command("makepkg", "-si", "--noconfirm")
    if output, err := cmdMakePkg.CombinedOutput(); err != nil {
        logger.Errorf("error building package with makepkg: %s, %v", output, err)
        return "", fmt.Errorf("error building package with makepkg: %s, %v", output, err)
    }

    // Extract the version from PKGBUILD
    version, err := ExtractVersionFromPKGBUILD(repoDir)
    if err != nil {
        logger.Errorf("error extracting version from PKGBUILD: %v", err)
        return "", fmt.Errorf("error extracting version from PKGBUILD: %v", err)
    }

    // Confirm before installing
    if !skipConfirmation && !confirmAction("Do you want to install the built package " + repoName + "?") {
        logger.Warnf("user aborted the installation")
        return "", fmt.Errorf("user aborted the installation")
    }

    if err := LogInstallation(repoName, "aur", version); err != nil {
        logger.Errorf("error logging installation")
        return "", fmt.Errorf("error logging installation: %v", err)
    }

    return version, nil
}

// InstallSnap installs Snap manually from the AUR
func InstallSnap() error {
    version, err := CloneAndInstallFromAUR("https://aur.archlinux.org/snapd.git", true)
    if err != nil {
        logger.Errorf("error installing Snap: %v", err)
        return fmt.Errorf("error installing Snap: %v", err)
    }

    if err := LogInstallation("snapd", "aur", version); err != nil {
        logger.Errorf("error logging installation")
        return fmt.Errorf("error logging installation: %v", err)
    }
    return nil
}

// InstallGit installs Git using Pacman
func InstallGit() error {
    if err := InstallPackagePacman("git"); err != nil {
        logger.Errorf("error installing Git: %v", err)
        return fmt.Errorf("error installing Git: %v", err)
    }
    return nil
}

// InstallBaseDevel installs the base-devel group using Pacman
func InstallBaseDevel() error {
    if err := InstallPackagePacman("base-devel"); err != nil {
        logger.Errorf("error installing base-devel: %v", err)
        return fmt.Errorf("error installing base-devel: %v", err)
    }
    return nil
}

// InstallFlatpak installs the Flatpak package using Pacman
func InstallFlatpak() error {
    if err := InstallPackagePacman("flatpak"); err != nil {
        logger.Errorf("error installing flatpak: %v", err)
        return fmt.Errorf("error installing flatpak: %v", err)
    }
    return nil
}
