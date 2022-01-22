package logger_test

import (
	"context"
	"sync"
	"testing"

	"github.com/jacekolszak/yala/logger"
	"github.com/stretchr/testify/assert"
)

func TestConcurrency(t *testing.T) {
	t.Run("global log functions", func(t *testing.T) {
		adapter := &concurrencySafeAdapter{}
		logger.SetAdapter(adapter)

		var wg sync.WaitGroup

		for i := 0; i < 1000; i++ {
			wg.Add(1)

			go func() {
				// when
				logger.Debug(ctx, message)
				logger.Info(ctx, message)
				logger.Warn(ctx, message)
				logger.Error(ctx, message)
				logger.With(ctx, "k", "v").Info(message)
				logger.WithError(ctx, ErrSome).Error(message)
				wg.Done()
			}()
		}

		wg.Wait()
		// then
		assert.Equal(t, adapter.Count(), 6000)
	})

	t.Run("local log functions", func(t *testing.T) {
		adapter := &concurrencySafeAdapter{}
		localLogger := logger.Local(adapter)

		var wg sync.WaitGroup

		for i := 0; i < 1000; i++ {
			wg.Add(1)

			go func() {
				// when
				localLogger.Debug(ctx, message)
				localLogger.Info(ctx, message)
				localLogger.Warn(ctx, message)
				localLogger.Error(ctx, message)
				localLogger.With(ctx, "k", "v").Info(message)
				localLogger.WithError(ctx, ErrSome).Error(message)
				wg.Done()
			}()
		}

		wg.Wait()
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