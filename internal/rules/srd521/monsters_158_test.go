package srd521

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func loadMonsters158(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_158.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #158 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Draft Horse's block starts on p. 347 and the Earth Elemental's
// starts on p. 282; each runs onto the next page.
var monsters158Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Dire Wolf", [2]int{347, 347}, "Large Beast, Unaligned|AC 14 Initiative +2 (12)|HP 22 (3d10 + 6)|Speed 50 ft.|Str 17 +3 +3 Dex 15 +2 +2 Con 15 +2 +2|Int 3 −4 −4 Wis 12 +1 +1 Cha 7 −2 −2|Skills Perception +5, Stealth +4|Senses Darkvision 60 ft.; Passive Perception 15|Languages None|CR 1 (XP 200; PB +2)"},
	{"Djinni", [2]int{280, 280}, "Large Elemental (Genie), Neutral|AC 17 Initiative +2 (12)|HP 218 (19d10 + 114)|Speed 30 ft., Fly 90 ft. (hover)|Str 21 +5 +5 Dex 15 +2 +6 Con 22 +6 +6|Int 15 +2 +2 Wis 16 +3 +7 Cha 20 +5 +5|Immunities Lightning, Thunder|Senses Darkvision 120 ft.; Passive Perception 13|Languages Primordial (Auran)|CR 11 (XP 7,200; PB +4)"},
	{"Doppelganger", [2]int{280, 280}, "Medium Monstrosity, Neutral|AC 14 Initiative +4 (14)|HP 52 (8d8 + 16)|Speed 30 ft.|Str 11 +0 +0 Dex 18 +4 +4 Con 14 +2 +2|Int 11 +0 +0 Wis 12 +1 +1 Cha 14 +2 +2|Skills Deception +6, Insight +3|Immunities Charmed|Senses Darkvision 60 ft.; Passive Perception 11|Languages Common plus three other languages|CR 3 (XP 700; PB +2)"},
	{"Draft Horse", [2]int{347, 348}, "Large Beast, Unaligned|AC 10 Initiative +0 (10)|HP 15 (2d10 + 4)|Speed 40 ft.|Str 18 +4 +4 Dex 10 +0 +0 Con 15 +2 +2|Int 2 −4 −4 Wis 11 +0 +0 Cha 7 −2 −2|Senses Passive Perception 10|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Dragon Turtle", [2]int{281, 281}, "Gargantuan Dragon, Neutral|AC 20 Initiative +6 (16)|HP 356 (23d20 + 115)|Speed 20 ft., Swim 50 ft.|Str 25 +7 +7 Dex 10 +0 +0 Con 20 +5 +11|Int 10 +0 +0 Wis 12 +1 +7 Cha 12 +1 +1|Resistances Fire|Senses Darkvision 120 ft.; Passive Perception 11|Languages Draconic, Primordial (Aquan)|CR 17 (XP 18,000; PB +6)"},
	{"Dretch", [2]int{281, 281}, "Small Fiend (Demon), Chaotic Evil|AC 11 Initiative +0 (10)|HP 18 (4d6 + 4)|Speed 20 ft.|Str 12 +1 +1 Dex 11 +0 +0 Con 12 +1 +1|Int 5 −3 −3 Wis 8 −1 −1 Cha 3 −4 −4|Resistances Cold, Fire, Lightning|Immunities Poison; Poisoned|Senses Darkvision 60 ft.; Passive Perception 9|Languages Abyssal; telepathy 60 ft. (works only with creatures that understand Abyssal)|CR 1/4 (XP 50; PB +2)"},
	{"Drider", [2]int{281, 281}, "Large Monstrosity, Chaotic Evil|AC 19 Initiative +4 (14)|HP 123 (13d10 + 52)|Speed 30 ft., Climb 30 ft.|Str 16 +3 +3 Dex 19 +4 +4 Con 18 +4 +4|Int 13 +1 +1 Wis 16 +3 +3 Cha 12 +1 +1|Skills Perception +6, Stealth +10|Senses Darkvision 120 ft.; Passive Perception 16|Languages Elvish, Undercommon|CR 6 (XP 2,300; PB +3)"},
	{"Druid", [2]int{282, 282}, "Medium or Small Humanoid (Druid), Neutral|AC 13 Initiative +1 (11)|HP 44 (8d8 + 8)|Speed 30 ft.|Str 10 +0 +0 Dex 12 +1 +1 Con 13 +1 +1|Int 12 +1 +1 Wis 16 +3 +3 Cha 11 +0 +0|Skills Medicine +5, Nature +3, Perception +5|Gear Studded Leather Armor|Senses Passive Perception 15|Languages Common, Druidic, Sylvan|CR 2 (XP 450; PB +2)"},
	{"Dryad", [2]int{282, 282}, "Medium Fey, Neutral|AC 16 Initiative +1 (11)|HP 22 (5d8)|Speed 30 ft.|Str 10 +0 +0 Dex 12 +1 +1 Con 11 +0 +0|Int 14 +2 +2 Wis 15 +2 +2 Cha 18 +4 +4|Skills Perception +4, Stealth +5|Senses Darkvision 60 ft.; Passive Perception 14|Languages Elvish, Sylvan|CR 1 (XP 200; PB +2)"},
	{"Dust Mephit", [2]int{307, 307}, "Small Elemental, Neutral Evil|AC 12 Initiative +2 (12)|HP 17 (5d6)|Speed 30 ft., Fly 30 ft.|Str 5 −3 −3 Dex 14 +2 +2 Con 10 +0 +0|Int 9 −1 −1 Wis 11 +0 +0 Cha 10 +0 +0|Skills Perception +2, Stealth +4|Vulnerabilities Fire|Immunities Poison; Exhaustion, Poisoned|Senses Darkvision 60 ft.; Passive Perception 12|Languages Primordial (Auran, Terran)|CR 1/2 (XP 100; PB +2)"},
	{"Eagle", [2]int{348, 348}, "Small Beast, Unaligned|AC 12 Initiative +2 (12)|HP 4 (1d6 + 1)|Speed 10 ft., Fly 60 ft.|Str 6 −2 −2 Dex 15 +2 +2 Con 12 +1 +1|Int 2 −4 −4 Wis 14 +2 +2 Cha 7 −2 −2|Skills Perception +6|Senses Passive Perception 16|Languages None|CR 0 (XP 10; PB +2)"},
	{"Earth Elemental", [2]int{282, 283}, "Large Elemental, Neutral|AC 17 Initiative −1 (9)|HP 147 (14d10 + 70)|Speed 30 ft., Burrow 30 ft.|Str 20 +5 +5 Dex 8 −1 −1 Con 20 +5 +5|Int 5 −3 −3 Wis 10 +0 +0 Cha 5 −3 −3|Vulnerabilities Thunder|Immunities Poison; Exhaustion, Paralyzed, Petrified, Poisoned, Unconscious|Senses Darkvision 60 ft., Tremorsense 60 ft.; Passive Perception 10|Languages Primordial (Terran)|CR 5 (XP 1,800; PB +3)"},
}

// Typed values of every trait, action and bonus action, reviewed against the
// same pages: limited uses, attack bonus/reach/range/damage, save ability and
// DC, target shape size/width/range, dice and dependency keys. A ranged attack
// printed with one range repeats it. A Death Burst, a one-off spell cast and a
// whirlwind keep their printed text; the keys name what a runtime would have to
// apply, and none is implemented.
var monsters158Entries = map[string][]string{
	"Dire Wolf": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage, requirement/ally-within-5-feet, requirement/ally-not-incapacitated",
		"action|bite|Bite; melee +5 reach 5 range 0/0 8 (1d10 + 3) piercing; effect/piercing-damage 8 (1d10 + 3), condition/prone, requirement/target-size",
	},
	"Djinni": {
		"trait|trait/elemental-restoration|Elemental Restoration; effect/revive, effect/revival-delay",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"trait|trait/wishes|Wishes; spell/wish, effect/wish-chance, effect/wish-for-non-genie, effect/no-wish-stress, effect/wish-limit",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|storm-blade|Storm Blade; melee +9 reach 5 range 0/0 12 (2d6 + 5) slashing; effect/slashing-damage 12 (2d6 + 5), effect/lightning-damage 7 (2d6)",
		"action|storm-bolt|Storm Bolt; ranged +9 reach 0 range 120/120 13 (3d8) thunder; effect/thunder-damage 13 (3d8), condition/prone, requirement/target-size",
		"action|create-whirlwind|Create Whirlwind; str DC 17; effect/whirlwind, effect/concentration, effect/whirlwind-move, save/str, save/once-per-turn, condition/restrained, effect/thunder-damage 21 (6d6), save/repeat",
		"action|spellcasting|Spellcasting; cha DC 17 +0; at will: Detect Evil and Good, Detect Magic; 2/day: Create Food and Water (can create wine instead of water), Tongues, Wind Walk; 1/day: Creation, Gaseous Form, Invisibility, Major Image, Plane Shift; spellcasting/no-material-components, spell/detect-evil-and-good [spellcasting/at-will], spell/detect-magic [spellcasting/at-will], spell/create-food-and-water [limit/per-day], spell/tongues [limit/per-day], spell/wind-walk [limit/per-day], spell/creation [limit/per-day], spell/gaseous-form [limit/per-day], spell/invisibility [limit/per-day], spell/major-image [limit/per-day], spell/plane-shift [limit/per-day]",
	},
	"Doppelganger": {
		"action|multiattack|Multiattack; action/melee-attack, multiattack/additional-action",
		"action|slam|Slam; melee +6 reach 5 range 0/0 11 (2d6 + 4) bludgeoning; effect/bludgeoning-damage 11 (2d6 + 4), attack/advantage [requirement/first-combat-round]",
		"action|read-thoughts|Read Thoughts; action/spellcasting, spell/detect-thoughts, spellcasting/no-components, spellcasting/fixed-save-dc",
		"action|unsettling-visage|Unsettling Visage (Recharge 6); wis DC 12; emanation 15/0/0; save/wis, area/emanation [requirement/target-sees-source], limit/recharge, condition/frightened, save/repeat, effect/timed-automatic-success",
		"bonus|shape-shift|Shape-Shift; effect/shape-shift, effect/shapechange-form-limits, effect/equipment-not-transformed",
	},
	"Draft Horse": {
		"action|hooves|Hooves; melee +6 reach 5 range 0/0 6 (1d4 + 4) bludgeoning; effect/bludgeoning-damage 6 (1d4 + 4)",
	},
	"Dragon Turtle": {
		"trait|trait/amphibious|Amphibious; effect/breathe-air-and-water",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack",
		"action|bite|Bite; melee +13 reach 15 range 0/0 23 (3d10 + 7) piercing; effect/piercing-damage 23 (3d10 + 7), effect/fire-damage 7 (2d6), effect/underwater-no-resistance",
		"action|tail|Tail; melee +13 reach 15 range 0/0 18 (2d10 + 7) bludgeoning; effect/bludgeoning-damage 18 (2d10 + 7), condition/prone, requirement/target-size",
		"action|steam-breath|Steam Breath (Recharge 5–6); con DC 19; cone 60/0/0; save/con, area/cone, limit/recharge, effect/fire-damage 56 (16d6), save/half-damage, effect/underwater-no-resistance",
	},
	"Dretch": {
		"action|rend|Rend; melee +3 reach 5 range 0/0 4 (1d6 + 1) slashing; effect/slashing-damage 4 (1d6 + 1)",
		"action|fetid-cloud|Fetid Cloud (1/Day); con DC 11; emanation 10/0/0; save/con, area/emanation, limit/per-day, condition/poisoned, effect/action-or-bonus-action, effect/no-reactions",
	},
	"Drider": {
		"trait|trait/spider-climb|Spider Climb; effect/climb-without-check",
		"trait|trait/sunlight-sensitivity|Sunlight Sensitivity; requirement/in-sunlight, check/disadvantage, attack/disadvantage",
		"trait|trait/web-walker|Web Walker; effect/ignore-web-restrictions, effect/know-web-contacts",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|foreleg|Foreleg; melee +7 reach 10 range 0/0 13 (2d8 + 4) piercing; effect/piercing-damage 13 (2d8 + 4)",
		"action|poison-burst|Poison Burst; ranged +6 reach 0 range 120/120 13 (3d6 + 3) poison; effect/poison-damage 13 (3d6 + 3)",
		"bonus|magic-of-the-spider-queen|Magic of the Spider Queen (Recharge 5–6); limit/recharge, action/spellcasting, spell/darkness, spell/faerie-fire, spell/web, spellcasting/no-material-components, spellcasting/fixed-save-dc",
	},
	"Druid": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|vine-staff|Vine Staff; melee +5 reach 5 range 0/0 7 (1d8 + 3) bludgeoning; effect/bludgeoning-damage 7 (1d8 + 3), effect/poison-damage 2 (1d4)",
		"action|verdant-wisp|Verdant Wisp; ranged +5 reach 0 range 90/90 10 (3d6) radiant; effect/radiant-damage 10 (3d6)",
		"action|spellcasting|Spellcasting; wis DC 13 +0; at will: Druidcraft, Speak with Animals; 2/day: Entangle, Thunderwave; 1/day: Animal Messenger, Longstrider, Moonbeam; spell/druidcraft [spellcasting/at-will], spell/speak-with-animals [spellcasting/at-will], spell/entangle [limit/per-day], spell/thunderwave [limit/per-day], spell/animal-messenger [limit/per-day], spell/longstrider [limit/per-day], spell/moonbeam [limit/per-day]",
	},
	"Dryad": {
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"trait|trait/speak-with-beasts-and-plants|Speak with Beasts and Plants; effect/communicate-with-beasts-and-plants",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/additional-action, action/spellcasting",
		"action|vine-lash|Vine Lash; melee +6 reach 10 range 0/0 8 (1d8 + 4) slashing; effect/slashing-damage 8 (1d8 + 4)",
		"action|thorn-burst|Thorn Burst; ranged +6 reach 0 range 60/60 7 (1d6 + 4) piercing; effect/piercing-damage 7 (1d6 + 4)",
		"action|spellcasting|Spellcasting; cha DC 14 +0; at will: Animal Friendship, Charm Monster (lasts 24 hours; ends early if the dryad casts the spell again), Druidcraft; 1/day: Entangle, Pass without Trace; spellcasting/no-material-components, spell/animal-friendship [spellcasting/at-will], spell/charm-monster [spellcasting/at-will], spell/druidcraft [spellcasting/at-will], spell/entangle [limit/per-day], spell/pass-without-trace [limit/per-day]",
		"bonus|tree-stride|Tree Stride; effect/teleport, requirement/near-large-tree, requirement/unoccupied-space",
	},
	"Dust Mephit": {
		"trait|trait/death-burst|Death Burst; dex DC 10; emanation 5/0/0; 5 (2d4); trigger/death, save/dex, area/emanation, effect/bludgeoning-damage, save/half-damage",
		"action|claw|Claw; melee +4 reach 5 range 0/0 4 (1d4 + 2) slashing; effect/slashing-damage 4 (1d4 + 2)",
		"action|blinding-breath|Blinding Breath (Recharge 6); dex DC 10; cone 15/0/0; save/dex, area/cone, limit/recharge, condition/blinded",
		"action|sleep|Sleep (1/Day); limit/per-day, action/spellcasting, spell/sleep, spellcasting/no-components, spellcasting/fixed-save-dc",
	},
	"Eagle": {
		"action|talons|Talons; melee +4 reach 5 range 0/0 4 (1d4 + 2) slashing; effect/slashing-damage 4 (1d4 + 2)",
	},
	"Earth Elemental": {
		"trait|trait/earth-glide|Earth Glide; movement/burrow-through-earth-and-stone, effect/no-disturbance",
		"trait|trait/siege-monster|Siege Monster; effect/double-damage-to-objects",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|slam|Slam; melee +8 reach 10 range 0/0 14 (2d8 + 5) bludgeoning; effect/bludgeoning-damage 14 (2d8 + 5)",
		"action|rock-launch|Rock Launch; ranged +8 reach 0 range 60/60 8 (1d6 + 5) bludgeoning; effect/bludgeoning-damage 8 (1d6 + 5), condition/prone, requirement/target-size",
	},
}

// monsterDice158 prints a fixed hit, which the source prints as a bare number
// (none in this batch), and otherwise the source's "N (XdY + Z)" notation.
func monsterDice158(d MonsterHP) string {
	switch {
	case d.Dice == 0:
		return fmt.Sprint(d.Average)
	case d.Modifier < 0:
		return fmt.Sprintf("%d (%dd%d − %d)", d.Average, d.Dice, d.Sides, -d.Modifier)
	}
	return monsterDice151(d)
}

// monsterHeader158 prints typed fields in the source's stat-block notation.
func monsterHeader158(m Monster) string {
	kind := monsterTitle151(m.Type)
	if len(m.Tags) > 0 {
		kind += " (" + strings.Join(m.Tags, ", ") + ")"
	}
	sizes := make([]string, len(m.Sizes))
	for i, size := range m.Sizes {
		sizes[i] = monsterTitle151(size)
	}
	lines := []string{
		strings.Join(sizes, " or ") + " " + kind + ", " + m.Alignment,
		fmt.Sprintf("AC %d Initiative %s (%d)", m.AC, monsterSigned151(m.Initiative.Modifier), m.Initiative.Score),
		"HP " + monsterDice158(m.HP),
	}
	speed := fmt.Sprintf("Speed %d ft.", m.Speeds["walk"])
	for _, mode := range monsterSorted151(m.Speeds) {
		if mode != "walk" {
			speed += fmt.Sprintf(", %s %d ft.", monsterTitle151(mode), m.Speeds[mode])
		}
	}
	if m.Hover {
		speed += " (hover)"
	}
	lines = append(lines, speed)
	for _, row := range [][]string{{"str", "dex", "con"}, {"int", "wis", "cha"}} {
		cells := []string{}
		for _, key := range row {
			a := m.Abilities[key]
			cells = append(cells, fmt.Sprintf("%s %d %s %s", monsterTitle151(key), a.Score, monsterSigned151(a.Modifier), monsterSigned151(a.Save)))
		}
		lines = append(lines, strings.Join(cells, " "))
	}
	if m.Skills != nil {
		skills := []string{}
		for _, skill := range monsterSorted151(m.Skills) {
			skills = append(skills, monsterTitle151(skill)+" "+monsterSigned151(m.Skills[skill]))
		}
		lines = append(lines, "Skills "+strings.Join(skills, ", "))
	}
	if m.Vulnerabilities != nil {
		lines = append(lines, "Vulnerabilities "+monsterTitles151(m.Vulnerabilities))
	}
	if m.Resistances != nil {
		lines = append(lines, "Resistances "+monsterTitles151(m.Resistances))
	}
	if immunities := slices.DeleteFunc([]string{monsterTitles151(m.Immunities), monsterTitles151(m.ConditionImmunities)}, func(s string) bool { return s == "" }); len(immunities) > 0 {
		lines = append(lines, "Immunities "+strings.Join(immunities, "; "))
	}
	if len(m.Gear) > 0 {
		lines = append(lines, "Gear "+strings.Join(m.Gear, ", "))
	}
	senses := []string{}
	for _, sense := range monsterSorted151(m.Senses.Special) {
		printed := fmt.Sprintf("%s %d ft.", monsterTitle151(sense), m.Senses.Special[sense])
		if note, ok := m.Senses.Notes[sense]; ok {
			printed += " (" + note + ")"
		}
		senses = append(senses, printed)
	}
	passive := fmt.Sprintf("Passive Perception %d", m.Senses.PassivePerception)
	if len(senses) > 0 {
		passive = strings.Join(senses, ", ") + "; " + passive
	}
	languages := "Languages None"
	if len(m.Languages) > 0 {
		languages = "Languages " + strings.Join(m.Languages, ", ")
	}
	if m.Telepathy > 0 {
		languages += fmt.Sprintf("; telepathy %d ft.", m.Telepathy)
		if m.TelepathyNote != "" {
			languages += " (" + m.TelepathyNote + ")"
		}
	}
	lines = append(lines, "Senses "+passive, languages, fmt.Sprintf("CR %s (XP %s; PB +%d)", m.CR, monsterThousands151(m.XP), m.ProficiencyBonus))
	return strings.Join(lines, "|")
}

func monsterEntry158(section, id string, x MonsterMechanic, attack *MonsterAttack, sc *MonsterSpellcasting, effects []MonsterMechanic) string {
	head := x.Name
	if u := x.Uses; u != nil {
		if u.Recharge == 6 {
			head += " (Recharge 6)"
		} else if u.Recharge > 0 {
			head += fmt.Sprintf(" (Recharge %d–6)", u.Recharge)
		} else {
			head += fmt.Sprintf(" (%d/Day)", u.PerDay)
		}
	}
	parts := []string{head}
	if attack != nil {
		parts = append(parts, fmt.Sprintf("%s +%d reach %d range %d/%d %s %s", attack.Kind, attack.Bonus, attack.Reach, attack.Range, attack.LongRange, monsterDice158(attack.Damage), attack.DamageType))
	}
	if x.Save != nil {
		parts = append(parts, fmt.Sprintf("%s DC %d", x.Save.Ability, x.Save.DC))
	}
	if t := x.Target; t != nil {
		parts = append(parts, fmt.Sprintf("%s %d/%d/%d", t.Shape, t.Size, t.Width, t.Range))
	}
	if x.Dice != nil {
		parts = append(parts, monsterDice158(*x.Dice))
	}
	for _, d := range x.MoreDice {
		parts = append(parts, d.Key+" "+monsterDice158(d.Dice))
	}
	if sc != nil {
		parts = append(parts, fmt.Sprintf("%s DC %d +%d", sc.Ability, sc.SaveDC, sc.AttackBonus), "at will: "+strings.Join(sc.AtWill, ", "))
		for _, group := range sc.PerDay {
			parts = append(parts, fmt.Sprintf("%d/day: %s", group.Uses, strings.Join(group.Spells, ", ")))
		}
	}
	keys := slices.Clone(x.Requires)
	for _, e := range effects {
		key := e.Key
		if e.Dice != nil {
			key += " " + monsterDice158(*e.Dice)
		}
		if len(e.Requires) > 0 {
			key += " [" + strings.Join(e.Requires, ", ") + "]"
		}
		keys = append(keys, key)
	}
	return section + "|" + id + "|" + strings.Join(append(parts, strings.Join(keys, ", ")), "; ")
}

func monsterEntries158(m Monster) []string {
	out := []string{}
	for _, trait := range m.Traits {
		out = append(out, monsterEntry158("trait", trait.Key, trait, nil, nil, nil))
	}
	for _, section := range []struct {
		label   string
		actions []MonsterAction
	}{{"action", m.Actions}, {"bonus", m.BonusActions}, {"reaction", m.Reactions}} {
		for _, a := range section.actions {
			out = append(out, monsterEntry158(section.label, a.ID, MonsterMechanic{Name: a.Name, Uses: a.Uses, Save: a.Save, Target: a.Target}, a.Attack, a.Spellcasting, a.Effects))
		}
	}
	return out
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters158SourceValues(t *testing.T) {
	c := loadMonsters158(t)
	if len(c.Monsters) != 12 || len(monsters158Headers) != 12 || len(monsters158Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #158") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters158Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader158(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries158(m); !slices.Equal(got, monsters158Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters158Entries[want.name])
			}
			if m.Reactions != nil || m.Hover != (m.Name == "Djinni") || m.LairXP != 0 || m.LegendaryActions != nil || m.Legendary != nil {
				t.Fatal("source lists no reactions, lair XP or legendary actions, and only the Djinni hovers")
			}
			if (m.Telepathy != 0) != (m.Name == "Dretch") || (m.TelepathyNote != "") != (m.Name == "Dretch") {
				t.Fatal("only the Dretch lists telepathy, with its printed condition")
			}
		})
		entries += len(monsters158Entries[want.name])
	}
	// Bold-italic entry headings counted on PDF pp. 280–283, 307 and 347–348.
	if entries != 51 {
		t.Fatalf("reviewed %d entries, want 51", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters158MatchesDecisionInventory(t *testing.T) {
	b, err := os.ReadFile("../../../docs/investigations/issue-82-world-inventory.json")
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Records []struct {
			ID, Name, Family, Status string
			SourcePage               int `json:"source_page"`
			CatalogueIssue           int `json:"catalogue_issue"`
		}
	}
	if err := json.Unmarshal(b, &inventory); err != nil {
		t.Fatal(err)
	}
	type identity struct {
		id, name string
		page     int
	}
	want, got, families := []identity{}, []identity{}, []string{}
	for _, r := range inventory.Records {
		if r.CatalogueIssue == 158 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters158(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Dire Wolf, Draft Horse and Eagle under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("animals monsters monsters animals monsters monsters monsters monsters monsters monsters animals monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #158 cites the page each block starts on, in this order.
	names := []string{"Dire Wolf", "Djinni", "Doppelganger", "Draft Horse", "Dragon Turtle", "Dretch", "Drider", "Druid", "Dryad", "Dust Mephit", "Eagle", "Earth Elemental"}
	pages := []int{347, 280, 280, 347, 281, 281, 281, 282, 282, 307, 348, 282}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic that holds for these twelve blocks (PDF pp. 254–256):
// XP and PB follow CR, Hit Dice follow size and Constitution, saves and skills
// add PB to the ability modifier, and attack rolls add PB to the listed ability.
// Initiative adds PB (Dragon Turtle) or nothing.
func TestMonsters158DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"0": 2, "1/4": 2, "1/2": 2, "1": 2, "2": 2, "3": 2, "5": 3, "6": 3, "11": 4, "17": 6}
	skillAbility := map[string]string{"deception": "cha", "insight": "wis", "medicine": "wis", "nature": "int", "perception": "wis", "stealth": "dex"}
	hitDie := map[string]int{"tiny": 4, "small": 6, "medium": 8, "large": 10, "huge": 12, "gargantuan": 20}
	// Attacks that do not use the better of Strength and Dexterity, and hits
	// whose printed damage adds no ability modifier.
	attackRule := map[string]struct {
		ability string
		plain   bool
	}{
		"Druid/vine-staff": {"wis", false}, "Druid/verdant-wisp": {"wis", true},
		"Dryad/vine-lash": {"cha", false}, "Dryad/thorn-burst": {"cha", false},
		"Drider/poison-burst": {"str", false}, "Djinni/storm-bolt": {"", true},
	}
	for _, m := range loadMonsters158(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			p := m.ProficiencyBonus
			if !encounterValidXP(m.CR, int64(m.XP)) || pb[m.CR] != p {
				t.Fatal("XP or PB does not follow CR")
			}
			if m.HP.Sides != hitDie[m.Sizes[0]] || m.HP.Modifier != m.HP.Dice*m.Abilities["con"].Modifier {
				t.Fatal("Hit Dice do not follow size and Constitution")
			}
			for key, a := range m.Abilities {
				if a.Save != a.Modifier && a.Save != a.Modifier+p {
					t.Fatalf("%s save is not the modifier with optional PB", key)
				}
			}
			if extra := m.Initiative.Modifier - m.Abilities["dex"].Modifier; extra != 0 && extra != p {
				t.Fatal("initiative is not Dexterity with optional PB")
			}
			passive := 10 + m.Abilities["wis"].Modifier
			for skill, bonus := range m.Skills {
				ability, ok := skillAbility[skill]
				if extra := bonus - m.Abilities[ability].Modifier; !ok || extra != p && extra != 2*p {
					t.Fatalf("%s is not an ability modifier plus PB or twice PB", skill)
				}
				if skill == "perception" {
					passive = 10 + bonus
				}
			}
			if m.Senses.PassivePerception != passive {
				t.Fatal("Passive Perception differs from 10 plus Perception")
			}
			for _, a := range m.allActions() {
				if a.Attack == nil {
					continue
				}
				mod := max(m.Abilities["str"].Modifier, m.Abilities["dex"].Modifier)
				rule, special := attackRule[m.Name+"/"+a.ID]
				if rule.ability != "" {
					mod = m.Abilities[rule.ability].Modifier
				}
				// A hit printed without a modifier adds no ability modifier.
				damage := mod
				if rule.plain && special {
					damage = 0
				}
				if a.Attack.Bonus != mod+p || a.Attack.Damage.Modifier != damage {
					t.Fatalf("%s attack does not use the ability modifier and PB", a.Name)
				}
			}
		})
	}
}

// Every listed key is used and required, every effect quotes its own entry, and
// every dice expression is typed. Removing any one key invalidates the record.
func TestMonsters158DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters158(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			if m.RequiredMechanics[0] != "monster/stat-block" {
				t.Fatal("stat-block dependency must come first")
			}
			quoted := func(parent string, x MonsterMechanic) {
				if !strings.Contains(parent, strings.TrimSuffix(x.Text, ".")) {
					t.Fatalf("%q is not quoted from its entry", x.Text)
				}
				if x.Dice != nil && !strings.Contains(x.Text, monsterDice158(*x.Dice)) {
					t.Fatalf("%q does not print its typed dice", x.Text)
				}
			}
			for _, trait := range m.Traits {
				if !strings.HasPrefix(trait.Text, trait.Name) || trait.Key != "trait/"+strings.ReplaceAll(strings.ToLower(trait.Name), " ", "-") {
					t.Fatalf("trait %q key or text differs", trait.Name)
				}
				quoted(trait.Text, trait)
			}
			for _, a := range m.allActions() {
				if !strings.HasPrefix(a.Text, a.Name) || a.ID != strings.ReplaceAll(strings.ToLower(a.Name), " ", "-") {
					t.Fatalf("action %q ID or text differs", a.Name)
				}
				for _, e := range a.Effects {
					quoted(a.Text, e)
				}
				if a.Attack == nil || m.ID == BanditID { // the #143 Bandit revision types damage on the attack only
					continue
				}
				first := a.Effects[0]
				if fixedMonsterDamage(a.Attack.Damage) {
					if first.Dice != nil || !slices.Equal(first.Requires, []string{"damage/fixed"}) || !strings.Contains(a.Text, fmt.Sprintf("Hit: %d %s damage.", a.Attack.Damage.Average, monsterTitle151(a.Attack.DamageType))) {
						t.Fatalf("%s: fixed hit must be printed without dice and require damage/fixed", a.Name)
					}
				} else if first.Dice == nil || *first.Dice != a.Attack.Damage {
					t.Fatalf("%s: first effect must carry the attack damage", a.Name)
				}
			}
			for i, key := range m.RequiredMechanics {
				without := m
				without.RequiredMechanics = slices.Delete(slices.Clone(m.RequiredMechanics), i, i+1)
				if without.Validate() == nil {
					t.Fatalf("dependency %q is unused or unenforced", key)
				}
			}
		})
	}
}

// Catalogue presence grants nothing: no key is implemented, the Bandit fixture's
// keys cover only the record #143 already pinned, and the register still reports
// both families unsupported.
func TestMonsters158RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters158(t).Monsters {
		if m.Available() || m.availableWith(bandit) != (m.ID == BanditID) {
			t.Fatalf("%s must stay unavailable", m.Name)
		}
		all := map[string]bool{}
		for _, key := range m.RequiredMechanics[1:] {
			all[key] = true
		}
		if m.availableWith(all) {
			t.Fatalf("%s: one missing mechanic must block the whole record", m.Name)
		}
	}
	r, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range r.Entries {
		if (e.Family == "monsters" || e.Family == "animals") && (e.Status != "unsupported" || len(e.Fixtures) != 0) {
			t.Fatal("catalogue data promoted the monster or animal family")
		}
	}
}

func TestMonsters158MalformedFixtures(t *testing.T) {
	// Indexes: 0 Dire Wolf, 1 Djinni, 2 Doppelganger, 3 Draft Horse, 4 Dragon Turtle,
	// 5 Dretch, 6 Drider, 7 Druid, 8 Dryad, 9 Dust Mephit, 10 Eagle, 11 Earth Elemental.
	tests := []struct {
		name   string
		mutate func(*MonsterCatalogue)
	}{
		{"duplicate ID", func(c *MonsterCatalogue) { c.Monsters[1].ID = c.Monsters[2].ID }},
		{"animal block in another family", func(c *MonsterCatalogue) { c.Monsters[0].ID = "srd-5.2.1:content/spells/dire-wolf" }},
		{"missing languages list", func(c *MonsterCatalogue) { c.Monsters[0].Languages = nil }},
		{"sense note without sense", func(c *MonsterCatalogue) {
			c.Monsters[0].Senses.Notes = map[string]string{"blindsight": "unimpeded by magical Darkness"}
		}},
		{"blank sense note", func(c *MonsterCatalogue) { c.Monsters[0].Senses.Notes = map[string]string{"darkvision": " "} }},
		{"unlisted sense note", func(c *MonsterCatalogue) {
			c.Monsters[0].Senses.Notes = map[string]string{"darkvision": "impeded by magical Darkness"}
		}},
		{"empty sense notes", func(c *MonsterCatalogue) { c.Monsters[0].Senses.Notes = map[string]string{} }},
		{"telepathy note without telepathy", func(c *MonsterCatalogue) {
			m := &c.Monsters[0]
			m.TelepathyNote = "works only with creatures that understand Abyssal"
			m.RequiredMechanics = append(m.RequiredMechanics, "language/telepathy-condition")
		}},
		{"telepathy note without mechanic", func(c *MonsterCatalogue) {
			m := &c.Monsters[5]
			m.RequiredMechanics = slices.DeleteFunc(m.RequiredMechanics, func(key string) bool { return key == "language/telepathy-condition" })
		}},
		{"padded telepathy note", func(c *MonsterCatalogue) {
			c.Monsters[5].TelepathyNote = " works only with creatures that understand Abyssal"
		}},
		{"blank telepathy note", func(c *MonsterCatalogue) { c.Monsters[5].TelepathyNote = " " }},
		{"melee range", func(c *MonsterCatalogue) { c.Monsters[0].Actions[0].Attack.Range = 5 }},
		{"further dice without first", func(c *MonsterCatalogue) {
			c.Monsters[0].Traits[0].MoreDice = []MonsterKeyedDice{{Key: "effect/future-damage", Dice: MonsterHP{Average: 3, Dice: 1, Sides: 4, Modifier: 1}}}
		}},
		{"XP off CR", func(c *MonsterCatalogue) { c.Monsters[1].XP = 7201 }},
		{"save off modifier", func(c *MonsterCatalogue) { c.Monsters[1].Abilities["wis"] = MonsterAbility{16, 3, 10} }},
		{"modifier", func(c *MonsterCatalogue) { c.Monsters[0].Abilities["str"] = MonsterAbility{17, 4, 4} }},
		{"HP average", func(c *MonsterCatalogue) { c.Monsters[10].HP.Average = 5 }},
		{"hover without flight", func(c *MonsterCatalogue) { c.Monsters[1].Speeds["fly"] = 0 }},
		{"empty bonus actions", func(c *MonsterCatalogue) { c.Monsters[0].BonusActions = []MonsterAction{} }},
		{"empty reactions", func(c *MonsterCatalogue) { c.Monsters[0].Reactions = []MonsterAction{} }},
		{"duplicate action across sections", func(c *MonsterCatalogue) { c.Monsters[2].BonusActions[0].ID = "slam" }},
		{"recharge", func(c *MonsterCatalogue) { c.Monsters[4].Actions[3].Uses.Recharge = 1 }},
		{"recharge above six", func(c *MonsterCatalogue) { c.Monsters[2].Actions[3].Uses.Recharge = 7 }},
		{"daily uses with recharge", func(c *MonsterCatalogue) { c.Monsters[4].Actions[3].Uses.PerDay = 1 }},
		{"save DC off entry text", func(c *MonsterCatalogue) { c.Monsters[4].Actions[3].Save.DC = 20 }},
		{"save ability off entry text", func(c *MonsterCatalogue) { c.Monsters[4].Actions[3].Save.Ability = "dex" }},
		{"cone width", func(c *MonsterCatalogue) { c.Monsters[4].Actions[3].Target.Width = 5 }},
		{"emanation width", func(c *MonsterCatalogue) { c.Monsters[5].Actions[1].Target.Width = 5 }},
		{"unknown target shape", func(c *MonsterCatalogue) {
			c.Monsters[1].Actions[3].Target = &MonsterTarget{Shape: "cylinder", Size: 20}
		}},
		{"effect dice", func(c *MonsterCatalogue) { c.Monsters[1].Actions[1].Effects[1].Dice.Average = 9 }},
		{"unlisted action effect", func(c *MonsterCatalogue) { c.Monsters[5].Actions[1].Effects[3].Key = "effect/future" }},
		{"untyped spellcasting", func(c *MonsterCatalogue) { c.Monsters[1].Actions[4].Spellcasting = nil }},
		{"misplaced spellcasting", func(c *MonsterCatalogue) {
			c.Monsters[1].Actions[3].Spellcasting = c.Monsters[1].Actions[4].Spellcasting
		}},
		{"spell ability", func(c *MonsterCatalogue) { c.Monsters[1].Actions[4].Spellcasting.Ability = "" }},
		{"spell DC", func(c *MonsterCatalogue) { c.Monsters[1].Actions[4].Spellcasting.SaveDC = 0 }},
		{"negative spell attack bonus", func(c *MonsterCatalogue) { c.Monsters[7].Actions[3].Spellcasting.AttackBonus = -1 }},
		{"unlisted spell", func(c *MonsterCatalogue) { c.Monsters[1].Actions[4].Spellcasting.AtWill[0] = "Meteor Swarm" }},
		{"duplicate spell", func(c *MonsterCatalogue) { c.Monsters[1].Actions[4].Spellcasting.AtWill[0] = "Tongues" }},
		{"daily spell uses", func(c *MonsterCatalogue) { c.Monsters[7].Actions[3].Spellcasting.PerDay[0].Uses = 0 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters158(t)
			tt.mutate(&c)
			if c.Validate() == nil {
				t.Fatal("accepted malformed catalogue")
			}
			data, err := json.Marshal(c)
			if err != nil {
				t.Fatal(err)
			}
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			if strings.HasPrefix(tt.name, "empty ") {
				return
			}
			if decoded, err := DecodeMonsters(bytes.NewReader(data)); err == nil || len(decoded.Monsters) != 0 {
				t.Fatal("decoder must return no catalogue on invalid data")
			}
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_158.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, edit := range [][2]string{
		{`"lair_per_day": 0,`, ``},
		{`"sides": 10,`, `"sides": null,`},
		{`"long_range": 0,`, `"long_range": null,`},
		{`"recharge": 5`, `"recharge": null`},
		{`"per_day": 0,`, ``},
		{`"dc": 17`, `"dc": 17, "dc": 16`},
		{`"width": 0,`, ``},
		{`"special": {}`, `"special": {"darkvision": null}`},
		{`"shape": "cone",`, `"shape": "cone", "height": 5,`},
		{`"hover": true,`, `"hover": true, "flies": true,`},
		{`"telepathy_note": "`, `"telepathy_notes": "`},
	} {
		if !bytes.Contains(raw, []byte(edit[0])) {
			t.Fatalf("fixture lost %s", edit[0])
		}
		if _, err := DecodeMonsters(bytes.NewReader(bytes.Replace(raw, []byte(edit[0]), []byte(edit[1]), 1))); err == nil {
			t.Fatalf("accepted %s -> %s", edit[0], edit[1])
		}
	}
}

// The batch composes with the #143 Bandit catalogue and stat-block batches 01
// to 07 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters158ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters158(t)
	if len(base.Monsters) != 1 || base.Monsters[0].ID != BanditID {
		t.Fatal("batch data must not enter the embedded runtime catalogue")
	}
	b02, err := LoadMonsterBatch02()
	if err != nil {
		t.Fatal(err)
	}
	var second []Monster
	for _, sb := range b02.StatBlocks {
		second = append(second, sb.Monster)
	}
	first, third, fourth := loadMonsters151(t), loadMonsters153(t), loadMonsters154(t)
	fifth, sixth, seventh := loadMonsters155(t), loadMonsters156(t), loadMonsters157(t)
	for _, c := range []MonsterCatalogue{base, first, third, fourth, fifth, sixth, seventh} {
		a, b := c.Source, batch.Source
		a.Changes, b.Changes = "", ""
		if a != b {
			t.Fatal("batches must share one source pin")
		}
	}
	earlier := slices.Concat(base.Monsters, first.Monsters, second, third.Monsters, slices.Concat(fourth.Monsters[:4], fourth.Monsters[5:]), fifth.Monsters, sixth.Monsters, seventh.Monsters)
	combined := MonsterCatalogue{SchemaVersion: 1, Source: batch.Source, Monsters: slices.Concat(earlier, batch.Monsters)}
	if err := combined.Validate(); err != nil || len(combined.Monsters) != len(earlier)+12 {
		t.Fatal("batch must compose with the earlier catalogues:", err)
	}
	for _, m := range batch.Monsters {
		dup := MonsterCatalogue{SchemaVersion: 1, Source: batch.Source, Monsters: slices.Concat(earlier, batch.Monsters, []Monster{m})}
		if dup.Validate() == nil {
			t.Fatalf("composition accepted a second %s record", m.Name)
		}
	}
	_, version, err := BanditRecord()
	if err != nil || version != "bandit-v1:1169c77601cf83a9b96985fc761a6b8be9c901dc1a86e6a9758dcc399707c3d4" {
		t.Fatal("contract extension changed the pinned Bandit revision:", version, err)
	}
	for _, id := range []string{"bite", "claw", "foreleg", "hooves", "rend", "rock-launch", "slam", "storm-blade", "storm-bolt", "tail", "talons", "vine-lash", "vine-staff"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
