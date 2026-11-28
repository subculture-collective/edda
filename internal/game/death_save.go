package game

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"git.subcult.tv/subculture-collective/edda/internal/combat"
	"git.subcult.tv/subculture-collective/edda/internal/domain"
	"git.subcult.tv/subculture-collective/edda/internal/rules/srd521"
	statedb "git.subcult.tv/subculture-collective/edda/internal/state/sqlc"
)

// ZeroHPState returns the 0-HP state a client shows, or nil while the
// character has Hit Points and is alive. State the strict decoder refuses is
// left out; every action preview gives the reason.
func ZeroHPState(stats []byte, hp int) *srd521.DeathSaves {
	saves, err := srd521.DecodeDeathSaves(stats)
	if err != nil || (hp > 0 && !saves.Dead) {
		return nil
	}
	return &saves
}

// deadReason is shown for every choice once the character is dead.
const deadReason = "Your character is dead. No supported rule revives them."

// missingUnconsciousReason names a character an older build left at 0 HP
// without the Unconscious condition. The death-save action repairs it.
const missingUnconsciousReason = "Your character is at 0 HP without the Unconscious condition, which earlier versions of Edda did not apply. Use the death-save action (web action menu, TUI /death-save) to apply it, then make Death Saving Throws."

// ZeroHPTurnReason says why the player character cannot take a free-text
// turn: it is dead, at 0 HP, or Unconscious. It reads the same death saves
// and conditions that gate rules actions; an active encounter owns the
// player's conditions. It returns "" when the character can act.
func ZeroHPTurnReason(stats []byte, hp int, encounter *combat.Encounter) string {
	saves, err := srd521.DecodeDeathSaves(stats)
	if err != nil {
		return err.Error()
	}
	if saves.Dead {
		return deadReason
	}
	conditions, err := srd521.DecodeConditions(stats)
	active := encounter != nil && encounter.Status == "active"
	if active {
		for _, a := range encounter.Actors {
			if a.Player {
				conditions, err = a.Conditions, a.Validate()
			}
		}
	}
	if err != nil {
		return err.Error()
	}
	switch {
	case hp > 0 && !conditions.Unconscious:
		return ""
	case hp > 0:
		return "Your character is Unconscious and cannot act."
	case active:
		return "Your character is Unconscious at 0 HP and cannot act. End your turn in the encounter; the Death Saving Throw is rolled at the start of your next turn."
	case !conditions.Unconscious:
		return missingUnconsciousReason
	case saves.Stable:
		return "Your character is Stable at 0 HP but still Unconscious, and cannot act. Use the death-save action (web action menu, TUI /death-save) to wait 1d4 hours and regain 1 HP."
	}
	return "Your character is Unconscious at 0 HP and cannot act. Use the death-save action (web action menu, TUI /death-save) to make a Death Saving Throw."
}

// prepareDeathSave previews the explicit Death Saving Throw outside
// encounters, or a Stable character's 1d4-hour recovery. Edda keeps no turn
// order outside encounters, so the player confirms each save, and each one
// advances campaign time by one 6-second round. Inside an active encounter
// the save is rolled at the start of the player's turn instead.
// Preparation never rolls.
func prepareDeathSave(ctx context.Context, database statedb.DBTX, campaign domain.Campaign, pc domain.PlayerCharacter, intent ActionIntent) (*StoredAction, error) {
	p, err := prepareCharacterAction(ctx, campaign, pc, ActionIntent{ActorID: intent.ActorID, Kind: srd521.AbilityCheck, Ability: "strength"}, true)
	if err != nil {
		return nil, err
	}
	result := &StoredAction{Preview: p}
	p.Intent = intent
	p.Modifier = nil
	p.Options, p.SourceIDs = []string{}, []string{}
	if !p.Available {
		return result, nil
	}
	p.Available = false
	if intent.EntryID != "" || intent.Ability != "" || intent.Skill != "" || intent.SlotLevel != 0 {
		p.Reason = "A Death Saving Throw takes no choices."
		return result, nil
	}
	if pc.HP != 0 {
		p.Reason = "Death Saving Throws are made only at 0 HP."
		return result, nil
	}
	encounter, err := ReadEncounter(ctx, database, campaign.ID)
	if err != nil {
		return nil, err
	}
	if encounter != nil && encounter.Status == "active" {
		p.Reason = "In an encounter, the Death Saving Throw is rolled at the start of your turn. End your turn."
		return result, nil
	}
	saves, err := srd521.DecodeDeathSaves(pc.Stats)
	if err != nil {
		p.Reason = err.Error()
		return result, nil
	}
	conditions, err := srd521.DecodeConditions(pc.Stats)
	if err != nil {
		p.Reason = err.Error()
		return result, nil
	}
	if result.Clock, err = ReadSpellClock(ctx, database, campaign.ID); err != nil {
		return nil, err
	}
	encoded, err := json.Marshal([]any{p.ExpectedRevision, saves, result.Clock})
	if err != nil {
		return nil, err
	}
	p.ExpectedRevision = fmt.Sprintf("death-save-v1:%x", sha256.Sum256(encoded))
	if !conditions.Unconscious {
		// An older build left the character at 0 HP without Unconscious.
		// Confirming applies it once; the next confirmation rolls the save.
		p.ExpectedRevision = fmt.Sprintf("unconscious-repair-v1:%x", sha256.Sum256(encoded))
		p.Options = []string{"Your character is at 0 HP without the Unconscious condition, which earlier versions of Edda did not apply. Confirming applies Unconscious now: you fall Prone, drop held weapons and a Shield, and Concentration ends. Nothing is rolled and no time passes; Death Saving Throws follow with this action."}
		p.SourceIDs = []string{"srd-5.2.1:content/damage/falling-unconscious", "srd-5.2.1:content/rules-glossary/unconscious"}
	} else if saves.Stable {
		p.Options = []string{"Stable at 0 HP and Unconscious. Wait to recover: after 1d4 hours you regain 1 HP; campaign time advances that long. Unconscious ends and Prone remains."}
		p.SourceIDs = []string{"srd-5.2.1:content/damage/stabilizing-a-character"}
	} else {
		p.Options = []string{fmt.Sprintf("Death Saving Throw: d20, no modifier. 10 or higher succeeds; a 1 counts as two failures; a 20 regains 1 HP. Three successes make you Stable; three failures kill you. Now %d successes, %d failures. Campaign time advances one 6-second round.", saves.Successes, saves.Failures)}
		p.SourceIDs = []string{"srd-5.2.1:content/damage/death-saving-throws", "srd-5.2.1:content/rules-glossary/death-saving-throw"}
	}
	p.Available = true
	p.Reason = ""
	return result, nil
}
