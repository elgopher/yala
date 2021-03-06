// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logadapter_test

import (
	"context"
	"testing"

	"github.com/elgopher/yala/adapter/logadapter"
	"github.com/elgopher/yala/logger"
	"github.com/stretchr/testify/assert"
)

const message = "message"

func TestAdapter_Log(t *testing.T) {
	ctx := context.Background()

	t.Run("should not panic when entry is nil", func(t *testing.T) {
		adapter := logadapter.Adapter(nil)
		assert.NotPanics(t, func() {
			adapter.Log(ctx, logger.Entry{
				Level:   logger.InfoLevel,
				Message: message,
			})
		})
	})
}
