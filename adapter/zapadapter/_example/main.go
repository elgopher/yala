package main

import (
	"context"
	"errors"

	"github.com/jacekolszak/yala/adapter/zapadapter"
	"github.com/jacekolszak/yala/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ErrSome = errors.New("ErrSome")

// This example shows how to use yala with zap adapter
func main() {
	ctx := context.Background()

	zapLogger := newZapLogger()
	adapter := zapadapter.Adapter{Logger: zapLogger} // create logger.Adapter for zap
	logger.SetAdapter(adapter)                       // set it globally

	logger.Debug(ctx, "Hello zap")
	logger.With(ctx, "field_name", "field_value").Info("Some info")
	logger.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	logger.WithError(ctx, ErrSome).Error("Some error")
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
