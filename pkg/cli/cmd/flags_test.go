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
		expectedFormat cmd.OutputFormat
		expectedFlags  []string
	}{
		{
			name:           "Default format",
			args:           []string{},
			expectedFormat: cmd.OutputFormatPlain,
			expectedFlags:  []string{},
		},
		{
			name:           "JSON format",
			args:           []string{"--json"},
			expectedFormat: cmd.OutputFormatJSON,
			expectedFlags:  []string{"json"},
		},
		{
			name:           "YAML format",
			args:           []string{"--yaml"},
			expectedFormat: cmd.OutputFormatYAML,
			expectedFlags:  []string{"yaml"},
		},
		{
			name:           "Markdown format",
			args:           []string{"--markdown"},
			expectedFormat: cmd.OutputFormatMarkdown,
			expectedFlags:  []string{"markdown"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cobra.Command{}
			ec := cmd.NewExecutionContext("test", "test", "0.0.0", os.Stdin, os.Stdout, os.Stderr)
			ec.PFlags.JSONEnabled = true
			ec.PFlags.YAMLEnabled = true
			ec.PFlags.MarkdownEnabled = true

			os.Args = append([]string{"cmd"}, tt.args...)
			cmd.AddPersistentFlags(c, ec)

			assert.Equal(t, tt.expectedFormat, ec.OutputFormat)
		})
	}
}
