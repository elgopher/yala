package main

import (
	"context"

	"github.com/elgopher/yala/adapter/console"
	"github.com/elgopher/yala/logger"
)

// This advanced example shows how to filter logs by level (if the library of your choice does not support that,
// or you want to use different filtering for different loggers).
func main() {
	adapter := console.StdoutAdapter()

	// creates an adapter which filters logs by level
	filterAdapter := FilterByLevel{
		MinLevel:    logger.WarnLevel,
		NextAdapter: adapter,
	}
	l := logger.Local{Adapter: filterAdapter}

	ctx := context.Background()

	l.Info(ctx, "will be filtered out")
	l.Warn(ctx, "will be logged")
	l.Error(ctx, "this too will be logged")
}

// FilterByLevel is a middleware (decorator) which filters out entries
// with level less severe than MinLevel
type FilterByLevel struct {
	MinLevel    logger.Level
	NextAdapter logger.Adapter
}

func (a FilterByLevel) Log(ctx context.Context, entry logger.Entry) {
	// It is always better to use MoreSevereThan method instead of directly comparing integers,
	// because it hides out how levels are sorted
	if a.MinLevel.MoreSevereThan(entry.Level) {
		return
	}

	entry.SkippedCallerFrames++ // each middleware adapter must additionally skip one frame
	a.NextAdapter.Log(ctx, entry)
}
