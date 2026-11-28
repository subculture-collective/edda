package engine

import (
	"context"
	"encoding/json"
	"errors"

	"git.subcult.tv/subculture-collective/edda/internal/combat"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ErrPlayerCannotAct matches PlayerCannotActError.
var ErrPlayerCannotAct = errors.New("the player character cannot act")

// PlayerCannotActError refuses a free-text turn for a player character that
// cannot act (dead, at 0 HP or Unconscious). Reason is safe to show players.
type PlayerCannotActError struct{ Reason string }

func (e PlayerCannotActError) Error() string { return e.Reason }

// Is makes errors.Is(err, ErrPlayerCannotAct) report this error.
func (e PlayerCannotActError) Is(target error) bool { return target == ErrPlayerCannotAct }

// repairUnconscious applies the Unconscious condition that an older build
// left off a character at 0 HP, through ApplyConditionChange so held objects
// drop and Concentration ends. It returns nil when no repair is due.
func repairUnconscious(ctx context.Context, tx pgx.Tx, campaignID, actorID uuid.UUID, dice srd521.Dice) (*srd521.ZeroHPReceipt, error) {
	var hp, maxHP int
	var stats []byte
	if err := tx.QueryRow(ctx, `SELECT hp,max_hp,stats FROM player_characters WHERE id=$1 AND campaign_id=$2 FOR UPDATE`, actorID, campaignID).Scan(&hp, &maxHP, &stats); err != nil {
		return nil, err
	}
	saves, err := srd521.DecodeDeathSaves(stats)
	if err != nil {
		return nil, err
	}
	conditions, err := srd521.DecodeConditions(stats)
	if err != nil {
		return nil, err
	}
	if hp != 0 || saves.Dead || conditions.Unconscious {
		return nil, nil
	}
	r := &srd521.ZeroHPReceipt{RulesVersion: "srd-5.2.1", Event: srd521.ZeroHPUnconsciousRepair, MaxHP: maxHP, Before: saves, After: saves,
		SourceIDs: []string{"srd-5.2.1:content/damage/falling-unconscious", "srd-5.2.1:content/rules-glossary/unconscious"},
		Reason:    "At 0 HP without the Unconscious condition, which earlier versions of Edda did not apply: now Unconscious and Prone. Death Saving Throws follow."}
	if r.Condition, err = ApplyConditionChange(ctx, tx, campaignID, actorID, nil, srd521.ConditionChange{Condition: srd521.ConditionUnconscious, Apply: true}, dice); err != nil {
		return nil, err
	}
	return r, nil
}

// playerDamage is what one damage source did to the player character.
type playerDamage struct {
	HP            int
	Concentration []*srd521.EffectReceipt
	ZeroHP        *srd521.ZeroHPReceipt
}

// damagePlayer applies one resolved damage source to the player character in
// the caller's transaction. It saves HP and the 0-HP state, applies the
// damage's Concentration event, and on dropping to 0 HP without dying applies
// Unconscious through ApplyConditionChange. In an active encounter the player
// actor follows, and the player's death completes the encounter. Nothing here
// reads client input: source and damage come from a server-side roll.
func damagePlayer(ctx context.Context, tx pgx.Tx, campaignID, actorID uuid.UUID, damage int, critical bool, source string, encounter *combat.Encounter, dice srd521.Dice) (*playerDamage, error) {
	var hp, maxHP int
	var stats []byte
	if err := tx.QueryRow(ctx, `SELECT hp,max_hp,stats FROM player_characters WHERE id=$1 AND campaign_id=$2 FOR UPDATE`, actorID, campaignID).Scan(&hp, &maxHP, &stats); err != nil {
		return nil, err
	}
	saves, err := srd521.DecodeDeathSaves(stats)
	if err != nil {
		return nil, err
	}
	after, zero, err := srd521.ResolvePlayerDamage(hp, maxHP, damage, critical, saves)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `UPDATE player_characters SET hp=$1,updated_at=clock_timestamp() WHERE id=$2 AND campaign_id=$3`, after, actorID, campaignID); err != nil {
		return nil, err
	}
	result := &playerDamage{HP: after, ZeroHP: zero}
	if result.Concentration, err = ApplyConcentrationEvents(ctx, tx, campaignID, actorID, hp, []srd521.ConcentrationEvent{{Kind: srd521.ConcentrationDamageSource, SourceID: source, Damage: damage, HPAfter: after}}, dice); err != nil {
		return nil, err
	}
	if zero != nil {
		zero.HPAfter = after
		if err := saveDeathSaves(ctx, tx, campaignID, actorID, zero.After); err != nil {
			return nil, err
		}
		if zero.Event == srd521.ZeroHPDropped {
			if zero.Condition, err = ApplyConditionChange(ctx, tx, campaignID, actorID, encounter, srd521.ConditionChange{Condition: srd521.ConditionUnconscious, Apply: true}, dice); err != nil {
				return nil, err
			}
		}
	}
	if encounter != nil && encounter.Status == "active" {
		encounter.SetHP(actorID, after)
		if zero != nil && zero.After.Dead {
			encounter.Status = "completed"
		}
	}
	return result, nil
}

// rollDeathSave rolls the player's Death Saving Throw, or outside an active
// encounter a Stable character's 1d4-hour recovery. It returns nil when none
// is due: the character has HP, is dead, or is Stable inside an encounter.
// Regaining 1 HP ends Unconscious and clears the counts; in an encounter the
// player's turn budget restarts with its Speed. Death completes the encounter.
func rollDeathSave(ctx context.Context, tx pgx.Tx, campaignID, actorID uuid.UUID, encounter *combat.Encounter, dice srd521.Dice) (*srd521.ZeroHPReceipt, error) {
	var hp, maxHP int
	var stats []byte
	if err := tx.QueryRow(ctx, `SELECT hp,max_hp,stats FROM player_characters WHERE id=$1 AND campaign_id=$2 FOR UPDATE`, actorID, campaignID).Scan(&hp, &maxHP, &stats); err != nil {
		return nil, err
	}
	saves, err := srd521.DecodeDeathSaves(stats)
	if err != nil {
		return nil, err
	}
	active := encounter != nil && encounter.Status == "active"
	if hp != 0 || saves.Dead || (saves.Stable && active) {
		return nil, nil
	}
	var r *srd521.ZeroHPReceipt
	if saves.Stable {
		r, err = srd521.ResolveStableRecovery(saves, maxHP, dice)
	} else {
		r, err = srd521.ResolveDeathSave(saves, maxHP, dice)
	}
	if err != nil {
		return nil, err
	}
	if err := saveDeathSaves(ctx, tx, campaignID, actorID, r.After); err != nil {
		return nil, err
	}
	if r.Regained() {
		if _, err := tx.Exec(ctx, `UPDATE player_characters SET hp=$1,updated_at=clock_timestamp() WHERE id=$2 AND campaign_id=$3`, r.HPAfter, actorID, campaignID); err != nil {
			return nil, err
		}
		if r.Condition, err = ApplyConditionChange(ctx, tx, campaignID, actorID, encounter, srd521.ConditionChange{Condition: srd521.ConditionUnconscious, Apply: false}, dice); err != nil {
			return nil, err
		}
	}
	if active {
		encounter.SetHP(actorID, r.HPAfter)
		if r.Regained() {
			encounter.ResetBudget()
		}
		if r.After.Dead {
			encounter.Status = "completed"
		}
	}
	return r, nil
}

// saveDeathSaves replaces stats.srd521_death_saves; the zero state removes it.
func saveDeathSaves(ctx context.Context, tx pgx.Tx, campaignID, actorID uuid.UUID, saves srd521.DeathSaves) error {
	if saves == (srd521.DeathSaves{}) {
		_, err := tx.Exec(ctx, `UPDATE player_characters SET stats=stats-'srd521_death_saves',updated_at=clock_timestamp() WHERE id=$1 AND campaign_id=$2`, actorID, campaignID)
		return err
	}
	encoded, err := json.Marshal(saves)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE player_characters SET stats=jsonb_set(stats,'{srd521_death_saves}',$1::jsonb),updated_at=clock_timestamp() WHERE id=$2 AND campaign_id=$3`, encoded, actorID, campaignID)
	return err
}

// PreviewDeathSave prepares the campaign owner's explicit Death Saving Throw
// (or Stable recovery) outside encounters for the TUI. It takes no choices.
func (e *Engine) PreviewDeathSave(ctx context.Context, campaignID uuid.UUID) (*game.ActionPreview, error) {
	if e.baseDB == nil {
		return nil, ErrActionRefresh
	}
	c, pc, err := e.ownerSnapshot(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	prepared, err := game.PrepareStoredAction(ctx, e.baseDB, c, pc, game.ActionIntent{ActorID: pc.ID, Kind: srd521.DeathSave}, false)
	if err != nil {
		return nil, err
	}
	return prepared.Preview, nil
}
