package logging

import (
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/log"
)

// TestBridgeWriterReadsLevelsFromCharmLogger feeds the bridge real
// charmbracelet/log output. With a timestamp, which is what the commands use,
// the message stays the whole line. Without one, the level token is dropped.
func TestBridgeWriterReadsLevelsFromCharmLogger(t *testing.T) {
	t.Parallel()

	const stamp = "2026/11/28 12:34:56 "
	levels := []struct {
		level log.Level
		token string
		want  slog.Level
	}{
		{log.DebugLevel, "DEBU", slog.LevelDebug},
		{log.InfoLevel, "INFO", slog.LevelInfo},
		{log.WarnLevel, "WARN", slog.LevelWarn},
		{log.ErrorLevel, "ERRO", slog.LevelError},
		{log.FatalLevel, "FATA", slog.LevelError},
	}
	shapes := []struct {
		name    string
		prefix  string
		msg     string
		keyvals []any
		rest    string // the line after the level token
	}{
		{name: "message", msg: "load config: boom", rest: "load config: boom"},
		{name: "prefix", prefix: "server", msg: "boom", rest: "server: boom"},
		{name: "key values", msg: "boom", keyvals: []any{"path", "/x", "reason", "no such file"}, rest: `boom path=/x reason="no such file"`},
		{name: "prefix and key values", prefix: "server", msg: "boom", keyvals: []any{"status", 500}, rest: "server: boom status=500"},
		{name: "key values only", keyvals: []any{"status", 500}, rest: "status=500"},
		{name: "multiline value", msg: "boom", keyvals: []any{"detail", "one\ntwo"}, rest: "boom\n  detail=\n  │ one\n  │ two"},
		{name: "multiline value only", keyvals: []any{"detail", "one\ntwo"}, rest: "\n  detail=\n  │ one\n  │ two"},
		{name: "message naming another level", msg: "ERRO WARN DEBU INFO", rest: "ERRO WARN DEBU INFO"},
	}

	for _, timestamp := range []bool{true, false} {
		for _, level := range levels {
			for _, shape := range shapes {
				name := level.token + "/" + shape.name
				if timestamp {
					name += "/timestamp"
				}
				t.Run(name, func(t *testing.T) {
					t.Parallel()

					got := bridgeCharmRecord(t, timestamp, func(logger *log.Logger) {
						if shape.prefix != "" {
							logger = logger.WithPrefix(shape.prefix)
						}
						var msg any
						if shape.msg != "" {
							msg = shape.msg
						}
						logger.Log(level.level, msg, shape.keyvals...)
					})

					wantMsg := strings.TrimSpace(shape.rest)
					if timestamp {
						sep := " "
						if strings.HasPrefix(shape.rest, "\n") {
							sep = ""
						}
						wantMsg = stamp + level.token + sep + shape.rest
					}
					if got.Level != level.want {
						t.Errorf("level: got %v, want %v (message %q)", got.Level, level.want, got.Message)
					}
					if got.Message != wantMsg {
						t.Errorf("message: got %q, want %q", got.Message, wantMsg)
					}
				})
			}
		}
	}

	// Print has no level token, so a message must not be read as one.
	for _, msg := range []string{"Information about the world", "Errors: none", "warning shown"} {
		for _, timestamp := range []bool{true, false} {
			t.Run("print/"+msg, func(t *testing.T) {
				t.Parallel()

				got := bridgeCharmRecord(t, timestamp, func(logger *log.Logger) { logger.Print(msg) })

				wantMsg := msg
				if timestamp {
					wantMsg = stamp + msg
				}
				if got.Level != slog.LevelInfo || got.Message != wantMsg {
					t.Errorf("got %v %q, want INFO %q", got.Level, got.Message, wantMsg)
				}
			})
		}
	}
}

// CLICOLOR_FORCE makes lipgloss color output that is not a terminal.
func TestBridgeWriterReadsLevelsWhenColorIsForced(t *testing.T) {
	t.Setenv("CLICOLOR_FORCE", "1")

	got := bridgeCharmRecord(t, true, func(logger *log.Logger) { logger.Error("boom", "status", 500) })

	const want = "2026/11/28 12:34:56 ERRO boom status=500"
	if got.Level != slog.LevelError || got.Message != want {
		t.Fatalf("got %v %q, want ERROR %q", got.Level, got.Message, want)
	}
}

// bridgeCharmRecord logs one record through a charmbracelet/log logger built as
// the commands build it, and returns what the bridge recorded.
func bridgeCharmRecord(t *testing.T, timestamp bool, write func(*log.Logger)) LogEntry {
	t.Helper()

	buf := NewRingBuffer(4)
	handler := NewHandler(HandlerOptions{
		MinLevel:   slog.LevelDebug,
		RingBuffer: buf,
		Stderr:     io.Discard,
		File:       io.Discard,
	})
	logger := NewCharmLogger(NewBridgeWriter(slog.New(handler)))
	logger.SetReportTimestamp(timestamp)
	logger.SetLevel(log.DebugLevel)
	logger.SetTimeFunction(func(time.Time) time.Time {
		return time.Date(2026, 11, 28, 12, 34, 56, 0, time.UTC)
	})

	write(logger)

	entries := buf.Snapshot()
	if len(entries) != 1 {
		t.Fatalf("entries: got %d, want 1", len(entries))
	}
	return entries[0]
}
