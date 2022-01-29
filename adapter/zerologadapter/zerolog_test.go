package zerologadapter_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/elgopher/yala/adapter/zerologadapter"
	"github.com/elgopher/yala/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var entry = logger.Entry{
	Level:   logger.InfoLevel,
	Message: "message",
}

func TestAdapter_Log(t *testing.T) {
	ctx := context.Background()

	t.Run("should not panic for zero-value adapter", func(t *testing.T) {
		assert.NotPanics(t, func() {
			var adapter zerologadapter.Adapter
			adapter.Log(ctx, entry)
		})
	})

	t.Run("should log message with no fields and no error", func(t *testing.T) {
		levels := map[string]logger.Level{
			"debug": logger.DebugLevel,
			"info":  logger.InfoLevel,
			"warn":  logger.WarnLevel,
			"error": logger.ErrorLevel,
		}
		for zerologLevel, level := range levels {
			t.Run(zerologLevel, func(t *testing.T) {
				var writer strings.Builder
				adapter := newAdapter(&writer)
				e := logger.Entry{
					Level:   level,
					Message: "message",
				}
				// when
				adapter.Log(ctx, e)
				// then
				out := unmarshalMessage(t, &writer)
				assert.Equal(t,
					zerologMessage{
						Level:   zerologLevel,
						Message: entry.Message,
					},
					out,
				)
			})
		}
	})

	t.Run("should log info for unknown level", func(t *testing.T) {
		var writer strings.Builder
		adapter := newAdapter(&writer)
		e := entry
		e.Level = 100
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

	t.Run("should log message with field", func(t *testing.T) {
		fields := map[string]struct {
			value interface{}
			get   func(zerologMessage) interface{}
		}{
			"StringField": {
				value: "value",
				get: func(message zerologMessage) interface{} {
					return message.StringField
				},
			},
			"IntField": {
				value: 1,
				get: func(message zerologMessage) interface{} {
					return message.IntField
				},
			},
			"Int64Field": {
				value: int64(1),
				get: func(message zerologMessage) interface{} {
					return message.Int64Field
				},
			},
			"Float32Field": {
				value: float32(1.1),
				get: func(message zerologMessage) interface{} {
					return message.Float32Field
				},
			},
			"Float64Field": {
				value: 1.1,
				get: func(message zerologMessage) interface{} {
					return message.Float64Field
				},
			},
			"TimeField": {
				value: time.Unix(1000, 0).UTC(), // avoid zero value
				get: func(message zerologMessage) interface{} {
					return message.TimeField
				},
			},
			"InterfaceField": {
				value: interfaceField{NestedField: "nested"},
				get: func(message zerologMessage) interface{} {
					return message.InterfaceField
				},
			},
		}
		for name, field := range fields {
			var writer strings.Builder
			adapter := newAdapter(&writer)
			entryWithField := entry.With(
				logger.Field{
					Key:   name,
					Value: field.value,
				},
			)
			// when
			adapter.Log(ctx, entryWithField)
			// then
			out := unmarshalMessage(t, &writer)
			assert.Equal(t, field.value, field.get(out))
		}
	})

	t.Run("should log message with two fields", func(t *testing.T) {
		var writer strings.Builder
		adapter := newAdapter(&writer)
		const (
			stringFieldValue = "string"
			intFieldValue    = 2
		)
		entryWithFields := entry.With(
			logger.Field{
				Key:   "StringField",
				Value: stringFieldValue,
			},
		)
		entryWithFields = entryWithFields.With(
			logger.Field{
				Key:   "IntField",
				Value: intFieldValue,
			},
		)
		// when
		adapter.Log(ctx, entryWithFields)
		// then
		out := unmarshalMessage(t, &writer)
		assert.Equal(t, stringFieldValue, out.StringField)
		assert.Equal(t, intFieldValue, out.IntField)
	})
}

func newAdapter(writer io.Writer) zerologadapter.Adapter {
	return zerologadapter.Adapter{
		Logger: zerolog.New(writer),
	}
}

func unmarshalMessage(t *testing.T, stringer fmt.Stringer) zerologMessage {
	t.Helper()

	out := zerologMessage{}
	err := json.Unmarshal([]byte(stringer.String()), &out)
	require.NoError(t, err)

	return out
}

type zerologMessage struct {
	Level   string
	Message string

	// fields
	Error          string
	StringField    string
	IntField       int
	Int64Field     int64
	Float32Field   float32
	Float64Field   float64
	TimeField      time.Time
	InterfaceField interfaceField
}

type interfaceField struct {
	NestedField string
}
