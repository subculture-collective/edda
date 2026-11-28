//go:build integration

package statedb_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

// The NPC writes match the campaign as well as the ID: another campaign's
// NPC is treated as missing and left unchanged.
func TestIntegrationNPCWritesRequireCampaign(t *testing.T) {
	ctx := context.Background()
	q := newTx(t)
	user := createUser(t, q, "npc-scope-user")
	current := createCampaign(t, q, user.ID)
	other := createCampaign(t, q, user.ID)
	npc := createNPC(t, q, other.ID)

	writes := map[string]func(campaign statedb.Campaign) (statedb.Npc, error){
		"KillNPC": func(c statedb.Campaign) (statedb.Npc, error) {
			return q.KillNPC(ctx, statedb.KillNPCParams{ID: npc.ID, CampaignID: c.ID})
		},
		"UpdateNPCDisposition": func(c statedb.Campaign) (statedb.Npc, error) {
			return q.UpdateNPCDisposition(ctx, statedb.UpdateNPCDispositionParams{ID: npc.ID, CampaignID: c.ID, Disposition: 90})
		},
		"UpdateNPC": func(c statedb.Campaign) (statedb.Npc, error) {
			return q.UpdateNPC(ctx, statedb.UpdateNPCParams{ID: npc.ID, CampaignID: c.ID, Name: "Renamed", Disposition: 90, Alive: false})
		},
	}
	for name, write := range writes {
		if _, err := write(current); !errors.Is(err, pgx.ErrNoRows) {
			t.Fatalf("%s from another campaign: err = %v, want no rows", name, err)
		}
		got, err := q.GetNPCByID(ctx, statedb.GetNPCByIDParams{ID: npc.ID})
		if err != nil {
			t.Fatal(err)
		}
		if got.Name != npc.Name || got.Disposition != npc.Disposition || got.Alive != npc.Alive {
			t.Fatalf("%s from another campaign changed the NPC: %+v", name, got)
		}
	}
	for name, write := range writes {
		if _, err := write(other); err != nil {
			t.Fatalf("%s in its own campaign: %v", name, err)
		}
	}
}
