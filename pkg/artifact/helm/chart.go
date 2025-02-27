package helm

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/neticdk/go-common/pkg/artifact"
	"github.com/neticdk/go-common/pkg/helm"
)

type ChartOptions struct{}

// PullChart pulls a helm chart to the destination directory
func PullChart(ctx context.Context, a *artifact.Artifact, _ *ChartOptions) (*artifact.PullResult, error) {
	res, err := helm.PullChart(ctx, a.Repository, a.Name, a.DestDir())
	if err != nil {
		return nil, fmt.Errorf("pulling helm chart: %w", err)
	}

	return &artifact.PullResult{
		Dir:     filepath.Join(a.BaseDir, a.Name),
		Version: res.Version,
	}, nil
}
