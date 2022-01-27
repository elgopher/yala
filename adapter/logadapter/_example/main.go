package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/elgopher/yala/adapter/logadapter"
	"github.com/elgopher/yala/logger"
)

var ErrSome = errors.New("ErrSome")

// This example shows how to use yala with standard `log` package
func main() {
	ctx := context.Background()

	// log using standard log package
	standardLog := log.New(os.Stdout, "", log.LstdFlags)
	adapter := logadapter.Adapter(standardLog)
	yalaLogger := logger.Local{Adapter: adapter}

	yalaLogger.Debug(ctx, "Hello standard log")
	yalaLogger.With(ctx, "f1", "v1").With("f2", "f2").Info("Some info")
	yalaLogger.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	yalaLogger.WithError(ctx, ErrSome).Error("Some error")
}