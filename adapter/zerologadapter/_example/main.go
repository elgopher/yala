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
	logger.With(ctx, "field_name", "field_value").Info("Some info")
	logger.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	logger.WithError(ctx, errors.New("ss")).Error("Some error")
}
