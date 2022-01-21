package logger_test

import (
	"testing"

	"github.com/jacekolszak/yala/logger"
)

func BenchmarkInfo(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "msg") // 30ns, 0 allocs
	}
}

func BenchmarkLogger_Info(b *testing.B) {
	loggerWithField := logger.With(ctx, "k", "v")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		loggerWithField.Info("msg") // 12ns, 0 allocs
	}
}

func BenchmarkWith(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = logger.With(ctx, "k", "v") // 55ns, 1 alloc
	}
}
