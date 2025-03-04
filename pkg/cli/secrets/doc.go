/*
Package secrets provides a flexible way to retrieve secrets from various providers.

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

Alternatively, create the identifier directly:

	identifier := secrets.NewIdentifier(secrets.ProviderFile, "/path/to/secret")
	secret, err := identifier.Provider.RetrieveSecret()
	if err != nil {
		return err
	}

# Integration with Cobra CLI

Example of adding flags to a Cobra command that accept secrets:

	import (
		"github.com/spf13/cobra"
		"github.com/neticdk/go-common/pkg/cli/secrets"
	)

	func NewCommand() *cobra.Command {
		var secretValue string

		cmd := &cobra.Command{
			Use:   "your-command",
			Short: "Command that uses secrets",
			RunE: func(cmd *cobra.Command, args []string) error {
				// Parse the secret identifier
				value, err := secrets.GetSecretValue(secretValue)
				if err != nil {
					return err
				}

				// Use the secret value in your command
				// ...

				return nil
			},
		}

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
