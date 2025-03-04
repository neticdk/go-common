package secrets

// Secret stores a secret value along with some data about the secret.
type Secret struct {
	// Value is the secret value.
	Value []byte

	// Provider is the type of secret provider.
	Provider ProviderID

	// Location is the location of the secret within the provider.
	Location Location

	// Data is additional data about the secret.
	Data map[string]string
}

// String returns a string representation of the secret.
func (s *Secret) String() string {
	return string(s.Value)
}

// SecretOption is a function that configures a secret.
type SecretOption func(*Secret)

// NewSecret creates a new secret.
func NewSecret(value []byte, opts ...SecretOption) *Secret {
	s := &Secret{
		Value:    value,
		Provider: ProviderUnknown,
		Location: "",
		Data:     make(map[string]string),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithProvider sets the provider of the secret.
func WithProvider(provider ProviderID) SecretOption {
	return func(s *Secret) {
		s.Provider = provider
	}
}

// WithLocation sets the location of the secret.
func WithLocation(location Location) SecretOption {
	return func(s *Secret) {
		s.Location = location
	}
}
