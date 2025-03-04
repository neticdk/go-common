package secrets

import (
	"errors"
	"fmt"
	"regexp"
)

var secretScheme = regexp.MustCompile(`^([a-z]+):\/\/(.+)$`)

// Parse parses a secret identifier
// The identifier is given as an url-like string in the form of:
// PROVIDER://LOCATION
func Parse(id string) (*Identifier, error) {
	p, l, err := parseSecretIdentifier(id)
	if err != nil {
		return nil, fmt.Errorf("parsing secret identifier %q: %w", id, err)
	}

	return NewIdentifier(p, l), nil
}

func parseSecretIdentifier(id string) (ProviderID, Location, error) {
	m := secretScheme.FindStringSubmatch(id)
	if m == nil {
		return "", "", errors.New("invalid identifier")
	}

	provider, err := ParseProvider(m[1])
	if err != nil {
		return "", "", fmt.Errorf("parsing provider %q: %w", m[1], err)
	}

	return provider, Location(m[2]), nil
}
