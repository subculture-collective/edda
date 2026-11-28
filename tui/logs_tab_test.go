package tui

import (
	"log/slog"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"git.subcult.tv/subculture-collective/edda/internal/logging"
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
