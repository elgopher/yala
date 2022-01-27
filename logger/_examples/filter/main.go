package main

import (
	"context"
	"strings"

	"github.com/elgopher/yala/adapter/console"
	"github.com/elgopher/yala/logger"
)

// This advanced example shows how to filter out messages starting with given prefix
func main() {
	adapter := console.StdoutAdapter()

	// creates an adapter which filters out messages
	filterAdapter := FilterOutMessages{
		Prefix:      "example:",
		NextAdapter: adapter,
	}
	l := logger.Local{Adapter: filterAdapter}

	ctx := context.Background()

	l.Info(ctx, "message without prefix")
	l.Info(ctx, "example: message which will be filtered out")
	l.Info(ctx, "another message without prefix")
}

// FilterOutMessages is a middleware (decorator) which filters out entries
// with message starting with prefix
type FilterOutMessages struct {
	Prefix      string
	NextAdapter logger.Adapter
}

func (a FilterOutMessages) Log(ctx context.Context, entry logger.Entry) {
	if strings.HasPrefix(entry.Message, a.Prefix) {
		return
	}

	entry.SkippedCallerFrames++ // each middleware adapter must additionally skip one frame
	a.NextAdapter.Log(ctx, entry)
}
