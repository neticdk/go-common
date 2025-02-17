package ui

import "github.com/pterm/pterm"

// DisableColor disables color output
func DisableColor() {
	pterm.DisableColor()
}

// EnableColor enables color output
func EnableColor() {
	pterm.EnableColor()
}

func IsColorEnabled() bool {
	return pterm.PrintColor
}
