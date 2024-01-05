package packagemanager

import (
    "fmt"
	"sync"
	"pixelridgesoftworks.com/AllPac/pkg/logger"
)

// UpdateAllPackages updates all packages on the system
func UpdateAllPackages() error {
    pkgList, err := readPackageList()
    if err != nil {
		logger.Errorf("Failed to load config: %v", err)
        return fmt.Errorf("error reading package list: %v", err)
    }

    // Categorize packages by their source
    pacmanPackages, aurPackages, snapPackages, flatpakPackages := separatePackagesBySource(pkgList)

    // Check and collect packages that need updating for each category
    pacmanToUpdate := checkPackagesForUpdate(pacmanPackages, "pacman")
    snapToUpdate := checkPackagesForUpdate(snapPackages, "snap")
    flatpakToUpdate := checkPackagesForUpdate(flatpakPackages, "flatpak")

    // Perform batch updates
    if err := UpdatePacmanPackages(pacmanToUpdate...); err != nil {
        logger.Errorf("Error updating Pacman packages: %v\n", err)
    }
    if err := UpdateSnapPackages(snapToUpdate...); err != nil {
        logger.Errorf("Error updating Snap packages: %v\n", err)
    }
    if err := UpdateFlatpakPackages(flatpakToUpdate...); err != nil {
        logger.Errorf("Error updating Flatpak packages: %v\n", err)
    }

    // Update AUR packages (can be done concurrently)
    updateAURPackagesConcurrently(aurPackages)

    fmt.Println("All packages have been updated.")
	logger.Info("All packages have been updated.")
    return nil
}

// separatePackagesBySource categorizes package names by their source
func separatePackagesBySource(pkgList PackageList) ([]string, []string, []string, []string) {
    var pacmanPackages, aurPackages, snapPackages, flatpakPackages []string
    for pkgName, pkgInfo := range pkgList {
        switch pkgInfo.Source {
        case "pacman":
            pacmanPackages = append(pacmanPackages, pkgName)
        case "aur":
            aurPackages = append(aurPackages, pkgName)
        case "snap":
            snapPackages = append(snapPackages, pkgName)
        case "flatpak":
            flatpakPackages = append(flatpakPackages, pkgName)
        }
    }
    return pacmanPackages, aurPackages, snapPackages, flatpakPackages
}

// checkPackagesForUpdate checks which packages need updating and returns their names
func checkPackagesForUpdate(packageNames []string, source string) []string {
    var toUpdate []string
    for _, name := range packageNames {
        if needsUpdate, _ := checkIfPackageNeedsUpdate(name, source); needsUpdate {
            toUpdate = append(toUpdate, name)
        }
    }
    return toUpdate
}

// checkIfPackageNeedsUpdate checks if a given package needs an update
func checkIfPackageNeedsUpdate(name, source string) (bool, error) {
    var currentVersion, latestVersion string
    var err error

    // Retrieve the current version of the package from the package list
    pkgList, err := readPackageList()
    if err != nil {
		logger.Errorf("error reading package list: %v", err)
        return false, fmt.Errorf("error reading package list: %v", err)
    }
    if pkgInfo, exists := pkgList[name]; exists {
        currentVersion = pkgInfo.Version
    } else {
		logger.Errorf("package %s not found in package list", name)
        return false, fmt.Errorf("package %s not found in package list", name)
    }

    // Get the latest version based on the source
    switch source {
    case "pacman":
        latestVersion, err = GetPacmanPackageVersion(name)
    case "aur":
        latestVersion, err = GetAURPackageVersion(name)
    case "snap":
        latestVersion, err = GetSnapPackageVersion(name)
    case "flatpak":
        latestVersion, err = GetFlatpakPackageVersion(name)
    default:
		logger.Errorf("unknown package source for %s", name)
        return false, fmt.Errorf("unknown package source for %s", name)
    }

    if err != nil {
		logger.Errorf("An error has occured:", err)
        return false, err
    }

    // Compare the current version with the latest version
    return currentVersion != latestVersion, nil
}

// updateAURPackagesConcurrently updates AUR packages using concurrency
func updateAURPackagesConcurrently(packageNames []string) {
    var wg sync.WaitGroup
    for _, pkgName := range packageNames {
        wg.Add(1)
        go func(name string) {
            defer wg.Done()
            if err := UpdateAURPackages(name); err != nil {
                logger.Errorf("Error updating AUR package %s: %v\n", name, err)
            }
        }(pkgName)
    }
    wg.Wait()
}
