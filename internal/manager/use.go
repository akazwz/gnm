package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use [version]",
	Short: "Use a specific Node.js version",
	Long: `Switch to a specific version of Node.js.
If version is "lts", the latest installed LTS version will be used.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		useVersion(version)
	},
}

// useVersion switches to the specified Node.js version
func useVersion(version string) {
	// If version is "lts", find the latest installed LTS version
	if version == "lts" {
		// Get all installed versions
		entries, err := os.ReadDir(VersionsDir)
		if err != nil {
			fmt.Printf("Error reading versions directory: %v\n", err)
			return
		}

		// Get available versions for LTS info
		availableVersions, err := FetchAvailableVersions()
		if err != nil {
			fmt.Printf("Error fetching available versions: %v\n", err)
			return
		}

		// Create a map of LTS versions
		ltsVersions := make(map[string]bool)
		for _, v := range availableVersions {
			if IsLTS(v.LTS) {
				ltsVersions[v.Version] = true
			}
		}

		// Find the latest installed LTS version
		var latestLTS string
		for _, entry := range entries {
			if entry.IsDir() {
				entryVersion := entry.Name()
				if ltsVersions[entryVersion] {
					if latestLTS == "" || entryVersion > latestLTS {
						latestLTS = entryVersion
					}
				}
			}
		}

		if latestLTS == "" {
			fmt.Println("No LTS version installed")
			return
		}

		version = latestLTS
		fmt.Printf("Using latest installed LTS version: %s\n", version)
	}

	version = NormalizeVersion(version)

	// Check if version is installed
	if !IsVersionInstalled(version) {
		fmt.Printf("Node.js %s is not installed. Install it first with 'gnm install %s'\n", version, version)
		return
	}

	// Get paths to binaries
	versionDir := filepath.Join(VersionsDir, version)
	nodeBin := filepath.Join(versionDir, "bin", "node")
	npmBin := filepath.Join(versionDir, "bin", "npm")
	npxBin := filepath.Join(versionDir, "bin", "npx")

	// Check if binaries exist
	if _, err := os.Stat(nodeBin); os.IsNotExist(err) {
		fmt.Printf("Node binary not found in %s\n", nodeBin)
		return
	}

	// Remove existing symlinks
	os.Remove(filepath.Join(BinDir, "node"))
	os.Remove(filepath.Join(BinDir, "npm"))
	os.Remove(filepath.Join(BinDir, "npx"))

	// Create new symlinks
	if err := os.Symlink(nodeBin, filepath.Join(BinDir, "node")); err != nil {
		fmt.Printf("Error creating symlink for node: %v\n", err)
		return
	}

	if _, err := os.Stat(npmBin); err == nil {
		if err := os.Symlink(npmBin, filepath.Join(BinDir, "npm")); err != nil {
			fmt.Printf("Error creating symlink for npm: %v\n", err)
		}
	}

	if _, err := os.Stat(npxBin); err == nil {
		if err := os.Symlink(npxBin, filepath.Join(BinDir, "npx")); err != nil {
			fmt.Printf("Error creating symlink for npx: %v\n", err)
		}
	}

	fmt.Printf("Now using Node.js %s\n", version)
	fmt.Println("Add this to your shell profile to use gnm:")
	fmt.Printf("  export PATH=\"%s:$PATH\"\n", BinDir)
}
