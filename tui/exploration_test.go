package tui

import (
	"context"
	"slices"
	"strings"
	"testing"

	"git.subcult.tv/subculture-collective/edda/internal/combat"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"git.subcult.tv/subculture-collective/edda/tui/narrative"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

type explorationMock struct {
	tacticalMock
	wood, peak uuid.UUID
	prone      bool
	previews   []string
}

func (m *explorationMock) PreviewDoor(_ context.Context, _ uuid.UUID, number int) ([]srd521.SheetDoor, *game.ActionPreview, error) {
	m.previews = append(m.previews, "door")
	doors := []srd521.SheetDoor{{ID: "d1", Name: "Oak door", State: "closed"}, {ID: "d2", Name: "Iron door"}}
	switch number {
	case 1:
		return doors, &game.ActionPreview{Available: true, ExpectedRevision: "door", Intent: game.ActionIntent{Kind: srd521.OpenDoor, EntryID: "d1"}, Options: []string{"Open Oak door."}}, nil
	case 2:
		return doors, &game.ActionPreview{Reason: "Locked doors are not supported."}, nil
	}
	return doors, nil, nil
}

func (m *explorationMock) PreviewTravel(_ context.Context, _ uuid.UUID, pace string) (*game.ActionPreview, error) {
	m.previews = append(m.previews, pace)
	return &game.ActionPreview{Available: true, ExpectedRevision: "travel", Intent: game.ActionIntent{Kind: srd521.TravelHour, EntryID: pace}, Options: []string{"Travel one segment."}, Targets: []game.ActionTarget{{ID: m.wood, Label: "Wood: 1.00 miles"}, {ID: m.peak, Label: "Peak: 2.00 miles", UnavailableReason: "mountain terrain allows at most a slow pace"}}}, nil
}

func (m *explorationMock) PreviewStand(context.Context, uuid.UUID) (*game.ActionPreview, error) {
	m.previews = append(m.previews, "stand")
	if !m.prone {
		return &game.ActionPreview{Reason: "No supported encounter is active."}, nil
	}
	return &game.ActionPreview{Available: true, ExpectedRevision: "stand", Intent: game.ActionIntent{Kind: combat.EncounterStand}, Options: []string{"Stand up from Prone."}}, nil
}

type explorationTUI struct {
	t *testing.T
	a App
}

func newExplorationTUI(t *testing.T, m *explorationMock) *explorationTUI {
	a := NewAppWithEngine(testCfg, testCampaign, context.Background(), m, nil)
	model, _ := a.Update(tea.WindowSizeMsg{Width: 200, Height: 60})
	return &explorationTUI{t: t, a: model.(App)}
}

func (u *explorationTUI) send(input string) {
	u.t.Helper()
	model, cmd := u.a.Update(narrative.SubmitMsg{Input: input})
	u.a = model.(App)
	if cmd != nil {
		model, _ = u.a.Update(cmd())
		u.a = model.(App)
	}
}

func (u *explorationTUI) shows(want string) {
	u.t.Helper()
	if !strings.Contains(u.a.narrativeView().View(), want) {
		u.t.Fatalf("missing %q:\n%s", want, u.a.narrativeView().View())
	}
}

func TestTUIDoorListPreviewCancelConfirmAndRetry(t *testing.T) {
	m := &explorationMock{}
	m.fail = true
	u := newExplorationTUI(t, m)
	// A door fact is never an argument.
	for _, bad := range []string{"/confirm-door", "/door open", "/door 0", "/door 1 open", "/cancel-door now"} {
		u.send(bad)
	}
	if u.a.doorRequest != nil || len(m.previews) != 0 || len(m.requests) != 0 {
		t.Fatalf("bad input reached the engine: %v", m.previews)
	}
	u.send("/door")
	u.shows("1. Oak door (closed)")
	u.shows("2. Iron door")
	u.send("/door 3")
	u.shows("That door number is not offered.")
	u.send("/door 2")
	u.shows("Locked doors are not supported.")
	if u.a.doorRequest != nil {
		t.Fatal("refused door pending")
	}
	u.send("/door 1")
	u.shows("Preview only: Open Oak door.")
	u.send("/door 2")
	u.shows("Confirm/retry or cancel the pending door action first.")
	u.send("/cancel-door")
	if u.a.doorRequest != nil || len(m.requests) != 0 {
		t.Fatal("cancel committed")
	}
	u.send("/door 1")
	u.send("/confirm-door")
	if u.a.doorRequest == nil || len(m.requests) != 1 {
		t.Fatal("uncertain request lost")
	}
	m.fail = false
	u.send("/confirm-door")
	if len(m.requests) != 2 || m.requests[0] != m.requests[1] || m.requests[0].Intent.EntryID != "d1" || m.requests[0].TargetID != uuid.Nil || u.a.doorRequest != nil {
		t.Fatalf("retry identity: %+v", m.requests)
	}
}

func TestTUITravelPaceDestinationAndConfirm(t *testing.T) {
	m := &explorationMock{wood: uuid.New(), peak: uuid.New()}
	u := newExplorationTUI(t, m)
	// No distance or terrain is accepted.
	for _, bad := range []string{"/confirm-travel", "/travel", "/travel gallop", "/travel normal forest", "/travel normal 1 5280"} {
		u.send(bad)
	}
	if u.a.travelRequest != nil || len(m.previews) != 0 {
		t.Fatalf("bad input reached the engine: %v", m.previews)
	}
	u.send("/travel Normal")
	u.shows("1. Wood: 1.00 miles")
	u.shows("2. Peak: 2.00 miles. Unavailable: mountain terrain allows at most a slow pace")
	u.send("/travel normal 3")
	u.shows("That destination number is not offered")
	u.send("/travel normal 2")
	if u.a.travelRequest != nil {
		t.Fatal("unavailable destination pending")
	}
	u.send("/travel normal 1")
	u.shows("Destination: Wood: 1.00 miles. /confirm-travel to commit")
	if r := u.a.travelRequest; r == nil || r.TargetID != m.wood || r.Intent.EntryID != srd521.PaceNormal || len(m.requests) != 0 {
		t.Fatalf("travel preview: %+v", r)
	}
	u.send("/confirm-travel")
	if len(m.requests) != 1 || m.requests[0].TargetID != m.wood || u.a.travelRequest != nil {
		t.Fatalf("travel confirm: %+v", m.requests)
	}
	if !slices.Equal(m.previews, []string{"normal", "normal", "normal", "normal"}) {
		t.Fatalf("previewed paces: %v", m.previews)
	}
}

func TestTUIStandPreviewAndConfirm(t *testing.T) {
	m := &explorationMock{}
	u := newExplorationTUI(t, m)
	u.send("/confirm-stand")
	u.send("/stand now")
	u.send("/stand")
	u.shows("No supported encounter is active.")
	if u.a.standRequest != nil || len(m.previews) != 1 {
		t.Fatalf("stand prepared: %v", m.previews)
	}
	m.prone = true
	u.send("/stand")
	u.shows("Preview only: Stand up from Prone.")
	u.send("/confirm-stand")
	if len(m.requests) != 1 || m.requests[0].Intent.Kind != combat.EncounterStand || m.requests[0].TargetID != uuid.Nil || u.a.standRequest != nil {
		t.Fatalf("stand confirm: %+v", m.requests)
	}
}
