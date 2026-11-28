package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// captureLogs sends slog's default logger to a plain JSON handler, which does
// not redact, so the output shows what the caller passed to it. Tests that use
// it do not run in parallel.
func captureLogs(t *testing.T) *bytes.Buffer {
	t.Helper()
	previous, writer, flags := slog.Default(), log.Writer(), log.Flags()
	t.Cleanup(func() {
		slog.SetDefault(previous)
		log.SetOutput(writer)
		log.SetFlags(flags)
	})
	var out bytes.Buffer
	slog.SetDefault(slog.New(slog.NewJSONHandler(&out, nil)))
	return &out
}

// assertOneStreamFailure checks that the logs hold exactly one record, an error
// naming the cause and holding no fake credential.
func assertOneStreamFailure(t *testing.T, out *bytes.Buffer, service, cause string) {
	t.Helper()
	text := strings.TrimSpace(out.String())
	if text == "" || strings.Count(text, "\n") != 0 {
		t.Fatalf("want one log record, got %q", text)
	}
	if strings.Contains(text, "fake") {
		t.Errorf("log holds a credential: %s", text)
	}
	var record map[string]any
	if err := json.Unmarshal([]byte(text), &record); err != nil {
		t.Fatal(err)
	}
	if record["level"] != "ERROR" || record["msg"] != "stream failed" {
		t.Errorf("record = %s", text)
	}
	group, _ := record[service].(map[string]any)
	if got, _ := group["error"].(string); !strings.Contains(got, cause) {
		t.Errorf("%s error = %q, want it to hold %q", service, got, cause)
	}
}

// streamServer replies with wire. With cutOff, it promises a longer body than
// it sends, so reading the reply fails.
func streamServer(t *testing.T, contentType, wire string, cutOff bool) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", contentType)
		if cutOff {
			w.Header().Set("Content-Length", fmt.Sprint(len(wire)+100))
		}
		_, _ = io.WriteString(w, wire)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func drainStream(t *testing.T, ctx context.Context, p Provider) (done bool) {
	t.Helper()
	ch, err := p.Stream(ctx, []Message{{Role: RoleUser, Content: "hi"}}, nil)
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}
	for chunk := range ch {
		done = done || chunk.Done
	}
	return done
}

const cutOffURL = "http://fake-user:fake-pa"

func TestOllamaStreamFailureIsLoggedOnce(t *testing.T) {
	const sse, ndjson = "text/event-stream", "application/x-ndjson"
	const partial = `{"message":{"content":"start"},"done":false}` + "\n"
	for name, tc := range map[string]struct {
		contentType, wire, cause string
		cutOff                   bool
	}{
		"invalid chunk":        {ndjson, partial + "{not json}\n", "decode stream chunk", false},
		"invalid broker event": {sse, `data: {"request_id":"r1","status":"queued","position":"x"}` + "\n\n", "decode broker status event", false},
		"invalid tool call":    {ndjson, `{"message":{"tool_calls":[{"function":{"name":"roll","arguments":"{bad"}}]}}` + "\n", "decode ollama tool arguments", false},
		"done sentinel first":  {sse, "data: " + partial + "\ndata: [DONE]\n\n", "[DONE] before a final reply", false},
		"ended early":          {ndjson, partial, "ended before a final reply", false},
		"read error":           {ndjson, partial, "unexpected EOF", true},
		"broker unavailable":   {sse, `data: {"request_id":"r1","status":"ollama_unavailable","message":"cannot reach ` + cutOffURL + `"}` + "\n\n", "ollama_unavailable", false},
		"broker dropped":       {sse, `data: {"request_id":"r1","status":"dropped_by_admin"}` + "\n\n", "dropped request by admin", false},
	} {
		t.Run(name, func(t *testing.T) {
			srv := streamServer(t, tc.contentType, tc.wire, tc.cutOff)
			out := captureLogs(t)
			if drainStream(t, context.Background(), NewOllamaClient(srv.URL, "m")) {
				t.Fatal("failed stream completed")
			}
			assertOneStreamFailure(t, out, "ollama", tc.cause)
		})
	}
}

func TestOpenRouterStreamFailureIsLoggedOnce(t *testing.T) {
	const partial = `data: {"choices":[{"delta":{"content":"start"}}]}` + "\n\n"
	for name, tc := range map[string]struct {
		wire, cause string
		cutOff      bool
	}{
		"error event":            {partial + `data: {"error":{"message":"upstream failed at http://fake-user:fake-pass@h.example/?x=1&token=fake","code":503}}` + "\n\ndata: [DONE]\n\n", "upstream failed at http://h.example/?x=1 (code 503)", false},
		"cut off error event":    {partial + `data: {"error":{"message":"upstream failed at ` + cutOffURL + `"}}` + "\n\n", "upstream failed at http://xxxxx", false},
		"error with a text code": {partial + `data: {"error":{"message":"overloaded","code":"server_error"}}` + "\n\n", "overloaded (code server_error)", false},
		"error of another shape": {partial + `data: {"error":"cannot reach ` + cutOffURL + `"}` + "\n\n", "cannot reach http://xxxxx", false},
		"finish reason error":    {partial + `data: {"choices":[{"delta":{},"finish_reason":"error"}]}` + "\n\n", `finish_reason "error"`, false},
		"invalid tool call":      {`data: {"choices":[{"delta":{"tool_calls":[{"index":0,"id":"c","function":{"name":"roll","arguments":"{bad"}}]},"finish_reason":"tool_calls"}]}` + "\n\n", `arguments for "roll"`, false},
		"ended early":            {partial, "ended before [DONE] or a finish_reason", false},
		"only invalid events":    {"data: {not json}\n\n", "invalid character", false},
		"read error":             {partial, "unexpected EOF", true},
	} {
		t.Run(name, func(t *testing.T) {
			srv := streamServer(t, "text/event-stream", tc.wire, tc.cutOff)
			c := NewOpenRouterClient("k", "m")
			c.baseURL = srv.URL
			out := captureLogs(t)
			if drainStream(t, context.Background(), c) {
				t.Fatal("failed stream completed")
			}
			assertOneStreamFailure(t, out, "openrouter", tc.cause)
		})
	}
}

func TestChatGPTStreamFailureIsLoggedOnce(t *testing.T) {
	const partial = `data: {"type":"response.output_text.delta","delta":"start"}` + "\n\n"
	for name, tc := range map[string]struct {
		wire, cause string
		cutOff      bool
	}{
		"failed response":     {partial + `data: {"type":"response.failed","response":{"error":{"code":"server_error","message":"upstream failed at http://fake-user:fake-pass@h.example/?x=1&token=fake"}}}` + "\n\n", "server_error: upstream failed at http://h.example/?x=1", false},
		"cut off error event": {partial + `data: {"type":"error","code":"server_error","message":"upstream failed at ` + cutOffURL + `"}` + "\n\n", "upstream failed at http://xxxxx", false},
		"incomplete response": {partial + `data: {"type":"response.incomplete","response":{"incomplete_details":{"reason":"max_output_tokens"}}}` + "\n\n", "response incomplete: max_output_tokens", false},
		"invalid event":       {partial + "data: {not json}\n\n", "decode stream event", false},
		"invalid output item": {`data: {"type":"response.output_item.done","item":"text"}` + "\n\n", "decode output item", false},
		"invalid tool call":   {`data: {"type":"response.output_item.done","item":{"type":"function_call","call_id":"c","name":"roll","arguments":"{bad"}}` + "\n\n", "decode arguments for roll", false},
		"ended early":         {partial, "ended before response.completed", false},
		"read error":          {partial, "unexpected EOF", true},
	} {
		t.Run(name, func(t *testing.T) {
			srv := streamServer(t, "text/event-stream", tc.wire, tc.cutOff)
			c := testChatGPTClient(t, srv.URL)
			out := captureLogs(t)
			if drainStream(t, context.Background(), c) {
				t.Fatal("failed stream completed")
			}
			assertOneStreamFailure(t, out, "chatgpt", tc.cause)
		})
	}
}

func testChatGPTClient(t *testing.T, base string) *ChatGPTClient {
	t.Helper()
	tokens, err := LoadChatGPTTokenStore(writeTestCreds(t, 3600))
	if err != nil {
		t.Fatal(err)
	}
	c := NewChatGPTClient(tokens, "m", "")
	c.baseURL = base
	return c
}

func TestClaudeStreamFailureIsLoggedOnce(t *testing.T) {
	event := func(name, data string) string { return "event: " + name + "\ndata: " + data + "\n\n" }
	partial := event("content_block_start", `{"index":0,"content_block":{"type":"text","text":""}}`) +
		event("content_block_delta", `{"index":0,"delta":{"type":"text_delta","text":"start"}}`)
	toolStart := event("content_block_start", `{"index":1,"content_block":{"type":"tool_use","id":"c","name":"roll","input":{}}}`)
	badArgs := event("content_block_delta", `{"index":1,"delta":{"type":"input_json_delta","partial_json":"{bad"}}`) +
		event("content_block_stop", `{"index":1}`)
	stop := event("message_stop", `{}`)
	for name, tc := range map[string]struct {
		wire, cause string
		cutOff      bool
	}{
		"error event":             {partial + event("error", `{"type":"error","error":{"type":"api_error","message":"upstream failed at http://fake-user:fake-pass@h.example/?x=1&token=fake"}}`) + stop, "api_error: upstream failed at http://h.example/?x=1", false},
		"cut off error event":     {partial + event("error", `{"type":"error","error":{"type":"overloaded_error","message":"upstream failed at `+cutOffURL+`"}}`), "overloaded_error: upstream failed at http://xxxxx", false},
		"error of another shape":  {partial + event("error", `{"type":"error","error":"cannot reach `+cutOffURL+`"}`), "cannot reach http://xxxxx", false},
		"refusal":                 {partial + event("message_delta", `{"delta":{"stop_reason":"refusal"}}`) + stop, "stop_reason refusal", false},
		"invalid message delta":   {partial + event("message_delta", `{not json}`) + stop, "decode message_delta event", false},
		"invalid block start":     {event("content_block_start", `{"index":0,"content_block":"text"}`), "ended before message_stop; skipped an event that did not decode", false},
		"invalid tool call":       {toolStart + badArgs + stop, `decode arguments for "roll"`, false},
		"stop inside a tool call": {toolStart + stop, "message_stop before a tool_use block ended", false},
		"ended early":             {partial, "ended before message_stop", false},
		"only invalid events":     {event("content_block_delta", `{not json}`), "invalid character", false},
		"read error":              {partial, "unexpected EOF", true},
	} {
		t.Run(name, func(t *testing.T) {
			srv := streamServer(t, "text/event-stream", tc.wire, tc.cutOff)
			out := captureLogs(t)
			if drainStream(t, context.Background(), NewClaudeClient(srv.URL, "k", "m")) {
				t.Fatal("failed stream completed")
			}
			assertOneStreamFailure(t, out, "claude", tc.cause)
		})
	}
}

// A finished stream and one the caller gave up on are not failures. A passed
// deadline is.
func TestStreamLoggingOnSuccessCancelAndDeadline(t *testing.T) {
	providers := map[string]struct {
		contentType, first, final string
		build                     func(base string) Provider
	}{
		"ollama": {
			"application/x-ndjson",
			`{"message":{"content":"a"},"done":false}` + "\n",
			`{"message":{"content":"b"},"done":true}` + "\n",
			func(base string) Provider { return NewOllamaClient(base, "m") },
		},
		"openrouter": {
			"text/event-stream",
			`data: {"choices":[{"delta":{"content":"a"}}]}` + "\n\n",
			"data: [DONE]\n\n",
			func(base string) Provider {
				c := NewOpenRouterClient("k", "m")
				c.baseURL = base
				return c
			},
		},
		"claude": {
			"text/event-stream",
			"event: content_block_delta\n" + `data: {"index":0,"delta":{"type":"text_delta","text":"a"}}` + "\n\n",
			"event: message_stop\ndata: {}\n\n",
			func(base string) Provider { return NewClaudeClient(base, "k", "m") },
		},
		"chatgpt": {
			"text/event-stream",
			`data: {"type":"response.output_text.delta","delta":"a"}` + "\n\n",
			"data: " + completedEvent + "\n\n",
			func(base string) Provider { return testChatGPTClient(t, base) },
		},
	}
	for name, tc := range providers {
		t.Run(name+"/success", func(t *testing.T) {
			srv := streamServer(t, tc.contentType, tc.first+tc.final, false)
			out := captureLogs(t)
			if !drainStream(t, context.Background(), tc.build(srv.URL)) {
				t.Fatal("stream did not complete")
			}
			if out.Len() != 0 {
				t.Errorf("a finished stream was logged: %s", out)
			}
		})
		stalled := func(t *testing.T) *httptest.Server {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tc.contentType)
				_, _ = io.WriteString(w, tc.first)
				w.(http.Flusher).Flush()
				<-r.Context().Done()
			}))
			t.Cleanup(srv.Close)
			return srv
		}
		t.Run(name+"/canceled", func(t *testing.T) {
			srv := stalled(t)
			out := captureLogs(t)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			ch, err := tc.build(srv.URL).Stream(ctx, nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			<-ch
			cancel()
			for range ch {
			}
			if out.Len() != 0 {
				t.Errorf("a canceled stream was logged: %s", out)
			}
		})
		t.Run(name+"/deadline", func(t *testing.T) {
			srv := stalled(t)
			out := captureLogs(t)
			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()
			if drainStream(t, ctx, tc.build(srv.URL)) {
				t.Fatal("failed stream completed")
			}
			assertOneStreamFailure(t, out, name, "deadline exceeded")
		})
	}
}
