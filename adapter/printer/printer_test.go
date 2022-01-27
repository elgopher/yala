// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package printer_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/elgopher/yala/adapter/printer"
	"github.com/elgopher/yala/logger"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

const message = "message"

func TestAdapter_Log(t *testing.T) {
	tests := map[string]struct {
		entry           logger.Entry
		expectedMessage string
	}{
		"message debug": {
			entry: logger.Entry{
				Level:   logger.DebugLevel,
				Message: message,
			},
			expectedMessage: "DEBUG message\n",
		},
		"message warn": {
			entry: logger.Entry{
				Level:   logger.WarnLevel,
				Message: message,
			},
			expectedMessage: "WARN message\n",
		},
		"fields": {
			entry: logger.Entry{
				Level:   logger.InfoLevel,
				Message: message,
				Fields:  []logger.Field{{Key: "k", Value: "v"}},
			},
			expectedMessage: "INFO message k=v\n",
		},
		"error": {
			entry: logger.Entry{
				Level:   logger.ErrorLevel,
				Message: message,
				Error:   stringError("err message"),
			},
			expectedMessage: "ERROR message error=\"err message\"\n",
		},
		"fields and error": {
			entry: logger.Entry{
				Level:   logger.ErrorLevel,
				Message: message,
				Fields:  []logger.Field{{Key: "k", Value: "v"}},
				Error:   stringError("err message"),
			},
			expectedMessage: "ERROR message k=v error=\"err message\"\n",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var actual strings.Builder
			adapter := printer.Adapter{Printer: stringPrinter{&actual}}
			// when
			adapter.Log(ctx, test.entry)
			// then
			assert.Equal(t, test.expectedMessage, actual.String())
		})
	}

	t.Run("should not panic when printer is nil", func(t *testing.T) {
		adapter := printer.Adapter{Printer: nil}
		assert.NotPanics(t, func() {
			adapter.Log(ctx, logger.Entry{
				Level:   logger.InfoLevel,
				Message: message,
			})
		})
	})
}

type stringError string

func (e stringError) Error() string {
	return string(e)
}

type stringPrinter struct {
	io.StringWriter
}

func (p stringPrinter) Println(i ...interface{}) {
	s := fmt.Sprintln(i...)
	_, _ = p.WriteString(s)
}
