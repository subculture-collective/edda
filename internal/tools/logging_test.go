package tools

import (
	"io"
	"log/slog"
	"testing"

	"github.com/charmbracelet/log"

	"git.subcult.tv/subculture-collective/edda/internal/logging"
)

// This test replaces charmbracelet/log's default, so it does not run in parallel.
func TestDefaultLoggersAreNamedTools(t *testing.T) {
	previous := log.Default()
	t.Cleanup(func() { log.SetDefault(previous) })
	buf := logging.NewRingBuffer(4)
	handler := logging.NewHandler(logging.HandlerOptions{RingBuffer: buf, Stderr: io.Discard})
	logging.NewCharmLogger(logging.NewBridgeWriter(slog.New(handler).WithGroup("app"))).SetDefault()

	NewAutoEmbedder(nil, nil, nil).logger.Error("embedder")
	NewPostValidator(nil, nil).logger.Error("validator")

	entries := buf.Snapshot()
	if len(entries) != 2 {
		t.Fatalf("entries: got %d, want 2", len(entries))
	}
	for _, entry := range entries {
		if entry.Service != "tools" {
			t.Errorf("%s: service = %q, want tools", entry.Message, entry.Service)
		}
	}
}
