package tui

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/engine"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"git.subcult.tv/subculture-collective/edda/tui/narrative"
)

type detectMagicEngine interface {
	PreviewDetectMagic(context.Context, uuid.UUID, game.ActionIntent) (*game.StoredAction, error)
	ConfirmAction(context.Context, uuid.UUID, engine.PreparedAction) (*engine.TurnResult, error)
}

// detectMagicPreparedMsg answers one Detect Magic command.
type detectMagicPreparedMsg struct {
	stored *game.StoredAction
	status bool
	err    error
}

type detectMagicConfirmedMsg struct {
	result *engine.TurnResult
	err    error
}

const detectMagicUsage = "Use /detect-magic for active concentration, /detect-magic cast to preview the Magic Initiate free use, or /study-detect-magic to preview the study; then /confirm-detect-magic or /cancel-detect-magic."

// updateDetectMagic handles the Detect Magic commands and their results. It
// reports whether the message belonged to Detect Magic.
func (a *App) updateDetectMagic(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case detectMagicPreparedMsg:
		a.turnBusy = false
		a.detectMagicPrepared(msg)
		return nil, true
	case detectMagicConfirmedMsg:
		a.turnBusy = false
		text := "Delivery uncertain; /confirm-detect-magic retries the same action ID."
		if errors.Is(msg.err, engine.ErrActionRefresh) {
			a.detectMagicRequest = nil
			text = "Action changed; prepare Detect Magic again."
		}
		if msg.err == nil && msg.result != nil {
			a.detectMagicRequest = nil
			text = msg.result.Narrative
		}
		a.systemSay(text)
		return nil, true
	case narrative.SubmitMsg:
		fields := strings.Fields(msg.Input)
		if len(fields) == 0 || !slices.Contains([]string{"/detect-magic", "/study-detect-magic", "/confirm-detect-magic", "/cancel-detect-magic"}, fields[0]) {
			return nil, false
		}
		return a.submitDetectMagic(fields), true
	}
	return nil, false
}

func (a *App) submitDetectMagic(fields []string) tea.Cmd {
	if a.turnBusy {
		return nil
	}
	se, ok := a.engine.(detectMagicEngine)
	if !ok {
		return nil
	}
	ctx, campaign := a.ctx, dbutil.FromPgtype(a.campaign.ID)
	command, args := fields[0], fields[1:]
	cast := command == "/detect-magic" && len(args) == 1 && strings.EqualFold(args[0], "cast")
	if len(args) != 0 && !cast {
		a.systemSay(detectMagicUsage)
		return nil
	}
	switch command {
	case "/cancel-detect-magic":
		a.detectMagicRequest = nil
		a.systemSay("Detect Magic confirmation dismissed.")
		return nil
	case "/confirm-detect-magic":
		if a.detectMagicRequest == nil {
			a.systemSay(detectMagicUsage)
			return nil
		}
		request := *a.detectMagicRequest
		a.turnBusy = true
		return func() tea.Msg {
			result, err := se.ConfirmAction(ctx, campaign, request)
			return detectMagicConfirmedMsg{result, err}
		}
	}

	// The read-only study preview carries any live concentration.
	intent := game.ActionIntent{Kind: srd521.DetectMagicStudy}
	status := command == "/detect-magic" && !cast
	if cast {
		intent = game.ActionIntent{Kind: srd521.SpellCast, EntryID: srd521.DetectMagicID}
	}
	if a.detectMagicRequest != nil && !status {
		a.systemSay("Confirm/retry or cancel the pending Detect Magic action first.")
		return nil
	}
	a.turnBusy = true
	return func() tea.Msg {
		stored, err := se.PreviewDetectMagic(ctx, campaign, intent)
		return detectMagicPreparedMsg{stored, status, err}
	}
}

func (a *App) detectMagicPrepared(msg detectMagicPreparedMsg) {
	if msg.err != nil || msg.stored == nil || msg.stored.Preview == nil {
		a.systemSay("Detect Magic is unavailable for this campaign's character.")
		return
	}
	preview := msg.stored.Preview
	if msg.status {
		text := "No active concentration shown."
		if c := msg.stored.Concentration; c != nil {
			text = fmt.Sprintf("Active concentration: %s, expires %s.", srd521.ConcentrationLabel(c.Spell), game.FormatCampaignTime(c.ExpiresAt))
		}
		if !preview.Available {
			text += " " + preview.Reason
		}
		a.systemSay(text + " " + detectMagicUsage)
		return
	}
	if !preview.Available {
		a.systemSay(preview.Reason)
		return
	}
	request := engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: preview.ExpectedRevision, Intent: preview.Intent}
	// A cast targets the caster, the only offered target; a study has none.
	if len(preview.Targets) == 1 {
		request.TargetID = preview.Targets[0].ID
	}
	a.detectMagicRequest = &request
	a.systemSay("Preview only: " + strings.Join(preview.Options, " ") + " /confirm-detect-magic to commit; /cancel-detect-magic to dismiss.")
}
