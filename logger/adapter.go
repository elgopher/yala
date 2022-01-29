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

// Entry is a logging entry created by logger and passed to adapter.
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

// Level is a severity level of message. Use Level.MoreSevereThan to compare two levels.
type Level int8

const (
	// DebugLevel level is usually enabled only when debugging (disabled in production). Very verbose logging.
	DebugLevel Level = iota - 1
	// InfoLevel is used for informational messages, for confirmation that the program is working as expected.
	InfoLevel
	// WarnLevel is used for non-critical entries that deserve eyes.
	WarnLevel
	// ErrorLevel is used for errors that should definitely be noted. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
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

// MoreSevereThan returns true if level is more severe than the argument.
func (l Level) MoreSevereThan(other Level) bool {
	return l > other
}

// Field contains key-value pair.
type Field struct {
	Key   string
	Value interface{}
}
