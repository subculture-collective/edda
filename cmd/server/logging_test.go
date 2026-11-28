package main

import (
	"io"
	stdlog "log"
	"log/slog"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/charmbracelet/log"

	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/logging"
)

func TestRunLogsAsServer(t *testing.T) {
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	stderr, charmDefault := os.Stderr, log.Default()
	slogDefault, stdFlags, stdOutput := slog.Default(), stdlog.Flags(), stdlog.Writer()
	t.Cleanup(func() {
		os.Stderr = stderr
		log.SetDefault(charmDefault)
		slog.SetDefault(slogDefault)
		stdlog.SetFlags(stdFlags)
		stdlog.SetOutput(stdOutput)
		_ = write.Close()
		_ = read.Close()
	})
	os.Stderr = write // one short record fits in the pipe
	t.Chdir(t.TempDir())
	if err := os.WriteFile("bad.yaml", []byte("llm: [\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	if code := run([]string{"-config", "bad.yaml"}); code != 1 {
		t.Fatalf("run() = %d, want 1", code)
	}

	file, err := os.ReadFile(".logs/edda.jsonl")
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	entry, err := logging.ParseJSONLEntry([]byte(strings.TrimSpace(string(file))))
	if err != nil {
		t.Fatalf("parse %q: %v", file, err)
	}
	if entry.Service != "server" {
		t.Errorf("service = %q, want server", entry.Service)
	}
}

func TestLogStartupUsesFields(t *testing.T) {
	buf := logging.NewRingBuffer(4)
	handler := logging.NewHandler(logging.HandlerOptions{RingBuffer: buf, Stderr: io.Discard})
	logger := logging.NewCharmLogger(logging.NewBridgeWriter(slog.New(handler).WithGroup("app")))
	var cfg config.Config
	cfg.LLM.Provider = "claude"

	logStartup(logger.WithPrefix("server"), ":8080", cfg)

	entries := buf.Snapshot()
	if len(entries) != 1 {
		t.Fatalf("entries: got %d, want 1", len(entries))
	}
	got := entries[0]
	if got.Service != "server" || got.Message != "starting HTTP server" {
		t.Errorf("got %q %q, want server %q", got.Service, got.Message, "starting HTTP server")
	}
	want := []logging.Field{{Key: "addr", Value: ":8080"}, {Key: "provider", Value: "claude"}}
	if !slices.Equal(got.Fields, want) {
		t.Errorf("fields: got %+v, want %+v", got.Fields, want)
	}
}
