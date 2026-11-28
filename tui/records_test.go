package tui

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/tui/narrative"
)

// recordsMock is an engine that also reads the player's saved campaign data.
type recordsMock struct {
	*mockGameEngine
	mu      sync.Mutex
	records *game.PlayerRecords
	err     error
	asked   []uuid.UUID
}

func (m *recordsMock) PlayerRecords(_ context.Context, id uuid.UUID) (*game.PlayerRecords, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.asked = append(m.asked, id)
	return m.records, m.err
}

func (m *recordsMock) set(records *game.PlayerRecords, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records, m.err = records, err
}

// campaignRecords is a campaign whose character, item and quest names all hold
// marker, so a test can tell one saved state from another.
func campaignRecords(marker string) *game.PlayerRecords {
	return &game.PlayerRecords{
		Character: &game.CharacterRecord{
			PlayerCharacter: domain.PlayerCharacter{Name: "Mara " + marker, Level: 1, HP: 7, MaxHP: 10, Status: "active"},
			RulesVersion:    "legacy",
		},
		Inventory: []domain.Item{{Name: "Lantern " + marker, ItemType: domain.ItemTypeMisc, Rarity: "common", Quantity: 1}},
		Quests: []game.QuestRecord{{
			Quest:      domain.Quest{Title: "Debt " + marker, Status: domain.QuestStatusActive},
			Objectives: []domain.QuestObjective{{Description: "Find the ledger " + marker}},
		}},
	}
}

// recordReads runs cmd, and what it batches, side by side as Bubble Tea does
// and returns the first n results of reading the campaign data.
func recordReads(t *testing.T, cmd tea.Cmd, n int) []tea.Msg {
	t.Helper()
	messages := make(chan tea.Msg, 16)
	var run func(tea.Cmd)
	run = func(cmd tea.Cmd) {
		if cmd != nil {
			go func() { messages <- cmd() }()
		}
	}
	run(cmd)
	var reads []tea.Msg
	deadline := time.After(5 * time.Second)
	for len(reads) < n {
		select {
		case msg := <-messages:
			switch msg := msg.(type) {
			case tea.BatchMsg:
				for _, c := range msg {
					run(c)
				}
			case recordsLoadedMsg:
				reads = append(reads, msg)
			}
		case <-deadline:
			t.Fatalf("got %d reads of the campaign data, want %d", len(reads), n)
		}
	}
	return reads
}

// reread runs cmd and gives the app the campaign data it reads.
func reread(t *testing.T, app tea.Model, cmd tea.Cmd) App {
	t.Helper()
	next, _ := app.Update(recordReads(t, cmd, 1)[0])
	return next.(App)
}

// startedRecordsApp is an app that has run its start-up commands.
func startedRecordsApp(t *testing.T, mock *recordsMock) App {
	t.Helper()
	app := resize(t, testAppWithEngine(mock), 100, 48)
	return reread(t, app, app.Init())
}

// tabViews opens the three data tabs in turn and returns what each shows.
func tabViews(app App) map[ViewState]string {
	views := map[ViewState]string{}
	for _, state := range []ViewState{ViewCharacterSheet, ViewInventory, ViewQuestLog} {
		app.router.GoToTab(int(state))
		views[state] = app.View()
	}
	app.router.GoToTab(int(ViewNarrative))
	return views
}

func assertTabs(t *testing.T, app App, want map[ViewState][]string) {
	t.Helper()
	views := tabViews(app)
	for state, texts := range want {
		for _, text := range texts {
			if !strings.Contains(views[state], text) {
				t.Errorf("tab %d does not show %q:\n%s", state, text, views[state])
			}
		}
	}
	for state, view := range views {
		for _, sample := range []string{"Missing Merchant", "innkeeper", "future update"} {
			if strings.Contains(view, sample) {
				t.Errorf("tab %d shows sample content %q:\n%s", state, sample, view)
			}
		}
	}
}

func shows(marker string) map[ViewState][]string {
	return map[ViewState][]string{
		ViewCharacterSheet: {"Mara " + marker, "HP 7/10"},
		ViewInventory:      {"Lantern " + marker},
		ViewQuestLog:       {"Debt " + marker, "Find the ledger " + marker},
	}
}

func TestTabsShowTheLoadedCampaignsData(t *testing.T) {
	mock := &recordsMock{mockGameEngine: &mockGameEngine{}, records: campaignRecords("SAVED")}
	launcher := resize(t, newTestLauncherWithEngine(mock.mockGameEngine), 100, 48)
	launcher.engine = mock
	c := statedb.Campaign{ID: dbutil.ToPgtype(uuid.New()), Name: "Ashfall"}
	next, cmd := launcher.Update(campaignLoadedMsg{c: c})
	app := reread(t, next, cmd)

	assertTabs(t, app, shows("SAVED"))
	if len(mock.asked) != 1 || mock.asked[0] != dbutil.FromPgtype(c.ID) {
		t.Errorf("read campaigns %v, want only %v", mock.asked, dbutil.FromPgtype(c.ID))
	}
}

func TestTabsSayWhatACampaignDoesNotHaveYet(t *testing.T) {
	mock := &recordsMock{mockGameEngine: &mockGameEngine{}, records: &game.PlayerRecords{}}
	assertTabs(t, startedRecordsApp(t, mock), map[ViewState][]string{
		ViewCharacterSheet: {"No character yet"},
		ViewInventory:      {"No items yet"},
		ViewQuestLog:       {"No quests yet"},
	})
}

func TestTabsBeforeTheDataArrives(t *testing.T) {
	mock := &recordsMock{mockGameEngine: &mockGameEngine{}, records: campaignRecords("SAVED")}
	app := resize(t, testAppWithEngine(mock), 100, 48)
	assertTabs(t, app, map[ViewState][]string{
		ViewCharacterSheet: {"Loading"},
		ViewInventory:      {"Loading"},
		ViewQuestLog:       {"Loading"},
	})
	for state, view := range tabViews(app) {
		if strings.Contains(view, "yet") {
			t.Errorf("tab %d says the campaign is empty before it was read:\n%s", state, view)
		}
	}
}

func TestTabsRefreshWhenATurnOrConfirmedActionEnds(t *testing.T) {
	failed := errors.New("delivery uncertain")
	done := &engine.TurnResult{Narrative: "Done."}
	for name, msg := range map[string]tea.Msg{
		"turn":                         turnProcessedMsg{result: done},
		"turn without narrative":       turnProcessedMsg{result: &engine.TurnResult{}},
		"failed turn":                  turnProcessedMsg{err: failed},
		"action surge":                 surgeConfirmedMsg{result: done},
		"tactical mind":                tacticalConfirmedMsg{result: done},
		"shield of faith":              faithConfirmedMsg{result: done},
		"light":                        lightConfirmedMsg{result: done},
		"detect magic":                 detectMagicConfirmedMsg{result: done},
		"fall":                         fallConfirmedMsg{result: done},
		"death save":                   deathSaveConfirmedMsg{result: done},
		"exploration":                  explorationConfirmedMsg{control: doorControl, result: done},
		"action with uncertain result": lightConfirmedMsg{err: failed},
	} {
		t.Run(name, func(t *testing.T) {
			mock := &recordsMock{mockGameEngine: &mockGameEngine{}, records: campaignRecords("BEFORE")}
			app := startedRecordsApp(t, mock)
			mock.set(campaignRecords("AFTER"), nil)

			next, cmd := app.Update(msg)
			assertTabs(t, reread(t, next, cmd), shows("AFTER"))
		})
	}
}

func TestPreviewsAndTypingDoNotRereadTheData(t *testing.T) {
	mock := &recordsMock{mockGameEngine: &mockGameEngine{}, records: campaignRecords("SAVED")}
	var m tea.Model = startedRecordsApp(t, mock)
	for _, msg := range []tea.Msg{
		typed("a"), key(tea.KeyTab), key(tea.KeyTab), key(tea.KeyPgDown), key(tea.KeyEsc),
		narrative.SubmitMsg{Input: "/light"}, lightPreparedMsg{}, narrativeChunkMsg{chunk: "x"}, narrativeStreamDoneMsg{},
	} {
		m, _ = m.Update(msg)
	}
	if len(mock.asked) != 1 {
		t.Errorf("read the campaign data %d times, want once", len(mock.asked))
	}
}

func TestAnOlderReadDoesNotReplaceANewerOne(t *testing.T) {
	mock := &recordsMock{mockGameEngine: &mockGameEngine{}, records: campaignRecords("OLDER")}
	app := resize(t, testAppWithEngine(mock), 100, 48)
	older := recordReads(t, app.Init(), 1)[0]

	mock.set(campaignRecords("NEWER"), nil)
	next, cmd := app.Update(turnProcessedMsg{result: &engine.TurnResult{}})
	after := reread(t, next, cmd)
	next, _ = after.Update(older)

	assertTabs(t, next.(App), shows("NEWER"))
}

func TestAFailedReadIsSaidAndRetriedWhenATabOpens(t *testing.T) {
	mock := &recordsMock{mockGameEngine: &mockGameEngine{}, err: errors.New("dial tcp: password=hunter2")}
	app := startedRecordsApp(t, mock)
	assertTabs(t, app, map[ViewState][]string{
		ViewCharacterSheet: {"Could not load"},
		ViewInventory:      {"Could not load"},
		ViewQuestLog:       {"Could not load"},
	})
	for state, view := range tabViews(app) {
		if strings.Contains(view, "hunter2") || strings.Contains(view, "yet") {
			t.Errorf("tab %d shows the raw error or an empty campaign:\n%s", state, view)
		}
	}

	mock.set(campaignRecords("SAVED"), nil)
	next, cmd := app.Update(key(tea.KeyTab))
	retried := reread(t, next, cmd)
	assertTabs(t, retried, shows("SAVED"))

	// A failed refresh keeps what was read before and says it may be old.
	mock.set(nil, errors.New("connection reset"))
	next, cmd = retried.Update(turnProcessedMsg{result: &engine.TurnResult{}})
	stale := reread(t, next, cmd)
	want := shows("SAVED")
	for state := range want {
		want[state] = append(want[state], "Could not refresh")
	}
	assertTabs(t, stale, want)

	mock.set(campaignRecords("AGAIN"), nil)
	next, cmd = stale.Update(key(tea.KeyTab))
	fresh := reread(t, next, cmd)
	assertTabs(t, fresh, shows("AGAIN"))
	for state, view := range tabViews(fresh) {
		if strings.Contains(view, "Could not") {
			t.Errorf("tab %d still reports the failure:\n%s", state, view)
		}
	}
}

func TestAFailedSheetLoadIsRetriedWhenATabOpens(t *testing.T) {
	records := campaignRecords("SAVED")
	records.Character.RulesVersion = "srd-5.2.1"
	records.Character.SheetError = fmt.Errorf("%w: connection reset", game.ErrSheetLoadFailed)
	mock := &recordsMock{mockGameEngine: &mockGameEngine{}, records: records}
	app := startedRecordsApp(t, mock)
	assertTabs(t, app, map[ViewState][]string{
		ViewCharacterSheet: {"Mara SAVED", "could not be loaded"},
		ViewQuestLog:       {"Debt SAVED"},
	})

	mock.set(campaignRecords("AGAIN"), nil)
	next, cmd := app.Update(key(tea.KeyTab))
	assertTabs(t, reread(t, next, cmd), shows("AGAIN"))
}

// An unsupported build stays unsupported, so opening a tab does not read again.
func TestAnUnsupportedBuildIsNotReadAgain(t *testing.T) {
	records := campaignRecords("SAVED")
	records.Character.SheetError = errors.New("six explicit ability scores required")
	mock := &recordsMock{mockGameEngine: &mockGameEngine{}, records: records}
	app := startedRecordsApp(t, mock)
	if _, cmd := app.Update(key(tea.KeyTab)); cmd != nil {
		t.Error("opening a tab started a command")
	}
}

// A real program runs Init, so the tabs fill without any other message.
func TestTeatest_TabsFillAtStartUp(t *testing.T) {
	records := campaignRecords("SAVED")
	records.Character.Stats = json.RawMessage(`{"concept":"ledger keeper"}`)
	mock := &recordsMock{mockGameEngine: &mockGameEngine{}, records: records}
	tm := teatest.NewTestModel(t, testAppWithEngine(mock), teatest.WithInitialTermSize(100, 48))
	tm.Send(key(tea.KeyTab))
	teatest.WaitFor(t, tm.Output(), func(out []byte) bool {
		return bytes.Contains(out, []byte("ledger keeper"))
	}, waitDuration)
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, finalTimeout)
}
