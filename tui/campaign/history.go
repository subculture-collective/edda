package campaign

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"git.subcult.tv/subculture-collective/edda/tui/styles"
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

// latestOffset is the history line on which the last message starts.
func latestOffset(messages []string, width int) int {
	if len(messages) < 2 {
		return 0
	}
	// A blank line separates the last message from the ones before it.
	return lipgloss.Height(wrapHistory(messages[:len(messages)-1], width)) + 1
}

// showLatest fills the history and scrolls to the start of the last message, so
// that a long reply is read from the top. The viewport stops at the bottom, so
// a short one leaves the newest lines in view.
func showLatest(history *viewport.Model, messages []string) {
	history.SetContent(wrapHistory(messages, history.Width))
	history.SetYOffset(latestOffset(messages, history.Width))
}

// resizeHistory wraps the history to a new size and keeps the reader's place:
// the start of the last message, the bottom, or the same line. The launcher
// sizes its view after every message, so an unchanged size must not scroll.
func resizeHistory(history *viewport.Model, messages []string, width, height int) {
	if width == history.Width && height == history.Height {
		return
	}
	atLatest := history.YOffset == latestOffset(messages, history.Width)
	atBottom := history.AtBottom()
	history.Width = width
	history.Height = height
	history.SetContent(wrapHistory(messages, width))
	switch {
	case atLatest:
		history.SetYOffset(latestOffset(messages, width))
	case atBottom:
		history.GotoBottom()
	default:
		history.SetYOffset(history.YOffset) // a taller box may need a smaller offset
	}
}

// moreBelow is the line between the history and the input. It says when part of
// the history is below the screen.
func moreBelow(history viewport.Model) string {
	if history.AtBottom() {
		return ""
	}
	return styles.Muted.Render(ansi.Truncate("↓ more below: ↓ or PgDn", history.Width, "…"))
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
