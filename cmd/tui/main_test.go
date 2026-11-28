package main

import (
	"io"
	stdlog "log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/log"

	"git.subcult.tv/subculture-collective/edda/internal/logging"
)

func TestParseConfigPathUsesDefaultWhenFlagMissing(t *testing.T) {
	defaultPath := "/tmp/default.yaml"

	got, err := parseConfigPath(nil, defaultPath)
	if err != nil {
		t.Fatalf("parseConfigPath() error = %v", err)
	}
	if got != defaultPath {
		t.Fatalf("expected default path %q, got %q", defaultPath, got)
	}
}

func TestParseConfigPathUsesFlagOverride(t *testing.T) {
	defaultPath := "/tmp/default.yaml"
	overridePath := "/tmp/override.yaml"

	got, err := parseConfigPath([]string{"--config", overridePath}, defaultPath)
	if err != nil {
		t.Fatalf("parseConfigPath() error = %v", err)
	}
	if got != overridePath {
		t.Fatalf("expected override path %q, got %q", overridePath, got)
	}
}

func TestParseConfigPathReturnsErrorForUnknownFlag(t *testing.T) {
	if _, err := parseConfigPath([]string{"--unknown"}, ""); err == nil {
		t.Fatal("expected error for unknown flag, got nil")
	}
}

func TestRunQuietlyKeepsLogRecordsOffStderr(t *testing.T) {
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	stderr, defaultLogger, stdFlags, stdOutput := os.Stderr, slog.Default(), stdlog.Flags(), stdlog.Writer()
	t.Cleanup(func() {
		os.Stderr = stderr
		slog.SetDefault(defaultLogger)
		stdlog.SetFlags(stdFlags)
		stdlog.SetOutput(stdOutput)
	})
	os.Stderr = write

	logPath := filepath.Join(t.TempDir(), "edda.jsonl")
	result, err := logging.Setup(logPath, slog.LevelDebug)
	if err != nil {
		t.Fatalf("Setup: %v", err)
	}
	t.Cleanup(result.Cleanup)
	charm := log.NewWithOptions(result.BridgeWriter, log.Options{})

	slog.Info("before the program")
	if err := runQuietly(result, func() error {
		slog.Error("slog during the program")
		stdlog.Print("standard log during the program")
		charm.Error("charm log during the program")
		return nil
	}); err != nil {
		t.Fatalf("runQuietly: %v", err)
	}
	slog.Info("after the program")

	if err := write.Close(); err != nil {
		t.Fatalf("close pipe: %v", err)
	}
	console, err := io.ReadAll(read)
	if err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	file, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	var panel strings.Builder
	for _, entry := range result.RingBuffer.Snapshot() {
		panel.WriteString(entry.Message + "\n")
	}

	if strings.Contains(string(console), "during the program") {
		t.Errorf("stderr was written while the program ran:\n%s", console)
	}
	for _, want := range []string{"before the program", "after the program"} {
		if !strings.Contains(string(console), want) {
			t.Errorf("stderr is missing %q:\n%s", want, console)
		}
	}
	for _, want := range []string{"slog during", "standard log during", "charm log during"} {
		if !strings.Contains(string(file), want) {
			t.Errorf("log file is missing %q:\n%s", want, file)
		}
		if !strings.Contains(panel.String(), want) {
			t.Errorf("log panel buffer is missing %q:\n%s", want, panel.String())
		}
	}
}
