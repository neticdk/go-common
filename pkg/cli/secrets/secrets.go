package secrets

import (
	"context"
	"fmt"
	"time"
)

var DefaultTimeout = 10 * time.Second

// GetSecret retrieves a secret.
// The secret locator is in the form: provider://location
// It sets a default timeout for the operation. See DefaultTimeout.
func GetSecret(rawSL string) (*Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	return fetchSecretWithContext(ctx, rawSL)
}

// GetSecretWithContext retrieves a secret.
// The secret locator is in the form: provider://location
func GetSecretWithContext(ctx context.Context, rawSL string) (*Secret, error) {
	return fetchSecretWithContext(ctx, rawSL)
}

// GetSecretValue retrieves the value of a secret.
// The secret locator is in the form: provider://location
// It sets a default timeout for the operation. See DefaultTimeout.
func GetSecretValue(rawSL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	return fetchSecretValueWithContext(ctx, rawSL)
}

// GetSecretValueWithContext retrieves the value of a secret.
// The secret locator is in the form: provider://location
func GetSecretValueWithContext(ctx context.Context, rawSL string) (string, error) {
	return fetchSecretValueWithContext(ctx, rawSL)
}

func fetchSecretWithContext(ctx context.Context, rawSL string) (*Secret, error) {
	sl, err := Parse(rawSL)
	if err != nil {
		return nil, err
	}
	return sl.Provider.RetrieveSecret(ctx, sl.Location)
}

func fetchSecretValueWithContext(ctx context.Context, rawSL string) (string, error) {
	sl, err := Parse(rawSL)
	if err != nil {
		return "", err
	}
	secret, err := sl.Provider.RetrieveSecret(ctx, sl.Location)
	if err != nil {
		return "", fmt.Errorf("retrieving secret: %w", err)
	}
	return secret.String(), nil
}
