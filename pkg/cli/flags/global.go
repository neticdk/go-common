package flags

import (
	"fmt"
	"os"

	"github.com/neticdk/go-common/pkg/cli/context"
	"github.com/neticdk/go-common/pkg/cli/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// AddPersistentFlags adds global flags to the command and does some initialization
func AddPersistentFlags(cmd *cobra.Command, ec *context.ExecutionContext) *pflag.FlagSet {
	var plain, json, yaml, markdown bool

	f := cmd.PersistentFlags()

	logFormats := newEnum(log.AllFormatsStr(), log.DefaultFormat.String())
	f.Var(logFormats, "log-format", fmt.Sprintf("Log format (%s)", log.AllFormatsJoined()))

	logLevels := newEnum(log.AllLevelsStr(), log.DefaultLevel.String())
	f.Var(logLevels, "log-level", fmt.Sprintf("Log level (%s)", log.AllLevelsJoined()))

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

	_ = cmd.PersistentFlags().Parse(os.Args[1:])

	if logFormat, err := cmd.PersistentFlags().GetString("log-format"); err == nil {
		ec.PFlags.LogFormat = log.Format(logFormat)
	}
	if logLevel, err := cmd.PersistentFlags().GetString("log-level"); err == nil {
		ec.PFlags.LogLevel = log.Level(logLevel)
	}

	if plainArg, err := cmd.PersistentFlags().GetBool("plain"); err == nil && plainArg {
		ec.OutputFormat = context.OutputFormatPlain
		outputFlags = append(outputFlags, "plain")
	}

	if jsonArg, err := cmd.PersistentFlags().GetBool("json"); err == nil && jsonArg {
		ec.OutputFormat = context.OutputFormatJSON
		outputFlags = append(outputFlags, "json")
	}

	if yamlArg, err := cmd.PersistentFlags().GetBool("yaml"); err == nil && yamlArg {
		ec.OutputFormat = context.OutputFormatYAML
		outputFlags = append(outputFlags, "yaml")
	}

	if markdownArg, err := cmd.PersistentFlags().GetBool("markdown"); err == nil && markdownArg {
		ec.OutputFormat = context.OutputFormatMarkdown
		outputFlags = append(outputFlags, "markdown")
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
