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

type lightEngine interface {
	PreviewLight(context.Context, uuid.UUID, game.ActionIntent) (*game.StoredAction, error)
	ConfirmAction(context.Context, uuid.UUID, engine.PreparedAction) (*engine.TurnResult, error)
}

// lightPreparedMsg answers one /light command. object is the 1-based carried
// object chosen for a cast; zero lists the objects or reports status.
type lightPreparedMsg struct {
	stored *game.StoredAction
	status bool
	object int
	err    error
}

type lightConfirmedMsg struct {
	result *engine.TurnResult
	err    error
}

const lightUsage = "Use /light for status, /light <color> to list carried objects, then /light <color> <number>; /cover-light, /uncover-light or /dismiss-light for the active Light; then /confirm-light or /cancel-light."

// updateLight handles the Light commands and their results. It reports whether
// the message belonged to Light.
func (a *App) updateLight(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case lightPreparedMsg:
		a.turnBusy = false
		a.lightPrepared(msg)
		return nil, true
	case lightConfirmedMsg:
		a.turnBusy = false
		text := "Delivery uncertain; /confirm-light retries the same action ID."
		if errors.Is(msg.err, engine.ErrActionRefresh) {
			a.lightRequest = nil
			text = "Action changed; prepare Light again."
		}
		if msg.err == nil && msg.result != nil {
			a.lightRequest = nil
			text = msg.result.Narrative
		}
		a.lightSay(text)
		return nil, true
	case narrative.SubmitMsg:
		fields := strings.Fields(msg.Input)
		if len(fields) == 0 || !slices.Contains([]string{"/light", "/confirm-light", "/cancel-light", "/dismiss-light", "/cover-light", "/uncover-light"}, fields[0]) {
			return nil, false
		}
		return a.submitLight(fields), true
	}
	return nil, false
}

func (a *App) submitLight(fields []string) tea.Cmd {
	if a.turnBusy {
		return nil
	}
	se, ok := a.engine.(lightEngine)
	if !ok {
		return nil
	}
	ctx, campaign := a.ctx, dbutil.FromPgtype(a.campaign.ID)
	command, args := fields[0], fields[1:]
	if command != "/light" && len(args) != 0 {
		a.lightSay(lightUsage)
		return nil
	}
	switch command {
	case "/cancel-light":
		a.lightRequest = nil
		a.lightSay("Light confirmation dismissed.")
		return nil
	case "/confirm-light":
		if a.lightRequest == nil {
			a.lightSay(lightUsage)
			return nil
		}
		request := *a.lightRequest
		a.turnBusy = true
		return func() tea.Msg {
			result, err := se.ConfirmAction(ctx, campaign, request)
			return lightConfirmedMsg{result, err}
		}
	}

	prepared := lightPreparedMsg{}
	var intent game.ActionIntent
	switch {
	case command == "/dismiss-light":
		intent.Kind = srd521.DismissLight
	case command == "/cover-light":
		intent.Kind = srd521.CoverLight
	case command == "/uncover-light":
		intent.Kind = srd521.UncoverLight
	case len(args) == 0:
		// The read-only dismissal preview carries the active Light.
		intent.Kind, prepared.status = srd521.DismissLight, true
	case len(args) <= 2:
		intent = game.ActionIntent{Kind: srd521.SpellCast, EntryID: game.LightEntryPrefix + strings.ToLower(args[0])}
		if len(args) == 2 {
			number, err := strconv.Atoi(args[1])
			if err != nil || number < 1 {
				a.lightSay(lightUsage)
				return nil
			}
			prepared.object = number
		}
	default:
		a.lightSay(lightUsage)
		return nil
	}
	if a.lightRequest != nil && !prepared.status {
		a.lightSay("Confirm/retry or cancel the pending Light action first.")
		return nil
	}
	a.turnBusy = true
	return func() tea.Msg {
		prepared.stored, prepared.err = se.PreviewLight(ctx, campaign, intent)
		return prepared
	}
}

func (a *App) lightPrepared(msg lightPreparedMsg) {
	if msg.err != nil || msg.stored == nil || msg.stored.Preview == nil {
		a.lightSay("Light is unavailable for this campaign's character.")
		return
	}
	preview := msg.stored.Preview
	if msg.status {
		text := preview.Reason
		if light := msg.stored.Light; light != nil {
			covered := ""
			if light.Covered {
				covered = " Covered: it sheds no light until uncovered."
			}
			text = fmt.Sprintf("Active Light on %s (%s), expires %s.%s", msg.stored.LightItem, light.Color, campaignTime(light.ExpiresAt), covered)
			if !preview.Available {
				text += " " + preview.Reason
			}
		}
		a.lightSay(text + " Colors: " + strings.Join(srd521.LightColors, ", ") + ". " + lightUsage)
		return
	}
	if !preview.Available {
		a.lightSay(preview.Reason)
		return
	}
	request := engine.PreparedAction{ActionID: uuid.New(), ExpectedRevision: preview.ExpectedRevision, Intent: preview.Intent}
	object := ""
	if preview.Intent.Kind == srd521.SpellCast {
		if msg.object == 0 {
			lines := []string{"Carried objects:"}
			for i, target := range preview.Targets {
				lines = append(lines, fmt.Sprintf("%d. %s", i+1, target.Label))
			}
			color, _ := game.LightColor(preview.Intent.EntryID)
			a.lightSay(strings.Join(lines, "\n") + "\nType /light " + color + " <number> to preview the cast. Nothing spent.")
			return
		}
		if msg.object > len(preview.Targets) {
			a.lightSay("That object number is not offered; list the carried objects again with /light <color>.")
			return
		}
		target := preview.Targets[msg.object-1]
		request.TargetID, object = target.ID, " Object: "+target.Label+"."
	}
	a.lightRequest = &request
	a.lightSay("Preview only: " + strings.Join(preview.Options, " ") + object + " /confirm-light to commit; /cancel-light to dismiss.")
}

func (a *App) lightSay(text string) {
	if nv := a.narrativeView(); nv != nil {
		nv.AddEntry(narrative.Entry{Kind: narrative.KindSystem, Text: text})
	}
}

// campaignTime matches the web client's campaign clock format.
func campaignTime(seconds int) string {
	return fmt.Sprintf("Day %d, %02d:%02d:%02d", seconds/86400+1, seconds/3600%24, seconds/60%60, seconds%60)
}
