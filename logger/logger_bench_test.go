package logger_test

import (
	"context"
	"testing"

	"github.com/jacekolszak/yala/adapter/logrusadapter"
	"github.com/jacekolszak/yala/logger"
	"github.com/sirupsen/logrus"
)

func BenchmarkInfo(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "msg") // 5ns
	}
}

func BenchmarkLogger_Info(b *testing.B) {
	ctx := context.Background()
	l := logger.With(ctx, "k", "v")

	for i := 0; i < b.N; i++ {
		l.Info("msg") // 12ns
	}
}

func BenchmarkWith(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		_ = logger.With(ctx, "k", "v") // 55ns
	}
}

type discardWriter struct{}

func (d discardWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func BenchmarkLogrus(b *testing.B) {
	ctx := context.Background()

	logrusLogger := logrus.New()
	logrusLogger.SetOutput(discardWriter{})

	adapter := logrusadapter.Service{
		Entry: logrus.NewEntry(logrusLogger),
	}
	logger.SetAdapter(adapter)

	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "msg") // 1200ns
	}
}
