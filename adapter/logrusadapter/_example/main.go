package main

import (
	"context"
	"errors"

	"github.com/elgopher/yala/adapter/logrusadapter"
	"github.com/elgopher/yala/logger"
	"github.com/sirupsen/logrus"
)

var ErrSome = errors.New("some error")

// This example shows how to use yala with logrus adapter
func main() {
	ctx := context.Background()

	// First create a logrus logger.
	l := newLogrusLogger()
	// Then create a logger.Adapter.
	adapter := logrusadapter.Adapter{
		Logger: l, // you can also pass *logrus.Entry
	}
	// Create yala logger
	log := logger.WithAdapter(adapter)

	log.Debug(ctx, "Hello logrus ")

	log.Info(ctx,
		"Some info",
		"field_name", "field_value",
		"another", "ccc",
	)

	log.Warn(ctx,
		"Deprecated configuration parameter. It will be removed.",
		"parameter", "some",
	)

	log.WithError(ErrSome).
		Error(ctx, "Some error")
}

func newLogrusLogger() *logrus.Logger {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	return l
}
