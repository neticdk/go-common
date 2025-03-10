package secrets

import (
	"context"
	"fmt"
)

// Location is a string that represents the location of a secret.
type Location string

// SecretLocator is a reference to a secret.
type SecretLocator struct {
	// Provider implements the provider for the secret.
	Provider Provider

	// Location is the location of the secret within the provider.
	Location Location
}

// String returns the string representation of the secret locator.
func (sl *SecretLocator) String() string {
	var provider string
	if sl.Provider != nil {
		provider = sl.Provider.String()
	}
	return fmt.Sprintf("%s://%s", provider, sl.Location)
}

// GetSecret retrieves the secret from the provider.
// It's a convenience method for retrieving the secret.
func (sl *SecretLocator) GetSecret(ctx context.Context) (*Secret, error) {
	if sl.Provider == nil {
		return nil, fmt.Errorf("missing provider")
	}
	return sl.Provider.RetrieveSecret(ctx, sl.Location)
}

// GetSecretValue retrieves the secret value from the provider.
func (sl *SecretLocator) GetSecretValue(ctx context.Context) (string, error) {
	if sl.Provider == nil {
		return "", fmt.Errorf("missing provider")
	}
	secret, err := sl.Provider.RetrieveSecret(ctx, sl.Location)
	if err != nil {
		return "", fmt.Errorf("retrieving secret: %w", err)
	}
	return secret.String(), nil
}

// NewSecretLocator creates a new secret locator.
func NewSecretLocator(scheme string, location Location) (*SecretLocator, error) {
	provider, err := NewProvider(scheme, location)
	if err != nil {
		return nil, fmt.Errorf("creating provider: %w", err)
	}
	return &SecretLocator{
		Location: location,
		Provider: provider,
	}, nil
}

// Validate validates the secret locator.
func (i *SecretLocator) Validate() error {
	if i.Provider == nil {
		return fmt.Errorf("missing provider")
	}

	if i.Location == "" {
		return fmt.Errorf("missing location")
	}
	return nil
}
