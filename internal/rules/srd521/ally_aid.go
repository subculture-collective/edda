package srd521

import "fmt"

// AllyAid is a trusted ally's turn spent helping the player character at
// 0 HP outside encounters (SRD 5.2.1 PDF pp. 17-18 "Stabilizing a
// Character", p. 99 "Potion of Healing" and p. 204 "Using a Potion"; adapted, see
// Attribution). The Unconscious player takes no action: the ally does.
const AllyAid = "ally-aid"

// AllyAidMedicine is the Help action to stabilize: DC 10 Wisdom (Medicine).
const AllyAidMedicine = "medicine"

// StabilizeDC is the Wisdom (Medicine) DC to stabilize a creature at 0 HP.
const StabilizeDC = 10

// Zero-HP receipt events for help from another creature.
const (
	ZeroHPAllyStabilize = "ally-stabilize"
	ZeroHPAllyHealing   = "ally-healing"
)

// ResolveAllyStabilize is a helper's Help action to stabilize a dying
// character: a d20 plus the helper's Wisdom (Medicine) modifier against DC 10.
// Success makes the character Stable and resets the counts. Failure changes
// nothing. Ability checks have no automatic success on a 20 or failure on a 1.
func ResolveAllyStabilize(saves DeathSaves, maxHP, modifier int, dice Dice) (*ZeroHPReceipt, error) {
	if err := saves.Validate(); err != nil {
		return nil, err
	}
	if saves.Dead || saves.Stable || maxHP < 1 || dice == nil {
		return nil, fmt.Errorf("the character is not dying")
	}
	roll := dice.Intn(20)
	if roll < 0 || roll >= 20 {
		return nil, fmt.Errorf("invalid d20")
	}
	m := modifier
	r := &ZeroHPReceipt{RulesVersion: "srd-5.2.1", Event: ZeroHPAllyStabilize, Die: roll + 1, Modifier: &m, Total: roll + 1 + modifier, DC: StabilizeDC, MaxHP: maxHP, Before: saves, After: saves, SourceIDs: []string{prefix + "damage/stabilizing-a-character", prefix + "actions/help"}}
	r.Reason = fmt.Sprintf("Help to stabilize: Wisdom (Medicine) d20 %d %+d = %d against DC %d", r.Die, modifier, r.Total, StabilizeDC)
	if r.Total >= StabilizeDC {
		r.Outcome, r.After = "stable", DeathSaves{Stable: true}
		r.SourceIDs = append(r.SourceIDs, prefix+"rules-glossary/stable")
		r.Reason += ": success. Stable at 0 HP, still Unconscious; no more Death Saving Throws."
		return r, nil
	}
	r.Outcome = "failure"
	r.Reason += ": failure. Still dying (" + saves.counts() + ")."
	return r, nil
}

// ResolveAdministeredPotion is a basic Potion of Healing that another creature
// administers. Unlike drinking one, the target may be at 0 HP: any healing
// ends Unconscious and resets the Death Saving Throw counts, which the caller
// applies.
func ResolveAdministeredPotion(hp, maxHP int, dice Dice) (*HealingReceipt, error) {
	return resolveHealingPotion(hp, maxHP, 0, dice)
}
