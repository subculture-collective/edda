//go:build integration

package handlers_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/config"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
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
