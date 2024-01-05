package packagemanager

import (
    "bufio"
    "os"
    "path/filepath"
    "strings"
    "fmt"
    "pixelridgesoftworks.com/AllPac/pkg/logger"
)

// extractVersionFromPKGBUILD reads the PKGBUILD file and extracts the package version
func ExtractVersionFromPKGBUILD(repoDir string) (string, error) {
    pkgbuildPath := filepath.Join(repoDir, "PKGBUILD")
    file, err := os.Open(pkgbuildPath)
    if err != nil {
        logger.Errorf("An error has occured:", err)
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
        logger.Errorf("An error has occured:", err)
        return "", err
    }

    logger.Errorf("pkgver not found in PKGBUILD")
    return "", fmt.Errorf("pkgver not found in PKGBUILD")
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
            logger.Errorf("Error reading response: %v", err)
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
