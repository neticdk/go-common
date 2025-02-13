package ui

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
)

// Select displays a select prompt with the given header and options
func Select(header string, options []string) (string, error) {
	selectedOption, err := pterm.DefaultInteractiveSelect.WithOptions(options).Show(header)
	if err != nil {
		return "", fmt.Errorf("selecting option: %w", err)
	}
	if selectedOption != "" {
		return strings.Split(selectedOption, ":")[0], nil
	}
	return "", nil
}
