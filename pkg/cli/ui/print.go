package ui

import (
	"io"
)

// SetOutput sets the default output writer for the UI
func SetDefaultOutput(w io.Writer) {
	updateWriters(w)
}
