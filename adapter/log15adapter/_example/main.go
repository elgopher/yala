package main

import (
	"context"
	"errors"

	"github.com/inconshreveable/log15"
	"github.com/jacekolszak/yala/adapter/log15adapter"
	"github.com/jacekolszak/yala/logger"
)

var ErrSome = errors.New("some error")

func main() {
	ctx := context.Background()

	l := log15.New()                           // create log15 logger
	service := log15adapter.Service{Logger: l} // create logger.Service for log15
	logger.SetService(service)                 // set log15 it globally

	logger.Debug(ctx, "Hello log15")
	logger.With(ctx, "tag", "bbb").Info("Some info")
	logger.WithError(ctx, ErrSome).Error("Some error")

}
