package cmd

import (
	"strings"

	"pokemon/internal/pokemon"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/spf13/cobra"
)

func newSearchCmd(ac *pokemon.Context) *cobra.Command {
	o := &cmd.NoopRunner[*pokemon.Context]{}
	c := cmd.NewSubCommand("search", o, ac).
		WithShortDesc("Search for things").
		WithLongDesc("Searches for things using pokeapi.co").
		WithExample(searchCmdExample()).
		WithMinArgs(1).
		Build()
	c.Aliases = []string{"find"}

	c.Args = cobra.NoArgs
	c.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	}

	c.AddCommand(
		newSearchPokemonCmd(ac),
	)

	c.AddGroup(
		&cobra.Group{
			ID:    groupSearch,
			Title: "Search Commands:",
		},
	)

	return c
}

func searchCmdExample() string {
	b := strings.Builder{}

	b.WriteString("  # Search for things\n")
	b.WriteString("  $ pokemon search THING [flags]\n")

	return b.String()
}
