package zerologadapter_test

import (
	"context"
	"testing"

	"github.com/jacekolszak/yala/adapter/zerologadapter"
	"github.com/jacekolszak/yala/logger"
	"github.com/rs/zerolog"
)

func BenchmarkZerolog(b *testing.B) {
	ctx := context.Background()

	zerologLogger := zerolog.New(discardWriter{})

	adapter := zerologadapter.Adapter{
		Logger: zerologLogger,
	}
	logger.SetAdapter(adapter)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "msg") // 105ns, 0 allocs
	}
}

type discardWriter struct{}

func (d discardWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
