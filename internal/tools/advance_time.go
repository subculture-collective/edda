package tools

import (
	"context"
	"errors"
	"fmt"

	"git.subcult.tv/subculture-collective/edda/internal/db"
	"git.subcult.tv/subculture-collective/edda/internal/llm"
)

const advanceTimeToolName = "advance_time"

// TimeStore is the database interface used by time-related tools.
type TimeStore = db.DBTX

// AdvanceTimeTool returns the advance_time tool definition and JSON schema.
func AdvanceTimeTool() llm.Tool {
	return llm.Tool{
		Name:        advanceTimeToolName,
		Description: "Advance the in-game clock by the given hours and/or minutes. Handles day rollover automatically.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"hours": map[string]any{
					"type":        "integer",
					"description": "Number of hours to advance (0 or more).",
				},
				"minutes": map[string]any{
					"type":        "integer",
					"description": "Number of minutes to advance (0 or more, optional, defaults to 0).",
				},
			},
			"required":             []string{"hours"},
			"additionalProperties": false,
		},
	}
}

// RegisterAdvanceTime registers the advance_time tool and handler.
func RegisterAdvanceTime(reg *Registry, db TimeStore) error {
	if db == nil {
		return errors.New("advance_time database connection is required")
	}
	return reg.Register(AdvanceTimeTool(), NewAdvanceTimeHandler(db).Handle)
}

// AdvanceTimeHandler executes advance_time tool calls.
type AdvanceTimeHandler struct {
	db TimeStore
}

// NewAdvanceTimeHandler creates a new advance_time handler.
func NewAdvanceTimeHandler(db TimeStore) *AdvanceTimeHandler {
	return &AdvanceTimeHandler{db: db}
}

// Handle executes the advance_time tool.
func (h *AdvanceTimeHandler) Handle(ctx context.Context, args map[string]any) (*ToolResult, error) {
	if h == nil {
		return nil, errors.New("advance_time handler is nil")
	}

	hours, err := parseIntArg(args, "hours")
	if err != nil {
		return nil, err
	}
	if hours < 0 {
		return nil, errors.New("hours must be 0 or more")
	}

	minutes := 0
	if _, ok := args["minutes"]; ok {
		minutes, err = parseIntArg(args, "minutes")
		if err != nil {
			return nil, err
		}
		if minutes < 0 {
			return nil, errors.New("minutes must be 0 or more")
		}
	}

	if hours == 0 && minutes == 0 {
		return nil, errors.New("at least one of hours or minutes must be greater than 0")
	}

	campaignID, ok := CurrentCampaignIDFromContext(ctx)
	if !ok {
		return nil, errors.New("advance_time requires current campaign id in context")
	}

	day, hour, minute, err := db.AdvanceCampaignTime(ctx, h.db, campaignID, hours, minutes)
	if err != nil {
		return nil, fmt.Errorf("update campaign time: %w", err)
	}

	narrative := fmt.Sprintf("%s pass. It is now Day %d, %02d:%02d.", formatDuration(hours, minutes), day, hour, minute)

	return &ToolResult{
		Success: true,
		Data: map[string]any{
			"day":    day,
			"hour":   hour,
			"minute": minute,
		},
		Narrative: narrative,
	}, nil
}

func formatDuration(hours, minutes int) string {
	parts := make([]string, 0, 2)
	if hours > 0 {
		if hours == 1 {
			parts = append(parts, "1 hour")
		} else {
			parts = append(parts, fmt.Sprintf("%d hours", hours))
		}
	}
	if minutes > 0 {
		if minutes == 1 {
			parts = append(parts, "1 minute")
		} else {
			parts = append(parts, fmt.Sprintf("%d minutes", minutes))
		}
	}
	if len(parts) == 0 {
		return "0 minutes"
	}
	if len(parts) == 2 {
		return parts[0] + " and " + parts[1]
	}
	return parts[0]
}
