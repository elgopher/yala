package main

import (
	"context"

	"github.com/elgopher/yala/adapter/printer"
	"github.com/elgopher/yala/logger"
)

const tag = "tag"

// This advanced example shows how to log messages with additional field taken from context.Context
func main() {
	adapter := printer.StdoutAdapter()

	// creates an adapter which adds field from context to each logged message.
	addFieldAdapter := AddFieldFromContextAdapter{Adapter: adapter}
	l := logger.Local(addFieldAdapter)

	ctx := context.Background()
	// add tag to context
	ctx = context.WithValue(ctx, tag, "value")

	l.Info(ctx, "tagged message")                // INFO tagged message tag=value
	l.With(ctx, "k", "v").Info("tagged message") // INFO tagged message k=v tag=value
}

// AddFieldFromContextAdapter is a middleware (decorator) which adds
// a new field to logger.Entry from the tag stored in the context.
type AddFieldFromContextAdapter struct {
	Adapter logger.Adapter
}

func (a AddFieldFromContextAdapter) Log(ctx context.Context, entry logger.Entry) {
	// entry.With creates an entry and adds a new field to it
	newEntry := entry.With(
		logger.Field{
			Key:   tag,
			Value: ctx.Value(tag),
		},
	)
	a.Adapter.Log(ctx, newEntry)
}
