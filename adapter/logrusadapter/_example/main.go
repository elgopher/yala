package main

import (
	"context"
	"errors"

	"github.com/jacekolszak/yala/adapter/logrusadapter"
	"github.com/jacekolszak/yala/logger"
	"github.com/sirupsen/logrus"
)

var ErrSome = errors.New("some error")

func main() {
	ctx := context.Background()

	// First create logrus logger
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	// Then create a logger.Service
	service := logrusadapter.Service{
		Entry: logrus.NewEntry(l),
	}
	// And use it globally
	logger.SetService(service)

	logger.Debug(ctx, "Hello logrus ")
	logger.With(ctx, "tag", "bbb").With("another", "ccc").Info("Some info")
	logger.Warnf(ctx, "Be careful with %s", "hot water")
	logger.WithError(ctx, ErrSome).Error("Some error")
}
