package srd521

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"
)

func loadMonsters163(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_163.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #163 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Gold Dragon Wyrmling's and the Gorgon's blocks run onto the next
// page. Two values are printed in a form the typed fields do not repeat: the Gray
// Ooze's "Initiative −2 (13)" (p. 293; the record keeps the score 8 that the
// modifier gives) and the Gold Dragon Wyrmling's "CR 3 (700 XP; PB +2)" (p. 290).
var monsters163Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Gnoll Warrior", [2]int{289, 289}, "Medium Fiend, Chaotic Evil|AC 15 Initiative +1 (11)|HP 27 (6d8)|Speed 30 ft.|Str 14 +2 +2 Dex 12 +1 +1 Con 11 +0 +0|Int 6 −2 −2 Wis 10 +0 +0 Cha 7 −2 −2|Senses Darkvision 60 ft.; Passive Perception 10|Languages Gnoll|CR 1/2 (XP 100; PB +2)"},
	{"Goat", [2]int{355, 355}, "Medium Beast, Unaligned|AC 10 Initiative +0 (10)|HP 4 (1d8)|Speed 40 ft., Climb 30 ft.|Str 11 +0 +2 Dex 10 +0 +0 Con 11 +0 +0|Int 2 −4 −4 Wis 10 +0 +0 Cha 5 −3 −3|Skills Perception +2|Senses Darkvision 60 ft.; Passive Perception 12|Languages None|CR 0 (XP 10; PB +2)"},
	{"Goblin Boss", [2]int{290, 290}, "Small Fey (Goblinoid), Chaotic Neutral|AC 17 Initiative +2 (12)|HP 21 (6d6)|Speed 30 ft.|Str 10 +0 +0 Dex 15 +2 +2 Con 10 +0 +0|Int 10 +0 +0 Wis 8 −1 −1 Cha 10 +0 +0|Skills Stealth +6|Gear Chain Shirt, Scimitar, Shield, Shortbow|Senses Darkvision 60 ft.; Passive Perception 9|Languages Common, Goblin|CR 1 (XP 200; PB +2)"},
	{"Goblin Minion", [2]int{290, 290}, "Small Fey (Goblinoid), Chaotic Neutral|AC 12 Initiative +2 (12)|HP 7 (2d6)|Speed 30 ft.|Str 8 −1 −1 Dex 15 +2 +2 Con 10 +0 +0|Int 10 +0 +0 Wis 8 −1 −1 Cha 8 −1 −1|Skills Stealth +6|Gear Daggers (3)|Senses Darkvision 60 ft.; Passive Perception 9|Languages Common, Goblin|CR 1/8 (XP 25; PB +2)"},
	{"Goblin Warrior", [2]int{290, 290}, "Small Fey (Goblinoid), Chaotic Neutral|AC 15 Initiative +2 (12)|HP 10 (3d6)|Speed 30 ft.|Str 8 −1 −1 Dex 15 +2 +2 Con 10 +0 +0|Int 10 +0 +0 Wis 8 −1 −1 Cha 8 −1 −1|Skills Stealth +6|Gear Leather Armor, Scimitar, Shield, Shortbow|Senses Darkvision 60 ft.; Passive Perception 9|Languages Common, Goblin|CR 1/4 (XP 50; PB +2)"},
	{"Gold Dragon Wyrmling", [2]int{290, 291}, "Medium Dragon (Metallic), Lawful Good|AC 17 Initiative +4 (14)|HP 60 (8d8 + 24)|Speed 30 ft., Fly 60 ft., Swim 30 ft.|Str 19 +4 +4 Dex 14 +2 +4 Con 17 +3 +3|Int 14 +2 +2 Wis 11 +0 +2 Cha 16 +3 +3|Skills Perception +4, Stealth +4|Immunities Fire|Senses Blindsight 10 ft., Darkvision 60 ft.; Passive Perception 14|Languages Draconic|CR 3 (700 XP; PB +2)"},
	{"Gorgon", [2]int{292, 293}, "Large Construct, Unaligned|AC 19 Initiative +0 (10)|HP 114 (12d10 + 48)|Speed 40 ft.|Str 20 +5 +5 Dex 11 +0 +0 Con 18 +4 +4|Int 2 −4 −4 Wis 12 +1 +1 Cha 7 −2 −2|Skills Perception +7|Immunities Exhaustion, Petrified|Senses Darkvision 60 ft.; Passive Perception 17|Languages None|CR 5 (XP 1,800; PB +3)"},
	{"Gray Ooze", [2]int{293, 293}, "Medium Ooze, Unaligned|AC 9 Initiative −2 (13)|HP 22 (3d8 + 9)|Speed 10 ft., Climb 10 ft.|Str 12 +1 +1 Dex 6 −2 −2 Con 16 +3 +3|Int 1 −5 −5 Wis 6 −2 −2 Cha 2 −4 −4|Skills Stealth +2|Resistances Acid, Cold, Fire|Immunities Blinded, Charmed, Deafened, Exhaustion, Frightened, Grappled, Prone, Restrained|Senses Blindsight 60 ft.; Passive Perception 8|Languages None|CR 1/2 (XP 100; PB +2)"},
	{"Green Dragon Wyrmling", [2]int{293, 293}, "Medium Dragon (Chromatic), Lawful Evil|AC 17 Initiative +3 (13)|HP 38 (7d8 + 7)|Speed 30 ft., Fly 60 ft., Swim 30 ft.|Str 15 +2 +2 Dex 12 +1 +3 Con 13 +1 +1|Int 14 +2 +2 Wis 11 +0 +2 Cha 13 +1 +1|Skills Perception +4, Stealth +3|Immunities Poison; Poisoned|Senses Blindsight 10 ft., Darkvision 60 ft.; Passive Perception 14|Languages Draconic|CR 2 (XP 450; PB +2)"},
	{"Green Hag", [2]int{295, 295}, "Medium Fey, Neutral Evil|AC 17 Initiative +1 (11)|HP 82 (11d8 + 33)|Speed 30 ft., Swim 30 ft.|Str 18 +4 +4 Dex 12 +1 +1 Con 16 +3 +3|Int 13 +1 +1 Wis 14 +2 +2 Cha 14 +2 +2|Skills Arcana +5, Deception +4, Perception +4, Stealth +3|Senses Darkvision 60 ft.; Passive Perception 14|Languages Common, Elvish, Sylvan|CR 3 (XP 700; PB +2)"},
	{"Grick", [2]int{295, 295}, "Medium Aberration, Unaligned|AC 14 Initiative +2 (12)|HP 54 (12d8)|Speed 30 ft., Climb 30 ft.|Str 14 +2 +2 Dex 14 +2 +2 Con 11 +0 +0|Int 3 −4 −4 Wis 14 +2 +2 Cha 5 −3 −3|Skills Stealth +4|Senses Darkvision 60 ft.; Passive Perception 12|Languages None|CR 2 (XP 450; PB +2)"},
	{"Griffon", [2]int{295, 295}, "Large Monstrosity, Unaligned|AC 12 Initiative +2 (12)|HP 59 (7d10 + 21)|Speed 30 ft., Fly 80 ft.|Str 18 +4 +4 Dex 15 +2 +2 Con 16 +3 +3|Int 2 −4 −4 Wis 13 +1 +1 Cha 8 −1 −1|Skills Perception +5|Senses Darkvision 60 ft.; Passive Perception 15|Languages None|CR 2 (XP 450; PB +2)"},
}

// Typed values of every trait, action, bonus action and reaction, reviewed
// against the same pages: limited uses, attack bonus/reach/range/damage, save
// ability and DC, target shape size/width/range, dice and dependency keys. A
// ranged attack printed with one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text: the Goblin Minion's Dagger is a "Melee or Ranged Attack Roll" (key
// action/melee-or-ranged-attack, damage typed on its effect, as for the #153
// Archmage); the Goat's Ram prints a fixed 1 and a charge hit of 2 (1d4), the
// second kept as effect/charge-bludgeoning-damage; the Goblin's extra 2 (1d4)
// damage with Advantage is a second typed effect; the Green Hag's Spellcasting
// prints only an At Will list, so its typed daily groups are empty.
var monsters163Entries = map[string][]string{
	"Gnoll Warrior": {
		"action|rend|Rend; melee +4 reach 5 range 0/0 5 (1d6 + 2) piercing; effect/piercing-damage 5 (1d6 + 2)",
		"action|bone-bow|Bone Bow; ranged +3 reach 0 range 150/600 6 (1d10 + 1) piercing; effect/piercing-damage 6 (1d10 + 1)",
		"bonus|rampage|Rampage (1/Day); limit/per-day, requirement/target-already-bloodied, effect/move-half-speed, action/melee-attack",
	},
	"Goat": {
		"action|ram|Ram; melee +2 reach 5 range 0/0 1 bludgeoning; effect/bludgeoning-damage [damage/fixed], effect/charge-bludgeoning-damage 2 (1d4), requirement/charge-distance",
	},
	"Goblin Boss": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack",
		"action|scimitar|Scimitar; melee +4 reach 5 range 0/0 5 (1d6 + 2) slashing; effect/slashing-damage 5 (1d6 + 2), effect/advantage-slashing-damage 2 (1d4), requirement/attack-roll-advantage",
		"action|shortbow|Shortbow; ranged +4 reach 0 range 80/320 5 (1d6 + 2) piercing; effect/piercing-damage 5 (1d6 + 2), effect/advantage-piercing-damage 2 (1d4), requirement/attack-roll-advantage",
		"bonus|nimble-escape|Nimble Escape; action/disengage, action/hide",
		"reaction|redirect-attack|Redirect Attack; reaction/trigger, requirement/ally-size-small-or-medium, effect/swap-places, effect/redirect-attack",
	},
	"Goblin Minion": {
		"action|dagger|Dagger; effect/piercing-damage 4 (1d4 + 2)",
		"bonus|nimble-escape|Nimble Escape; action/disengage, action/hide",
	},
	"Goblin Warrior": {
		"action|scimitar|Scimitar; melee +4 reach 5 range 0/0 5 (1d6 + 2) slashing; effect/slashing-damage 5 (1d6 + 2), effect/advantage-slashing-damage 2 (1d4), requirement/attack-roll-advantage",
		"action|shortbow|Shortbow; ranged +4 reach 0 range 80/320 5 (1d6 + 2) piercing; effect/piercing-damage 5 (1d6 + 2), effect/advantage-piercing-damage 2 (1d4), requirement/attack-roll-advantage",
		"bonus|nimble-escape|Nimble Escape; action/disengage, action/hide",
	},
	"Gold Dragon Wyrmling": {
		"trait|trait/amphibious|Amphibious; effect/breathe-air-and-water",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +6 reach 5 range 0/0 9 (1d10 + 4) slashing; effect/slashing-damage 9 (1d10 + 4)",
		"action|fire-breath|Fire Breath (Recharge 5–6); dex DC 13; cone 15/0/0; save/dex, area/cone, limit/recharge, effect/fire-damage 22 (4d10), save/half-damage",
		"action|weakening-breath|Weakening Breath; str DC 13; cone 15/0/0; save/str, area/cone, requirement/not-already-affected, effect/strength-d20-disadvantage, effect/damage-roll-penalty 2 (1d4), save/repeat, effect/duration-one-minute",
	},
	"Gorgon": {
		"action|gore|Gore; melee +8 reach 5 range 0/0 18 (2d12 + 5) piercing; effect/piercing-damage 18 (2d12 + 5), requirement/target-size, requirement/charge-distance, condition/prone",
		"action|petrifying-breath|Petrifying Breath (Recharge 5–6); con DC 15; cone 30/0/0; save/con, area/cone, limit/recharge, save/first-failure, condition/restrained, save/repeat, save/second-failure, condition/petrified",
		"bonus|trample|Trample; dex DC 16; creature 0/0/5; save/dex, target/creature, requirement/target-prone, effect/bludgeoning-damage 16 (2d10 + 5), save/half-damage",
	},
	"Gray Ooze": {
		"trait|trait/amorphous|Amorphous; effect/squeeze-through-narrow-space",
		"trait|trait/corrosive-form|Corrosive Form; effect/ammunition-destroyed-on-hit, effect/weapon-attack-penalty, effect/weapon-destroyed-at-penalty-limit, effect/mending-removes-penalty, effect/eat-through-metal-or-wood",
		"action|pseudopod|Pseudopod; melee +3 reach 5 range 0/0 10 (2d8 + 1) acid; effect/acid-damage 10 (2d8 + 1), effect/armor-ac-penalty, effect/armor-destroyed-at-ac-10, effect/mending-removes-penalty",
	},
	"Green Dragon Wyrmling": {
		"trait|trait/amphibious|Amphibious; effect/breathe-air-and-water",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +4 reach 5 range 0/0 7 (1d10 + 2) slashing; effect/slashing-damage 7 (1d10 + 2), effect/poison-damage 3 (1d6)",
		"action|poison-breath|Poison Breath (Recharge 5–6); con DC 11; cone 15/0/0; save/con, area/cone, limit/recharge, effect/poison-damage 21 (6d6), save/half-damage",
	},
	"Green Hag": {
		"trait|trait/amphibious|Amphibious; effect/breathe-air-and-water",
		"trait|trait/coven-magic|Coven Magic; effect/coven-casting, requirement/two-hag-allies-within-30-feet, spellcasting/no-material-components, spellcasting/normal-casting-time, limit/long-rest-per-spell, spell/augury, spell/find-familiar, spell/identify, spell/locate-object, spell/scrying, spell/unseen-servant",
		"trait|trait/mimicry|Mimicry; effect/mimic-sounds, check/wisdom-insight",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|claw|Claw; melee +6 reach 5 range 0/0 8 (1d8 + 4) slashing; effect/slashing-damage 8 (1d8 + 4), effect/poison-damage 3 (1d6)",
		"action|spellcasting|Spellcasting; wis DC 12 +4; at will: Dancing Lights, Disguise Self (24-hour duration), Invisibility (self only, and the hag leaves no tracks while Invisible), Minor Illusion, Ray of Sickness (level 3 version); spellcasting/no-material-components, spellcasting/spell-attack-bonus, spell/dancing-lights [spellcasting/at-will], spell/disguise-self [spellcasting/at-will], spell/invisibility [spellcasting/at-will], spell/minor-illusion [spellcasting/at-will], spell/ray-of-sickness [spellcasting/at-will]",
	},
	"Grick": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|beak|Beak; melee +4 reach 5 range 0/0 9 (2d6 + 2) piercing; effect/piercing-damage 9 (2d6 + 2)",
		"action|tentacles|Tentacles; melee +4 reach 5 range 0/0 7 (1d10 + 2) slashing; effect/slashing-damage 7 (1d10 + 2), requirement/target-size, condition/grappled",
	},
	"Griffon": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +6 reach 5 range 0/0 8 (1d8 + 4) piercing; effect/piercing-damage 8 (1d8 + 4), requirement/target-size, condition/grappled",
	},
}

// monsterDice163 prints a fixed hit, which the source prints as a bare number
// (the Goat's Ram, PDF p. 355), and otherwise the source's "N (XdY + Z)" notation.
func monsterDice163(d MonsterHP) string {
	switch {
	case d.Dice == 0:
		return fmt.Sprint(d.Average)
	case d.Modifier < 0:
		return fmt.Sprintf("%d (%dd%d − %d)", d.Average, d.Dice, d.Sides, -d.Modifier)
	}
	return monsterDice151(d)
}

// printedInitiative163 is the Initiative score the PDF prints where it differs
// from 10 plus the modifier, which the record keeps.
var printedInitiative163 = map[string]int{"Gray Ooze": 13}

// monsterHeader163 prints typed fields in the source's stat-block notation.
func monsterHeader163(m Monster) string {
	score := m.Initiative.Score
	if printed, ok := printedInitiative163[m.Name]; ok {
		score = printed
	}
	cr := fmt.Sprintf("CR %s (XP %s; PB +%d)", m.CR, monsterThousands151(m.XP), m.ProficiencyBonus)
	if m.Name == "Gold Dragon Wyrmling" {
		cr = fmt.Sprintf("CR %s (%s XP; PB +%d)", m.CR, monsterThousands151(m.XP), m.ProficiencyBonus)
	}
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
		fmt.Sprintf("AC %d Initiative %s (%d)", m.AC, monsterSigned151(m.Initiative.Modifier), score),
		"HP " + monsterDice163(m.HP),
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
	lines = append(lines, "Senses "+passive, languages, cr)
	return strings.Join(lines, "|")
}

func monsterEntry163(section, id string, x MonsterMechanic, attack *MonsterAttack, sc *MonsterSpellcasting, effects []MonsterMechanic) string {
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
		parts = append(parts, fmt.Sprintf("%s +%d reach %d range %d/%d %s %s", attack.Kind, attack.Bonus, attack.Reach, attack.Range, attack.LongRange, monsterDice163(attack.Damage), attack.DamageType))
	}
	if x.Save != nil {
		parts = append(parts, fmt.Sprintf("%s DC %d", x.Save.Ability, x.Save.DC))
	}
	if t := x.Target; t != nil {
		parts = append(parts, fmt.Sprintf("%s %d/%d/%d", t.Shape, t.Size, t.Width, t.Range))
	}
	if x.Dice != nil {
		parts = append(parts, monsterDice163(*x.Dice))
	}
	for _, d := range x.MoreDice {
		parts = append(parts, d.Key+" "+monsterDice163(d.Dice))
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
			key += " " + monsterDice163(*e.Dice)
		}
		if len(e.Requires) > 0 {
			key += " [" + strings.Join(e.Requires, ", ") + "]"
		}
		keys = append(keys, key)
	}
	return section + "|" + id + "|" + strings.Join(append(parts, strings.Join(keys, ", ")), "; ")
}

func monsterEntries163(m Monster) []string {
	out := []string{}
	for _, trait := range m.Traits {
		out = append(out, monsterEntry163("trait", trait.Key, trait, nil, nil, nil))
	}
	for _, section := range []struct {
		label   string
		actions []MonsterAction
	}{{"action", m.Actions}, {"bonus", m.BonusActions}, {"reaction", m.Reactions}} {
		for _, a := range section.actions {
			out = append(out, monsterEntry163(section.label, a.ID, MonsterMechanic{Name: a.Name, Uses: a.Uses, Save: a.Save, Target: a.Target}, a.Attack, a.Spellcasting, a.Effects))
		}
	}
	return out
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters163SourceValues(t *testing.T) {
	c := loadMonsters163(t)
	if len(c.Monsters) != 12 || len(monsters163Headers) != 12 || len(monsters163Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #163") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters163Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader163(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries163(m); !slices.Equal(got, monsters163Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters163Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.Hover || m.Telepathy != 0 || m.LairXP != 0 || m.LegendaryActions != nil || m.Legendary != nil {
				t.Fatal("source lists no vulnerabilities, hover, telepathy, lair XP or legendary actions")
			}
		})
		entries += len(monsters163Entries[want.name])
	}
	// Bold-italic entry headings counted on PDF pp. 289–295 and 355.
	if entries != 40 {
		t.Fatalf("reviewed %d entries, want 40", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters163MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 163 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters163(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the p. 355 Goat under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("monsters animals monsters monsters monsters monsters monsters monsters monsters monsters monsters monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #163 cites the page each block starts on, in this order.
	names := []string{"Gnoll Warrior", "Goat", "Goblin Boss", "Goblin Minion", "Goblin Warrior", "Gold Dragon Wyrmling", "Gorgon", "Gray Ooze", "Green Dragon Wyrmling", "Green Hag", "Grick", "Griffon"}
	pages := []int{289, 355, 290, 290, 290, 290, 292, 293, 293, 295, 295, 295}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic that holds for these twelve blocks (PDF pp. 254–256):
// XP and PB follow CR, Hit Dice follow size and Constitution, saves and skills
// add PB to the ability modifier, and attack rolls add PB to the listed ability.
// Initiative adds PB (Gold and Green Dragon Wyrmlings) or nothing.
func TestMonsters163DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"0": 2, "1/8": 2, "1/4": 2, "1/2": 2, "1": 2, "2": 2, "3": 2, "5": 3}
	skillAbility := map[string]string{"arcana": "int", "deception": "cha", "insight": "wis", "perception": "wis", "persuasion": "cha", "religion": "int", "stealth": "dex"}
	hitDie := map[string]int{"small": 6, "medium": 8, "large": 10}
	// Attacks that do not use the better of Strength and Dexterity.
	attackAbility := map[string]string{"Gnoll Warrior/bone-bow": "dex"}
	for _, m := range loadMonsters163(t).Monsters {
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
func TestMonsters163DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters163(t).Monsters {
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
// keys cover only the record #143 already pinned, and the register still reports
// both families unsupported.
func TestMonsters163RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters163(t).Monsters {
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

func TestMonsters163MalformedFixtures(t *testing.T) {
	// Indexes: 0 Gnoll Warrior, 1 Goat, 2 Goblin Boss, 3 Goblin Minion,
	// 4 Goblin Warrior, 5 Gold Dragon Wyrmling, 6 Gorgon, 7 Gray Ooze,
	// 8 Green Dragon Wyrmling, 9 Green Hag, 10 Grick, 11 Griffon.
	tests := []struct {
		name   string
		mutate func(*MonsterCatalogue)
	}{
		{"duplicate ID", func(c *MonsterCatalogue) { c.Monsters[1].ID = c.Monsters[2].ID }},
		{"animal block in another family", func(c *MonsterCatalogue) { c.Monsters[1].ID = "srd-5.2.1:content/spells/goat" }},
		{"missing languages list", func(c *MonsterCatalogue) { c.Monsters[1].Languages = nil }},
		{"fixed damage with sides", func(c *MonsterCatalogue) { c.Monsters[1].Actions[0].Attack.Damage.Sides = 4 }},
		{"fixed damage with modifier", func(c *MonsterCatalogue) { c.Monsters[1].Actions[0].Attack.Damage.Modifier = 1 }},
		{"fixed damage of zero", func(c *MonsterCatalogue) { c.Monsters[1].Actions[0].Attack.Damage.Average = 0 }},
		{"fixed damage without mechanic", func(c *MonsterCatalogue) {
			m := &c.Monsters[1]
			m.Actions[0].Effects[0].Requires = nil
			m.RequiredMechanics = slices.DeleteFunc(m.RequiredMechanics, func(key string) bool { return key == "damage/fixed" })
		}},
		{"fixed effect dice", func(c *MonsterCatalogue) { c.Monsters[1].Actions[0].Effects[0].Dice = &MonsterHP{Average: 1} }},
		{"sense note without sense", func(c *MonsterCatalogue) {
			c.Monsters[0].Senses.Notes = map[string]string{"blindsight": "unimpeded by magical Darkness"}
		}},
		{"blank sense note", func(c *MonsterCatalogue) { c.Monsters[2].Senses.Notes = map[string]string{"darkvision": " "} }},
		{"unlisted sense note", func(c *MonsterCatalogue) {
			c.Monsters[2].Senses.Notes = map[string]string{"darkvision": "impeded by magical Darkness"}
		}},
		{"empty sense notes", func(c *MonsterCatalogue) { c.Monsters[2].Senses.Notes = map[string]string{} }},
		{"melee range", func(c *MonsterCatalogue) { c.Monsters[0].Actions[0].Attack.Range = 5 }},
		{"further dice without first", func(c *MonsterCatalogue) {
			c.Monsters[5].Traits[0].MoreDice = []MonsterKeyedDice{{Key: "effect/future-damage", Dice: MonsterHP{Average: 3, Dice: 1, Sides: 4, Modifier: 1}}}
		}},
		{"XP off CR", func(c *MonsterCatalogue) { c.Monsters[11].XP = 451 }},
		{"save off modifier", func(c *MonsterCatalogue) { c.Monsters[11].Abilities["wis"] = MonsterAbility{13, 1, 10} }},
		{"modifier", func(c *MonsterCatalogue) { c.Monsters[0].Abilities["str"] = MonsterAbility{14, 3, 3} }},
		{"HP average", func(c *MonsterCatalogue) { c.Monsters[1].HP.Average = 5 }},
		{"Initiative score off modifier", func(c *MonsterCatalogue) { c.Monsters[7].Initiative.Score = 13 }},
		{"hover without flight", func(c *MonsterCatalogue) { c.Monsters[0].Hover = true }},
		{"empty bonus actions", func(c *MonsterCatalogue) { c.Monsters[0].BonusActions = []MonsterAction{} }},
		{"empty reactions", func(c *MonsterCatalogue) { c.Monsters[2].Reactions = []MonsterAction{} }},
		{"duplicate action across sections", func(c *MonsterCatalogue) { c.Monsters[0].BonusActions[0].ID = "rend" }},
		{"recharge", func(c *MonsterCatalogue) { c.Monsters[5].Actions[2].Uses.Recharge = 1 }},
		{"daily uses with recharge", func(c *MonsterCatalogue) { c.Monsters[5].Actions[2].Uses.PerDay = 1 }},
		{"daily uses without count", func(c *MonsterCatalogue) { c.Monsters[0].BonusActions[0].Uses.PerDay = 0 }},
		{"save DC off entry text", func(c *MonsterCatalogue) { c.Monsters[5].Actions[2].Save.DC = 14 }},
		{"save ability off entry text", func(c *MonsterCatalogue) { c.Monsters[5].Actions[2].Save.Ability = "con" }},
		{"creature size", func(c *MonsterCatalogue) { c.Monsters[6].BonusActions[0].Target.Size = 5 }},
		{"cone width", func(c *MonsterCatalogue) { c.Monsters[5].Actions[2].Target.Width = 5 }},
		{"effect dice", func(c *MonsterCatalogue) { c.Monsters[5].Actions[2].Effects[3].Dice.Average = 23 }},
		{"unlisted action effect", func(c *MonsterCatalogue) { c.Monsters[6].Actions[0].Effects[1].Key = "effect/future" }},
		{"untyped spellcasting", func(c *MonsterCatalogue) { c.Monsters[9].Actions[2].Spellcasting = nil }},
		{"misplaced spellcasting", func(c *MonsterCatalogue) {
			c.Monsters[9].Actions[1].Spellcasting = c.Monsters[9].Actions[2].Spellcasting
		}},
		{"spell ability", func(c *MonsterCatalogue) { c.Monsters[9].Actions[2].Spellcasting.Ability = "" }},
		{"spell DC", func(c *MonsterCatalogue) { c.Monsters[9].Actions[2].Spellcasting.SaveDC = 0 }},
		{"negative spell attack bonus", func(c *MonsterCatalogue) { c.Monsters[9].Actions[2].Spellcasting.AttackBonus = -1 }},
		{"unlisted spell", func(c *MonsterCatalogue) { c.Monsters[9].Actions[2].Spellcasting.AtWill[0] = "Wish" }},
		{"duplicate spell", func(c *MonsterCatalogue) { c.Monsters[9].Actions[2].Spellcasting.AtWill[0] = "Minor Illusion" }},
		{"missing daily spell groups", func(c *MonsterCatalogue) { c.Monsters[9].Actions[2].Spellcasting.PerDay = nil }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters163(t)
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
	raw, err := os.ReadFile("testdata/monsters_163.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, edit := range [][2]string{
		{`"dice": 0,`, ``},
		{`"sides": 0,`, `"sides": null,`},
		{`"long_range": 0,`, `"long_range": null,`},
		{`"recharge": 5`, `"recharge": null`},
		{`"per_day": 0,`, ``},
		{`"dc": 15`, `"dc": 15, "dc": 16`},
		{`"width": 0,`, ``},
		{`"save_dc": 12,`, ``},
		{`"attack_bonus": 4`, `"attack_bonus": 4, "attack_bonus": 5`},
		{`"darkvision": 60`, `"darkvision": null`},
		{`"shape": "cone",`, `"shape": "cone", "height": 5,`},
		{`"sizes": [`, `"colour": "grey", "sizes": [`},
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
func TestMonsters163ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters163(t)
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
	for _, id := range []string{"beak", "bone-bow", "claw", "dagger", "gore", "pseudopod", "ram", "rend", "shortbow", "tentacles"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
