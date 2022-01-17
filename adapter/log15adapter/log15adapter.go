package log15adapter

import (
	"context"

	"github.com/inconshreveable/log15"
	"github.com/jacekolszak/yala/logger"
)

// Service is a logger.Adapter implementation, which is using `log15` package
// (https://github.com/inconshreveable/log15).
type Service struct {
	Logger log15.Logger
}

func (s Service) Log(ctx context.Context, entry logger.Entry) {
	if s.Logger == nil {
		return
	}

	log15Logger := s.Logger
	for _, field := range entry.Fields {
		log15Logger = log15Logger.New(field.Key, field.Value)
	}

	if entry.Error != nil {
		log15Logger = log15Logger.New("error", entry.Error)
	}

	switch entry.Level {
	case logger.DebugLevel:
		log15Logger.Debug(entry.Message)
	case logger.InfoLevel:
		log15Logger.Info(entry.Message)
	case logger.WarnLevel:
		log15Logger.Warn(entry.Message)
	case logger.ErrorLevel:
		log15Logger.Error(entry.Message)
	default:
		log15Logger.Info(entry.Message)
	}
}
