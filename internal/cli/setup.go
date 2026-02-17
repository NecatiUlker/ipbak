package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"ipbak/internal/config"
	"ipbak/internal/update"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup for ipbak",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("GeoLite2 database setup")
		fmt.Println("This tool uses the free GeoLite2 database provided by MaxMind.")
		fmt.Println("")
		fmt.Println("1) Create a free account:")
		fmt.Println("   https://www.maxmind.com")
		fmt.Println("")
		fmt.Println("2) Generate a license key:")
		fmt.Println("   https://www.maxmind.com/en/accounts/license-key")
		fmt.Println("")

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter MaxMind Account ID: ")
		accountID, _ := reader.ReadString('\n')
		accountID = strings.TrimSpace(accountID)

		fmt.Print("Enter MaxMind License Key: ")
		licenseKey, _ := reader.ReadString('\n')
		licenseKey = strings.TrimSpace(licenseKey)

		if accountID == "" || licenseKey == "" {
			return fmt.Errorf("Account ID and License Key are required")
		}

		cfg := &config.Config{
			MaxMindAccountID:  accountID,
			MaxMindLicenseKey: licenseKey,
		}

		if err := config.SaveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		if err := update.GenerateGeoIPConf(accountID, licenseKey); err != nil {
			return fmt.Errorf("failed to generate GeoIP.conf: %w", err)
		}

		fmt.Println("\nConfiguration saved.")
		fmt.Println("Downloading databases... (this may take a moment)")

		if err := update.RunGeoIPUpdate(true); err != nil {
             fmt.Println("\nError downloading databases.")
             fmt.Println("Please ensure 'geoipupdate' is installed and your credentials are correct.")
             fmt.Println("Install instructions: https://github.com/maxmind/geoipupdate#installation")
			return err
		}

		fmt.Println("Setup complete! You can now use ipbak.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
