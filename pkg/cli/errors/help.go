// Deprecated: errors is deprecated and has been moved to github.com/go-common/pkg/cli/cmd
//
// It has been deprecated due to the fact that it conflicts with the standard library errors package.
// Please use directly from cmd package instead.
// Expect the errors package to be removed in later versions.
package errors

// ErrorWithHelp interface is used for errors that can provide help
//
// Deprecated: use github.com/neticdk/go-common/pkg/cli/cmd instead.
type ErrorWithHelp interface {
	error

	// Help returns a help message for the error
	Help() string

	// Unwrap returns the underlying error
	Unwrap() error // Optional: for wrapped errors

	// Code returns the error code
	Code() int // Optional: for error codes
}
