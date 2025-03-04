package secrets

import (
	"errors"
	"fmt"
	"regexp"
)

var secretScheme = regexp.MustCompile(`^([a-z]+):\/\/(.+)$`)

// Parse parses a secret identifier string into an Identifier struct.
func Parse(identifier string) (*Identifier, error) {
	provider, location, err := parseSecretIdentifier(identifier)
	if err != nil {
		return nil, fmt.Errorf("parsing secret identifier %q: %w", identifier, err)
	}

	return NewIdentifier(provider, location)
}

func parseSecretIdentifier(identifier string) (string, Location, error) {
	m := secretScheme.FindStringSubmatch(identifier)
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
