package gorm

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/neticdk/go-common/pkg/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	gtrace "gorm.io/plugin/opentelemetry/tracing"
)

// ConfigureDatabase creates [gorm.DB] instance with slog based logger and OTEL tracing. Supports "sqlite" and "postgres" drivers.
func ConfigureDatabase(dbDriver, dbDSN string) (*gorm.DB, error) {
	var dialect gorm.Dialector
	switch dbDriver {
	case "sqlite":
		dialect = sqlite.Open(dbDSN)
	case "postgres":
		dialect = postgres.Open(dbDSN)
	default:
		return nil, fmt.Errorf("unsupported database driver type: %s", dbDriver)
	}

	db, err := gorm.Open(dialect, &gorm.Config{Logger: &zl{}})
	if err != nil {
		return nil, err
	}

	if err := db.Use(gtrace.NewPlugin()); err != nil {
		return nil, err
	}

	return db, nil
}

type zl struct{}

func (l *zl) LogMode(gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *zl) Error(ctx context.Context, msg string, opts ...any) {
	logger := log.FromContext(ctx)
	if logger.Enabled(ctx, slog.LevelError) {
		logger.ErrorContext(ctx, fmt.Sprintf(msg, opts...))
	}
}

func (l *zl) Warn(ctx context.Context, msg string, opts ...any) {
	logger := log.FromContext(ctx)
	if logger.Enabled(ctx, slog.LevelWarn) {
		logger.WarnContext(ctx, fmt.Sprintf(msg, opts...))
	}
}

func (l *zl) Info(ctx context.Context, msg string, opts ...any) {
	logger := log.FromContext(ctx)
	if logger.Enabled(ctx, slog.LevelInfo) {
		logger.InfoContext(ctx, fmt.Sprintf(msg, opts...))
	}
}

func (l *zl) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	logger := log.FromContext(ctx)

	if logger.Enabled(ctx, slog.LevelDebug) {
		attrs := []any{}

		dur := time.Since(begin)
		attrs = append(attrs, slog.Duration("duration", dur))

		sql, rows := fc()
		if rows != -1 {
			attrs = append(attrs, slog.Int64("rows", rows))
		}

		if err != nil {
			attrs = append(attrs, log.Error(err))
		}

		logger.DebugContext(ctx, sql, attrs...)
	}
}
