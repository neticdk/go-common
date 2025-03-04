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
		expect         *Identifier
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
			expect: &Identifier{
				Location: "GITHUB_TOKEN",
			},
			expectProvider: NewEnvProvider(Location("GITHUB_TOKEN")),
			expectError:    false,
		},
		{
			desc: "valid uri with file provider",
			raw:  "file:///path/to/secret.txt",
			expect: &Identifier{
				Location: "/path/to/secret.txt",
			},
			expectProvider: NewFileProvider(Location("/path/to/secret.txt")),
			expectError:    false,
		},
		{
			desc: "valid uri with cmd provider",
			raw:  `cmd://gh auth token`,
			expect: &Identifier{
				Location: "gh auth token",
			},
			expectProvider: NewLastPassProvider(Location("gh auth token")),
			expectError:    false,
		},
		{
			desc: "valid uri with lp provider",
			raw:  `lp://123456`,
			expect: &Identifier{
				Location: "123456",
			},
			expectProvider: NewLastPassProvider(Location("123456")),
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
				providerID, err := ParseProvider(secretScheme.FindStringSubmatch(c.raw)[1])

				assert.NoError(t, err)
				expectedProvider := NewProvider(providerID, actual.Location)
				assert.Equal(t, expectedProvider, actual.Provider)

			}
		})
	}
}
