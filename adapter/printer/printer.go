// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package printer

import (
	"context"
	"strings"

	"github.com/elgopher/yala/adapter/logfmt"
	"github.com/elgopher/yala/logger"
)

// Adapter is a logger.Adapter implementation, which is using Printer interface. This interface is implemented for
// example by log.Logger from the Go standard library.
//
// This adapter prints fields and error in logfmt format. For example:
//
// 		message key=value error=message
type Adapter struct {
	Printer Printer
}

type Printer interface {
	Println(...interface{})
}

func (f Adapter) Log(ctx context.Context, entry logger.Entry) {
	if f.Printer == nil {
		return
	}

	var builder strings.Builder

	builder.WriteString(string(entry.Level))
	builder.WriteByte(' ')
	builder.WriteString(entry.Message)

	if len(entry.Fields) > 0 {
		builder.WriteByte(' ')
		logfmt.WriteFields(&builder, entry.Fields)
	}

	if entry.Error != nil {
		builder.WriteByte(' ')
		logfmt.WriteField(&builder, logger.Field{Key: "error", Value: entry.Error})
	}

	f.Printer.Println(builder.String())
}
