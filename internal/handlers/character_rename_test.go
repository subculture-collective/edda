package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/pkg/api"
)

type renameQuerier struct {
	statedb.Querier
	renamed []statedb.RenamePlayerCharacterParams
	rows    int64
	err     error
	name    string
}

func (q *renameQuerier) RenamePlayerCharacter(_ context.Context, arg statedb.RenamePlayerCharacterParams) (int64, error) {
	q.renamed = append(q.renamed, arg)
	return q.rows, q.err
}

func (q *renameQuerier) GetPlayerCharacterByCampaign(_ context.Context, campaignID pgtype.UUID) ([]statedb.PlayerCharacter, error) {
	return []statedb.PlayerCharacter{{ID: dbutil.ToPgtype(uuid.New()), CampaignID: campaignID, Name: q.name}}, nil
}

func (q *renameQuerier) GetCampaignByID(_ context.Context, id pgtype.UUID) (statedb.Campaign, error) {
	return statedb.Campaign{ID: id, RulesVersion: "legacy"}, nil
}

func characterRouter(q statedb.Querier) http.Handler {
	h := &CharacterHandlers{Queries: q, Logger: log.Default()}
	r := chi.NewRouter()
	r.Get("/campaigns/{id}/character", h.GetCharacter)
	r.Put("/campaigns/{id}/character/name", h.RenameCharacter)
	return r
}

func TestRenameCharacter(t *testing.T) {
	campaignID, characterID := uuid.New(), uuid.New()
	body := func(id, name string) string {
		return string(mustStartupJSON(t, api.CharacterRenameRequest{CharacterID: id, Name: name}))
	}
	for _, tc := range []struct {
		name, body string
		rows       int64
		err        error
		status     int
		message    string
		saved      string
	}{
		{name: "saves the trimmed name", body: body(characterID.String(), " \tIona \n"), rows: 1, status: http.StatusNoContent, saved: "Iona"},
		{name: "empty name", body: body(characterID.String(), ""), status: http.StatusBadRequest, message: "character name is required"},
		{name: "blank name", body: body(characterID.String(), " \t　"), status: http.StatusBadRequest, message: "character name is required"},
		{name: "invalid character id", body: body("nobody", "Iona"), status: http.StatusBadRequest, message: "invalid character id"},
		{name: "invalid body", body: "{", status: http.StatusBadRequest, message: "invalid request body"},
		{name: "not in this campaign", body: body(characterID.String(), "Iona"), status: http.StatusNotFound, message: "character not found", saved: "Iona"},
		{name: "failed write", body: body(characterID.String(), "Iona"), err: errors.New("private database detail"), status: http.StatusInternalServerError, message: "failed to rename character", saved: "Iona"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			q := &renameQuerier{rows: tc.rows, err: tc.err}
			rec := httptest.NewRecorder()
			characterRouter(q).ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/campaigns/"+campaignID.String()+"/character/name", strings.NewReader(tc.body)))
			if rec.Code != tc.status {
				t.Fatalf("status = %d, body = %s, want %d", rec.Code, rec.Body.String(), tc.status)
			}
			if tc.message != "" {
				if got := decodeStartupJSON[map[string]any](t, rec)["error"]; got != tc.message {
					t.Errorf("error = %v, want %q", got, tc.message)
				}
			}
			if tc.saved == "" {
				if len(q.renamed) != 0 {
					t.Errorf("wrote %+v, want no write", q.renamed)
				}
				return
			}
			want := statedb.RenamePlayerCharacterParams{ID: dbutil.ToPgtype(characterID), CampaignID: dbutil.ToPgtype(campaignID), Name: tc.saved}
			if len(q.renamed) != 1 || q.renamed[0] != want {
				t.Errorf("writes = %+v, want %+v", q.renamed, want)
			}
		})
	}
}

// needs_name tells the web client both that the name is blank and that this API
// can save a new one. An API without the rename route never sends it.
func TestGetCharacterFlagsABlankName(t *testing.T) {
	for name, want := range map[string]bool{"": true, " \t\n": true, "Mara": false} {
		rec := httptest.NewRecorder()
		characterRouter(&renameQuerier{name: name}).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/campaigns/"+uuid.NewString()+"/character", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("name %q: status = %d, body = %s", name, rec.Code, rec.Body.String())
		}
		flag, present := decodeStartupJSON[map[string]any](t, rec)["needs_name"]
		if present != want || want && flag != true {
			t.Errorf("name %q: needs_name = %v (present %v), want %v", name, flag, present, want)
		}
	}
}
