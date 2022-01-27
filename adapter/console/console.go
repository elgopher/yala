// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package console

import (
	"fmt"
	"io"
	"os"

	"github.com/elgopher/yala/adapter/printer"
	"github.com/elgopher/yala/logger"
)

// StdoutAdapter returns a logger.Adapter implementation which prints log messages to stdout.
func StdoutAdapter() logger.Adapter { // nolint
	return printer.Adapter{Printer: WriterPrinter{os.Stdout}}
}

// StderrAdapter returns a logger.Adapter implementation which prints log messages to stderr.
func StderrAdapter() logger.Adapter { // nolint
	return printer.Adapter{Printer: WriterPrinter{os.Stderr}}
}

type WriterPrinter struct {
	Writer io.Writer
}

func (p WriterPrinter) Println(args ...interface{}) {
	if p.Writer == nil {
		return
	}

	_, _ = fmt.Fprintln(p.Writer, args...)
}
