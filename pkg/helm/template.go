package helm

import (
	"context"
	"fmt"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
)

const (
	DefaultNamespace = "default"
)

// TemplateChart renders a Helm chart with the given namespace and values
// If namespace is empty, it defaults to "default"
func TemplateChart(ctx context.Context, chart *chart.Chart, namespace string, vals map[string]interface{}) (string, error) {
	if chart == nil {
		return "", fmt.Errorf("chart is required")
	}
	if vals == nil {
		vals = map[string]interface{}{}
	}
	if namespace == "" {
		namespace = DefaultNamespace
	}
	mem := driver.NewMemory()
	mem.SetNamespace(namespace)
	cfg := &action.Configuration{
		Releases: storage.Init(mem),
	}
	install := action.NewInstall(cfg)
	install.ClientOnly = true
	install.Namespace = namespace
	install.ReleaseName = chart.Name()
	install.IncludeCRDs = false
	// Since we're not applying the manifests we're setting kubeVersion to ignore the check
	chart.Metadata.KubeVersion = ""

	release, err := install.RunWithContext(ctx, chart, vals)
	if err != nil {
		return "", fmt.Errorf("failed to install Helm chart: %w", err)
	}

	return release.Manifest, nil
}
