package server

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/llm"
	"git.subcult.tv/subculture-collective/edda/internal/logging"
)

type failingProvider struct{}

func (failingProvider) Complete(context.Context, []llm.Message, []llm.Tool) (*llm.Response, error) {
	return nil, errors.New("model offline")
}

func (failingProvider) Stream(context.Context, []llm.Message, []llm.Tool) (<-chan llm.StreamChunk, error) {
	return nil, errors.New("model offline")
}

// bridgedLogger returns a logger set up as cmd/server sets it up, and the buffer
// that receives its records.
func bridgedLogger() (*logging.CharmLogger, *logging.RingBuffer) {
	buf := logging.NewRingBuffer(8)
	handler := logging.NewHandler(logging.HandlerOptions{RingBuffer: buf, Stderr: io.Discard})
	bridge := logging.NewBridgeWriter(slog.New(handler).WithGroup("app"))
	return logging.NewCharmLogger(bridge).WithPrefix("server"), buf
}

func TestRequestLogUsesFields(t *testing.T) {
	t.Parallel()

	logger, buf := bridgedLogger()
	router := NewRouterWithProvider(logger, nil, nil, nil, nil, uuid.Nil, config.Config{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz?token=secret", nil)
	req.Header.Set("X-Request-Id", "req-1")
	router.ServeHTTP(httptest.NewRecorder(), req)

	entries := buf.Snapshot()
	if len(entries) != 1 {
		t.Fatalf("entries: got %d, want 1", len(entries))
	}
	got := entries[0]
	if got.Service != "server" || got.Message != "http request" || got.Level != slog.LevelInfo {
		t.Errorf("got %v %q %q, want INFO server %q", got.Level, got.Service, got.Message, "http request")
	}

	var keys []string
	values := map[string]string{}
	for _, f := range got.Fields {
		keys = append(keys, f.Key)
		values[f.Key] = f.Value
	}
	if want := []string{"method", "path", "status", "bytes", "duration", "request_id"}; !slices.Equal(keys, want) {
		t.Fatalf("keys: got %q, want %q", keys, want)
	}
	for key, want := range map[string]string{"method": "GET", "path": "/healthz", "status": "200", "bytes": "2", "request_id": "req-1"} {
		if values[key] != want {
			t.Errorf("%s: got %q, want %q", key, values[key], want)
		}
	}
	if !strings.HasSuffix(values["duration"], "s") {
		t.Errorf("duration: got %q, want a duration", values["duration"])
	}
}

func TestRouterNamesTheHandlerService(t *testing.T) {
	t.Parallel()

	logger, buf := bridgedLogger()
	router := NewRouterWithProvider(logger, nil, nil, failingProvider{}, nil, uuid.Nil, config.Config{}, nil)

	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/api/v1/campaigns/start/campaign-interview", nil))

	var got []string
	for _, entry := range buf.Snapshot() {
		got = append(got, entry.Service+" "+entry.Level.String())
	}
	if want := []string{"handlers ERROR", "server INFO"}; !slices.Equal(got, want) {
		t.Errorf("records: got %q, want %q", got, want)
	}
}
