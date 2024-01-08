package packagemanager

import (
	"fmt"
)

// UpdatePackageByName updates a specific package by its name
func UpdatePackageByName(packageName string) error {
    pkgList, err := ReadPackageList()
    if err != nil {
        return fmt.Errorf("error reading package list: %v", err)
    }

    pkgInfo, exists := pkgList[packageName]
    if !exists {
        return fmt.Errorf("package %s not found in package list", packageName)
    }

    switch pkgInfo.Source {
    case "pacman":
        return UpdatePacmanPackages(packageName)
    case "aur":
        return UpdateAURPackages(packageName)
    case "snap":
        return UpdateSnapPackages(packageName)
    case "flatpak":
        return UpdateFlatpakPackages(packageName)
    default:
        return fmt.Errorf("unknown source for package %s", packageName)
    }
}
