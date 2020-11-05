package cmd

import (
	log "github.com/sirupsen/logrus"
)

func init() {
	rootCmd.PersistentFlags().StringP("log-level", "l", log.InfoLevel.String(), "log level")
}
