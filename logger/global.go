package logger

import (
	"context"
	"sync/atomic"
)

type global struct {
	logger atomic.Value
}

func (g *global) SetAdapter(adapter Adapter) {
	if adapter == nil {
		adapter = noopLogger{}
	}

	g.logger.Store(LocalLogger{adapter: adapter})
}

func (g *global) getLogger() LocalLogger {
	return g.logger.Load().(LocalLogger)
}

// Debug logs message using globally configured logger.Adapter.
func (g *global) Debug(ctx context.Context, msg string) {
	g.loggerWithSkippedCallerFrame(ctx).Debug(msg)
}

func (g *global) loggerWithSkippedCallerFrame(ctx context.Context) Logger {
	return g.getLogger().WithSkippedCallerFrame(ctx).WithSkippedCallerFrame()
}

// Info logs message using globally configured logger.Adapter.
func (g *global) Info(ctx context.Context, msg string) {
	g.loggerWithSkippedCallerFrame(ctx).Info(msg)
}

// Warn logs message using globally configured logger.Adapter.
func (g *global) Warn(ctx context.Context, msg string) {
	g.loggerWithSkippedCallerFrame(ctx).Warn(msg)
}

// Error logs message using globally configured logger.Adapter.
func (g *global) Error(ctx context.Context, msg string) {
	g.loggerWithSkippedCallerFrame(ctx).Error(msg)
}

// With creates a new Logger with field and using globally configured logger.Adapter.
func (g *global) With(ctx context.Context, key string, value interface{}) Logger {
	return g.getLogger().With(ctx, key, value)
}

// WithError creates a new Logger with error and using globally configured logger.Adapter.
func (g *global) WithError(ctx context.Context, err error) Logger {
	return g.getLogger().WithError(ctx, err)
}
