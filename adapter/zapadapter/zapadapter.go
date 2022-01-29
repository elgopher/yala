// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

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

func (a Adapter) Log(_ context.Context, entry logger.Entry) {
	if a.Logger == nil {
		return
	}

	zapLogger := a.Logger

	if entry.Error != nil {
		zapLogger = zapLogger.With(zap.Error(entry.Error))
	}

	fields := make([]zap.Field, len(entry.Fields))

	for i, f := range entry.Fields {
		fields[i] = zap.Any(f.Key, f.Value)
	}

	zapLogger = zapLogger.WithOptions(zap.AddCallerSkip(entry.SkippedCallerFrames + 1))

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
