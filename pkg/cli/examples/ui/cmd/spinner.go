package cmd

import (
	"context"
	"ui/internal/ui"
	"ui/internal/usecases/spinner"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/spf13/cobra"
)

func newSpinnerCmd(ac *ui.Context) *cobra.Command {
	o := &spinnerOptions{}
	c := cmd.NewSubCommand("spinner", o, ac).
		WithShortDesc("Spinner demo").
		Build()

	return c
}

type spinnerOptions struct{}

func (o *spinnerOptions) Complete(_ context.Context, ac *ui.Context) error { return nil }

func (o *spinnerOptions) Validate(_ context.Context, _ *ui.Context) error { return nil }

func (o *spinnerOptions) Run(_ context.Context, ac *ui.Context) error {
	return spinner.Render(ac)
}
