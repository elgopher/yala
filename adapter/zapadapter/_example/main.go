package main

import (
	"context"
	"errors"

	"github.com/elgopher/yala/adapter/zapadapter"
	"github.com/elgopher/yala/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ErrSome = errors.New("ErrSome")

// This example shows how to use yala with zap adapter
func main() {
	ctx := context.Background()

	zapLogger := newZapLogger()
	adapter := zapadapter.Adapter{Logger: zapLogger} // create logger.Adapter for zap
	yalaLogger := logger.Local(adapter)              // Create yala logger

	yalaLogger.Debug(ctx, "Hello zap")
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
