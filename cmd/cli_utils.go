package main

import (
	"embed"
	"fmt"
	"io/fs"

	"pixelridgesoftworks.com/AllPac/pkg/logger"
	"pixelridgesoftworks.com/AllPac/pkg/packagemanager"
)

func handleUpdateError(updateOption string, err error) {
    if err != nil {
        fmt.Printf("Error occurred during '%s' update: %v\n", updateOption, err)
    } else {
        fmt.Printf("Update '%s' completed successfully.\n", updateOption)
    }
}

//go:embed .version
var versionFS embed.FS

func handleVersion(args []string) {
    content, err := fs.ReadFile(versionFS, ".version")
    if err != nil {
        logger.Errorf("Error reading version file: %v", err)
    }
    fmt.Println(string(content))
}

func handleRepair(args []string) {
    // Assuming GetPkgListPath() returns a string path
    pkgListPath, _ := packagemanager.GetPkgListPath()

    err :=packagemanager.InitializePkgListFile(pkgListPath)
    if err != nil {
        logger.Errorf("Error initializing version file: %v", err)
    }
}
