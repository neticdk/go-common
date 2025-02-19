package cmd

import (
	"context"

	"ui/internal/ui"
	"ui/internal/usecases/selekt"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/spf13/cobra"
)

func newSelectCmd(ac *ui.Context) *cobra.Command {
	o := &selectOptions{}
	c := cmd.NewSubCommand("select", o, ac).
		WithShortDesc("Select demo").
		Build()

	return c
}

type selectOptions struct{}

func (o *selectOptions) Complete(_ context.Context, ac *ui.Context) error { return nil }

func (o *selectOptions) Validate(_ context.Context, _ *ui.Context) error { return nil }

func (o *selectOptions) Run(_ context.Context, ac *ui.Context) error {
	return selekt.Render(ac)
}
