package secrets

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	ProviderFile      = "file"
	worldReadablePerm = 0o004
)

type fileProvider struct {
	path string
}

// NewFileProvider creates a new file provider.
func NewFileProvider(location Location) *fileProvider {
	p := &fileProvider{path: string(location)}
	p.clean()
	return p
}

// RetrieveSecret retrieves the secret from a file.
func (p *fileProvider) RetrieveSecret(_ context.Context) (*Secret, error) {
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

	return NewSecret(content,
		WithProvider(ProviderFile),
		WithLocation(Location(p.path))), nil
}

// String returns the provider ID.
func (p *fileProvider) String() string {
	return ProviderFile
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
