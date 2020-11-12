package cmd

import (
	"github.com/spf13/cobra"

	"github.com/TheoBrigitte/kimsufi-notifier/cmd/check"
	"github.com/TheoBrigitte/kimsufi-notifier/cmd/order"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               "kimsufi-notifier",
	Short:             "kimsufi availability notifier",
	Long:              `Send notification when kimsufi server are available.`,
	PersistentPreRunE: logLevel,
	SilenceUsage:      true,
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(check.Cmd)
	rootCmd.AddCommand(order.Cmd)
}
