package tools

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

// requireCurrentCampaignID returns the campaign the server set for this turn.
// Tools check every ID the model supplies against it, never against a
// campaign ID the model supplies.
func requireCurrentCampaignID(ctx context.Context, toolName string) (uuid.UUID, error) {
	campaignID, ok := CurrentCampaignIDFromContext(ctx)
	if !ok {
		return uuid.Nil, fmt.Errorf("%s requires current campaign id in context", toolName)
	}
	return campaignID, nil
}

// campaignEntityLookup loads the entities a quest, lore entry or relationship
// can point at.
type campaignEntityLookup interface {
	GetLocationByID(ctx context.Context, id pgtype.UUID) (statedb.Location, error)
	GetNPCByID(ctx context.Context, id pgtype.UUID) (statedb.Npc, error)
	GetFactionByID(ctx context.Context, id pgtype.UUID) (statedb.Faction, error)
	GetPlayerCharacterByID(ctx context.Context, id pgtype.UUID) (statedb.PlayerCharacter, error)
	GetItemByID(ctx context.Context, id pgtype.UUID) (statedb.Item, error)
}

// validateEntityInCampaign checks that an entity exists in the campaign. An
// entity in another campaign gets the same error as a missing one, so the
// error does not reveal other campaigns' IDs.
func validateEntityInCampaign(ctx context.Context, store campaignEntityLookup, campaignID uuid.UUID, entityType string, entityID uuid.UUID) error {
	id := dbutil.ToPgtype(entityID)
	var owner pgtype.UUID
	var err error
	switch domain.EntityType(entityType) {
	case domain.EntityTypeNPC:
		var row statedb.Npc
		row, err = store.GetNPCByID(ctx, id)
		owner = row.CampaignID
	case domain.EntityTypeLocation:
		var row statedb.Location
		row, err = store.GetLocationByID(ctx, id)
		owner = row.CampaignID
	case domain.EntityTypeFaction:
		var row statedb.Faction
		row, err = store.GetFactionByID(ctx, id)
		owner = row.CampaignID
	case domain.EntityTypePlayerCharacter:
		var row statedb.PlayerCharacter
		row, err = store.GetPlayerCharacterByID(ctx, id)
		owner = row.CampaignID
	case domain.EntityTypeItem:
		var row statedb.Item
		row, err = store.GetItemByID(ctx, id)
		owner = row.CampaignID
	default:
		return fmt.Errorf("entity_type must be one of: npc, location, faction, player_character, player, item")
	}
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && dbutil.FromPgtype(owner) != campaignID) {
		return fmt.Errorf("entity not found: %s %s", entityType, entityID)
	}
	if err != nil {
		return fmt.Errorf("lookup %s: %w", entityType, err)
	}
	return nil
}
