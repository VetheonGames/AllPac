package packagemanager

// This package is responsible for handling updating and uninstalling flatpak applications

import (
    "os/exec"
    "encoding/json"
    "io/ioutil"
    "fmt"
    "path/filepath"
    "os/user"
)

// PackageList represents the mapping of installed packages to their sources
type PackageList map[string]string

// getPkgListPath returns the file path for the package list
func getPkgListPath() (string, error) {
    usr, err := user.Current()
    if err != nil {
        return "", fmt.Errorf("error getting current user: %v", err)
    }
    return filepath.Join(usr.HomeDir, ".allpac", "pkg.list"), nil
}

// readPackageList reads the package list from the file
func readPackageList() (PackageList, error) {
    pkgListPath, err := getPkgListPath()
    if err != nil {
        return nil, err
    }

    file, err := ioutil.ReadFile(pkgListPath)
    if err != nil {
        return nil, fmt.Errorf("error reading package list file: %v", err)
    }

    var pkgList PackageList
    err = json.Unmarshal(file, &pkgList)
    if err != nil {
        return nil, fmt.Errorf("error decoding package list: %v", err)
    }

    return pkgList, nil
}

// UpdateFlatpakPackages updates specified Flatpak packages or all if no specific package is provided
func UpdateFlatpakPackages(packageNames ...string) error {
    var cmd *exec.Cmd
    if len(packageNames) == 0 {
        cmd = exec.Command("flatpak", "update", "-y")
    } else {
        args := append([]string{"update", "-y"}, packageNames...)
        cmd = exec.Command("flatpak", args...)
    }

    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("error updating Flatpak packages: %s, %v", output, err)
    }
    return nil
}

// UninstallFlatpakPackage uninstalls a specified Flatpak package
func UninstallFlatpakPackage(packageName string) error {
    cmd := exec.Command("flatpak", "uninstall", "-y", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("error uninstalling Flatpak package: %s, %v", output, err)
    }
    return nil
}
