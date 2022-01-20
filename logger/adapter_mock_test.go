package logger_test

import (
	"context"
	"testing"

	"github.com/jacekolszak/yala/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type adapterMock struct {
	entries []logger.Entry
}

func (a *adapterMock) Log(_ context.Context, entry logger.Entry) {
	a.entries = append(a.entries, entry)
}

func (a *adapterMock) HasExactlyOneEntry(t *testing.T, expected logger.Entry) {
	t.Helper()

	require.Len(t, a.entries, 1)
	actual := a.entries[0]
	assert.Equal(t, expected, actual)
}

func (a *adapterMock) HasExactlyOneEntryWithFields(t *testing.T, expected []logger.Field) {
	t.Helper()

	require.Len(t, a.entries, 1)
	actual := a.entries[0].Fields
	assert.Equal(t, expected, actual)
}

func (a *adapterMock) HasExactlyOneEntryWithError(t *testing.T, expected error) {
	t.Helper()

	require.Len(t, a.entries, 1)
	actual := a.entries[0].Error
	assert.Equal(t, expected, actual)
}
