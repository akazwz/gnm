package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	// BaseURL is the Node.js distribution URL
	BaseURL = "https://nodejs.org/dist/"
	// VersionListURL is the URL to fetch available Node.js versions
	VersionListURL = "https://nodejs.org/dist/index.json"
	// HomeDir is the directory name for GNM in the user's home directory
	HomeDir = ".gnm"
)

var (
	// GnmDir is the full path to the GNM home directory
	GnmDir string
	// VersionsDir is the directory where Node.js versions are installed
	VersionsDir string
	// BinDir is the directory with symlinks to the active Node.js version
	BinDir string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "gnm",
	Short: "Go Node Manager (gnm) - A Node.js version manager written in Go",
	Long: `Go Node Manager (gnm) is a simple and fast Node.js version manager 
written in Go. It allows you to install, manage, and switch between 
multiple Node.js versions.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Initialize directories
	initDirs()

	// Add commands to root command
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(lsRemoteCmd)
	rootCmd.AddCommand(uninstallCmd)
}

// initDirs initializes the directory structure for GNM
func initDirs() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}

	GnmDir = filepath.Join(home, HomeDir)
	VersionsDir = filepath.Join(GnmDir, "versions")
	BinDir = filepath.Join(GnmDir, "bin")

	// Create directories if they don't exist
	os.MkdirAll(VersionsDir, 0755)
	os.MkdirAll(BinDir, 0755)
}
