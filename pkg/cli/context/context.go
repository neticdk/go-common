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

	// CommandArgs are the arguments passed to the command
	CommandArgs []string

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

	// PFlags are the persistent flag configuration
	PFlags PFlags

	// OutputFormat is the format used for outputting data
	// Flags: --plain, --json, --yaml, --markdown, etc
	OutputFormat OutputFormat

	// for changing log level
	LogLevel *slog.LevelVar
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
		PFlags: PFlags{
			LogFormat: log.DefaultFormat,
			LogLevel:  log.DefaultLevel,
		},
		LogLevel: new(slog.LevelVar),
	}

	ec.initInput()
	ec.initOutput()
	ec.initLogger()
	ec.initErrorHandler()

	return ec
}

// SetLogLevel sets the ec.Logger log level
func (ec *ExecutionContext) SetLogLevel() {
	ui.Logger.Level = ui.ParseLevel(ec.PFlags.LogLevel)
	ec.LogLevel.Set(log.ParseLevel(ec.PFlags.LogLevel))
}

// SetColor sets weather color should be used in the output
// If the output is not a terminal, color is disabled
// If the --no-color flag is set, color is disabled
// If the --no-input flag is set, color is disabled
func (ec *ExecutionContext) SetColor(noColor bool) {
	if !ec.IsTerminal || ec.PFlags.NoInput || noColor {
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
	ec.PFlags.NoInput = !ec.IsTerminal
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
		ec.PFlags.LogFormat = log.FormatJSON
	}
	ec.LogLevel.Set(log.ParseLevel(log.DefaultLevel))
	handler := ui.NewHandler(ec.Stderr, ec.PFlags.LogFormat, log.DefaultLevel)
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
