package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/TheoBrigitte/kimsufi-notifier/cmd/flag"
)

// logLevel set the logger log level using the value of the flag
func logLevel(cmd *cobra.Command, args []string) error {
	level, err := log.ParseLevel(cmd.Flag(flag.LogLevelFlagName).Value.String())
	if err != nil {
		return err
	}
	log.SetLevel(level)

	return nil
}
