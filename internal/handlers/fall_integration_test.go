//go:build integration

package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/combat"
	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/handlers"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/pkg/api"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// fallFaces returns queued die faces (1-based) for each roll, d20s and d6s
// alike; an empty queue returns an invalid value so a failure is injected.
type fallFaces struct {
	faces []int
	calls int
}

func (d *fallFaces) Intn(n int) int {
	d.calls++
	if len(d.faces) == 0 {
		return n
	}
	face := d.faces[0]
	d.faces = d.faces[1:]
	return face - 1
}

type fallFixture struct {
	*spatialFixture
	faces *fallFaces
}

func (f *fallFixture) declare(fall string) {
	f.t.Helper()
	_, err := testPool.Exec(context.Background(), `UPDATE player_characters SET stats=jsonb_set(stats,'{srd521_fall}',$2::jsonb) WHERE id=$1`, f.actor, fall)
	f.must(err)
}

func (f *fallFixture) setHP(hp int) {
	f.t.Helper()
	_, err := testPool.Exec(context.Background(), `UPDATE player_characters SET hp=$2 WHERE id=$1`, f.actor, hp)
	f.must(err)
}

func (f *fallFixture) preview(kind, entry, skill string) *game.ActionPreview {
	f.t.Helper()
	ctx := context.Background()
	c, pc, err := game.ReadActionSnapshot(ctx, testPool, f.campaign, f.actor, testUserID, false)
	f.must(err)
	a, err := game.PrepareStoredAction(ctx, testPool, c, pc, game.ActionIntent{ActorID: f.actor, Kind: kind, EntryID: entry, Skill: skill}, false)
	f.must(err)
	return a.Preview
}

func (f *fallFixture) confirmWith(p *game.ActionPreview, faces ...int) (engine.PreparedAction, *engine.TurnResult, error) {
	f.t.Helper()
	f.faces.faces, f.faces.calls = faces, 0
	e, err := engine.New(testPool, nil, config.LLMConfig{}, engine.WithActionDice(f.faces))
	f.must(err)
	r := engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: p.ExpectedRevision, Intent: p.Intent}
	if p.Intent.Kind == combat.EncounterStart || p.Intent.Kind == srd521.MartialAttack {
		r.TargetID = f.target
	}
	result, err := e.ConfirmAction(context.Background(), f.campaign, r)
	return r, result, err
}

// state reads the actor's HP, whether a fall is still declared, and saved Prone.
func (f *fallFixture) state() (hp int, declared, prone bool) {
	f.t.Helper()
	f.must(testPool.QueryRow(context.Background(), `SELECT hp,stats ? 'srd521_fall',COALESCE((stats->'srd521_conditions'->>'prone')::boolean,false) FROM player_characters WHERE id=$1`, f.actor).Scan(&hp, &declared, &prone))
	return hp, declared, prone
}

func fallDecl(id string, height int, landing string) string {
	b, _ := json.Marshal(map[string]any{"id": id, "height_feet": height, "landing": landing, "flying": false, "onto_creature": false, "fall_prevention": false})
	return string(b)
}

func fallReceipt(t *testing.T, result *engine.TurnResult) srd521.FallReceipt {
	t.Helper()
	if len(result.AppliedToolCalls) != 1 || result.AppliedToolCalls[0].Tool != "srd_fall" {
		t.Fatalf("tool calls: %+v", result.AppliedToolCalls)
	}
	var payload struct {
		Receipt srd521.FallReceipt `json:"receipt"`
	}
	if err := json.Unmarshal(result.AppliedToolCalls[0].Result, &payload); err != nil {
		t.Fatal(err)
	}
	return payload.Receipt
}

func TestIntegrationDeclaredFall(t *testing.T) {
	f := &fallFixture{spatialFixture: newSpatialFixture(t), faces: &fallFaces{}}
	ctx := context.Background()
	unavailable := func(p *game.ActionPreview, reason string) {
		t.Helper()
		if p.Available || !strings.Contains(p.Reason, reason) {
			t.Fatalf("want unavailable %q: %+v", reason, p)
		}
	}
	unavailable(f.preview(srd521.FallResolve, "f1", ""), "No declared fall")

	// Excluded and unknown facts fail closed.
	f.declare(`{"id":"f1","height_feet":30,"landing":"solid","flying":true,"onto_creature":false,"fall_prevention":false}`)
	unavailable(f.preview(srd521.FallResolve, "f1", ""), "require adjudication")
	f.declare(`{"id":"f1","height_feet":30,"landing":"solid"}`)
	unavailable(f.preview(srd521.FallResolve, "f1", ""), "Clarify whether")
	f.declare(fallDecl("f1", 30, "lava"))
	unavailable(f.preview(srd521.FallResolve, "f1", ""), "requires adjudication")
	f.declare(fallDecl("f1", 30, "solid"))
	unavailable(f.preview(srd521.FallResolve, "f1", "athletics"), "Choose no check")
	unavailable(f.preview(srd521.FallResolve, "other", ""), "No declared fall")

	// Clients learn only a declaration's ID, plus height and landing once the
	// rules accept it. The TUI preview takes no ID, height or damage.
	e, err := engine.New(testPool, nil, config.LLMConfig{}, engine.WithActionDice(f.faces))
	f.must(err)
	tuiPreview := func() *game.ActionPreview {
		t.Helper()
		p, err := e.PreviewFall(ctx, f.campaign, "")
		f.must(err)
		return p
	}
	// pending reads the fall from the character sheet the web client loads.
	characters := &handlers.CharacterHandlers{Queries: statedb.New(testPool), Pool: testPool, Logger: log.New(io.Discard)}
	router := chi.NewRouter()
	router.Get("/campaigns/{id}/character", characters.GetCharacter)
	pending := func() *srd521.PendingFall {
		t.Helper()
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/campaigns/"+f.campaign.String()+"/character", nil))
		var response api.CharacterResponse
		f.must(json.Unmarshal(recorder.Body.Bytes(), &response))
		if response.Sheet == nil {
			t.Fatalf("sheet unavailable: %d %s", recorder.Code, recorder.Body)
		}
		return response.Sheet.Fall
	}
	f.declare(fallDecl("f1", 30, "lava"))
	if got := pending(); got == nil || got.ID != "f1" || got.HeightFeet != nil || got.Landing != "" {
		t.Fatalf("refused declaration projected: %+v", got)
	}
	unavailable(tuiPreview(), "requires adjudication")
	f.declare(fallDecl("f1", 30, "solid"))
	if got := pending(); got == nil || got.HeightFeet == nil || *got.HeightFeet != 30 || got.Landing != "solid" {
		t.Fatalf("pending fall: %+v", got)
	}
	if p := tuiPreview(); !p.Available || p.Intent.EntryID != "f1" || p.Intent.ActorID != f.actor {
		t.Fatalf("TUI preview: %+v", p)
	}

	// A changed height invalidates an open preview without rolling.
	stale := f.preview(srd521.FallResolve, "f1", "")
	f.declare(fallDecl("f1", 40, "solid"))
	if _, _, err := f.confirmWith(stale, 6, 6, 6); !errors.Is(err, engine.ErrActionRefresh) || f.faces.calls != 0 {
		t.Fatalf("stale height confirmed: %v calls=%d", err, f.faces.calls)
	}

	// An injected dice failure and a history failure both commit nothing.
	f.declare(fallDecl("f1", 35, "solid"))
	p := f.preview(srd521.FallResolve, "f1", "")
	if !p.Available || !strings.Contains(p.Options[0], "3d6 Bludgeoning") || p.Modifier != nil {
		t.Fatalf("35-foot preview: %+v", p)
	}
	if _, _, err := f.confirmWith(p, 2, 5); err == nil {
		t.Fatal("invalid third die committed")
	}
	if hp, declared, prone := f.state(); hp != 12 || !declared || prone {
		t.Fatalf("dice failure partially committed: hp=%d declared=%t prone=%t", hp, declared, prone)
	}
	_, err = testPool.Exec(ctx, `ALTER TABLE session_logs ADD CONSTRAINT fall_history_failure CHECK (player_input NOT LIKE 'Confirm Fall%') NOT VALID`)
	f.must(err)
	t.Cleanup(func() {
		_, _ = testPool.Exec(ctx, `ALTER TABLE session_logs DROP CONSTRAINT IF EXISTS fall_history_failure`)
	})
	if _, _, err := f.confirmWith(p, 2, 5, 1); err == nil {
		t.Fatal("history failure committed")
	}
	_, err = testPool.Exec(ctx, `ALTER TABLE session_logs DROP CONSTRAINT fall_history_failure`)
	f.must(err)
	var attempts int
	f.must(testPool.QueryRow(ctx, `SELECT count(*) FROM turn_attempts WHERE campaign_id=$1`, f.campaign).Scan(&attempts))
	if hp, declared, prone := f.state(); hp != 12 || !declared || prone || attempts != 0 {
		t.Fatalf("history failure partially committed: hp=%d declared=%t prone=%t attempts=%d", hp, declared, prone, attempts)
	}

	// 35 feet: three dice, damage and Prone commit once with the receipt.
	request, result, err := f.confirmWith(p, 2, 5, 1)
	f.must(err)
	r := fallReceipt(t, result)
	if r.DiceCount != 3 || r.Rolled != 8 || r.Damage != 8 || r.HPBefore != 12 || r.HPAfter != 4 || !r.ProneAfter || r.ProneBefore || r.DamageType != "bludgeoning" {
		t.Fatalf("receipt: %+v", r)
	}
	if !strings.Contains(result.Narrative, "3d6 (2, 5, 1) = 8") || !strings.Contains(result.Narrative, "Prone condition") {
		t.Fatalf("narrative: %s", result.Narrative)
	}
	if hp, declared, prone := f.state(); hp != 4 || declared || !prone {
		t.Fatalf("saved state: hp=%d declared=%t prone=%t", hp, declared, prone)
	}
	if pending() != nil {
		t.Fatal("consumed fall is still pending")
	}
	unavailable(tuiPreview(), "No declared fall")
	// Replay returns the saved result without rolling; the declaration is gone.
	f.faces.calls = 0
	replayed, err := e.ConfirmAction(ctx, f.campaign, request)
	f.must(err)
	if f.faces.calls != 0 || replayed.Narrative != result.Narrative {
		t.Fatalf("replay rolled or changed: %+v", replayed)
	}
	if hp, _, _ := f.state(); hp != 4 {
		t.Fatal("replay applied damage again")
	}
	unavailable(f.preview(srd521.FallResolve, "f1", ""), "No declared fall")
	if _, _, err := f.confirmWith(p, 6, 6, 6); !errors.Is(err, engine.ErrActionRefresh) {
		t.Fatalf("consumed fall confirmed again: %v", err)
	}

	// Standing outside an encounter ends the saved Prone condition.
	stand := f.preview(combat.EncounterStand, "", "")
	if !stand.Available {
		t.Fatalf("stand: %+v", stand)
	}
	_, result, err = f.confirmWith(stand)
	f.must(err)
	if _, _, prone := f.state(); prone || result.AppliedToolCalls[0].Tool != "srd_condition" {
		t.Fatalf("stand did not end Prone: %+v", result)
	}
	unavailable(f.preview(combat.EncounterStand, "", ""), "No supported encounter is active")

	// Under 10 feet: no dice, no damage, not Prone.
	f.declare(fallDecl("f2", 9, "solid"))
	_, result, err = f.confirmWith(f.preview(srd521.FallResolve, "f2", ""))
	f.must(err)
	if r := fallReceipt(t, result); r.DiceCount != 0 || r.Damage != 0 || r.ProneAfter || f.faces.calls != 0 {
		t.Fatalf("9-foot fall: %+v calls=%d", r, f.faces.calls)
	}
	if hp, declared, prone := f.state(); hp != 4 || declared || prone {
		t.Fatalf("9-foot state: hp=%d declared=%t prone=%t", hp, declared, prone)
	}

	// Water: a DC 15 Strength (Athletics) Reaction halves the damage.
	f.setHP(12)
	f.declare(fallDecl("f3", 20, "water"))
	if p := f.preview(srd521.FallResolve, "f3", ""); !p.Available || !strings.Contains(p.Options[1], "DC 15") {
		t.Fatalf("water preview: %+v", p)
	}
	p = f.preview(srd521.FallResolve, "f3", "athletics")
	if !p.Available || p.Modifier == nil || p.Modifier.Value != 5 || len(p.Costs) != 1 || p.Costs[0].ResourceID != "reaction" {
		t.Fatalf("athletics preview: %+v", p)
	}
	_, result, err = f.confirmWith(p, 10, 6, 5) // d20 10 + 5 = 15 meets the DC.
	f.must(err)
	r = fallReceipt(t, result)
	if r.WaterCheck == nil || r.WaterCheck.Total != 15 || r.WaterCheck.Outcome != "success" || !r.Halved || r.Rolled != 11 || r.Damage != 5 || r.HPAfter != 7 || !r.ProneAfter {
		t.Fatalf("water receipt: %+v check=%+v", r, r.WaterCheck)
	}

	// Maximum dice: 500 feet rolls 20d6, and HP stops at 0.
	f.setHP(12)
	f.declare(fallDecl("f4", 500, "solid"))
	p = f.preview(srd521.FallResolve, "f4", "")
	if !p.Available || !strings.Contains(p.Options[0], "20d6") || !strings.Contains(p.Options[0], "capped at 20d6") {
		t.Fatalf("max preview: %+v", p)
	}
	ones := make([]int, 20)
	for i := range ones {
		ones[i] = 1
	}
	_, result, err = f.confirmWith(p, ones...)
	f.must(err)
	if r := fallReceipt(t, result); r.DiceCount != 20 || len(r.Dice) != 20 || r.Damage != 20 || r.HPAfter != 0 || !strings.Contains(result.Narrative, "require adjudication") {
		t.Fatalf("max receipt: %+v %s", r, result.Narrative)
	}

	// A later encounter takes the saved Prone onto the player actor.
	f.setHP(12)
	if _, _, prone := f.state(); !prone {
		t.Fatal("max fall did not leave Prone")
	}
	_, _, err = f.confirmWith(f.prepare(combat.EncounterStart), 15, 5)
	f.must(err)
	encounter, err := game.ReadEncounter(ctx, testPool, f.campaign)
	f.must(err)
	for _, actor := range encounter.Actors {
		if actor.Player != actor.Prone {
			t.Fatalf("encounter Prone: %+v", encounter.Actors)
		}
	}
	if _, _, prone := f.state(); prone {
		t.Fatal("Prone saved twice: on the actor and in the encounter")
	}
	f.declare(fallDecl("f5", 20, "solid"))
	unavailable(f.preview(srd521.FallResolve, "f5", ""), "during an encounter")

	// The Prone player crawls into reach and defeats the foe; the completed
	// encounter hands Prone back to the actor for a later stand.
	for range 2 {
		_, _, err = f.confirmWith(f.prepare(combat.EncounterApproach))
		f.must(err)
	}
	_, err = testPool.Exec(ctx, `UPDATE npcs SET hp=1 WHERE id=$1`, f.target)
	f.must(err)
	_, result, err = f.confirmWith(f.prepare(srd521.MartialAttack), 20, 20, 8, 8)
	f.must(err)
	encounter, err = game.ReadEncounter(ctx, testPool, f.campaign)
	f.must(err)
	if encounter.Status != "completed" || !strings.Contains(result.Narrative, "Disadvantage") {
		t.Fatalf("encounter did not complete: %s %+v", result.Narrative, encounter)
	}
	if _, _, prone := f.state(); !prone {
		t.Fatal("completed encounter dropped Prone")
	}
	if p := f.preview(combat.EncounterStand, "", ""); !p.Available {
		t.Fatalf("stand after encounter: %+v", p)
	}
}
