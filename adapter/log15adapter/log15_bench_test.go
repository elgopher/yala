package log15adapter_test

import (
	"testing"

	"github.com/inconshreveable/log15"
	"github.com/jacekolszak/yala/adapter/internal/benchmark"
	"github.com/jacekolszak/yala/adapter/log15adapter"
)

func BenchmarkLog15(b *testing.B) {
	log15logger := log15.New()
	log15logger.SetHandler(log15.DiscardHandler())

	adapter := log15adapter.Adapter{
		Logger: log15logger,
	}

	benchmark.Adapter(b, adapter)
}
