package packagemanager

// This package is responsible for handling updating and uninstalling snapd applications

import (
    "os/exec"
    "fmt"
)

// UpdateSnapPackages updates specified Snap packages or all if no specific package is provided
func UpdateSnapPackages(packageNames ...string) error {
    var cmd *exec.Cmd
    if len(packageNames) == 0 {
        cmd = exec.Command("sudo", "snap", "refresh")
    } else {
        args := append([]string{"refresh"}, packageNames...)
        cmd = exec.Command(args[0], args[1:]...)
    }

    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("error updating Snap packages: %s, %v", string(output), err)
    }
    return nil
}

// UninstallSnapPackage uninstalls a specified Snap package
func UninstallSnapPackage(packageName string) error {
    cmd := exec.Command("sudo", "snap", "remove", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("error uninstalling Snap package: %s, %v", string(output), err)
    }
    return nil
}
