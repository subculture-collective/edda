package game

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"git.subcult.tv/subculture-collective/edda/internal/dbutil"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// AidPlayerProfile is the trusted srd521_ally ruling that an NPC is willing
// and able to help this player character at 0 HP: the shared ally facts plus
// its Wisdom (Medicine) modifier. Model tools cannot author it.
const AidPlayerProfile = "aid-player-v1"

// prepareAllyAid previews a trusted ally's turn spent helping the Unconscious
// player outside encounters: the Help action to stabilize (DC 10 Wisdom
// (Medicine)), or administering one of the player's Potions of Healing. The
// player takes no action; the confirmation stands for one 6-second round in
// which the ally acts and, if the player is still dying, the player's next
// turn starts with a Death Saving Throw. Preparation never rolls.
func prepareAllyAid(ctx context.Context, database statedb.DBTX, campaign domain.Campaign, pc domain.PlayerCharacter, intent ActionIntent, lock bool) (*StoredAction, error) {
	choices := intent.Ability != "" || intent.Skill != "" || intent.SlotLevel != 0
	result, saves, _, ok, err := prepareZeroHPAction(ctx, database, campaign, pc, intent, choices, zeroHPReasons{
		Choices:     "Ally aid takes only the kind of help: medicine or one of your Potions of Healing.",
		NotZero:     "An ally's aid at 0 HP is for a character at 0 HP.",
		Encounter:   "The two-actor encounter has no ally who can act. End your turn.",
		Unconscious: "This character has 0 HP without the Unconscious condition; that state requires adjudication.",
	})
	if err != nil || !ok {
		return result, err
	}
	p := result.Preview
	var item []any
	potion := intent.EntryID != srd521.AllyAidMedicine
	if !potion {
		if saves.Stable {
			p.Reason = "You are already Stable. A Potion of Healing can still bring you back."
			return result, nil
		}
	} else {
		itemID, err := uuid.Parse(intent.EntryID)
		if err != nil {
			p.Reason = "Choose medicine or one of your Potions of Healing."
			return result, nil
		}
		locking := ""
		if lock {
			locking = " FOR UPDATE"
		}
		var properties []byte
		var updated time.Time
		err = database.QueryRow(ctx, `SELECT properties,updated_at,quantity FROM items WHERE id=$1 AND campaign_id=$2 AND player_character_id=$3 AND item_type='consumable' AND quantity>0`+locking, dbutil.ToPgtype(itemID), dbutil.ToPgtype(campaign.ID), dbutil.ToPgtype(pc.ID)).Scan(&properties, &updated, &result.Quantity)
		if err != nil && err != pgx.ErrNoRows {
			return nil, err
		}
		var props struct {
			ID string `json:"srd521_id"`
		}
		if err == pgx.ErrNoRows || json.Unmarshal(properties, &props) != nil || props.ID != "potion-of-healing" {
			p.Reason = "Only a Potion of Healing you carry can be administered to you."
			return result, nil
		}
		item = []any{itemID, properties, updated, result.Quantity}
	}
	if result.HealingTargets, err = readAllies(ctx, database, campaign, pc, AidPlayerProfile, lock); err != nil {
		return nil, err
	}
	if len(result.HealingTargets) == 0 {
		p.Reason = "No ally here is known to be willing and able to help you. An ally needs a trusted aid ruling within 5 feet of you."
		return result, nil
	}
	if result.Clock, err = ReadSpellClock(ctx, database, campaign.ID); err != nil {
		return nil, err
	}
	encoded, err := json.Marshal([]any{p.ExpectedRevision, saves, result.Clock, item, result.HealingTargets})
	if err != nil {
		return nil, err
	}
	p.ExpectedRevision = fmt.Sprintf("ally-aid-v1:%x", sha256.Sum256(encoded))
	round := "Campaign time advances one 6-second round."
	if saves.Stable {
		round = "Campaign time advances one 6-second round. You are Stable, so no Death Saving Throw follows."
	} else {
		round += fmt.Sprintf(" If you are still dying afterwards, your next turn starts with a Death Saving Throw (now %d successes, %d failures).", saves.Successes, saves.Failures)
	}
	for _, ally := range result.HealingTargets {
		label := fmt.Sprintf("%s (Wisdom (Medicine) %+d; within 5 feet)", ally.Name, ally.Medicine)
		if potion {
			label = ally.Name + " (administers the potion; within 5 feet)"
		}
		p.Targets = append(p.Targets, ActionTarget{ID: ally.ID, Label: label, DistanceFeet: ally.Space.DistanceFeet, ReachFeet: 5})
	}
	if potion {
		p.Options = []string{fmt.Sprintf("The chosen ally takes one of your Potions of Healing (quantity %d) and, as a Bonus Action, gives it to you: you regain 2d4 + 2 HP, capped at %d. Unconscious ends, Prone remains and the Death Saving Throw counts reset. %s", result.Quantity, pc.MaxHP, round)}
		p.Costs = []ActionCost{{ResourceID: "item:" + intent.EntryID, Amount: 1}}
		p.SourceIDs = []string{"srd-5.2.1:content/items/potion-of-healing", "srd-5.2.1:content/damage/healing", "srd-5.2.1:content/damage/falling-unconscious"}
	} else {
		p.Options = []string{"The chosen ally takes the Help action to stabilize you: a Wisdom (Medicine) check against DC 10. On a success you are Stable at 0 HP, still Unconscious. On a failure nothing changes. " + round}
		p.SourceIDs = []string{"srd-5.2.1:content/damage/stabilizing-a-character", "srd-5.2.1:content/actions/help"}
	}
	p.SourceIDs = append(p.SourceIDs, "srd-5.2.1:content/damage/death-saving-throws")
	p.CostStatus = "fixture-only"
	p.Available = true
	p.Reason = ""
	return result, nil
}
