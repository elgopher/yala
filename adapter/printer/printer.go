// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package printer

import (
	"context"
	"strings"

	"github.com/elgopher/yala/adapter/logfmt"
	"github.com/elgopher/yala/logger"
)

// Adapter is a logger.Adapter implementation, which is using Printer interface.
//
// This adapter prints fields and error in logfmt format. For example:
//
// 		message key=value error=message
type Adapter struct {
	Printer Printer
}

// Printer is someone who can print lines.
type Printer interface {
	// Println prints line. It can use skipCallerFrames to print information about caller.
	Println(skipCallerFrames int, msg string)
}

// Log logs the entry using Printer. Message is formatted using logfmt.
func (f Adapter) Log(ctx context.Context, entry logger.Entry) {
	if f.Printer == nil {
		return
	}

	var builder strings.Builder

	builder.WriteString(entry.Level.String())
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

	f.Printer.Println(entry.SkippedCallerFrames+1, builder.String())
}
