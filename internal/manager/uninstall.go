package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:     "uninstall [version]",
	Aliases: []string{"remove"},
	Short:   "Uninstall a specific Node.js version",
	Long:    `Uninstall a specific version of Node.js.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		version := args[0]
		uninstallVersion(version)
	},
}

// uninstallVersion uninstalls the specified Node.js version
func uninstallVersion(version string) {
	version = NormalizeVersion(version)

	// Check if version is installed
	if !IsVersionInstalled(version) {
		fmt.Printf("Node.js %s is not installed\n", version)
		return
	}

	// Check if this is the current version
	currentVersion, err := GetCurrentVersion()
	if err == nil && currentVersion == version {
		fmt.Printf("Cannot uninstall the currently active version. Switch to another version first.\n")
		return
	}

	// Remove the version directory
	versionDir := filepath.Join(VersionsDir, version)
	if err := os.RemoveAll(versionDir); err != nil {
		fmt.Printf("Error removing directory %s: %v\n", versionDir, err)
		return
	}

	fmt.Printf("Node.js %s uninstalled successfully\n", version)
}
