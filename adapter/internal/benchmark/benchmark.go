// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package benchmark

import (
	"context"
	"testing"
	"time"

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

	fields := map[string]interface{}{
		"string":  "str",
		"int":     1,
		"int64":   int64(64),     // nolint
		"float64": 1.64,          // nolint
		"float32": float32(1.32), // nolint
		"time":    time.Time{},
	}

	for fieldType, fieldValue := range fields {
		b.Run(fieldType, func(b *testing.B) {
			b.Run("local logger with field", func(b *testing.B) {
				localLogger := logger.Local{Adapter: adapter}

				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					localLogger.With("a", fieldValue).Info(ctx, "msg")
				}
			})
		})
	}
}
