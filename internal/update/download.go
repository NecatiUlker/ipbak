package update

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"ipbak/internal/config"
)

const (
	geoipUpdateVersion = "7.0.1"
	githubReleasesURL  = "https://github.com/maxmind/geoipupdate/releases/download/v%s/%s"
)

func EnsureGeoIPUpdate() (string, error) {
	// First, check if we already have it
	path, err := findGeoIPUpdate()
	if err == nil {
		return path, nil
	}

	// If not found, download it
	fmt.Println("geoipupdate not found. Downloading automatically...")
	return downloadGeoIPUpdate()
}

func downloadGeoIPUpdate() (string, error) {
	dbDir, err := config.GetDatabaseDir()
	if err != nil {
		return "", err
	}

	// Ensure directory exists
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return "", err
	}

	osStr := runtime.GOOS
	archStr := runtime.GOARCH

	var fileName string
	var extractFunc func(string, string) (string, error)

	if osStr == "windows" {
		fileName = fmt.Sprintf("geoipupdate_%s_%s_%s.zip", geoipUpdateVersion, osStr, archStr)
		extractFunc = extractZip
	} else {
		fileName = fmt.Sprintf("geoipupdate_%s_%s_%s.tar.gz", geoipUpdateVersion, osStr, archStr)
		extractFunc = extractTarGz
	}

	url := fmt.Sprintf(githubReleasesURL, geoipUpdateVersion, fileName)
	tmpFile := filepath.Join(dbDir, fileName)

	if err := downloadFile(url, tmpFile); err != nil {
		return "", fmt.Errorf("failed to download geoipupdate: %w", err)
	}
	defer os.Remove(tmpFile) // Clean up archive

	exePath, err := extractFunc(tmpFile, dbDir)
	if err != nil {
		return "", fmt.Errorf("failed to extract geoipupdate: %w", err)
	}

    // Make executable on Linux/macOS
    if osStr != "windows" {
        if err := os.Chmod(exePath, 0755); err != nil {
            return "", err
        }
    }

	fmt.Printf("geoipupdate downloaded to: %s\n", exePath)
	return exePath, nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractZip(src, destDir string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var exePath string

	for _, f := range r.File {
		// We encounter a directory structure inside the zip, e.g., geoipupdate_7.0.1_windows_amd64/geoipupdate.exe
		// We want to extract just the executable to destDir
		baseName := filepath.Base(f.Name)
		if baseName == "geoipupdate.exe" {
			params := f.Name
            // sanitize path to prevent Zip Slip (although we filter rigidly)
            if strings.Contains(params, "..") {
                continue
            }
            
			exePath = filepath.Join(destDir, baseName)
			outFile, err := os.OpenFile(exePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return "", err
			}
			
			rc, err := f.Open()
			if err != nil {
				outFile.Close()
				return "", err
			}

			_, err = io.Copy(outFile, rc)
			outFile.Close()
			rc.Close()
            break 
		}
	}
    
    if exePath == "" {
        return "", fmt.Errorf("geoipupdate.exe not found in zip")
    }

	return exePath, nil
}

func extractTarGz(src, destDir string) (string, error) {
	f, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	var exePath string

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		baseName := filepath.Base(header.Name)
		if baseName == "geoipupdate" {
			exePath = filepath.Join(destDir, baseName)
			outFile, err := os.OpenFile(exePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return "", err
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return "", err
			}
			outFile.Close()
            break
		}
	}
    
    if exePath == "" {
        return "", fmt.Errorf("geoipupdate binary not found in tar.gz")
    }

	return exePath, nil
}
