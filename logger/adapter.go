package logger

import (
	"context"
	"sync/atomic"
)

// SetAdapter sets a global adapter implementation used by logging functions in the logger package,
// such as `logger.Info`. By default, nothing is logged.
func SetAdapter(adapter Adapter) {
	if adapter == nil {
		adapter = noopLogger{}
	}

	globalLogger.Store(
		LocalLogger{adapter: adapter},
	)
}

// Adapter is an interface to be implemented by logger adapters.
type Adapter interface {
	Log(context.Context, Entry)
}

type Entry struct {
	Level   Level
	Message string
	Fields  []Field // Fields can be nil
	Error   error   // Error can be nil
	// SkippedCallerFrames can be used by logger.Adapter to extract caller information (file and line number)
	SkippedCallerFrames int
}

type Level string

const (
	DebugLevel Level = "DEBUG"
	InfoLevel  Level = "INFO"
	WarnLevel  Level = "WARN"
	ErrorLevel Level = "ERROR"
)

type Field struct {
	Key   string
	Value interface{}
}

func init() {
	SetAdapter(noopLogger{})
}

var globalLogger atomic.Value

func getGlobalLogger() LocalLogger {
	return globalLogger.Load().(LocalLogger)
}
