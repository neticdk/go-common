package secrets

import (
	"context"
	"fmt"
)

const (
	SchemeUnknown = "unknown"
)

// Provider is an interface that provides a secret.
type Provider interface {
	// RetrieveSecret retrieves a secret from the provider.
	// It does the actual work of retrieving the secret..
	RetrieveSecret(context.Context, Location) (*Secret, error)

	// Scheme returns the scheme for the provider.
	Scheme() Scheme
}

// Factory function type for creating providers
type ProviderFactory func() Provider

// Registry to store provider factories
var providerRegistry = make(map[Scheme]ProviderFactory)

// Register a new provider type
func RegisterProvider(scheme Scheme, factory ProviderFactory) {
	providerRegistry[scheme] = factory
}

// NewProvider creates a new provider instance
func NewProvider(scheme Scheme) (Provider, error) {
	factory, exists := providerRegistry[scheme]
	if !exists {
		return nil, fmt.Errorf("unknown provider scheme: %s", scheme)
	}
	return factory(), nil
}

func init() {
	RegisterProvider(SchemeFile, func() Provider {
		return NewFileProvider()
	})

	RegisterProvider(SchemeEnv, func() Provider {
		return NewEnvProvider()
	})

	RegisterProvider(SchemeCmd, func() Provider {
		return NewCmdProvider()
	})

	RegisterProvider(SchemeLastPass, func() Provider {
		return NewLastPassProvider()
	})
}
