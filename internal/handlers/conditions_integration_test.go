//go:build integration

package handlers_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/combat"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
)

// TestIntegrationConditionAdapter applies and ends Incapacitated and
// Unconscious on a concentrating Cleric through the trusted adapter (#487).
func TestIntegrationConditionAdapter(t *testing.T) {
	f := newConcentrationFixture(t)
	ctx := context.Background()
	dice := &lifecycleIntegrationDice{}
	item := func(name, kind, key string, equipped bool) uuid.UUID {
		t.Helper()
		var id uuid.UUID
		f.must(testPool.QueryRow(ctx, `INSERT INTO items(campaign_id,player_character_id,name,item_type,rarity,equipped,properties) VALUES($1,$2,$3,$4,'common',$5,jsonb_build_object('srd521_id',$6::text)) RETURNING id`, f.campaign, f.caster, name, kind, equipped, key).Scan(&id))
		return id
	}
	mace, shield, mail := item("Mace", "weapon", "mace", true), item("Shield", "armor", "shield", true), item("Chain Mail", "armor", "chain-mail", true)
	rope := item("Rope", "misc", "rope", true)
	_, err := testPool.Exec(ctx, `UPDATE player_characters SET stats=jsonb_set(stats,'{srd521,shield}','true'::jsonb) WHERE id=$1`, f.caster)
	f.must(err)

	apply := func(campaign uuid.UUID, encounter *combat.Encounter, change srd521.ConditionChange, commit bool) (*srd521.ConditionReceipt, error) {
		t.Helper()
		tx, err := testPool.Begin(ctx)
		f.must(err)
		defer func() { _ = tx.Rollback(ctx) }()
		r, err := engine.ApplyConditionChange(ctx, tx, campaign, f.caster, encounter, change, dice)
		if err == nil && commit {
			f.must(tx.Commit(ctx))
		}
		return r, err
	}
	type snapshot struct {
		conditions, concentration, custom string
		links, otherLinks                 int
		equipped                          map[uuid.UUID]bool
		shieldBuilt                       bool
	}
	read := func() snapshot {
		t.Helper()
		s := snapshot{equipped: map[uuid.UUID]bool{}}
		f.must(testPool.QueryRow(ctx, `SELECT COALESCE(stats->>'srd521_conditions',''),COALESCE(stats#>>'{srd521,casting,concentration,id}',''),COALESCE(stats->>'custom',''),COALESCE((stats#>>'{srd521,shield}')::boolean,false) FROM player_characters WHERE id=$1`, f.caster).Scan(&s.conditions, &s.concentration, &s.custom, &s.shieldBuilt))
		f.must(testPool.QueryRow(ctx, `SELECT count(*) FILTER (WHERE cast_id=$1),count(*) FILTER (WHERE cast_id=$2) FROM concentration_recipients`, f.cast, f.otherCast).Scan(&s.links, &s.otherLinks))
		for _, id := range []uuid.UUID{mace, shield, mail, rope} {
			var e bool
			f.must(testPool.QueryRow(ctx, `SELECT equipped FROM items WHERE id=$1`, id).Scan(&e))
			s.equipped[id] = e
		}
		return s
	}
	initial := read()
	unchanged := func(when string) {
		t.Helper()
		got := read()
		if got.conditions != initial.conditions || got.concentration != initial.concentration || got.custom != "keep" || got.links != 1 || got.otherLinks != 1 || !got.shieldBuilt || dice.calls != 0 {
			t.Fatalf("%s changed state: %+v calls=%d", when, got, dice.calls)
		}
		for id, e := range initial.equipped {
			if got.equipped[id] != e {
				t.Fatalf("%s changed item %s", when, id)
			}
		}
	}
	preview := func(intent game.ActionIntent) *game.ActionPreview {
		t.Helper()
		c, pc, err := game.ReadActionSnapshot(ctx, testPool, f.campaign, f.caster, testUserID, false)
		f.must(err)
		intent.ActorID = f.caster
		a, err := game.PrepareStoredAction(ctx, testPool, c, pc, intent, false)
		f.must(err)
		return a.Preview
	}
	check := game.ActionIntent{Kind: srd521.AbilityCheck, Ability: "wisdom"}
	stand := game.ActionIntent{Kind: combat.EncounterStand}
	unconscious := srd521.ConditionChange{Condition: srd521.ConditionUnconscious, Apply: true}
	if p := preview(check); !p.Available {
		t.Fatalf("baseline check unavailable: %+v", p)
	}

	// A foreign campaign, an unknown condition and an equipped item that may
	// be held all fail closed with no change and no roll.
	if _, err := apply(uuid.New(), nil, unconscious, true); !errors.Is(err, engine.ErrConditionActor) {
		t.Fatalf("foreign campaign: %v", err)
	}
	unchanged("foreign campaign")
	if _, err := apply(f.campaign, nil, srd521.ConditionChange{Condition: "stunned", Apply: true}, true); err == nil {
		t.Fatal("unknown condition accepted")
	}
	unchanged("unknown condition")
	if _, err := apply(f.campaign, nil, unconscious, true); err == nil || !strings.Contains(err.Error(), "requires adjudication") {
		t.Fatalf("equipped misc item: %v", err)
	}
	unchanged("equipped misc item")
	_, err = testPool.Exec(ctx, `UPDATE items SET equipped=false WHERE id=$1`, rope)
	f.must(err)
	initial = read()

	// A transaction that rolls back after the adapter leaves nothing behind.
	if _, err := apply(f.campaign, nil, unconscious, false); err != nil {
		t.Fatal(err)
	}
	unchanged("rolled-back transaction")

	// Unconscious: Prone too, Concentration ends unrolled, the mace and the
	// shield drop (and leave the build), worn chain mail stays.
	r, err := apply(f.campaign, nil, unconscious, true)
	f.must(err)
	if r.After != (srd521.Conditions{Prone: true, Unconscious: true}) || len(r.Dropped) != 2 || len(r.Concentration) != 1 || r.Concentration[0].Mechanic != srd521.ConcentrationUnconscious || dice.calls != 0 {
		t.Fatalf("receipt %+v calls=%d", r, dice.calls)
	}
	got := read()
	if got.conditions != `{"prone": true, "unconscious": true}` || got.concentration != "" || got.links != 0 || got.otherLinks != 1 || got.custom != "keep" || got.shieldBuilt || got.equipped[mace] || got.equipped[shield] || !got.equipped[mail] {
		t.Fatalf("after Unconscious: %+v", got)
	}
	for _, intent := range []game.ActionIntent{check, stand, {Kind: srd521.SpellCast, EntryID: "shield-of-faith"}} {
		if p := preview(intent); p.Available || !strings.Contains(p.Reason, "Unconscious") {
			t.Fatalf("%s while Unconscious: %+v", intent.Kind, p)
		}
	}
	// Re-applying is a no-op.
	if r, err := apply(f.campaign, nil, unconscious, true); err != nil || r.Changed || r.Dropped != nil {
		t.Fatalf("repeat: %+v %v", r, err)
	}

	// Waking leaves Prone, so standing is the next open choice.
	r, err = apply(f.campaign, nil, srd521.ConditionChange{Condition: srd521.ConditionUnconscious}, true)
	f.must(err)
	if r.After != (srd521.Conditions{Prone: true}) || read().conditions != `{"prone": true}` {
		t.Fatalf("woke: %+v", r)
	}
	if p := preview(stand); !p.Available {
		t.Fatalf("stand after waking: %+v", p)
	}

	// Incapacitated blocks actions but not standing.
	_, err = apply(f.campaign, nil, srd521.ConditionChange{Condition: srd521.ConditionIncapacitated, Apply: true}, true)
	f.must(err)
	if p := preview(check); p.Available || !strings.Contains(p.Reason, "Incapacitated") {
		t.Fatalf("check while Incapacitated: %+v", p)
	}
	if p := preview(stand); !p.Available {
		t.Fatalf("stand while Incapacitated: %+v", p)
	}

	// Inside an active encounter the player actor owns the state; stats stay.
	before := read().conditions
	encounter := &combat.Encounter{Status: "active", Round: 1, Actors: []combat.EncounterActor{{ID: f.caster, Player: true, HP: 10, Conditions: srd521.Conditions{Prone: true, Incapacitated: true}}, {ID: uuid.New(), HP: 7}}}
	r, err = apply(f.campaign, encounter, srd521.ConditionChange{Condition: srd521.ConditionIncapacitated}, true)
	f.must(err)
	if encounter.Actors[0].Conditions != (srd521.Conditions{Prone: true}) || read().conditions != before || r.After != encounter.Actors[0].Conditions {
		t.Fatalf("encounter actor %+v stats %s", encounter.Actors[0], read().conditions)
	}
	if _, err := apply(f.campaign, &combat.Encounter{Status: "active", Actors: []combat.EncounterActor{{ID: uuid.New(), Player: true}}}, unconscious, true); !errors.Is(err, engine.ErrConditionActor) {
		t.Fatalf("actor outside encounter: %v", err)
	}

	// Unknown saved condition state fails closed in previews and the adapter.
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET stats=jsonb_set(stats,'{srd521_conditions}','{"stunned":true}'::jsonb) WHERE id=$1`, f.caster)
	f.must(err)
	if p := preview(check); p.Available || !strings.Contains(p.Reason, "adjudication") {
		t.Fatalf("unknown conditions: %+v", p)
	}
	if _, err := apply(f.campaign, nil, unconscious, true); err == nil {
		t.Fatal("adapter accepted unknown condition state")
	}
}
