package zerologadapter

import (
	"context"

	"github.com/jacekolszak/yala/logger"
	"github.com/rs/zerolog"
)

// Adapter is a logger.Adapter implementation, which is using `zerolog` module (https://github.com/rs/zerolog).
type Adapter struct {
	Logger zerolog.Logger
}

func (l Adapter) Log(ctx context.Context, entry logger.Entry) {
	event := l.Logger.WithLevel(convertLevel(entry.Level))

	for _, field := range entry.Fields {
		event = event.Interface(field.Key, field.Value)
	}

	if entry.Error != nil {
		event = event.Err(entry.Error)
	}

	event.Msg(entry.Message)
}

func convertLevel(level logger.Level) zerolog.Level {
	switch level {
	case logger.DebugLevel:
		return zerolog.DebugLevel
	case logger.InfoLevel:
		return zerolog.InfoLevel
	case logger.WarnLevel:
		return zerolog.WarnLevel
	case logger.ErrorLevel:
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
