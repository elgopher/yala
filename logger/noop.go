// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logger

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

type noopAdapter struct{}

func (n noopAdapter) Log(context.Context, Entry) {}

type initialGlobalNoopAdapter struct {
	once sync.Once
}

func (g *initialGlobalNoopAdapter) Log(_ context.Context, entry Entry) {
	if entry.Level == WarnLevel || entry.Level == ErrorLevel {
		g.once.Do(func() {
			const framesToSkip = 4
			_, file, line, _ := runtime.Caller(entry.SkippedCallerFrames + framesToSkip)
			fmt.Printf("%s:%d cannot log message with level %s. Please configure the global logger.\n", file, line, entry.Level) // nolint
		})
	}
}
