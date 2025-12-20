// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package benchmark

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/elgopher/yala/logger"
)

var ErrSome = errors.New("ErrSome")

// Adapter runs benchmarks on any implementation of logger.Adapter.
func Adapter(b *testing.B, adapter logger.Adapter) { //nolint:funlen
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

	b.Run("global logger info with three fields", func(b *testing.B) {
		var global logger.Global
		global.SetAdapter(adapter)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			global.InfoFields(ctx, "msg", logger.Fields{
				"field1": "value",
				"field2": "value",
				"field3": "value",
			})
		}
	})

	b.Run("normal logger info", func(b *testing.B) {
		log := logger.WithAdapter(adapter)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			log.Info(ctx, "msg")
		}
	})

	fields := map[string]interface{}{
		"string":  "str",
		"int":     1,
		"int64":   int64(64),
		"float64": 1.64,
		"float32": float32(1.32),
		"time":    time.Time{},
	}

	for fieldType, fieldValue := range fields {
		b.Run(fieldType, func(b *testing.B) {
			b.Run("normal logger with field", func(b *testing.B) {
				log := logger.WithAdapter(adapter)

				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					log.With("a", fieldValue).Info(ctx, "msg")
				}
			})
		})
	}

	b.Run("normal logger error with cause and two fields", func(b *testing.B) {
		log := logger.WithAdapter(adapter)

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			log.ErrorCauseFields(ctx, "msg", ErrSome, logger.Fields{
				"field1": "value1",
				"field2": "value2",
			})
		}
	})
}
