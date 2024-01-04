# AllPac - Unified Package Manager for Arch Linux

## Overview

AllPac is a command-line utility designed to simplify package management on Arch Linux by combining various package managers into one cohesive tool. With AllPac, users can seamlessly interact with packages from the Snap Store, Flatpak, Pacman, and the Arch User Repository (AUR) using a single interface. This eliminates the need to juggle multiple package managers and provides a unified solution for installing, updating, uninstalling, and searching for packages.

## Features

### 1. Unified Package Management

AllPac consolidates package management tasks from different sources, allowing users to handle Snap packages, Flatpaks, Pacman packages, and AUR packages all in one place.

### 2. Installer

Easily install packages from various sources with a straightforward installation command. AllPac intelligently recognizes the package type and fetches it from the appropriate repository.

```bash
allpac install <package_name>
```

### 3. Updater

Keep all your installed packages up-to-date with a single command. AllPac checks for updates across different repositories and ensures your system is current.

```bash
allpac update
```

### 4. Uninstaller

Remove packages cleanly and efficiently, regardless of their origin. AllPac ensures a consistent uninstallation process for Snap, Flatpak, Pacman, and AUR packages.

```bash
allpac uninstall <package_name>
```

### 5. Package Search

Quickly find packages across Snap Store, Flatpak, Pacman, and AUR using the integrated search feature.

```bash
allpac search <query>
```

## Installation

To install AllPac on your Arch Linux system, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://git.pixelridgesoftworks.com/PixelRidge-Softworks/AllPac.git
   ```

2. Navigate to the AllPac directory:
   ```bash
   cd AllPac
   ```

3. Run the installer script:
   ```bash
   ./install.sh
   ```

## Usage

Once installed, you can use AllPac with the following commands:

- Install a package:
  ```bash
  allpac install <package_name>
  ```

- Update all installed packages:
  ```bash
  allpac update
  ```

- Update a specific installed package or packages:
  ```bash
  allpac update {package_name}
  ```
  or
  ```bash
  allpac update {packagename1} {packagename2} {packagename3}
  ```

- Uninstall a package:
  ```bash
  allpac uninstall <package_name>
  ```

- Search for packages:
  ```bash
  allpac search <query>
  ```

## Feedback and Contributions

Feedback, bug reports, and contributions are welcome! Feel free to open issues on the [Git repository](https://git.pixelridgesoftworks.com/PixelRidge-Softworks/AllPac/issues) or submit pull requests.

## License

This project is licensed under the PixelRidge-BEGPULSE License. See the [LICENSE](LICENSE) file for details.
