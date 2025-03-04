package secrets

import "fmt"

type ProviderID string

const (
	ProviderFile     ProviderID = "file"
	ProviderEnv      ProviderID = "env"
	ProviderCmd      ProviderID = "cmd"
	ProviderLastPass ProviderID = "lp"
	ProviderUnknown  ProviderID = "unknown"
)

// String returns the string representation of the provider.
func (p ProviderID) String() string {
	return string(p)
}

// ParseProvider parses a provider from a string.
func ParseProvider(v string) (ProviderID, error) {
	switch ProviderID(v) {
	case ProviderFile, ProviderEnv, ProviderCmd, ProviderLastPass:
		return ProviderID(v), nil
	default:
		return "", fmt.Errorf("unknown provider: %s", v)
	}
}

// Provider is an interface that provides a secret.
type Provider interface {
	// RetrieveSecret retrieves the secret from the provider.
	// It does the actual work of retrieving the secret..
	RetrieveSecret() (*Secret, error)

	// String returns the string representation of the provider.
	String() string
}

// NewProvider creates a new provider.
func NewProvider(id ProviderID, location Location) Provider {
	switch id {
	case ProviderFile:
		return NewFileProvider(location)
	case ProviderEnv:
		return NewEnvProvider(location)
	case ProviderCmd:
		return NewCmdProvider(location)
	case ProviderLastPass:
		return NewLastPassProvider(location)
	default:
		return nil
	}
}
