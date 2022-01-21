package glogadapter_test

import (
	"context"
	"testing"

	"github.com/jacekolszak/yala/adapter/glogadapter"
	"github.com/jacekolszak/yala/logger"
)

func BenchmarkAdapter(b *testing.B) {
	ctx := context.Background()
	adapter := glogadapter.Adapter{}

	b.Run("global logger info", func(b *testing.B) {
		logger.SetAdapter(adapter)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			logger.Info(ctx, "msg") // 1000ns, 4 allocs
		}
	})

	b.Run("local logger info", func(b *testing.B) {
		localLogger := logger.Local(adapter)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			localLogger.Info(ctx, "msg") // 970ns, 4 allocs
		}
	})
}
