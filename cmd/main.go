package main

// This file is our main entrypoint, and build point for AllPac

import (
    "flag"
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

    // Define flag sets for different commands
    commandHandlers := map[string]func(*flag.FlagSet){
        "update":     handleUpdate,
        "install":    handleInstall,
        "uninstall":  handleUninstall,
        "search":     func(cmd *flag.FlagSet) { handleSearch(cmd, os.Args[2:]) },
        "rebuild":    handleRebuild,
        "clean-aur":  handleCleanAur,
        "toolcheck":  handleToolCheck,
    }

    if len(os.Args) < 2 {
        fmt.Println("Expected 'update', 'install', 'uninstall', 'search', 'rebuild', 'clean-aur', or 'toolcheck' subcommands")
        os.Exit(1)
    }

    if handler, ok := commandHandlers[os.Args[1]]; ok {
        cmd := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
        handler(cmd)
    } else {
        fmt.Printf("Unknown subcommand: %s\n", os.Args[1])
        os.Exit(1)
    }
}

func handleUpdate(cmd *flag.FlagSet) {
    updateFlags := map[string]*bool{
        "everything": cmd.Bool("everything", false, "Update all packages on the system"),
        "snap":       cmd.Bool("snap", false, "Update all Snap packages"),
        "aur":        cmd.Bool("aur", false, "Update all AUR packages"),
        "arch":       cmd.Bool("arch", false, "Update all Arch packages"),
        "flats":      cmd.Bool("flats", false, "Update all Flatpak packages"),
    }
    cmd.Parse(os.Args[2:])

    updateFuncs := map[string]func() error{
        "everything": packagemanager.UpdateAllPackages,
        "snap":       func() error { return packagemanager.UpdateSnapPackages() },
        "aur":        func() error { return packagemanager.UpdateAURPackages() },
        "arch":       func() error { return packagemanager.UpdatePacmanPackages() },
        "flats":      func() error { return packagemanager.UpdateFlatpakPackages() },
    }

    for flagName, flagValue := range updateFlags {
        if *flagValue {
            if updateFunc, ok := updateFuncs[flagName]; ok {
                if err := updateFunc(); err != nil {
                    fmt.Printf("Error occurred during '%s' update: %v\n", flagName, err)
                }
                return
            }
        }
    }

    fmt.Println("No update option specified or unrecognized option")
}

// handles the install command for packages
func handleInstall(cmd *flag.FlagSet) {
    packageNames := cmd.String("packages", "", "Comma-separated list of packages to install")
    cmd.Parse(os.Args[2:])

    if *packageNames == "" {
        fmt.Println("You must specify at least one package name.")
        cmd.Usage()
        return
    }

    packages := strings.Split(*packageNames, ",")
    searchResults, err := packagemanager.SearchAllSources(packages)
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
func handleUninstall(cmd *flag.FlagSet) {
    packageNames := cmd.String("packages", "", "Comma-separated list of packages to uninstall")

    cmd.Parse(os.Args[2:])

    if *packageNames == "" {
        fmt.Println("You must specify at least one package name.")
        cmd.Usage()
        return
    }

    packages := strings.Split(*packageNames, ",")

    // Call the function to uninstall the packages
    err := packagemanager.UninstallPackages(packages)
    if err != nil {
        fmt.Printf("Error uninstalling packages: %v\n", err)
    } else {
        fmt.Println("Requested packages uninstalled successfully.")
    }
}

// handles the search command for packages across different package managers
func handleSearch(cmd *flag.FlagSet, args []string) {
    if len(args) < 1 {
        fmt.Println("You must specify a package name.")
        cmd.Usage()
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
func handleRebuild(cmd *flag.FlagSet) {
    packageName := cmd.String("package", "", "Name of the AUR package to rebuild")

    cmd.Parse(os.Args[2:])

    if *packageName == "" {
        fmt.Println("You must specify a package name.")
        cmd.Usage()
        return
    }

    err := packagemanager.RebuildAndReinstallAURPackage(*packageName)
    if err != nil {
        fmt.Printf("Error rebuilding package %s: %v\n", *packageName, err)
    } else {
        fmt.Printf("Package %s rebuilt and reinstalled successfully.\n", *packageName)
    }
}

// handles the cleaning of AUR cache
func handleCleanAur(cmd *flag.FlagSet) {
    // Parse the command flags if needed
    cmd.Parse(os.Args[2:])

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

func handleToolCheck(cmd *flag.FlagSet) {
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
