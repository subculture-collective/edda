package logging

import (
	"io"
	"log/slog"
	"testing"
)

// A line that is not a charmbracelet/log JSON record is kept whole.
func TestBridgeWriterKeepsOtherLinesWhole(t *testing.T) {
	t.Parallel()

	buf := NewRingBuffer(8)
	logger := slog.New(NewHandler(HandlerOptions{
		MinLevel:   slog.LevelDebug,
		RingBuffer: buf,
		Stderr:     io.Discard,
		File:       io.Discard,
	}))
	bridge := NewBridgeWriter(logger)

	lines := []string{
		"ERRO plain text\r\n",
		"{not json}\n",
		`{"msg":"cut off"` + "\n",
		`{"msg":"one"} trailing` + "\n",
		`["msg"]` + "\n",
		"\n",
	}
	for _, line := range lines {
		if n, err := bridge.Write([]byte(line)); err != nil {
			t.Fatalf("write %q: %v", line, err)
		} else if n != len(line) {
			t.Fatalf("write len for %q: got %d, want %d", line, n, len(line))
		}
	}

	want := []string{"ERRO plain text", "{not json}", `{"msg":"cut off"`, `{"msg":"one"} trailing`, `["msg"]`}
	entries := buf.Snapshot()
	if len(entries) != len(want) {
		t.Fatalf("snapshot len: got %d, want %d", len(entries), len(want))
	}
	for i := range want {
		if entries[i].Level != slog.LevelInfo || entries[i].Message != want[i] || len(entries[i].Fields) != 0 {
			t.Errorf("entry[%d]: got %v %q %q, want INFO %q without fields", i, entries[i].Level, entries[i].Message, entries[i].Fields, want[i])
		}
	}
}
