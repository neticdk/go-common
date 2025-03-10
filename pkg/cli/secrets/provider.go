package secrets

import (
	"context"
	"fmt"
)

const (
	ProviderUnknown = "unknown"
)

// Provider is an interface that provides a secret.
type Provider interface {
	// RetrieveSecret retrieves a secret from the provider.
	// It does the actual work of retrieving the secret..
	RetrieveSecret(context.Context, Location) (*Secret, error)

	// String returns the string representation of the provider.
	String() string
}

// Factory function type for creating providers
type ProviderFactory func(location Location) Provider

// Registry to store provider factories
var providerRegistry = make(map[string]ProviderFactory)

// Register a new provider type
func RegisterProvider(scheme string, factory ProviderFactory) {
	providerRegistry[scheme] = factory
}

// NewProvider creates a new provider instance
func NewProvider(scheme string, location Location) (Provider, error) {
	factory, exists := providerRegistry[scheme]
	if !exists {
		return nil, fmt.Errorf("unknown provider scheme: %s", scheme)
	}
	return factory(location), nil
}

func init() {
	RegisterProvider("file", func(_ Location) Provider {
		return NewFileProvider()
	})

	RegisterProvider("env", func(_ Location) Provider {
		return NewEnvProvider()
	})

	RegisterProvider("cmd", func(_ Location) Provider {
		return NewCmdProvider()
	})

	RegisterProvider("lp", func(_ Location) Provider {
		return NewLastPassProvider()
	})
}
