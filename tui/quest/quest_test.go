package quest

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/tui/msgs"
)

// Compile-time checks: *Model must satisfy tea.Model and the sub-model
// interface (tea.Model + SetSize) expected by the Router. Importing the parent
// tui package here would create a cycle, so the interface is inlined.
var _ tea.Model = (*Model)(nil)
var _ interface {
	tea.Model
	SetSize(width, height int)
} = (*Model)(nil)

func TestNewReturnsModel(t *testing.T) {
	m := New()
	if m.width != 0 || m.height != 0 {
		t.Fatal("expected zero dimensions on a freshly created model")
	}
}

func TestSetSizeUpdatesDimensions(t *testing.T) {
	m := New()
	m.SetSize(80, 24)
	if m.width != 80 || m.height != 24 {
		t.Fatalf("expected 80x24, got %dx%d", m.width, m.height)
	}
}

func TestInitReturnsNil(t *testing.T) {
	m := New()
	if m.Init() != nil {
		t.Fatal("Init() should return nil")
	}
}

func TestEscEmitsNavigateBackMsg(t *testing.T) {
	m := New()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected a command on Esc, got nil")
	}
	msg := cmd()
	if _, ok := msg.(msgs.NavigateBackMsg); !ok {
		t.Fatalf("expected NavigateBackMsg, got %T", msg)
	}
}

func TestOtherKeyReturnsNoCmd(t *testing.T) {
	m := New()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Fatal("expected no command for non-Esc keys")
	}
}

func TestViewContainsQuestLogHeader(t *testing.T) {
	m := New()
	m.SetSize(80, 24)
	got := m.View()
	if !strings.Contains(got, "Quest Log") {
		t.Fatalf("expected View to contain %q, got:\n%s", "Quest Log", got)
	}
}

func TestViewShowsActiveQuestsFirstWithTheirObjectives(t *testing.T) {
	m := New()
	m.SetSize(100, 100)
	m.Show([]game.QuestRecord{
		{Quest: domain.Quest{Title: "Old Debt", Status: domain.QuestStatusFailed},
			Objectives: []domain.QuestObjective{{Description: "Pay the ferryman"}}},
		{Quest: domain.Quest{Title: "The Ledger", Description: "Someone tore out a page.", Status: domain.QuestStatusActive},
			Objectives: []domain.QuestObjective{{Description: "Ask the archivist", Completed: true}, {Description: "Find the page"}}},
	})
	view := m.View()
	for _, text := range []string{"◆ The Ledger", "Someone tore out a page.", "✓ Ask the archivist", "○ Find the page", "◆ Old Debt  (failed)", "✗ Pay the ferryman"} {
		if !strings.Contains(view, text) {
			t.Errorf("the quest log does not show %q:\n%s", text, view)
		}
	}
	if strings.Index(view, "The Ledger") > strings.Index(view, "Old Debt") {
		t.Errorf("the active quest is not first:\n%s", view)
	}
	if strings.Contains(view, "No quests yet") {
		t.Errorf("unexpected text:\n%s", view)
	}
}

func TestViewSaysWhenThereAreNoQuests(t *testing.T) {
	m := New()
	m.SetSize(80, 24)
	if view := m.View(); strings.Contains(view, "No quests") || strings.Contains(view, "Missing Merchant") {
		t.Fatalf("an unread quest log is called empty or shows a sample:\n%s", view)
	}
	m.Show(nil)
	if view := m.View(); !strings.Contains(view, "No quests yet") {
		t.Fatalf("want %q:\n%s", "No quests yet", view)
	}
}
