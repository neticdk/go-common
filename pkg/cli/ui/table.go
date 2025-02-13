package ui

import (
	"io"

	"github.com/pterm/pterm"
)

// NewTable creates a new table with the given headers
func NewTable(w io.Writer, headers []string) *pterm.TablePrinter {
	t := pterm.DefaultTable.WithWriter(w).WithHasHeader(true)
	t.Data = pterm.TableData{headers}
	return t
}
