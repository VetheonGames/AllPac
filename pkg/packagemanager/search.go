package packagemanager

// This package is responsible for searching various sources for the availability of the requested package

import (
    "encoding/json"
    "io/ioutil"
    "fmt"
    "os/user"
	"os/exec"
    "path/filepath"
	"strings"
	"net/http"
    "pixelridgesoftworks.com/AllPac/pkg/logger"
)

// uninstalls the provided packages
func UninstallPackages(packageNames []string) error {
    pkgList, err := readPackageList()
    if err != nil {
        logger.Errorf("An error has occured:", err)
        return err
    }

    for _, packageName := range packageNames {
        pkgInfo, exists := pkgList[packageName]
        if !exists {
            logger.Warnf("Package %s not found in installed packages list\n", packageNames)
            fmt.Printf("Package %s not found in installed packages list\n", packageName)
            continue
        }

        switch pkgInfo.Source {
        case "pacman":
            err = UninstallPacmanPackage(packageName)
        case "snap":
            err = UninstallSnapPackage(packageName)
        case "flatpak":
            err = UninstallFlatpakPackage(packageName)
        case "aur":
            err = UninstallAURPackage(packageName)
        default:
            logger.Warnf("Unknown source for package %s\n", packageNames)
            fmt.Printf("Unknown source for package %s\n", packageName)
            continue
        }

        if err != nil {
            logger.Warnf("Error uninstalling package %s: %v\n", packageName, err)
            fmt.Printf("Error uninstalling package %s: %v\n", packageName, err)
        } else {
            logger.Infof("Successfully uninstalled package %s\n", packageName)
            fmt.Printf("Successfully uninstalled package %s\n", packageName)
        }
    }

    return nil
}

// reads the package list from the pkg.list file
func readPackageList() (PackageList, error) {
    usr, err := user.Current()
    if err != nil {
        logger.Errorf("error getting current user: %v", err)
        return nil, fmt.Errorf("error getting current user: %v", err)
    }
    pkgListPath := filepath.Join(usr.HomeDir, ".allpac", "pkg.list")

    file, err := ioutil.ReadFile(pkgListPath)
    if err != nil {
        logger.Errorf("error reading package list file: %v", err)
        return nil, fmt.Errorf("error reading package list file: %v", err)
    }

    var pkgList PackageList
    err = json.Unmarshal(file, &pkgList)
    if err != nil {
        logger.Errorf("error decoding package list: %v", err)
        return nil, fmt.Errorf("error decoding package list: %v", err)
    }

    return pkgList, nil
}

// represents the structure of the response from AUR RPC
type AURResponse struct {
    Version     int `json:"version"`
    Type        string `json:"type"`
    ResultCount int `json:"resultcount"`
    Results     []AURPackage `json:"results"`
}

// represents a package in the AUR
type AURPackage struct {
    ID       int `json:"ID"`
    Name     string `json:"Name"`
    Version  string `json:"Version"`
    Description string `json:"Description"`
    URL      string `json:"URL"`
    // Add other fields as needed
}

// searches for a package in the Pacman repositories
func SearchPacman(packageName string) ([]string, error) {
    cmd := exec.Command("pacman", "-Ss", packageName)
    output, err := cmd.CombinedOutput()

    // Check if the error is due to no results found
    if err != nil && len(output) == 0 {
        return nil, nil // No results is not an error in this context
    } else if err != nil {
        // Other errors are still treated as errors
        logger.Errorf("Error searching Pacman: %v", err)
        return nil, fmt.Errorf("error searching Pacman: %v", err)
    }

    return parsePacmanOutput(string(output)), nil
}

// searches for a package in the Snap store
func SearchSnap(packageName string) ([]string, error) {
    cmd := exec.Command("snap", "find", packageName)
    output, err := cmd.CombinedOutput()
    if err != nil {
        logger.Errorf("error searching Snap: %v", err)
        return nil, fmt.Errorf("error searching Snap: %v", err)
    }
    return strings.Split(string(output), "\n"), nil
}

// searches for a package in Flatpak repositories
func SearchFlatpak(packageName string) ([]string, error) {
    cmd := exec.Command("flatpak", "search", packageName)
    output, err := cmd.CombinedOutput()
    if err != nil {
        logger.Errorf("error searching Flatpak: %v", err)
        return nil, fmt.Errorf("error searching Flatpak: %v", err)
    }
    return strings.Split(string(output), "\n"), nil
}

// searches the AUR for the given term
func SearchAUR(searchTerm string) ([]AURPackage, error) {
    url := fmt.Sprintf("https://aur.archlinux.org/rpc/?v=5&type=search&arg=%s", searchTerm)
    resp, err := http.Get(url)
    if err != nil {
        logger.Errorf("error making request to AUR: %v", err)
        return nil, fmt.Errorf("error making request to AUR: %v", err)
    }
    defer resp.Body.Close()

    var aurResponse AURResponse
    if err := json.NewDecoder(resp.Body).Decode(&aurResponse); err != nil {
        logger.Errorf("error decoding AUR response: %v", err)
        return nil, fmt.Errorf("error decoding AUR response: %v", err)
    }

    return aurResponse.Results, nil
}

// parses the output from Pacman search command
func parsePacmanOutput(output string) []string {
    sections := strings.Split(output, "\n\n")

    var packages []string
    for _, section := range sections {
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

// returns the version of a package in the Pacman repositories
func GetPacmanPackageVersion(packageName string) (string, error) {
    searchResults, err := SearchPacman(packageName)
    if err != nil {
        logger.Errorf("An error has occured:", err)
        return "", err
    }

    for _, result := range searchResults {
        if strings.Contains(result, packageName) {
            return extractVersionFromPacmanResult(result), nil
        }
    }

    logger.Errorf("package %s not found in Pacman", packageName)
    return "", fmt.Errorf("package %s not found in Pacman", packageName)
}

// extracts the version from a Pacman search result string
func extractVersionFromPacmanResult(result string) string {
    // Assuming the result is in the format "packageName version description"
    parts := strings.Fields(result)
    if len(parts) >= 2 {
        return parts[1]
    }
    return ""
}

// fetches package information from the AUR
func fetchAURPackageInfo(packageName string) (*AURPackageInfo, error) {
    url := fmt.Sprintf("https://aur.archlinux.org/rpc/?v=5&type=info&arg[]=%s", packageName)
    resp, err := http.Get(url)
    if err != nil {
        logger.Errorf("An error has occured:", err)
        return nil, err
    }
    defer resp.Body.Close()

    var result struct {
        Results []AURPackageInfo `json:"results"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        logger.Errorf("An error has occured:", err)
        return nil, err
    }

    if len(result.Results) == 0 {
        logger.Errorf("package %s not found in AUR", packageName)
        return nil, fmt.Errorf("package %s not found in AUR", packageName)
    }

    return &result.Results[0], nil
}

// returns the version of a package in the Snap store
func GetSnapPackageVersion(packageName string) (string, error) {
    cmd := exec.Command("snap", "info", packageName)
    output, err := cmd.CombinedOutput()
    if err != nil {
        logger.Errorf("error getting Snap package info: %v", err)
        return "", fmt.Errorf("error getting Snap package info: %v", err)
    }

    return parseSnapInfoOutput(string(output)), nil
}

// parses the output from the Snap info command to extract the version
func parseSnapInfoOutput(output string) string {
    lines := strings.Split(output, "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "installed:") {
            parts := strings.Fields(line)
            if len(parts) >= 2 {
                return strings.TrimSpace(parts[1])
            }
        }
    }
    return ""
}

// returns the version of a package in Flatpak repositories
func GetFlatpakPackageVersion(packageName string) (string, error) {
    cmd := exec.Command("flatpak", "info", packageName)
    output, err := cmd.CombinedOutput()
    if err != nil {
        logger.Errorf("error getting Flatpak package info: %v", err)
        return "", fmt.Errorf("error getting Flatpak package info: %v", err)
    }

    return parseFlatpakInfoOutput(string(output)), nil
}

// parses the output from the Flatpak info command to extract the version
func parseFlatpakInfoOutput(output string) string {
    lines := strings.Split(output, "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "Version:") {
            parts := strings.Fields(line)
            if len(parts) >= 2 {
                return strings.TrimSpace(parts[1])
            }
        }
    }
    return ""
}

// returns the version of a package in the AUR
func GetAURPackageVersion(packageName string) (string, error) {
    aurInfo, err := fetchAURPackageInfo(packageName)
    if err != nil {
        logger.Errorf("An error has occured:", err)
        return "", err
    }
    return aurInfo.Version, nil
}

// represents the search result from a specific source
type SourceResult struct {
    Source  string
    Results []string
}

// represents the search results for a package across different sources
type PackageSearchResult struct {
    PackageName string
    Results     []SourceResult
}

// searches for packages across Pacman, Snap, Flatpak, and AUR
func SearchAllSources(packageNames []string) ([]PackageSearchResult, error) {
    var allPackageResults []PackageSearchResult

    for _, packageName := range packageNames {
        var packageResults PackageSearchResult
        packageResults.PackageName = packageName

        // Search in Pacman
        pacmanResults, err := SearchPacman(packageName)
        if err == nil {
            packageResults.Results = append(packageResults.Results, SourceResult{"Pacman", pacmanResults})
        }

        // Search in Snap
        snapResults, err := SearchSnap(packageName)
        if err == nil {
            packageResults.Results = append(packageResults.Results, SourceResult{"Snap", snapResults})
        }

        // Search in Flatpak
        flatpakResults, err := SearchFlatpak(packageName)
        if err == nil {
            packageResults.Results = append(packageResults.Results, SourceResult{"Flatpak", flatpakResults})
        }

        // Search in AUR
        aurResults, err := SearchAUR(packageName)
        if err == nil {
            var aurResultStrings []string
            for _, result := range aurResults {
                aurResultStrings = append(aurResultStrings, fmt.Sprintf("%s - %s", result.Name, result.Version))
            }
            packageResults.Results = append(packageResults.Results, SourceResult{"AUR", aurResultStrings})
        }

        allPackageResults = append(allPackageResults, packageResults)
    }

    return allPackageResults, nil
}
