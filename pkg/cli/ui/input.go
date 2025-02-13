package ui

import "github.com/pterm/pterm"

// Prompt displays a prompt to the user and returns the user's input
func Prompt(header, defaultValue string) (string, error) {
	pterm.Println(header)
	pterm.Println()
	p := pterm.DefaultInteractiveTextInput.WithDefaultValue(defaultValue).WithDelimiter(" ")
	return p.Show(">")
}
