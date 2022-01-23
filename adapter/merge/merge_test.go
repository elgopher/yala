package merge_test

import (
	"context"
	"testing"

	"github.com/elgopher/yala/adapter/merge"
	"github.com/elgopher/yala/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var entry = logger.Entry{Message: "message"}

func TestAdapter_Log(t *testing.T) {
	ctx := context.Background()

	t.Run("should not panic when zero value", func(t *testing.T) {
		var adapter merge.Adapter
		adapter.Log(ctx, entry)
	})

	t.Run("should delegate to wrapped adapter without merging when MergeFunc is nil", func(t *testing.T) {
		wrappedAdapter := &adapterMock{}
		adapter := merge.Adapter{
			Adapter:   wrappedAdapter,
			MergeFunc: nil,
		}
		// when
		adapter.Log(ctx, entry)
		// then
		wrappedAdapter.assertHasOneEntry(t, entry)
	})

	t.Run("should pass a merged entry to wrapped adapter", func(t *testing.T) {
		wrappedAdapter := &adapterMock{}

		expectedEntry := entry
		expectedEntry.Message = "altered"

		adapter := merge.Adapter{
			Adapter: wrappedAdapter,
			MergeFunc: func(ctx context.Context, entry logger.Entry) logger.Entry {
				entry.Message = "altered"

				return entry
			},
		}
		// when
		adapter.Log(ctx, entry)
		// then
		wrappedAdapter.assertHasOneEntry(t, expectedEntry)
	})
}

type adapterMock struct {
	entries []logger.Entry
}

func (a *adapterMock) Log(ctx context.Context, entry logger.Entry) {
	a.entries = append(a.entries, entry)
}

func (a *adapterMock) assertHasOneEntry(t *testing.T, expected logger.Entry) {
	t.Helper()

	require.Len(t, a.entries, 1)
	assert.Equal(t, expected, a.entries[0])
}
