package main

import (
	"context"

	"github.com/elgopher/yala/adapter/console"
	"github.com/elgopher/yala/logger"
)

// This example shows how to reuse loggers
func main() {
	ctx := context.Background()

	log := logger.WithAdapter(console.StdoutAdapter())

	// requestLogger will log all messages with at least two fields: request_id and user
	requestLogger := log.With("request_id", "123").With("user", "elgopher")

	requestLogger.Debug(ctx, "request started")

	requestLogger.Debug(ctx,
		"sql update executed",
		"rows_updated", 3,
		"table", "gophers")

	requestLogger.Debug(ctx, "request finished")
}
