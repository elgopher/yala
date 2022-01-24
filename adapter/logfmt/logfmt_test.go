// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logfmt_test

import (
	"strings"
	"testing"

	"github.com/elgopher/yala/adapter/logfmt"
	"github.com/elgopher/yala/logger"
	"github.com/stretchr/testify/assert"
)

func TestWriteField(t *testing.T) {
	t.Run("should format field using logfmt", func(t *testing.T) {
		tests := map[string]struct {
			field    logger.Field
			expected string
		}{
			"one field": {
				field:    field("key", "value"),
				expected: "key=value",
			},
			"nil value": {
				field:    field("key", nil),
				expected: "key=nil",
			},
			"nil string value": {
				field:    field("key", "nil"),
				expected: "key=\"nil\"",
			},
			"space": {
				field:    field("key", "v v"),
				expected: `key="v v"`,
			},
			"=": {
				field:    field("key", "="),
				expected: `key="="`,
			},
			"space and =": {
				field:    field("key", " ="),
				expected: `key=" ="`,
			},
			`"`: {
				field:    field("key", `"`),
				expected: `key=\"`,
			},
			`\`: {
				field:    field("key", `\`),
				expected: `key=\\`,
			},
			`\"`: {
				field:    field("key", `\"`),
				expected: `key=\\\"`,
			},
			`"quoted with spaces"`: {
				field:    field("k", `"quoted with spaces"`),
				expected: `k="\"quoted with spaces\""`,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				var builder strings.Builder
				// when
				logfmt.WriteField(&builder, test.field)
				// then
				assert.Equal(t, test.expected, builder.String())
			})
		}
	})
}

func field(k string, v interface{}) logger.Field {
	return logger.Field{Key: k, Value: v}
}

func TestWriteFields(t *testing.T) {
	t.Run("should format fields using logfmt", func(t *testing.T) {
		tests := map[string]struct {
			fields   []logger.Field
			expected string
		}{
			"single field": {
				fields:   []logger.Field{field("k", "v")},
				expected: "k=v",
			},
			"two fields": {
				fields:   []logger.Field{field("k1", "v1"), field("k2", "v2")},
				expected: "k1=v1 k2=v2",
			},
			"three fields": {
				fields:   []logger.Field{field("k1", "v1"), field("k2", "v2"), field("k3", "v3")},
				expected: "k1=v1 k2=v2 k3=v3",
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				var builder strings.Builder
				// when
				logfmt.WriteFields(&builder, test.fields)
				// then
				assert.Equal(t, test.expected, builder.String())
			})
		}
	})
}
