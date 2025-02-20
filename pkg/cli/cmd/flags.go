package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// PFlags represents is presistent/global flags
type PFlags struct {
	// LogFormat is the log format used for the logger
	// The ForFormat flag is always enabled
	// Flag: --log-format [plain|json]
	LogFormat LogFormat

	// LogLevel is the log level used for the logger
	// The LogLevel flag is always enabled
	// Flag: --log-level [debug|info|warn|error]
	LogLevel LogLevel

	// ForceEnabled is used to enable the Force flag
	ForceEnabled bool

	// Force is used to force actions
	// Flags: --force, -f
	Force bool

	// DryRun is used to simulate actions
	// Flags: --dry-run
	DryRun bool

	// DryRunEnabled is used to enable the DryRun flag
	DryRunEnabled bool

	// NoINput can be used to disable interactive mode
	// Flags: --no-input
	NoInput bool

	// NoInputEnabled is used to enable the NoInput flag
	NoInputEnabled bool

	// NoColor is used to control whether color is used in output
	// The NoColor flag is always enabled
	// Flags: --no-color
	NoColor bool

	// Quiet is used to control whether output is suppressed
	// Flags: --quiet, -q
	Quiet bool

	// QuietEnabled is used to enable the Quiet flag
	QuietEnabled bool

	// Debug is used for debugging
	// Usually this implies verbose output
	// The Debug flag is always enabled
	// Flags: --debug, -d
	Debug bool

	// OutputFormatEnabled is used to enable the OutputFormat flag
	PlainEnabled    bool
	JSONEnabled     bool
	YAMLEnabled     bool
	MarkdownEnabled bool
	TableEnabled    bool

	// NoHeaders is used to control whether headers are printed
	// Flag: --no-headers
	NoHeaders bool

	// NoHeadersEnabled is used to enable the NoHeaders flag
	NoHeadersEnabled bool
}

// AddPersistentFlags adds global flags to the command and does some initialization
func AddPersistentFlags(cmd *cobra.Command, ec *ExecutionContext) *pflag.FlagSet {
	var plain, json, yaml, markdown, table bool

	f := cmd.PersistentFlags()

	logFormats := NewEnum(AllLogFormatsStr(), LogFormatDefault.String())
	f.Var(logFormats, "log-format", fmt.Sprintf("Log format (%s)", AllLogFormatsJoined()))

	logLevels := NewEnum(AllLogLevelsStr(), LogLevelDefault.String())
	f.Var(logLevels, "log-level", fmt.Sprintf("Log level (%s)", AllLogLevelsJoined()))

	if ec.PFlags.ForceEnabled {
		f.BoolVarP(&ec.PFlags.Force, "force", "f", false, "Force actions")
	}
	if ec.PFlags.DryRunEnabled {
		f.BoolVar(&ec.PFlags.DryRun, "dry-run", false, "Simulate action when possible")
	}
	if ec.PFlags.NoInputEnabled {
		f.BoolVar(&ec.PFlags.NoInput, "no-input", false, "Assume non-interactive mode")
	}
	if ec.PFlags.QuietEnabled {
		f.BoolVarP(&ec.PFlags.Quiet, "quiet", "q", false, "Suppress non-essential output")
	}
	f.BoolVar(&ec.PFlags.NoColor, "no-color", false, "Do not print color")
	f.BoolVarP(&ec.PFlags.Debug, "debug", "d", false, "Debug mode")

	outputFlags := []string{}
	if ec.PFlags.PlainEnabled {
		f.BoolVar(&plain, "plain", false, "Output in plain format")
	}
	if ec.PFlags.JSONEnabled {
		f.BoolVar(&json, "json", false, "Output in JSON format")
	}
	if ec.PFlags.YAMLEnabled {
		f.BoolVar(&yaml, "yaml", false, "Output in YAML format")
	}
	if ec.PFlags.MarkdownEnabled {
		f.BoolVar(&markdown, "markdown", false, "Output in Markdown format")
	}
	if ec.PFlags.TableEnabled {
		f.BoolVar(&table, "table", false, "Output in table format")
	}

	_ = cmd.PersistentFlags().Parse(os.Args[1:])

	if logFormat, err := cmd.PersistentFlags().GetString("log-format"); err == nil {
		ec.PFlags.LogFormat = LogFormat(logFormat)
	}
	if logLevel, err := cmd.PersistentFlags().GetString("log-level"); err == nil {
		ec.PFlags.LogLevel = LogLevel(logLevel)
	}

	if plainArg, err := cmd.PersistentFlags().GetBool("plain"); err == nil && plainArg {
		ec.OutputFormat = OutputFormatPlain
		outputFlags = append(outputFlags, "plain")
	}

	if jsonArg, err := cmd.PersistentFlags().GetBool("json"); err == nil && jsonArg {
		ec.OutputFormat = OutputFormatJSON
		outputFlags = append(outputFlags, "json")
	}

	if yamlArg, err := cmd.PersistentFlags().GetBool("yaml"); err == nil && yamlArg {
		ec.OutputFormat = OutputFormatYAML
		outputFlags = append(outputFlags, "yaml")
	}

	if markdownArg, err := cmd.PersistentFlags().GetBool("markdown"); err == nil && markdownArg {
		ec.OutputFormat = OutputFormatMarkdown
		outputFlags = append(outputFlags, "markdown")
	}

	if tableArg, err := cmd.PersistentFlags().GetBool("table"); err == nil && tableArg {
		ec.OutputFormat = OutputFormatTable
		outputFlags = append(outputFlags, "table")
	}

	cmd.MarkFlagsMutuallyExclusive(outputFlags...)

	if noColor, err := cmd.PersistentFlags().GetBool("no-color"); err == nil {
		ec.SetColor(noColor)
	}

	if ec.PFlags.NoHeadersEnabled {
		f.BoolVar(&ec.PFlags.NoHeaders, "no-headers", false, "Do not print headers")
	}

	cmd.Flags().SortFlags = false
	cmd.PersistentFlags().SortFlags = false

	return f
}
