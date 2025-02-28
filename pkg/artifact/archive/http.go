package archive

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/neticdk/go-common/pkg/artifact"
)

type HTTPOptions struct {
	Downloader artifact.Downloader
	Uncompress func(src, dest string) error
}

// PullHTTP downloads a compressed file from a URL and uncompresses it
func PullHTTP(_ context.Context, a *artifact.Artifact, opts *HTTPOptions) (*artifact.PullResult, error) {
	if opts == nil {
		opts = &HTTPOptions{}
	}

	if opts.Downloader == nil {
		return nil, fmt.Errorf("downloader is required")
	}

	if opts.Uncompress == nil {
		opts.Uncompress = Uncompress
	}

	assetName := path.Base(a.URL)
	if !IsArchive(assetName) {
		return nil, fmt.Errorf("not a compressed file: %s", assetName)
	}

	tmpDir, err := os.MkdirTemp("", "go-common-")
	if err != nil {
		return nil, fmt.Errorf("creating temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	assetDestFile := filepath.Join(tmpDir, assetName)
	if _, err := opts.Downloader.Download(a.URL, assetDestFile); err != nil {
		return nil, err
	}

	tmpDstDir := filepath.Join(tmpDir, a.Name)
	if err := opts.Uncompress(assetDestFile, tmpDstDir); err != nil {
		return nil, fmt.Errorf(`uncompressing asset %q: %w`, assetDestFile, err)
	}

	// Rename the directory to include the version if it's not already included
	if err := os.Rename(
		tmpDstDir,
		a.DestDir()); err != nil {
		return nil, fmt.Errorf("renaming directory: %w", err)
	}

	return &artifact.PullResult{
		Dir:     a.DestDir(),
		Version: a.Version,
	}, nil
}
