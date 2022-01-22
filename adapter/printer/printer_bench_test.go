package printer_test

import (
	"testing"

	"github.com/jacekolszak/yala/adapter/internal/benchmark"
	"github.com/jacekolszak/yala/adapter/printer"
)

func BenchmarkPrinter(b *testing.B) {
	adapter := printer.Adapter{Printer: printer.WriterPrinter{Writer: discardWriter{}}}
	benchmark.Adapter(b, adapter)
}

type discardWriter struct{}

func (d discardWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
