// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logrusadapter_test

import (
	"testing"

	"github.com/elgopher/yala/adapter/internal/benchmark"
	"github.com/elgopher/yala/adapter/logrusadapter"
	"github.com/sirupsen/logrus"
)

func BenchmarkLogrus(b *testing.B) {
	l := logrus.New()
	l.SetOutput(benchmark.DiscardWriter{})

	adapter := logrusadapter.Adapter{Logger: l}

	benchmark.Adapter(b, adapter)
}
