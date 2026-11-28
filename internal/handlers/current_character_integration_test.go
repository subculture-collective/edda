//go:build integration

package handlers_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/export"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/journal"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

// Readers that show "the" character of a campaign pick the one the turn loader
// uses, even when several characters share a created_at (rows inserted in one
// transaction get the same now()).
func TestIntegrationCurrentCharacterReadersAgreeWithLoader(t *testing.T) {
	ctx := context.Background()
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	var campaign uuid.UUID
	must(testPool.QueryRow(ctx, `INSERT INTO campaigns(name,status,created_by) VALUES('Current character','active',$1) RETURNING id`, testUserID).Scan(&campaign))
	t.Cleanup(func() { _, _ = testPool.Exec(ctx, `DELETE FROM campaigns WHERE id=$1`, campaign) })

	// The UUIDs are fixed so the tie-break by id is tested both ways round.
	low, high := uuid.MustParse("00000000-0000-4000-8000-000000000001"), uuid.MustParse("ffffffff-ffff-4fff-bfff-ffffffffffff")
	locations := map[uuid.UUID]uuid.UUID{}
	for _, id := range []uuid.UUID{low, high} {
		var loc uuid.UUID
		must(testPool.QueryRow(ctx, `INSERT INTO locations(campaign_id,name,player_known) VALUES($1,$2,true) RETURNING id`, campaign, "Place "+id.String()[:4]).Scan(&loc))
		locations[id] = loc
	}
	_, err := testPool.Exec(ctx, `INSERT INTO player_characters(id,campaign_id,user_id,name,stats,hp,max_hp,experience,level,status,abilities,current_location_id)
		VALUES ($1,$3,$4,'Low','{}',5,5,0,1,'active','[]',$5), ($2,$3,$4,'High','{}',5,5,0,1,'active','[]',$6)`,
		low, high, campaign, testUserID, locations[low], locations[high])
	must(err)

	state, err := game.NewStateLoader(statedb.New(testPool), testPool).Load(ctx, campaign)
	must(err)
	want := state.Player.ID
	if want != high {
		t.Fatalf("loader picked %s; this test expects the higher id on a created_at tie", want)
	}

	exported, err := export.NewStore(testPool).GetPlayerCharacter(ctx, campaign)
	must(err)
	if exported.Name != "High" {
		t.Errorf("export picked %q, want the loader's character High", exported.Name)
	}

	reminder, err := journal.NewStore(testPool).ReadReturnReminder(ctx, campaign)
	must(err)
	if reminder.Scene == nil || reminder.Scene.SourceID != "known-location:"+locations[want].String() {
		t.Errorf("reminder scene = %+v, want the loader's character's location %s", reminder.Scene, locations[want])
	}
}
