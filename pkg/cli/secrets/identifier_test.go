package secrets

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifierString(t *testing.T) {
	tests := []struct {
		provider Provider
		location Location
		expected string
	}{
		{
			provider: NewProvider(ProviderEnv, Location("MY_SECRET")),
			location: Location("MY_SECRET"),
			expected: "env://MY_SECRET",
		},
		{
			provider: NewProvider(ProviderFile, Location("/path/to/secret.txt")),
			location: Location("/path/to/secret.txt"),
			expected: "file:///path/to/secret.txt",
		},
		{
			provider: NewProvider(ProviderCmd, Location("get secret")),
			location: Location("get secret"),
			expected: "cmd://get secret",
		},
		{
			provider: NewProvider(ProviderLastPass, Location("12345")),
			location: Location("12345"),
			expected: "lp://12345",
		},
		{
			provider: NewProvider(ProviderUnknown, Location("")),
			location: Location(""),
			expected: "://",
		},
	}

	for _, test := range tests {
		identifier := &Identifier{Provider: test.provider, Location: test.location}
		assert.Equal(t, test.expected, identifier.String())
	}
}

func TestNewIdentifier(t *testing.T) {
	tests := []struct {
		providerID ProviderID
		location   Location
		isValid    bool
	}{
		{
			providerID: ProviderEnv,
			location:   Location("MY_SECRET"),
			isValid:    true,
		},
		{
			providerID: ProviderFile,
			location:   Location("/path/to/secret.txt"),
			isValid:    true,
		},
		{
			providerID: ProviderCmd,
			location:   Location("get secret"),
			isValid:    true,
		},
		{
			providerID: ProviderLastPass,
			location:   Location("12345"),
			isValid:    true,
		},
		{
			providerID: ProviderID(""),
			location:   Location("another-secret"),
			isValid:    false,
		},
	}
	for _, test := range tests {
		identifier := NewIdentifier(test.providerID, test.location)

		if test.isValid {
			assert.NotNil(t, identifier.Provider)
		} else {
			assert.Nil(t, identifier.Provider)
		}
	}
}

func TestIdentifierValidate(t *testing.T) {
	tests := []struct {
		identifier *Identifier
		expected   error
	}{
		{
			identifier: &Identifier{},
			expected:   fmt.Errorf("missing provider"),
		},
		{
			identifier: &Identifier{Provider: NewProvider(ProviderEnv, Location(""))},
			expected:   fmt.Errorf("missing location"),
		},
		{
			identifier: &Identifier{Provider: NewProvider(ProviderEnv, Location("MY_SECRET")), Location: "MY_SECRET"},
			expected:   nil,
		},
	}
	for _, test := range tests {
		err := test.identifier.Validate()
		assert.Equal(t, test.expected, err)
	}
}
