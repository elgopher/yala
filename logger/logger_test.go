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

type loggerMethod func(l logger.Logger, msg string)

var loggerMethods = map[string]loggerMethod{
	"Debug": logger.Logger.Debug,
	"Info":  logger.Logger.Info,
	"Warn":  logger.Logger.Warn,
	"Error": logger.Logger.Error,
}

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

func TestWith(t *testing.T) {
	type newLoggerWithField func(adapter logger.Adapter, fieldName string, fieldValue interface{}) logger.Logger

	loggersWithField := map[string]newLoggerWithField{
		"global": func(adapter logger.Adapter, fieldName string, fieldValue interface{}) logger.Logger {
			logger.SetAdapter(adapter)

			return logger.With(ctx, fieldName, fieldValue)
		},
		"local": func(adapter logger.Adapter, fieldName string, fieldValue interface{}) logger.Logger {
			return logger.Local(adapter).With(ctx, fieldName, fieldValue)
		},
	}

	const (
		field1Name  = "field1_name"
		field1Value = "field1_value"
		field2Name  = "field2_name"
		field2Value = "field2_value"
	)

	for name, newLogger := range loggersWithField {
		t.Run(name, func(t *testing.T) {
			for methodName, logMessage := range loggerMethods {
				t.Run(methodName, func(t *testing.T) {
					t.Run("should log message with field", func(t *testing.T) {
						adapter := &adapterMock{}
						l := newLogger(adapter, field1Name, field1Value)
						// when
						logMessage(l, message)
						// then
						expectedFields := []logger.Field{
							{Key: field1Name, Value: field1Value},
						}
						adapter.HasExactlyOneEntryWithFields(t, expectedFields)
					})

					t.Run("should log message with two fields", func(t *testing.T) {
						adapter := &adapterMock{}
						l := newLogger(adapter, field1Name, field1Value).With(field2Name, field2Value)
						// when
						logMessage(l, message)
						// then
						expectedFields := []logger.Field{
							{Key: field1Name, Value: field1Value},
							{Key: field2Name, Value: field2Value},
						}
						adapter.HasExactlyOneEntryWithFields(t, expectedFields)
					})
				})
			}
		})
	}
}

func TestWithError(t *testing.T) {
	type newLoggerWithError func(adapter logger.Adapter, err error) logger.Logger

	loggersWithError := map[string]newLoggerWithError{
		"global": func(adapter logger.Adapter, err error) logger.Logger {
			logger.SetAdapter(adapter)

			return logger.WithError(ctx, err)
		},
		"local": func(adapter logger.Adapter, err error) logger.Logger {
			return logger.Local(adapter).WithError(ctx, err)
		},
	}
	for name, newLogger := range loggersWithError {
		t.Run(name, func(t *testing.T) {
			for methodName, logMessage := range loggerMethods {
				t.Run(methodName, func(t *testing.T) {
					t.Run("should log message with error", func(t *testing.T) {
						adapter := &adapterMock{}
						l := newLogger(adapter, ErrSome)
						// when
						logMessage(l, message)
						// then
						adapter.HasExactlyOneEntryWithError(t, ErrSome)
					})
				})
			}

			t.Run("adding error should not modify existing logger but create a new one", func(t *testing.T) {
				adapter := &adapterMock{}
				loggerWithSomeError := newLogger(adapter, ErrSome)
				// when
				loggerWithAnotherError := loggerWithSomeError.WithError(ErrAnother)
				// then
				loggerWithSomeError.Error(message)
				loggerWithAnotherError.Error(message)
				require.Len(t, adapter.entries, 2)
				assert.Same(t, adapter.entries[0].Error, ErrSome)
				assert.Same(t, adapter.entries[1].Error, ErrAnother)
			})
		})
	}
}
