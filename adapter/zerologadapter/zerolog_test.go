// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package zerologadapter_test

import (
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/elgopher/yala/adapter/internal/adaptertest"
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

	adaptertest.Run(t, adaptertest.Subject{
		NewAdapter: func(writer io.Writer) logger.Adapter {
			return zerologadapter.Adapter{Logger: zerolog.New(writer)}
		},
		UnmarshalMessage: unmarshalMessage,
	})
}

var levelsMapping = map[string]logger.Level{
	"debug": logger.DebugLevel,
	"info":  logger.InfoLevel,
	"warn":  logger.WarnLevel,
	"error": logger.ErrorLevel,
}

func unmarshalMessage(t *testing.T, line string) adaptertest.Message {
	t.Helper()

	out := zerologMessage{}
	bytes := []byte(line)
	err := json.Unmarshal(bytes, &out)
	require.NoError(t, err)

	return adaptertest.Message{
		Level:          levelsMapping[out.Level],
		Message:        out.Message,
		Error:          out.Error,
		StringField:    out.StringField,
		IntField:       out.IntField,
		Int64Field:     out.Int64Field,
		Float32Field:   out.Float32Field,
		Float64Field:   out.Float64Field,
		TimeField:      out.TimeField,
		InterfaceField: out.InterfaceField,
	}
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
	InterfaceField adaptertest.InterfaceField
}
