package logrusadapter_test

import (
	"context"
	"testing"

	"github.com/jacekolszak/yala/adapter/logrusadapter"
	"github.com/jacekolszak/yala/logger"
	"github.com/sirupsen/logrus"
)

func BenchmarkAdapter(b *testing.B) {
	ctx := context.Background()

	l := logrus.New()
	l.SetOutput(discardWriter{})
	logrusEntry := logrus.NewEntry(l)

	adapter := logrusadapter.Adapter{Entry: logrusEntry}
	logger.SetAdapter(adapter)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "msg") // 1200ns, 14 allocs
	}
}

type discardWriter struct{}

func (d discardWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
