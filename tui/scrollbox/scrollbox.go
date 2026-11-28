// Package scrollbox draws a tab's content in a bordered box that scrolls when
// the content needs more rows than the tab has.
package scrollbox

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"git.subcult.tv/subculture-collective/edda/tui/styles"
)

// Model keeps the scroll position. The caller passes the content and the size
// to Update and View, so both can change between them. The box may use width
// columns and height rows, and is only as tall as its content. A zero height
// leaves the rows unlimited. The zero value is ready.
type Model struct {
	view viewport.Model
}

// Update scrolls content with the viewport's keys and the mouse wheel.
func (m *Model) Update(msg tea.Msg, content string, width, height int) tea.Cmd {
	m.layout(content, width, height)
	var cmd tea.Cmd
	m.view, cmd = m.view.Update(msg)
	return cmd
}

// View renders content from the current scroll position.
func (m Model) View(content string, width, height int) string {
	if width <= 0 || height <= 0 {
		return styles.Frame(styles.Container, width, 0).Render(content)
	}
	scrolls := m.layout(content, width, height) // m is a copy, so the tab is unchanged
	body := m.view.View()
	if scrolls && rows(height) > 1 {
		hint := fmt.Sprintf("↑/↓ PgUp/PgDn: scroll (%.0f%%)", m.view.ScrollPercent()*100)
		body += "\n" + styles.Muted.Render(ansi.Truncate(hint, m.view.Width, "…"))
	}
	return styles.Frame(styles.Container, width, 0).Render(body)
}

// rows is the number of rows inside a box height rows tall.
func rows(height int) int {
	return max(height-styles.Container.GetVerticalFrameSize(), 1)
}

// layout wraps content to the box and reports whether it has to scroll. A
// scrolling box gives its last row to the hint.
func (m *Model) layout(content string, width, height int) (scrolls bool) {
	if width <= 0 || height <= 0 {
		return false
	}
	inner := max(width-styles.Container.GetHorizontalFrameSize(), 1)
	wrapped := lipgloss.NewStyle().Width(inner).Render(content)
	shown, available := lipgloss.Height(wrapped), rows(height)
	if scrolls = shown > available; scrolls {
		shown = max(available-1, 1)
	}
	if m.view.Width == 0 {
		m.view = viewport.New(inner, shown) // sets the key and wheel bindings
	}
	m.view.Width = inner
	m.view.Height = shown
	m.view.SetContent(wrapped)
	m.view.SetYOffset(m.view.YOffset) // a taller box may need a smaller offset
	return scrolls
}
