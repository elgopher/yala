package glogadapter_test

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/elgopher/yala/adapter/glogadapter"
	"github.com/elgopher/yala/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdapter_Log(t *testing.T) {
	t.Run("should log caller", func(t *testing.T) {
		stderr := useFakeStderr(t)
		defer stderr.Release()

		adapter := glogadapter.Adapter{}
		// when
		adapter.Log(context.Background(), logger.Entry{
			Level:   logger.ErrorLevel,
			Message: "message",
		})
		// then
		msg := unmarshalLine(t, stderr.FirstLine(t))
		const expectedPrefix = "glog_test.go:"
		assert.Truef(t,
			strings.HasPrefix(msg.caller, expectedPrefix),
			"caller %s has no prefix %s", msg.caller, expectedPrefix)
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

	return glogMessage{
		caller:  caller,
		level:   level,
		message: parts[4],
		fields:  parts[5],
	}
}

type glogMessage struct {
	caller  string
	level   string
	message string
	fields  string
}
