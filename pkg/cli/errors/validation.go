// Deprecated: errors is deprecated and has been moved to github.com/go-common/pkg/cli/cmd
//
// It has been deprecated due to the fact that it conflicts with the standard library errors package.
// Please use directly from cmd package instead.
// Expect the errors package to be removed in later versions.
package errors

import (
	"fmt"
	"strings"
)

// InvalidArgumentError is returned when an argument is invalid
//
// Deprecated: use github.com/neticdk/go-common/pkg/cli/cmd instead.
type InvalidArgumentError struct {
	Flag     string
	Val      string
	OneOf    []string
	SeeOther string
	Context  string
}

// Deprecated: use github.com/neticdk/go-common/pkg/cli/cmd instead.
func (e *InvalidArgumentError) Error() string { return "Invalid argument" }

// Deprecated: use github.com/neticdk/go-common/pkg/cli/cmd instead.
func (e *InvalidArgumentError) Unwrap() error { return nil }

// Deprecated: use github.com/neticdk/go-common/pkg/cli/cmd instead.
func (e *InvalidArgumentError) Code() int { return 0 }

// Help returns the error message
//
// Deprecated: use github.com/neticdk/go-common/pkg/cli/cmd instead.
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
