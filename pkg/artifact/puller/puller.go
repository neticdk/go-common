package puller

import (
	"context"
	"fmt"
	"os"

	"github.com/neticdk/go-common/pkg/artifact"
	"github.com/neticdk/go-common/pkg/artifact/archive"
	"github.com/neticdk/go-common/pkg/artifact/git"
	"github.com/neticdk/go-common/pkg/artifact/github"
	"github.com/neticdk/go-common/pkg/artifact/helm"
	logger "github.com/neticdk/go-common/pkg/tui/logger/charm"
)

// PullMethod is the method used to pull an artifact
type PullMethod int

const (
	PullMethodUnknown PullMethod = iota
	PullMethodHelmChart
	PullMethodGithubRelease
	PullMethodHTTPArchive
	PullMethodGit
)

// Puller is the interface for pulling artifacts
type Puller interface {
	// Pull pulls an artifact from a source to a destination
	Pull(context.Context, PullMethod, *artifact.Artifact, ...PullOption) (*artifact.PullResult, error)
}

type puller struct {
	logger               logger.Logger
	helmChartOptions     *helm.ChartOptions
	githubReleaseOptions *github.ReleaseOptions
	archiveHTTPOptions   *archive.HTTPOptions
	gitRepositoryOptions *git.RepositoryOptions
}

// NewPuller creates a new puller
func NewPuller(opts ...PullerOption) Puller {
	p := &puller{}
	for _, opt := range opts {
		opt(p)
	}
	if p.logger == nil {
		p.logger = logger.New(os.Stderr, "info")
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
	default:
		return nil, fmt.Errorf("uknown pull method: %d", m)
	}
}
