// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package logger provides tiny structured logging abstraction or facade for various logging libraries, allowing the end
// user to plug in the desired logging library in main.go.
package logger

import (
	"context"
)

// Logger is an immutable logger to log messages or create new loggers with fields or error.
//
// You can't update the adapter once created.
//
// It is safe to use it concurrently.
type Logger struct {
	adapter Adapter
	entry   Entry
}

// WithAdapter creates a new Logger.
func WithAdapter(adapter Adapter) Logger {
	return Logger{adapter: adapter}
}

// Debug logs a message at DebugLevel.
func (l Logger) Debug(ctx context.Context, msg string) {
	l.log(ctx, DebugLevel, msg)
}

func (l Logger) log(ctx context.Context, lvl Level, msg string) {
	if l.adapter == nil {
		return
	}

	e := l.entry
	e.Level = lvl
	e.Message = msg
	e.SkippedCallerFrames += 2

	l.adapter.Log(ctx, e)
}

// Info logs a message at InfoLevel.
func (l Logger) Info(ctx context.Context, msg string) {
	l.log(ctx, InfoLevel, msg)
}

// Warn logs a message at WarnLevel.
func (l Logger) Warn(ctx context.Context, msg string) {
	l.log(ctx, WarnLevel, msg)
}

// Error logs a message at ErrorLevel.
func (l Logger) Error(ctx context.Context, msg string) {
	l.log(ctx, ErrorLevel, msg)
}

// With creates a new logger with field.
func (l Logger) With(key string, value interface{}) Logger {
	l.entry = l.entry.With(Field{key, value})

	return l
}

// WithError creates a new logger with error.
func (l Logger) WithError(err error) Logger {
	l.entry.Error = err

	return l
}

// WithSkippedCallerFrame creates a new logger with one more skipped caller frame. This function is handy when you
// want to write your own logging helpers.
func (l Logger) WithSkippedCallerFrame() Logger {
	l.entry.SkippedCallerFrames++

	return l
}
