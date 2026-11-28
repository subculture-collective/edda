package engine

import (
	"context"
	"encoding/json"

	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// confirmLight saves one Light cast or dismissal under the caller's locks and
// transaction. It writes only stats.srd521.casting.light, so concentration,
// slots and unknown stored fields are unchanged; the receipt and history are
// written by the caller in the same transaction.
func confirmLight(ctx context.Context, tx pgx.Tx, campaign uuid.UUID, pc domain.PlayerCharacter, sheet *srd521.CharacterSheet, prepared *game.StoredAction, request PreparedAction) (*srd521.LightReceipt, error) {
	var envelope struct {
		Build srd521.Build `json:"srd521"`
	}
	if err := json.Unmarshal(pc.Stats, &envelope); err != nil {
		return nil, err
	}
	if request.Intent.Kind == srd521.DismissLight {
		receipt, err := srd521.EndLight(&envelope.Build, prepared.LightItem, prepared.Clock)
		if err != nil {
			return nil, ErrActionRefresh
		}
		if _, err := tx.Exec(ctx, `UPDATE player_characters SET stats=stats #- '{srd521,casting,light}',updated_at=clock_timestamp() WHERE id=$1 AND campaign_id=$2`, pc.ID, campaign); err != nil {
			return nil, err
		}
		return receipt, nil
	}
	object := prepared.LightObject(request.TargetID)
	color, ok := game.LightColor(request.Intent.EntryID)
	if object == nil || !ok {
		return nil, ErrActionRefresh
	}
	receipt, err := srd521.CastLight(&envelope.Build, sheet, object.ID.String(), object.Name, color, request.ActionID.String(), prepared.Clock)
	if err != nil {
		return nil, ErrActionRefresh
	}
	if receipt.Before != nil {
		receipt.Reason = "Light on " + object.Name + " (" + color + "); your previous Light on " + prepared.LightItem + " ended. Bright light 20 feet and dim light 20 feet further for one hour; Magic action, no slot."
	}
	effect, err := json.Marshal(receipt.After)
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(ctx, `UPDATE player_characters SET stats=jsonb_set(stats,'{srd521,casting,light}',$1::jsonb),updated_at=clock_timestamp() WHERE id=$2 AND campaign_id=$3`, effect, pc.ID, campaign)
	return receipt, err
}
