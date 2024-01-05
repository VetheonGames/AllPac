package packagemanager

// This package is responsible for handling updating and uninstalling pacman packages

import (
	"fmt"
	"os/exec"

	"pixelridgesoftworks.com/AllPac/pkg/logger"
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
		logger.Errorf("error updating Pacman packages: %s, %v", string(output), err)
        return fmt.Errorf("error updating Pacman packages: %s, %v", string(output), err)
    }
    return nil
}

// UninstallPacmanPackage uninstalls a specified Pacman package
func UninstallPacmanPackage(packageName string) error {
    cmd := exec.Command("sudo", "pacman", "-Rns", "--noconfirm", packageName)
    if output, err := cmd.CombinedOutput(); err != nil {
		logger.Errorf("error uninstalling Pacman package: %s, %v", output, err)
        return fmt.Errorf("error uninstalling Pacman package: %s, %v", output, err)
    }
    return nil
}
