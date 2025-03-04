package secrets

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdProvider_GetSecret(t *testing.T) {
	testCases := []struct {
		location      Location
		expectedValue []byte
		expectedErr   error
		env           map[string]string
	}{
		{
			location:      "echo hello",
			expectedValue: []byte("hello"),
			expectedErr:   nil,
			env:           nil,
		},
		{
			location:      "echo 'hello world'",
			expectedValue: []byte("hello world"),
			expectedErr:   nil,
			env:           nil,
		},
		{
			location:      `foo '`,
			expectedValue: nil,
			expectedErr:   fmt.Errorf("parsing command \"foo '\": invalid command line string"),
			env:           nil,
		},
		{
			location:      "nonexistent_command",
			expectedValue: nil,
			expectedErr:   fmt.Errorf("executing command \"nonexistent_command\": exec: \"nonexistent_command\": executable file not found in $PATH"),
			env:           nil,
		},
		{
			location:      "",
			expectedValue: nil,
			expectedErr:   fmt.Errorf("validating command \"\": command is empty"),
			env:           nil,
		},
		{
			location:      "echo \"hello world\"",
			expectedValue: []byte("hello world"),
			expectedErr:   nil,
			env:           nil,
		},
		{
			location:      "echo $FOO", // Test environment variable expansion
			expectedValue: []byte("bar"),
			expectedErr:   nil,
			env:           map[string]string{"FOO": "bar"},
		},
	}

	for _, tc := range testCases {
		if tc.env != nil {
			for k, v := range tc.env {
				os.Setenv(k, v)
			}
		}
		p := NewCmdProvider(tc.location)
		secret, err := p.GetSecret()

		if tc.expectedErr != nil {
			assert.EqualError(t, err, tc.expectedErr.Error())
		} else {
			assert.NoError(t, err)
			if secret != nil {
				assert.Equal(t, tc.expectedValue, secret.Value)
			}
		}
	}
}
