package packagemanager

import (
    "fmt"
    "os/exec"
)

// AURPackageInfo represents the package information from the AUR
type AURPackageInfo struct {
    Version string `json:"Version"`
    // Add other relevant fields
}

// UpdateAURPackages updates specified AUR packages or all if no specific package is provided
func UpdateAURPackages(packageNames ...string) error {
    pkgList, err := ReadPackageList()
    if err != nil {
        return fmt.Errorf("error reading package list: %v", err)
    }

    for _, packageName := range packageNames {
        aurInfo, err := fetchAURPackageInfo(packageName)
        if err != nil {
            return fmt.Errorf("error fetching AUR package info for %s: %v", packageName, err)
        }

        installedInfo, ok := pkgList[packageName]
        if !ok || installedInfo.Version != aurInfo.Version {
            _, err := CloneAndInstallFromAUR("https://aur.archlinux.org/" + packageName + ".git", true)
            if err != nil {
                return fmt.Errorf("error updating AUR package %s: %v", packageName, err)
            }
        }
    }
    return nil
}

// UninstallAURPackage uninstalls a specified AUR package
func UninstallAURPackage(packageName string) error {
    // Uninstalling an AUR package is typically done with pacman
    cmd := exec.Command("sudo", "pacman", "-Rns", "--noconfirm", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("error uninstalling AUR package: %s, %v", output, err)
    }
    return nil
}
