package tui

import (
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/logging"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/internal/world"
	"git.subcult.tv/subculture-collective/edda/tui/campaign"
	"git.subcult.tv/subculture-collective/edda/tui/narrative"
)

var layoutSizes = []struct{ width, height int }{
	{80, 24},
	{100, 30},
	{160, 48},
}

// assertFits fails when view needs more columns or rows than the terminal has.
// Bubble Tea drops the top lines of a view that is too tall.
func assertFits(t *testing.T, view string, width, height int) {
	t.Helper()
	lines := strings.Split(view, "\n")
	widest := 0
	for _, line := range lines {
		widest = max(widest, lipgloss.Width(line))
	}
	if widest > width || len(lines) > height {
		t.Errorf("view is %d columns by %d lines, terminal is %dx%d\n%s", widest, len(lines), width, height, view)
	}
}

func resize[M tea.Model](t *testing.T, m M, width, height int) M {
	t.Helper()
	next, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: height})
	resized, ok := next.(M)
	if !ok {
		t.Fatalf("unexpected model type %T", next)
	}
	return resized
}

func TestLayout_AppViewsFitTerminal(t *testing.T) {
	long := strings.Repeat("The innkeeper lowers their voice and points east. ", 12)
	choose := func(texts ...string) []engine.Choice {
		choices := make([]engine.Choice, len(texts))
		for i, text := range texts {
			choices[i] = engine.Choice{ID: fmt.Sprint(i), Text: text}
		}
		return choices
	}
	narrativeStates := []struct {
		name  string
		setup func(*narrative.Model)
		want  []string // text that must be on screen
	}{
		{"welcome", func(*narrative.Model) {}, []string{"Welcome to Edda"}},
		{"long entries", func(nv *narrative.Model) {
			for i := range 10 {
				nv.AddEntry(narrative.Entry{Kind: narrative.EntryKind(i % 3), Speaker: "Edda", Text: long})
			}
			nv.AddEntry(narrative.Entry{Kind: narrative.KindSystem, Text: "newest entry"})
		}, []string{"newest entry"}},
		{"choices after long entries", func(nv *narrative.Model) {
			for range 10 {
				nv.AddEntry(narrative.Entry{Kind: narrative.KindNPC, Speaker: "Edda", Text: long})
			}
			nv.AddEntry(narrative.Entry{Kind: narrative.KindSystem, Text: "newest entry"})
			nv.SetChoices(choose(long, "Leave", "Wait"))
		}, []string{"newest entry", "Leave", "Wait"}},
		{"two choices", func(nv *narrative.Model) {
			nv.SetChoices(choose(long, "Leave"))
		}, []string{"Leave"}},
		{"more choices than shown while loading", func(nv *narrative.Model) {
			nv.SetChoices(choose(long, "Wait", "Hide", "Run", "Sleep"))
			nv.SetLoading(true)
		}, []string{"Thinking", "Hide", "••"}},
		{"long typed input", func(nv *narrative.Model) {
			nv.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(long + "typed last")})
		}, []string{"typed last"}},
	}

	// A turn in progress with choices on screen needs the most rows, so it is
	// also checked on terminals shorter than the other states need.
	choicesWhileThinking := func(nv *narrative.Model) {
		nv.AddEntry(narrative.Entry{Kind: narrative.KindSystem, Text: "newest entry"})
		nv.SetChoices(choose(long, "Leave", "Wait"))
		nv.SetLoading(true)
	}
	for _, size := range []struct{ width, height int }{{80, 24}, {100, 30}, {80, 16}, {100, 16}} {
		t.Run(fmt.Sprintf("%dx%d/narrative/choices while thinking", size.width, size.height), func(t *testing.T) {
			app := resize(t, testApp(), size.width, size.height)
			choicesWhileThinking(app.narrativeView())
			view := app.View()
			assertFits(t, view, size.width, size.height)
			for _, want := range []string{"newest entry", "Thinking", "Leave", "Wait", "What do you do?", "Narrative", "Views:"} {
				if !strings.Contains(view, want) {
					t.Errorf("%q is not on screen:\n%s", want, view)
				}
			}
		})
	}

	for _, size := range layoutSizes {
		for _, state := range narrativeStates {
			t.Run(fmt.Sprintf("%dx%d/narrative/%s", size.width, size.height, state.name), func(t *testing.T) {
				app := resize(t, testApp(), size.width, size.height)
				state.setup(app.narrativeView())
				view := app.View()
				assertFits(t, view, size.width, size.height)
				for _, want := range append(state.want, "Narrative", "Views:") {
					if !strings.Contains(view, want) {
						t.Errorf("%q is not on screen:\n%s", want, view)
					}
				}
			})
		}

		buf := logging.NewRingBuffer(100)
		for range 80 {
			buf.Write(logging.LogEntry{
				Time: time.Now(), Level: slog.LevelInfo, Service: "a-very-long-service-name", Message: long,
			})
		}
		app := NewAppWithEngine(testCfg, testCampaign, t.Context(), nil, buf)
		app = resize(t, app, size.width, size.height)
		for idx, tab := range app.router.Tabs() {
			t.Run(fmt.Sprintf("%dx%d/tab/%s", size.width, size.height, tab.Name), func(t *testing.T) {
				app.router.GoToTab(idx)
				if ViewState(idx) == ViewLogs {
					// "s" filters to one service, which the panel header names.
					app.router.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")})
				}
				assertFits(t, app.View(), size.width, size.height)
			})
		}

		t.Run(fmt.Sprintf("%dx%d/long campaign name", size.width, size.height), func(t *testing.T) {
			campaign := testCampaign
			campaign.Name = strings.Repeat("The Very Long Campaign ", 10)
			assertFits(t, resize(t, NewApp(testCfg, campaign), size.width, size.height).View(), size.width, size.height)
		})
	}
}

func TestLayout_LauncherViewsFitTerminal(t *testing.T) {
	long := strings.Repeat("A drowned kingdom keeps its debts and its ghosts. ", 3)
	campaigns := make([]statedb.Campaign, 20)
	for i := range campaigns {
		campaigns[i] = statedb.Campaign{ID: dbutil.ToPgtype(uuid.New()), Name: strings.Repeat("Long Name ", 20)}
	}
	proposals := []world.CampaignProposal{{Name: long, Summary: long}, {Name: "Two", Summary: long}, {Name: "Three", Summary: long}}
	profile := world.CampaignProfile{Genre: "Fantasy", Tone: "Mysterious", Themes: []string{"debt", "memory"}}
	character := world.CharacterProfile{Name: "Mara", Concept: long, Background: long}

	views := []struct {
		name string
		view func() View
	}{
		{"picker", func() View { v := campaign.NewPicker(campaigns); return &v }},
		{"method", func() View { v := campaign.NewMethodModel(); return &v }},
		{"interview", func() View {
			v := campaign.NewInterviewModel(t.Context(), nil)
			return &v
		}},
		{"attributes", func() View { v := campaign.NewAttributesModel(); return &v }},
		{"attributes with error", func() View { v := campaign.NewAttributesModelWithState("", "", "", long); return &v }},
		{"proposals", func() View { v := campaign.NewProposalsModel(proposals); return &v }},
		{"character method", func() View { v := campaign.NewCharMethodModel(); return &v }},
		{"character interview", func() View {
			v := campaign.NewCharInterviewModel(t.Context(), nil, &profile)
			return &v
		}},
		{"character form", func() View { v := campaign.NewCharFormModel(); return &v }},
		{"confirmation", func() View {
			v := campaign.NewConfirmationModel("The Drowned Crown", long, profile, character)
			return &v
		}},
		{"world build", func() View {
			v := campaign.NewWorldBuildModel(t.Context(), nil, world.OrchestratorInput{})
			return &v
		}},
	}

	for _, size := range layoutSizes {
		newLauncher := func(t *testing.T) Launcher {
			return resize(t, NewLauncher(testCfg, t.Context(), nil), size.width, size.height)
		}
		for _, tc := range views {
			t.Run(fmt.Sprintf("%dx%d/%s", size.width, size.height, tc.name), func(t *testing.T) {
				l, _ := newLauncher(t).activateView(launcherSelecting, tc.view())
				assertFits(t, l.View(), size.width, size.height)

				// A launcher error takes rows above the view. With no LLM provider,
				// ConfirmedMsg reports one and leaves the view in place.
				next, _ := l.Update(campaign.ConfirmedMsg{})
				withError := next.(Launcher)
				if !strings.Contains(withError.View(), "LLM provider unavailable") {
					t.Fatalf("expected the launcher error above the view:\n%s", withError.View())
				}
				assertFits(t, withError.View(), size.width, size.height)
			})
		}

		t.Run(fmt.Sprintf("%dx%d/bootstrap error", size.width, size.height), func(t *testing.T) {
			next, _ := newLauncher(t).Update(bootstrapDoneMsg{err: errors.New(long)})
			assertFits(t, next.(Launcher).View(), size.width, size.height)
		})

		t.Run(fmt.Sprintf("%dx%d/app after launcher", size.width, size.height), func(t *testing.T) {
			// Bubble Tea sends no new WindowSizeMsg when the launcher hands over.
			next, _ := newLauncher(t).Update(campaignLoadedMsg{c: testCampaign})
			app, ok := next.(App)
			if !ok {
				t.Fatalf("expected App, got %T", next)
			}
			view := app.View()
			assertFits(t, view, size.width, size.height)
			if got := lipgloss.Width(view); got != size.width {
				t.Errorf("app after launcher is %d columns wide, want %d", got, size.width)
			}
		})
	}
}

// longField returns about 700 characters ending in marker, so a test can tell
// whether the end of the text can be reached.
func longField(marker string) string {
	return strings.Repeat("A drowned kingdom keeps its debts and its ghosts. ", 14) + marker
}

// assertScrollsTo pages down through l's view. Every page must fit the terminal
// and keep onEveryPage on screen, and each marker must appear on some page.
func assertScrollsTo(t *testing.T, l Launcher, width, height int, onEveryPage, markers []string) {
	t.Helper()
	missing := slices.Clone(markers)
	previous := ""
	for range 100 {
		view := l.View()
		if view == previous {
			break
		}
		previous = view
		assertFits(t, view, width, height)
		for _, want := range onEveryPage {
			if !strings.Contains(view, want) {
				t.Fatalf("%q is not on screen:\n%s", want, view)
			}
		}
		missing = slices.DeleteFunc(missing, func(m string) bool { return strings.Contains(view, m) })
		next, _ := l.Update(tea.KeyMsg{Type: tea.KeyPgDown})
		l = next.(Launcher)
	}
	if len(missing) > 0 {
		t.Errorf("paging down never showed %v; last page:\n%s", missing, previous)
	}
}

func TestLayout_LongLauncherTextScrolls(t *testing.T) {
	for _, size := range layoutSizes[:2] {
		newLauncher := func(t *testing.T) Launcher {
			return resize(t, NewLauncher(testCfg, t.Context(), nil), size.width, size.height)
		}

		t.Run(fmt.Sprintf("%dx%d/confirmation", size.width, size.height), func(t *testing.T) {
			v := campaign.NewConfirmationModel("The Drowned Crown", longField("SUMMARYEND"),
				world.CampaignProfile{Genre: "Fantasy", Tone: "Mysterious"},
				world.CharacterProfile{Name: "Mara", Concept: longField("CONCEPTEND"), Background: longField("BACKGROUNDEND")})
			l, _ := newLauncher(t).activateView(launcherConfirmation, &v)
			assertScrollsTo(t, l, size.width, size.height,
				[]string{"Begin Adventure", "Go Back"},
				[]string{"The Drowned Crown", "SUMMARYEND", "CONCEPTEND", "BACKGROUNDEND"})
		})

		t.Run(fmt.Sprintf("%dx%d/attributes after failed proposals", size.width, size.height), func(t *testing.T) {
			err := errors.New(longField("") + longField("ERROREND"))
			next, _ := newLauncher(t).Update(proposalsGeneratedMsg{err: err})
			assertScrollsTo(t, next.(Launcher), size.width, size.height,
				[]string{"Campaign Attributes", "Genre", "Fantasy", "Steampunk"},
				[]string{"Generate proposals failed", "ERROREND"})
		})
	}
}
