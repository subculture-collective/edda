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

type fallEngine interface {
	PreviewFall(context.Context, uuid.UUID, string) (*game.ActionPreview, error)
	ConfirmAction(context.Context, uuid.UUID, engine.PreparedAction) (*engine.TurnResult, error)
}

type fallPreparedMsg struct {
	preview *game.ActionPreview
	err     error
}

type fallConfirmedMsg struct {
	result *engine.TurnResult
	err    error
}

const fallUsage = "Use /fall to preview the declared fall, or /fall athletics or /fall acrobatics to use your Reaction when it lands in water; then /confirm-fall or /cancel-fall."

// updateFall handles the declared-fall commands and their results. It reports
// whether the message belonged to the fall.
func (a *App) updateFall(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case fallPreparedMsg:
		a.turnBusy = false
		a.fallPrepared(msg)
		return nil, true
	case fallConfirmedMsg:
		a.turnBusy = false
		text := "Delivery uncertain; /confirm-fall retries the same action ID."
		if errors.Is(msg.err, engine.ErrActionRefresh) {
			a.fallRequest = nil
			text = "Action changed; prepare the fall again."
		}
		if msg.err == nil && msg.result != nil {
			a.fallRequest = nil
			text = msg.result.Narrative
		}
		a.systemSay(text)
		return nil, true
	case narrative.SubmitMsg:
		fields := strings.Fields(msg.Input)
		if len(fields) == 0 || !slices.Contains([]string{"/fall", "/confirm-fall", "/cancel-fall"}, fields[0]) {
			return nil, false
		}
		return a.submitFall(fields), true
	}
	return nil, false
}

// submitFall accepts no height or damage: the only argument names the
// water-entry skill.
func (a *App) submitFall(fields []string) tea.Cmd {
	if a.turnBusy {
		return nil
	}
	se, ok := a.engine.(fallEngine)
	if !ok {
		return nil
	}
	ctx, campaign := a.ctx, dbutil.FromPgtype(a.campaign.ID)
	command, args := fields[0], fields[1:]
	skill := ""
	if command == "/fall" && len(args) == 1 {
		skill = strings.ToLower(args[0])
	}
	if len(args) != 0 && skill != "athletics" && skill != "acrobatics" {
		a.systemSay(fallUsage)
		return nil
	}
	switch command {
	case "/cancel-fall":
		a.fallRequest = nil
		a.systemSay("Fall confirmation dismissed. The fall stays declared.")
		return nil
	case "/confirm-fall":
		if a.fallRequest == nil {
			a.systemSay(fallUsage)
			return nil
		}
		request := *a.fallRequest
		a.turnBusy = true
		return func() tea.Msg {
			result, err := se.ConfirmAction(ctx, campaign, request)
			return fallConfirmedMsg{result, err}
		}
	}
	if a.fallRequest != nil {
		a.systemSay("Confirm/retry or cancel the pending fall first.")
		return nil
	}
	a.turnBusy = true
	return func() tea.Msg {
		preview, err := se.PreviewFall(ctx, campaign, skill)
		return fallPreparedMsg{preview, err}
	}
}

func (a *App) fallPrepared(msg fallPreparedMsg) {
	if msg.err != nil || msg.preview == nil {
		a.systemSay("The fall is unavailable for this campaign's character.")
		return
	}
	preview := msg.preview
	if !preview.Available {
		a.systemSay(preview.Reason + " " + fallUsage)
		return
	}
	a.fallRequest = &engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: preview.ExpectedRevision, Intent: preview.Intent}
	a.systemSay("Preview only: " + strings.Join(preview.Options, " ") + " /confirm-fall to commit; /cancel-fall to dismiss.")
}
