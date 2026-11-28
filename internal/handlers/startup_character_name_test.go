package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/charmbracelet/log"

	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/pkg/api"
)

func TestBuildWorld_BlankCharacterNameReturnsBadRequest(t *testing.T) {
	for _, version := range []string{"", "legacy", "srd-5.2.1"} {
		for _, name := range []string{"", "   ", "\t\n"} {
			queries := &startupStubQuerier{createCampaignFn: func(context.Context, statedb.CreateCampaignParams) (statedb.Campaign, error) {
				t.Error("a campaign was created")
				return statedb.Campaign{}, nil
			}}
			// The provider has no responses, so any model request fails the test.
			h := NewStartupHandlers(newStartupScriptedProvider(t), queries, log.Default(), nil)
			body := mustStartupJSON(t, api.WorldBuildRequest{
				Name: "Ashfall", Profile: &api.CampaignProfile{Genre: "fantasy"}, RulesVersion: version,
				CharacterProfile: &api.CharacterProfile{Name: name, Concept: "human fighter", Background: "Soldier"},
			})
			req := httptest.NewRequest(http.MethodPost, "/api/v1/campaigns/start/world", bytes.NewReader(body))
			rec := httptest.NewRecorder()
			newStartupRouter(h, true).ServeHTTP(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Errorf("rules %q, name %q: status = %d, body = %s", version, name, rec.Code, rec.Body.String())
				continue
			}
			if got := decodeStartupJSON[map[string]any](t, rec)["error"]; got != "character name is required" {
				t.Errorf("rules %q, name %q: error = %v", version, name, got)
			}
		}
	}
}
