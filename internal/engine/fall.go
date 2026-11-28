package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"git.subcult.tv/subculture-collective/edda/internal/combat"
	"git.subcult.tv/subculture-collective/edda/internal/game"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// saveConditions replaces the actor's saved out-of-encounter conditions.
func saveConditions(ctx context.Context, tx pgx.Tx, campaignID, actorID uuid.UUID, conditions game.Conditions) error {
	if !conditions.Prone {
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
	tag, err := tx.Exec(ctx, `UPDATE player_characters SET hp=$1,stats=CASE WHEN $2::jsonb='{}'::jsonb THEN stats-'srd521_fall'-'srd521_conditions' ELSE jsonb_set(stats-'srd521_fall','{srd521_conditions}',$2::jsonb) END,updated_at=clock_timestamp() WHERE id=$3 AND campaign_id=$4 AND stats->'srd521_fall'->>'id'=$5`, r.HPAfter, encoded, actorID, campaignID, fall.ID)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() != 1 {
		return nil, ErrActionRefresh
	}
	if r.Damage > 0 && sheet.Spellcasting != nil {
		r.Concentration, err = concentrationAfterDamage(ctx, tx, campaignID, &combat.EncounterActor{ID: actorID, HP: hp}, sheet, r.Damage, r.HPAfter, dice)
		if err != nil {
			return nil, err
		}
	}
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
	if r.HPAfter == 0 {
		b.WriteString(" At 0 HP, falling Unconscious and death saving throws require adjudication.")
	}
	if r.Concentration != nil {
		b.WriteString(" " + r.Concentration.Reason)
	}
	return b.String()
}
