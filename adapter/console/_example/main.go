package main

import (
	"context"
	"errors"

	"github.com/elgopher/yala/adapter/console"
	"github.com/elgopher/yala/logger"
)

var ErrSome = errors.New("ErrSome")

// This example shows how to use yala with minimal console adapter
func main() {
	ctx := context.Background()

	// log to console, stdout
	log := logger.WithAdapter(console.StdoutAdapter())

	log.Debug(ctx, "Hello fmt")

	log.Info(ctx,
		"Some info",
		logger.Field{Key: "field_name", Value: "field_value"},
	)

	log.Warn(ctx,
		"Deprecated configuration parameter. It will be removed.",
		logger.Field{Key: "parameter", Value: "some value"},
	)

	log.WithError(ErrSome).
		Error(ctx, "Some error")
}
