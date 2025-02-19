package confirm

import (
	"ui/internal/ui"

	cliui "github.com/neticdk/go-common/pkg/cli/ui"
)

func Render(ac *ui.Context) error {
	confirm := cliui.Confirm("Do you want to continue?")
	if !confirm {
		cliui.Info.Println("You selected no")
		return nil
	}
	cliui.Info.Println("You selected yes")

	return nil
}
