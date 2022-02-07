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

	log.InfoFields(ctx, "Some info", logger.Fields{
		"field_name": "field_value",
		"other_name": "field_value",
	})

	log.WarnFields(ctx, "Deprecated configuration parameter. It will be removed.", logger.Fields{
		"parameter": "some",
	})

	log.ErrorCause(ctx, "Some error", ErrSome)
}
