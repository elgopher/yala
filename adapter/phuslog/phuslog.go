package phuslog

import (
	"context"

	"github.com/jacekolszak/yala/logger"
	"github.com/phuslu/log"
)

// Adapter is a logger.Adapter implementation, which is using `phuslog` module (https://github.com/phuslu/log).
type Adapter struct {
	Logger log.Logger
}

func (l Adapter) Log(ctx context.Context, entry logger.Entry) {
	phuslogEntry := l.Logger.Log()
	phuslogEntry.Level = phuslogLevel(entry.Level)

	for _, field := range entry.Fields {
		phuslogEntry = phuslogEntry.Interface(field.Key, field.Value)
	}

	if entry.Error != nil {
		phuslogEntry = phuslogEntry.Err(entry.Error)
	}

	phuslogEntry.Msg(entry.Message)
}

func phuslogLevel(level logger.Level) log.Level {
	switch level {
	case logger.DebugLevel:
		return log.DebugLevel
	case logger.InfoLevel:
		return log.InfoLevel
	case logger.WarnLevel:
		return log.WarnLevel
	case logger.ErrorLevel:
		return log.ErrorLevel
	default:
		return log.InfoLevel
	}
}
