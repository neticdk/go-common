package logger

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const (
	LogLevelInfo  = "info"
	LogLevelDebug = "debug"
	LogLevelTrace = "trace"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

var (
	iconBaseStyle = lipgloss.NewStyle().
			Bold(true).
			MaxWidth(4)

	debugLevelStyle = lipgloss.NewStyle().
			SetString("🪲").
			Foreground(lipgloss.Color("63")).
			Inherit(iconBaseStyle)

	infoLevelStyle = lipgloss.NewStyle().
			Inherit(iconBaseStyle).
			SetString("🦶").
			Foreground(lipgloss.Color("86"))

	warnLevelStyle = lipgloss.NewStyle().
			SetString("⚠️").
			Foreground(lipgloss.Color("192")).
			Inherit(iconBaseStyle)

	errorLevelStyle = lipgloss.NewStyle().
			SetString("✗").
			Foreground(lipgloss.Color("204")).
			Inherit(iconBaseStyle)

	fatalLevelStyle = lipgloss.NewStyle().
			SetString("💀").
			Foreground(lipgloss.Color("134")).
			Inherit(iconBaseStyle)
)

// Logger defines an abstract logger that can be used to log to the output
type Logger interface {
	// Set the logger level
	SetLevel(level string) error

	// Level returns the logger level
	Level() string

	// Set the logger output
	SetOutput(w io.Writer)

	// Output returns the logger output
	Output() io.Writer

	// Print prints a log message
	Print(message string, keyvals ...any)
	// Info logs to info level
	Info(message string, keyvals ...any)
	// Debug logs to debug level
	Debug(message string, keyvals ...any)
	// Error logs to error level
	Error(message string, keyvals ...any)
	// Fatal logs to fatal level
	Fatal(message string, keyvals ...any)
	// Warn logs to warn level
	Warn(message string, keyvals ...any)
	// Trace logs to trace level
	Trace(message string, keyvals ...any)
	// Infof logs formatted info level
	Infof(format string, keyvals ...any)
	// Debugf logs formatted debug level
	Debugf(format string, keyvals ...any)

	// StandardWriter returns a writer that can be used to write logs
	StandardWriter() io.Writer

	// IsLevel returns true if the logger is set to the given level
	IsInfo() bool
	// IsDebug returns true if the logger is set to debug level
	IsDebug() bool
	// IsError returns true if the logger is set to error level
	IsError() bool
	// IsTrace returns true if the logger is set to trace level
	IsTrace() bool
	// IsWarn returns true if the logger is set to warn level
	IsWarn() bool

	// WithPrefix returns a new logger with the given prefix
	WithPrefix(string) Logger

	// SetInteractive sets the logger to use fancy styles if isTerminal is true
	SetInteractive(string, bool)
}

// CharmLogger is a charms/lipgloss based logger
type CharmLogger struct {
	internal *log.Logger
	writer   io.Writer
	level    string
}

// New creates a new CharmLogger
func New(w io.Writer, level string) Logger {
	l := log.New(w)
	if parsedLevel, err := log.ParseLevel(level); err == nil {
		l.SetLevel(parsedLevel)
	}

	return &CharmLogger{l, w, level}
}

// WithPrefix returns a new logger with the given prefix
func (l *CharmLogger) WithPrefix(prefix string) Logger {
	newLogger := l.internal.WithPrefix(prefix)
	return &CharmLogger{newLogger, l.writer, l.level}
}

// SetOutput sets the logger output
func (l *CharmLogger) SetOutput(w io.Writer) {
	l.writer = w
	l.internal.SetOutput(w)
}

// Output returns the logger output
func (l *CharmLogger) Output() io.Writer {
	return l.writer
}

// IsLevel returns true if the logger is set to the given level
func (l *CharmLogger) IsInfo() bool {
	return l.level == LogLevelInfo
}

// IsDebug returns true if the logger is set to debug level
func (l *CharmLogger) IsDebug() bool {
	return l.level == LogLevelDebug
}

// IsError returns true if the logger is set to error level
func (l *CharmLogger) IsError() bool {
	return l.level == LogLevelError
}

// IsWarn returns true if the logger is set to warn level
func (l *CharmLogger) IsWarn() bool {
	return l.level == LogLevelWarn
}

// IsTrace returns true if the logger is set to trace level
func (l *CharmLogger) IsTrace() bool {
	return l.level == LogLevelTrace
}

// SetLevel sets the logger level
func (l *CharmLogger) SetLevel(level string) error {
	l.level = level
	parsedLevel, err := log.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("invalid log level: %s", level)
	}
	l.internal.SetLevel(parsedLevel)
	return nil
}

// Level returns the logger level
func (l *CharmLogger) Level() string {
	return l.level
}

// StadardWriter returns a writer that can be used to write logs
func (l *CharmLogger) StandardWriter() io.Writer {
	return l.internal.StandardLog(log.StandardLogOptions{ForceLevel: log.DebugLevel}).Writer()
}

// Print prints a log message
func (l *CharmLogger) Print(message string, keyvals ...any) {
	l.internal.Print(message, keyvals...)
}

// Info logs to info level
func (l *CharmLogger) Info(message string, keyvals ...any) {
	l.internal.Info(message, keyvals...)
}

// Debug logs to debug level
func (l *CharmLogger) Debug(message string, keyvals ...any) {
	l.internal.Debug(message, keyvals...)
}

// Error logs to error level
func (l *CharmLogger) Error(message string, keyvals ...any) {
	l.internal.Error(message, keyvals...)
}

// Fatal logs to fatal level
func (l *CharmLogger) Fatal(message string, keyvals ...any) {
	l.internal.Fatal(message, keyvals...)
}

// Warn logs to warn level
func (l *CharmLogger) Warn(message string, keyvals ...any) {
	l.internal.Warn(message, keyvals...)
}

// Trace logs to trace level
func (l *CharmLogger) Trace(message string, keyvals ...any) {
	l.internal.Debug(message, keyvals...)
}

// Infof logs formatted info level
func (l *CharmLogger) Infof(format string, keyvals ...any) {
	l.internal.Info(fmt.Sprintf(format, keyvals...))
}

// Debugf logs formatted debug level
func (l *CharmLogger) Debugf(format string, keyvals ...any) {
	l.internal.Debug(fmt.Sprintf(format, keyvals...))
}

// SetInteractive sets the logger to use fancy styles if isTerminal is true
func (l *CharmLogger) SetInteractive(interactive string, isTerminal bool) { //nolint:revive
	switch interactive {
	case "auto":
		if isTerminal {
			l.setFancyStyle()
		} else {
			l.setDefaultStyle()
		}
	case "yes":
		l.setFancyStyle()
	default:
		l.setDefaultStyle()
	}
}

func (l *CharmLogger) setDefaultStyle() {
	l.internal.SetStyles(log.DefaultStyles())
	l.internal.SetReportTimestamp(true)
	l.internal.SetReportCaller(l.internal.GetLevel() < log.InfoLevel)
}

func (l *CharmLogger) setFancyStyle() {
	styles := log.DefaultStyles()
	styles.Levels[log.DebugLevel] = debugLevelStyle
	styles.Levels[log.InfoLevel] = infoLevelStyle
	styles.Levels[log.WarnLevel] = warnLevelStyle
	styles.Levels[log.ErrorLevel] = errorLevelStyle
	styles.Levels[log.FatalLevel] = fatalLevelStyle
	styles.Keys["err"] = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))

	l.internal.SetStyles(styles)
	l.internal.SetReportTimestamp(l.internal.GetLevel() < log.InfoLevel)
	l.internal.SetReportCaller(l.internal.GetLevel() < log.InfoLevel)
}

// NullLogger is a logger implementation that captures log output in a string builder.
// It embeds the Logger interface and provides a LogOutput field to access the captured logs.
// Useful for tests
type NullLogger struct {
	Logger
	LogOutput *strings.Builder
}

// Logger that sends all output to a string buffer the captured log output
// can be retrieved by accessing the string buffer at LogOutput
// In the instance of a test failure, the log output is written to StdOut
func NewTestLogger(t *testing.T) Logger {
	sb := &strings.Builder{}
	cl := New(sb, LogLevelDebug)

	t.Cleanup(func() {
		if t.Failed() {
			fmt.Println(sb.String())
		}
	})

	return &NullLogger{cl, sb}
}
