package secrets

import (
	"context"
	"fmt"
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProvider(t *testing.T) {
	tests := []struct {
		scheme        string
		expected      Provider
		location      Location
		expectedError error
	}{
		{
			scheme:        ProviderEnv,
			expected:      NewEnvProvider(),
			location:      Location("MY_SECRET"),
			expectedError: nil,
		},
		{
			scheme:        ProviderFile,
			expected:      NewFileProvider(),
			location:      Location("/path/to/secret.txt"),
			expectedError: nil,
		},
		{
			scheme:        ProviderCmd,
			expected:      NewCmdProvider(),
			location:      Location("get secret"),
			expectedError: nil,
		},
		{
			scheme:        ProviderLastPass,
			expected:      NewLastPassProvider(),
			location:      Location("123456"),
			expectedError: nil,
		},
		{
			scheme:        ProviderUnknown,
			expected:      nil,
			location:      Location(""),
			expectedError: fmt.Errorf("unknown provider scheme: %s", ProviderUnknown),
		},
	}
	for _, tt := range tests {
		t.Run(tt.scheme, func(t *testing.T) {
			got, err := NewProvider(tt.scheme, tt.location)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.IsType(t, tt.expected, got)
		})
	}
}

func TestRegisterProvider(t *testing.T) {
	// Save the original registry to restore it later
	originalRegistry := make(map[string]ProviderFactory)
	maps.Copy(originalRegistry, providerRegistry)

	// Clear the registry for this test
	for k := range providerRegistry {
		delete(providerRegistry, k)
	}

	// Verify registry is empty
	assert.Empty(t, providerRegistry)

	mockFactory := func(loc Location) Provider {
		return &mockProvider{loc: string(loc)}
	}

	// Register the mock provider
	RegisterProvider("mock", mockFactory)

	// Verify the provider was registered
	assert.Contains(t, providerRegistry, "mock")
	assert.NotNil(t, providerRegistry["mock"])

	// Create a provider using the factory
	provider, err := NewProvider("mock", Location("test"))
	assert.NoError(t, err)
	assert.IsType(t, &mockProvider{}, provider)
	assert.Equal(t, "test", provider.(*mockProvider).loc)

	// Restore the original registry
	providerRegistry = originalRegistry
}

// Create a mock provider
type mockProvider struct {
	loc string
}

func (p *mockProvider) RetrieveSecret(ctx context.Context, loc Location) (*Secret, error) {
	p.loc = "test"
	secret := "test"
	return NewSecret([]byte(secret),
		WithProvider(ProviderCmd),
		WithLocation(p.loc)), nil
}

func (p *mockProvider) String() string {
	return "mock"
}
