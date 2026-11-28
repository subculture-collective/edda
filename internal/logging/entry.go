package logging

import (
	"log/slog"
	"time"
)

// Field is a structured logging key/value pair prepared for display and JSONL output.
// Older log files name these "Key" and "Value", which still decode because
// encoding/json ignores case.
type Field struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// LogEntry is a single structured log event captured for UI display and streaming.
type LogEntry struct {
	Time    time.Time
	Level   slog.Level
	Service string
	Message string
	Fields  []Field
}
