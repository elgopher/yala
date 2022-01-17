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

const contextAdapterKey key = "contextAdapter"

// This example shows how to pass zap logger in the context.Context
func main() {
	ctx := context.Background()

	zapLogger := newZapLogger()
	defaultContextAdapter := zapadapter.Adapter{Logger: zapLogger}
	adapter := contextadapter.New(contextAdapterKey, defaultContextAdapter) // create logger.Adapter implementation which will look for logger.Adapter in the context
	logger.SetAdapter(adapter)                                              // set it globally

	contextLogger := zapLogger.With(zap.String("request_field", "value")) // create zap logger which will be bound to context.Context
	contextAdapter := zapadapter.Adapter{Logger: contextLogger}
	ctx = context.WithValue(ctx, contextAdapterKey, contextAdapter)

	logger.Debug(ctx, "Hello zap from ctx")
	logger.With(ctx, "field_name", "field_value").Info("Some info")
	logger.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	logger.WithError(ctx, errors.New("ss")).Error("Some error")
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
