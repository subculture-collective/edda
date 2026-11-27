package srd521

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func loadMonsters170(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_170.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #170 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Purple Worm’s and the Rakshasa’s blocks run onto the next
// page.
var monsters170Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Polar Bear", [2]int{359, 359}, "Large Beast, Unaligned|AC 12 Initiative +2 (12)|HP 42 (5d10 + 15)|Speed 40 ft., Swim 40 ft.|Str 20 +5 +5 Dex 14 +2 +2 Con 16 +3 +3|Int 2 −4 −4 Wis 13 +1 +1 Cha 7 −2 −2|Skills Perception +5, Stealth +4|Resistances Cold|Senses Darkvision 60 ft.; Passive Perception 15|Languages None|CR 2 (XP 450; PB +2)"},
	{"Pony", [2]int{359, 359}, "Medium Beast, Unaligned|AC 10 Initiative +0 (10)|HP 11 (2d8 + 2)|Speed 40 ft.|Str 15 +2 +4 Dex 10 +0 +0 Con 13 +1 +1|Int 2 −4 −4 Wis 11 +0 +0 Cha 7 −2 −2|Senses Passive Perception 10|Languages None|CR 1/8 (XP 25; PB +2)"},
	{"Priest", [2]int{316, 316}, "Medium or Small Humanoid (Cleric), Neutral|AC 13 Initiative +0 (10)|HP 38 (7d8 + 7)|Speed 30 ft.|Str 16 +3 +3 Dex 10 +0 +0 Con 12 +1 +1|Int 13 +1 +1 Wis 16 +3 +3 Cha 13 +1 +1|Skills Medicine +7, Perception +5, Religion +5|Gear Chain Shirt, Holy Symbol, Mace|Senses Passive Perception 15|Languages Common plus one other language|CR 2 (XP 450; PB +2)"},
	{"Priest Acolyte", [2]int{316, 316}, "Medium or Small Humanoid (Cleric), Neutral|AC 13 Initiative +0 (10)|HP 11 (2d8 + 2)|Speed 30 ft.|Str 14 +2 +2 Dex 10 +0 +0 Con 12 +1 +1|Int 10 +0 +0 Wis 14 +2 +2 Cha 11 +0 +0|Skills Medicine +4, Religion +2|Gear Chain Shirt, Holy Symbol, Mace|Senses Passive Perception 12|Languages Common|CR 1/4 (XP 50; PB +2)"},
	{"Pseudodragon", [2]int{316, 316}, "Tiny Dragon, Neutral Good|AC 14 Initiative +2 (12)|HP 10 (3d4 + 3)|Speed 15 ft., Fly 60 ft.|Str 6 −2 −2 Dex 15 +2 +2 Con 13 +1 +1|Int 10 +0 +0 Wis 12 +1 +1 Cha 10 +0 +0|Skills Perception +5, Stealth +4|Senses Blindsight 10 ft., Darkvision 60 ft.; Passive Perception 15|Languages Understands Common and Draconic but can’t speak|CR 1/4 (XP 50; PB +2)"},
	{"Pteranodon", [2]int{359, 359}, "Medium Beast (Dinosaur), Unaligned|AC 13 Initiative +2 (12)|HP 13 (3d8)|Speed 10 ft., Fly 60 ft.|Str 12 +1 +1 Dex 15 +2 +2 Con 10 +0 +0|Int 2 −4 −4 Wis 9 −1 −1 Cha 5 −3 −3|Skills Perception +1|Senses Passive Perception 11|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Purple Worm", [2]int{316, 317}, "Gargantuan Monstrosity, Unaligned|AC 18 Initiative +3 (13)|HP 247 (15d20 + 90)|Speed 50 ft., Burrow 50 ft.|Str 28 +9 +9 Dex 7 −2 −2 Con 22 +6 +11|Int 1 −5 −5 Wis 8 −1 +4 Cha 4 −3 −3|Senses Blindsight 30 ft., Tremorsense 60 ft.; Passive Perception 9|Languages None|CR 15 (XP 13,000; PB +5)"},
	{"Quasit", [2]int{317, 317}, "Tiny Fiend (Demon), Chaotic Evil|AC 13 Initiative +3 (13)|HP 25 (10d4)|Speed 40 ft.|Str 5 −3 −3 Dex 17 +3 +3 Con 10 +0 +0|Int 7 −2 −2 Wis 10 +0 +0 Cha 10 +0 +0|Skills Stealth +5|Resistances Cold, Fire, Lightning|Immunities Poison; Poisoned|Senses Darkvision 120 ft.; Passive Perception 10|Languages Abyssal, Common|CR 1 (XP 200; PB +2)"},
	{"Rakshasa", [2]int{317, 318}, "Medium Fiend, Lawful Evil|AC 17 Initiative +8 (18)|HP 221 (26d8 + 104)|Speed 40 ft.|Str 14 +2 +2 Dex 17 +3 +3 Con 18 +4 +4|Int 13 +1 +1 Wis 16 +3 +3 Cha 20 +5 +5|Skills Deception +10, Insight +8, Perception +8|Vulnerabilities Piercing damage from weapons wielded by creatures under the effect of a Bless spell|Immunities Charmed, Frightened|Senses Truesight 60 ft.; Passive Perception 18|Languages Common, Infernal|CR 13 (XP 10,000; PB +5)"},
	{"Rat", [2]int{359, 359}, "Tiny Beast, Unaligned|AC 10 Initiative +0 (10)|HP 1 (1d4 − 1)|Speed 20 ft., Climb 20 ft.|Str 2 −4 −4 Dex 11 +0 +0 Con 9 −1 −1|Int 2 −4 −4 Wis 10 +0 +0 Cha 4 −3 −3|Skills Perception +2|Senses Darkvision 30 ft.; Passive Perception 12|Languages None|CR 0 (XP 10; PB +2)"},
	{"Raven", [2]int{359, 359}, "Tiny Beast, Unaligned|AC 12 Initiative +2 (12)|HP 2 (1d4)|Speed 10 ft., Fly 50 ft.|Str 2 −4 −4 Dex 14 +2 +2 Con 10 +0 +0|Int 5 −3 −3 Wis 13 +1 +1 Cha 6 −2 −2|Skills Perception +3|Senses Passive Perception 13|Languages None|CR 0 (XP 10; PB +2)"},
	{"Red Dragon Wyrmling", [2]int{318, 318}, "Medium Dragon (Chromatic), Chaotic Evil|AC 17 Initiative +2 (12)|HP 75 (10d8 + 30)|Speed 30 ft., Climb 30 ft., Fly 60 ft.|Str 19 +4 +4 Dex 10 +0 +2 Con 17 +3 +3|Int 12 +1 +1 Wis 11 +0 +2 Cha 15 +2 +2|Skills Perception +4, Stealth +2|Immunities Fire|Senses Blindsight 10 ft., Darkvision 60 ft.; Passive Perception 14|Languages Draconic|CR 4 (XP 1,100; PB +2)"},
}

// Typed values of every trait and action entry, reviewed against the same pages:
// limited uses, attack bonus/reach/range/damage, save ability and DC, target
// shape size/width/range, dice and dependency keys. A ranged attack printed with
// one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text. The Priest Acolyte's Spellcasting prints no spell save DC, so its typed
// DC is zero.
var monsters170Entries = map[string][]string{
	"Polar Bear": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +7 reach 5 range 0/0 9 (1d8 + 5) slashing; effect/slashing-damage 9 (1d8 + 5)",
	},
	"Pony": {
		"action|hooves|Hooves; melee +4 reach 5 range 0/0 4 (1d4 + 2) bludgeoning; effect/bludgeoning-damage 4 (1d4 + 2)",
	},
	"Priest": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|mace|Mace; melee +5 reach 5 range 0/0 6 (1d6 + 3) bludgeoning; effect/bludgeoning-damage 6 (1d6 + 3), effect/radiant-damage 5 (2d4)",
		"action|radiant-flame|Radiant Flame; ranged +5 reach 0 range 60/60 11 (2d10) radiant; effect/radiant-damage 11 (2d10)",
		"action|spellcasting|Spellcasting; wis DC 13 +0; at will: Light, Thaumaturgy; 1/day: Spirit Guardians; spell/light [spellcasting/at-will], spell/thaumaturgy [spellcasting/at-will], spell/spirit-guardians [limit/per-day]",
		"bonus|divine-aid|Divine Aid (3/Day); limit/per-day, action/spellcasting, spell/bless, spell/dispel-magic, spell/healing-word, spell/lesser-restoration",
	},
	"Priest Acolyte": {
		"action|mace|Mace; melee +4 reach 5 range 0/0 5 (1d6 + 2) bludgeoning; effect/bludgeoning-damage 5 (1d6 + 2), effect/radiant-damage 2 (1d4)",
		"action|radiant-flame|Radiant Flame; ranged +4 reach 0 range 60/60 7 (2d6) radiant; effect/radiant-damage 7 (2d6)",
		"action|spellcasting|Spellcasting; wis DC 0 +0; at will: Light, Thaumaturgy; spell/light [spellcasting/at-will], spell/thaumaturgy [spellcasting/at-will]",
		"bonus|divine-aid|Divine Aid (1/Day); limit/per-day, action/spellcasting, spell/bless, spell/healing-word, spell/sanctuary",
	},
	"Pseudodragon": {
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +4 reach 5 range 0/0 4 (1d4 + 2) piercing; effect/piercing-damage 4 (1d4 + 2)",
		"action|sting|Sting; con DC 12; creature 0/0/5; save/con, target/creature, effect/poison-damage 5 (2d4), condition/poisoned, effect/failure-by-five-or-more, condition/unconscious, effect/unconscious-ends-early",
	},
	"Pteranodon": {
		"trait|trait/flyby|Flyby; effect/no-opportunity-attack-flying-out",
		"action|bite|Bite; melee +4 reach 5 range 0/0 6 (1d8 + 2) piercing; effect/piercing-damage 6 (1d8 + 2)",
	},
	"Purple Worm": {
		"trait|trait/tunneler|Tunneler; movement/tunnel-through-rock",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +14 reach 10 range 0/0 22 (3d8 + 9) piercing; effect/piercing-damage 22 (3d8 + 9), requirement/target-size, condition/grappled, condition/restrained",
		"action|tail-stinger|Tail Stinger; melee +14 reach 10 range 0/0 16 (2d6 + 9) piercing; effect/piercing-damage 16 (2d6 + 9), effect/poison-damage 35 (10d6)",
		"bonus|swallow|Swallow; str DC 19; creature 0/0/0; save/str, target/creature, requirement/target-size, requirement/grappled-creature, limit/three-swallowed-creatures, effect/swallow, condition/blinded, condition/restrained, effect/total-cover, effect/acid-damage 17 (5d6), effect/regurgitation-save, condition/prone, effect/escape-corpse",
	},
	"Quasit": {
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|rend|Rend; melee +5 reach 5 range 0/0 5 (1d4 + 3) slashing; effect/slashing-damage 5 (1d4 + 3), condition/poisoned",
		"action|invisibility|Invisibility; action/spellcasting, spell/invisibility, spellcasting/no-components, spellcasting/charisma-ability",
		"action|scare|Scare (1/Day); wis DC 10; creature 0/0/20; limit/per-day, save/wis, target/creature, condition/frightened, save/repeat, save/automatic-success",
		"action|shape-shift|Shape-Shift; effect/shape-shift, effect/shapechange-form-limits, effect/equipment-not-transformed",
	},
	"Rakshasa": {
		"trait|trait/greater-magic-resistance|Greater Magic Resistance; save/automatic-success-against-magic, effect/spell-attacks-miss, effect/scrying-immunity",
		"trait|trait/fiendish-restoration|Fiendish Restoration; effect/revive",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|cursed-touch|Cursed Touch; melee +10 reach 5 range 0/0 12 (2d6 + 5) slashing; effect/slashing-damage 12 (2d6 + 5), effect/necrotic-damage 19 (3d12), requirement/target-creature, effect/curse, effect/no-rest-benefit",
		"action|baleful-command|Baleful Command (Recharge 5–6); wis DC 18; emanation 30/0/0; limit/recharge, save/wis, area/emanation, target/enemies-only, effect/psychic-damage 28 (8d6), condition/frightened, condition/incapacitated",
		"action|spellcasting|Spellcasting; cha DC 18 +0; at will: Detect Magic, Detect Thoughts, Disguise Self, Mage Hand, Minor Illusion; 1/day: Fly, Invisibility, Major Image, Plane Shift; spellcasting/no-material-components, spell/detect-magic [spellcasting/at-will], spell/detect-thoughts [spellcasting/at-will], spell/disguise-self [spellcasting/at-will], spell/mage-hand [spellcasting/at-will], spell/minor-illusion [spellcasting/at-will], spell/fly [limit/per-day], spell/invisibility [limit/per-day], spell/major-image [limit/per-day], spell/plane-shift [limit/per-day]",
	},
	"Rat": {
		"trait|trait/agile|Agile; effect/no-opportunity-attack-moving-out",
		"action|bite|Bite; melee +2 reach 5 range 0/0 1 piercing; effect/piercing-damage [damage/fixed]",
	},
	"Raven": {
		"trait|trait/mimicry|Mimicry; effect/mimic-sounds, check/wis-insight",
		"action|beak|Beak; melee +4 reach 5 range 0/0 1 piercing; effect/piercing-damage [damage/fixed]",
	},
	"Red Dragon Wyrmling": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +6 reach 5 range 0/0 9 (1d10 + 4) slashing; effect/slashing-damage 9 (1d10 + 4), effect/fire-damage 3 (1d6)",
		"action|fire-breath|Fire Breath (Recharge 5–6); dex DC 13; cone 15/0/0; limit/recharge, save/dex, area/cone, effect/fire-damage 24 (7d6), save/half-damage",
	},
}

// printedSkill170 lists skill bonuses the PDF prints that no ability modifier
// plus PB or twice PB explains. The record keeps the printed value: the Shambling
// Mound's Stealth +3 (p. 323) with Dexterity −1 and PB +3.
var printedSkill170 = map[string]int{"Shambling Mound/stealth": 3}

// xpBeforeLabel170 prints "(450 XP; PB +2)" instead of "(XP 450; PB +2)"
// (p. 324), as the Gold Dragon Wyrmling does.
var xpBeforeLabel170 = map[string]bool{"Silver Dragon Wyrmling": true}

// monsterHeader170 prints typed fields in the source's stat-block notation. It
// extends monsterHeader164 with a swarm type line, retained notes, a hover mark
// inside the speed list and the lower-case word of "Sleight of Hand".
func monsterHeader170(m Monster) string {
	kind := monsterTitle151(m.Type)
	if m.SwarmOf != "" {
		plural := map[string]string{"beast": "Beasts", "undead": "Undead"}[m.Type]
		kind = "Swarm of " + monsterTitle151(m.SwarmOf) + " " + plural
	}
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
		"HP " + monsterDice163(m.HP),
	}
	speed := fmt.Sprintf("Speed %d ft.", m.Speeds["walk"])
	for _, mode := range monsterSorted151(m.Speeds) {
		if mode != "walk" {
			speed += fmt.Sprintf(", %s %d ft.", monsterTitle151(mode), m.Speeds[mode])
		}
		if mode == "fly" && m.Hover {
			speed += " (hover)"
		}
	}
	if m.SpeedNote != "" {
		speed += ", " + m.SpeedNote
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
			name := strings.ReplaceAll(monsterTitle151(skill), " Of ", " of ")
			skills = append(skills, name+" "+monsterSigned151(m.Skills[skill]))
		}
		lines = append(lines, "Skills "+strings.Join(skills, ", "))
	}
	for _, line := range []struct{ label, value string }{
		{"Vulnerabilities", monsterTitles151(m.Vulnerabilities) + m.VulnerabilityNote},
		{"Resistances", monsterTitles151(m.Resistances) + m.ResistanceNote},
	} {
		if line.value != "" {
			lines = append(lines, line.label+" "+line.value)
		}
	}
	if immunities := slices.DeleteFunc([]string{monsterTitles151(m.Immunities), monsterTitles151(m.ConditionImmunities) + m.ConditionImmunityNote}, func(s string) bool { return s == "" }); len(immunities) > 0 {
		lines = append(lines, "Immunities "+strings.Join(immunities, "; "))
	}
	if len(m.Gear) > 0 {
		lines = append(lines, "Gear "+strings.Join(m.Gear, ", "))
	}
	senses := []string{}
	for _, sense := range monsterSorted151(m.Senses.Special) {
		senses = append(senses, fmt.Sprintf("%s %d ft.", monsterTitle151(sense), m.Senses.Special[sense]))
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
	}
	xp := "XP " + monsterThousands151(m.XP)
	if xpBeforeLabel170[m.Name] {
		xp = monsterThousands151(m.XP) + " XP"
	}
	if m.LairXP > 0 {
		xp += ", or " + monsterThousands151(m.LairXP) + " in lair"
	}
	lines = append(lines, "Senses "+passive, languages, fmt.Sprintf("CR %s (%s; PB +%d)", m.CR, xp, m.ProficiencyBonus))
	return strings.Join(lines, "|")
}

// attackAbility170 is the ability behind a typed attack that does not use the
// higher of Strength and Dexterity, or whose hit adds no ability modifier.
type attackAbility170 struct {
	ability string
	damage  bool // the printed damage adds the same modifier
}

var skillAbility170 = map[string]string{
	"acrobatics": "dex", "animal-handling": "wis", "arcana": "int", "athletics": "str", "deception": "cha", "history": "int",
	"insight": "wis", "intimidation": "cha", "investigation": "int", "medicine": "wis", "nature": "int", "perception": "wis",
	"performance": "cha", "persuasion": "cha", "religion": "int", "sleight-of-hand": "dex", "stealth": "dex", "survival": "wis",
}

// monsterConsistency170 checks the source arithmetic of one block (PDF
// pp. 254–256): XP and PB follow CR, Hit Dice follow size and Constitution,
// saves add nothing or PB to the ability modifier, skills add PB or twice PB,
// and attack rolls add PB to the listed ability. Initiative adds nothing, PB or
// twice PB to Dexterity. A lair XP value is the next rating's XP.
func monsterConsistency170(t *testing.T, m Monster, pb map[string]int, attacks map[string]attackAbility170) {
	t.Helper()
	p := m.ProficiencyBonus
	if !encounterValidXP(m.CR, int64(m.XP)) || pb[m.CR] != p {
		t.Fatal("XP or PB does not follow CR")
	}
	if m.LairXP > 0 {
		rating, err := strconv.Atoi(m.CR)
		if err != nil || !encounterValidXP(strconv.Itoa(rating+1), int64(m.LairXP)) {
			t.Fatal("lair XP is not the next rating's XP")
		}
	}
	hitDie := map[string]int{"tiny": 4, "small": 6, "medium": 8, "large": 10, "huge": 12, "gargantuan": 20}
	if m.HP.Sides != hitDie[m.Sizes[0]] || m.HP.Modifier != m.HP.Dice*m.Abilities["con"].Modifier {
		t.Fatal("Hit Dice do not follow size and Constitution")
	}
	for key, a := range m.Abilities {
		if a.Save != a.Modifier && a.Save != a.Modifier+p {
			t.Fatalf("%s save is not the modifier with optional PB", key)
		}
	}
	if extra := m.Initiative.Modifier - m.Abilities["dex"].Modifier; extra != 0 && extra != p && extra != 2*p {
		t.Fatal("initiative is not Dexterity with optional PB or twice PB")
	}
	passive := 10 + m.Abilities["wis"].Modifier
	for skill, bonus := range m.Skills {
		ability, ok := skillAbility170[skill]
		extra := bonus - m.Abilities[ability].Modifier
		if printed, odd := printedSkill170[m.Name+"/"+skill]; odd {
			if bonus != printed || extra == p || extra == 2*p {
				t.Fatalf("%s no longer holds the printed exception", skill)
			}
		} else if !ok || extra != p && extra != 2*p {
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
		damage := mod
		if other, ok := attacks[m.Name+"/"+a.Name]; ok {
			mod = m.Abilities[other.ability].Modifier
			damage = 0
			if other.damage {
				damage = mod
			}
		}
		// A fixed hit adds no ability modifier.
		if fixedMonsterDamage(a.Attack.Damage) {
			damage = 0
		}
		if a.Attack.Bonus != mod+p || a.Attack.Damage.Modifier != damage {
			t.Fatalf("%s attack does not use the ability modifier and PB", a.Name)
		}
	}
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters170SourceValues(t *testing.T) {
	c := loadMonsters170(t)
	if len(c.Monsters) != 12 || len(monsters170Headers) != 12 || len(monsters170Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #170") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters170Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader170(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries164(m); !slices.Equal(got, monsters170Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters170Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.ResistanceNote != "" || m.ConditionImmunityNote != "" || m.TelepathyNote != "" || m.Senses.Notes != nil || m.SwarmOf != "" || m.SpeedNote != "" || m.LairXP != 0 || m.Legendary != nil || m.LegendaryActions != nil {
				t.Fatal("source lists no vulnerabilities, resistance note, condition immunity note, telepathy note, sense notes, swarm line, speed note, lair XP, legendary actions")
			}
		})
		for _, line := range monsters170Entries[want.name] {
			if !strings.HasPrefix(line, "legendary-uses|") {
				entries++
			}
		}
	}
	// Bold-italic entry headings counted on PDF pp. 316–318 and 359.
	if entries != 41 {
		t.Fatalf("reviewed %d entries, want 41", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters170MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 170 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters170(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Animals appendix blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("animals animals monsters monsters monsters animals monsters monsters monsters animals animals monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #170 cites the page each block starts on, in this order.
	names := []string{"Polar Bear", "Pony", "Priest", "Priest Acolyte", "Pseudodragon", "Pteranodon", "Purple Worm", "Quasit", "Rakshasa", "Rat", "Raven", "Red Dragon Wyrmling"}
	pages := []int{359, 359, 316, 316, 316, 359, 316, 317, 317, 359, 359, 318}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic for these twelve blocks; see monsterConsistency170.
func TestMonsters170DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"2": 2, "1/8": 2, "1/4": 2, "15": 5, "1": 2, "13": 5, "0": 2, "4": 2}
	attacks := map[string]attackAbility170{
		"Priest/Radiant Flame":         {"wis", false},
		"Priest Acolyte/Radiant Flame": {"wis", false},
		"Rakshasa/Cursed Touch":        {"cha", true},
	}
	for _, m := range loadMonsters170(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			monsterConsistency170(t, m, pb, attacks)
		})
	}
}

// Every listed key is used and required, every effect quotes its own entry, and
// every dice expression is typed. Removing any one key invalidates the record.
func TestMonsters170DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters170(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			if m.RequiredMechanics[0] != "monster/stat-block" {
				t.Fatal("stat-block dependency must come first")
			}
			quoted := func(parent string, x MonsterMechanic) {
				if !strings.Contains(parent, strings.TrimSuffix(x.Text, ".")) {
					t.Fatalf("%q is not quoted from its entry", x.Text)
				}
				if x.Dice != nil && !strings.Contains(x.Text, monsterDice163(*x.Dice)) {
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
				if !strings.HasPrefix(a.Text, a.Name) || a.ID != strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(a.Name), "’", ""), " ", "-") {
					t.Fatalf("action %q ID or text differs", a.Name)
				}
				for _, e := range a.Effects {
					quoted(a.Text, e)
				}
				if a.Attack == nil {
					continue
				}
				first := a.Effects[0]
				if fixedMonsterDamage(a.Attack.Damage) {
					if first.Dice != nil || !slices.Equal(first.Requires, []string{"damage/fixed"}) || !regexp.MustCompile(fmt.Sprintf(`Hit: %d %s damage[.,]`, a.Attack.Damage.Average, monsterTitle151(a.Attack.DamageType))).MatchString(a.Text) {
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
// keys cover none of these records, and the register still reports
// both families unsupported.
func TestMonsters170RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters170(t).Monsters {
		if m.Available() || m.availableWith(bandit) {
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

func TestMonsters170MalformedFixtures(t *testing.T) {
	tests := []struct {
		name, record string
		mutate       func(*Monster)
	}{
		{"missing languages list", "Rat", func(m *Monster) { m.Languages = nil }},
		{"animal or monster block in another family", "Raven", func(m *Monster) {
			m.ID = strings.Replace(strings.Replace(m.ID, "/animals/", "/spells/", 1), "/monsters/", "/spells/", 1)
		}},
		{"XP off CR", "Polar Bear", func(m *Monster) { m.XP++ }},
		{"save off modifier", "Polar Bear", func(m *Monster) {
			m.Abilities["wis"] = MonsterAbility{m.Abilities["wis"].Score, m.Abilities["wis"].Modifier, m.Abilities["wis"].Modifier + 11}
		}},
		{"modifier", "Priest", func(m *Monster) {
			m.Abilities["str"] = MonsterAbility{m.Abilities["str"].Score, m.Abilities["str"].Modifier + 1, m.Abilities["str"].Modifier + 1}
		}},
		{"HP average", "Priest", func(m *Monster) { m.HP.Average++ }},
		{"vulnerability note beside a typed list", "Rakshasa", func(m *Monster) { m.Vulnerabilities = []string{"piercing"} }},
		{"padded vulnerability note", "Rakshasa", func(m *Monster) { m.VulnerabilityNote += " " }},
		{"recharge", "Red Dragon Wyrmling", func(m *Monster) { m.Actions[2].Uses.Recharge = 7 }},
		{"save DC off entry text", "Quasit", func(m *Monster) { m.Actions[2].Save.DC = 11 }},
		{"cone width", "Red Dragon Wyrmling", func(m *Monster) { m.Actions[2].Target.Width = 5 }},
		{"unlisted spell", "Priest", func(m *Monster) { m.Actions[3].Spellcasting.AtWill[0] = "Wish" }},
		{"spell DC off entry text", "Priest", func(m *Monster) { m.Actions[3].Spellcasting.SaveDC = 14 }},
		{"spell DC the source does not print", "Priest Acolyte", func(m *Monster) { m.Actions[2].Spellcasting.SaveDC = 12 }},
		{"daily uses without count", "Priest", func(m *Monster) { m.BonusActions[0].Uses.PerDay = 0 }},
		{"fixed damage of zero", "Rat", func(m *Monster) { m.Actions[0].Attack.Damage.Average = 0 }},
		{"empty bonus actions", "Pony", func(m *Monster) { m.BonusActions = []MonsterAction{} }},
		{"empty reactions", "Pony", func(m *Monster) { m.Reactions = []MonsterAction{} }},
	}
	record := func(t *testing.T, c *MonsterCatalogue, name string) *Monster {
		i := slices.IndexFunc(c.Monsters, func(m Monster) bool { return m.Name == name })
		if i < 0 {
			t.Fatalf("no %s record", name)
		}
		return &c.Monsters[i]
	}
	rejected := func(t *testing.T, c MonsterCatalogue, encodes bool) {
		if c.Validate() == nil {
			t.Fatal("accepted malformed catalogue")
		}
		if !encodes {
			return
		}
		data, err := json.Marshal(c)
		if err != nil {
			t.Fatal(err)
		}
		if decoded, err := DecodeMonsters(bytes.NewReader(data)); err == nil || len(decoded.Monsters) != 0 {
			t.Fatal("decoder must return no catalogue on invalid data")
		}
	}
	t.Run("duplicate ID", func(t *testing.T) {
		c := loadMonsters170(t)
		c.Monsters[1].ID = c.Monsters[2].ID
		rejected(t, c, true)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters170(t)
			tt.mutate(record(t, &c, tt.record))
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			rejected(t, c, !strings.HasPrefix(tt.name, "empty "))
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_170.json")
	if err != nil {
		t.Fatal(err)
	}
	// The fixtures are indented; compact them so each edit matches one field.
	var compact bytes.Buffer
	if err := json.Compact(&compact, raw); err != nil {
		t.Fatal(err)
	}
	for _, edit := range [][2]string{
		{`"dice":0,`, ``},
		{`"sides":0,`, `"sides":null,`},
		{`"long_range":0,`, `"long_range":null,`},
		{`"recharge":5`, `"recharge":null`},
		{`"per_day":0,`, ``},
		{`"dc":12`, `"dc":12,"dc":22`},
		{`"width":0,`, ``},
		{`"save_dc":13,`, ``},
		{`"attack_bonus":0`, `"attack_bonus":0,"attack_bonus":5`},
		{`"darkvision":60`, `"darkvision":null`},
		{`"shape":"creature",`, `"shape":"creature","height":5,`},
		{`"sizes":[`, `"colour":"grey","sizes":[`},
	} {
		if !bytes.Contains(compact.Bytes(), []byte(edit[0])) {
			t.Fatalf("fixture lost %s", edit[0])
		}
		if _, err := DecodeMonsters(bytes.NewReader(bytes.Replace(compact.Bytes(), []byte(edit[0]), []byte(edit[1]), 1))); err == nil {
			t.Fatalf("accepted %s -> %s", edit[0], edit[1])
		}
	}
}

// The batch composes with the #143 Bandit catalogue and stat-block batches 01
// to 19 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters170ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters170(t)
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
	catalogues := []MonsterCatalogue{base, first, third, fourth, fifth, sixth, seventh}
	var later []Monster
	for _, load := range []func(*testing.T) MonsterCatalogue{loadMonsters158, loadMonsters159, loadMonsters160, loadMonsters161, loadMonsters162, loadMonsters163, loadMonsters164, loadMonsters165, loadMonsters166, loadMonsters167, loadMonsters168, loadMonsters169} {
		c := load(t)
		catalogues = append(catalogues, c)
		later = append(later, c.Monsters...)
	}
	for _, c := range catalogues {
		a, b := c.Source, batch.Source
		a.Changes, b.Changes = "", ""
		if a != b {
			t.Fatal("batches must share one source pin")
		}
	}
	earlier := slices.Concat(base.Monsters, first.Monsters, second, third.Monsters, slices.Concat(fourth.Monsters[:4], fourth.Monsters[5:]), fifth.Monsters, sixth.Monsters, seventh.Monsters, later)
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
	for _, id := range []string{"beak", "bite", "cursed-touch", "hooves", "mace", "radiant-flame", "rend", "tail-stinger"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
