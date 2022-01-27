package logadapter

import (
	"log"

	"github.com/elgopher/yala/adapter/printer"
	"github.com/elgopher/yala/logger"
)

func Adapter(l *log.Logger) logger.Adapter { // nolint
	return printer.Adapter{Printer: printerLogger{l}}
}

type printerLogger struct {
	*log.Logger
}

func (p printerLogger) Println(skipCallerFrames int, msg string) {
	_ = p.Logger.Output(skipCallerFrames+1, msg)
}
