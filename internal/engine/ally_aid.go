package engine

import (
	"context"
	"fmt"

	"git.subcult.tv/subculture-collective/edda/internal/db"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// resolveAllyAid runs one round outside encounters in which a prepared,
// locked ally helps the Unconscious player: the Help action to stabilize, or
// a Potion of Healing from the player's inventory. Healing from 0 HP ends
// Unconscious through ApplyConditionChange (Prone remains) and clears the
// counts. A player still dying afterwards rolls the Death Saving Throw that
// starts their next turn. The caller advances the clock and records history.
func resolveAllyAid(ctx context.Context, tx pgx.Tx, campaignID uuid.UUID, pc domain.PlayerCharacter, prepared *game.StoredAction, request PreparedAction, dice srd521.Dice) (*srd521.ZeroHPReceipt, error) {
	helper := prepared.HealingTarget(request.TargetID)
	if helper == nil || pc.HP != 0 {
		return nil, ErrActionRefresh
	}
	saves, err := srd521.DecodeDeathSaves(pc.Stats)
	if err != nil {
		return nil, ErrActionRefresh
	}
	var r *srd521.ZeroHPReceipt
	if request.Intent.EntryID == srd521.AllyAidMedicine {
		if r, err = srd521.ResolveAllyStabilize(saves, pc.MaxHP, helper.Medicine, dice); err != nil {
			return nil, err
		}
		r.Reason = helper.Name + " takes the Help action. " + r.Reason
		if err := saveDeathSaves(ctx, tx, campaignID, pc.ID, r.After); err != nil {
			return nil, err
		}
	} else {
		itemID, err := uuid.Parse(request.Intent.EntryID)
		if err != nil {
			return nil, ErrActionRefresh
		}
		healing, err := srd521.ResolveAdministeredPotion(0, pc.MaxHP, dice)
		if err != nil {
			return nil, err
		}
		tag, err := tx.Exec(ctx, `UPDATE items SET quantity=quantity-1,updated_at=clock_timestamp() WHERE id=$1 AND campaign_id=$2 AND player_character_id=$3 AND quantity=$4 AND quantity>0`, db.ToPgUUID(itemID), db.ToPgUUID(campaignID), db.ToPgUUID(pc.ID), prepared.Quantity)
		if err != nil {
			return nil, err
		}
		if tag.RowsAffected() != 1 {
			return nil, ErrActionRefresh
		}
		healing.ItemID, healing.TargetID, healing.Target = itemID.String(), pc.ID.String(), pc.Name
		healing.QuantityBefore, healing.QuantityAfter = prepared.Quantity, prepared.Quantity-1
		if _, err := tx.Exec(ctx, `UPDATE player_characters SET hp=$1,updated_at=clock_timestamp() WHERE id=$2 AND campaign_id=$3`, healing.HPAfter, pc.ID, campaignID); err != nil {
			return nil, err
		}
		if err := saveDeathSaves(ctx, tx, campaignID, pc.ID, srd521.DeathSaves{}); err != nil {
			return nil, err
		}
		r = &srd521.ZeroHPReceipt{RulesVersion: "srd-5.2.1", Event: srd521.ZeroHPAllyHealing, Outcome: "healed", HPAfter: healing.HPAfter, MaxHP: pc.MaxHP, Before: saves, Healing: healing, SourceIDs: []string{"srd-5.2.1:content/items/potion-of-healing", "srd-5.2.1:content/damage/healing", "srd-5.2.1:content/damage/falling-unconscious"}}
		r.Reason = fmt.Sprintf("%s administers a Potion of Healing: %s regains %d HP (0 to %d). Unconscious ends and Prone remains. One potion consumed; %d remain.", helper.Name, pc.Name, healing.Healing, healing.HPAfter, healing.QuantityAfter)
		if r.Condition, err = ApplyConditionChange(ctx, tx, campaignID, pc.ID, nil, srd521.ConditionChange{Condition: srd521.ConditionUnconscious, Apply: false}, dice); err != nil {
			return nil, err
		}
	}
	r.HelperID, r.Helper = helper.ID.String(), helper.Name
	r.TimeBefore, r.TimeAfter = prepared.Clock, prepared.Clock+6
	if r.HPAfter == 0 && !r.After.Stable {
		if r.DeathSave, err = rollDeathSave(ctx, tx, campaignID, pc.ID, nil, dice); err != nil {
			return nil, err
		}
		if r.DeathSave != nil {
			r.Reason += " Your next turn starts: " + r.DeathSave.Reason
		}
	}
	if err := saveSpellClock(ctx, tx, campaignID, r.TimeAfter); err != nil {
		return nil, err
	}
	return r, nil
}
