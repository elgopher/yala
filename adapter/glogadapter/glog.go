// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package glogadapter

import (
	"context"
	"strings"

	"github.com/elgopher/yala/adapter/logfmt"
	"github.com/elgopher/yala/logger"
	"github.com/golang/glog"
)

// Adapter is a logger.Adapter implementation, which is using `glog` package (https://github.com/golang/glog).
type Adapter struct{}

func (a Adapter) Log(_ context.Context, entry logger.Entry) {
	var fieldsAndError strings.Builder

	if len(entry.Fields) > 0 {
		logfmt.WriteFields(&fieldsAndError, entry.Fields)
	}

	if entry.Error != nil {
		if len(entry.Fields) > 0 {
			fieldsAndError.WriteByte(' ')
		}

		logfmt.WriteField(&fieldsAndError, logger.Field{Key: "error", Value: entry.Error})
	}

	fieldsAndErrorString := fieldsAndError.String()
	message := entry.Message
	depth := entry.SkippedCallerFrames + 1

	switch entry.Level {
	case logger.DebugLevel:
		glog.InfoDepth(depth, message, " ", fieldsAndErrorString)
	case logger.InfoLevel:
		glog.InfoDepth(depth, message, " ", fieldsAndErrorString)
	case logger.WarnLevel:
		glog.WarningDepth(depth, message, " ", fieldsAndErrorString)
	case logger.ErrorLevel:
		glog.ErrorDepth(depth, message, " ", fieldsAndErrorString)
	default:
		glog.InfoDepth(depth, message, " ", fieldsAndErrorString)
	}
}
