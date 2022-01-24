// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package printer_test

import (
	"testing"

	"github.com/elgopher/yala/adapter/internal/benchmark"
	"github.com/elgopher/yala/adapter/printer"
)

func BenchmarkPrinter(b *testing.B) {
	adapter := printer.Adapter{Printer: printer.WriterPrinter{Writer: benchmark.DiscardWriter{}}}
	benchmark.Adapter(b, adapter)
}
