package inventory

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
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

func TestViewContainsInventoryHeader(t *testing.T) {
	m := New()
	m.SetSize(80, 24)
	got := m.View()
	if !strings.Contains(got, "Inventory") {
		t.Fatalf("expected View to contain %q, got:\n%s", "Inventory", got)
	}
}

func TestViewContainsEscHint(t *testing.T) {
	m := New()
	m.SetSize(80, 24)
	got := m.View()
	if !strings.Contains(got, "Esc") {
		t.Fatalf("expected View to contain Esc key hint, got:\n%s", got)
	}
}

func TestViewListsItemsEquippedThenRarestThenByName(t *testing.T) {
	m := New()
	m.SetSize(100, 100)
	m.Show([]domain.Item{
		{Name: "Rope", ItemType: domain.ItemTypeMisc, Rarity: "common", Quantity: 1},
		{Name: "Javelin", ItemType: domain.ItemTypeWeapon, Rarity: "common", Quantity: 8},
		{Name: "Moon Charm", ItemType: domain.ItemTypeMisc, Rarity: "Rare", Quantity: 1, Description: "Cold to the touch"},
		{Name: "Shield", ItemType: domain.ItemTypeArmor, Rarity: "common", Quantity: 1, Equipped: true},
	})
	view := m.View()
	order := []string{"Shield", "Moon Charm", "Cold to the touch", "Javelin ×8", "Rope"}
	last := -1
	for _, text := range order {
		at := strings.Index(view, text)
		if at <= last {
			t.Fatalf("want %v in order, %q is out of place:\n%s", order, text, view)
		}
		last = at
	}
	for _, text := range []string{"armor · common · equipped", "misc · Rare"} {
		if !strings.Contains(view, text) {
			t.Errorf("the inventory does not show %q:\n%s", text, view)
		}
	}
	if strings.Contains(view, "No items yet") || strings.Contains(view, "×1") {
		t.Errorf("unexpected text:\n%s", view)
	}
}

func TestViewSaysWhenThereAreNoItems(t *testing.T) {
	m := New()
	m.SetSize(80, 24)
	if view := m.View(); strings.Contains(view, "No items") {
		t.Fatalf("an unread inventory is called empty:\n%s", view)
	}
	m.Show(nil)
	if view := m.View(); !strings.Contains(view, "No items yet") {
		t.Fatalf("want %q:\n%s", "No items yet", view)
	}
}
