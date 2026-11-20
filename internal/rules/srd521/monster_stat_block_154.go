package srd521

import (
	"fmt"
	"slices"
	"strings"
)

// The fields below extend the stat-block contract for notation first met in
// #154: a second dice expression in one trait, a parenthetical after a special
// sense, and an attack whose hit is a fixed number. Like every stat-block field,
// each requires a dependency key and no key is implemented.

// MonsterKeyedDice is one further dice expression of an entry and the mechanic
// it belongs to, such as the second damage type of a trait.
type MonsterKeyedDice struct {
	Key  string    `json:"key"`
	Dice MonsterHP `json:"dice"`
}

// MonsterSenseNoteKey returns the dependency key of a special sense's printed
// parenthetical, such as "unimpeded by magical Darkness".
func MonsterSenseNoteKey(sense, note string) string {
	return "sense/" + sense + "-" + strings.Trim(monsterSlugGap.ReplaceAllString(strings.ToLower(note), "-"), "-")
}

// fixedMonsterDamage reports a hit the source prints as a number without dice,
// for example "Hit: 1 Piercing damage". Its dice fields are explicit zeros.
func fixedMonsterDamage(d MonsterHP) bool {
	return d.Average > 0 && d.Dice == 0 && d.Sides == 0 && d.Modifier == 0
}

// validAttackDamage accepts rolled damage, or fixed damage when the record
// names the mechanic that would have to apply it.
func (m Monster) validAttackDamage(d MonsterHP) bool {
	return validMonsterDice(d) || fixedMonsterDamage(d) && slices.Contains(m.RequiredMechanics, "damage/fixed")
}

// validateStatBlockNotes checks sense notes and further dice. A note must
// belong to a listed special sense, and further dice must follow a typed first
// expression, be printed by their own entry and name a required mechanic.
func (m Monster) validateStatBlockNotes() error {
	if m.Senses.Notes != nil && len(m.Senses.Notes) == 0 {
		return fmt.Errorf("empty sense notes")
	}
	for sense, note := range m.Senses.Notes {
		key := MonsterSenseNoteKey(sense, note)
		if m.Senses.Special[sense] <= 0 || note != strings.TrimSpace(note) || !validMonsterKey(key) {
			return fmt.Errorf("invalid note for sense %q", sense)
		}
		if !slices.Contains(m.RequiredMechanics, key) {
			return fmt.Errorf("stat-block field requires mechanic %q", key)
		}
	}
	more := func(x MonsterMechanic) error {
		if x.MoreDice == nil {
			return nil
		}
		if len(x.MoreDice) == 0 || x.Dice == nil {
			return fmt.Errorf("further dice for %q need a first dice expression", x.Name)
		}
		seen := map[string]bool{}
		for _, d := range x.MoreDice {
			if seen[d.Key] {
				return fmt.Errorf("duplicate further dice key %q for %q", d.Key, x.Name)
			}
			seen[d.Key] = true
			if !validMonsterDice(d.Dice) || !validMonsterKey(d.Key) || !slices.Contains(m.RequiredMechanics, d.Key) {
				return fmt.Errorf("invalid further dice for %q", x.Name)
			}
			printed := fmt.Sprintf("%d (%dd%d)", d.Dice.Average, d.Dice.Dice, d.Dice.Sides)
			if d.Dice.Modifier > 0 {
				printed = fmt.Sprintf("%d (%dd%d + %d)", d.Dice.Average, d.Dice.Dice, d.Dice.Sides, d.Dice.Modifier)
			} else if d.Dice.Modifier < 0 {
				printed = fmt.Sprintf("%d (%dd%d − %d)", d.Dice.Average, d.Dice.Dice, d.Dice.Sides, -d.Dice.Modifier)
			}
			if !strings.Contains(x.Text, printed) {
				return fmt.Errorf("further dice for %q are not printed in its entry", x.Name)
			}
		}
		return nil
	}
	for _, trait := range m.Traits {
		if err := more(trait); err != nil {
			return err
		}
	}
	for _, a := range m.allActions() {
		for _, effect := range a.Effects {
			if err := more(effect); err != nil {
				return err
			}
		}
	}
	return m.validateTelepathyNote()
}
