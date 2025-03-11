package jsonnetbundler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/neticdk/go-common/pkg/artifact"
	"github.com/pkg/errors"

	"github.com/jsonnet-bundler/jsonnet-bundler/pkg"
	"github.com/jsonnet-bundler/jsonnet-bundler/pkg/jsonnetfile"
	specv1 "github.com/jsonnet-bundler/jsonnet-bundler/spec/v1"
	"github.com/jsonnet-bundler/jsonnet-bundler/spec/v1/deps"
)

const filePermissionNewFile = 0o640

type PackageOptions struct {
	JsonnetHome string
}

// Install installs a package and its dependencies in the vendor directory.
// It mimics https://github.com/jsonnet-bundler/jsonnet-bundler/blob/master/cmd/jb/install.go
func Install(_ context.Context, a *artifact.Artifact, opts *PackageOptions) (*artifact.PullResult, error) {
	if opts == nil {
		opts = &PackageOptions{}
	}

	if opts.JsonnetHome == "" {
		opts.JsonnetHome = "vendor"
	}

	jbfilebytes, err := os.ReadFile(filepath.Join(a.BaseDir, jsonnetfile.File))
	if err != nil {
		return nil, fmt.Errorf("reading %q: %w", jsonnetfile.File, err)
	}

	jsonnetFile, err := jsonnetfile.Unmarshal(jbfilebytes)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling %q: %w", jsonnetfile.File, err)
	}

	jblockfilebytes, err := os.ReadFile(filepath.Join(a.BaseDir, jsonnetfile.LockFile))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading %q: %w", jsonnetfile.LockFile, err)
	}

	lockFile, err := jsonnetfile.Unmarshal(jblockfilebytes)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling %q: %w", jsonnetfile.LockFile, err)
	}

	if err := os.MkdirAll(filepath.Join(a.BaseDir, opts.JsonnetHome, ".tmp"), os.ModePerm); err != nil {
		return nil, fmt.Errorf("creating %q: %w", filepath.Join(a.BaseDir, opts.JsonnetHome, ".tmp"), err)
	}

	d := deps.Parse(a.BaseDir, a.URL)
	if d == nil {
		return nil, fmt.Errorf("parsing url %q: %w", a.URL, err)
	}
	a.Version = d.Version

	pkg.GitQuiet = true
	jd, _ := jsonnetFile.Dependencies.Get(d.Name())
	if !depEqual(jd, *d) {
		// the dep passed on the cli is different from the jsonnetFile
		jsonnetFile.Dependencies.Set(d.Name(), *d)

		// we want to install the passed version (ignore the lock)
		lockFile.Dependencies.Delete(d.Name())
	}

	jsonnetPkgHomeDir := filepath.Join(a.BaseDir, opts.JsonnetHome)
	locked, err := pkg.Ensure(jsonnetFile, jsonnetPkgHomeDir, lockFile.Dependencies)
	if err != nil {
		return nil, fmt.Errorf("ensuring package: %w", err)
	}

	if err := writeChangedJsonnetFile(jbfilebytes, &jsonnetFile, filepath.Join(a.BaseDir, jsonnetfile.File)); err != nil {
		return nil, fmt.Errorf("writing %q: %w", jsonnetfile.File, err)
	}

	if err := writeChangedJsonnetFile(jblockfilebytes, &specv1.JsonnetFile{Dependencies: locked}, filepath.Join(a.BaseDir, jsonnetfile.LockFile)); err != nil {
		return nil, fmt.Errorf("writing %q: %w", jsonnetfile.LockFile, err)
	}

	return &artifact.PullResult{
		Dir:     a.DestDir(),
		Version: a.Version,
	}, nil
}

func writeJSONFile(name string, d any) error {
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return errors.Wrap(err, "encoding json")
	}
	b = append(b, []byte("\n")...)

	return os.WriteFile(name, b, filePermissionNewFile)
}

func writeChangedJsonnetFile(originalBytes []byte, modified *specv1.JsonnetFile, path string) error {
	origJsonnetFile, err := jsonnetfile.Unmarshal(originalBytes)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(origJsonnetFile, *modified) {
		return nil
	}

	return writeJSONFile(path, *modified)
}

func depEqual(d1, d2 deps.Dependency) bool {
	name := d1.Name() == d2.Name()
	version := d1.Version == d2.Version
	source := reflect.DeepEqual(d1.Source, d2.Source)

	return name && version && source
}
