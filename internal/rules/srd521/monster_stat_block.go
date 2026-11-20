package srd521

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// The types below extend the #143 monster contract for full stat-block batches
// (first used by #151). They describe source values only. Every typed field
// requires a dependency key in RequiredMechanics, and no key is implemented.

// MonsterUses retains limited-use notation: N/Day, N/Day in a lair, and the
// lowest d6 roll of a "Recharge X–6" entry. Unused parts are explicit zeros.
type MonsterUses struct {
	PerDay     int `json:"per_day"`
	LairPerDay int `json:"lair_per_day"`
	Recharge   int `json:"recharge"`
}

// MonsterSave retains saving throw effect notation without resolving it.
type MonsterSave struct {
	Ability string `json:"ability"`
	DC      int    `json:"dc"`
}

// MonsterTarget retains who makes a save. Shape is creature, line, cone, sphere
// or emanation. Size is line length, cone length or radius in feet; Width is a
// line's width; Range is the distance in feet to the creature or point of origin.
type MonsterTarget struct {
	Shape string `json:"shape"`
	Size  int    `json:"size"`
	Width int    `json:"width"`
	Range int    `json:"range"`
}

// MonsterLegendary retains the Legendary Action Uses entry. LairUses is zero
// when the source gives no lair value.
type MonsterLegendary struct {
	Uses     int    `json:"uses"`
	LairUses int    `json:"lair_uses"`
	Text     string `json:"text"`
}

// MonsterSpellcasting retains a Spellcasting action's ability, DC, spell attack
// bonus (zero when the source gives none) and spell lists as printed, including
// any parenthetical version or restriction. It never registers a spell.
type MonsterSpellcasting struct {
	Ability     string             `json:"ability"`
	SaveDC      int                `json:"save_dc"`
	AttackBonus int                `json:"attack_bonus"`
	AtWill      []string           `json:"at_will"`
	PerDay      []MonsterSpellUses `json:"per_day"`
}

// MonsterSpellUses is one "N/Day" or "N/Day Each" spell group.
type MonsterSpellUses struct {
	Uses   int      `json:"uses"`
	Spells []string `json:"spells"`
}

// UnmarshalJSON requires every limited-use number, including explicit zeros.
func (u *MonsterUses) UnmarshalJSON(data []byte) error {
	type plain MonsterUses
	return decodeMonsterNumbers(data, (*plain)(u), "per_day", "lair_per_day", "recharge")
}

// UnmarshalJSON requires an explicit DC.
func (s *MonsterSave) UnmarshalJSON(data []byte) error {
	type plain MonsterSave
	return decodeMonsterNumbers(data, (*plain)(s), "dc")
}

// UnmarshalJSON requires every dimension, including explicit zeros.
func (t *MonsterTarget) UnmarshalJSON(data []byte) error {
	type plain MonsterTarget
	return decodeMonsterNumbers(data, (*plain)(t), "size", "width", "range")
}

// UnmarshalJSON requires both use counts, including an explicit zero lair value.
func (l *MonsterLegendary) UnmarshalJSON(data []byte) error {
	type plain MonsterLegendary
	return decodeMonsterNumbers(data, (*plain)(l), "uses", "lair_uses")
}

// UnmarshalJSON requires the DC and an explicit spell attack bonus.
func (s *MonsterSpellcasting) UnmarshalJSON(data []byte) error {
	type plain MonsterSpellcasting
	return decodeMonsterNumbers(data, (*plain)(s), "save_dc", "attack_bonus")
}

// UnmarshalJSON requires an explicit use count.
func (s *MonsterSpellUses) UnmarshalJSON(data []byte) error {
	type plain MonsterSpellUses
	return decodeMonsterNumbers(data, (*plain)(s), "uses")
}

var (
	monsterDamageTypes = strings.Fields("acid bludgeoning cold fire force lightning necrotic piercing poison psychic radiant slashing thunder")
	monsterConditions  = strings.Fields("blinded charmed deafened exhaustion frightened grappled incapacitated invisible paralyzed petrified poisoned prone restrained stunned unconscious")
	monsterSkills      = strings.Fields("acrobatics animal-handling arcana athletics deception history insight intimidation investigation medicine nature perception performance persuasion religion sleight-of-hand stealth survival")
	monsterAbilities   = strings.Fields("str dex con int wis cha")
	monsterAbilityName = map[string]string{"str": "Strength", "dex": "Dexterity", "con": "Constitution", "int": "Intelligence", "wis": "Wisdom", "cha": "Charisma"}
	monsterShapes      = strings.Fields("creature line cone sphere emanation")
	monsterSpellNote   = regexp.MustCompile(` \(.*\)$`)
	monsterSlugGap     = regexp.MustCompile(`[^a-z0-9]+`)
)

// MonsterSpellKey returns the dependency key of a printed spell entry, ignoring
// its parenthetical version or restriction.
func MonsterSpellKey(printed string) string {
	name := strings.ToLower(monsterSpellNote.ReplaceAllString(strings.TrimSpace(printed), ""))
	return "spell/" + strings.Trim(monsterSlugGap.ReplaceAllString(name, "-"), "-")
}

// allActions lists every action section so IDs, names and dependency keys are
// checked across the whole block.
func (m Monster) allActions() []MonsterAction {
	return slices.Concat(m.Actions, m.BonusActions, m.Reactions, m.LegendaryActions)
}

// validateStatBlock checks the optional stat-block fields. Each present field
// must name its dependency key, which keeps the whole block unavailable. It also
// rejects numbers that contradict the source's own rules: XP must follow CR
// (PDF p. 256), a saving throw is the ability modifier with or without the
// Proficiency Bonus, and a typed save must be printed in its own entry.
func (m Monster) validateStatBlock() error {
	if printed, ok := monsterPrintedXP[m.ID]; ok {
		if printed != [2]string{m.CR, fmt.Sprint(m.XP)} {
			return fmt.Errorf("XP %d or CR %s differs from the printed stat block", m.XP, m.CR)
		}
	} else if !encounterValidXP(m.CR, int64(m.XP)) {
		return fmt.Errorf("XP %d does not follow CR %s", m.XP, m.CR)
	}
	for _, key := range monsterAbilities {
		if a := m.Abilities[key]; a.Save != a.Modifier && a.Save != a.Modifier+m.ProficiencyBonus {
			return fmt.Errorf("%s save is not the modifier with optional proficiency", key)
		}
	}
	need := func(present bool, key string) error {
		if present && !slices.Contains(m.RequiredMechanics, key) {
			return fmt.Errorf("stat-block field requires mechanic %q", key)
		}
		return nil
	}
	listed := func(values, allowed []string, key string) error {
		if values == nil {
			return nil
		}
		if len(values) == 0 || !uniqueMonsterStrings(values) {
			return fmt.Errorf("empty or duplicate list for %q", key)
		}
		for _, v := range values {
			if allowed != nil && !slices.Contains(allowed, v) {
				return fmt.Errorf("invalid value %q for %q", v, key)
			}
		}
		return need(key != "", key)
	}
	for _, c := range []struct {
		values, allowed []string
		key             string
	}{
		{m.Tags, nil, ""},
		{m.Vulnerabilities, monsterDamageTypes, "defense/damage-vulnerability"},
		{m.Resistances, monsterDamageTypes, "defense/damage-resistance"},
		{m.Immunities, monsterDamageTypes, "defense/damage-immunity"},
		{m.ConditionImmunities, monsterConditions, "defense/condition-immunity"},
	} {
		if err := listed(c.values, c.allowed, c.key); err != nil {
			return err
		}
	}
	if m.Skills != nil && len(m.Skills) == 0 {
		return fmt.Errorf("empty skills")
	}
	for skill, bonus := range m.Skills {
		if !slices.Contains(monsterSkills, skill) || bonus < -5 || bonus > 30 {
			return fmt.Errorf("invalid skill %q", skill)
		}
	}
	for mode := range m.Speeds {
		if err := need(mode != "walk", "movement/"+mode); err != nil {
			return err
		}
	}
	for sense := range m.Senses.Special {
		if err := need(true, "sense/"+sense); err != nil {
			return err
		}
	}
	if m.Hover && m.Speeds["fly"] <= 0 || m.Telepathy < 0 || m.LairXP != 0 && m.LairXP <= m.XP {
		return fmt.Errorf("invalid hover, telepathy or lair XP")
	}
	if (m.Legendary != nil) != (len(m.LegendaryActions) > 0) || m.LegendaryActions != nil && len(m.LegendaryActions) == 0 ||
		m.BonusActions != nil && len(m.BonusActions) == 0 || m.Reactions != nil && len(m.Reactions) == 0 {
		return fmt.Errorf("empty action section or legendary uses without actions")
	}
	if l := m.Legendary; l != nil && (l.Uses < 1 || l.LairUses != 0 && l.LairUses <= l.Uses || strings.TrimSpace(l.Text) == "") {
		return fmt.Errorf("invalid legendary action uses")
	}
	for _, c := range []struct {
		present bool
		key     string
	}{
		{m.Skills != nil, "monster/skills"},
		{m.Hover, "movement/hover"},
		{m.Telepathy > 0, "language/telepathy"},
		{m.LairXP > 0, "lair/xp"},
		{m.BonusActions != nil, "action/bonus-action"},
		{m.Reactions != nil, "action/reaction"},
		{m.Legendary != nil, "legendary/actions"},
		{m.Legendary != nil && m.Legendary.LairUses > 0, "lair/uses"},
	} {
		if err := need(c.present, c.key); err != nil {
			return err
		}
	}
	mechanic := func(x MonsterMechanic) error {
		if u := x.Uses; u != nil {
			if u.PerDay < 0 || u.Recharge != 0 && (u.Recharge < 2 || u.Recharge > 6) || (u.PerDay == 0) == (u.Recharge == 0) ||
				u.LairPerDay != 0 && u.LairPerDay <= u.PerDay {
				return fmt.Errorf("invalid uses for %q", x.Name)
			}
			if err := need(u.PerDay > 0, "limit/per-day"); err != nil {
				return err
			}
			if err := need(u.Recharge > 0, "limit/recharge"); err != nil {
				return err
			}
			if err := need(u.LairPerDay > 0, "lair/uses"); err != nil {
				return err
			}
		}
		if s := x.Save; s != nil {
			if !slices.Contains(monsterAbilities, s.Ability) || s.DC < 1 || s.DC > 30 {
				return fmt.Errorf("invalid save for %q", x.Name)
			}
			if !strings.Contains(x.Text, fmt.Sprintf("%s Saving Throw: DC %d", monsterAbilityName[s.Ability], s.DC)) {
				return fmt.Errorf("typed save for %q is not printed in its entry", x.Name)
			}
			if err := need(true, "save/"+s.Ability); err != nil {
				return err
			}
		}
		if t := x.Target; t != nil {
			if !slices.Contains(monsterShapes, t.Shape) || t.Size < 0 || t.Width < 0 || t.Range < 0 ||
				(t.Shape == "creature") != (t.Size == 0) || (t.Shape == "line") != (t.Width > 0) {
				return fmt.Errorf("invalid target for %q", x.Name)
			}
			key := "area/" + t.Shape
			if t.Shape == "creature" {
				key = "target/creature"
			}
			if err := need(true, key); err != nil {
				return err
			}
		}
		if x.Dice != nil && !validMonsterDice(*x.Dice) {
			return fmt.Errorf("invalid dice for %q", x.Name)
		}
		if x.Requires != nil && (len(x.Requires) == 0 || !uniqueMonsterStrings(x.Requires)) {
			return fmt.Errorf("empty or duplicate requirements for %q", x.Name)
		}
		for _, key := range x.Requires {
			if !validMonsterKey(key) {
				return fmt.Errorf("invalid requirement %q", key)
			}
			if err := need(true, key); err != nil {
				return err
			}
		}
		return nil
	}
	for _, trait := range m.Traits {
		if err := mechanic(trait); err != nil {
			return err
		}
	}
	for _, a := range m.allActions() {
		if err := mechanic(MonsterMechanic{Name: a.Name, Text: a.Text, Uses: a.Uses, Save: a.Save, Target: a.Target}); err != nil {
			return err
		}
		for _, effect := range a.Effects {
			if err := mechanic(effect); err != nil {
				return err
			}
		}
		if (a.Key == "action/spellcasting") != (a.Spellcasting != nil) {
			return fmt.Errorf("spellcasting action requires typed spell lists")
		}
		if s := a.Spellcasting; s != nil {
			// A zero DC means the source prints none (the Giant Owl); a printed DC must match.
			printed := strings.Contains(a.Text, "spell save DC")
			if !slices.Contains(monsterAbilities, s.Ability) || s.SaveDC < 0 || s.SaveDC > 30 || s.AttackBonus < 0 || s.AtWill == nil || s.PerDay == nil ||
				printed != (s.SaveDC > 0) || printed && !strings.Contains(a.Text, fmt.Sprintf("spell save DC %d", s.SaveDC)) {
				return fmt.Errorf("invalid spellcasting")
			}
			spells := slices.Clone(s.AtWill)
			for _, group := range s.PerDay {
				if group.Uses < 1 || len(group.Spells) == 0 {
					return fmt.Errorf("invalid daily spell group")
				}
				spells = append(spells, group.Spells...)
			}
			if len(spells) == 0 || !uniqueMonsterStrings(spells) {
				return fmt.Errorf("empty or duplicate spell list")
			}
			for _, spell := range spells {
				key := MonsterSpellKey(spell)
				if !validMonsterKey(key) {
					return fmt.Errorf("invalid spell %q", spell)
				}
				if err := need(true, key); err != nil {
					return err
				}
			}
		}
	}
	return m.validateStatBlockNotes()
}
