package printer_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jacekolszak/yala/adapter/printer"
	"github.com/jacekolszak/yala/logger"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

const message = "message"

func TestFormattingFields(t *testing.T) {
	t.Run("should format fields using logfmt", func(t *testing.T) {
		tests := map[string]struct {
			fields   []logger.Field
			expected string
		}{
			"one field": {
				fields:   singleField("key", "value"),
				expected: "key=value",
			},
			"two fields": {
				fields:   []logger.Field{{Key: "k1", Value: "v1"}, {Key: "k2", Value: "v2"}},
				expected: "k1=v1 k2=v2",
			},
			"nil value": {
				fields:   singleField("key", nil),
				expected: "key=nil",
			},
			"nil value and then one field": {
				fields:   []logger.Field{{Key: "k1", Value: nil}, {Key: "k2", Value: "v2"}},
				expected: "k1=nil k2=v2",
			},
			"nil string value": {
				fields:   singleField("key", "nil"),
				expected: "key=\"nil\"",
			},
			"space": {
				fields:   singleField("key", "v v"),
				expected: `key="v v"`,
			},
			"=": {
				fields:   singleField("key", "="),
				expected: `key="="`,
			},
			"space and =": {
				fields:   singleField("key", " ="),
				expected: `key=" ="`,
			},
			`"`: {
				fields:   singleField("key", `"`),
				expected: `key=\"`,
			},
			`\`: {
				fields:   singleField("key", `\`),
				expected: `key=\\`,
			},
			`\"`: {
				fields:   singleField("key", `\"`),
				expected: `key=\\\"`,
			},
			`"quoted with spaces"`: {
				fields:   singleField("k", `"quoted with spaces"`),
				expected: `k="\"quoted with spaces\""`,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				var builder strings.Builder
				adapter := printer.Adapter{
					Printer: printer.WriterPrinter{Writer: &builder},
				}
				// when
				adapter.Log(ctx, logger.Entry{
					Level:   logger.InfoLevel,
					Message: message,
					Fields:  test.fields,
				})
				// then
				expectedLine := fmt.Sprintf("INFO message %s\n", test.expected)
				assert.Equal(t, expectedLine, builder.String())
			})
		}
	})
}

func singleField(k string, v interface{}) []logger.Field {
	return []logger.Field{{Key: k, Value: v}}
}

type stringError string

func (e stringError) Error() string {
	return string(e)
}

func TestFormattingError(t *testing.T) {
	tests := map[string]struct {
		error    error
		expected string
	}{
		"nil": {
			error:    nil,
			expected: "",
		},
		"err": {
			error:    stringError("err"),
			expected: " error=err",
		},
		"space": {
			error:    stringError("some error"),
			expected: ` error="some error"`,
		},
		"=": {
			error:    stringError("="),
			expected: ` error="="`,
		},
		`"`: {
			error:    stringError(`"`),
			expected: ` error=\"`,
		},
		`\`: {
			error:    stringError(`\`),
			expected: ` error=\\`,
		},
		`\"`: {
			error:    stringError(`\"`),
			expected: ` error=\\\"`,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var builder strings.Builder
			adapter := printer.Adapter{
				Printer: printer.WriterPrinter{Writer: &builder},
			}
			// when
			adapter.Log(ctx, logger.Entry{
				Level:   logger.ErrorLevel,
				Message: message,
				Error:   test.error,
			})
			// then
			expectedLine := fmt.Sprintf("ERROR message%s\n", test.expected)
			assert.Equal(t, expectedLine, builder.String())
		})
	}
}
