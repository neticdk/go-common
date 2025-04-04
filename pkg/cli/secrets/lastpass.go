package secrets

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

const SchemeLastPass = "lp"

// execCommand is used for mocking in tests.
var execCommand = exec.CommandContext

type lastPassProvider struct {
	id string
}

// NewLastPassProvider creates a new LastPass provider.
func NewLastPassProvider() *lastPassProvider {
	return &lastPassProvider{}
}

// RetrieveSecret retrieves a secret from LastPass using the password field.
func (p *lastPassProvider) RetrieveSecret(ctx context.Context, loc Location) (*Secret, error) {
	p.id = string(loc)
	p.clean()

	cmd := execCommand(ctx, "lpass", "show", "--password", p.id)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("executing lpass command %q gave output %q: %w", cmd.String(), string(output), err)
	}
	secret := strings.Trim(string(output), "\n")

	sl, err := NewSecretLocator(SchemeLastPass, loc)
	if err != nil {
		return nil, fmt.Errorf("creating secret locator: %w", err)
	}

	return NewSecret([]byte(secret), WithLocator(sl)), nil
}

// String returns the scheme for the provider.
func (p *lastPassProvider) Scheme() Scheme {
	return SchemeLastPass
}

func (p *lastPassProvider) clean() {
	p.id = strings.TrimSpace(p.id)
}
