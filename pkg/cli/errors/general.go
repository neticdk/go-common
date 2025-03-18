package errors

import "fmt"

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
