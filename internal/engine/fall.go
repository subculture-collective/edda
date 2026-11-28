package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// saveConditions replaces the actor's saved out-of-encounter conditions.
func saveConditions(ctx context.Context, tx pgx.Tx, campaignID, actorID uuid.UUID, conditions game.Conditions) error {
	if conditions == (game.Conditions{}) {
		_, err := tx.Exec(ctx, `UPDATE player_characters SET stats=stats-'srd521_conditions',updated_at=clock_timestamp() WHERE id=$1 AND campaign_id=$2`, actorID, campaignID)
		return err
	}
	encoded, err := json.Marshal(conditions)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE player_characters SET stats=jsonb_set(stats,'{srd521_conditions}',$1::jsonb),updated_at=clock_timestamp() WHERE id=$2 AND campaign_id=$3`, encoded, actorID, campaignID)
	return err
}

// PreviewFall prepares the declared fall of the campaign owner's character for
// the TUI. The caller chooses only the water-entry skill: the saved declaration
// names the fall and carries its height and landing.
func (e *Engine) PreviewFall(ctx context.Context, campaignID uuid.UUID, skill string) (*game.ActionPreview, error) {
	if e.baseDB == nil {
		return nil, ErrActionRefresh
	}
	c, pc, err := e.ownerSnapshot(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	intent := game.ActionIntent{ActorID: pc.ID, Kind: srd521.FallResolve, Skill: skill}
	if pending := game.PendingFall(pc.Stats); pending != nil {
		intent.EntryID = pending.ID
	}
	prepared, err := game.PrepareStoredAction(ctx, e.baseDB, c, pc, intent, false)
	if err != nil {
		return nil, err
	}
	return prepared.Preview, nil
}

// resolveDeclaredFall rolls the optional water-entry check, then the falling
// dice, and in the same transaction consumes the declaration, saves HP and
// Prone, and applies a Concentration save. Nothing here accepts a damage amount.
func resolveDeclaredFall(ctx context.Context, tx pgx.Tx, campaignID, actorID uuid.UUID, hp int, stats []byte, fall *game.DeclaredFall, skill string, sheet *srd521.CharacterSheet, dice srd521.Dice) (*srd521.FallReceipt, error) {
	if fall == nil || fall.HeightFeet == nil {
		return nil, ErrActionRefresh
	}
	conditions := game.ReadConditions(stats)
	var water *srd521.CheckReceipt
	halved := false
	if skill != "" {
		ability := map[string]string{"athletics": "strength", "acrobatics": "dexterity"}[skill]
		if ability == "" || fall.Landing != srd521.FallWater {
			return nil, ErrActionRefresh
		}
		check, err := srd521.ResolveCheck(sheet, srd521.Check{Kind: srd521.AbilityCheck, Ability: ability, Skill: skill}, srd521.Difficulty{Value: srd521.FallWaterDC, PlayerKnown: true, OutcomeKnown: true}, dice)
		if err != nil {
			return nil, err
		}
		halved = check.Success()
		receipt := check.PlayerReceipt()
		water = &receipt
	}
	r, err := srd521.ResolveFall(srd521.Fall{HeightFeet: *fall.HeightFeet, Landing: fall.Landing, Halved: halved, HP: hp, ProneBefore: conditions.Prone}, dice)
	if err != nil {
		return nil, err
	}
	r.WaterCheck = water
	if water != nil {
		r.SourceIDs = append(r.SourceIDs, water.SourceIDs...)
	}
	conditions.Prone = r.ProneAfter
	encoded, err := json.Marshal(conditions)
	if err != nil {
		return nil, err
	}
	// The declaration is consumed exactly once: a second confirmation finds no row.
	tag, err := tx.Exec(ctx, `UPDATE player_characters SET stats=CASE WHEN $1::jsonb='{}'::jsonb THEN stats-'srd521_fall'-'srd521_conditions' ELSE jsonb_set(stats-'srd521_fall','{srd521_conditions}',$1::jsonb) END,updated_at=clock_timestamp() WHERE id=$2 AND campaign_id=$3 AND stats->'srd521_fall'->>'id'=$4`, encoded, actorID, campaignID, fall.ID)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() != 1 {
		return nil, ErrActionRefresh
	}
	// The fall is one damage source: one Concentration save, or an ending at
	// zero HP, where the character falls Unconscious or dies.
	damage, err := damagePlayer(ctx, tx, campaignID, actorID, r.Damage, false, "fall-"+fall.ID, nil, dice)
	if err != nil {
		return nil, err
	}
	if len(damage.Concentration) > 0 {
		r.Concentration = damage.Concentration[0]
	}
	r.ZeroHP = damage.ZeroHP
	return r, nil
}

func fallNarrative(name string, r *srd521.FallReceipt) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s falls %d feet into %s.", name, r.HeightFeet, map[string]string{srd521.FallSolid: "solid ground", srd521.FallWater: "water"}[r.Landing])
	if c := r.WaterCheck; c != nil {
		fmt.Fprintf(&b, " Reaction: DC %d %s check, total %d: %s.", srd521.FallWaterDC, c.Skill, c.Total, c.Outcome)
	}
	if r.DiceCount == 0 {
		b.WriteString(" Under 10 feet: no damage.")
	} else {
		dice := make([]string, len(r.Dice))
		for i, d := range r.Dice {
			dice[i] = fmt.Sprint(d)
		}
		fmt.Fprintf(&b, " %dd6 (%s) = %d", r.DiceCount, strings.Join(dice, ", "), r.Rolled)
		if r.Halved {
			fmt.Fprintf(&b, ", halved to %d", r.Damage)
		}
		fmt.Fprintf(&b, " Bludgeoning damage; HP %d to %d.", r.HPBefore, r.HPAfter)
	}
	switch {
	case r.ProneAfter && !r.ProneBefore:
		fmt.Fprintf(&b, " %s has the Prone condition.", name)
	case r.ProneAfter:
		fmt.Fprintf(&b, " %s is still Prone.", name)
	case r.DiceCount > 0:
		fmt.Fprintf(&b, " %s takes no damage and is not Prone.", name)
	}
	if r.Concentration != nil {
		b.WriteString(" " + r.Concentration.Reason)
	}
	if r.ZeroHP != nil {
		b.WriteString(" " + r.ZeroHP.Reason)
	}
	return b.String()
}
