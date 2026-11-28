package tools

import (
	"context"
	"errors"
	"fmt"

	"git.subcult.tv/subculture-collective/edda/internal/db"
	"git.subcult.tv/subculture-collective/edda/internal/llm"
)

const restToolName = "rest"

// RestTool returns the rest tool definition and JSON schema.
func RestTool() llm.Tool {
	return llm.Tool{
		Name:        restToolName,
		Description: "Rest the player character for a number of hours. Short rest (1 hour) restores 25% max HP; long rest (up to 8 hours) restores full HP. Advances in-game time accordingly.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"hours": map[string]any{
					"type":        "integer",
					"description": "Number of hours to rest (1-8).",
					"minimum":     1,
					"maximum":     8,
				},
				"short_rest": map[string]any{
					"type":        "boolean",
					"description": "If true, perform a short rest (1 hour, 25% HP recovery). Defaults to false (long rest, full HP recovery).",
				},
			},
			"required":             []string{"hours"},
			"additionalProperties": false,
		},
	}
}

// RegisterRest registers the rest tool and handler.
func RegisterRest(reg *Registry, db TimeStore) error {
	if db == nil {
		return errors.New("rest database connection is required")
	}
	return reg.Register(RestTool(), NewRestHandler(db).Handle)
}

// RestHandler executes rest tool calls.
type RestHandler struct {
	db TimeStore
}

// NewRestHandler creates a new rest handler.
func NewRestHandler(db TimeStore) *RestHandler {
	return &RestHandler{db: db}
}

// restPlayerHPSQL heals one character from its own row in one statement. A
// short rest ($3) restores a quarter of maximum HP, at least 1; a long rest
// restores all of it.
//
// It takes the campaign row lock before the character row, the order the SRD
// rules actions use. The lock mode matches db.AdvanceCampaignTime, so the clock
// advance that follows in the same transaction does not upgrade it.
const restPlayerHPSQL = `
WITH locked AS (
  SELECT id FROM campaigns WHERE id = $1 FOR NO KEY UPDATE
)
UPDATE player_characters AS p
SET hp = CASE WHEN $3::boolean THEN LEAST(p.max_hp, p.hp + GREATEST(p.max_hp / 4, 1)) ELSE p.max_hp END
FROM locked
WHERE p.id = $2 AND p.campaign_id = locked.id
RETURNING p.hp, p.max_hp
`

// Handle executes the rest tool.
func (h *RestHandler) Handle(ctx context.Context, args map[string]any) (*ToolResult, error) {
	if h == nil {
		return nil, errors.New("rest handler is nil")
	}

	hours, err := parseIntArg(args, "hours")
	if err != nil {
		return nil, err
	}
	if hours < 1 || hours > 8 {
		return nil, errors.New("hours must be between 1 and 8")
	}

	shortRest, err := parseBoolArg(args, "short_rest")
	if err != nil {
		return nil, err
	}

	campaignID, ok := CurrentCampaignIDFromContext(ctx)
	if !ok {
		return nil, errors.New("rest requires current campaign id in context")
	}

	playerCharacterID, ok := CurrentPlayerCharacterIDFromContext(ctx)
	if !ok {
		return nil, errors.New("rest requires current player character id in context")
	}

	var newHP, maxHP int
	err = h.db.QueryRow(ctx, restPlayerHPSQL, db.ToPgUUID(campaignID), db.ToPgUUID(playerCharacterID), shortRest).Scan(&newHP, &maxHP)
	if err != nil {
		return nil, fmt.Errorf("update player HP: %w", err)
	}

	// Advance time.
	day, hour, minute, err := db.AdvanceCampaignTime(ctx, h.db, campaignID, hours, 0)
	if err != nil {
		return nil, fmt.Errorf("update campaign time: %w", err)
	}

	restType := "long rest"
	if shortRest {
		restType = "short rest"
	}

	narrative := fmt.Sprintf("You take a %s for %s. HP restored to %d/%d. It is now Day %d, %02d:%02d.",
		restType, formatDuration(hours, 0), newHP, maxHP, day, hour, minute)

	return &ToolResult{
		Success: true,
		Data: map[string]any{
			"hp":        newHP,
			"max_hp":    maxHP,
			"rest_type": restType,
			"day":       day,
			"hour":      hour,
			"minute":    minute,
		},
		Narrative: narrative,
	}, nil
}
