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
	standardLog := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	adapter := logadapter.Adapter(standardLog)
	log := logger.WithAdapter(adapter)

	log.Debug(ctx, "Hello standard log")

	log.Info(ctx,
		"Some info",
		"f1", "v1",
		"f2", "f2",
	)

	log.Warn(ctx,
		"Deprecated configuration parameter. It will be removed.",
		"parameter", "some",
	)

	log.WithError(ErrSome).
		Error(ctx, "Some error")
}
