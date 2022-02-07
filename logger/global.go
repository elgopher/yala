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
//
//		var log logger.Global // define global logger, no need to initialize (by default nothing is logged)
//
//		func SetLoggerAdapter(adapter logger.Adapter) {
//			log.SetAdapter(adapter)
//		}
//
// It is safe to use it concurrently.
//
// Please do not copy logger.Global instance. If you want to create a child logger, please use With
// and WithError methods. These methods will create *logger.Global using shared adapter.
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

	g.adapterValue().Store(adapterWrapper{Adapter: adapter})
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
	g.log(ctx, DebugLevel, msg, g.entry.Error, nil)
}

// DebugFields logs a message at DebugLevel with fields.
func (g *Global) DebugFields(ctx context.Context, msg string, fields Fields) {
	g.log(ctx, DebugLevel, msg, g.entry.Error, fields)
}

func (g *Global) log(ctx context.Context, level Level, msg string, cause error, fields Fields) {
	newEntry := g.entry.WithFields(fields)
	newEntry.Level = level
	newEntry.Message = msg
	newEntry.Error = cause
	newEntry.SkippedCallerFrames += 2

	g.getAdapter().Log(ctx, newEntry)
}

// Info logs a message at InfoLevel.
func (g *Global) Info(ctx context.Context, msg string) {
	g.log(ctx, InfoLevel, msg, g.entry.Error, nil)
}

// InfoFields logs a message at InfoLevel with fields.
func (g *Global) InfoFields(ctx context.Context, msg string, fields Fields) {
	g.log(ctx, InfoLevel, msg, g.entry.Error, fields)
}

// Warn logs a message at WarnLevel.
func (g *Global) Warn(ctx context.Context, msg string) {
	g.log(ctx, WarnLevel, msg, g.entry.Error, nil)
}

// WarnFields logs a message at WarnLevel with fields.
func (g *Global) WarnFields(ctx context.Context, msg string, fields Fields) {
	g.log(ctx, WarnLevel, msg, g.entry.Error, fields)
}

// Error logs a message at ErrorLevel.
func (g *Global) Error(ctx context.Context, msg string) {
	g.log(ctx, ErrorLevel, msg, g.entry.Error, nil)
}

// ErrorCause logs a message at ErrorLevel with cause.
func (g *Global) ErrorCause(ctx context.Context, msg string, cause error) {
	g.log(ctx, ErrorLevel, msg, cause, nil)
}

// ErrorFields logs a message at ErrorLevel with fields.
func (g *Global) ErrorFields(ctx context.Context, msg string, fields Fields) {
	g.log(ctx, ErrorLevel, msg, g.entry.Error, fields)
}

// ErrorCauseFields logs a message at ErrorLevel with cause and fields.
func (g *Global) ErrorCauseFields(ctx context.Context, msg string, cause error, fields Fields) {
	g.log(ctx, ErrorLevel, msg, cause, fields)
}

// With creates a new child logger with additional field.
func (g *Global) With(key string, value interface{}) *Global {
	newEntry := g.entry.With(Field{Key: key, Value: value})

	return &Global{
		entry:       newEntry,
		rootAdapter: g.adapterValue(),
	}
}

// WithFields creates a new logger with additional fields.
func (g *Global) WithFields(fields Fields) *Global {
	newEntry := g.entry.WithFields(fields)

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

// WithSkippedCallerFrame creates a new child logger with one more skipped caller frame. This function is handy when you
// want to write your own logging helpers.
func (g *Global) WithSkippedCallerFrame() *Global {
	newEntry := g.entry
	newEntry.SkippedCallerFrames++

	return &Global{
		entry:       newEntry,
		rootAdapter: g.adapterValue(),
	}
}
