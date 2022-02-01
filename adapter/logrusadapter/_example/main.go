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
	log.With("field_name", "field_value").With("another", "ccc").Info(ctx, "Some info")
	log.With("parameter", "some").Warn(ctx, "Deprecated configuration parameter. It will be removed.")
	log.WithError(ErrSome).Error(ctx, "Some error")
}

func newLogrusLogger() *logrus.Logger {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	return l
}
