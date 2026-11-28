package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/llm"
	"git.subcult.tv/subculture-collective/edda/pkg/api"
)

// providerFailure is a provider error that names an internal endpoint and
// quotes the upstream reply.
var providerFailure = &llm.ErrConnection{
	URL: "http://models.internal:11434/api/chat",
	Err: fmt.Errorf("upstream replied: backend.internal unreachable"),
}

func TestProcessAction_LLMOutputErrorHidesProviderDetail(t *testing.T) {
	for _, sentinel := range []error{engine.ErrUnresolvedDurableClaims, engine.ErrEmptyTurnResponse} {
		t.Run(sentinel.Error(), func(t *testing.T) {
			// The shape the durable-claim repair step returns.
			eng := &stubEngine{turnErr: fmt.Errorf("%w: repair call failed: %v", sentinel, providerFailure)}
			h := &ActionHandlers{Engine: eng, Logger: log.New(io.Discard)}
			body, _ := json.Marshal(api.ActionRequest{Input: "listen"})
			req := httptest.NewRequest(http.MethodPost, "/campaigns/"+uuid.NewString()+"/action", bytes.NewReader(body))
			rec := httptest.NewRecorder()
			newActionRouter(h).ServeHTTP(rec, req)

			if rec.Code != http.StatusBadGateway {
				t.Fatalf("status = %d, want 502; body: %s", rec.Code, rec.Body.String())
			}
			var resp map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if resp["error"] != sentinel.Error() {
				t.Errorf("error = %q, want %q", resp["error"], sentinel)
			}
		})
	}
}

func TestEngineStateChangesToAPIHidesWarningDetail(t *testing.T) {
	warning := "ollama embed connection error (http://models.internal:11434/api/embed): refused"
	factID := uuid.NewString()
	changes := engine.StateChangesFromAppliedToolCalls([]engine.AppliedToolCall{{
		Tool:   "establish_fact",
		Result: mustJSON(t, map[string]any{"id": factID, "fact": "The well is dry.", "memory_warning": warning, "time_warning": ""}),
	}})
	got := engineStateChangesToAPI(changes)
	if len(got) != 1 {
		t.Fatalf("state changes = %d, want 1", len(got))
	}
	data, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "models.internal") {
		t.Errorf("response holds the warning detail: %s", data)
	}
	newValue, _ := got[0].Details["new_value"].(map[string]any)
	for name, details := range map[string]map[string]any{"details": got[0].Details, "new_value": newValue} {
		if details["memory_warning"] != withheldWarning {
			t.Errorf("%s memory_warning = %v, want %q", name, details["memory_warning"], withheldWarning)
		}
		if details["time_warning"] != "" {
			t.Errorf("%s time_warning = %v, want it left empty", name, details["time_warning"])
		}
		if details["fact"] != "The well is dry." {
			t.Errorf("%s fact = %v", name, details["fact"])
		}
	}
}
