package main

import (
	"context"
	"errors"

	"github.com/inconshreveable/log15"
	"github.com/jacekolszak/yala/adapter/log15adapter"
	"github.com/jacekolszak/yala/logger"
)

var ErrSome = errors.New("some error")

// This example shows how to use yala with log15 adapter
func main() {
	ctx := context.Background()

	l := log15.New()                           // create log15 logger
	adapter := log15adapter.Adapter{Logger: l} // create logger.Adapter for log15
	logger.SetAdapter(adapter)                 // set log15 it globally

	logger.Debug(ctx, "Hello log15")
	logger.With(ctx, "field_name", "field_value").Info("Some info")
	logger.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	logger.WithError(ctx, ErrSome).Error("Some error")
}
