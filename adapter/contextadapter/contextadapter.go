package contextadapter

import (
	"context"

	"github.com/jacekolszak/yala/logger"
)

func New(contextKey interface{}, a logger.Adapter) logger.Adapter { // nolint
	return adapter{
		adapter:    a,
		contextKey: contextKey,
	}
}

type adapter struct {
	adapter    logger.Adapter
	contextKey interface{}
}

func (p adapter) Log(ctx context.Context, entry logger.Entry) {
	loggerAdapter, ok := ctx.Value(p.contextKey).(logger.Adapter)
	if !ok {
		loggerAdapter = p.adapter
	}

	newEntry := entry
	newEntry.SkippedCallerFrames++

	loggerAdapter.Log(ctx, newEntry)
}
