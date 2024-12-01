package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/TheoBrigitte/kimsufi-notifier/cmd/check"
	"github.com/TheoBrigitte/kimsufi-notifier/cmd/flag"
	"github.com/TheoBrigitte/kimsufi-notifier/cmd/list"
	"github.com/TheoBrigitte/kimsufi-notifier/cmd/order"
	"github.com/TheoBrigitte/kimsufi-notifier/cmd/version"
)

// rootCmd represents the base command when called without any arguments
var rootCmd = &cobra.Command{
	Use:               "kimsufi-notifier",
	Short:             "kimsufi availability notifier",
	Long:              "List, check availability and order OVH Eco (including kimsufi) servers.",
	PersistentPreRunE: logLevel,
	SilenceUsage:      true,
}

// init registers all subcommands and global flags
func init() {
	// Global flags
	flag.Bind(rootCmd)

	// Subcommands
	rootCmd.AddCommand(check.Cmd)
	rootCmd.AddCommand(order.Cmd)
	rootCmd.AddCommand(list.Cmd)
	rootCmd.AddCommand(version.Cmd)
}

// Execute is the main entry point for the CLI
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
