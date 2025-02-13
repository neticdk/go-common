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

	f.BoolVarP(&ec.Force, "force", "f", false, "Force actions")
	f.BoolVar(&ec.DryRun, "dry-run", false, "Simulate action when possible")
	f.BoolVar(&ec.NoInput, "no-input", false, "Assume non-interactive mode")
	f.BoolVarP(&ec.Quiet, "quiet", "q", false, "Suppress non-essential output")
	f.BoolVarP(&ec.Debug, "debug", "d", false, "Debug mode")

	f.BoolVar(&plain, "plain", false, "Output in plain format")
	f.BoolVar(&json, "json", false, "Output in JSON format")
	f.BoolVar(&yaml, "yaml", false, "Output in YAML format")
	f.BoolVar(&markdown, "markdown", false, "Output in Markdown format")

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

	f.BoolVar(&ec.NoHeaders, "no-headers", false, "Do not print headers")
	f.BoolVar(&ec.NoColor, "no-color", false, "Do not print color")

	cmd.MarkFlagsMutuallyExclusive("plain", "json", "yaml", "markdown")

	_ = cmd.PersistentFlags().Parse(os.Args[1:])
	if noColor, err := cmd.PersistentFlags().GetBool("no-color"); err == nil {
		ec.SetColor(noColor)
	}

	cmd.Flags().SortFlags = false
	cmd.PersistentFlags().SortFlags = false

	return f
}
