//go:build integration

package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"git.subcult.tv/subculture-collective/edda/internal/auth"
	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/db"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/saves"
	"git.subcult.tv/subculture-collective/edda/internal/server"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/internal/world"
	"git.subcult.tv/subculture-collective/edda/pkg/api"
)

// Creation and rename trim the name, and only the campaign's owner can rename
// that campaign's characters, in both rules versions.
func TestIntegrationCharacterNames(t *testing.T) {
	ctx := context.Background()
	q := statedb.New(testPool)
	newUser := func() pgtype.UUID {
		user, err := q.CreateUser(ctx, "names-"+uuid.NewString())
		if err != nil {
			t.Fatal(err)
		}
		return user.ID
	}
	owner, stranger := newUser(), newUser()
	t.Cleanup(func() {
		_, _ = testPool.Exec(context.Background(), `DELETE FROM campaigns WHERE created_by = $1`, owner)
	})
	cfg := config.Config{}
	cfg.Server.JWTSecret = "integration-character-names"
	router := func() http.Handler {
		// A new provider answers one world build.
		return server.NewRouterWithProvider(log.New(io.Discard), &stubIntegrationEngine{}, q, &creationProvider{t: t}, testPool, testUserID, cfg, saves.NewStore(testPool))
	}
	request := func(user pgtype.UUID, method, path string, body any, want int) []byte {
		t.Helper()
		token, err := auth.GenerateToken(db.FromPgUUID(user), cfg.Server.JWTSecret, time.Hour)
		if err != nil {
			t.Fatal(err)
		}
		encoded, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		req := httptest.NewRequest(method, "/api/v1/campaigns"+path, strings.NewReader(string(encoded)))
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		router().ServeHTTP(rec, req)
		if rec.Code != want {
			t.Fatalf("%s %s: %d %s, want %d", method, path, rec.Code, rec.Body.String(), want)
		}
		return rec.Body.Bytes()
	}

	for _, version := range []domain.RulesVersion{domain.RulesVersionLegacy, domain.RulesVersionSRD521} {
		t.Run(string(version), func(t *testing.T) {
			var built api.WorldBuildResponse
			if err := json.Unmarshal(request(owner, http.MethodPost, "/start/world", api.WorldBuildRequest{
				Name: "Padded " + string(version), Profile: &api.CampaignProfile{Genre: "fantasy", Tone: "hopeful", Themes: []string{}},
				CharacterProfile: &world.CharacterProfile{Name: " \t Mara \n", Concept: "human fighter", Background: "Soldier"},
				RulesVersion:     string(version),
			}, http.StatusOK), &built); err != nil {
				t.Fatal(err)
			}
			base := "/" + built.Campaign.ID
			character := func() api.CharacterResponse {
				t.Helper()
				var got api.CharacterResponse
				if err := json.Unmarshal(request(owner, http.MethodGet, base+"/character", nil, http.StatusOK), &got); err != nil {
					t.Fatal(err)
				}
				return got
			}
			created := character()
			if created.Name != "Mara" || created.NeedsName {
				t.Fatalf("created name = %q, needs_name = %v", created.Name, created.NeedsName)
			}

			// Blank the name as the broken form did.
			if _, err := testPool.Exec(ctx, `UPDATE player_characters SET name = '' WHERE id = $1`, created.ID); err != nil {
				t.Fatal(err)
			}
			if !character().NeedsName {
				t.Fatal("a blank-named character is not flagged")
			}
			rename := api.CharacterRenameRequest{CharacterID: created.ID, Name: "  Iona  "}
			request(stranger, http.MethodPut, base+"/character/name", rename, http.StatusNotFound)
			// The owner cannot reach the character through another campaign.
			other, err := q.CreateCampaign(ctx, statedb.CreateCampaignParams{Name: "Other", Status: "active", CreatedBy: owner, Themes: []string{}})
			if err != nil {
				t.Fatal(err)
			}
			request(owner, http.MethodPut, "/"+db.FromPgUUID(other.ID).String()+"/character/name", rename, http.StatusNotFound)
			request(owner, http.MethodPut, base+"/character/name", api.CharacterRenameRequest{CharacterID: created.ID, Name: " \n"}, http.StatusBadRequest)
			if got := character(); got.Name != "" {
				t.Fatalf("refused renames changed the name to %q", got.Name)
			}
			request(owner, http.MethodPut, base+"/character/name", rename, http.StatusNoContent)
			if got := character(); got.Name != "Iona" || got.NeedsName {
				t.Fatalf("renamed name = %q, needs_name = %v", got.Name, got.NeedsName)
			}
		})
	}
}
