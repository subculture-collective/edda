//go:build integration

package handlers_test

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/db"
	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/internal/tools"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// pausedDice holds the first roll until released. A Short Rest rolls its Hit
// Die after the rules code has read the campaign clock and before it writes it.
type pausedDice struct {
	once            sync.Once
	rolled, release chan struct{}
}

func (d *pausedDice) Intn(int) int {
	d.once.Do(func() {
		close(d.rolled)
		<-d.release
	})
	return 0
}

type clockRaceFixture struct {
	campaign, actor, bystander, gate, wood uuid.UUID
}

// newClockRaceFixture creates a Cleric who can take a Short Rest, a second
// character whose row the rules transaction does not lock, and a 15 minute road.
func newClockRaceFixture(t *testing.T) clockRaceFixture {
	t.Helper()
	ctx := context.Background()
	q := statedb.New(testPool)
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	creation, err := srd521.SupportedBuildForConcept("Human Cleric", "Acolyte")
	must(err)
	stats, err := creation.Stats()
	must(err)
	c, err := q.CreateCampaign(ctx, statedb.CreateCampaignParams{Name: "Clock race", RulesVersion: "srd-5.2.1", Status: "active", CreatedBy: dbutil.ToPgtype(testUserID)})
	must(err)
	t.Cleanup(func() { _, _ = testPool.Exec(ctx, `DELETE FROM campaigns WHERE id=$1`, c.ID) })
	f := clockRaceFixture{campaign: dbutil.FromPgtype(c.ID)}
	character := func(name string, stats []byte) uuid.UUID {
		t.Helper()
		pc, err := q.CreatePlayerCharacter(ctx, statedb.CreatePlayerCharacterParams{CampaignID: c.ID, UserID: dbutil.ToPgtype(testUserID), Name: name, Stats: stats, Hp: 2, MaxHp: 10, Level: 1, Status: "active", Abilities: []byte(`[]`)})
		must(err)
		return dbutil.FromPgtype(pc.ID)
	}
	f.actor = character("Mira", stats)
	f.bystander = character("Pell", []byte(`{}`))
	for name, id := range map[string]*uuid.UUID{"Gate": &f.gate, "Wood": &f.wood} {
		must(testPool.QueryRow(ctx, `INSERT INTO locations(campaign_id,name,player_known)VALUES($1,$2,true)RETURNING id`, f.campaign, name).Scan(id))
	}
	_, err = testPool.Exec(ctx, `INSERT INTO location_connections(from_location_id,to_location_id,bidirectional,travel_time,campaign_id)VALUES($1,$2,true,'15 minutes',$3)`, f.gate, f.wood, f.campaign)
	must(err)
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET current_location_id=$1 WHERE campaign_id=$2`, f.gate, f.campaign)
	must(err)
	return f
}

func (f clockRaceFixture) clock(t *testing.T) int {
	t.Helper()
	now, err := game.ReadSpellClock(context.Background(), testPool, f.campaign)
	if err != nil {
		t.Fatal(err)
	}
	return now
}

func (f clockRaceFixture) advance(ctx context.Context, database db.DBTX, minutes int) error {
	_, err := tools.NewAdvanceTimeHandler(database).Handle(tools.WithCurrentCampaignID(ctx, f.campaign), map[string]any{"hours": 0, "minutes": minutes})
	return err
}

func (f clockRaceFixture) rest(ctx context.Context, database db.DBTX) error {
	_, err := tools.NewRestHandler(database).Handle(tools.WithCurrentCampaignID(ctx, f.campaign), map[string]any{"hours": 1, "short_rest": true})
	return err
}

func (f clockRaceFixture) move(ctx context.Context, database db.DBTX) error {
	res, err := game.NewStateStore(statedb.New(testPool), database).MovePlayer(ctx, game.MovePlayerCommand{CampaignID: f.campaign, PlayerCharacterID: f.bystander, CurrentLocationID: f.gate, TargetLocationID: f.wood})
	if err == nil && res.TimeWarning != "" {
		return &timeWarning{res.TimeWarning}
	}
	return err
}

type timeWarning struct{ text string }

func (w *timeWarning) Error() string { return "move time warning: " + w.text }

// TestIntegrationLegacyClockWriterSerializesWithRulesClock interleaves a
// legacy writer between the rules clock's read and write. Both advances must
// survive, whether or not the campaign already has a clock row.
func TestIntegrationLegacyClockWriterSerializesWithRulesClock(t *testing.T) {
	writers := []struct {
		name  string
		delta int
		run   func(clockRaceFixture, context.Context) error
	}{
		{"advance_time", 7 * 60, func(f clockRaceFixture, ctx context.Context) error { return f.advance(ctx, testPool, 7) }},
		{"move_player", 15 * 60, func(f clockRaceFixture, ctx context.Context) error { return f.move(ctx, testPool) }},
	}
	for _, writer := range writers {
		for _, seeded := range []bool{true, false} {
			name := writer.name + " without a clock row"
			if seeded {
				name = writer.name + " with a clock row"
			}
			t.Run(name, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				f := newClockRaceFixture(t)
				if seeded {
					if _, err := testPool.Exec(ctx, `INSERT INTO campaign_time(campaign_id,day,hour,minute,second)VALUES($1,1,23,30,30)`, f.campaign); err != nil {
						t.Fatal(err)
					}
				}
				start := f.clock(t)

				dice := &pausedDice{rolled: make(chan struct{}), release: make(chan struct{})}
				release := sync.OnceFunc(func() { close(dice.release) })
				defer release()
				e, err := engine.New(testPool, nil, config.LLMConfig{}, engine.WithActionDice(dice))
				if err != nil {
					t.Fatal(err)
				}
				campaign, pc, err := game.ReadActionSnapshot(ctx, testPool, f.campaign, f.actor, testUserID, false)
				if err != nil {
					t.Fatal(err)
				}
				prepared, err := game.PrepareStoredAction(ctx, testPool, campaign, pc, game.ActionIntent{ActorID: f.actor, Kind: srd521.ShortRest, EntryID: "hit-die"}, false)
				if err != nil {
					t.Fatal(err)
				}
				if !prepared.Preview.Available {
					t.Fatalf("short rest unavailable: %+v", prepared.Preview)
				}
				rules := make(chan error, 1)
				go func() {
					_, err := e.ConfirmAction(ctx, f.campaign, engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: prepared.Preview.ExpectedRevision, Intent: prepared.Preview.Intent})
					rules <- err
				}()
				select {
				case <-dice.rolled:
				case err := <-rules:
					t.Fatalf("rest finished without rolling: %v", err)
				case <-ctx.Done():
					t.Fatal("rest never rolled")
				}

				// The rules transaction has read the clock. Let the legacy writer
				// either finish or queue behind a database lock, then release it.
				legacy := make(chan error, 1)
				go func() { legacy <- writer.run(f, ctx) }()
				finished := false
				for blocked := false; !finished && !blocked; {
					select {
					case err := <-legacy:
						if err != nil {
							t.Fatal(err)
						}
						finished = true
					case <-ctx.Done():
						t.Fatal("legacy writer neither finished nor blocked")
					case <-time.After(10 * time.Millisecond):
						if err := testPool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM pg_stat_activity WHERE datname=current_database() AND wait_event_type='Lock')`).Scan(&blocked); err != nil {
							t.Fatal(err)
						}
					}
				}
				release()
				if err := <-rules; err != nil {
					t.Fatal(err)
				}
				if !finished {
					if err := <-legacy; err != nil {
						t.Fatal(err)
					}
				}
				if got, want := f.clock(t), start+3600+writer.delta; got != want {
					t.Fatalf("clock = %d, want %d: lost %d seconds", got, want, want-got)
				}
			})
		}
	}
}

// afterClockQuery runs then once, right after the first campaign_time
// statement sent through it returns.
type afterClockQuery struct {
	db.DBTX
	once sync.Once
	then func()
}

func (d *afterClockQuery) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	row := d.DBTX.QueryRow(ctx, query, args...)
	if !strings.Contains(query, "campaign_time") {
		return row
	}
	return scanThen{row, func() { d.once.Do(d.then) }}
}

type scanThen struct {
	pgx.Row
	then func()
}

func (r scanThen) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	r.then()
	return err
}

// TestIntegrationRulesClockWriteSurvivesLegacyWriter commits a whole rules
// advance right after a legacy writer's first clock statement. A writer that
// reads first and writes later would overwrite the rules minute.
func TestIntegrationRulesClockWriteSurvivesLegacyWriter(t *testing.T) {
	writers := []struct {
		name  string
		delta int
		run   func(clockRaceFixture, context.Context, db.DBTX) error
	}{
		{"advance_time", 7 * 60, func(f clockRaceFixture, ctx context.Context, d db.DBTX) error { return f.advance(ctx, d, 7) }},
		{"rest", 3600, clockRaceFixture.rest},
		{"move_player", 15 * 60, clockRaceFixture.move},
	}
	for _, writer := range writers {
		t.Run(writer.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			f := newClockRaceFixture(t)
			if _, err := testPool.Exec(ctx, `INSERT INTO campaign_time(campaign_id,day,hour,minute,second)VALUES($1,1,23,30,30)`, f.campaign); err != nil {
				t.Fatal(err)
			}
			start := f.clock(t)
			e, err := engine.New(testPool, nil, config.LLMConfig{})
			if err != nil {
				t.Fatal(err)
			}
			waited := false
			database := &afterClockQuery{DBTX: testPool, then: func() {
				campaign, pc, err := game.ReadActionSnapshot(ctx, testPool, f.campaign, f.actor, testUserID, false)
				if err != nil {
					t.Fatal(err)
				}
				prepared, err := game.PrepareStoredAction(ctx, testPool, campaign, pc, game.ActionIntent{ActorID: f.actor, Kind: srd521.WaitMinute}, false)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := e.ConfirmAction(ctx, f.campaign, engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: prepared.Preview.ExpectedRevision, Intent: prepared.Preview.Intent}); err != nil {
					t.Fatal(err)
				}
				waited = true
			}}
			if err := writer.run(f, ctx, database); err != nil {
				t.Fatal(err)
			}
			if !waited {
				t.Fatal("rules wait never ran")
			}
			if got, want := f.clock(t), start+60+writer.delta; got != want {
				t.Fatalf("clock = %d, want %d: lost %d seconds", got, want, want-got)
			}
		})
	}
}

// TestIntegrationConcurrentLegacyClockWritersKeepEveryAdvance starts all three
// legacy writers together; each advance must be added exactly once.
func TestIntegrationConcurrentLegacyClockWritersKeepEveryAdvance(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	f := newClockRaceFixture(t)
	// rest updates HP on every character in the campaign, and concurrent
	// multi-row updates can deadlock on their own. Keep one character here.
	if _, err := testPool.Exec(ctx, `DELETE FROM player_characters WHERE id=$1`, f.bystander); err != nil {
		t.Fatal(err)
	}
	f.bystander = f.actor
	start := f.clock(t)

	const rounds = 8
	writers := []func() error{
		func() error { return f.advance(ctx, testPool, 7) },
		func() error { return f.rest(ctx, testPool) },
		// Every move names its own current location, so repeats stay valid.
		func() error { return f.move(ctx, testPool) },
	}
	var wg sync.WaitGroup
	barrier := make(chan struct{})
	errs := make(chan error, rounds*len(writers))
	for range rounds {
		for _, write := range writers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-barrier
				errs <- write()
			}()
		}
	}
	close(barrier)
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	if got, want := f.clock(t), start+rounds*(7*60+3600+15*60); got != want {
		t.Fatalf("clock = %d, want %d: lost %d seconds", got, want, want-got)
	}
}
