package types

import (
	"fmt"
	"strings"
)

// Problem is simple implementation of [RFC9457]
//
// [RFC9457]: https://datatracker.ietf.org/doc/html/rfc9457
type Problem struct {
	// Type identify problem type RFC-9457#3.1.1
	//schema:format uri
	Type string `json:"type,omitempty"`

	// Status is the http status code and must be consistent with the server status code RFC-9457#3.1.2
	Status *int `json:"status,omitempty"`

	// Title is short humanreadable summary RFC-9457#3.1.3
	Title string `json:"title,omitempty"`

	// Detail is humanreadable explanation of the specific occurence of the problem RFC-9457#3.1.4
	Detail string `json:"detail,omitempty"`

	// Instance identifies the specific instance of the problem RFC-9457#3.1.5
	Instance string `json:"instance,omitempty"`

	// Err is containing wrapped error and will not be serialized to JSON
	Err error `json:"-"`
}

// Error implements the [error.Error] function to let [Problem] act as an [error]
func (p *Problem) Error() string {
	sb := strings.Builder{}

	if p.Type != "" {
		sb.WriteString(fmt.Sprintf("%s: ", p.Type))
	}

	if p.Detail != "" {
		sb.WriteString(p.Detail)
	} else if p.Title != "" {
		sb.WriteString(p.Title)
	}

	if p.Err != nil {
		sb.WriteString(fmt.Sprintf(": %s", p.Err.Error()))
	}

	return sb.String()
}

// Unwrap allows for nested errors to be unwrapped
func (p *Problem) Unwrap() error {
	return p.Err
}

// IntPointer will return a pointer to the given int
func IntPointer(i int) *int {
	return &i
}
