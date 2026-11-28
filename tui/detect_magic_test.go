package tui

import (
	"context"
	"strings"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"git.subcult.tv/subculture-collective/edda/tui/narrative"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

type detectMagicMock struct {
	tacticalMock
	self          uuid.UUID
	concentration *srd521.Concentration
	intents       []game.ActionIntent
}

func (m *detectMagicMock) PreviewDetectMagic(_ context.Context, _ uuid.UUID, intent game.ActionIntent) (*game.StoredAction, error) {
	m.intents = append(m.intents, intent)
	p := &game.ActionPreview{Available: true, ExpectedRevision: "detect", Intent: intent, Options: []string{"Detect Magic ruling for " + intent.Kind + "."}}
	if intent.Kind == srd521.SpellCast {
		p.Targets = []game.ActionTarget{{ID: m.self, Label: "Mira (self)"}}
	} else if m.concentration == nil {
		p.Available, p.Reason = false, "Detect Magic concentration is required."
	}
	return &game.StoredAction{Preview: p, Concentration: m.concentration}, nil
}

func TestTUIDetectMagicCastRetryStudyAndStatus(t *testing.T) {
	m := &detectMagicMock{self: uuid.New()}
	m.fail = true
	a := NewAppWithEngine(testCfg, testCampaign, context.Background(), m, nil)
	model, _ := a.Update(tea.WindowSizeMsg{Width: 200, Height: 60})
	a = model.(App)
	send := func(input string) {
		t.Helper()
		model, cmd := a.Update(narrative.SubmitMsg{Input: input})
		a = model.(App)
		if cmd != nil {
			model, _ = a.Update(cmd())
			a = model.(App)
		}
	}
	shows := func(want string) {
		t.Helper()
		if !strings.Contains(a.narrativeView().View(), want) {
			t.Fatalf("missing %q:\n%s", want, a.narrativeView().View())
		}
	}
	send("/confirm-detect-magic")
	send("/detect-magic")
	shows("No active concentration shown. Detect Magic concentration is required.")
	send("/study-detect-magic")
	for _, bad := range []string{"/detect-magic ritual", "/detect-magic cast 1", "/study-detect-magic now", "/confirm-detect-magic now"} {
		send(bad)
	}
	if a.detectMagicRequest != nil || len(m.requests) != 0 || len(m.intents) != 2 {
		t.Fatalf("unavailable or bad input prepared: %+v", m.intents)
	}
	send("/detect-magic cast")
	cast := m.intents[2]
	if a.detectMagicRequest == nil || cast.Kind != srd521.SpellCast || cast.EntryID != srd521.DetectMagicID || cast.SlotLevel != 0 || a.detectMagicRequest.TargetID != m.self || len(m.requests) != 0 {
		t.Fatalf("cast preview: %+v %+v", cast, a.detectMagicRequest)
	}
	send("/study-detect-magic")
	if a.detectMagicRequest.Intent.Kind != srd521.SpellCast || len(m.intents) != 3 {
		t.Fatal("pending request replaced")
	}
	send("/confirm-detect-magic")
	if a.detectMagicRequest == nil || len(m.requests) != 1 {
		t.Fatal("uncertain request lost")
	}
	m.fail = false
	send("/confirm-detect-magic")
	if len(m.requests) != 2 || m.requests[0] != m.requests[1] || a.detectMagicRequest != nil {
		t.Fatal("retry identity changed")
	}

	m.concentration = &srd521.Concentration{Spell: srd521.DetectMagicID, ExpiresAt: 32400}
	send("/detect-magic")
	shows("Active concentration: Detect Magic, expires Day 1, 09:00:00.")
	if a.detectMagicRequest != nil {
		t.Fatal("status prepared an action")
	}
	send("/study-detect-magic")
	if a.detectMagicRequest == nil || a.detectMagicRequest.Intent.Kind != srd521.DetectMagicStudy || a.detectMagicRequest.TargetID != uuid.Nil {
		t.Fatalf("study preview: %+v", a.detectMagicRequest)
	}
	send("/cancel-detect-magic")
	if a.detectMagicRequest != nil || len(m.requests) != 2 {
		t.Fatal("cancel committed")
	}
}
