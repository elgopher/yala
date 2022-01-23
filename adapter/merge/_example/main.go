package main

import (
	"context"

	"github.com/elgopher/yala/adapter/merge"
	"github.com/elgopher/yala/adapter/printer"
	"github.com/elgopher/yala/logger"
)

const tag = "tag"

// This advanced example shows how to log messages with additional field taken from context.Context
func main() {
	adapter := printer.StdoutAdapter()

	// creates an adapter which adds field from context to each logged message.
	mergeAdapter := merge.Adapter{Adapter: adapter, MergeFunc: addFieldFromContext}
	l := logger.Local(mergeAdapter)

	ctx := context.Background()
	// add tag to context
	ctx = context.WithValue(ctx, tag, "value")

	l.Info(ctx, "tagged message")                // INFO tagged message tag=value
	l.With(ctx, "k", "v").Info("tagged message") // INFO tagged message k=v tag=value
}

func addFieldFromContext(ctx context.Context, entry logger.Entry) logger.Entry {
	tagValue := ctx.Value(tag)

	return entry.With(
		logger.Field{
			Key:   tag,
			Value: tagValue,
		})
}
