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
		global.With("k1", "v").With("k2", "v").With("k3", "v").Info(ctx, message) // 235-315ns, 3 allocs
	}
}

func BenchmarkWithError(b *testing.B) {
	b.ReportAllocs()

	var global logger.Global

	for i := 0; i < b.N; i++ {
		_ = global.WithError(ErrSome) // 2.6ns, 0 allocs
	}
}

func BenchmarkGlobalInfo3Fields(b *testing.B) {
	var global logger.Global

	global.SetAdapter(discardAdapter{})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		global.Info(ctx,
			message,
			"k1", "v",
			"k2", "v",
			"k3", "v") // 136ns, 1 allocs
	}
}

func BenchmarkNormalInfo3Fields(b *testing.B) {
	log := logger.WithAdapter(discardAdapter{})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		log.Info(ctx,
			message,
			"k1", "v",
			"k2", "v",
			"k3", "v") // 136ns, 1 allocs
	}
}

type discardAdapter struct{}

func (d discardAdapter) Log(context.Context, logger.Entry) {}
