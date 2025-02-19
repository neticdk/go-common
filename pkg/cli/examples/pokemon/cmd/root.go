package cmd

import (
	"os"

	"pokemon/internal/pokemon"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/spf13/cobra"
)

const (
	AppName   = "pokemon"
	ShortDesc = "The pokemon app!"
	LongDesc  = "Search pokeapi.co for information about pokemon"

	groupSearch = "group-search"
)

// newRootCmd creates the root command
func newRootCmd(ac *pokemon.Context) *cobra.Command {
	c := cmd.NewRootCommand(ac.EC).
		WithInitFunc(func(_ *cobra.Command, _ []string) error {
			ac.SetupPokemonService()
			return nil
		}).
		Build()

	c.AddCommand(
		newSearchCmd(ac),
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
	ac := pokemon.NewContext()
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
