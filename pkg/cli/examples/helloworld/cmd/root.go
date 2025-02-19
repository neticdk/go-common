package cmd

import (
	"hello-world/internal/helloworld"
	"os"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/spf13/cobra"
)

const (
	AppName   = "hello-world"
	ShortDesc = "A greeting app"
	LongDesc  = `This application greets the user with a friendly messages`
)

// newRootCmd creates the root command
func newRootCmd(ac *helloworld.Context) *cobra.Command {
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
	ac := helloworld.NewContext()
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
