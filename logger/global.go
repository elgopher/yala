package logger

import (
	"context"
	"sync/atomic"
)

// Global is a logger shared globally. You can use it to define global logger for your package:
//
//		package yourpackage
//		import "github.com/jacekolszak/yala/logger"
//		var Logger logger.Global // define global logger, no need to initialize (by default nothing is logged)
//
//
// It is safe to use it concurrently.
type Global struct {
	logger atomic.Value
}

// SetAdapter updates adapter implementation. By default, nothing is logged.
//
// It can be run anytime. Please note though that this method is meant to be used by end user, configuring logging
// from the central place (such as main.go or any other package setting up the entire application).
func (g *Global) SetAdapter(adapter Adapter) {
	if adapter == nil {
		adapter = noopLogger{}
	}

	g.logger.Store(LocalLogger{adapter: adapter})
}

func (g *Global) getLogger() LocalLogger {
	logger, ok := g.logger.Load().(LocalLogger)
	if !ok {
		g.logger.Store(LocalLogger{adapter: &initialGlobalNoopLogger{}})

		return g.getLogger()
	}

	return logger
}

// Debug logs message using globally configured logger.Adapter.
func (g *Global) Debug(ctx context.Context, msg string) {
	g.loggerWithSkippedCallerFrame(ctx).Debug(msg)
}

func (g *Global) loggerWithSkippedCallerFrame(ctx context.Context) Logger {
	return g.getLogger().WithSkippedCallerFrame(ctx)
}

// Info logs message using globally configured logger.Adapter.
func (g *Global) Info(ctx context.Context, msg string) {
	g.loggerWithSkippedCallerFrame(ctx).Info(msg)
}

// Warn logs message using globally configured logger.Adapter.
func (g *Global) Warn(ctx context.Context, msg string) {
	g.loggerWithSkippedCallerFrame(ctx).Warn(msg)
}

// Error logs message using globally configured logger.Adapter.
func (g *Global) Error(ctx context.Context, msg string) {
	g.loggerWithSkippedCallerFrame(ctx).Error(msg)
}

// With creates a new Logger with field and using globally configured logger.Adapter.
func (g *Global) With(ctx context.Context, key string, value interface{}) Logger {
	return g.getLogger().With(ctx, key, value)
}

// WithError creates a new Logger with error and using globally configured logger.Adapter.
func (g *Global) WithError(ctx context.Context, err error) Logger {
	return g.getLogger().WithError(ctx, err)
}
