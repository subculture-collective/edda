package game

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
)

func TestActiveCombatStateFromSessionLogsUsesLatestActiveCombatState(t *testing.T) {
	campaignID := uuid.New()
	firstCombatID := uuid.New()
	latestCombatID := uuid.New()

	logs := []domain.SessionLog{
		combatLog(t, campaignID, 1, "initiate_combat", firstCombatID, "active", 1),
		combatLog(t, campaignID, 2, "combat_round", latestCombatID, "active", 2),
	}

	got := activeCombatStateFromSessionLogs(logs)
	if len(got) == 0 {
		t.Fatal("expected active combat state")
	}
	if !strings.Contains(string(got), latestCombatID.String()) {
		t.Fatalf("expected latest combat state %s, got %s", latestCombatID, got)
	}
	if combatStateStatus(got) != "active" {
		t.Fatalf("expected active status, got %q", combatStateStatus(got))
	}
}

func TestActiveCombatStateFromSessionLogsClearsOnResolveCombat(t *testing.T) {
	campaignID := uuid.New()
	combatID := uuid.New()

	logs := []domain.SessionLog{
		combatLog(t, campaignID, 1, "initiate_combat", combatID, "active", 1),
		combatLog(t, campaignID, 2, "resolve_combat", combatID, "completed", 2),
	}

	if got := activeCombatStateFromSessionLogs(logs); len(got) != 0 {
		t.Fatalf("expected resolved combat to clear active state, got %s", got)
	}
}

func combatLog(t *testing.T, campaignID uuid.UUID, turn int, tool string, combatID uuid.UUID, status string, round int) domain.SessionLog {
	t.Helper()
	result := map[string]any{
		"combat_state": map[string]any{
			"id":               combatID.String(),
			"campaign_id":      campaignID.String(),
			"status":           status,
			"round_number":     round,
			"combatants":       []map[string]any{},
			"initiative_order": []string{},
			"environment": map[string]any{
				"description": "test battlefield",
			},
		},
	}
	resultJSON, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	appliedJSON, err := json.Marshal([]appliedToolLog{{Tool: tool, Result: resultJSON}})
	if err != nil {
		t.Fatalf("marshal applied: %v", err)
	}
	return domain.SessionLog{CampaignID: campaignID, TurnNumber: turn, ToolCalls: appliedJSON}
}
