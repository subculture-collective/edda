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

// confirmLight saves one Light cast, cover ruling or dismissal under the caller's locks and
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
	if kind := request.Intent.Kind; kind == srd521.CoverLight || kind == srd521.UncoverLight {
		receipt, err := srd521.SetLightCovered(&envelope.Build, prepared.LightItem, kind, prepared.Clock)
		if err != nil {
			return nil, ErrActionRefresh
		}
		return receipt, saveLight(ctx, tx, campaign, pc.ID, receipt.After)
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
	return receipt, saveLight(ctx, tx, campaign, pc.ID, receipt.After)
}

func saveLight(ctx context.Context, tx pgx.Tx, campaign, actor uuid.UUID, light *srd521.LightEffect) error {
	effect, err := json.Marshal(light)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE player_characters SET stats=jsonb_set(stats,'{srd521,casting,light}',$1::jsonb),updated_at=clock_timestamp() WHERE id=$2 AND campaign_id=$3`, effect, actor, campaign)
	return err
}

// PreviewLight prepares one Light intent for the campaign owner's character for
// the TUI. The stored action also carries the active Light and its object name.
func (e *Engine) PreviewLight(ctx context.Context, campaignID uuid.UUID, intent game.ActionIntent) (*game.StoredAction, error) {
	if !game.IsLightIntent(intent) {
		return nil, ErrActionRefresh
	}
	return e.previewOwnerAction(ctx, campaignID, intent)
}

// PreviewDetectMagic prepares the Detect Magic free-use cast or its study for
// the campaign owner's character for the TUI. A study preview also carries any
// live concentration.
func (e *Engine) PreviewDetectMagic(ctx context.Context, campaignID uuid.UUID, intent game.ActionIntent) (*game.StoredAction, error) {
	if intent.Kind != srd521.DetectMagicStudy && (intent.Kind != srd521.SpellCast || intent.EntryID != srd521.DetectMagicID) {
		return nil, ErrActionRefresh
	}
	return e.previewOwnerAction(ctx, campaignID, intent)
}

// previewOwnerAction prepares an intent, without writing, for the campaign
// owner's first active character.
func (e *Engine) previewOwnerAction(ctx context.Context, campaignID uuid.UUID, intent game.ActionIntent) (*game.StoredAction, error) {
	if e.baseDB == nil {
		return nil, ErrActionRefresh
	}
	c, pc, err := e.ownerSnapshot(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	intent.ActorID = pc.ID
	return game.PrepareStoredAction(ctx, e.baseDB, c, pc, intent, false)
}

// ownerSnapshot reads the campaign owner's first active character.
func (e *Engine) ownerSnapshot(ctx context.Context, campaignID uuid.UUID) (domain.Campaign, domain.PlayerCharacter, error) {
	var actorID uuid.UUID
	if err := e.baseDB.QueryRow(ctx, `SELECT p.id FROM player_characters p JOIN campaigns c ON c.id=p.campaign_id WHERE c.id=$1 AND p.user_id=c.created_by AND p.status='active' ORDER BY p.id LIMIT 1`, campaignID).Scan(&actorID); err != nil {
		return domain.Campaign{}, domain.PlayerCharacter{}, err
	}
	return game.ReadActionSnapshot(ctx, e.baseDB, campaignID, actorID, uuid.Nil, false)
}
