package log15adapter_test

import (
	"testing"

	"github.com/inconshreveable/log15"
	"github.com/jacekolszak/yala/adapter/internal/benchmark"
	"github.com/jacekolszak/yala/adapter/log15adapter"
)

func BenchmarkLog15(b *testing.B) {
	log15logger := log15.New()
	log15logger.SetHandler(log15.StreamHandler(discardWriter{}, log15.LogfmtFormat()))

	adapter := log15adapter.Adapter{
		Logger: log15logger,
	}

	benchmark.Adapter(b, adapter)
}

type discardWriter struct{}

func (d discardWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
