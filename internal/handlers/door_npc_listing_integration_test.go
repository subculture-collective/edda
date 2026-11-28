//go:build integration

package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strings"
	"testing"

	clog "github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/auth"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/handlers"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/pkg/api"
)

// Doors and magic objects are object rows in npcs. They must not appear where
// NPCs are treated as characters: the model's "NPCs Present" context and the
// NPC listings. A creature that bears magic is still a character.
func TestIntegrationSceneObjectRowsAreNotNPCs(t *testing.T) {
	f := newSpatialFixture(t)
	ctx := context.Background()
	door := f.addDoor("Oak door", true)
	f.doorFacts(door, f.actor, `"hidden":false,"state":"closed","locked":false,"distance_feet":5`)

	magicRow := func(name, bearer string) uuid.UUID {
		t.Helper()
		id := f.addDoor(name, true)
		_, err := testPool.Exec(ctx, `UPDATE npcs SET properties=jsonb_build_object('srd521_magic',jsonb_build_object('profile','detect-magic-v1','bearer',$1::text)) WHERE id=$2`, bearer, id)
		f.must(err)
		return id
	}
	chest := magicRow("Warded chest", "object")
	envoy := magicRow("Charmed envoy", "creature")
	characters := []string{f.target.String(), envoy.String()}
	slices.Sort(characters)
	objects := []string{"Oak door", door.String(), "Warded chest", chest.String()}

	state, err := game.NewStateLoader(statedb.New(testPool), testPool).Load(ctx, f.campaign)
	f.must(err)
	present := []string{}
	for _, npc := range state.NearbyNPCs {
		present = append(present, npc.ID.String())
	}
	slices.Sort(present)
	if !slices.Equal(present, characters) {
		t.Fatalf("prompt NPCs present: want the practice foe and the envoy, got %+v", state.NearbyNPCs)
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
		for _, object := range objects {
			if strings.Contains(body, object) {
				t.Fatalf("GET %s lists object %q as an NPC: %s", path, object, body)
			}
		}
		listed := []string{}
		for _, npc := range decodeJSON[[]api.NPCResponse](t, rec) {
			listed = append(listed, npc.ID)
		}
		slices.Sort(listed)
		if !slices.Equal(listed, characters) {
			t.Fatalf("GET %s: want the practice foe and the envoy, got %v", path, listed)
		}
	}

	// The door rules still read the row.
	if p := f.openDoor(door); !p.Available || len(p.Targets) != 1 || p.Targets[0].Label != "Oak door" {
		t.Fatalf("door preview after listing filter: %+v", p)
	}
}
