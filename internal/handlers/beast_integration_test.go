//go:build integration

package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/combat"
	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"github.com/google/uuid"
)

// beastFixture is one player and one pinned Wolf or Brown Bear in a scene.
type beastFixture struct {
	t                *testing.T
	ctx              context.Context
	dice             *banditDice
	campaign, actor  uuid.UUID
	target, location uuid.UUID
	weapon           uuid.UUID
	facts            string
}

func newBeastFixture(t *testing.T, profile, action, extra string) *beastFixture {
	t.Helper()
	ctx := context.Background()
	f := &beastFixture{t: t, ctx: ctx, dice: &banditDice{}}
	q := statedb.New(testPool)
	stats, err := os.ReadFile("../rules/srd521/testdata/fighter1.json")
	f.must(err)
	c, err := q.CreateCampaign(ctx, statedb.CreateCampaignParams{Name: "Beast fixture", RulesVersion: "srd-5.2.1", Status: "active", CreatedBy: dbutil.ToPgtype(testUserID)})
	f.must(err)
	t.Cleanup(func() { _, _ = testPool.Exec(ctx, `DELETE FROM campaigns WHERE id=$1`, c.ID) })
	pc, err := q.CreatePlayerCharacter(ctx, statedb.CreatePlayerCharacterParams{CampaignID: c.ID, UserID: dbutil.ToPgtype(testUserID), Name: "Mira", Stats: stats, Hp: 12, MaxHp: 12, Level: 1, Status: "active", Abilities: []byte(`[]`)})
	f.must(err)
	f.campaign, f.actor = dbutil.FromPgtype(c.ID), dbutil.FromPgtype(pc.ID)
	f.must(testPool.QueryRow(ctx, `INSERT INTO locations(campaign_id,name)VALUES($1,'Forest')RETURNING id`, c.ID).Scan(&f.location))
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET current_location_id=$1 WHERE id=$2`, f.location, pc.ID)
	f.must(err)
	f.must(testPool.QueryRow(ctx, `INSERT INTO items(campaign_id,player_character_id,name,item_type,rarity,equipped,properties)VALUES($1,$2,'Longsword','weapon','common',true,'{"srd521_id":"longsword"}')RETURNING id`, c.ID, pc.ID).Scan(&f.weapon))
	_, version, err := srd521.BeastRecord(profile)
	f.must(err)
	f.facts = fmt.Sprintf(`{"srd521_target":{"profile":"%s","actor_id":"%s","distance_feet":5,"line_of_sight":true,"cover":"none"%s,"beast":{"version":"%s","action_id":"%s","target_size":"medium"}}}`, profile, f.actor, extra, version, action)
	f.must(testPool.QueryRow(ctx, `INSERT INTO npcs(campaign_id,name,hp,location_id,properties)VALUES($1,'Practice beast',100,$2,$3)RETURNING id`, c.ID, f.location, f.facts).Scan(&f.target))
	_, err = testPool.Exec(ctx, `INSERT INTO session_logs(campaign_id,turn_number,player_input,input_type,llm_response,npcs_involved)VALUES($1,1,'Arrive','game_action','A beast waits.',ARRAY[$2]::uuid[])`, c.ID, f.target)
	f.must(err)
	return f
}

func (f *beastFixture) must(err error) {
	f.t.Helper()
	if err != nil {
		f.t.Fatal(err)
	}
}

func (f *beastFixture) engine() *engine.Engine {
	e, err := engine.New(testPool, nil, config.LLMConfig{}, engine.WithActionDice(f.dice))
	f.must(err)
	return e
}

func (f *beastFixture) prepare(kind, entry string) *game.ActionPreview {
	f.t.Helper()
	c, p, err := game.ReadActionSnapshot(f.ctx, testPool, f.campaign, f.actor, testUserID, false)
	f.must(err)
	a, err := game.PrepareStoredAction(f.ctx, testPool, c, p, game.ActionIntent{ActorID: f.actor, Kind: kind, EntryID: entry}, false)
	f.must(err)
	return a.Preview
}

func (f *beastFixture) request(p *game.ActionPreview, target uuid.UUID) engine.PreparedAction {
	f.t.Helper()
	if !p.Available {
		f.t.Fatalf("unavailable: %+v", p)
	}
	return engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: p.ExpectedRevision, Intent: p.Intent, TargetID: target}
}

func (f *beastFixture) confirm(r engine.PreparedAction, rolls ...int) *engine.TurnResult {
	f.t.Helper()
	f.dice.values = rolls
	result, err := f.engine().ConfirmAction(f.ctx, f.campaign, r)
	f.must(err)
	return result
}

func (f *beastFixture) encounter() *combat.Encounter {
	f.t.Helper()
	e, err := game.ReadEncounter(f.ctx, testPool, f.campaign)
	f.must(err)
	return e
}

func (f *beastFixture) hp() int {
	f.t.Helper()
	var hp int
	f.must(testPool.QueryRow(f.ctx, `SELECT hp FROM player_characters WHERE id=$1`, f.actor).Scan(&hp))
	return hp
}

func (f *beastFixture) setFacts(mutation string) {
	f.t.Helper()
	_, err := testPool.Exec(f.ctx, mutation, f.target)
	f.must(err)
}

func (f *beastFixture) reset() {
	f.t.Helper()
	_, err := testPool.Exec(f.ctx, `UPDATE npcs SET location_id=$1,properties=$2 WHERE id=$3`, f.location, f.facts, f.target)
	f.must(err)
}

func npcReceipts(t *testing.T, result *engine.TurnResult) []srd521.AttackReceipt {
	t.Helper()
	var payload struct {
		Receipts []srd521.AttackReceipt `json:"npc_receipts"`
	}
	if err := json.Unmarshal(result.AppliedToolCalls[0].Result, &payload); err != nil {
		t.Fatal(err)
	}
	return payload.Receipts
}

// refuseEndTurn checks that a changed fact invalidates a prepared End turn
// before any die is rolled, then restores the pinned facts.
func (f *beastFixture) refuseEndTurn(mutations ...string) {
	f.t.Helper()
	for _, mutation := range mutations {
		request := f.request(f.prepare(combat.EncounterEndTurn, ""), uuid.Nil)
		f.setFacts(mutation)
		f.dice.values = []int{15, 15, 3, 15, 1}
		before := f.dice.calls
		if _, err := f.engine().ConfirmAction(f.ctx, f.campaign, request); !errors.Is(err, engine.ErrActionRefresh) {
			f.t.Fatalf("mutation accepted: %s: %v", mutation, err)
		}
		if f.dice.calls != before {
			f.t.Fatalf("mutation rolled: %s", mutation)
		}
		f.reset()
	}
}

// failAtomically injects a dice failure and then a history failure; neither
// may persist HP, condition, budget, cursor or receipt. It returns the request
// so the caller can retry the same identity.
func (f *beastFixture) failAtomically(badDice, rolls []int) engine.PreparedAction {
	f.t.Helper()
	request := f.request(f.prepare(combat.EncounterEndTurn, ""), uuid.Nil)
	before, hp := f.encounter(), f.hp()
	f.dice.values = badDice
	if _, err := f.engine().ConfirmAction(f.ctx, f.campaign, request); err == nil {
		f.t.Fatal("invalid dice committed")
	}
	if !reflect.DeepEqual(before, f.encounter()) || f.hp() != hp {
		f.t.Fatal("dice failure leaked state")
	}
	_, err := testPool.Exec(f.ctx, `CREATE FUNCTION reject_beast_history() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN IF NEW.player_input='Confirm combat-end-turn combat-end-turn' THEN RAISE EXCEPTION 'fixture history failure'; END IF; RETURN NEW; END; $$; CREATE TRIGGER reject_beast_history BEFORE INSERT ON session_logs FOR EACH ROW EXECUTE FUNCTION reject_beast_history()`)
	f.must(err)
	f.dice.values = append([]int(nil), rolls...)
	_, failure := f.engine().ConfirmAction(f.ctx, f.campaign, request)
	_, err = testPool.Exec(f.ctx, `DROP TRIGGER reject_beast_history ON session_logs; DROP FUNCTION reject_beast_history()`)
	f.must(err)
	if failure == nil {
		f.t.Fatal("history failure ignored")
	}
	var attempts int
	f.must(testPool.QueryRow(f.ctx, `SELECT count(*) FROM turn_attempts WHERE id=$1`, request.ActionID).Scan(&attempts))
	if attempts != 0 || f.hp() != hp || !reflect.DeepEqual(before, f.encounter()) {
		f.t.Fatal("failed history leaked HP, condition, budget, cursor or receipt")
	}
	return request
}

// replay confirms that a completed identity returns its saved result without
// rolling or applying damage or Prone again.
func (f *beastFixture) replay(request engine.PreparedAction, result *engine.TurnResult) {
	f.t.Helper()
	calls, hp, saved := f.dice.calls, f.hp(), f.encounter()
	again, err := f.engine().ConfirmAction(f.ctx, f.campaign, request)
	f.must(err)
	if f.dice.calls != calls || !sameTurnResult(result, again) || f.hp() != hp || !reflect.DeepEqual(saved, f.encounter()) {
		f.t.Fatal("replay rerolled or reapplied effects")
	}
	receipt, err := f.engine().GetActionReceipt(f.ctx, f.campaign, request.ActionID)
	f.must(err)
	if !sameTurnResult(result, receipt.Result) {
		f.t.Fatal("receipt reload changed")
	}
}

func playerProne(e *combat.Encounter) bool {
	for _, a := range e.Actors {
		if a.Player {
			return a.Prone
		}
	}
	return false
}

func TestIntegrationWolfPackTactics(t *testing.T) {
	f := newBeastFixture(t, "wolf-v1", "bite", `,"pack_ally":{"present":true,"distance_feet":5,"incapacitated":false}`)
	// Unknown size or ally facts need clarification before the encounter starts.
	for mutation, reason := range map[string]string{
		`UPDATE npcs SET properties=properties #- '{srd521_target,beast,target_size}' WHERE id=$1`:       "size",
		`UPDATE npcs SET properties=properties #- '{srd521_target,pack_ally}' WHERE id=$1`:               "ally",
		`UPDATE npcs SET properties=properties #- '{srd521_target,pack_ally,incapacitated}' WHERE id=$1`: "Incapacitated",
	} {
		f.setFacts(mutation)
		p := f.prepare(combat.EncounterStart, f.weapon.String())
		if p.Available || len(p.Targets) != 1 || !strings.Contains(p.Targets[0].UnavailableReason, reason) {
			t.Fatalf("%s: %+v", mutation, p)
		}
		f.reset()
	}
	start := f.prepare(combat.EncounterStart, f.weapon.String())
	public, err := json.Marshal(start)
	f.must(err)
	if f.dice.calls != 0 || strings.Contains(string(public), `"hp":100`) || strings.Contains(string(public), "required_mechanics") || !strings.Contains(string(public), "Pack Tactics") || !strings.Contains(string(public), "Medium or smaller") {
		t.Fatalf("public preview=%s", public)
	}
	f.confirm(f.request(start, f.target), 18, 0)
	for _, a := range f.encounter().Actors {
		if !a.Player && (a.Initiative != 3 || a.Speed != 40) {
			t.Fatalf("actor=%+v", a)
		}
	}
	f.refuseEndTurn(
		`UPDATE npcs SET properties=jsonb_set(properties,'{srd521_target,hidden}','true') WHERE id=$1`,
		`UPDATE npcs SET properties=jsonb_set(properties,'{srd521_target,beast,version}','"wolf-v1:stale"') WHERE id=$1`,
		`UPDATE npcs SET properties=jsonb_set(properties,'{srd521_target,beast,target_size}','"large"') WHERE id=$1`,
		`UPDATE npcs SET properties=properties #- '{srd521_target,beast,target_size}' WHERE id=$1`,
		`UPDATE npcs SET properties=properties #- '{srd521_target,pack_ally}' WHERE id=$1`,
		`UPDATE npcs SET properties=jsonb_set(properties,'{srd521_target,distance_feet}','10') WHERE id=$1`,
	)
	end := f.prepare(combat.EncounterEndTurn, "")
	if !strings.Contains(strings.Join(end.Options, " "), "Advantage on each attack: Pack Tactics") {
		t.Fatalf("options=%v", end.Options)
	}
	// Advantage d20s 2 and 16 hit; 1d6 (3) + 2 Piercing; a Medium target is knocked Prone.
	request := f.failAtomically([]int{1, -1}, []int{1, 15, 2})
	result := f.confirm(request, 1, 15, 2)
	r := npcReceipts(t, result)
	if len(r) != 1 || r[0].Weapon != "Bite" || r[0].Die != 16 || !reflect.DeepEqual(r[0].D20s, []int{2, 16}) || len(r[0].Advantage) != 1 || r[0].Damage != 5 || r[0].Condition != "prone-applied" || f.hp() != 7 || !playerProne(f.encounter()) {
		t.Fatalf("receipts=%+v hp=%d", r, f.hp())
	}
	f.replay(request, result)

	// Prone: the longsword attack has Disadvantage, shown in the preview.
	attack := f.prepare(srd521.MartialAttack, f.weapon.String())
	if !strings.Contains(strings.Join(attack.Options, " "), "Disadvantage on this attack roll (attacker is Prone)") {
		t.Fatalf("attack options=%v", attack.Options)
	}
	request = f.request(attack, f.target)
	before := f.encounter()
	f.dice.values = []int{19, -1}
	if _, err := f.engine().ConfirmAction(f.ctx, f.campaign, request); err == nil {
		t.Fatal("invalid second d20 committed")
	}
	if !reflect.DeepEqual(before, f.encounter()) {
		t.Fatal("dice failure leaked state")
	}
	// d20s 20 and 4: the lower die is kept, so no critical hit; 4 + 5 misses AC 13.
	result = f.confirm(request, 19, 3)
	var payload struct {
		Receipt srd521.AttackReceipt `json:"receipt"`
	}
	f.must(json.Unmarshal(result.AppliedToolCalls[0].Result, &payload))
	if a := payload.Receipt; a.Die != 4 || !reflect.DeepEqual(a.D20s, []int{20, 4}) || !reflect.DeepEqual(a.Disadvantage, []string{"attacker is Prone"}) || a.Advantage != nil || a.Outcome != "miss" || !strings.Contains(result.Narrative, "Disadvantage (attacker is Prone): rolled 20 and 4, kept 4.") {
		t.Fatalf("attack=%+v narrative=%q", a, result.Narrative)
	}
	f.replay(request, result)
	// Crawling: 5 feet costs 10 feet of movement and keeps Prone.
	move := f.prepare(combat.EncounterMove, "")
	if !move.Available || len(move.Costs) != 1 || move.Costs[0].Amount != 10 {
		t.Fatalf("crawl=%+v", move)
	}
	f.confirm(f.request(move, uuid.Nil))
	if e := f.encounter(); !playerProne(e) || e.Budget.Movement != 20 || e.Budget.Action {
		t.Fatalf("crawl=%+v", e)
	}
	if stand := f.prepare(combat.EncounterStand, ""); !stand.Available || len(stand.Costs) != 1 || stand.Costs[0].Amount != 15 {
		t.Fatalf("stand=%+v", stand)
	}
	// Staying Prone: Pack Tactics and the Prone target still grant one Advantage.
	// The Bite hits again but does not apply Prone twice.
	result = f.confirm(f.request(f.prepare(combat.EncounterEndTurn, ""), uuid.Nil), 15, 0, 0)
	r = npcReceipts(t, result)
	if len(r) != 1 || len(r[0].Advantage) != 2 || len(r[0].D20s) != 2 || r[0].Damage != 3 || r[0].Condition != "already-prone" || f.hp() != 4 {
		t.Fatalf("receipts=%+v hp=%d", r, f.hp())
	}
	f.confirm(f.request(f.prepare(combat.EncounterStand, ""), uuid.Nil))
	if e := f.encounter(); playerProne(e) || e.Budget.Movement != 15 || !e.Budget.Action {
		t.Fatalf("stand=%+v", e)
	}
	// An Incapacitated ally is a new ruling: no Advantage, one d20.
	f.facts = strings.Replace(f.facts, `"incapacitated":false`, `"incapacitated":true`, 1)
	f.reset()
	result = f.confirm(f.request(f.prepare(combat.EncounterEndTurn, ""), uuid.Nil), 15, 0)
	r = npcReceipts(t, result)
	if len(r) != 1 || r[0].D20s != nil || r[0].Advantage != nil || r[0].Damage != 3 || r[0].Condition != "prone-applied" || f.hp() != 1 || !playerProne(f.encounter()) {
		t.Fatalf("receipts=%+v hp=%d", r, f.hp())
	}
}

func TestIntegrationBrownBearMultiattack(t *testing.T) {
	f := newBeastFixture(t, "brown-bear-v1", "multiattack", "")
	start := f.prepare(combat.EncounterStart, f.weapon.String())
	f.confirm(f.request(start, f.target), 18, 0)
	for _, a := range f.encounter().Actors {
		if !a.Player && (a.Initiative != 2 || a.Speed != 40) {
			t.Fatalf("actor=%+v", a)
		}
	}
	f.refuseEndTurn(
		`UPDATE npcs SET properties=jsonb_set(properties,'{srd521_target,line_of_sight}','false') WHERE id=$1`,
		`UPDATE npcs SET properties=jsonb_set(properties,'{srd521_target,beast,action_id}','"claw"') WHERE id=$1`,
		`UPDATE npcs SET properties=properties #- '{srd521_target,beast,target_size}' WHERE id=$1`,
		`UPDATE npcs SET location_id=NULL WHERE id=$1`,
	)
	end := f.prepare(combat.EncounterEndTurn, "")
	options := strings.Join(end.Options, "\n")
	if strings.Index(options, "1. Bite") < 0 || strings.Index(options, "2. Claw") < strings.Index(options, "1. Bite") || !strings.Contains(options, "ends the sequence") {
		t.Fatalf("options=%s", options)
	}
	// The Bite hits (7 damage written in the transaction) before the Claw's d20
	// fails, so nothing from either attack may persist.
	request := f.failAtomically([]int{15, 3, -1}, []int{0, 15, 1})
	result := f.confirm(request, 0, 15, 1)
	r := npcReceipts(t, result)
	if len(r) != 2 || r[0].Weapon != "Bite" || r[0].Outcome != "miss" || r[1].Weapon != "Claw" || r[1].Damage != 5 || r[1].DamageType != "slashing" || r[1].Condition != "prone-applied" || f.hp() != 7 {
		t.Fatalf("receipts=%+v hp=%d", r, f.hp())
	}
	e := f.encounter()
	if !playerProne(e) || !e.Active().Player || e.Round != 2 || !e.Budget.Action {
		t.Fatalf("encounter=%+v", e)
	}
	f.replay(request, result)
	// Against the Prone target each attack has Advantage. The Bite's 11 damage
	// defeats the player, so the Claw is never rolled.
	calls := f.dice.calls
	result = f.confirm(f.request(f.prepare(combat.EncounterEndTurn, ""), uuid.Nil), 15, 15, 7, 15, 1)
	r = npcReceipts(t, result)
	if len(r) != 1 || r[0].Weapon != "Bite" || r[0].Damage != 11 || len(r[0].D20s) != 2 || f.dice.calls != calls+3 || f.hp() != 0 || f.encounter().Status != "completed" {
		t.Fatalf("receipts=%+v hp=%d calls=%d", r, f.hp(), f.dice.calls-calls)
	}
}

// A Prone player's opportunity attack has Disadvantage, and crawling back into
// reach costs double movement.
func TestIntegrationProneOpportunityAttack(t *testing.T) {
	f := newBeastFixture(t, "wolf-v1", "bite", `,"pack_ally":{"present":false},"tactic":"withdraw"`)
	f.confirm(f.request(f.prepare(combat.EncounterStart, f.weapon.String()), f.target), 18, 0)
	// The Bite hits (d20 16), deals 1d6 (3) + 2 and knocks the player Prone;
	// the Wolf then starts to withdraw.
	f.confirm(f.request(f.prepare(combat.EncounterEndTurn, ""), uuid.Nil), 15, 2)
	if e := f.encounter(); !playerProne(e) || e.Reaction.Pending == nil {
		t.Fatalf("encounter=%+v", e)
	}
	react := f.prepare(combat.EncounterReact, f.weapon.String())
	if !react.Available || !strings.Contains(strings.Join(react.Options, " "), "Disadvantage on this attack roll (attacker is Prone)") {
		t.Fatalf("react=%+v", react)
	}
	// d20s 20 and 10: the lower die is kept, so the hit is not critical.
	request := f.request(react, f.target)
	result := f.confirm(request, 19, 9, 3)
	var payload struct {
		Receipt srd521.AttackReceipt `json:"receipt"`
	}
	f.must(json.Unmarshal(result.AppliedToolCalls[0].Result, &payload))
	if a := payload.Receipt; a.Die != 10 || a.Outcome != "hit" || len(a.DamageDice) != 1 || !reflect.DeepEqual(a.D20s, []int{20, 10}) || len(a.Disadvantage) != 1 {
		t.Fatalf("opportunity attack=%+v", a)
	}
	f.replay(request, result)
	e := f.encounter()
	if !playerProne(e) || !e.Active().Player || *e.Space.DistanceFeet != 10 || e.Budget.Movement != 30 {
		t.Fatalf("after withdrawal=%+v", e)
	}
	approach := f.prepare(combat.EncounterApproach, "")
	if !approach.Available || approach.Costs[0].Amount != 10 || !strings.Contains(strings.Join(approach.Options, " "), "crawl") {
		t.Fatalf("approach=%+v", approach)
	}
	f.confirm(f.request(approach, uuid.Nil))
	if e := f.encounter(); !playerProne(e) || *e.Space.DistanceFeet != 5 || e.Budget.Movement != 20 {
		t.Fatalf("crawl=%+v", e)
	}
}
