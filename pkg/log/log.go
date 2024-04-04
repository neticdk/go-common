// Package log contains utility functions around Go slog package
package log

import (
	"context"
	"log/slog"
)

type ctxKey struct{}

// WithLogger returns a context containing a reference to the given logger
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

// FromContext retrieves logger from the given context. If no logger is associated with the given context
// it will return value from [slog.Default].
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

// Error returns an [slog.Attr] for an [error]
func Error(err error) slog.Attr {
	return slog.Any("error", err)
}
