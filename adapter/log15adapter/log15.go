// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package log15adapter

import (
	"context"

	"github.com/elgopher/yala/logger"
	"github.com/inconshreveable/log15"
)

// Adapter is a logger.Adapter implementation, which is using `log15` package
// (https://github.com/inconshreveable/log15).
type Adapter struct {
	Logger log15.Logger
}

// Log logs the entry using log15 package.
func (a Adapter) Log(ctx context.Context, entry logger.Entry) {
	if a.Logger == nil {
		return
	}

	log15ctx := toCtx(entry)

	switch entry.Level {
	case logger.DebugLevel:
		a.Logger.Debug(entry.Message, log15ctx...)
	case logger.InfoLevel:
		a.Logger.Info(entry.Message, log15ctx...)
	case logger.WarnLevel:
		a.Logger.Warn(entry.Message, log15ctx...)
	case logger.ErrorLevel:
		a.Logger.Error(entry.Message, log15ctx...)
	default:
		a.Logger.Info(entry.Message, log15ctx...)
	}
}

func toCtx(entry logger.Entry) []interface{} {
	entryError := entry.Error

	const lengthOfField = 2

	length := len(entry.Fields) * lengthOfField
	if entryError != nil {
		length += lengthOfField
	}

	log15ctx := make([]interface{}, length)

	for i, field := range entry.Fields {
		log15ctx[i*lengthOfField] = field.Key
		log15ctx[i*lengthOfField+1] = field.Value
	}

	if entryError != nil {
		log15ctx[length-lengthOfField] = "error"
		log15ctx[length-1] = entryError
	}

	return log15ctx
}
