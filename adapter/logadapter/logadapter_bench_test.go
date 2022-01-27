// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logadapter_test

import (
	"log"
	"testing"

	"github.com/elgopher/yala/adapter/internal/benchmark"
	"github.com/elgopher/yala/adapter/printer"
)

func BenchmarkLog(b *testing.B) {
	standardLog := log.New(benchmark.DiscardWriter{}, "", log.LstdFlags)
	adapter := printer.Adapter{Printer: standardLog}
	benchmark.Adapter(b, adapter)
}
