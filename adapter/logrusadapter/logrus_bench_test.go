package logrusadapter_test

import (
	"testing"

	"github.com/jacekolszak/yala/adapter/internal/benchmark"
	"github.com/jacekolszak/yala/adapter/logrusadapter"
	"github.com/sirupsen/logrus"
)

func BenchmarkLogrus(b *testing.B) {
	l := logrus.New()
	l.SetOutput(discardWriter{})
	logrusEntry := logrus.NewEntry(l)

	adapter := logrusadapter.Adapter{Entry: logrusEntry}

	benchmark.Adapter(b, adapter)
}

type discardWriter struct{}

func (d discardWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
