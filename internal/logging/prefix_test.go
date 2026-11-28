package logging

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/charmbracelet/log"
)

func TestWithPrefixNamesTheService(t *testing.T) {
	t.Parallel()

	got := bridgeCharmRecord(t, func(logger *CharmLogger) {
		WithPrefix(logger.WithPrefix("server"), "handlers").Info("ready", "n", 1)
	})

	if got.Service != "handlers" || got.Message != "ready" {
		t.Errorf("got service %q message %q, want handlers ready", got.Service, got.Message)
	}
}

func TestWithPrefixOnAPlainCharmLogger(t *testing.T) {
	t.Parallel()

	var out bytes.Buffer
	WithPrefix(log.New(&out), "handlers").Info("ready")

	if !strings.Contains(out.String(), "handlers: ready") {
		t.Errorf("output = %q, want the handlers prefix", out.String())
	}
}

type otherLogger struct{ Logger }

func TestWithPrefixKeepsOtherLoggers(t *testing.T) {
	t.Parallel()

	logger := &otherLogger{log.New(io.Discard)}

	if got := WithPrefix(logger, "handlers"); got != Logger(logger) {
		t.Errorf("got %T, want the logger unchanged", got)
	}
}
