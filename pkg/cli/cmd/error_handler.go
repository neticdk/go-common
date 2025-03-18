package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mitchellh/go-wordwrap"
	"github.com/pterm/pterm"
)

const DefaultWrapWidth = 80

// ErrorHandler is an interface for handling errors.
type ErrorHandler interface {
	// HandleError handles the given error.
	HandleError(err error)
	// NewGeneralError creates a new GeneralError with the specified message, help message, error, and code.
	NewGeneralError(message, helpMsg string, err error, code int) *GeneralError
	// SetWrap sets whether to wrap the error message and help message.
	SetWrap(wrap bool)
	// SetWrapWidth sets the width to wrap the error message and help message.
	SetWrapWidth(width int)
}

// DefaultErrorHandler handles errors providing colored output to a specified writer
// (defaults to stderr).
type DefaultErrorHandler struct {
	Output    io.Writer
	wrap      bool
	wrapWidth int
}

// NewDefaultHandler creates a new DefaultErrorHandler with the specified output writer.
// If no writer is provided, it defaults to stderr.
func NewDefaultHandler(output io.Writer) *DefaultErrorHandler {
	if output == nil {
		output = os.Stderr
	}
	return &DefaultErrorHandler{Output: output, wrap: true, wrapWidth: DefaultWrapWidth}
}

func (h *DefaultErrorHandler) NewGeneralError(message, helpMsg string, err error, code int) *GeneralError {
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
func (h *DefaultErrorHandler) HandleError(err error) {
	if userErr, ok := err.(ErrorWithHelp); ok {
		h.printErrorWithHelp(userErr)
	} else {
		h.printGenericError(err)
	}
}

func (h *DefaultErrorHandler) SetWrap(wrap bool) {
	h.wrap = wrap
}

func (h *DefaultErrorHandler) SetWrapWidth(width int) {
	h.wrapWidth = width
}

// printUserFriendlyError prints the error and help messages in a colored format.
func (h *DefaultErrorHandler) printErrorWithHelp(err ErrorWithHelp) {
	var (
		errorText     string
		errorCode     int
		errorCodeText string
		detailsText   string
		helpText      string
	)

	maxWidth := uint(min(h.wrapWidth, pterm.GetTerminalWidth())) //nolint:gosec  // if you have a terminal width larger than uint, you have other problems
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
		if h.wrap {
			detailsText = wordwrap.WrapString(detailsText, maxWidth)
		}
		paragraph.Print(detailsText)
		pterm.Fprintln(h.Output)
	}

	if helpText != "" {
		pterm.Fprintln(h.Output)
		pterm.Fprintln(h.Output, pterm.Magenta("Explanation")+":")
		if h.wrap {
			helpText = wordwrap.WrapString(helpText, maxWidth)
		}
		paragraph.Println(helpText)
	}
}

// printGenericError prints a generic error message in red.
func (h *DefaultErrorHandler) printGenericError(err error) {
	maxWidth := uint(min(h.wrapWidth, pterm.GetTerminalWidth())) //nolint:gosec  // if you have a terminal width larger than uint, you have other problems
	paragraph := pterm.DefaultBasicText.WithWriter(h.Output)

	pterm.Error.WithWriter(h.Output).Println("An unexpected error occurred")
	pterm.Fprintln(h.Output)
	pterm.Fprintln(h.Output, pterm.Yellow("Details")+": ")
	detailsText := fmt.Sprintf("%v", err)
	if h.wrap {
		detailsText = wordwrap.WrapString(detailsText, maxWidth)
	}
	paragraph.Print(detailsText)
	pterm.Fprintln(h.Output)
	pterm.Fprintln(h.Output)
}
