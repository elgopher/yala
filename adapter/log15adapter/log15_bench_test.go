// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package log15adapter_test

import (
	"testing"

	"github.com/elgopher/yala/adapter/internal/benchmark"
	"github.com/elgopher/yala/adapter/log15adapter"
	"github.com/inconshreveable/log15"
)

func BenchmarkLog15(b *testing.B) {
	log15logger := log15.New()
	log15logger.SetHandler(log15.StreamHandler(benchmark.DiscardWriter{}, log15.LogfmtFormat()))

	adapter := log15adapter.Adapter{
		Logger: log15logger,
	}

	benchmark.Adapter(b, adapter)
}
