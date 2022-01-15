package contextservice

import (
	"context"

	"github.com/jacekolszak/yala/logger"
)

func New(contextKey interface{}, s logger.Service) logger.Service { // nolint
	return service{
		service:    s,
		contextKey: contextKey,
	}
}

type service struct {
	service    logger.Service
	contextKey interface{}
}

func (p service) Log(ctx context.Context, entry logger.Entry) {
	loggerService, ok := ctx.Value(p.contextKey).(logger.Service)
	if !ok {
		loggerService = p.service
	}

	newEntry := entry
	newEntry.SkippedCallerFrames++

	loggerService.Log(ctx, newEntry)
}
