package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestConfigureTelemetryTracesExporterNone(t *testing.T) {
	t.Setenv("OTEL_TRACES_EXPORTER", "none")

	shutdown, err := ConfigureTelemetry(0, "test", WithoutMetricsServer())
	require.NoError(t, err)
	assert.NoError(t, shutdown(context.Background()))

	_, isSDK := otel.GetTracerProvider().(*trace.TracerProvider)
	assert.False(t, isSDK, "tracer provider should stay the no-op default")
}
