package main

import (
	"context"

	"github.com/elgopher/yala/adapter/zapadapter"
	"github.com/elgopher/yala/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type key string

const contextLoggerKey key = "contextLogger"

// This advanced example shows how to make use of zap logger passed in the context.Context
func main() {
	ctx := context.Background()

	defaultZapLogger := newZapLogger()
	// this adapter will look for zap logger in the context and will wrap it with zapadapter.Adapter
	adapter := ZapContextAdapter{DefaultZapLogger: defaultZapLogger}
	// create logger
	log := logger.WithAdapter(adapter)

	contextLogger := defaultZapLogger.With(zap.String("tag", "value"))
	// bind zap logger to ctx, so all messages will be logged with tag
	ctx = context.WithValue(ctx, contextLoggerKey, contextLogger)

	log.Debug(ctx, "Hello zap from context.Context") // Hello zap from context.Context   {"tag": "value"}
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

type ZapContextAdapter struct {
	DefaultZapLogger *zap.Logger
}

func (a ZapContextAdapter) Log(ctx context.Context, entry logger.Entry) {
	adapter := zapadapter.Adapter{Logger: a.DefaultZapLogger}

	contextLogger := ctx.Value(contextLoggerKey)
	if zapLogger, ok := contextLogger.(*zap.Logger); ok {
		adapter = zapadapter.Adapter{Logger: zapLogger}
	}

	entry.SkippedCallerFrames++ // each middleware adapter must additionally skip one frame

	adapter.Log(ctx, entry)
}
