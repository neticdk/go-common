package secrets

import (
	"fmt"
	"os"
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

func TestIdentifierGetSecret(t *testing.T) {
	// Test for nil provider
	t.Run("nil provider", func(t *testing.T) {
		identifier := &Identifier{
			Provider: nil,
			Location: "test-location",
		}

		secret, err := identifier.GetSecret()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing provider")
		assert.Nil(t, secret)
	})

	// Test with actual file provider
	t.Run("file provider with existing file", func(t *testing.T) {
		// Create a temporary file with test content
		tempFile, err := os.CreateTemp("", "secret-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		// Write test content
		testContent := "test-secret-value"
		if _, err := tempFile.WriteString(testContent); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tempFile.Close()

		// Create identifier with file provider
		identifier := &Identifier{
			Provider: NewProvider(ProviderFile, Location(tempFile.Name())),
			Location: Location(tempFile.Name()),
		}

		// Test secret retrieval
		secret, err := identifier.GetSecret()

		assert.NoError(t, err)
		assert.NotNil(t, secret)
		assert.Equal(t, testContent, secret.String())
	})

	// Test with file provider pointing to non-existent file
	t.Run("file provider with non-existent file", func(t *testing.T) {
		nonExistentPath := "/path/to/nonexistent/file"

		identifier := &Identifier{
			Provider: NewProvider(ProviderFile, Location(nonExistentPath)),
			Location: Location(nonExistentPath),
		}

		secret, err := identifier.GetSecret()

		assert.Error(t, err)
		assert.Nil(t, secret)
	})
}

func TestIdentifierGetSecretValue(t *testing.T) {
	// Test for nil provider
	t.Run("nil provider", func(t *testing.T) {
		identifier := &Identifier{
			Provider: nil,
			Location: "test-location",
		}

		value, err := identifier.GetSecretValue()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing provider")
		assert.Equal(t, "", value)
	})

	// Test with actual file provider
	t.Run("file provider with existing file", func(t *testing.T) {
		// Create a temporary file with test content
		tempFile, err := os.CreateTemp("", "secret-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		// Write test content
		testContent := "test-secret-value"
		if _, err := tempFile.WriteString(testContent); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tempFile.Close()

		// Create identifier with file provider
		identifier := &Identifier{
			Provider: NewProvider(ProviderFile, Location(tempFile.Name())),
			Location: Location(tempFile.Name()),
		}

		// Test secret value retrieval
		value, err := identifier.GetSecretValue()

		assert.NoError(t, err)
		assert.Equal(t, testContent, value)
	})

	// Test with file provider pointing to non-existent file
	t.Run("file provider with non-existent file", func(t *testing.T) {
		nonExistentPath := "/path/to/nonexistent/file"

		identifier := &Identifier{
			Provider: NewProvider(ProviderFile, Location(nonExistentPath)),
			Location: Location(nonExistentPath),
		}

		value, err := identifier.GetSecretValue()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "retrieving secret:")
		assert.Equal(t, "", value)
	})
}
