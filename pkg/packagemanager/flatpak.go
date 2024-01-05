package packagemanager

// This package is responsible for handling updating and uninstalling flatpak applications

import (
    "os/exec"
    "fmt"
	"strings"
	"pixelridgesoftworks.com/AllPac/pkg/logger"
)

// UpdateFlatpakPackages updates specified Flatpak packages or all if no specific package is provided
func UpdateFlatpakPackages(packageNames ...string) error {
    // Read the current package list
    pkgList, err := ReadPackageList()
    if err != nil {
        logger.Errorf("error reading package list: %v", err)
        return fmt.Errorf("error reading package list: %v", err)
    }

    // Determine which packages need updating
    var packagesToUpdate []string
    for _, packageName := range packageNames {
        installedInfo, ok := pkgList[packageName]
        if !ok {
            logger.Infof("Package %s not managed by AllPac, skipping", packageName)
            continue
        }

        availableVersion, err := GetFlatpakPackageVersion(packageName)
        if err != nil {
            logger.Errorf("error getting available version for Flatpak package %s: %v", packageName, err)
            continue
        }

        if installedInfo.Version != availableVersion {
            packagesToUpdate = append(packagesToUpdate, packageName)
        }
    }

    // Update the packages
    if len(packagesToUpdate) > 0 {
        args := append([]string{"update", "-y"}, packagesToUpdate...)
        cmd := exec.Command("flatpak", args...)
        if output, err := cmd.CombinedOutput(); err != nil {
            logger.Errorf("error updating Flatpak packages: %s, %v", output, err)
            return fmt.Errorf("error updating Flatpak packages: %s, %v", output, err)
        }

        // Update the package list with the new versions
        for _, packageName := range packagesToUpdate {
            newVersion, err := GetFlatpakPackageVersion(packageName)
            if err != nil {
                logger.Errorf("error getting new version for Flatpak package %s after update: %v", packageName, err)
                continue
            }
            if err := UpdatePackageInList(packageName, "flatpak", newVersion); err != nil {
                logger.Errorf("error updating package list for %s: %v", packageName, err)
                return fmt.Errorf("error updating package list for %s: %v", packageName, err)
            }
        }
    } else {
        logger.Info("No Flatpak packages need updating")
    }

    return nil
}

// UninstallFlatpakPackage uninstalls a specified Flatpak package
func UninstallFlatpakPackage(packageName string) error {
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

    // Uninstalling the Flatpak package
    cmd := exec.Command("flatpak", "uninstall", "-y", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        logger.Errorf("error uninstalling Flatpak package: %s, %v", output, err)
        return fmt.Errorf("error uninstalling Flatpak package: %s, %v", output, err)
    }

    // Remove the package from the list after successful uninstallation
    if err := RemovePackageFromList(packageName); err != nil {
        logger.Errorf("An error has occurred while removing the package from the list: %v", err)
        return err
    }

    logger.Infof("Package %s successfully uninstalled and removed from the package list", packageName)
    return nil
}

// GetVersionFromFlatpak gets the installed version of a Flatpak package
func GetVersionFromFlatpak(applicationID string) (string, error) {
    cmd := exec.Command("flatpak", "info", applicationID)
    output, err := cmd.CombinedOutput()
    if err != nil {
		logger.Errorf("error getting flatpak package info: %v", err)
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
	logger.Errorf("version not found for flatpak package: %s", applicationID)
    return "", fmt.Errorf("version not found for flatpak package: %s", applicationID)
}
