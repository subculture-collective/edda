package srd521

import (
	"fmt"
	"strings"
)

const MartialAttack = "martial-attack"

// AttackReceipt describes supported ordinary weapon and fixture attacks.
// It never contains a target's private stats or raw stored properties.
type AttackReceipt struct {
	RulesVersion   string   `json:"rules_version"`
	Kind           string   `json:"kind"`
	Weapon         string   `json:"weapon"`
	TargetID       string   `json:"target_id"`
	Target         string   `json:"target"`
	Die            int      `json:"die"`
	Modifier       Value    `json:"modifier"`
	Total          int      `json:"total"`
	Outcome        string   `json:"outcome"`
	DamageDice     []int    `json:"damage_dice"`
	DamageModifier int      `json:"damage_modifier"`
	Damage         int      `json:"damage"`
	DamageType     string   `json:"damage_type"`
	SourceIDs      []string `json:"source_ids"`

	// Advantage and Disadvantage list the durable facts that granted each.
	// When only one applies, D20s holds both rolled dice and Die the kept one
	// (higher for Advantage, lower for Disadvantage). When both apply they
	// cancel: one d20 is rolled and D20s is empty.
	Advantage    []string `json:"advantage,omitempty"`
	Disadvantage []string `json:"disadvantage,omitempty"`
	D20s         []int    `json:"d20s,omitempty"`
	// Condition records a size-gated condition outcome on a hit: "prone-applied",
	// "already-prone" or "size-exempt". It is empty when no condition applies.
	Condition string `json:"condition,omitempty"`
}

// LongswordModifier uses Fighter martial proficiency and Strength.
func LongswordModifier(sheet *CharacterSheet) Value {
	return derived(term("Strength", sheet.Abilities["strength"].Modifier.Value, "character/ability-modifiers"), term("Fighter martial weapon proficiency", sheet.Proficiency.Value, "fighter/core-traits"))
}

// ResolveLongsword handles an ordinary attack with no mastery, conditions,
// resistance or other situational effects. Callers must validate that scope.
func ResolveLongsword(sheet *CharacterSheet, armorClass int, dice Dice) (*AttackReceipt, error) {
	return ResolveLongswordRoll(sheet, armorClass, nil, nil, dice)
}

// ResolveLongswordRoll is ResolveLongsword with the durable facts that give
// the roll Advantage or Disadvantage, such as the attacker's Prone condition.
func ResolveLongswordRoll(sheet *CharacterSheet, armorClass int, advantage, disadvantage []string, dice Dice) (*AttackReceipt, error) {
	if sheet == nil || sheet.RulesVersion != "srd-5.2.1" || armorClass < 1 || dice == nil {
		return nil, fmt.Errorf("supported sheet, target and dice required")
	}
	return resolveAttackRoll(sheet.RulesVersion, "Longsword (one hand)", LongswordModifier(sheet), sheet.Abilities["strength"].Modifier.Value, 8, "slashing", armorClass, advantage, disadvantage, dice)
}

// ResolveFixtureUnarmed resolves the declared Strength 10, proficiency +2 NPC
// fixture. Its ordinary Unarmed Strike deals 1 bludgeoning damage on a hit.
func ResolveFixtureUnarmed(armorClass int, dice Dice) (*AttackReceipt, error) {
	if armorClass < 1 || dice == nil {
		return nil, fmt.Errorf("target and dice required")
	}
	return resolveWeaponAttack("srd-5.2.1", "Unarmed Strike", derived(term("Fixture attack proficiency", 2, "d20-tests/attack-rolls")), 1, 0, "bludgeoning", armorClass, dice)
}

func resolveWeaponAttack(version, weapon string, modifier Value, damageModifier, damageDie int, damageType string, armorClass int, dice Dice) (*AttackReceipt, error) {
	return resolveAttackRoll(version, weapon, modifier, damageModifier, damageDie, damageType, armorClass, nil, nil, dice)
}

// resolveAttackRoll rolls two d20s when only Advantage or only Disadvantage
// names a reason, keeping the higher or lower die. Neither stacks: several
// reasons still roll twice. When both name a reason they cancel and one d20
// is rolled, however many reasons each side has.
func resolveAttackRoll(version, weapon string, modifier Value, damageModifier, damageDie int, damageType string, armorClass int, advantage, disadvantage []string, dice Dice) (*AttackReceipt, error) {
	roll := dice.Intn(20)
	if roll < 0 || roll >= 20 {
		return nil, fmt.Errorf("invalid d20")
	}
	r := &AttackReceipt{RulesVersion: version, Kind: MartialAttack, Weapon: weapon, Die: roll + 1, Modifier: modifier, Outcome: "miss", DamageType: damageType, DamageDice: []int{}, SourceIDs: []string{prefix + "d20-tests/attack-rolls", prefix + "weapons/longsword", prefix + "damage/critical-hits"}}
	if len(advantage) > 0 {
		r.Advantage = append([]string(nil), advantage...)
	}
	if len(disadvantage) > 0 {
		r.Disadvantage = append([]string(nil), disadvantage...)
	}
	if (len(advantage) > 0) != (len(disadvantage) > 0) {
		second := dice.Intn(20)
		if second < 0 || second >= 20 {
			return nil, fmt.Errorf("invalid d20")
		}
		r.D20s, r.Die = []int{roll + 1, second + 1}, max(roll, second)+1
		if len(disadvantage) > 0 {
			r.Die = min(roll, second) + 1
		}
	}
	if damageDie == 0 {
		r.SourceIDs = []string{prefix + "d20-tests/attack-rolls", prefix + "rules-glossary/unarmed-strike"}
	}
	if len(advantage) > 0 || len(disadvantage) > 0 {
		r.SourceIDs = append(r.SourceIDs, prefix+"d20-tests/advantage-disadvantage")
	}
	r.Total = r.Die + r.Modifier.Value
	if r.Die == 1 || (r.Die != 20 && r.Total < armorClass) {
		return r, nil
	}
	r.Outcome = "hit"
	count := 1
	if r.Die == 20 {
		count = 2
		r.Outcome = "critical-hit"
	}
	r.DamageModifier = damageModifier
	r.Damage = r.DamageModifier
	for i := 0; i < count && damageDie > 0; i++ {
		value := dice.Intn(damageDie)
		if value < 0 || value >= damageDie {
			return nil, fmt.Errorf("invalid damage die")
		}
		r.DamageDice = append(r.DamageDice, value+1)
		r.Damage += value + 1
	}
	if r.Damage < 0 {
		r.Damage = 0
	}
	return r, nil
}

// RollNote explains how the d20 was rolled, for narration. It is empty for an
// ordinary single roll.
func (r *AttackReceipt) RollNote() string {
	switch {
	case len(r.Advantage) > 0 && len(r.Disadvantage) > 0:
		return fmt.Sprintf("Advantage (%s) and Disadvantage (%s) cancel: one d20.", strings.Join(r.Advantage, "; "), strings.Join(r.Disadvantage, "; "))
	case len(r.D20s) == 2 && len(r.Disadvantage) > 0:
		return fmt.Sprintf("Disadvantage (%s): rolled %d and %d, kept %d.", strings.Join(r.Disadvantage, "; "), r.D20s[0], r.D20s[1], r.Die)
	case len(r.D20s) == 2:
		return fmt.Sprintf("Advantage (%s): rolled %d and %d.", strings.Join(r.Advantage, "; "), r.D20s[0], r.D20s[1])
	}
	return ""
}
