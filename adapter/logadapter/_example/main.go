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
	yalaLogger := logger.WithAdapter(adapter)

	yalaLogger.Debug(ctx, "Hello standard log")

	yalaLogger.InfoFields(ctx, "Some info", logger.Fields{
		"f1": "v1",
		"f2": "v2",
	})

	yalaLogger.WarnFields(ctx, "Deprecated configuration parameter. It will be removed.", logger.Fields{
		"parameter": "some",
	})

	yalaLogger.ErrorCause(ctx, "Some error", ErrSome)
}
