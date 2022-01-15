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
	logger.Error(ctx, "Some error")

	// log using standard log package
	standardLog := log.New(os.Stdout, "", log.LstdFlags)
	logger.SetService(printer.Service{Printer: standardLog})

	logger.Debug(ctx, "Hello stdlog")
	logger.With(ctx, "tag", "bbb").Info("Some info")
	logger.Error(ctx, "Some error")
}
