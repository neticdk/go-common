package log

import (
	"context"
	"io"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type wrapper struct {
	handler slog.Handler
}

// NewJSONTraceIDHandler creates new [slog.JSONHandler] which adds a `TraceID` attribute to all loggings if trace id is present.
func NewJSONTraceIDHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	h := &wrapper{
		handler: slog.NewJSONHandler(w, opts),
	}
	return h
}

func (h *wrapper) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *wrapper) Handle(ctx context.Context, record slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().HasTraceID() {
		record = record.Clone()
		record.AddAttrs(slog.String("TraceID", span.SpanContext().TraceID().String()))
	}
	return h.handler.Handle(ctx, record)
}

func (h *wrapper) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &wrapper{
		handler: h.handler.WithAttrs(attrs),
	}
}

func (h *wrapper) WithGroup(name string) slog.Handler {
	return &wrapper{
		handler: h.handler.WithGroup(name),
	}
}
