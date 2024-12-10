package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/neticdk/go-common/pkg/artifact"
	"github.com/neticdk/go-common/pkg/git"
)

type RepositoryOptions struct {
	Auth       transport.AuthMethod
	Repository git.Repository
}

// PullRepository pulls a git repository to the destination directory
func PullRepository(ctx context.Context, a *artifact.Artifact, opts *RepositoryOptions) (*artifact.PullResult, error) {
	if opts == nil {
		opts = &RepositoryOptions{}
	}

	if opts.Repository == nil {
		return nil, fmt.Errorf("repository is required")
	}

	cloneOptions := &gogit.CloneOptions{
		URL:           a.URL,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName("main"),
		Depth:         1,
	}

	if a.CommitHash == "" {
		if a.Tag != "" {
			cloneOptions.ReferenceName = plumbing.NewTagReferenceName(a.Tag)
		} else if a.Branch != "" {
			cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(a.Branch)
		}
	}

	if opts.Auth != nil {
		cloneOptions.Auth = opts.Auth
	}

	tmpDir, err := os.MkdirTemp("", "go-common-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpDstDir := filepath.Join(tmpDir, a.Name)
	if err := os.MkdirAll(tmpDstDir, 0o750); err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	repo, err := opts.Repository.PlainCloneContext(ctx, tmpDstDir, cloneOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	if a.CommitHash != "" {
		w, err := repo.Worktree()
		if err != nil {
			return nil, fmt.Errorf("failed to get worktree: %w", err)
		}
		hash := plumbing.NewHash(a.CommitHash)
		if err := w.Checkout(&gogit.CheckoutOptions{
			Hash:  hash,
			Force: true,
		}); err != nil {
			return nil, fmt.Errorf("failed to checkout commit: %w", err)
		}
	}
	if err := os.RemoveAll(filepath.Join(tmpDstDir, ".git")); err != nil {
		return nil, fmt.Errorf("failed to remove .git directory: %w", err)
	}

	if a.SubDir != "" {
		if !filepath.IsLocal(a.SubDir) {
			return nil, fmt.Errorf("subdir must be a local path")
		}
		tmpDstDir = filepath.Join(tmpDstDir, a.SubDir)
	}

	// Move the directory to the final destination
	if err := os.Rename(
		tmpDstDir,
		a.DestDir()); err != nil {
		return nil, fmt.Errorf("failed to rename directory: %w", err)
	}

	artifactVersion := artifact.FirstVersionOrLatest(a.Version)

	return &artifact.PullResult{
		Dir:     a.DestDir(),
		Version: artifactVersion,
	}, nil
}