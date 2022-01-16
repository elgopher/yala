package logger

import (
	"context"
	"fmt"
)

const localLoggerSkippedCallerFrames = 2

type LocalLogger struct {
	service Service
}

func Local(service Service) LocalLogger {
	if service == nil {
		service = noopLogger{}
	}

	return LocalLogger{
		service: service,
	}
}

func (l LocalLogger) Debug(ctx context.Context, msg string) {
	l.service.Log(ctx, Entry{Level: DebugLevel, Message: msg, SkippedCallerFrames: localLoggerSkippedCallerFrames})
}

func (l LocalLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	l.WithSkippedCallerFrame(ctx).Debug(fmt.Sprintf(format, args...))
}

func (l LocalLogger) Info(ctx context.Context, msg string) {
	l.service.Log(ctx, Entry{Level: InfoLevel, Message: msg, SkippedCallerFrames: localLoggerSkippedCallerFrames})
}

func (l LocalLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	l.WithSkippedCallerFrame(ctx).Info(fmt.Sprintf(format, args...))
}

func (l LocalLogger) Error(ctx context.Context, msg string) {
	l.service.Log(ctx, Entry{Level: ErrorLevel, Message: msg, SkippedCallerFrames: localLoggerSkippedCallerFrames})
}

func (l LocalLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.WithSkippedCallerFrame(ctx).Error(fmt.Sprintf(format, args...))
}

// With creates a new Logger with field.
func (l LocalLogger) With(ctx context.Context, key string, value interface{}) Logger {
	return l.fromContext(ctx).With(key, value)
}

func (l LocalLogger) fromContext(ctx context.Context) Logger {
	return Logger{
		service: l.service,
		ctx:     ctx,
	}
}

// WithError creates a new Logger with error.
func (l LocalLogger) WithError(ctx context.Context, err error) Logger {
	return l.fromContext(ctx).WithError(err)
}

func (l LocalLogger) WithSkippedCallerFrame(ctx context.Context) Logger {
	return l.fromContext(ctx).WithSkippedCallerFrame()
}
