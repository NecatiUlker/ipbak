package update

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"ipbak/internal/config"
)

const (
	geoipConfTemplate = `
AccountID {{.AccountID}}
LicenseKey {{.LicenseKey}}
EditionIDs GeoLite2-ASN GeoLite2-City GeoLite2-Country
`
)

type Metadata struct {
	LastUpdate time.Time `json:"last_update"`
}

func GenerateGeoIPConf(accountID, licenseKey string) error {
	path, err := config.GetGeoIPConfPath()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	tmpl, err := template.New("geoip.conf").Parse(geoipConfTemplate)
	if err != nil {
		return err
	}

    file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	data := struct {
		AccountID  string
		LicenseKey string
	}{
		AccountID:  accountID,
		LicenseKey: licenseKey,
	}

	return tmpl.Execute(file, data)
}

func RunGeoIPUpdate(force bool) error {
	dbDir, err := config.GetDatabaseDir()
	if err != nil {
		return err
	}

	// Ensure database directory exists
	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return err
	}

	if !force {
		metaPath := filepath.Join(dbDir, "metadata.json")
		if data, err := os.ReadFile(metaPath); err == nil {
			var meta Metadata
			if err := json.Unmarshal(data, &meta); err == nil {
				if time.Since(meta.LastUpdate) < 24*time.Hour {
					return fmt.Errorf("database is up to date (last update: %s), use --force to override", meta.LastUpdate.Format(time.RFC3339))
				}
			}
		}
	}

	confPath, err := config.GetGeoIPConfPath()
	if err != nil {
		return err
	}

	// Check if GeoIP.conf exists
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		return fmt.Errorf("GeoIP.conf not found at %s. Please run 'ipbak setup' first", confPath)
	}

	var cmdName string
	if path, err := EnsureGeoIPUpdate(); err == nil {
		cmdName = path
	} else {
		return fmt.Errorf("geoipupdate not found and download failed: %w", err)
	}

	cmd := exec.Command(cmdName, "-d", dbDir, "-f", confPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("geoipupdate failed: %w", err)
	}

	// Update metadata
	meta := Metadata{
		LastUpdate: time.Now(),
	}
	data, err := json.Marshal(meta)
	if err == nil {
		os.WriteFile(filepath.Join(dbDir, "metadata.json"), data, 0600)
	}

	return nil
}

func findGeoIPUpdate() (string, error) {
	// 1. Check current working directory
	if _, err := os.Stat("geoipupdate"); err == nil {
		return filepath.Abs("geoipupdate")
	}
	if _, err := os.Stat("geoipupdate.exe"); err == nil {
		return filepath.Abs("geoipupdate.exe")
	}

	// 2. Check executable directory
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		if _, err := os.Stat(filepath.Join(exeDir, "geoipupdate")); err == nil {
			return filepath.Join(exeDir, "geoipupdate"), nil
		}
		if _, err := os.Stat(filepath.Join(exeDir, "geoipupdate.exe")); err == nil {
			return filepath.Join(exeDir, "geoipupdate.exe"), nil
		}
	}

	// 3. Check PATH
	return exec.LookPath("geoipupdate")
}

func GetMetadata() (*Metadata, error) {
	dbDir, err := config.GetDatabaseDir()
	if err != nil {
		return nil, err
	}
	metaPath := filepath.Join(dbDir, "metadata.json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, err
	}
	var meta Metadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}
