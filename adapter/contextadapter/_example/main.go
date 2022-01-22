package main

import (
	"context"
	"errors"

	"github.com/jacekolszak/yala/adapter/contextadapter"
	"github.com/jacekolszak/yala/adapter/zapadapter"
	"github.com/jacekolszak/yala/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type key string

const contextLoggerKey key = "contextLogger"

var ErrSome = errors.New("ErrSome")

// This advanced example shows how to make use of zap logger passed in the context.Context
func main() {
	ctx := context.Background()

	defaultZapLogger := newZapLogger()

	// this adapter will look for zap logger in the context and will wrap it with zapadapter.Adapter
	adapter := contextadapter.New(contextLoggerKey, func(loggerOrNil interface{}) logger.Adapter {
		if zapLogger, ok := loggerOrNil.(*zap.Logger); ok {
			return zapadapter.Adapter{Logger: zapLogger}
		}

		return zapadapter.Adapter{Logger: defaultZapLogger}
	})
	// create local logger
	yalaLogger := logger.Local(adapter)

	contextLogger := defaultZapLogger.With(zap.String("tag", "value"))
	// bind zap logger to ctx, so all messages will be logged with tag
	ctx = context.WithValue(ctx, contextLoggerKey, contextLogger)

	yalaLogger.Debug(ctx, "Hello zap from ctx")
	yalaLogger.With(ctx, "field_name", "field_value").Info("Some info")
	yalaLogger.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	yalaLogger.WithError(ctx, ErrSome).Error("Some error")
}

func newZapLogger() *zap.Logger {
	zapCfg := zap.NewDevelopmentConfig()
	zapCfg.DisableStacktrace = true
	zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLogger, err := zapCfg.Build()
	if err != nil {
		panic(err)
	}
	return zapLogger
}
