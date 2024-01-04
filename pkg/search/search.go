package search

// This package is responsible for searching various sources for the availability of the requested package

import (
    "pixelridgesoftworks.com/AllPac/pkg/packagemanager"
    "encoding/json"
    "io/ioutil"
    "fmt"
    "os/user"
	"os/exec"
    "path/filepath"
	"strings"
	"net/http"
)

// PackageList represents the mapping of installed packages to their sources
type PackageList map[string]string

// UninstallPackages uninstalls the provided packages
func UninstallPackages(packageNames []string) error {
    pkgList, err := readPackageList()
    if err != nil {
        return err
    }

    for _, packageName := range packageNames {
        source, exists := pkgList[packageName]
        if !exists {
            fmt.Printf("Package %s not found in installed packages list\n", packageName)
            continue
        }

        switch source {
        case "pacman":
            err = packagemanager.UninstallPacmanPackage(packageName)
        case "snap":
            err = packagemanager.UninstallSnapPackage(packageName)
        case "flatpak":
            err = packagemanager.UninstallFlatpakPackage(packageName)
        // Add cases for other package managers if necessary
        default:
            fmt.Printf("Unknown source for package %s\n", packageName)
            continue
        }

        if err != nil {
            fmt.Printf("Error uninstalling package %s: %v\n", packageName, err)
        } else {
            fmt.Printf("Successfully uninstalled package %s\n", packageName)
        }
    }

    return nil
}

// readPackageList reads the package list from the pkg.list file
func readPackageList() (PackageList, error) {
    usr, err := user.Current()
    if err != nil {
        return nil, fmt.Errorf("error getting current user: %v", err)
    }
    pkgListPath := filepath.Join(usr.HomeDir, ".allpac", "pkg.list")

    file, err := ioutil.ReadFile(pkgListPath)
    if err != nil {
        return nil, fmt.Errorf("error reading package list file: %v", err)
    }

    var pkgList PackageList
    err = json.Unmarshal(file, &pkgList)
    if err != nil {
        return nil, fmt.Errorf("error decoding package list: %v", err)
    }

    return pkgList, nil
}

// AURResponse represents the structure of the response from AUR RPC
type AURResponse struct {
    Version     int `json:"version"`
    Type        string `json:"type"`
    ResultCount int `json:"resultcount"`
    Results     []AURPackage `json:"results"`
}

// AURPackage represents a package in the AUR
type AURPackage struct {
    ID       int `json:"ID"`
    Name     string `json:"Name"`
    Version  string `json:"Version"`
    Description string `json:"Description"`
    URL      string `json:"URL"`
    // Add other fields as needed
}

// SearchPacman searches for a package in the Pacman repositories
func SearchPacman(packageName string) ([]string, error) {
    cmd := exec.Command("pacman", "-Ss", packageName)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("error searching Pacman: %v", err)
    }
    return parsePacmanOutput(string(output)), nil
}

// SearchSnap searches for a package in the Snap store
func SearchSnap(packageName string) ([]string, error) {
    cmd := exec.Command("snap", "find", packageName)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("error searching Snap: %v", err)
    }
    return strings.Split(string(output), "\n"), nil
}

// SearchFlatpak searches for a package in Flatpak repositories
func SearchFlatpak(packageName string) ([]string, error) {
    cmd := exec.Command("flatpak", "search", packageName)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("error searching Flatpak: %v", err)
    }
    return strings.Split(string(output), "\n"), nil
}

// SearchAUR searches the AUR for the given term
func SearchAUR(searchTerm string) ([]AURPackage, error) {
    url := fmt.Sprintf("https://aur.archlinux.org/rpc/?v=5&type=search&arg=%s", searchTerm)
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("error making request to AUR: %v", err)
    }
    defer resp.Body.Close()

    var aurResponse AURResponse
    if err := json.NewDecoder(resp.Body).Decode(&aurResponse); err != nil {
        return nil, fmt.Errorf("error decoding AUR response: %v", err)
    }

    return aurResponse.Results, nil
}

// parsePacmanOutput parses the output from Pacman search command
func parsePacmanOutput(output string) []string {
    // Split the output into sections, each representing a package
    sections := strings.Split(output, "\n\n")

    var packages []string
    for _, section := range sections {
        // Split each section into lines
        lines := strings.Split(section, "\n")

        // The first line should contain the package name and version
        if len(lines) > 0 {
            packageNameLine := lines[0]

            // Check if the package is installed
            if strings.Contains(packageNameLine, "[installed]") {
                packageNameLine += " (Installed)"
            }

            packages = append(packages, packageNameLine)
        }
    }

    return packages
}
