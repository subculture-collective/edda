//go:build integration

package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"github.com/google/uuid"
)

// concentrationFixture is a concentrating Cleric (Shield of Faith on an ally)
// beside a second concentrating Cleric whose cast must never be touched.
type concentrationFixture struct {
	t                      *testing.T
	campaign, caster, ally uuid.UUID
	cast, otherCast        uuid.UUID
}

func newConcentrationFixture(t *testing.T) *concentrationFixture {
	t.Helper()
	ctx := context.Background()
	q := statedb.New(testPool)
	f := &concentrationFixture{t: t, cast: uuid.New(), otherCast: uuid.New()}
	creation, err := srd521.SupportedBuildForConcept("Human Cleric", "Acolyte")
	f.must(err)
	stats, err := creation.Stats()
	f.must(err)
	c, err := q.CreateCampaign(ctx, statedb.CreateCampaignParams{Name: "Concentration events fixture", RulesVersion: "srd-5.2.1", Status: "active", CreatedBy: dbutil.ToPgtype(testUserID)})
	f.must(err)
	t.Cleanup(func() { _, _ = testPool.Exec(ctx, `DELETE FROM campaigns WHERE id=$1`, c.ID) })
	f.campaign = dbutil.FromPgtype(c.ID)
	create := func(name string) uuid.UUID {
		pc, err := q.CreatePlayerCharacter(ctx, statedb.CreatePlayerCharacterParams{CampaignID: c.ID, UserID: dbutil.ToPgtype(testUserID), Name: name, Stats: stats, Hp: 10, MaxHp: 10, Level: 1, Status: "active", Abilities: []byte(`[]`)})
		f.must(err)
		return dbutil.FromPgtype(pc.ID)
	}
	f.caster, f.ally = create("Mira"), create("Ansel")
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET stats=stats||'{"custom":"keep"}'::jsonb WHERE id=$1`, f.caster)
	f.must(err)
	f.concentrate(f.caster, f.cast, f.ally)
	f.concentrate(f.ally, f.otherCast, f.ally)
	return f
}

func (f *concentrationFixture) must(err error) {
	f.t.Helper()
	if err != nil {
		f.t.Fatal(err)
	}
}

func (f *concentrationFixture) concentrate(caster, cast, recipient uuid.UUID) {
	f.t.Helper()
	ctx := context.Background()
	c, _ := json.Marshal(srd521.Concentration{ID: cast.String(), Spell: "shield-of-faith", RecipientID: recipient.String(), StartedAt: 28800, ExpiresAt: 29400})
	_, err := testPool.Exec(ctx, `UPDATE player_characters SET hp=10,stats=jsonb_set(stats,'{srd521,casting,concentration}',$2::jsonb) WHERE id=$1`, caster, c)
	f.must(err)
	_, err = testPool.Exec(ctx, `INSERT INTO concentration_casts(id,campaign_id,caster_id,started_at,expires_at) VALUES($1,$2,$3,28800,29400)`, cast, f.campaign, caster)
	f.must(err)
	_, err = testPool.Exec(ctx, `INSERT INTO concentration_recipients(cast_id,recipient_id) VALUES($1,$2)`, cast, recipient)
	f.must(err)
}

// state reads the caster's HP, saved concentration ID, links for the
// caster's cast and the other caster's cast, and the unrelated stats key.
func (f *concentrationFixture) state() (hp int, concentration string, links, otherLinks int, custom string) {
	f.t.Helper()
	ctx := context.Background()
	f.must(testPool.QueryRow(ctx, `SELECT hp,COALESCE(stats#>>'{srd521,casting,concentration,id}',''),COALESCE(stats->>'custom','') FROM player_characters WHERE id=$1`, f.caster).Scan(&hp, &concentration, &custom))
	f.must(testPool.QueryRow(ctx, `SELECT count(*) FILTER (WHERE cast_id=$1),count(*) FILTER (WHERE cast_id=$2) FROM concentration_recipients`, f.cast, f.otherCast).Scan(&links, &otherLinks))
	return
}

// apply runs the adapter in its own transaction, committing only on success
// and when commit is true.
func (f *concentrationFixture) apply(campaign uuid.UUID, hp int, events []srd521.ConcentrationEvent, dice srd521.Dice, commit bool) ([]*srd521.EffectReceipt, error) {
	f.t.Helper()
	ctx := context.Background()
	tx, err := testPool.Begin(ctx)
	f.must(err)
	defer func() { _ = tx.Rollback(ctx) }()
	receipts, err := engine.ApplyConcentrationEvents(ctx, tx, campaign, f.caster, hp, events, dice)
	if err == nil && commit {
		f.must(tx.Commit(ctx))
	}
	return receipts, err
}

func TestIntegrationConcentrationEventsAdapter(t *testing.T) {
	f := newConcentrationFixture(t)
	dice := &lifecycleIntegrationDice{}
	unchanged := func(when string) {
		t.Helper()
		if hp, c, links, other, custom := f.state(); hp != 10 || c != f.cast.String() || links != 1 || other != 1 || custom != "keep" || dice.calls != 0 {
			t.Fatalf("%s changed state: hp=%d concentration=%q links=%d other=%d custom=%q calls=%d", when, hp, c, links, other, custom, dice.calls)
		}
	}
	twoHitsThenUnconscious := []srd521.ConcentrationEvent{
		{Kind: srd521.ConcentrationDamageSource, SourceID: "attack-1", Damage: 3, HPAfter: 7},
		{Kind: srd521.ConcentrationDamageSource, SourceID: "attack-2", Damage: 2, HPAfter: 5},
		{Kind: srd521.ConcentrationUnconscious},
	}

	// Foreign campaigns and unknown event keys fail closed with no roll.
	if _, err := f.apply(uuid.New(), 10, twoHitsThenUnconscious, dice, true); !errors.Is(err, engine.ErrConcentrationCaster) {
		t.Fatalf("foreign campaign: %v", err)
	}
	unchanged("foreign campaign")
	if _, err := f.apply(f.campaign, 10, []srd521.ConcentrationEvent{{Kind: "concentration.end.asleep"}}, dice, true); err == nil || !strings.Contains(err.Error(), "unsupported concentration event") {
		t.Fatalf("unknown key: %v", err)
	}
	unchanged("unknown key")

	// A transaction that rolls back after the adapter leaves nothing behind.
	dice.values = []int{7, 7}
	if _, err := f.apply(f.campaign, 10, twoHitsThenUnconscious, dice, false); err != nil {
		t.Fatal(err)
	}
	dice.calls = 0
	unchanged("rolled-back transaction")

	// The exact fixture: two saves (DC 10 each, both succeed), then the
	// trusted Unconscious transition at positive HP ends the cast unrolled.
	dice.values = []int{7, 7}
	receipts, err := f.apply(f.campaign, 10, twoHitsThenUnconscious, dice, true)
	f.must(err)
	if dice.calls != 2 || len(receipts) != 3 || receipts[0].Check == nil || receipts[1].Check == nil || receipts[2].Check != nil || receipts[2].Mechanic != srd521.ConcentrationUnconscious || receipts[2].Before.ID != f.cast.String() {
		t.Fatalf("calls=%d receipts=%+v", dice.calls, receipts)
	}
	if hp, c, links, other, custom := f.state(); hp != 5 || c != "" || links != 0 || other != 1 || custom != "keep" {
		t.Fatalf("after: hp=%d concentration=%q links=%d other=%d custom=%q", hp, c, links, other, custom)
	}
	var castKept bool
	f.must(testPool.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM concentration_casts WHERE id=$1)`, f.cast).Scan(&castKept))
	if !castKept {
		t.Fatal("cast identity deleted")
	}

	// Healing and waking never restore the ended cast, and further events on
	// the caster roll nothing.
	_, err = testPool.Exec(context.Background(), `UPDATE player_characters SET hp=10 WHERE id=$1`, f.caster)
	f.must(err)
	dice.calls = 0
	receipts, err = f.apply(f.campaign, 10, []srd521.ConcentrationEvent{{Kind: srd521.ConcentrationDamageSource, SourceID: "attack-3", Damage: 4, HPAfter: 6}, {Kind: srd521.ConcentrationDead}}, dice, true)
	if err != nil || len(receipts) != 0 || dice.calls != 0 {
		t.Fatalf("ended cast resolved again: %+v calls=%d err=%v", receipts, dice.calls, err)
	}
	if _, c, links, _, _ := f.state(); c != "" || links != 0 {
		t.Fatal("ended cast came back")
	}
}

// Fall damage reaches the Concentration save through the shared adapter
// (tracker #2164: #147 reused the NPC-attack helper with no fall test).
func TestIntegrationConcentrationAfterFall(t *testing.T) {
	f := newConcentrationFixture(t)
	ctx := context.Background()
	faces := &fallFaces{}
	declare := func(id string) {
		t.Helper()
		_, err := testPool.Exec(ctx, `UPDATE player_characters SET stats=jsonb_set(stats,'{srd521_fall}',$2::jsonb) WHERE id=$1`, f.caster, fallDecl(id, 20, "solid"))
		f.must(err)
	}
	preview := func(id string) *game.ActionPreview {
		t.Helper()
		c, pc, err := game.ReadActionSnapshot(ctx, testPool, f.campaign, f.caster, testUserID, false)
		f.must(err)
		a, err := game.PrepareStoredAction(ctx, testPool, c, pc, game.ActionIntent{ActorID: f.caster, Kind: srd521.FallResolve, EntryID: id}, false)
		f.must(err)
		if !a.Preview.Available {
			t.Fatalf("fall unavailable: %+v", a.Preview)
		}
		return a.Preview
	}
	confirm := func(r engine.PreparedAction, values ...int) (*engine.TurnResult, error) {
		faces.faces, faces.calls = values, 0
		e, err := engine.New(testPool, nil, config.LLMConfig{}, engine.WithActionDice(faces))
		f.must(err)
		return e.ConfirmAction(ctx, f.campaign, r)
	}
	request := func(p *game.ActionPreview) engine.PreparedAction {
		return engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: p.ExpectedRevision, Intent: p.Intent}
	}

	// 20 feet: 2d6 (3, 4) = 7 damage, HP 10 to 3, then a DC 10 save.
	declare("f1")
	p := preview("f1")
	if faces.calls != 0 {
		t.Fatal("preview rolled")
	}
	// A history-write failure rolls back the fall, the save and the cast.
	_, err := testPool.Exec(ctx, `ALTER TABLE session_logs ADD CONSTRAINT concentration_fall_failure CHECK (player_input NOT LIKE 'Confirm Fall%') NOT VALID`)
	f.must(err)
	t.Cleanup(func() {
		_, _ = testPool.Exec(ctx, `ALTER TABLE session_logs DROP CONSTRAINT IF EXISTS concentration_fall_failure`)
	})
	if _, err := confirm(request(p), 3, 4, 1); err == nil {
		t.Fatal("history failure committed")
	}
	_, err = testPool.Exec(ctx, `ALTER TABLE session_logs DROP CONSTRAINT concentration_fall_failure`)
	f.must(err)
	if hp, c, links, other, custom := f.state(); hp != 10 || c != f.cast.String() || links != 1 || other != 1 || custom != "keep" {
		t.Fatalf("rollback left: hp=%d concentration=%q links=%d other=%d", hp, c, links, other)
	}

	// Success: d20 8 + 2 = 10 meets DC 10; the cast and its link stay.
	r := request(p)
	result, err := confirm(r, 3, 4, 8)
	f.must(err)
	receipt := fallReceipt(t, result)
	if faces.calls != 3 || receipt.Damage != 7 || receipt.HPAfter != 3 || receipt.Concentration == nil || receipt.Concentration.Check == nil || *receipt.Concentration.Check.DC != 10 || receipt.Concentration.Check.Outcome != "success" || receipt.Concentration.After == nil {
		t.Fatalf("calls=%d receipt=%+v concentration=%+v", faces.calls, receipt, receipt.Concentration)
	}
	if !strings.Contains(result.Narrative, "Concentration maintained after damage.") {
		t.Fatalf("narrative: %s", result.Narrative)
	}
	if hp, c, links, _, _ := f.state(); hp != 3 || c != f.cast.String() || links != 1 {
		t.Fatalf("after success: hp=%d concentration=%q links=%d", hp, c, links)
	}
	// Replaying the same action ID returns the receipt with no extra roll.
	replay, err := confirm(r, 6, 6, 1)
	f.must(err)
	if faces.calls != 0 || replay.Narrative != result.Narrative {
		t.Fatalf("replay rolled or changed: calls=%d %s", faces.calls, replay.Narrative)
	}

	// Failure: a second fall from 3 HP, 2d6 (1, 1) = 2, d20 1 + 2 fails; only
	// this caster's link ends.
	declare("f2")
	result, err = confirm(request(preview("f2")), 1, 1, 1)
	f.must(err)
	receipt = fallReceipt(t, result)
	if receipt.Concentration == nil || receipt.Concentration.Check.Outcome != "failure" || receipt.Concentration.After != nil {
		t.Fatalf("failed save receipt: %+v", receipt.Concentration)
	}
	if hp, c, links, other, custom := f.state(); hp != 1 || c != "" || links != 0 || other != 1 || custom != "keep" {
		t.Fatalf("after failure: hp=%d concentration=%q links=%d other=%d custom=%q", hp, c, links, other, custom)
	}

	// No concentration left: a further fall rolls only its damage dice.
	declare("f3")
	result, err = confirm(request(preview("f3")), 1, 1)
	f.must(err)
	if receipt := fallReceipt(t, result); receipt.Concentration != nil || faces.calls != 2 {
		t.Fatalf("save without concentration: %+v calls=%d", receipt.Concentration, faces.calls)
	}
}
