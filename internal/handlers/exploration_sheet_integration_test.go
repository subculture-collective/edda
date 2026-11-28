//go:build integration

package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/handlers"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/pkg/api"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// characterSheet reads the campaign's sheet through the web API handler.
func characterSheet(t *testing.T, campaign uuid.UUID) *srd521.CharacterSheet {
	t.Helper()
	router := chi.NewRouter()
	characters := &handlers.CharacterHandlers{Queries: statedb.New(testPool), Pool: testPool, Logger: log.New(t.Output())}
	router.Get("/campaigns/{id}/character", characters.GetCharacter)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/campaigns/"+campaign.String()+"/character", nil))
	var response api.CharacterResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil || response.Sheet == nil {
		t.Fatalf("sheet unavailable: %d %s %v", recorder.Code, recorder.Body, err)
	}
	return response.Sheet
}
