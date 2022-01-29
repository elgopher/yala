// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package logger provides structured logging abstraction or facade, to be used by code which is not aware what logging
// library is used by end user.
//
// Each message logged has a level, which was modeled after http://tools.ietf.org/html/rfc5424 severity levels:
//
// 	* Debug - Information useful to developers for debugging the application.
// 	* Info  - Normal operational messages that require no action.
// 	* Warn  - May indicate that an error will occur if action is not taken.
// 	* Error - Non-urgent failures - these should be relayed to developers or admins; each item must be resolved within
//	          a given time.
//
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

func (l Logger) WithSkippedCallerFrame() Logger {
	l.entry.SkippedCallerFrames++

	return l
}

func (l Logger) Debug(ctx context.Context, msg string) {
	l.log(ctx, DebugLevel, msg)
}

func (l Logger) Info(ctx context.Context, msg string) {
	l.log(ctx, InfoLevel, msg)
}

func (l Logger) Warn(ctx context.Context, msg string) {
	l.log(ctx, WarnLevel, msg)
}

func (l Logger) Error(ctx context.Context, msg string) {
	l.log(ctx, ErrorLevel, msg)
}

func (l Logger) log(ctx context.Context, level Level, msg string) {
	e := l.entry
	e.Level = level
	e.Message = msg
	e.SkippedCallerFrames += 3

	l.adapter.Log(ctx, e)
}
