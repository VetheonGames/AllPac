package packagemanager

import (
	"pixelridgesoftworks.com/AllPac/pkg/logger"
	"os"
	"os/user"
	"fmt"
    "path/filepath"
	"encoding/json"
)

type PackageInfo struct {
    Source  string `json:"source"`
    Version string `json:"version"`
}

type PackageList map[string]PackageInfo

const pkgListFilename = "pkg.list"

// returns the file path for the package list
func GetPkgListPath() (string, error) {
    usr, err := user.Current()
    if err != nil {
        logger.Errorf("error getting current user: %v", err)
        return "", fmt.Errorf("error getting current user: %v", err)
    }

    pkgListDir := filepath.Join(usr.HomeDir, ".allpac")
    pkgListPath := filepath.Join(pkgListDir, pkgListFilename)

    logger.Infof("Checking directory: %s", pkgListDir)

    // Ensure the directory exists
    if err := os.MkdirAll(pkgListDir, 0755); err != nil {
        logger.Errorf("error creating directory: %v", err)
        return "", fmt.Errorf("error creating directory: %v", err)
    }

    // Check if the pkg.list file exists
    if _, err := os.Stat(pkgListPath); os.IsNotExist(err) {
        logger.Infof("pkg.list file does not exist, initializing: %s", pkgListPath)
        // Create and initialize the file if it doesn't exist
        if err := initializePkgListFile(pkgListPath); err != nil {
            return "", err // Error already logged in initializePkgListFile
        }
    } else if err != nil {
        logger.Errorf("error checking pkg.list file: %v", err)
        return "", fmt.Errorf("error checking pkg.list file: %v", err)
    } else {
        logger.Infof("pkg.list file exists: %s", pkgListPath)
    }

    return pkgListPath, nil
}

// creates a new pkg.list file with an empty JSON object
func initializePkgListFile(filePath string) error {
    file, err := os.Create(filePath)
    if err != nil {
        logger.Errorf("error creating package list file: %v", err)
        return fmt.Errorf("error creating package list file: %v", err)
    }
    defer file.Close()

    logger.Infof("Writing empty JSON object to pkg.list file: %s", filePath)

    if _, err := file.WriteString("{}"); err != nil {
        logger.Errorf("error initializing package list file: %v", err)
        return fmt.Errorf("error initializing package list file: %v", err)
    }

    logger.Infof("pkg.list file initialized successfully: %s", filePath)

    return nil
}

// reads the package list from the file
func ReadPackageList() (PackageList, error) {
    pkgListPath, err := GetPkgListPath()
    if err != nil {
        logger.Errorf("An error has occurred: %v", err)
        return nil, err
    }

    // Ensure the directory exists
    if err := os.MkdirAll(filepath.Dir(pkgListPath), 0755); err != nil {
        logger.Errorf("error creating directory: %v", err)
        return nil, fmt.Errorf("error creating directory: %v", err)
    }

    // Open or create the file
    file, err := os.OpenFile(pkgListPath, os.O_RDWR|os.O_CREATE, 0600)
    if err != nil {
        logger.Errorf("error opening or creating package list file: %v", err)
        return nil, fmt.Errorf("error opening or creating package list file: %v", err)
    }
    defer file.Close()

    // Check if the file is empty
    fileInfo, err := file.Stat()
    if err != nil {
        logger.Errorf("error getting file info: %v", err)
        return nil, fmt.Errorf("error getting file info: %v", err)
    }

    if fileInfo.Size() == 0 {
        // Initialize file with an empty JSON object
        if _, err := file.WriteString("{}"); err != nil {
            logger.Errorf("error initializing package list file: %v", err)
            return nil, fmt.Errorf("error initializing package list file: %v", err)
        }
        if _, err := file.Seek(0, 0); err != nil { // Reset file pointer to the beginning
            logger.Errorf("error seeking in package list file: %v", err)
            return nil, fmt.Errorf("error seeking in package list file: %v", err)
        }
    }

    var pkgList PackageList
    err = json.NewDecoder(file).Decode(&pkgList)
    if err != nil {
        logger.Errorf("error decoding package list: %v", err)
        return nil, fmt.Errorf("error decoding package list: %v", err)
    }

    return pkgList, nil
}

// writes the package list to the file
func writePackageList(pkgList PackageList) error {
    pkgListPath, err := GetPkgListPath()
    if err != nil {
        logger.Errorf("An error has occured:", err)
        return err
    }

    file, err := os.Create(pkgListPath)
    if err != nil {
        logger.Errorf("error creating package list file: %v", err)
        return fmt.Errorf("error creating package list file: %v", err)
    }
    defer file.Close()

    err = json.NewEncoder(file).Encode(pkgList)
    if err != nil {
        logger.Errorf("error encoding package list: %v", err)
        return fmt.Errorf("error encoding package list: %v", err)
    }

    return nil
}

// logs the package installation details
func LogInstallation(packageName, source, version string) error {
    pkgList, err := readPackageList()
    if err != nil {
        logger.Errorf("An error has occured:", err)
        return err
    }

    pkgList[packageName] = PackageInfo{
        Source:  source,
        Version: version,
    }

    return writePackageList(pkgList)
}

// removes a package from the package list file
func RemovePackageFromList(packageName string) error {
    // Read the current package list
    pkgList, err := ReadPackageList()
    if err != nil {
        logger.Errorf("An error has occurred while reading the package list: %v", err)
        return err
    }

    // Check if the package exists in the list
    if _, exists := pkgList[packageName]; !exists {
        logger.Infof("Package %s not found in the package list, no action taken", packageName)
        return nil // No need to update the file if the package isn't there
    }

    // Remove the package from the list
    delete(pkgList, packageName)
    logger.Infof("Package %s removed from the package list", packageName)

    // Write the updated list back to the file
    if err := writePackageList(pkgList); err != nil {
        logger.Errorf("An error has occurred while writing the updated package list: %v", err)
        return err
    }

    return nil
}

// updates the details of a package in the package list file
func UpdatePackageInList(packageName, source, newVersion string) error {
    // Read the current package list
    pkgList, err := ReadPackageList()
    if err != nil {
        logger.Errorf("An error has occurred while reading the package list: %v", err)
        return err
    }

    // Check if the package exists in the list
    if pkgInfo, exists := pkgList[packageName]; exists {
        // Update the package details
        pkgInfo.Source = source
        pkgInfo.Version = newVersion
        pkgList[packageName] = pkgInfo
        logger.Infof("Package %s updated in the package list", packageName)
    } else {
        logger.Infof("Package %s not found in the package list, adding new entry", packageName)
        // If the package is not found, add it as a new entry
        pkgList[packageName] = PackageInfo{
            Source:  source,
            Version: newVersion,
        }
    }

    // Write the updated list back to the file
    if err := writePackageList(pkgList); err != nil {
        logger.Errorf("An error has occurred while writing the updated package list: %v", err)
        return err
    }

    return nil
}
