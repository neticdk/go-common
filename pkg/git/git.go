package git

import (
	"context"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	ssh2 "golang.org/x/crypto/ssh"
)

// Repository is an interface for interacting with git repositories
type Repository interface {
	// Repo returns the git repository
	Repo() *git.Repository

	// Init initializes a new git repository in the specified path with the
	Init(path string, branch string) (*git.Repository, error)

	// PlainCloneContext clones a git repository to the specified path
	PlainCloneContext(ctx context.Context, path string, o *git.CloneOptions) (*git.Repository, error)

	// CreateRemote creates a new remote with the specified name and url
	CreateRemote(name string, url string) (*git.Remote, error)

	// SetUpstream sets the upstream for the specified local branch
	// to the specified remote branch.
	SetUpstream(local, remote, branch string) error

	// Add adds the specified files to the repository
	Add(paths ...string) error

	// Commit commits the changes to the repository
	Commit(message string, opts *git.CommitOptions) (plumbing.Hash, error)

	// Push pushes the changes to the remote
	Push(o *git.PushOptions) error

	// InitAndCommit initializes a git repository, adds all files in the directory,
	InitAndCommit(dir string, url string, cfg *config.Config, opts *git.CommitOptions) error

	// Config returns the repository configuration
	Config() (*config.Config, error)
}

type gitRepository struct {
	repo *git.Repository
}

// NewGitRepository returns a new gitRepository
func NewRepository(repo *git.Repository) Repository {
	r := &gitRepository{}
	if repo != nil {
		r.repo = repo
	}
	return r
}

// Repo returns the git repository
func (g *gitRepository) Repo() *git.Repository {
	return g.repo
}

// Init initializes a new git repository in the specified path with the
// specified branch
func (g *gitRepository) Init(path string, branch string) (*git.Repository, error) {
	ref := plumbing.NewBranchReferenceName(branch)
	if err := ref.Validate(); err != nil {
		return nil, err
	}
	repo, err := git.PlainInitWithOptions(path, &git.PlainInitOptions{
		InitOptions: git.InitOptions{
			DefaultBranch: ref,
		},
		Bare: false,
	})
	if err != nil {
		return nil, fmt.Errorf("initializing git repository: %w", err)
	}
	g.repo = repo
	return repo, nil
}

// PlainCloneContext clones a git repository to the specified path
func (g *gitRepository) PlainCloneContext(ctx context.Context, path string, o *git.CloneOptions) (*git.Repository, error) {
	repo, err := git.PlainCloneContext(ctx, path, false, o)
	if err != nil {
		return nil, fmt.Errorf("cloning repository: %w", err)
	}
	g.repo = repo
	return repo, nil
}

// CreateRemote creates a new remote with the specified name and url
func (g *gitRepository) CreateRemote(name string, url string) (*git.Remote, error) {
	if g.repo == nil {
		return nil, fmt.Errorf("repository is not yet initialized")
	}
	remote, err := g.repo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{url},
	})
	if err != nil {
		return nil, fmt.Errorf("creating remote: %w", err)
	}
	return remote, nil
}

// SetUpstream sets the upstream for the specified local branch
// to the specified branch for the remote repository.
// Example:
//
//	// Set local branch "main" to track remote branch "main"
//	gitRepository.SetUpstream("main", "origin", "main")
//	// Set local branch "fix/issue" to track remote branch "develop"
//	gitRepository.SetUpstream("fix/issue", "origin", "develop")
func (g *gitRepository) SetUpstream(local, remote, branch string) error {
	remoteRef := plumbing.NewRemoteReferenceName(remote, branch)
	ref, err := g.Repo().Reference(remoteRef, true)
	if err != nil {
		return fmt.Errorf("getting reference: %w", err)
	}

	mergeRef := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))
	if err := g.Repo().CreateBranch(&config.Branch{Name: local, Remote: remote, Merge: mergeRef}); err != nil {
		return fmt.Errorf("creating branch: %w", err)
	}

	localRef := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", local))
	if err := g.Repo().Storer.SetReference(plumbing.NewHashReference(localRef, ref.Hash())); err != nil {
		return fmt.Errorf("setting reference: %w", err)
	}
	return nil
}

// Add adds the specified files to the repository
func (g *gitRepository) Add(paths ...string) error {
	if g.repo == nil {
		return fmt.Errorf("repository is not yet initialized")
	}
	worktree, err := g.repo.Worktree()
	if err != nil {
		return fmt.Errorf("getting worktree: %w", err)
	}
	for _, path := range paths {
		_, err = worktree.Add(path)
		if err != nil {
			return fmt.Errorf("adding files: %w", err)
		}
	}
	return nil
}

// Commit commits the changes to the repository
func (g *gitRepository) Commit(message string, opts *git.CommitOptions) (plumbing.Hash, error) {
	if g.repo == nil {
		return plumbing.ZeroHash, fmt.Errorf("repository is not yet initialized")
	}
	if opts == nil {
		opts = &git.CommitOptions{}
	}
	worktree, err := g.repo.Worktree()
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("getting worktree: %w", err)
	}
	hash, err := worktree.Commit(message, opts)
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("committing changes: %w", err)
	}
	return hash, nil
}

// Push pushes the changes to the remote
func (g *gitRepository) Push(o *git.PushOptions) error {
	if g.repo == nil {
		return fmt.Errorf("repository is not yet initialized")
	}
	err := g.repo.Push(o)
	if err != nil {
		return fmt.Errorf("pushing changes: %w", err)
	}
	return nil
}

// InitAndCommit initializes a git repository, adds all files in the directory,
// and commits them
func (g *gitRepository) InitAndCommit(path string, url string, cfg *config.Config, opts *git.CommitOptions) error {
	repo, err := g.Init(path, "main")
	if err != nil {
		return err
	}
	if cfg != nil {
		if err := repo.SetConfig(cfg); err != nil {
			return err
		}
	}
	g.repo = repo
	if _, err := g.CreateRemote("origin", url); err != nil {
		return err
	}
	if err := g.Add("."); err != nil {
		return err
	}
	if _, err := g.Commit("Initial commit", opts); err != nil {
		return err
	}
	return nil
}

// Config returns the repository configuration
func (g *gitRepository) Config() (*config.Config, error) {
	if g.repo == nil {
		return nil, fmt.Errorf("repository is not yet initialized")
	}
	cfg, err := g.repo.Config()
	if err != nil {
		return nil, fmt.Errorf("getting git config: %w", err)
	}
	return cfg, nil
}

func SSHKeyAuth(privateKey []byte, user string, ignoreHostKeys bool) (transport.AuthMethod, error) { //nolint:revive
	if user == "" {
		user = "git"
	}
	v, err := ssh.NewPublicKeys(user, privateKey, "")
	if err != nil {
		return nil, fmt.Errorf("creating SSH public keys from private key: %w", err)
	}
	if ignoreHostKeys {
		//nolint:gosec
		v.HostKeyCallback = ssh2.InsecureIgnoreHostKey()
	}
	return v, nil
}
