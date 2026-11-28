package tui

import (
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"git.subcult.tv/subculture-collective/edda/internal/logging"
	"git.subcult.tv/subculture-collective/edda/tui/logpanel"
)

func TestLogsTabRefreshesWhileAnotherTabIsOpen(t *testing.T) {
	buf := logging.NewRingBuffer(16)
	next, cmd := resize(t, NewLauncherWithEngine(testCfg, t.Context(), nil, nil, WithLogBuffer(buf)), 100, 30).
		Update(campaignLoadedMsg{c: testCampaign})
	app, ok := next.(App)
	if !ok {
		t.Fatalf("expected App, got %T", next)
	}
	if cmd == nil {
		t.Fatal("handing over to the app did not start the log panel's timer")
	}

	// Two ticks: the first must schedule the second.
	for _, text := range []string{"first late entry", "second late entry"} {
		buf.Write(logging.LogEntry{Time: time.Now(), Level: slog.LevelInfo, Service: "app", Message: text})
		if app.ActiveViewState() != ViewNarrative {
			t.Fatalf("expected the narrative tab, got %d", app.ActiveViewState())
		}
		var updated tea.Model
		updated, cmd = app.Update(cmd())
		app = updated.(App)
		if cmd == nil {
			t.Fatalf("timer stopped after %q", text)
		}

		// Look without delivering anything else to the panel.
		app.router.GoToTab(int(ViewLogs))
		if view := app.View(); !strings.Contains(view, text) {
			t.Fatalf("%q is not on the Logs tab:\n%s", text, view)
		}
		app.router.GoToTab(int(ViewNarrative))
	}
}

func TestLogsTabScrollsWithTheMouseWheel(t *testing.T) {
	buf := logging.NewRingBuffer(64)
	logLine := func(text string) {
		buf.Write(logging.LogEntry{Time: time.Now(), Level: slog.LevelInfo, Service: "app", Message: text})
	}
	for i := range 40 {
		logLine(fmt.Sprintf("entry %02d.", i))
	}
	app := resize(t, NewAppWithEngine(testCfg, testCampaign, t.Context(), nil, buf), 80, 24)
	app.router.GoToTab(int(ViewLogs))
	var m tea.Model = app
	wheel := func(button tea.MouseButton, times int) {
		for range times {
			m, _ = m.Update(tea.MouseMsg{Button: button, Action: tea.MouseActionPress})
		}
	}
	check := func(step, shown, hidden string) {
		t.Helper()
		view := m.View()
		if !strings.Contains(view, shown) || strings.Contains(view, hidden) {
			t.Fatalf("%s: want %q and not %q:\n%s", step, shown, hidden, view)
		}
	}

	check("at first", "entry 39.", "entry 00.")
	wheel(tea.MouseButtonWheelUp, 40)
	check("after wheeling up", "entry 00.", "entry 39.")

	// A new entry must not pull a player who scrolled up back to the end.
	logLine("late entry")
	m, _ = m.Update(logpanel.TickMsg{})
	check("after a new entry", "entry 00.", "late entry")

	// Back at the end, the panel follows new entries again.
	wheel(tea.MouseButtonWheelDown, 40)
	check("after wheeling down", "late entry", "entry 00.")
	logLine("later entry")
	m, _ = m.Update(logpanel.TickMsg{})
	check("after another entry", "later entry", "entry 00.")
}

// A wider panel has fewer lines. A player who scrolled up must still see a full
// panel, not blank rows under the last line.
func TestLogsTabStaysFullWhenWidenedWhileScrolledUp(t *testing.T) {
	buf := logging.NewRingBuffer(64)
	for range 13 {
		buf.Write(logging.LogEntry{Time: time.Now(), Level: slog.LevelInfo, Service: "app", Message: strings.Repeat("wrapped words ", 20)})
	}
	panel := logpanel.New(buf)
	panel.SetSize(40, 16)
	// Scroll to the top, then down to a line the wider panel also has.
	for range 20 {
		panel.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	}
	for range 8 {
		panel.Update(tea.KeyMsg{Type: tea.KeyDown})
	}

	panel.SetSize(400, 16)
	// The header, then 12 rows of entries, then the footer, inside the border.
	lines := strings.Split(panel.View(), "\n")
	for _, line := range lines[2 : len(lines)-2] {
		if !strings.Contains(line, "wrapped words") {
			t.Fatalf("the panel has a row with no entry:\n%s", panel.View())
		}
	}
}
