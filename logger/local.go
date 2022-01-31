// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logger

import (
	"context"
)

// Local is an immutable struct to log messages or create new loggers with fields or error.
//
// It is safe to use it concurrently.
type Local struct {
	Adapter Adapter
	entry   Entry
}

// Debug logs a message at DebugLevel.
func (l Local) Debug(ctx context.Context, msg string) {
	l.log(ctx, DebugLevel, msg)
}

func (l Local) log(ctx context.Context, lvl Level, msg string) {
	if l.Adapter == nil {
		return
	}

	e := l.entry
	e.Level = lvl
	e.Message = msg
	e.SkippedCallerFrames += 2

	l.Adapter.Log(ctx, e)
}

// Info logs a message at InfoLevel.
func (l Local) Info(ctx context.Context, msg string) {
	l.log(ctx, InfoLevel, msg)
}

// Warn logs a message at WarnLevel.
func (l Local) Warn(ctx context.Context, msg string) {
	l.log(ctx, WarnLevel, msg)
}

// Error logs a message at ErrorLevel.
func (l Local) Error(ctx context.Context, msg string) {
	l.log(ctx, ErrorLevel, msg)
}

// With creates a new logger with field.
func (l Local) With(key string, value interface{}) Local {
	l.entry = l.entry.With(Field{key, value})

	return l
}

// WithError creates a new logger with error.
func (l Local) WithError(err error) Local {
	l.entry.Error = err

	return l
}

// WithSkippedCallerFrame creates a new logger with one more skipped caller frame. This function is handy when you
// want to write your own logging helpers.
func (l Local) WithSkippedCallerFrame() Local {
	l.entry.SkippedCallerFrames++

	return l
}
