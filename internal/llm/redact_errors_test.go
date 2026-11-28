package llm

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
// credentials, either in its message or in its own fields. Every fake
// credential in these tests contains "fake".
func assertNoSecret(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error")
	}
	var walk func(error)
	walk = func(e error) {
		if e == nil {
			return
		}
		for _, out := range []string{e.Error(), fmt.Sprintf("%#v", e)} {
			if strings.Contains(out, "fake") {
				t.Errorf("%T holds a credential: %s", e, out)
			}
		}
		switch w := e.(type) {
		case interface{ Unwrap() error }:
			walk(w.Unwrap())
		case interface{ Unwrap() []error }:
			for _, inner := range w.Unwrap() {
				walk(inner)
			}
		}
	}
	walk(err)
}

// withSecrets adds fake credentials to base as userinfo and a query parameter.
func withSecrets(base string) string {
	return strings.Replace(base, "://", "://fake-user:fake-pass@", 1) + "?api_key=fake-key"
}

const echoedSecretURL = "http://fake-user:fake-pass@backend.example/run?token=fake-token"

// secretProviders builds every HTTP provider against base.
func secretProviders(t *testing.T) map[string]func(base string) Provider {
	t.Helper()
	chatGPT := func(base string, expiresIn int) *ChatGPTClient {
		tokens, err := LoadChatGPTTokenStore(writeTestCreds(t, expiresIn))
		if err != nil {
			t.Fatal(err)
		}
		c := NewChatGPTClient(tokens, "m", "")
		c.baseURL = base
		return c
	}
	return map[string]func(string) Provider{
		"ollama": func(base string) Provider { return NewOllamaClient(base, "m") },
		"claude": func(base string) Provider { return NewClaudeClient(base, "k", "m") },
		"openrouter": func(base string) Provider {
			c := NewOpenRouterClient("k", "m")
			c.baseURL = base
			return c
		},
		"chatgpt": func(base string) Provider { return chatGPT(base, 3600) },
		"chatgpt token refresh": func(base string) Provider {
			c := chatGPT("http://127.0.0.1:1", -1)
			c.tokens.tokenURL = base
			return c
		},
		"claudecode": func(base string) Provider { return NewClaudeCodeClient(base, "t") },
	}
}

func assertProviderErrorsHoldNoSecret(t *testing.T, ctx context.Context, p Provider) {
	t.Helper()
	messages := []Message{{Role: RoleUser, Content: "hi"}}
	_, err := p.Complete(ctx, messages, nil)
	assertNoSecret(t, err)
	ch, err := p.Stream(ctx, messages, nil)
	if err != nil {
		assertNoSecret(t, err)
		return
	}
	for range ch {
	}
}

func TestProviderErrorsRedactRequestURL(t *testing.T) {
	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	unreachable := map[string]struct {
		ctx  context.Context
		base string
	}{
		"connection refused":            {context.Background(), withSecrets("http://127.0.0.1:1")},
		"token as user name":            {context.Background(), "http://fake-token@127.0.0.1:1"},
		"canceled":                      {canceled, withSecrets("http://127.0.0.1:1")},
		"invalid port":                  {context.Background(), "http://fake-user:fake-pass@h.example:bad?api_key=fake-key"},
		"invalid fragment":              {context.Background(), "http://fake-user:fake-pass@h.example/#%zz"},
		"scheme left out":               {context.Background(), "fake-user:fake-pass@h.example:11434"},
		"slash in the password":         {context.Background(), "http://fake-user:fake/pass@h.example"},
		"question mark in the password": {context.Background(), "http://fake-user:fake?pass@h.example"},
		"digits then slash in password": {context.Background(), "http://fake-user:1/fake-pass@127.0.0.1:1"},
	}
	for name, build := range secretProviders(t) {
		for caseName, tc := range unreachable {
			t.Run(name+"/"+caseName, func(t *testing.T) {
				assertProviderErrorsHoldNoSecret(t, tc.ctx, build(tc.base))
			})
		}
	}
}

func TestProviderStatusErrorsRedactRequestURL(t *testing.T) {
	statuses := []int{200, 400, 401, 403, 404, 429, 500, 502, 503, 504, 418}
	for name, build := range secretProviders(t) {
		for _, status := range statuses {
			t.Run(fmt.Sprintf("%s/%d", name, status), func(t *testing.T) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(status)
					if status == http.StatusOK {
						_, _ = fmt.Fprint(w, "not json")
						return
					}
					// A proxy that quotes the address it failed to reach.
					_, _ = fmt.Fprintf(w, `{"error":{"type":"x","message":"model not found at %s"}}`, echoedSecretURL)
				}))
				defer srv.Close()
				assertProviderErrorsHoldNoSecret(t, context.Background(), build(withSecrets(srv.URL)))
			})
		}
		t.Run(name+"/plain text body", func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadGateway)
				_, _ = fmt.Fprintf(w, "cannot reach %s", echoedSecretURL)
			}))
			defer srv.Close()
			assertProviderErrorsHoldNoSecret(t, context.Background(), build(withSecrets(srv.URL)))
		})
	}
}

func TestProviderTransportErrorsRedactRedirectTarget(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Location", "http://fake-user:fake-pass@h.example:bad/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}))
	defer srv.Close()
	for name, build := range secretProviders(t) {
		t.Run(name, func(t *testing.T) {
			assertProviderErrorsHoldNoSecret(t, context.Background(), build(withSecrets(srv.URL)))
		})
	}
}

func TestOllamaBrokerErrorRedactsMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = fmt.Fprintf(w, "data: {\"request_id\":\"r1\",\"status\":\"ollama_unavailable\",\"message\":\"cannot reach %s\"}\n\n", echoedSecretURL)
	}))
	defer srv.Close()
	_, err := NewOllamaClient(withSecrets(srv.URL), "m").Complete(context.Background(), []Message{{Role: RoleUser, Content: "hi"}}, nil)
	assertNoSecret(t, err)
	var connErr *ErrConnection
	if !errors.As(err, &connErr) {
		t.Fatalf("error = %T, want *ErrConnection", err)
	}
}

func TestOllamaEndpointValidationErrorsRedactEndpoint(t *testing.T) {
	failing := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, "denied for %s", echoedSecretURL)
	}))
	defer failing.Close()

	for name, endpoint := range map[string]string{
		"unreachable":           withSecrets("http://127.0.0.1:1"),
		"http status":           withSecrets(failing.URL),
		"invalid port":          "http://fake-user:fake-pass@h.example:bad",
		"slash in the password": "http://fake-user:fake/pass@h.example",
		"wrong scheme":          "ftp://fake-user:fake-pass@h.example?token=fake-token",
		"scheme left out":       "fake-user:fake-pass@h.example:11434",
		"relative":              "h.example?api_key=fake-key",
	} {
		t.Run(name, func(t *testing.T) {
			assertNoSecret(t, validateOllamaEndpoint(endpoint, ""))
		})
	}
}

// Redaction must leave the error chain usable.
func TestRedactedProviderErrorsKeepTheirCause(t *testing.T) {
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	for name, build := range secretProviders(t) {
		t.Run(name, func(t *testing.T) {
			_, err := build(withSecrets("http://127.0.0.1:1")).Complete(canceled, []Message{{Role: RoleUser, Content: "hi"}}, nil)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("errors.Is(%v, context.Canceled) = false", err)
			}
			var timeout *ErrTimeout
			if !errors.As(err, &timeout) {
				t.Errorf("error = %T, want *ErrTimeout", err)
			}
		})
	}
}

// A reply cut off inside a URL has no @ left to mark its credentials.
func TestProviderErrorsRedactCutOffReplies(t *testing.T) {
	const cut = "proxy failed for http://fake-user:fake-pa"
	replies := map[string]string{
		"plain text":   cut,
		"cut off json": `{"error":{"message":"model not found: ` + cut,
		// A sender that cut its message short before encoding it. Only a client
		// that decodes the message can tell.
		"cut off message": `{"error":{"type":"x","message":"model not found: ` + cut + `"}}`,
	}
	decodesMessage := map[string]bool{"claude": true, "openrouter": true, "chatgpt": true}
	for name, build := range secretProviders(t) {
		for replyName, reply := range replies {
			if replyName == "cut off message" && !decodesMessage[name] {
				continue
			}
			for _, status := range []int{404, 502} {
				t.Run(fmt.Sprintf("%s/%s/%d", name, replyName, status), func(t *testing.T) {
					srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.WriteHeader(status)
						_, _ = fmt.Fprint(w, reply)
					}))
					defer srv.Close()
					assertProviderErrorsHoldNoSecret(t, context.Background(), build(srv.URL))
				})
			}
		}
	}
	t.Run("broker message", func(t *testing.T) {
		assertNoSecret(t, brokerTerminalError("http://h.example", brokerStatusEvent{Status: brokerStatusOllamaUnavail, Message: cut}))
	})
	t.Run("start-up check", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = fmt.Fprint(w, cut)
		}))
		defer srv.Close()
		assertNoSecret(t, validateOllamaEndpoint(srv.URL, ""))
	})
}
