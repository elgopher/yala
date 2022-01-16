package main

import (
	"context"
	"errors"

	"github.com/jacekolszak/yala/adapter/contextservice"
	"github.com/jacekolszak/yala/adapter/zapadapter"
	"github.com/jacekolszak/yala/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type key string

const contextServiceKey key = "contextService"

// This example shows how to pass zap logger in the context.Context
func main() {
	ctx := context.Background()

	zapLogger := newZapLogger()
	defaultContextService := zapadapter.Service{Logger: zapLogger}
	service := contextservice.New(contextServiceKey, defaultContextService) // create logger.Service implementation which will look for logger.Service in the context
	logger.SetService(service)                                              // set it globally

	contextLogger := zapLogger.With(zap.String("request_field", "value")) // create zap logger which will be bound to context.Context
	contextService := zapadapter.Service{Logger: contextLogger}
	ctx = context.WithValue(ctx, contextServiceKey, contextService)

	logger.Debug(ctx, "Hello zap from ctx")
	logger.With(ctx, "tag", "bbb").Info("Some info")
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
