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
	l.log(ctx, DebugLevel, msg, l.entry.Error, nil)
}

// DebugFields logs a message at DebugLevel with fields.
func (l Logger) DebugFields(ctx context.Context, msg string, fields Fields) {
	l.log(ctx, DebugLevel, msg, l.entry.Error, fields)
}

func (l Logger) log(ctx context.Context, lvl Level, msg string, cause error, fields Fields) {
	if l.adapter == nil {
		return
	}

	newEntry := l.entry.WithFields(fields)
	newEntry.Error = cause
	newEntry.Level = lvl
	newEntry.Message = msg
	newEntry.SkippedCallerFrames += 2

	l.adapter.Log(ctx, newEntry)
}

// Info logs a message at InfoLevel.
func (l Logger) Info(ctx context.Context, msg string) {
	l.log(ctx, InfoLevel, msg, l.entry.Error, nil)
}

// InfoFields logs a message at InfoLevel with fields.
func (l Logger) InfoFields(ctx context.Context, msg string, fields Fields) {
	l.log(ctx, InfoLevel, msg, l.entry.Error, fields)
}

// Warn logs a message at WarnLevel.
func (l Logger) Warn(ctx context.Context, msg string) {
	l.log(ctx, WarnLevel, msg, l.entry.Error, nil)
}

// WarnFields logs a message at WarnLevel with fields.
func (l Logger) WarnFields(ctx context.Context, msg string, fields Fields) {
	l.log(ctx, WarnLevel, msg, l.entry.Error, fields)
}

// Error logs a message at ErrorLevel.
func (l Logger) Error(ctx context.Context, msg string) {
	l.log(ctx, ErrorLevel, msg, l.entry.Error, nil)
}

// ErrorCause logs a message at ErrorLevel with cause.
func (l Logger) ErrorCause(ctx context.Context, msg string, cause error) {
	l.log(ctx, ErrorLevel, msg, cause, nil)
}

// ErrorFields logs a message at ErrorLevel with fields.
func (l Logger) ErrorFields(ctx context.Context, msg string, fields Fields) {
	l.log(ctx, ErrorLevel, msg, l.entry.Error, fields)
}

// ErrorCauseFields logs a message at ErrorLevel with cause and fields.
func (l Logger) ErrorCauseFields(ctx context.Context, msg string, cause error, fields Fields) {
	l.log(ctx, ErrorLevel, msg, cause, fields)
}

// With creates a new logger with additional field.
func (l Logger) With(key string, value interface{}) Logger {
	l.entry = l.entry.With(Field{key, value})

	return l
}

// WithFields creates a new logger with additional fields.
func (l Logger) WithFields(fields Fields) Logger {
	l.entry = l.entry.WithFields(fields)

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
