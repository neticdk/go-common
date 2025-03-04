package secrets

import (
	"fmt"
	"os/exec"
	"strings"
)

// execCommand is used for mocking in tests.
var execCommand = exec.Command

type lastPassProvider struct {
	id string
}

// NewLastPassProvider creates a new LastPass provider.
func NewLastPassProvider(location Location) *lastPassProvider {
	p := &lastPassProvider{id: string(location)}
	p.clean()
	return p
}

// GetSecret retrieves the secret from LastPass using the password field.
func (p *lastPassProvider) GetSecret() (*Secret, error) {
	cmd := execCommand("lpass", "show", "--password", p.id)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("executing command %q: %w", cmd.String(), err)
	}
	secret := strings.Trim(string(output), "\n")

	return NewSecret([]byte(secret),
		WithProvider(ProviderLastPass),
		WithLocation(Location(p.id))), nil
}

// String returns the provider ID.
func (p *lastPassProvider) String() string {
	return ProviderLastPass.String()
}

func (p *lastPassProvider) clean() {
	p.id = strings.TrimSpace(p.id)
}
