package table

import (
	"ui/internal/ui"

	cliui "github.com/neticdk/go-common/pkg/cli/ui"
)

func Render(ac *ui.Context) error {
	t := cliui.NewTable(ac.EC.Stdout, []string{"Name", "Element"})

	data := [][]string{
		{"Bulbasaur", "Grass/Poison"},
		{"Charmander", "Fire"},
		{"Squirtle", "Water"},
		{"Pikachu", "Electric"},
		{"Jigglypuff", "Normal/Fairy"},
		{"Meowth", "Normal"},
		{"Psyduck", "Water"},
		{"Machop", "Fighting"},
		{"Magnemite", "Electric/Steel"},
		{"Gengar", "Ghost/Poison"},
	}
	return t.WithData(data).Render()
}
