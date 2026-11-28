package logging

import (
	"errors"
	"io"
	"log/slog"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/log"
)

// TestBridgeWriterReadsRecordsFromCharmLogger feeds the bridge real output from
// a logger built as the commands build it. The message carries no timestamp,
// level or prefix, key/value pairs become fields and the prefix is the service.
func TestBridgeWriterReadsRecordsFromCharmLogger(t *testing.T) {
	t.Parallel()

	levels := []struct {
		level log.Level
		name  string
		want  slog.Level
	}{
		{log.DebugLevel, "debug", slog.LevelDebug},
		{log.InfoLevel, "info", slog.LevelInfo},
		{log.WarnLevel, "warn", slog.LevelWarn},
		{log.ErrorLevel, "error", slog.LevelError},
		{log.FatalLevel, "fatal", slog.LevelError},
	}
	shapes := []struct {
		name    string
		prefix  string
		with    []any // fields set on the logger
		msg     string
		keyvals []any
		fields  []Field
		service string // empty means the fallback, "app"
	}{
		{name: "message", msg: "load config: boom"},
		{name: "message that looks like key values", msg: `starting TUI (provider=ollama timeout=2m0s) reason="x y"`},
		{name: "message naming another level", msg: "ERRO WARN DEBU INFO"},
		{name: "multiline message", msg: "one\ntwo=2"},
		{
			name:    "unquoted values",
			msg:     "boom",
			keyvals: []any{"path", "/x", "status", 500, "retry", true},
			fields:  []Field{{"path", "/x"}, {"status", "500"}, {"retry", "true"}},
		},
		{
			name:    "values with spaces and equals signs",
			msg:     "boom",
			keyvals: []any{"reason", "no such file", "query", "a=b&c=d", "detail", "status=500 retry=true", "pad", "  x  "},
			fields: []Field{
				{"reason", "no such file"}, {"query", "a=b&c=d"},
				{"detail", "status=500 retry=true"}, {"pad", "  x  "},
			},
		},
		{
			name:    "quotes, backslashes, tabs and newlines",
			msg:     `say "hi"`,
			keyvals: []any{"quoted", `say "hi"`, "path", `C:\tmp\new dir`, "tab", "a\tb", "lines", "one\ntwo", "empty", ""},
			fields: []Field{
				{"quoted", `say "hi"`}, {"path", `C:\tmp\new dir`}, {"tab", "a\tb"},
				{"lines", "one\ntwo"}, {"empty", ""},
			},
		},
		{
			name:    "typed values",
			msg:     "boom",
			keyvals: []any{"error", errors.New("dial tcp: <refused> & gone"), "took", 1500 * time.Millisecond, "id", int64(9007199254740993), "none", nil},
			fields: []Field{
				{"error", "dial tcp: <refused> & gone"}, {"took", "1.5s"},
				{"id", "9007199254740993"}, {"none", "<nil>"},
			},
		},
		{
			name:    "composite values",
			msg:     "boom",
			keyvals: []any{"tools", []string{"a", "b c"}},
			fields:  []Field{{"tools", "[a b c]"}},
		},
		{
			name:    "repeated keys and keys named like record parts",
			msg:     "boom",
			keyvals: []any{"msg", "inner", "prefix", "not a service", "a", 1, "a", 2},
			fields:  []Field{{"msg", "inner"}, {"prefix", "not a service"}, {"a", "1"}, {"a", "2"}},
		},
		{
			name:    "missing value",
			msg:     "boom",
			keyvals: []any{"orphan"},
			fields:  []Field{{"orphan", "missing value"}},
		},
		{
			name:    "key values only",
			keyvals: []any{"status", 500},
			fields:  []Field{{"status", "500"}},
		},
		{
			name:    "logger fields come first",
			with:    []any{"request_id", "r 1"},
			msg:     "boom",
			keyvals: []any{"status", 500},
			fields:  []Field{{"request_id", "r 1"}, {"status", "500"}},
		},
		{name: "prefix", prefix: "server", msg: "boom", service: "server"},
		{name: "prefix with a space and a colon", prefix: "auto embed: v2", msg: "boom", service: "auto embed: v2"},
		{
			name:    "prefix and key values",
			prefix:  "server",
			msg:     "load config: boom",
			keyvals: []any{"status", 500},
			fields:  []Field{{"status", "500"}},
			service: "server",
		},
		{
			name:    "prefix without a message",
			prefix:  "server",
			keyvals: []any{"status", 500},
			fields:  []Field{{"status", "500"}},
			service: "server",
		},
		// A service attribute names the service, as it does in a slog record.
		{
			name:    "service key",
			prefix:  "server",
			msg:     "boom",
			keyvals: []any{"service", "embed"},
			fields:  []Field{{"service", "embed"}},
			service: "embed",
		},
	}

	for _, level := range levels {
		for _, shape := range shapes {
			t.Run(level.name+"/"+shape.name, func(t *testing.T) {
				t.Parallel()

				got := bridgeCharmRecord(t, func(logger *CharmLogger) {
					if shape.prefix != "" {
						logger = logger.WithPrefix(shape.prefix)
					}
					if shape.with != nil {
						logger = logger.With(shape.with...)
					}
					var msg any
					if shape.msg != "" {
						msg = shape.msg
					}
					logger.Log(level.level, msg, shape.keyvals...)
				})

				service := shape.service
				if service == "" {
					service = "app"
				}
				if got.Level != level.want {
					t.Errorf("level: got %v, want %v", got.Level, level.want)
				}
				if got.Message != shape.msg {
					t.Errorf("message: got %q, want %q", got.Message, shape.msg)
				}
				if got.Service != service {
					t.Errorf("service: got %q, want %q", got.Service, service)
				}
				if !slices.Equal(got.Fields, shape.fields) {
					t.Errorf("fields: got %q, want %q", got.Fields, shape.fields)
				}
			})
		}
	}

	// Print has no level, so a message must not be read as one.
	for _, msg := range []string{"Information about the world", "Errors: none", "warning shown"} {
		t.Run("print/"+msg, func(t *testing.T) {
			t.Parallel()

			got := bridgeCharmRecord(t, func(logger *CharmLogger) { logger.base.Print(msg) })

			if got.Level != slog.LevelInfo || got.Message != msg || got.Service != "app" || len(got.Fields) != 0 {
				t.Errorf("got %v %q %q %q, want INFO app %q without fields", got.Level, got.Service, got.Message, got.Fields, msg)
			}
		})
	}
}

// CLICOLOR_FORCE makes lipgloss color output that is not a terminal.
func TestBridgeWriterKeepsEscapeCodesOutWhenColorIsForced(t *testing.T) {
	t.Setenv("CLICOLOR_FORCE", "1")

	got := bridgeCharmRecord(t, func(logger *CharmLogger) {
		logger.WithPrefix("server").Error("boom", "status", 500)
	})

	line := formatHumanLine(got)
	if strings.Contains(line, "\x1b") || got.Level != slog.LevelError || got.Message != "boom" || got.Service != "server" {
		t.Fatalf("got %v %q %q in %q", got.Level, got.Service, got.Message, line)
	}
}

// bridgeCharmRecord logs one record through a charmbracelet/log logger built as
// the commands build it, and returns what the bridge recorded.
func bridgeCharmRecord(t *testing.T, write func(*CharmLogger)) LogEntry {
	t.Helper()

	buf := NewRingBuffer(4)
	handler := NewHandler(HandlerOptions{
		MinLevel:   slog.LevelDebug,
		RingBuffer: buf,
		Stderr:     io.Discard,
		File:       io.Discard,
	})
	// Setup gives the bridge a logger in the "app" group.
	logger := NewCharmLogger(NewBridgeWriter(slog.New(handler).WithGroup("app")))
	logger.base.SetLevel(log.DebugLevel)

	write(logger)

	entries := buf.Snapshot()
	if len(entries) != 1 {
		t.Fatalf("entries: got %d, want 1", len(entries))
	}
	return entries[0]
}
