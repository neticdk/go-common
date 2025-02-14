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
		outputFlags = append(outputFlags, "plain")
	}
	if ec.PFlags.JSONEnabled {
		f.BoolVar(&json, "json", false, "Output in JSON format")
		outputFlags = append(outputFlags, "json")
	}
	if ec.PFlags.YAMLEnabled {
		f.BoolVar(&yaml, "yaml", false, "Output in YAML format")
		outputFlags = append(outputFlags, "yaml")
	}
	if ec.PFlags.MarkdownEnabled {
		f.BoolVar(&markdown, "markdown", false, "Output in Markdown format")
		outputFlags = append(outputFlags, "markdown")
	}

	switch {
	case plain:
		ec.OutputFormat = context.OutputFormatPlain
	case json:
		ec.OutputFormat = context.OutputFormatJSON
	case yaml:
		ec.OutputFormat = context.OutputFormatYAML
	case markdown:
		ec.OutputFormat = context.OutputFormatMarkdown
	}

	if ec.PFlags.NoHeadersEnabled {
		f.BoolVar(&ec.PFlags.NoHeaders, "no-headers", false, "Do not print headers")
	}

	cmd.MarkFlagsMutuallyExclusive(outputFlags...)

	_ = cmd.PersistentFlags().Parse(os.Args[1:])
	if noColor, err := cmd.PersistentFlags().GetBool("no-color"); err == nil {
		ec.SetColor(noColor)
	}

	cmd.Flags().SortFlags = false
	cmd.PersistentFlags().SortFlags = false

	return f
}
