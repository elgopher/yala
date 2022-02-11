// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package zapadapter provides yala adapter which leverages zap module (https://github.com/uber-go/zap).
package zapadapter

import (
	"context"

	"github.com/elgopher/yala/logger"
	"go.uber.org/zap"
)

// Adapter is a logger.Adapter implementation, which is using `zap` module (https://github.com/uber-go/zap).
type Adapter struct {
	Logger *zap.Logger
}

// Log logs the entry using zap module.
func (a Adapter) Log(_ context.Context, entry logger.Entry) {
	if a.Logger == nil {
		return
	}

	zapLogger := a.Logger.WithOptions(zap.AddCallerSkip(entry.SkippedCallerFrames + 1))
	fields := zapFields(entry)

	switch entry.Level {
	case logger.DebugLevel:
		zapLogger.Debug(entry.Message, fields...)
	case logger.InfoLevel:
		zapLogger.Info(entry.Message, fields...)
	case logger.WarnLevel:
		zapLogger.Warn(entry.Message, fields...)
	case logger.ErrorLevel:
		zapLogger.Error(entry.Message, fields...)
	default:
		zapLogger.Info(entry.Message, fields...)
	}
}

func zapFields(entry logger.Entry) []zap.Field {
	length := len(entry.Fields)
	if entry.Error != nil {
		length++
	}

	fields := make([]zap.Field, length)

	for i, f := range entry.Fields {
		fields[i] = zap.Any(f.Key, f.Value)
	}

	if entry.Error != nil {
		fields[length-1] = zap.Error(entry.Error)
	}

	return fields
}
