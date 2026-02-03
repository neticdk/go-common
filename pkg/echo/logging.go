package echo

import (
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/neticdk/go-common/pkg/log"
)

// SlogContext is a middleware function to add the given logger is to the
// request context. The logger may then be retrieved from the context using
// [log.FromContext]
func SlogContext(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx := log.WithLogger(c.Request().Context(), logger)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)
			return next(c)
		}
	}
}

// RequestLogger is a middleware function to log requests in echo using slog.
// It will try to retrieve a logger from the [context.Context] but will fall
// back to the default logger.
func RequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogMethod:    true,
		LogHost:      true,
		LogUserAgent: true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []any{}
			attrs = append(attrs, slog.String("URI", v.URI))
			attrs = append(attrs, slog.Int("status", v.Status))
			attrs = append(attrs, slog.String("method", v.Method))
			attrs = append(attrs, slog.String("host", v.Host))
			attrs = append(attrs, slog.String("user_agent", v.UserAgent))
			attrs = append(attrs, slog.Duration("latency", v.Latency))
			attrs = append(attrs, slog.String("remote_ip", v.RemoteIP))

			if v.Error != nil {
				attrs = append(attrs, log.Error(v.Error))
			}

			logger := log.FromContext(c.Request().Context())
			logger.InfoContext(c.Request().Context(), "request", attrs...)

			return nil
		},
	})
}
