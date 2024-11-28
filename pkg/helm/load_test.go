package helm

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadChartFromFS(t *testing.T) {
	testChartPath := "testdata/test-chart"

	filesystem := os.DirFS(testChartPath)
	loadedChart, err := LoadChartFromFS(filesystem)

	assert.NoError(t, err)
	assert.NotNil(t, loadedChart)
	assert.Equal(t, "test-chart", loadedChart.Metadata.Name)
}
