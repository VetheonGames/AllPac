package packagemanager

// This package is responsible for handling our actual install logic. We could have probably gotten away with
// implementing this into the packagemanager package, but this seems like a better way
// because this provides a single interface for all our install functions

import (
    "fmt"
    "time"
    "os"
    "os/exec"
    "os/user"
    "strings"
    "path/filepath"
    "pixelridgesoftworks.com/AllPac/pkg/logger"
)

// installs a package using Pacman and logs the installation
func InstallPackagePacman(packageName string) error {
    cmd := exec.Command("sudo", "pacman", "-Syu", "--noconfirm", packageName)
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

// installs a package using Snap and logs the installation
func InstallPackageSnap(packageName string) error {
    cmd := exec.Command("sudo", "snap", "install", packageName)
    output, err := cmd.CombinedOutput()

    if err != nil {
        outputStr := string(output)
        logger.Errorf("error installing package with Snap: %s, %v", outputStr, err)

        // Check if the error is due to the need for classic confinement
        if strings.Contains(outputStr, "using classic") {
            fmt.Println("This package requires installation in classic mode, which may perform arbitrary system changes outside of the security sandbox. Do you want to proceed? (yes/no)")
            var response string
            fmt.Scanln(&response)
            if strings.ToLower(response) == "yes" {
                // Retry installation with --classic flag
                classicCmd := exec.Command("sudo", "snap", "install", "--classic", packageName)
                if classicOutput, classicErr := classicCmd.CombinedOutput(); classicErr != nil {
                    logger.Errorf("error installing package with Snap in classic mode: %s, %v", classicOutput, classicErr)
                    return fmt.Errorf("error installing package with Snap in classic mode: %s, %v", classicOutput, classicErr)
                }
            } else {
                return fmt.Errorf("installation aborted by user")
            }
        } else {
            return fmt.Errorf("error installing package with Snap: %s, %v", outputStr, err)
        }
    }

    version, err := GetVersionFromSnap(packageName)
    if err != nil {
        logger.Errorf("An error has occurred:", err)
        return err
    }

    if err := LogInstallation(packageName, "snap", version); err != nil {
        logger.Errorf("error logging installation: %v", err)
        return fmt.Errorf("error logging installation: %v", err)
    }
    return nil
}

// installs a package using Flatpak and logs the installation
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

// clones the given AUR repository and installs it
func CloneAndInstallFromAUR(repoURL string, skipConfirmation bool) (string, error) {
    // System update
    if !skipConfirmation && !confirmAction("Do you want to update the system before proceeding? (skipping this step may result in partial updates, and break your system)") {
        logger.Warnf("user aborted the system update")
        return "", fmt.Errorf("user aborted the system update")
    }

    cmdUpdate := exec.Command("sudo", "pacman", "-Syu", "--noconfirm")
    if output, err := cmdUpdate.CombinedOutput(); err != nil {
        logger.Errorf("error updating system: %s, %v", output, err)
        return "", fmt.Errorf("error updating system: %s, %v", output, err)
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

    // Determine the name of the package from the repo URL
    repoName := filepath.Base(repoURL)

    // Remove .git suffix
    repoName = strings.TrimSuffix(repoName, ".git")

    // Get the current date in YYYYMMDD format
    currentDate := time.Now().Format("20060102")

    // Define the directory for this specific package clone
    cloneDir := filepath.Join(usr.HomeDir, ".allpac", "cache", repoName+"-"+currentDate)

    // Ensure the clone directory exists
    if err := os.MkdirAll(cloneDir, 0755); err != nil {
        logger.Errorf("error creating clone directory: %v", err)
        return "", fmt.Errorf("error creating clone directory: %v", err)
    }

    // Define the base directory for AllPac cache
    baseDir := filepath.Join(usr.HomeDir, ".allpac", "cache")

    // Ensure the base directory exists
    if err := os.MkdirAll(baseDir, 0755); err != nil {
        logger.Errorf("error creating base directory: %v", err)
        return "", fmt.Errorf("error creating base directory: %v", err)
    }

    // Clone the repository
    cmdGitClone := exec.Command("git", "clone", repoURL, cloneDir)
    if output, err := cmdGitClone.CombinedOutput(); err != nil {
        logger.Errorf("error cloning AUR repo: %s, %v", output, err)
        return "", fmt.Errorf("error cloning AUR repo: %s, %v", output, err)
    }

    // Change directory to the cloned repository
    if err := os.Chdir(cloneDir); err != nil {
        logger.Errorf("error changing directory: %v", err)
        return "", fmt.Errorf("error changing directory: %v", err)
    }

    // Get the username of the user who invoked sudo
    sudoUser := os.Getenv("SUDO_USER")
    if sudoUser == "" {
        logger.Errorf("cannot determine the non-root user to run makepkg")
        return "", fmt.Errorf("cannot determine the non-root user to run makepkg")
    }

    // Build the package using makepkg as the non-root user
    cmdMakePkg := exec.Command("makepkg", "-si", "--noconfirm")
    cmdMakePkg.Env = []string{"HOME=" + usr.HomeDir, "USER=" + usr.Username, "LOGNAME=" + usr.Username}
    if output, err := cmdMakePkg.CombinedOutput(); err != nil {
        logger.Errorf("error building package with makepkg: %s, %v", output, err)
        return "", fmt.Errorf("error building package with makepkg: %s, %v", output, err)
    }

    // Extract the version from PKGBUILD
    version, err := ExtractVersionFromPKGBUILD(cloneDir)
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

// installs Snap manually from the AUR
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

// installs Git using Pacman
func InstallGit() error {
    if err := InstallPackagePacman("git"); err != nil {
        logger.Errorf("error installing Git: %v", err)
        return fmt.Errorf("error installing Git: %v", err)
    }
    return nil
}

// installs the base-devel group using Pacman
func InstallBaseDevel() error {
    if err := InstallPackagePacman("base-devel"); err != nil {
        logger.Errorf("error installing base-devel: %v", err)
        return fmt.Errorf("error installing base-devel: %v", err)
    }
    return nil
}

// installs the Flatpak package using Pacman
func InstallFlatpak() error {
    if err := InstallPackagePacman("flatpak"); err != nil {
        logger.Errorf("error installing flatpak: %v", err)
        return fmt.Errorf("error installing flatpak: %v", err)
    }
    return nil
}
