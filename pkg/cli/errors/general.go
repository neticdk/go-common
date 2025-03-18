// Deprecated: errors is deprecated and has been moved to github.com/neticdk/go-common/pkg/cli/cmd
//
// It has been deprecated due to the fact that it conflicts with the standard library errors package.
// Please use directly from cmd package instead.
// Expect the errors package to be removed in later versions.
package errors

// Deprecated: use github.com/neticdk/go-common/pkg/cli/cmd instead.
type GeneralError struct {
	Message string
	HelpMsg string
	CodeVal int
	Err     error
}

// Deprecated: use github.com/neticdk/go-common/pkg/cli/cmd instead.
func (e *GeneralError) Error() string { return e.Message }

// Deprecated: use github.com/neticdk/go-common/pkg/cli/cmd instead.
func (e *GeneralError) Help() string { return e.HelpMsg }

// Deprecated: use github.com/neticdk/go-common/pkg/cli/cmd instead.
func (e *GeneralError) Unwrap() error { return e.Err }

// Deprecated: use github.com/neticdk/go-common/pkg/cli/cmd instead.
func (e *GeneralError) Code() int { return e.CodeVal }
