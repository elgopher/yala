package logrusadapter

import (
	"context"

	"github.com/jacekolszak/yala/logger"
	"github.com/sirupsen/logrus"
)

type Service struct {
	*logrus.Entry
}

func (s Service) Log(ctx context.Context, entry logger.Entry) {
	if s.Entry == nil {
		return
	}

	lvl := logrus.InfoLevel

	switch entry.Level {
	case logger.DebugLevel:
		lvl = logrus.DebugLevel
	case logger.InfoLevel:
		lvl = logrus.InfoLevel
	case logger.ErrorLevel:
		lvl = logrus.ErrorLevel
	}

	logrusEntry := s.Entry

	for _, f := range entry.Fields {
		logrusEntry = logrusEntry.WithField(f.Key, f.Value)
	}

	if entry.Error != nil {
		logrusEntry = logrusEntry.WithError(entry.Error)
	}

	logrusEntry.Log(lvl, entry.Message)
}
