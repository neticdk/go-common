package spinner

import (
	"errors"
	"time"

	"ui/internal/ui"

	cliui "github.com/neticdk/go-common/pkg/cli/ui"
)

func Render(ac *ui.Context) error {
	if err := cliui.Spin(ac.EC.Spinner, "Waiting for things to settle", func(s cliui.Spinner) error {
		time.Sleep(5 * time.Second)
		cliui.UpdateSpinnerText(s, "Almost there")
		time.Sleep(3 * time.Second)
		cliui.UpdateSpinnerText(s, "Done")
		return nil
	}); err != nil {
		cliui.Error.Println("Something went wrong")
	}

	cliui.Info.Println("Well that went well. Let's try something else.")

	if err := cliui.Spin(ac.EC.Spinner, "Searching for your favorite Pokemon", func(s cliui.Spinner) error {
		time.Sleep(3 * time.Second)
		return errors.New("Pokemon not found")
	}); err != nil {
		cliui.Error.Println("Something went wrong")
	}

	return nil
}
