package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/pkg/api"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type concentratingQuerier struct {
	statedb.Querier
	stats []byte
}

func (q concentratingQuerier) GetPlayerCharacterByCampaign(context.Context, pgtype.UUID) ([]statedb.PlayerCharacter, error) {
	return []statedb.PlayerCharacter{{ID: dbutil.ToPgtype(uuid.New()), Stats: q.stats, Hp: 10, MaxHp: 10, Level: 1}}, nil
}

func (q concentratingQuerier) GetCampaignByID(context.Context, pgtype.UUID) (statedb.Campaign, error) {
	return statedb.Campaign{RulesVersion: "srd-5.2.1"}, nil
}

// The sheet names the concentration spell, which saved state does not.
func TestCharacterSheetNamesConcentrationSpell(t *testing.T) {
	saved, err := os.ReadFile("../rules/srd521/testdata/cleric1.json")
	if err != nil {
		t.Fatal(err)
	}
	stats := strings.Replace(string(saved), `"somatic": true`, `"somatic": true, "concentration": {"id": "cast", "spell": "shield-of-faith", "started_at": 0, "expires_at": 600}`, 1)
	h := &CharacterHandlers{Queries: concentratingQuerier{stats: []byte(stats)}, Logger: log.New(io.Discard)}
	router := chi.NewRouter()
	router.Get("/campaigns/{id}/character", h.GetCharacter)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/campaigns/"+uuid.NewString()+"/character", nil))
	var response api.CharacterResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil || response.Sheet == nil || response.Sheet.Concentration == nil {
		t.Fatalf("status=%d body=%s err=%v", rec.Code, rec.Body.String(), err)
	}
	if got := response.Sheet.Concentration.SpellLabel; got != "Shield of Faith" {
		t.Fatalf("spell_label=%q", got)
	}
}
