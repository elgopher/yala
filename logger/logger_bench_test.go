// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logger_test

import (
	"testing"

	"github.com/elgopher/yala/logger"
)

func BenchmarkInfo(b *testing.B) {
	b.ReportAllocs()

	var global logger.Global

	for i := 0; i < b.N; i++ {
		global.Info(ctx, "msg") // 27.5ns, 0 allocs
	}
}

func BenchmarkParallelInfo(b *testing.B) {
	var global logger.Global

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			global.Info(ctx, "msg") // 3.6 ns/op (8 goroutines, 8 cores)
		}
	})
}

func BenchmarkLogger_Info(b *testing.B) {
	var global logger.Global
	loggerWithField := global.With(ctx, "k", "v")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		loggerWithField.Info("msg") // 12ns, 0 allocs
	}
}

func BenchmarkWith(b *testing.B) {
	b.ReportAllocs()

	var global logger.Global

	for i := 0; i < b.N; i++ {
		_ = global.With(ctx, "k", "v") // 55ns, 1 alloc
	}
}

func BenchmarkWithError(b *testing.B) {
	b.ReportAllocs()

	var global logger.Global

	for i := 0; i < b.N; i++ {
		_ = global.WithError(ctx, ErrSome) // 15ns, 0 allocs
	}
}
