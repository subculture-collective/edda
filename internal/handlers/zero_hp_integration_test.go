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
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
)

// zeroState is the player's saved 0-HP state.
type zeroState struct {
	hp         int
	saves      srd521.DeathSaves
	conditions srd521.Conditions
	equipped   bool
	clock      int
}

func newZeroFixture(t *testing.T) *fallFixture {
	return &fallFixture{spatialFixture: newSpatialFixture(t), faces: &fallFaces{}}
}

func (f *fallFixture) zero() zeroState {
	f.t.Helper()
	ctx := context.Background()
	var s zeroState
	var stats []byte
	// equipped is the fixture's Longsword, whatever other items a test adds.
	f.must(testPool.QueryRow(ctx, `SELECT p.hp,p.stats,i.equipped FROM player_characters p JOIN items i ON i.player_character_id=p.id AND i.id=$2 WHERE p.id=$1`, f.actor, f.weapon).Scan(&s.hp, &stats, &s.equipped))
	var err error
	s.saves, err = srd521.DecodeDeathSaves(stats)
	f.must(err)
	s.conditions, err = srd521.DecodeConditions(stats)
	f.must(err)
	s.clock, err = game.ReadSpellClock(ctx, testPool, f.campaign)
	f.must(err)
	return s
}

// dropByFall drops the 12-HP fighter to 0 HP: 100 feet, ten 2s, 20 damage,
// 8 remaining, under the 12 HP maximum.
func (f *fallFixture) dropByFall(id string) *engine.TurnResult {
	f.t.Helper()
	f.declare(fallDecl(id, 100, "solid"))
	_, result, err := f.confirmWith(f.preview(srd521.FallResolve, id, ""), 2, 2, 2, 2, 2, 2, 2, 2, 2, 2)
	f.must(err)
	return result
}

// replayEngine counts any roll on the fixture's dice.
func (f *fallFixture) replayEngine() *engine.Engine {
	f.t.Helper()
	e, err := engine.New(testPool, nil, config.LLMConfig{}, engine.WithActionDice(f.faces))
	f.must(err)
	return e
}

func deathSaveReceipt(t *testing.T, result *engine.TurnResult) srd521.ZeroHPReceipt {
	t.Helper()
	if len(result.AppliedToolCalls) != 1 || result.AppliedToolCalls[0].Tool != "srd_death_save" {
		t.Fatalf("tool calls: %+v", result.AppliedToolCalls)
	}
	var payload struct {
		Receipt srd521.ZeroHPReceipt `json:"receipt"`
	}
	if err := json.Unmarshal(result.AppliedToolCalls[0].Result, &payload); err != nil {
		t.Fatal(err)
	}
	return payload.Receipt
}

func TestIntegrationZeroHPOutsideEncounter(t *testing.T) {
	ctx := context.Background()
	f := newZeroFixture(t)
	unavailable := func(p *game.ActionPreview, reason string) {
		t.Helper()
		if p.Available || !strings.Contains(p.Reason, reason) {
			t.Fatalf("want unavailable %q: %+v", reason, p)
		}
	}
	unavailable(f.preview(srd521.DeathSave, "", ""), "only at 0 HP")

	// Dropping to 0 HP by a fall: Unconscious and Prone, the held Longsword
	// dropped, and no counts yet.
	result := f.dropByFall("f1")
	r := fallReceipt(t, result)
	if r.HPAfter != 0 || r.ZeroHP == nil || r.ZeroHP.Event != srd521.ZeroHPDropped || r.ZeroHP.Remaining != 8 || r.ZeroHP.Condition == nil || !r.ZeroHP.Condition.After.Unconscious {
		t.Fatalf("fall receipt: %+v zero=%+v", r, r.ZeroHP)
	}
	if strings.Contains(result.Narrative, "adjudication") || !strings.Contains(result.Narrative, "Unconscious") {
		t.Fatalf("narrative: %s", result.Narrative)
	}
	s := f.zero()
	if s.hp != 0 || !s.conditions.Unconscious || !s.conditions.Prone || s.equipped || s.saves != (srd521.DeathSaves{}) {
		t.Fatalf("after drop: %+v", s)
	}
	unavailable(f.preview(srd521.AbilityCheck, "", "athletics"), "Unconscious")
	unavailable(f.preview(combat.EncounterStand, "", ""), "Unconscious")

	// The preview never rolls; a confirmation rolls one d20 and advances one round.
	start := s.clock
	p := f.preview(srd521.DeathSave, "", "")
	if !p.Available || p.Modifier != nil || !strings.Contains(p.Options[0], "0 successes, 0 failures") {
		t.Fatalf("death save preview: %+v", p)
	}
	request, result, err := f.confirmWith(p, 9)
	f.must(err)
	if d := deathSaveReceipt(t, result); d.Die != 9 || d.Outcome != "failure" || d.After != (srd521.DeathSaves{Failures: 1}) || d.TimeAfter-d.TimeBefore != 6 {
		t.Fatalf("failure receipt: %+v", d)
	}
	if s := f.zero(); s.saves != (srd521.DeathSaves{Failures: 1}) || s.clock != start+6 {
		t.Fatalf("after failure: %+v", s)
	}
	// Replays never roll again; the stale preview needs a refresh.
	f.faces.calls = 0
	replayed, err := f.replayEngine().ConfirmAction(ctx, f.campaign, request)
	f.must(err)
	if f.faces.calls != 0 || replayed.Narrative != result.Narrative || f.zero().saves.Failures != 1 {
		t.Fatalf("replay rolled: %+v", replayed)
	}
	if _, _, err := f.confirmWith(p, 10); !errors.Is(err, engine.ErrActionRefresh) || f.faces.calls != 0 {
		t.Fatalf("stale death save: %v", err)
	}

	// A history failure commits nothing.
	_, err = testPool.Exec(ctx, `ALTER TABLE session_logs ADD CONSTRAINT death_save_history_failure CHECK (player_input NOT LIKE 'Confirm Death saving throw%') NOT VALID`)
	f.must(err)
	if _, _, err := f.confirmWith(f.preview(srd521.DeathSave, "", ""), 15); err == nil {
		t.Fatal("history failure committed")
	}
	_, err = testPool.Exec(ctx, `ALTER TABLE session_logs DROP CONSTRAINT death_save_history_failure`)
	f.must(err)
	if s := f.zero(); s.saves.Failures != 1 || s.saves.Successes != 0 || s.clock != start+6 {
		t.Fatalf("history failure partially committed: %+v", s)
	}

	// A natural 20 regains 1 HP, ends Unconscious, keeps Prone and clears the counts.
	_, result, err = f.confirmWith(f.preview(srd521.DeathSave, "", ""), 20)
	f.must(err)
	if d := deathSaveReceipt(t, result); d.Outcome != "regain-1-hp" || d.HPAfter != 1 || d.Condition == nil || d.Condition.After.Unconscious {
		t.Fatalf("natural 20: %+v", d)
	}
	if s := f.zero(); s.hp != 1 || s.conditions != (srd521.Conditions{Prone: true}) || s.saves != (srd521.DeathSaves{}) {
		t.Fatalf("after natural 20: %+v", s)
	}
	unavailable(f.preview(srd521.DeathSave, "", ""), "only at 0 HP")

	// Three successes make the character Stable; a Stable character recovers
	// 1 HP after 1d4 hours.
	f.setHP(12)
	f.dropByFall("f2")
	for _, die := range []int{10, 12, 19} {
		_, _, err = f.confirmWith(f.preview(srd521.DeathSave, "", ""), die)
		f.must(err)
	}
	if s := f.zero(); s.saves != (srd521.DeathSaves{Stable: true}) || !s.conditions.Unconscious {
		t.Fatalf("after three successes: %+v", s)
	}
	p = f.preview(srd521.DeathSave, "", "")
	if !p.Available || !strings.Contains(p.Options[0], "1d4 hours") {
		t.Fatalf("stable preview: %+v", p)
	}
	before := f.zero().clock
	_, result, err = f.confirmWith(p, 2)
	f.must(err)
	if d := deathSaveReceipt(t, result); d.Event != srd521.ZeroHPStableRecovery || d.RecoveryHours != 2 {
		t.Fatalf("recovery: %+v", d)
	}
	if s := f.zero(); s.hp != 1 || s.conditions.Unconscious || s.saves != (srd521.DeathSaves{}) || s.clock != before+7200 {
		t.Fatalf("after recovery: %+v", s)
	}

	// A 1 counts as two failures; the third failure kills the character, and
	// nothing is available afterwards.
	f.setHP(12)
	f.dropByFall("f3")
	_, _, err = f.confirmWith(f.preview(srd521.DeathSave, "", ""), 1)
	f.must(err)
	_, result, err = f.confirmWith(f.preview(srd521.DeathSave, "", ""), 5)
	f.must(err)
	if d := deathSaveReceipt(t, result); d.Outcome != "dead" || !d.After.Dead || d.After.Failures != 3 {
		t.Fatalf("death: %+v", d)
	}
	unavailable(f.preview(srd521.DeathSave, "", ""), "dead")
	unavailable(f.preview(srd521.AbilityCheck, "", "athletics"), "dead")
	if got := game.ZeroHPState([]byte(`{"srd521_death_saves":{"failures":3,"dead":true}}`), 0); got == nil || !got.Dead {
		t.Fatalf("sheet state: %+v", got)
	}
}

func TestIntegrationZeroHPMassiveDamage(t *testing.T) {
	f := newZeroFixture(t)
	// 200 feet, twenty 2s: 40 damage, 28 remaining at a 12 HP maximum.
	f.declare(fallDecl("f1", 200, "solid"))
	twos := make([]int, 20)
	for i := range twos {
		twos[i] = 2
	}
	_, result, err := f.confirmWith(f.preview(srd521.FallResolve, "f1", ""), twos...)
	f.must(err)
	if r := fallReceipt(t, result); r.ZeroHP == nil || r.ZeroHP.Event != srd521.ZeroHPMassiveDamage || !r.ZeroHP.After.Dead || r.ZeroHP.Condition != nil {
		t.Fatalf("massive damage: %+v", r.ZeroHP)
	}
	if s := f.zero(); !s.saves.Dead || s.hp != 0 {
		t.Fatalf("after massive damage: %+v", s)
	}
	if p := f.preview(srd521.DeathSave, "", ""); p.Available || !strings.Contains(p.Reason, "dead") {
		t.Fatalf("death save after death: %+v", p)
	}
}

func TestIntegrationZeroHPBanditEncounter(t *testing.T) {
	ctx := context.Background()
	f := newZeroFixture(t)
	_, version, err := srd521.BanditRecord()
	f.must(err)
	_, err = testPool.Exec(ctx, `UPDATE npcs SET properties=$1 WHERE id=$2`, fmt.Sprintf(`{"srd521_target":{"profile":"bandit-v1","actor_id":"%s","distance_feet":5,"line_of_sight":true,"cover":"none","bandit":{"version":"%s","action_id":"scimitar"}}}`, f.actor, version), f.target)
	f.must(err)
	f.setHP(1)

	// The Bandit wins Initiative (1 against 20) and hits with its Scimitar
	// (18 + 3; 1d6 4 + 1): 1 HP to 0, Unconscious. The player's turn then
	// starts with a Death Saving Throw: 12, one success.
	_, result, err := f.confirmWith(f.prepare(combat.EncounterStart), 1, 20, 18, 4, 12)
	f.must(err)
	encounter, err := game.ReadEncounter(ctx, testPool, f.campaign)
	f.must(err)
	player := encounter.Actors[1]
	if encounter.Status != "active" || !encounter.Active().Player || !player.Unconscious || !player.Prone || player.HP != 0 || encounter.Budget.Movement != 0 {
		t.Fatalf("encounter after drop: %+v", encounter)
	}
	for _, want := range []string{"Dropped to 0 HP", "Death Saving Throw: d20 12, success"} {
		if !strings.Contains(result.Narrative, want) {
			t.Fatalf("narrative missing %q: %s", want, result.Narrative)
		}
	}
	if s := f.zero(); s.hp != 0 || s.saves != (srd521.DeathSaves{Successes: 1}) || s.equipped {
		t.Fatalf("saved state: %+v", s)
	}
	if p := f.prepare(srd521.MartialAttack); p.Available {
		t.Fatalf("attack at 0 HP: %+v", p)
	}
	if p := f.preview(srd521.DeathSave, "", ""); p.Available || !strings.Contains(p.Reason, "start of your turn") {
		t.Fatalf("explicit save in encounter: %+v", p)
	}

	// Ending the turn no longer stalls. The Bandit attacks with Advantage
	// (5 and 18 keep 18, + 3 hits) and, within 5 feet, scores an automatic
	// Critical Hit (2d6 3 + 3, + 1): two failures. The next turn's save, 2,
	// is the third failure.
	end := f.prepare(combat.EncounterEndTurn)
	request, result, err := f.confirmWith(end, 5, 18, 3, 3, 2)
	f.must(err)
	var payload struct {
		Encounter   combat.Encounter        `json:"encounter"`
		NPCReceipts []*srd521.AttackReceipt `json:"npc_receipts"`
		DeathSave   *srd521.ZeroHPReceipt   `json:"death_save"`
	}
	f.must(json.Unmarshal(result.AppliedToolCalls[0].Result, &payload))
	if len(payload.NPCReceipts) != 1 {
		t.Fatalf("npc receipts: %+v", payload.NPCReceipts)
	}
	attack := payload.NPCReceipts[0]
	if len(attack.Advantage) != 2 || len(attack.D20s) != 2 || attack.Die != 18 || attack.Outcome != "critical-hit" || attack.AutoCritical == "" || attack.Damage != 7 || attack.ZeroHP == nil || attack.ZeroHP.Outcome != "two-failures" {
		t.Fatalf("attack: %+v zero=%+v", attack, attack.ZeroHP)
	}
	if payload.DeathSave == nil || payload.DeathSave.Outcome != "dead" || payload.Encounter.Status != "completed" {
		t.Fatalf("death save: %+v status=%s", payload.DeathSave, payload.Encounter.Status)
	}
	if !strings.Contains(result.Narrative, "Automatic Critical Hit") || !strings.Contains(result.Narrative, "the character dies") {
		t.Fatalf("narrative: %s", result.Narrative)
	}
	// The completed encounter hands Unconscious back; the dead character has
	// no choices, and a replay rolls nothing.
	if s := f.zero(); !s.saves.Dead || !s.conditions.Unconscious {
		t.Fatalf("after death: %+v", s)
	}
	if p := f.preview(srd521.AbilityCheck, "", "athletics"); p.Available || !strings.Contains(p.Reason, "dead") {
		t.Fatalf("dead character acted: %+v", p)
	}
	f.faces.calls = 0
	replayed, err := f.replayEngine().ConfirmAction(ctx, f.campaign, request)
	f.must(err)
	if f.faces.calls != 0 || replayed.Narrative != result.Narrative {
		t.Fatalf("replay rolled: %+v", replayed)
	}
}
