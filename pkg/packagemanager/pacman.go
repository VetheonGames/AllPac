package packagemanager

// This package is responsible for handling updating and uninstalling pacman packages

import (
	"fmt"
	"os/exec"
    "strings"
	"pixelridgesoftworks.com/AllPac/pkg/logger"
)

// updates specified Pacman packages or all if no specific package is provided
func UpdatePacmanPackages(packageNames ...string) error {
    // If no specific packages are provided, update all packages
    if len(packageNames) == 0 {
        logger.Info("No specific package names provided, updating all Pacman packages")
        cmd := exec.Command("sudo", "pacman", "-Syu", "--noconfirm")
        if output, err := cmd.CombinedOutput(); err != nil {
            logger.Errorf("error updating all Pacman packages: %s, %v", string(output), err)
            return fmt.Errorf("error updating all Pacman packages: %s, %v", string(output), err)
        }
        return nil
    }
    // Read the current package list
    pkgList, err := ReadPackageList()
    if err != nil {
        logger.Errorf("error reading package list: %v", err)
        return fmt.Errorf("error reading package list: %v", err)
    }

    var packagesToUpdate []string
    for _, packageName := range packageNames {
        installedInfo, ok := pkgList[packageName]
        if !ok {
            logger.Infof("Package %s not managed by AllPac, skipping", packageName)
            continue
        }

        latestVersion, err := GetPacmanLatestVersion(packageName)
        if err != nil {
            logger.Errorf("error getting latest version for Pacman package %s: %v", packageName, err)
            continue
        }

        if installedInfo.Version != latestVersion {
            packagesToUpdate = append(packagesToUpdate, packageName)
        }
    }

    if len(packagesToUpdate) > 0 {
        args := append([]string{"sudo", "pacman", "-S", "--noconfirm"}, packagesToUpdate...)
        cmd := exec.Command(args[0], args[1:]...)
        if output, err := cmd.CombinedOutput(); err != nil {
            logger.Errorf("error updating Pacman packages: %s, %v", string(output), err)
            return fmt.Errorf("error updating Pacman packages: %s, %v", string(output), err)
        }

        // Update the package list with the new versions
        for _, packageName := range packagesToUpdate {
            newVersion, err := GetPacmanLatestVersion(packageName)
            if err != nil {
                logger.Errorf("error getting new version for Pacman package %s after update: %v", packageName, err)
                continue
            }
            if err := UpdatePackageInList(packageName, "pacman", newVersion); err != nil {
                logger.Errorf("error updating package list for %s: %v", packageName, err)
                return fmt.Errorf("error updating package list for %s: %v", packageName, err)
            }
        }
    } else {
        logger.Info("No Pacman packages need updating")
    }

    return nil
}

// uninstalls a specified Pacman package
func UninstallPacmanPackage(packageName string) error {
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

    // Uninstalling the Pacman package
    cmd := exec.Command("sudo", "pacman", "-Rns", "--noconfirm", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        logger.Errorf("error uninstalling Pacman package: %s, %v", output, err)
        return fmt.Errorf("error uninstalling Pacman package: %s, %v", output, err)
    }

    // Remove the package from the list after successful uninstallation
    if err := RemovePackageFromList(packageName); err != nil {
        logger.Errorf("An error has occurred while removing the package from the list: %v", err)
        return err
    }

    logger.Infof("Package %s successfully uninstalled and removed from the package list", packageName)
    return nil
}

// retrieves the latest available version of a package from Pacman
func GetPacmanLatestVersion(packageName string) (string, error) {
    cmd := exec.Command("pacman", "-Si", packageName)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("error getting package info from Pacman: %v", err)
    }
    // Parse the output to find the version
    versionLine := strings.Split(string(output), "\n")[2]
    version := strings.Fields(versionLine)[2]
    return version, nil
}
