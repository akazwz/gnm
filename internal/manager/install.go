package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Install a specific Node.js version",
	Long: `Install a specific version of Node.js.
If version is "lts", the latest LTS version will be installed.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]

		// If version is "lts", find the latest LTS version
		if version == "lts" {
			versions, err := FetchAvailableVersions()
			if err != nil {
				fmt.Printf("Error fetching available versions: %v\n", err)
				return
			}

			// Find the latest LTS version
			for _, v := range versions {
				if IsLTS(v.LTS) {
					version = v.Version
					break
				}
			}

			if version == "lts" {
				fmt.Println("Could not find latest LTS version")
				return
			}

			fmt.Printf("Latest LTS version is %s\n", version)
		}

		version = NormalizeVersion(version)

		// Check if already installed
		if IsVersionInstalled(version) {
			fmt.Printf("Node.js %s is already installed\n", version)
			return
		}

		// Create temp directory
		tmpDir, err := os.MkdirTemp("", "gnm-*")
		if err != nil {
			fmt.Printf("Error creating temporary directory: %v\n", err)
			return
		}
		defer os.RemoveAll(tmpDir)

		// Define target file paths
		osType := runtime.GOOS
		arch := GetNodeArch()
		fileName := fmt.Sprintf("node-%s-%s-%s.tar.gz", version, osType, arch)
		url := fmt.Sprintf("%s%s/%s", BaseURL, version, fileName)
		tarballPath := filepath.Join(tmpDir, fileName)

		fmt.Printf("Downloading Node.js %s (%s %s)...\n", version, osType, arch)

		// Download tarball
		if err := DownloadFile(url, tarballPath); err != nil {
			fmt.Printf("Error downloading Node.js %s: %v\n", version, err)
			return
		}

		// Create version directory
		versionDir := filepath.Join(VersionsDir, version)

		fmt.Printf("Extracting Node.js %s...\n", version)

		// Extract tarball
		if err := ExtractTarGz(tarballPath, versionDir); err != nil {
			fmt.Printf("Error extracting tarball: %v\n", err)
			return
		}

		fmt.Printf("Node.js %s installed successfully\n", version)

		// Ask if user wants to use this version
		fmt.Print("Do you want to use this version now? (y/n): ")
		var answer string
		fmt.Scanln(&answer)

		if answer == "y" || answer == "Y" {
			useVersion(version)
		}
	},
}
