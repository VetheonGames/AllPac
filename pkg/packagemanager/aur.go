package packagemanager

import (
    "fmt"
    "os/exec"
	"os"
	"path/filepath"
	"pixelridgesoftworks.com/AllPac/pkg/logger"
)

// AURPackageInfo represents the package information from the AUR
type AURPackageInfo struct {
    Version string `json:"Version"`
}

// UpdateAURPackages updates specified AUR packages or all if no specific package is provided
func UpdateAURPackages(packageNames ...string) error {
    pkgList, err := ReadPackageList()
    if err != nil {
        logger.Errorf("error reading package list: %v", err)
        return fmt.Errorf("error reading package list: %v", err)
    }

    // If no specific packages are provided, update all AUR packages in the list
    if len(packageNames) == 0 {
        for packageName, pkgInfo := range pkgList {
            if pkgInfo.Source == "aur" {
                packageNames = append(packageNames, packageName)
            }
        }
    }

    for _, packageName := range packageNames {
        aurInfo, err := fetchAURPackageInfo(packageName)
        if err != nil {
            logger.Errorf("error fetching AUR package info for %s: %v", packageName, err)
            return fmt.Errorf("error fetching AUR package info for %s: %v", packageName, err)
        }

        installedInfo, ok := pkgList[packageName]
        if !ok || installedInfo.Version != aurInfo.Version {
            _, err := CloneAndInstallFromAUR("https://aur.archlinux.org/" + packageName + ".git", true)
            if err != nil {
                logger.Errorf("error updating AUR package %s: %v", packageName, err)
                return fmt.Errorf("error updating AUR package %s: %v", packageName, err)
            }

            // Update the package list with the new version
            if err := UpdatePackageInList(packageName, "aur", aurInfo.Version); err != nil {
                logger.Errorf("error updating package list for %s: %v", packageName, err)
                return fmt.Errorf("error updating package list for %s: %v", packageName, err)
            }
        }
    }
    return nil
}

// UninstallAURPackage uninstalls a specified AUR package
func UninstallAURPackage(packageName string) error {
    // Read the current package list
    pkgList, err := ReadPackageList()
    if err != nil {
        logger.Errorf("An error has occurred while reading the package list: %v", err)
        return err
    }

    // Check if the package is managed by AllPac
    if _, exists := pkgList[packageName]; !exists {
        logger.Infof("Package %s not found in the package list, may not be managed by AllPac. Skipping uninstallation.", packageName)
        return nil // Skip this package as it's not managed by AllPac
    }

    // Uninstalling an AUR package is typically done with pacman
    cmd := exec.Command("sudo", "pacman", "-Rns", "--noconfirm", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        logger.Errorf("error uninstalling AUR package: %s, %v", output, err)
        return fmt.Errorf("error uninstalling AUR package: %s, %v", output, err)
    }

    // Remove the package from the list after successful uninstallation
    if err := RemovePackageFromList(packageName); err != nil {
        logger.Errorf("An error has occurred while removing the package from the list: %v", err)
        return err
    }

    logger.Infof("Package %s successfully uninstalled and removed from the package list", packageName)
    return nil
}

// ClearAllPacCache clears the contents of the ~/.allpac/cache/ directory
func ClearAllPacCache() error {
    cacheDir, err := getCacheDir()
    if err != nil {
		logger.Errorf("An error has occured:", err)
        return err
    }

    // Remove the directory and its contents
    err = os.RemoveAll(cacheDir)
    if err != nil {
		logger.Errorf("An error has occured:", err)
        return err
    }

    // Optionally, recreate the cache directory after clearing it
    return os.MkdirAll(cacheDir, 0755)
}

// getCacheDir returns the path to the ~/.allpac/cache/ directory
func getCacheDir() (string, error) {
    userHomeDir, err := os.UserHomeDir()
    if err != nil {
		logger.Errorf("An error has occured:", err)
        return "", err
    }
    return filepath.Join(userHomeDir, ".allpac", "cache"), nil
}

// RebuildAndReinstallAURPackage rebuilds and reinstalls the specified AUR package
func RebuildAndReinstallAURPackage(packageName string) error {
    // Read the package list
    pkgList, err := readPackageList()
    if err != nil {
		logger.Errorf("error reading package list: %v", err)
        return fmt.Errorf("error reading package list: %v", err)
    }

    // Check if the package is in the list and is an AUR package
    pkgInfo, found := pkgList[packageName]
    if !found || pkgInfo.Source != "aur" {
		logger.Errorf("package %s is not found or not an AUR package", packageName)
        return fmt.Errorf("package %s is not found or not an AUR package", packageName)
    }

    // Rebuild and reinstall the package
    _, err = CloneAndInstallFromAUR(fmt.Sprintf("https://aur.archlinux.org/%s.git", packageName), false)
    logger.Errorf("An error has occured:", err)
	return err
}
