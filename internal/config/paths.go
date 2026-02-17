package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetConfigDir() (string, error) {
	var configDir string
	var err error

	switch runtime.GOOS {
	case "windows":
		configDir, err = os.UserConfigDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(configDir, "ipbak")
	default:
		// Linux/macOS
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(homeDir, ".config", "ipbak")
	}

	return configDir, nil
}

func GetConfigPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func GetDatabaseDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			return "", fmt.Errorf("LOCALAPPDATA environment variable not set")
		}
		return filepath.Join(localAppData, "ipbak"), nil
	default:
		// Linux/macOS
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, ".local", "share", "ipbak"), nil
	}
}

func GetGeoIPConfPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "GeoIP.conf"), nil
}
