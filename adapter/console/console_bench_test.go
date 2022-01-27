// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package console_test

import (
	"testing"

	"github.com/elgopher/yala/adapter/console"
	"github.com/elgopher/yala/adapter/internal/benchmark"
	"github.com/elgopher/yala/adapter/printer"
)

func BenchmarkConsole(b *testing.B) {
	adapter := printer.Adapter{
		Printer: console.WriterPrinter{
			Writer: benchmark.DiscardWriter{},
		},
	}
	benchmark.Adapter(b, adapter)
}
