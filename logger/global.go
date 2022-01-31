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
	entry       Entry
	adapter     atomic.Value  // not used when logger is a child
	rootAdapter *atomic.Value // not nil, if logger is a child
}

type adapterWrapper struct{ Adapter } // stored in atomic.Value

// SetAdapter updates adapter implementation. By default, nothing is logged.
//
// It can be run anytime. Please note though that this method is meant to be used by end user, configuring logging
// from the central place (such as main.go or any other package setting up the entire application).
//
// If this method is called on an instance created using With* methods, then all parent and child loggers
// are updated too.
func (g *Global) SetAdapter(adapter Adapter) {
	if adapter == nil {
		adapter = noopAdapter{}
	}

	if g.rootAdapter != nil {
		g.rootAdapter.Store(adapterWrapper{Adapter: adapter})

		return
	}

	g.adapter.Store(adapterWrapper{Adapter: adapter})
}

func (g *Global) getAdapter() Adapter { // nolint:ireturn
	value := g.adapterValue()

	adapter, ok := value.Load().(Adapter)
	if !ok {
		value.CompareAndSwap(nil, adapterWrapper{Adapter: &initialGlobalNoopAdapter{}})

		return g.getAdapter()
	}

	return adapter
}

func (g *Global) adapterValue() *atomic.Value {
	if g.rootAdapter != nil {
		return g.rootAdapter
	}

	return &g.adapter
}

// Debug logs a message at DebugLevel.
func (g *Global) Debug(ctx context.Context, msg string) {
	g.log(ctx, DebugLevel, msg)
}

func (g *Global) log(ctx context.Context, level Level, msg string) {
	e := g.entry
	e.Level = level
	e.Message = msg
	e.SkippedCallerFrames += 2

	g.getAdapter().Log(ctx, e)
}

// Info logs a message at InfoLevel.
func (g *Global) Info(ctx context.Context, msg string) {
	g.log(ctx, InfoLevel, msg)
}

// Warn logs a message at WarnLevel.
func (g *Global) Warn(ctx context.Context, msg string) {
	g.log(ctx, WarnLevel, msg)
}

// Error logs a message at ErrorLevel.
func (g *Global) Error(ctx context.Context, msg string) {
	g.log(ctx, ErrorLevel, msg)
}

// With creates a new child logger with field.
func (g *Global) With(key string, value interface{}) *Global {
	newEntry := g.entry.With(Field{Key: key, Value: value})

	return &Global{
		entry:       newEntry,
		rootAdapter: g.adapterValue(),
	}
}

// WithError creates a new child logger with error.
func (g *Global) WithError(err error) *Global {
	newEntry := g.entry
	newEntry.Error = err

	return &Global{
		entry:       newEntry,
		rootAdapter: g.adapterValue(),
	}
}
