package ui

import (
	"io"
	"strings"

	"github.com/pterm/pterm"
)

var Logger *pterm.Logger

func NewHandler(w io.Writer, format string, level string) *pterm.SlogHandler {
	Logger = &pterm.DefaultLogger
	Logger.Writer = w
	Logger.Level = ParseLevel(level)
	if strings.ToLower(format) == "json" {
		Logger.Formatter = pterm.LogFormatterJSON
	}
	return pterm.NewSlogHandler(Logger)
}

func ParseLevel(level string) pterm.LogLevel {
	logLevelMap := map[string]pterm.LogLevel{
		"debug": pterm.LogLevelDebug,
		"info":  pterm.LogLevelInfo,
		"warn":  pterm.LogLevelWarn,
		"error": pterm.LogLevelError,
	}
	if level, ok := logLevelMap[strings.ToLower(level)]; ok {
		return level
	}
	return pterm.LogLevelInfo
}
