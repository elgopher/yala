package printer

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jacekolszak/yala/adapter/logfmt"
	"github.com/jacekolszak/yala/logger"
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

// StderrAdapter returns a logger.Adapter implementation which prints log messages to stderr using `fmt` package.
func StderrAdapter() Adapter {
	return Adapter{Printer: WriterPrinter{os.Stderr}}
}

// StdoutAdapter returns a logger.Adapter implementation which prints log messages to stdout using `fmt` package.
func StdoutAdapter() Adapter {
	return Adapter{Printer: WriterPrinter{os.Stdout}}
}

type WriterPrinter struct {
	io.Writer
}

func (p WriterPrinter) Println(args ...interface{}) {
	if p.Writer == nil {
		return
	}

	_, _ = fmt.Fprintln(p.Writer, args...)
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
