package glogadapter

import (
	"context"

	"github.com/golang/glog"
	"github.com/jacekolszak/yala/logger"
)

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
