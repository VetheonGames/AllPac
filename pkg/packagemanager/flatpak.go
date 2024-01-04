package packagemanager

// This package is responsible for handling updating and uninstalling flatpak applications

import (
    "os/exec"
    "fmt"
)

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
