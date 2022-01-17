package logrusadapter

import (
	"context"

	"github.com/jacekolszak/yala/logger"
	"github.com/sirupsen/logrus"
)

// Adapter is a logger.Adapter implementation, which is using `logrus` module (https://github.com/sirupsen/logrus).
type Adapter struct {
	Entry *logrus.Entry
}

func (a Adapter) Log(ctx context.Context, entry logger.Entry) {
	if a.Entry == nil {
		return
	}

	lvl := logrus.InfoLevel

	switch entry.Level {
	case logger.DebugLevel:
		lvl = logrus.DebugLevel
	case logger.InfoLevel:
		lvl = logrus.InfoLevel
	case logger.WarnLevel:
		lvl = logrus.WarnLevel
	case logger.ErrorLevel:
		lvl = logrus.ErrorLevel
	}

	logrusEntry := a.Entry

	for _, f := range entry.Fields {
		logrusEntry = logrusEntry.WithField(f.Key, f.Value)
	}

	if entry.Error != nil {
		logrusEntry = logrusEntry.WithError(entry.Error)
	}

	logrusEntry.Log(lvl, entry.Message)
}
