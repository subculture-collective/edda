// Package inventory provides the inventory view for the TUI.
// It renders a placeholder inventory using the shared styles package.
// Real inventory data will be populated in the inventory management epic.
package inventory

import (
	tea "github.com/charmbracelet/bubbletea"

	"git.subcult.tv/subculture-collective/edda/tui/msgs"
	"git.subcult.tv/subculture-collective/edda/tui/scrollbox"
	"git.subcult.tv/subculture-collective/edda/tui/styles"
)

// Model is the Bubble Tea model for the inventory view.
type Model struct {
	width, height int
	box           scrollbox.Model
}

// New returns a freshly initialised inventory Model.
func New() Model {
	return Model{}
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
// sheet that does not fit.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok && key.Type == tea.KeyEsc {
		return m, func() tea.Msg { return msgs.NavigateBackMsg{} }
	}
	return m, m.box.Update(msg, m.content(), m.width, m.height)
}

// View implements tea.Model and renders the inventory placeholder.
func (m Model) View() string {
	return m.box.View(m.content(), m.width, m.height)
}

// content renders everything inside the box.
func (m Model) content() string {
	title := styles.Header.Render("🎒 Inventory")

	placeholder := styles.SystemMessage.Render(
		"Inventory data will be available in a future update.",
	)
	hint := styles.Muted.Render("Press Esc to return to the narrative view.")

	content := styles.JoinVertical(placeholder, "", hint)
	return styles.JoinVertical(title, "", content)
}
