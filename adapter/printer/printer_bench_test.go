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
