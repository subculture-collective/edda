package memory

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// assertNoSecret fails when err, or any error it wraps, holds one of the fake
// credentials, all of which contain "fake".
func assertNoSecret(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error")
	}
	for e := err; e != nil; e = errors.Unwrap(e) {
		for _, out := range []string{e.Error(), fmt.Sprintf("%#v", e)} {
			if strings.Contains(out, "fake") {
				t.Errorf("%T holds a credential: %s", e, out)
			}
		}
	}
}

func TestOllamaEmbedderErrorsRedactEndpoint(t *testing.T) {
	failing := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, "denied for http://fake-user:fake-pass@backend.example/?token=fake-token")
	}))
	defer failing.Close()
	broker := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = fmt.Fprint(w, "data: {\"request_id\":\"r1\",\"status\":\"ollama_unavailable\",\"message\":\"cannot reach http://fake-user:fake-pass@backend.example/\"}\n\n")
	}))
	defer broker.Close()
	withSecrets := func(base string) string {
		return strings.Replace(base, "://", "://fake-user:fake-pass@", 1) + "?api_key=fake-key"
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	for name, tc := range map[string]struct {
		ctx  context.Context
		base string
	}{
		"connection refused":    {context.Background(), withSecrets("http://127.0.0.1:1")},
		"canceled":              {canceled, withSecrets("http://127.0.0.1:1")},
		"invalid port":          {context.Background(), "http://fake-user:fake-pass@h.example:bad"},
		"scheme left out":       {context.Background(), "fake-user:fake-pass@h.example:11434"},
		"invalid fragment":      {context.Background(), "http://fake-user:fake-pass@h.example/#%zz"},
		"slash in the password": {context.Background(), "http://fake-user:fake/pass@h.example"},
		"http status":           {context.Background(), withSecrets(failing.URL)},
		"broker":                {context.Background(), withSecrets(broker.URL)},
	} {
		t.Run(name, func(t *testing.T) {
			embedder := NewOllamaEmbedder(tc.base, "m")
			_, err := embedder.Embed(tc.ctx, "text")
			assertNoSecret(t, err)
			_, err = embedder.BatchEmbed(tc.ctx, []string{"a", "b"})
			assertNoSecret(t, err)
		})
	}
}
