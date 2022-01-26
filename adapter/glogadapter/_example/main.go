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
	// create local logger
	log := logger.Local{Adapter: glogadapter.Adapter{}}

	log.Debug(ctx, "Hello glog ") // Debug will be logged as Info
	log.With(ctx, "field_name", "field_value").Info("Some info")
	log.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	log.WithError(ctx, ErrSome).Error("Error occurred")
}
