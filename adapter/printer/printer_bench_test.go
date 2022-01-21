package printer_test

import (
	"testing"

	"github.com/jacekolszak/yala/adapter/internal/benchmark"
	"github.com/jacekolszak/yala/adapter/printer"
)

func BenchmarkPrinter(b *testing.B) {
	adapter := printer.Adapter{Printer: discardPrinter{}}
	benchmark.Adapter(b, adapter)
}

type discardPrinter struct{}

func (d discardPrinter) Println(...interface{}) {}
