package manager

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List installed Node.js versions",
	Long:    `List all installed Node.js versions, showing the currently active version.`,
	Run: func(cmd *cobra.Command, args []string) {
		listInstalledVersions()
	},
}

// listInstalledVersions lists all installed Node.js versions
func listInstalledVersions() {
	entries, err := os.ReadDir(VersionsDir)
	if err != nil {
		fmt.Printf("Error reading versions directory: %v\n", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("No Node.js versions installed")
		return
	}

	// Get current version
	currentVersion, err := GetCurrentVersion()
	if err != nil {
		currentVersion = ""
	}

	fmt.Println("Installed Node.js versions:")
	for _, entry := range entries {
		if entry.IsDir() {
			version := entry.Name()
			active := ""
			if version == currentVersion {
				active = " (active)"
			}

			// Check if this is an LTS version
			ltsStatus := ""
			availableVersions, _ := FetchAvailableVersions()
			for _, v := range availableVersions {
				if strings.TrimPrefix(version, "v") == strings.TrimPrefix(v.Version, "v") && IsLTS(v.LTS) {
					ltsStatus = " (LTS)"
					break
				}
			}

			fmt.Printf("  %s%s%s\n", version, ltsStatus, active)
		}
	}
}
