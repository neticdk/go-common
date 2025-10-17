package cmd

import (
	"context"
	"ui/internal/ui"
	"ui/internal/usecases/prefix"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/spf13/cobra"
)

func newPrefixCmd(ac *ui.Context) *cobra.Command {
	o := &prefixOptions{}
	c := cmd.NewSubCommand("prefix", o, ac).
		WithShortDesc("Prefix demo").
		Build()

	return c
}

type prefixOptions struct{}

func (o *prefixOptions) Complete(_ context.Context, ac *ui.Context) error { return nil }

func (o *prefixOptions) Validate(_ context.Context, _ *ui.Context) error { return nil }

func (o *prefixOptions) Run(_ context.Context, ac *ui.Context) error {
	return prefix.Render(ac)
}
