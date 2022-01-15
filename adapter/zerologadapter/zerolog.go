package zerologadapter

import (
	"context"

	"github.com/jacekolszak/yala/logger"
	"github.com/rs/zerolog"
)

type Service struct {
	zerolog.Logger
}

func (l Service) Log(ctx context.Context, entry logger.Entry) {
	event := l.WithLevel(convertLevel(entry.Level))

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
	case logger.ErrorLevel:
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
