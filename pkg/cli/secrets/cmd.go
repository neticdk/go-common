package secrets

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mattn/go-shellwords"
)

const SchemeCmd = "cmd"

type cmdProvider struct {
	command string
}

// NewCmdProvider creates a new command provider.
func NewCmdProvider() *cmdProvider {
	return &cmdProvider{}
}

// RetrieveSecret retrieves a secret from the command output.
func (p *cmdProvider) RetrieveSecret(ctx context.Context, loc Location) (*Secret, error) {
	p.command = string(loc)
	p.clean()

	if err := p.validate(); err != nil {
		return nil, fmt.Errorf("validating command %q: %w", p.command, err)
	}

	cmd, err := splitCommandShellwords(ctx, p.command)
	if err != nil {
		return nil, fmt.Errorf("parsing command %q: %w", p.command, err)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("executing command %q: %w", p.command, err)
	}
	secret := strings.Trim(string(output), "\n")

	sl, err := NewSecretLocator(SchemeCmd, loc)
	if err != nil {
		return nil, fmt.Errorf("creating secret locator: %w", err)
	}

	return NewSecret([]byte(secret), WithLocator(sl)), nil
}

// String returns the scheme for the provider.
func (p *cmdProvider) Scheme() Scheme {
	return SchemeCmd
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
func splitCommandShellwords(ctx context.Context, command string) (*exec.Cmd, error) {
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
	cmd := exec.CommandContext(ctx, fields[0], fields[1:]...)
	return cmd, nil
}
