package cmd

import (
	"bytes"
	"log/slog"
	"os"
	"testing"

	"github.com/neticdk/go-common/pkg/cli/errors"
	"github.com/neticdk/go-common/pkg/cli/ui"
	"github.com/neticdk/go-common/pkg/tui/term"
	"github.com/stretchr/testify/assert"
)

func TestNewExecutionContext(t *testing.T) {
	stdin := bytes.NewBufferString("")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	ec := NewExecutionContext("testapp", "A test app", "v1.0.0")
	ec.Stdin = stdin
	ec.Stdout = stdout
	ec.Stderr = stderr

	assert.Equal(t, "testapp", ec.AppName)
	assert.Equal(t, "A test app", ec.ShortDescription)
	assert.Equal(t, "v1.0.0", ec.Version)
	assert.Equal(t, stdin, ec.Stdin)
	assert.Equal(t, stdout, ec.Stdout)
	assert.Equal(t, stderr, ec.Stderr)
	assert.Equal(t, OutputFormatPlain, ec.PFlags.OutputFormat)
	assert.NotNil(t, ec.Logger)
	assert.NotNil(t, ec.ErrorHandler)
	assert.NotNil(t, ec.Spinner)
}

func TestSetLogLevel(t *testing.T) {
	ec := &ExecutionContext{
		PFlags: PFlags{
			LogLevel: LogLevelInfo,
		},
		LogLevel: new(slog.LevelVar),
	}

	ec.SetLogLevel()

	assert.Equal(t, slog.Level(0), ec.LogLevel.Level())
}

func TestSetColor(t *testing.T) {
	testCases := []struct {
		name       string
		isTerminal bool
		noInput    bool
		noColor    bool
		wantColor  bool
	}{
		{
			name:       "terminal, no-input=false, no-color=false",
			isTerminal: true,
			noInput:    false,
			noColor:    false,
			wantColor:  true,
		},
		{
			name:       "terminal, no-input=false, no-color=true",
			isTerminal: true,
			noInput:    false,
			noColor:    true,
			wantColor:  false,
		},
		{
			name:       "terminal, no-input=true, no-color=false",
			isTerminal: true,
			noInput:    true,
			noColor:    false,
			wantColor:  false,
		},
		{
			name:       "terminal, no-input=true, no-color=true",
			isTerminal: true,
			noInput:    true,
			noColor:    true,
			wantColor:  false,
		},
		{
			name:       "no-terminal, no-input=false, no-color=false",
			isTerminal: false,
			noInput:    false,
			noColor:    false,
			wantColor:  false,
		},
		{
			name:       "no-terminal, no-input=false, no-color=true",
			isTerminal: false,
			noInput:    false,
			noColor:    true,
			wantColor:  false,
		},
		{
			name:       "no-terminal, no-input=true, no-color=false",
			isTerminal: false,
			noInput:    true,
			noColor:    false,
			wantColor:  false,
		},
		{
			name:       "no-terminal, no-input=true, no-color=true",
			isTerminal: false,
			noInput:    true,
			noColor:    true,
			wantColor:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ec := &ExecutionContext{
				Stdout:     &bytes.Buffer{},
				Stderr:     &bytes.Buffer{},
				IsTerminal: tc.isTerminal,
				PFlags:     PFlags{NoInput: tc.noInput},
			}

			ui.EnableColor() // reset before each test to ensure no side effects

			ec.SetColor(tc.noColor)

			assert.Equal(t, tc.wantColor, ui.IsColorEnabled())
		})
	}
}

func TestInitInput(t *testing.T) {
	ec := &ExecutionContext{}
	ec.initInput()

	assert.Equal(t, os.Stdin, ec.Stdin)
	assert.Equal(t, term.IsTerminal(os.Stdout), ec.IsTerminal)
	assert.Equal(t, !ec.IsTerminal, ec.PFlags.NoInput)
}

func TestInitOutput(t *testing.T) {
	ec := &ExecutionContext{}
	ec.initOutput()

	assert.Equal(t, os.Stdout, ec.Stdout)
	assert.Equal(t, os.Stderr, ec.Stderr)
}

func TestInitLogger(t *testing.T) {
	stderr := &bytes.Buffer{}
	ec := &ExecutionContext{
		Stderr: stderr,
		PFlags: PFlags{
			LogFormat: LogFormatDefault,
		},
		IsTerminal: false,
		LogLevel:   new(slog.LevelVar),
	}

	ec.initLogger()

	assert.NotNil(t, ec.Logger)
	assert.Equal(t, LogFormatJSON, ec.PFlags.LogFormat)
}

func TestInitErrorHandler(t *testing.T) {
	stderr := &bytes.Buffer{}
	ec := &ExecutionContext{
		Stderr: stderr,
	}

	ec.initErrorHandler()

	assert.NotNil(t, ec.ErrorHandler)
	assert.IsType(t, &errors.DefaultHandler{}, ec.ErrorHandler)
}

func TestInitSpinner(t *testing.T) {
	stdout := &bytes.Buffer{}
	ec := &ExecutionContext{
		Stdout: stdout,
	}

	ec.initSpinner()

	assert.NotNil(t, ec.Spinner)
}
