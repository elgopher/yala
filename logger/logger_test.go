// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logger_test

import (
	"context"
	"errors"
	"testing"

	"github.com/elgopher/yala/logger"
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
			testName := lvl.String()

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
		localLogger := logger.Local{Adapter: nil}
		localLogger.Info(ctx, message)
	})

	t.Run("using zero value should not panic", func(t *testing.T) {
		var localLogger logger.Local
		assert.NotPanics(t, func() {
			localLogger.Info(ctx, message)
		})
	})

	t.Run("should log message using adapter", func(t *testing.T) {
		type functionUnderTest func(l logger.Local, ctx context.Context, msg string)
		tests := map[logger.Level]functionUnderTest{
			logger.DebugLevel: logger.Local.Debug,
			logger.InfoLevel:  logger.Local.Info,
			logger.WarnLevel:  logger.Local.Warn,
			logger.ErrorLevel: logger.Local.Error,
		}

		for lvl, log := range tests {
			testName := lvl.String()

			t.Run(testName, func(t *testing.T) {
				adapter := &adapterMock{}
				localLogger := logger.Local{Adapter: adapter}
				// when
				log(localLogger, context.Background(), message)
				// then
				adapter.HasExactlyOneEntry(t,
					logger.Entry{
						Level:               lvl,
						Message:             message,
						SkippedCallerFrames: 3,
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
			return logger.Local{Adapter: adapter}.With(ctx, field.Key, field.Value)
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
			return logger.Local{Adapter: adapter}.WithError(ctx, err)
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

func TestEntry_With(t *testing.T) {
	field1 := logger.Field{Key: "k1", Value: "v1"}
	field2 := logger.Field{Key: "k2", Value: "v2"}

	t.Run("should add field to empty entry", func(t *testing.T) {
		entry := logger.Entry{}
		newEntry := entry.With(field1)
		assert.Empty(t, entry.Fields)
		require.Len(t, newEntry.Fields, 1)
		assert.Equal(t, newEntry.Fields[0], field1)
	})

	t.Run("should add field to entry with one field", func(t *testing.T) {
		entry := logger.Entry{}.With(field1)
		newEntry := entry.With(field2)
		require.Len(t, entry.Fields, 1)
		require.Len(t, newEntry.Fields, 2)
		assert.Equal(t, entry.Fields[0], field1)
		assert.Equal(t, newEntry.Fields[0], field1)
		assert.Equal(t, newEntry.Fields[1], field2)
	})
}

func TestLevel_MoreSevereThan(t *testing.T) {
	t.Run("should return true", func(t *testing.T) {
		assert.True(t, logger.InfoLevel.MoreSevereThan(logger.DebugLevel))
		assert.True(t, logger.WarnLevel.MoreSevereThan(logger.DebugLevel))
		assert.True(t, logger.ErrorLevel.MoreSevereThan(logger.DebugLevel))

		assert.True(t, logger.WarnLevel.MoreSevereThan(logger.InfoLevel))
		assert.True(t, logger.ErrorLevel.MoreSevereThan(logger.InfoLevel))

		assert.True(t, logger.ErrorLevel.MoreSevereThan(logger.WarnLevel))
	})

	t.Run("should return false", func(t *testing.T) {
		assert.False(t, logger.DebugLevel.MoreSevereThan(logger.DebugLevel))

		assert.False(t, logger.DebugLevel.MoreSevereThan(logger.InfoLevel))
		assert.False(t, logger.InfoLevel.MoreSevereThan(logger.InfoLevel))

		assert.False(t, logger.DebugLevel.MoreSevereThan(logger.WarnLevel))
		assert.False(t, logger.InfoLevel.MoreSevereThan(logger.WarnLevel))
		assert.False(t, logger.WarnLevel.MoreSevereThan(logger.WarnLevel))

		assert.False(t, logger.DebugLevel.MoreSevereThan(logger.ErrorLevel))
		assert.False(t, logger.InfoLevel.MoreSevereThan(logger.ErrorLevel))
		assert.False(t, logger.WarnLevel.MoreSevereThan(logger.ErrorLevel))
		assert.False(t, logger.ErrorLevel.MoreSevereThan(logger.ErrorLevel))
	})
}

func TestLevel_String(t *testing.T) {
	t.Run("should convert to string", func(t *testing.T) {
		assert.Equal(t, "DEBUG", logger.DebugLevel.String())
		assert.Equal(t, "INFO", logger.InfoLevel.String())
		assert.Equal(t, "WARN", logger.WarnLevel.String())
		assert.Equal(t, "ERROR", logger.ErrorLevel.String())
	})
}
