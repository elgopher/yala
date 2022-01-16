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
	t.Run("passing nil global service should disable logger", func(t *testing.T) {
		logger.SetService(nil)
		logger.Info(ctx, message)
	})

	t.Run("should log message using global service", func(t *testing.T) {
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
				service := &serviceMock{}
				logger.SetService(service)
				// when
				log(ctx, message)
				// then
				service.HasExactlyOneEntry(t,
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
		service := &serviceMock{}
		logger.SetService(service)
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
		service.HasExactlyOneEntryWithFields(t, expectedFields)
	})

	t.Run("should log message with error", func(t *testing.T) {
		service := &serviceMock{}
		logger.SetService(service)
		// when
		logger.WithError(ctx, ErrSome).Error("message")
		// then
		service.HasExactlyOneEntryWithError(t, ErrSome)
	})

	t.Run("adding error should not modify existing logger but create a new one", func(t *testing.T) {
		service := &serviceMock{}
		logger.SetService(service)
		loggerWithSomeError := logger.WithError(ctx, ErrSome)
		// when
		loggerWithAnotherError := loggerWithSomeError.WithError(ErrAnother)
		// then
		loggerWithSomeError.Error(message)
		loggerWithAnotherError.Error(message)
		require.Len(t, service.entries, 2)
		assert.Same(t, service.entries[0].Error, ErrSome)
		assert.Same(t, service.entries[1].Error, ErrAnother)
	})
}

func TestLocalLogger(t *testing.T) {
	t.Run("passing nil service should disable logger", func(t *testing.T) {
		localLogger := logger.Local(nil)
		localLogger.Info(ctx, message)
	})

	t.Run("should log message using service", func(t *testing.T) {
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
				service := &serviceMock{}
				localLogger := logger.Local(service)
				// when
				log(localLogger, context.Background(), message)
				// then
				service.HasExactlyOneEntry(t,
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

type serviceMock struct {
	entries []logger.Entry
}

func (s *serviceMock) Log(_ context.Context, entry logger.Entry) {
	s.entries = append(s.entries, entry)
}

func (s *serviceMock) HasExactlyOneEntry(t *testing.T, expected logger.Entry) {
	t.Helper()

	require.Len(t, s.entries, 1)
	actual := s.entries[0]
	assert.Equal(t, expected, actual)
}

func (s *serviceMock) HasExactlyOneEntryWithFields(t *testing.T, expected []logger.Field) {
	t.Helper()

	require.Len(t, s.entries, 1)
	actual := s.entries[0].Fields
	assert.Equal(t, expected, actual)
}

func (s *serviceMock) HasExactlyOneEntryWithError(t *testing.T, expected error) {
	t.Helper()

	require.Len(t, s.entries, 1)
	actual := s.entries[0].Error
	assert.Equal(t, expected, actual)
}
