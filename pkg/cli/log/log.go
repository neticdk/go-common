package log

import (
	"log/slog"
	"strings"

	"github.com/neticdk/go-common/pkg/slices"
)

type Format string

const (
	FormatPlain Format = "plain"
	FormatJSON  Format = "json"

	DefaultFormat = FormatPlain
)

// String returns the string representation of the format
func (f Format) String() string {
	return string(f)
}

// AllFormats returns all formats
func AllFormats() []Format {
	return []Format{
		FormatPlain,
		FormatJSON,
	}
}

// AllFormatsStr returns all formats as strings
func AllFormatsStr() []string {
	return slices.Map(AllFormats(), func(f Format) string {
		return f.String()
	})
}

// AllFormatsJoined returns all formats joined by "|"
func AllFormatsJoined() string {
	return strings.Join([]string{
		FormatPlain.String(),
		FormatJSON.String(),
	}, "|")
}

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"

	DefaultLevel = LevelInfo
)

// String returns the string representation of the level
func (l Level) String() string {
	return string(l)
}

// AllLevels returns all levels
func AllLevels() []Level {
	return []Level{
		LevelDebug,
		LevelInfo,
		LevelWarn,
		LevelError,
	}
}

// AllLevelsStr returns all levels as strings
func AllLevelsStr() []string {
	return slices.Map(AllLevels(), func(l Level) string {
		return l.String()
	})
}

// AllLevelsJoined returns all levels joined by "|"
func AllLevelsJoined() string {
	return strings.Join([]string{
		LevelDebug.String(),
		LevelInfo.String(),
		LevelWarn.String(),
		LevelError.String(),
	}, "|")
}

// ParseLevel converts a log level to a slog level
func ParseLevel(level Level) slog.Level {
	logLevelMap := map[Level]slog.Level{
		LevelDebug: slog.LevelDebug,
		LevelInfo:  slog.LevelInfo,
		LevelWarn:  slog.LevelWarn,
		LevelError: slog.LevelError,
	}
	if level, ok := logLevelMap[level]; ok {
		return level
	}
	return slog.LevelInfo
}
