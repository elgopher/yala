package main

import (
	"context"
	"errors"

	"github.com/elgopher/yala/adapter/log15adapter"
	"github.com/elgopher/yala/logger"
	"github.com/inconshreveable/log15"
)

var ErrSome = errors.New("some error")

// This example shows how to use yala with log15 adapter
func main() {
	ctx := context.Background()

	l := log15.New()                           // create log15 logger
	adapter := log15adapter.Adapter{Logger: l} // create logger.Adapter for log15
	log := logger.WithAdapter(adapter)         // create yala logger

	log.Debug(ctx, "Hello log15")

	log.Info(ctx,
		"Some info",
		logger.Field{Key: "field_name", Value: "field_value"},
	)

	log.Warn(ctx,
		"Deprecated configuration parameter. It will be removed.",
		logger.Field{Key: "parameter", Value: "some"},
	)

	log.WithError(ErrSome).
		Error(ctx, "Some error")
}
