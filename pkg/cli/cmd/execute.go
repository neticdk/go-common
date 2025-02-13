package cmd

import (
	"os"

	"github.com/neticdk/go-common/pkg/cli/context"
	"github.com/spf13/cobra"
)

// Execute runs the root command and returns the exit code
func Execute(cmd *cobra.Command, appName, shortDesc, longDesc, version string) int {
	ec := context.NewExecutionContext(appName, shortDesc, version, os.Stdin, os.Stdout, os.Stderr)
	ec.ShortDescription = shortDesc
	err := cmd.Execute()
	_ = ec.Spinner.Stop()
	if err != nil {
		ec.ErrorHandler.HandleError(err)
		return 1
	}
	return 0
}
