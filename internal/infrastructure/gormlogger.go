package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _ logger.Interface = (*gormLogger)(nil)

type gormLogger struct {
	logger               *slog.Logger
	logLevel             logger.LogLevel
	slowThreshold        time.Duration
	ignoreRecordNotFound bool
}

func newGormLogger(l *slog.Logger, level logger.LogLevel) *gormLogger {
	if l == nil {
		l = slog.Default()
	}

	return &gormLogger{
		logger:               l,
		logLevel:             level,
		slowThreshold:        200 * time.Millisecond,
		ignoreRecordNotFound: false,
	}
}

func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	clone := *l
	clone.logLevel = level
	return &clone
}

func (l *gormLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.logLevel < logger.Info {
		return
	}

	l.logger.InfoContext(ctx, fmt.Sprintf(msg, args...))
}

func (l *gormLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.logLevel < logger.Warn {
		return
	}

	l.logger.WarnContext(ctx, fmt.Sprintf(msg, args...))
}

func (l *gormLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.logLevel < logger.Error {
		return
	}

	l.logger.ErrorContext(ctx, fmt.Sprintf(msg, args...))
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel == logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil && l.logLevel >= logger.Error {
		if l.ignoreRecordNotFound && errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}

		l.logger.ErrorContext(ctx, "Gorm query failed",
			slog.String("sql", sql),
			slog.Duration("elapsed", elapsed),
			slog.Int64("rows", rows),
			slog.String("error", err.Error()),
		)

		return
	}

	if l.slowThreshold > 0 && elapsed > l.slowThreshold && l.logLevel >= logger.Warn {
		l.logger.WarnContext(ctx, "Gorm slow query",
			slog.String("sql", sql),
			slog.Duration("elapsed", elapsed),
			slog.Duration("threshold", l.slowThreshold),
			slog.Int64("rows", rows),
		)
		return
	}

	if l.logLevel >= logger.Info {
		l.logger.DebugContext(ctx, "Gorm query",
			slog.String("sql", sql),
			slog.Duration("elapsed", elapsed),
			slog.Int64("rows", rows),
		)
	}

}
