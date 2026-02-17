package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"ipbak/internal/config"
	"ipbak/internal/geo"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check the health of the installation",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking ipbak health...")
		fmt.Println("-----------------------")

		// 1. Check Config
		configPath, err := config.GetConfigPath()
		if err == nil {
			if _, err := os.Stat(configPath); err == nil {
				fmt.Printf("[OK] Config file found: %s\n", configPath)
				// Check permissions? Windows permissions are tricky, maybe skip strict check for now or basic check.
			} else {
				fmt.Printf("[FAIL] Config file missing: %s. Run 'ipbak setup'\n", configPath)
			}
		} else {
			fmt.Printf("[FAIL] Could not determine config path: %v\n", err)
		}

		// 2. Check GeoIP Database Directory
		dbDir, err := config.GetDatabaseDir()
		if err == nil {
			if _, err := os.Stat(dbDir); err == nil {
				fmt.Printf("[OK] Database directory found: %s\n", dbDir)
				
				// Check for MMDB files
				asnPath := filepath.Join(dbDir, "GeoLite2-ASN.mmdb")
				if _, err := os.Stat(asnPath); err == nil {
					fmt.Printf("[OK] GeoLite2-ASN.mmdb found\n")
				} else {
					fmt.Printf("[FAIL] GeoLite2-ASN.mmdb missing\n")
				}

				cityPath := filepath.Join(dbDir, "GeoLite2-City.mmdb")
				if _, err := os.Stat(cityPath); err == nil {
					fmt.Printf("[OK] GeoLite2-City.mmdb found\n")
				} else {
					fmt.Printf("[FAIL] GeoLite2-City.mmdb missing\n")
				}
			} else {
				fmt.Printf("[FAIL] Database directory missing: %s\n", dbDir)
			}
		} else {
			fmt.Printf("[FAIL] Could not determine database directory: %v\n", err)
		}

		// 3. Check geoipupdate
		path, err := exec.LookPath("geoipupdate")
		if err == nil {
			fmt.Printf("[OK] geoipupdate found: %s\n", path)
		} else {
			fmt.Printf("[WARN] geoipupdate not found in PATH. Updates will fail.\n")
		}

		// 4. Try loading GeoService
		service, err := geo.NewGeoService()
		if err == nil {
			fmt.Printf("[OK] GeoService initialized successfully\n")
			service.Close()
		} else {
			fmt.Printf("[FAIL] GeoService initialization failed: %v\n", err)
		}

		fmt.Println("-----------------------")
		fmt.Println("Doctor check complete.")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
