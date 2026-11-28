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

type lightMock struct {
	tacticalMock
	torch, mace uuid.UUID
	light       *srd521.LightEffect
	intents     []game.ActionIntent
}

func (m *lightMock) PreviewLight(_ context.Context, _ uuid.UUID, intent game.ActionIntent) (*game.StoredAction, error) {
	m.intents = append(m.intents, intent)
	p := &game.ActionPreview{Available: true, ExpectedRevision: "light", Intent: intent, Options: []string{"Light ruling for " + intent.Kind + "."}}
	if intent.Kind == srd521.SpellCast {
		p.Targets = []game.ActionTarget{{ID: m.torch, Label: "Torch (carried; touch)"}, {ID: m.mace, Label: "Mace (carried; touch)"}}
	} else if m.light == nil {
		p.Available, p.Reason = false, "No active Light to dismiss."
	}
	return &game.StoredAction{Preview: p, Light: m.light, LightItem: "Mace"}, nil
}

func TestTUILightPickConfirmRetryCoverAndDismiss(t *testing.T) {
	m := &lightMock{torch: uuid.New(), mace: uuid.New()}
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
	send("/confirm-light")
	send("/light")
	shows("No active Light")
	send("/light Blue")
	shows("2. Mace (carried; touch)")
	if a.lightRequest != nil || len(m.requests) != 0 || m.intents[1].EntryID != "light:blue" {
		t.Fatalf("listing prepared or committed: %+v", m.intents)
	}
	for _, bad := range []string{"/light blue 3", "/light blue 0", "/light blue torch", "/light blue 1 extra", "/dismiss-light now"} {
		send(bad)
		if a.lightRequest != nil {
			t.Fatalf("%q prepared a cast", bad)
		}
	}
	send("/light blue 2")
	if a.lightRequest == nil || a.lightRequest.TargetID != m.mace || len(m.requests) != 0 {
		t.Fatalf("cast preview: %+v", a.lightRequest)
	}
	send("/light red 1")
	if a.lightRequest.TargetID != m.mace {
		t.Fatal("pending request replaced")
	}
	send("/confirm-light")
	if a.lightRequest == nil || len(m.requests) != 1 {
		t.Fatal("uncertain request lost")
	}
	m.fail = false
	send("/confirm-light")
	if len(m.requests) != 2 || m.requests[0] != m.requests[1] || a.lightRequest != nil {
		t.Fatal("retry identity changed")
	}

	m.light = &srd521.LightEffect{ItemID: m.mace.String(), Color: "blue", ExpiresAt: 32400, Covered: true}
	send("/light")
	shows("Active Light on Mace (blue), expires Day 1, 09:00:00. Covered")
	if a.lightRequest != nil {
		t.Fatal("status prepared an action")
	}
	for command, kind := range map[string]string{"/cover-light": srd521.CoverLight, "/uncover-light": srd521.UncoverLight, "/dismiss-light": srd521.DismissLight} {
		send(command)
		if a.lightRequest == nil || a.lightRequest.Intent.Kind != kind || a.lightRequest.TargetID != uuid.Nil {
			t.Fatalf("%s: %+v", command, a.lightRequest)
		}
		send("/cancel-light")
		if a.lightRequest != nil || len(m.requests) != 2 {
			t.Fatal("cancel committed")
		}
	}
}
