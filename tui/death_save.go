package tui

import (
	"context"
	"errors"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/tui/narrative"
)

type deathSaveEngine interface {
	PreviewDeathSave(context.Context, uuid.UUID) (*game.ActionPreview, error)
	ConfirmAction(context.Context, uuid.UUID, engine.PreparedAction) (*engine.TurnResult, error)
}

type deathSavePreparedMsg struct {
	preview *game.ActionPreview
	err     error
}

type deathSaveConfirmedMsg struct {
	result *engine.TurnResult
	err    error
}

const deathSaveUsage = "At 0 HP outside an encounter, use /death-save to preview your Death Saving Throw (or a Stable character's recovery); then /confirm-death-save or /cancel-death-save. In an encounter it is rolled at the start of your turn."

// updateDeathSave handles the Death Saving Throw commands and their results.
// It reports whether the message belonged to them.
func (a *App) updateDeathSave(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case deathSavePreparedMsg:
		a.turnBusy = false
		if msg.err != nil || msg.preview == nil {
			a.systemSay("The Death Saving Throw is unavailable for this campaign's character.")
			return nil, true
		}
		if !msg.preview.Available {
			a.systemSay(msg.preview.Reason + " " + deathSaveUsage)
			return nil, true
		}
		a.deathSaveRequest = &engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: msg.preview.ExpectedRevision, Intent: msg.preview.Intent}
		a.systemSay("Preview only: " + strings.Join(msg.preview.Options, " ") + " /confirm-death-save to roll; /cancel-death-save to dismiss.")
		return nil, true
	case deathSaveConfirmedMsg:
		a.turnBusy = false
		text := "Delivery uncertain; /confirm-death-save retries the same action ID."
		if errors.Is(msg.err, engine.ErrActionRefresh) {
			a.deathSaveRequest = nil
			text = "Action changed; prepare the Death Saving Throw again."
		}
		if msg.err == nil && msg.result != nil {
			a.deathSaveRequest = nil
			text = msg.result.Narrative
		}
		a.systemSay(text)
		return nil, true
	case narrative.SubmitMsg:
		fields := strings.Fields(msg.Input)
		if len(fields) == 0 || !slices.Contains([]string{"/death-save", "/confirm-death-save", "/cancel-death-save"}, fields[0]) {
			return nil, false
		}
		return a.submitDeathSave(fields), true
	}
	return nil, false
}

// submitDeathSave takes no arguments: the server holds the counts and rolls.
func (a *App) submitDeathSave(fields []string) tea.Cmd {
	if a.turnBusy {
		return nil
	}
	se, ok := a.engine.(deathSaveEngine)
	if !ok {
		return nil
	}
	if len(fields) != 1 {
		a.systemSay(deathSaveUsage)
		return nil
	}
	ctx, campaign := a.ctx, dbutil.FromPgtype(a.campaign.ID)
	switch fields[0] {
	case "/cancel-death-save":
		a.deathSaveRequest = nil
		a.systemSay("Death Saving Throw dismissed; nothing was rolled.")
		return nil
	case "/confirm-death-save":
		if a.deathSaveRequest == nil {
			a.systemSay(deathSaveUsage)
			return nil
		}
		request := *a.deathSaveRequest
		a.turnBusy = true
		return func() tea.Msg {
			result, err := se.ConfirmAction(ctx, campaign, request)
			return deathSaveConfirmedMsg{result, err}
		}
	}
	if a.deathSaveRequest != nil {
		a.systemSay("Confirm/retry or cancel the pending Death Saving Throw first.")
		return nil
	}
	a.turnBusy = true
	return func() tea.Msg {
		preview, err := se.PreviewDeathSave(ctx, campaign)
		return deathSavePreparedMsg{preview, err}
	}
}
