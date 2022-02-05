// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logger_test

import (
	"context"
	"fmt"
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
		// when
		runConcurrently(1000, func() {
			global.Debug(ctx, message)
			global.Info(ctx, message)
			global.Warn(ctx, message)
			global.Error(ctx, message)
			global.With("k", "v").Info(ctx, message)
			global.WithError(ErrSome).Error(ctx, message)
		})
		// then
		assert.Equal(t, adapter.Count(), 6000)
	})

	t.Run("normal log functions", func(t *testing.T) {
		adapter := &concurrencySafeAdapter{}
		log := logger.WithAdapter(adapter)
		// when
		runConcurrently(1000, func() {
			log.Debug(ctx, message)
			log.Info(ctx, message)
			log.Warn(ctx, message)
			log.Error(ctx, message)
			log.With("k", "v").Info(ctx, message)
			log.WithError(ErrSome).Error(ctx, message)
		})
		// then
		assert.Equal(t, adapter.Count(), 6000)
	})

	t.Run("With should not data race when -race flag is used", func(t *testing.T) {
		type newLogger func(numberOfFields int) func()

		tests := map[string]newLogger{
			"normal": func(numberOfFields int) func() {
				log := logger.WithAdapter(&concurrencySafeAdapter{})
				for i := 0; i < numberOfFields; i++ {
					log = log.With("k", "v")
				}

				return func() {
					log.With("k", "v").Info(ctx, message)
				}
			},
			"global": func(numberOfFields int) func() {
				var log logger.Global
				log.SetAdapter(&concurrencySafeAdapter{})

				theLog := &log

				for i := 0; i < numberOfFields; i++ {
					theLog = log.With("k", "v")
				}

				return func() {
					theLog.With("k", "v").Info(ctx, message)
				}
			},
		}

		for loggerName, createNewLogger := range tests {
			t.Run(loggerName, func(t *testing.T) {
				for fields := 0; fields < 16; fields++ {
					testName := fmt.Sprintf("log with %d fields", fields)

					t.Run(testName, func(t *testing.T) {
						f := createNewLogger(fields)
						runConcurrently(1000, func() {
							f()
						})
					})
				}
			})
		}
	})
}

func runConcurrently(goroutines int, code func()) {
	var waitGroup sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		waitGroup.Add(1)

		go func() {
			code()
			waitGroup.Done()
		}()
	}

	waitGroup.Wait()
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
