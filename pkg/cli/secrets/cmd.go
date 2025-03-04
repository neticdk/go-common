package secrets

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
)

type cmdProvider struct {
	command string
}

// NewCmdProvider creates a new command provider.
func NewCmdProvider(location Location) *cmdProvider {
	p := &cmdProvider{command: string(location)}
	p.clean()
	return p
}

// GetSecret retrieves the secret from the command output.
func (p *cmdProvider) GetSecret() (*Secret, error) {
	if err := p.validate(); err != nil {
		return nil, fmt.Errorf("validating command %q: %w", p.command, err)
	}

	cmd, err := splitCommandShellwords(p.command)
	if err != nil {
		return nil, fmt.Errorf("parsing command %q: %w", p.command, err)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("executing command %q: %w", p.command, err)
	}
	secret := strings.Trim(string(output), "\n")

	return NewSecret([]byte(secret),
		WithProvider(ProviderCmd),
		WithLocation(Location(p.command))), nil
}

// String returns the provider ID.
func (p *cmdProvider) String() string {
	return ProviderCmd.String()
}

func (p *cmdProvider) clean() {
	p.command = strings.TrimSpace(p.command)
}

func (p *cmdProvider) validate() error {
	if p.command == "" {
		return fmt.Errorf("command is empty")
	}
	return nil
}

// splitCommandShellwords splits a command into fields using shellwords.
// shellwords handles stuff like quotes and escapes.
func splitCommandShellwords(command string) (*exec.Cmd, error) {
	p := shellwords.NewParser()
	p.ParseEnv = true
	fields, err := p.Parse(command)
	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		return nil, nil
	}
	// #nosec G204 - This is a deliberate design choice to allow flexible
	// command execution. The user should be able to run whatever they want in
	// order to retrieve the secret.
	cmd := exec.Command(fields[0], fields[1:]...)
	return cmd, nil
}
