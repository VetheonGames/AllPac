package search

// This package is responsible for searching various sources for the availability of the requested package

import (
    "os/exec"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
)

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
