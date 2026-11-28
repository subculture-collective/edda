//go:build integration

package handlers_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/auth"
	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/handlers"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// TestIntegrationGMRulings declares route facts and a fall through the
// owner-only HTTP routes, then reads them back through the travel and fall
// previews. Every refusal leaves the facts and the session log unchanged.
func TestIntegrationGMRulings(t *testing.T) {
	ctx := context.Background()
	q := statedb.New(testPool)
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	campaignFor := func(name, rules string) (uuid.UUID, uuid.UUID) {
		t.Helper()
		c, err := q.CreateCampaign(ctx, statedb.CreateCampaignParams{Name: name, RulesVersion: rules, Status: "active", CreatedBy: dbutil.ToPgtype(testUserID)})
		must(err)
		campaign := dbutil.FromPgtype(c.ID)
		t.Cleanup(func() { _, _ = testPool.Exec(ctx, `DELETE FROM campaigns WHERE id=$1`, campaign) })
		build, err := srd521.SupportedBuildForConcept("Human Cleric", "Acolyte")
		must(err)
		stats, err := build.Stats()
		must(err)
		pc, err := q.CreatePlayerCharacter(ctx, statedb.CreatePlayerCharacterParams{CampaignID: c.ID, UserID: dbutil.ToPgtype(testUserID), Name: "Mira", Stats: stats, Hp: 10, MaxHp: 10, Level: 1, Status: "active", Abilities: []byte(`[]`)})
		must(err)
		return campaign, dbutil.FromPgtype(pc.ID)
	}
	location := func(campaign uuid.UUID, name string) uuid.UUID {
		t.Helper()
		var id uuid.UUID
		must(testPool.QueryRow(ctx, `INSERT INTO locations(campaign_id,name,player_known)VALUES($1,$2,true)RETURNING id`, campaign, name).Scan(&id))
		return id
	}
	connect := func(campaign, from, to uuid.UUID, both bool) {
		t.Helper()
		_, err := testPool.Exec(ctx, `INSERT INTO location_connections(from_location_id,to_location_id,bidirectional,travel_time,campaign_id)VALUES($1,$2,$3,'an hour',$4)`, from, to, both, campaign)
		must(err)
	}

	campaign, actor := campaignFor("Ruling fixture", "srd-5.2.1")
	gate, wood, ford, marsh := location(campaign, "Gate"), location(campaign, "Wood"), location(campaign, "Ford"), location(campaign, "Marsh")
	connect(campaign, gate, wood, false)
	connect(campaign, ford, gate, true) // declared from the Gate side
	connect(campaign, gate, marsh, false)
	connect(campaign, marsh, gate, true) // two connections lead there: ambiguous
	_, err := testPool.Exec(ctx, `UPDATE player_characters SET current_location_id=$1 WHERE id=$2`, gate, actor)
	must(err)
	foreign, foreignActor := campaignFor("Foreign fixture", "srd-5.2.1")
	foreignA, foreignB := location(foreign, "Elsewhere"), location(foreign, "Beyond")
	connect(foreign, foreignA, foreignB, false)
	legacy, legacyActor := campaignFor("Legacy fixture", "legacy")
	legacyA, legacyB := location(legacy, "Old"), location(legacy, "New")
	connect(legacy, legacyA, legacyB, false)

	e, err := engine.New(testPool, nil, config.LLMConfig{})
	must(err)
	logger := log.New(io.Discard)
	routerFor := func(user uuid.UUID) http.Handler {
		owner := &handlers.CampaignHandlers{Queries: q, Logger: logger}
		ah := &handlers.ActionHandlers{Engine: e, Queries: q, Logger: logger}
		router := chi.NewRouter()
		router.Use(auth.NewNoOpMiddleware(user).Authenticate)
		router.Route("/campaigns/{id}", func(r chi.Router) {
			r.Use(owner.RequireOwnership)
			r.Put("/locations/{lid}/routes/{to}", ah.DeclareRoute)
			r.Post("/character/fall", ah.DeclareFall)
		})
		return router
	}
	ownerRouter, strangerRouter := routerFor(testUserID), routerFor(uuid.New())
	call := func(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
		t.Helper()
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(method, path, strings.NewReader(body)))
		return rec
	}
	routePath := func(c, from, to uuid.UUID) string {
		return "/campaigns/" + c.String() + "/locations/" + from.String() + "/routes/" + to.String()
	}
	fallPath := func(c uuid.UUID) string { return "/campaigns/" + c.String() + "/character/fall" }
	type snapshot struct {
		Routes string
		Stats  string
		Logs   int
	}
	state := func() snapshot {
		t.Helper()
		var s snapshot
		must(testPool.QueryRow(ctx, `SELECT COALESCE((SELECT string_agg(COALESCE(srd521_route::text,'-'),',' ORDER BY id) FROM location_connections),''),COALESCE((SELECT string_agg(stats::text,',' ORDER BY id) FROM player_characters),''),(SELECT count(*) FROM session_logs)`).Scan(&s.Routes, &s.Stats, &s.Logs))
		return s
	}
	refused := func(rec *httptest.ResponseRecorder, code int, reason string, before snapshot) {
		t.Helper()
		if rec.Code != code || !strings.Contains(rec.Body.String(), reason) {
			t.Fatalf("want %d %q, got %d %s", code, reason, rec.Code, rec.Body.String())
		}
		if after := state(); after != before {
			t.Fatalf("refusal changed state:\n%+v\n%+v", before, after)
		}
	}
	forest := `{"profile":"travel-route-v1","distance_feet":5280,"terrain":"forest","good_road":false}`
	fall := func(id uuid.UUID, rest string) string {
		return `{"actor_id":"` + id.String() + `",` + rest + `}`
	}
	solid30 := `"height_feet":30,"landing":"solid","flying":false,"onto_creature":false,"fall_prevention":false`

	// Refusals: other users, foreign rows, malformed facts.
	before := state()
	refused(call(strangerRouter, http.MethodPut, routePath(campaign, gate, wood), forest), 404, "campaign not found", before)
	refused(call(strangerRouter, http.MethodPost, fallPath(campaign), fall(actor, solid30)), 404, "campaign not found", before)
	refused(call(ownerRouter, http.MethodPut, routePath(campaign, foreignA, foreignB), forest), 404, "not found", before)
	refused(call(ownerRouter, http.MethodPut, routePath(campaign, wood, gate), forest), 404, "not found", before) // one-way connection
	refused(call(ownerRouter, http.MethodPost, fallPath(campaign), fall(foreignActor, solid30)), 404, "not found", before)
	refused(call(ownerRouter, http.MethodPut, routePath(campaign, gate, marsh), forest), 409, "ambiguous", before)
	refused(call(ownerRouter, http.MethodPut, routePath(legacy, legacyA, legacyB), forest), 422, "SRD 5.2.1", before)
	refused(call(ownerRouter, http.MethodPost, fallPath(legacy), fall(legacyActor, solid30)), 422, "SRD 5.2.1", before)
	for body, reason := range map[string]string{
		`{"profile":"travel-route-v1","distance_feet":5280,"terrain":"forest"}`:                                "not declared",
		`{"profile":"travel-route-v1","distance_feet":0,"terrain":"forest","good_road":false}`:                 "positive",
		`{"profile":"travel-route-v1","distance_feet":5280,"terrain":"lava","good_road":false}`:                "Travel Terrain",
		`{"profile":"travel-route-v1","distance_feet":5280,"terrain":"waterborne","good_road":false}`:          "vehicle",
		`{"profile":"travel-route-v2","distance_feet":5280,"terrain":"forest","good_road":false}`:              "not declared",
		`{"profile":"travel-route-v1","distance_feet":5280,"terrain":"forest","good_road":false,"mount":true}`: "unknown or malformed",
		`{"profile":"travel-route-v1","distance_feet":"far","terrain":"forest","good_road":false}`:             "unknown or malformed",
		forest + `{}`: "unknown or malformed",
	} {
		refused(call(ownerRouter, http.MethodPut, routePath(campaign, gate, wood), body), 422, reason, before)
	}
	for body, reason := range map[string]string{
		fall(actor, `"landing":"solid","flying":false,"onto_creature":false,"fall_prevention":false`):                  "exact height",
		fall(actor, `"height_feet":-5,"landing":"solid","flying":false,"onto_creature":false,"fall_prevention":false`): "exact height",
		fall(actor, `"height_feet":30,"landing":"lava","flying":false,"onto_creature":false,"fall_prevention":false`):  "adjudication",
		fall(actor, `"height_feet":30,"landing":"solid","onto_creature":false,"fall_prevention":false`):                "Clarify whether",
		fall(actor, `"height_feet":30,"landing":"solid","flying":true,"onto_creature":false,"fall_prevention":false`):  "adjudication",
		fall(actor, solid30+`,"damage":12`):   "invalid fall",
		fall(actor, `"id":"chosen",`+solid30): "server assigns",
		`{` + solid30 + `}`:                   "actor_id",
	} {
		refused(call(ownerRouter, http.MethodPost, fallPath(campaign), body), 422, reason, before)
	}

	// A route declared from either end of a bidirectional connection is stored
	// canonically and listed by the travel preview.
	prepare := func(intent game.ActionIntent) *game.ActionPreview {
		t.Helper()
		intent.ActorID = actor
		c, pc, err := game.ReadActionSnapshot(ctx, testPool, campaign, actor, testUserID, false)
		must(err)
		a, err := game.PrepareStoredAction(ctx, testPool, c, pc, intent, false)
		must(err)
		return a.Preview
	}
	if rec := call(ownerRouter, http.MethodPut, routePath(campaign, gate, wood), forest); rec.Code != 200 || !strings.Contains(rec.Body.String(), `"terrain":"forest"`) {
		t.Fatalf("declare route: %d %s", rec.Code, rec.Body.String())
	}
	if rec := call(ownerRouter, http.MethodPut, routePath(campaign, gate, ford), `{"good_road":true,"terrain":"grassland","distance_feet":10560,"profile":"travel-route-v1"}`); rec.Code != 200 {
		t.Fatalf("declare reverse route: %d %s", rec.Code, rec.Body.String())
	}
	var stored string
	must(testPool.QueryRow(ctx, `SELECT srd521_route::text FROM location_connections WHERE from_location_id=$1 AND to_location_id=$2`, ford, gate).Scan(&stored))
	if stored != `{"profile": "travel-route-v1", "terrain": "grassland", "good_road": true, "distance_feet": 10560}` {
		t.Fatalf("stored route: %s", stored)
	}
	normal := prepare(game.ActionIntent{Kind: srd521.TravelHour, EntryID: srd521.PaceNormal})
	reasons := map[uuid.UUID]string{}
	for _, target := range normal.Targets {
		reasons[target.ID] = target.UnavailableReason
	}
	if !normal.Available || reasons[wood] != "" || reasons[ford] != "" || !strings.Contains(reasons[marsh], "ambiguous") {
		t.Fatalf("travel preview: %+v", normal)
	}
	var logs []string
	rows, err := testPool.Query(ctx, `SELECT input_type||':'||llm_response||':'||tool_calls::text FROM session_logs WHERE campaign_id=$1 ORDER BY turn_number`, campaign)
	must(err)
	for rows.Next() {
		var line string
		must(rows.Scan(&line))
		logs = append(logs, line)
	}
	must(rows.Err())
	if len(logs) != 2 || !strings.HasPrefix(logs[0], "meta:GM ruling: the route from Gate to Wood") || !strings.Contains(logs[0], `"Tool": "gm_declare_route"`) {
		t.Fatalf("route audit: %q", logs)
	}

	// A declared fall gets a server ID and resolves through the fall preview.
	rec := call(ownerRouter, http.MethodPost, fallPath(campaign), fall(actor, solid30))
	var declared game.DeclaredFall
	if rec.Code != 200 || json.Unmarshal(rec.Body.Bytes(), &declared) != nil || !strings.HasPrefix(declared.ID, "fall-") || declared.HeightFeet == nil || *declared.HeightFeet != 30 {
		t.Fatalf("declare fall: %d %s", rec.Code, rec.Body.String())
	}
	preview := prepare(game.ActionIntent{Kind: srd521.FallResolve, EntryID: declared.ID})
	if !preview.Available || !strings.Contains(preview.Options[0], "Fall 30 feet onto solid ground") {
		t.Fatalf("fall preview: %+v", preview)
	}
	if travel := prepare(game.ActionIntent{Kind: srd521.TravelHour, EntryID: srd521.PaceNormal}); travel.Available {
		t.Fatal("travel ignored the declared fall")
	}
	var fallLog string
	must(testPool.QueryRow(ctx, `SELECT input_type||':'||llm_response FROM session_logs WHERE campaign_id=$1 ORDER BY turn_number DESC LIMIT 1`, campaign).Scan(&fallLog))
	if fallLog != "meta:GM ruling: Mira falls 30 feet onto solid ground. Resolve the fall to roll its damage." {
		t.Fatalf("fall audit: %q", fallLog)
	}
	// A second declaration never replaces an unresolved one.
	refused(call(ownerRouter, http.MethodPost, fallPath(campaign), fall(actor, solid30)), 409, "Resolve the declared fall", state())
}
