package logger_test

import (
	"testing"

	"github.com/elgopher/yala/logger"
)

func BenchmarkInfo(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var global logger.Global

	for i := 0; i < b.N; i++ {
		global.Info(ctx, "msg") // 30ns, 0 allocs
	}
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
	b.ResetTimer()

	var global logger.Global

	for i := 0; i < b.N; i++ {
		_ = global.With(ctx, "k", "v") // 55ns, 1 alloc
	}
}

func BenchmarkWithError(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var global logger.Global

	for i := 0; i < b.N; i++ {
		_ = global.WithError(ctx, ErrSome) // 15ns, 0 allocs
	}
}
