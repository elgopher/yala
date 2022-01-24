// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package zapadapter_test

import (
	"net/url"
	"testing"

	"github.com/elgopher/yala/adapter/internal/benchmark"
	"github.com/elgopher/yala/adapter/zapadapter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func BenchmarkZap(b *testing.B) {
	_ = zap.RegisterSink("discard", func(url *url.URL) (zap.Sink, error) {
		return discardSink{}, nil
	})
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"discard://"}
	zapLogger, err := cfg.Build()
	require.NoError(b, err)

	adapter := zapadapter.Adapter{
		Logger: zapLogger,
	}

	benchmark.Adapter(b, adapter)
}

type discardSink struct {
	benchmark.DiscardWriter
}

func (d discardSink) Sync() error {
	return nil
}

func (d discardSink) Close() error {
	return nil
}
