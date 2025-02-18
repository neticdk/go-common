package cmd

import (
	"context"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/neticdk/go-common/pkg/cli/ui"
	"github.com/spf13/cobra"
)

func HelloCmd(ac *AppContext) *cobra.Command {
	o := &helloOptions{}
	c := cmd.NewSubCommand("hello", o, ac).
		WithShortDesc("Say hello!").
		Build()

	return c
}

type helloOptions struct {
	who string
}

func (o *helloOptions) Complete(_ context.Context, ac *AppContext) {
	if len(ac.EC.CommandArgs) > 0 {
		o.who = ac.EC.CommandArgs[0]
	} else {
		o.who = "World"
	}
}

func (o *helloOptions) Validate(_ context.Context, _ *AppContext) error { return nil }

func (o *helloOptions) Run(_ context.Context, ac *AppContext) error {
	ui.Info.Printf("Hello, %s!\n", o.who)
	return nil
}
