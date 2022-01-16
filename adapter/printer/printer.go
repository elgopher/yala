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

	fieldsAsString := fieldsToString(entry.Fields)
	if entry.Error == nil {
		f.Printer.Println(entry.Level, entry.Message, fieldsAsString)
	} else {
		f.Printer.Println(entry.Level, entry.Message, fieldsAsString, "error:", entry.Error)
	}
}

func fieldsToString(fields []logger.Field) string {
	var b strings.Builder

	for i, f := range fields {
		b.WriteString(f.Key)
		b.WriteRune('=')
		b.WriteString(fmt.Sprintf("%s", f.Value))

		notLast := i < len(fields)-1
		if notLast {
			b.WriteRune(',')
		}
	}

	return b.String()
}
