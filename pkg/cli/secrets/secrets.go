package secrets

import (
	"context"
	"fmt"
	"time"
)

var DefaultTimeout = 10 * time.Second

// GetSecret retrieves a secret using a provider.
// It sets a default timeout for the operation. See DefaultTimeout.
// The identifier is given as an url-like string in the form of:
// PROVIDER://LOCATION
func GetSecret(identifierString string) (*Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	return fetchSecretWithContext(ctx, identifierString)
}

// GetSecretWithContext retrieves a secret using a provider.
// It uses the provided context.
// The identifier is given as an url-like string in the form of:
// PROVIDER://LOCATION
func GetSecretWithContext(ctx context.Context, identifierString string) (*Secret, error) {
	return fetchSecretWithContext(ctx, identifierString)
}

// GetSecretValue retrieves the value of a secret using a provider.
// It sets a default timeout for the operation. See DefaultTimeout.
// The identifier is given as an url-like string in the form of:
// PROVIDER://LOCATION
func GetSecretValue(identifierString string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	return fetchSecretValueWithContext(ctx, identifierString)
}

// GetSecretValueWithContext retrieves the value of a secret using a provider.
// It uses the provided context.
// The identifier is given as an url-like string in the form of:
// PROVIDER://LOCATION
func GetSecretValueWithContext(ctx context.Context, identifierString string) (string, error) {
	return fetchSecretValueWithContext(ctx, identifierString)
}

func fetchSecretWithContext(ctx context.Context, identifierString string) (*Secret, error) {
	identifier, err := Parse(identifierString)
	if err != nil {
		return nil, err
	}
	return identifier.Provider.RetrieveSecret(ctx)
}

func fetchSecretValueWithContext(ctx context.Context, identifierString string) (string, error) {
	identifier, err := Parse(identifierString)
	if err != nil {
		return "", err
	}
	secret, err := identifier.Provider.RetrieveSecret(ctx)
	if err != nil {
		return "", fmt.Errorf("retrieving secret: %w", err)
	}
	return secret.String(), nil
}
