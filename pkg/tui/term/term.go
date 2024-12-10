package term

import (
	"os"

	"github.com/mattn/go-isatty"
)

// IsTerminal is used to determine if the program is running is a terminal
func IsTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}

// IsInteractive is used to determine if the program is running in an
// interactive terminal
func IsInteractive() bool {
	return IsTerminal(os.Stdin) && IsTerminal(os.Stdout)
}
