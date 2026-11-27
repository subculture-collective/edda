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

func loadMonsters164(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_164.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #164 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Half-Dragon's Resistances line names no damage type (p. 297); the record
// keeps it as resistance_note. The Hippopotamus's block runs onto p. 356.
var monsters164Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Grimlock", [2]int{296, 296}, "Medium Aberration, Neutral Evil|AC 11 Initiative +1 (11)|HP 11 (2d8 + 2)|Speed 30 ft., Climb 30 ft.|Str 16 +3 +3 Dex 12 +1 +1 Con 12 +1 +1|Int 9 −1 −1 Wis 8 −1 −1 Cha 6 −2 −2|Skills Athletics +5, Perception +3, Stealth +5|Senses Blindsight 30 ft.; Passive Perception 13|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Guard", [2]int{296, 296}, "Medium or Small Humanoid, Neutral|AC 16 Initiative +1 (11)|HP 11 (2d8 + 2)|Speed 30 ft.|Str 13 +1 +1 Dex 12 +1 +1 Con 12 +1 +1|Int 10 +0 +0 Wis 11 +0 +0 Cha 10 +0 +0|Skills Perception +2|Gear Chain Shirt, Shield, Spear|Senses Passive Perception 12|Languages Common|CR 1/8 (XP 25; PB +2)"},
	{"Guard Captain", [2]int{296, 296}, "Medium or Small Humanoid, Neutral|AC 18 Initiative +4 (14)|HP 75 (10d8 + 30)|Speed 30 ft.|Str 18 +4 +4 Dex 14 +2 +2 Con 16 +3 +3|Int 12 +1 +1 Wis 14 +2 +2 Cha 13 +1 +1|Skills Athletics +6, Perception +4|Gear Breastplate, Javelins (6), Longsword, Shield|Senses Passive Perception 14|Languages Common|CR 4 (XP 1,100; PB +2)"},
	{"Guardian Naga", [2]int{296, 296}, "Large Celestial, Lawful Good|AC 18 Initiative +4 (14)|HP 136 (16d10 + 48)|Speed 40 ft., Climb 40 ft., Swim 40 ft.|Str 19 +4 +4 Dex 18 +4 +8 Con 16 +3 +7|Int 16 +3 +7 Wis 19 +4 +8 Cha 18 +4 +8|Skills Arcana +11, History +11, Religion +11|Immunities Poison; Charmed, Paralyzed, Poisoned, Restrained|Senses Darkvision 60 ft.; Passive Perception 14|Languages Celestial, Common|CR 10 (XP 5,900; PB +4)"},
	{"Half-Dragon", [2]int{297, 297}, "Medium Dragon, Neutral|AC 18 Initiative +5 (15)|HP 105 (14d8 + 42)|Speed 40 ft.|Str 19 +4 +4 Dex 14 +2 +5 Con 16 +3 +3|Int 10 +0 +0 Wis 15 +2 +5 Cha 14 +2 +2|Skills Athletics +7, Perception +5, Stealth +5|Resistances Damage type chosen for the Draconic Origin trait below|Senses Blindsight 10 ft., Darkvision 60 ft.; Passive Perception 15|Languages Common, Draconic|CR 5 (XP 1,800; PB +3)"},
	{"Harpy", [2]int{297, 297}, "Medium Monstrosity, Chaotic Evil|AC 11 Initiative +1 (11)|HP 38 (7d8 + 7)|Speed 20 ft., Fly 40 ft.|Str 12 +1 +1 Dex 13 +1 +1 Con 12 +1 +1|Int 7 −2 −2 Wis 10 +0 +0 Cha 13 +1 +1|Senses Passive Perception 10|Languages Common|CR 1 (XP 200; PB +2)"},
	{"Hawk", [2]int{355, 355}, "Tiny Beast, Unaligned|AC 13 Initiative +3 (13)|HP 1 (1d4 − 1)|Speed 10 ft., Fly 60 ft.|Str 5 −3 −3 Dex 16 +3 +3 Con 8 −1 −1|Int 2 −4 −4 Wis 14 +2 +2 Cha 6 −2 −2|Skills Perception +6|Senses Passive Perception 16|Languages None|CR 0 (XP 10; PB +2)"},
	{"Hell Hound", [2]int{297, 297}, "Medium Fiend, Lawful Evil|AC 15 Initiative +1 (11)|HP 58 (9d8 + 18)|Speed 50 ft.|Str 17 +3 +3 Dex 12 +1 +1 Con 14 +2 +2|Int 6 −2 −2 Wis 13 +1 +1 Cha 6 −2 −2|Skills Perception +5|Immunities Fire|Senses Darkvision 60 ft.; Passive Perception 15|Languages Understands Infernal but can’t speak|CR 3 (XP 700; PB +2)"},
	{"Hezrou", [2]int{298, 298}, "Large Fiend (Demon), Chaotic Evil|AC 18 Initiative +6 (16)|HP 157 (15d10 + 75)|Speed 30 ft.|Str 19 +4 +7 Dex 17 +3 +3 Con 20 +5 +8|Int 5 −3 −3 Wis 12 +1 +4 Cha 13 +1 +1|Resistances Cold, Fire, Lightning|Immunities Poison; Poisoned|Senses Darkvision 120 ft.; Passive Perception 11|Languages Abyssal; telepathy 120 ft.|CR 8 (XP 3,900; PB +3)"},
	{"Hill Giant", [2]int{298, 298}, "Huge Giant, Chaotic Evil|AC 13 Initiative +2 (12)|HP 105 (10d12 + 40)|Speed 40 ft.|Str 21 +5 +5 Dex 8 −1 −1 Con 19 +4 +4|Int 5 −3 −3 Wis 9 −1 −1 Cha 6 −2 −2|Skills Perception +2|Senses Passive Perception 12|Languages Giant|CR 5 (XP 1,800; PB +3)"},
	{"Hippogriff", [2]int{298, 298}, "Large Monstrosity, Unaligned|AC 11 Initiative +1 (11)|HP 26 (4d10 + 4)|Speed 40 ft., Fly 60 ft.|Str 17 +3 +3 Dex 13 +1 +1 Con 13 +1 +1|Int 2 −4 −4 Wis 12 +1 +1 Cha 8 −1 −1|Skills Perception +5|Senses Passive Perception 15|Languages None|CR 1 (XP 200; PB +2)"},
	{"Hippopotamus", [2]int{355, 356}, "Large Beast, Unaligned|AC 14 Initiative −2 (8)|HP 82 (11d10 + 22)|Speed 30 ft., Swim 30 ft.|Str 21 +5 +7 Dex 7 −2 −2 Con 15 +2 +2|Int 2 −4 −4 Wis 12 +1 +1 Cha 4 −3 −3|Skills Perception +3|Senses Passive Perception 13|Languages None|CR 4 (XP 1,100; PB +2)"},
}

// Typed values of every trait and action entry, reviewed against the same pages:
// limited uses, attack bonus/reach/range/damage, save ability and DC, target
// shape size/width/range, dice and dependency keys. A ranged attack printed with
// one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text. The Guard's Spear and Guard Captain's Javelin are "Melee or Ranged Attack
// Roll" entries (key action/melee-or-ranged-attack, damage typed on their effect, as for the
// #153 Archmage). The Half-Dragon's chosen damage type stays effect/chosen-type-damage.
var monsters164Entries = map[string][]string{
	"Grimlock": {
		"action|bone-cudgel|Bone Cudgel; melee +5 reach 5 range 0/0 6 (1d6 + 3) bludgeoning; effect/bludgeoning-damage 6 (1d6 + 3), effect/psychic-damage 2 (1d4)",
	},
	"Guard": {
		"action|spear|Spear; effect/piercing-damage 4 (1d6 + 1)",
	},
	"Guard Captain": {
		"action|multiattack|Multiattack; action/melee-attack, action/melee-or-ranged-attack, multiattack/any-combination",
		"action|javelin|Javelin; effect/piercing-damage 14 (3d6 + 4)",
		"action|longsword|Longsword; melee +6 reach 5 range 0/0 15 (2d10 + 4) slashing; effect/slashing-damage 15 (2d10 + 4)",
	},
	"Guardian Naga": {
		"trait|trait/celestial-restoration|Celestial Restoration; effect/revive, effect/revival-delay, requirement/not-under-dispel-evil-and-good, spell/dispel-evil-and-good",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, action/saving-throw",
		"action|bite|Bite; melee +8 reach 10 range 0/0 17 (2d12 + 4) piercing; effect/piercing-damage 17 (2d12 + 4), effect/poison-damage 22 (4d10)",
		"action|poisonous-spittle|Poisonous Spittle; con DC 16; creature 0/0/60; save/con, target/creature, effect/poison-damage 31 (7d8), condition/blinded, save/half-damage",
		"action|spellcasting|Spellcasting; wis DC 16 +0; at will: Thaumaturgy; 1/day: Clairvoyance, Cure Wounds (level 6 version), Flame Strike (level 6 version), Geas, True Seeing; spellcasting/no-somatic-or-material-components, spell/thaumaturgy [spellcasting/at-will], spell/clairvoyance [limit/per-day], spell/cure-wounds [limit/per-day], spell/flame-strike [limit/per-day], spell/geas [limit/per-day], spell/true-seeing [limit/per-day]",
	},
	"Half-Dragon": {
		"trait|trait/draconic-origin|Draconic Origin; effect/gm-chosen-damage-type, effect/chosen-type-modifies-stat-block",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|claw|Claw; melee +7 reach 10 range 0/0 6 (1d4 + 4) slashing; effect/slashing-damage 6 (1d4 + 4), effect/chosen-type-damage 7 (2d6)",
		"action|dragons-breath|Dragon’s Breath (Recharge 5–6); dex DC 14; cone 30/0/0; limit/recharge, save/dex, area/cone, effect/chosen-type-damage 28 (8d6), save/half-damage",
		"bonus|leap|Leap; effect/jump",
	},
	"Harpy": {
		"action|claw|Claw; melee +3 reach 5 range 0/0 6 (2d4 + 1) slashing; effect/slashing-damage 6 (2d4 + 1)",
		"action|luring-song|Luring Song; wis DC 11; emanation 300/0/0; effect/magical-song, effect/concentration, save/wis, area/emanation, target/humanoids-and-giants, condition/charmed, save/repeat, condition/incapacitated, effect/ignore-other-harpy-songs, effect/forced-movement-toward-source, effect/no-opportunity-attack-avoidance, save/repeat-before-hazard-or-on-damage, save/success-immunity",
	},
	"Hawk": {
		"action|talons|Talons; melee +5 reach 5 range 0/0 1 slashing; effect/slashing-damage [damage/fixed]",
	},
	"Hell Hound": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage, requirement/ally-within-5-feet, requirement/ally-not-incapacitated",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +5 reach 5 range 0/0 7 (1d8 + 3) piercing; effect/piercing-damage 7 (1d8 + 3), effect/fire-damage 3 (1d6)",
		"action|fire-breath|Fire Breath (Recharge 5–6); dex DC 12; cone 15/0/0; limit/recharge, save/dex, area/cone, effect/fire-damage 17 (5d6), save/half-damage",
	},
	"Hezrou": {
		"trait|trait/demonic-restoration|Demonic Restoration; effect/revive",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"trait|trait/stench|Stench; con DC 16; emanation 10/0/0; save/con, area/emanation, target/creatures-starting-turn-in-area, condition/poisoned",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +7 reach 5 range 0/0 6 (1d4 + 4) slashing; effect/slashing-damage 6 (1d4 + 4), effect/poison-damage 9 (2d8)",
		"bonus|leap|Leap; effect/jump",
	},
	"Hill Giant": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|tree-club|Tree Club; melee +8 reach 10 range 0/0 18 (3d8 + 5) bludgeoning; effect/bludgeoning-damage 18 (3d8 + 5), requirement/target-size, condition/prone",
		"action|trash-lob|Trash Lob; ranged +8 reach 0 range 60/240 16 (2d10 + 5) bludgeoning; effect/bludgeoning-damage 16 (2d10 + 5), condition/poisoned",
	},
	"Hippogriff": {
		"trait|trait/flyby|Flyby; effect/no-opportunity-attack-flying-out",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +5 reach 5 range 0/0 7 (1d8 + 3) slashing; effect/slashing-damage 7 (1d8 + 3)",
	},
	"Hippopotamus": {
		"trait|trait/hold-breath|Hold Breath; effect/hold-breath",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +7 reach 5 range 0/0 16 (2d10 + 5) piercing; effect/piercing-damage 16 (2d10 + 5)",
	},
}

// printedInitiative164 lists Initiative scores the PDF prints that differ from
// 10 plus the modifier, which batches 164–169 keep.
var printedInitiative164 = map[string]int{"Invisible Stalker": 22}

// monsterHeader164 prints typed fields in the source's stat-block notation.
func monsterHeader164(m Monster) string {
	score := m.Initiative.Score
	if printed, ok := printedInitiative164[m.Name]; ok {
		score = printed
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
	if m.Vulnerabilities != nil {
		lines = append(lines, "Vulnerabilities "+monsterTitles151(m.Vulnerabilities))
	}
	if m.Resistances != nil {
		lines = append(lines, "Resistances "+monsterTitles151(m.Resistances))
	}
	if m.ResistanceNote != "" {
		lines = append(lines, "Resistances "+m.ResistanceNote)
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
	if m.TelepathyNote != "" {
		languages += " (" + m.TelepathyNote + ")"
	}
	xp := monsterThousands151(m.XP)
	if m.LairXP > 0 {
		xp += ", or " + monsterThousands151(m.LairXP) + " in lair"
	}
	lines = append(lines, "Senses "+passive, languages, fmt.Sprintf("CR %s (XP %s; PB +%d)", m.CR, xp, m.ProficiencyBonus))
	return strings.Join(lines, "|")
}

// monsterEntries164 summarises every entry as monsterEntries163 does, then
// adds legendary actions and their use counts.
func monsterEntries164(m Monster) []string {
	out := monsterEntries163(m)
	for _, a := range m.LegendaryActions {
		out = append(out, monsterEntry163("legendary", a.ID, MonsterMechanic{Name: a.Name, Uses: a.Uses, Save: a.Save, Target: a.Target}, a.Attack, a.Spellcasting, a.Effects))
	}
	if m.Legendary != nil {
		out = append(out, fmt.Sprintf("legendary-uses|%d|%d", m.Legendary.Uses, m.Legendary.LairUses))
	}
	// Lair daily uses are pinned beside the ordinary count; traits come first.
	for i, trait := range m.Traits {
		if u := trait.Uses; u != nil && u.LairPerDay > 0 {
			out[i] = strings.Replace(out[i], fmt.Sprintf("(%d/Day)", u.PerDay), fmt.Sprintf("(%d/Day, or %d/Day in Lair)", u.PerDay, u.LairPerDay), 1)
		}
	}
	// A "Recharge 6" entry prints that number alone, not a range.
	for i := range out {
		out[i] = strings.Replace(out[i], "(Recharge 6–6)", "(Recharge 6)", 1)
	}
	return out
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters164SourceValues(t *testing.T) {
	c := loadMonsters164(t)
	if len(c.Monsters) != 12 || len(monsters164Headers) != 12 || len(monsters164Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #164") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters164Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader164(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries164(m); !slices.Equal(got, monsters164Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters164Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.Hover || m.LairXP != 0 || m.Legendary != nil || m.LegendaryActions != nil {
				t.Fatal("source lists no vulnerabilities, hover, lair XP, legendary actions")
			}
		})
		for _, line := range monsters164Entries[want.name] {
			if !strings.HasPrefix(line, "legendary-uses|") {
				entries++
			}
		}
	}
	// Bold-italic entry headings counted on PDF pp. 296–298 and 355–356.
	if entries != 37 {
		t.Fatalf("reviewed %d entries, want 37", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters164MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 164 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters164(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Animals appendix blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("monsters monsters monsters monsters monsters monsters animals monsters monsters monsters monsters animals")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #164 cites the page each block starts on, in this order.
	names := []string{"Grimlock", "Guard", "Guard Captain", "Guardian Naga", "Half-Dragon", "Harpy", "Hawk", "Hell Hound", "Hezrou", "Hill Giant", "Hippogriff", "Hippopotamus"}
	pages := []int{296, 296, 296, 296, 297, 297, 355, 297, 298, 298, 298, 355}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic that holds for these twelve blocks (PDF pp. 254–256):
// XP and PB follow CR, Hit Dice follow size and Constitution, saves add nothing
// or PB to the ability modifier, skills add PB or twice PB, and attack rolls add
// PB to the listed ability. Initiative adds nothing, PB or twice PB to Dexterity.
func TestMonsters164DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"1/4": 2, "1/8": 2, "4": 2, "10": 4, "5": 3, "1": 2, "0": 2, "3": 2, "8": 3}
	skillAbility := map[string]string{"arcana": "int", "athletics": "str", "history": "int", "perception": "wis", "religion": "int", "stealth": "dex"}
	hitDie := map[string]int{"tiny": 4, "medium": 8, "large": 10, "huge": 12}
	for _, m := range loadMonsters164(t).Monsters {
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
			if extra := m.Initiative.Modifier - m.Abilities["dex"].Modifier; extra != 0 && extra != p && extra != 2*p {
				t.Fatal("initiative is not Dexterity with optional PB or twice PB")
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
func TestMonsters164DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters164(t).Monsters {
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
func TestMonsters164RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters164(t).Monsters {
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

func TestMonsters164MalformedFixtures(t *testing.T) {
	tests := []struct {
		name, record string
		mutate       func(*Monster)
	}{
		{"missing languages list", "Guard", func(m *Monster) { m.Languages = nil }},
		{"animal block in another family", "Guard", func(m *Monster) {
			m.ID = strings.Replace(strings.Replace(m.ID, "/animals/", "/spells/", 1), "/monsters/", "/spells/", 1)
		}},
		{"XP off CR", "Hippopotamus", func(m *Monster) { m.XP++ }},
		{"save off modifier", "Hippopotamus", func(m *Monster) {
			m.Abilities["wis"] = MonsterAbility{m.Abilities["wis"].Score, m.Abilities["wis"].Modifier, m.Abilities["wis"].Modifier + 7}
		}},
		{"modifier", "Grimlock", func(m *Monster) {
			m.Abilities["str"] = MonsterAbility{m.Abilities["str"].Score, m.Abilities["str"].Modifier + 1, m.Abilities["str"].Modifier + 1}
		}},
		{"HP average", "Grimlock", func(m *Monster) { m.HP.Average++ }},
		{"resistance note with typed list", "Half-Dragon", func(m *Monster) { m.Resistances = []string{"fire"} }},
		{"padded resistance note", "Half-Dragon", func(m *Monster) { m.ResistanceNote += " " }},
		{"resistance note without mechanic", "Half-Dragon", func(m *Monster) {
			m.RequiredMechanics = slices.DeleteFunc(m.RequiredMechanics, func(key string) bool { return key == "defense/chosen-damage-resistance" })
		}},
		{"sense note elsewhere", "Grimlock", func(m *Monster) { m.Senses.Notes = map[string]string{"darkvision": "unimpeded by magical Darkness"} }},
		{"recharge", "Half-Dragon", func(m *Monster) { m.Actions[2].Uses.Recharge = 1 }},
		{"save DC off entry text", "Half-Dragon", func(m *Monster) { m.Actions[2].Save.DC = 15 }},
		{"creature size", "Guardian Naga", func(m *Monster) { m.Actions[2].Target.Size = 5 }},
		{"unlisted spell", "Guardian Naga", func(m *Monster) { m.Actions[3].Spellcasting.AtWill[0] = "Wish" }},
		{"untyped spellcasting", "Guardian Naga", func(m *Monster) { m.Actions[3].Spellcasting = nil }},
		{"spell DC off entry text", "Guardian Naga", func(m *Monster) { m.Actions[3].Spellcasting.SaveDC = 17 }},
		{"fixed damage with sides", "Hawk", func(m *Monster) { m.Actions[0].Attack.Damage.Sides = 4 }},
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
		c := loadMonsters164(t)
		c.Monsters[1].ID = c.Monsters[2].ID
		rejected(t, c, true)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters164(t)
			tt.mutate(record(t, &c, tt.record))
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			rejected(t, c, !strings.HasPrefix(tt.name, "empty "))
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_164.json")
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
		{`"dc":16`, `"dc":16,"dc":17`},
		{`"width":0,`, ``},
		{`"save_dc":16,`, ``},
		{`"attack_bonus":0`, `"attack_bonus":0,"attack_bonus":5`},
		{`"darkvision":60`, `"darkvision":null`},
		{`"shape":"cone",`, `"shape":"cone","height":5,`},
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
// to 13 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters164ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters164(t)
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
	for _, load := range []func(*testing.T) MonsterCatalogue{loadMonsters158, loadMonsters159, loadMonsters160, loadMonsters161, loadMonsters162, loadMonsters163} {
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
	for _, id := range []string{"bite", "bone-cudgel", "claw", "longsword", "rend", "talons", "trash-lob", "tree-club"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
