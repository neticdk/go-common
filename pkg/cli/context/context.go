package context

import (
	"io"
	"log/slog"
	"os"

	"github.com/neticdk/go-common/pkg/cli/errors"
	"github.com/neticdk/go-common/pkg/cli/log"
	"github.com/neticdk/go-common/pkg/cli/ui"
	"github.com/neticdk/go-common/pkg/tui/term"
	"github.com/spf13/cobra"
)

// ExecutionContext holds configuration that can be used (and modified) across
// the application
type ExecutionContext struct {
	// AppName is the executable app name
	// Keep it in lower case letters and use dashes for multi-word app names
	AppName string

	// ShortDescription is a short description of the app
	ShortDescription string

	// LongDescription is a long description of the app
	LongDescription string

	// Stdin is the input stream
	Stdin io.Reader

	// Stdout is the output stream
	Stdout io.Writer

	// Stderr is the error stream
	Stderr io.Writer

	// Command is the current command
	Command *cobra.Command

	// Logger is the global logger
	Logger *slog.Logger

	// ErrorHandler
	ErrorHandler errors.Handler

	// Spinner is the global spinner object used to show progress across the cli
	Spinner ui.Spinner

	// IsTerminal indicates whether the current session is a terminal or not
	IsTerminal bool

	// Version is the CLI version
	// Flag: --version
	Version string

	// LogFormat is the log format used for the logger
	// Flag: --log-format [plain|json]
	LogFormat log.Format

	// LogLevel is the log level used for the logger
	// Flag: --log-level [debug|info|warn|error]
	LogLevel log.Level

	// for changing log level
	logLevel *slog.LevelVar

	// Force is used to force actions
	// Flags: --force, -f
	Force bool

	// DryRun is used to simulate actions
	// Flags: --dry-run
	DryRun bool

	// NoINput can be used to disable interactive mode
	// Flags: --no-input
	NoInput bool

	// NoColor is used to control whether color is used in output
	// Flags: --no-color
	NoColor bool

	// Quiet is used to control whether output is suppressed
	// Flags: --quiet, -q
	Quiet bool

	// Debug is used for debugging
	// Usually this implies verbose output
	// Flags: --debug, -d
	Debug bool

	// OutputFormat is the format used for outputting data
	// Flags: --plain, --json, --yaml, --markdown, etc
	OutputFormat OutputFormat

	// NoHeaders is used to control whether headers are printed
	// Flag: --no-headers
	NoHeaders bool
}

// NewExecutionContext creates a new ExecutionContext
func NewExecutionContext(appName, shortDesc, version string, stdin io.Reader, stdout, stderr io.Writer) *ExecutionContext {
	ec := &ExecutionContext{
		AppName:          appName,
		ShortDescription: shortDesc,
		Version:          version,
		Stdin:            stdin,
		Stdout:           stdout,
		Stderr:           stderr,
		OutputFormat:     OutputFormatPlain,
		LogLevel:         log.DefaultLevel,
		logLevel:         new(slog.LevelVar),
		LogFormat:        log.DefaultFormat,
	}

	ec.initInput()
	ec.initOutput()
	ec.initLogger()
	ec.initErrorHandler()

	return ec
}

// SetLogLevel sets the ec.Logger log level
func (ec *ExecutionContext) SetLogLevel() {
	ui.Logger.Level = ui.ParseLevel(ec.LogLevel)
	ec.logLevel.Set(log.ParseLevel(ec.LogLevel))
}

// SetColor sets weather color should be used in the output
// If the output is not a terminal, color is disabled
// If the --no-color flag is set, color is disabled
// If the --no-input flag is set, color is disabled
func (ec *ExecutionContext) SetColor(noColor bool) {
	if !ec.IsTerminal || ec.NoInput || noColor {
		ui.DisableColor()
	}
}

func (ec *ExecutionContext) initInput() {
	if ec.Stdin == nil {
		ec.Stdin = os.Stdin
	}

	stdout, ok := ec.Stdout.(*os.File)
	if !ok {
		stdout = os.Stdout
	}
	ec.IsTerminal = term.IsTerminal(stdout)
	ec.NoInput = !ec.IsTerminal
}

func (ec *ExecutionContext) initOutput() {
	if ec.Stdout == nil {
		ec.Stdout = os.Stdout
	}
	if ec.Stderr == nil {
		ec.Stderr = os.Stderr
	}

	ui.SetDefaultOutput(ec.Stdout)

	ec.initSpinner()
}

func (ec *ExecutionContext) initLogger() {
	if !ec.IsTerminal {
		ec.LogFormat = log.FormatJSON
	}
	ec.logLevel.Set(log.ParseLevel(log.DefaultLevel))
	handler := ui.NewHandler(ec.Stderr, ec.LogFormat, log.DefaultLevel)
	ec.Logger = slog.New(handler)
}

func (ec *ExecutionContext) initErrorHandler() {
	ec.ErrorHandler = errors.NewDefaultHandler(ec.Stderr)
}

func (ec *ExecutionContext) initSpinner() {
	if ec.Spinner == nil {
		ec.Spinner = ui.DefaultSpinner.WithWriter(ec.Stdout)
	}
}
