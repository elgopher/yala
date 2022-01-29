// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package logger provides structured logging abstraction or facade, to be used by code which is not aware what logging
// library is used by end user.
package logger

import (
	"context"
)

// Logger is an immutable struct to log messages or create new loggers with fields or error.
//
// It is safe to use it concurrently.
type Logger struct {
	entry   Entry
	adapter Adapter
}

// With creates a new Logger with field.
func (l Logger) With(key string, value interface{}) Logger {
	l.entry = l.entry.With(Field{key, value})

	return l
}

// WithError creates a new Logger with error.
func (l Logger) WithError(err error) Logger {
	l.entry.Error = err

	return l
}

// WithSkippedCallerFrame creates a new Logger with one more skipped caller frame. This function is handy when you
// want to write your own logging helpers.
func (l Logger) WithSkippedCallerFrame() Logger {
	l.entry.SkippedCallerFrames++

	return l
}

// Debug logs a message at DebugLevel.
func (l Logger) Debug(ctx context.Context, msg string) {
	l.log(ctx, DebugLevel, msg)
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

func (l Logger) log(ctx context.Context, level Level, msg string) {
	e := l.entry
	e.Level = level
	e.Message = msg
	e.SkippedCallerFrames += 2

	l.adapter.Log(ctx, e)
}
