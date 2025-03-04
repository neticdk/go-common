package secrets

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type envProvider struct {
	variable string
}

// NewEnvProvider creates a new environment variable provider.
func NewEnvProvider(location Location) *envProvider {
	p := &envProvider{variable: string(location)}
	p.clean()
	return p
}

// RetrieveSecret retrieves the secret from an environment variable.
func (p *envProvider) RetrieveSecret() (*Secret, error) {
	if err := p.validate(); err != nil {
		return nil, fmt.Errorf("validating environment variable %q: %w", p.variable, err)
	}

	v := os.Getenv(p.variable)
	if v == "" {
		return nil, fmt.Errorf("missing environment variable %q", p.variable)
	}
	return NewSecret([]byte(v),
		WithProvider(ProviderEnv),
		WithLocation(Location(p.variable))), nil
}

// String returns the provider ID.
func (p *envProvider) String() string {
	return ProviderEnv.String()
}

func (p *envProvider) clean() {
	p.variable = strings.TrimSpace(p.variable)
}

func (p *envProvider) validate() error {
	matched, _ := regexp.MatchString(`^\w+$`, p.variable)
	if !matched {
		return fmt.Errorf("invalid environment variable name: %s", p.variable)
	}
	return nil
}
