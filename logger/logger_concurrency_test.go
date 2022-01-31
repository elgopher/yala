// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logger_test

import (
	"context"
	"sync"
	"testing"

	"github.com/elgopher/yala/logger"
	"github.com/stretchr/testify/assert"
)

func TestConcurrency(t *testing.T) {
	t.Run("global log functions", func(t *testing.T) {
		adapter := &concurrencySafeAdapter{}
		var global logger.Global
		global.SetAdapter(adapter)

		var waitGroup sync.WaitGroup

		for i := 0; i < 1000; i++ {
			waitGroup.Add(1)

			go func() {
				// when
				global.Debug(ctx, message)
				global.Info(ctx, message)
				global.Warn(ctx, message)
				global.Error(ctx, message)
				global.With("k", "v").Info(ctx, message)
				global.WithError(ErrSome).Error(ctx, message)
				waitGroup.Done()
			}()
		}

		waitGroup.Wait()
		// then
		assert.Equal(t, adapter.Count(), 6000)
	})

	t.Run("normal log functions", func(t *testing.T) {
		adapter := &concurrencySafeAdapter{}
		log := logger.WithAdapter(adapter)

		var waitGroup sync.WaitGroup

		for i := 0; i < 1000; i++ {
			waitGroup.Add(1)

			go func() {
				// when
				log.Debug(ctx, message)
				log.Info(ctx, message)
				log.Warn(ctx, message)
				log.Error(ctx, message)
				log.With("k", "v").Info(ctx, message)
				log.WithError(ErrSome).Error(ctx, message)
				waitGroup.Done()
			}()
		}

		waitGroup.Wait()
		// then
		assert.Equal(t, adapter.Count(), 6000)
	})
}

type concurrencySafeAdapter struct {
	mutex   sync.Mutex
	entries int
}

func (c *concurrencySafeAdapter) Log(context.Context, logger.Entry) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries++
}

func (c *concurrencySafeAdapter) Count() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.entries
}
