package git

import (
	"context"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

// Repository is an interface for interacting with git repositories
type Repository interface {
	Repo() *git.Repository
	Init(path string, branch string) (*git.Repository, error)
	PlainCloneContext(ctx context.Context, path string, o *git.CloneOptions) (*git.Repository, error)
	CreateRemote(name string, url string) (*git.Remote, error)
	Add(paths ...string) error
	Commit(message string) (plumbing.Hash, error)
	Push(o *git.PushOptions) error
	InitAndCommit(dir string, url string, cfg *config.Config) error
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
func (g *gitRepository) Commit(message string) (plumbing.Hash, error) {
	if g.repo == nil {
		return plumbing.ZeroHash, fmt.Errorf("repository is not yet initialized")
	}
	worktree, err := g.repo.Worktree()
	if err != nil {
		return plumbing.ZeroHash, fmt.Errorf("getting worktree: %w", err)
	}
	hash, err := worktree.Commit(message, &git.CommitOptions{})
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
func (g *gitRepository) InitAndCommit(path string, url string, cfg *config.Config) error {
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
	if _, err := g.Commit("Initial commit"); err != nil {
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
