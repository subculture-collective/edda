//go:build integration

package handlers_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"git.subcult.tv/subculture-collective/edda/internal/game"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

// MovePlayer in a journalled tool transaction locks the campaign row before
// the character row, the order the SRD rules actions use, so a transaction
// that holds the campaign and then locks the character does not deadlock
// with it.
func TestIntegrationMovePlayerLocksCampaignBeforeCharacter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	f := newClockRaceFixture(t)
	before := f.clock(t)
	holder, err := testPool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = holder.Rollback(context.Background()) }()
	if _, err := holder.Exec(ctx, `SELECT 1 FROM campaigns WHERE id=$1 FOR UPDATE`, f.campaign); err != nil {
		t.Fatal(err)
	}

	done := make(chan error, 1)
	go func() {
		done <- pgx.BeginFunc(ctx, testPool, func(tx pgx.Tx) error {
			res, err := game.NewStateStore(statedb.New(tx), tx).MovePlayer(ctx, game.MovePlayerCommand{CampaignID: f.campaign, PlayerCharacterID: f.bystander, CurrentLocationID: f.gate, TargetLocationID: f.wood})
			if err == nil && res.TimeWarning != "" {
				return &timeWarning{res.TimeWarning}
			}
			return err
		})
	}()
	for blocked := false; !blocked; {
		select {
		case err := <-done:
			t.Fatalf("move finished while the campaign was locked: %v", err)
		case <-ctx.Done():
			t.Fatal("move never waited for a lock")
		case <-time.After(10 * time.Millisecond):
			if err := testPool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM pg_stat_activity WHERE datname=current_database() AND wait_event_type='Lock')`).Scan(&blocked); err != nil {
				t.Fatal(err)
			}
		}
	}

	if _, err := holder.Exec(ctx, `SELECT 1 FROM player_characters WHERE id=$1 FOR UPDATE`, f.bystander); err != nil {
		t.Fatalf("lock the character under the campaign lock: %v", err)
	}
	if err := holder.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	var location string
	if err := testPool.QueryRow(ctx, `SELECT current_location_id::text FROM player_characters WHERE id=$1`, f.bystander).Scan(&location); err != nil {
		t.Fatal(err)
	}
	if location != f.wood.String() {
		t.Fatalf("character at %s, want %s", location, f.wood)
	}
	if got := f.clock(t); got != before+15*60 {
		t.Fatalf("clock = %d, want %d", got, before+15*60)
	}
}
