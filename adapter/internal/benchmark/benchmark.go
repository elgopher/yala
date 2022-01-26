// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package benchmark

import (
	"context"
	"testing"

	"github.com/elgopher/yala/logger"
)

func Adapter(b *testing.B, adapter logger.Adapter) {
	b.Helper()

	ctx := context.Background()

	b.Run("global logger info", func(b *testing.B) {
		var global logger.Global
		global.SetAdapter(adapter)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			global.Info(ctx, "msg")
		}
	})

	b.Run("local logger info", func(b *testing.B) {
		localLogger := logger.Local{Adapter: adapter}

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			localLogger.Info(ctx, "msg")
		}
	})
}
