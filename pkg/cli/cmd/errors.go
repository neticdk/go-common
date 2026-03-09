package cmd

import (
	"fmt"
	"strings"
)

// InvalidArgumentError is returned when an argument is invalid
type InvalidArgumentError struct {
	Flag     string
	Val      string
	OneOf    []string
	SeeOther string
	Context  string
}

func (e *InvalidArgumentError) Error() string { return "Invalid argument" }
func (e *InvalidArgumentError) Unwrap() error { return nil }
func (e *InvalidArgumentError) Code() int     { return 0 }

// Help returns the error message
func (e *InvalidArgumentError) Help() string {
	var msg strings.Builder

	fmt.Fprintf(&msg, "The argument %q is not valid for flag %q.", e.Val, e.Flag)
	if len(e.OneOf) > 0 {
		msg.WriteString("\n\nValid choices:\n\n")
		for _, choice := range e.OneOf {
			fmt.Fprintf(&msg, "  - %s\n", choice)
		}
	}
	if e.Context != "" {
		fmt.Fprintf(&msg, "\n\n%s", e.Context)
	}
	if e.SeeOther != "" {
		fmt.Fprintf(&msg, "\n\nSee also: %q", e.SeeOther)
	}
	return msg.String()
}

type GeneralError struct {
	Message string
	HelpMsg string
	CodeVal int
	Err     error
}

func (e *GeneralError) Error() string { return e.Message }
func (e *GeneralError) Help() string  { return e.HelpMsg }
func (e *GeneralError) Unwrap() error { return e.Err }
func (e *GeneralError) Code() int     { return e.CodeVal }

type GeneralErrorBuilder struct {
	*GeneralError
}

// NewGeneralError creates a new GeneralError with support for formatting the message and help message.
// The error is provided as an argument since it should always be present.
func NewGeneralError(err error) *GeneralErrorBuilder {
	return &GeneralErrorBuilder{&GeneralError{Err: err}}
}

// WithMessage formats according to a format specifier and assigns the result to the message.
func (b *GeneralErrorBuilder) WithMessage(format string, a ...any) *GeneralErrorBuilder {
	b.Message = fmt.Sprintf(format, a...)
	return b
}

// WithHelp formats according to a format specifier and assigns the result to the helpMsg.
func (b *GeneralErrorBuilder) WithHelp(format string, a ...any) *GeneralErrorBuilder {
	b.HelpMsg = fmt.Sprintf(format, a...)
	return b
}

// WithCode assigns the exit code to the codeVal.
func (b *GeneralErrorBuilder) WithCode(code int) *GeneralErrorBuilder {
	b.CodeVal = code
	return b
}

// Build returns the GeneralError.
func (b *GeneralErrorBuilder) Build() *GeneralError {
	return b.GeneralError
}

// ErrorWithHelp interface is used for errors that can provide help
type ErrorWithHelp interface {
	error

	// Help returns a help message for the error
	Help() string

	// Unwrap returns the underlying error
	Unwrap() error // Optional: for wrapped errors

	// Code returns the error code
	Code() int // Optional: for error codes
}
