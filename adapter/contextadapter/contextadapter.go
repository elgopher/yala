// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package contextadapter provides logger.Adapter implementation which lookups logger in the context.Context and wrap
// it inside adapter. It is useful when your app is using context.Context to store logger instance.
package contextadapter

import (
	"context"

	"github.com/elgopher/yala/logger"
)

func New(contextKey interface{}, adapterFromContextLogger func(loggerOrNil interface{}) logger.Adapter) logger.Adapter { // nolint
	return adapter{
		contextKey:               contextKey,
		adapterFromContextLogger: adapterFromContextLogger,
	}
}

type adapter struct {
	contextKey               interface{}
	adapterFromContextLogger func(loggerOrNil interface{}) logger.Adapter
}

func (p adapter) Log(ctx context.Context, entry logger.Entry) {
	contextLogger := ctx.Value(p.contextKey)
	loggerAdapter := p.adapterFromContextLogger(contextLogger)

	entry.SkippedCallerFrames++

	loggerAdapter.Log(ctx, entry)
}
