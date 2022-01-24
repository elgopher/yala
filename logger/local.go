// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logger

import (
	"context"
)

const localLoggerSkippedCallerFrames = 3

// LocalLogger is an immutable struct to log messages or create new loggers with fields or error.
//
// It is safe to use it concurrently.
type LocalLogger struct {
	adapter Adapter
}

// Local creates a new LocalLogger.
func Local(adapter Adapter) LocalLogger {
	if adapter == nil {
		adapter = noopAdapter{}
	}

	return LocalLogger{
		adapter: adapter,
	}
}

func (l LocalLogger) Debug(ctx context.Context, msg string) {
	l.log(ctx, DebugLevel, msg)
}

func (l LocalLogger) log(ctx context.Context, lvl Level, msg string) {
	if l.adapter == nil {
		return
	}

	l.adapter.Log(ctx, Entry{Level: lvl, Message: msg, SkippedCallerFrames: localLoggerSkippedCallerFrames})
}

func (l LocalLogger) Info(ctx context.Context, msg string) {
	l.log(ctx, InfoLevel, msg)
}

func (l LocalLogger) Warn(ctx context.Context, msg string) {
	l.log(ctx, WarnLevel, msg)
}

func (l LocalLogger) Error(ctx context.Context, msg string) {
	l.log(ctx, ErrorLevel, msg)
}

// With creates a new Logger with field.
func (l LocalLogger) With(ctx context.Context, key string, value interface{}) Logger {
	return l.logger(ctx).With(key, value)
}

func (l LocalLogger) logger(ctx context.Context) Logger {
	return Logger{
		adapter: l.adapter,
		ctx:     ctx,
	}
}

// WithError creates a new Logger with error.
func (l LocalLogger) WithError(ctx context.Context, err error) Logger {
	return l.logger(ctx).WithError(err)
}

func (l LocalLogger) WithSkippedCallerFrame(ctx context.Context) Logger {
	return l.logger(ctx).WithSkippedCallerFrame()
}
