package tui

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"git.subcult.tv/subculture-collective/edda/tui/narrative"
)

type explorationEngine interface {
	PreviewDoor(context.Context, uuid.UUID, int) ([]srd521.SheetDoor, *game.ActionPreview, error)
	PreviewTravel(context.Context, uuid.UUID, string) (*game.ActionPreview, error)
	PreviewStand(context.Context, uuid.UUID) (*game.ActionPreview, error)
	ConfirmAction(context.Context, uuid.UUID, engine.PreparedAction) (*engine.TurnResult, error)
}

// The three exploration controls: opening a known door, one hour of travel and
// standing from a saved Prone condition. Each keeps its own pending request.
const (
	doorControl   = "door"
	travelControl = "travel"
	standControl  = "stand"
)

var explorationUsage = map[string]string{
	doorControl:   "Use /door to list the doors you know here, /door <number> to preview opening one; then /confirm-door or /cancel-door.",
	travelControl: "Use /travel slow, /travel normal or /travel fast to list destinations, then /travel <pace> <number> to preview one hour of travel; then /confirm-travel or /cancel-travel.",
	standControl:  "Use /stand to preview standing from Prone outside an encounter; then /confirm-stand or /cancel-stand.",
}

var explorationLabel = map[string]string{doorControl: "Opening the door", travelControl: "Travel", standControl: "Standing up"}

// explorationPreparedMsg answers one preview command. number is the 1-based
// door or destination chosen; zero lists the choices.
type explorationPreparedMsg struct {
	control string
	number  int
	doors   []srd521.SheetDoor
	preview *game.ActionPreview
	err     error
}

type explorationConfirmedMsg struct {
	control string
	result  *engine.TurnResult
	err     error
}

func (a *App) explorationRequest(control string) **engine.PreparedAction {
	switch control {
	case doorControl:
		return &a.doorRequest
	case travelControl:
		return &a.travelRequest
	}
	return &a.standRequest
}

// updateExploration handles the door, travel and stand commands and their
// results. It reports whether the message belonged to them.
func (a *App) updateExploration(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case explorationPreparedMsg:
		a.turnBusy = false
		a.explorationPrepared(msg)
		return nil, true
	case explorationConfirmedMsg:
		a.turnBusy = false
		pending := a.explorationRequest(msg.control)
		text := "Delivery uncertain; /confirm-" + msg.control + " retries the same action ID."
		if errors.Is(msg.err, engine.ErrActionRefresh) {
			*pending = nil
			text = "Action changed; prepare it again with /" + msg.control + "."
		}
		if msg.err == nil && msg.result != nil {
			*pending = nil
			text = msg.result.Narrative
		}
		a.systemSay(text)
		return nil, true
	case narrative.SubmitMsg:
		fields := strings.Fields(msg.Input)
		if len(fields) == 0 {
			return nil, false
		}
		for _, control := range []string{doorControl, travelControl, standControl} {
			if slices.Contains([]string{"/" + control, "/confirm-" + control, "/cancel-" + control}, fields[0]) {
				return a.submitExploration(control, fields[0], fields[1:]), true
			}
		}
	}
	return nil, false
}

// submitExploration accepts no door state, distance or terrain: the only
// arguments name a pace and a number from the server's list.
func (a *App) submitExploration(control, command string, args []string) tea.Cmd {
	if a.turnBusy {
		return nil
	}
	se, ok := a.engine.(explorationEngine)
	if !ok {
		return nil
	}
	ctx, campaign, usage := a.ctx, dbutil.FromPgtype(a.campaign.ID), explorationUsage[control]
	pending := a.explorationRequest(control)
	switch command {
	case "/cancel-" + control:
		if len(args) != 0 {
			a.systemSay(usage)
			return nil
		}
		*pending = nil
		a.systemSay(explorationLabel[control] + " confirmation dismissed; nothing changed.")
		return nil
	case "/confirm-" + control:
		if len(args) != 0 || *pending == nil {
			a.systemSay(usage)
			return nil
		}
		request := **pending
		a.turnBusy = true
		return func() tea.Msg {
			result, err := se.ConfirmAction(ctx, campaign, request)
			return explorationConfirmedMsg{control, result, err}
		}
	}

	prepared := explorationPreparedMsg{control: control}
	pace, numbers := "", args
	if control == travelControl {
		if len(args) == 0 {
			a.systemSay(usage)
			return nil
		}
		pace, numbers = strings.ToLower(args[0]), args[1:]
		if !slices.Contains(srd521.TravelPaces, pace) {
			a.systemSay(usage)
			return nil
		}
	}
	if len(numbers) > 1 || (control == standControl && len(numbers) != 0) {
		a.systemSay(usage)
		return nil
	}
	if len(numbers) == 1 {
		number, err := strconv.Atoi(numbers[0])
		if err != nil || number < 1 {
			a.systemSay(usage)
			return nil
		}
		prepared.number = number
	}
	// Listing is read-only, so it stays open while a request is pending.
	if *pending != nil && (control == standControl || prepared.number > 0) {
		a.systemSay("Confirm/retry or cancel the pending " + control + " action first.")
		return nil
	}
	a.turnBusy = true
	return func() tea.Msg {
		switch control {
		case doorControl:
			prepared.doors, prepared.preview, prepared.err = se.PreviewDoor(ctx, campaign, prepared.number)
		case travelControl:
			prepared.preview, prepared.err = se.PreviewTravel(ctx, campaign, pace)
		default:
			prepared.preview, prepared.err = se.PreviewStand(ctx, campaign)
		}
		return prepared
	}
}

func (a *App) explorationPrepared(msg explorationPreparedMsg) {
	usage := explorationUsage[msg.control]
	if msg.err == nil && msg.control == doorControl && msg.preview == nil {
		a.systemSay(doorList(msg.doors, msg.number) + " " + usage)
		return
	}
	if msg.err != nil || msg.preview == nil {
		a.systemSay(explorationLabel[msg.control] + " is unavailable for this campaign's character.")
		return
	}
	preview := msg.preview
	var target uuid.UUID
	destination := ""
	if msg.control == travelControl {
		if msg.number == 0 {
			a.systemSay(destinationList(preview) + usage)
			return
		}
		if msg.number > len(preview.Targets) {
			a.systemSay("That destination number is not offered; list the destinations again with /travel <pace>.")
			return
		}
		chosen := preview.Targets[msg.number-1]
		if chosen.UnavailableReason != "" {
			a.systemSay(chosen.Label + ". Unavailable: " + chosen.UnavailableReason)
			return
		}
		target, destination = chosen.ID, " Destination: "+chosen.Label+"."
	}
	if !preview.Available {
		a.systemSay(preview.Reason + " " + usage)
		return
	}
	*a.explorationRequest(msg.control) = &engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: preview.ExpectedRevision, Intent: preview.Intent, TargetID: target}
	a.systemSay("Preview only: " + strings.Join(preview.Options, " ") + destination + " /confirm-" + msg.control + " to commit; /cancel-" + msg.control + " to dismiss.")
}

func doorList(doors []srd521.SheetDoor, number int) string {
	if len(doors) == 0 {
		return "You know no door here."
	}
	lines := []string{"Doors you know here:"}
	if number > len(doors) {
		lines = []string{"That door number is not offered. Doors you know here:"}
	}
	for i, door := range doors {
		state := ""
		if door.State != "" {
			state = " (" + door.State + ")"
		}
		lines = append(lines, fmt.Sprintf("%d. %s%s", i+1, door.Name, state))
	}
	return strings.Join(lines, "\n") + "\n"
}

func destinationList(preview *game.ActionPreview) string {
	lines := []string{}
	if !preview.Available {
		lines = append(lines, preview.Reason)
	}
	for i, target := range preview.Targets {
		line := fmt.Sprintf("%d. %s", i+1, target.Label)
		if target.UnavailableReason != "" {
			line += ". Unavailable: " + target.UnavailableReason
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}
