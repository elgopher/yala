// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logrusadapter_test

import (
	"context"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/elgopher/yala/adapter/internal/adaptertest"
	"github.com/elgopher/yala/adapter/logrusadapter"
	"github.com/elgopher/yala/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const message = "message"

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

	adaptertest.Run(t, adaptertest.Subject{
		NewAdapter:       newAdapter,
		UnmarshalMessage: unmarshalMessage,
	})
}

func newAdapter(writer io.Writer) logger.Adapter { // nolint
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})
	logrusLogger.SetOutput(writer)
	logrusLogger.SetLevel(logrus.DebugLevel)

	return logrusadapter.Adapter{Entry: logrus.NewEntry(logrusLogger)}
}

var levelMapping = map[string]logger.Level{
	"debug":   logger.DebugLevel,
	"info":    logger.InfoLevel,
	"warning": logger.WarnLevel,
	"error":   logger.ErrorLevel,
}

func unmarshalMessage(t *testing.T, line string) adaptertest.Message {
	t.Helper()

	out := logrusMessage{}
	err := json.Unmarshal([]byte(line), &out)
	require.NoError(t, err)

	return adaptertest.Message{
		Level:          levelMapping[out.Level],
		Message:        out.Msg,
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

type logrusMessage struct {
	Level          string
	Msg            string
	Error          string
	StringField    string
	IntField       int
	Int64Field     int64
	Float32Field   float32
	Float64Field   float64
	TimeField      time.Time
	InterfaceField adaptertest.InterfaceField
}
