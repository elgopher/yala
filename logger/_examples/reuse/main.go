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
	requestLogger := log.WithFields(logger.Fields{
		"request_id": "123",
		"user":       "elgopher",
	})

	requestLogger.Debug(ctx, "request started")

	requestLogger.
		DebugFields(ctx, "sql update executed", logger.Fields{
			"rows_updated": 3,
			"table":        "gophers",
		})

	requestLogger.Debug(ctx, "request finished")
}
