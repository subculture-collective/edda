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

type banditDice struct {
	values []int
	calls  int
}

func (d *banditDice) Intn(n int) int {
	d.calls++
	if len(d.values) == 0 {
		return -1
	}
	v := d.values[0]
	d.values = d.values[1:]
	return v
}
func TestIntegrationBandit(t *testing.T) {
	for _, action := range []string{"scimitar", "light-crossbow"} {
		t.Run(action, func(t *testing.T) {
			ctx := context.Background()
			must := func(err error) {
				t.Helper()
				if err != nil {
					t.Fatal(err)
				}
			}
			q := statedb.New(testPool)
			stats, err := os.ReadFile("../rules/srd521/testdata/fighter1.json")
			must(err)
			c, err := q.CreateCampaign(ctx, statedb.CreateCampaignParams{Name: "Encounter fixture", RulesVersion: "srd-5.2.1", Status: "active", CreatedBy: dbutil.ToPgtype(testUserID)})
			must(err)
			t.Cleanup(func() { _, _ = testPool.Exec(ctx, `DELETE FROM campaigns WHERE id=$1`, c.ID) })
			pc, err := q.CreatePlayerCharacter(ctx, statedb.CreatePlayerCharacterParams{CampaignID: c.ID, UserID: dbutil.ToPgtype(testUserID), Name: "Mira", Stats: stats, Hp: 12, MaxHp: 12, Level: 1, Status: "active", Abilities: []byte(`[]`)})
			must(err)
			campaign, actor := dbutil.FromPgtype(c.ID), dbutil.FromPgtype(pc.ID)
			var location, target, weapon uuid.UUID
			must(testPool.QueryRow(ctx, `INSERT INTO locations(campaign_id,name)VALUES($1,'Court')RETURNING id`, c.ID).Scan(&location))
			_, err = testPool.Exec(ctx, `UPDATE player_characters SET current_location_id=$1 WHERE id=$2`, location, pc.ID)
			must(err)
			must(testPool.QueryRow(ctx, `INSERT INTO items(campaign_id,player_character_id,name,item_type,rarity,equipped,properties)VALUES($1,$2,'Longsword','weapon','common',true,'{"srd521_id":"longsword"}')RETURNING id`, c.ID, pc.ID).Scan(&weapon))
			_, version, err := srd521.BanditRecord()
			must(err)
			distance := 5
			if action == "light-crossbow" {
				distance = 30
			}
			facts := fmt.Sprintf(`{"srd521_target":{"profile":"bandit-v1","actor_id":"%s","distance_feet":%d,"line_of_sight":true,"cover":"none","bandit":{"version":"%s","action_id":"%s"}}}`, actor, distance, version, action)
			must(testPool.QueryRow(ctx, `INSERT INTO npcs(campaign_id,name,hp,location_id,properties)VALUES($1,'Practice foe',100,$2,$3)RETURNING id`, c.ID, location, facts).Scan(&target))
			_, err = testPool.Exec(ctx, `INSERT INTO session_logs(campaign_id,turn_number,player_input,input_type,llm_response,npcs_involved)VALUES($1,1,'Arrive','game_action','Practice foe waits.',ARRAY[$2]::uuid[])`, c.ID, target)
			must(err)
			dice := &banditDice{values: []int{18, 0}}
			newEngine := func() *engine.Engine {
				e, err := engine.New(testPool, nil, config.LLMConfig{}, engine.WithActionDice(dice))
				must(err)
				return e
			}
			prepare := func(kind, entry string) *game.ActionPreview {
				t.Helper()
				c, p, err := game.ReadActionSnapshot(ctx, testPool, campaign, actor, testUserID, false)
				must(err)
				a, err := game.PrepareStoredAction(ctx, testPool, c, p, game.ActionIntent{ActorID: actor, Kind: kind, EntryID: entry}, false)
				must(err)
				return a.Preview
			}
			confirm := func(p *game.ActionPreview, targetID uuid.UUID) (engine.PreparedAction, *engine.TurnResult) {
				t.Helper()
				if !p.Available {
					t.Fatalf("unavailable: %+v", p)
				}
				r := engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: p.ExpectedRevision, Intent: p.Intent, TargetID: targetID}
				result, err := newEngine().ConfirmAction(ctx, campaign, r)
				must(err)
				return r, result
			}

			start := prepare(combat.EncounterStart, weapon.String())
			if dice.calls != 0 || !start.Available {
				t.Fatalf("preview=%+v dice=%d", start, dice.calls)
			}
			public, err := json.Marshal(start)
			must(err)
			if strings.Contains(string(public), `"hp":100`) || strings.Contains(string(public), "required_mechanics") || !strings.Contains(string(public), "Scimitar") || !strings.Contains(string(public), "Light Crossbow") {
				t.Fatalf("public preview=%s", public)
			}
			confirm(start, target)
			encounter, err := game.ReadEncounter(ctx, testPool, campaign)
			must(err)
			for _, a := range encounter.Actors {
				if !a.Player && (a.Initiative != 2 || a.Speed != 30) {
					t.Fatalf("actor=%+v", a)
				}
			}
			// Each mutation must invalidate confirmation before rolling, even after reload.
			for _, mutation := range []string{
				`UPDATE npcs SET location_id=NULL WHERE id=$1`,
				`UPDATE npcs SET properties=jsonb_set(properties,'{srd521_target,actor_id}','"another-actor"') WHERE id=$1`,
				`UPDATE npcs SET properties=jsonb_set(properties,'{srd521_target,hidden}','true') WHERE id=$1`,
				`UPDATE npcs SET properties=jsonb_set(properties,'{srd521_target,bandit,version}','"stale"') WHERE id=$1`,
				`UPDATE npcs SET properties=jsonb_set(properties,'{srd521_target,bandit,action_id}','"multiattack"') WHERE id=$1`,
				`UPDATE npcs SET properties=jsonb_set(properties,'{srd521_target,line_of_sight}','false') WHERE id=$1`,
			} {
				preview := prepare(combat.EncounterEndTurn, "")
				request := engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: preview.ExpectedRevision, Intent: preview.Intent}
				_, err = testPool.Exec(ctx, mutation, target)
				must(err)
				before := dice.calls
				if _, err = newEngine().ConfirmAction(ctx, campaign, request); !errors.Is(err, engine.ErrActionRefresh) {
					t.Fatalf("mutation accepted: %s: %v", mutation, err)
				}
				if dice.calls != before {
					t.Fatal("negative rolled")
				}
				_, err = testPool.Exec(ctx, `UPDATE npcs SET location_id=$1,properties=$2 WHERE id=$3`, location, facts, target)
				must(err)
			}
			hp := 12
			for _, tc := range []struct {
				outcome string
				rolls   []int
				damage  int
			}{
				{"hit", []int{18, 2}, 4}, {"critical-hit", []int{19, 1, 3}, 7}, {"miss", []int{0}, 0},
			} {
				t.Run(tc.outcome, func(t *testing.T) {
					preview := prepare(combat.EncounterEndTurn, "")
					request := engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: preview.ExpectedRevision, Intent: preview.Intent}
					before, err := game.ReadEncounter(ctx, testPool, campaign)
					must(err)
					dice.values = []int{-1}
					if _, err := newEngine().ConfirmAction(ctx, campaign, request); err == nil {
						t.Fatal("invalid dice committed")
					}
					afterDice, err := game.ReadEncounter(ctx, testPool, campaign)
					must(err)
					if !reflect.DeepEqual(before, afterDice) {
						t.Fatal("dice failure leaked turn state")
					}
					// Fail history after HP/budget/receipt mutation, then retry the same identity.
					_, err = testPool.Exec(ctx, `CREATE FUNCTION reject_bandit_history() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN IF NEW.player_input='Confirm combat-end-turn combat-end-turn' THEN RAISE EXCEPTION 'fixture history failure'; END IF; RETURN NEW; END; $$; CREATE TRIGGER reject_bandit_history BEFORE INSERT ON session_logs FOR EACH ROW EXECUTE FUNCTION reject_bandit_history()`)
					must(err)
					dice.values = append([]int(nil), tc.rolls...)
					_, failure := newEngine().ConfirmAction(ctx, campaign, request)
					_, err = testPool.Exec(ctx, `DROP TRIGGER reject_bandit_history ON session_logs; DROP FUNCTION reject_bandit_history()`)
					must(err)
					if failure == nil {
						t.Fatal("history failure ignored")
					}
					var storedHP, attempts int
					must(testPool.QueryRow(ctx, `SELECT hp FROM player_characters WHERE id=$1`, pc.ID).Scan(&storedHP))
					must(testPool.QueryRow(ctx, `SELECT count(*) FROM turn_attempts WHERE id=$1`, request.ActionID).Scan(&attempts))
					after, err := game.ReadEncounter(ctx, testPool, campaign)
					must(err)
					if storedHP != hp || attempts != 0 || !reflect.DeepEqual(before, after) {
						t.Fatal("failed history leaked HP, budget, cursor or receipt")
					}
					dice.values = append([]int(nil), tc.rolls...)
					result, err := newEngine().ConfirmAction(ctx, campaign, request)
					must(err)
					var payload struct {
						Receipts []srd521.AttackReceipt `json:"npc_receipts"`
					}
					must(json.Unmarshal(result.AppliedToolCalls[0].Result, &payload))
					if len(payload.Receipts) != 1 || payload.Receipts[0].Outcome != tc.outcome || payload.Receipts[0].Damage != tc.damage || payload.Receipts[0].Modifier.Value != 3 {
						t.Fatalf("receipt=%+v", payload)
					}
					if action == "scimitar" && payload.Receipts[0].DamageType != "slashing" || action == "light-crossbow" && payload.Receipts[0].DamageType != "piercing" {
						t.Fatal("wrong damage type")
					}
					hp -= tc.damage
					must(testPool.QueryRow(ctx, `SELECT hp FROM player_characters WHERE id=$1`, pc.ID).Scan(&storedHP))
					if storedHP != hp {
						t.Fatalf("HP=%d want %d", storedHP, hp)
					}
					calls := dice.calls
					replay, err := newEngine().ConfirmAction(ctx, campaign, request)
					must(err)
					if dice.calls != calls || !sameTurnResult(result, replay) {
						t.Fatal("replay rerolled")
					}
					saved, err := newEngine().GetActionReceipt(ctx, campaign, request.ActionID)
					must(err)
					if !sameTurnResult(result, saved.Result) {
						t.Fatal("receipt reload changed")
					}
					reloaded, err := game.ReadEncounter(ctx, testPool, campaign)
					must(err)
					if reloaded.Round != before.Round+1 || reloaded.Active().ID != actor || !reloaded.Budget.Action || reloaded.Budget.Movement != 30 {
						t.Fatalf("budget=%+v", reloaded)
					}
				})
			}
			// At melee distance a total of 11 misses the Bandit's pinned AC 12.
			if action == "scimitar" {
				dice.values = []int{5}
				_, result := confirm(prepare(srd521.MartialAttack, weapon.String()), target)
				var p struct {
					Receipt srd521.AttackReceipt `json:"receipt"`
				}
				must(json.Unmarshal(result.AppliedToolCalls[0].Result, &p))
				if p.Receipt.Total != 11 || p.Receipt.Outcome != "miss" {
					t.Fatalf("AC was not 12: %+v", p)
				}
				if prepare(srd521.MartialAttack, weapon.String()).Available {
					t.Fatal("double attack")
				}
			}
		})
	}
}

// newBanditScene reuses the Beast fixture's scene with a pinned Bandit.
func newBanditScene(t *testing.T, action string, distance int) *beastFixture {
	f := newBeastFixture(t, "wolf-v1", "bite", "")
	_, version, err := srd521.BanditRecord()
	f.must(err)
	f.facts = fmt.Sprintf(`{"srd521_target":{"profile":"bandit-v1","actor_id":"%s","distance_feet":%d,"line_of_sight":true,"cover":"none","bandit":{"version":"%s","action_id":"%s"}}}`, f.actor, distance, version, action)
	f.reset()
	return f
}

// newBanditFixture also starts the encounter; the player acts first.
func newBanditFixture(t *testing.T, action string, distance int) *beastFixture {
	f := newBanditScene(t, action, distance)
	f.confirm(f.request(f.prepare(combat.EncounterStart, f.weapon.String()), f.target), 18, 0)
	return f
}

// endBanditTurn previews and confirms End turn; want is in the preview.
func (f *beastFixture) endBanditTurn(want string, rolls ...int) []srd521.AttackReceipt {
	f.t.Helper()
	end := f.prepare(combat.EncounterEndTurn, "")
	if !end.Available || !strings.Contains(strings.Join(end.Options, " "), want) {
		f.t.Fatalf("end=%+v", end)
	}
	return npcReceipts(f.t, f.confirm(f.request(end, uuid.Nil), rolls...))
}

func TestIntegrationBanditOutOfScimitarReach(t *testing.T) {
	f := newBanditFixture(t, "scimitar", 5)
	for range 2 {
		f.confirm(f.request(f.prepare(combat.EncounterRetreat, ""), uuid.Nil))
	}
	// At 15 feet the Bandit shoots its Light Crossbow from where it stands.
	r := f.endBanditTurn("Light Crossbow", 18, 2)
	e := f.encounter()
	if len(r) != 1 || r[0].DamageType != "piercing" || r[0].Damage != 4 || *e.Space.DistanceFeet != 15 || !e.Active().Player || e.Round != 2 || f.hp() != 8 {
		t.Fatalf("receipts=%+v encounter=%+v", r, e)
	}
}

func TestIntegrationBanditCrossbowWithinFiveFeet(t *testing.T) {
	f := newBanditFixture(t, "light-crossbow", 30)
	for range 5 {
		f.confirm(f.request(f.prepare(combat.EncounterApproach, ""), uuid.Nil))
	}
	// Within 5 feet the Bandit uses its Scimitar instead.
	r := f.endBanditTurn("Scimitar", 18, 2)
	e := f.encounter()
	if len(r) != 1 || r[0].DamageType != "slashing" || r[0].Damage != 4 || *e.Space.DistanceFeet != 5 || !e.Active().Player || e.Round != 2 {
		t.Fatalf("receipts=%+v encounter=%+v", r, e)
	}
}

func TestIntegrationBanditApproachesBeyondNormalRange(t *testing.T) {
	f := newBanditFixture(t, "light-crossbow", 80)
	for range 6 {
		f.confirm(f.request(f.prepare(combat.EncounterRetreat, ""), uuid.Nil))
	}
	// From 110 feet it walks 30 feet to 80, then shoots; a d20 of 1 misses.
	r := f.endBanditTurn("moves 30 feet", 0)
	e := f.encounter()
	if len(r) != 1 || r[0].Outcome != "miss" || *e.Space.DistanceFeet != 80 || !e.Active().Player || e.Round != 2 {
		t.Fatalf("receipts=%+v encounter=%+v", r, e)
	}
	// Dash and retreat to 140 feet: after 30 feet it is still at 110 and does
	// not attack, so nothing is rolled.
	f.confirm(f.request(f.prepare(combat.EncounterDash, ""), uuid.Nil))
	for range 12 {
		f.confirm(f.request(f.prepare(combat.EncounterRetreat, ""), uuid.Nil))
	}
	calls := f.dice.calls
	r = f.endBanditTurn("does not attack")
	e = f.encounter()
	if len(r) != 0 || f.dice.calls != calls || *e.Space.DistanceFeet != 110 || !e.Active().Player || e.Round != 3 || f.hp() != 12 {
		t.Fatalf("receipts=%+v encounter=%+v", r, e)
	}
}

// The turn picks the weapon by distance, so the configured weapon need not be
// legal for the encounter to start.
func TestIntegrationBanditStartsBeyondScimitarReach(t *testing.T) {
	f := newBanditFixture(t, "scimitar", 30)
	if e := f.encounter(); e.Bandit.ActionID != "scimitar" || *e.Space.DistanceFeet != 30 || !e.Active().Player {
		t.Fatalf("encounter=%+v", e)
	}
	r := f.endBanditTurn("Light Crossbow", 18, 2)
	if len(r) != 1 || r[0].DamageType != "piercing" || f.hp() != 8 {
		t.Fatalf("receipts=%+v hp=%d", r, f.hp())
	}
}

func TestIntegrationBanditStartRefusesUnknownFacts(t *testing.T) {
	f := newBanditScene(t, "scimitar", 30)
	for _, mutation := range []string{
		`UPDATE npcs SET properties=properties #- '{srd521_target,distance_feet}' WHERE id=$1`,
		`UPDATE npcs SET properties=properties #- '{srd521_target,line_of_sight}' WHERE id=$1`,
		`UPDATE npcs SET properties=properties #- '{srd521_target,cover}' WHERE id=$1`,
		`UPDATE npcs SET properties=jsonb_set(properties,'{srd521_target,bandit,version}','"stale"') WHERE id=$1`,
		`UPDATE npcs SET properties=jsonb_set(properties,'{srd521_target,bandit,action_id}','"multiattack"') WHERE id=$1`,
	} {
		f.reset()
		start := f.prepare(combat.EncounterStart, f.weapon.String())
		request := f.request(start, f.target)
		f.setFacts(mutation)
		if p := f.prepare(combat.EncounterStart, f.weapon.String()); p.Available {
			t.Fatalf("start offered: %s: %+v", mutation, p)
		}
		calls := f.dice.calls
		if _, err := f.engine().ConfirmAction(f.ctx, f.campaign, request); !errors.Is(err, engine.ErrActionRefresh) || f.dice.calls != calls {
			t.Fatalf("start accepted or rolled: %s: %v", mutation, err)
		}
	}
	// A declared fact that rules out every attack is known, so the encounter
	// starts and the Bandit's turn passes without a roll.
	f.reset()
	f.setFacts(`UPDATE npcs SET properties=jsonb_set(properties,'{srd521_target,cover}','"total"') WHERE id=$1`)
	f.confirm(f.request(f.prepare(combat.EncounterStart, f.weapon.String()), f.target), 18, 0)
	calls := f.dice.calls
	if r := f.endBanditTurn("does not attack"); len(r) != 0 || f.dice.calls != calls || f.hp() != 12 {
		t.Fatalf("receipts=%+v hp=%d", r, f.hp())
	}
}

// Beyond every weapon's range the start preview says the Bandit will close in,
// rather than listing each weapon as unavailable.
func TestIntegrationBanditStartPreviewBeyondEveryRange(t *testing.T) {
	f := newBanditScene(t, "scimitar", 100)
	text := strings.Join(f.prepare(combat.EncounterStart, f.weapon.String()).Options, "\n")
	if !strings.Contains(text, "100 feet away, beyond its 80-foot range; it moves closer on its turn, up to its 30-foot speed") || strings.Contains(text, "Unavailable:") {
		t.Fatalf("options=%s", text)
	}
}
