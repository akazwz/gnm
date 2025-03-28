package manager

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// NodeVersion represents a Node.js version from the remote API
type NodeVersion struct {
	Version  string      `json:"version"`
	Date     string      `json:"date"`
	Files    []string    `json:"files"`
	LTS      interface{} `json:"lts"`
	Security bool        `json:"security"`
}

// NormalizeVersion ensures the version string has a "v" prefix
func NormalizeVersion(version string) string {
	if !strings.HasPrefix(version, "v") {
		return "v" + version
	}
	return version
}

// GetNodeArch converts Go architecture to Node.js architecture
func GetNodeArch() string {
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		return "x64"
	case "386":
		return "x86"
	case "arm64":
		return "arm64"
	default:
		return arch
	}
}

// GetCurrentVersion returns the currently active Node.js version
func GetCurrentVersion() (string, error) {
	nodePath, err := os.Readlink(filepath.Join(BinDir, "node"))
	if err != nil {
		return "", err
	}

	// Extract version from path
	// The path will be something like: .gnm/versions/v16.14.0/bin/node
	parts := strings.Split(nodePath, string(os.PathSeparator))
	for i, part := range parts {
		if strings.HasPrefix(part, "v") && i+1 < len(parts) && parts[i+1] == "bin" {
			return part, nil
		}
	}

	return "", fmt.Errorf("couldn't determine current version")
}

// DownloadFile downloads a file from the given URL to the specified path
func DownloadFile(url, destPath string) error {
	// Create the output file
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Send GET request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Copy the response body to the output file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// ExtractTarGz extracts a .tar.gz file to the specified directory
func ExtractTarGz(tarballPath, destDir string) error {
	file, err := os.Open(tarballPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	// Create destination directory
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// Find the root directory in the tarball
	var rootDir string

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		parts := strings.SplitN(header.Name, "/", 2)
		if len(parts) > 0 && rootDir == "" {
			rootDir = parts[0]
		}
	}

	// Reset the tarball reader
	file.Seek(0, 0)
	gzr, _ = gzip.NewReader(file)
	tr = tar.NewReader(gzr)

	// Extract files
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Skip the root directory
		if header.Name == rootDir+"/" {
			continue
		}

		// Remove the root directory from the path
		relPath := strings.TrimPrefix(header.Name, rootDir+"/")
		if relPath == "" {
			continue
		}

		target := filepath.Join(destDir, relPath)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			dir := filepath.Dir(target)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			f.Close()
		case tar.TypeSymlink:
			dir := filepath.Dir(target)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}
			os.Symlink(header.Linkname, target)
		}
	}

	return nil
}

// FetchAvailableVersions fetches the list of available Node.js versions
func FetchAvailableVersions() ([]NodeVersion, error) {
	// Send GET request
	resp, err := http.Get(VersionListURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Decode JSON
	var versions []NodeVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, err
	}

	return versions, nil
}

// IsVersionInstalled checks if a Node.js version is installed
func IsVersionInstalled(version string) bool {
	version = NormalizeVersion(version)
	versionDir := filepath.Join(VersionsDir, version)

	_, err := os.Stat(versionDir)
	return err == nil
}

// IsLTS checks if a version is LTS
func IsLTS(lts interface{}) bool {
	// LTS can be a string codename or boolean
	switch v := lts.(type) {
	case bool:
		return v
	case string:
		return v != ""
	default:
		return false
	}
}
