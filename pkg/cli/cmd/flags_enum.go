package cmd

import (
	"fmt"
	"slices"
	"strings"
)

type enum struct {
	Allowed []string
	Value   string
}

// NewEnum gives a list of allowed flag parameters, where the second argument is the default
func NewEnum(allowed []string, d string) *enum {
	return &enum{
		Allowed: allowed,
		Value:   d,
	}
}

func (a enum) String() string {
	return a.Value
}

// Set sets the value of the flag
func (a *enum) Set(p string) error {
	if !slices.Contains(a.Allowed, p) {
		return fmt.Errorf("%s is not included in %s", p, strings.Join(a.Allowed, ","))
	}
	a.Value = p
	return nil
}

func (a *enum) Type() string {
	return "string"
}
