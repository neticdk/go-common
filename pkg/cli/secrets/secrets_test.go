package secrets

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSecret(t *testing.T) {
	tests := []struct {
		name          string
		identifier    string
		expectedError string
		env           map[string]string
	}{
		{
			name:          "invalid identifier",
			identifier:    "invalid-format",
			expectedError: "invalid identifier",
			env:           nil,
		},
		{
			name:          "unknown provider",
			identifier:    "unknown://location",
			expectedError: "unknown provider",
			env:           nil,
		},
		{
			name:          "env provider",
			identifier:    "env://MY_SECRET",
			expectedError: "",
			env:           map[string]string{"MY_SECRET": "secret-value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				_ = os.Setenv(k, v)
			}
			secret, err := GetSecret(tt.identifier)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, secret)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetSecretValue(t *testing.T) {
	tests := []struct {
		name          string
		identifier    string
		expectedError string
		expectedValue string
		env           map[string]string
	}{
		{
			name:          "invalid identifier",
			identifier:    "invalid-format",
			expectedError: "invalid identifier",
			expectedValue: "",
		},
		{
			name:          "unknown provider",
			identifier:    "unknown://location",
			expectedError: "unknown provider",
			expectedValue: "",
		},
		{
			name:          "env provider",
			identifier:    "env://MY_SECRET",
			expectedError: "",
			env:           map[string]string{"MY_SECRET": "secret-value"},
			expectedValue: "secret-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				_ = os.Setenv(k, v)
			}
			value, err := GetSecretValue(tt.identifier)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Equal(t, "", value)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValue, value)
			}
		})
	}
}
