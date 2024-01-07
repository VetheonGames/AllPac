package packagemanager

// This package is responsible for handling updating and uninstalling snapd applications

import (
    "os/exec"
    "fmt"
	"strings"
	"pixelridgesoftworks.com/AllPac/pkg/logger"
)

// UpdateSnapPackages updates specified Snap packages or all if no specific package is provided
func UpdateSnapPackages(packageNames ...string) error {
    // Read the current package list
    pkgList, err := ReadPackageList()
    if err != nil {
        logger.Errorf("error reading package list: %v", err)
        return fmt.Errorf("error reading package list: %v", err)
    }

    // If no specific packages are provided, update all Snap packages in the list
    if len(packageNames) == 0 {
        for packageName, pkgInfo := range pkgList {
            if pkgInfo.Source == "snap" {
                packageNames = append(packageNames, packageName)
            }
        }
    }

    var cmd *exec.Cmd
    if len(packageNames) == 0 {
        cmd = exec.Command("sudo", "snap", "refresh")
    } else {
        args := append([]string{"snap", "refresh"}, packageNames...)
        cmd = exec.Command("sudo", args...)
    }

    if output, err := cmd.CombinedOutput(); err != nil {
        logger.Errorf("error updating Snap packages: %s, %v", string(output), err)
        return fmt.Errorf("error updating Snap packages: %s, %v", string(output), err)
    }

    // Update the package list with the new versions
    for _, packageName := range packageNames {
        newVersion, err := GetVersionFromSnap(packageName)
        if err != nil {
            logger.Errorf("error getting new version for Snap package %s after update: %v", packageName, err)
            continue
        }
        if err := UpdatePackageInList(packageName, "snap", newVersion); err != nil {
            logger.Errorf("error updating package list for %s: %v", packageName, err)
            return fmt.Errorf("error updating package list for %s: %v", packageName, err)
        }
    }

    return nil
}

// UninstallSnapPackage uninstalls a specified Snap package
func UninstallSnapPackage(packageName string) error {
    // Read the current package list
    pkgList, err := ReadPackageList()
    if err != nil {
        logger.Errorf("An error has occurred while reading the package list: %v", err)
        return err
    }

    // Check if the package is managed by AllPac
    if _, exists := pkgList[packageName]; !exists {
        logger.Infof("Package %s not found in the package list, may not be managed by AllPac", packageName)
        return nil
    }

    // Uninstalling the Snap package
    cmd := exec.Command("sudo", "snap", "remove", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        logger.Errorf("error uninstalling Snap package: %s, %v", string(output), err)
        return fmt.Errorf("error uninstalling Snap package: %s, %v", string(output), err)
    }

    // Remove the package from the list after successful uninstallation
    if err := RemovePackageFromList(packageName); err != nil {
        logger.Errorf("An error has occurred while removing the package from the list: %v", err)
        return err
    }

    logger.Infof("Package %s successfully uninstalled and removed from the package list", packageName)
    return nil
}

// GetVersionFromSnap gets the installed version of a Snap package
func GetVersionFromSnap(packageName string) (string, error) {
    cmd := exec.Command("snap", "info", packageName)
    output, err := cmd.CombinedOutput()
    if err != nil {
		logger.Errorf("error getting snap package info: %v", err)
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
	logger.Errorf("version not found for snap package: %s", packageName)
    return "", fmt.Errorf("version not found for snap package: %s", packageName)
}
