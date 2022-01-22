package zerologadapter_test

import (
	"testing"

	"github.com/jacekolszak/yala/adapter/internal/benchmark"
	"github.com/jacekolszak/yala/adapter/zerologadapter"
	"github.com/rs/zerolog"
)

func BenchmarkZerolog(b *testing.B) {
	zerologLogger := zerolog.New(benchmark.DiscardWriter{})
	adapter := zerologadapter.Adapter{
		Logger: zerologLogger,
	}

	benchmark.Adapter(b, adapter)
}
