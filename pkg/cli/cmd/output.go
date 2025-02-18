package cmd

// OutputFormat represents the output format for the command
type OutputFormat string

const (
	OutputFormatPlain    OutputFormat = "plain"
	OutputFormatJSON     OutputFormat = "json"
	OutputFormatYAML     OutputFormat = "yaml"
	OutputFormatMarkdown OutputFormat = "markdown"
)

func (o OutputFormat) String() string {
	return string(o)
}
