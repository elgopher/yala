package main

import (
	"context"
	"errors"

	"github.com/jacekolszak/yala/adapter/zapadapter"
	"github.com/jacekolszak/yala/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()

	zapLogger := newZapLogger()
	service := zapadapter.Service{Logger: zapLogger} // create logger.Service for zap
	logger.SetService(service)                       // set it globally

	logger.Debug(ctx, "Hello zap")
	logger.With(ctx, "tag", "bbb").Info("Some info")
	logger.Warnf(ctx, "Be careful with %s", "hot water")
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
