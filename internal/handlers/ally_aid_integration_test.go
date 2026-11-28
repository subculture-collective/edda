//go:build integration

package handlers_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
)

// aidState is the player's saved state after an ally's aid.
type aidState struct {
	hp, quantity, clock int
	saves               srd521.DeathSaves
	conditions          srd521.Conditions
}

func TestIntegrationAllyAidAtZeroHP(t *testing.T) {
	ctx := context.Background()
	f := newZeroFixture(t)
	var location uuid.UUID
	f.must(testPool.QueryRow(ctx, `SELECT current_location_id FROM player_characters WHERE id=$1`, f.actor).Scan(&location))
	var potion, ally uuid.UUID
	f.must(testPool.QueryRow(ctx, `INSERT INTO items(campaign_id,player_character_id,name,item_type,rarity,quantity,properties) VALUES($1,$2,'Potion of Healing','consumable','common',2,'{"srd521_id":"potion-of-healing"}') RETURNING id`, f.campaign, f.actor).Scan(&potion))
	f.must(testPool.QueryRow(ctx, `INSERT INTO npcs(campaign_id,name,hp,location_id) VALUES($1,'Rowan',8,$2) RETURNING id`, f.campaign, location).Scan(&ally))
	_, err := testPool.Exec(ctx, `INSERT INTO session_logs(campaign_id,turn_number,player_input,input_type,llm_response,npcs_involved) VALUES($1,2,'Meet','game_action','Rowan waits.',ARRAY[$2]::uuid[])`, f.campaign, ally)
	f.must(err)
	facts := func(extra string) string {
		return fmt.Sprintf(`{"profile":"aid-player-v1","actor_id":"%s","location_id":"%s","allied":true,"willing":true,"conscious":true,"hidden":false,"conditions":"ordinary","max_hp":8,"distance_feet":5,"line_of_sight":true,"cover":"none"%s}`, f.actor, location, extra)
	}
	setFacts := func(value string) {
		t.Helper()
		_, err := testPool.Exec(ctx, `UPDATE npcs SET properties=jsonb_build_object('srd521_ally',$1::jsonb,'secret','unpublished villain') WHERE id=$2`, value, ally)
		f.must(err)
	}
	state := func() aidState {
		t.Helper()
		var s aidState
		var stats []byte
		f.must(testPool.QueryRow(ctx, `SELECT p.hp,p.stats,i.quantity FROM player_characters p JOIN items i ON i.id=$2 WHERE p.id=$1`, f.actor, potion).Scan(&s.hp, &stats, &s.quantity))
		var err error
		s.saves, err = srd521.DecodeDeathSaves(stats)
		f.must(err)
		s.conditions, err = srd521.DecodeConditions(stats)
		f.must(err)
		s.clock, err = game.ReadSpellClock(ctx, testPool, f.campaign)
		f.must(err)
		return s
	}
	confirm := func(p *game.ActionPreview, target uuid.UUID, faces ...int) (engine.PreparedAction, *engine.TurnResult, error) {
		t.Helper()
		f.faces.faces, f.faces.calls = faces, 0
		e, err := engine.New(testPool, nil, config.LLMConfig{}, engine.WithActionDice(f.faces))
		f.must(err)
		r := engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: p.ExpectedRevision, Intent: p.Intent, TargetID: target}
		result, err := e.ConfirmAction(ctx, f.campaign, r)
		return r, result, err
	}
	unavailable := func(p *game.ActionPreview, reason string) {
		t.Helper()
		if p.Available || !strings.Contains(p.Reason, reason) {
			t.Fatalf("want unavailable %q: %+v", reason, p)
		}
	}
	medicine := func() *game.ActionPreview { return f.preview(srd521.AllyAid, srd521.AllyAidMedicine, "") }
	setFacts(facts(`,"medicine_modifier":2`))
	unavailable(medicine(), "at 0 HP")

	f.dropByFall("aid-fall")
	start := state()
	if start.hp != 0 || !start.conditions.Unconscious {
		t.Fatalf("not dropped: %+v", start)
	}
	// The Unconscious player still cannot drink the potion themselves.
	unavailable(f.preview(srd521.ConsumableUse, potion.String(), ""), "Unconscious")

	// Every fact is required, exact and ordinary; the profile is the aid one.
	for name, value := range map[string]string{
		"missing medicine":   facts(""),
		"float medicine":     facts(`,"medicine_modifier":2.5`),
		"string medicine":    facts(`,"medicine_modifier":"2"`),
		"out of range":       facts(`,"medicine_modifier":29`),
		"extra key":          facts(`,"medicine_modifier":2,"dc":5`),
		"unconscious helper": strings.Replace(facts(`,"medicine_modifier":2`), `"conscious":true`, `"conscious":false`, 1),
		"unwilling":          strings.Replace(facts(`,"medicine_modifier":2`), `"willing":true`, `"willing":false`, 1),
		"out of reach":       strings.Replace(facts(`,"medicine_modifier":2`), `"distance_feet":5`, `"distance_feet":10`, 1),
		"other player":       strings.Replace(facts(`,"medicine_modifier":2`), f.actor.String(), uuid.NewString(), 1),
		"potion profile":     strings.Replace(facts(`,"medicine_modifier":2`), `aid-player-v1`, `potion-healing-ally-v1`, 1),
	} {
		setFacts(value)
		if p := medicine(); p.Available || !strings.Contains(p.Reason, "No ally") {
			t.Fatalf("%s offered: %+v", name, p)
		}
	}
	setFacts(facts(`,"medicine_modifier":2`))

	// Preview never rolls and offers only the prepared ally; no hidden facts.
	f.faces.calls = 0
	p := medicine()
	if !p.Available || len(p.Targets) != 1 || p.Targets[0].ID != ally || !strings.Contains(p.Targets[0].Label, "(Medicine) +2") || len(p.Costs) != 0 || f.faces.calls != 0 {
		t.Fatalf("medicine preview: %+v", p)
	}
	if strings.Contains(fmt.Sprint(p), "unpublished") {
		t.Fatalf("preview leaked: %+v", p)
	}
	// Only the prepared ally crosses confirmation; the player cannot be the helper.
	for _, target := range []uuid.UUID{uuid.Nil, f.actor, f.target} {
		if _, _, err := confirm(p, target, 20); !errors.Is(err, engine.ErrActionRefresh) || f.faces.calls != 0 {
			t.Fatalf("target %s: %v", target, err)
		}
	}

	// A failed check (5 + 2 = 7) changes nothing; the next turn's Death Saving
	// Throw (9) fails. One round passes.
	request, result, err := confirm(p, ally, 5, 9)
	f.must(err)
	r := deathSaveReceipt(t, result)
	if r.Event != srd521.ZeroHPAllyStabilize || r.Outcome != "failure" || r.Total != 7 || r.Helper != "Rowan" || r.DeathSave == nil || r.DeathSave.Die != 9 || r.DeathSave.After != (srd521.DeathSaves{Failures: 1}) || r.TimeAfter-r.TimeBefore != 6 || f.faces.calls != 2 {
		t.Fatalf("failed aid: %+v save=%+v", r, r.DeathSave)
	}
	if s := state(); s.saves != (srd521.DeathSaves{Failures: 1}) || s.clock != start.clock+6 || s.hp != 0 {
		t.Fatalf("after failed aid: %+v", s)
	}
	f.faces.calls = 0
	replayed, err := f.replayEngine().ConfirmAction(ctx, f.campaign, request)
	f.must(err)
	if f.faces.calls != 0 || replayed.Narrative != result.Narrative || state().saves.Failures != 1 {
		t.Fatalf("replay rolled: %+v", replayed)
	}
	if _, _, err := confirm(p, ally, 18); !errors.Is(err, engine.ErrActionRefresh) || f.faces.calls != 0 {
		t.Fatalf("stale aid: %v", err)
	}
	// A changed ally (here its HP) needs a fresh preview.
	p = medicine()
	_, err = testPool.Exec(ctx, `UPDATE npcs SET hp=7,updated_at=clock_timestamp() WHERE id=$1`, ally)
	f.must(err)
	if _, _, err := confirm(p, ally, 18); !errors.Is(err, engine.ErrActionRefresh) || f.faces.calls != 0 {
		t.Fatalf("changed ally: %v", err)
	}

	// A history failure commits nothing.
	_, err = testPool.Exec(ctx, `ALTER TABLE session_logs ADD CONSTRAINT ally_aid_history_failure CHECK (player_input NOT LIKE 'Confirm Ally aid%') NOT VALID`)
	f.must(err)
	before := state()
	if _, _, err := confirm(medicine(), ally, 18); err == nil {
		t.Fatal("history failure committed")
	}
	_, err = testPool.Exec(ctx, `ALTER TABLE session_logs DROP CONSTRAINT ally_aid_history_failure`)
	f.must(err)
	if s := state(); s != before {
		t.Fatalf("rollback: %+v want %+v", s, before)
	}

	// Success (8 + 2 = 10): Stable and the counts reset; no Death Saving Throw.
	_, result, err = confirm(medicine(), ally, 8)
	f.must(err)
	if r := deathSaveReceipt(t, result); r.Outcome != "stable" || r.After != (srd521.DeathSaves{Stable: true}) || r.DeathSave != nil || f.faces.calls != 1 {
		t.Fatalf("stabilized: %+v", r)
	}
	if s := state(); s.saves != (srd521.DeathSaves{Stable: true}) || !s.conditions.Unconscious || s.clock != start.clock+12 {
		t.Fatalf("after stabilizing: %+v", s)
	}
	unavailable(medicine(), "already Stable")

	// The ally gives the Stable player a potion (2 + 3 + 2 = 7 HP): Unconscious
	// ends, Prone remains, the state clears and one potion is used.
	f.faces.calls = 0
	p = f.preview(srd521.AllyAid, potion.String(), "")
	if !p.Available || len(p.Targets) != 1 || len(p.Costs) != 1 || p.Costs[0].ResourceID != "item:"+potion.String() || f.faces.calls != 0 {
		t.Fatalf("potion preview: %+v", p)
	}
	_, result, err = confirm(p, ally, 2, 3)
	f.must(err)
	r = deathSaveReceipt(t, result)
	if r.Event != srd521.ZeroHPAllyHealing || r.HPAfter != 7 || r.Healing == nil || r.Healing.QuantityAfter != 1 || r.Condition == nil || r.Condition.After.Unconscious || !r.Condition.After.Prone || r.DeathSave != nil || r.After != (srd521.DeathSaves{}) {
		t.Fatalf("potion: %+v", r)
	}
	if s := state(); s.hp != 7 || s.quantity != 1 || s.saves != (srd521.DeathSaves{}) || s.conditions.Unconscious || !s.conditions.Prone || s.clock != start.clock+18 {
		t.Fatalf("after potion: %+v", s)
	}
	unavailable(medicine(), "at 0 HP")
	// The ally's HP and the player's own potion entry are unchanged rules.
	if p := f.preview(srd521.ConsumableUse, potion.String(), ""); !p.Available {
		t.Fatalf("conscious self-use: %+v", p)
	}
}
