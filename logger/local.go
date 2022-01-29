// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logger

import (
	"context"
)

const localLoggerSkippedCallerFrames = 2

// Local is an immutable struct to log messages or create new loggers with fields or error.
//
// It is safe to use it concurrently.
type Local struct {
	Adapter Adapter
}

// Debug logs a message at DebugLevel.
func (l Local) Debug(ctx context.Context, msg string) {
	l.log(ctx, DebugLevel, msg)
}

func (l Local) log(ctx context.Context, lvl Level, msg string) {
	if l.Adapter == nil {
		return
	}

	l.Adapter.Log(ctx, Entry{Level: lvl, Message: msg, SkippedCallerFrames: localLoggerSkippedCallerFrames})
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

// With creates a new Logger with field.
func (l Local) With(key string, value interface{}) Logger {
	return l.logger().With(key, value)
}

func (l Local) logger() Logger {
	return Logger{
		adapter: l.Adapter,
	}
}

// WithError creates a new Logger with error.
func (l Local) WithError(err error) Logger {
	return l.logger().WithError(err)
}

// WithSkippedCallerFrame creates a new Logger with one more skipped caller frame. This function is handy when you
// want to write your own logging helpers.
func (l Local) WithSkippedCallerFrame() Logger {
	return l.logger().WithSkippedCallerFrame()
}
