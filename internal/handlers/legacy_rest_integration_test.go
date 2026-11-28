//go:build integration

package handlers_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/internal/tools"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type legacyRestFixture struct {
	campaign, mira, pell uuid.UUID
}

// newLegacyRestFixture creates a legacy campaign with two player characters
// whose HP and maximum HP differ.
func newLegacyRestFixture(t *testing.T) legacyRestFixture {
	t.Helper()
	ctx := context.Background()
	q := statedb.New(testPool)
	c, err := q.CreateCampaign(ctx, statedb.CreateCampaignParams{Name: "Legacy rest", RulesVersion: "legacy", Status: "active", CreatedBy: dbutil.ToPgtype(testUserID)})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = testPool.Exec(ctx, `DELETE FROM campaigns WHERE id=$1`, c.ID) })
	character := func(name string, hp, maxHP int32) uuid.UUID {
		t.Helper()
		pc, err := q.CreatePlayerCharacter(ctx, statedb.CreatePlayerCharacterParams{CampaignID: c.ID, UserID: dbutil.ToPgtype(testUserID), Name: name, Stats: []byte(`{}`), Hp: hp, MaxHp: maxHP, Level: 1, Status: "active", Abilities: []byte(`[]`)})
		if err != nil {
			t.Fatal(err)
		}
		return dbutil.FromPgtype(pc.ID)
	}
	return legacyRestFixture{campaign: dbutil.FromPgtype(c.ID), mira: character("Mira", 2, 20), pell: character("Pell", 7, 40)}
}

func (f legacyRestFixture) rest(ctx context.Context, database tools.TimeStore, actor uuid.UUID, short bool) (*tools.ToolResult, error) {
	ctx = tools.WithCurrentPlayerCharacterID(tools.WithCurrentCampaignID(ctx, f.campaign), actor)
	return tools.NewRestHandler(database).Handle(ctx, map[string]any{"hours": 1, "short_rest": short})
}

// restInTx runs a rest in its own transaction, as the engine's journalled tool
// transaction does, so row locks last until the commit.
func (f legacyRestFixture) restInTx(ctx context.Context, actor uuid.UUID, short bool) (*tools.ToolResult, error) {
	var result *tools.ToolResult
	err := pgx.BeginFunc(ctx, testPool, func(tx pgx.Tx) (err error) {
		result, err = f.rest(ctx, tx, actor, short)
		return err
	})
	return result, err
}

func (f legacyRestFixture) assertHP(t *testing.T, step string, mira, pell int) {
	t.Helper()
	var gotMira, gotPell int
	if err := testPool.QueryRow(context.Background(), `SELECT (SELECT hp FROM player_characters WHERE id=$1),(SELECT hp FROM player_characters WHERE id=$2)`, f.mira, f.pell).Scan(&gotMira, &gotPell); err != nil {
		t.Fatal(err)
	}
	if gotMira != mira || gotPell != pell {
		t.Fatalf("%s: HP Mira=%d Pell=%d, want Mira=%d Pell=%d", step, gotMira, gotPell, mira, pell)
	}
}

func TestIntegrationLegacyRestHealsOnlyTheCurrentCharacterFromItsOwnHP(t *testing.T) {
	ctx := context.Background()
	f := newLegacyRestFixture(t)
	steps := []struct {
		name        string
		actor       uuid.UUID
		short       bool
		hp, maxHP   int
		mira, other int
	}{
		{"short rest Pell", f.pell, true, 17, 40, 2, 17},
		{"short rest Mira", f.mira, true, 7, 20, 7, 17},
		{"long rest Mira", f.mira, false, 20, 20, 20, 17},
		{"short rest Mira at full HP", f.mira, true, 20, 20, 20, 17},
		{"long rest Pell", f.pell, false, 40, 40, 20, 40},
	}
	for _, step := range steps {
		result, err := f.rest(ctx, testPool, step.actor, step.short)
		if err != nil {
			t.Fatalf("%s: %v", step.name, err)
		}
		f.assertHP(t, step.name, step.mira, step.other)
		if result.Data["hp"] != step.hp || result.Data["max_hp"] != step.maxHP {
			t.Fatalf("%s: result HP %v/%v, want %d/%d", step.name, result.Data["hp"], result.Data["max_hp"], step.hp, step.maxHP)
		}
	}
}

// A rest that names no character in this campaign must not heal anyone or move
// the clock.
func TestIntegrationLegacyRestRequiresACharacterInTheCampaign(t *testing.T) {
	ctx := context.Background()
	f := newLegacyRestFixture(t)
	outsider := newLegacyRestFixture(t)

	if _, err := tools.NewRestHandler(testPool).Handle(tools.WithCurrentCampaignID(ctx, f.campaign), map[string]any{"hours": 8}); err == nil {
		t.Fatal("rest without a current player character succeeded")
	}
	if _, err := f.rest(ctx, testPool, outsider.mira, false); err == nil {
		t.Fatal("rest for another campaign's character succeeded")
	}
	f.assertHP(t, "refused rests", 2, 7)
	outsider.assertHP(t, "outsider", 2, 7)
	var clocks int
	if err := testPool.QueryRow(ctx, `SELECT count(*) FROM campaign_time WHERE campaign_id IN ($1,$2)`, f.campaign, outsider.campaign).Scan(&clocks); err != nil {
		t.Fatal(err)
	}
	if clocks != 0 {
		t.Fatalf("refused rests wrote %d clock rows", clocks)
	}
}

// TestIntegrationConcurrentLegacyRestsDoNotDeadlock releases several rests for
// two characters together. Every rest must succeed and heal its own character
// exactly once. Three short rests stay below either maximum, so a lost heal
// shows up. Single-row writes run alongside, like a move would: they change the
// order in which a statement over several characters finds their rows, which is
// what lets two such statements deadlock.
func TestIntegrationConcurrentLegacyRestsDoNotDeadlock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	f := newLegacyRestFixture(t)
	const rounds, perCharacter, otherWrites = 100, 3, 24
	for round := range rounds {
		if _, err := testPool.Exec(ctx, `UPDATE player_characters SET hp=CASE WHEN id=$2 THEN 1 ELSE 2 END, max_hp=CASE WHEN id=$2 THEN 4000 ELSE 2000 END WHERE campaign_id=$1`, f.campaign, f.mira); err != nil {
			t.Fatal(err)
		}
		var wg sync.WaitGroup
		barrier := make(chan struct{})
		errs := make(chan error, 2*perCharacter+otherWrites)
		for range perCharacter {
			for _, actor := range []uuid.UUID{f.mira, f.pell} {
				wg.Add(1)
				go func() {
					defer wg.Done()
					<-barrier
					_, err := f.restInTx(ctx, actor, true)
					errs <- err
				}()
			}
		}
		for i := range otherWrites {
			id := f.mira
			if i%3 == 0 {
				id = f.pell
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-barrier
				_, err := testPool.Exec(ctx, `UPDATE player_characters SET status=status WHERE id=$1`, id)
				errs <- err
			}()
		}
		close(barrier)
		wg.Wait()
		close(errs)
		for err := range errs {
			if err != nil {
				t.Fatalf("round %d: %v", round, err)
			}
		}
		f.assertHP(t, "concurrent short rests", 1+perCharacter*1000, 2+perCharacter*500)
	}
}

// TestIntegrationLegacyRestLocksCampaignBeforeCharacter holds the campaign row
// the way a rules action does, then locks and damages the resting character. A
// rest that locked characters first would deadlock here. The queued rest must
// also heal from the HP committed while it waited.
func TestIntegrationLegacyRestLocksCampaignBeforeCharacter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	f := newLegacyRestFixture(t)
	holder, err := testPool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = holder.Rollback(context.Background()) }()
	if _, err := holder.Exec(ctx, `SELECT 1 FROM campaigns WHERE id=$1 FOR UPDATE`, f.campaign); err != nil {
		t.Fatal(err)
	}

	type outcome struct {
		result *tools.ToolResult
		err    error
	}
	done := make(chan outcome, 1)
	go func() {
		result, err := f.restInTx(ctx, f.mira, true)
		done <- outcome{result, err}
	}()
	for blocked := false; !blocked; {
		select {
		case got := <-done:
			t.Fatalf("rest finished while the campaign was locked: %v", got.err)
		case <-ctx.Done():
			t.Fatal("rest never waited for a lock")
		case <-time.After(10 * time.Millisecond):
			if err := testPool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM pg_stat_activity WHERE datname=current_database() AND wait_event_type='Lock')`).Scan(&blocked); err != nil {
				t.Fatal(err)
			}
		}
	}

	if _, err := holder.Exec(ctx, `SELECT 1 FROM player_characters WHERE campaign_id=$1 ORDER BY id FOR UPDATE`, f.campaign); err != nil {
		t.Fatalf("lock characters under the campaign lock: %v", err)
	}
	if _, err := holder.Exec(ctx, `UPDATE player_characters SET hp=3 WHERE id=$1`, f.mira); err != nil {
		t.Fatal(err)
	}
	if err := holder.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	got := <-done
	if got.err != nil {
		t.Fatal(got.err)
	}
	if got.result.Data["hp"] != 8 {
		t.Fatalf("result HP %v, want 8 from the committed 3", got.result.Data["hp"])
	}
	f.assertHP(t, "queued short rest", 8, 7)
}
