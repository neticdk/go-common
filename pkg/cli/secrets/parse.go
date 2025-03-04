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

	scheme := m[1]
	location := Location(m[2])

	// Check if the scheme is registered
	if _, exists := providerRegistry[scheme]; !exists {
		return "", "", fmt.Errorf("unknown provider scheme: %s", scheme)
	}

	return scheme, location, nil
}
