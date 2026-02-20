package helm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/neticdk/go-stdlib/xstrings"
	"helm.sh/helm/v3/pkg/action"
	helmChart "helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/registry"
)

type pullOption struct {
	RegistryClient *registry.Client
	Version        string
}

type PullOption func(*pullOption)

type PullResult struct {
	Chart   *helmChart.Chart
	Version string
}

// WithRegistryClient sets the registry client to use
func WithRegistryClient(client *registry.Client) PullOption {
	return func(o *pullOption) {
		o.RegistryClient = client
	}
}

// WithVersion sets the version of the chart to pull
func WithVersion(v string) PullOption {
	return func(o *pullOption) {
		o.Version = v
	}
}

// PullChart pulls a helm chart to the destination directory
// Version defaults to "latest"
func PullChart(_ context.Context, repository, chartName, dstDir string, opts ...PullOption) (*PullResult, error) {
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	if chartName == "" {
		return nil, fmt.Errorf("chart reference is required")
	}

	opt := &pullOption{}
	for _, o := range opts {
		o(opt)
	}

	tmpDir, err := os.MkdirTemp("", "go-common-helm-")
	if err != nil {
		return nil, fmt.Errorf("creating temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	settings := cli.New()
	actionConfig := &action.Configuration{}
	if opt.RegistryClient != nil {
		actionConfig.RegistryClient = opt.RegistryClient
	}
	client := action.NewPullWithOpts(action.WithConfig(actionConfig))
	client.Settings = settings

	client.Untar = true
	client.UntarDir = tmpDir

	chartRef := chartName
	if registry.IsOCI(repository) {
		chartRef = repository + "/" + chartName
		if opt.RegistryClient == nil {
			registryClient, err := registry.NewClient()
			if err != nil {
				return nil, err
			}
			actionConfig.RegistryClient = registryClient
		}
	} else {
		client.RepoURL = repository
	}
	client.Version = opt.Version

	if _, err := client.Run(chartRef); err != nil {
		return nil, fmt.Errorf("pulling helm chart: %w", err)
	}

	tmpDstDir := filepath.Join(tmpDir, chartName)
	// Load and inspect the chart:
	chart, err := loader.Load(tmpDstDir)
	if err != nil {
		return nil, fmt.Errorf("loading chart: %w", err)
	}

	// Move the directory to the final destination
	if err := os.Rename(
		tmpDstDir,
		dstDir); err != nil {
		return nil, fmt.Errorf("renaming directory: %w", err)
	}

	artifactVersion := xstrings.Coalesce(chart.Metadata.Version, opt.Version, "latest")

	return &PullResult{
		Chart:   chart,
		Version: artifactVersion,
	}, nil
}
