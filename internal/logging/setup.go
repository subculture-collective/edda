package logging

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// SetupResult holds the outputs of a successful Setup call.
type SetupResult struct {
	RingBuffer   *RingBuffer
	BridgeWriter io.Writer
	Cleanup      func()

	console *consoleWriter
}

// SetConsole turns the stderr copy of log records on or off. A full-screen
// terminal UI turns it off while it runs, because those writes paint over it.
// Records still reach the ring buffer and the JSONL file.
func (r *SetupResult) SetConsole(enabled bool) {
	r.console.muted.Store(!enabled)
}

// Setup initializes the unified logging stack: ring buffer for TUI consumption,
// JSONL file for CLI streaming, and stderr human-readable output. It sets
// slog.Default to the new handler and returns a BridgeWriter suitable for
// charmbracelet/log.
func Setup(logPath string, minLevel slog.Leveler) (*SetupResult, error) {
	if logPath == "" {
		return nil, errors.New("logging: logPath must not be empty")
	}

	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	buf := NewRingBuffer(2048)

	console := &consoleWriter{w: os.Stderr}

	handler := NewHandler(HandlerOptions{
		MinLevel:   minLevel,
		RingBuffer: buf,
		Stderr:     console,
		File:       file,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	bridgeWriter := NewBridgeWriter(logger.WithGroup("app"))

	return &SetupResult{
		RingBuffer:   buf,
		BridgeWriter: bridgeWriter,
		Cleanup:      func() { file.Close() },
		console:      console,
	}, nil
}
