// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package zerologadapter_test

import (
	"testing"

	"github.com/elgopher/yala/adapter/internal/benchmark"
	"github.com/elgopher/yala/adapter/zerologadapter"
	"github.com/rs/zerolog"
)

func BenchmarkZerolog(b *testing.B) {
	zerologLogger := zerolog.New(benchmark.DiscardWriter{})
	adapter := zerologadapter.Adapter{
		Logger: zerologLogger,
	}

	benchmark.Adapter(b, adapter)
}
