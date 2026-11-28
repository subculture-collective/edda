//go:build integration

package handlers_test

import (
	"context"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/saves"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/internal/tools"
	"github.com/google/uuid"
)

// TestIntegrationLegacyClockWritersKeepSeconds checks the invariant the SRD
// rules clock relies on: writers that only know day, hour and minute move the
// absolute campaign second by exactly their whole-minute delta, so an expiry
// stored as an absolute second (Light, concentration) still resolves exactly.
func TestIntegrationLegacyClockWritersKeepSeconds(t *testing.T) {
	ctx := context.Background()
	q := statedb.New(testPool)
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	c, err := q.CreateCampaign(ctx, statedb.CreateCampaignParams{Name: "Clock seconds", RulesVersion: "srd-5.2.1", Status: "active", CreatedBy: dbutil.ToPgtype(testUserID)})
	must(err)
	campaign := dbutil.FromPgtype(c.ID)
	t.Cleanup(func() { _, _ = testPool.Exec(ctx, `DELETE FROM campaigns WHERE id=$1`, campaign) })
	pc, err := q.CreatePlayerCharacter(ctx, statedb.CreatePlayerCharacterParams{CampaignID: c.ID, UserID: dbutil.ToPgtype(testUserID), Name: "Mira", Stats: []byte(`{}`), Hp: 5, MaxHp: 10, Level: 1, Status: "active", Abilities: []byte(`[]`)})
	must(err)
	location := func(name string) uuid.UUID {
		t.Helper()
		var id uuid.UUID
		must(testPool.QueryRow(ctx, `INSERT INTO locations(campaign_id,name,player_known)VALUES($1,$2,true)RETURNING id`, campaign, name).Scan(&id))
		return id
	}
	gate, wood := location("Gate"), location("Wood")
	_, err = testPool.Exec(ctx, `INSERT INTO location_connections(from_location_id,to_location_id,bidirectional,travel_time,campaign_id)VALUES($1,$2,true,'2 hours 15 minutes',$3)`, gate, wood, campaign)
	must(err)
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET current_location_id=$1 WHERE id=$2`, gate, pc.ID)
	must(err)

	// The same absolute-second expression the rules code reads.
	now := func() int {
		t.Helper()
		var s int
		must(testPool.QueryRow(ctx, `SELECT (day-1)*86400+hour*3600+minute*60+second FROM campaign_time WHERE campaign_id=$1`, campaign).Scan(&s))
		return s
	}
	// Day 1, 23:59:30: a Light cast 59m30s ago has 30 seconds left.
	_, err = testPool.Exec(ctx, `INSERT INTO campaign_time(campaign_id,day,hour,minute,second)VALUES($1,1,23,59,30)`, campaign)
	must(err)
	start := now()
	lightExpires := start + 30

	toolCtx := tools.WithCurrentCampaignID(ctx, campaign)
	step := func(name string, delta int, run func() error) {
		t.Helper()
		before := now()
		must(run())
		if got := now(); got != before+delta {
			t.Fatalf("%s moved clock %d -> %d, want +%d", name, before, got, delta)
		}
	}
	step("advance_time 1 minute across midnight", 60, func() error {
		_, err := tools.NewAdvanceTimeHandler(testPool).Handle(toolCtx, map[string]any{"hours": 0, "minutes": 1})
		return err
	})
	if now() < lightExpires {
		t.Fatalf("Light expiring at %d still active at %d", lightExpires, now())
	}
	if got := now() - lightExpires; got != 30 {
		t.Fatalf("Light expired %d seconds ago, want 30", got)
	}
	step("advance_time 1 hour", 3600, func() error {
		_, err := tools.NewAdvanceTimeHandler(testPool).Handle(toolCtx, map[string]any{"hours": 1})
		return err
	})
	step("rest 1 hour", 3600, func() error {
		_, err := tools.NewRestHandler(testPool).Handle(toolCtx, map[string]any{"hours": 1, "short_rest": true})
		return err
	})
	step("move with prose travel time", 2*3600+15*60, func() error {
		res, err := game.NewStateStore(q, testPool).MovePlayer(ctx, game.MovePlayerCommand{CampaignID: campaign, PlayerCharacterID: dbutil.FromPgtype(pc.ID), CurrentLocationID: gate, TargetLocationID: wood})
		if err == nil && res.TimeWarning != "" {
			t.Fatalf("move time warning: %s", res.TimeWarning)
		}
		return err
	})
	var second int
	must(testPool.QueryRow(ctx, `SELECT second FROM campaign_time WHERE campaign_id=$1`, campaign).Scan(&second))
	if second != 30 {
		t.Fatalf("second = %d, want 30 kept by minute-only writers", second)
	}

	// Setting an absolute day/hour/minute means that minute's first second.
	_, err = saves.NewStore(testPool).UpsertCampaignTime(ctx, campaign, 3, 6, 0)
	must(err)
	if got, want := now(), 2*86400+6*3600; got != want {
		t.Fatalf("absolute set gives %d, want %d", got, want)
	}

	// Start over removes the row; the next legacy write starts at 08:00:00.
	must(saves.NewStore(testPool).StartOver(ctx, campaign))
	_, err = tools.NewAdvanceTimeHandler(testPool).Handle(toolCtx, map[string]any{"hours": 1})
	must(err)
	if got, want := now(), 9*3600; got != want {
		t.Fatalf("advance after start over gives %d, want %d", got, want)
	}
}
