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

	flag.Parse()                             // glog does not work without parsing the flags first
	logger.SetService(glogadapter.Service{}) // set glog service globally

	logger.Debug(ctx, "Hello glog ")
	logger.With(ctx, "tag", "bbb").Info("Some info")
	logger.WithError(ctx, ErrSome).Error("Error occurred")
}
