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

	log.InfoFields(ctx, "Some info", logger.Fields{
		"field_name": "field_value",
		"other_name": "field_value",
	})

	log.WarnFields(ctx, "Deprecated configuration parameter. It will be removed.", logger.Fields{
		"parameter": "some",
	})

	log.ErrorCause(ctx, "Some error", ErrSome)
}
