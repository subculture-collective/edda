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

type deathSaveMock struct {
	tacticalMock
	atZero   bool
	previews int
}

func (m *deathSaveMock) PreviewDeathSave(context.Context, uuid.UUID) (*game.ActionPreview, error) {
	m.previews++
	if !m.atZero {
		return &game.ActionPreview{Reason: "Death Saving Throws are made only at 0 HP."}, nil
	}
	return &game.ActionPreview{Available: true, ExpectedRevision: "death-save", Intent: game.ActionIntent{Kind: srd521.DeathSave}, Options: []string{"Death Saving Throw: d20. Now 1 successes, 2 failures."}}, nil
}

func TestTUIDeathSavePreviewCancelConfirmAndRetry(t *testing.T) {
	m := &deathSaveMock{}
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
	send("/confirm-death-save")
	send("/death-save")
	shows("only at 0 HP")
	m.atZero = true
	send("/death-save 20")
	if a.deathSaveRequest != nil || m.previews != 1 {
		t.Fatal("an argument prepared a save")
	}
	send("/death-save")
	shows("Now 1 successes, 2 failures.")
	send("/cancel-death-save")
	if a.deathSaveRequest != nil || len(m.requests) != 0 {
		t.Fatal("cancel committed")
	}
	send("/death-save")
	send("/confirm-death-save")
	if a.deathSaveRequest == nil || len(m.requests) != 1 {
		t.Fatal("uncertain request lost")
	}
	m.fail = false
	send("/confirm-death-save")
	if len(m.requests) != 2 || m.requests[0] != m.requests[1] || a.deathSaveRequest != nil || m.requests[0].Intent.Kind != srd521.DeathSave {
		t.Fatal("retry identity changed")
	}
}
