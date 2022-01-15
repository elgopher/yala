package logger

import (
	"context"
)

type noopLogger struct{}

func (n noopLogger) Log(context.Context, Entry) {}
