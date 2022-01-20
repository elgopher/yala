package contextadapter

import (
	"context"

	"github.com/jacekolszak/yala/logger"
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

	newEntry := entry
	newEntry.SkippedCallerFrames++

	loggerAdapter.Log(ctx, newEntry)
}
