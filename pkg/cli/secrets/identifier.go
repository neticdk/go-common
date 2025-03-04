package secrets

import (
	"fmt"
)

// Location is a string that represents the location of a secret.
type Location string

// Identifier is a reference to a secret.
type Identifier struct {
	// Provider implements the provider for the identifier.
	Provider Provider

	// Location is the location of the secret within the provider.
	Location Location
}

// String returns the string representation of the identifier.
func (i *Identifier) String() string {
	var provider string
	if i.Provider != nil {
		provider = i.Provider.String()
	}
	return fmt.Sprintf("%s://%s", provider, i.Location)
}

// GetSecret retrieves the secret from the provider.
// It's a convenience method for retrieving the secret.
func (i *Identifier) GetSecret() (*Secret, error) {
	if i.Provider == nil {
		return nil, fmt.Errorf("missing provider")
	}
	return i.Provider.RetrieveSecret()
}

// GetSecretValue retrieves the secret value from the provider.
func (i *Identifier) GetSecretValue() (string, error) {
	if i.Provider == nil {
		return "", fmt.Errorf("missing provider")
	}
	secret, err := i.Provider.RetrieveSecret()
	if err != nil {
		return "", fmt.Errorf("retrieving secret: %w", err)
	}
	return secret.String(), nil
}

// NewIdentifier creates a new identifier.
func NewIdentifier(providerID ProviderID, location Location) *Identifier {
	provider := NewProvider(providerID, location)
	i := &Identifier{
		Location: location,
		Provider: provider,
	}
	return i
}

// Validate validates the identifier.
func (i *Identifier) Validate() error {
	if i.Provider == nil {
		return fmt.Errorf("missing provider")
	}

	if i.Location == "" {
		return fmt.Errorf("missing location")
	}
	return nil
}
