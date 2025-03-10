package secrets

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const SchemeEnv = "env"

type envProvider struct {
	variable string
}

// NewEnvProvider creates a new environment variable provider.
func NewEnvProvider() *envProvider {
	return &envProvider{}
}

// RetrieveSecret retrieves a secret from an environment variable.
func (p *envProvider) RetrieveSecret(_ context.Context, loc Location) (*Secret, error) {
	p.variable = string(loc)
	p.clean()

	if err := p.validate(); err != nil {
		return nil, fmt.Errorf("validating environment variable %q: %w", p.variable, err)
	}

	v := os.Getenv(p.variable)
	if v == "" {
		return nil, fmt.Errorf("missing environment variable %q", p.variable)
	}

	sl, err := NewSecretLocator(SchemeEnv, loc)
	if err != nil {
		return nil, fmt.Errorf("creating secret locator: %w", err)
	}

	return NewSecret([]byte(v), WithLocator(sl)), nil
}

// String returns the scheme for the provider.
func (p *envProvider) Scheme() Scheme {
	return SchemeEnv
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
