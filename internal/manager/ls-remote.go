package manager

import (
	"fmt"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

// lsRemoteCmd represents the ls-remote command
var lsRemoteCmd = &cobra.Command{
	Use:     "ls-remote",
	Aliases: []string{"list-remote"},
	Short:   "List available Node.js versions",
	Long:    `List all available Node.js versions that can be installed.`,
	Run: func(cmd *cobra.Command, args []string) {
		ltsOnly, _ := cmd.Flags().GetBool("lts")
		all, _ := cmd.Flags().GetBool("all")
		listRemoteVersions(ltsOnly, all)
	},
}

func init() {
	lsRemoteCmd.Flags().BoolP("lts", "l", false, "List only LTS versions")
	lsRemoteCmd.Flags().BoolP("all", "a", false, "Show all versions")
}

// listRemoteVersions lists all available Node.js versions from the remote server
func listRemoteVersions(ltsOnly bool, showAll bool) {
	fmt.Println("Fetching available Node.js versions...")

	versions, err := FetchAvailableVersions()
	if err != nil {
		fmt.Printf("Error fetching available versions: %v\n", err)
		return
	}

	// Sort versions by date (newest first)
	sort.Slice(versions, func(i, j int) bool {
		// Parse dates - assume YYYY-MM-DD format
		// If parsing fails, use version string comparison as fallback
		iDate, iErr := time.Parse("2006-01-02", versions[i].Date)
		jDate, jErr := time.Parse("2006-01-02", versions[j].Date)

		if iErr != nil || jErr != nil {
			return versions[i].Version > versions[j].Version
		}

		return iDate.After(jDate)
	})

	fmt.Println("Available Node.js versions:")

	now := time.Now()

	count := 0
	for _, v := range versions {
		// Skip non-LTS versions if ltsOnly is true
		isLTS := IsLTS(v.LTS)
		if ltsOnly && !isLTS {
			continue
		}

		// Check if this version is installed
		installed := ""
		if IsVersionInstalled(v.Version) {
			installed = " (installed)"
		}

		// Check if this is an LTS version
		ltsStatus := ""
		if isLTS {
			ltsStatus = " [LTS]"
		}

		// Add "recent" tag for versions released in the last 30 days
		recent := ""
		releaseDate, err := time.Parse("2006-01-02", v.Date)
		if err == nil {
			if now.Sub(releaseDate).Hours() < 24*30 {
				recent = " [recent]"
			}
		}

		// Add security flag
		security := ""
		if v.Security {
			security = " [security]"
		}

		fmt.Printf("  %s%s%s%s%s (%s)\n",
			v.Version,
			ltsStatus,
			recent,
			security,
			installed,
			v.Date)

		count++

		// Limit to 20 versions to avoid overwhelming output
		if count >= 20 && !showAll {
			fmt.Println("  ...")
			fmt.Println("  (Use --all to show all versions)")
			break
		}
	}
}
