package game

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/internal/tools"
)

// npcWriteRecorder lets NPC writes succeed for any row, so a test can tell
// whether a tool reached them.
type npcWriteRecorder struct {
	*mockQuerier
	npcWrites []string
}

func (r *npcWriteRecorder) UpdateNPC(_ context.Context, arg statedb.UpdateNPCParams) (statedb.Npc, error) {
	r.npcWrites = append(r.npcWrites, "update")
	return r.npcByID[arg.ID.Bytes].npc, nil
}

func (r *npcWriteRecorder) UpdateNPCDisposition(_ context.Context, arg statedb.UpdateNPCDispositionParams) (statedb.Npc, error) {
	r.npcWrites = append(r.npcWrites, "disposition")
	return r.npcByID[arg.ID.Bytes].npc, nil
}

func (r *npcWriteRecorder) KillNPC(_ context.Context, arg statedb.KillNPCParams) (statedb.Npc, error) {
	record, ok := r.npcByID[arg.ID.Bytes]
	if !ok {
		return statedb.Npc{}, pgx.ErrNoRows
	}
	r.npcWrites = append(r.npcWrites, "kill")
	return record.npc, nil
}

func TestNPCToolsTreatSceneObjectsAsUnknown(t *testing.T) {
	campaignID, locationID, playerID, questID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	objects := map[string]string{
		"door":         `{"srd521_door":{"version":"door-v1"}}`,
		"magic object": `{"srd521_magic":{"version":"detect-magic-v1","bearer":"object"}}`,
	}

	combatState := func(npcID uuid.UUID) map[string]any {
		return map[string]any{
			"id":          uuid.New().String(),
			"campaign_id": campaignID.String(),
			"status":      "active",
			"combatants": []any{map[string]any{
				"entity_id": npcID.String(), "entity_type": "npc", "name": "Target",
				"hp": 1, "max_hp": 1, "status": "alive",
			}},
		}
	}

	calls := map[string]func(context.Context, *npcWriteRecorder, uuid.UUID) (*tools.ToolResult, error){
		"npc_dialogue": func(ctx context.Context, q *npcWriteRecorder, id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewNPCDialogueHandler(NewNPCService(q)).Handle(ctx, map[string]any{
				"npc_id": id.String(), "dialogue": "Who goes there?",
			})
		},
		"update_npc": func(ctx context.Context, q *npcWriteRecorder, id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewUpdateNPCHandler(NewNPCService(q)).Handle(ctx, map[string]any{
				"npc_id": id.String(), "description": "A talkative door.",
			})
		},
		"resolve_combat victory": func(ctx context.Context, q *npcWriteRecorder, id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewResolveCombatHandler(NewCombatService(q)).Handle(ctx, map[string]any{
				"outcome_type": "victory", "combat_state": combatState(id),
				"dead_npc_ids": []any{id.String()},
			})
		},
		"resolve_combat surrender": func(ctx context.Context, q *npcWriteRecorder, id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewResolveCombatHandler(NewCombatService(q)).Handle(ctx, map[string]any{
				"outcome_type": "surrender", "combat_state": combatState(id),
				"surrender_npc_ids": []any{id.String()}, "disposition_change": 10,
			})
		},
		"establish_relationship": func(ctx context.Context, q *npcWriteRecorder, id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewEstablishRelationshipHandler(NewWorldService(q)).Handle(ctx, map[string]any{
				"source_entity_type": "npc", "source_entity_id": id.String(),
				"target_entity_type": "location", "target_entity_id": locationID.String(),
				"relationship_type": "guards", "description": "Guards the hall.",
			})
		},
		"link_quest_entity": func(ctx context.Context, q *npcWriteRecorder, id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewLinkQuestEntityHandler(NewWorldService(q)).Handle(ctx, map[string]any{
				"quest_id": questID.String(), "entity_type": "npc", "entity_id": id.String(),
				"link_description": "Knows the way.",
			})
		},
		"create_lore": func(ctx context.Context, q *npcWriteRecorder, id uuid.UUID) (*tools.ToolResult, error) {
			return tools.NewCreateLoreHandler(NewWorldService(q), nil, nil).Handle(ctx, map[string]any{
				"content": "An old tale.", "category": "legend",
				"related_entities": []any{map[string]any{"entity_type": "npc", "entity_id": id.String()}},
			})
		},
	}

	// run returns the tool's outcome with the NPC ID masked, so the outcomes for
	// an object row and for a missing row can be compared exactly.
	run := func(t *testing.T, name, properties string, id uuid.UUID) (string, []string) {
		t.Helper()
		q := &npcWriteRecorder{mockQuerier: newMockQuerier()}
		q.location = statedb.Location{ID: dbutil.ToPgtype(locationID), CampaignID: dbutil.ToPgtype(campaignID)}
		q.quests = []statedb.Quest{{ID: dbutil.ToPgtype(questID), CampaignID: dbutil.ToPgtype(campaignID)}}
		q.playerCharacters = []statedb.PlayerCharacter{{ID: dbutil.ToPgtype(playerID), CampaignID: dbutil.ToPgtype(campaignID)}}
		q.playerCharacter = q.playerCharacters[0]
		if properties != "" {
			q.npcByID[id] = mockNPCRecord{npc: statedb.Npc{
				ID: dbutil.ToPgtype(id), CampaignID: dbutil.ToPgtype(campaignID),
				LocationID: dbutil.ToPgtype(locationID), Name: "Oak door", Alive: true,
				Properties: []byte(properties),
			}}
		}
		ctx := tools.WithCurrentLocationID(tools.WithCurrentCampaignID(context.Background(), campaignID), locationID)
		ctx = tools.WithCurrentPlayerCharacterID(ctx, playerID)

		result, err := calls[name](ctx, q, id)
		outcome := ""
		switch {
		case err != nil:
			outcome = "error: " + err.Error()
		default:
			outcome = "success: " + result.Narrative
			if updated, ok := result.Data["updated_npcs"].([]map[string]any); ok && len(updated) > 0 {
				outcome += " (updated npcs)"
			}
		}
		return strings.ReplaceAll(outcome, id.String(), "<id>"), q.npcWrites
	}

	for name := range calls {
		unknown, _ := run(t, name, "", uuid.New())
		for kind, properties := range objects {
			t.Run(name+"/"+kind, func(t *testing.T) {
				got, writes := run(t, name, properties, uuid.New())
				if got != unknown {
					t.Errorf("outcome for object row = %q, want the unknown-NPC outcome %q", got, unknown)
				}
				if len(writes) != 0 {
					t.Errorf("npc writes = %v, want none", writes)
				}
			})
		}
	}
}
