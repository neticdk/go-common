package helm

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v4/pkg/chart"
)

func TestTemplateChart(t *testing.T) {
	// Define the path to the test chart
	testChartPath := "testdata/test-chart"

	// Create a filesystem from the test chart directory
	filesystem := os.DirFS(testChartPath)

	// Load the chart using the LoadChartFromFS function
	loadedChart, err := LoadChartFromFS(filesystem)
	assert.NoError(t, err)
	assert.NotNil(t, loadedChart)

	// Define test cases
	tests := []struct {
		name      string
		namespace string
		values    map[string]interface{}
		expectErr bool
	}{
		{
			name:      "default values",
			namespace: "default",
			values:    map[string]interface{}{},
			expectErr: false,
		},
		{
			name:      "custom namespace",
			namespace: "custom-namespace",
			values:    map[string]interface{}{},
			expectErr: false,
		},
		{
			name:      "with values",
			namespace: "default",
			values: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			expectErr: false,
		},
		{
			name:      "nil chart",
			namespace: "default",
			values:    map[string]interface{}{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var chart *chart.Chart
			if tt.name != "nil chart" {
				chart = loadedChart
			}

			// Render the chart template with the provided values
			rendered, err := TemplateChart(context.Background(), chart, tt.namespace, tt.values)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, rendered)
			}
		})
	}
}
