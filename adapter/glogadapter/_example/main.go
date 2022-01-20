package main

import (
	"context"
	"errors"
	"flag"

	"github.com/jacekolszak/yala/adapter/glogadapter"
	"github.com/jacekolszak/yala/logger"
)

var ErrSome = errors.New("ErrSome")

func main() {
	ctx := context.Background()

	flag.Parse()                             // glog will pick command line options like -stderrthreshold=[INFO|WARNING|ERROR]
	logger.SetAdapter(glogadapter.Adapter{}) // set glog adapter globally

	logger.Debug(ctx, "Hello glog ") // Debug will be logged as Info
	logger.With(ctx, "field_name", "field_value").Info("Some info")
	logger.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	logger.WithError(ctx, ErrSome).Error("Error occurred")
}
