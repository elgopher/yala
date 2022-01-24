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
	logrusEntry := logrus.NewEntry(l)

	adapter := logrusadapter.Adapter{Entry: logrusEntry}

	benchmark.Adapter(b, adapter)
}
