// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package logrusadapter provides yala adapter which leverages logrus module (https://github.com/sirupsen/logrus).
package logrusadapter

import (
	"context"

	"github.com/elgopher/yala/logger"
	"github.com/sirupsen/logrus"
)

// Adapter is a logger.Adapter implementation, which is using `logrus` module (https://github.com/sirupsen/logrus).
type Adapter struct {
	Entry *logrus.Entry
}

// Log logs the entry using logrus module.
func (a Adapter) Log(ctx context.Context, entry logger.Entry) {
	if a.Entry == nil {
		return
	}

	logrusEntry := a.Entry

	for _, f := range entry.Fields {
		logrusEntry = logrusEntry.WithField(f.Key, f.Value)
	}

	if entry.Error != nil {
		logrusEntry = logrusEntry.WithError(entry.Error)
	}

	logrusEntry.Log(logrusLevel(entry), entry.Message)
}

func logrusLevel(entry logger.Entry) logrus.Level {
	switch entry.Level {
	case logger.DebugLevel:
		return logrus.DebugLevel
	case logger.InfoLevel:
		return logrus.InfoLevel
	case logger.WarnLevel:
		return logrus.WarnLevel
	case logger.ErrorLevel:
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}
