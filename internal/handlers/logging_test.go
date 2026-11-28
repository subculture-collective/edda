package handlers

import (
	"io"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/charmbracelet/log"

	"git.subcult.tv/subculture-collective/edda/internal/logging"
)

// These tests replace charmbracelet/log's default, so they do not run in parallel.
func TestDefaultLoggingIsNamedHandlers(t *testing.T) {
	tests := []struct {
		name  string
		write func()
	}{
		{"startup handlers without a logger", func() { NewStartupHandlers(nil, nil, nil, nil).Logger.Errorf("boom") }},
		{"JSON encode failure", func() { writeJSON(httptest.NewRecorder(), 200, make(chan int)) }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			previous := log.Default()
			t.Cleanup(func() { log.SetDefault(previous) })
			buf := logging.NewRingBuffer(4)
			handler := logging.NewHandler(logging.HandlerOptions{RingBuffer: buf, Stderr: io.Discard})
			logging.NewCharmLogger(logging.NewBridgeWriter(slog.New(handler).WithGroup("app"))).SetDefault()

			tc.write()

			entries := buf.Snapshot()
			if len(entries) != 1 || entries[0].Service != "handlers" {
				t.Errorf("entries: got %+v, want one from handlers", entries)
			}
		})
	}
}
