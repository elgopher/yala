package logger_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jacekolszak/yala/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const message = "message"

var ErrSome = errors.New("some error")
var ErrAnother = errors.New("another error")

var ctx = context.Background()

func TestGlobalLogging(t *testing.T) {
	t.Run("passing nil global adapter should disable logger", func(t *testing.T) {
		logger.SetAdapter(nil)
		logger.Info(ctx, message)
	})

	t.Run("should log message using global adapter", func(t *testing.T) {
		type functionUnderTest func(ctx context.Context, msg string)
		tests := map[logger.Level]functionUnderTest{
			logger.DebugLevel: logger.Debug,
			logger.InfoLevel:  logger.Info,
			logger.WarnLevel:  logger.Warn,
			logger.ErrorLevel: logger.Error,
		}

		for lvl, log := range tests {
			testName := string(lvl)

			t.Run(testName, func(t *testing.T) {
				adapter := &adapterMock{}
				logger.SetAdapter(adapter)
				// when
				log(ctx, message)
				// then
				adapter.HasExactlyOneEntry(t,
					logger.Entry{
						Level:               lvl,
						Message:             message,
						SkippedCallerFrames: 4,
					},
				)
			})
		}
	})

	t.Run("should log message with field", func(t *testing.T) {
		adapter := &adapterMock{}
		logger.SetAdapter(adapter)
		const (
			fieldName  = "field_name"
			fieldValue = "field_value"
		)
		// when
		logger.With(ctx, fieldName, fieldValue).Info("message")
		// then
		expectedFields := []logger.Field{
			{Key: fieldName, Value: fieldValue},
		}
		adapter.HasExactlyOneEntryWithFields(t, expectedFields)
	})

	t.Run("should log message with error", func(t *testing.T) {
		adapter := &adapterMock{}
		logger.SetAdapter(adapter)
		// when
		logger.WithError(ctx, ErrSome).Error("message")
		// then
		adapter.HasExactlyOneEntryWithError(t, ErrSome)
	})

	t.Run("adding error should not modify existing logger but create a new one", func(t *testing.T) {
		adapter := &adapterMock{}
		logger.SetAdapter(adapter)
		loggerWithSomeError := logger.WithError(ctx, ErrSome)
		// when
		loggerWithAnotherError := loggerWithSomeError.WithError(ErrAnother)
		// then
		loggerWithSomeError.Error(message)
		loggerWithAnotherError.Error(message)
		require.Len(t, adapter.entries, 2)
		assert.Same(t, adapter.entries[0].Error, ErrSome)
		assert.Same(t, adapter.entries[1].Error, ErrAnother)
	})
}

func TestLocalLogger(t *testing.T) {
	t.Run("passing nil adapter should disable logger", func(t *testing.T) {
		localLogger := logger.Local(nil)
		localLogger.Info(ctx, message)
	})

	t.Run("should log message using adapter", func(t *testing.T) {
		type functionUnderTest func(l logger.LocalLogger, ctx context.Context, msg string)
		tests := map[logger.Level]functionUnderTest{
			logger.DebugLevel: logger.LocalLogger.Debug,
			logger.InfoLevel:  logger.LocalLogger.Info,
			logger.WarnLevel:  logger.LocalLogger.Warn,
			logger.ErrorLevel: logger.LocalLogger.Error,
		}

		for lvl, log := range tests {
			testName := string(lvl)

			t.Run(testName, func(t *testing.T) {
				adapter := &adapterMock{}
				localLogger := logger.Local(adapter)
				// when
				log(localLogger, context.Background(), message)
				// then
				adapter.HasExactlyOneEntry(t,
					logger.Entry{
						Level:               lvl,
						Message:             message,
						SkippedCallerFrames: 2,
					},
				)
			})
		}
	})
}

type adapterMock struct {
	entries []logger.Entry
}

func (a *adapterMock) Log(_ context.Context, entry logger.Entry) {
	a.entries = append(a.entries, entry)
}

func (a *adapterMock) HasExactlyOneEntry(t *testing.T, expected logger.Entry) {
	t.Helper()

	require.Len(t, a.entries, 1)
	actual := a.entries[0]
	assert.Equal(t, expected, actual)
}

func (a *adapterMock) HasExactlyOneEntryWithFields(t *testing.T, expected []logger.Field) {
	t.Helper()

	require.Len(t, a.entries, 1)
	actual := a.entries[0].Fields
	assert.Equal(t, expected, actual)
}

func (a *adapterMock) HasExactlyOneEntryWithError(t *testing.T, expected error) {
	t.Helper()

	require.Len(t, a.entries, 1)
	actual := a.entries[0].Error
	assert.Equal(t, expected, actual)
}
