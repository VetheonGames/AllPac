package packagemanager

// This package is responsible for handling updating and uninstalling pacman packages

import (
    "os/exec"
    "fmt"
	"strings"
)

// UpdatePacmanPackages updates specified Pacman packages or all if no specific package is provided
func UpdatePacmanPackages(packageNames ...string) error {
    var cmd *exec.Cmd
    if len(packageNames) == 0 {
        cmd = exec.Command("sudo", "pacman", "-Syu")
    } else {
        args := append([]string{"sudo", "pacman", "-S", "--noconfirm"}, packageNames...)
        cmd = exec.Command(args[0], args[1:]...)
    }

    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("error updating Pacman packages: %s, %v", string(output), err)
    }
    return nil
}

// UninstallPacmanPackage uninstalls a specified Pacman package
func UninstallPacmanPackage(packageName string) error {
    cmd := exec.Command("sudo", "pacman", "-Rns", "--noconfirm", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
        return fmt.Errorf("error uninstalling Pacman package: %s, %v", output, err)
    }
    return nil
}

// getVersionFromPacman gets the installed version of a package using Pacman
func GetVersionFromPacman(packageName string) (string, error) {
    cmd := exec.Command("pacman", "-Qi", packageName)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("error getting package version: %v", err)
    }

    for _, line := range strings.Split(string(output), "\n") {
        if strings.HasPrefix(line, "Version") {
            return strings.Fields(line)[2], nil
        }
    }
    return "", fmt.Errorf("version not found for package: %s", packageName)
}
