package engine

import (
	"context"
	"fmt"
	"time"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
)

// TurnCompleter owns the deterministic post-LLM turn completion work: turning
// raw narrative plus applied tool calls into one persisted, player-facing turn.
type TurnCompleter struct {
	engine *Engine
}

func NewTurnCompleter(engine *Engine) *TurnCompleter {
	return &TurnCompleter{engine: engine}
}

func (c *TurnCompleter) CompleteAndPersist(ctx context.Context, tc *TurnContext) error {
	if c == nil || c.engine == nil {
		return fmt.Errorf("turn completer is not configured")
	}

	cleaned, choices, err := extractChoicesStrict(tc.Narrative)
	if err != nil {
		tc.Logger.Error("process turn failed during choice extraction",
			"campaign_id", tc.CampaignID,
			"duration_ms", time.Since(tc.Started).Milliseconds(),
			"error", err,
		)
		return fmt.Errorf("complete turn: %w", err)
	}
	tc.Narrative, tc.Choices = cleaned, choices
	tc.CombatActive = combatActiveAfterApplied(tc.State.CombatActive, tc.Applied)
	tc.StateChanges = StateChangesFromAppliedToolCalls(tc.Applied)

	toolCallsJSON, err := marshalAppliedToolCalls(tc.Applied)
	if err != nil {
		tc.Logger.Error("process turn failed during tool-call marshal",
			"campaign_id", tc.CampaignID,
			"duration_ms", time.Since(tc.Started).Milliseconds(),
			"error", err,
		)
		return fmt.Errorf("marshal applied tool calls: %w", err)
	}
	tc.ToolCallsJSON = toolCallsJSON

	tc.TurnNumber = nextTurnNumber(tc.RecentLogs)
	log := domain.SessionLog{
		CampaignID:  tc.CampaignID,
		TurnNumber:  tc.TurnNumber,
		PlayerInput: tc.PlayerInput,
		InputType:   domain.Classify(tc.PlayerInput),
		LLMResponse: tc.Narrative,
		ToolCalls:   toolCallsJSON,
		LocationID:  finalLocationIDFromApplied(tc.State.Player.CurrentLocationID, tc.Applied),
	}
	if err := c.engine.state.SaveSessionLog(ctx, log); err != nil {
		tc.Logger.Error("process turn failed during session-log save",
			"campaign_id", tc.CampaignID,
			"duration_ms", time.Since(tc.Started).Milliseconds(),
			"error", err,
		)
		return fmt.Errorf("save session log: %w", err)
	}

	c.engine.snapshotQuestsIfNeeded(ctx, tc.CampaignID, tc.Applied)
	c.engine.autoSaveIfNeeded(ctx, tc.CampaignID, tc.TurnNumber)
	c.engine.autoSummarizeIfNeeded(ctx, tc.CampaignID, tc.TurnNumber)

	return nil
}

func combatActiveAfterApplied(initial bool, applied []AppliedToolCall) bool {
	combatActive := initial
	for _, atc := range applied {
		switch atc.Tool {
		case "initiate_combat":
			combatActive = true
		case "resolve_combat":
			combatActive = false
		}
	}
	return combatActive
}
