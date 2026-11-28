// Package inventory provides the inventory view for the TUI. It lists what the
// current character carries.
package inventory

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/tui/msgs"
	"git.subcult.tv/subculture-collective/edda/tui/scrollbox"
	"git.subcult.tv/subculture-collective/edda/tui/styles"
)

// Model is the Bubble Tea model for the inventory view.
type Model struct {
	width, height int
	box           scrollbox.Model

	items  []domain.Item
	loaded bool
	notice string
}

// New returns an inventory Model that has nothing to show yet.
func New() Model {
	return Model{}
}

// Show replaces the list with items, in the web app's order: equipped first,
// then the rarest, then by name.
func (m *Model) Show(items []domain.Item) {
	items = slices.Clone(items)
	slices.SortStableFunc(items, func(a, b domain.Item) int {
		if a.Equipped != b.Equipped {
			if a.Equipped {
				return -1
			}
			return 1
		}
		return cmp.Or(cmp.Compare(rarityRank(b.Rarity), rarityRank(a.Rarity)),
			cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)))
	})
	m.items, m.loaded, m.notice = items, true, ""
}

// Notify says why the list is missing or may be out of date. Items that are
// already shown stay.
func (m *Model) Notify(notice string) { m.notice = notice }

func rarityRank(rarity string) int {
	return slices.Index([]string{"common", "uncommon", "rare", "epic", "legendary"}, strings.ToLower(rarity)) + 1
}

// SetSize updates the viewport dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd { return nil }

// Update implements tea.Model. Pressing Escape emits NavigateBackMsg so the
// parent App can return focus to the narrative view. Other messages scroll a
// list that does not fit.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEsc {
		return m, func() tea.Msg { return msgs.NavigateBackMsg{} }
	}
	return m, m.box.Update(msg, m.content(), m.width, m.height)
}

// View implements tea.Model and renders the inventory.
func (m Model) View() string {
	return m.box.View(m.content(), m.width, m.height)
}

// content renders everything inside the box.
func (m Model) content() string {
	lines := []string{styles.Header.Render("🎒 Inventory"), ""}
	if m.notice != "" {
		lines = append(lines, styles.SystemMessage.Render(m.notice))
	}
	switch {
	case !m.loaded:
	case len(m.items) == 0:
		lines = append(lines, styles.Muted.Render("No items yet."))
	default:
		for _, item := range m.items {
			lines = append(lines, itemLines(item)...)
		}
	}
	lines = append(lines, "", styles.Muted.Render("Press Esc to return to the narrative view."))
	return styles.JoinVertical(lines...)
}

func itemLines(item domain.Item) []string {
	name := "◆ " + item.Name
	if item.Quantity != 1 {
		name += fmt.Sprintf(" ×%d", item.Quantity)
	}
	details := []string{string(item.ItemType)}
	if item.Rarity != "" {
		details = append(details, item.Rarity)
	}
	if item.Equipped {
		details = append(details, "equipped")
	}
	lines := []string{styles.SubHeader.Render(name) + styles.Muted.Render("  "+strings.Join(details, " · "))}
	if item.Description != "" {
		lines = append(lines, styles.Body.Render("  "+item.Description))
	}
	return lines
}
