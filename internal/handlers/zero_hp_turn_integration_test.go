//go:build integration

package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/combat"
	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
)

// narrativeTurnError runs one free-text turn through the TUI entry point and
// the web stream, and returns the error both give. Neither may call the model.
func (f *fallFixture) narrativeTurnError() error {
	f.t.Helper()
	ctx := context.Background()
	provider := &qualificationProvider{}
	e, err := engine.New(testPool, provider, config.LLMConfig{})
	f.must(err)
	_, turnErr := e.ProcessTurn(ctx, f.campaign, "I get up and run")
	stream, err := e.ProcessTurnStream(ctx, f.campaign, "I get up and run")
	f.must(err)
	var streamErr error
	for event := range stream {
		if event.Type == "error" {
			streamErr = event.Err
		}
	}
	if provider.calls.Load() != 0 || turnErr == nil || streamErr == nil || turnErr.Error() != streamErr.Error() {
		f.t.Fatalf("turn=%v stream=%v model calls=%d", turnErr, streamErr, provider.calls.Load())
	}
	return turnErr
}

func TestIntegrationZeroHPRefusesNarrativeTurns(t *testing.T) {
	f := newZeroFixture(t)
	cannotAct := func(want string) {
		t.Helper()
		err := f.narrativeTurnError()
		if !errors.Is(err, engine.ErrPlayerCannotAct) || !strings.Contains(err.Error(), want) {
			t.Fatalf("want cannot act %q: %v", want, err)
		}
	}
	// A character with HP keeps the existing SRD guard against legacy turns.
	if err := f.narrativeTurnError(); !errors.Is(err, domain.ErrLegacyRulesMutation) || errors.Is(err, engine.ErrPlayerCannotAct) {
		t.Fatalf("healthy turn: %v", err)
	}
	f.dropByFall("f1")
	cannotAct("Unconscious at 0 HP")
	cannotAct("death-save")
	for _, die := range []int{10, 12, 19} {
		_, _, err := f.confirmWith(f.preview(srd521.DeathSave, "", ""), die)
		f.must(err)
	}
	cannotAct("Stable")
	_, err := testPool.Exec(context.Background(), `UPDATE player_characters SET stats=jsonb_set(stats,'{srd521_death_saves}','{"failures":3,"dead":true}') WHERE id=$1`, f.actor)
	f.must(err)
	cannotAct("dead")
}

// Older builds could leave a character at 0 HP without Unconscious. The
// death-save action repairs it once: Unconscious through the condition
// adapter, with its receipt in history, and no roll or time.
func TestIntegrationZeroHPRepairsMissingUnconscious(t *testing.T) {
	ctx := context.Background()
	f := newZeroFixture(t)
	f.setHP(0)
	before := f.zero()
	if before.conditions.Unconscious || !before.equipped {
		t.Fatalf("fixture: %+v", before)
	}
	if err := f.narrativeTurnError(); !errors.Is(err, engine.ErrPlayerCannotAct) || !strings.Contains(err.Error(), "without the Unconscious condition") {
		t.Fatalf("turn: %v", err)
	}
	if p := f.preview(srd521.AbilityCheck, "", "athletics"); p.Available || !strings.Contains(p.Reason, "death-save") {
		t.Fatalf("athletics at 0 HP: %+v", p)
	}
	p := f.preview(srd521.DeathSave, "", "")
	if !p.Available || !strings.Contains(p.Options[0], "applies Unconscious") || strings.Contains(p.Reason, "adjudication") {
		t.Fatalf("repair preview: %+v", p)
	}
	request, result, err := f.confirmWith(p)
	f.must(err)
	d := deathSaveReceipt(t, result)
	if d.Event != srd521.ZeroHPUnconsciousRepair || d.Die != 0 || d.HPBefore != 0 || d.HPAfter != 0 || d.Condition == nil || !d.Condition.After.Unconscious || len(d.Condition.Dropped) != 1 || f.faces.calls != 0 {
		t.Fatalf("repair receipt: %+v condition=%+v rolls=%d", d, d.Condition, f.faces.calls)
	}
	after := f.zero()
	if !after.conditions.Unconscious || !after.conditions.Prone || after.equipped || after.saves != (srd521.DeathSaves{}) || after.clock != before.clock {
		t.Fatalf("after repair: %+v", after)
	}
	// The replay returns the saved receipt; the next preview is an ordinary save.
	replayed, err := f.replayEngine().ConfirmAction(ctx, f.campaign, request)
	f.must(err)
	if replayed.Narrative != result.Narrative || f.faces.calls != 0 {
		t.Fatalf("replay: %+v", replayed)
	}
	if p := f.preview(srd521.DeathSave, "", ""); !p.Available || !strings.Contains(p.Options[0], "0 successes, 0 failures") {
		t.Fatalf("death save after repair: %+v", p)
	}
	cannot := f.narrativeTurnError()
	if !strings.Contains(cannot.Error(), "Unconscious at 0 HP") {
		t.Fatalf("turn after repair: %v", cannot)
	}
}

// Inside an active encounter the repair happens when the player's turn
// starts: Unconscious through the condition adapter, with its receipt, before
// that turn's Death Saving Throw, which then follows the ordinary rules.
func TestIntegrationZeroHPRepairsMissingUnconsciousAtEncounterTurn(t *testing.T) {
	ctx := context.Background()
	f := newZeroFixture(t)
	_, version, err := srd521.BanditRecord()
	f.must(err)
	_, err = testPool.Exec(ctx, `UPDATE npcs SET properties=$1 WHERE id=$2`, fmt.Sprintf(`{"srd521_target":{"profile":"bandit-v1","actor_id":"%s","distance_feet":5,"line_of_sight":true,"cover":"none","bandit":{"version":"%s","action_id":"scimitar"}}}`, f.actor, version), f.target)
	f.must(err)

	// The Bandit wins Initiative (1 against 20) and misses (d20 1).
	_, _, err = f.confirmWith(f.prepare(combat.EncounterStart), 1, 20, 1)
	f.must(err)
	encounter, err := game.ReadEncounter(ctx, testPool, f.campaign)
	f.must(err)
	if encounter.Status != "active" || !encounter.Active().Player || encounter.Actors[1].HP == 0 {
		t.Fatalf("encounter: %+v", encounter)
	}

	// An older build left the player at 0 HP without Unconscious.
	f.setHP(0)
	_, err = testPool.Exec(ctx, `UPDATE session_logs SET tool_calls=jsonb_set(tool_calls,'{0,Result,encounter,actors,1,hp}','0') WHERE campaign_id=$1 AND turn_number=(SELECT MAX(turn_number) FROM session_logs WHERE campaign_id=$1)`, f.campaign)
	f.must(err)
	if s := f.zero(); s.hp != 0 || s.conditions.Unconscious || !s.equipped {
		t.Fatalf("fixture: %+v", s)
	}

	// Ending the turn: the Bandit misses (d20 1). The player's next turn starts
	// with the repair, which rolls nothing, then the save: 12, one success.
	request, result, err := f.confirmWith(f.prepare(combat.EncounterEndTurn), 1, 12)
	f.must(err)
	var payload struct {
		Encounter combat.Encounter      `json:"encounter"`
		Repair    *srd521.ZeroHPReceipt `json:"unconscious_repair"`
		DeathSave *srd521.ZeroHPReceipt `json:"death_save"`
	}
	f.must(json.Unmarshal(result.AppliedToolCalls[0].Result, &payload))
	if r := payload.Repair; r == nil || r.Event != srd521.ZeroHPUnconsciousRepair || r.Die != 0 || r.Condition == nil || !r.Condition.After.Unconscious || len(r.Condition.Dropped) != 1 {
		t.Fatalf("repair receipt: %+v", payload.Repair)
	}
	if d := payload.DeathSave; d == nil || d.Die != 12 || d.After != (srd521.DeathSaves{Successes: 1}) {
		t.Fatalf("death save: %+v", payload.DeathSave)
	}
	player := payload.Encounter.Actors[1]
	if payload.Encounter.Status != "active" || !payload.Encounter.Active().Player || !player.Unconscious || !player.Prone || f.faces.calls != 2 {
		t.Fatalf("encounter after repair: %+v rolls=%d", payload.Encounter, f.faces.calls)
	}
	if !strings.Contains(result.Narrative, "without the Unconscious condition") || !strings.Contains(result.Narrative, "Death Saving Throw: d20 12, success") {
		t.Fatalf("narrative: %s", result.Narrative)
	}
	// While the encounter is active its actor holds the conditions; the held
	// Longsword is already dropped.
	if s := f.zero(); s.hp != 0 || s.saves != (srd521.DeathSaves{Successes: 1}) || s.equipped {
		t.Fatalf("saved state: %+v", s)
	}
	f.faces.calls = 0
	replayed, err := f.replayEngine().ConfirmAction(ctx, f.campaign, request)
	f.must(err)
	if f.faces.calls != 0 || replayed.Narrative != result.Narrative {
		t.Fatalf("replay rolled: %+v", replayed)
	}
}

// legacyZeroHPInEncounter leaves the player as an older build did: 0 HP
// without Unconscious, in the character row and the saved encounter.
func (f *fallFixture) legacyZeroHPInEncounter() {
	f.t.Helper()
	ctx := context.Background()
	encounter, err := game.ReadEncounter(ctx, testPool, f.campaign)
	f.must(err)
	player := 0
	if encounter.Actors[1].Player {
		player = 1
	}
	f.setHP(0)
	_, err = testPool.Exec(ctx, `UPDATE session_logs SET tool_calls=jsonb_set(tool_calls,ARRAY['0','Result','encounter','actors',$2::text,'hp'],'0') WHERE campaign_id=$1 AND turn_number=(SELECT MAX(turn_number) FROM session_logs WHERE campaign_id=$1)`, f.campaign, fmt.Sprint(player))
	f.must(err)
	encounter, err = game.ReadEncounter(ctx, testPool, f.campaign)
	f.must(err)
	if s := f.zero(); s.hp != 0 || s.conditions.Unconscious || encounter.Status != "active" || encounter.Actors[player].HP != 0 || encounter.Actors[player].Unconscious {
		f.t.Fatalf("fixture: %+v %+v", s, encounter)
	}
}

// Mid-turn in an active encounter, a character at 0 HP without Unconscious has
// the same gate as outside one: only ending the turn stays open, and every
// other choice says why.
func TestIntegrationZeroHPGatesMissingUnconsciousMidEncounterTurn(t *testing.T) {
	f := newZeroFixture(t)
	f.facts(`"distance_feet":5,"line_of_sight":true,"cover":"none"`)
	// The player wins Initiative (20 against 1).
	_, _, err := f.confirmWith(f.prepare(combat.EncounterStart), 20, 1)
	f.must(err)
	attack := f.prepare(srd521.MartialAttack)
	move := f.prepare(combat.EncounterMove)
	if !attack.Available || !move.Available {
		t.Fatalf("healthy turn: %+v %+v", attack, move)
	}
	f.legacyZeroHPInEncounter()

	const want = "without the Unconscious condition"
	for _, kind := range []string{srd521.MartialAttack, combat.EncounterMove, combat.EncounterDash, combat.EncounterApproach, combat.EncounterRetreat, combat.EncounterStand, srd521.ActionSurge, srd521.AbilityCheck} {
		entry := ""
		if kind == srd521.MartialAttack {
			entry = f.weapon.String()
		}
		if p := f.preview(kind, entry, ""); p.Available || !strings.Contains(p.Reason, want) || !strings.Contains(p.Reason, "End your turn") || strings.Contains(p.Reason, "death-save") {
			t.Errorf("%s: available=%t reason=%q", kind, p.Available, p.Reason)
		}
	}
	// Confirming a preview taken while healthy commits nothing.
	for _, p := range []*game.ActionPreview{attack, move} {
		if _, _, err := f.confirmWith(p, 15, 4); !errors.Is(err, engine.ErrActionRefresh) || f.faces.calls != 0 {
			t.Errorf("%s confirmed at 0 HP: %v rolls=%d", p.Intent.Kind, err, f.faces.calls)
		}
	}
	if p := f.prepare(combat.EncounterEndTurn); !p.Available {
		t.Fatalf("end turn: %+v", p)
	}
}

// An older build offered a reaction to a character it had left at 0 HP,
// because without Unconscious it was not Incapacitated. Ending the turn waits
// for that reaction, so declining it stays open; the next turn then repairs.
func TestIntegrationZeroHPMissingUnconsciousCanDeclinePendingReaction(t *testing.T) {
	ctx := context.Background()
	f := newZeroFixture(t)
	f.facts(`"distance_feet":5,"line_of_sight":true,"cover":"none","tactic":"withdraw"`)
	_, _, err := f.confirmWith(f.prepare(combat.EncounterStart), 20, 1)
	f.must(err)
	// The foe misses (d20 1) and starts to withdraw.
	_, _, err = f.confirmWith(f.prepare(combat.EncounterEndTurn), 1)
	f.must(err)
	encounter, err := game.ReadEncounter(ctx, testPool, f.campaign)
	f.must(err)
	if encounter.Reaction.Pending == nil {
		t.Fatalf("no pending reaction: %+v", encounter)
	}
	f.legacyZeroHPInEncounter()

	if p := f.prepare(combat.EncounterReact); p.Available || !strings.Contains(p.Reason, "without the Unconscious condition") {
		t.Errorf("reaction at 0 HP: %+v", p)
	}
	p := f.prepare(combat.EncounterDecline)
	if !p.Available {
		t.Fatalf("decline: %+v", p)
	}
	// The player's next turn starts with the repair, then the save: 12.
	_, result, err := f.confirmWith(p, 12)
	f.must(err)
	var payload struct {
		Encounter combat.Encounter      `json:"encounter"`
		Repair    *srd521.ZeroHPReceipt `json:"unconscious_repair"`
		DeathSave *srd521.ZeroHPReceipt `json:"death_save"`
	}
	f.must(json.Unmarshal(result.AppliedToolCalls[0].Result, &payload))
	if payload.Repair == nil || payload.DeathSave == nil || payload.DeathSave.Die != 12 || !payload.Encounter.Active().Player || !payload.Encounter.Active().Unconscious {
		t.Fatalf("after decline: %s", result.AppliedToolCalls[0].Result)
	}
}
