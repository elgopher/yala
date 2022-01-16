package main

import (
	"context"
	"log"
	"os"

	"github.com/jacekolszak/yala/adapter/printer"
	"github.com/jacekolszak/yala/logger"
)

func main() {
	ctx := context.Background()

	// log using fmt.Println
	logger.SetService(printer.StdoutService())

	logger.Debug(ctx, "Hello fmt")
	logger.With(ctx, "tag", "bbb").Info("Some info")
	logger.Warnf(ctx, "Be careful with %s", "hot water")
	logger.Error(ctx, "Some error")

	// log using standard log package
	standardLog := log.New(os.Stdout, "", log.LstdFlags)
	service := printer.Service{Printer: standardLog}
	logger.SetService(service)

	logger.Debug(ctx, "Hello stdlog")
	logger.With(ctx, "tag", "bbb").Info("Some info")
	logger.Warnf(ctx, "Be careful with %s", "hot water")
	logger.Error(ctx, "Some error")
}
