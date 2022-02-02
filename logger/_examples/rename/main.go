package main

import (
	"context"

	"github.com/elgopher/yala/adapter/console"
	"github.com/elgopher/yala/logger"
)

// This example shows how to rename all fields matching given string.
func main() {
	ctx := context.Background()

	adapter := RenameFieldsAdapter{
		From:        "this",
		To:          "that",
		NextAdapter: console.StdoutAdapter(),
	}

	log := logger.WithAdapter(adapter)

	log.With("this", "value").
		Info(ctx, "this field will be replaced with that")
}

// RenameFieldsAdapter is a middleware (decorator) renaming all fields equal to From into To.
type RenameFieldsAdapter struct {
	From, To    string
	NextAdapter logger.Adapter
}

func (r RenameFieldsAdapter) Log(ctx context.Context, entry logger.Entry) {
	fields := make([]logger.Field, len(entry.Fields)) // Create a new slice in order to be concurrency-safe

	for i, field := range entry.Fields {
		key := field.Key
		if key == r.From {
			key = r.To
		}

		fields[i] = logger.Field{Key: key, Value: field.Value}
	}

	entry.Fields = fields
	entry.SkippedCallerFrames++ // each middleware adapter must additionally skip one frame

	r.NextAdapter.Log(ctx, entry)
}
