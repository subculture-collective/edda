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

func loadMonsters162(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_162.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #162 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Giant Spider's block starts on p. 353 and the Giant Wasp's on
// p. 354; each runs onto the next page. The Giant Wolf Spider's Senses line
// wraps in the source.
var monsters162Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Giant Seahorse", [2]int{353, 353}, "Large Beast, Unaligned|AC 14 Initiative +1 (11)|HP 16 (3d10)|Speed 5 ft., Swim 40 ft.|Str 15 +2 +2 Dex 12 +1 +1 Con 11 +0 +0|Int 2 −4 −4 Wis 12 +1 +1 Cha 5 −3 −3|Senses Passive Perception 11|Languages None|CR 1/2 (XP 100; PB +2)"},
	{"Giant Shark", [2]int{353, 353}, "Huge Beast, Unaligned|AC 13 Initiative +3 (13)|HP 92 (8d12 + 40)|Speed 5 ft., Swim 60 ft.|Str 23 +6 +6 Dex 11 +0 +0 Con 21 +5 +5|Int 1 −5 −5 Wis 10 +0 +0 Cha 5 −3 −3|Skills Perception +3|Senses Blindsight 60 ft.; Passive Perception 13|Languages None|CR 5 (XP 1,800; PB +3)"},
	{"Giant Spider", [2]int{353, 354}, "Large Beast, Unaligned|AC 14 Initiative +3 (13)|HP 26 (4d10 + 4)|Speed 30 ft., Climb 30 ft.|Str 14 +2 +2 Dex 16 +3 +3 Con 12 +1 +1|Int 2 −4 −4 Wis 11 +0 +0 Cha 4 −3 −3|Skills Perception +4, Stealth +7|Senses Darkvision 60 ft.; Passive Perception 14|Languages None|CR 1 (XP 200; PB +2)"},
	{"Giant Toad", [2]int{354, 354}, "Large Beast, Unaligned|AC 11 Initiative +1 (11)|HP 39 (6d10 + 6)|Speed 30 ft., Swim 30 ft.|Str 15 +2 +2 Dex 13 +1 +1 Con 13 +1 +1|Int 2 −4 −4 Wis 10 +0 +0 Cha 3 −4 −4|Senses Darkvision 60 ft.; Passive Perception 10|Languages None|CR 1 (XP 200; PB +2)"},
	{"Giant Venomous Snake", [2]int{354, 354}, "Medium Beast, Unaligned|AC 14 Initiative +4 (14)|HP 11 (2d8 + 2)|Speed 40 ft., Swim 40 ft.|Str 10 +0 +0 Dex 18 +4 +4 Con 13 +1 +1|Int 2 −4 −4 Wis 10 +0 +0 Cha 3 −4 −4|Skills Perception +2|Senses Blindsight 10 ft.; Passive Perception 12|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Giant Vulture", [2]int{354, 354}, "Large Monstrosity, Neutral Evil|AC 10 Initiative +0 (10)|HP 25 (3d10 + 9)|Speed 10 ft., Fly 60 ft.|Str 15 +2 +2 Dex 10 +0 +0 Con 16 +3 +3|Int 6 −2 −2 Wis 12 +1 +1 Cha 7 −2 −2|Skills Perception +3|Resistances Necrotic|Senses Darkvision 60 ft.; Passive Perception 13|Languages Understands Common but can’t speak|CR 1 (XP 200; PB +2)"},
	{"Giant Wasp", [2]int{354, 355}, "Medium Beast, Unaligned|AC 13 Initiative +2 (12)|HP 22 (5d8)|Speed 10 ft., Fly 50 ft.|Str 10 +0 +0 Dex 14 +2 +2 Con 10 +0 +0|Int 1 −5 −5 Wis 10 +0 +0 Cha 3 −4 −4|Senses Passive Perception 10|Languages None|CR 1/2 (XP 100; PB +2)"},
	{"Giant Weasel", [2]int{355, 355}, "Medium Beast, Unaligned|AC 13 Initiative +3 (13)|HP 9 (2d8)|Speed 40 ft., Climb 30 ft.|Str 11 +0 +0 Dex 17 +3 +3 Con 10 +0 +0|Int 4 −3 −3 Wis 12 +1 +1 Cha 5 −3 −3|Skills Acrobatics +5, Perception +3, Stealth +5|Senses Darkvision 60 ft.; Passive Perception 13|Languages None|CR 1/8 (XP 25; PB +2)"},
	{"Giant Wolf Spider", [2]int{355, 355}, "Medium Beast, Unaligned|AC 13 Initiative +3 (13)|HP 11 (2d8 + 2)|Speed 40 ft., Climb 40 ft.|Str 12 +1 +1 Dex 16 +3 +3 Con 13 +1 +1|Int 3 −4 −4 Wis 12 +1 +1 Cha 4 −3 −3|Skills Perception +3, Stealth +7|Senses Blindsight 10 ft., Darkvision 60 ft.; Passive Perception 13|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Gibbering Mouther", [2]int{288, 288}, "Medium Aberration, Chaotic Neutral|AC 9 Initiative −1 (9)|HP 52 (7d8 + 21)|Speed 20 ft., Swim 20 ft.|Str 10 +0 +0 Dex 8 −1 −1 Con 16 +3 +3|Int 3 −4 −4 Wis 10 +0 +0 Cha 6 −2 −2|Immunities Prone|Senses Darkvision 60 ft.; Passive Perception 10|Languages None|CR 2 (XP 450; PB +2)"},
	{"Glabrezu", [2]int{289, 289}, "Large Fiend (Demon), Chaotic Evil|AC 17 Initiative +6 (16)|HP 189 (18d10 + 90)|Speed 40 ft.|Str 20 +5 +9 Dex 15 +2 +2 Con 21 +5 +9|Int 19 +4 +4 Wis 17 +3 +7 Cha 16 +3 +7|Skills Deception +7, Perception +7|Resistances Cold, Fire, Lightning|Immunities Poison; Poisoned|Senses Truesight 120 ft.; Passive Perception 17|Languages Abyssal; telepathy 120 ft.|CR 9 (XP 5,000; PB +4)"},
	{"Gladiator", [2]int{289, 289}, "Medium or Small Humanoid, Neutral|AC 16 Initiative +5 (15)|HP 112 (15d8 + 45)|Speed 30 ft.|Str 18 +4 +7 Dex 15 +2 +5 Con 16 +3 +6|Int 10 +0 +0 Wis 12 +1 +4 Cha 15 +2 +2|Skills Athletics +10, Performance +5|Gear Shield, Spears (3), Studded Leather Armor|Senses Passive Perception 11|Languages Common|CR 5 (XP 1,800; PB +3)"},
}

// Typed values of every trait, action, bonus action and reaction, reviewed
// against the same pages: limited uses, attack bonus/reach/range/damage, save
// ability and DC, target shape size/width/range, dice and dependency keys. A
// ranged attack printed with one range repeats it.
//
// Printed values the contract cannot type stay in the entry's retained source
// text: the Gladiator's Spear is a "Melee or Ranged Attack Roll" (key
// action/melee-or-ranged-attack, damage typed on its effect, as for the #153
// Archmage); the Giant Shark's "(with Advantage ...)" after the attack bonus,
// escape DCs, the Giant Toad's swallow rules and the Giant Spider's web
// statistics are keyed effects that quote their entry.
var monsters162Entries = map[string][]string{
	"Giant Seahorse": {
		"trait|trait/water-breathing|Water Breathing; effect/breathe-only-underwater",
		"action|ram|Ram; melee +4 reach 5 range 0/0 9 (2d6 + 2) bludgeoning; effect/bludgeoning-damage 9 (2d6 + 2), effect/charge-bludgeoning-damage 11 (2d8 + 2), requirement/charge-distance",
		"bonus|bubble-dash|Bubble Dash; requirement/underwater, effect/move-half-swim-speed",
	},
	"Giant Shark": {
		"trait|trait/water-breathing|Water Breathing; effect/breathe-only-underwater",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +9 reach 5 range 0/0 22 (3d10 + 6) piercing; effect/piercing-damage 22 (3d10 + 6), attack/advantage-against-target-missing-hit-points",
	},
	"Giant Spider": {
		"trait|trait/spider-climb|Spider Climb; movement/climb-difficult-surfaces, check/no-ability-check-to-climb",
		"trait|trait/web-walker|Web Walker; movement/ignore-web-restrictions, effect/know-web-contact-creature-locations",
		"action|bite|Bite; melee +5 reach 5 range 0/0 7 (1d8 + 3) piercing; effect/piercing-damage 7 (1d8 + 3), effect/poison-damage 7 (2d6)",
		"action|web|Web (Recharge 5–6); dex DC 13; creature 0/0/60; save/dex, target/creature, limit/recharge, condition/restrained, effect/web-armor-class, effect/web-hit-points, effect/web-vulnerability, effect/web-immunity",
	},
	"Giant Toad": {
		"trait|trait/amphibious|Amphibious; effect/breathe-air-and-water",
		"trait|trait/standing-leap|Standing Leap; movement/long-jump-20-feet, movement/high-jump-10-feet, movement/jump-without-running-start",
		"action|bite|Bite; melee +4 reach 5 range 0/0 5 (1d6 + 2) piercing; effect/piercing-damage 5 (1d6 + 2), effect/poison-damage 5 (2d4), requirement/target-size, condition/grappled",
		"action|swallow|Swallow; effect/swallow, requirement/target-size, requirement/grappled-creature, condition/blinded, condition/restrained, effect/total-cover, effect/acid-damage 10 (3d6), limit/one-swallowed-creature, effect/no-bite-while-swallowing, effect/escape-corpse, condition/prone",
	},
	"Giant Venomous Snake": {
		"action|bite|Bite; melee +6 reach 10 range 0/0 6 (1d4 + 4) piercing; effect/piercing-damage 6 (1d4 + 4), effect/poison-damage 4 (1d8)",
	},
	"Giant Vulture": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage, requirement/ally-within-5-feet, requirement/ally-not-incapacitated",
		"action|gouge|Gouge; melee +4 reach 5 range 0/0 9 (2d6 + 2) piercing; effect/piercing-damage 9 (2d6 + 2), condition/poisoned",
	},
	"Giant Wasp": {
		"trait|trait/flyby|Flyby; effect/no-opportunity-attack-flying-out-of-reach",
		"action|sting|Sting; melee +4 reach 5 range 0/0 5 (1d6 + 2) piercing; effect/piercing-damage 5 (1d6 + 2), effect/poison-damage 5 (2d4)",
	},
	"Giant Weasel": {
		"action|bite|Bite; melee +5 reach 5 range 0/0 5 (1d4 + 3) piercing; effect/piercing-damage 5 (1d4 + 3)",
	},
	"Giant Wolf Spider": {
		"trait|trait/spider-climb|Spider Climb; movement/climb-difficult-surfaces, check/no-ability-check-to-climb",
		"action|bite|Bite; melee +5 reach 5 range 0/0 5 (1d4 + 3) piercing; effect/piercing-damage 5 (1d4 + 3), effect/poison-damage 5 (2d4)",
	},
	"Gibbering Mouther": {
		"trait|trait/aberrant-ground|Aberrant Ground; area/emanation, effect/difficult-terrain",
		"trait|trait/gibbering|Gibbering; wis DC 10; save/wis, condition/incapacitated, effect/random-action-roll, effect/do-nothing, effect/random-direction-movement, effect/random-melee-attack",
		"action|bite|Bite; melee +2 reach 5 range 0/0 7 (2d6) piercing; effect/piercing-damage 7 (2d6), requirement/target-size, condition/prone, effect/die-at-zero-hit-points, effect/absorb-body",
		"action|blinding-spittle|Blinding Spittle (Recharge 5–6); dex DC 10; sphere 10/0/30; save/dex, area/sphere, limit/recharge, effect/radiant-damage 7 (2d6), condition/blinded",
	},
	"Glabrezu": {
		"trait|trait/demonic-restoration|Demonic Restoration; effect/revive",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/use-pummel-or-spellcasting",
		"action|pincer|Pincer; melee +9 reach 10 range 0/0 16 (2d10 + 5) slashing; effect/slashing-damage 16 (2d10 + 5), requirement/target-size, condition/grappled",
		"action|pummel|Pummel; dex DC 17; creature 0/0/0; save/dex, target/creature, requirement/grappled-creature, effect/bludgeoning-damage 15 (3d6 + 5), save/half-damage",
		"action|spellcasting|Spellcasting; int DC 16 +0; at will: Darkness, Detect Magic, Dispel Magic; 1/day: Confusion, Fly, Power Word Stun; spellcasting/no-material-components, spell/darkness [spellcasting/at-will], spell/detect-magic [spellcasting/at-will], spell/dispel-magic [spellcasting/at-will], spell/confusion [limit/per-day], spell/fly [limit/per-day], spell/power-word-stun [limit/per-day]",
	},
	"Gladiator": {
		"action|multiattack|Multiattack; action/melee-or-ranged-attack, multiattack/replace-attack",
		"action|spear|Spear; effect/piercing-damage 11 (2d6 + 4)",
		"action|shield-bash|Shield Bash; str DC 15; creature 0/0/5; save/str, target/creature, effect/bludgeoning-damage 9 (2d4 + 4), requirement/target-size, condition/prone",
		"reaction|parry|Parry; reaction/trigger, requirement/holding-weapon, effect/ac-bonus",
	},
}

// monsterDice162 prints a fixed hit, which the source prints as a bare number
// (the Cat's Scratch), and otherwise the source's "N (XdY + Z)" notation.
func monsterDice162(d MonsterHP) string {
	switch {
	case d.Dice == 0:
		return fmt.Sprint(d.Average)
	case d.Modifier < 0:
		return fmt.Sprintf("%d (%dd%d − %d)", d.Average, d.Dice, d.Sides, -d.Modifier)
	}
	return monsterDice151(d)
}

// monsterHeader162 prints typed fields in the source's stat-block notation.
func monsterHeader162(m Monster) string {
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
		"HP " + monsterDice162(m.HP),
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

func monsterEntry162(section, id string, x MonsterMechanic, attack *MonsterAttack, sc *MonsterSpellcasting, effects []MonsterMechanic) string {
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
		parts = append(parts, fmt.Sprintf("%s +%d reach %d range %d/%d %s %s", attack.Kind, attack.Bonus, attack.Reach, attack.Range, attack.LongRange, monsterDice162(attack.Damage), attack.DamageType))
	}
	if x.Save != nil {
		parts = append(parts, fmt.Sprintf("%s DC %d", x.Save.Ability, x.Save.DC))
	}
	if t := x.Target; t != nil {
		parts = append(parts, fmt.Sprintf("%s %d/%d/%d", t.Shape, t.Size, t.Width, t.Range))
	}
	if x.Dice != nil {
		parts = append(parts, monsterDice162(*x.Dice))
	}
	for _, d := range x.MoreDice {
		parts = append(parts, d.Key+" "+monsterDice162(d.Dice))
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
			key += " " + monsterDice162(*e.Dice)
		}
		if len(e.Requires) > 0 {
			key += " [" + strings.Join(e.Requires, ", ") + "]"
		}
		keys = append(keys, key)
	}
	return section + "|" + id + "|" + strings.Join(append(parts, strings.Join(keys, ", ")), "; ")
}

func monsterEntries162(m Monster) []string {
	out := []string{}
	for _, trait := range m.Traits {
		out = append(out, monsterEntry162("trait", trait.Key, trait, nil, nil, nil))
	}
	for _, section := range []struct {
		label   string
		actions []MonsterAction
	}{{"action", m.Actions}, {"bonus", m.BonusActions}, {"reaction", m.Reactions}} {
		for _, a := range section.actions {
			out = append(out, monsterEntry162(section.label, a.ID, MonsterMechanic{Name: a.Name, Uses: a.Uses, Save: a.Save, Target: a.Target}, a.Attack, a.Spellcasting, a.Effects))
		}
	}
	return out
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters162SourceValues(t *testing.T) {
	c := loadMonsters162(t)
	if len(c.Monsters) != 12 || len(monsters162Headers) != 12 || len(monsters162Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #162") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters162Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader162(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries162(m); !slices.Equal(got, monsters162Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters162Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.Hover || m.LairXP != 0 || m.LegendaryActions != nil || m.Legendary != nil {
				t.Fatal("source lists no vulnerabilities, lair XP or legendary actions, and no block hovers")
			}
		})
		entries += len(monsters162Entries[want.name])
	}
	// Bold-italic entry headings counted on PDF pp. 288–289 and 353–355.
	if entries != 36 {
		t.Fatalf("reviewed %d entries, want 36", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters162MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 162 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters162(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the nine Giant animals under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("animals animals animals animals animals animals animals animals animals monsters monsters monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #162 cites the page each block starts on, in this order.
	names := []string{"Giant Seahorse", "Giant Shark", "Giant Spider", "Giant Toad", "Giant Venomous Snake", "Giant Vulture", "Giant Wasp", "Giant Weasel", "Giant Wolf Spider", "Gibbering Mouther", "Glabrezu", "Gladiator"}
	pages := []int{353, 353, 353, 354, 354, 354, 354, 355, 355, 288, 289, 289}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic that holds for these twelve blocks (PDF pp. 254–256):
// XP and PB follow CR, Hit Dice follow size and Constitution, saves and skills
// add PB to the ability modifier, and attack rolls add PB to the listed ability.
// Initiative adds PB (Giant Shark, Glabrezu, Gladiator) or nothing.
func TestMonsters162DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"1/8": 2, "1/4": 2, "1/2": 2, "1": 2, "2": 2, "5": 3, "9": 4}
	skillAbility := map[string]string{"acrobatics": "dex", "athletics": "str", "deception": "cha", "perception": "wis", "performance": "cha", "stealth": "dex"}
	hitDie := map[string]int{"small": 6, "medium": 8, "large": 10, "huge": 12}
	for _, m := range loadMonsters162(t).Monsters {
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
			// The Gladiator's Spear has no typed attack, so its printed bonus is read from the text.
			for _, a := range m.allActions() {
				if a.Key == "action/melee-or-ranged-attack" && !strings.Contains(a.Text, fmt.Sprintf("Attack Roll: +%d, ", max(m.Abilities["str"].Modifier, m.Abilities["dex"].Modifier)+p)) {
					t.Fatalf("%s bonus does not use the ability modifier and PB", a.Name)
				}
			}
			for _, a := range m.allActions() {
				if a.Attack == nil {
					continue
				}
				mod := max(m.Abilities["str"].Modifier, m.Abilities["dex"].Modifier)
				if a.Attack.Bonus != mod+p || a.Attack.Damage.Modifier != mod {
					t.Fatalf("%s attack does not use the ability modifier and PB", a.Name)
				}
			}
		})
	}
}

// Every listed key is used and required, every effect quotes its own entry, and
// every dice expression is typed. Removing any one key invalidates the record.
func TestMonsters162DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters162(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			if m.RequiredMechanics[0] != "monster/stat-block" {
				t.Fatal("stat-block dependency must come first")
			}
			quoted := func(parent string, x MonsterMechanic) {
				if !strings.Contains(parent, strings.TrimSuffix(x.Text, ".")) {
					t.Fatalf("%q is not quoted from its entry", x.Text)
				}
				if x.Dice != nil && !strings.Contains(x.Text, monsterDice162(*x.Dice)) {
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
func TestMonsters162RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters162(t).Monsters {
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

func TestMonsters162MalformedFixtures(t *testing.T) {
	// Indexes: 0 Giant Seahorse, 1 Giant Shark, 2 Giant Spider, 3 Giant Toad,
	// 4 Giant Venomous Snake, 5 Giant Vulture, 6 Giant Wasp, 7 Giant Weasel,
	// 8 Giant Wolf Spider, 9 Gibbering Mouther, 10 Glabrezu, 11 Gladiator.
	tests := []struct {
		name   string
		mutate func(*MonsterCatalogue)
	}{
		{"duplicate ID", func(c *MonsterCatalogue) { c.Monsters[1].ID = c.Monsters[2].ID }},
		{"animal block in another family", func(c *MonsterCatalogue) { c.Monsters[0].ID = "srd-5.2.1:content/spells/giant-seahorse" }},
		{"missing languages list", func(c *MonsterCatalogue) { c.Monsters[0].Languages = nil }},
		{"sense note without sense", func(c *MonsterCatalogue) {
			c.Monsters[6].Senses.Notes = map[string]string{"blindsight": "unimpeded by magical Darkness"}
		}},
		{"blank sense note", func(c *MonsterCatalogue) { c.Monsters[2].Senses.Notes = map[string]string{"darkvision": " "} }},
		{"unlisted sense note", func(c *MonsterCatalogue) {
			c.Monsters[2].Senses.Notes = map[string]string{"darkvision": "impeded by magical Darkness"}
		}},
		{"empty sense notes", func(c *MonsterCatalogue) { c.Monsters[2].Senses.Notes = map[string]string{} }},
		{"melee range", func(c *MonsterCatalogue) { c.Monsters[0].Actions[0].Attack.Range = 5 }},
		{"further dice without first", func(c *MonsterCatalogue) {
			c.Monsters[0].Traits[0].MoreDice = []MonsterKeyedDice{{Key: "effect/future-damage", Dice: MonsterHP{Average: 3, Dice: 1, Sides: 4, Modifier: 1}}}
		}},
		{"XP off CR", func(c *MonsterCatalogue) { c.Monsters[10].XP = 5001 }},
		{"save off modifier", func(c *MonsterCatalogue) { c.Monsters[10].Abilities["wis"] = MonsterAbility{17, 3, 8} }},
		{"modifier", func(c *MonsterCatalogue) { c.Monsters[0].Abilities["str"] = MonsterAbility{15, 3, 3} }},
		{"HP average", func(c *MonsterCatalogue) { c.Monsters[4].HP.Average = 12 }},
		{"hover without flight", func(c *MonsterCatalogue) { c.Monsters[0].Hover = true }},
		{"empty bonus actions", func(c *MonsterCatalogue) { c.Monsters[1].BonusActions = []MonsterAction{} }},
		{"empty reactions", func(c *MonsterCatalogue) { c.Monsters[1].Reactions = []MonsterAction{} }},
		{"duplicate action across sections", func(c *MonsterCatalogue) { c.Monsters[0].BonusActions[0].ID = "ram" }},
		{"recharge", func(c *MonsterCatalogue) { c.Monsters[2].Actions[1].Uses.Recharge = 1 }},
		{"daily uses with recharge", func(c *MonsterCatalogue) { c.Monsters[2].Actions[1].Uses.PerDay = 1 }},
		{"save DC off entry text", func(c *MonsterCatalogue) { c.Monsters[2].Actions[1].Save.DC = 14 }},
		{"save ability off entry text", func(c *MonsterCatalogue) { c.Monsters[2].Actions[1].Save.Ability = "str" }},
		{"trait save DC off entry text", func(c *MonsterCatalogue) { c.Monsters[9].Traits[1].Save.DC = 11 }},
		{"creature size", func(c *MonsterCatalogue) { c.Monsters[2].Actions[1].Target.Size = 5 }},
		{"sphere width", func(c *MonsterCatalogue) { c.Monsters[9].Actions[1].Target.Width = 5 }},
		{"unlisted target shape", func(c *MonsterCatalogue) { c.Monsters[9].Actions[1].Target.Shape = "cube" }},
		{"effect dice", func(c *MonsterCatalogue) { c.Monsters[3].Actions[1].Effects[6].Dice.Average = 9 }},
		{"unlisted action effect", func(c *MonsterCatalogue) { c.Monsters[10].Actions[1].Effects[1].Key = "effect/future" }},
		{"untyped spellcasting", func(c *MonsterCatalogue) { c.Monsters[10].Actions[3].Spellcasting = nil }},
		{"misplaced spellcasting", func(c *MonsterCatalogue) {
			c.Monsters[10].Actions[2].Spellcasting = c.Monsters[10].Actions[3].Spellcasting
		}},
		{"spell ability", func(c *MonsterCatalogue) { c.Monsters[10].Actions[3].Spellcasting.Ability = "" }},
		{"spell DC", func(c *MonsterCatalogue) { c.Monsters[10].Actions[3].Spellcasting.SaveDC = 0 }},
		{"negative spell attack bonus", func(c *MonsterCatalogue) { c.Monsters[10].Actions[3].Spellcasting.AttackBonus = -1 }},
		{"unlisted spell", func(c *MonsterCatalogue) { c.Monsters[10].Actions[3].Spellcasting.AtWill[0] = "Wish" }},
		{"duplicate spell", func(c *MonsterCatalogue) { c.Monsters[10].Actions[3].Spellcasting.AtWill[0] = "Confusion" }},
		{"daily spell uses", func(c *MonsterCatalogue) { c.Monsters[10].Actions[3].Spellcasting.PerDay[0].Uses = 0 }},
		{"weapon action without typed attack", func(c *MonsterCatalogue) { c.Monsters[1].Actions[1].Attack = nil }},
		{"attack damage type off effect", func(c *MonsterCatalogue) { c.Monsters[1].Actions[1].Attack.DamageType = "slashing" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters162(t)
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
	raw, err := os.ReadFile("testdata/monsters_162.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, edit := range [][2]string{
		{`"lair_per_day": 0,`, ``},
		{`"long_range": 0,`, `"long_range": null,`},
		{`"range": 0,`, ``},
		{`"recharge": 5`, `"recharge": null`},
		{`"per_day": 0,`, ``},
		{`"dc": 13`, `"dc": 13, "dc": 14`},
		{`"width": 0,`, ``},
		{`"special": {}`, `"special": {"darkvision": null}`},
		{`"shape": "sphere",`, `"shape": "sphere", "height": 5,`},
		{`"passive_perception": 11,`, ``},
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
func TestMonsters162ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters162(t)
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
	for _, id := range []string{"bite", "claw", "club", "crush", "gouge", "holy-mace", "pact-blade", "pincer", "ram", "rend", "ritual-sickle", "sting"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}

// Printed values the #143 contract has no typed member for stay in each entry's
// retained source text, read here as literal PDF strings: the Spear's melee or
// ranged attack, the Shark's attack-roll condition, escape DCs, the Toad's
// swallow rules, the Spider's web statistics and the Mouther's random actions.
var monsters162Printed = []struct{ name, entry, text string }{
	{"Giant Seahorse", "ram", "Hit: 9 (2d6 + 2) Bludgeoning damage, or 11 (2d8 + 2) Bludgeoning damage if the seahorse moved 20+ feet straight toward the target immediately before the hit."},
	{"Giant Shark", "bite", "Melee Attack Roll: +9 (with Advantage if the target doesn’t have all its Hit Points), reach 5 ft."},
	{"Giant Spider", "web", "The target has the Restrained condition until the web is destroyed (AC 10; HP 5; Vulnerability to Fire damage; Immunity to Poison and Psychic damage)."},
	{"Giant Toad", "bite", "it has the Grappled condition (escape DC 12)."},
	{"Giant Toad", "swallow", "the target takes 10 (3d6) Acid damage at the end of each of the toad’s turns. The toad can have only one target swallowed at a time, and it can’t use Bite while it has a swallowed target."},
	{"Giant Venomous Snake", "bite", "Melee Attack Roll: +6, reach 10 ft."},
	{"Gibbering Mouther", "trait/gibbering", "rolls 1d8 to determine what it does during the current turn: 1–4. The target does nothing. 5–6. The target takes no action or Bonus Action and uses all its movement to move in a random direction. 7–8. The target makes a melee attack against a randomly determined creature within its reach or does nothing if it can’t make such an attack."},
	{"Glabrezu", "pincer", "it has the Grappled condition (escape DC 15) from one of two pincers."},
	{"Gladiator", "spear", "Spear. Melee or Ranged Attack Roll: +7, reach 5 ft. or range 20/60 ft. Hit: 11 (2d6 + 4) Piercing damage."},
}

func TestMonsters162PrintedSourceText(t *testing.T) {
	for _, want := range monsters162Printed {
		var text string
		for _, m := range loadMonsters162(t).Monsters {
			if m.Name != want.name {
				continue
			}
			for _, trait := range m.Traits {
				if trait.Key == want.entry {
					text = trait.Text
				}
			}
			for _, a := range m.allActions() {
				if a.ID == want.entry {
					text = a.Text
				}
			}
		}
		if !strings.Contains(text, want.text) {
			t.Fatalf("%s %s lost its printed text %q", want.name, want.entry, want.text)
		}
	}
	// The Spear carries no typed attack; the catalogue must not invent one.
	for _, m := range loadMonsters162(t).Monsters {
		for _, a := range m.allActions() {
			if (a.Key == "action/melee-or-ranged-attack") != (m.Name == "Gladiator" && a.ID == "spear") || a.Key == "action/melee-or-ranged-attack" && a.Attack != nil {
				t.Fatalf("%s %s: only the Gladiator's Spear is an untyped melee or ranged attack", m.Name, a.Name)
			}
		}
	}
}
