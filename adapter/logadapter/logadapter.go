package logadapter

import (
	"context"
	"log"

	"github.com/elgopher/yala/adapter/printer"
	"github.com/elgopher/yala/logger"
)

func Adapter(l *log.Logger) logger.Adapter {
	if l == nil {
		return noopAdapter{}
	}

	return printer.Adapter{Printer: printerLogger{l}}
}

type printerLogger struct {
	*log.Logger
}

func (p printerLogger) Println(skipCallerFrames int, msg string) {
	_ = p.Logger.Output(skipCallerFrames+2, msg) //nolint
}

type noopAdapter struct{}

func (n noopAdapter) Log(context.Context, logger.Entry) {}
