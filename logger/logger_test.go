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
	Info(context.Context, string, ...logger.Field)
	Debug(context.Context, string, ...logger.Field)
	Warn(context.Context, string, ...logger.Field)
	Error(context.Context, string, ...logger.Field)
}

type loggerMethod func(l anyLogger, ctx context.Context, msg string, fields ...logger.Field)

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
}

func TestNormalLogger(t *testing.T) {
	t.Run("passing nil adapter should disable logger", func(t *testing.T) {
		log := logger.WithAdapter(nil)
		log.Info(ctx, message)
	})

	t.Run("using zero value should not panic", func(t *testing.T) {
		var log logger.Logger
		assert.NotPanics(t, func() {
			log.Info(ctx, message)
		})
	})
}

func TestLogging(t *testing.T) {
	type newLogger func(adapter logger.Adapter) anyLogger

	loggers := map[string]newLogger{
		"normal": func(adapter logger.Adapter) anyLogger {
			return logger.WithAdapter(adapter)
		},
		"global": func(adapter logger.Adapter) anyLogger {
			var global logger.Global
			global.SetAdapter(adapter)

			return &global
		},
	}

	for loggerName, newLogger := range loggers {
		t.Run(loggerName, func(t *testing.T) {
			t.Run("should log message using adapter", func(t *testing.T) {
				type functionUnderTest func(l anyLogger, ctx context.Context, msg string, fields ...logger.Field)
				tests := map[logger.Level]functionUnderTest{
					logger.DebugLevel: anyLogger.Debug,
					logger.InfoLevel:  anyLogger.Info,
					logger.WarnLevel:  anyLogger.Warn,
					logger.ErrorLevel: anyLogger.Error,
				}

				for lvl, log := range tests {
					testName := lvl.String()

					t.Run(testName, func(t *testing.T) {
						adapter := &adapterMock{}
						l := newLogger(adapter)
						// when
						log(l, context.Background(), message)
						// then
						adapter.HasExactlyOneEntry(t,
							logger.Entry{
								Level:               lvl,
								Message:             message,
								SkippedCallerFrames: 2,
							},
						)
					})

					t.Run(testName+" with field", func(t *testing.T) {
						adapter := &adapterMock{}
						l := newLogger(adapter)
						// when
						log(l, context.Background(), message, logger.Field{Key: "key", Value: "value"})
						// then
						adapter.HasExactlyOneEntry(t,
							logger.Entry{
								Level:               lvl,
								Message:             message,
								Fields:              []logger.Field{{Key: "key", Value: "value"}},
								SkippedCallerFrames: 2,
							},
						)
					})

					t.Run(testName+" with two fields", func(t *testing.T) {
						adapter := &adapterMock{}
						l := newLogger(adapter)
						// when
						log(l, context.Background(),
							message,
							logger.Field{Key: "k1", Value: "v1"},
							logger.Field{Key: "k2", Value: "v2"},
						)
						// then
						adapter.HasExactlyOneEntry(t,
							logger.Entry{
								Level:               lvl,
								Message:             message,
								Fields:              []logger.Field{{"k1", "v1"}, {"k2", "v2"}},
								SkippedCallerFrames: 2,
							},
						)
					})
				}
			})
		})
	}
}

func TestWith(t *testing.T) {
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
		"normal": {
			NewLoggerWithField: func(adapter logger.Adapter, field logger.Field) anyLogger {
				return logger.WithAdapter(adapter).With(field.Key, field.Value)
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

func TestWithError(t *testing.T) {
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
		"normal": {
			NewLoggerWithError: func(adapter logger.Adapter, err error) anyLogger {
				return logger.WithAdapter(adapter).WithError(err)
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

func TestWithSkippedCallerFrame(t *testing.T) {
	tests := map[string]struct {
		newLogger    func(logger.Adapter) anyLogger
		skipOneFrame func(anyLogger, logger.Adapter) anyLogger
	}{
		"normal": {
			newLogger: func(adapter logger.Adapter) anyLogger {
				return logger.WithAdapter(adapter)
			},
			skipOneFrame: func(l anyLogger, adapter logger.Adapter) anyLogger {
				return l.(logger.Logger).WithSkippedCallerFrame() // nolint:forcetypeassert // no generics still in Go
			},
		},
		"global": {
			newLogger: func(adapter logger.Adapter) anyLogger {
				var global logger.Global
				global.SetAdapter(adapter)

				return &global
			},
			skipOneFrame: func(l anyLogger, adapter logger.Adapter) anyLogger {
				return l.(*logger.Global).WithSkippedCallerFrame() // nolint:forcetypeassert // no generics still in Go
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			defaultFrames := defaultSkippedCallerFrames(test.newLogger)

			t.Run("should skip one caller frame", func(t *testing.T) {
				adapter := &adapterMock{}
				log := test.newLogger(adapter)
				// when
				log = test.skipOneFrame(log, adapter)
				// then
				require.NotNil(t, log)
				log.Info(ctx, message)
				adapter.HasExactlyOneEntryWithSkippedCallerFrames(t, defaultFrames+1)
			})

			t.Run("should skip two frames", func(t *testing.T) {
				adapter := &adapterMock{}
				log := test.newLogger(adapter)
				// when
				log = test.skipOneFrame(log, adapter)
				log = test.skipOneFrame(log, adapter)
				// then
				require.NotNil(t, log)
				log.Info(ctx, message)
				adapter.HasExactlyOneEntryWithSkippedCallerFrames(t, defaultFrames+2)
			})

			t.Run("each created logger should be a copy", func(t *testing.T) {
				adapter := &adapterMock{}
				log := test.newLogger(adapter)
				// when
				log1 := test.skipOneFrame(log, adapter)  // +1 frame
				log2 := test.skipOneFrame(log1, adapter) // +2 frames
				// then
				require.NotNil(t, log1)
				require.NotNil(t, log2)
				log1.Info(ctx, message)
				log2.Info(ctx, message)
				// then
				require.Len(t, adapter.entries, 2)
				assert.Equal(t, defaultFrames+1, adapter.entries[0].SkippedCallerFrames)
				assert.Equal(t, defaultFrames+2, adapter.entries[1].SkippedCallerFrames)
			})
		})
	}
}

func defaultSkippedCallerFrames(newLogger func(logger.Adapter) anyLogger) int {
	adapter := &adapterMock{}
	log := newLogger(adapter)
	log.Info(ctx, message)

	return adapter.entries[0].SkippedCallerFrames
}
