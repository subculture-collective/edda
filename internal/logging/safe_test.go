package logging

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

type nilErr struct{ msg *string }

func (e *nilErr) Error() string { return *e.msg }

type nilStringer struct{ s *string }

func (n *nilStringer) String() string { return *n.s }

type panicErr struct{}

func (panicErr) Error() string { panic("error boom") }

type panicStringer struct{}

func (panicStringer) String() string { panic("stringer boom") }

type panicJSON struct{}

func (panicJSON) MarshalJSON() ([]byte, error) { panic("json boom") }

type panicText struct{}

func (panicText) MarshalText() ([]byte, error) { panic("text boom") }

type panicLogValuer struct{}

func (panicLogValuer) LogValue() slog.Value { panic("log value boom") }

// charmbracelet/log's JSON formatter calls methods on keys and values without
// recovering. Every way of logging a value has to survive one that panics, and
// has to write it as fmt does.
func TestCharmLoggerNilValuesDoNotPanic(t *testing.T) {
	t.Parallel()

	var e *nilErr
	var s *nilStringer
	values := []struct {
		name  string
		value any
	}{
		{"nil error", error(e)},
		{"nil stringer", s},
		{"joined nil error", errors.Join(error(e))},
		{"panicking error", panicErr{}},
		{"panicking stringer", panicStringer{}},
		{"panicking MarshalJSON", panicJSON{}},
		{"panicking MarshalText", panicText{}},
		{"panicking LogValue", panicLogValuer{}},
		{"nested panicking MarshalJSON", []any{panicJSON{}}},
		{"channel", make(chan int)},
		{"func", func() {}},
	}
	calls := []struct {
		name string
		log  func(l *CharmLogger, value any)
	}{
		{"value", func(l *CharmLogger, v any) { l.Error("boom", "k", v) }},
		{"key", func(l *CharmLogger, v any) { l.Error("boom", v, "k") }},
		{"missing value", func(l *CharmLogger, v any) { l.Error("boom", v) }},
		{"message", func(l *CharmLogger, v any) { l.Error(v) }},
		{"format argument", func(l *CharmLogger, v any) { l.Errorf("boom %v", v) }},
		{"logger field", func(l *CharmLogger, v any) { l.With("k", v).Error("boom") }},
		{"logger field with prefix", func(l *CharmLogger, v any) { l.WithPrefix("p").With(v, v).Info("boom") }},
		{"debug", func(l *CharmLogger, v any) { l.Debug("boom", "k", v) }},
		{"warn", func(l *CharmLogger, v any) { l.Warn("boom", "k", v) }},
	}

	for _, value := range values {
		for _, call := range calls {
			t.Run(value.name+"/"+call.name, func(t *testing.T) {
				t.Parallel()

				var buf bytes.Buffer
				call.log(NewCharmLogger(&buf), value.value)

				// Debug is below the default level.
				if call.name == "debug" {
					return
				}
				if !json.Valid(buf.Bytes()) {
					t.Fatalf("output is not one JSON record: %q", buf.String())
				}
				if strings.Contains(buf.String(), "invalid value") {
					t.Errorf("value was lost: %q", buf.String())
				}
			})
		}

		t.Run(value.name+"/written as fmt writes it", func(t *testing.T) {
			t.Parallel()

			got := bridgeCharmRecord(t, func(logger *CharmLogger) { logger.Error("boom", "k", value.value) })

			want := []Field{{"k", fmt.Sprintf("%+v", value.value)}}
			if len(got.Fields) != 1 || got.Fields[0] != want[0] {
				t.Errorf("fields: got %q, want %q", got.Fields, want)
			}
		})
	}
}

// A nil pointer reads as <nil>, as it did with the text formatter.
func TestCharmLoggerWritesNilPointersAsNil(t *testing.T) {
	t.Parallel()

	var e *nilErr
	var s *nilStringer

	got := bridgeCharmRecord(t, func(logger *CharmLogger) { logger.Error("boom", "err", error(e), "s", s) })

	want := []Field{{"err", "<nil>"}, {"s", "<nil>"}}
	if len(got.Fields) != 2 || got.Fields[0] != want[0] || got.Fields[1] != want[1] {
		t.Errorf("fields: got %q, want %q", got.Fields, want)
	}
}
