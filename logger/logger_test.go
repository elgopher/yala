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

	t.Run("should log message using adapter", func(t *testing.T) {
		type functionUnderTest func(l logger.Logger, ctx context.Context, msg string)
		tests := map[logger.Level]functionUnderTest{
			logger.DebugLevel: logger.Logger.Debug,
			logger.InfoLevel:  logger.Logger.Info,
			logger.WarnLevel:  logger.Logger.Warn,
			logger.ErrorLevel: logger.Logger.Error,
		}

		for lvl, log := range tests {
			testName := lvl.String()

			t.Run(testName, func(t *testing.T) {
				adapter := &adapterMock{}
				normalLogger := logger.WithAdapter(adapter)
				// when
				log(normalLogger, context.Background(), message)
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

func TestLogFields(t *testing.T) {
	type fieldsLogger interface {
		DebugFields(context.Context, string, logger.Fields)
		InfoFields(context.Context, string, logger.Fields)
		WarnFields(context.Context, string, logger.Fields)
		ErrorFields(context.Context, string, logger.Fields)
	}

	type newLogger func(adapter logger.Adapter, loggerFields logger.Fields) fieldsLogger

	tests := map[string]newLogger{
		"normal logger": func(adapter logger.Adapter, loggerFields logger.Fields) fieldsLogger {
			log := logger.WithAdapter(adapter)

			return log.WithFields(loggerFields)
		},
		"global logger": func(adapter logger.Adapter, loggerFields logger.Fields) fieldsLogger {
			var log logger.Global
			log.SetAdapter(adapter)

			return log.WithFields(loggerFields)
		},
	}

	levelToMethodMapping := map[logger.Level]func(fieldsLogger, context.Context, string, logger.Fields){
		logger.DebugLevel: fieldsLogger.DebugFields,
		logger.InfoLevel:  fieldsLogger.InfoFields,
		logger.WarnLevel:  fieldsLogger.WarnFields,
		logger.ErrorLevel: fieldsLogger.ErrorFields,
	}

	for name, createLogger := range tests {
		//
		t.Run(name, func(t *testing.T) {
			for lvl, logFields := range levelToMethodMapping {
				//
				t.Run(lvl.String(), func(t *testing.T) {
					//
					t.Run("should log message with empty fields", func(t *testing.T) {
						adapter := &adapterMock{}
						log := createLogger(adapter, nil)
						// when
						logFields(log, ctx, message, logger.Fields{})
						// then
						adapter.HasExactlyOneEntry(t, logger.Entry{
							Level:               lvl,
							Message:             message,
							SkippedCallerFrames: 2,
						})
					})

					t.Run("should log message with nil fields", func(t *testing.T) {
						adapter := &adapterMock{}
						log := createLogger(adapter, nil)
						// when
						logFields(log, ctx, message, nil)
						// then
						adapter.HasExactlyOneEntry(t, logger.Entry{
							Level:               lvl,
							Message:             message,
							SkippedCallerFrames: 2,
						})
					})

					t.Run("should log message with one field", func(t *testing.T) {
						adapter := &adapterMock{}
						log := createLogger(adapter, nil)
						// when
						logFields(log, ctx, message, logger.Fields{
							"k": "v",
						})
						// then
						adapter.HasExactlyOneEntry(t, logger.Entry{
							Level:               lvl,
							Fields:              []logger.Field{{"k", "v"}},
							Message:             message,
							SkippedCallerFrames: 2,
						})
					})

					t.Run("should log message with two fields", func(t *testing.T) {
						adapter := &adapterMock{}
						log := createLogger(adapter, nil)
						// when
						logFields(log, ctx, message, logger.Fields{
							"k1": "v1",
							"k2": "v2",
						})
						// then
						assert.ElementsMatch(t,
							[]logger.Field{
								{"k1", "v1"},
								{"k2", "v2"},
							},
							adapter.entries[0].Fields)
					})

					t.Run("should append fields to existing logger fields", func(t *testing.T) {
						adapter := &adapterMock{}
						log := createLogger(adapter, logger.Fields{
							"k1": "v1",
							"k2": "v2",
						})
						// when
						logFields(log, ctx, message, logger.Fields{
							"k3": "v3",
							"k4": "v4",
						})
						// then
						assert.ElementsMatch(t,
							[]logger.Field{
								{"k1", "v1"},
								{"k2", "v2"},
								{"k3", "v3"},
								{"k4", "v4"},
							},
							adapter.entries[0].Fields)
					})
				})
			}
		})
	}
}

func TestLogCause(t *testing.T) {
	type errorLogger interface {
		ErrorCause(context.Context, string, error)
		ErrorCauseFields(context.Context, string, error, logger.Fields)
	}

	type newLogger func(adapter logger.Adapter, originalError error) errorLogger

	loggers := map[string]newLogger{
		"normal logger": func(adapter logger.Adapter, originalError error) errorLogger {
			log := logger.WithAdapter(adapter)

			return log.WithError(originalError)
		},
		"global logger": func(adapter logger.Adapter, originalError error) errorLogger {
			var log logger.Global
			log.SetAdapter(adapter)

			return log.WithError(originalError)
		},
	}

	for loggerName, createNewLogger := range loggers {
		t.Run(loggerName, func(t *testing.T) {
			t.Run("should log cause", func(t *testing.T) {
				adapter := &adapterMock{}
				log := createNewLogger(adapter, nil)
				// when
				log.ErrorCause(ctx, message, ErrSome)
				// then
				adapter.HasExactlyOneEntryWithError(t, ErrSome)
			})

			t.Run("should log cause and field", func(t *testing.T) {
				adapter := &adapterMock{}
				log := createNewLogger(adapter, nil)
				// when
				log.ErrorCauseFields(ctx, message, ErrSome, logger.Fields{
					"k": "v",
				})
				// then
				adapter.HasExactlyOneEntry(t, logger.Entry{
					Level:               logger.ErrorLevel,
					Message:             message,
					Fields:              []logger.Field{{"k", "v"}},
					Error:               ErrSome,
					SkippedCallerFrames: 2,
				})
			})

			methods := map[string]func(errorLogger, error){
				"ErrorCauseFields": func(errorLogger errorLogger, cause error) {
					errorLogger.ErrorCauseFields(ctx, message, cause, logger.Fields{})
				},
				"ErrorCause": func(errorLogger errorLogger, cause error) {
					errorLogger.ErrorCause(ctx, message, cause)
				},
			}

			t.Run("should override original error", func(t *testing.T) {
				for methodName, runLoggerMethodWithCause := range methods {
					t.Run(methodName, func(t *testing.T) {
						adapter := &adapterMock{}
						log := createNewLogger(adapter, ErrSome)
						// when
						runLoggerMethodWithCause(log, ErrAnother)
						// then
						adapter.HasExactlyOneEntryWithError(t, ErrAnother)
					})
				}
			})

			t.Run("should override error with nil", func(t *testing.T) {
				for methodName, runLoggerMethodWithCause := range methods {
					t.Run(methodName, func(t *testing.T) {
						adapter := &adapterMock{}
						log := createNewLogger(adapter, ErrSome)
						// when
						runLoggerMethodWithCause(log, nil)
						// then
						adapter.HasExactlyOneEntryWithError(t, nil)
					})
				}
			})
		})
	}
}

func TestWith(t *testing.T) {
	globalWith := func(l anyLogger, k string, v interface{}) anyLogger {
		return l.(*logger.Global).With(k, v) //nolint:forcetypeassert // no generics still in Go
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
				return l.(logger.Logger).With(k, v) //nolint:forcetypeassert // no generics still in Go
			},
		},
		"global, WithFields": {
			NewLoggerWithField: func(adapter logger.Adapter, field logger.Field) anyLogger {
				var global logger.Global
				global.SetAdapter(adapter)

				return global.With(field.Key, field.Value)
			},
			With: func(l anyLogger, k string, v interface{}) anyLogger {
				return l.(*logger.Global).WithFields(logger.Fields{k: v}) //nolint:forcetypeassert // no generics still in Go
			},
		},
		"normal, WithFields": {
			NewLoggerWithField: func(adapter logger.Adapter, field logger.Field) anyLogger {
				return logger.WithAdapter(adapter).With(field.Key, field.Value)
			},
			With: func(l anyLogger, k string, v interface{}) anyLogger {
				return l.(logger.Logger).WithFields(logger.Fields{k: v}) //nolint:forcetypeassert // no generics still in Go
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
				assert.Equal(t, []logger.Field{field1}, adapter.entries[0].Fields)
				assert.Equal(t, []logger.Field{field1, field2}, adapter.entries[1].Fields)
			})
		})
	}
}

func TestLogger_WithFields(t *testing.T) {
	t.Run("should create a new logger with fields", func(t *testing.T) {
		adapter := &adapterMock{}
		newLogger := logger.WithAdapter(adapter)
		// when
		newLogger = newLogger.WithFields(logger.Fields{
			"k1": "v1",
			"k2": "v2",
		})
		// then
		newLogger.Info(ctx, message)
		require.Len(t, adapter.entries, 1)
		assert.ElementsMatch(t,
			[]logger.Field{{"k1", "v1"}, {"k2", "v2"}},
			adapter.entries[0].Fields,
		)
	})
}

func TestGlobal_WithFields(t *testing.T) {
	t.Run("should create a new logger with fields", func(t *testing.T) {
		adapter := &adapterMock{}
		var log logger.Global
		log.SetAdapter(adapter)
		// when
		newLogger := log.WithFields(logger.Fields{
			"k1": "v1",
			"k2": "v2",
		})
		// then
		newLogger.Info(ctx, message)
		require.Len(t, adapter.entries, 1)
		assert.ElementsMatch(t,
			[]logger.Field{{"k1", "v1"}, {"k2", "v2"}},
			adapter.entries[0].Fields,
		)
	})
}

func TestWithError(t *testing.T) {
	globalWithError := func(l anyLogger, err error) anyLogger {
		return l.(*logger.Global).WithError(err) //nolint:forcetypeassert // no generics still in Go
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
				return l.(logger.Logger).WithError(err) //nolint:forcetypeassert // no generics still in Go
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
				return l.(logger.Logger).WithSkippedCallerFrame() //nolint:forcetypeassert // no generics still in Go
			},
		},
		"global": {
			newLogger: func(adapter logger.Adapter) anyLogger {
				var global logger.Global
				global.SetAdapter(adapter)

				return &global
			},
			skipOneFrame: func(l anyLogger, adapter logger.Adapter) anyLogger {
				return l.(*logger.Global).WithSkippedCallerFrame() //nolint:forcetypeassert // no generics still in Go
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
