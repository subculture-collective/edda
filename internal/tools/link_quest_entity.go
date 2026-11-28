package tools

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/llm"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

const linkQuestEntityToolName = "link_quest_entity"

// LinkQuestEntityStore validates quests and entities, then creates the linking relationship.
type LinkQuestEntityStore interface {
	GetLocationByID(ctx context.Context, id pgtype.UUID) (statedb.Location, error)
	GetQuestByID(ctx context.Context, id pgtype.UUID) (statedb.Quest, error)
	GetNPCByID(ctx context.Context, id pgtype.UUID) (statedb.Npc, error)
	GetFactionByID(ctx context.Context, id pgtype.UUID) (statedb.Faction, error)
	GetPlayerCharacterByID(ctx context.Context, id pgtype.UUID) (statedb.PlayerCharacter, error)
	GetItemByID(ctx context.Context, id pgtype.UUID) (statedb.Item, error)
	CreateRelationship(ctx context.Context, arg statedb.CreateRelationshipParams) (statedb.EntityRelationship, error)
	GetRelationshipsByEntity(ctx context.Context, arg statedb.GetRelationshipsByEntityParams) ([]statedb.EntityRelationship, error)
}

// LinkQuestEntityTool returns the link_quest_entity tool definition and JSON schema.
func LinkQuestEntityTool() llm.Tool {
	return llm.Tool{
		Name:        linkQuestEntityToolName,
		Description: "Link an existing entity (NPC, location, faction, player character, or item) to a quest.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"quest_id": map[string]any{
					"type":        "string",
					"description": "UUID of the quest to link.",
				},
				"entity_type": map[string]any{
					"type":        "string",
					"description": "Entity type (npc, location, faction, player_character, item). You may also provide player as an alias for player_character.",
				},
				"entity_id": map[string]any{
					"type":        "string",
					"description": "UUID of the entity to link.",
				},
				"link_description": map[string]any{
					"type":        "string",
					"description": "Narrative description of how this entity relates to the quest.",
				},
				"link_type": map[string]any{
					"type":        "string",
					"description": "Optional relationship category. Defaults to quest_related.",
				},
			},
			"required":             []string{"quest_id", "entity_type", "entity_id", "link_description"},
			"additionalProperties": false,
		},
	}
}

// RegisterLinkQuestEntity registers the link_quest_entity tool and handler.
func RegisterLinkQuestEntity(reg *Registry, store LinkQuestEntityStore) error {
	if store == nil {
		return errors.New("link_quest_entity store is required")
	}
	return reg.Register(LinkQuestEntityTool(), NewLinkQuestEntityHandler(store).Handle)
}

// LinkQuestEntityHandler executes link_quest_entity tool calls.
type LinkQuestEntityHandler struct {
	store LinkQuestEntityStore
}

// NewLinkQuestEntityHandler creates a new link_quest_entity handler.
func NewLinkQuestEntityHandler(store LinkQuestEntityStore) *LinkQuestEntityHandler {
	return &LinkQuestEntityHandler{store: store}
}

// Handle executes the link_quest_entity tool.
func (h *LinkQuestEntityHandler) Handle(ctx context.Context, args map[string]any) (*ToolResult, error) {
	if h == nil {
		return nil, errors.New("link_quest_entity handler is nil")
	}
	if h.store == nil {
		return nil, errors.New("link_quest_entity store is required")
	}

	questID, err := parseUUIDArg(args, "quest_id")
	if err != nil {
		return nil, err
	}
	entityType, err := parseRelationshipEntityTypeArg(args, "entity_type")
	if err != nil {
		return nil, err
	}
	entityID, err := parseUUIDArg(args, "entity_id")
	if err != nil {
		return nil, err
	}
	linkDescription, err := parseStringArg(args, "link_description")
	if err != nil {
		return nil, err
	}
	trimmedDescription := strings.TrimSpace(linkDescription)
	if trimmedDescription == "" {
		return nil, errors.New("link_description must not be empty or whitespace")
	}

	linkType := "quest_related"
	if raw, ok := args["link_type"]; ok {
		if s, ok := raw.(string); ok {
			trimmed := strings.TrimSpace(s)
			if trimmed != "" {
				linkType = trimmed
			}
		}
	}

	// Resolve campaign from current location.
	currentLocationID, ok := CurrentLocationIDFromContext(ctx)
	if !ok {
		return nil, errors.New("link_quest_entity requires current location id in context")
	}
	currentLocation, err := h.store.GetLocationByID(ctx, dbutil.ToPgtype(currentLocationID))
	if err != nil {
		return nil, fmt.Errorf("resolve campaign from current location: %w", err)
	}
	campaignID := dbutil.FromPgtype(currentLocation.CampaignID)

	// Validate the quest exists in the campaign; another campaign's quest
	// looks missing.
	quest, err := h.store.GetQuestByID(ctx, dbutil.ToPgtype(questID))
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && dbutil.FromPgtype(quest.CampaignID) != campaignID) {
		return nil, fmt.Errorf("quest not found: %s", questID)
	}
	if err != nil {
		return nil, fmt.Errorf("lookup quest: %w", err)
	}

	// Validate entity exists and belongs to the same campaign.
	if err := validateEntityInCampaign(ctx, h.store, campaignID, entityType, entityID); err != nil {
		return nil, err
	}

	relationship, err := h.store.CreateRelationship(ctx, statedb.CreateRelationshipParams{
		CampaignID:       dbutil.ToPgtype(campaignID),
		SourceEntityType: "quest",
		SourceEntityID:   quest.ID,
		TargetEntityType: entityType,
		TargetEntityID:   dbutil.ToPgtype(entityID),
		RelationshipType: linkType,
		Description:      pgtype.Text{String: trimmedDescription, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("create quest-entity link: %w", err)
	}

	summary := fmt.Sprintf("Linked %s %s to quest %q via %q.",
		entityType, entityID.String(), quest.Title, linkType)

	return &ToolResult{
		Success: true,
		Data: map[string]any{
			"id":                 dbutil.FromPgtype(relationship.ID).String(),
			"campaign_id":        campaignID.String(),
			"quest_id":           questID.String(),
			"source_entity_type": relationship.SourceEntityType,
			"source_entity_id":   dbutil.FromPgtype(relationship.SourceEntityID).String(),
			"target_entity_type": relationship.TargetEntityType,
			"target_entity_id":   dbutil.FromPgtype(relationship.TargetEntityID).String(),
			"relationship_type":  relationship.RelationshipType,
			"description":        relationship.Description.String,
			"summary":            summary,
		},
		Narrative: summary,
	}, nil
}
