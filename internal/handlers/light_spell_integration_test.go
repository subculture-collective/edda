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
	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"github.com/google/uuid"
)

func TestIntegrationLightOnCarriedObject(t *testing.T) {
	ctx := context.Background()
	q := statedb.New(testPool)
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	c, err := q.CreateCampaign(ctx, statedb.CreateCampaignParams{Name: "Light fixture", RulesVersion: "srd-5.2.1", Status: "active", CreatedBy: dbutil.ToPgtype(testUserID)})
	must(err)
	campaign := dbutil.FromPgtype(c.ID)
	t.Cleanup(func() { _, _ = testPool.Exec(ctx, `DELETE FROM campaigns WHERE id=$1`, campaign) })
	build, err := srd521.SupportedBuildForConcept("Human Cleric", "Acolyte")
	must(err)
	stats, err := build.Stats()
	must(err)
	created, err := q.CreatePlayerCharacter(ctx, statedb.CreatePlayerCharacterParams{CampaignID: c.ID, UserID: dbutil.ToPgtype(testUserID), Name: "Mira", Stats: stats, Hp: 10, MaxHp: 10, Level: 1, Status: "active", Abilities: []byte(`[]`)})
	must(err)
	actor := dbutil.FromPgtype(created.ID)
	other, err := q.CreatePlayerCharacter(ctx, statedb.CreatePlayerCharacterParams{CampaignID: c.ID, UserID: dbutil.ToPgtype(testUserID), Name: "Iona", Stats: stats, Hp: 10, MaxHp: 10, Level: 1, Status: "active", Abilities: []byte(`[]`)})
	must(err)
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET stats=stats||'{"srd521_senses":{"darkvision_feet":0},"custom":{"keep":7}}'::jsonb WHERE id=$1`, actor)
	must(err)
	var location uuid.UUID
	must(testPool.QueryRow(ctx, `INSERT INTO locations(campaign_id,name)VALUES($1,'Crypt')RETURNING id`, campaign).Scan(&location))
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET current_location_id=$1 WHERE campaign_id=$2`, location, campaign)
	must(err)
	item := func(owner any, name, kind string, quantity int) uuid.UUID {
		t.Helper()
		var id uuid.UUID
		must(testPool.QueryRow(ctx, `INSERT INTO items(campaign_id,player_character_id,name,item_type,rarity,quantity)VALUES($1,$2,$3,$4,'common',$5)RETURNING id`, campaign, owner, name, kind, quantity).Scan(&id))
		return id
	}
	torch := item(actor, "Torch", "misc", 1)
	mace := item(actor, "Mace", "weapon", 1)
	item(actor, "Candles", "misc", 5)
	item(actor, "Potion of Healing", "consumable", 1)
	held := item(dbutil.FromPgtype(other.ID), "Iona's lantern", "misc", 1)
	unattended := item(nil, "Loose stone", "misc", 1)
	var target uuid.UUID
	must(testPool.QueryRow(ctx, `INSERT INTO npcs(campaign_id,name,hp,location_id)VALUES($1,'Lurker',9,$2)RETURNING id`, campaign, location).Scan(&target))
	_, err = testPool.Exec(ctx, `INSERT INTO session_logs(campaign_id,turn_number,player_input,input_type,llm_response,npcs_involved)VALUES($1,1,'Arrive','game_action','Something waits.',ARRAY[$2]::uuid[])`, campaign, target)
	must(err)
	facts := func(distance int) {
		t.Helper()
		props := fmt.Sprintf(`{"srd521_target":{"profile":"ordinary-spatial-fixture-v1","actor_id":"%s","distance_feet":%d,"line_of_sight":true,"cover":"none","light":"dark","obscurement":"none"}}`, actor, distance)
		_, err := testPool.Exec(ctx, `UPDATE npcs SET properties=$1 WHERE id=$2`, props, target)
		must(err)
	}
	facts(15)
	dice := &lifecycleIntegrationDice{}
	newEngine := func() *engine.Engine {
		e, err := engine.New(testPool, nil, config.LLMConfig{}, engine.WithActionDice(dice))
		must(err)
		return e
	}
	prepare := func(intent game.ActionIntent) *game.ActionPreview {
		t.Helper()
		intent.ActorID = actor
		c, pc, err := game.ReadActionSnapshot(ctx, testPool, campaign, actor, testUserID, false)
		must(err)
		a, err := game.PrepareStoredAction(ctx, testPool, c, pc, intent, false)
		must(err)
		return a.Preview
	}
	cast := func(color string) *game.ActionPreview {
		t.Helper()
		return prepare(game.ActionIntent{Kind: srd521.SpellCast, EntryID: game.LightEntryPrefix + color})
	}
	sight := func() *game.ActionPreview {
		t.Helper()
		return prepare(game.ActionIntent{Kind: srd521.AbilityCheck, EntryID: game.SightPerceptionEntryPrefix + target.String()})
	}
	confirm := func(p *game.ActionPreview, targetID uuid.UUID) (engine.PreparedAction, *engine.TurnResult, error) {
		t.Helper()
		r := engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: p.ExpectedRevision, Intent: p.Intent, TargetID: targetID}
		result, err := newEngine().ConfirmAction(ctx, campaign, r)
		return r, result, err
	}
	saved := func() (*srd521.LightEffect, *srd521.Concentration, int, int) {
		t.Helper()
		var raw []byte
		must(testPool.QueryRow(ctx, `SELECT stats FROM player_characters WHERE id=$1`, actor).Scan(&raw))
		var envelope struct {
			Build  srd521.Build `json:"srd521"`
			Custom struct {
				Keep int `json:"keep"`
			} `json:"custom"`
		}
		must(json.Unmarshal(raw, &envelope))
		return envelope.Build.Casting.Light, envelope.Build.Casting.Concentration, envelope.Build.Casting.Slots, envelope.Custom.Keep
	}
	receipt := func(result *engine.TurnResult) srd521.LightReceipt {
		t.Helper()
		var payload struct {
			Receipt srd521.LightReceipt `json:"receipt"`
		}
		if len(result.AppliedToolCalls) != 1 || result.AppliedToolCalls[0].Tool != "srd_light" {
			t.Fatalf("calls: %+v", result.AppliedToolCalls)
		}
		must(json.Unmarshal(result.AppliedToolCalls[0].Result, &payload))
		return payload.Receipt
	}
	unavailable := func(p *game.ActionPreview, reason string) {
		t.Helper()
		if p.Available || !strings.Contains(p.Reason, reason) {
			t.Fatalf("want unavailable %q: %+v", reason, p)
		}
	}
	attempts := func() int {
		t.Helper()
		var n int
		must(testPool.QueryRow(ctx, `SELECT count(*) FROM turn_attempts WHERE campaign_id=$1`, campaign).Scan(&n))
		return n
	}

	// Darkness hides the target before any Light.
	unavailable(sight(), "can't see")
	unavailable(prepare(game.ActionIntent{Kind: srd521.SpellCast, EntryID: srd521.LightSpell}), "explicit Light color")
	unavailable(cast("ultraviolet"), "explicit Light color")
	unavailable(prepare(game.ActionIntent{Kind: srd521.SpellCast, EntryID: game.LightEntryPrefix + "blue", SlotLevel: 1}), "Choose Light")
	unavailable(prepare(game.ActionIntent{Kind: srd521.DismissLight}), "No active Light")

	// Shield of Faith concentration must survive every Light change.
	faith := prepare(game.ActionIntent{Kind: srd521.SpellCast, EntryID: "shield-of-faith", SlotLevel: 1})
	_, _, err = confirm(faith, actor)
	must(err)
	_, concentration, slots, _ := saved()
	if concentration == nil {
		t.Fatal("Shield of Faith not saved")
	}

	// Preview lists only objects the caster carries and is inert.
	beforeLight, _, _, _ := saved()
	count := attempts()
	p := cast("blue")
	if !p.Available || len(p.Targets) != 2 || len(p.Costs) != 1 || p.Costs[0].ResourceID != "magic-action" {
		t.Fatalf("Light preview: %+v", p)
	}
	for _, option := range p.Targets {
		if option.ID != torch && option.ID != mace {
			t.Fatalf("ineligible object offered: %+v", option)
		}
	}
	if after, _, _, _ := saved(); after != beforeLight || attempts() != count {
		t.Fatal("preview changed state")
	}
	// Objects held by someone else, unattended objects and the caster are refused.
	for _, id := range []uuid.UUID{held, unattended, actor, uuid.Nil} {
		if _, _, err := confirm(p, id); !errors.Is(err, engine.ErrActionRefresh) {
			t.Fatalf("target %s accepted: %v", id, err)
		}
	}
	if attempts() != count || dice.calls != 0 {
		t.Fatal("refused Light changed state or rolled")
	}

	// Atomic: a failed history write leaves no Light and no attempt.
	_, err = testPool.Exec(ctx, `ALTER TABLE session_logs ADD CONSTRAINT light_history_failure CHECK (player_input NOT LIKE 'Confirm Light%') NOT VALID`)
	must(err)
	t.Cleanup(func() {
		_, _ = testPool.Exec(ctx, `ALTER TABLE session_logs DROP CONSTRAINT IF EXISTS light_history_failure`)
	})
	if _, _, err := confirm(p, torch); err == nil {
		t.Fatal("history failure committed")
	}
	if light, _, _, _ := saved(); light != nil || attempts() != count {
		t.Fatal("partial Light commit")
	}
	_, err = testPool.Exec(ctx, `ALTER TABLE session_logs DROP CONSTRAINT light_history_failure`)
	must(err)

	now, err := game.ReadSpellClock(ctx, testPool, campaign)
	must(err)
	request, result, err := confirm(p, torch)
	must(err)
	first := receipt(result)
	want := srd521.LightEffect{ID: request.ActionID.String(), ItemID: torch.String(), Color: "blue", StartedAt: now, ExpiresAt: now + 3600, BrightFeet: 20, DimFeet: 20}
	light, kept, slotsAfter, custom := saved()
	if first.Before != nil || first.After == nil || *first.After != want || light == nil || *light != want || *kept != *concentration || slotsAfter != slots || custom != 7 {
		t.Fatalf("Light cast: %+v saved=%+v", first, light)
	}
	// Replay reads the saved result without casting again.
	replayed, err := newEngine().ConfirmAction(ctx, campaign, request)
	must(err)
	if replayed.Narrative != result.Narrative || attempts() != count+1 {
		t.Fatal("replay changed Light")
	}
	// History reconstruction exposes the typed receipt.
	var history string
	must(testPool.QueryRow(ctx, `SELECT tool_calls::text FROM session_logs WHERE campaign_id=$1 ORDER BY turn_number DESC LIMIT 1`, campaign).Scan(&history))
	if !strings.Contains(history, `"srd_light"`) {
		t.Fatal("Light receipt missing from history")
	}

	// The carried Light raises darkness at the target through #145 sight.
	lit := sight()
	if !lit.Available || strings.Contains(lit.Options[0], "Disadvantage") || !strings.Contains(lit.Options[0], "lit by your Light on Torch") {
		t.Fatalf("bright radius: %+v", lit)
	}
	facts(40)
	if p := sight(); !p.Available || !strings.Contains(p.Options[0], "Disadvantage") {
		t.Fatalf("dim radius edge: %+v", p)
	}
	facts(45)
	unavailable(sight(), "can't see")
	facts(15)
	// Handing the object to someone else ends its effect on sight and
	// invalidates an open preview.
	stale := sight()
	_, err = testPool.Exec(ctx, `UPDATE items SET player_character_id=$1 WHERE id=$2`, dbutil.FromPgtype(other.ID), torch)
	must(err)
	if _, _, err := confirm(stale, uuid.Nil); !errors.Is(err, engine.ErrActionRefresh) {
		t.Fatalf("stale lit preview confirmed: %v", err)
	}
	unavailable(sight(), "can't see")
	_, err = testPool.Exec(ctx, `UPDATE items SET player_character_id=$1 WHERE id=$2`, actor, torch)
	must(err)

	// Recast removes only the previous Light; concentration is unchanged.
	recast := cast("red")
	if !recast.Available || !strings.Contains(recast.Options[0], "Ends your current Light on Torch") {
		t.Fatalf("recast preview: %+v", recast)
	}
	request, result, err = confirm(recast, mace)
	must(err)
	second := receipt(result)
	light, kept, _, _ = saved()
	if second.Before == nil || *second.Before != want || light.ID != request.ActionID.String() || light.ItemID != mace.String() || light.Color != "red" || *kept != *concentration {
		t.Fatalf("recast: %+v saved=%+v", second, light)
	}
	expires := light.ExpiresAt

	// Time advancement keeps the exact absolute expiry.
	_, _, err = confirm(prepare(game.ActionIntent{Kind: srd521.WaitMinute}), uuid.Nil)
	must(err)
	if light, _, _, _ = saved(); light == nil || light.ExpiresAt != expires {
		t.Fatalf("wait changed expiry: %+v", light)
	}
	if p := sight(); !p.Available {
		t.Fatalf("Light ended early: %+v", p)
	}
	_, _, err = confirm(prepare(game.ActionIntent{Kind: srd521.ShortRest}), uuid.Nil)
	must(err)
	now, err = game.ReadSpellClock(ctx, testPool, campaign)
	must(err)
	if light, _, _, _ = saved(); light == nil || light.ExpiresAt != expires || now < expires {
		t.Fatalf("rest changed expiry: %+v now=%d", light, now)
	}
	unavailable(sight(), "can't see")
	unavailable(prepare(game.ActionIntent{Kind: srd521.DismissLight}), "No active Light")

	// Voluntary dismissal ends only the Light.
	request, result, err = confirm(cast("white"), torch)
	must(err)
	if r := receipt(result); r.Before != nil {
		t.Fatal("expired Light reported as ended by recast")
	}
	_, _, err = confirm(prepare(game.ActionIntent{Kind: srd521.SpellCast, EntryID: "shield-of-faith", SlotLevel: 1}), actor)
	must(err)
	_, kept, _, _ = saved()
	if kept == nil {
		t.Fatal("second Shield of Faith not saved")
	}
	dismiss := prepare(game.ActionIntent{Kind: srd521.DismissLight})
	if !dismiss.Available || len(dismiss.Costs) != 0 || !strings.Contains(dismiss.Options[0], "Torch") {
		t.Fatalf("dismiss preview: %+v", dismiss)
	}
	_, result, err = confirm(dismiss, uuid.Nil)
	must(err)
	ended := receipt(result)
	_, keptAfter, _, _ := saved()
	if light, _, _, _ = saved(); light != nil || ended.Kind != srd521.DismissLight || ended.Before.ID != request.ActionID.String() {
		t.Fatalf("dismissal: %+v saved=%+v", ended, light)
	}
	if keptAfter == nil || *keptAfter != *kept {
		t.Fatal("dismissal changed concentration")
	}
	unavailable(prepare(game.ActionIntent{Kind: srd521.DismissLight}), "No active Light")

	// Light is outside encounters only.
	facts(5)
	_, err = testPool.Exec(ctx, `UPDATE npcs SET properties=jsonb_set(properties,'{srd521_target,light}','"bright"') WHERE id=$1`, target)
	must(err)
	_, _, err = confirm(prepare(game.ActionIntent{Kind: combat.EncounterStart}), target)
	must(err)
	unavailable(cast("blue"), "outside encounters")
}
