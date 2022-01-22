package logger

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

type noopLogger struct{}

func (n noopLogger) Log(context.Context, Entry) {}

type initialGlobalNoopLogger struct {
	once sync.Once
}

func (g *initialGlobalNoopLogger) Log(_ context.Context, entry Entry) {
	if entry.Level == WarnLevel || entry.Level == ErrorLevel {
		g.once.Do(func() {
			const framesToSkip = 7
			_, file, line, _ := runtime.Caller(framesToSkip)
			fmt.Printf("%s:%d cannot log message with level %s. Please configure the global logger.\n", file, line, entry.Level) // nolint
		})
	}
}
