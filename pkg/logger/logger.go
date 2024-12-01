package logger

import (
	log "github.com/sirupsen/logrus"
)

func AllLevelsString() (levels []string) {
	for _, level := range log.AllLevels {
		levels = append(levels, level.String())
	}
	return
}
