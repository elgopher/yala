// Package logger provides API for logging messages, to be used by code which is not aware what logging library is used.
package logger

import (
	"context"
	"fmt"
)

// Debug logs message using globally configured logger.Service.
func Debug(ctx context.Context, msg string) {
	globalLoggerWithSkippedCallerFrame(ctx).Debug(msg)
}

func globalLoggerWithSkippedCallerFrame(ctx context.Context) Logger {
	return getGlobalLogger().WithSkippedCallerFrame(ctx)
}

// Debugf logs message using globally configured logger.Service.
func Debugf(ctx context.Context, format string, args ...interface{}) {
	globalLoggerWithSkippedCallerFrame(ctx).Debugf(format, args...)
}

// Info logs message using globally configured logger.Service.
func Info(ctx context.Context, msg string) {
	globalLoggerWithSkippedCallerFrame(ctx).Info(msg)
}

// Infof logs message using globally configured logger.Service.
func Infof(ctx context.Context, format string, args ...interface{}) {
	globalLoggerWithSkippedCallerFrame(ctx).Infof(format, args...)
}

// Warn logs message using globally configured logger.Service.
func Warn(ctx context.Context, msg string) {
	globalLoggerWithSkippedCallerFrame(ctx).Warn(msg)
}

// Warnf logs message using globally configured logger.Service.
func Warnf(ctx context.Context, format string, args ...interface{}) {
	globalLoggerWithSkippedCallerFrame(ctx).Warnf(format, args...)
}

// Error logs message using globally configured logger.Service.
func Error(ctx context.Context, msg string) {
	globalLoggerWithSkippedCallerFrame(ctx).Error(msg)
}

// Errorf logs message using globally configured logger.Service.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	globalLoggerWithSkippedCallerFrame(ctx).Errorf(format, args...)
}

// With creates a new Logger with field and using globally configured logger.Service.
func With(ctx context.Context, key string, value interface{}) Logger {
	return getGlobalLogger().With(ctx, key, value)
}

// WithError creates a new Logger with error and using globally configured logger.Service.
func WithError(ctx context.Context, err error) Logger {
	return getGlobalLogger().WithError(ctx, err)
}

type Logger struct {
	entry   Entry
	service Service
	ctx     context.Context
}

// With creates a new Logger with field.
func (l Logger) With(key string, value interface{}) Logger {
	c := l
	c.entry.Fields = make([]Field, len(c.entry.Fields), cap(c.entry.Fields))
	copy(c.entry.Fields, l.entry.Fields)
	c.entry.Fields = append(c.entry.Fields, Field{key, value})

	return c
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

func (l Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l Logger) Info(msg string) {
	l.log(InfoLevel, msg)
}

func (l Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l Logger) Warn(msg string) {
	l.log(WarnLevel, msg)
}

func (l Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

func (l Logger) Error(msg string) {
	l.log(ErrorLevel, msg)
}

func (l Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l Logger) log(level Level, msg string) {
	e := l.entry
	e.Level = level
	e.Message = msg
	e.SkippedCallerFrames += 3
	l.service.Log(l.ctx, e)
}
