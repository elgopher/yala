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

// WriterPrinter implements printer.Printer by adapting io.Writer. Should be used with care, because it discards all
// errors returned during writing.
type WriterPrinter struct {
	Writer io.Writer
}

// Println prints the msg using io.Writer. Errors are discarded.
func (p WriterPrinter) Println(skipCallerFrames int, msg string) {
	if p.Writer == nil {
		return
	}

	_, _ = fmt.Fprintln(p.Writer, msg)
}
