package glogadapter

import (
	"context"

	"github.com/golang/glog"
	"github.com/jacekolszak/yala/logger"
)

// Service is a logger.Service implementation, which is using `glog` package (https://github.com/golang/glog).
type Service struct{}

type log func(args ...interface{})

func (s Service) Log(_ context.Context, entry logger.Entry) {
	var logMessage log

	switch entry.Level {
	case logger.DebugLevel:
		logMessage = glog.Infoln
	case logger.InfoLevel:
		logMessage = glog.Infoln
	case logger.ErrorLevel:
		logMessage = glog.Errorln
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
