// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logrusadapter_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/elgopher/yala/adapter/logrusadapter"
	"github.com/elgopher/yala/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const message = "message"

var entry = logger.Entry{
	Level:   logger.InfoLevel,
	Message: "message",
}

func TestAdapter_Log(t *testing.T) {
	ctx := context.Background()

	t.Run("should not panic when entry is nil", func(t *testing.T) {
		adapter := logrusadapter.Adapter{Entry: nil}
		assert.NotPanics(t, func() {
			adapter.Log(ctx, logger.Entry{
				Level:   logger.InfoLevel,
				Message: message,
			})
		})
	})

	t.Run("should log message with no fields and no error", func(t *testing.T) {
		levels := map[string]logger.Level{
			"debug":   logger.DebugLevel,
			"info":    logger.InfoLevel,
			"warning": logger.WarnLevel,
			"error":   logger.ErrorLevel,
		}
		for logrusLevel, level := range levels {
			t.Run(logrusLevel, func(t *testing.T) {
				var writer strings.Builder
				adapter := newAdapter(&writer)
				entryWithLevel := logger.Entry{
					Level:   level,
					Message: "message",
				}
				// when
				adapter.Log(ctx, entryWithLevel)
				// then
				out := unmarshalMessage(t, &writer)
				assert.Equal(t,
					logrusMessage{
						Level: logrusLevel,
						Msg:   entryWithLevel.Message,
					},
					out,
				)
			})
		}
	})

	t.Run("should log info for unknown level", func(t *testing.T) {
		var writer strings.Builder
		adapter := newAdapter(&writer)
		e := logger.Entry{
			Level:   -100,
			Message: "message",
		}
		// when
		adapter.Log(ctx, e)
		// then
		out := unmarshalMessage(t, &writer)
		assert.Equal(t, "info", out.Level)
	})

	t.Run("should add error field when Entry has an error", func(t *testing.T) {
		var writer strings.Builder
		adapter := newAdapter(&writer)
		e := entry
		const err = "err"
		e.Error = errors.New(err) // nolint
		// when
		adapter.Log(ctx, e)
		// then
		out := unmarshalMessage(t, &writer)
		assert.Equal(t, err, out.Error)
	})

	t.Run("should add field", func(t *testing.T) {
		var writer strings.Builder
		adapter := newAdapter(&writer)
		e := entry.With(logger.Field{
			Key:   "Field1",
			Value: "field1",
		})
		// when
		adapter.Log(ctx, e)
		// then
		out := unmarshalMessage(t, &writer)
		assert.Equal(t, "field1", out.Field1)
	})

	t.Run("should add two fields", func(t *testing.T) {
		var writer strings.Builder
		adapter := newAdapter(&writer)
		entryWithTwoFields := entry.With(logger.Field{
			Key:   "Field1",
			Value: "field1",
		}).With(logger.Field{
			Key:   "Field2",
			Value: "field2",
		})
		// when
		adapter.Log(ctx, entryWithTwoFields)
		// then
		out := unmarshalMessage(t, &writer)
		assert.Equal(t, "field1", out.Field1)
		assert.Equal(t, "field2", out.Field2)
	})
}

func newAdapter(writer io.Writer) logrusadapter.Adapter {
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})
	logrusLogger.SetOutput(writer)
	logrusLogger.SetLevel(logrus.DebugLevel)

	return logrusadapter.Adapter{Entry: logrus.NewEntry(logrusLogger)}
}

func unmarshalMessage(t *testing.T, stringer fmt.Stringer) logrusMessage {
	t.Helper()

	out := logrusMessage{}
	err := json.Unmarshal([]byte(stringer.String()), &out)
	require.NoError(t, err)

	return out
}

type logrusMessage struct {
	Level  string
	Msg    string
	Error  string
	Field1 string
	Field2 string
}
