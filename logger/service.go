package logger

import (
	"context"
	"sync/atomic"
)

// SetService sets a global service implementation used by logging functions in the logger package,
// such as `logger.Info`. By default, nothing is logged.
func SetService(service Service) {
	if service == nil {
		service = noopLogger{}
	}

	globalLogger.Store(
		LocalLogger{service: service},
	)
}

type Service interface {
	Log(context.Context, Entry)
}

type Entry struct {
	Level   Level
	Message string
	Fields  []Field
	Error   error
	// SkippedCallerFrames can be used by logger.Service to caller information (file and line number)
	SkippedCallerFrames int
}

type Level string

const (
	DebugLevel Level = "DEBUG"
	InfoLevel  Level = "INFO"
	ErrorLevel Level = "ERROR"
)

type Field struct {
	Key   string
	Value interface{}
}

func init() {
	SetService(noopLogger{})
}

var globalLogger atomic.Value

func getGlobalLogger() LocalLogger {
	return globalLogger.Load().(LocalLogger)
}
