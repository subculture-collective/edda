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

// rotate renames path aside and returns the new name of the old file.
func rotate(t *testing.T, path, suffix string) string {
	t.Helper()
	old := path + suffix
	if err := os.Rename(path, old); err != nil {
		t.Fatalf("rename: %v", err)
	}
	return old
}

func create(t *testing.T, path, text string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(text), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func assertInOrder(t *testing.T, out string, messages ...string) {
	t.Helper()
	last := -1
	for _, message := range messages {
		if n := strings.Count(out, message); n != 1 {
			t.Fatalf("%q printed %d times, want once; output %q", message, n, out)
		}
		at := strings.Index(out, message)
		if at < last {
			t.Fatalf("%q printed out of order; output %q", message, out)
		}
		last = at
	}
}

func TestTailFollowsRotationByRename(t *testing.T) {
	path := filepath.Join(t.TempDir(), "edda.jsonl")
	// The new file is longer than the old one, so a size check alone cannot
	// notice the change.
	create(t, path, jsonl("line one"))
	out := startTail(t, path)
	waitFor(t, out, "line one")

	rotate(t, path, ".1")
	create(t, path, jsonl("line two, from the new file")+jsonl("line three"))
	waitFor(t, out, "line three")

	appendTo(t, path, jsonl("line four"))
	waitFor(t, out, "line four")

	rotate(t, path, ".2")
	create(t, path, jsonl("line five"))
	waitFor(t, out, "line five")

	assertInOrder(t, out.String(), "line one", "line two", "line three", "line four", "line five")
}

func TestTailDrainsOldFileBeforeReopening(t *testing.T) {
	path := filepath.Join(t.TempDir(), "edda.jsonl")
	create(t, path, jsonl("line one"))
	out := startTail(t, path)
	waitFor(t, out, "line one")

	// All of this normally lands within one poll, so tail sees the new file
	// before it has read the old file's last lines. The last has no newline.
	appendTo(t, path, jsonl("line two")+strings.TrimSuffix(jsonl("line three"), "\n"))
	rotate(t, path, ".1")
	create(t, path, jsonl("line four"))
	waitFor(t, out, "line four")

	assertInOrder(t, out.String(), "line one", "line two", "line three", "line four")
}

func TestTailKeepsOldFileUntilNewOneAppears(t *testing.T) {
	path := filepath.Join(t.TempDir(), "edda.jsonl")
	create(t, path, jsonl("line one"))
	out := startTail(t, path)
	waitFor(t, out, "line one")

	old := rotate(t, path, ".1")
	time.Sleep(300 * time.Millisecond) // nothing at path
	appendTo(t, old, jsonl("line two"))
	waitFor(t, out, "line two")

	create(t, path, jsonl("line three"))
	waitFor(t, out, "line three")

	assertInOrder(t, out.String(), "line one", "line two", "line three")
}
