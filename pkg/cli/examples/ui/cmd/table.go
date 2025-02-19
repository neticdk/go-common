package cmd

import (
	"context"

	"ui/internal/ui"
	"ui/internal/usecases/table"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/spf13/cobra"
)

func newTableCmd(ac *ui.Context) *cobra.Command {
	o := &tableOptions{}
	c := cmd.NewSubCommand("table", o, ac).
		WithShortDesc("Table demo").
		Build()

	return c
}

type tableOptions struct{}

func (o *tableOptions) Complete(_ context.Context, ac *ui.Context) error { return nil }

func (o *tableOptions) Validate(_ context.Context, _ *ui.Context) error { return nil }

func (o *tableOptions) Run(_ context.Context, ac *ui.Context) error {
	return table.Render(ac)
}
