package logging

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
)

// BridgeWriter adapts charmbracelet/log text output into slog records.
type BridgeWriter struct {
	logger *slog.Logger
}

// NewBridgeWriter creates a writer that forwards parsed log lines into logger.
func NewBridgeWriter(logger *slog.Logger) io.Writer {
	return &BridgeWriter{logger: logger}
}

// NewCharmLogger returns the charmbracelet/log logger that the commands point at
// a bridge writer. parseBridgeLine relies on this setup: the text format, the
// default time format and no color.
func NewCharmLogger(w io.Writer) *log.Logger {
	logger := log.NewWithOptions(w, log.Options{ReportTimestamp: true})
	// CLICOLOR_FORCE would otherwise put escape codes around the level token.
	logger.SetColorProfile(termenv.Ascii)
	return logger
}

// Write parses a single text log line, infers its level, and forwards it to slog.
func (w *BridgeWriter) Write(p []byte) (int, error) {
	line := strings.TrimSuffix(string(p), "\n")
	line = strings.TrimSuffix(line, "\r")
	if strings.TrimSpace(line) == "" || w.logger == nil {
		return len(p), nil
	}

	level, msg := parseBridgeLine(line)
	w.logger.Log(context.Background(), level, msg)
	return len(p), nil
}

// bridgeLevels maps the level tokens of charmbracelet/log's text format. slog
// has no fatal level, so FATA is recorded as an error.
var bridgeLevels = map[string]slog.Level{
	"DEBU": slog.LevelDebug,
	"INFO": slog.LevelInfo,
	"WARN": slog.LevelWarn,
	"ERRO": slog.LevelError,
	"FATA": slog.LevelError,
}

// parseBridgeLine reads the level token, which is the first word of the line or
// the one after the timestamp. A timestamped line is kept whole as the message;
// otherwise the token is dropped. A line without a token is recorded as info.
func parseBridgeLine(line string) (slog.Level, string) {
	trimmed := strings.TrimSpace(line)
	rest, timestamped := cutBridgeTimestamp(trimmed)

	token, msg := rest, ""
	if end := strings.IndexAny(rest, " \n"); end >= 0 {
		token, msg = rest[:end], rest[end:]
	}

	level, ok := bridgeLevels[token]
	switch {
	case !ok:
		return slog.LevelInfo, trimmed
	case timestamped:
		return level, trimmed
	default:
		return level, strings.TrimSpace(msg)
	}
}

// cutBridgeTimestamp removes a leading timestamp in the default time format.
func cutBridgeTimestamp(line string) (string, bool) {
	const width = len(log.DefaultTimeFormat)
	if len(line) <= width || line[width] != ' ' {
		return line, false
	}
	if _, err := time.Parse(log.DefaultTimeFormat, line[:width]); err != nil {
		return line, false
	}
	return line[width+1:], true
}
