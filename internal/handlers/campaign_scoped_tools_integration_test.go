//go:build integration

package handlers_test

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/game"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/internal/tools"
)

// A model turn in one campaign must not read or change another campaign's
// rows by ID. Each tool answers another campaign's ID exactly as it answers an
// unknown ID, and a rejected call saves nothing.
func TestIntegrationToolsRefuseOtherCampaignIDs(t *testing.T) {
	ctx := context.Background()
	q := statedb.New(testPool)
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	newCampaign := func(name string) (campaign, player, location uuid.UUID) {
		t.Helper()
		must(testPool.QueryRow(ctx, `INSERT INTO campaigns(name,status,created_by) VALUES($1,'active',$2) RETURNING id`, name, testUserID).Scan(&campaign))
		t.Cleanup(func() {
			for _, sql := range []string{
				`DELETE FROM entity_relationships WHERE campaign_id=$1`,
				`DELETE FROM quests WHERE campaign_id=$1`,
				`DELETE FROM items WHERE campaign_id=$1`,
				`DELETE FROM npcs WHERE campaign_id=$1`,
				`DELETE FROM factions WHERE campaign_id=$1`,
				`DELETE FROM player_characters WHERE campaign_id=$1`,
				`DELETE FROM locations WHERE campaign_id=$1`,
				`DELETE FROM campaigns WHERE id=$1`,
			} {
				_, _ = testPool.Exec(ctx, sql, campaign)
			}
		})
		must(testPool.QueryRow(ctx, `INSERT INTO locations(campaign_id,name) VALUES($1,'Square') RETURNING id`, campaign).Scan(&location))
		must(testPool.QueryRow(ctx, `INSERT INTO player_characters(campaign_id,user_id,name,stats,hp,max_hp,experience,level,status,abilities,current_location_id)
			VALUES($1,$2,'Mira','{}',12,12,0,1,'in_combat','[]',$3) RETURNING id`, campaign, testUserID, location).Scan(&player))
		return campaign, player, location
	}

	campaign, player, location := newCampaign("Current campaign")
	other, otherPlayer, otherLocation := newCampaign("Other campaign")
	var foe, deadFoe, otherNPC, otherFaction, otherQuest, otherItem uuid.UUID
	must(testPool.QueryRow(ctx, `INSERT INTO npcs(campaign_id,name,location_id,disposition) VALUES($1,'Bandit',$2,0) RETURNING id`, campaign, location).Scan(&foe))
	must(testPool.QueryRow(ctx, `INSERT INTO npcs(campaign_id,name,location_id,alive) VALUES($1,'Fallen bandit',$2,false) RETURNING id`, campaign, location).Scan(&deadFoe))
	// The other campaign's NPC stands in the current location, so a location
	// check alone cannot tell the campaigns apart.
	must(testPool.QueryRow(ctx, `INSERT INTO npcs(campaign_id,name,description,location_id,disposition) VALUES($1,'Envoy','other secret',$2,0) RETURNING id`, other, location).Scan(&otherNPC))
	must(testPool.QueryRow(ctx, `INSERT INTO factions(campaign_id,name) VALUES($1,'Guild') RETURNING id`, other).Scan(&otherFaction))
	must(testPool.QueryRow(ctx, `INSERT INTO quests(campaign_id,title,quest_type,status) VALUES($1,'Other secret quest','short_term','active') RETURNING id`, other).Scan(&otherQuest))
	must(testPool.QueryRow(ctx, `INSERT INTO items(campaign_id,player_character_id,name,item_type,rarity,quantity) VALUES($1,$2,'Other ring','misc','common',2) RETURNING id`, other, otherPlayer).Scan(&otherItem))

	toolCtx := tools.WithCurrentCampaignID(ctx, campaign)
	toolCtx = tools.WithCurrentPlayerCharacterID(toolCtx, player)
	toolCtx = tools.WithCurrentLocationID(toolCtx, location)

	combatState := func(npcs ...uuid.UUID) map[string]any {
		combatants := []any{map[string]any{
			"entity_id": player.String(), "entity_type": "player", "name": "Mira",
			"hp": 3, "max_hp": 99, "status": "alive",
		}}
		for _, id := range npcs {
			combatants = append(combatants, map[string]any{
				"entity_id": id.String(), "entity_type": "npc", "name": "Foe",
				"hp": 0, "max_hp": 5, "status": "dead",
			})
		}
		return map[string]any{
			"id": uuid.NewString(), "campaign_id": campaign.String(),
			"status": "active", "combatants": combatants,
		}
	}

	type snapshot struct {
		hp, maxHP, xp     int
		status            string
		playerLocation    uuid.UUID
		items, quests     int
		relationships     int
		factions          int
		locations         int
		otherItemQuantity int
		otherItemProps    string
		otherDescription  string
		otherDisposition  int
		otherAlive        bool
		foeAlive          bool
		otherPlayerStatus string
	}
	read := func() snapshot {
		t.Helper()
		var s snapshot
		must(testPool.QueryRow(ctx, `SELECT hp,max_hp,experience,status,current_location_id FROM player_characters WHERE id=$1`, player).Scan(&s.hp, &s.maxHP, &s.xp, &s.status, &s.playerLocation))
		must(testPool.QueryRow(ctx, `SELECT count(*) FROM items WHERE campaign_id=$1`, campaign).Scan(&s.items))
		must(testPool.QueryRow(ctx, `SELECT count(*) FROM quests WHERE campaign_id=$1`, campaign).Scan(&s.quests))
		must(testPool.QueryRow(ctx, `SELECT count(*) FROM entity_relationships WHERE campaign_id=$1`, campaign).Scan(&s.relationships))
		must(testPool.QueryRow(ctx, `SELECT count(*) FROM factions WHERE campaign_id=$1`, campaign).Scan(&s.factions))
		must(testPool.QueryRow(ctx, `SELECT count(*) FROM locations WHERE campaign_id=$1`, campaign).Scan(&s.locations))
		must(testPool.QueryRow(ctx, `SELECT quantity,coalesce(properties::text,'') FROM items WHERE id=$1`, otherItem).Scan(&s.otherItemQuantity, &s.otherItemProps))
		must(testPool.QueryRow(ctx, `SELECT coalesce(description,''),disposition,alive FROM npcs WHERE id=$1`, otherNPC).Scan(&s.otherDescription, &s.otherDisposition, &s.otherAlive))
		must(testPool.QueryRow(ctx, `SELECT alive FROM npcs WHERE id=$1`, foe).Scan(&s.foeAlive))
		must(testPool.QueryRow(ctx, `SELECT status FROM player_characters WHERE id=$1`, otherPlayer).Scan(&s.otherPlayerStatus))
		return s
	}
	before := read()

	npcSvc := game.NewNPCService(q)
	combatSvc := game.NewCombatService(q)
	worldSvc := game.NewWorldService(q)
	invSvc := game.NewInventoryService(q)
	calls := map[string]func(id uuid.UUID) (*tools.ToolResult, error){
		"update_npc": func(id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewUpdateNPCHandler(npcSvc).Handle(toolCtx, map[string]any{
				"npc_id": id.String(), "description": "Rewritten by another campaign.", "disposition": 100,
			})
		},
		"npc_dialogue": func(id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewNPCDialogueHandler(npcSvc).Handle(toolCtx, map[string]any{
				"npc_id": id.String(), "dialogue": "Tell me your secrets.",
			})
		},
		"resolve_combat victory": func(id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewResolveCombatHandler(combatSvc).Handle(toolCtx, map[string]any{
				"outcome_type": "victory", "combat_state": combatState(foe, id), "xp_earned": 500,
				"dead_npc_ids": []any{foe.String(), id.String()},
				"loot":         []any{map[string]any{"name": "Coin", "description": "A coin.", "item_type": "misc"}},
			})
		},
		"resolve_combat surrender": func(id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewResolveCombatHandler(combatSvc).Handle(toolCtx, map[string]any{
				"outcome_type": "surrender", "combat_state": combatState(id),
				"surrender_npc_ids": []any{id.String()}, "disposition_change": 50,
			})
		},
		"create_quest": func(id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewCreateQuestHandler(worldSvc).Handle(toolCtx, map[string]any{
				"title": "Find the envoy", "description": "Track them down.", "quest_type": "short_term",
				"objectives":       []any{map[string]any{"description": "Ask around", "order_index": 0}},
				"related_entities": []any{map[string]any{"entity_type": "npc", "entity_id": id.String()}},
			})
		},
		"link_quest_entity": func(id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewLinkQuestEntityHandler(worldSvc).Handle(toolCtx, map[string]any{
				"quest_id": id.String(), "entity_type": "npc", "entity_id": foe.String(),
				"link_description": "Knows the way.",
			})
		},
		"create_faction": func(id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewCreateFactionHandler(worldSvc, nil, nil).Handle(toolCtx, map[string]any{
				"name": "Wardens", "description": "Keepers.", "agenda": "Guard.", "territory": "Hills",
				"properties":    map[string]any{},
				"relationships": []any{map[string]any{"faction_id": id.String(), "type": "rival", "description": "Old feud."}},
			})
		},
		"create_location": func(id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewCreateLocationHandler(worldSvc, nil, nil).Handle(toolCtx, map[string]any{
				"name": "Tower", "description": "Tall.", "region": "Hills", "location_type": "building",
				"connections": []any{map[string]any{"location_id": id.String(), "description": "Path", "bidirectional": true}},
			})
		},
		"remove_item": func(id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewRemoveItemHandler(invSvc).Handle(toolCtx, map[string]any{"item_id": id.String(), "quantity": 1})
		},
		"modify_item": func(id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewModifyItemHandler(invSvc).Handle(toolCtx, map[string]any{
				"item_id": id.String(), "properties": map[string]any{"charges": 9},
			})
		},
	}
	otherIDs := map[string]uuid.UUID{
		"link_quest_entity": otherQuest,
		"create_faction":    otherFaction,
		"create_location":   otherLocation,
		"remove_item":       otherItem,
		"modify_item":       otherItem,
	}

	outcome := func(name string, id uuid.UUID) string {
		result, err := calls[name](id)
		if err != nil {
			return "error: " + strings.ReplaceAll(err.Error(), id.String(), "<id>")
		}
		return "success: " + result.Narrative
	}
	for name := range calls {
		t.Run(name, func(t *testing.T) {
			id, ok := otherIDs[name]
			if !ok {
				id = otherNPC
			}
			unknown := outcome(name, uuid.New())
			if !strings.HasPrefix(unknown, "error: ") {
				t.Fatalf("unknown id accepted: %s", unknown)
			}
			if got := outcome(name, id); got != unknown {
				t.Fatalf("other campaign's id: got %q, want the unknown-id outcome %q", got, unknown)
			}
			if got := read(); got != before {
				t.Fatalf("rejected call changed state:\n got %+v\nwant %+v", got, before)
			}
		})
	}

	rejected := map[string]map[string]any{
		// A valid foe and XP do not save when another listed NPC is already dead.
		"dead npc": {
			"outcome_type": "victory", "combat_state": combatState(foe, deadFoe), "xp_earned": 500,
			"dead_npc_ids": []any{foe.String(), deadFoe.String()},
		},
		"other campaign's location": {
			"outcome_type": "flee", "combat_state": combatState(foe), "flee_location_id": otherLocation.String(),
		},
		"other campaign's combat state": {
			"outcome_type": "victory", "xp_earned": 500, "dead_npc_ids": []any{foe.String()},
			"combat_state": func() map[string]any { s := combatState(foe); s["campaign_id"] = other.String(); return s }(),
		},
		"other campaign's player": {
			"outcome_type": "victory", "xp_earned": 500, "dead_npc_ids": []any{foe.String()},
			"combat_state": func() map[string]any {
				s := combatState(foe)
				s["combatants"].([]any)[0].(map[string]any)["entity_id"] = otherPlayer.String()
				return s
			}(),
		},
	}
	for name, args := range rejected {
		t.Run("resolve_combat rejects "+name, func(t *testing.T) {
			if _, err := tools.NewResolveCombatHandler(combatSvc).Handle(toolCtx, args); err == nil {
				t.Fatal("accepted")
			}
			if got := read(); got != before {
				t.Fatalf("rejected call changed state:\n got %+v\nwant %+v", got, before)
			}
		})
	}

	// Related entities in the current campaign still link.
	if _, err := tools.NewCreateQuestHandler(worldSvc).Handle(toolCtx, map[string]any{
		"title": "Defeat the bandit", "description": "End the raids.", "quest_type": "short_term",
		"objectives":       []any{map[string]any{"description": "Find the camp", "order_index": 0}},
		"related_entities": []any{map[string]any{"entity_type": "npc", "entity_id": foe.String()}},
	}); err != nil {
		t.Fatalf("create_quest with a current-campaign entity: %v", err)
	}
	if _, err := tools.NewCreateQuestHandler(worldSvc).Handle(toolCtx, map[string]any{
		"title": "Guild work", "description": "Join.", "quest_type": "short_term",
		"objectives":       []any{map[string]any{"description": "Apply", "order_index": 0}},
		"related_entities": []any{map[string]any{"entity_type": "faction", "entity_id": otherFaction.String()}},
	}); err == nil {
		t.Fatal("create_quest linked another campaign's faction")
	}

	// A valid victory saves the clamped HP, the XP, the loot and the kill together.
	if _, err := tools.NewResolveCombatHandler(combatSvc).Handle(toolCtx, map[string]any{
		"outcome_type": "victory", "combat_state": combatState(foe), "xp_earned": 50,
		"dead_npc_ids": []any{foe.String()},
		"loot":         []any{map[string]any{"name": "Coin", "description": "A coin.", "item_type": "misc"}},
	}); err != nil {
		t.Fatalf("valid victory: %v", err)
	}
	after := read()
	if after.hp != 3 || after.maxHP != 12 || after.xp != 50 || after.items != 1 || after.foeAlive || after.status != "active" {
		t.Fatalf("valid victory saved %+v", after)
	}
	_ = otherLocation
}
