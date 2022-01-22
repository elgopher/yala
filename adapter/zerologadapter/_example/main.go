package main

import (
	"context"
	"errors"
	"os"

	"github.com/jacekolszak/yala/adapter/zerologadapter"
	"github.com/jacekolszak/yala/logger"
	"github.com/rs/zerolog"
)

var ErrSome = errors.New("ErrSome")

// This example shows how to use yala with zerolog adapter
func main() {
	ctx := context.Background()

	l := zerolog.New(os.Stdout)                  // create zerolog logger
	adapter := zerologadapter.Adapter{Logger: l} // create logger.Adapter for zerolog
	yalaLogger := logger.Local(adapter)          // Create yala logger

	yalaLogger.Debug(ctx, "Hello zerolog")
	yalaLogger.With(ctx, "field_name", "field_value").Info("Some info")
	yalaLogger.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	yalaLogger.WithError(ctx, ErrSome).Error("Some error")
}
