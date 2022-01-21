package zerologadapter_test

import (
	"testing"

	"github.com/jacekolszak/yala/adapter/internal/benchmark"
	"github.com/jacekolszak/yala/adapter/zerologadapter"
	"github.com/rs/zerolog"
)

func BenchmarkZerolog(b *testing.B) {
	zerologLogger := zerolog.New(discardWriter{})
	adapter := zerologadapter.Adapter{
		Logger: zerologLogger,
	}

	benchmark.Adapter(b, adapter)
}

type discardWriter struct{}

func (d discardWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
