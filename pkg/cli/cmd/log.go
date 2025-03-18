package cmd

import (
	"log/slog"
	"strings"

	"github.com/neticdk/go-common/pkg/slice"
)

type LogFormat string

const (
	LogFormatPlain LogFormat = "plain"
	LogFormatJSON  LogFormat = "json"

	LogFormatDefault = LogFormatPlain
)

// String returns the string representation of the format
func (f LogFormat) String() string {
	return string(f)
}

// AllLogFormats returns all formats
func AllLogFormats() []LogFormat {
	return []LogFormat{
		LogFormatPlain,
		LogFormatJSON,
	}
}

// AllLogFormatsStr returns all formats as strings
func AllLogFormatsStr() []string {
	return slice.Map(AllLogFormats(), func(f LogFormat) string {
		return f.String()
	})
}

// AllLogFormatsJoined returns all formats joined by "|"
func AllLogFormatsJoined() string {
	return strings.Join([]string{
		LogFormatPlain.String(),
		LogFormatJSON.String(),
	}, "|")
}

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"

	LogLevelDefault = LogLevelInfo
)

// String returns the string representation of the level
func (l LogLevel) String() string {
	return string(l)
}

// AllLogLevels returns all levels
func AllLogLevels() []LogLevel {
	return []LogLevel{
		LogLevelDebug,
		LogLevelInfo,
		LogLevelWarn,
		LogLevelError,
	}
}

// AllLogLevelsStr returns all levels as strings
func AllLogLevelsStr() []string {
	return slice.Map(AllLogLevels(), func(l LogLevel) string {
		return l.String()
	})
}

// AllLogLevelsJoined returns all levels joined by "|"
func AllLogLevelsJoined() string {
	return strings.Join([]string{
		LogLevelDebug.String(),
		LogLevelInfo.String(),
		LogLevelWarn.String(),
		LogLevelError.String(),
	}, "|")
}

// ParseLogLevel converts a log level to a slog level
func ParseLogLevel(level LogLevel) slog.Level {
	logLevelMap := map[LogLevel]slog.Level{
		LogLevelDebug: slog.LevelDebug,
		LogLevelInfo:  slog.LevelInfo,
		LogLevelWarn:  slog.LevelWarn,
		LogLevelError: slog.LevelError,
	}
	if level, ok := logLevelMap[level]; ok {
		return level
	}
	return slog.LevelInfo
}
