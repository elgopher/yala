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
	log := logger.WithAdapter(adapter)           // Create yala logger

	log.Debug(ctx, "Hello zerolog")

	log.With("field_name", "field_value").
		Info(ctx, "Some info")

	log.With("parameter", "some").
		Warn(ctx, "Deprecated configuration parameter. It will be removed.")

	log.WithError(ErrSome).
		Error(ctx, "Some error")
}
