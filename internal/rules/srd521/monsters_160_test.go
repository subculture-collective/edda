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

func loadMonsters160(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_160.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #160 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Gelatinous Cube's, Ghost's and Giant Centipede's blocks run
// onto the next page.
var monsters160Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Gargoyle", [2]int{286, 286}, "Medium Elemental, Chaotic Evil|AC 15 Initiative +2 (12)|HP 67 (9d8 + 27)|Speed 30 ft., Fly 60 ft.|Str 15 +2 +2 Dex 11 +0 +0 Con 16 +3 +3|Int 6 −2 −2 Wis 11 +0 +0 Cha 7 −2 −2|Skills Stealth +4|Immunities Poison; Exhaustion, Petrified, Poisoned|Senses Darkvision 60 ft.; Passive Perception 10|Languages Primordial (Terran)|CR 2 (XP 450; PB +2)"},
	{"Gelatinous Cube", [2]int{286, 287}, "Large Ooze, Unaligned|AC 6 Initiative −4 (6)|HP 63 (6d10 + 30)|Speed 15 ft.|Str 14 +2 +2 Dex 3 −4 −4 Con 20 +5 +5|Int 1 −5 −5 Wis 6 −2 −2 Cha 1 −5 −5|Immunities Acid; Blinded, Charmed, Deafened, Exhaustion, Frightened, Prone|Senses Blindsight 60 ft.; Passive Perception 8|Languages None|CR 2 (XP 450; PB +2)"},
	{"Ghast", [2]int{287, 287}, "Medium Undead, Chaotic Evil|AC 13 Initiative +3 (13)|HP 36 (8d8)|Speed 30 ft.|Str 16 +3 +3 Dex 17 +3 +3 Con 10 +0 +0|Int 11 +0 +0 Wis 10 +0 +2 Cha 8 −1 −1|Resistances Necrotic|Immunities Poison; Charmed, Exhaustion, Poisoned|Senses Darkvision 60 ft.; Passive Perception 10|Languages Common|CR 2 (XP 450; PB +2)"},
	{"Ghost", [2]int{287, 288}, "Medium Undead, Neutral|AC 11 Initiative +1 (11)|HP 45 (10d8)|Speed 5 ft., Fly 40 ft. (hover)|Str 7 −2 −2 Dex 13 +1 +1 Con 10 +0 +0|Int 10 +0 +0 Wis 12 +1 +1 Cha 17 +3 +3|Resistances Acid, Bludgeoning, Cold, Fire, Lightning, Piercing, Slashing, Thunder|Immunities Necrotic, Poison; Charmed, Exhaustion, Frightened, Grappled, Paralyzed, Petrified, Poisoned, Prone, Restrained|Senses Darkvision 60 ft.; Passive Perception 11|Languages Common plus one other language|CR 4 (XP 1,100; PB +2)"},
	{"Ghoul", [2]int{288, 288}, "Medium Undead, Chaotic Evil|AC 12 Initiative +2 (12)|HP 22 (5d8)|Speed 30 ft.|Str 13 +1 +1 Dex 15 +2 +2 Con 10 +0 +0|Int 7 −2 −2 Wis 10 +0 +0 Cha 6 −2 −2|Immunities Poison; Charmed, Exhaustion, Poisoned|Senses Darkvision 60 ft.; Passive Perception 10|Languages Common|CR 1 (XP 200; PB +2)"},
	{"Giant Ape", [2]int{349, 349}, "Huge Beast, Unaligned|AC 12 Initiative +5 (15)|HP 168 (16d12 + 64)|Speed 40 ft., Climb 40 ft.|Str 23 +6 +6 Dex 14 +2 +2 Con 18 +4 +4|Int 5 −3 −3 Wis 12 +1 +1 Cha 7 −2 −2|Skills Athletics +9, Perception +4, Survival +4|Senses Passive Perception 14|Languages None|CR 7 (XP 2,900; PB +3)"},
	{"Giant Badger", [2]int{349, 349}, "Medium Beast, Unaligned|AC 13 Initiative +0 (10)|HP 15 (2d8 + 6)|Speed 30 ft., Burrow 10 ft.|Str 13 +1 +1 Dex 10 +0 +0 Con 17 +3 +3|Int 2 −4 −4 Wis 12 +1 +1 Cha 5 −3 −3|Skills Perception +3|Resistances Poison|Senses Darkvision 60 ft.; Passive Perception 13|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Giant Bat", [2]int{349, 349}, "Large Beast, Unaligned|AC 13 Initiative +3 (13)|HP 22 (4d10)|Speed 10 ft., Fly 60 ft.|Str 15 +2 +2 Dex 16 +3 +3 Con 11 +0 +0|Int 2 −4 −4 Wis 12 +1 +1 Cha 6 −2 −2|Senses Blindsight 120 ft.; Passive Perception 11|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Giant Boar", [2]int{349, 349}, "Large Beast, Unaligned|AC 13 Initiative +0 (10)|HP 42 (5d10 + 15)|Speed 40 ft.|Str 17 +3 +5 Dex 10 +0 +0 Con 16 +3 +3|Int 2 −4 −4 Wis 7 −2 −2 Cha 5 −3 −3|Senses Passive Perception 8|Languages None|CR 2 (XP 450; PB +2)"},
	{"Giant Centipede", [2]int{349, 350}, "Small Beast, Unaligned|AC 14 Initiative +2 (12)|HP 9 (2d6 + 2)|Speed 30 ft., Climb 30 ft.|Str 5 −3 −3 Dex 14 +2 +2 Con 12 +1 +1|Int 1 −5 −5 Wis 7 −2 −2 Cha 3 −4 −4|Senses Blindsight 30 ft.; Passive Perception 8|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Giant Constrictor Snake", [2]int{350, 350}, "Huge Beast, Unaligned|AC 12 Initiative +2 (12)|HP 60 (8d12 + 8)|Speed 30 ft., Swim 30 ft.|Str 19 +4 +4 Dex 14 +2 +2 Con 12 +1 +1|Int 1 −5 −5 Wis 10 +0 +0 Cha 3 −4 −4|Skills Perception +2|Senses Blindsight 10 ft.; Passive Perception 12|Languages None|CR 2 (XP 450; PB +2)"},
	{"Giant Crab", [2]int{350, 350}, "Medium Beast, Unaligned|AC 15 Initiative +1 (11)|HP 13 (3d8)|Speed 30 ft., Swim 30 ft.|Str 13 +1 +1 Dex 13 +1 +1 Con 11 +0 +0|Int 1 −5 −5 Wis 9 −1 −1 Cha 3 −4 −4|Skills Stealth +3|Senses Blindsight 30 ft.; Passive Perception 9|Languages None|CR 1/8 (XP 25; PB +2)"},
}

// Typed values of every trait, action, bonus action and reaction, reviewed
// against the same pages: limited uses, attack bonus/reach/range/damage, save
// ability and DC, target shape size/width/range, dice and dependency keys. A
// printed "Recharge 6" is typed as recharge 6.
//
// Printed values the contract cannot type stay in the entry's retained source
// text: the Giant Boar's and Giant Crab's conditional rider, the Gelatinous
// Cube's capacity and the Ghost's Possession duration are named by dependency
// keys only, and the Ghost's Etherealness is a spell key, not a Spellcasting action.
var monsters160Entries = map[string][]string{
	"Gargoyle": {
		"trait|trait/flyby|Flyby; effect/no-opportunity-attack-leaving-reach-flying",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|claw|Claw; melee +4 reach 5 range 0/0 7 (2d4 + 2) slashing; effect/slashing-damage 7 (2d4 + 2)",
	},
	"Gelatinous Cube": {
		"trait|trait/ooze-cube|Ooze Cube; 10 (3d6); effect/fills-space, movement/enter-creature-space, save/disadvantage-against-engulf, effect/total-cover, effect/creature-capacity, effect/pull-out, check/str-athletics, effect/acid-damage",
		"trait|trait/transparent|Transparent; effect/transparent, check/wis-perception",
		"action|pseudopod|Pseudopod; melee +4 reach 5 range 0/0 12 (3d6 + 2) acid; effect/acid-damage 12 (3d6 + 2)",
		"action|engulf|Engulf; dex DC 12; effect/move-without-opportunity-attacks, movement/through-creature-spaces, save/dex, target/creatures-entered-spaces, effect/acid-damage 10 (3d6), effect/engulf, effect/suffocating, effect/no-verbal-spells, condition/restrained, effect/start-of-turn-damage 10 (3d6), effect/move-with-cube, effect/escape-engulf, check/str-athletics, effect/escape-success, save/half-damage, effect/success-move, effect/no-space-fails-save",
	},
	"Ghast": {
		"trait|trait/stench|Stench; con DC 10; emanation 5/0/0; save/con, area/emanation, target/creatures-starting-turn-in-area, condition/poisoned, save/success-immunity",
		"action|bite|Bite; melee +5 reach 5 range 0/0 7 (1d8 + 3) piercing; effect/piercing-damage 7 (1d8 + 3), effect/necrotic-damage 9 (2d8)",
		"action|claw|Claw; melee +5 reach 5 range 0/0 10 (2d6 + 3) slashing; con DC 10; effect/slashing-damage 10 (2d6 + 3), requirement/target-non-undead, save/con, condition/paralyzed",
	},
	"Ghost": {
		"trait|trait/ethereal-sight|Ethereal Sight; effect/see-into-ethereal-plane",
		"trait|trait/incorporeal-movement|Incorporeal Movement; 5 (1d10); movement/through-creature-spaces, movement/through-objects, effect/force-damage",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|withering-touch|Withering Touch; melee +5 reach 5 range 0/0 19 (3d10 + 3) necrotic; effect/necrotic-damage 19 (3d10 + 3)",
		"action|etherealness|Etherealness; spell/etherealness, spellcasting/no-components, spellcasting/charisma-ability, effect/visible-on-both-planes, effect/no-cross-plane-effects",
		"action|horrific-visage|Horrific Visage; wis DC 13; cone 60/0/0; save/wis, area/cone, requirement/not-undead, effect/psychic-damage 10 (2d6 + 3), condition/frightened, save/success-immunity",
		"action|possession|Possession (Recharge 6); cha DC 13; creature 0/0/5; limit/recharge, save/cha, target/creature, requirement/target-humanoid, effect/possession, condition/incapacitated, effect/ghost-controls-body, effect/untargetable-except-undead, effect/possessed-statistics, effect/possession-duration, effect/possession-ends, save/success-immunity",
	},
	"Ghoul": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +4 reach 5 range 0/0 5 (1d6 + 2) piercing; effect/piercing-damage 5 (1d6 + 2), effect/necrotic-damage 3 (1d6)",
		"action|claw|Claw; melee +4 reach 5 range 0/0 4 (1d4 + 2) slashing; con DC 10; effect/slashing-damage 4 (1d4 + 2), requirement/target-not-undead-or-elf, save/con, condition/paralyzed",
	},
	"Giant Ape": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|fist|Fist; melee +9 reach 10 range 0/0 22 (3d10 + 6) bludgeoning; effect/bludgeoning-damage 22 (3d10 + 6)",
		"action|boulder-toss|Boulder Toss (Recharge 6); dex DC 17; sphere 5/0/90; limit/recharge, effect/hurl-boulder, save/dex, area/sphere, effect/bludgeoning-damage 24 (7d6), requirement/target-size, condition/prone, save/half-damage-only",
		"bonus|leap|Leap; effect/jump",
	},
	"Giant Badger": {
		"action|bite|Bite; melee +3 reach 5 range 0/0 6 (2d4 + 1) piercing; effect/piercing-damage 6 (2d4 + 1)",
	},
	"Giant Bat": {
		"action|bite|Bite; melee +5 reach 5 range 0/0 6 (1d6 + 3) piercing; effect/piercing-damage 6 (1d6 + 3)",
	},
	"Giant Boar": {
		"trait|trait/bloodied-fury|Bloodied Fury; requirement/bloodied, attack/advantage",
		"action|gore|Gore; melee +5 reach 5 range 0/0 10 (2d6 + 3) piercing; effect/piercing-damage 10 (2d6 + 3), requirement/target-size, requirement/charge-distance, effect/extra-piercing-damage 7 (2d6), condition/prone",
	},
	"Giant Centipede": {
		"action|bite|Bite; melee +4 reach 5 range 0/0 4 (1d4 + 2) piercing; effect/piercing-damage 4 (1d4 + 2), condition/poisoned",
	},
	"Giant Constrictor Snake": {
		"action|multiattack|Multiattack; action/melee-attack, multiattack/additional-action",
		"action|bite|Bite; melee +6 reach 10 range 0/0 11 (2d6 + 4) piercing; effect/piercing-damage 11 (2d6 + 4)",
		"action|constrict|Constrict; str DC 14; creature 0/0/10; save/str, target/creature, requirement/target-size, effect/bludgeoning-damage 13 (2d8 + 4), condition/grappled",
	},
	"Giant Crab": {
		"trait|trait/amphibious|Amphibious; effect/breathe-air-and-water",
		"action|claw|Claw; melee +3 reach 5 range 0/0 4 (1d6 + 1) bludgeoning; effect/bludgeoning-damage 4 (1d6 + 1), requirement/target-size, condition/grappled, limit/grapple-from-one-of-two-claws",
	},
}

// monsterDice160 prints a fixed hit, which the source prints as a bare number
// (the Cat's Scratch), and otherwise the source's "N (XdY + Z)" notation.
func monsterDice160(d MonsterHP) string {
	switch {
	case d.Dice == 0:
		return fmt.Sprint(d.Average)
	case d.Modifier < 0:
		return fmt.Sprintf("%d (%dd%d − %d)", d.Average, d.Dice, d.Sides, -d.Modifier)
	}
	return monsterDice151(d)
}

// monsterHeader160 prints typed fields in the source's stat-block notation.
func monsterHeader160(m Monster) string {
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
		"HP " + monsterDice160(m.HP),
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

func monsterEntry160(section, id string, x MonsterMechanic, attack *MonsterAttack, sc *MonsterSpellcasting, effects []MonsterMechanic) string {
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
		parts = append(parts, fmt.Sprintf("%s +%d reach %d range %d/%d %s %s", attack.Kind, attack.Bonus, attack.Reach, attack.Range, attack.LongRange, monsterDice160(attack.Damage), attack.DamageType))
	}
	if x.Save != nil {
		parts = append(parts, fmt.Sprintf("%s DC %d", x.Save.Ability, x.Save.DC))
	}
	if t := x.Target; t != nil {
		parts = append(parts, fmt.Sprintf("%s %d/%d/%d", t.Shape, t.Size, t.Width, t.Range))
	}
	if x.Dice != nil {
		parts = append(parts, monsterDice160(*x.Dice))
	}
	for _, d := range x.MoreDice {
		parts = append(parts, d.Key+" "+monsterDice160(d.Dice))
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
			key += " " + monsterDice160(*e.Dice)
		}
		if len(e.Requires) > 0 {
			key += " [" + strings.Join(e.Requires, ", ") + "]"
		}
		keys = append(keys, key)
	}
	return section + "|" + id + "|" + strings.Join(append(parts, strings.Join(keys, ", ")), "; ")
}

func monsterEntries160(m Monster) []string {
	out := []string{}
	for _, trait := range m.Traits {
		out = append(out, monsterEntry160("trait", trait.Key, trait, nil, nil, nil))
	}
	for _, section := range []struct {
		label   string
		actions []MonsterAction
	}{{"action", m.Actions}, {"bonus", m.BonusActions}, {"reaction", m.Reactions}} {
		for _, a := range section.actions {
			out = append(out, monsterEntry160(section.label, a.ID, MonsterMechanic{Name: a.Name, Uses: a.Uses, Save: a.Save, Target: a.Target}, a.Attack, a.Spellcasting, a.Effects))
		}
	}
	return out
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters160SourceValues(t *testing.T) {
	c := loadMonsters160(t)
	if len(c.Monsters) != 12 || len(monsters160Headers) != 12 || len(monsters160Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #160") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters160Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader160(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries160(m); !slices.Equal(got, monsters160Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters160Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.Hover != (m.Name == "Ghost") || m.LairXP != 0 || m.LegendaryActions != nil || m.Legendary != nil {
				t.Fatal("source lists no vulnerabilities, lair XP or legendary actions, and only the Ghost hovers")
			}
		})
		entries += len(monsters160Entries[want.name])
	}
	// Bold-italic entry headings counted on PDF pp. 286–288 and 349–350.
	if entries != 34 {
		t.Fatalf("reviewed %d entries, want 34", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters160MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 160 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters160(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the five p. 286-288 blocks under monsters and the seven p. 349-350 animals under animals.
	if !slices.Equal(families, strings.Fields("monsters monsters monsters monsters monsters animals animals animals animals animals animals animals")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #160 cites the page each block starts on, in this order.
	names := []string{"Gargoyle", "Gelatinous Cube", "Ghast", "Ghost", "Ghoul", "Giant Ape", "Giant Badger", "Giant Bat", "Giant Boar", "Giant Centipede", "Giant Constrictor Snake", "Giant Crab"}
	pages := []int{286, 286, 287, 287, 288, 349, 349, 349, 349, 349, 350, 350}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic that holds for these twelve blocks (PDF pp. 254–256):
// XP and PB follow CR, Hit Dice follow size and Constitution, saves and skills
// add PB to the ability modifier, and attack rolls add PB to the listed ability.
// Initiative adds PB (Gargoyle, Giant Ape) or nothing.
func TestMonsters160DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"1/8": 2, "1/4": 2, "1": 2, "2": 2, "4": 2, "7": 3}
	skillAbility := map[string]string{"athletics": "str", "perception": "wis", "stealth": "dex", "survival": "wis"}
	hitDie := map[string]int{"small": 6, "medium": 8, "large": 10, "huge": 12}
	// Attacks that do not use the better of Strength and Dexterity.
	attackAbility := map[string]string{"Ghost/withering-touch": "cha"}
	for _, m := range loadMonsters160(t).Monsters {
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
func TestMonsters160DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters160(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			if m.RequiredMechanics[0] != "monster/stat-block" {
				t.Fatal("stat-block dependency must come first")
			}
			quoted := func(parent string, x MonsterMechanic) {
				if !strings.Contains(parent, strings.TrimSuffix(x.Text, ".")) {
					t.Fatalf("%q is not quoted from its entry", x.Text)
				}
				if x.Dice != nil && !strings.Contains(x.Text, monsterDice160(*x.Dice)) {
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
func TestMonsters160RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters160(t).Monsters {
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

func TestMonsters160MalformedFixtures(t *testing.T) {
	// Indexes: 0 Gargoyle, 1 Gelatinous Cube, 2 Ghast, 3 Ghost, 4 Ghoul, 5 Giant Ape,
	// 6 Giant Badger, 7 Giant Bat, 8 Giant Boar, 9 Giant Centipede,
	// 10 Giant Constrictor Snake, 11 Giant Crab.
	tests := []struct {
		name   string
		mutate func(*MonsterCatalogue)
	}{
		{"duplicate ID", func(c *MonsterCatalogue) { c.Monsters[1].ID = c.Monsters[2].ID }},
		{"animal block in another family", func(c *MonsterCatalogue) { c.Monsters[10].ID = "srd-5.2.1:content/spells/giant-constrictor-snake" }},
		{"missing languages list", func(c *MonsterCatalogue) { c.Monsters[10].Languages = nil }},
		{"sense note without sense", func(c *MonsterCatalogue) {
			c.Monsters[0].Senses.Notes = map[string]string{"blindsight": "unimpeded by magical Darkness"}
		}},
		{"blank sense note", func(c *MonsterCatalogue) { c.Monsters[0].Senses.Notes = map[string]string{"darkvision": " "} }},
		{"unlisted sense note", func(c *MonsterCatalogue) {
			c.Monsters[0].Senses.Notes = map[string]string{"darkvision": "impeded by magical Darkness"}
		}},
		{"empty sense notes", func(c *MonsterCatalogue) { c.Monsters[0].Senses.Notes = map[string]string{} }},
		{"melee range", func(c *MonsterCatalogue) { c.Monsters[0].Actions[1].Attack.Range = 5 }},
		{"further dice without first", func(c *MonsterCatalogue) {
			c.Monsters[0].Traits[0].MoreDice = []MonsterKeyedDice{{Key: "effect/future-damage", Dice: MonsterHP{Average: 3, Dice: 1, Sides: 4, Modifier: 1}}}
		}},
		{"XP off CR", func(c *MonsterCatalogue) { c.Monsters[3].XP = 1101 }},
		{"save off modifier", func(c *MonsterCatalogue) { c.Monsters[3].Abilities["wis"] = MonsterAbility{20, 5, 10} }},
		{"modifier", func(c *MonsterCatalogue) { c.Monsters[0].Abilities["str"] = MonsterAbility{15, 3, 3} }},
		{"HP average", func(c *MonsterCatalogue) { c.Monsters[11].HP.Average = 14 }},
		{"hover without flight", func(c *MonsterCatalogue) { c.Monsters[3].Speeds["fly"] = 0 }},
		{"empty bonus actions", func(c *MonsterCatalogue) { c.Monsters[5].BonusActions = []MonsterAction{} }},
		{"empty reactions", func(c *MonsterCatalogue) { c.Monsters[0].Reactions = []MonsterAction{} }},
		{"duplicate action across sections", func(c *MonsterCatalogue) { c.Monsters[5].BonusActions[0].ID = "fist" }},
		{"recharge", func(c *MonsterCatalogue) { c.Monsters[5].Actions[2].Uses.Recharge = 1 }},
		{"daily uses with recharge", func(c *MonsterCatalogue) { c.Monsters[5].Actions[2].Uses.PerDay = 1 }},
		{"trait save DC off entry text", func(c *MonsterCatalogue) { c.Monsters[2].Traits[0].Save.DC = 11 }},
		{"save DC off entry text", func(c *MonsterCatalogue) { c.Monsters[10].Actions[2].Save.DC = 16 }},
		{"save ability off entry text", func(c *MonsterCatalogue) { c.Monsters[10].Actions[2].Save.Ability = "dex" }},
		{"creature size", func(c *MonsterCatalogue) { c.Monsters[10].Actions[2].Target.Size = 5 }},
		{"cone width", func(c *MonsterCatalogue) { c.Monsters[3].Actions[3].Target.Width = 5 }},
		{"trait dice", func(c *MonsterCatalogue) { c.Monsters[1].Traits[0].Dice.Average = 9 }},
		{"effect dice", func(c *MonsterCatalogue) { c.Monsters[1].Actions[1].Effects[4].Dice.Average = 9 }},
		{"unlisted action effect", func(c *MonsterCatalogue) { c.Monsters[2].Actions[1].Effects[1].Key = "effect/future" }},
		{"untyped spellcasting on a spell-casting entry", func(c *MonsterCatalogue) {
			c.Monsters[3].Actions[2].Spellcasting = &MonsterSpellcasting{Ability: "cha", SaveDC: 13, AtWill: []string{"Etherealness"}, PerDay: []MonsterSpellUses{}}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters160(t)
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
	raw, err := os.ReadFile("testdata/monsters_160.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, edit := range [][2]string{
		{`"long_range": 0,`, `"long_range": null,`},
		{`"range": 0,`, ``},
		{`"recharge": 6`, `"recharge": null`},
		{`"per_day": 0,`, ``},
		{`"lair_per_day": 0,`, ``},
		{`"dc": 13`, `"dc": 13, "dc": 14`},
		{`"width": 0,`, ``},
		{`"special": {}`, `"special": {"darkvision": null}`},
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
// to 07 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters160ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters160(t)
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
	for _, id := range []string{"bite", "claw", "fist", "gore", "pseudopod", "withering-touch"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
