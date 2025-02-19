package cmd

import (
	"os"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/spf13/cobra"
)

const (
	AppName   = "hello-world"
	ShortDesc = "A greeting app"
	LongDesc  = `This application greets the user with a friendly messages`
)

type AppContext struct {
	EC *cmd.ExecutionContext
}

func newAppContext() *AppContext {
	return &AppContext{}
}

// newRootCmd creates the root command
func newRootCmd(ac *AppContext) *cobra.Command {
	c := cmd.NewRootCommand(ac.EC).
		Build()

	c.AddCommand(
		newHelloCmd(ac),
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
	ac := newAppContext()
	ac.EC = ec
	ec.LongDescription = LongDesc
	rootCmd := NewRootCmd(ac)
	err := rootCmd.Execute()
	_ = ec.Spinner.Stop()
	if err != nil {
		ec.ErrorHandler.HandleError(err)
		return 1
	}
	return 0
}
