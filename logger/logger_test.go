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

type anyLogger interface {
	Info(context.Context, string)
	Debug(context.Context, string)
	Warn(context.Context, string)
	Error(context.Context, string)
}

type loggerMethod func(l anyLogger, ctx context.Context, msg string)

var loggerMethods = map[string]loggerMethod{
	"Debug": anyLogger.Debug,
	"Info":  anyLogger.Info,
	"Warn":  anyLogger.Warn,
	"Error": anyLogger.Error,
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
			global.With("k", "v").Warn(ctx, message)
		})

		t.Run("WithError", func(t *testing.T) {
			var global logger.Global
			global.WithError(ErrSome).Warn(ctx, message)
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
						SkippedCallerFrames: 2,
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
						SkippedCallerFrames: 2,
					},
				)
			})
		}
	})
}

func TestGlobal_With(t *testing.T) {
	globalWith := func(l anyLogger, k string, v interface{}) anyLogger {
		return l.(*logger.Global).With(k, v) // nolint:forcetypeassert // no generics still in Go
	}

	tests := map[string]struct {
		With               func(anyLogger, string, interface{}) anyLogger
		NewLoggerWithField func(adapter logger.Adapter, field logger.Field) anyLogger
	}{
		"global": {
			NewLoggerWithField: func(adapter logger.Adapter, field logger.Field) anyLogger {
				var global logger.Global
				global.SetAdapter(adapter)

				return global.With(field.Key, field.Value)
			},
			With: globalWith,
		},
		"global, adapter set on child logger": {
			NewLoggerWithField: func(adapter logger.Adapter, field logger.Field) anyLogger {
				var global logger.Global
				log := global.With(field.Key, field.Value)

				log.SetAdapter(adapter)

				return log
			},
			With: globalWith,
		},
		"global, adapter set after logger was created": {
			NewLoggerWithField: func(adapter logger.Adapter, field logger.Field) anyLogger {
				var global logger.Global
				log := global.With(field.Key, field.Value)

				global.SetAdapter(adapter)

				return log
			},
			With: globalWith,
		},
		"local": {
			NewLoggerWithField: func(adapter logger.Adapter, field logger.Field) anyLogger {
				return logger.Local{Adapter: adapter}.With(field.Key, field.Value)
			},
			With: func(l anyLogger, k string, v interface{}) anyLogger {
				return l.(logger.Logger).With(k, v) // nolint:forcetypeassert // no generics still in Go
			},
		},
	}

	field1 := logger.Field{Key: "field1_name", Value: "field1_value"}
	field2 := logger.Field{Key: "field2_name", Value: "field2_value"}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for methodName, logMessage := range loggerMethods {
				t.Run(methodName, func(t *testing.T) {
					t.Run("should log message with field", func(t *testing.T) {
						adapter := &adapterMock{}
						l := test.NewLoggerWithField(adapter, field1)
						// when
						logMessage(l, ctx, message)
						// then
						expectedFields := []logger.Field{field1}
						adapter.HasExactlyOneEntryWithFields(t, expectedFields)
					})

					t.Run("should log message with two fields", func(t *testing.T) {
						adapter := &adapterMock{}
						l := test.With(test.NewLoggerWithField(adapter, field1), field2.Key, field2.Value)
						// when
						logMessage(l, ctx, message)
						// then
						expectedFields := []logger.Field{field1, field2}
						adapter.HasExactlyOneEntryWithFields(t, expectedFields)
					})
				})
			}

			t.Run("adding field should not modify existing logger but create a new one", func(t *testing.T) {
				adapter := &adapterMock{}
				loggerWithField1 := test.NewLoggerWithField(adapter, field1)
				// when
				loggerWithBothFields := test.With(loggerWithField1, field2.Key, field2.Value)
				// then
				loggerWithField1.Info(ctx, message)
				loggerWithBothFields.Info(ctx, message)
				require.Len(t, adapter.entries, 2)
				assert.Equal(t, adapter.entries[0].Fields, []logger.Field{field1})
				assert.Equal(t, adapter.entries[1].Fields, []logger.Field{field1, field2})
			})
		})
	}
}

func TestGlobal_WithError(t *testing.T) {
	globalWithError := func(l anyLogger, err error) anyLogger {
		return l.(*logger.Global).WithError(err) // nolint:forcetypeassert // no generics still in Go
	}

	loggersWithError := map[string]struct {
		NewLoggerWithError func(adapter logger.Adapter, err error) anyLogger
		WithError          func(l anyLogger, err error) anyLogger
	}{
		"global": {
			NewLoggerWithError: func(adapter logger.Adapter, err error) anyLogger {
				var global logger.Global
				global.SetAdapter(adapter)

				return global.WithError(err)
			},
			WithError: globalWithError,
		},
		"global, adapter set on child logger": {
			NewLoggerWithError: func(adapter logger.Adapter, err error) anyLogger {
				var global logger.Global
				log := global.WithError(err)

				log.SetAdapter(adapter)

				return log
			},
			WithError: globalWithError,
		},
		"global, adapter set after logger was created": {
			NewLoggerWithError: func(adapter logger.Adapter, err error) anyLogger {
				var global logger.Global
				log := global.WithError(err)

				global.SetAdapter(adapter)

				return log
			},
			WithError: globalWithError,
		},
		"local": {
			NewLoggerWithError: func(adapter logger.Adapter, err error) anyLogger {
				return logger.Local{Adapter: adapter}.WithError(err)
			},
			WithError: func(l anyLogger, err error) anyLogger {
				return l.(logger.Logger).WithError(err) // nolint:forcetypeassert // no generics still in Go
			},
		},
	}

	for name, test := range loggersWithError {
		t.Run(name, func(t *testing.T) {
			for methodName, logMessage := range loggerMethods {
				t.Run(methodName, func(t *testing.T) {
					t.Run("should log message with error", func(t *testing.T) {
						adapter := &adapterMock{}
						l := test.NewLoggerWithError(adapter, ErrSome)
						// when
						logMessage(l, ctx, message)
						// then
						adapter.HasExactlyOneEntryWithError(t, ErrSome)
					})
				})
			}

			t.Run("adding error should not modify existing logger but create a new one", func(t *testing.T) {
				adapter := &adapterMock{}
				loggerWithSomeError := test.NewLoggerWithError(adapter, ErrSome)
				// when
				loggerWithAnotherError := test.WithError(loggerWithSomeError, ErrAnother)
				// then
				loggerWithSomeError.Error(ctx, message)
				loggerWithAnotherError.Error(ctx, message)
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
		assert.Equal(t, "10", logger.Level(10).String())
	})
}
