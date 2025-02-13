package errors

import (
	"fmt"
	"io"
	"os"

	"github.com/mitchellh/go-wordwrap"
	"github.com/pterm/pterm"
)

// Handler is an interface for handling errors.
type Handler interface {
	HandleError(err error)
	NewGeneralError(message, helpMsg string, err error, code int) *GeneralError
}

// Handler handles errors providing colored output to a specified writer
// (defaults to stderr).
type DefaultHandler struct {
	Output io.Writer
}

// NewDefaultHandler creates a new DefaultErrorHandler with the specified output writer.
// If no writer is provided, it defaults to stderr.
func NewDefaultHandler(output io.Writer) *DefaultHandler {
	if output == nil {
		output = os.Stderr
	}
	return &DefaultHandler{Output: output}
}

func (h *DefaultHandler) NewGeneralError(message, helpMsg string, err error, code int) *GeneralError {
	return &GeneralError{
		Message: message,
		HelpMsg: helpMsg,
		CodeVal: code,
		Err:     err,
	}
}

// HandleError checks if the given error implements the ErrorWithHelp interface.
// If it does, it prints the error message and help message in a colored format
// to the configured output writer.  Otherwise, it prints a generic error message.
func (h *DefaultHandler) HandleError(err error) {
	if userErr, ok := err.(ErrorWithHelp); ok {
		h.printErrorWithHelp(userErr)
	} else {
		h.printGenericError(err)
	}
}

// printUserFriendlyError prints the error and help messages in a colored format.
func (h *DefaultHandler) printErrorWithHelp(err ErrorWithHelp) {
	var (
		errorText     string
		errorCode     int
		errorCodeText string
		detailsText   string
		helpText      string
	)

	maxWidth := uint(min(80, pterm.GetTerminalWidth())) //nolint:gosec  // if you have a terminal width larger than uint, you have other problems
	paragraph := pterm.DefaultBasicText.WithWriter(h.Output)

	wrappedErr := err.Unwrap()

	errorText = err.Error()
	errorCode = err.Code()
	helpText = err.Help()

	// If the wrapped error is an ErrorWithHelp, we want to print its help message
	if wrappedErr != nil {
		detailsText = fmt.Sprintf("%v", wrappedErr)
		if wrappedAsEWH, ok := wrappedErr.(ErrorWithHelp); ok {
			helpText = wrappedAsEWH.Help()
			errorText = wrappedAsEWH.Error()
			errorCode = wrappedAsEWH.Code()
			detailsText = fmt.Sprintf("%v: %v", wrappedErr, wrappedAsEWH.Unwrap())
		}
	}

	if errorCode != 0 {
		errorCodeText = fmt.Sprintf(" (%d)", errorCode)
	}

	pterm.Error.WithWriter(h.Output).Println(errorText + errorCodeText)

	if wrappedErr != nil {
		pterm.Fprintln(h.Output)
		pterm.Fprintln(h.Output, pterm.Yellow("Details")+": ")
		paragraph.Print(wordwrap.WrapString(detailsText, maxWidth))
		pterm.Fprintln(h.Output)
	}

	if helpText != "" {
		pterm.Fprintln(h.Output)
		pterm.Fprintln(h.Output, pterm.Magenta("Explanation")+":")
		paragraph.Println(wordwrap.WrapString(helpText, maxWidth))
	}
}

// printGenericError prints a generic error message in red.
func (h *DefaultHandler) printGenericError(err error) {
	maxWidth := uint(min(80, pterm.GetTerminalWidth())) //nolint:gosec  // if you have a terminal width larger than uint, you have other problems
	paragraph := pterm.DefaultBasicText.WithWriter(h.Output)

	pterm.Error.WithWriter(h.Output).Println("An unexpected error occurred")
	pterm.Fprintln(h.Output)
	pterm.Fprintln(h.Output, pterm.Yellow("Details")+": ")
	detailsText := fmt.Sprintf("%v", err)
	paragraph.Print(wordwrap.WrapString(detailsText, maxWidth))
	pterm.Fprintln(h.Output)
	pterm.Fprintln(h.Output)
}
