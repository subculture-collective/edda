//go:build integration

package handlers_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/game"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/internal/tools"
)

// link_quest_entity stores a quest-sourced relationship and branch_quest
// copies it to the branch, alongside its own quest-to-quest link.
func TestIntegrationQuestRelationshipsLinkAndBranch(t *testing.T) {
	ctx := context.Background()
	q := statedb.New(testPool)
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	var campaign, location, npc, quest uuid.UUID
	must(testPool.QueryRow(ctx, `INSERT INTO campaigns(name,status,created_by) VALUES('Quest links','active',$1) RETURNING id`, testUserID).Scan(&campaign))
	t.Cleanup(func() { _, _ = testPool.Exec(ctx, `DELETE FROM campaigns WHERE id=$1`, campaign) })
	must(testPool.QueryRow(ctx, `INSERT INTO locations(campaign_id,name) VALUES($1,'Square') RETURNING id`, campaign).Scan(&location))
	must(testPool.QueryRow(ctx, `INSERT INTO npcs(campaign_id,name,location_id) VALUES($1,'Guide',$2) RETURNING id`, campaign, location).Scan(&npc))
	must(testPool.QueryRow(ctx, `INSERT INTO quests(campaign_id,title,quest_type,status) VALUES($1,'Find the pass','short_term','active') RETURNING id`, campaign).Scan(&quest))

	toolCtx := tools.WithCurrentCampaignID(ctx, campaign)
	toolCtx = tools.WithCurrentLocationID(toolCtx, location)
	worldSvc := game.NewWorldService(q)

	if _, err := tools.NewLinkQuestEntityHandler(worldSvc).Handle(toolCtx, map[string]any{
		"quest_id": quest.String(), "entity_type": "npc", "entity_id": npc.String(),
		"link_description": "Knows the way.",
	}); err != nil {
		t.Fatalf("link_quest_entity: %v", err)
	}

	type rel struct {
		sourceType, targetType, relType string
		source, target                  uuid.UUID
	}
	relsOf := func(questID uuid.UUID) []rel {
		t.Helper()
		rows, err := testPool.Query(ctx, `SELECT source_entity_type,source_entity_id,target_entity_type,target_entity_id,relationship_type
			FROM entity_relationships
			WHERE (source_entity_type='quest' AND source_entity_id=$1) OR (target_entity_type='quest' AND target_entity_id=$1)
			ORDER BY relationship_type, source_entity_id, target_entity_id`, questID)
		must(err)
		defer rows.Close()
		var out []rel
		for rows.Next() {
			var r rel
			must(rows.Scan(&r.sourceType, &r.source, &r.targetType, &r.target, &r.relType))
			out = append(out, r)
		}
		must(rows.Err())
		return out
	}

	branch := func(title string) uuid.UUID {
		t.Helper()
		result, err := tools.NewBranchQuestHandler(worldSvc).Handle(toolCtx, map[string]any{
			"original_quest_id": quest.String(), "branch_title": title, "branch_description": "Another way.",
			"branch_objectives": []any{map[string]any{"description": "Climb", "order_index": 0}},
			"reason":            "The bridge fell.",
		})
		if err != nil {
			t.Fatalf("branch_quest %q: %v", title, err)
		}
		if got := result.Data["copied_relationships"]; got != 1 {
			t.Fatalf("branch_quest %q copied %v relationships, want only the NPC link", title, got)
		}
		id, err := uuid.Parse(result.Data["id"].(string))
		must(err)
		return id
	}

	for _, title := range []string{"Climb the cliff", "Take the tunnel"} {
		id := branch(title)
		want := []rel{
			{sourceType: "quest", source: id, targetType: "quest", target: quest, relType: "quest_branch"},
			{sourceType: "quest", source: id, targetType: "npc", target: npc, relType: "quest_related"},
		}
		got := relsOf(id)
		if len(got) != len(want) {
			t.Fatalf("branch %q relationships = %+v, want %+v", title, got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("branch %q relationships = %+v, want %+v", title, got, want)
			}
		}
	}
}
