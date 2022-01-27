// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package zerologadapter

import (
	"context"
	"runtime"
	"time"

	"github.com/elgopher/yala/logger"
	"github.com/rs/zerolog"
)

// Adapter is a logger.Adapter implementation, which is using `zerolog` module (https://github.com/rs/zerolog).
type Adapter struct {
	Logger zerolog.Logger
	// ReportCaller adds file and line fields which have caller information. Useful, but very slow.
	ReportCaller bool
}

func (l Adapter) Log(ctx context.Context, entry logger.Entry) {
	event := l.Logger.WithLevel(convertLevel(entry.Level))

	for _, field := range entry.Fields {
		event = eventWithField(event, field)
	}

	if entry.Error != nil {
		event = event.Err(entry.Error)
	}

	if l.ReportCaller {
		if _, file, line, ok := runtime.Caller(entry.SkippedCallerFrames); ok {
			event = eventWithField(event, logger.Field{Key: "file", Value: file})
			event = eventWithField(event, logger.Field{Key: "line", Value: line})
		}
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

func eventWithField(event *zerolog.Event, field logger.Field) *zerolog.Event {
	switch value := field.Value.(type) {
	case string:
		event = event.Str(field.Key, value)
	case int:
		event = event.Int(field.Key, value)
	case int64:
		event = event.Int64(field.Key, value)
	case float64:
		event = event.Float64(field.Key, value)
	case float32:
		event = event.Float32(field.Key, value)
	case time.Time:
		event = event.Time(field.Key, value)
	default:
		event = event.Interface(field.Key, field.Value)
	}

	return event
}
