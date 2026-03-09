package helm

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/neticdk/go-common/pkg/artifact"
	"github.com/neticdk/go-common/pkg/helm"
	"helm.sh/helm/v3/pkg/registry"
)

type ChartOptions struct {
	RegistryClient *registry.Client
	Version        string
}

// PullChart pulls a helm chart to the destination directory
func PullChart(ctx context.Context, a *artifact.Artifact, opts *ChartOptions) (*artifact.PullResult, error) {
	if opts == nil {
		opts = &ChartOptions{}
	}
	res, err := helm.PullChart(ctx, a.Repository, a.Name, a.DestDir(), helm.WithRegistryClient(opts.RegistryClient), helm.WithVersion(opts.Version))
	if err != nil {
		return nil, fmt.Errorf("pulling helm chart: %w", err)
	}

	return &artifact.PullResult{
		Dir:     filepath.Join(a.BaseDir, a.Name),
		Version: res.Version,
	}, nil
}
