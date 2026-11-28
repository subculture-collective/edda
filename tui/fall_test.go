package tui

import (
	"context"
	"slices"
	"strings"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"git.subcult.tv/subculture-collective/edda/tui/narrative"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

type fallMock struct {
	tacticalMock
	declared bool
	skills   []string
}

func (m *fallMock) PreviewFall(_ context.Context, _ uuid.UUID, skill string) (*game.ActionPreview, error) {
	m.skills = append(m.skills, skill)
	if !m.declared {
		return &game.ActionPreview{Reason: "No declared fall to resolve."}, nil
	}
	return &game.ActionPreview{Available: true, ExpectedRevision: "fall", Intent: game.ActionIntent{Kind: srd521.FallResolve, EntryID: "f1", Skill: skill}, Options: []string{"Fall 20 feet onto water."}}, nil
}

func TestTUIFallPreviewCancelConfirmAndRetry(t *testing.T) {
	m := &fallMock{}
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
	send("/confirm-fall")
	send("/fall")
	shows("No declared fall to resolve.")
	m.declared = true
	for _, bad := range []string{"/fall 30", "/fall athletics 30", "/confirm-fall 30", "/cancel-fall now"} {
		send(bad)
		if a.fallRequest != nil || len(m.skills) != 1 {
			t.Fatalf("%q prepared a fall: %v", bad, m.skills)
		}
	}
	send("/fall")
	shows("Preview only: Fall 20 feet onto water.")
	send("/fall athletics")
	if a.fallRequest == nil || a.fallRequest.Intent.Skill != "" || len(m.skills) != 2 {
		t.Fatalf("pending request replaced: %+v", a.fallRequest)
	}
	send("/cancel-fall")
	if a.fallRequest != nil || len(m.requests) != 0 {
		t.Fatal("cancel committed")
	}
	send("/fall Athletics")
	if a.fallRequest == nil || a.fallRequest.Intent.EntryID != "f1" || a.fallRequest.Intent.Skill != "athletics" || a.fallRequest.TargetID != uuid.Nil || len(m.requests) != 0 {
		t.Fatalf("water preview: %+v", a.fallRequest)
	}
	send("/confirm-fall")
	if a.fallRequest == nil || len(m.requests) != 1 {
		t.Fatal("uncertain request lost")
	}
	m.fail = false
	send("/confirm-fall")
	if len(m.requests) != 2 || m.requests[0] != m.requests[1] || a.fallRequest != nil {
		t.Fatal("retry identity changed")
	}
	if !slices.Equal(m.skills, []string{"", "", "athletics"}) {
		t.Fatalf("previewed skills: %v", m.skills)
	}
}
