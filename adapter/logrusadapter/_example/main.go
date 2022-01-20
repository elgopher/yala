package main

import (
	"context"
	"errors"

	"github.com/jacekolszak/yala/adapter/logrusadapter"
	"github.com/jacekolszak/yala/logger"
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
	// And use adapter globally
	logger.SetAdapter(adapter)

	logger.Debug(ctx, "Hello logrus ")
	logger.With(ctx, "field_name", "field_value").With("another", "ccc").Info("Some info")
	logger.With(ctx, "parameter", "some").Warn("Deprecated configuration parameter. It will be removed.")
	logger.WithError(ctx, ErrSome).Error("Some error")
}

func newLogrusEntry() *logrus.Entry {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	return logrus.NewEntry(l)
}
