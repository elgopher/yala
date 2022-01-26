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

	// First create a logrus logger entry. This is basically a Logger with some optional fields.
	entry := newLogrusEntry()
	// Then create a logger.Adapter
	adapter := logrusadapter.Adapter{
		Entry: entry, // inject logrus
	}
	// Create yala logger
	log := logger.Local{Adapter: adapter}

	log.Debug(ctx, "Hello logrus ")
	log.With(ctx, "field_name", "field_value").With("another", "ccc").Info("Some info")
	log.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	log.WithError(ctx, ErrSome).Error("Some error")
}

func newLogrusEntry() *logrus.Entry {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	return logrus.NewEntry(l)
}
