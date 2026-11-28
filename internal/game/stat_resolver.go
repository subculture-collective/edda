package game

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"git.subcult.tv/subculture-collective/edda/internal/tools"
)

// statModifierResolver resolves stat modifiers by reading the player
// character's Stats JSONB from the database.
type statModifierResolver struct {
	queries statedb.Querier
}

// NewStatModifierResolver creates a domain.StatModifierResolver backed by the
// given Querier.
func NewStatModifierResolver(q statedb.Querier) domain.StatModifierResolver {
	return &statModifierResolver{queries: q}
}

func (r *statModifierResolver) GetStatModifier(ctx context.Context, characterID uuid.UUID, skill string) (int, error) {
	campaignID, ok := tools.CurrentCampaignIDFromContext(ctx)
	if !ok {
		return 0, fmt.Errorf("stat modifier requires current campaign id in context")
	}
	pc, err := r.queries.GetPlayerCharacterByID(ctx, dbutil.ToPgtype(characterID))
	if err == nil && pc.CampaignID != dbutil.ToPgtype(campaignID) {
		// A character in another campaign looks the same as an unknown one.
		err = pgx.ErrNoRows
	}
	if err != nil {
		return 0, fmt.Errorf("get player character: %w", err)
	}

	if len(pc.Stats) == 0 {
		return 0, nil
	}

	var stats map[string]any
	if err := json.Unmarshal(pc.Stats, &stats); err != nil {
		return 0, fmt.Errorf("unmarshal stats: %w", err)
	}

	key := strings.ToLower(skill)
	for k, v := range stats {
		if strings.ToLower(k) != key {
			continue
		}
		switch val := v.(type) {
		case float64:
			return abilityModifier(int(val)), nil
		case json.Number:
			n, err := val.Int64()
			if err != nil {
				return 0, fmt.Errorf("stat %q is not a valid integer: %w", skill, err)
			}
			return abilityModifier(int(n)), nil
		default:
			return 0, fmt.Errorf("stat %q has non-numeric value", skill)
		}
	}

	// Stat not found — return 0 modifier (no bonus/penalty).
	return 0, nil
}

// abilityModifier computes a D&D-style modifier: (score - 10) / 2, rounded down.
func abilityModifier(score int) int {
	return (score - 10) / 2
}
