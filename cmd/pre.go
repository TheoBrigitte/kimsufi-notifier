package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// logLevel set the level of the logger.
func logLevel(cmd *cobra.Command, args []string) error {
	level, err := log.ParseLevel(cmd.Flag("log-level").Value.String())
	if err != nil {
		return err
	}
	log.SetLevel(level)

	return nil
}
