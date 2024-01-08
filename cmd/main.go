package main

// This file is our main entrypoint, and build point for AllPac

import (
    "fmt"
    "os"
	"strings"
    "pixelridgesoftworks.com/AllPac/pkg/packagemanager"
	"pixelridgesoftworks.com/AllPac/pkg/logger"
    "pixelridgesoftworks.com/AllPac/pkg/toolcheck"
	"path/filepath"
    "regexp"
)

func main() {
    // Initialize the logger
    logFilePath := filepath.Join(os.Getenv("HOME"), ".allpac", "logs", "allpac.log")
    if err := logger.Init(logFilePath); err != nil {
        logger.Errorf("Failed to initialize logger: %v", err)
    }

    if len(os.Args) < 2 {
        fmt.Println("Expected 'update', 'install', 'uninstall', 'search', 'rebuild', 'clean-aur', or 'toolcheck' subcommands")
        os.Exit(1)
    }

    command := os.Args[1]
    args := os.Args[2:]

    switch command {
    case "update":
        handleUpdate(args)
    case "install":
        handleInstall(args)
    case "uninstall":
        handleUninstall(args)
    case "search":
        handleSearch(args)
    case "rebuild":
        handleRebuild(args)
    case "clean-aur":
        handleCleanAur(args)
    case "toolcheck":
        handleToolCheck(args)
    default:
        fmt.Printf("Unknown subcommand: %s\n", command)
        os.Exit(1)
    }
}

func handleUpdate(args []string) {
    if len(args) == 0 {
        fmt.Println("You must specify an update option: 'everything', 'snap', 'aur', 'arch', 'flats', or a specific package name.")
        return
    }

    updateFuncs := map[string]func() error{
        "everything": packagemanager.UpdateAllPackages,
        "snap":       func() error { return packagemanager.UpdateSnapPackages() },
        "aur":        func() error { return packagemanager.UpdateAURPackages() },
        "arch":       func() error { return packagemanager.UpdatePacmanPackages() },
        "flats":      func() error { return packagemanager.UpdateFlatpakPackages() },
    }

    updateOption := args[0]
    if updateFunc, ok := updateFuncs[updateOption]; ok {
        err := updateFunc()
        handleUpdateError(updateOption, err)
    } else {
        err := packagemanager.UpdatePackageByName(updateOption)
        handleUpdateError(updateOption, err)
    }
}

// handles the install command for packages
func handleInstall(args []string) {
    if len(args) == 0 {
        fmt.Println("You must specify at least one package name.")
        return
    }

    // Join the args to form a single string, then split by comma
    packagesInput := strings.Join(args, " ")
    packageNames := strings.Split(packagesInput, ",")

    // Trim whitespace from each package name
    for i, pkg := range packageNames {
        packageNames[i] = strings.TrimSpace(pkg)
    }

    searchResults, err := packagemanager.SearchAllSources(packageNames)
    if err != nil {
        fmt.Printf("Error searching for packages: %v\n", err)
        return
    }

    installFuncs := map[string]func(string) error{
        "Pacman": packagemanager.InstallPackagePacman,
        "Snap":   packagemanager.InstallPackageSnap,
        "Flatpak": packagemanager.InstallPackageFlatpak,
        "AUR": func(pkgName string) error {
            _, err := packagemanager.CloneAndInstallFromAUR(fmt.Sprintf("https://aur.archlinux.org/%s.git", pkgName), false)
            return err
        },
    }

    for _, result := range searchResults {
        fmt.Printf("Searching for package: %s\n", result.PackageName)
        exactMatches := filterExactMatches(result.PackageName, result.Results)

        if len(exactMatches) == 0 {
            fmt.Println("No exact matches found for package.")
            continue
        }

        selectedSource := getSelectedSource(exactMatches)
        if selectedSource == "" {
            continue
        }

        // Display the actual package name that will be installed
        fmt.Printf("Available package(s) for installation from %s:\n", selectedSource)
        for _, match := range exactMatches {
            if match.Source == selectedSource {
                for _, pkg := range match.Results {
                    fmt.Println(pkg)
                }
            }
        }

        fmt.Printf("Installing %s from %s...\n", result.PackageName, selectedSource)
        if installFunc, ok := installFuncs[selectedSource]; ok {
            if err := installFunc(result.PackageName); err != nil {
                fmt.Printf("Error installing package %s from %s: %v\n", result.PackageName, selectedSource, err)
            } else {
                fmt.Printf("Package %s installed successfully from %s.\n", result.PackageName, selectedSource)
            }
        } else {
            fmt.Printf("Unknown source for package %s\n", result.PackageName)
        }
    }
}

func getSelectedSource(exactMatches []packagemanager.SourceResult) string {
    if len(exactMatches) == 1 {
        return exactMatches[0].Source
    }

    sourceIndex := promptUserForSource(exactMatches)
    if sourceIndex < 0 || sourceIndex >= len(exactMatches) {
        fmt.Println("Invalid selection. Skipping package.")
        return ""
    }
    return exactMatches[sourceIndex].Source
}

// filters the search results to include only those with an exact match
func filterExactMatches(packageName string, sourceResults []packagemanager.SourceResult) []packagemanager.SourceResult {
    var exactMatches []packagemanager.SourceResult
    for _, sourceResult := range sourceResults {
        var filteredResults []string
        for _, result := range sourceResult.Results {
            if isExactMatch(packageName, result) {
                filteredResults = append(filteredResults, result)
            }
        }
        if len(filteredResults) > 0 {
            exactMatches = append(exactMatches, packagemanager.SourceResult{Source: sourceResult.Source, Results: filteredResults})
        }
    }
    return exactMatches
}

// checks if the given result string is an exact match for the package name
func isExactMatch(packageName, result string) bool {
    pattern := fmt.Sprintf("^%s(?:-\\d+|\\-dev)?(?: - [^ ]+)?", regexp.QuoteMeta(packageName))
    matched, _ := regexp.MatchString(pattern, result)
    return matched
}

// handles the uninstall command for packages
func handleUninstall(args []string) {
    if len(args) == 0 {
        fmt.Println("You must specify at least one package name.")
        return
    }

    // Join the args to form a single string, then split by comma
    packagesInput := strings.Join(args, " ")
    packageNames := strings.Split(packagesInput, ",")

    // Trim whitespace from each package name
    for i, pkg := range packageNames {
        packageNames[i] = strings.TrimSpace(pkg)
    }

    // Call the function to uninstall the packages
    err := packagemanager.UninstallPackages(packageNames)
    if err != nil {
        fmt.Printf("Error uninstalling packages: %v\n", err)
    } else {
        fmt.Println("Requested packages uninstalled successfully.")
    }
}

// handles the search command for packages across different package managers
func handleSearch(args []string) {
    if len(args) < 1 {
        fmt.Println("You must specify a package name.")
        return
    }

    packageName := args[0]

    // Search across all sources
    searchResults, err := packagemanager.SearchAllSources([]string{packageName})
    if err != nil {
        logger.Errorf("Error searching for package %s: %v", packageName, err)
        return
    }

    if len(searchResults) == 0 {
        fmt.Println("No results found for package:", packageName)
        return
    }

    // Iterate over the search results and print them
    for _, result := range searchResults {
        if len(result.Results) == 0 {
            fmt.Printf("%s: No results found\n", result.PackageName)
            continue
        }

        for _, sourceResult := range result.Results {
            fmt.Printf("%s Results from %s:\n", result.PackageName, sourceResult.Source)
            for _, res := range sourceResult.Results {
                fmt.Println(res)
            }
        }
    }
}

// handles the rebuild command for an AUR package
func handleRebuild(args []string) {
    if len(args) == 0 {
        fmt.Println("You must specify the name of an AUR package to rebuild.")
        return
    }

    packageName := args[0]

    pkgList, err := packagemanager.ReadPackageList()
    if err != nil {
        fmt.Printf("Error reading package list: %v\n", err)
        return
    }

    if _, exists := pkgList[packageName]; !exists {
        fmt.Printf("Package %s is not managed by AllPac or not installed.\n", packageName)
        return
    }

    cacheDir := filepath.Join(os.Getenv("HOME"), ".allpac", "cache", packageName)
    if err := os.RemoveAll(cacheDir); err != nil {
        fmt.Printf("Error removing old build directory: %v\n", err)
        return
    }

    err = packagemanager.RebuildAndReinstallAURPackage(packageName)
    if err != nil {
        fmt.Printf("Error rebuilding package %s: %v\n", packageName, err)
    } else {
        fmt.Printf("Package %s rebuilt and reinstalled successfully.\n", packageName)
    }
}

// handles the cleaning of AUR cache
func handleCleanAur(args []string) {
    // Call the function to clear the AUR cache
    err := packagemanager.ClearAllPacCache()
    if err != nil {
        fmt.Printf("Error clearing AllPac cache: %v\n", err)
        return
    }

    fmt.Println("AllPac cache cleared successfully.")
}

// prompts the user to select a source for installation
func promptUserForSource(sources []packagemanager.SourceResult) int {
    for i, source := range sources {
        fmt.Printf("%d: %s\n", i, source.Source)
    }
    fmt.Print("Select the source number to install from: ")
    var choice int
    fmt.Scan(&choice)
    return choice
}

func handleToolCheck(args []string) {
    checks := []struct {
        Name string
        Func func() error
    }{
        {"Pacman", toolcheck.EnsurePacman},
        {"Base-devel", toolcheck.EnsureBaseDevel},
        {"Git", toolcheck.EnsureGit},
        {"Snap", toolcheck.EnsureSnap},
        {"Flatpak", toolcheck.EnsureFlatpak},
    }

    for _, check := range checks {
        if err := check.Func(); err != nil {
            fmt.Printf("%s check failed: %v\n", check.Name, err)
        } else {
            fmt.Printf("%s is installed and available.\n", check.Name)
        }
    }
}
