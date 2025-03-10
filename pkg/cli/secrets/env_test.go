package secrets

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvProvider_GetSecret(t *testing.T) {
	setEnv := func(key, value string) {
		os.Setenv(key, value)
		t.Cleanup(func() {
			os.Unsetenv(key)
		})
	}

	tests := []struct {
		name        string
		location    Location
		expected    *Secret
		expectedErr error
	}{
		{
			name:     "Valid environment variable",
			location: "TEST_SECRET",
			expected: NewSecret([]byte("test_value"),
				WithLocator(&SecretLocator{Scheme: SchemeEnv, Location: "TEST_SECRET"})),
			expectedErr: nil,
		},
		{
			name:        "Missing environment variable",
			location:    "MISSING_SECRET",
			expected:    nil,
			expectedErr: fmt.Errorf("missing environment variable \"MISSING_SECRET\""),
		},
		{
			name:        "Invalid environment variable name",
			location:    "INVALID SECRET",
			expected:    nil,
			expectedErr: fmt.Errorf("validating environment variable \"INVALID SECRET\": invalid environment variable name: INVALID SECRET"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expected != nil {
				setEnv(string(test.location), string(test.expected.Value))
			}

			p := NewEnvProvider()
			secret, err := p.RetrieveSecret(t.Context(), test.location)

			if test.expectedErr != nil {
				assert.Error(t, err)
				if err != nil { // Check to avoid panic if err is nil unexpectedly
					assert.EqualError(t, err, test.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected.GetLocation(), secret.GetLocation())
				assert.Equal(t, test.expected.GetScheme(), secret.GetScheme())
				assert.Equal(t, test.expected.Value, secret.Value)
			}
		})
	}
}
