package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"ipbak/internal/geo"
	"ipbak/internal/output"
)

var batchCmd = &cobra.Command{
	Use:   "batch [file]",
	Short: "Process a batch of IPs from a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		service, err := geo.NewGeoService()
		if err != nil {
			return fmt.Errorf("%w\n\nRun 'ipbak setup' to initialize the database.", err)
		}
		defer service.Close()

		var printer output.Printer
		if jsonOutput {
			printer = &output.JsonPrinter{}
		} else {
			printer = &output.TextPrinter{Advanced: advanced}
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			ip := strings.TrimSpace(scanner.Text())
			if ip == "" {
				continue
			}

			if err := processIP(service, ip, printer); err != nil {
				// Log error to stderr but continue processing
				fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", ip, err)
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(batchCmd)
}
