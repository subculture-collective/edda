package srd521

import "fmt"

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

	// Advantage lists the durable facts that granted Advantage; D20s then holds
	// both rolled dice and Die the higher one. Both are empty for a single roll.
	Advantage []string `json:"advantage,omitempty"`
	D20s      []int    `json:"d20s,omitempty"`
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
	if sheet == nil || sheet.RulesVersion != "srd-5.2.1" || armorClass < 1 || dice == nil {
		return nil, fmt.Errorf("supported sheet, target and dice required")
	}
	return resolveWeaponAttack(sheet.RulesVersion, "Longsword (one hand)", LongswordModifier(sheet), sheet.Abilities["strength"].Modifier.Value, 8, "slashing", armorClass, dice)
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
	return resolveAttackRoll(version, weapon, modifier, damageModifier, damageDie, damageType, armorClass, nil, dice)
}

// resolveAttackRoll rolls two d20s and keeps the higher when advantage names
// at least one reason. Advantage never stacks: two reasons still roll twice.
func resolveAttackRoll(version, weapon string, modifier Value, damageModifier, damageDie int, damageType string, armorClass int, advantage []string, dice Dice) (*AttackReceipt, error) {
	roll := dice.Intn(20)
	if roll < 0 || roll >= 20 {
		return nil, fmt.Errorf("invalid d20")
	}
	r := &AttackReceipt{RulesVersion: version, Kind: MartialAttack, Weapon: weapon, Die: roll + 1, Modifier: modifier, Outcome: "miss", DamageType: damageType, DamageDice: []int{}, SourceIDs: []string{prefix + "d20-tests/attack-rolls", prefix + "weapons/longsword", prefix + "damage/critical-hits"}}
	if len(advantage) > 0 {
		second := dice.Intn(20)
		if second < 0 || second >= 20 {
			return nil, fmt.Errorf("invalid d20")
		}
		r.Advantage, r.D20s, r.Die = append([]string(nil), advantage...), []int{roll + 1, second + 1}, max(roll, second)+1
	}
	if damageDie == 0 {
		r.SourceIDs = []string{prefix + "d20-tests/attack-rolls", prefix + "rules-glossary/unarmed-strike"}
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
