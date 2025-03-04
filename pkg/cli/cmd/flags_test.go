package cmd_test

import (
	"os"
	"testing"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestAddPersistentFlags(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedFormat string
	}{
		{
			name:           "Default format",
			args:           []string{},
			expectedFormat: cmd.OutputFormatPlain,
		},
		{
			name:           "JSON format",
			args:           []string{"--output", "json"},
			expectedFormat: cmd.OutputFormatJSON,
		},
		{
			name:           "YAML format",
			args:           []string{"--output", "yaml"},
			expectedFormat: cmd.OutputFormatYAML,
		},
		{
			name:           "Markdown format",
			args:           []string{"--output", "markdown"},
			expectedFormat: cmd.OutputFormatMarkdown,
		},
		{
			name:           "Table format",
			args:           []string{"--output", "table"},
			expectedFormat: cmd.OutputFormatTable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cobra.Command{}
			ec := cmd.NewExecutionContext("test", "test", "0.0.0")

			os.Args = append([]string{"cmd"}, tt.args...)
			cmd.AddPersistentFlags(c, ec)

			assert.Equal(t, tt.expectedFormat, ec.OutputFormat)
		})
	}
}
