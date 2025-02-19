package prompt

import (
	"ui/internal/ui"

	cliui "github.com/neticdk/go-common/pkg/cli/ui"
)

func Render(ac *ui.Context) error {
	prompt, err := cliui.Prompt("What is your name?", "")
	if err != nil {
		return err
	}
	if prompt == "" {
		cliui.Info.Println("You didn't enter a name, so I'll just call you Bulbasaur.")
	} else {
		cliui.Info.Printf("Hi %s!\n", prompt)
	}

	return nil
}
