/*
Package secrets provides a flexible way to retrieve secrets from various
providers with context-aware functionality for timeout and cancellation support.

The package supports retrieving secrets from:
  - Environment variables (env://)
  - Files (file://)
  - Command output (cmd://)
  - LastPass password manager (lp://)

Each secret is identified by a URL-like string in the format "provider://location" where:

  - provider: is one of "env", "file", "cmd", or "lp"

  - location: is specific to the provider (environment variable name, file path, command, or LastPass ID)

# Basic Usage

Retrieve a secret using the GetSecret function:

	// Get a secret from an environment variable
	secret, err := secrets.GetSecret("env://API_KEY")
	if err != nil {
		return err
	}

Retrieve a secret value using the GetSecretValue function:

	// Get a secret from an environment variable
	apiKey, err := secrets.GetSecretValue("env://API_KEY")
	if err != nil {
		return err
	}

# Context-Aware Functions

For more control over timeouts and cancellation, use the WithContext variants:

	// Use a custom context for secret retrieval
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	secret, err := secrets.GetSecretWithContext(ctx, "env://API_KEY")
	if err != nil {
		return err
	}

	// Or with direct value retrieval
	apiKey, err := secrets.GetSecretValueWithContext(ctx, "env://API_KEY")
	if err != nil {
		return err
	}

Note: The non-context functions use a default timeout of 10 seconds. The default
timeout can be changed by setting the `DefaultTimeout` variable.

# Direct Identifier Usage

Alternatively, create the identifier directly:

	identifier := secrets.NewIdentifier(secrets.ProviderFile, "/path/to/secret")

	// Use with context
	ctx := context.Background()
	secret, err := identifier.Provider.RetrieveSecret(ctx)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retrieve the secret with the custom timeout
	secret, err := secrets.GetSecretWithContext(ctx, "lp://My-Secret-Note/Password")
	if err != nil {
		// Handle timeout or other errors
		if errors.Is(err, context.DeadlineExceeded) {
			log.Fatal("Secret retrieval timed out")
		}
		log.Fatalf("Failed to retrieve secret: %v", err)
	}

## Cancellation

Cancel the operation based on external events:

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Set up cancellation based on some condition
	go func() {
		// Monitor some condition
		<-someConditionChannel
		// Cancel the secret retrieval when condition is met
		cancel()
	}()

	// Attempt to retrieve the secret
	secret, err := secrets.GetSecretWithContext(ctx, "cmd://aws secretsmanager get-secret-value")
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Println("Secret retrieval was cancelled")
		}
		// Handle other errors
	}

## Inheriting Context from Parent Operations

Propagate context from parent operations:

	func processWithSecrets(ctx context.Context) error {
		// Use the parent context for secret retrieval
		apiKey, err := secrets.GetSecretValueWithContext(ctx, "env://API_KEY")
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		// Use the secret in subsequent operations
		return callExternalAPI(ctx, apiKey)
	}

	secret, err := identifier.Provider.RetrieveSecret()
	if err != nil {
		return err
	}

# Integration with Cobra CLI

Example of adding flags to a Cobra command that accept secrets:

	import (

		"context"
		"errors"
		"time"

		"github.com/spf13/cobra"
		"github.com/neticdk/go-common/pkg/cli/secrets"

	)

	func NewCommand() *cobra.Command {
		var secretValue string
		var timeout int

		cmd := &cobra.Command{
			Use:   "your-command",
			Short: "Command that uses secrets",
			RunE: func(cmd *cobra.Command, args []string) error {
				// Use default timeout behavior
				value, err := secrets.GetSecretValue(secretValue)
				if err != nil {
					return err
				}

				// Or use context for custom timeout
				ctx, cancel := context.WithTimeout(cmd.Context(), time.Duration(timeout)*time.Second)
				defer cancel()

				value, err = secrets.GetSecretValueWithContext(ctx, secretValue)
				if err != nil {
					return err
				}

				// Use the secret value in your command
				// ...

				return nil
			},
		}

		// Add timeout flag
		cmd.Flags().IntVar(&timeout, "timeout", 10, "Timeout in seconds for secret retrieval")

		// Add flags that accept secrets
		cmd.Flags().StringVar(&secretValue, "api-key", "",
			`API key (supports secret references like "env://API_KEY", "file:///path/to/secret",
			"cmd://command to execute", or "lp://lastpass-id")`)

		return cmd
	}

# Provider-Specific Details

Environment Variables (env://):
  - Location is the name of the environment variable
  - Example: env://API_KEY

Files (file://):
  - Location is the absolute path to the file
  - Example: file:///etc/secrets/api-key

Commands (cmd://):
  - Location is the command to execute
  - The command's output to stdout is used as the secret
  - Example: cmd://aws secretsmanager get-secret-value --secret-id my-secret --query SecretString --output text

LastPass (lp://):
  - Location is the LastPass item ID
  - Requires the LastPass CLI (lpass) to be installed and configured
  - Example: lp://My-Secret-Note/Password
*/
package secrets
