// Package logger provides API for logging messages, to be used by code which is not aware what logging library is used.
package logger

import (
	"context"
)

// Debug logs message using globally configured logger.Adapter.
func Debug(ctx context.Context, msg string) {
	globalLoggerWithSkippedCallerFrame(ctx).Debug(msg)
}

func globalLoggerWithSkippedCallerFrame(ctx context.Context) Logger {
	return getGlobalLogger().WithSkippedCallerFrame(ctx)
}

// Info logs message using globally configured logger.Adapter.
func Info(ctx context.Context, msg string) {
	globalLoggerWithSkippedCallerFrame(ctx).Info(msg)
}

// Warn logs message using globally configured logger.Adapter.
func Warn(ctx context.Context, msg string) {
	globalLoggerWithSkippedCallerFrame(ctx).Warn(msg)
}

// Error logs message using globally configured logger.Adapter.
func Error(ctx context.Context, msg string) {
	globalLoggerWithSkippedCallerFrame(ctx).Error(msg)
}

// With creates a new Logger with field and using globally configured logger.Adapter.
func With(ctx context.Context, key string, value interface{}) Logger {
	return getGlobalLogger().With(ctx, key, value)
}

// WithError creates a new Logger with error and using globally configured logger.Adapter.
func WithError(ctx context.Context, err error) Logger {
	return getGlobalLogger().WithError(ctx, err)
}

type Logger struct {
	entry   Entry
	adapter Adapter
	ctx     context.Context
}

// With creates a new Logger with field.
func (l Logger) With(key string, value interface{}) Logger {
	loggerCopy := l

	newLen := len(l.entry.Fields) + 1
	loggerCopy.entry.Fields = make([]Field, newLen)
	copy(loggerCopy.entry.Fields, l.entry.Fields)
	loggerCopy.entry.Fields[newLen-1] = Field{key, value}

	return loggerCopy
}

// WithError creates a new Logger with error.
func (l Logger) WithError(err error) Logger {
	c := l
	c.entry.Error = err

	return c
}

func (l Logger) WithSkippedCallerFrame() Logger {
	c := l
	c.entry.SkippedCallerFrames++

	return c
}

func (l Logger) Debug(msg string) {
	l.log(DebugLevel, msg)
}

func (l Logger) Info(msg string) {
	l.log(InfoLevel, msg)
}

func (l Logger) Warn(msg string) {
	l.log(WarnLevel, msg)
}

func (l Logger) Error(msg string) {
	l.log(ErrorLevel, msg)
}

func (l Logger) log(level Level, msg string) {
	e := l.entry
	e.Level = level
	e.Message = msg
	e.SkippedCallerFrames += 3

	l.adapter.Log(l.ctx, e)
}
