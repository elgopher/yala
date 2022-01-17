package glogadapter

import (
	"context"

	"github.com/golang/glog"
	"github.com/jacekolszak/yala/logger"
)

// Adapter is a logger.Adapter implementation, which is using `glog` package (https://github.com/golang/glog).
type Adapter struct{}

type log func(args ...interface{})

func (s Adapter) Log(_ context.Context, entry logger.Entry) {
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

	var args []interface{}
	args = append(args, entry.Message)

	if len(entry.Fields) > 0 {
		args = append(args, "fields:", entry.Fields)
	}

	if entry.Error != nil {
		args = append(args, "error:", entry.Error)
	}

	logMessage(args...)
}
