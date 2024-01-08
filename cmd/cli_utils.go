package main

import (
	"fmt"
	"embed"
	"io/fs"
	"pixelridgesoftworks.com/AllPac/pkg/logger"
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
