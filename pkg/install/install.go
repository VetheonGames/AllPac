package install

import (
    "os"
    "os/exec"
    "os/user"
    "path/filepath"
    "fmt"
)

// InstallPackagePacman installs a package using Pacman
func InstallPackagePacman(packageName string) error {
    cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("error installing package with Pacman: %s, %v", output, err)
    }
    return nil
}

// InstallPackageYay installs a package using Yay (AUR)
func InstallPackageYay(packageName string) error {
    cmd := exec.Command("yay", "-S", "--noconfirm", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("error installing package with Yay: %s, %v", output, err)
    }
    return nil
}

// InstallPackageSnap installs a package using Snap
func InstallPackageSnap(packageName string) error {
    cmd := exec.Command("sudo", "snap", "install", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("error installing package with Snap: %s, %v", output, err)
    }
    return nil
}

// InstallPackageFlatpak installs a package using Flatpak
func InstallPackageFlatpak(packageName string) error {
    cmd := exec.Command("flatpak", "install", "-y", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("error installing package with Flatpak: %s, %v", output, err)
    }
    return nil
}

// InstallSnap installs Snap manually from the AUR
func InstallSnap() error {
    if err := cloneAndInstallFromAUR("https://aur.archlinux.org/snapd.git"); err != nil {
        return fmt.Errorf("error installing Snap: %v", err)
    }
    return nil
}

// cloneAndInstallFromAUR clones the given AUR repository and installs it
func cloneAndInstallFromAUR(repoURL string) error {
    // Get the current user's home directory
    usr, err := user.Current()
    if err != nil {
        return fmt.Errorf("error getting current user: %v", err)
    }

    // Define the base directory for AllPac cache
    baseDir := filepath.Join(usr.HomeDir, ".allpac", "cache")

    // Ensure the base directory exists
    if err := os.MkdirAll(baseDir, 0755); err != nil {
        return fmt.Errorf("error creating base directory: %v", err)
    }

    // Clone the repository
    cmdGitClone := exec.Command("git", "clone", repoURL, baseDir)
    if output, err := cmdGitClone.CombinedOutput(); err != nil {
        return fmt.Errorf("error cloning AUR repo: %s, %v", output, err)
    }

    // Determine the name of the created directory
    repoName := filepath.Base(repoURL)
    repoDir := filepath.Join(baseDir, repoName)

    // Change directory to the cloned repository
    if err := os.Chdir(repoDir); err != nil {
        return fmt.Errorf("error changing directory: %v", err)
    }

    // Build the package using makepkg
    cmdMakePkg := exec.Command("makepkg", "-si", "--noconfirm")
    if output, err := cmdMakePkg.CombinedOutput(); err != nil {
        return fmt.Errorf("error building package with makepkg: %s, %v", output, err)
    }

    return nil
}
