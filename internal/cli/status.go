package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"ipbak/internal/update"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show database status",
	Run: func(cmd *cobra.Command, args []string) {
		meta, err := update.GetMetadata()
		if err != nil {
			fmt.Println("Status: No metadata found or error reading metadata.")
			fmt.Println("Run 'ipbak setup' or 'ipbak update' to initialize.")
			return
		}

		fmt.Printf("Last Update: %s\n", meta.LastUpdate.Format(time.RFC1123))
		
		nextUpdate := meta.LastUpdate.Add(24 * time.Hour)
		timeUntil := time.Until(nextUpdate)
		
		if timeUntil > 0 {
			fmt.Printf("Next eligible update: %s (in %s)\n", nextUpdate.Format(time.RFC1123), timeUntil.Round(time.Minute))
		} else {
			fmt.Println("Next eligible update: Now")
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
