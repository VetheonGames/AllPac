package packagemanager

import (
    "bufio"
    "os"
    "path/filepath"
    "strings"
    "fmt"
    "pixelridgesoftworks.com/AllPac/pkg/logger"
    "os/exec"
    "os/user"
    "strconv"
    "syscall"
)

// reads the PKGBUILD file and extracts the package version
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

// prompts the user with a yes/no question and returns true if the answer is yes
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

// this is unused, just incase I need to do it this way since makepkg is being a pain in the neck
func RunMakepkgAsUser(username string) error {
    // Lookup the non-root user
    usr, err := user.Lookup(username)
    if err != nil {
        return err
    }

    // Convert UID and GID to integers
    uid, _ := strconv.Atoi(usr.Uid)
    gid, _ := strconv.Atoi(usr.Gid)

    // Set UID and GID of the process
    err = syscall.Setgid(gid)
    if err != nil {
        return err
    }
    err = syscall.Setuid(uid)
    if err != nil {
        return err
    }

    // Now run makepkg as the non-root user
    cmd := exec.Command("makepkg", "-si", "--noconfirm")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}
