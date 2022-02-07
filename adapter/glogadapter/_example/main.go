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

	log.InfoFields(ctx, "Some info", logger.Fields{
		"field_name": "field_value",
		"other_name": "field_value",
	})

	log.WarnFields(ctx, "Deprecated configuration parameter. It will be removed.", logger.Fields{
		"parameter": "some",
	})

	log.ErrorCause(ctx, "Error occurred", ErrSome)
}
