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

	log.InfoFields(ctx, "Some info", logger.Fields{
		"field_name": "field_value",
		"other_name": "field_value",
	})

	log.WarnFields(ctx, "Deprecated configuration parameter. It will be removed.", logger.Fields{
		"parameter": "some value",
	})

	log.ErrorCause(ctx, "Some error", ErrSome)
}
