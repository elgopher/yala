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

	for _, f := range entry.Fields {
		zapLogger = zapLogger.With(zap.Any(f.Key, f.Value))
	}

	zapLogger = zapLogger.WithOptions(zap.AddCallerSkip(entry.SkippedCallerFrames))

	switch entry.Level {
	case logger.DebugLevel:
		zapLogger.Debug(entry.Message)
	case logger.InfoLevel:
		zapLogger.Info(entry.Message)
	case logger.WarnLevel:
		zapLogger.Warn(entry.Message)
	case logger.ErrorLevel:
		zapLogger.Error(entry.Message)
	default:
		zapLogger.Info(entry.Message)
	}
}
