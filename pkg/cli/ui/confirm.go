package ui

import "github.com/pterm/pterm"

// Confirm asks the user for confirmation
func Confirm(text string) bool {
	c := pterm.DefaultInteractiveConfirm
	result, _ := c.Show(text)
	pterm.Println()
	return result
}
