//go:build integration

package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/combat"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
)

func (f *spatialFixture) senses(senses string) {
	f.t.Helper()
	query := `UPDATE player_characters SET stats=stats-'srd521_senses' WHERE id=$1`
	args := []any{f.actor}
	if senses != "" {
		query, args = `UPDATE player_characters SET stats=jsonb_set(stats,'{srd521_senses}',$2::jsonb) WHERE id=$1`, append(args, senses)
	}
	_, err := testPool.Exec(context.Background(), query, args...)
	f.must(err)
}

func (f *spatialFixture) sight(target uuid.UUID) *game.ActionPreview {
	f.t.Helper()
	ctx := context.Background()
	c, pc, err := game.ReadActionSnapshot(ctx, testPool, f.campaign, f.actor, testUserID, false)
	f.must(err)
	a, err := game.PrepareStoredAction(ctx, testPool, c, pc, game.ActionIntent{ActorID: f.actor, Kind: srd521.AbilityCheck, EntryID: game.SightPerceptionEntryPrefix + target.String()}, false)
	f.must(err)
	return a.Preview
}

func TestIntegrationLightAndDarkvisionVisibility(t *testing.T) {
	f := newSpatialFixture(t)
	unavailable := func(p *game.ActionPreview, reason string) {
		t.Helper()
		if p.Available || !strings.Contains(p.Reason, reason) {
			t.Fatalf("want unavailable %q: %+v", reason, p)
		}
	}
	// Legacy facts declare no light: sight is unresolved, attacks keep #67 rules.
	unavailable(f.sight(f.target), "Clarify the light")
	f.facts(`"distance_feet":30,"line_of_sight":true,"cover":"none","light":"dim","obscurement":"none"`)
	unavailable(f.sight(f.target), "Clarify the observer's senses")

	f.senses(`{"darkvision_feet":0}`)
	p := f.sight(f.target)
	if !p.Available || len(p.Targets) != 1 || p.Targets[0].Label != "Practice foe" || !strings.Contains(p.Options[0], "Disadvantage") {
		t.Fatalf("dim without Darkvision: %+v", p)
	}
	f.dice.d20 = []int{4, 15}
	request, result := f.confirm(p)
	var payload struct {
		Receipt srd521.CheckReceipt `json:"receipt"`
		Sight   combat.Sight        `json:"sight"`
	}
	f.must(json.Unmarshal(result.AppliedToolCalls[0].Result, &payload))
	if payload.Receipt.Mode != "disadvantage" || len(payload.Receipt.Dice) != 2 || payload.Sight.Circumstance != combat.SightDisadvantage || f.dice.calls != 2 {
		t.Fatalf("disadvantage receipt: %+v calls=%d", payload, f.dice.calls)
	}
	// A replayed confirmation returns the saved receipt without rolling.
	replayed, err := f.engine().ConfirmAction(context.Background(), f.campaign, request)
	f.must(err)
	if f.dice.calls != 2 || replayed.Narrative != result.Narrative {
		t.Fatalf("replay rolled or changed: %+v", replayed)
	}

	// Darkvision 60 feet: dim light within range is seen as bright.
	f.senses(`{"darkvision_feet":60}`)
	p = f.sight(f.target)
	if !p.Available || strings.Contains(p.Options[0], "Disadvantage") || !strings.Contains(p.Options[0], "Darkvision") {
		t.Fatalf("dim with Darkvision: %+v", p)
	}
	f.facts(`"distance_feet":60,"line_of_sight":true,"cover":"none","light":"dark","obscurement":"none"`)
	if p := f.sight(f.target); !p.Available || !strings.Contains(p.Options[0], "Disadvantage") {
		t.Fatalf("darkness at the Darkvision boundary: %+v", p)
	}
	stale := f.sight(f.target)
	f.facts(`"distance_feet":65,"line_of_sight":true,"cover":"none","light":"dark","obscurement":"none"`)
	unavailable(f.sight(f.target), "can't see")
	_, err = f.engine().ConfirmAction(context.Background(), f.campaign, engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: stale.ExpectedRevision, Intent: stale.Intent})
	if !errors.Is(err, engine.ErrActionRefresh) || f.dice.calls != 2 {
		t.Fatalf("stale light ruling confirmed: %v calls=%d", err, f.dice.calls)
	}
	// A changed sense ruling also invalidates an open preview.
	f.facts(`"distance_feet":30,"line_of_sight":true,"cover":"none","light":"dark","obscurement":"none"`)
	stale = f.sight(f.target)
	f.senses(`{"darkvision_feet":0}`)
	_, err = f.engine().ConfirmAction(context.Background(), f.campaign, engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: stale.ExpectedRevision, Intent: stale.Intent})
	if !errors.Is(err, engine.ErrActionRefresh) || f.dice.calls != 2 {
		t.Fatalf("stale senses confirmed: %v", err)
	}
	unavailable(f.sight(f.target), "can't see")

	// Total cover and heavy fog block sight even in bright light or Darkvision.
	f.senses(`{"darkvision_feet":60}`)
	f.facts(`"distance_feet":10,"line_of_sight":true,"cover":"total","light":"bright","obscurement":"none"`)
	unavailable(f.sight(f.target), "can't see")
	f.facts(`"distance_feet":10,"line_of_sight":true,"cover":"none","light":"dim","obscurement":"heavy"`)
	unavailable(f.sight(f.target), "can't see")
	f.facts(`"distance_feet":10,"line_of_sight":true,"cover":"none","light":"bright"`)
	unavailable(f.sight(f.target), "Clarify whether fog")

	// Attacks honor declared light: darkness hides the target without Darkvision.
	f.facts(`"distance_feet":5,"line_of_sight":true,"cover":"none","light":"dark","obscurement":"none"`)
	if p := f.prepare(srd521.MartialAttack); !p.Available {
		t.Fatalf("Darkvision attack unavailable: %+v", p)
	}
	f.senses(`{"darkvision_feet":0}`)
	if p := f.prepare(srd521.MartialAttack); p.Available || !strings.Contains(p.Targets[0].UnavailableReason, "can't see") {
		t.Fatalf("unseen attack admitted: %+v", p)
	}

	// The ruling survives reload: a fresh snapshot derives the same result.
	f.senses(`{"darkvision_feet":60}`)
	first := f.sight(f.target)
	if again := f.sight(f.target); again.ExpectedRevision != first.ExpectedRevision || again.Options[0] != first.Options[0] {
		t.Fatal("reload changed the derivation")
	}

	// History failure rolls back the receipt and attempt.
	_, err = testPool.Exec(context.Background(), `ALTER TABLE session_logs ADD CONSTRAINT sight_history_failure CHECK (player_input NOT LIKE 'Confirm Sight-based%') NOT VALID`)
	f.must(err)
	t.Cleanup(func() {
		_, _ = testPool.Exec(context.Background(), `ALTER TABLE session_logs DROP CONSTRAINT IF EXISTS sight_history_failure`)
	})
	failed := engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: first.ExpectedRevision, Intent: first.Intent}
	if _, err := f.engine().ConfirmAction(context.Background(), f.campaign, failed); err == nil {
		t.Fatal("history failure committed")
	}
	var attempts int
	f.must(testPool.QueryRow(context.Background(), `SELECT count(*) FROM turn_attempts WHERE id=$1`, failed.ActionID).Scan(&attempts))
	if attempts != 0 {
		t.Fatal("partial commit after history failure")
	}
	_, err = testPool.Exec(context.Background(), `ALTER TABLE session_logs DROP CONSTRAINT sight_history_failure`)
	f.must(err)

	// Option labels never reveal a hidden or unknown actor.
	unknown := f.sight(uuid.New())
	f.facts(`"distance_feet":5,"line_of_sight":true,"cover":"none","light":"bright","obscurement":"none","hidden":true`)
	hidden := f.sight(f.target)
	encoded, err := json.Marshal(hidden)
	f.must(err)
	if hidden.Available || hidden.Reason != unknown.Reason || len(hidden.Targets) != 0 || strings.Contains(string(encoded), "Practice foe") || strings.Contains(string(encoded), "unpublished villain") {
		t.Fatalf("hidden actor revealed: %s", encoded)
	}

	// Sight checks are outside encounters only.
	f.facts(`"distance_feet":5,"line_of_sight":true,"cover":"none","light":"bright","obscurement":"none"`)
	f.confirm(f.prepare(combat.EncounterStart))
	if p := f.sight(f.target); p.Available || len(p.Targets) != 0 {
		t.Fatalf("sight check during encounter: %+v", p)
	}
}
