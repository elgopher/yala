package printer_test

import (
	"context"
	"strings"
	"testing"

	"github.com/jacekolszak/yala/adapter/printer"
	"github.com/jacekolszak/yala/logger"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

const message = "message"

func TestAdapter_Log(t *testing.T) {
	tests := map[string]struct {
		entry           logger.Entry
		expectedMessage string
	}{
		"message alone": {
			entry: logger.Entry{
				Level:   logger.DebugLevel,
				Message: message,
			},
			expectedMessage: "DEBUG message\n",
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
			var builder strings.Builder
			adapter := printer.Adapter{
				Printer: printer.WriterPrinter{Writer: &builder},
			}
			// when
			adapter.Log(ctx, test.entry)
			// then
			assert.Equal(t, test.expectedMessage, builder.String())
		})
	}
}

type stringError string

func (e stringError) Error() string {
	return string(e)
}
