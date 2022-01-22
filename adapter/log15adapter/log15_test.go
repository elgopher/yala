package log15adapter_test

import (
	"context"
	"testing"

	"github.com/inconshreveable/log15"
	"github.com/jacekolszak/yala/adapter/log15adapter"
	"github.com/jacekolszak/yala/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const message = "message"

func TestAdapter_Log(t *testing.T) {
	ctx := context.Background()
	err := stringError("err message")

	t.Run("should log message with proper severity level", func(t *testing.T) {
		tests := map[logger.Level]log15.Lvl{
			logger.DebugLevel: log15.LvlDebug,
			logger.InfoLevel:  log15.LvlInfo,
			logger.WarnLevel:  log15.LvlWarn,
			logger.ErrorLevel: log15.LvlError,
		}

		for level, log15level := range tests {
			t.Run(string(level), func(t *testing.T) {
				log15Logger := log15.New()
				handler := &handlerMock{}
				log15Logger.SetHandler(handler)
				adapter := log15adapter.Adapter{Logger: log15Logger}
				// when
				adapter.Log(ctx, logger.Entry{Level: level, Message: message})
				// then
				require.Len(t, handler.records, 1)
				record := handler.records[0]
				assert.Equal(t, log15level, record.Lvl)
				assert.Equal(t, message, record.Msg)
			})
		}
	})

	t.Run("should log message with fields and error", func(t *testing.T) {
		tests := map[string]struct {
			entry          logger.Entry
			expectedRecord log15.Record
		}{
			"fields": {
				entry: logger.Entry{
					Level:   logger.InfoLevel,
					Message: message,
					Fields:  []logger.Field{{Key: "k", Value: "v"}},
				},
				expectedRecord: log15.Record{
					Lvl: log15.LvlInfo,
					Msg: message,
					Ctx: []interface{}{"k", "v"},
				},
			},
			"error": {
				entry: logger.Entry{
					Level:   logger.ErrorLevel,
					Message: message,
					Error:   err,
				},
				expectedRecord: log15.Record{
					Lvl: log15.LvlError,
					Msg: message,
					Ctx: []interface{}{"error", err},
				},
			},
			"fields and error": {
				entry: logger.Entry{
					Level:   logger.ErrorLevel,
					Message: message,
					Fields:  []logger.Field{{Key: "k", Value: "v"}},
					Error:   stringError("err message"),
				},
				expectedRecord: log15.Record{
					Lvl: log15.LvlError,
					Msg: message,
					Ctx: []interface{}{"k", "v", "error", err},
				},
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				log15Logger := log15.New()
				handler := &handlerMock{}
				log15Logger.SetHandler(handler)
				adapter := log15adapter.Adapter{Logger: log15Logger}
				// when
				adapter.Log(ctx, test.entry)
				// then
				require.Len(t, handler.records, 1)
				assert.ElementsMatch(t, test.expectedRecord.Ctx, handler.records[0].Ctx)
			})
		}
	})

	t.Run("should not panic when logger is nil", func(t *testing.T) {
		adapter := log15adapter.Adapter{Logger: nil}
		assert.NotPanics(t, func() {
			adapter.Log(ctx, logger.Entry{
				Level:   logger.InfoLevel,
				Message: message,
			})
		})
	})
}

type handlerMock struct {
	records []*log15.Record
}

func (h *handlerMock) Log(r *log15.Record) error {
	h.records = append(h.records, r)

	return nil
}

type stringError string

func (e stringError) Error() string {
	return string(e)
}
