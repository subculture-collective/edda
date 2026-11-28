package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/tools"
)

// gatherStage loads the game state and recent session logs for the campaign.
// It also enriches tc.Ctx with campaign, player, and location IDs so that
// downstream stages (especially tool execution) can access them.
func (e *Engine) gatherStage() Stage {
	return func(ctx context.Context, tc *TurnContext) error {
		state, err := e.state.GatherState(ctx, tc.CampaignID)
		if err != nil {
			tc.Logger.Error("process turn failed during state gather",
				"campaign_id", tc.CampaignID,
				"duration_ms", time.Since(tc.Started).Milliseconds(),
				"error", err,
			)
			return fmt.Errorf("gather state: %w", err)
		}
		// An SRD 5.2.1 character that is dead, at 0 HP or Unconscious cannot
		// act, so no free-text turn reaches the model. The reason replaces the
		// generic guard below, which refuses SRD campaigns' free-text turns.
		if state.Campaign.RulesVersion == domain.RulesVersionSRD521 && state.Player.ID != (uuid.UUID{}) && e.baseDB != nil {
			encounter, err := game.ReadEncounter(ctx, e.baseDB, tc.CampaignID)
			if err != nil {
				return fmt.Errorf("read encounter: %w", err)
			}
			if reason := game.ZeroHPTurnReason(state.Player.Stats, state.Player.HP, encounter); reason != "" {
				return PlayerCannotActError{Reason: reason}
			}
		}
		if err := domain.RequireLegacyRules(state.Campaign.RulesVersion); err != nil {
			return err
		}
		tc.State = state
		tc.Logger.Debug("state gathered",
			"campaign_id", tc.CampaignID,
			"player_id", state.Player.ID,
			"has_location", state.Player.CurrentLocationID != nil,
		)

		// Enrich context with campaign/player/location IDs for tool execution.
		enriched := tools.WithCurrentCampaignID(ctx, tc.CampaignID)
		if state.Player.ID != (uuid.UUID{}) {
			enriched = tools.WithCurrentPlayerCharacterID(enriched, state.Player.ID)
		}
		if state.Player.CurrentLocationID != nil {
			enriched = tools.WithCurrentLocationID(enriched, *state.Player.CurrentLocationID)
		}
		tc.Ctx = enriched

		recentTurns, err := e.state.ListRecentSessionLogs(enriched, tc.CampaignID, recentTurnLimit)
		if err != nil {
			tc.Logger.Error("process turn failed during session-log fetch",
				"campaign_id", tc.CampaignID,
				"duration_ms", time.Since(tc.Started).Milliseconds(),
				"error", err,
			)
			return fmt.Errorf("list recent session logs: %w", err)
		}
		tc.RecentLogs = recentTurns

		return e.startJournal(enriched, tc)
	}
}
