package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"ipbak/internal/geo"
	"ipbak/internal/output"
)

var (
	jsonOutput bool
	advanced   bool
)

var rootCmd = &cobra.Command{
	Use:   "ipbak [ip]",
	Short: "Professional security-grade IP intelligence utility",
	Long: `ipbak is a production-grade, cross-platform CLI tool for IP intelligence.
It uses MaxMind GeoLite2 databases to provide detailed information about IP addresses.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return runQuery(args[0])
	},
	SilenceUsage: true, // Don't show usage on error
    SilenceErrors: true, // we handle errors
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&advanced, "advanced", false, "Show advanced details")
}

func runQuery(ip string) error {
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

	return processIP(service, ip, printer)
}

func processIP(service *geo.GeoService, ip string, printer output.Printer) error {
	loc, err := service.Lookup(ip)
	if err != nil {
		// Instead of failing the whole batch, maybe print error to stderr?
		// For single query, we want to return error.
		// For batch, we probably want to continue.
		// Let's return error and let caller decide.
		return err
	}
	return printer.Print(loc)
}
