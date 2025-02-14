package context

import "github.com/neticdk/go-common/pkg/cli/log"

type PFlags struct {
	// LogFormat is the log format used for the logger
	// The ForFormat flag is always enabled
	// Flag: --log-format [plain|json]
	LogFormat log.Format

	// LogLevel is the log level used for the logger
	// The LogLevel flag is always enabled
	// Flag: --log-level [debug|info|warn|error]
	LogLevel log.Level

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

	// NoHeaders is used to control whether headers are printed
	// Flag: --no-headers
	NoHeaders bool

	// NoHeadersEnabled is used to enable the NoHeaders flag
	NoHeadersEnabled bool
}
