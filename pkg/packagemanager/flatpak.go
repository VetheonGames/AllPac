package packagemanager

// This package is responsible for handling updating and uninstalling flatpak applications

import (
    "os/exec"
    "fmt"
	"strings"
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

// GetVersionFromFlatpak gets the installed version of a Flatpak package
func GetVersionFromFlatpak(applicationID string) (string, error) {
    cmd := exec.Command("flatpak", "info", applicationID)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("error getting flatpak package info: %v", err)
    }

    lines := strings.Split(string(output), "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "Version:") {
            parts := strings.Fields(line)
            if len(parts) >= 2 {
                return parts[1], nil
            }
            break
        }
    }
    return "", fmt.Errorf("version not found for flatpak package: %s", applicationID)
}
