package campaign

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// wrapHistory joins an interview's messages and wraps them to width. The
// viewport cuts off anything wider.
func wrapHistory(messages []string, width int) string {
	content := strings.Join(messages, "\n\n")
	if width <= 0 {
		return content
	}
	return lipgloss.NewStyle().Width(width).Render(content)
}

// scrollsHistory reports whether the chat history should see msg. The viewport
// also binds letters and space, which the player types into the input.
func scrollsHistory(msg tea.Msg) bool {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return true
	}
	switch key.Type {
	case tea.KeyUp, tea.KeyDown, tea.KeyPgUp, tea.KeyPgDown:
		return true
	default:
		return false
	}
}
