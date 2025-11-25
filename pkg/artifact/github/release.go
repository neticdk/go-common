package github

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	gh "github.com/google/go-github/v79/github"
	"github.com/neticdk/go-common/pkg/artifact"
	"github.com/neticdk/go-common/pkg/artifact/archive"
	"github.com/neticdk/go-common/pkg/github"
	"github.com/neticdk/go-stdlib/file"
)

type ReleaseOptions struct {
	Client     *gh.Client
	Downloader artifact.Downloader
	Uncompress func(src, dest string) error
	Owner      string
	Repository string
	AssetName  string
}

// PullRelease pulls a release from a github repository to the destination directory
func PullRelease(ctx context.Context, a *artifact.Artifact, opts *ReleaseOptions) (*artifact.PullResult, error) {
	if opts == nil {
		opts = &ReleaseOptions{}
	}

	if opts.Client == nil {
		return nil, fmt.Errorf("client is required")
	}
	client := opts.Client

	if opts.Downloader == nil {
		return nil, fmt.Errorf("downloader is required")
	}

	if opts.Uncompress == nil {
		opts.Uncompress = archive.Uncompress
	}

	tmpDir, err := os.MkdirTemp("", "go-common-")
	if err != nil {
		return nil, fmt.Errorf("creating temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	release, _, err := github.GetReleaseByTagOrLatest(ctx, client, opts.Owner, opts.Repository, a.Version)
	if err != nil {
		return nil, err
	}

	withExt := func(baseName, ext string) string {
		return fmt.Sprintf("%s.%s", baseName, ext)
	}

	var assetURL, assetName string
	for _, asset := range release.Assets {
		assetName = *asset.Name
		if assetName == withExt(opts.AssetName, "zip") {
			assetURL = *asset.BrowserDownloadURL
			break
		}
		if assetName == withExt(opts.AssetName, "tgz") {
			assetURL = *asset.BrowserDownloadURL
			break
		}
		if assetName == withExt(opts.AssetName, "tar.gz") {
			assetURL = *asset.BrowserDownloadURL
			break
		}
		if assetName == opts.AssetName {
			assetURL = *asset.BrowserDownloadURL
			break
		}
	}
	if assetURL == "" {
		return nil, fmt.Errorf("finding asset %q", opts.AssetName)
	}

	assetDestFile := filepath.Join(tmpDir, assetName)
	if _, err = opts.Downloader.Download(assetURL, assetDestFile); err != nil {
		return nil, fmt.Errorf(`downloading asset %q: %w`, assetURL, err)
	}

	tmpDstDir := filepath.Join(tmpDir, a.Name)
	if err := os.MkdirAll(tmpDstDir, file.FileModeNewDirectory); err != nil {
		return nil, fmt.Errorf("creating temporary directory: %w", err)
	}

	if err := opts.Uncompress(assetDestFile, tmpDstDir); err != nil {
		return nil, fmt.Errorf(`uncompressing asset %q: %w`, assetDestFile, err)
	}

	artifactDirName := opts.AssetName

	// Remove the extension from the artifact directory name
	artifactDirName = strings.TrimSuffix(artifactDirName, ".zip")
	artifactDirName = strings.TrimSuffix(artifactDirName, ".gz")
	artifactDirName = strings.TrimSuffix(artifactDirName, ".tar")
	artifactDirName = strings.TrimSuffix(artifactDirName, ".tgz")

	// Regular expression to match semver-like versions at the end of the string
	re := regexp.MustCompile(`[-_]?v?(\d+\.\d+\.\d+(-\w+(\.\d+)?)?)$`)
	artifactDirName = re.ReplaceAllString(artifactDirName, "")

	// Move the directory to the final destination
	if err := os.Rename(
		tmpDstDir,
		filepath.Join(a.BaseDir, artifactDirName)); err != nil {
		return nil, fmt.Errorf("renaming directory: %w", err)
	}

	artifactVersion := artifact.FirstVersionOrLatest(a.Version)

	return &artifact.PullResult{
		Dir:     filepath.Join(a.BaseDir, artifactDirName),
		Version: artifactVersion,
	}, nil
}
