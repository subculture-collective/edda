package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PreviewTravel prepares one hour of travel at pace for the campaign owner's
// character for the TUI. The destinations and their route facts come from the
// preview.
func (e *Engine) PreviewTravel(ctx context.Context, campaignID uuid.UUID, pace string) (*game.ActionPreview, error) {
	prepared, err := e.previewOwnerAction(ctx, campaignID, game.ActionIntent{Kind: srd521.TravelHour, EntryID: pace})
	if err != nil {
		return nil, err
	}
	return prepared.Preview, nil
}

// confirmTravel moves the actor along one prepared route and advances the
// campaign clock once, under the caller's locks and transaction. Durable
// expiries (Light, concentration) keep their absolute seconds; concentration
// that has ended by the arrival second is removed, as a wait or rest does.
func confirmTravel(ctx context.Context, tx pgx.Tx, campaign uuid.UUID, pc domain.PlayerCharacter, prepared *game.StoredAction, request PreparedAction) (*srd521.TravelReceipt, error) {
	option := prepared.TravelOption(request.TargetID)
	if option == nil || pc.CurrentLocationID == nil {
		return nil, ErrActionRefresh
	}
	day, err := game.ReadTravelDay(pc.Stats)
	if err != nil {
		return nil, ErrActionRefresh
	}
	plan := option.Plan
	before := day.TravelledOn(prepared.Clock)
	after := prepared.Clock + plan.Seconds
	if before+plan.Seconds > srd521.TravelDaySeconds {
		return nil, ErrActionRefresh
	}
	record, err := json.Marshal(srd521.TravelDay{Day: prepared.Clock/86400 + 1, Seconds: before + plan.Seconds})
	if err != nil {
		return nil, err
	}
	var from string
	if err := tx.QueryRow(ctx, `SELECT name FROM locations WHERE id=$1 AND campaign_id=$2`, *pc.CurrentLocationID, campaign).Scan(&from); err != nil {
		return nil, err
	}
	tag, err := tx.Exec(ctx, `UPDATE player_characters SET current_location_id=$1,stats=jsonb_set(stats,'{srd521_travel}',$2::jsonb),updated_at=clock_timestamp() WHERE id=$3 AND campaign_id=$4 AND current_location_id=$5`, option.ToID, record, pc.ID, campaign, *pc.CurrentLocationID)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() != 1 {
		return nil, ErrActionRefresh
	}
	if _, err := tx.Exec(ctx, `UPDATE locations SET player_visited=true,player_known=true,updated_at=clock_timestamp() WHERE id=$1 AND campaign_id=$2`, option.ToID, campaign); err != nil {
		return nil, err
	}
	if err := saveSpellClock(ctx, tx, campaign, after); err != nil {
		return nil, err
	}
	receipt := &srd521.TravelReceipt{RulesVersion: "srd-5.2.1", Kind: srd521.TravelHour, RouteID: option.RouteID.String(), FromID: pc.CurrentLocationID.String(), From: from, ToID: option.ToID.String(), To: option.To, Plan: *plan, TimeBefore: prepared.Clock, TimeAfter: after, DayTravelBefore: before, DayTravelAfter: before + plan.Seconds, SourceIDs: srd521.TravelSourceIDs()}
	var envelope struct {
		Build srd521.Build `json:"srd521"`
	}
	if err := json.Unmarshal(pc.Stats, &envelope); err != nil {
		return nil, err
	}
	if c := envelope.Build.Casting; c != nil && c.Concentration != nil && c.Concentration.ExpiresAt <= after {
		receipt.Expired = c.Concentration
		c.Concentration = nil
		if err := saveCaster(ctx, tx, campaign, pc.ID, &envelope.Build, pc.HP); err != nil {
			return nil, err
		}
	}
	road := "no road"
	if plan.GoodRoad {
		road = "a good road"
	}
	receipt.Reason = fmt.Sprintf("Travelled %.2f miles from %s to %s at a %s pace (%s terrain, %s): %d min %02d s; campaign time advanced once. Pace effect while travelling: %s. No checks were made during the segment.", float64(plan.DistanceFeet)/5280, from, option.To, plan.Pace, plan.Terrain, road, plan.Seconds/60, plan.Seconds%60, strings.Join(plan.Effects, "; "))
	if receipt.Expired != nil {
		receipt.Reason += " " + srd521.ConcentrationLabel(receipt.Expired.Spell) + " expired."
	}
	return receipt, nil
}
