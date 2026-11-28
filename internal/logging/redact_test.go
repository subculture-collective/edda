package logging

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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

		// Secret query parameter names: the suffix rule, in any case.
		{"suffix rule", "http://h.example/x?x-api-key=fake1&subscription-key=fake2&private_token=fake3&X-Amz-Signature=fake4&model=a", "http://h.example/x?model=a"},
		{"suffix rule in any case", "http://h.example/x?SessionToken=fake1&CLIENT_SECRET=fake2&DbPassword=fake3&Ocp-Apim-Subscription-KEY=fake4&X-Amz-Security-Token=fake5&page=2", "http://h.example/x?page=2"},
		{"escaped query name", "http://h.example/x?x%2Dapi%2Dkey=fake1&page=2", "http://h.example/x?page=2"},
		{"exact names", "http://h.example/x?passwd=fake1&pass=fake2&auth=fake3&Authorization=fake4&sig=fake5&page=2", "http://h.example/x?page=2"},
		{"ordinary params survive", "http://h.example/x?model=a&page=2&max_tokens=5&keyword=b&tokenizer=c&author=d&signed=e", "http://h.example/x?model=a&page=2&max_tokens=5&keyword=b&tokenizer=c&author=d&signed=e"},

		// Unescaped userinfo that url.Parse reads as host, port and path.
		{"digits then slash in password", "http://fake-user:12/fake34@h.example", "http://h.example"},
		{"digits then slash in password, with path", "http://fake-user:12/fake34@h.example:11434/api/chat", "http://h.example:11434/api/chat"},
		{"digits then question mark in password", "http://fake-user:12?fake34@h.example", "http://xxxxx"},
		{"digits then hash in password", "http://fake-user:12#fake34@h.example", "http://xxxxx"},
		{"at sign then slash in password", "http://fake-user:fake@ss/fake-word@h.example/x", "http://h.example/x"},
		{"slash in a token used as user name", "http://fake/fake-token@h.example", "http://h.example"},
		{"slash in the user name", "http://fake.us/er:fake-pass@h.example", "http://h.example"},
		{"postgres digits then slash in password", "postgres://edda:12/fake@db.example/edda", "postgres://db.example/edda"},
		{"scheme left out, digits then slash", "fake-user:12/fake34@h.example:11434", "h.example:11434"},
		{"valid userinfo with a path", "https://fake-user:fake-pass@h.example/a/b?page=2", "https://h.example/a/b?page=2"},
		{"at sign in a path", "https://cdn.example/npm/pkg@1.2.3/file.js", "https://cdn.example/npm/pkg@1.2.3/file.js"},
		{"at sign in a localhost path", "http://localhost/users/a@b.example", "http://localhost/users/a@b.example"},
		{"at sign in a query", "http://h.example/search?q=a@b.example&page=2", "http://h.example/search?q=a@b.example&page=2"},
		// Accepted over-redaction: this valid URL has the same shape as a dotted
		// user name with a digit password (the rows below), so it is cut too.
		{"at sign in a path after a port", "http://h.example:8080/users/a@b.example", "http://b.example"},
		{"dotted user name, digits and a slash", "http://fake.user:12/fake34@h.example", "http://h.example"},
		{"dotted user name, digits and a question mark", "http://fake.user:12?fake34@h.example", "http://xxxxx"},
		{"dotted database user, digits and a hash", "postgres://fake.user:2024#fake@db.example/edda", "postgres://xxxxx"},
		{"dotted user name, scheme left out", "fake.user:12/fake34@h.example", "h.example"},

		// Fragments take the same parameter rule as queries.
		{"token in a fragment", "http://h.example/cb#access_token=fake", "http://h.example/cb"},
		{"other fragment pairs survive", "http://h.example/cb#access_token=fake1&token_type=bearer&id_token=fake2&expires_in=3600", "http://h.example/cb#token_type=bearer&expires_in=3600"},
		{"fragment after a route", "https://h.example/#/callback?access_token=fake&state=a", "https://h.example/#/callback?state=a"},
		{"fragment right after a route slash", "https://h.example/#/access_token=fake&state=a", "https://h.example/#state=a"},
		{"escaped fragment name", "http://h.example/#access%5Ftoken=fake&a=b%26c", "http://h.example/#a=b%26c"},
		{"semicolons in a fragment", "http://h.example/#a=1;client_secret=fake;b=2", "http://h.example/#a=1;b=2"},
		{"question mark inside a fragment secret", "http://h.example/#access_token=fake1?fake2", "http://h.example/"},
		{"fragment and query together", "http://h.example/x?key=fake1&page=2#token=fake2", "http://h.example/x?page=2"},
		{"fragment kept when the query changes", "http://h.example/x?key=fake1#section-2", "http://h.example/x#section-2"},
		{"anchor survives", "http://h.example/docs#section-key", "http://h.example/docs#section-key"},
		{"ordinary fragment pairs survive", "http://h.example/docs#page=2&line=5", "http://h.example/docs#page=2&line=5"},
		{"scheme left out with a fragment", "h.example/cb#access_token=fake", "h.example/cb"},
		{"userinfo and a fragment", "https://fake-user:fake-pass@h.example/cb#access_token=fake", "https://h.example/cb"},

		// A semicolon separates query pairs too.
		{"semicolon before a secret", "http://h.example/x?model=a;token=fake;page=2", "http://h.example/x?model=a;page=2"},
		{"semicolon after a secret", "http://h.example/x?api_key=fake;model=a&page=2", "http://h.example/x?model=a&page=2"},
		{"semicolon inside a secret value", "http://h.example/x?password=fake1;fake2;fake3&page=2", "http://h.example/x?page=2"},
		{"mixed separators", "http://h.example/x?a=1&token=fake1;sig=fake2&b=2;key=fake3", "http://h.example/x?a=1&b=2"},
		{"semicolons without secrets survive", "http://h.example/x?a=1;b=2", "http://h.example/x?a=1;b=2"},

		{"amazon credential", "https://h.example/o?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=fake%2F20260101%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20260101T000000Z", "https://h.example/o?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20260101T000000Z"},
		{"google credential", "https://h.example/o?x-goog-credential=fake&page=2", "https://h.example/o?page=2"},
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
		{"url with a token in its fragment", "redirected to http://h.example/cb#access_token=fake: refused", "redirected to http://h.example/cb refused"},
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

// A reply from a provider can be cut off anywhere, so a URL that ends it
// without reaching the end of its host may be the start of a password.
func TestRedactReply(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		{"cut in a password", "cannot reach http://fake-user:fake-pa", "cannot reach http://xxxxx"},
		{"cut after a digit password", "cannot reach http://fake-user:12", "cannot reach http://xxxxx"},
		{"cut in a token", "cannot reach https://fake-tok", "cannot reach https://xxxxx"},
		{"cut in a database url", "open postgres://edda:fake-pa", "open postgres://xxxxx"},
		{"cut inside json", `{"error":{"message":"proxy failed for http://fake-user:fake-pa`, `{"error":{"message":"proxy failed for http://xxxxx`},
		{"cut after the scheme", "cannot reach http://", "cannot reach http://"},
		{"trailing newline", "cannot reach http://fake-user:fake-pa\n", "cannot reach http://xxxxx\n"},
		{"an ellipsis added by the sender", "cannot reach http://fake-user:fake-pa...", "cannot reach http://xxxxx..."},
		{"only the last url is cut", "a http://u:fake1@a.example/x b http://fake-user:fake2", "a http://a.example/x b http://xxxxx"},
		{"cut after a path", "cannot reach http://fake-user:fake-pass@h.example/ap", "cannot reach http://h.example/ap"},
		{"cut in a query value", "get http://h.example/x?private_token=fake-to", "get http://h.example/x"},
		{"host followed by a path", "cannot reach http://h.example/", "cannot reach http://h.example/"},
		{"host followed by text", "http://h.example refused", "http://h.example refused"},
		// A string in a reply can have been cut before it was encoded, so a URL
		// that ends one is masked as a URL that ends the reply is.
		{"host ending a json string", `{"url":"http://h.example"}`, `{"url":"http://xxxxx"}`},
		{"path ending a json string", `{"url":"http://h.example/"}`, `{"url":"http://h.example/"}`},
		{"cut in nested json", `{"error":{"detail":"x http://fake-user:fake-pa"}}`, `{"error":{"detail":"x http://xxxxx"}}`},
		{"cut deeper in arrays", `{"a":[{"b":["x http://fake-user:12","y http://fake-tok"]}]}`, `{"a":[{"b":["x http://xxxxx","y http://xxxxx"]}]}`},
		{"cut in a json key", `{"x http://fake-user:fake-pa":1}`, `{"x http://xxxxx":1}`},
		{"cut under a repeated key", `{"d":"x http://fake-tok","d":"ok"}`, `{"d":"x http://xxxxx","d":"ok"}`},
		{"cut in json that is itself cut", `{"error":{"detail":"x http://fake-user:fake-pa"},"trace":"ab`, `{"error":{"detail":"x http://xxxxx"},"trace":"ab`},
		{"cut in json after other text", `upstream said {"detail":"x http://fake-user:fake-pa…"}`, `upstream said {"detail":"x http://xxxxx…"}`},
		// In the raw text the escape reads as part of the host, and is masked too.
		{"cut before an encoded newline", `{"detail":"x http://fake-user:fake-pa\n"}`, `{"detail":"x http://xxxxx"}`},
		{"cut in json inside a json string", `{"body":"{\"detail\":\"x http://fake-user:fake-pa\"}"}`, `{"body":"{\"detail\":\"x http://xxxxx\"}"}`},
		{"cut in encoded json inside cut json", `{"body":"{\"detail\":\"x http://fake-user:fake-pa\"}`, `{"body":"{\"detail\":\"x http://xxxxx\"}`},
		{"escaped slashes", `{"detail":"x http:\/\/fake-user:fake-pass@h.example\/ap"}`, `{"detail":"x http://h.example/ap"}`},
		{"escaped ampersand", `{"detail":"get http://h.example/x?a=1\u0026token=fake"}`, `{"detail":"get http://h.example/x?a=1"}`},
		{"escaped at sign", `{"detail":"http://fake-user:fake-pass\u0040h.example/<a>"}`, `{"detail":"http://h.example/<a>"}`},
		{"json with nothing to redact keeps its escapes", `{"detail":"a\u0026b \/ http:\/\/h.example\/x"}`, `{"detail":"a\u0026b \/ http:\/\/h.example\/x"}`},
		{"no url", "model not found", "model not found"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RedactReply(tt.in); got != tt.want {
				t.Errorf("RedactReply(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// net/http accepts "user:12/34@host" as host "user" and port "12", so its cause
// names both pieces.
func TestRedactURLErrorAmbiguousUserinfo(t *testing.T) {
	raw := "http://fake-user:12/fake34@h.example"
	got := RedactURLError(&url.Error{Op: "Get", URL: raw, Err: errors.New("dial tcp: lookup fake-user: no such host")})
	if strings.Contains(got.Error(), "fake") {
		t.Errorf("error = %q holds a credential", got)
	}
	got = RedactURLError(&url.Error{Op: "Get", URL: raw, Err: context.Canceled})
	if !errors.Is(got, context.Canceled) {
		t.Errorf("error = %q no longer matches context.Canceled", got)
	}
	got = RedactURLError(&url.Error{Op: "Get", URL: raw, Err: fmt.Errorf("dial fake-user: %w", context.DeadlineExceeded)})
	if !errors.Is(got, context.DeadlineExceeded) || strings.Contains(got.Error(), "fake") {
		t.Errorf("error = %q, want a bare deadline error", got)
	}
}

// A reply can be cut inside a password after a /, ? or #, or can end a cut-off
// URL with a single quote. Each row cuts "http://fake-user:...@h.example".
func TestRedactReplyCutPasswords(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		{"single quote", `'x http://fake-user:fake-pa'`, `'x http://xxxxx'`},
		{"single quote without a colon", `'x http://fake-tok'`, `'x http://xxxxx'`},
		{"escaped single quote", `'x http://fake-tok\'`, `'x http://xxxxx\'`},
		{"single quote after an ellipsis", `'x http://fake-tok...'`, `'x http://xxxxx...'`},
		{"single quote in json", `{"d":"got 'http://fake-tok'"}`, `{"d":"got 'http://xxxxx'"}`},
		{"single quote inside the password", `x http://fake-user:fa'fake`, `x http://xxxxx'fake`},
		{"slash, in json", `{"d":"http://fake-user:fake/pa"}`, `{"d":"http://xxxxx"}`},
		{"slash at the end", "x http://fake-user:fake/pa", "x http://xxxxx"},
		{"nothing after the slash", "x http://fake-user:fake/", "x http://xxxxx"},
		{"question mark at the end", "x http://fake-user:fa?ke", "x http://xxxxx"},
		{"hash at the end", "x http://fake-user:fa#ke", "x http://xxxxx"},
		{"question mark, in json", `{"d":"x http://fake-user:fa?ke"}`, `{"d":"x http://xxxxx"}`},
		{"hash before a single quote", `'x http://fake-user:fa#ke'`, `'x http://xxxxx'`},
		{"slash before a space", "x http://fake-user:fa/ke refused", "x http://xxxxx refused"},
		{"slash before a space, no more text", "x http://fake-user:fa/ke\n", "x http://xxxxx\n"},
		{"cut before a space", "x http://fake-user:fake-pa refused", "x http://xxxxx refused"},
		{"slash, in cut json", `{"d":"x http://fake-user:fa/ke`, `{"d":"x http://xxxxx`},
		{"slash, in nested json", `{"body":"{\"d\":\"http://fake-user:fa/ke\"}"}`, `{"body":"{\"d\":\"http://xxxxx\"}"}`},
		{"slash then an ellipsis", "x http://fake-user:fa/ke…", "x http://xxxxx…"},
		{"digits too big for a port", "x http://fake.user:99999/fake", "x http://xxxxx"},
		{"digits then letters", "x http://fake.user:12fake/pa", "x http://xxxxx"},
		{"digits without a dot in the user name", "x http://fake-user:12/fake", "x http://xxxxx"},
		{"database url", "open postgres://fake-user:fa/ke", "open postgres://xxxxx"},
		// What was masked before stays masked: the new rules run last, on the
		// finished result. Run earlier, they would show part of each of these.
		{"encoded single quote in the password", `{"d":"http://fake-user:pa\u0027fake"}`, `{"d":"http://xxxxx"}`},
		{"encoded angle bracket in the password", `{"d":"http://fake-user:pa\u003cfake"}`, `{"d":"http://xxxxx"}`},
		{"keyword value past the url", `x http://fake-user:fa/password="fake"`, "x http://xxxxx"},
		// Accepted over-redaction: the check runs after the userinfo is removed.
		{"userinfo before a host without a dot", "get http://fake-user:fake-pass@ollama:11434/api now", "get http://xxxxx now"},
		// Accepted over-redaction: a host with no dot has the shape of a user name.
		{"host without a dot, with a port and a path", "x http://ollama:11434/api/chat refused", "x http://xxxxx refused"},
		{"complete url in single quotes", "cannot reach 'http://h.example'", "cannot reach 'http://xxxxx'"},

		// Ordinary complete URLs do not change.
		{"path", "get https://h.example/a/b?page=2#top refused", "get https://h.example/a/b?page=2#top refused"},
		{"port and path", "get http://h.example:8080/v1/models refused", "get http://h.example:8080/v1/models refused"},
		{"port and path at the end", "get http://h.example:65535/v1/models", "get http://h.example:65535/v1/models"},
		{"ip address, port and path", `{"d":"get http://127.0.0.1:11434/api/chat"}`, `{"d":"get http://127.0.0.1:11434/api/chat"}`},
		{"localhost, port and query", "get http://localhost:11434?page=2", "get http://localhost:11434?page=2"},
		{"ipv6, port and path", "get http://[::1]:11434/api/chat", "get http://[::1]:11434/api/chat"},
		{"ipv6 and path", "get http://[2001:db8::1]/api", "get http://[2001:db8::1]/api"},
		{"port before text", "http://h.example:8080 refused", "http://h.example:8080 refused"},
		{"port before punctuation", "dial (http://h.example:8080), http://localhost:11434: refused", "dial (http://h.example:8080), http://localhost:11434: refused"},
		{"port in single quotes with a path", "get 'http://h.example:8080/'", "get 'http://h.example:8080/'"},
		{"port before an escaped slash", `{"d":"http://h.example:8080\/x"} http://h.example:8080\/x`, `{"d":"http://h.example:8080\/x"} http://h.example:8080\/x`},
		{"colon in a path", "see https://h.example/wiki/Help:Contents now", "see https://h.example/wiki/Help:Contents now"},
		{"colon in a path, host without a dot", "see http://ollama/library/llama3:8b now", "see http://ollama/library/llama3:8b now"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RedactReply(tt.in); got != tt.want {
				t.Errorf("RedactReply(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
