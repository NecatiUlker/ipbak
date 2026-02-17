package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"ipbak/internal/update"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the GeoLite2 database",
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		if err := update.RunGeoIPUpdate(force); err != nil {
			return err
		}
		fmt.Println("Database update successful.")
		return nil
	},
}

func init() {
	updateCmd.Flags().Bool("force", false, "Force update even if cooldown is active")
	rootCmd.AddCommand(updateCmd)
}
