package errors

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
