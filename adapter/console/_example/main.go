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
	log := logger.Local{Adapter: console.StdoutAdapter()}

	log.Debug(ctx, "Hello fmt")
	log.With(ctx, "field_name", "field_value").Info("Some info")
	log.With(ctx, "parameter", "some value").Warn("Deprecated configuration parameter. It will be removed.")
	log.WithError(ctx, ErrSome).Error("Some error")
}
