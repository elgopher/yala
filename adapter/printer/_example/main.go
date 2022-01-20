package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/jacekolszak/yala/adapter/printer"
	"github.com/jacekolszak/yala/logger"
)

var ErrSome = errors.New("ErrSome")

// This example shows how to use yala with standard fmt.Println and standard `log` package
func main() {
	ctx := context.Background()

	// log using fmt.Println
	logger.SetAdapter(printer.StdoutAdapter())

	logger.Debug(ctx, "Hello fmt")
	logger.With(ctx, "field_name", "field_value").Info("Some info")
	logger.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	logger.WithError(ctx, ErrSome).Error("Some error")

	// log using standard log package
	standardLog := log.New(os.Stdout, "", log.LstdFlags)
	adapter := printer.Adapter{Printer: standardLog}
	logger.SetAdapter(adapter)

	logger.Debug(ctx, "Hello stdlog")
	logger.With(ctx, "f1", "v1").With("f2", "f2").Info("Some info")
	logger.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	logger.WithError(ctx, ErrSome).Error("Some error")
}
