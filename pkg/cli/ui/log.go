package ui

import (
	"io"

	"github.com/neticdk/go-common/pkg/cli/log"
	"github.com/pterm/pterm"
)

var Logger *pterm.Logger

func NewHandler(w io.Writer, format log.Format, level log.Level) *pterm.SlogHandler {
	Logger = &pterm.DefaultLogger
	Logger.Writer = w
	Logger.Level = ParseLevel(level)
	if format == log.FormatJSON {
		Logger.Formatter = pterm.LogFormatterJSON
	}
	return pterm.NewSlogHandler(Logger)
}

func ParseLevel(level log.Level) pterm.LogLevel {
	logLevelMap := map[log.Level]pterm.LogLevel{
		log.LevelDebug: pterm.LogLevelDebug,
		log.LevelInfo:  pterm.LogLevelInfo,
		log.LevelWarn:  pterm.LogLevelWarn,
		log.LevelError: pterm.LogLevelError,
	}
	if level, ok := logLevelMap[level]; ok {
		return level
	}
	return pterm.LogLevelInfo
}
