package tui

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/bootstrap"
	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/llm"
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

// feed runs cmd and gives its messages, and those of the commands that follow,
// to m. Spinner ticks are dropped because each one schedules the next.
func feed(t *testing.T, m tea.Model, cmd tea.Cmd) tea.Model {
	t.Helper()
	queue := []tea.Cmd{cmd}
	for steps := 0; len(queue) > 0; steps++ {
		if steps > 100 {
			t.Fatal("commands did not settle")
		}
		next := queue[0]
		queue = queue[1:]
		if next == nil {
			continue
		}
		msg := next()
		if _, ok := msg.(spinner.TickMsg); ok || msg == nil {
			continue
		}
		// Bubble Tea does not export the message tea.Sequence returns.
		if v := reflect.ValueOf(msg); v.Kind() == reflect.Slice && v.Type().Elem() == reflect.TypeFor[tea.Cmd]() {
			for i := range v.Len() {
				queue = append(queue, v.Index(i).Interface().(tea.Cmd))
			}
			continue
		}
		var follow tea.Cmd
		m, follow = m.Update(msg)
		queue = append(queue, follow)
	}
	return m
}

// send gives msg to m and runs what it starts.
func send(t *testing.T, m tea.Model, msg tea.Msg) tea.Model {
	t.Helper()
	next, cmd := m.Update(msg)
	return feed(t, next, cmd)
}

func key(k tea.KeyType) tea.KeyMsg { return tea.KeyMsg{Type: k} }

func typed(text string) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(text)} }

var shortSizes = []struct{ width, height int }{{80, 24}, {80, 16}}

func TestLayout_LongLauncherErrorsAreCapped(t *testing.T) {
	err := errors.New("ERRORSTART " + longField("") + longField("") + longField("ERROREND"))
	user := bootstrapDoneMsg{result: bootstrap.Result{Campaigns: []statedb.Campaign{makeTestCampaign(1, "Saved Campaign")}}}

	failures := []struct {
		name  string
		msgs  []tea.Msg
		below []string // text from the view under the banner
	}{
		{"bootstrap", []tea.Msg{bootstrapDoneMsg{err: err}}, nil},
		{"world build", []tea.Msg{campaign.WorldErrorMsg{Err: err}, user}, []string{"Select Campaign", "Saved Campaign"}},
		{"load campaign", []tea.Msg{campaignLoadedMsg{err: err}, user}, []string{"Select Campaign", "Saved Campaign"}},
		{"campaign name", []tea.Msg{campaignNamedMsg{err: err}}, []string{"How would you like", "Describe Your World"}},
	}

	for _, size := range shortSizes {
		for _, tc := range failures {
			t.Run(fmt.Sprintf("%dx%d/%s", size.width, size.height, tc.name), func(t *testing.T) {
				var m tea.Model = resize(t, NewLauncher(testCfg, t.Context(), nil), size.width, size.height)
				for _, msg := range tc.msgs {
					// The commands would query the database, so they are not run.
					m, _ = m.Update(msg)
				}
				closed := m.View()
				assertFits(t, closed, size.width, size.height)
				for _, want := range append([]string{"ERRORSTART", "? for the full error"}, tc.below...) {
					if !strings.Contains(closed, want) {
						t.Fatalf("%q is not on screen:\n%s", want, closed)
					}
				}

				m, _ = m.Update(typed("?"))
				assertScrollsTo(t, m.(Launcher), size.width, size.height,
					[]string{"Esc: close"}, []string{"ERRORSTART", "ERROREND"})

				m, _ = m.Update(key(tea.KeyEsc))
				for _, closeKey := range []tea.KeyMsg{key(tea.KeyEsc), typed("?"), typed("q")} {
					again, _ := m.Update(typed("?"))
					again, cmd := again.Update(closeKey)
					if cmd != nil {
						t.Fatalf("%q did more than close the error", closeKey)
					}
					if got := again.View(); got != closed {
						t.Fatalf("%q did not return to the view:\n%s", closeKey, got)
					}
				}

				// A resize lays out the open error again.
				open, _ := m.Update(typed("?"))
				open, _ = open.Update(tea.WindowSizeMsg{Width: 60, Height: 12})
				assertScrollsTo(t, open.(Launcher), 60, 12, []string{"Esc: close"}, []string{"ERRORSTART", "ERROREND"})
			})
		}
	}

	t.Run("a short error leaves ? to the view", func(t *testing.T) {
		var m tea.Model = resize(t, NewLauncher(testCfg, t.Context(), nil), 80, 24)
		iv := campaign.NewInterviewModel(t.Context(), nil)
		l, _ := m.(Launcher).activateView(launcherInterviewing, &iv)
		m, _ = l.showError("LLM provider unavailable")
		if strings.Contains(m.View(), "full error") {
			t.Fatalf("a short error offers details:\n%s", m.View())
		}
		m, _ = m.Update(typed("why?"))
		if view := m.View(); !strings.Contains(view, "why?") {
			t.Fatalf("the input did not get the question mark:\n%s", view)
		}
	})
}

func TestLayout_AttributesFormFitsShortTerminals(t *testing.T) {
	steps := []struct {
		title   string
		options []string
	}{
		{"Genre", []string{"Fantasy", "Sci-Fi", "Horror", "Historical", "Modern", "Post-Apocalyptic", "Steampunk"}},
		{"Setting Style", []string{"Open Wilderness", "Urban Sprawl", "Island Archipelago", "Underground", "War-Torn Kingdom", "Peaceful Realm"}},
		{"Tone", []string{"Gritty and Dark", "Light-Hearted", "Epic and Grand", "Mysterious", "Humorous", "Tense and Suspenseful"}},
	}
	failed := proposalsGeneratedMsg{err: errors.New(longField("") + longField("ERROREND"))}

	states := []struct {
		name string
		open func(t *testing.T, m tea.Model) tea.Model
		want []string
	}{
		{"no error", func(t *testing.T, m tea.Model) tea.Model {
			return send(t, m, campaign.MethodChosenMsg{Method: campaign.MethodAttributes})
		}, nil},
		{"failed proposals", func(t *testing.T, m tea.Model) tea.Model {
			return send(t, m, failed)
		}, []string{"Generate proposals"}},
		{"failed proposals and a banner", func(t *testing.T, m tea.Model) tea.Model {
			// With no LLM provider, submitting the form again adds a banner.
			return send(t, send(t, m, failed), campaign.AttributesReadyMsg{})
		}, []string{"Generate proposals", "LLM provider unavailable"}},
	}

	for _, height := range []int{24, 16, 14, 12, 10, 8} {
		for _, state := range states {
			t.Run(fmt.Sprintf("80x%d/%s", height, state.name), func(t *testing.T) {
				var m tea.Model = resize(t, NewLauncher(testCfg, t.Context(), nil), 80, height)
				m = state.open(t, m)

				check := func(wants ...string) {
					t.Helper()
					view := m.View()
					assertFits(t, view, 80, height)
					for _, want := range append(wants, state.want...) {
						if !strings.Contains(view, want) {
							t.Fatalf("%q is not on screen:\n%s", want, view)
						}
					}
				}
				for _, step := range steps {
					for i, option := range step.options {
						if i > 0 {
							m = send(t, m, key(tea.KeyDown))
						}
						check(step.title, "> "+option)
					}
					// Growing and shrinking the terminal keeps the selection.
					last := "> " + step.options[len(step.options)-1]
					m = send(t, m, tea.WindowSizeMsg{Width: 80, Height: 30})
					for _, option := range step.options {
						if !strings.Contains(m.View(), option) {
							t.Fatalf("a tall terminal does not show %q:\n%s", option, m.View())
						}
					}
					m = send(t, m, tea.WindowSizeMsg{Width: 80, Height: height})
					check(step.title, last)
					for i := len(step.options) - 2; i >= 0; i-- {
						m = send(t, m, key(tea.KeyUp))
						check(step.title, "> "+step.options[i])
					}
					if step.title != "Tone" {
						m = send(t, m, key(tea.KeyEnter))
					}
				}
			})
		}
	}
}

// unwrapped joins view's lines without their borders and spaces, so that a test
// can find text that wrapping split.
func unwrapped(view string) string {
	var b strings.Builder
	for _, line := range strings.Split(view, "\n") {
		b.WriteString(strings.ReplaceAll(strings.Trim(line, "│║ "), " ", ""))
	}
	return b.String()
}

func TestLayout_LongLogLinesWrap(t *testing.T) {
	for _, size := range shortSizes {
		t.Run(fmt.Sprintf("%dx%d", size.width, size.height), func(t *testing.T) {
			buf := logging.NewRingBuffer(16)
			buf.Write(logging.LogEntry{
				Time: time.Now(), Level: slog.LevelError, Service: "engine",
				Message: strings.Repeat("the model refused the turn ", 5) + "MESSAGEEND",
				Fields:  []logging.Field{{Key: "url", Value: "http://model.invalid/" + strings.Repeat("path/", 24) + "FIELDEND"}},
			})
			app := resize(t, NewAppWithEngine(testCfg, testCampaign, t.Context(), nil, buf), size.width, size.height)
			app.router.GoToTab(int(ViewLogs))

			assertFits(t, app.View(), size.width, size.height)

			// Sizing the panel alone checks that a narrower width wraps again.
			rows := lipgloss.Height(app.router.View())
			for _, width := range []int{size.width, 60, size.width} {
				app.router.SetSize(width, rows)
				view := app.router.View()
				assertFits(t, view, width, rows)
				for _, want := range []string{"modelrefused", "MESSAGEEND", "url=http://model.invalid/", "FIELDEND"} {
					if !strings.Contains(unwrapped(view), want) {
						t.Errorf("at %d columns %q is not on screen:\n%s", width, want, view)
					}
				}
			}
		})
	}
}

// replyProvider answers every interview turn with its next reply.
type replyProvider struct {
	stubLLMProvider
	replies *[]string
}

func (p replyProvider) Complete(context.Context, []llm.Message, []llm.Tool) (*llm.Response, error) {
	reply := (*p.replies)[0]
	*p.replies = (*p.replies)[1:]
	return &llm.Response{Content: reply}, nil
}

func TestLayout_LongInterviewReplyStartsAtItsTop(t *testing.T) {
	reply := func(name string) string {
		return name + "START " + strings.Repeat("Tell me which debts the drowned kingdom still keeps. ", 50) + name + "END"
	}
	interviews := []struct {
		name  string
		start tea.Msg
	}{
		{"campaign", campaign.MethodChosenMsg{Method: campaign.MethodDescribe}},
		{"character", campaign.CharMethodChosenMsg{Method: campaign.MethodDescribe}},
	}

	for _, size := range shortSizes {
		for _, iv := range interviews {
			t.Run(fmt.Sprintf("%dx%d/%s", size.width, size.height, iv.name), func(t *testing.T) {
				replies := []string{reply("FIRST"), reply("SECOND"), "SHORTREPLY"}
				launcher := newTestLauncherWithProvider(replyProvider{replies: &replies})
				profile := testCampaignProfile()
				launcher.campaignProfile = &profile
				var m tea.Model = resize(t, launcher, size.width, size.height)
				m = send(t, m, iv.start)

				check := func(step string, shown, hidden []string) {
					t.Helper()
					view := m.View()
					assertFits(t, view, size.width, size.height)
					for _, want := range shown {
						if !strings.Contains(view, want) {
							t.Fatalf("%s: %q is not on screen:\n%s", step, want, view)
						}
					}
					for _, unwanted := range hidden {
						if strings.Contains(view, unwanted) {
							t.Fatalf("%s: %q is on screen:\n%s", step, unwanted, view)
						}
					}
				}
				// pageTo pages with k until marker shows.
				pageTo := func(k tea.KeyType, marker string) {
					t.Helper()
					for range 20 {
						if strings.Contains(m.View(), marker) {
							return
						}
						m, _ = m.Update(key(k))
					}
					t.Fatalf("paging never showed %q:\n%s", marker, m.View())
				}

				check("first reply", []string{"GM: FIRSTSTART", "more below"}, []string{"FIRSTEND"})
				m, _ = m.Update(typed("a"))
				check("after typing", []string{"GM: FIRSTSTART"}, []string{"FIRSTEND"})
				pageTo(tea.KeyPgDown, "FIRSTEND")
				m, _ = m.Update(typed("b"))
				check("after paging down", []string{"FIRSTEND"}, []string{"more below", "FIRSTSTART"})
				pageTo(tea.KeyPgUp, "FIRSTSTART")
				pageTo(tea.KeyPgDown, "FIRSTEND")

				m = send(t, m, key(tea.KeyEnter))
				check("second reply", []string{"GM: SECONDSTART", "more below"}, []string{"SECONDEND"})
				// A narrower terminal wraps the history again and keeps its place.
				m, _ = m.Update(tea.WindowSizeMsg{Width: 60, Height: size.height})
				if view := m.View(); !strings.Contains(view, "GM: SECONDSTART") {
					t.Fatalf("resizing lost the start of the reply:\n%s", view)
				}
				m, _ = m.Update(tea.WindowSizeMsg{Width: size.width, Height: size.height})
				check("after resizing back", []string{"GM: SECONDSTART"}, []string{"SECONDEND"})

				m, _ = m.Update(typed("ok"))
				m = send(t, m, key(tea.KeyEnter))
				check("short reply", []string{"SECONDEND", "SHORTREPLY"}, []string{"more below"})
			})
		}
	}
}
