package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSecret(t *testing.T) {
	tests := []struct {
		name    string
		value   []byte
		opts    []SecretOption
		want    *Secret
		wantErr bool
	}{
		{
			name:  "NoOptions",
			value: []byte("secret"),
			want: &Secret{
				Value:    []byte("secret"),
				Provider: ProviderUnknown,
				Location: "",
				Data:     map[string]string{},
			},
		},
		{
			name:  "WithProvider",
			value: []byte("secret"),
			opts: []SecretOption{
				WithProvider(ProviderEnv),
			},
			want: &Secret{
				Value:    []byte("secret"),
				Provider: ProviderEnv,
				Location: "",
				Data:     map[string]string{},
			},
		},
		{
			name:  "WithLocation",
			value: []byte("secret"),
			opts: []SecretOption{
				WithLocation("path/to/secret"),
			},
			want: &Secret{
				Value:    []byte("secret"),
				Provider: ProviderUnknown,
				Location: "path/to/secret",
				Data:     map[string]string{},
			},
		},
		{
			name:  "WithProviderAndLocation",
			value: []byte("secret"),
			opts: []SecretOption{
				WithProvider(ProviderFile),
				WithLocation("path/to/secret"),
			},
			want: &Secret{
				Value:    []byte("secret"),
				Provider: ProviderFile,
				Location: "path/to/secret",
				Data:     map[string]string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSecret(tt.value, tt.opts...)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWithProvider(t *testing.T) {
	s := &Secret{}
	WithProvider(ProviderEnv)(s)

	assert.Equal(t, ProviderEnv, s.Provider)
}

func TestWithLocation(t *testing.T) {
	s := &Secret{}
	WithLocation("path/to/secret")(s)

	assert.Equal(t, Location("path/to/secret"), s.Location)
}
