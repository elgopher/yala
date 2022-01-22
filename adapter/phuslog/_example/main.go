package main

import (
	"context"
	"errors"
	"os"

	"github.com/jacekolszak/yala/adapter/phuslog"
	"github.com/jacekolszak/yala/logger"
	"github.com/phuslu/log"
)

var ErrSome = errors.New("ErrSome")

// This example shows how to use yala with phuslog adapter
func main() {
	ctx := context.Background()

	// create phuslog logger
	l := log.Logger{
		Level:  log.DebugLevel,
		Writer: log.IOWriter{Writer: os.Stderr},
	}
	adapter := phuslog.Adapter{Logger: l} // create logger.Adapter for phuslog
	yalaLogger := logger.Local(adapter)   // Create yala logger

	yalaLogger.Debug(ctx, "Hello phuslog")
	yalaLogger.With(ctx, "field_name", "field_value").Info("Some info")
	yalaLogger.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	yalaLogger.WithError(ctx, ErrSome).Error("Some error")
}
