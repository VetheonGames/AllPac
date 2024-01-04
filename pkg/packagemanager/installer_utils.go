package packagemanager

import (
    "bufio"
    "os"
	"os/user"
    "path/filepath"
    "strings"
    "fmt"
	"encoding/json"
)

// extractVersionFromPKGBUILD reads the PKGBUILD file and extracts the package version
func ExtractVersionFromPKGBUILD(repoDir string) (string, error) {
    pkgbuildPath := filepath.Join(repoDir, "PKGBUILD")
    file, err := os.Open(pkgbuildPath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "pkgver=") {
            return strings.TrimPrefix(line, "pkgver="), nil
        }
    }

    if err := scanner.Err(); err != nil {
        return "", err
    }

    return "", fmt.Errorf("pkgver not found in PKGBUILD")
}

type PackageInfo struct {
    Source  string `json:"source"`
    Version string `json:"version"`
}

type PackageList map[string]PackageInfo

const pkgListFilename = "pkg.list"

// getPkgListPath returns the file path for the package list
func GetPkgListPath() (string, error) {
    usr, err := user.Current()
    if err != nil {
        return "", fmt.Errorf("error getting current user: %v", err)
    }
    return filepath.Join(usr.HomeDir, ".allpac", pkgListFilename), nil
}

// readPackageList reads the package list from the file
func ReadPackageList() (PackageList, error) {
    pkgListPath, err := GetPkgListPath()
    if err != nil {
        return nil, err
    }

    file, err := os.Open(pkgListPath)
    if err != nil {
        if os.IsNotExist(err) {
            return PackageList{}, nil // Return an empty list if file doesn't exist
        }
        return nil, fmt.Errorf("error opening package list file: %v", err)
    }
    defer file.Close()

    var pkgList PackageList
    err = json.NewDecoder(file).Decode(&pkgList)
    if err != nil {
        return nil, fmt.Errorf("error decoding package list: %v", err)
    }

    return pkgList, nil
}

// writePackageList writes the package list to the file
func writePackageList(pkgList PackageList) error {
    pkgListPath, err := GetPkgListPath()
    if err != nil {
        return err
    }

    file, err := os.Create(pkgListPath)
    if err != nil {
        return fmt.Errorf("error creating package list file: %v", err)
    }
    defer file.Close()

    err = json.NewEncoder(file).Encode(pkgList)
    if err != nil {
        return fmt.Errorf("error encoding package list: %v", err)
    }

    return nil
}

// logInstallation logs the package installation details
func LogInstallation(packageName, source, version string) error {
    pkgList, err := readPackageList()
    if err != nil {
        return err
    }

    pkgList[packageName] = PackageInfo{
        Source:  source,
        Version: version,
    }

    return writePackageList(pkgList)
}

// requestRootPermissions prompts the user for root permissions
func requestRootPermissions() bool {
    fmt.Println("Root permissions are required to install AUR packages.")
    return confirmAction("Do you want to continue with root permissions?")
}

// confirmAction prompts the user with a yes/no question and returns true if the answer is yes
func confirmAction(question string) bool {
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Printf("%s [Y/n]: ", question)
        response, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading response:", err)
            return false
        }
        response = strings.ToLower(strings.TrimSpace(response))

        if response == "y" || response == "yes" {
            return true
        } else if response == "n" || response == "no" {
            return false
        }
    }
}
