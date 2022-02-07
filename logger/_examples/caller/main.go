package main

import (
	"context"
	"os"
	"runtime"

	"github.com/elgopher/yala/adapter/zerologadapter"
	"github.com/elgopher/yala/logger"
	"github.com/rs/zerolog"
)

// This example shows how to add caller information to each message in a form of two new fields - `file` and `line`.
func main() {
	ctx := context.Background()

	zerologAdapter := zerologadapter.Adapter{Logger: zerolog.New(os.Stdout)}

	// create middleware (decorator) adapter which adds caller information to each message,
	// before reaching zerolog module:
	adapter := ReportCallerAdapter{NextAdapter: zerologAdapter}

	log := logger.WithAdapter(adapter)

	// The chain of execution will look like this:
	// log.Info() -> ReportCallerAdapter -> ZerologAdapter -> zerolog
	log.Info(ctx, "Message with file and line fields")
}

// ReportCallerAdapter is a middleware (decorator) adapter which adds caller information to each logged message. It can
// be very useful if the selected logging library does not support reporting caller at all (such as zerolog), or is
// somehow limited (like in logrus, where the number of skipped caller frames is fixed).
//
// Please note though, that such functionality adds a significant overhead. For example, in case of zerolog, this
// functionality might slow down logging even multiple times. So, please be careful when making the decision whether to
// use it on production. For other module, logrus, the overhead is much less significant (mostly because logrus is not
// optimized as well as zerolog).
type ReportCallerAdapter struct {
	NextAdapter logger.Adapter
}

func (a ReportCallerAdapter) Log(ctx context.Context, entry logger.Entry) {
	entry.SkippedCallerFrames++ // each middleware adapter must additionally skip one frame (at least)

	if _, file, line, ok := runtime.Caller(entry.SkippedCallerFrames); ok {
		entry = entry.WithFields(logger.Fields{
			"file": file,
			"line": line,
		})
	}

	a.NextAdapter.Log(ctx, entry)
}
