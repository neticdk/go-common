package prefix

import (
	"ui/internal/ui"

	cliui "github.com/neticdk/go-common/pkg/cli/ui"
)

func Render(ac *ui.Context) error {
	cliui.Info.Println("Welcome to the Pokemon selector!")
	cliui.InfoI.Println("We got a lot of Pokemon for you to choose from!")

	cliui.Warning.Println("Be caerful, you can only select one Pokemon!")
	cliui.WarningI.Println("You can't change your selection once you've made it!")

	cliui.Success.Println("Good luck!")
	cliui.SuccessI.Println("You can do it!")

	cliui.Error.Println("Oh no! Something went wrong!")
	cliui.ErrorI.Println("Please try again!")

	cliui.Debug.Println("Debugging information")
	cliui.DebugI.Println("Debugging information with icon!")

	defer func() {
		_ = recover()
		cliui.Description.Println("Recovered from panic")
	}()

	cliui.Fatal.Println("Fatal error!")
	cliui.FatalI.Println("Fatal error with icon!")

	return nil
}
