package tui

import (
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/logging"
	"git.subcult.tv/subculture-collective/edda/tui/character"
	"git.subcult.tv/subculture-collective/edda/tui/inventory"
	"git.subcult.tv/subculture-collective/edda/tui/quest"
)

// recordsEngine reads what the Character, Inventory and Quests tabs show.
type recordsEngine interface {
	PlayerRecords(context.Context, uuid.UUID) (*game.PlayerRecords, error)
}

// recordsLoadedMsg answers one read of the campaign's records.
type recordsLoadedMsg struct {
	// read numbers the reads, so a late answer cannot replace a newer one.
	read    int
	records *game.PlayerRecords
	err     error
}

const (
	recordsUnavailableNotice = "Campaign data is unavailable without a game engine."
	recordsLoadingNotice     = "Loading campaign data…"
	recordsFailedNotice      = "Could not load campaign data. Open this tab again to retry; the Logs tab says why."
	recordsStaleNotice       = "Could not refresh, so this may be out of date. Open this tab again to retry; the Logs tab says why."
)

// recordViews are the tabs that show the campaign's records.
type recordViews struct {
	character *character.Model
	inventory *inventory.Model
	quests    *quest.Model
}

func (v recordViews) show(records *game.PlayerRecords) {
	v.character.Show(records.Character)
	v.inventory.Show(records.Inventory)
	v.quests.Show(records.Quests)
}

func (v recordViews) notify(notice string) {
	v.character.Notify(notice)
	v.inventory.Notify(notice)
	v.quests.Notify(notice)
}

// changesRecords reports whether msg ends a turn or a confirmed action. A
// failed one counts too: its effects may have been saved before the failure.
func changesRecords(msg tea.Msg) bool {
	switch msg.(type) {
	case turnProcessedMsg, surgeConfirmedMsg, tacticalConfirmedMsg, faithConfirmedMsg,
		lightConfirmedMsg, detectMagicConfirmedMsg, fallConfirmedMsg, deathSaveConfirmedMsg,
		explorationConfirmedMsg:
		return true
	}
	return false
}

// readRecords reads the campaign's records for the data tabs. It is nil when
// the engine cannot read them.
func (a App) readRecords() tea.Cmd {
	reader, ok := a.engine.(recordsEngine)
	if !ok {
		return nil
	}
	ctx, campaign, read := a.ctx, dbutil.FromPgtype(a.campaign.ID), a.recordsRead
	return func() tea.Msg {
		records, err := reader.PlayerRecords(ctx, campaign)
		return recordsLoadedMsg{read: read, records: records, err: err}
	}
}

// rereadRecords starts a read that replaces any still on its way.
func (a *App) rereadRecords() tea.Cmd {
	a.recordsRead++
	a.recordsFailed = false
	return a.readRecords()
}

// showRecords puts the answer to the latest read on the data tabs.
func (a *App) showRecords(msg recordsLoadedMsg) {
	if msg.read != a.recordsRead {
		return
	}
	if msg.err != nil {
		a.recordsFailed = true
		a.logger().Error("read campaign records", "campaign_id", dbutil.FromPgtype(a.campaign.ID), "error", logging.Redact(msg.err.Error()))
		notice := recordsFailedNotice
		if a.recordsShown {
			notice = recordsStaleNotice
		}
		a.records.notify(notice)
		return
	}
	a.recordsShown, a.recordsFailed = true, false
	// A sheet the server failed to load is worth reading again. An unsupported
	// build is not, and the sheet itself says what is wrong with it.
	if c := msg.records.Character; c != nil && errors.Is(c.SheetError, game.ErrSheetLoadFailed) {
		a.recordsFailed = true
		a.logger().Error("load character sheet", "campaign_id", dbutil.FromPgtype(a.campaign.ID), "error", logging.Redact(c.SheetError.Error()))
	}
	a.records.show(msg.records)
}

// opensRecordsTab reports whether the open tab changed from previous to one
// that shows the campaign's records.
func (a App) opensRecordsTab(previous int) bool {
	switch open := a.ActiveViewState(); open {
	case ViewCharacterSheet, ViewInventory, ViewQuestLog:
		return int(open) != previous
	}
	return false
}
