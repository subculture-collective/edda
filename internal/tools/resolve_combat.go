package tools

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/combat"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/llm"
)

const (
	resolveCombatToolName = "resolve_combat"

	outcomeVictory   = "victory"
	outcomeDefeat    = "defeat"
	outcomeFlee      = "flee"
	outcomeSurrender = "surrender"

	playerStatusActive   = "active"
	playerStatusDefeated = "defeated"
)

// lootInput represents a single item to create as combat loot.
type lootInput struct {
	Name        string
	Description string
	ItemType    string
	Quantity    int
}

// ResolveCombatStore provides persistence needed for resolve_combat.
type ResolveCombatStore interface {
	GetPlayerCharacterByID(ctx context.Context, playerCharacterID uuid.UUID) (*domain.PlayerCharacter, error)
	UpdatePlayerHP(ctx context.Context, playerCharacterID uuid.UUID, hp, maxHP int) error
	UpdatePlayerStatus(ctx context.Context, playerCharacterID uuid.UUID, status string) error
	UpdatePlayerLocation(ctx context.Context, playerCharacterID uuid.UUID, locationID uuid.UUID) error
	AddPlayerExperience(ctx context.Context, playerCharacterID uuid.UUID, xpAmount int) error
	CreatePlayerItem(ctx context.Context, playerCharacterID uuid.UUID, name, description, itemType, rarity string, quantity int) (uuid.UUID, error)
	MarkNPCDead(ctx context.Context, npcID uuid.UUID) error
	// GetNPCByID returns nil for an NPC outside the campaign.
	GetNPCByID(ctx context.Context, campaignID, npcID uuid.UUID) (*domain.NPC, error)
	LocationExistsInCampaign(ctx context.Context, locationID, campaignID uuid.UUID) (bool, error)
	UpdateNPCDisposition(ctx context.Context, npcID uuid.UUID, newDisposition int) error
}

// ResolveCombatTool returns the resolve_combat tool definition and JSON schema.
func ResolveCombatTool() llm.Tool {
	return llm.Tool{
		Name: resolveCombatToolName,
		Description: "End combat and resolve the outcome: distribute XP and loot on victory, " +
			"handle player defeat or capture, move the player on flee, or apply NPC disposition changes " +
			"on surrender. Persists updated player HP back to the character record and clears combat mode.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"combat_state": map[string]any{
					"type":        "object",
					"description": "Current combat state from the last combat round.",
				},
				"outcome_type": map[string]any{
					"type":        "string",
					"description": "How combat ended. One of: victory, defeat, flee, surrender.",
				},
				"xp_earned": map[string]any{
					"type":        "integer",
					"description": "Experience points to award the player. Used on victory.",
				},
				"loot": map[string]any{
					"type":        "array",
					"description": "Items to add to the player's inventory. Used on victory.",
					"items": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"name":        map[string]any{"type": "string", "description": "Item name."},
							"description": map[string]any{"type": "string", "description": "Item description."},
							"item_type":   map[string]any{"type": "string", "description": "Item type: weapon, armor, consumable, quest, or misc."},
							"quantity":    map[string]any{"type": "integer", "description": "Quantity. Defaults to 1."},
						},
						"required":             []string{"name", "description", "item_type"},
						"additionalProperties": false,
					},
				},
				"dead_npc_ids": map[string]any{
					"type":        "array",
					"description": "UUIDs of NPCs killed in this combat. Used on victory.",
					"items":       map[string]any{"type": "string"},
				},
				"flee_location_id": map[string]any{
					"type":        "string",
					"description": "Destination location UUID the player flees to. Used on flee.",
				},
				"surrender_npc_ids": map[string]any{
					"type":        "array",
					"description": "UUIDs of NPCs whose disposition changes on surrender.",
					"items":       map[string]any{"type": "string"},
				},
				"disposition_change": map[string]any{
					"type":        "integer",
					"description": "Disposition delta applied to each surrendering NPC (-100 to 100). Used on surrender.",
				},
				"consequences": map[string]any{
					"type":        "string",
					"description": "Narrative description of the lasting combat consequences.",
				},
			},
			"required":             []string{"combat_state", "outcome_type"},
			"additionalProperties": false,
		},
	}
}

// RegisterResolveCombat registers the resolve_combat tool and handler.
func RegisterResolveCombat(reg *Registry, store ResolveCombatStore) error {
	if store == nil {
		return errors.New("resolve_combat store is required")
	}
	return reg.Register(ResolveCombatTool(), NewResolveCombatHandler(store).Handle)
}

// ResolveCombatHandler executes resolve_combat tool calls.
type ResolveCombatHandler struct {
	store ResolveCombatStore
}

// NewResolveCombatHandler creates a new resolve_combat handler.
func NewResolveCombatHandler(store ResolveCombatStore) *ResolveCombatHandler {
	return &ResolveCombatHandler{store: store}
}

// combatResolution is a validated resolve_combat call. Handle builds it with
// reads only, then applies it, so a rejected argument saves nothing. The
// engine runs each tool call in one transaction, so a failed write rolls back
// the writes before it.
type combatResolution struct {
	xpEarned       int
	deadNPCIDs     []uuid.UUID
	loot           []lootInput
	fleeLocationID *uuid.UUID
	surrender      []surrenderUpdate
	playerStatus   string
}

type surrenderUpdate struct {
	npcID          uuid.UUID
	oldDisposition int
	newDisposition int
}

// Handle executes the resolve_combat tool.
func (h *ResolveCombatHandler) Handle(ctx context.Context, args map[string]any) (*ToolResult, error) {
	if h == nil {
		return nil, errors.New("resolve_combat handler is nil")
	}
	if h.store == nil {
		return nil, errors.New("resolve_combat store is required")
	}

	playerCharacterID, ok := CurrentPlayerCharacterIDFromContext(ctx)
	if !ok {
		return nil, errors.New("resolve_combat requires current player character id in context")
	}
	campaignID, err := requireCurrentCampaignID(ctx, resolveCombatToolName)
	if err != nil {
		return nil, err
	}

	state, err := parseCombatStateArgWithCampaignFallback(args, "combat_state", campaignID)
	if err != nil {
		return nil, err
	}
	if state.CampaignID != campaignID {
		return nil, errors.New("combat_state.campaign_id must be the current campaign")
	}

	outcomeType, err := parseStringArg(args, "outcome_type")
	if err != nil {
		return nil, err
	}
	outcomeType = strings.ToLower(strings.TrimSpace(outcomeType))
	switch outcomeType {
	case outcomeVictory, outcomeDefeat, outcomeFlee, outcomeSurrender:
	default:
		return nil, fmt.Errorf("outcome_type must be one of: victory, defeat, flee, surrender")
	}

	player, err := h.store.GetPlayerCharacterByID(ctx, playerCharacterID)
	if err != nil {
		return nil, fmt.Errorf("get player character: %w", err)
	}
	if player == nil || player.CampaignID != campaignID {
		return nil, errors.New("player character not found")
	}

	// The model's combat state is not trusted: its player must be the current
	// player, and only the current HP is kept, within the stored maximum.
	playerCombatant := firstCombatantByType(state, combat.CombatantTypePlayer)
	if playerCombatant != nil && playerCombatant.EntityID != playerCharacterID {
		return nil, errors.New("combat_state player combatant must be the current player character")
	}

	resolution, err := h.validateOutcome(ctx, outcomeType, args, state, campaignID)
	if err != nil {
		return nil, err
	}

	// Every argument is valid; write.
	if playerCombatant != nil {
		hp := min(max(playerCombatant.HP, 0), player.MaxHP)
		if err := h.store.UpdatePlayerHP(ctx, playerCharacterID, hp, player.MaxHP); err != nil {
			return nil, fmt.Errorf("persist player hp: %w", err)
		}
	}
	summary, err := h.applyOutcome(ctx, outcomeType, resolution, playerCharacterID)
	if err != nil {
		return nil, err
	}

	// Update combat state status to reflect the resolved outcome.
	state.Status = combatStatusFromOutcome(outcomeType)
	summary["combat_state"] = combatStateToMap(state)
	summary["outcome_type"] = outcomeType

	consequences, _ := parseOptionalNonEmptyStringArg(args, "consequences")
	if consequences != nil {
		summary["consequences"] = *consequences
	}

	return &ToolResult{
		Success:   true,
		Data:      summary,
		Narrative: buildResolveCombatNarrative(outcomeType, player.Name, summary),
	}, nil
}

// validateOutcome parses and checks the outcome arguments without writing.
func (h *ResolveCombatHandler) validateOutcome(
	ctx context.Context,
	outcomeType string,
	args map[string]any,
	state *combat.CombatState,
	campaignID uuid.UUID,
) (*combatResolution, error) {
	resolution := &combatResolution{playerStatus: playerStatusActive}
	switch outcomeType {
	case outcomeVictory:
		xpEarned, _, err := parseOptionalIntArg(args, "xp_earned")
		if err != nil {
			return nil, err
		}
		if xpEarned < 0 {
			return nil, errors.New("xp_earned must be greater than or equal to 0")
		}
		resolution.xpEarned = xpEarned
		deadNPCs, err := h.combatNPCsArg(ctx, args, "dead_npc_ids", state, campaignID)
		if err != nil {
			return nil, err
		}
		for _, npc := range deadNPCs {
			resolution.deadNPCIDs = append(resolution.deadNPCIDs, npc.ID)
		}
		if resolution.loot, err = parseLootArg(args, "loot"); err != nil {
			return nil, err
		}
	case outcomeDefeat:
		resolution.playerStatus = playerStatusDefeated
	case outcomeFlee:
		fleeLocationID, fleeLocationSet, err := parseOptionalUUIDArg(args, "flee_location_id")
		if err != nil {
			return nil, err
		}
		if fleeLocationSet {
			exists, err := h.store.LocationExistsInCampaign(ctx, fleeLocationID, campaignID)
			if err != nil {
				return nil, fmt.Errorf("check flee location: %w", err)
			}
			if !exists {
				return nil, errors.New("flee_location_id does not reference an existing location")
			}
			resolution.fleeLocationID = &fleeLocationID
		}
	case outcomeSurrender:
		dispositionChange, _, err := parseOptionalIntArg(args, "disposition_change")
		if err != nil {
			return nil, err
		}
		surrendering, err := h.combatNPCsArg(ctx, args, "surrender_npc_ids", state, campaignID)
		if err != nil {
			return nil, err
		}
		for _, npc := range surrendering {
			resolution.surrender = append(resolution.surrender, surrenderUpdate{
				npcID:          npc.ID,
				oldDisposition: npc.Disposition,
				newDisposition: clampDispositionValue(npc.Disposition + dispositionChange),
			})
		}
	default:
		return nil, fmt.Errorf("unknown outcome_type: %s", outcomeType)
	}
	return resolution, nil
}

// combatNPCsArg loads the NPCs an ID array names. Each must be an NPC
// combatant in the combat state and a living NPC in the campaign; an NPC in
// another campaign gets the unknown-NPC error.
func (h *ResolveCombatHandler) combatNPCsArg(
	ctx context.Context,
	args map[string]any,
	key string,
	state *combat.CombatState,
	campaignID uuid.UUID,
) ([]*domain.NPC, error) {
	ids, err := parseUUIDArrayArg(args, key)
	if err != nil {
		return nil, err
	}
	combatants := npcCombatantIDSet(state)
	seen := make(map[uuid.UUID]bool, len(ids))
	npcs := make([]*domain.NPC, 0, len(ids))
	for _, id := range ids {
		if seen[id] {
			continue
		}
		seen[id] = true
		if !combatants[id] {
			return nil, fmt.Errorf("%s: %s is not an NPC combatant in this combat", key, id)
		}
		npc, err := h.store.GetNPCByID(ctx, campaignID, id)
		if err != nil {
			return nil, fmt.Errorf("get npc %s: %w", id, err)
		}
		if npc == nil {
			return nil, fmt.Errorf("%s: %s does not reference an existing npc", key, id)
		}
		if !npc.Alive {
			return nil, fmt.Errorf("%s: %s is already dead", key, id)
		}
		npcs = append(npcs, npc)
	}
	return npcs, nil
}

// applyOutcome writes a validated resolution and returns its summary.
func (h *ResolveCombatHandler) applyOutcome(
	ctx context.Context,
	outcomeType string,
	r *combatResolution,
	playerCharacterID uuid.UUID,
) (map[string]any, error) {
	summary := map[string]any{}
	switch outcomeType {
	case outcomeVictory:
		if r.xpEarned > 0 {
			if err := h.store.AddPlayerExperience(ctx, playerCharacterID, r.xpEarned); err != nil {
				return nil, fmt.Errorf("add player experience: %w", err)
			}
		}
		summary["xp_earned"] = r.xpEarned
		killed := make([]string, 0, len(r.deadNPCIDs))
		for _, npcID := range r.deadNPCIDs {
			if err := h.store.MarkNPCDead(ctx, npcID); err != nil {
				return nil, fmt.Errorf("mark npc %s dead: %w", npcID, err)
			}
			killed = append(killed, npcID.String())
		}
		summary["dead_npc_ids"] = killed
		createdLoot := make([]map[string]any, 0, len(r.loot))
		for _, item := range r.loot {
			itemID, err := h.store.CreatePlayerItem(ctx, playerCharacterID, item.Name, item.Description, item.ItemType, defaultItemRarity, item.Quantity)
			if err != nil {
				return nil, fmt.Errorf("create loot item %q: %w", item.Name, err)
			}
			createdLoot = append(createdLoot, map[string]any{
				"item_id":     itemID.String(),
				"name":        item.Name,
				"description": item.Description,
				"item_type":   item.ItemType,
				"quantity":    item.Quantity,
			})
		}
		summary["loot"] = createdLoot
	case outcomeFlee:
		if r.fleeLocationID != nil {
			if err := h.store.UpdatePlayerLocation(ctx, playerCharacterID, *r.fleeLocationID); err != nil {
				return nil, fmt.Errorf("move player on flee: %w", err)
			}
			summary["flee_location_id"] = r.fleeLocationID.String()
		}
	case outcomeSurrender:
		updatedNPCs := make([]map[string]any, 0, len(r.surrender))
		for _, u := range r.surrender {
			if err := h.store.UpdateNPCDisposition(ctx, u.npcID, u.newDisposition); err != nil {
				return nil, fmt.Errorf("update npc %s disposition: %w", u.npcID, err)
			}
			updatedNPCs = append(updatedNPCs, map[string]any{
				"npc_id":          u.npcID.String(),
				"old_disposition": u.oldDisposition,
				"new_disposition": u.newDisposition,
			})
		}
		summary["updated_npcs"] = updatedNPCs
	}

	if err := h.store.UpdatePlayerStatus(ctx, playerCharacterID, r.playerStatus); err != nil {
		if r.playerStatus == playerStatusDefeated {
			return nil, fmt.Errorf("set defeated status: %w", err)
		}
		return nil, fmt.Errorf("exit combat mode: %w", err)
	}
	summary["player_status"] = r.playerStatus
	return summary, nil
}

// parseLootArg parses the "loot" argument into a slice of lootInput.
func parseLootArg(args map[string]any, key string) ([]lootInput, error) {
	raw, ok := args[key]
	if !ok {
		return nil, nil
	}
	items, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("%s must be an array", key)
	}

	out := make([]lootInput, 0, len(items))
	for i, item := range items {
		obj, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("%s[%d] must be an object", key, i)
		}
		name, err := parseStringArg(obj, "name")
		if err != nil {
			return nil, fmt.Errorf("%s[%d].%w", key, i, err)
		}
		description, err := parseStringArg(obj, "description")
		if err != nil {
			return nil, fmt.Errorf("%s[%d].%w", key, i, err)
		}
		itemType, err := parseStringArg(obj, "item_type")
		if err != nil {
			return nil, fmt.Errorf("%s[%d].%w", key, i, err)
		}
		if _, allowed := allowedItemTypes[itemType]; !allowed {
			return nil, fmt.Errorf("%s[%d].item_type must be one of: weapon, armor, consumable, quest, misc", key, i)
		}
		quantity, err := parsePositiveIntArgWithDefault(obj, "quantity", defaultAddQuantity)
		if err != nil {
			return nil, fmt.Errorf("%s[%d].%w", key, i, err)
		}
		out = append(out, lootInput{
			Name:        name,
			Description: description,
			ItemType:    itemType,
			Quantity:    quantity,
		})
	}
	return out, nil
}

// combatStatusFromOutcome maps an outcome type to the final CombatStatus.
func combatStatusFromOutcome(outcomeType string) combat.CombatStatus {
	if outcomeType == outcomeFlee {
		return combat.CombatStatusFled
	}
	return combat.CombatStatusCompleted
}

// npcCombatantIDSet returns the set of entity IDs for all NPC combatants in state.
func npcCombatantIDSet(state *combat.CombatState) map[uuid.UUID]bool {
	ids := make(map[uuid.UUID]bool, len(state.Combatants))
	for _, c := range state.Combatants {
		if c.EntityType == combat.CombatantTypeNPC {
			ids[c.EntityID] = true
		}
	}
	return ids
}

// buildResolveCombatNarrative constructs a narrative string for the combat outcome.
func buildResolveCombatNarrative(outcomeType, playerName string, summary map[string]any) string {
	switch outcomeType {
	case outcomeVictory:
		xpEarned, _ := summary["xp_earned"].(int)
		loot, _ := summary["loot"].([]map[string]any)
		if xpEarned > 0 && len(loot) > 0 {
			return fmt.Sprintf("Victory! %s has defeated all enemies, earning %d XP and %d item(s) of loot.", playerName, xpEarned, len(loot))
		}
		if xpEarned > 0 {
			return fmt.Sprintf("Victory! %s has defeated all enemies, earning %d XP.", playerName, xpEarned)
		}
		if len(loot) > 0 {
			return fmt.Sprintf("Victory! %s has defeated all enemies and claimed %d item(s) of loot.", playerName, len(loot))
		}
		return fmt.Sprintf("Victory! %s has defeated all enemies.", playerName)
	case outcomeDefeat:
		return fmt.Sprintf("%s has been defeated in combat.", playerName)
	case outcomeFlee:
		if loc, ok := summary["flee_location_id"].(string); ok && loc != "" {
			return fmt.Sprintf("%s fled from combat and moved to a new location.", playerName)
		}
		return fmt.Sprintf("%s fled from combat.", playerName)
	case outcomeSurrender:
		return "The combat has ended in surrender. Terms have been negotiated."
	}
	return fmt.Sprintf("Combat resolved with outcome: %s.", outcomeType)
}
