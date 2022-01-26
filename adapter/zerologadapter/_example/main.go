package main

import (
	"context"
	"errors"
	"os"

	"github.com/elgopher/yala/adapter/zerologadapter"
	"github.com/elgopher/yala/logger"
	"github.com/rs/zerolog"
)

var ErrSome = errors.New("ErrSome")

// This example shows how to use yala with zerolog adapter
func main() {
	ctx := context.Background()

	l := zerolog.New(os.Stdout)                  // create zerolog logger
	adapter := zerologadapter.Adapter{Logger: l} // create logger.Adapter for zerolog
	log := logger.Local{Adapter: adapter}        // Create yala logger

	log.Debug(ctx, "Hello zerolog")
	log.With(ctx, "field_name", "field_value").Info("Some info")
	log.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	log.WithError(ctx, ErrSome).Error("Some error")
}
