package logging

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"strings"

	"github.com/charmbracelet/log"
)

// BridgeWriter adapts charmbracelet/log JSON output into slog records.
type BridgeWriter struct {
	logger *slog.Logger
}

// NewBridgeWriter creates a writer that forwards parsed log lines into logger.
func NewBridgeWriter(logger *slog.Logger) io.Writer {
	return &BridgeWriter{logger: logger}
}

// NewCharmLogger returns the charmbracelet/log logger that the commands point at
// a bridge writer. parseBridgeLine relies on this setup: the JSON format and no
// timestamp, since the slog record carries the time. The text format does not
// escape the message or the prefix, so its parts cannot be told apart.
func NewCharmLogger(w io.Writer) *log.Logger {
	return log.NewWithOptions(w, log.Options{Formatter: log.JSONFormatter})
}

// Write forwards one charmbracelet/log record to slog. A prefix becomes a group,
// which names the service. Any other line is recorded whole, as info.
func (w *BridgeWriter) Write(p []byte) (int, error) {
	line := strings.TrimSpace(string(p))
	if line == "" || w.logger == nil {
		return len(p), nil
	}

	record, ok := parseBridgeLine(line)
	if !ok {
		record = bridgeRecord{level: slog.LevelInfo, message: line}
	}
	logger := w.logger
	if record.prefix != "" {
		logger = logger.WithGroup(record.prefix)
	}
	logger.LogAttrs(context.Background(), record.level, record.message, record.attrs...)
	return len(p), nil
}

// bridgeLevels maps charmbracelet/log's level names. slog has no fatal level, so
// fatal is recorded as an error.
var bridgeLevels = map[string]slog.Level{
	log.DebugLevel.String(): slog.LevelDebug,
	log.InfoLevel.String():  slog.LevelInfo,
	log.WarnLevel.String():  slog.LevelWarn,
	log.ErrorLevel.String(): slog.LevelError,
	log.FatalLevel.String(): slog.LevelError,
}

type bridgeRecord struct {
	level   slog.Level
	prefix  string
	message string
	attrs   []slog.Attr
}

// parseBridgeLine reads one JSON record. charmbracelet/log writes the level,
// prefix and message first, in that order, leaving out any that it lacks, and
// then the key/value pairs. A pair may reuse one of those three names, so each
// is read only from its place. A record without a level, from Print, is info.
func parseBridgeLine(line string) (bridgeRecord, bool) {
	record := bridgeRecord{level: slog.LevelInfo}
	if !strings.HasPrefix(line, "{") || !json.Valid([]byte(line)) {
		return record, false
	}

	dec := json.NewDecoder(strings.NewReader(line))
	if _, err := dec.Token(); err != nil {
		return record, false
	}
	const (
		atLevel = iota
		atPrefix
		atMessage
		atFields
	)
	// next is the earliest part that the following key can still be.
	next := atLevel
	for dec.More() {
		token, err := dec.Token()
		if err != nil {
			return record, false
		}
		key, _ := token.(string)
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			return record, false
		}
		value := bridgeValue(raw)

		level, isLevel := bridgeLevels[value]
		switch {
		case next == atLevel && key == log.LevelKey && isLevel:
			record.level, next = level, atPrefix
		case next <= atPrefix && key == log.PrefixKey:
			record.prefix, next = value, atMessage
		case next <= atMessage && key == log.MessageKey:
			record.message, next = value, atFields
		default:
			record.attrs = append(record.attrs, slog.String(key, value))
			next = atFields
		}
	}
	return record, true
}

// bridgeValue returns a string value as it was logged and any other value as its
// JSON text. null is written as slog writes a nil value.
func bridgeValue(raw json.RawMessage) string {
	var text string
	switch {
	case string(raw) == "null":
		return "<nil>"
	case json.Unmarshal(raw, &text) == nil:
		return text
	default:
		return string(raw)
	}
}
