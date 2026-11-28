package game

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

// PlayerRecords is the saved campaign data a player reads beside the narrative:
// the current character, what that character carries and the campaign's quests.
type PlayerRecords struct {
	// Character is nil when the campaign has no character.
	Character *CharacterRecord
	Inventory []domain.Item
	// Quests holds every quest, finished ones too, oldest first.
	Quests []QuestRecord
}

// CharacterRecord is the current character as its sheet shows it.
type CharacterRecord struct {
	domain.PlayerCharacter
	RulesVersion string
	// Location is the current location's name, or empty without one.
	Location string
	// Sheet is the server-derived sheet of an SRD 5.2.1 character. It is nil
	// for other rules, and when SheetError says why it could not be derived.
	Sheet      *srd521.CharacterSheet
	SheetError error
}

// QuestRecord is a quest with its objectives in order.
type QuestRecord struct {
	domain.Quest
	Objectives []domain.QuestObjective
}

// LoadPlayerRecords reads the campaign's records. database may be nil, which
// leaves an SRD sheet without its saved effects and exploration choices.
func LoadPlayerRecords(ctx context.Context, queries statedb.Querier, database statedb.DBTX, campaignID uuid.UUID) (*PlayerRecords, error) {
	pgID := dbutil.ToPgtype(campaignID)
	campaign, err := queries.GetCampaignByID(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("read campaign: %w", err)
	}
	records := &PlayerRecords{}

	characters, err := queries.GetPlayerCharacterByCampaign(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("read character: %w", err)
	}
	if len(characters) > 0 {
		// The same current character the turn loader acts as: the newest.
		pc := characters[len(characters)-1]
		character := &CharacterRecord{PlayerCharacter: playerCharacterToDomain(pc), RulesVersion: campaign.RulesVersion}
		if pc.CurrentLocationID.Valid {
			location, err := queries.GetLocationByID(ctx, statedb.GetLocationByIDParams{ID: pc.CurrentLocationID, CampaignID: pgID})
			if err != nil {
				return nil, fmt.Errorf("read character location: %w", err)
			}
			character.Location = location.Name
		}
		if campaign.RulesVersion == string(domain.RulesVersionSRD521) {
			character.Sheet, character.SheetError = ProjectCharacterSheet(ctx, database, pc)
		}
		records.Character = character

		items, err := queries.ListItemsByPlayer(ctx, statedb.ListItemsByPlayerParams{CampaignID: pgID, PlayerCharacterID: pc.ID})
		if err != nil {
			return nil, fmt.Errorf("read inventory: %w", err)
		}
		for _, item := range items {
			records.Inventory = append(records.Inventory, itemToDomain(item))
		}
	}

	quests, err := queries.ListQuestsByCampaign(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("read quests: %w", err)
	}
	if len(quests) == 0 {
		return records, nil
	}
	ids := make([]pgtype.UUID, len(quests))
	for i, quest := range quests {
		ids[i] = quest.ID
	}
	objectives, err := queries.ListObjectivesByQuests(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("read quest objectives: %w", err)
	}
	byQuest := make(map[uuid.UUID][]domain.QuestObjective, len(quests))
	for _, objective := range objectives {
		id := dbutil.FromPgtype(objective.QuestID)
		byQuest[id] = append(byQuest[id], questObjectiveToDomain(objective))
	}
	for _, quest := range quests {
		record := QuestRecord{Quest: questToDomain(quest)}
		record.Objectives = byQuest[record.ID]
		records.Quests = append(records.Quests, record)
	}
	return records, nil
}

// FormatCampaignTime writes a campaign clock value, in seconds, as the web
// client does.
func FormatCampaignTime(seconds int) string {
	return fmt.Sprintf("Day %d, %02d:%02d:%02d", seconds/86400+1, seconds/3600%24, seconds/60%60, seconds%60)
}

// ErrSheetLoadFailed marks a ProjectCharacterSheet error from inside the
// server. Its text can name tables, hosts and credentials, so it is for logs
// only. Every other error from there says what is wrong with the saved build
// and may be shown to the player.
var ErrSheetLoadFailed = errors.New("load character sheet")

// ProjectCharacterSheet derives the sheet of a saved SRD 5.2.1 character as
// every client shows it: the rules projection, the saved conditions and, when
// database is not nil, the active spell effects and exploration choices. It
// returns no sheet with an error.
func ProjectCharacterSheet(ctx context.Context, database statedb.DBTX, pc statedb.PlayerCharacter) (*srd521.CharacterSheet, error) {
	sheet, err := srd521.Project(pc.Stats, int(pc.Level), int(pc.Hp), int(pc.MaxHp), int(pc.Experience))
	if errors.Is(err, srd521.ErrCoverageUnavailable) {
		return nil, fmt.Errorf("%w: %w", ErrSheetLoadFailed, err)
	}
	if err != nil {
		return nil, err
	}
	if c := sheet.Concentration; c != nil {
		c.SpellLabel = srd521.ConcentrationLabel(c.Spell)
	}
	sheet.Fall = PendingFall(pc.Stats)
	sheet.DeathSaves = ZeroHPState(pc.Stats, int(pc.Hp))
	sheet.Prone = ReadConditions(pc.Stats).Prone
	if database != nil {
		if err := applySavedEffects(ctx, database, pc, sheet); err != nil {
			return nil, fmt.Errorf("%w: %w", ErrSheetLoadFailed, err)
		}
	}
	return sheet, nil
}

// applySavedEffects adds the effects and choices at the campaign's current time.
func applySavedEffects(ctx context.Context, database statedb.DBTX, pc statedb.PlayerCharacter, sheet *srd521.CharacterSheet) error {
	campaignID, actor := dbutil.FromPgtype(pc.CampaignID), playerCharacterToDomain(pc)
	now, err := ReadSpellClock(ctx, database, campaignID)
	if err != nil {
		return fmt.Errorf("read spell clock: %w", err)
	}
	if err := ApplyFaithAC(ctx, database, campaignID, actor.ID, now, sheet); err != nil {
		return fmt.Errorf("apply Shield of Faith: %w", err)
	}
	if err := ApplyLight(ctx, database, campaignID, actor.ID, now, sheet); err != nil {
		return fmt.Errorf("apply Light: %w", err)
	}
	if err := ApplyExploration(ctx, database, domain.Campaign{ID: campaignID}, actor, sheet); err != nil {
		return fmt.Errorf("read exploration choices: %w", err)
	}
	return nil
}
