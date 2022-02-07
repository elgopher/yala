package logger_test

import (
	"testing"

	"github.com/elgopher/yala/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntry_With(t *testing.T) {
	field1 := logger.Field{Key: "k1", Value: "v1"}
	field2 := logger.Field{Key: "k2", Value: "v2"}

	t.Run("should add field to empty entry", func(t *testing.T) {
		entry := logger.Entry{}
		newEntry := entry.With(field1)
		assert.Empty(t, entry.Fields)
		require.Len(t, newEntry.Fields, 1)
		assert.Equal(t, newEntry.Fields[0], field1)
	})

	t.Run("should add field to entry with one field", func(t *testing.T) {
		entry := logger.Entry{}.With(field1)
		newEntry := entry.With(field2)
		require.Len(t, entry.Fields, 1)
		require.Len(t, newEntry.Fields, 2)
		assert.Equal(t, entry.Fields[0], field1)
		assert.Equal(t, newEntry.Fields[0], field1)
		assert.Equal(t, newEntry.Fields[1], field2)
	})
}

func TestEntry_WithFields(t *testing.T) {
	t.Run("should copy entry when fields is nil", func(t *testing.T) {
		entry := logger.Entry{}
		newEntry := entry.WithFields(nil)
		assert.Equal(t, entry, newEntry)
	})

	t.Run("should copy entry when fields is empty", func(t *testing.T) {
		entry := logger.Entry{}
		newEntry := entry.WithFields(logger.Fields{})
		assert.Equal(t, entry, newEntry)
	})

	t.Run("should create new entry with one additional field", func(t *testing.T) {
		entry := logger.Entry{}
		// when
		newEntry := entry.WithFields(logger.Fields{
			"k": "v",
		})
		// then
		assert.Equal(t,
			[]logger.Field{{Key: "k", Value: "v"}},
			newEntry.Fields,
		)
		// and leave original entry unchanged
		assert.Empty(t, entry.Fields)
	})

	t.Run("should create new entry with two additional fields", func(t *testing.T) {
		entry := logger.Entry{}
		// when
		newEntry := entry.WithFields(logger.Fields{
			"k1": "v1",
			"k2": "v2",
		})
		// then
		assert.ElementsMatch(t,
			[]logger.Field{
				{Key: "k1", Value: "v1"},
				{Key: "k2", Value: "v2"},
			},
			newEntry.Fields,
		)
	})

	t.Run("should create new entry with existing field and one additional field", func(t *testing.T) {
		existingField := logger.Field{Key: "k1", Value: "v1"}
		entry := logger.Entry{
			Fields: []logger.Field{existingField},
		}
		// when
		newEntry := entry.WithFields(logger.Fields{
			"k2": "v2",
		})
		// then
		assert.Equal(t,
			[]logger.Field{
				existingField,
				{Key: "k2", Value: "v2"},
			},
			newEntry.Fields,
		)
	})

	t.Run("should create new entry with existing fields and one additional field", func(t *testing.T) {
		existingFields := []logger.Field{
			{Key: "k1", Value: "v1"},
			{Key: "k2", Value: "v2"},
		}
		entry := logger.Entry{
			Fields: existingFields,
		}
		// when
		newEntry := entry.WithFields(logger.Fields{
			"k3": "v3",
		})
		// then
		assert.Equal(t,
			[]logger.Field{
				existingFields[0],
				existingFields[1],
				{Key: "k3", Value: "v3"},
			},
			newEntry.Fields,
		)
	})
}

func TestLevel_MoreSevereThan(t *testing.T) {
	t.Run("should return true", func(t *testing.T) {
		assert.True(t, logger.InfoLevel.MoreSevereThan(logger.DebugLevel))
		assert.True(t, logger.WarnLevel.MoreSevereThan(logger.DebugLevel))
		assert.True(t, logger.ErrorLevel.MoreSevereThan(logger.DebugLevel))

		assert.True(t, logger.WarnLevel.MoreSevereThan(logger.InfoLevel))
		assert.True(t, logger.ErrorLevel.MoreSevereThan(logger.InfoLevel))

		assert.True(t, logger.ErrorLevel.MoreSevereThan(logger.WarnLevel))
	})

	t.Run("should return false", func(t *testing.T) {
		assert.False(t, logger.DebugLevel.MoreSevereThan(logger.DebugLevel))

		assert.False(t, logger.DebugLevel.MoreSevereThan(logger.InfoLevel))
		assert.False(t, logger.InfoLevel.MoreSevereThan(logger.InfoLevel))

		assert.False(t, logger.DebugLevel.MoreSevereThan(logger.WarnLevel))
		assert.False(t, logger.InfoLevel.MoreSevereThan(logger.WarnLevel))
		assert.False(t, logger.WarnLevel.MoreSevereThan(logger.WarnLevel))

		assert.False(t, logger.DebugLevel.MoreSevereThan(logger.ErrorLevel))
		assert.False(t, logger.InfoLevel.MoreSevereThan(logger.ErrorLevel))
		assert.False(t, logger.WarnLevel.MoreSevereThan(logger.ErrorLevel))
		assert.False(t, logger.ErrorLevel.MoreSevereThan(logger.ErrorLevel))
	})
}

func TestLevel_String(t *testing.T) {
	t.Run("should convert to string", func(t *testing.T) {
		assert.Equal(t, "DEBUG", logger.DebugLevel.String())
		assert.Equal(t, "INFO", logger.InfoLevel.String())
		assert.Equal(t, "WARN", logger.WarnLevel.String())
		assert.Equal(t, "ERROR", logger.ErrorLevel.String())
		assert.Equal(t, "10", logger.Level(10).String())
	})
}
