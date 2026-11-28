package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// syncBuffer is a strings.Builder that tail can write while the test reads.
type syncBuffer struct {
	mu sync.Mutex
	b  strings.Builder
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.String()
}

func jsonl(message string) string {
	return `{"time":"2026-11-28T12:00:00Z","level":"INFO","service":"server","message":"` + message + `"}` + "\n"
}

func appendTo(t *testing.T, path, text string) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		t.Fatalf("open for append: %v", err)
	}
	defer f.Close()
	if _, err := f.WriteString(text); err != nil {
		t.Fatalf("append: %v", err)
	}
}

func waitFor(t *testing.T, out *syncBuffer, want string) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if strings.Contains(out.String(), want) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("output never contained %q; got %q", want, out.String())
}

// startTail tails path from its start and returns the output buffer.
func startTail(t *testing.T, path string) *syncBuffer {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	out := &syncBuffer{}
	done := make(chan struct{})
	go func() {
		defer close(done)
		tail(ctx, file, out, slog.LevelDebug, nil)
	}()
	t.Cleanup(func() {
		cancel()
		<-done
		_ = file.Close()
	})
	return out
}

func TestTailPrintsLinesAppendedAfterEOF(t *testing.T) {
	path := filepath.Join(t.TempDir(), "edda.jsonl")
	if err := os.WriteFile(path, []byte(jsonl("first")), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	out := startTail(t, path)
	waitFor(t, out, "first")

	appendTo(t, path, jsonl("second"))
	waitFor(t, out, "second")

	appendTo(t, path, jsonl("third"))
	waitFor(t, out, "third")
}

func TestTailJoinsLineWrittenInParts(t *testing.T) {
	path := filepath.Join(t.TempDir(), "edda.jsonl")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	out := startTail(t, path)

	line := jsonl("split line")
	appendTo(t, path, line[:20])
	time.Sleep(300 * time.Millisecond) // let tail reach EOF mid-line
	appendTo(t, path, line[20:])
	waitFor(t, out, "split line")
}

func TestTailRestartsAfterTruncation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "edda.jsonl")
	if err := os.WriteFile(path, []byte(jsonl("before truncation, a longer line")), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
	out := startTail(t, path)
	waitFor(t, out, "before truncation")

	if err := os.WriteFile(path, []byte(jsonl("after")), 0o600); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	waitFor(t, out, "after")
}
