package context

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
