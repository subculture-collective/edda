package engine

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// completeStage finalizes a processed turn: choice cleanup, state projections,
// session-log persistence, quest snapshots, auto-save, and auto-summary.
func (e *Engine) completeStage() Stage {
	return func(ctx context.Context, tc *TurnContext) error {
		started := time.Now()
		if err := NewTurnCompleter(e, e.choiceProvider).CompleteAndPersist(ctx, tc); err != nil {
			return err
		}
		tc.Logger.Debug("turn completion finished", "campaign_id", tc.CampaignID, "duration_ms", time.Since(started).Milliseconds())
		return nil
	}
}

func finalLocationIDFromApplied(defaultLocationID *uuid.UUID, applied []AppliedToolCall) *uuid.UUID {
	locationID := defaultLocationID
	for _, call := range applied {
		var data map[string]any
		if err := json.Unmarshal(call.Result, &data); err != nil {
			continue
		}

		switch call.Tool {
		case "move_player":
			if id, ok := uuidStringField(data, "location_id"); ok {
				locationID = &id
			}
		case "create_location":
			if moved, _ := data["move_player_here"].(bool); moved {
				if id, ok := uuidStringField(data, "location_id"); ok {
					locationID = &id
					continue
				}
				if id, ok := uuidStringField(data, "id"); ok {
					locationID = &id
				}
			}
		}
	}
	return locationID
}

func uuidStringField(data map[string]any, key string) (uuid.UUID, bool) {
	raw, ok := data[key]
	if !ok {
		return uuid.UUID{}, false
	}
	str, ok := raw.(string)
	if !ok {
		return uuid.UUID{}, false
	}
	id, err := uuid.Parse(str)
	return id, err == nil
}
