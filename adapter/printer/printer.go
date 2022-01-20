package printer

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jacekolszak/yala/logger"
)

// Adapter is a logger.Adapter implementation, which is using Printer interface. This interface is implemented for
// example by log.Logger from the Go standard library.
//
// This adapter prints fields and error in logfmt format.
type Adapter struct {
	Printer Printer
}

type Printer interface {
	Println(...interface{})
}

// StdErrorAdapter returns a logger.Adapter implementation which prints log messages to stderr using `fmt` package.
func StdErrorAdapter() Adapter {
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
	_, _ = fmt.Fprintln(p.Writer, args...)
}

func (f Adapter) Log(ctx context.Context, entry logger.Entry) {
	if f.Printer == nil {
		return
	}

	var builder strings.Builder

	builder.WriteString(string(entry.Level))
	builder.WriteRune(' ')
	builder.WriteString(entry.Message)

	if len(entry.Fields) > 0 {
		builder.WriteRune(' ')
		writeFields(&builder, entry.Fields)
	}

	if entry.Error != nil {
		builder.WriteString(" error=")
		builder.WriteString(fmt.Sprintf("%s", entry.Error))
	}

	f.Printer.Println(builder.String())
}

func writeFields(builder *strings.Builder, fields []logger.Field) {
	for i, f := range fields {
		builder.WriteString(f.Key)
		builder.WriteRune('=')
		builder.WriteString(fmt.Sprintf("%s", f.Value))

		notLast := i < len(fields)-1
		if notLast {
			builder.WriteRune(' ')
		}
	}
}
