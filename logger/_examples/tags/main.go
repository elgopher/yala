package main

import (
	"context"

	"github.com/elgopher/yala/adapter/console"
	"github.com/elgopher/yala/logger"
)

const tag = "tag"

// This advanced example shows how to log messages with additional field taken from context.Context
func main() {
	adapter := console.StdoutAdapter()

	// creates an adapter which adds field from context to each logged message.
	addFieldAdapter := AddFieldFromContextAdapter{NextAdapter: adapter}
	log := logger.WithAdapter(addFieldAdapter)

	ctx := context.Background()
	// add tag to context
	ctx = context.WithValue(ctx, tag, "value")

	// The chain of execution will look like this:
	// log.Info() -> AddFieldFromContextAdapter -> console adapter
	log.Info(ctx, "tagged message") // INFO tagged message tag=value

	log.With("k", "v").
		Info(ctx, "tagged message") // INFO tagged message k=v tag=value
}

// AddFieldFromContextAdapter is a middleware (decorator) which adds
// a new field to logger.Entry from the tag stored in the context.
type AddFieldFromContextAdapter struct {
	NextAdapter logger.Adapter
}

func (a AddFieldFromContextAdapter) Log(ctx context.Context, entry logger.Entry) {
	entry.Fields = append(entry.Fields,
		logger.Field{
			Key:   tag,
			Value: ctx.Value(tag),
		},
	)
	entry.SkippedCallerFrames++ // each middleware adapter must additionally skip one frame
	a.NextAdapter.Log(ctx, entry)
}
