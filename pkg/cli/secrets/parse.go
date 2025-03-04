package secrets

import (
	"errors"
	"fmt"
	"regexp"
)

var secretScheme = regexp.MustCompile(`^([a-z]+):\/\/(.+)$`)

// Parse parses a secret identifier string into an Identifier struct.
func Parse(identifierString string) (*Identifier, error) {
	p, l, err := parseSecretIdentifier(identifierString)
	if err != nil {
		return nil, fmt.Errorf("parsing secret identifier %q: %w", identifierString, err)
	}

	return NewIdentifier(p, l)
}

func parseSecretIdentifier(identifierString string) (string, Location, error) {
	m := secretScheme.FindStringSubmatch(identifierString)
	if m == nil {
		return "", "", errors.New("invalid identifier")
	}

	return m[1], Location(m[2]), nil
}
