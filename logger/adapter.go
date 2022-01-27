// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package logger

import (
	"context"
	"strconv"
)

// Adapter is an interface to be implemented by logger adapters.
type Adapter interface {
	Log(context.Context, Entry)
}

type Entry struct {
	Level   Level
	Message string
	Fields  []Field // Fields can be nil
	Error   error   // Error can be nil
	// SkippedCallerFrames can be used by logger.Adapter to extract caller information (file and line number)
	SkippedCallerFrames int
}

// With creates a new entry with additional field.
func (e Entry) With(field Field) Entry {
	newLen := len(e.Fields) + 1
	fields := make([]Field, newLen)
	copy(fields, e.Fields)
	e.Fields = fields
	e.Fields[newLen-1] = field

	return e
}

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return strconv.Itoa(int(l))
	}
}

func (l Level) MoreSevereThan(other Level) bool {
	return l > other
}

type Field struct {
	Key   string
	Value interface{}
}
