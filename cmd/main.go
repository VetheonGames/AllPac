package main

// This file is our main entrypoint, and build point for AllPac

import (
    "flag"
    "fmt"
    "os"
	"strings"
    "pixelridgesoftworks.com/AllPac/pkg/packagemanager"
)

func main() {
    // Define flags for different commands
    updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
    installCmd := flag.NewFlagSet("install", flag.ExitOnError)
    uninstallCmd := flag.NewFlagSet("uninstall", flag.ExitOnError)
    searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	aurRebuildCmd := flag.NewFlagSet("rebuild", flag.ExitOnError)
	aurCleanCmd := flag.NewFlagSet("clean-aur", flag.ExitOnError)

    if len(os.Args) < 2 {
        fmt.Println("Expected 'update', 'install', 'uninstall', 'search', 'rebuild', or 'clean-aur' subcommands")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "update":
        handleUpdate(updateCmd)
    case "install":
        handleInstall(installCmd)
    case "uninstall":
        handleUninstall(uninstallCmd)
    case "search":
        handleSearch(searchCmd)
	case "rebuild":
		handleRebuild(aurRebuildCmd)
	case "clean-aur":
		handleCleanAur(aurCleanCmd)
    default:
        fmt.Printf("Unknown subcommand: %s\n", os.Args[1])
        os.Exit(1)
    }
}

func handleUpdate(cmd *flag.FlagSet) {
    everythingFlag := cmd.Bool("everything", false, "Update all packages on the system")
    snapFlag := cmd.Bool("snap", false, "Update all Snap packages")
    aurFlag := cmd.Bool("aur", false, "Update all AUR packages")
    archFlag := cmd.Bool("arch", false, "Update all Arch packages")
    flatsFlag := cmd.Bool("flats", false, "Update all Flatpak packages")
    cmd.Parse(os.Args[2:])

    if *everythingFlag {
        // Call function to update all packages
        packagemanager.UpdateAllPackages()
    } else if *snapFlag {
        // Call function to update Snap packages
        packagemanager.UpdateSnapPackages()
    } else if *aurFlag {
        // Call function to update AUR packages
        packagemanager.UpdateAURPackages()
    } else if *archFlag {
        // Call function to update Arch packages
        packagemanager.UpdatePacmanPackages()
    } else if *flatsFlag {
        // Call function to update Flatpak packages
        packagemanager.UpdateFlatpakPackages()
    } else {
        fmt.Println("No update option specified or unrecognized option")
    }
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

    for _, result := range searchResults {
        fmt.Printf("Searching for package: %s\n", result.PackageName)
        if len(result.Results) == 0 {
            fmt.Println("No sources found for package.")
            continue
        }

        sourceIndex := promptUserForSource(result.Results)
        if sourceIndex < 0 || sourceIndex >= len(result.Results) {
            fmt.Println("Invalid selection. Skipping package.")
            continue
        }

        selectedSource := result.Results[sourceIndex].Source
        fmt.Printf("Installing %s from %s...\n", result.PackageName, selectedSource)

        switch selectedSource {
        case "pacman":
            err = packagemanager.InstallPackagePacman(result.PackageName)
        case "snap":
            err = packagemanager.InstallPackageSnap(result.PackageName)
        case "flatpak":
            err = packagemanager.InstallPackageFlatpak(result.PackageName)
        case "aur":
            _, err = packagemanager.CloneAndInstallFromAUR(fmt.Sprintf("https://aur.archlinux.org/%s.git", result.PackageName), false)
        default:
            fmt.Printf("Unknown source for package %s\n", result.PackageName)
            continue
        }

        if err != nil {
            fmt.Printf("Error installing package %s from %s: %v\n", result.PackageName, selectedSource, err)
        } else {
            fmt.Printf("Package %s installed successfully from %s.\n", result.PackageName, selectedSource)
        }
    }
}

// handleUninstall handles the uninstall command for packages
func handleUninstall(cmd *flag.FlagSet) {
    // Define a flag for accepting multiple package names
    packageNames := cmd.String("packages", "", "Comma-separated list of packages to uninstall")

    // Parse the command line arguments
    cmd.Parse(os.Args[2:])

    // Check if the package names were provided
    if *packageNames == "" {
        fmt.Println("You must specify at least one package name.")
        cmd.Usage()
        return
    }

    // Split the package names and convert to a slice
    packages := strings.Split(*packageNames, ",")

    // Call the function to uninstall the packages
    err := packagemanager.UninstallPackages(packages)
    if err != nil {
        fmt.Printf("Error uninstalling packages: %v\n", err)
    } else {
        fmt.Println("Requested packages uninstalled successfully.")
    }
}

// handleSearch handles the search command for packages across different package managers
func handleSearch(cmd *flag.FlagSet) {
    // Define a flag for the package name
    packageName := cmd.String("package", "", "Name of the package to search")

    // Parse the command line arguments
    cmd.Parse(os.Args[2:])

    // Check if the package name was provided
    if *packageName == "" {
        fmt.Println("You must specify a package name.")
        cmd.Usage()
        return
    }

    // Search in Pacman
    pacmanResults, err := packagemanager.SearchPacman(*packageName)
    if err != nil {
        fmt.Printf("Error searching in Pacman: %v\n", err)
    } else {
        fmt.Println("Pacman Results:")
        for _, result := range pacmanResults {
            fmt.Println(result)
        }
    }

    // Search in Snap
    snapResults, err := packagemanager.SearchSnap(*packageName)
    if err != nil {
        fmt.Printf("Error searching in Snap: %v\n", err)
    } else {
        fmt.Println("Snap Results:")
        for _, result := range snapResults {
            fmt.Println(result)
        }
    }

    // Search in Flatpak
    flatpakResults, err := packagemanager.SearchFlatpak(*packageName)
    if err != nil {
        fmt.Printf("Error searching in Flatpak: %v\n", err)
    } else {
        fmt.Println("Flatpak Results:")
        for _, result := range flatpakResults {
            fmt.Println(result)
        }
    }

    // Search in AUR
    aurResults, err := packagemanager.SearchAUR(*packageName)
    if err != nil {
        fmt.Printf("Error searching in AUR: %v\n", err)
    } else {
        fmt.Println("AUR Results:")
        for _, result := range aurResults {
            fmt.Printf("%s - %s\n", result.Name, result.Version)
        }
    }
}

// handleRebuild handles the rebuild command for an AUR package
func handleRebuild(cmd *flag.FlagSet) {
    // Define a flag for the package name
    packageName := cmd.String("package", "", "Name of the AUR package to rebuild")

    // Parse the command line arguments
    cmd.Parse(os.Args[2:])

    // Check if the package name was provided
    if *packageName == "" {
        fmt.Println("You must specify a package name.")
        cmd.Usage()
        return
    }

    // Call the function to rebuild and reinstall the AUR package
    err := packagemanager.RebuildAndReinstallAURPackage(*packageName)
    if err != nil {
        fmt.Printf("Error rebuilding package %s: %v\n", *packageName, err)
    } else {
        fmt.Printf("Package %s rebuilt and reinstalled successfully.\n", *packageName)
    }
}

// handleCleanAur handles the cleaning of AUR cache
func handleCleanAur(cmd *flag.FlagSet) {
    // Parse the command flags if needed
    cmd.Parse(os.Args[2:])

    // Call the function to clear the AUR cache
    err := packagemanager.ClearAllPacCache()
    if err != nil {
        fmt.Printf("Error clearing AUR cache: %v\n", err)
        return
    }

    fmt.Println("AUR cache cleared successfully.")
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
