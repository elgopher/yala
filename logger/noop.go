package logger

import (
	"context"
	"fmt"
	"sync"
)

type noopLogger struct{}

func (n noopLogger) Log(context.Context, Entry) {}

var once sync.Once

type initialGlobalNoopLogger struct{}

func (g *initialGlobalNoopLogger) Log(_ context.Context, entry Entry) {
	if entry.Level == WarnLevel || entry.Level == ErrorLevel {
		once.Do(func() {
			fmt.Printf("github.com/jacekolszak/yala/logger cannot log message with level %s. Please set the global logging adapter. For example: logger.SetAdapter(printer.StdoutAdapter()) to log all messages to stdout.\n", entry.Level) // nolint
		})
	}
}
