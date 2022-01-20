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
	logger.SetAdapter(adapter)                   // set it globally

	logger.Debug(ctx, "Hello zerolog")
	logger.With(ctx, "field_name", "field_value").Info("Some info")
	logger.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	logger.WithError(ctx, ErrSome).Error("Some error")
}
