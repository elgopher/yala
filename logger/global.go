// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logger

import (
	"context"
	"sync/atomic"
)

// Global is a logger shared globally. You can use it to define global logger for your package:
//
//		package yourpackage
//		import "github.com/elgopher/yala/logger"
//		var log logger.Global // define global logger, no need to initialize (by default nothing is logged)
//
//		func SetLoggerAdapter(adapter logger.Adapter) {
//			log.SetAdapter(adapter)
//		}
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
		adapter = noopAdapter{}
	}

	g.logger.Store(Local{Adapter: adapter})
}

func (g *Global) getLogger() Local {
	logger, ok := g.logger.Load().(Local)
	if !ok {
		g.logger.CompareAndSwap(nil, Local{Adapter: &initialGlobalNoopAdapter{}})

		return g.getLogger()
	}

	return logger
}

// Debug logs message using globally configured logger.Adapter.
func (g *Global) Debug(ctx context.Context, msg string) {
	g.loggerWithSkippedCallerFrame().Debug(ctx, msg)
}

func (g *Global) loggerWithSkippedCallerFrame() Logger {
	return g.getLogger().WithSkippedCallerFrame()
}

// Info logs message using globally configured logger.Adapter.
func (g *Global) Info(ctx context.Context, msg string) {
	g.loggerWithSkippedCallerFrame().Info(ctx, msg)
}

// Warn logs message using globally configured logger.Adapter.
func (g *Global) Warn(ctx context.Context, msg string) {
	g.loggerWithSkippedCallerFrame().Warn(ctx, msg)
}

// Error logs message using globally configured logger.Adapter.
func (g *Global) Error(ctx context.Context, msg string) {
	g.loggerWithSkippedCallerFrame().Error(ctx, msg)
}

// With creates a new Logger with field and using globally configured logger.Adapter.
func (g *Global) With(key string, value interface{}) Logger {
	return g.getLogger().With(key, value)
}

// WithError creates a new Logger with error and using globally configured logger.Adapter.
func (g *Global) WithError(err error) Logger {
	return g.getLogger().WithError(err)
}
