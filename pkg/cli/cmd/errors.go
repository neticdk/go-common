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

	msg.WriteString(fmt.Sprintf("The argument %q is not valid for flag %q.", e.Val, e.Flag))
	if len(e.OneOf) > 0 {
		msg.WriteString("\n\nValid choices:\n\n")
		for _, choice := range e.OneOf {
			msg.WriteString(fmt.Sprintf("  - %s\n", choice))
		}
	}
	if e.Context != "" {
		msg.WriteString(fmt.Sprintf("\n\n%s", e.Context))
	}
	if e.SeeOther != "" {
		msg.WriteString(fmt.Sprintf("\n\nSee also: %q", e.SeeOther))
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
