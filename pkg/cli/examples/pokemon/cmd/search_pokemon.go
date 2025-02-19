package cmd

import (
	"context"
	"fmt"
	"strings"

	"pokemon/internal/pokemon"
	svc "pokemon/internal/service/pokemon"
	pkmn "pokemon/internal/usecases/pokemon"

	"github.com/neticdk/go-common/pkg/cli/cmd"
	"github.com/neticdk/go-common/pkg/cli/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func newSearchPokemonCmd(ac *pokemon.Context) *cobra.Command {
	o := &searchPokemonOptions{}
	c := cmd.NewSubCommand("pokemon", o, ac).
		WithShortDesc("Search for pokemons").
		WithLongDesc("Searches for pokemons using pokeapi.co").
		WithExample(searchPokemonCmdExample()).
		Build()

	o.bindFlags(c.Flags())
	return c
}

type searchPokemonOptions struct {
	name string
}

func (o *searchPokemonOptions) bindFlags(f *pflag.FlagSet) {
	f.StringVar(&o.name, "name", "", "Name of the pokemon to search for")
}

func (o *searchPokemonOptions) Complete(_ context.Context, ac *pokemon.Context) error {
	o.name = strings.ToLower(o.name)
	return nil
}

func (o *searchPokemonOptions) Validate(_ context.Context, ac *pokemon.Context) error {
	if o.name == "spiritomb" {
		return ac.EC.ErrorHandler.NewGeneralError(
			"Spiritomb is not a real pokemon",
			"Spiritomb is not a real pokemon. It is a Ghost/Dark-type Pokémon introduced in Generation IV.",
			nil,
			0)
	}
	return nil
}

func (o *searchPokemonOptions) Run(_ context.Context, ac *pokemon.Context) error {
	fmt.Println(printBanner())
	fmt.Println()
	if o.name == "" {
		var names []string
		if err := ui.Spin(ac.EC.Spinner, "Listing pokemons", func(s ui.Spinner) error {
			pokemons, err := pkmn.ListPokemons(ac.PokemonService)
			if err != nil {
				return err
			}

			for _, p := range pokemons {
				names = append(names, p.Name)
			}
			return nil
		}); err != nil {
			return ac.EC.ErrorHandler.NewGeneralError(
				"Failed to list pokemons",
				"Failed to list pokemons. Please try again.",
				err,
				0)
		}

		val, err := ui.Select("Select a pokemon", names)
		if err != nil {
			return ac.EC.ErrorHandler.NewGeneralError(
				"Failed to select pokemon",
				"Failed to select pokemon. Please try again.",
				err,
				0)
		}
		o.name = val
	}

	var (
		p   *svc.Pokemon
		err error
	)
	if err := ui.Spin(ac.EC.Spinner, fmt.Sprintf("Searching for pokemon %q", o.name), func(s ui.Spinner) error {
		p, err = pkmn.SearchPokemon(ac.PokemonService, o.name)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return ac.EC.ErrorHandler.NewGeneralError(
			"Failed to search for pokemon",
			"The pokemon could not be found. Maybe it doesn't exist?",
			err,
			0)
	}
	pkmn.PrintPokemon(p)
	return nil
}

func searchPokemonCmdExample() string {
	b := strings.Builder{}

	b.WriteString("  # Search for a pokemon by name\n")
	b.WriteString("  $ pokemon search pokemon --name pikachu\n")

	b.WriteString("\n")

	b.WriteString("  # Search for a pokemon interactively\n")
	b.WriteString("  $ pokemon search pokemon\n")

	return b.String()
}

func printBanner() string {
	ar := `
⠀⠀⠀⠀⠀⠘⣿⣶⡖⠤⢄⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⠤⠔⣶⣿⠟
⠀⠀⠀⠀⠀⠀⠈⠻⣷⠀⠀⠈⠳⠤⠔⠒⠉⠉⠉⠑⠒⠤⠔⠉⠀⠀⢀⠟⠁⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠈⠑⠤⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡤⠖⠁⠀⠀⠀
⢋⠒⠒⠤⠤⣀⣀⠀⠀⠀⠀⡘⢀⣤⢀⠀⠀⠀⠀⠀⢀⣤⣀⠈⡄⠀⠀⠀⠀⠀
⢸⠀⠀⠀⠀⠀⠀⠉⠉⠒⢲⠇⠈⠛⠛⠀⠀⢀⠀⠀⠈⠛⠋⠀⢃⠀⠀⠀⠀⠀
⠀⣇⡀⠀⠀⠀⠀⠀⠀⠀⢸⢸⣿⡇⠀⠠⢠⠤⡄⠀⠀⢠⣾⣷⢸⠀⠀⠀⠀⠀
⠀⠀⠈⠉⠉⠉⢹⠂⠀⢀⠎⢆⠉⠀⠀⠀⠈⠂⠁⠀⠀⠀⠉⢉⠎⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⢠⠃⠀⢠⠃⠀⡜⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⡄⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢠⠃⠀⠀⠣⢄⡸⢠⠀⠀⠀⢣⠀⠀⠀⢀⠄⠀⠀⢠⢱⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠉⠒⠢⠤⣀⠀⠇⠀⠣⡀⠀⠀⢆⠀⠀⢸⠀⠀⢀⠞⠈⡆⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠰⢅⣸⠀⠀⠀⠙⢤⠴⠚⠂⠀⢮⣤⣤⠎⠀⠀⡇⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠇⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡜⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢈⠗⠂⢤⠤⠤⠤⠤⢤⡤⠤⢤⡊⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⠤⠒⠁⠀⠀⠀⠀⠀⠈⠒⢤⣷⠀⠀⠀⠀⠀⠀
`
	return "\033[33m" + ar + "\033[0m"
}
