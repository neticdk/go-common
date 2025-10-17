package cmd

import (
	"context"
	"ui/internal/ui"
	"ui/internal/usecases/prompt"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/spf13/cobra"
)

func newPromptCmd(ac *ui.Context) *cobra.Command {
	o := &promptOptions{}
	c := cmd.NewSubCommand("prompt", o, ac).
		WithShortDesc("Prompt demo").
		Build()

	return c
}

type promptOptions struct{}

func (o *promptOptions) Complete(_ context.Context, ac *ui.Context) error { return nil }

func (o *promptOptions) Validate(_ context.Context, _ *ui.Context) error { return nil }

func (o *promptOptions) Run(_ context.Context, ac *ui.Context) error {
	return prompt.Render(ac)
}
