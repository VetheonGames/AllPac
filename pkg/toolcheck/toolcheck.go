package toolcheck

// This package is responsible for checking to ensure all our tools are available to us.
// Since we aren't hooking directly into the internals of anything, we require the availability of the packages
// on the system in order to use their CLIs.
// In the future, we might hook directly into the backends for pacman, flatpak, and snapd
// but for now, this is a perfectly fine way of going about it without introducing weird bugs

import (
    "os/exec"
    "fmt"
    "pixelridgesoftworks.com/AllPac/pkg/packagemanager"
)

// isCommandAvailable checks if a command exists
func isCommandAvailable(name string) bool {
    cmd := exec.Command("which", name)
    if err := cmd.Run(); err != nil {
        return false
    }
    return true
}

// EnsurePacman ensures that Pacman is installed and available
func EnsurePacman() error {
    if !isCommandAvailable("pacman") {
        // Pacman should always be available on Arch-based systems, handle this as an error or special case
        return fmt.Errorf("pacman is not available, which is required for AllPac to function")
    }
    return nil
}

// EnsureSnap ensures that Snap is installed and available
func EnsureSnap() error {
    if !isCommandAvailable("snap") {
        return packagemanager.InstallSnap()
    }
    return nil
}

// EnsureGit ensures that Git is installed and available
func EnsureGit() error {
    if !isCommandAvailable("git") {
        return packagemanager.InstallGit()
    }
    return nil
}

// EnsureBaseDevel ensures that the base-devel group is installed
func EnsureBaseDevel() error {
    if !isCommandAvailable("make") { // 'make' is part of base-devel, this is the best method to check
        return packagemanager.InstallBaseDevel()
    }
    return nil
}

// EnsureFlatpak ensures that Flatpak is installed and available
func EnsureFlatpak() error {
    if !isCommandAvailable("flatpak") {
        return packagemanager.InstallFlatpak()
    }
    return nil
}
