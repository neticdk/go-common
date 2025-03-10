package secrets

// Secret stores a secret value along with some data about the secret.
type Secret struct {
	// Value is the secret value.
	Value []byte

	// Provider is the type of secret provider.
	Provider string

	// Location is the location of the secret within the provider.
	Location string

	// Data is additional data about the secret.
	Data map[string]string
}

// String returns a string representation of the secret.
func (s *Secret) String() string {
	return string(s.Value)
}

// DestroySecret destroys the secret value.
func (s *Secret) DestroyValue() {
	if s.Value == nil {
		return
	}
	for i := range s.Value {
		s.Value[i] = 0
	}
	s.Value = nil
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
func WithProvider(scheme string) SecretOption {
	return func(s *Secret) {
		s.Provider = scheme
	}
}

// WithLocation sets the location of the secret.
func WithLocation(location string) SecretOption {
	return func(s *Secret) {
		s.Location = location
	}
}
