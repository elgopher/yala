package logadapter

import (
	"log"

	"github.com/elgopher/yala/adapter/printer"
	"github.com/elgopher/yala/logger"
)

func Adapter(l *log.Logger) logger.Adapter { // nolint
	return printer.Adapter{Printer: l}
}
