package secrets

import (
	"errors"
	"fmt"
	"regexp"
)

var slRe = regexp.MustCompile(`^([a-z]+):\/\/(.+)$`)

// Parse parses a secret locator string into a SecretLocator struct.
func Parse(rawSL string) (*SecretLocator, error) {
	scheme, location, err := parseSecretLocator(rawSL)
	if err != nil {
		return nil, fmt.Errorf("parsing secret identifier %q: %w", rawSL, err)
	}

	return NewSecretLocator(scheme, location)
}

func parseSecretLocator(rawSL string) (Scheme, Location, error) {
	m := slRe.FindStringSubmatch(rawSL)
	if m == nil {
		return "", "", errors.New("invalid identifier")
	}

	scheme := Scheme(m[1])
	location := Location(m[2])

	// Check if the scheme is registered
	if _, exists := providerRegistry[scheme]; !exists {
		return "", "", fmt.Errorf("unknown provider scheme: %s", scheme)
	}

	return scheme, location, nil
}
