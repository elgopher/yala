package printer

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jacekolszak/yala/logger"
)

// Service is a logger.Service implementation, which is using Printer interface. This interface is implemented for
// example by log.Logger from the Go standard library.
type Service struct {
	Printer
}

type Printer interface {
	Println(...interface{})
}

// StdErrorService returns a logger.Service implementation which prints log messages to stderr using `fmt` package.
func StdErrorService() Service {
	return Service{Printer: WriterPrinter{os.Stderr}}
}

// StdoutService returns a logger.Service implementation which prints log messages to stdout using `fmt` package.
func StdoutService() Service {
	return Service{Printer: WriterPrinter{os.Stdout}}
}

type WriterPrinter struct {
	io.Writer
}

func (p WriterPrinter) Println(args ...interface{}) {
	_, _ = fmt.Fprintln(p.Writer, args...)
}

func (f Service) Log(ctx context.Context, entry logger.Entry) {
	if f.Printer == nil {
		return
	}

	fieldsAsString := fieldsToString(entry.Fields)
	if entry.Error == nil {
		f.Printer.Println(entry.Level, entry.Message, fieldsAsString)
	} else {
		f.Printer.Println(entry.Level, entry.Message, fieldsAsString, entry.Error)
	}
}

func fieldsToString(fields []logger.Field) string {
	var b strings.Builder

	for _, f := range fields {
		b.WriteString(f.Key)
		b.WriteRune('=')
		b.WriteString(fmt.Sprintf("%s", f.Value))
	}

	return b.String()
}
