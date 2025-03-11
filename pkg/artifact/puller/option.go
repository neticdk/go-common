package puller

import (
	"log/slog"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/neticdk/go-common/pkg/artifact/archive"
	agit "github.com/neticdk/go-common/pkg/artifact/git"
	"github.com/neticdk/go-common/pkg/artifact/github"
	"github.com/neticdk/go-common/pkg/artifact/helm"
	"github.com/neticdk/go-common/pkg/artifact/jsonnetbundler"
	"github.com/neticdk/go-common/pkg/git"
)

type Option func(*puller)

// WithLogger sets the logger for the puller
func WithLogger(l *slog.Logger) Option {
	return func(p *puller) {
		p.logger = l
	}
}

// WithGitRepository sets the git repository for the puller
func WithGitRepository(r git.Repository) Option {
	return func(p *puller) {
		if p.gitRepositoryOptions == nil {
			p.gitRepositoryOptions = &agit.RepositoryOptions{}
		}
		p.gitRepositoryOptions.Repository = r
	}
}

type PullOption func(*puller)

// WithHelmOptions sets the helm options for the puller
func WithHelmOptions(o *helm.ChartOptions) PullOption {
	return func(p *puller) {
		p.helmChartOptions = o
	}
}

// WithHTTPOptions sets the HTTP options for the puller
func WithHTTPOptions(o *archive.HTTPOptions) PullOption {
	return func(p *puller) {
		p.archiveHTTPOptions = o
	}
}

// WithRepositoryAuth sets the auth transport on the git repository
func WithRepositoryAuth(a transport.AuthMethod) PullOption {
	return func(p *puller) {
		p.gitRepositoryOptions.Auth = a
	}
}

// WithGithubReleaseOptions sets the github release options for the puller
func WithGithubReleaseOptions(o *github.ReleaseOptions) PullOption {
	return func(p *puller) {
		p.githubReleaseOptions = o
	}
}

// WithJsonnetBundlerOptions sets the jsonnet bundler options for the puller
func WithJsonnetBundlerOptions(o *jsonnetbundler.PackageOptions) PullOption {
	return func(p *puller) {
		p.jsonnetBundlerOptions = o
	}
}
