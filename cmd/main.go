package main

// This file is our main entrypoint, and build point for AllPac

import (
    "flag"
    "fmt"
    "os"
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

func handleInstall(cmd *flag.FlagSet) {
    // Parse and handle the install command
    // Use functions from install to install packages
}

func handleUninstall(cmd *flag.FlagSet) {
    // Parse and handle the uninstall command
    // Use functions from packagemanager to uninstall packages
}

func handleSearch(cmd *flag.FlagSet) {
    // Parse and handle the search command
    // Use functions from search to search for packages
}

func handleRebuild(cmd *flag.FlagSet) {
    // Parse and handle the search command
    // Use functions from search to search for packages
}

func handleCleanAur(cmd *flag.FlagSet) {
    // Parse and handle the search command
    // Use functions from search to search for packages
}
