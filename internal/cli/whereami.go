package cli

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var whereAmICmd = &cobra.Command{
	Use:   "whereami",
	Short: "Find your public IP address and location",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := http.Get("https://api.ipify.org?format=text")
		if err != nil {
			return fmt.Errorf("failed to get public IP: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		ip := strings.TrimSpace(string(body))
		return runQuery(ip)
	},
}

func init() {
	rootCmd.AddCommand(whereAmICmd)
}
