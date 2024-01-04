package packagemanager

// This package is responsible for handling updating and uninstalling snapd applications

import (
    "os/exec"
    "fmt"
	"strings"
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

// GetVersionFromSnap gets the installed version of a Snap package
func GetVersionFromSnap(packageName string) (string, error) {
    cmd := exec.Command("snap", "info", packageName)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("error getting snap package info: %v", err)
    }

    lines := strings.Split(string(output), "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "installed:") {
            parts := strings.Fields(line)
            if len(parts) >= 2 {
                return parts[1], nil
            }
            break
        }
    }
    return "", fmt.Errorf("version not found for snap package: %s", packageName)
}
