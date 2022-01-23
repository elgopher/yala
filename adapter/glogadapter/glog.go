package glogadapter

import (
	"context"
	"strings"

	"github.com/elgopher/yala/adapter/logfmt"
	"github.com/elgopher/yala/logger"
	"github.com/golang/glog"
)

// Adapter is a logger.Adapter implementation, which is using `glog` package (https://github.com/golang/glog).
type Adapter struct{}

type log func(args ...interface{})

func (a Adapter) Log(_ context.Context, entry logger.Entry) {
	var logMessage log

	switch entry.Level {
	case logger.DebugLevel:
		logMessage = glog.Infoln
	case logger.InfoLevel:
		logMessage = glog.Infoln
	case logger.WarnLevel:
		logMessage = glog.Warningln
	case logger.ErrorLevel:
		logMessage = glog.Errorln
	default:
		logMessage = glog.Infoln
	}

	var fieldsAndError strings.Builder

	if len(entry.Fields) > 0 {
		logfmt.WriteFields(&fieldsAndError, entry.Fields)
	}

	if entry.Error != nil {
		if len(entry.Fields) > 0 {
			fieldsAndError.WriteByte(' ')
		}

		logfmt.WriteField(&fieldsAndError, logger.Field{Key: "error", Value: entry.Error})
	}

	logMessage(entry.Message, fieldsAndError.String())
}
