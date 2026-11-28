package logging

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/url"
	"strings"
	"testing"
)

func TestRedactURL(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		{"plain endpoint", "http://localhost:11434", "http://localhost:11434"},
		{"userinfo stripped", "https://user:fake-pass@ollama.example:11434/api", "https://ollama.example:11434/api"},
		{"username only stripped", "https://fake-token@bridge.example/", "https://bridge.example/"},
		{"postgres password masked", "postgres://edda:fake-pass@db.example:5432/edda?sslmode=disable", "postgres://edda:xxxxx@db.example:5432/edda?sslmode=disable"},
		{"postgresql password masked", "postgresql://edda:fake-pass@db.example/edda", "postgresql://edda:xxxxx@db.example/edda"},
		{"postgres without password", "postgres://edda@db.example/edda", "postgres://edda@db.example/edda"},
		{"secret query params dropped", "http://h.example/x?model=a&api_key=fake1&token=fake2&Password=fake3&b=2", "http://h.example/x?model=a&b=2"},
		{"postgres query password dropped", "postgres://db.example/edda?user=edda&password=fake-pass", "postgres://db.example/edda?user=edda"},
		{"keyword dsn password masked", "host=db.example user=edda password=fake-pass dbname=edda", "host=db.example user=edda password=xxxxx dbname=edda"},
		{"keyword dsn quoted password masked", "host=db.example password='fake pass' dbname=edda", "host=db.example password=xxxxx dbname=edda"},
		{"empty", "", ""},
		{"unparseable url still redacted", "http://user:fake%zzpass@h.example/", "http://h.example/"},
		{"unparseable url loses its query", "http://user:fake-pass@h.example:bad/x?api_key=fake-key#f", "http://h.example:bad/x"},
		{"unescaped slash in password", "http://user:fake/pass@h.example", "http://h.example"},
		{"unescaped question mark in password", "http://user:fake?pass@h.example", "http://xxxxx"},
		{"at sign inside a secret query value", "http://h.example:bad/?token=fake@fake2", "http://xxxxx"},
		{"scheme left out", "user:fake-pass@h.example:11434", "h.example:11434"},
		{"scheme left out with query", "h.example:11434/api?api_key=fake&b=2", "h.example:11434/api?b=2"},
		{"plain host without scheme", "h.example:11434", "h.example:11434"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RedactURL(tt.in)
			if got != tt.want {
				t.Errorf("RedactURL(%q) = %q, want %q", tt.in, got, tt.want)
			}
			if strings.Contains(got, "fake") {
				t.Errorf("RedactURL(%q) = %q still holds a credential", tt.in, got)
			}
		})
	}
}

func TestRedactText(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		{"no url", "load config: bad yaml", "load config: bad yaml"},
		{"quoted url in error", `ollama provider unavailable: invalid endpoint "http://u:fake-pass@h.example:11434": bad port`, `ollama provider unavailable: invalid endpoint "http://h.example:11434": bad port`},
		{"url followed by colon", "cannot reach http://u:fake-pass@h.example/api/tags: dial tcp: refused", "cannot reach http://h.example/api/tags: dial tcp: refused"},
		{"two urls", "a http://u:fake1@a.example b postgres://e:fake2@b.example/d", "a http://a.example b postgres://e:xxxxx@b.example/d"},
		{"keyword dsn in error", "open database: cannot parse `host=h password=fake-pass`: bad", "open database: cannot parse `host=h password=xxxxx`: bad"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Redact(tt.in); got != tt.want {
				t.Errorf("Redact(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestHandlerRedactsMessagesAndFields(t *testing.T) {
	var file, stderr bytes.Buffer
	buf := NewRingBuffer(4)
	logger := slog.New(NewHandler(HandlerOptions{RingBuffer: buf, Stderr: &stderr, File: &file}))

	logger.Error("open database: postgres://edda:fake-pass@db.example/edda",
		"endpoint", "http://u:fake-pass@h.example:11434",
		slog.Group("g", "dsn", "postgres://edda:fake-pass@db.example/edda"),
	)

	for name, out := range map[string]string{"file": file.String(), "stderr": stderr.String()} {
		if strings.Contains(out, "fake-pass") {
			t.Errorf("%s output holds a credential: %s", name, out)
		}
	}
	entries := buf.Snapshot()
	if len(entries) != 1 {
		t.Fatalf("entries: got %d, want 1", len(entries))
	}
	got := entries[0]
	if got.Message != "open database: postgres://edda:xxxxx@db.example/edda" {
		t.Errorf("message = %q", got.Message)
	}
	for _, f := range got.Fields {
		if strings.Contains(f.Value, "fake-pass") {
			t.Errorf("field %s = %q holds a credential", f.Key, f.Value)
		}
	}
}

func TestCharmLoggerRedactsErrorf(t *testing.T) {
	buf := NewRingBuffer(4)
	handler := NewHandler(HandlerOptions{RingBuffer: buf, Stderr: io.Discard})
	logger := NewCharmLogger(NewBridgeWriter(slog.New(handler).WithGroup("app")))

	logger.Errorf("initialize llm provider: invalid endpoint %q", "http://u:fake-pass@h.example")

	entries := buf.Snapshot()
	if len(entries) != 1 {
		t.Fatalf("entries: got %d, want 1", len(entries))
	}
	if strings.Contains(entries[0].Message, "fake-pass") {
		t.Errorf("message holds a credential: %q", entries[0].Message)
	}
}

func TestRedactURLError(t *testing.T) {
	cause := errors.New("connection refused")
	got := RedactURLError(&url.Error{Op: "Get", URL: "http://fake-token@h.example/?api_key=fake", Err: cause})
	if want := `Get "http://h.example/": connection refused`; got.Error() != want {
		t.Errorf("error = %q, want %q", got, want)
	}
	if !errors.Is(got, cause) {
		t.Error("the cause is no longer in the chain")
	}

	// net/url quotes pieces of an invalid URL, here part of the password.
	_, parseErr := url.Parse("http://user:fake/pass@h.example")
	got = RedactURLError(parseErr)
	if strings.Contains(got.Error(), "fake") {
		t.Errorf("error = %q holds a credential", got)
	}

	// A parse reason is kept when the URL held nothing to remove.
	_, parseErr = url.Parse("http://h.example:bad")
	if got = RedactURLError(parseErr); got.Error() != parseErr.Error() {
		t.Errorf("error = %q, want %q", got, parseErr)
	}

	got = RedactURLError(&url.Error{Op: "Get", URL: "fake-user:fake@h.example", Err: errors.New(`unsupported protocol scheme "fake-user"`)})
	if strings.Contains(got.Error(), "fake") {
		t.Errorf("error = %q holds a credential", got)
	}

	got = RedactURLError(&url.Error{Op: "Get", URL: "http://h.example", Err: errors.New("bad Location header http://u:fake@b.example")})
	if strings.Contains(errors.Unwrap(got).Error(), "fake") {
		t.Errorf("cause = %q holds a credential", errors.Unwrap(got))
	}

	other := errors.New("other")
	if got := RedactURLError(other); got != other {
		t.Errorf("a non-URL error changed: %v", got)
	}
	if RedactURLError(nil) != nil {
		t.Error("nil changed")
	}
}
