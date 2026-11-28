package main

import (
	stdlog "log"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/log"

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
