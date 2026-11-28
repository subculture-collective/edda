//go:build integration

package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	clog "github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"

	"git.subcult.tv/subculture-collective/edda/internal/auth"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/handlers"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/pkg/api"
)

// A door is an object row in npcs. It must not appear where NPCs are treated
// as characters: the model's "NPCs Present" context and the NPC listings.
func TestIntegrationDoorRowsAreNotNPCs(t *testing.T) {
	f := newSpatialFixture(t)
	ctx := context.Background()
	door := f.addDoor("Oak door", true)
	f.doorFacts(door, f.actor, `"hidden":false,"state":"closed","locked":false,"distance_feet":5`)

	state, err := game.NewStateLoader(statedb.New(testPool), testPool).Load(ctx, f.campaign)
	f.must(err)
	if len(state.NearbyNPCs) != 1 || state.NearbyNPCs[0].ID != f.target {
		t.Fatalf("prompt NPCs present: want only the practice foe, got %+v", state.NearbyNPCs)
	}

	logger := clog.New(os.Stderr)
	logger.SetLevel(clog.ErrorLevel)
	worldH := &handlers.WorldHandlers{Queries: statedb.New(testPool), Logger: logger}
	r := chi.NewRouter()
	r.Use(auth.NewNoOpMiddleware(testUserID).Authenticate)
	r.Get("/campaigns/{id}/npcs", worldH.ListNPCs)
	r.Get("/campaigns/{id}/npcs/encountered", worldH.ListEncounteredNPCs)

	for _, path := range []string{"/npcs", "/npcs/encountered"} {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/campaigns/"+f.campaign.String()+path, nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("GET %s: %d %s", path, rec.Code, rec.Body.String())
		}
		body := rec.Body.String()
		if strings.Contains(body, "Oak door") || strings.Contains(body, door.String()) {
			t.Fatalf("GET %s lists the door as an NPC: %s", path, body)
		}
		npcs := decodeJSON[[]api.NPCResponse](t, rec)
		if len(npcs) != 1 || npcs[0].ID != f.target.String() {
			t.Fatalf("GET %s: want only the practice foe, got %+v", path, npcs)
		}
	}

	// The door rules still read the row.
	if p := f.openDoor(door); !p.Available || len(p.Targets) != 1 || p.Targets[0].Label != "Oak door" {
		t.Fatalf("door preview after listing filter: %+v", p)
	}
}
