//go:build integration

package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/combat"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
)

// addDoor places one door object at the actor's location and makes it known
// through campaign history, like Detect Magic's object fixtures.
func (f *spatialFixture) addDoor(name string, known bool) uuid.UUID {
	f.t.Helper()
	ctx := context.Background()
	var door uuid.UUID
	f.must(testPool.QueryRow(ctx, `INSERT INTO npcs(campaign_id,name,hp,location_id) SELECT $1,$2,1,current_location_id FROM player_characters WHERE id=$3 RETURNING id`, f.campaign, name, f.actor).Scan(&door))
	if known {
		_, err := testPool.Exec(ctx, `INSERT INTO session_logs(campaign_id,turn_number,player_input,input_type,llm_response,npcs_involved) SELECT $1,COALESCE(MAX(turn_number),0)+1,'Look','game_action','A door.',ARRAY[$2]::uuid[] FROM session_logs WHERE campaign_id=$1`, f.campaign, door)
		f.must(err)
	}
	return door
}

func (f *spatialFixture) doorFacts(door uuid.UUID, actor uuid.UUID, facts string) {
	f.t.Helper()
	props := fmt.Sprintf(`{"srd521_door":{"profile":"ordinary-door-v1","actor_id":"%s",%s},"secret":"unpublished vault"}`, actor, facts)
	_, err := testPool.Exec(context.Background(), `UPDATE npcs SET properties=$1 WHERE id=$2`, props, door)
	f.must(err)
}

func (f *spatialFixture) openDoor(door uuid.UUID) *game.ActionPreview {
	f.t.Helper()
	ctx := context.Background()
	c, pc, err := game.ReadActionSnapshot(ctx, testPool, f.campaign, f.actor, testUserID, false)
	f.must(err)
	a, err := game.PrepareStoredAction(ctx, testPool, c, pc, game.ActionIntent{ActorID: f.actor, Kind: srd521.OpenDoor, EntryID: door.String()}, false)
	f.must(err)
	return a.Preview
}

func (f *spatialFixture) doorState(door uuid.UUID) string {
	f.t.Helper()
	var state string
	f.must(testPool.QueryRow(context.Background(), `SELECT properties->'srd521_door'->>'state' FROM npcs WHERE id=$1`, door).Scan(&state))
	return state
}

func TestIntegrationOpenOrdinaryDoor(t *testing.T) {
	f := newSpatialFixture(t)
	ctx := context.Background()
	unavailable := func(p *game.ActionPreview, reason string) {
		t.Helper()
		if p.Available || !strings.Contains(p.Reason, reason) {
			t.Fatalf("want unavailable %q: %+v", reason, p)
		}
	}
	closed := `"hidden":false,"state":"closed","locked":false,"distance_feet":5`
	door := f.addDoor("Oak door", true)
	f.doorFacts(door, f.actor, closed)
	f.senses(`{"darkvision_feet":0}`)

	// Unknown, hidden, foreign and never-encountered doors share one reason and
	// disclose no label or private facts.
	unknown := f.openDoor(uuid.New())
	if unknown.Available || len(unknown.Targets) != 0 {
		t.Fatalf("random door: %+v", unknown)
	}
	stranger := f.addDoor("Vault door", false)
	f.doorFacts(stranger, f.actor, closed)
	foreign := f.addDoor("Foreign door", true)
	f.doorFacts(foreign, uuid.New(), closed)
	hidden := f.addDoor("Secret door", true)
	f.doorFacts(hidden, f.actor, `"hidden":true,"state":"closed","locked":false,"distance_feet":5`)
	for name, id := range map[string]uuid.UUID{"Vault door": stranger, "Foreign door": foreign, "Secret door": hidden} {
		p := f.openDoor(id)
		encoded, err := json.Marshal(p)
		f.must(err)
		if p.Available || p.Reason != unknown.Reason || len(p.Targets) != 0 || strings.Contains(string(encoded), name) || strings.Contains(string(encoded), "unpublished vault") {
			t.Fatalf("unknown door disclosed: %s", encoded)
		}
	}

	// The sheet and the TUI offer only the known door, by name and state.
	want := []srd521.SheetDoor{{ID: door.String(), Name: "Oak door", State: srd521.DoorClosed}}
	if got := characterSheet(t, f.campaign).Doors; !slices.Equal(got, want) {
		t.Fatalf("sheet doors: %+v", got)
	}
	listed, none, err := f.engine().PreviewDoor(ctx, f.campaign, 0)
	f.must(err)
	_, tuiPreview, err := f.engine().PreviewDoor(ctx, f.campaign, 1)
	f.must(err)
	if !slices.Equal(listed, want) || none != nil || tuiPreview == nil || !tuiPreview.Available || tuiPreview.Intent.EntryID != door.String() || tuiPreview.Intent.ActorID != f.actor {
		t.Fatalf("TUI door: %+v %+v %+v", listed, none, tuiPreview)
	}
	// A refused door is still named, without a state it cannot be trusted for.
	f.doorFacts(door, f.actor, closed+`,"trap":"poison needle"`)
	if got := characterSheet(t, f.campaign).Doors; len(got) != 1 || got[0].State != "" {
		t.Fatalf("refused door: %+v", got)
	}
	f.doorFacts(door, f.actor, closed)

	// Range, locks and undeclared mechanics fail closed; nothing is invented.
	f.doorFacts(door, f.actor, `"hidden":false,"state":"closed","locked":false,"distance_feet":10`)
	unavailable(f.openDoor(door), "Out of reach")
	f.doorFacts(door, f.actor, `"hidden":false,"state":"closed","locked":true,"distance_feet":5`)
	unavailable(f.openDoor(door), "Locked doors are not supported")
	f.doorFacts(door, f.actor, `"hidden":false,"state":"closed","distance_feet":5`)
	unavailable(f.openDoor(door), "Clarify whether the door is locked")
	f.doorFacts(door, f.actor, closed+`,"trap":"poison needle"`)
	unavailable(f.openDoor(door), "does not support")
	f.doorFacts(door, f.actor, closed)

	// A target behind the closed door is not seen and cannot be attacked.
	f.facts(fmt.Sprintf(`"distance_feet":5,"line_of_sight":true,"cover":"none","light":"bright","obscurement":"none","door":"%s"`, door))
	unavailable(f.sight(f.target), "can't see")
	staleAttack := f.prepare(srd521.MartialAttack)
	if len(staleAttack.Targets) != 1 || !strings.Contains(staleAttack.Targets[0].UnavailableReason, "line of sight") {
		t.Fatalf("attack through a closed door: %+v", staleAttack)
	}
	staleSight := f.sight(f.target)

	p := f.openDoor(door)
	if !p.Available || len(p.Targets) != 1 || p.Targets[0].Label != "Oak door" || *p.Targets[0].DistanceFeet != 5 || p.Targets[0].ReachFeet != 5 || len(p.Costs) != 0 || !strings.Contains(p.Options[0], "no action and no campaign time") {
		t.Fatalf("door preview: %+v", p)
	}
	if again := f.openDoor(door); again.ExpectedRevision != p.ExpectedRevision {
		t.Fatal("reload changed the door preview")
	}
	competing := f.openDoor(door)

	// History failure rolls back the door state and the attempt.
	_, err = testPool.Exec(ctx, `ALTER TABLE session_logs ADD CONSTRAINT door_history_failure CHECK (player_input NOT LIKE 'Confirm Open door%') NOT VALID`)
	f.must(err)
	t.Cleanup(func() {
		_, _ = testPool.Exec(ctx, `ALTER TABLE session_logs DROP CONSTRAINT IF EXISTS door_history_failure`)
	})
	failed := engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: p.ExpectedRevision, Intent: p.Intent}
	if _, err := f.engine().ConfirmAction(ctx, f.campaign, failed); err == nil {
		t.Fatal("history failure committed")
	}
	var attempts int
	f.must(testPool.QueryRow(ctx, `SELECT count(*) FROM turn_attempts WHERE id=$1`, failed.ActionID).Scan(&attempts))
	if attempts != 0 || f.doorState(door) != "closed" {
		t.Fatalf("partial commit: attempts=%d state=%s", attempts, f.doorState(door))
	}
	_, err = testPool.Exec(ctx, `ALTER TABLE session_logs DROP CONSTRAINT door_history_failure`)
	f.must(err)

	// A target ID in the request, or a foreign door ID in the intent, refreshes.
	for _, bad := range []engine.PreparedAction{
		{ActionID: uuid.New(), ExpectedRevision: p.ExpectedRevision, Intent: p.Intent, TargetID: door},
		{ActionID: uuid.New(), ExpectedRevision: p.ExpectedRevision, Intent: game.ActionIntent{ActorID: f.actor, Kind: srd521.OpenDoor, EntryID: foreign.String()}},
	} {
		if _, err := f.engine().ConfirmAction(ctx, f.campaign, bad); !errors.Is(err, engine.ErrActionRefresh) {
			t.Fatalf("bad request confirmed: %v", err)
		}
	}

	// Two concurrent confirmations: exactly one commits.
	countHistory := func() int {
		t.Helper()
		var n int
		f.must(testPool.QueryRow(ctx, `SELECT count(*) FROM session_logs WHERE campaign_id=$1 AND player_input LIKE 'Confirm Open door%'`, f.campaign).Scan(&n))
		return n
	}
	requests := []engine.PreparedAction{
		{ActionID: uuid.New(), ExpectedRevision: p.ExpectedRevision, Intent: p.Intent},
		{ActionID: uuid.New(), ExpectedRevision: competing.ExpectedRevision, Intent: competing.Intent},
	}
	results := make([]*engine.TurnResult, 2)
	errs := make([]error, 2)
	e := f.engine()
	var wg sync.WaitGroup
	for i := range requests {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i], errs[i] = e.ConfirmAction(ctx, f.campaign, requests[i])
		}(i)
	}
	wg.Wait()
	winner := -1
	for i, err := range errs {
		if err == nil {
			if winner >= 0 {
				t.Fatal("both confirmations committed")
			}
			winner = i
		} else if !errors.Is(err, engine.ErrActionRefresh) && !errors.Is(err, engine.ErrTurnInProgress) {
			t.Fatalf("loser: %v", err)
		}
	}
	if winner < 0 {
		t.Fatalf("no confirmation committed: %v", errs)
	}
	if history := countHistory(); history != 1 || f.doorState(door) != "open" {
		t.Fatalf("history=%d state=%s", history, f.doorState(door))
	}
	var payload struct {
		Receipt srd521.DoorReceipt `json:"receipt"`
	}
	result := results[winner]
	f.must(json.Unmarshal(result.AppliedToolCalls[0].Result, &payload))
	if result.AppliedToolCalls[0].Tool != "srd_door" || payload.Receipt.Before != "closed" || payload.Receipt.After != "open" || payload.Receipt.DoorID != door.String() || !strings.Contains(result.Narrative, "Oak door is now open") {
		t.Fatalf("receipt: %+v %q", payload.Receipt, result.Narrative)
	}

	// The same action ID replays the saved receipt without writing again.
	replayed, err := f.engine().ConfirmAction(ctx, f.campaign, requests[winner])
	f.must(err)
	if replayed.Narrative != result.Narrative || replayed.TurnNumber != result.TurnNumber || countHistory() != 1 {
		t.Fatalf("replay changed history: %+v", replayed)
	}

	// Opening invalidates the spatial previews the door obstructed.
	for _, stale := range []*game.ActionPreview{staleAttack, staleSight} {
		request := engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: stale.ExpectedRevision, Intent: stale.Intent}
		if stale.Intent.Kind == srd521.MartialAttack {
			request.TargetID = f.target
		}
		if _, err := f.engine().ConfirmAction(ctx, f.campaign, request); !errors.Is(err, engine.ErrActionRefresh) {
			t.Fatalf("stale %s preview confirmed: %v", stale.Intent.Kind, err)
		}
	}
	seen := f.sight(f.target)
	if !seen.Available || seen.ExpectedRevision == staleSight.ExpectedRevision {
		t.Fatalf("target behind the open door: %+v", seen)
	}
	if attack := f.prepare(srd521.MartialAttack); !attack.Available || attack.Targets[0].UnavailableReason != "" {
		t.Fatalf("attack through the open door: %+v", attack)
	}
	// Reload: a fresh snapshot keeps the door open and the open preview stays refused.
	unavailable(f.openDoor(door), "already open")
	if f.sight(f.target).ExpectedRevision != seen.ExpectedRevision {
		t.Fatal("reload changed the spatial revision")
	}

	// A target naming a door the actor does not know stays unresolved.
	f.facts(fmt.Sprintf(`"distance_feet":5,"line_of_sight":true,"cover":"none","light":"bright","obscurement":"none","door":"%s"`, hidden))
	unavailable(f.sight(f.target), "Clarify the exact distance and line of sight")

	// Encounters: opening is unavailable, and a pending reaction is named first.
	f.facts(`"distance_feet":5,"line_of_sight":true,"cover":"none","light":"bright","obscurement":"none"`)
	f.doorFacts(door, f.actor, closed)
	f.confirm(f.prepare(combat.EncounterStart))
	unavailable(f.openDoor(door), "during an encounter")
	encounter, err := game.ReadEncounter(ctx, testPool, f.campaign)
	f.must(err)
	encounter.Reaction = &combat.ReactionState{Available: true, Pending: &combat.ReactionTrigger{ID: uuid.New(), TargetID: f.target, WeaponID: f.weapon}}
	saved, err := json.Marshal([]map[string]any{{"Tool": "srd_encounter", "Result": map[string]any{"encounter": encounter}}})
	f.must(err)
	_, err = testPool.Exec(ctx, `INSERT INTO session_logs(campaign_id,turn_number,player_input,input_type,llm_response,tool_calls) SELECT $1,MAX(turn_number)+1,'Reaction','game_action','A foe moves.',$2 FROM session_logs WHERE campaign_id=$1`, f.campaign, saved)
	f.must(err)
	unavailable(f.openDoor(door), "pending reaction")
	if f.doorState(door) != "closed" {
		t.Fatal("door changed during an encounter")
	}
}
