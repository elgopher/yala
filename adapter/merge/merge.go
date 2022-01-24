// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package merge provides logger.Adapter implementation which merges each logger.Entry with context.Context using
// provided function. It can be used for adding tags stored in the context.Context, adding custom fields or modifying
// messages.
package merge

import (
	"context"

	"github.com/elgopher/yala/logger"
)

type Func func(ctx context.Context, entry logger.Entry) logger.Entry

type Adapter struct {
	// Adapter will be executed once entry is merged
	Adapter   logger.Adapter
	MergeFunc Func
}

func (a Adapter) Log(ctx context.Context, entry logger.Entry) {
	if a.Adapter == nil {
		return
	}

	if a.MergeFunc != nil {
		entry = a.MergeFunc(ctx, entry)
	}

	a.Adapter.Log(ctx, entry)
}
