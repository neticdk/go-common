package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/neticdk/go-common/pkg/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

// config holds the configuration for telemetry.
type config struct {
	startServer bool
}

// Option is a functional option for configuring telemetry.
type Option func(*config)

// WithoutMetricsServer disables the startup of the metrics server.
func WithoutMetricsServer() Option {
	return func(c *config) {
		c.startServer = false
	}
}

// ConfigureTelemetry will configure OpenTelemetry to expose metrics and traces using Prometheus export for metrics and OTEL grpc exporter for traces. The
// given port is the port to expose metrics, the given serviceName is the OTEL service name attribute associated with all traces. The function will return
// a shutdown function that can be called when shutting down the process.
func ConfigureTelemetry(metricsPort int, serviceName string, opts ...Option) (func(context.Context) error, error) {
	cfg := &config{
		startServer: true,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	ctx := context.Background()
	logger := slog.Default()
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}
	otel.SetMeterProvider(metric.NewMeterProvider(metric.WithReader(exporter)))

	if cfg.startServer {
		go func() {
			logger.InfoContext(ctx, "starting metrics server")
			mux := http.NewServeMux()
			mux.Handle("/metrics", promhttp.Handler())
			s := http.Server{
				Addr:              fmt.Sprintf(":%d", metricsPort),
				Handler:           mux,
				ReadHeaderTimeout: 3 * time.Second,
			}
			if err := s.ListenAndServe(); err != nil {
				logger.ErrorContext(ctx, "metrics listener failed", log.Error(err))
			}
		}()
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("creating trace exporter: %w", err)
	}

	bsp := trace.NewBatchSpanProcessor(traceExporter)
	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return provider.Shutdown, nil
}
