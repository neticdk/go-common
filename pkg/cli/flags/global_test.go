package flags

import (
	"os"
	"testing"

	"github.com/neticdk/go-common/pkg/cli/context"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestAddPersistentFlags(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedFormat context.OutputFormat
		expectedFlags  []string
	}{
		{
			name:           "Default format",
			args:           []string{},
			expectedFormat: context.OutputFormatPlain,
			expectedFlags:  []string{},
		},
		{
			name:           "JSON format",
			args:           []string{"--json"},
			expectedFormat: context.OutputFormatJSON,
			expectedFlags:  []string{"json"},
		},
		{
			name:           "YAML format",
			args:           []string{"--yaml"},
			expectedFormat: context.OutputFormatYAML,
			expectedFlags:  []string{"yaml"},
		},
		{
			name:           "Markdown format",
			args:           []string{"--markdown"},
			expectedFormat: context.OutputFormatMarkdown,
			expectedFlags:  []string{"markdown"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			ec := context.NewExecutionContext("test", "test", "0.0.0", os.Stdin, os.Stdout, os.Stderr)
			ec.PFlags.JSONEnabled = true
			ec.PFlags.YAMLEnabled = true
			ec.PFlags.MarkdownEnabled = true

			os.Args = append([]string{"cmd"}, tt.args...)
			AddPersistentFlags(cmd, ec)

			assert.Equal(t, tt.expectedFormat, ec.OutputFormat)
		})
	}
}
