package log15adapter_test

import (
	"context"
	"testing"

	"github.com/inconshreveable/log15"
	"github.com/jacekolszak/yala/adapter/log15adapter"
	"github.com/jacekolszak/yala/logger"
)

func BenchmarkAdapter(b *testing.B) {
	ctx := context.Background()

	log15logger := log15.New()
	log15logger.SetHandler(log15.DiscardHandler())

	adapter := log15adapter.Adapter{
		Logger: log15logger,
	}

	b.Run("global logger info", func(b *testing.B) {
		logger.SetAdapter(adapter)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			logger.Info(ctx, "msg") // 1080ns, 3 allocs
		}
	})

	b.Run("local logger info", func(b *testing.B) {
		localLogger := logger.Local(adapter)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			localLogger.Info(ctx, "msg") // 960ns, 3 allocs
		}
	})
}
