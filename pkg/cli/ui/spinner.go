package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/neticdk/go-common/pkg/tui/term"
	"github.com/pterm/pterm"
)

type Spinner = *pterm.SpinnerPrinter

var DefaultSpinner = pterm.SpinnerPrinter{
	Sequence:            []string{" ⣾ ", " ⣽ ", " ⣻ ", " ⢿ ", " ⡿ ", " ⣟ ", " ⣯ ", " ⣷ "},
	Style:               &pterm.ThemeDefault.SpinnerStyle,
	Delay:               time.Millisecond * 100,
	ShowTimer:           false,
	TimerRoundingFactor: time.Second,
	TimerStyle:          &pterm.ThemeDefault.TimerStyle,
	MessageStyle:        &pterm.ThemeDefault.InfoMessageStyle,
	SuccessPrinter:      &SuccessI,
	FailPrinter:         &ErrorI,
	WarningPrinter:      &WarningI,
	InfoPrinter:         &InfoI,
}

// Hack - clear line when updating text - https://github.com/pterm/pterm/pull/656/files
func UpdateSpinnerText(s Spinner, text string) {
	if s == nil {
		Info.Println(text)
		return
	}
	s.UpdateText("\033[K" + text)
}

// Spin runs a spinner with the given text and function. If the terminal is not
// a TTY, the spinner is not shown.
func Spin(spinner *pterm.SpinnerPrinter, text string, fun func(Spinner) error) error {
	if !term.IsTerminal(os.Stdout) {
		Info.Println(text)
		if err := fun(nil); err != nil {
			Error.Println(fmt.Sprintf("%v", err))
			return err
		}
		Success.Println(text)
		return nil
	}

	s, _ := spinner.Start(text)
	defer func() {
		_ = s.Stop()
	}()
	if s == nil {
		Info.Println(text)
	}
	if err := fun(s); err != nil {
		if s != nil {
			s.Fail()
		}
		return err
	}
	if s != nil {
		// Clear line and move up one line if text is too long, otherwise
		// Success() will print on a new line
		if len(s.Text) > pterm.GetTerminalWidth() {
			pterm.Print("\033[K\033[A\r")
		}
		s.Success()
	}
	return nil
}
