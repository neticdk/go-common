package selekt

import (
	"ui/internal/ui"

	cliui "github.com/neticdk/go-common/pkg/cli/ui"
)

func Render(ac *ui.Context) error {
	pokemon, err := cliui.Select("Select your favorite Pokemon", []string{"Bulbasaur", "Charmander", "Squirtle", "Pikachu", "Jigglypuff", "Meowth", "Psyduck", "Machop", "Magnemite", "Gengar"})
	if err != nil {
		return err
	}
	cliui.Info.Printf("You selected: %s\n", pokemon)
	return nil
}
