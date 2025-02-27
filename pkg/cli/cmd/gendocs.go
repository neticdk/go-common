package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/neticdk/go-common/pkg/cli/ui"
	"github.com/neticdk/go-common/pkg/file"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

type genDocsOptions struct {
	dir string
}

func (o *genDocsOptions) Complete(_ context.Context, _ *ExecutionContext) error { return nil }
func (o *genDocsOptions) Validate(_ context.Context, _ *ExecutionContext) error { return nil }
func (o *genDocsOptions) Run(_ context.Context, ec *ExecutionContext) error {
	if err := os.MkdirAll(o.dir, file.FileModeNewDirectory); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", o.dir, err)
	}
	ui.Info.Printf("Generating documentation in: %s\n", o.dir)
	rootCmd := ec.Command.Root()
	return doc.GenMarkdownTree(rootCmd, o.dir)
}

func GenDocsCommand(ec *ExecutionContext) *cobra.Command {
	o := &genDocsOptions{}
	c := NewSubCommand("gendocs", o, ec).
		WithShortDesc("Generate documentation for the CLI").
		Build()
	c.Hidden = true

	c.Flags().StringVar(&o.dir, "dir", "docs", "The directory to write the documentation to")
	return c
}
