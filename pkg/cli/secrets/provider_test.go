package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseProvider(t *testing.T) {
	tests := []struct {
		input    string
		expected ProviderID
		wantErr  bool
	}{
		{"file", ProviderFile, false},
		{"env", ProviderEnv, false},
		{"cmd", ProviderCmd, false},
		{"lp", ProviderLastPass, false},
		{"unknown", ProviderUnknown, true},
		{"", ProviderUnknown, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseProvider(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, "", got.String())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestNewProvider(t *testing.T) {
	tests := []struct {
		id       ProviderID
		expected Provider
		location Location
	}{
		{
			id:       ProviderEnv,
			expected: NewEnvProvider(Location("MY_SECRET")),
			location: Location("MY_SECRET"),
		},
		{
			id:       ProviderFile,
			expected: NewFileProvider(Location("/path/to/secret.txt")),
			location: Location("/path/to/secret.txt"),
		},
		{
			id:       ProviderCmd,
			expected: NewCmdProvider(Location("get secret")),
			location: Location("get secret"),
		},
		{
			id:       ProviderLastPass,
			expected: NewLastPassProvider(Location("123456")),
			location: Location("123456"),
		},
		{
			id:       ProviderUnknown,
			expected: nil,
			location: Location(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			got := NewProvider(tt.id, tt.location)
			assert.IsType(t, tt.expected, got)
		})
	}
}
