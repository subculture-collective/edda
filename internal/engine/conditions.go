package engine

import (
	"context"
	"errors"
	"fmt"

	"git.subcult.tv/subculture-collective/edda/internal/combat"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ErrConditionActor reports a condition change for an actor outside the
// campaign or the given encounter, or in a campaign without SRD 5.2.1 rules.
var ErrConditionActor = errors.New("condition actor is not in this campaign")

// ApplyConditionChange is the trusted adapter that applies or ends Incapacitated
// or Unconscious on a player character inside the caller's transaction. Callers
// build change from a result they already resolved (for example dropping to 0
// HP); nothing here reads client input, and the caller writes the public
// history entry with the returned receipt.
//
// It locks the actor row. While encounter is active the player actor's
// conditions live on it and the caller persists the encounter; otherwise they
// are saved in stats.srd521_conditions. Becoming Incapacitated ends
// Concentration through ApplyConcentrationEvents. Becoming Unconscious also
// drops held objects: equipped weapons and a Shield, which also leaves the
// build. Any other equipped item that is not worn armor, unknown condition
// state, and a foreign actor fail closed with no change.
func ApplyConditionChange(ctx context.Context, tx pgx.Tx, campaignID, actorID uuid.UUID, encounter *combat.Encounter, change srd521.ConditionChange, dice srd521.Dice) (*srd521.ConditionReceipt, error) {
	var hp int
	var stats []byte
	err := tx.QueryRow(ctx, `SELECT p.hp,p.stats FROM player_characters p JOIN campaigns c ON c.id=p.campaign_id WHERE p.id=$1 AND p.campaign_id=$2 AND c.rules_version=$3 FOR UPDATE OF p`, actorID, campaignID, string(domain.RulesVersionSRD521)).Scan(&hp, &stats)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrConditionActor
	}
	if err != nil {
		return nil, err
	}
	var actor *combat.EncounterActor
	if encounter != nil && encounter.Status == "active" {
		for i := range encounter.Actors {
			if encounter.Actors[i].Player && encounter.Actors[i].ID == actorID {
				actor = &encounter.Actors[i]
			}
		}
		if actor == nil {
			return nil, ErrConditionActor
		}
	}
	var before srd521.Conditions
	if actor != nil {
		before = actor.Conditions
	} else if before, err = srd521.DecodeConditions(stats); err != nil {
		return nil, err
	}
	receipt, events, err := srd521.ResolveConditionChange(before, change)
	if err != nil {
		return nil, err
	}
	if !receipt.Changed {
		return receipt, nil
	}
	if receipt.After.Unconscious && !before.Unconscious {
		if receipt.Dropped, err = dropHeldObjects(ctx, tx, campaignID, actorID); err != nil {
			return nil, err
		}
	}
	if len(events) > 0 {
		if receipt.Concentration, err = ApplyConcentrationEvents(ctx, tx, campaignID, actorID, hp, events, dice); err != nil {
			return nil, err
		}
	}
	if actor != nil {
		actor.Conditions = receipt.After
		return receipt, nil
	}
	return receipt, saveConditions(ctx, tx, campaignID, actorID, receipt.After)
}

// dropHeldObjects unequips the actor's held items and returns their IDs. An
// equipped item that is neither a weapon, a Shield nor worn armor is not known
// to be held or worn, so the change fails closed.
func dropHeldObjects(ctx context.Context, tx pgx.Tx, campaignID, actorID uuid.UUID) ([]string, error) {
	rows, err := tx.Query(ctx, `SELECT id,item_type,COALESCE(properties->>'srd521_id','') FROM items WHERE campaign_id=$1 AND player_character_id=$2 AND equipped ORDER BY id FOR UPDATE`, campaignID, actorID)
	if err != nil {
		return nil, err
	}
	dropped := []string{}
	shield := false
	for rows.Next() {
		var id uuid.UUID
		var kind, key string
		if err := rows.Scan(&id, &kind, &key); err != nil {
			rows.Close()
			return nil, err
		}
		switch {
		case kind == "weapon":
		case kind == "armor" && key == "shield":
			shield = true
		case kind == "armor":
			continue
		default:
			rows.Close()
			return nil, fmt.Errorf("an equipped %s item may be held; dropping it requires adjudication", kind)
		}
		dropped = append(dropped, id.String())
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(dropped) == 0 {
		return nil, nil
	}
	if _, err := tx.Exec(ctx, `UPDATE items SET equipped=false,updated_at=clock_timestamp() WHERE campaign_id=$1 AND player_character_id=$2 AND id::text=ANY($3)`, campaignID, actorID, dropped); err != nil {
		return nil, err
	}
	if shield {
		// Only the authoritative shield field changes; other build state stays.
		if _, err := tx.Exec(ctx, `UPDATE player_characters SET stats=jsonb_set(stats,'{srd521,shield}','false'::jsonb),updated_at=clock_timestamp() WHERE id=$1 AND campaign_id=$2 AND stats ? 'srd521'`, actorID, campaignID); err != nil {
			return nil, err
		}
	}
	return dropped, nil
}
