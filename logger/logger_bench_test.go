// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logger_test

import (
	"context"
	"testing"

	"github.com/elgopher/yala/logger"
)

func BenchmarkInfo(b *testing.B) {
	b.ReportAllocs()

	var global logger.Global

	for i := 0; i < b.N; i++ {
		global.Info(ctx, "msg") // 12ns, 0 allocs
	}
}

func BenchmarkParallelInfo(b *testing.B) {
	var global logger.Global

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			global.Info(ctx, "msg") // 1 ns/op (8 goroutines, 8 cores)
		}
	})
}

func BenchmarkLogger_Info(b *testing.B) {
	var global logger.Global
	loggerWithField := global.With("k", "v")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		loggerWithField.Info(ctx, "msg") // 8ns, 0 allocs
	}
}

func BenchmarkWith(b *testing.B) {
	b.ReportAllocs()

	var global logger.Global

	for i := 0; i < b.N; i++ {
		_ = global.With("k", "v") // 41ns, 1 alloc
	}
}

func BenchmarkMultipleWith(b *testing.B) {
	b.ReportAllocs()

	var global logger.Global

	for i := 0; i < b.N; i++ {
		global.With("k1", "v").With("k2", "v").With("k3", "v").Info(ctx, "ads") // 530-700ns, 6 allocs
	}
}

func BenchmarkWithError(b *testing.B) {
	b.ReportAllocs()

	var global logger.Global

	for i := 0; i < b.N; i++ {
		_ = global.WithError(ErrSome) // 2.6ns, 0 allocs
	}
}

func BenchmarkGlobal_InfoFields(b *testing.B) {
	var log logger.Global

	log.SetAdapter(discardAdapter{})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		log.InfoFields(ctx, message, logger.Fields{
			"a": 1,
			"b": 2,
			"c": 3,
		})
	}
}

type discardAdapter struct{}

func (d discardAdapter) Log(context.Context, logger.Entry) {}
