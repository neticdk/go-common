package puller

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/neticdk/go-common/pkg/artifact"
	"github.com/neticdk/go-common/pkg/artifact/archive"
	"github.com/neticdk/go-common/pkg/artifact/git"
	"github.com/neticdk/go-common/pkg/artifact/github"
	"github.com/neticdk/go-common/pkg/artifact/helm"
	"github.com/neticdk/go-common/pkg/artifact/jsonnetbundler"
)

// PullMethod is the method used to pull an artifact
type PullMethod int

const (
	PullMethodUnknown PullMethod = iota
	PullMethodHelmChart
	PullMethodGithubRelease
	PullMethodHTTPArchive
	PullMethodGit
	PullMethodJsonnetBundler
)

// Puller is the interface for pulling artifacts
type Puller interface {
	// Pull pulls an artifact from a source to a destination
	Pull(context.Context, PullMethod, *artifact.Artifact, ...PullOption) (*artifact.PullResult, error)
}

type puller struct {
	logger                *slog.Logger
	helmChartOptions      *helm.ChartOptions
	githubReleaseOptions  *github.ReleaseOptions
	archiveHTTPOptions    *archive.HTTPOptions
	gitRepositoryOptions  *git.RepositoryOptions
	jsonnetBundlerOptions *jsonnetbundler.PackageOptions
}

// NewPuller creates a new puller
func NewPuller(opts ...Option) Puller {
	p := &puller{}
	for _, opt := range opts {
		opt(p)
	}
	if p.logger == nil {
		handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
		p.logger = slog.New(handler)
	}
	return p
}

// Pull pulls an artifact from a source to a destination
func (p *puller) Pull(ctx context.Context, m PullMethod, a *artifact.Artifact, opts ...PullOption) (*artifact.PullResult, error) {
	for _, opt := range opts {
		opt(p)
	}

	if err := a.DestDirExists(); err != nil {
		return nil, err
	}

	switch m {
	case PullMethodHelmChart:
		return helm.PullChart(ctx, a, p.helmChartOptions)
	case PullMethodGithubRelease:
		return github.PullRelease(ctx, a, p.githubReleaseOptions)
	case PullMethodHTTPArchive:
		return archive.PullHTTP(ctx, a, p.archiveHTTPOptions)
	case PullMethodGit:
		return git.PullRepository(ctx, a, p.gitRepositoryOptions)
	case PullMethodJsonnetBundler:
		return jsonnetbundler.Install(ctx, a, p.jsonnetBundlerOptions)
	default:
		return nil, fmt.Errorf("uknown pull method: %d", m)
	}
}
