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

func TestIntegrationTravelSegment(t *testing.T) {
	ctx := context.Background()
	q := statedb.New(testPool)
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	c, err := q.CreateCampaign(ctx, statedb.CreateCampaignParams{Name: "Travel fixture", RulesVersion: "srd-5.2.1", Status: "active", CreatedBy: dbutil.ToPgtype(testUserID)})
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
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET stats=stats||'{"custom":{"keep":7}}'::jsonb WHERE id=$1`, actor)
	must(err)
	location := func(name string, known bool) uuid.UUID {
		t.Helper()
		var id uuid.UUID
		must(testPool.QueryRow(ctx, `INSERT INTO locations(campaign_id,name,player_known)VALUES($1,$2,$3)RETURNING id`, campaign, name, known).Scan(&id))
		return id
	}
	gate, wood, peak, prose, secret := location("Gate", true), location("Wood", true), location("Peak", true), location("Old Road", true), location("Hidden Vale", false)
	connect := func(from, to uuid.UUID, both bool, travelTime string, facts any) uuid.UUID {
		t.Helper()
		var id uuid.UUID
		must(testPool.QueryRow(ctx, `INSERT INTO location_connections(from_location_id,to_location_id,bidirectional,travel_time,campaign_id,srd521_route)VALUES($1,$2,$3,NULLIF($4,''),$5,$6)RETURNING id`, from, to, both, travelTime, campaign, facts).Scan(&id))
		return id
	}
	woodRoute := connect(gate, wood, true, "", `{"profile":"travel-route-v1","distance_feet":5280,"terrain":"forest","good_road":false}`)
	connect(gate, peak, false, "", `{"profile":"travel-route-v1","distance_feet":10560,"terrain":"mountain","good_road":false}`)
	connect(gate, prose, false, "2 hours", nil)
	connect(gate, secret, false, "", `{"profile":"travel-route-v1","distance_feet":5280,"terrain":"grassland","good_road":true}`)
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET current_location_id=$1 WHERE id=$2`, gate, actor)
	must(err)
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
	travel := func(pace string) *game.ActionPreview {
		t.Helper()
		return prepare(game.ActionIntent{Kind: srd521.TravelHour, EntryID: pace})
	}
	confirm := func(p *game.ActionPreview, target uuid.UUID) (engine.PreparedAction, *engine.TurnResult, error) {
		t.Helper()
		r := engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: p.ExpectedRevision, Intent: p.Intent, TargetID: target}
		result, err := newEngine().ConfirmAction(ctx, campaign, r)
		return r, result, err
	}
	type state struct {
		Location uuid.UUID
		Clock    int
		Attempts int
		Travel   string
		Light    *srd521.LightEffect
		Faith    *srd521.Concentration
		Keep     int
	}
	snapshot := func() state {
		t.Helper()
		var s state
		var raw []byte
		var travelled *string
		must(testPool.QueryRow(ctx, `SELECT current_location_id,stats,stats->>'srd521_travel' FROM player_characters WHERE id=$1`, actor).Scan(&s.Location, &raw, &travelled))
		if travelled != nil {
			s.Travel = *travelled
		}
		var envelope struct {
			Build  srd521.Build `json:"srd521"`
			Custom struct {
				Keep int `json:"keep"`
			} `json:"custom"`
		}
		must(json.Unmarshal(raw, &envelope))
		s.Light, s.Faith, s.Keep = envelope.Build.Casting.Light, envelope.Build.Casting.Concentration, envelope.Custom.Keep
		s.Clock, err = game.ReadSpellClock(ctx, testPool, campaign)
		must(err)
		must(testPool.QueryRow(ctx, `SELECT count(*) FROM turn_attempts WHERE campaign_id=$1`, campaign).Scan(&s.Attempts))
		return s
	}
	same := func(a, b state) bool {
		return a.Location == b.Location && a.Clock == b.Clock && a.Attempts == b.Attempts && a.Travel == b.Travel
	}
	target := func(p *game.ActionPreview, id uuid.UUID) *game.ActionTarget {
		for i := range p.Targets {
			if p.Targets[i].ID == id {
				return &p.Targets[i]
			}
		}
		return nil
	}
	unavailable := func(p *game.ActionPreview, reason string) {
		t.Helper()
		if p.Available || !strings.Contains(p.Reason, reason) {
			t.Fatalf("want unavailable %q: %+v", reason, p)
		}
	}

	// Normal pace: only the declared forest route fits; the hidden destination
	// is never listed and travel_time prose is not a distance.
	before := snapshot()
	normal := travel(srd521.PaceNormal)
	if !normal.Available || len(normal.Targets) != 3 || target(normal, secret) != nil || len(normal.Costs) != 0 {
		t.Fatalf("normal preview: %+v", normal)
	}
	if w := target(normal, wood); w == nil || w.UnavailableReason != "" || !strings.Contains(w.Label, "1.00 miles (5280 feet), forest, no road; maximum pace normal; 20 min 00 s at a normal pace") {
		t.Fatalf("wood option: %+v", w)
	}
	if p := target(normal, peak); p == nil || !strings.Contains(p.UnavailableReason, "mountain terrain allows at most a slow pace") {
		t.Fatalf("peak option: %+v", p)
	}
	if p := target(normal, prose); p == nil || !strings.Contains(p.UnavailableReason, "travel time prose is not an SRD distance") || strings.Contains(p.Label, "hours") {
		t.Fatalf("prose option: %+v", p)
	}
	if !strings.Contains(normal.Options[0], "Disadvantage on Dexterity (Stealth) checks") {
		t.Fatalf("pace effect missing: %+v", normal.Options)
	}
	if !same(before, snapshot()) {
		t.Fatal("preview changed state")
	}
	unavailable(travel(srd521.PaceFast), "No known route can be travelled within one hour at a fast pace")
	unavailable(travel("forced-march"), "Choose a fast, normal or slow pace")
	unavailable(prepare(game.ActionIntent{Kind: srd521.TravelHour, EntryID: srd521.PaceNormal, Ability: "wisdom"}), "Choose a fast")
	slow := travel(srd521.PaceSlow)
	if !slow.Available || target(slow, peak).UnavailableReason != "" || !strings.Contains(target(slow, wood).Label, "30 min 00 s") || !strings.Contains(slow.Options[0], "Advantage on Wisdom (Perception or Survival) checks") {
		t.Fatalf("slow preview: %+v", slow)
	}

	// Unoffered destinations, unavailable routes and a missing target are refused.
	for _, id := range []uuid.UUID{secret, prose, peak, gate, uuid.Nil} {
		if _, _, err := confirm(normal, id); !errors.Is(err, engine.ErrActionRefresh) {
			t.Fatalf("destination %s accepted: %v", id, err)
		}
	}
	// A preview's revision binds its pace.
	switched := engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: normal.ExpectedRevision, Intent: slow.Intent, TargetID: peak}
	if _, err := newEngine().ConfirmAction(ctx, campaign, switched); !errors.Is(err, engine.ErrActionRefresh) {
		t.Fatalf("changed pace confirmed: %v", err)
	}
	// Changed route facts invalidate the open preview.
	_, err = testPool.Exec(ctx, `UPDATE location_connections SET srd521_route=jsonb_set(srd521_route,'{distance_feet}','15841') WHERE id=$1`, woodRoute)
	must(err)
	if _, _, err := confirm(normal, wood); !errors.Is(err, engine.ErrActionRefresh) {
		t.Fatalf("stale route confirmed: %v", err)
	}
	if w := target(travel(srd521.PaceNormal), wood); w == nil || !strings.Contains(w.UnavailableReason, "longer than 3 miles") {
		t.Fatalf("pace boundary: %+v", w)
	}
	_, err = testPool.Exec(ctx, `UPDATE location_connections SET srd521_route=jsonb_set(srd521_route,'{distance_feet}','15840') WHERE id=$1`, woodRoute)
	must(err)
	if w := target(travel(srd521.PaceNormal), wood); w == nil || w.UnavailableReason != "" || !strings.Contains(w.Label, "60 min 00 s") {
		t.Fatalf("pace boundary: %+v", w)
	}
	// An undeclared mechanic in the facts fails closed.
	_, err = testPool.Exec(ctx, `UPDATE location_connections SET srd521_route=jsonb_set(srd521_route,'{mounted}','true') WHERE id=$1`, woodRoute)
	must(err)
	if w := target(travel(srd521.PaceNormal), wood); w == nil || !strings.Contains(w.UnavailableReason, "unknown or malformed") {
		t.Fatalf("unknown route fact: %+v", w)
	}
	_, err = testPool.Exec(ctx, `UPDATE location_connections SET srd521_route=jsonb_set(srd521_route-'mounted','{distance_feet}','5280') WHERE id=$1`, woodRoute)
	must(err)
	if !same(before, snapshot()) || dice.calls != 0 {
		t.Fatal("refused travel changed state or rolled")
	}

	// Light (one hour) and Shield of Faith (ten minutes) are running.
	_, _, err = confirm(prepare(game.ActionIntent{Kind: srd521.SpellCast, EntryID: "shield-of-faith", SlotLevel: 1}), actor)
	must(err)
	var torch uuid.UUID
	must(testPool.QueryRow(ctx, `INSERT INTO items(campaign_id,player_character_id,name,item_type,rarity,quantity)VALUES($1,$2,'Torch','misc','common',1)RETURNING id`, campaign, actor).Scan(&torch))
	_, _, err = confirm(prepare(game.ActionIntent{Kind: srd521.SpellCast, EntryID: game.LightEntryPrefix + "white"}), torch)
	must(err)
	before = snapshot()
	if before.Light == nil || before.Faith == nil {
		t.Fatalf("effects not running: %+v", before)
	}

	// Atomic: a failed history write moves nothing and advances no time.
	normal = travel(srd521.PaceNormal)
	_, err = testPool.Exec(ctx, `ALTER TABLE session_logs ADD CONSTRAINT travel_history_failure CHECK (player_input NOT LIKE 'Confirm Travel%') NOT VALID`)
	must(err)
	t.Cleanup(func() {
		_, _ = testPool.Exec(ctx, `ALTER TABLE session_logs DROP CONSTRAINT IF EXISTS travel_history_failure`)
	})
	if _, _, err := confirm(normal, wood); err == nil {
		t.Fatal("history failure committed")
	}
	if after := snapshot(); !same(before, after) || after.Faith == nil {
		t.Fatalf("partial travel commit: %+v", after)
	}
	_, err = testPool.Exec(ctx, `ALTER TABLE session_logs DROP CONSTRAINT travel_history_failure`)
	must(err)

	request, result, err := confirm(normal, wood)
	must(err)
	if len(result.AppliedToolCalls) != 1 || result.AppliedToolCalls[0].Tool != "srd_travel" {
		t.Fatalf("calls: %+v", result.AppliedToolCalls)
	}
	var payload struct {
		Receipt srd521.TravelReceipt `json:"receipt"`
	}
	must(json.Unmarshal(result.AppliedToolCalls[0].Result, &payload))
	r := payload.Receipt
	after := snapshot()
	if r.Kind != srd521.TravelHour || r.ToID != wood.String() || r.FromID != gate.String() || r.RouteID != woodRoute.String() || r.Plan.Seconds != 1200 || r.TimeBefore != before.Clock || r.TimeAfter != before.Clock+1200 || r.DayTravelBefore != 0 || r.DayTravelAfter != 1200 {
		t.Fatalf("receipt: %+v", r)
	}
	if after.Location != wood || after.Clock != before.Clock+1200 || after.Attempts != before.Attempts+1 || after.Keep != 7 {
		t.Fatalf("travel state: %+v", after)
	}
	var travelled srd521.TravelDay
	must(json.Unmarshal([]byte(after.Travel), &travelled))
	if travelled != (srd521.TravelDay{Day: before.Clock/86400 + 1, Seconds: 1200}) {
		t.Fatalf("travel record: %+v", travelled)
	}
	// Light keeps its exact absolute expiry; Shield of Faith ended on the way.
	if after.Light == nil || *after.Light != *before.Light || !after.Light.Active(after.Clock) || after.Faith != nil || r.Expired == nil || *r.Expired != *before.Faith || !strings.Contains(result.Narrative, "Shield of Faith expired") {
		t.Fatalf("durable effects: before=%+v after=%+v receipt=%+v", before, after, r)
	}
	var visited bool
	must(testPool.QueryRow(ctx, `SELECT player_visited FROM locations WHERE id=$1`, wood).Scan(&visited))
	if !visited {
		t.Fatal("destination not marked visited")
	}
	var history string
	must(testPool.QueryRow(ctx, `SELECT tool_calls::text FROM session_logs WHERE campaign_id=$1 ORDER BY turn_number DESC LIMIT 1`, campaign).Scan(&history))
	if !strings.Contains(history, `"srd_travel"`) {
		t.Fatal("travel receipt missing from history")
	}

	// Same-ID replay returns the saved result and moves nothing again.
	replayed, err := newEngine().ConfirmAction(ctx, campaign, request)
	must(err)
	if replayed.Narrative != result.Narrative || !same(after, snapshot()) {
		t.Fatal("replay changed travel")
	}
	changed := request
	changed.Intent.EntryID = srd521.PaceSlow
	if _, err := newEngine().ConfirmAction(ctx, campaign, changed); !errors.Is(err, engine.ErrActionIDConflict) {
		t.Fatalf("reused action ID accepted: %v", err)
	}

	// Reload: a fresh engine and snapshot offer the way back on the
	// bidirectional route, and the old preview is stale.
	if _, _, err := confirm(normal, wood); !errors.Is(err, engine.ErrActionRefresh) {
		t.Fatalf("pre-travel preview confirmed after arrival: %v", err)
	}
	back := travel(srd521.PaceNormal)
	if !back.Available || len(back.Targets) != 1 || back.Targets[0].ID != gate {
		t.Fatalf("return preview: %+v", back)
	}

	// A saved Prone condition or a pending declared fall blocks travel.
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET stats=stats||'{"srd521_conditions":{"prone":true}}'::jsonb WHERE id=$1`, actor)
	must(err)
	unavailable(travel(srd521.PaceNormal), "Stand up from Prone")
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET stats=(stats-'srd521_conditions')||'{"srd521_fall":{"id":"f"}}'::jsonb WHERE id=$1`, actor)
	must(err)
	unavailable(travel(srd521.PaceNormal), "Resolve the declared fall")
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET stats=stats-'srd521_fall' WHERE id=$1`, actor)
	must(err)

	// No Extended Travel: a segment that would pass 8 hours today is refused.
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET stats=jsonb_set(stats,'{srd521_travel,seconds}',to_jsonb($1::int)) WHERE id=$2`, srd521.TravelDaySeconds-1199, actor)
	must(err)
	if p := travel(srd521.PaceNormal); p.Available || !strings.Contains(p.Targets[0].UnavailableReason, "8 hours a day") {
		t.Fatalf("day limit: %+v", p)
	}
	_, err = testPool.Exec(ctx, `UPDATE player_characters SET stats=jsonb_set(stats,'{srd521_travel,seconds}','"lots"') WHERE id=$1`, actor)
	must(err)
	unavailable(travel(srd521.PaceNormal), "travel record is unreadable")
}
