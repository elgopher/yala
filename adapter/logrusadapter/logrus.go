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
	Logger LogrusLogger
}

// LogrusLogger is either *logrus.Logger or *logrus.Entry.
type LogrusLogger interface {
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(logrus.Fields) *logrus.Entry
	WithError(error) *logrus.Entry
	Log(lvl logrus.Level, args ...interface{})
}

// Log logs the entry using logrus module.
func (a Adapter) Log(ctx context.Context, entry logger.Entry) {
	if a.Logger == nil {
		return
	}

	logrusLogger := loggerWithFields(a.Logger, entry)
	logrusLogger.Log(logrusLevel(entry), entry.Message)
}

func loggerWithFields(logrusLogger LogrusLogger, entry logger.Entry) LogrusLogger { // nolint:ireturn
	length := len(entry.Fields)
	if entry.Error != nil {
		length++
	}

	if length == 0 {
		return logrusLogger
	}

	fields := logrus.Fields{}
	for _, field := range entry.Fields {
		fields[field.Key] = field.Value
	}

	if entry.Error != nil {
		fields[logrus.ErrorKey] = entry.Error
	}

	return logrusLogger.WithFields(fields)
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
