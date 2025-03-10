package secrets

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	SchemeFile        = "file"
	worldReadablePerm = 0o004
)

type fileProvider struct {
	path string
}

// NewFileProvider creates a new file provider.
func NewFileProvider() *fileProvider {
	return &fileProvider{}
}

// RetrieveSecret retrieves a secret from a file.
func (p *fileProvider) RetrieveSecret(_ context.Context, loc Location) (*Secret, error) {
	p.path = string(loc)
	p.clean()

	if err := p.validate(); err != nil {
		return nil, fmt.Errorf("validating file %q: %w", p.path, err)
	}

	if err := p.checkFile(); err != nil {
		return nil, err
	}

	content, err := os.ReadFile(p.path)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", p.path, err)
	}

	sl, err := NewSecretLocator(SchemeFile, loc)
	if err != nil {
		return nil, fmt.Errorf("creating secret locator: %w", err)
	}

	return NewSecret(content, WithLocator(sl)), nil
}

// Scheme returns the scheme for the provider.
func (p *fileProvider) Scheme() Scheme {
	return SchemeFile
}

func (p *fileProvider) clean() {
	p.path = filepath.Clean(p.path)
}

func (p *fileProvider) validate() error {
	if filepath.IsAbs(p.path) {
		return nil
	}
	return fmt.Errorf("invalid file path: %s", p.path)
}

func (p *fileProvider) checkFile() error {
	info, err := os.Stat(p.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file %q does not exist", p.path)
		}
		return fmt.Errorf("checking if file %q exists: %w", p.path, err)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("file %q is not a regular file", p.path)
	}

	// Check if file is world-readable
	mode := info.Mode().Perm()
	if mode&worldReadablePerm != 0 {
		return fmt.Errorf("file %q has insecure permissions (world-readable): %v", p.path, mode)
	}
	return nil
}
