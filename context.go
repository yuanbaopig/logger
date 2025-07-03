package logger

import (
	"context"
	"go.uber.org/zap"
)

type key int

const (
	logContextKey key = iota
)

func WithContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, logContextKey, l)
}

// WithName adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func WithName(s string) *zap.Logger {
	return Log.Named(s)
}

// FromContext returns the value of the log key on the ctx.
func FromContext(ctx context.Context) *zap.Logger {
	if ctx != nil {
		logger := ctx.Value(logContextKey)
		if logger != nil {
			return logger.(*zap.Logger)
		}
	}

	return WithName("Unknown-Context")
}
