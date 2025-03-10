package secrets

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLastPassProvider_GetSecret(t *testing.T) {
	// Skip test if lpass is not installed
	if _, err := exec.LookPath("lpass"); err != nil {
		t.Skip("lpass is not installed, skipping LastPass provider tests")
	}

	// Create a temporary directory for gocov files
	tmpDir, err := os.MkdirTemp("", "gocov")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir) // Clean up the temporary directory

	// Mock the exec.Command function for testing
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	tests := []struct {
		name           string
		location       Location
		mockOutput     string
		mockError      error
		expectedSecret *Secret
		expectError    bool
		errorContains  string
	}{
		{
			name:       "successful retrieval",
			location:   "test/entry",
			mockOutput: "secretvalue\n",
			mockError:  nil,
			expectedSecret: &Secret{
				Value:   []byte("secretvalue"),
				locator: &SecretLocator{Scheme: SchemeLastPass, Location: "test/entry"},
			},
			expectError: false,
		},
		{
			name:       "with spaces in location",
			location:   "  entry with spaces  ",
			mockOutput: "another-secret\n",
			mockError:  nil,
			expectedSecret: &Secret{
				Value:   []byte("another-secret"),
				locator: &SecretLocator{Scheme: SchemeLastPass, Location: "  entry with spaces  "},
			},
			expectError: false,
		},
		{
			name:          "command error",
			location:      "nonexistent",
			mockOutput:    "",
			mockError:     fmt.Errorf("exit status 1"),
			expectError:   true,
			errorContains: "executing lpass command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock exec.Command to return our test data
			execCommand = func(ctx context.Context, command string, args ...string) *exec.Cmd {
				cs := []string{"-test.run=TestHelperProcess", "--", command}
				cs = append(cs, args...)
				cmd := exec.Command(os.Args[0], cs...)

				// Set environment variables to pass test case data to the helper process
				cmd.Env = []string{
					fmt.Sprintf("GOCOVERDIR=%s", tmpDir), // Set GOCOVERDIR to the temporary directory
					"GO_WANT_HELPER_PROCESS=1",
					fmt.Sprintf("GO_MOCK_OUTPUT=%s", tt.mockOutput),
					fmt.Sprintf("GO_MOCK_ERROR=%v", tt.mockError != nil),
				}
				return cmd
			}

			provider := NewLastPassProvider()
			secret, err := provider.RetrieveSecret(t.Context(), tt.location)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
				assert.Nil(t, secret)
			} else {
				require.NoError(t, err)
				require.NotNil(t, secret)
				assert.Equal(t, tt.expectedSecret.Value, secret.Value)
				assert.Equal(t, tt.expectedSecret.GetScheme(), secret.GetScheme())
				assert.Equal(t, tt.expectedSecret.GetLocation(), secret.GetLocation())
			}
		})
	}
}

// TestHelperProcess isn't a real test - it's used as a helper process for mocking exec.Command
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	// Get mock data from environment
	mockOutput := os.Getenv("GO_MOCK_OUTPUT")
	shouldError := os.Getenv("GO_MOCK_ERROR") == "true"

	if shouldError {
		os.Exit(1)
	}

	fmt.Fprint(os.Stdout, mockOutput)
	os.Exit(0)
}
