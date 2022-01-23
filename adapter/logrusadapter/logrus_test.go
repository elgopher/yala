package logrusadapter_test

import (
	"context"
	"testing"

	"github.com/elgopher/yala/adapter/logrusadapter"
	"github.com/elgopher/yala/logger"
	"github.com/stretchr/testify/assert"
)

const message = "message"

func TestAdapter_Log(t *testing.T) {
	ctx := context.Background()

	t.Run("should not panic when entry is nil", func(t *testing.T) {
		adapter := logrusadapter.Adapter{Entry: nil}
		assert.NotPanics(t, func() {
			adapter.Log(ctx, logger.Entry{
				Level:   logger.InfoLevel,
				Message: message,
			})
		})
	})
}
