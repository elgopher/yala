// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package zapadapter_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/elgopher/yala/adapter/internal/adaptertest"
	"github.com/elgopher/yala/adapter/zapadapter"
	"github.com/elgopher/yala/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

const message = "message"

func TestAdapter_Log(t *testing.T) {
	ctx := context.Background()

	t.Run("should not panic when logger is nil", func(t *testing.T) {
		adapter := zapadapter.Adapter{Logger: nil}
		assert.NotPanics(t, func() {
			adapter.Log(ctx, logger.Entry{
				Level:   logger.InfoLevel,
				Message: message,
			})
		})
	})

	t.Run("should log caller", func(t *testing.T) {
		var builder strings.Builder
		adapter := newAdapter(&builder)
		// when
		adapter.Log(ctx, logger.Entry{
			Level:   logger.InfoLevel,
			Message: message,
		})
		// then
		msg := unmarshalZapMessage(t, builder.String())
		const expectedPrefix = "zapadapter/zapadapter_test.go:"
		assert.Truef(t, strings.HasPrefix(msg.C, expectedPrefix), "caller %s has no prefix %s", msg.C, expectedPrefix)
	})

	adaptertest.Run(t, adaptertest.Subject{
		NewAdapter:       newAdapter,
		UnmarshalMessage: unmarshalMessage,
	})
}

func newAdapter(writer io.Writer) logger.Adapter { // nolint
	scheme := generateUniqueScheme() // Zap does not allow to override existing scheme
	_ = zap.RegisterSink(scheme, func(url *url.URL) (zap.Sink, error) {
		return sinkWriter{Writer: writer}, nil
	})

	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{scheme + "://"}
	cfg.Encoding = "json"

	zapLogger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return zapadapter.Adapter{Logger: zapLogger}
}

var sequence atomic.Uint32

func generateUniqueScheme() string {
	return fmt.Sprintf("s%d", sequence.Inc())
}

type sinkWriter struct {
	io.Writer
}

func (w sinkWriter) Sync() error {
	return nil
}

func (w sinkWriter) Close() error {
	return nil
}

var levelMapping = map[string]logger.Level{
	"DEBUG": logger.DebugLevel,
	"INFO":  logger.InfoLevel,
	"WARN":  logger.WarnLevel,
	"ERROR": logger.ErrorLevel,
}

func unmarshalMessage(t *testing.T, line string) adaptertest.Message {
	t.Helper()

	out := unmarshalZapMessage(t, line)

	return adaptertest.Message{
		Level:          levelMapping[out.L],
		Message:        out.M,
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

func unmarshalZapMessage(t *testing.T, line string) zapMessage {
	t.Helper()

	out := zapMessage{}
	err := json.Unmarshal([]byte(line), &out)
	require.NoError(t, err)

	return out
}

type zapMessage struct {
	L              string // level
	M              string // message
	C              string // caller
	Error          string
	StringField    string
	IntField       int
	Int64Field     int64
	Float32Field   float32
	Float64Field   float64
	TimeField      time.Time
	InterfaceField adaptertest.InterfaceField
}
