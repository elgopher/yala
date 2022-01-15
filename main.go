package main

import (
	"context"
	"errors"

	"github.com/jacekolszak/yala/adapter/printer"
	"github.com/jacekolszak/yala/logger"
)

var ErrSome = errors.New("some error")

func main() {
	ctx := context.Background()

	// avoid using global logger
	stdoutService := printer.StdoutService()
	localLogger := logger.Local(stdoutService)
	localLogger.Debug(ctx, "message from local logger")
	localLogger.WithError(ctx, ErrSome).Error("another error")
}
