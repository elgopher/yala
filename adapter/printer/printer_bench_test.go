package printer_test

import (
	"testing"

	"github.com/jacekolszak/yala/adapter/internal/benchmark"
	"github.com/jacekolszak/yala/adapter/printer"
)

func BenchmarkPrinter(b *testing.B) {
	adapter := printer.Adapter{Printer: printer.WriterPrinter{Writer: benchmark.DiscardWriter{}}}
	benchmark.Adapter(b, adapter)
}
