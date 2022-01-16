package main

import (
	"context"
	"errors"
	"os"

	"github.com/jacekolszak/yala/adapter/zerologadapter"
	"github.com/jacekolszak/yala/logger"
	"github.com/rs/zerolog"
)

func main() {
	ctx := context.Background()

	l := zerolog.New(os.Stdout)                  // create zerolog logger
	service := zerologadapter.Service{Logger: l} // create logger.Service for zerolog
	logger.SetService(service)                   // set it globally

	logger.Debug(ctx, "Hello zerolog")
	logger.With(ctx, "tag", "bbb").Info("Some info")
	logger.Warnf(ctx, "Be careful with %s", "hot water")
	logger.WithError(ctx, errors.New("ss")).Error("Some error")
}
