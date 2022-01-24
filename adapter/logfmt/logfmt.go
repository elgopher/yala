// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package logfmt provides functions encoding logger.Field using logfmt format, for example: "field=value".
package logfmt

import (
	"fmt"
	"strings"

	"github.com/elgopher/yala/logger"
)

func WriteField(builder *strings.Builder, f logger.Field) {
	builder.WriteString(f.Key)
	builder.WriteByte('=')
	writeValue(builder, f.Value)
}

func writeValue(builder *strings.Builder, value interface{}) {
	if value == nil {
		builder.WriteString("nil")

		return
	}

	if value == "nil" {
		builder.WriteString(`"nil"`)

		return
	}

	valueStr := fmt.Sprintf("%s", value)

	if strings.ContainsRune(valueStr, '\\') {
		valueStr = strings.ReplaceAll(valueStr, `\`, `\\`)
	}

	if strings.ContainsRune(valueStr, '"') {
		valueStr = strings.ReplaceAll(valueStr, `"`, `\"`)
	}

	requiresQuoting := false

	if strings.ContainsRune(valueStr, ' ') || strings.ContainsRune(valueStr, '=') {
		requiresQuoting = true
	}

	if requiresQuoting {
		builder.WriteByte('"')
	}

	builder.WriteString(valueStr)

	if requiresQuoting {
		builder.WriteByte('"')
	}
}

func WriteFields(builder *strings.Builder, fields []logger.Field) {
	for i, f := range fields {
		WriteField(builder, f)

		notLast := i < len(fields)-1
		if notLast {
			builder.WriteByte(' ')
		}
	}
}
