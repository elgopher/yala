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
func (l Logger) Debug(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, DebugLevel, msg, fields)
}

func (l Logger) log(ctx context.Context, lvl Level, msg string, fields []Field) {
	if l.adapter == nil {
		return
	}

	newEntry := l.entry
	newEntry.Level = lvl
	newEntry.Message = msg
	newEntry.SkippedCallerFrames += 2
	newEntry.Fields = mergeFields(newEntry.Fields, fields)

	l.adapter.Log(ctx, newEntry)
}

func mergeFields(fields []Field, newFields []Field) []Field {
	if len(fields) == 0 && len(newFields) == 0 {
		return nil
	}

	le := len(fields)
	merged := make([]Field, le+len(newFields))
	copy(merged, fields)
	copy(merged[le:], newFields)

	return merged
}

// Info logs a message at InfoLevel.
func (l Logger) Info(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, InfoLevel, msg, fields)
}

// Warn logs a message at WarnLevel.
func (l Logger) Warn(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, WarnLevel, msg, fields)
}

// Error logs a message at ErrorLevel.
func (l Logger) Error(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, ErrorLevel, msg, fields)
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
