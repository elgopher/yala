package main

import (
	"context"
	"errors"
	"flag"

	"github.com/elgopher/yala/adapter/glogadapter"
	"github.com/elgopher/yala/logger"
)

var ErrSome = errors.New("ErrSome")

// This example shows how to use yala with glog adapter
func main() {
	ctx := context.Background()

	flag.Parse() // glog will pick command line options like -stderrthreshold=[INFO|WARNING|ERROR]
	// create yala logger
	log := logger.WithAdapter(glogadapter.Adapter{})

	log.Debug(ctx, "Hello glog ") // Debug will be logged as Info

	log.Info(ctx,
		"Some info",
		logger.Field{Key: "field_name", Value: "field_value"},
	)

	log.Warn(ctx,
		"Deprecated configuration parameter. It will be removed.",
		logger.Field{Key: "parameter", Value: "some"},
	)

	log.WithError(ErrSome).
		Error(ctx, "Error occurred")
}
