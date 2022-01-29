// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package console_test

import (
	"context"
	"testing"

	"github.com/elgopher/yala/adapter/console"
	"github.com/elgopher/yala/adapter/internal/fake"
	"github.com/elgopher/yala/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriterPrinter_Println(t *testing.T) {
	t.Run("should not panic when Writer is nil", func(t *testing.T) {
		p := console.WriterPrinter{Writer: nil}
		assert.NotPanics(t, func() {
			p.Println(0, "")
		})
	})
}

func TestStderrAdapter(t *testing.T) {
	t.Run("should return adapter", func(t *testing.T) {
		adapter := console.StderrAdapter()
		require.NotNil(t, adapter)
	})

	t.Run("should log message", func(t *testing.T) {
		stderr := fake.UseFakeStderr(t)
		defer stderr.Release()

		adapter := console.StderrAdapter()
		// when
		adapter.Log(context.Background(), logger.Entry{
			Level:   logger.InfoLevel,
			Message: "message",
		})
		// then
		assert.Equal(t, "INFO message\n", stderr.FirstLine(t))
	})
}

func TestStdoutAdapter(t *testing.T) {
	t.Run("should return adapter", func(t *testing.T) {
		adapter := console.StdoutAdapter()
		require.NotNil(t, adapter)
	})
}
