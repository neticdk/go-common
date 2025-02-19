package cmd

import (
	"context"

	"ui/internal/ui"
	"ui/internal/usecases/confirm"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/spf13/cobra"
)

func newConfirmCmd(ac *ui.Context) *cobra.Command {
	o := &confirmOptions{}
	c := cmd.NewSubCommand("confirm", o, ac).
		WithShortDesc("Confirm demo").
		Build()

	return c
}

type confirmOptions struct{}

func (o *confirmOptions) Complete(_ context.Context, ac *ui.Context) error { return nil }

func (o *confirmOptions) Validate(_ context.Context, _ *ui.Context) error { return nil }

func (o *confirmOptions) Run(_ context.Context, ac *ui.Context) error {
	return confirm.Render(ac)
}
