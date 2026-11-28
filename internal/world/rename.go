package world

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

// UnnamedCharacter returns the campaign's first character with a blank name, or
// nil when every character has one. Earlier versions of the TUI's guided
// character form saved blank names.
func UnnamedCharacter(ctx context.Context, queries statedb.Querier, campaignID pgtype.UUID) (*statedb.PlayerCharacter, error) {
	characters, err := queries.GetPlayerCharacterByCampaign(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("list campaign characters: %w", err)
	}
	for _, character := range characters {
		if domain.ValidateCharacterName(character.Name) != nil {
			return &character, nil
		}
	}
	return nil, nil
}

// ErrCharacterNotFound refuses a rename of a character that is not in the
// campaign.
var ErrCharacterNotFound = errors.New("character not found in this campaign")

// RenameCharacter saves a new name for one of the campaign's characters.
func RenameCharacter(ctx context.Context, queries statedb.Querier, campaignID, characterID pgtype.UUID, name string) error {
	if err := domain.ValidateCharacterName(name); err != nil {
		return err
	}
	renamed, err := queries.RenamePlayerCharacter(ctx, statedb.RenamePlayerCharacterParams{
		ID: characterID, CampaignID: campaignID, Name: strings.TrimSpace(name),
	})
	if err != nil {
		return fmt.Errorf("rename character: %w", err)
	}
	if renamed != 1 {
		return fmt.Errorf("rename character: %w", ErrCharacterNotFound)
	}
	return nil
}
