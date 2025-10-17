package cmd

import (
	"os"
	"ui/internal/ui"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/spf13/cobra"
)

const (
	AppName   = "ui"
	ShortDesc = "UI Examples"
	LongDesc  = "This example app demonstrates how to use ui package to create simple UI elements"
)

// newRootCmd creates the root command
func newRootCmd(ac *ui.Context) *cobra.Command {
	c := cmd.NewRootCommand(ac.EC).
		Build()

	c.AddCommand(
		newTableCmd(ac),
		newSelectCmd(ac),
		newSpinnerCmd(ac),
		newConfirmCmd(ac),
		newPromptCmd(ac),
		newPrefixCmd(ac),
	)

	return c
}

// Execute runs the root command and returns the exit code
func Execute(version string) int {
	ec := cmd.NewExecutionContext(
		AppName,
		ShortDesc,
		version,
		os.Stdin,
		os.Stdout,
		os.Stderr)
	ac := ui.NewContext()
	ac.EC = ec
	ec.LongDescription = LongDesc
	rootCmd := newRootCmd(ac)
	err := rootCmd.Execute()
	_ = ec.Spinner.Stop()
	if err != nil {
		ec.ErrorHandler.HandleError(err)
		return 1
	}
	return 0
}
