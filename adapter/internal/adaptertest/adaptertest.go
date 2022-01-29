// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package adaptertest

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/elgopher/yala/logger"
	"github.com/stretchr/testify/assert"
)

type Subject struct {
	NewAdapter       func(io.Writer) logger.Adapter
	UnmarshalMessage func(t *testing.T, line string) Message
}

type Message struct {
	Level   logger.Level
	Message string
	// fields
	Error          string
	StringField    string
	IntField       int
	Int64Field     int64
	Float32Field   float32
	Float64Field   float64
	TimeField      time.Time
	InterfaceField InterfaceField
}

type InterfaceField struct {
	NestedField string
}

func Run(t *testing.T, subject Subject) { // nolint
	ctx := context.Background()

	var entry = logger.Entry{
		Level:   logger.InfoLevel,
		Message: "message",
	}

	t.Run("should log message with no fields and no error", func(t *testing.T) {
		levels := []logger.Level{
			logger.DebugLevel,
			logger.InfoLevel,
			logger.WarnLevel,
			logger.ErrorLevel,
		}
		for _, level := range levels {
			t.Run(level.String(), func(t *testing.T) {
				var builder strings.Builder
				adapter := subject.NewAdapter(&builder)
				const messageString = "message"
				e := logger.Entry{
					Level:   level,
					Message: messageString,
				}
				// when
				adapter.Log(ctx, e)
				// then
				message := subject.UnmarshalMessage(t, builder.String())
				assert.Equal(t, level, message.Level)
				assert.Equal(t, messageString, message.Message)
			})
		}
	})

	t.Run("should log info for unknown level", func(t *testing.T) {
		var builder strings.Builder
		adapter := subject.NewAdapter(&builder)
		e := entry
		e.Level = 100
		// when
		adapter.Log(ctx, e)
		// then
		out := subject.UnmarshalMessage(t, builder.String())
		assert.Equal(t, logger.InfoLevel, out.Level)
	})

	t.Run("should add error field when Entry has an error", func(t *testing.T) {
		var builder strings.Builder
		adapter := subject.NewAdapter(&builder)
		e := entry
		const err = "err"
		e.Error = errors.New(err) // nolint
		// when
		adapter.Log(ctx, e)
		// then
		out := subject.UnmarshalMessage(t, builder.String())
		assert.Equal(t, err, out.Error)
	})

	t.Run("should log message with field", func(t *testing.T) {
		fields := map[string]struct {
			value interface{}
			get   func(Message) interface{}
		}{
			"StringField": {
				value: "value",
				get: func(message Message) interface{} {
					return message.StringField
				},
			},
			"IntField": {
				value: 1,
				get: func(message Message) interface{} {
					return message.IntField
				},
			},
			"Int64Field": {
				value: int64(1),
				get: func(message Message) interface{} {
					return message.Int64Field
				},
			},
			"Float32Field": {
				value: float32(1.1),
				get: func(message Message) interface{} {
					return message.Float32Field
				},
			},
			"Float64Field": {
				value: 1.1,
				get: func(message Message) interface{} {
					return message.Float64Field
				},
			},
			"TimeField": {
				value: time.Unix(1000, 0).UTC(), // avoid zero value
				get: func(message Message) interface{} {
					return message.TimeField
				},
			},
			"InterfaceField": {
				value: InterfaceField{NestedField: "nested"},
				get: func(message Message) interface{} {
					return message.InterfaceField
				},
			},
		}

		for name, field := range fields {
			var builder strings.Builder
			adapter := subject.NewAdapter(&builder)
			entryWithField := entry.With(
				logger.Field{
					Key:   name,
					Value: field.value,
				},
			)
			// when
			adapter.Log(ctx, entryWithField)
			// then
			out := subject.UnmarshalMessage(t, builder.String())
			assert.Equal(t, field.value, field.get(out))
		}
	})

	t.Run("should log message with two fields", func(t *testing.T) {
		var builder strings.Builder
		adapter := subject.NewAdapter(&builder)
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
		out := subject.UnmarshalMessage(t, builder.String())
		assert.Equal(t, stringFieldValue, out.StringField)
		assert.Equal(t, intFieldValue, out.IntField)
	})
}
