package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Parallel()

	cases := []struct {
		desc           string
		raw            string
		expect         *SecretLocator
		expectProvider Provider
		expectError    bool
	}{
		{
			desc:        "empty uri",
			raw:         "",
			expectError: true,
		},
		{
			desc:        "invalid uri",
			raw:         ":foo",
			expectError: true,
		},
		{
			desc:        "empty provider",
			raw:         "://foo",
			expectError: true,
		},
		{
			desc:        "empty file",
			raw:         "file://",
			expectError: true,
		},
		{
			desc:        "empty env",
			raw:         "env://",
			expectError: true,
		},
		{
			desc:        "empty cmd",
			raw:         "cmd://",
			expectError: true,
		},
		{
			desc:        "empty lp",
			raw:         "lp://",
			expectError: true,
		},
		{
			desc: "valid uri with env provider",
			raw:  "env://GITHUB_TOKEN",
			expect: &SecretLocator{
				Location: "GITHUB_TOKEN",
			},
			expectProvider: NewEnvProvider(),
			expectError:    false,
		},
		{
			desc: "valid uri with file provider",
			raw:  "file:///path/to/secret.txt",
			expect: &SecretLocator{
				Location: "/path/to/secret.txt",
			},
			expectProvider: NewFileProvider(),
			expectError:    false,
		},
		{
			desc: "valid uri with cmd provider",
			raw:  `cmd://gh auth token`,
			expect: &SecretLocator{
				Location: "gh auth token",
			},
			expectProvider: NewLastPassProvider(),
			expectError:    false,
		},
		{
			desc: "valid uri with lp provider",
			raw:  `lp://123456`,
			expect: &SecretLocator{
				Location: "123456",
			},
			expectProvider: NewLastPassProvider(),
			expectError:    false,
		},
		{
			desc:        "valid uri with unknown provider",
			raw:         "foo://bar",
			expectError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			actual, err := Parse(c.raw)

			if c.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, c.expect.Location, actual.Location)

			if c.expect != nil {
				scheme := Scheme(slRe.FindStringSubmatch(c.raw)[1])

				assert.NoError(t, err)

				expectedProvider, err := NewProvider(scheme)
				assert.NoError(t, err)
				assert.Equal(t, expectedProvider, actual.Provider)

			}
		})
	}
}
