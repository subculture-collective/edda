//go:build integration

package statedb_test

import (
	"context"
	"testing"

	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

func TestIntegrationUpsertCampaignTimeResetsSecond(t *testing.T) {
	ctx := context.Background()
	tx, err := testPool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	q := statedb.New(tx)
	user, err := q.CreateUser(ctx, "clock-second-user")
	if err != nil {
		t.Fatal(err)
	}
	c, err := q.CreateCampaign(ctx, statedb.CreateCampaignParams{Name: "Absolute clock", Status: "active", CreatedBy: user.ID})
	if err != nil {
		t.Fatal(err)
	}
	set := func(day, hour int32) statedb.CampaignTime {
		t.Helper()
		got, err := q.UpsertCampaignTime(ctx, statedb.UpsertCampaignTimeParams{CampaignID: c.ID, Day: day, Hour: hour})
		if err != nil {
			t.Fatal(err)
		}
		return got
	}
	if got := set(1, 8); got.Second != 0 {
		t.Fatalf("new row second = %d, want 0", got.Second)
	}
	if _, err := tx.Exec(ctx, `UPDATE campaign_time SET second=30 WHERE campaign_id=$1`, c.ID); err != nil {
		t.Fatal(err)
	}
	if got := set(3, 6); got.Day != 3 || got.Hour != 6 || got.Minute != 0 || got.Second != 0 {
		t.Fatalf("absolute set = day %d %02d:%02d:%02d, want day 3 06:00:00", got.Day, got.Hour, got.Minute, got.Second)
	}
}
