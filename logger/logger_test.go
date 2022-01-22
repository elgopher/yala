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
		var global logger.Global
		global.SetAdapter(nil)
		global.Info(ctx, message)
	})

	t.Run("should log warning that global adapter was not set", func(t *testing.T) {
		t.Run("Warn", func(t *testing.T) {
			var global logger.Global
			global.Warn(ctx, message)
		})

		t.Run("With", func(t *testing.T) {
			var global logger.Global
			global.With(ctx, "k", "v").Warn(message)
		})

		t.Run("WithError", func(t *testing.T) {
			var global logger.Global
			global.WithError(ctx, ErrSome).Warn(message)
		})
	})

	t.Run("should log message using global adapter", func(t *testing.T) {
		var global logger.Global

		type functionUnderTest func(ctx context.Context, msg string)
		tests := map[logger.Level]functionUnderTest{
			logger.DebugLevel: global.Debug,
			logger.InfoLevel:  global.Info,
			logger.WarnLevel:  global.Warn,
			logger.ErrorLevel: global.Error,
		}

		for lvl, log := range tests {
			testName := string(lvl)

			t.Run(testName, func(t *testing.T) {
				adapter := &adapterMock{}
				global.SetAdapter(adapter)
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
	type newLoggerWithField func(adapter logger.Adapter, field logger.Field) logger.Logger

	loggersWithField := map[string]newLoggerWithField{
		"global": func(adapter logger.Adapter, field logger.Field) logger.Logger {
			var global logger.Global
			global.SetAdapter(adapter)

			return global.With(ctx, field.Key, field.Value)
		},
		"local": func(adapter logger.Adapter, field logger.Field) logger.Logger {
			return logger.Local(adapter).With(ctx, field.Key, field.Value)
		},
	}

	field1 := logger.Field{Key: "field1_name", Value: "field1_value"}
	field2 := logger.Field{Key: "field2_name", Value: "field2_value"}

	for name, newLogger := range loggersWithField {
		t.Run(name, func(t *testing.T) {
			for methodName, logMessage := range loggerMethods {
				t.Run(methodName, func(t *testing.T) {
					t.Run("should log message with field", func(t *testing.T) {
						adapter := &adapterMock{}
						l := newLogger(adapter, field1)
						// when
						logMessage(l, message)
						// then
						expectedFields := []logger.Field{field1}
						adapter.HasExactlyOneEntryWithFields(t, expectedFields)
					})

					t.Run("should log message with two fields", func(t *testing.T) {
						adapter := &adapterMock{}
						l := newLogger(adapter, field1).With(field2.Key, field2.Value)
						// when
						logMessage(l, message)
						// then
						expectedFields := []logger.Field{field1, field2}
						adapter.HasExactlyOneEntryWithFields(t, expectedFields)
					})
				})
			}

			t.Run("adding field should not modify existing logger but create a new one", func(t *testing.T) {
				adapter := &adapterMock{}
				loggerWithField1 := newLogger(adapter, field1)
				// when
				loggerWithBothFields := loggerWithField1.With(field2.Key, field2.Value)
				// then
				loggerWithField1.Info(message)
				loggerWithBothFields.Info(message)
				require.Len(t, adapter.entries, 2)
				assert.Equal(t, adapter.entries[0].Fields, []logger.Field{field1})
				assert.Equal(t, adapter.entries[1].Fields, []logger.Field{field1, field2})
			})
		})
	}
}

func TestWithError(t *testing.T) {
	type newLoggerWithError func(adapter logger.Adapter, err error) logger.Logger

	loggersWithError := map[string]newLoggerWithError{
		"global": func(adapter logger.Adapter, err error) logger.Logger {
			var global logger.Global
			global.SetAdapter(adapter)

			return global.WithError(ctx, err)
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
