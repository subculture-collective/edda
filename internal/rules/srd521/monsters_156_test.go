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

func loadMonsters156(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_156.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #156 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Chain Devil's block starts on p. 272 and the Chuul's on p. 273.
var monsters156Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Bugbear Warrior", [2]int{272, 272}, "Medium Fey (Goblinoid), Chaotic Evil|AC 14 Initiative +2 (12)|HP 33 (6d8 + 6)|Speed 30 ft.|Str 15 +2 +2 Dex 14 +2 +2 Con 13 +1 +1|Int 8 −1 −1 Wis 11 +0 +0 Cha 9 −1 −1|Skills Stealth +6, Survival +2|Gear Hide Armor, Light Hammers (3)|Senses Darkvision 60 ft.; Passive Perception 10|Languages Common, Goblin|CR 1 (XP 200; PB +2)"},
	{"Bulette", [2]int{272, 272}, "Large Monstrosity, Unaligned|AC 17 Initiative +0 (10)|HP 94 (9d10 + 45)|Speed 40 ft., Burrow 40 ft.|Str 19 +4 +4 Dex 11 +0 +0 Con 21 +5 +5|Int 2 −4 −4 Wis 10 +0 +0 Cha 5 −3 −3|Skills Perception +6|Senses Darkvision 60 ft., Tremorsense 120 ft.; Passive Perception 16|Languages None|CR 5 (XP 1,800; PB +3)"},
	{"Camel", [2]int{346, 346}, "Large Beast, Unaligned|AC 10 Initiative −1 (9)|HP 17 (2d10 + 6)|Speed 50 ft.|Str 15 +2 +2 Dex 8 −1 −1 Con 17 +3 +5|Int 2 −4 −4 Wis 11 +0 +0 Cha 5 −3 −3|Senses Darkvision 60 ft.; Passive Perception 10|Languages None|CR 1/8 (XP 25; PB +2)"},
	{"Cat", [2]int{346, 346}, "Tiny Beast, Unaligned|AC 12 Initiative +2 (12)|HP 2 (1d4)|Speed 40 ft., Climb 40 ft.|Str 3 −4 −4 Dex 15 +2 +4 Con 10 +0 +0|Int 3 −4 −4 Wis 12 +1 +1 Cha 7 −2 −2|Skills Perception +3, Stealth +4|Senses Darkvision 60 ft.; Passive Perception 13|Languages None|CR 0 (XP 10; PB +2)"},
	{"Centaur Trooper", [2]int{272, 272}, "Large Fey, Neutral Good|AC 16 Initiative +2 (12)|HP 45 (6d10 + 12)|Speed 50 ft.|Str 18 +4 +4 Dex 14 +2 +2 Con 14 +2 +2|Int 9 −1 −1 Wis 13 +1 +1 Cha 11 +0 +0|Skills Athletics +6, Perception +3|Gear Breastplate, Longbow, Pike|Senses Passive Perception 13|Languages Elvish, Sylvan|CR 2 (XP 450; PB +2)"},
	{"Chain Devil", [2]int{272, 273}, "Medium Fiend (Devil), Lawful Evil|AC 15 Initiative +5 (15)|HP 85 (10d8 + 40)|Speed 30 ft.|Str 18 +4 +4 Dex 15 +2 +2 Con 18 +4 +7|Int 11 +0 +0 Wis 12 +1 +4 Cha 14 +2 +2|Resistances Bludgeoning, Cold, Piercing, Slashing|Immunities Fire, Poison; Poisoned|Senses Darkvision 120 ft. (unimpeded by magical Darkness); Passive Perception 11|Languages Infernal; telepathy 120 ft.|CR 8 (XP 3,900; PB +3)"},
	{"Chimera", [2]int{273, 273}, "Large Monstrosity, Chaotic Evil|AC 14 Initiative +0 (10)|HP 114 (12d10 + 48)|Speed 30 ft., Fly 60 ft.|Str 19 +4 +4 Dex 11 +0 +0 Con 19 +4 +4|Int 3 −4 −4 Wis 14 +2 +2 Cha 10 +0 +0|Skills Perception +8|Senses Darkvision 60 ft.; Passive Perception 18|Languages Understands Draconic but can’t speak|CR 6 (XP 2,300; PB +3)"},
	{"Chuul", [2]int{273, 274}, "Large Aberration, Chaotic Evil|AC 16 Initiative +0 (10)|HP 76 (9d10 + 27)|Speed 30 ft., Swim 30 ft.|Str 19 +4 +4 Dex 10 +0 +0 Con 16 +3 +3|Int 5 −3 −3 Wis 11 +0 +0 Cha 5 −3 −3|Skills Perception +4|Immunities Poison; Poisoned|Senses Darkvision 60 ft.; Passive Perception 14|Languages Understands Deep Speech but can’t speak|CR 4 (XP 1,100; PB +2)"},
	{"Clay Golem", [2]int{274, 274}, "Large Construct, Unaligned|AC 14 Initiative +3 (13)|HP 123 (13d10 + 52)|Speed 30 ft.|Str 20 +5 +5 Dex 9 −1 −1 Con 18 +4 +4|Int 3 −4 −4 Wis 8 −1 −1 Cha 1 −5 −5|Resistances Bludgeoning, Piercing, Slashing|Immunities Acid, Poison, Psychic; Charmed, Exhaustion, Frightened, Paralyzed, Petrified, Poisoned|Senses Darkvision 60 ft.; Passive Perception 9|Languages Common plus one other language|CR 9 (XP 5,000; PB +4)"},
	{"Cloaker", [2]int{274, 274}, "Large Aberration, Chaotic Neutral|AC 14 Initiative +5 (15)|HP 91 (14d10 + 14)|Speed 10 ft., Fly 40 ft.|Str 17 +3 +3 Dex 15 +2 +2 Con 12 +1 +1|Int 13 +1 +1 Wis 14 +2 +2 Cha 7 −2 −2|Skills Stealth +5|Immunities Frightened|Senses Darkvision 120 ft.; Passive Perception 12|Languages Deep Speech, Undercommon|CR 8 (XP 3,900; PB +3)"},
	{"Cloud Giant", [2]int{275, 275}, "Huge Giant, Neutral|AC 14 Initiative +4 (14)|HP 200 (16d12 + 96)|Speed 40 ft., Fly 20 ft. (hover)|Str 27 +8 +8 Dex 10 +0 +0 Con 22 +6 +10|Int 12 +1 +1 Wis 16 +3 +7 Cha 16 +3 +3|Skills Insight +7, Perception +11|Senses Passive Perception 21|Languages Common, Giant|CR 9 (XP 5,000; PB +4)"},
	{"Cockatrice", [2]int{275, 275}, "Small Monstrosity, Unaligned|AC 11 Initiative +1 (11)|HP 22 (5d6 + 5)|Speed 20 ft., Fly 40 ft.|Str 6 −2 −2 Dex 12 +1 +1 Con 12 +1 +1|Int 2 −4 −4 Wis 13 +1 +1 Cha 5 −3 −3|Immunities Petrified|Senses Darkvision 60 ft.; Passive Perception 11|Languages None|CR 1/2 (XP 100; PB +2)"},
}

// Typed values of every trait, action, bonus action and reaction, reviewed
// against the same pages: limited uses, attack bonus/reach/range/damage, save
// ability and DC, target shape size/width/range, dice and dependency keys. A
// ranged attack printed with one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text: the Bugbear Warrior's Light Hammer is a "Melee or Ranged Attack Roll"
// (key action/melee-or-ranged-attack, damage typed on its effect, as for the
// #153 Archmage); the Cloaker's Phantasms prints "(Recharge after a Short or
// Long Rest)", kept as effect limit/recharge-after-rest instead of typed uses.
var monsters156Entries = map[string][]string{
	"Bugbear Warrior": {
		"trait|trait/abduct|Abduct; movement/no-extra-cost-moving-grappled",
		"action|grab|Grab; melee +4 reach 10 range 0/0 9 (2d6 + 2) bludgeoning; effect/bludgeoning-damage 9 (2d6 + 2), condition/grappled, requirement/target-size",
		"action|light-hammer|Light Hammer; effect/bludgeoning-damage 9 (3d4 + 2), attack/advantage-against-grappled-target",
	},
	"Bulette": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +7 reach 5 range 0/0 17 (2d12 + 4) piercing; effect/piercing-damage 17 (2d12 + 4)",
		"action|deadly-leap|Deadly Leap; dex DC 15; effect/jump, requirement/target-size, save/dex, target/creatures-in-destination-space, effect/bludgeoning-damage 19 (3d12), condition/prone, save/half-damage, effect/push",
		"bonus|leap|Leap; effect/jump",
	},
	"Camel": {
		"action|bite|Bite; melee +4 reach 5 range 0/0 4 (1d4 + 2) bludgeoning; effect/bludgeoning-damage 4 (1d4 + 2)",
	},
	"Cat": {
		"trait|trait/jumper|Jumper; movement/jump-uses-dexterity",
		"action|scratch|Scratch; melee +4 reach 5 range 0/0 1 slashing; effect/slashing-damage [damage/fixed]",
	},
	"Centaur Trooper": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|pike|Pike; melee +6 reach 10 range 0/0 9 (1d10 + 4) piercing; effect/piercing-damage 9 (1d10 + 4)",
		"action|longbow|Longbow; ranged +4 reach 0 range 150/600 6 (1d8 + 2) piercing; effect/piercing-damage 6 (1d8 + 2)",
		"bonus|trampling-charge|Trampling Charge (Recharge 5–6); str DC 14; limit/recharge, effect/move-without-opportunity-attacks, movement/through-creature-spaces, target/creatures-entered-spaces, save/str, effect/bludgeoning-damage 7 (1d6 + 4), condition/prone",
	},
	"Chain Devil": {
		"trait|trait/diabolical-restoration|Diabolical Restoration; effect/revive",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/additional-action",
		"action|chain|Chain; melee +7 reach 10 range 0/0 11 (2d6 + 4) slashing; effect/slashing-damage 11 (2d6 + 4), condition/grappled, condition/restrained, requirement/target-size",
		"action|conjure-infernal-chain|Conjure Infernal Chain; dex DC 15; creature 0/0/60; effect/conjure-chain, save/dex, target/creature, effect/fire-damage 9 (2d4 + 4), condition/restrained, requirement/target-size, effect/pull, effect/chain-disappears",
		"reaction|unnerving-gaze|Unnerving Gaze; wis DC 15; reaction/trigger, save/wis, target/triggering-creature, condition/frightened, save/success-immunity",
	},
	"Chimera": {
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack",
		"action|bite|Bite; melee +7 reach 5 range 0/0 11 (2d6 + 4) piercing; effect/piercing-damage 11 (2d6 + 4), effect/piercing-damage-with-advantage 18 (4d6 + 4)",
		"action|claw|Claw; melee +7 reach 5 range 0/0 7 (1d6 + 4) slashing; effect/slashing-damage 7 (1d6 + 4)",
		"action|ram|Ram; melee +7 reach 5 range 0/0 10 (1d12 + 4) bludgeoning; effect/bludgeoning-damage 10 (1d12 + 4), condition/prone, requirement/target-size",
		"action|fire-breath|Fire Breath (Recharge 5–6); dex DC 15; cone 15/0/0; save/dex, area/cone, limit/recharge, effect/fire-damage 31 (7d8), save/half-damage",
	},
	"Chuul": {
		"trait|trait/amphibious|Amphibious; effect/breathe-air-and-water",
		"trait|trait/sense-magic|Sense Magic; effect/sense-magic, spell/detect-magic",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/additional-action",
		"action|pincer|Pincer; melee +6 reach 10 range 0/0 9 (1d10 + 4) bludgeoning; effect/bludgeoning-damage 9 (1d10 + 4), condition/grappled, requirement/target-size",
		"action|paralyzing-tentacles|Paralyzing Tentacles; con DC 13; creature 0/0/0; save/con, target/creature, requirement/grappled-creature, condition/poisoned, save/repeat, save/automatic-success, condition/paralyzed",
	},
	"Clay Golem": {
		"trait|trait/acid-absorption|Acid Absorption; effect/absorb-acid-damage",
		"trait|trait/berserk|Berserk; requirement/bloodied, effect/berserk-roll, effect/berserk",
		"trait|trait/immutable-form|Immutable Form; effect/no-shape-shift",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/conditional-extra-attack",
		"action|slam|Slam; melee +9 reach 5 range 0/0 10 (1d10 + 5) bludgeoning; effect/bludgeoning-damage 10 (1d10 + 5), effect/acid-damage 6 (1d12), effect/hit-point-maximum-reduction",
		"bonus|hasten|Hasten (Recharge 5–6); limit/recharge, action/dash, action/disengage",
	},
	"Cloaker": {
		"trait|trait/light-sensitivity|Light Sensitivity; requirement/in-bright-light, attack/disadvantage",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|attach|Attach; melee +6 reach 5 range 0/0 13 (3d6 + 3) piercing; effect/piercing-damage 13 (3d6 + 3), requirement/target-size, effect/attach, condition/blinded, limit/one-attached-target, effect/halve-damage-taken, effect/share-damage-taken, effect/detach, check/str-athletics",
		"action|tail|Tail; melee +6 reach 10 range 0/0 8 (1d10 + 3) slashing; effect/slashing-damage 8 (1d10 + 3)",
		"bonus|moan|Moan; wis DC 13; emanation 60/0/0; save/wis, area/emanation, condition/frightened, save/success-immunity",
		"bonus|phantasms|Phantasms; limit/recharge-after-rest, action/spellcasting, spell/mirror-image, spellcasting/no-components, effect/spell-ends-in-bright-light",
	},
	"Cloud Giant": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination, multiattack/replace-attack, action/spellcasting",
		"action|thunderous-mace|Thunderous Mace; melee +12 reach 10 range 0/0 21 (3d8 + 8) bludgeoning; effect/bludgeoning-damage 21 (3d8 + 8), effect/thunder-damage 7 (2d6)",
		"action|thundercloud|Thundercloud; ranged +12 reach 0 range 240/240 18 (3d6 + 8) thunder; effect/thunder-damage 18 (3d6 + 8), condition/incapacitated",
		"action|spellcasting|Spellcasting; cha DC 15 +0; at will: Detect Magic, Fog Cloud, Light; 1/day: Control Weather, Gaseous Form, Telekinesis; spellcasting/no-material-components, spell/detect-magic [spellcasting/at-will], spell/fog-cloud [spellcasting/at-will], spell/light [spellcasting/at-will], spell/control-weather [limit/per-day], spell/gaseous-form [limit/per-day], spell/telekinesis [limit/per-day]",
		"bonus|misty-step|Misty Step; action/spellcasting, spell/misty-step",
	},
	"Cockatrice": {
		"action|petrifying-bite|Petrifying Bite; melee +3 reach 5 range 0/0 3 (1d4 + 1) piercing; con DC 11; effect/piercing-damage 3 (1d4 + 1), requirement/target-creature, save/con, condition/restrained, save/repeat, condition/petrified, save/second-failure",
	},
}

// monsterDice156 prints a fixed hit, which the source prints as a bare number
// (the Cat's Scratch), and otherwise the source's "N (XdY + Z)" notation.
func monsterDice156(d MonsterHP) string {
	switch {
	case d.Dice == 0:
		return fmt.Sprint(d.Average)
	case d.Modifier < 0:
		return fmt.Sprintf("%d (%dd%d − %d)", d.Average, d.Dice, d.Sides, -d.Modifier)
	}
	return monsterDice151(d)
}

// monsterHeader156 prints typed fields in the source's stat-block notation.
func monsterHeader156(m Monster) string {
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
		"HP " + monsterDice156(m.HP),
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
	}
	lines = append(lines, "Senses "+passive, languages, fmt.Sprintf("CR %s (XP %s; PB +%d)", m.CR, monsterThousands151(m.XP), m.ProficiencyBonus))
	return strings.Join(lines, "|")
}

func monsterEntry156(section, id string, x MonsterMechanic, attack *MonsterAttack, sc *MonsterSpellcasting, effects []MonsterMechanic) string {
	head := x.Name
	if u := x.Uses; u != nil {
		if u.Recharge > 0 {
			head += fmt.Sprintf(" (Recharge %d–6)", u.Recharge)
		} else {
			head += fmt.Sprintf(" (%d/Day)", u.PerDay)
		}
	}
	parts := []string{head}
	if attack != nil {
		parts = append(parts, fmt.Sprintf("%s +%d reach %d range %d/%d %s %s", attack.Kind, attack.Bonus, attack.Reach, attack.Range, attack.LongRange, monsterDice156(attack.Damage), attack.DamageType))
	}
	if x.Save != nil {
		parts = append(parts, fmt.Sprintf("%s DC %d", x.Save.Ability, x.Save.DC))
	}
	if t := x.Target; t != nil {
		parts = append(parts, fmt.Sprintf("%s %d/%d/%d", t.Shape, t.Size, t.Width, t.Range))
	}
	if x.Dice != nil {
		parts = append(parts, monsterDice156(*x.Dice))
	}
	for _, d := range x.MoreDice {
		parts = append(parts, d.Key+" "+monsterDice156(d.Dice))
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
			key += " " + monsterDice156(*e.Dice)
		}
		if len(e.Requires) > 0 {
			key += " [" + strings.Join(e.Requires, ", ") + "]"
		}
		keys = append(keys, key)
	}
	return section + "|" + id + "|" + strings.Join(append(parts, strings.Join(keys, ", ")), "; ")
}

func monsterEntries156(m Monster) []string {
	out := []string{}
	for _, trait := range m.Traits {
		out = append(out, monsterEntry156("trait", trait.Key, trait, nil, nil, nil))
	}
	for _, section := range []struct {
		label   string
		actions []MonsterAction
	}{{"action", m.Actions}, {"bonus", m.BonusActions}, {"reaction", m.Reactions}} {
		for _, a := range section.actions {
			out = append(out, monsterEntry156(section.label, a.ID, MonsterMechanic{Name: a.Name, Uses: a.Uses, Save: a.Save, Target: a.Target}, a.Attack, a.Spellcasting, a.Effects))
		}
	}
	return out
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters156SourceValues(t *testing.T) {
	c := loadMonsters156(t)
	if len(c.Monsters) != 12 || len(monsters156Headers) != 12 || len(monsters156Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #156") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters156Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader156(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries156(m); !slices.Equal(got, monsters156Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters156Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.Hover != (m.Name == "Cloud Giant") || m.LairXP != 0 || m.LegendaryActions != nil || m.Legendary != nil {
				t.Fatal("source lists no vulnerabilities, lair XP or legendary actions, and only the Cloud Giant hovers")
			}
		})
		entries += len(monsters156Entries[want.name])
	}
	// Bold-italic entry headings counted on PDF pp. 272–275 and 346.
	if entries != 49 {
		t.Fatalf("reviewed %d entries, want 49", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters156MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 156 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters156(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the two p. 346 blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("monsters monsters animals animals monsters monsters monsters monsters monsters monsters monsters monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #156 cites the page each block starts on, in this order.
	names := []string{"Bugbear Warrior", "Bulette", "Camel", "Cat", "Centaur Trooper", "Chain Devil", "Chimera", "Chuul", "Clay Golem", "Cloaker", "Cloud Giant", "Cockatrice"}
	pages := []int{272, 272, 346, 346, 272, 272, 273, 273, 274, 274, 275, 275}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic that holds for these twelve blocks (PDF pp. 254–256):
// XP and PB follow CR, Hit Dice follow size and Constitution, saves and skills
// add PB to the ability modifier, and attack rolls add PB to the listed ability.
// Initiative adds PB (Chain Devil, Clay Golem, Cloaker, Cloud Giant) or nothing.
func TestMonsters156DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"0": 2, "1/8": 2, "1/2": 2, "1": 2, "2": 2, "4": 2, "5": 3, "6": 3, "8": 3, "9": 4}
	skillAbility := map[string]string{"athletics": "str", "insight": "wis", "perception": "wis", "stealth": "dex", "survival": "wis"}
	hitDie := map[string]int{"tiny": 4, "small": 6, "medium": 8, "large": 10, "huge": 12}
	// Attacks that do not use the better of Strength and Dexterity.
	attackAbility := map[string]string{"Centaur Trooper/longbow": "dex"}
	for _, m := range loadMonsters156(t).Monsters {
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
				ability, special := attackAbility[m.Name+"/"+a.ID]
				if special {
					mod = m.Abilities[ability].Modifier
				}
				// A fixed hit adds no ability modifier.
				damage := mod
				if fixedMonsterDamage(a.Attack.Damage) {
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
func TestMonsters156DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters156(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			if m.RequiredMechanics[0] != "monster/stat-block" {
				t.Fatal("stat-block dependency must come first")
			}
			quoted := func(parent string, x MonsterMechanic) {
				if !strings.Contains(parent, strings.TrimSuffix(x.Text, ".")) {
					t.Fatalf("%q is not quoted from its entry", x.Text)
				}
				if x.Dice != nil && !strings.Contains(x.Text, monsterDice156(*x.Dice)) {
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
func TestMonsters156RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters156(t).Monsters {
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

func TestMonsters156MalformedFixtures(t *testing.T) {
	// Indexes: 0 Bugbear Warrior, 1 Bulette, 2 Camel, 3 Cat, 4 Centaur Trooper,
	// 5 Chain Devil, 6 Chimera, 7 Chuul, 8 Clay Golem, 9 Cloaker, 10 Cloud Giant.
	tests := []struct {
		name   string
		mutate func(*MonsterCatalogue)
	}{
		{"duplicate ID", func(c *MonsterCatalogue) { c.Monsters[2].ID = c.Monsters[3].ID }},
		{"animal block in another family", func(c *MonsterCatalogue) { c.Monsters[2].ID = "srd-5.2.1:content/spells/camel" }},
		{"missing languages list", func(c *MonsterCatalogue) { c.Monsters[2].Languages = nil }},
		{"fixed damage with sides", func(c *MonsterCatalogue) { c.Monsters[3].Actions[0].Attack.Damage.Sides = 4 }},
		{"fixed damage with modifier", func(c *MonsterCatalogue) { c.Monsters[3].Actions[0].Attack.Damage.Modifier = 1 }},
		{"fixed damage of zero", func(c *MonsterCatalogue) { c.Monsters[3].Actions[0].Attack.Damage.Average = 0 }},
		{"fixed damage without mechanic", func(c *MonsterCatalogue) {
			m := &c.Monsters[3]
			m.Actions[0].Effects[0].Requires = nil
			m.RequiredMechanics = slices.DeleteFunc(m.RequiredMechanics, func(key string) bool { return key == "damage/fixed" })
		}},
		{"fixed effect dice", func(c *MonsterCatalogue) { c.Monsters[3].Actions[0].Effects[0].Dice = &MonsterHP{Average: 1} }},
		{"sense note without sense", func(c *MonsterCatalogue) {
			c.Monsters[5].Senses.Notes = map[string]string{"blindsight": "unimpeded by magical Darkness"}
		}},
		{"blank sense note", func(c *MonsterCatalogue) { c.Monsters[5].Senses.Notes["darkvision"] = " " }},
		{"padded sense note", func(c *MonsterCatalogue) { c.Monsters[5].Senses.Notes["darkvision"] += " " }},
		{"unlisted sense note", func(c *MonsterCatalogue) { c.Monsters[5].Senses.Notes["darkvision"] = "impeded by magical Darkness" }},
		{"empty sense notes", func(c *MonsterCatalogue) { c.Monsters[5].Senses.Notes = map[string]string{} }},
		{"single range below range", func(c *MonsterCatalogue) { c.Monsters[4].Actions[2].Attack.LongRange = 100 }},
		{"ranged reach", func(c *MonsterCatalogue) { c.Monsters[4].Actions[2].Attack.Reach = 5 }},
		{"melee range", func(c *MonsterCatalogue) { c.Monsters[0].Actions[0].Attack.Range = 5 }},
		{"further dice without first", func(c *MonsterCatalogue) {
			c.Monsters[7].Traits[0].MoreDice = []MonsterKeyedDice{{Key: "effect/future-damage", Dice: MonsterHP{Average: 3, Dice: 1, Sides: 4, Modifier: 1}}}
		}},
		{"XP off CR", func(c *MonsterCatalogue) { c.Monsters[5].XP = 3901 }},
		{"save off modifier", func(c *MonsterCatalogue) { c.Monsters[5].Abilities["con"] = MonsterAbility{18, 4, 5} }},
		{"modifier", func(c *MonsterCatalogue) { c.Monsters[2].Abilities["str"] = MonsterAbility{15, 1, 1} }},
		{"HP average", func(c *MonsterCatalogue) { c.Monsters[3].HP.Average = 3 }},
		{"hover without flight", func(c *MonsterCatalogue) { c.Monsters[10].Speeds["fly"] = 0 }},
		{"empty bonus actions", func(c *MonsterCatalogue) { c.Monsters[2].BonusActions = []MonsterAction{} }},
		{"empty reactions", func(c *MonsterCatalogue) { c.Monsters[2].Reactions = []MonsterAction{} }},
		{"duplicate action across sections", func(c *MonsterCatalogue) { c.Monsters[1].BonusActions[0].ID = "bite" }},
		{"recharge", func(c *MonsterCatalogue) { c.Monsters[4].BonusActions[0].Uses.Recharge = 1 }},
		{"daily uses with recharge", func(c *MonsterCatalogue) { c.Monsters[4].BonusActions[0].Uses.PerDay = 1 }},
		{"save DC off entry text", func(c *MonsterCatalogue) { c.Monsters[1].Actions[2].Save.DC = 16 }},
		{"save ability off entry text", func(c *MonsterCatalogue) { c.Monsters[7].Actions[2].Save.Ability = "str" }},
		{"creature size", func(c *MonsterCatalogue) { c.Monsters[5].Actions[2].Target.Size = 5 }},
		{"cone width", func(c *MonsterCatalogue) { c.Monsters[6].Actions[4].Target.Width = 5 }},
		{"effect dice", func(c *MonsterCatalogue) { c.Monsters[6].Actions[1].Effects[1].Dice.Average = 19 }},
		{"unlisted reaction effect", func(c *MonsterCatalogue) { c.Monsters[5].Reactions[0].Effects[0].Key = "reaction/future" }},
		{"untyped spellcasting", func(c *MonsterCatalogue) { c.Monsters[10].Actions[3].Spellcasting = nil }},
		{"misplaced spellcasting", func(c *MonsterCatalogue) {
			c.Monsters[10].Actions[1].Spellcasting = c.Monsters[10].Actions[3].Spellcasting
		}},
		{"spell ability", func(c *MonsterCatalogue) { c.Monsters[10].Actions[3].Spellcasting.Ability = "" }},
		{"spell DC", func(c *MonsterCatalogue) { c.Monsters[10].Actions[3].Spellcasting.SaveDC = 0 }},
		{"unlisted spell", func(c *MonsterCatalogue) { c.Monsters[10].Actions[3].Spellcasting.AtWill[0] = "Wish" }},
		{"duplicate spell", func(c *MonsterCatalogue) { c.Monsters[10].Actions[3].Spellcasting.AtWill[0] = "Fog Cloud" }},
		{"daily spell uses", func(c *MonsterCatalogue) { c.Monsters[10].Actions[3].Spellcasting.PerDay[0].Uses = 0 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters156(t)
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
	raw, err := os.ReadFile("testdata/monsters_156.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, edit := range [][2]string{
		{`"dice": 0,`, ``},
		{`"sides": 0,`, `"sides": null,`},
		{`"long_range": 600,`, ``},
		{`"long_range": 0,`, `"long_range": null,`},
		{`"recharge": 5`, `"recharge": null`},
		{`"per_day": 0,`, ``},
		{`"dc": 15`, `"dc": 15, "dc": 16`},
		{`"width": 0,`, ``},
		{`"notes": {`, `"notes": {"truesight": "see invisible", `},
		{`"shape": "cone",`, `"shape": "cone", "height": 5,`},
		{`"hover": true,`, `"hover": true, "flies": true,`},
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
// to 04 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters156ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters156(t)
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
	for _, c := range []MonsterCatalogue{base, first, third, fourth} {
		a, b := c.Source, batch.Source
		a.Changes, b.Changes = "", ""
		if a != b {
			t.Fatal("batches must share one source pin")
		}
	}
	earlier := slices.Concat(base.Monsters, first.Monsters, second, third.Monsters, slices.Concat(fourth.Monsters[:4], fourth.Monsters[5:]))
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
	for _, id := range []string{"grab", "light-hammer", "slam", "trampling-charge", "unnerving-gaze"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
