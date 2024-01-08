# AllPac - Unified Package Manager for Arch Linux

## Overview

AllPac is a command-line utility designed to simplify package management on Arch Linux by combining various package managers into one cohesive tool. With AllPac, users can seamlessly interact with packages from the Snap Store, Flatpak, Pacman, and the Arch User Repository (AUR) using a single interface. This eliminates the need to juggle multiple package managers and provides a unified solution for installing, updating, uninstalling, and searching for packages.

## Installation

To install AllPac on your Arch Linux system, simply run the following command to run the install script ([Source](https://git.pixelridgesoftworks.com/PixelRidge-Softworks/Installers/src/branch/main/allpac)):
(if you don't want to use the install script, a pre-built binary can be found [here](https://git.pixelridgesoftworks.com/PixelRidge-Softworks/AllPac/releases), you will need to run `touch pkg.list && echo "{}" > ./pkg.list` where you want to run the binary from)
```bash
curl -s -o install.sh https://git.pixelridgesoftworks.com/PixelRidge-Softworks/Installers/raw/branch/main/allpac/install.sh
chmod +x install.sh
./install.sh
```

## Features

### 1. Unified Package Management

AllPac consolidates package management tasks from different sources, allowing users to handle Snap packages, Flatpaks, Pacman packages, and AUR packages all in one place.

### 2. Installer

Easily install packages from various sources with a straightforward installation command. AllPac intelligently recognizes the package type and fetches it from the appropriate repository.

```bash
allpac install
```

### 3. Updater

Keep all your installed packages up-to-date with a single command. AllPac checks for updates across different repositories and ensures your system is current.

```bash
allpac update
```

### 4. Uninstaller

Remove packages cleanly and efficiently, regardless of their origin. AllPac ensures a consistent uninstallation process for Snap, Flatpak, Pacman, and AUR packages.

```bash
allpac uninstall
```

### 5. Package Search

Quickly find packages across Snap Store, Flatpak, Pacman, and AUR using the integrated search feature.

```bash
allpac search
```

## Usage

Once installed, you can use AllPac with the following commands:

- Install a package:
  ```bash
  allpac install <package_name>
  ```

- Update all installed packages:
  ### WARNING: This will attempt to install all packages managed by AllPac across all sources! Be careful with this command!
  ```bash
  allpac update everything
  ```

- Update a specific installed package or packages:
  ```bash
  allpac update {package_name}
  ```
  or
  ```bash
  allpac update {packagename1} {packagename2} {packagename3}
  ```
  or
  ```bash
  allpac update {aur/flats/snaps/arch}
  ```

- Uninstall a package:
  ```bash
  allpac uninstall <package_name>
  ```

- Search for packages:
  ```bash
  allpac search <package_name>
  ```

## Logs and Cache

After you run things the first time (or you run the install script), all the logs, the package list, the binary, and the updater script will be contained here:
```bash
/home/{your_user}/.allpac/
```

## Uninstalling AllPac

To uninstall AllPac is quite simple. You just remove the `.allpac` directory. As the directory contains all the files associated with AllPac, removing the directory will completely remove AllPac.

### NOTE: UNINSTALLING AllPac WILL *NOT* UNINSTALL PACKAGES INSTALLED *BY* AllPac!

## Updating AllPac

If you used the Installer Script, updating is easy. Just run `allpac-update-system`.

If you did *NOT* use the Installer Script, updating is still super easy. Just use `wget` to pull down the [updater script](https://git.pixelridgesoftworks.com/PixelRidge-Softworks/Installers/raw/branch/main/allpac/update.sh):
```bash
wget -O ~/.allpac/bin/allpac-updater.sh "https://git.pixelridgesoftworks.com/PixelRidge-Softworks/Installers/raw/branch/main/allpac/update.sh"
```

Give it the needed permissions:
```bash
chmod u+rwx ~/.allpac/bin/allpac-updater.sh
```

Then run the updater script:
```bash
. ~/.allpac/bin/allpac-updater.sh
```

## Feedback and Contributions

Feedback, bug reports, and contributions are welcome! Feel free to open issues on the [Git repository](https://git.pixelridgesoftworks.com/PixelRidge-Softworks/AllPac/issues) or submit pull requests.

## License

This project is licensed under the PixelRidge-BEGPULSE License. See the [LICENSE](LICENSE) file for details.
