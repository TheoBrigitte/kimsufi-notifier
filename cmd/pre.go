package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/TheoBrigitte/kimsufi-notifier/cmd/flag"
	"github.com/TheoBrigitte/kimsufi-notifier/pkg/logger"
)

func commandPreRunE(cmd *cobra.Command, args []string) error {
	if err := setLogLevel(cmd, args); err != nil {
		return err
	}

	if err := setLogFormat(cmd, args); err != nil {
		return err
	}

	return nil
}

// logLevel set the logger log level using the value of the flag
func setLogLevel(cmd *cobra.Command, args []string) error {
	level, err := log.ParseLevel(cmd.Flag(flag.LogLevelFlagName).Value.String())
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	log.SetLevel(level)

	return nil
}

// logFormat set the logger format using the value of the flag
func setLogFormat(cmd *cobra.Command, args []string) error {
	logFormatter, err := logger.NewFormatter(cmd.Flag(flag.LogFormatFlagName).Value.String())
	if err != nil {
		return err
	}

	log.SetFormatter(logFormatter)

	return nil
}
