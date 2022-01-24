// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package glogadapter_test

import (
	"testing"

	"github.com/elgopher/yala/adapter/glogadapter"
	"github.com/elgopher/yala/adapter/internal/benchmark"
)

func BenchmarkGlog(b *testing.B) {
	adapter := glogadapter.Adapter{}
	benchmark.Adapter(b, adapter)
}
