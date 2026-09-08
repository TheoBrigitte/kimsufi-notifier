package logger

import (
	"fmt"
	"maps"
	"slices"

	log "github.com/sirupsen/logrus"
)

func AllLevelsString() (levels []string) {
	for _, level := range log.AllLevels {
		levels = append(levels, level.String())
	}
	return
}

type logFormat string

const (
	LogFormatText logFormat = "text"
	LogFormatJSON logFormat = "json"
)

func (l logFormat) String() string {
	return string(l)
}

var formats = map[string]logFormat{
	"text": LogFormatText,
	"json": LogFormatJSON,
}

func NewFormatter(format string) (log.Formatter, error) {
	logFormat, ok := formats[format]
	if !ok {
		return nil, fmt.Errorf("invalid log format: %q", format)
	}

	switch logFormat {
	case LogFormatText:
		return &log.TextFormatter{}, nil
	case LogFormatJSON:
		return &log.JSONFormatter{}, nil
	}

	return nil, fmt.Errorf("unsupported log format: %q", format)
}

func AllFormatsString() []string {
	f := slices.Collect(maps.Keys(formats))
	slices.Sort(f)
	return f
}
