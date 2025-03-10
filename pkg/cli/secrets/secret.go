package secrets

// Secret stores a secret value along with some data about the secret.
type Secret struct {
	// Value is the secret value.
	Value []byte

	// locator is a reference to the secret locator this this secret
	locator *SecretLocator
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

// GetScheme returns the scheme of the secret locator.
func (s *Secret) GetScheme() Scheme {
	if s.locator == nil {
		return ""
	}
	return s.locator.Scheme
}

// GetLocation returns the location of the secret locator.
func (s *Secret) GetLocation() Location {
	if s.locator == nil {
		return ""
	}
	return s.locator.Location
}

// SecretOption is a function that configures a secret.
type SecretOption func(*Secret)

// NewSecret creates a new secret.
func NewSecret(value []byte, opts ...SecretOption) *Secret {
	s := &Secret{
		Value: value,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithLocator sets the secret locator of the secret.
func WithLocator(sl *SecretLocator) SecretOption {
	return func(s *Secret) {
		s.locator = sl
	}
}
