package glogadapter_test

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/elgopher/yala/adapter/glogadapter"
	"github.com/elgopher/yala/adapter/internal/fake"
	"github.com/elgopher/yala/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdapter_Log(t *testing.T) {
	const message = "message"

	t.Run("should log caller", func(t *testing.T) {
		stderr := fake.UseFakeStderr(t)
		defer stderr.Release()

		adapter := glogadapter.Adapter{}
		// when
		adapter.Log(context.Background(), logger.Entry{
			Level:   logger.ErrorLevel,
			Message: message,
		})
		// then
		msg := unmarshalLine(t, stderr.String(t))
		const expectedPrefix = "glog_test.go:"
		assert.Truef(t,
			strings.HasPrefix(msg.caller, expectedPrefix),
			"caller %s has no prefix %s", msg.caller, expectedPrefix)
	})

	t.Run("should log message", func(t *testing.T) {
		stderr := fake.UseFakeStderr(t)
		defer stderr.Release()

		adapter := glogadapter.Adapter{}
		// when
		adapter.Log(context.Background(), logger.Entry{
			Level:   logger.ErrorLevel,
			Message: message,
		})
		// then
		msg := unmarshalLine(t, stderr.String(t))
		assert.Equal(t, "E", msg.level)
		assert.Equal(t, message, msg.message)
	})

	t.Run("should log message with field", func(t *testing.T) {
		stderr := fake.UseFakeStderr(t)
		defer stderr.Release()

		adapter := glogadapter.Adapter{}
		// when
		entry := logger.Entry{
			Level:   logger.ErrorLevel,
			Message: message,
		}.With(logger.Field{
			Key:   "k",
			Value: "v",
		})
		adapter.Log(context.Background(), entry)
		// then
		msg := unmarshalLine(t, stderr.String(t))
		assert.Equal(t, "k=v", msg.fields)
		assert.Equal(t, message, msg.message)
	})
}

func unmarshalLine(t *testing.T, line string) glogMessage {
	t.Helper()

	r := regexp.MustCompile(`\s+`)
	parts := r.Split(line, 6)
	require.Len(t, parts, 6)

	level := parts[0]
	level = level[:1]

	caller := parts[3]
	caller = caller[:len(caller)-1]

	message := parts[4]

	fields := parts[5]
	if len(fields) > 0 {
		fields = fields[:len(fields)-1]
	}

	return glogMessage{
		caller:  caller,
		level:   level,
		message: message,
		fields:  fields,
	}
}

type glogMessage struct {
	caller  string
	level   string
	message string
	fields  string
}
