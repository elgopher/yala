package glogadapter_test

import (
	"testing"

	"github.com/jacekolszak/yala/adapter/glogadapter"
	"github.com/jacekolszak/yala/adapter/internal/benchmark"
)

func BenchmarkGlog(b *testing.B) {
	adapter := glogadapter.Adapter{}
	benchmark.Adapter(b, adapter)
}
