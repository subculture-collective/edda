package world

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

// CharacterSpawnPackage describes starting state granted to a new character.
type CharacterSpawnPackage struct {
	Items         []StarterItem
	KnownFacts    []StarterKnownFact
	Relationships []StarterRelationship
}

// StarterItem describes a starting inventory item.
type StarterItem struct {
	Name        string
	Description string
	ItemType    string
	Rarity      string
	Properties  map[string]any
	Equipped    bool
	Quantity    int32
}

// StarterKnownFact describes a player-known fact granted at spawn.
type StarterKnownFact struct {
	Fact     string
	Category string
}

// StarterRelationship describes a player-aware relationship granted at spawn.
type StarterRelationship struct {
	TargetEntityType string
	TargetEntityID   uuid.UUID
	RelationshipType string
	Description      string
	Strength         *int32
}

func ApplySpawnPackage(ctx context.Context, q statedb.Querier, campaignID uuid.UUID, playerID uuid.UUID, pkg *CharacterSpawnPackage) error {
	if q == nil || pkg == nil {
		return nil
	}
	campaignPgID := dbutil.ToPgtype(campaignID)
	playerPgID := dbutil.ToPgtype(playerID)

	for _, item := range pkg.Items {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}
		properties := []byte(`{}`)
		if len(item.Properties) > 0 {
			encoded, err := json.Marshal(item.Properties)
			if err != nil {
				return fmt.Errorf("spawn package: encode item %q properties: %w", name, err)
			}
			properties = encoded
		}
		quantity := item.Quantity
		if quantity < 1 {
			quantity = 1
		}
		itemType := normalizeStarterItemType(item.ItemType)
		rarity := strings.TrimSpace(item.Rarity)
		if rarity == "" {
			rarity = "common"
		}
		desc := strings.TrimSpace(item.Description)
		_, err := q.CreateItem(ctx, statedb.CreateItemParams{CampaignID: campaignPgID, PlayerCharacterID: playerPgID, Name: name, Description: pgtype.Text{String: desc, Valid: desc != ""}, ItemType: itemType, Rarity: rarity, Properties: properties, Equipped: item.Equipped, Quantity: quantity})
		if err != nil {
			return fmt.Errorf("spawn package: create item %q: %w", name, err)
		}
	}

	for _, fact := range pkg.KnownFacts {
		text := strings.TrimSpace(fact.Fact)
		if text == "" {
			continue
		}
		category := strings.TrimSpace(fact.Category)
		if category == "" {
			category = "lore"
		}
		_, err := q.CreateFact(ctx, statedb.CreateFactParams{CampaignID: campaignPgID, Fact: text, Category: category, Source: "character_spawn", PlayerKnown: true})
		if err != nil {
			return fmt.Errorf("spawn package: create known fact %q: %w", text, err)
		}
	}

	for _, rel := range pkg.Relationships {
		targetType := strings.TrimSpace(rel.TargetEntityType)
		relType := strings.TrimSpace(rel.RelationshipType)
		if targetType == "" || rel.TargetEntityID == uuid.Nil || relType == "" {
			continue
		}
		if !ValidSpawnRelationshipEntityType(targetType) {
			return fmt.Errorf("spawn package: invalid relationship target_entity_type %q", targetType)
		}
		strength := pgtype.Int4{}
		if rel.Strength != nil {
			strength = pgtype.Int4{Int32: *rel.Strength, Valid: true}
		}
		desc := strings.TrimSpace(rel.Description)
		created, err := q.CreateRelationship(ctx, statedb.CreateRelationshipParams{CampaignID: campaignPgID, SourceEntityType: "player_character", SourceEntityID: playerPgID, TargetEntityType: targetType, TargetEntityID: dbutil.ToPgtype(rel.TargetEntityID), RelationshipType: relType, Description: pgtype.Text{String: desc, Valid: desc != ""}, Strength: strength})
		if err != nil {
			return fmt.Errorf("spawn package: create relationship %q: %w", relType, err)
		}
		if created.ID.Valid {
			if err := q.SetRelationshipPlayerAware(ctx, created.ID); err != nil {
				return fmt.Errorf("spawn package: reveal relationship %q: %w", relType, err)
			}
		}
	}

	return nil
}

func normalizeStarterItemType(itemType string) string {
	switch strings.TrimSpace(itemType) {
	case "weapon", "armor", "consumable", "quest", "misc":
		return strings.TrimSpace(itemType)
	default:
		return "misc"
	}
}

// ValidSpawnRelationshipEntityType reports whether entityType fits the
// entity_relationships target_entity_type database constraint.
func ValidSpawnRelationshipEntityType(entityType string) bool {
	switch strings.TrimSpace(entityType) {
	case "npc", "location", "faction", "player_character", "item":
		return true
	default:
		return false
	}
}
