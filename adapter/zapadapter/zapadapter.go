package zapadapter

import (
	"context"

	"github.com/jacekolszak/yala/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Service is a logger.Service implementation, which is using `zap` module (https://github.com/uber-go/zap).
type Service struct {
	*zap.Logger
}

func (s Service) Log(_ context.Context, entry logger.Entry) {
	if s.Logger == nil {
		return
	}

	zapLogger := s.Logger

	if entry.Error != nil {
		zapLogger = zapLogger.With(zap.Error(entry.Error))
	}

	for _, f := range entry.Fields {
		zapLogger = zapLogger.With(
			zap.Field{
				Key:       f.Key,
				Interface: f.Value,
				Type:      zapcore.StringType,
			},
		)
	}

	zapLogger = zapLogger.WithOptions(zap.AddCallerSkip(entry.SkippedCallerFrames))

	switch entry.Level {
	case logger.DebugLevel:
		zapLogger.Debug(entry.Message)
	case logger.InfoLevel:
		zapLogger.Info(entry.Message)
	case logger.ErrorLevel:
		zapLogger.Error(entry.Message)
	default:
		zapLogger.Info(entry.Message)
	}
}
