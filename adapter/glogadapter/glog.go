package glogadapter

import (
	"context"

	"github.com/golang/glog"
	"github.com/jacekolszak/yala/logger"
)

// Service is a logger.Service implementation, which is using `glog` package (https://github.com/golang/glog).
type Service struct{}

func (s Service) Log(_ context.Context, entry logger.Entry) {
	switch entry.Level {
	case logger.DebugLevel:
		glog.Infoln(entry.Message, entry.Fields)
	case logger.InfoLevel:
		glog.Infoln(entry.Message, entry.Fields)
	case logger.ErrorLevel:
		glog.Errorln(entry.Message, entry.Fields)
	}
}
