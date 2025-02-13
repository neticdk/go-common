package errors

// ErrorWithHelp interface is used for errors that can provide help
type ErrorWithHelp interface {
	error
	Help() string
	Unwrap() error // Optional: for wrapped errors
	Code() int     // Optional: for error codes
}
