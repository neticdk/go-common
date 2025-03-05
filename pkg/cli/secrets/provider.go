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
	// RetrieveSecret retrieves the secret from the provider.
	// It does the actual work of retrieving the secret..
	RetrieveSecret(context.Context) (*Secret, error)

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

func NewProvider(scheme string, location Location) (Provider, error) {
	factory, exists := providerRegistry[scheme]
	if !exists {
		return nil, fmt.Errorf("unknown provider scheme: %s", scheme)
	}
	return factory(location), nil
}

func init() {
	RegisterProvider("file", func(location Location) Provider {
		return NewFileProvider(location)
	})

	RegisterProvider("env", func(location Location) Provider {
		return NewEnvProvider(location)
	})

	RegisterProvider("cmd", func(location Location) Provider {
		return NewCmdProvider(location)
	})

	RegisterProvider("lp", func(location Location) Provider {
		return NewLastPassProvider(location)
	})
}
