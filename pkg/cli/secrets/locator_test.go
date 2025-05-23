package secrets

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifierString(t *testing.T) {
	tests := []struct {
		scheme        Scheme
		location      Location
		expected      string
		expectedError error
	}{
		{
			scheme:        SchemeEnv,
			location:      Location("MY_SECRET"),
			expected:      "env://MY_SECRET",
			expectedError: nil,
		},
		{
			scheme:        SchemeFile,
			location:      Location("/path/to/secret.txt"),
			expected:      "file:///path/to/secret.txt",
			expectedError: nil,
		},
		{
			scheme:        SchemeCmd,
			location:      Location("get secret"),
			expected:      "cmd://get secret",
			expectedError: nil,
		},
		{
			scheme:        SchemeLastPass,
			location:      Location("12345"),
			expected:      "lp://12345",
			expectedError: nil,
		},
		{
			scheme:        SchemeUnknown,
			location:      Location(""),
			expected:      "unknown://",
			expectedError: fmt.Errorf("unknown provider scheme: %s", SchemeUnknown),
		},
	}

	for _, tt := range tests {
		provider, err := NewProvider(tt.scheme)
		if tt.expectedError != nil {
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err)
		} else {
			assert.NoError(t, err)
		}
		identifier := &SecretLocator{Provider: provider, Scheme: tt.scheme, Location: tt.location}
		assert.Equal(t, tt.expected, identifier.String())
	}
}

func TestNewIdentifier(t *testing.T) {
	tests := []struct {
		scheme        Scheme
		location      Location
		isValid       bool
		expectedError error
	}{
		{
			scheme:        SchemeEnv,
			location:      Location("MY_SECRET"),
			isValid:       true,
			expectedError: nil,
		},
		{
			scheme:        SchemeFile,
			location:      Location("/path/to/secret.txt"),
			isValid:       true,
			expectedError: nil,
		},
		{
			scheme:        SchemeCmd,
			location:      Location("get secret"),
			isValid:       true,
			expectedError: nil,
		},
		{
			scheme:        SchemeLastPass,
			location:      Location("12345"),
			isValid:       true,
			expectedError: nil,
		},
		{
			scheme:        SchemeUnknown,
			location:      Location("another-secret"),
			isValid:       false,
			expectedError: fmt.Errorf("creating provider: unknown provider scheme: %s", SchemeUnknown),
		},
	}
	for _, tt := range tests {
		identifier, err := NewSecretLocator(tt.scheme, tt.location)

		if tt.expectedError != nil {
			assert.EqualError(t, err, tt.expectedError.Error())
		} else {
			assert.NoError(t, err)
			if tt.isValid {
				assert.NotNil(t, identifier.Provider)
			} else {
				assert.Nil(t, identifier.Provider)
			}
		}
	}
}

func TestIdentifierValidate(t *testing.T) {
	tests := []struct {
		name     string
		scheme   Scheme
		location Location
		expected error
	}{
		{
			name:     "unknown scheme",
			scheme:   SchemeUnknown,
			location: Location("TEST"),
			expected: fmt.Errorf("missing provider"),
		},
		{
			name:     "missing scheme",
			scheme:   Scheme(""),
			location: Location("MY_SECRET"),
			expected: fmt.Errorf("missing scheme"),
		},
		{
			name:     "missing location",
			scheme:   SchemeEnv,
			location: Location(""),
			expected: fmt.Errorf("missing location"),
		},
		{
			name:     "good",
			scheme:   SchemeEnv,
			location: Location("MY_SECRET"),
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, _ := NewProvider(tt.scheme)
			identifier := &SecretLocator{
				Provider: provider,
				Scheme:   tt.scheme,
				Location: tt.location,
			}
			err := identifier.Validate()
			assert.Equal(t, tt.expected, err)
		})
	}
}

func TestIdentifierGetSecret(t *testing.T) {
	// Test for nil provider
	t.Run("nil provider", func(t *testing.T) {
		identifier := &SecretLocator{
			Provider: nil,
			Location: "test-location",
		}

		secret, err := identifier.GetSecret(t.Context())

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

		provider, err := NewProvider(SchemeFile)
		assert.NoError(t, err)

		// Create identifier with file provider
		identifier := &SecretLocator{
			Provider: provider,
			Location: Location(tempFile.Name()),
		}

		// Test secret retrieval
		secret, err := identifier.GetSecret(t.Context())

		assert.NoError(t, err)
		assert.NotNil(t, secret)
		assert.Equal(t, testContent, secret.String())
	})

	// Test with file provider pointing to non-existent file
	t.Run("file provider with non-existent file", func(t *testing.T) {
		nonExistentPath := "/path/to/nonexistent/file"

		provider, err := NewProvider(SchemeFile)
		assert.NoError(t, err)

		identifier := &SecretLocator{
			Provider: provider,
			Location: Location(nonExistentPath),
		}

		secret, err := identifier.GetSecret(t.Context())

		assert.Error(t, err)
		assert.Nil(t, secret)
	})
}

func TestIdentifierGetSecretValue(t *testing.T) {
	// Test for nil provider
	t.Run("nil provider", func(t *testing.T) {
		identifier := &SecretLocator{
			Provider: nil,
			Location: "test-location",
		}

		value, err := identifier.GetSecretValue(t.Context())

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

		provider, err := NewProvider(SchemeFile)
		assert.NoError(t, err)

		// Create identifier with file provider
		identifier := &SecretLocator{
			Provider: provider,
			Location: Location(tempFile.Name()),
		}

		// Test secret value retrieval
		value, err := identifier.GetSecretValue(t.Context())

		assert.NoError(t, err)
		assert.Equal(t, testContent, value)
	})

	// Test with file provider pointing to non-existent file
	t.Run("file provider with non-existent file", func(t *testing.T) {
		nonExistentPath := "/path/to/nonexistent/file"

		provider, err := NewProvider(SchemeFile)
		assert.NoError(t, err)

		identifier := &SecretLocator{
			Provider: provider,
			Location: Location(nonExistentPath),
		}

		value, err := identifier.GetSecretValue(t.Context())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "retrieving secret:")
		assert.Equal(t, "", value)
	})
}
