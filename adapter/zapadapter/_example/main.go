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
	log := logger.WithAdapter(adapter)               // Create yala logger

	log.Debug(ctx, "Hello zap")
	log.With("field_name", "field_value").Info(ctx, "Some info")
	log.With("parameter", "some").Warn(ctx, "Deprecated configuration parameter. It will be removed.")
	log.WithError(ErrSome).Error(ctx, "Some error")
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
