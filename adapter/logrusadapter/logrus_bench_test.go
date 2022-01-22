package logrusadapter_test

import (
	"testing"

	"github.com/jacekolszak/yala/adapter/internal/benchmark"
	"github.com/jacekolszak/yala/adapter/logrusadapter"
	"github.com/sirupsen/logrus"
)

func BenchmarkLogrus(b *testing.B) {
	l := logrus.New()
	l.SetOutput(benchmark.DiscardWriter{})
	logrusEntry := logrus.NewEntry(l)

	adapter := logrusadapter.Adapter{Entry: logrusEntry}

	benchmark.Adapter(b, adapter)
}
