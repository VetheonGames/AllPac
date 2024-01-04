package install

// This package is responsible for handling our actual install logic. We could have probably gotten away with
// implementing this into the packagemanager package, but this seems like a better way
// because this provides a single interface for all our install functions

import (
    "encoding/json"
    "os"
    "os/exec"
    "os/user"
    "path/filepath"
    "fmt"
)

// PackageList represents the mapping of installed packages to their sources
type PackageList map[string]string

const pkgListFilename = "pkg.list"

// getPkgListPath returns the file path for the package list
func getPkgListPath() (string, error) {
    usr, err := user.Current()
    if err != nil {
        return "", fmt.Errorf("error getting current user: %v", err)
    }
    return filepath.Join(usr.HomeDir, ".allpac", pkgListFilename), nil
}

// readPackageList reads the package list from the file
func readPackageList() (PackageList, error) {
    pkgListPath, err := getPkgListPath()
    if err != nil {
        return nil, err
    }

    file, err := os.Open(pkgListPath)
    if err != nil {
        if os.IsNotExist(err) {
            return PackageList{}, nil // Return an empty list if file doesn't exist
        }
        return nil, fmt.Errorf("error opening package list file: %v", err)
    }
    defer file.Close()

    var pkgList PackageList
    err = json.NewDecoder(file).Decode(&pkgList)
    if err != nil {
        return nil, fmt.Errorf("error decoding package list: %v", err)
    }

    return pkgList, nil
}

// writePackageList writes the package list to the file
func writePackageList(pkgList PackageList) error {
    pkgListPath, err := getPkgListPath()
    if err != nil {
        return err
    }

    file, err := os.Create(pkgListPath)
    if err != nil {
        return fmt.Errorf("error creating package list file: %v", err)
    }
    defer file.Close()

    err = json.NewEncoder(file).Encode(pkgList)
    if err != nil {
        return fmt.Errorf("error encoding package list: %v", err)
    }

    return nil
}

// logInstallation logs the package installation details
func LogInstallation(packageName, source string) error {
    pkgList, err := readPackageList()
    if err != nil {
        return err
    }

    pkgList[packageName] = source

    return writePackageList(pkgList)
}

// InstallPackagePacman installs a package using Pacman
func InstallPackagePacman(packageName string) error {
    cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("error installing package with Pacman: %s, %v", output, err)
    }
    if err := LogInstallation(packageName, "pacman"); err != nil {
        return fmt.Errorf("error logging installation: %v", err)
    }
    return nil
}

// InstallPackageSnap installs a package using Snap
func InstallPackageSnap(packageName string) error {
    cmd := exec.Command("sudo", "snap", "install", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("error installing package with Snap: %s, %v", output, err)
    }
    if err := LogInstallation(packageName, "snap"); err != nil {
        return fmt.Errorf("error logging installation: %v", err)
    }
    return nil
}

// InstallPackageFlatpak installs a package using Flatpak
func InstallPackageFlatpak(packageName string) error {
    cmd := exec.Command("flatpak", "install", "-y", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("error installing package with Flatpak: %s, %v", output, err)
    }
    if err := LogInstallation(packageName, "flatpak"); err != nil {
        return fmt.Errorf("error logging installation: %v", err)
    }
    return nil
}

// cloneAndInstallFromAUR clones the given AUR repository and installs it
func CloneAndInstallFromAUR(repoURL string) error {
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

    // Determine the name of the created directory (and the package name)
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

    // Log the installation
    if err := LogInstallation(repoName, "aur"); err != nil {
        return fmt.Errorf("error logging installation: %v", err)
    }

    return nil
}

// InstallSnap installs Snap manually from the AUR
func InstallSnap() error {
    if err := CloneAndInstallFromAUR("https://aur.archlinux.org/snapd.git"); err != nil {
        return fmt.Errorf("error installing Snap: %v", err)
    }
    if err := LogInstallation("snapd", "aur"); err != nil {
        return fmt.Errorf("error logging installation: %v", err)
    }
    return nil
}

// InstallGit installs Git using Pacman
func InstallGit() error {
    if err := InstallPackagePacman("git"); err != nil {
        return fmt.Errorf("error installing Git: %v", err)
    }
    return nil
}

// InstallBaseDevel installs the base-devel group using Pacman
func InstallBaseDevel() error {
    if err := InstallPackagePacman("base-devel"); err != nil {
        return fmt.Errorf("error installing base-devel: %v", err)
    }
    return nil
}

// InstallFlatpak installs the Flatpak package using Pacman
func InstallFlatpak() error {
    if err := InstallPackagePacman("flatpak"); err != nil {
        return fmt.Errorf("error installing flatpak: %v", err)
    }
    return nil
}
