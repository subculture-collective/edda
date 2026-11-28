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

func loadMonsters176(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_176.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #176 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Warrior Infantry’s, the Werebear’s and the White Dragon Wyrmling’s blocks run onto the next
// page.
var monsters176Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Warhorse", [2]int{364, 364}, "Large Beast, Unaligned|AC 11 Initiative +1 (11)|HP 19 (3d10 + 3)|Speed 60 ft.|Str 18 +4 +4 Dex 12 +1 +1 Con 13 +1 +1|Int 2 −4 −4 Wis 12 +1 +3 Cha 7 −2 −2|Senses Passive Perception 11|Languages None|CR 1/2 (XP 100; PB +2)"},
	{"Warhorse Skeleton", [2]int{326, 326}, "Large Undead, Lawful Evil|AC 13 Initiative +1 (11)|HP 22 (3d10 + 6)|Speed 60 ft.|Str 18 +4 +4 Dex 12 +1 +1 Con 15 +2 +2|Int 2 −4 −4 Wis 8 −1 −1 Cha 5 −3 −3|Vulnerabilities Bludgeoning|Immunities Poison; Exhaustion, Poisoned|Senses Darkvision 60 ft.; Passive Perception 9|Languages None|CR 1/2 (XP 100; PB +2)"},
	{"Warrior Infantry", [2]int{336, 337}, "Medium or Small Humanoid, Neutral|AC 13 Initiative +0 (10)|HP 9 (2d8)|Speed 30 ft.|Str 13 +1 +1 Dex 11 +0 +0 Con 11 +0 +0|Int 8 −1 −1 Wis 11 +0 +0 Cha 8 −1 −1|Gear Chain Shirt, Spear|Senses Passive Perception 10|Languages Common|CR 1/8 (XP 25; PB +2)"},
	{"Warrior Veteran", [2]int{337, 337}, "Medium or Small Humanoid, Neutral|AC 17 Initiative +3 (13)|HP 65 (10d8 + 20)|Speed 30 ft.|Str 16 +3 +3 Dex 13 +1 +1 Con 14 +2 +2|Int 10 +0 +0 Wis 11 +0 +0 Cha 10 +0 +0|Skills Athletics +5, Perception +2|Gear Greatsword, Heavy Crossbow, Splint Armor|Senses Passive Perception 12|Languages Common plus one other language|CR 3 (XP 700; PB +2)"},
	{"Water Elemental", [2]int{337, 337}, "Large Elemental, Neutral|AC 14 Initiative +2 (12)|HP 114 (12d10 + 48)|Speed 30 ft., Swim 90 ft.|Str 18 +4 +4 Dex 14 +2 +2 Con 18 +4 +4|Int 5 −3 −3 Wis 10 +0 +0 Cha 8 −1 −1|Resistances Acid, Fire|Immunities Poison; Exhaustion, Grappled, Paralyzed, Petrified, Poisoned, Prone, Restrained, Unconscious|Senses Darkvision 60 ft.; Passive Perception 10|Languages Primordial (Aquan)|CR 5 (XP 1,800; PB +3)"},
	{"Weasel", [2]int{364, 364}, "Tiny Beast, Unaligned|AC 13 Initiative +3 (13)|HP 1 (1d4 − 1)|Speed 30 ft., Climb 30 ft.|Str 3 −4 −4 Dex 16 +3 +3 Con 8 −1 −1|Int 2 −4 −4 Wis 12 +1 +1 Cha 3 −4 −4|Skills Acrobatics +5, Perception +3, Stealth +5|Senses Darkvision 60 ft.; Passive Perception 13|Languages None|CR 0 (XP 10; PB +2)"},
	{"Werebear", [2]int{337, 338}, "Medium or Small Monstrosity (Lycanthrope), Neutral Good|AC 15 Initiative +3 (13)|HP 135 (18d8 + 54)|Speed 30 ft., 40 ft. (bear form only), Climb 30 ft. (bear form only)|Str 19 +4 +4 Dex 10 +0 +0 Con 17 +3 +3|Int 11 +0 +0 Wis 12 +1 +1 Cha 12 +1 +1|Skills Perception +7|Gear Handaxes (4)|Senses Darkvision 60 ft.; Passive Perception 17|Languages Common (can’t speak in bear form)|CR 5 (XP 1,800; PB +3)"},
	{"Wereboar", [2]int{338, 338}, "Medium or Small Monstrosity (Lycanthrope), Neutral Evil|AC 15 Initiative +2 (12)|HP 97 (15d8 + 30)|Speed 30 ft., 40 ft. (boar form only)|Str 17 +3 +3 Dex 10 +0 +0 Con 15 +2 +2|Int 10 +0 +0 Wis 11 +0 +0 Cha 8 −1 −1|Skills Perception +2|Gear Javelins (6)|Senses Passive Perception 12|Languages Common (can’t speak in boar form)|CR 4 (XP 1,100; PB +2)"},
	{"Wererat", [2]int{338, 338}, "Medium or Small Monstrosity (Lycanthrope), Lawful Evil|AC 13 Initiative +3 (13)|HP 60 (11d8 + 11)|Speed 30 ft., Climb 30 ft.|Str 10 +0 +0 Dex 16 +3 +3 Con 12 +1 +1|Int 11 +0 +0 Wis 10 +0 +0 Cha 8 −1 −1|Skills Perception +4, Stealth +5|Gear Hand Crossbow|Senses Darkvision 60 ft.; Passive Perception 14|Languages Common (can’t speak in rat form)|CR 2 (XP 450; PB +2)"},
	{"Weretiger", [2]int{339, 339}, "Medium or Small Monstrosity (Lycanthrope), Neutral|AC 12 Initiative +2 (12)|HP 120 (16d8 + 48)|Speed 30 ft., 40 ft. (tiger form only)|Str 17 +3 +3 Dex 15 +2 +2 Con 16 +3 +3|Int 10 +0 +0 Wis 13 +1 +1 Cha 11 +0 +0|Skills Perception +5, Stealth +4|Gear Longbow|Senses Darkvision 60 ft.; Passive Perception 15|Languages Common (can’t speak in tiger form)|CR 4 (XP 1,100; PB +2)"},
	{"Werewolf", [2]int{339, 339}, "Medium or Small Monstrosity (Lycanthrope), Chaotic Evil|AC 15 Initiative +4 (14)|HP 71 (11d8 + 22)|Speed 30 ft., 40 ft. (wolf form only)|Str 16 +3 +3 Dex 14 +2 +2 Con 14 +2 +2|Int 10 +0 +0 Wis 11 +0 +0 Cha 10 +0 +0|Skills Perception +4, Stealth +4|Gear Longbow|Senses Darkvision 60 ft.; Passive Perception 14|Languages Common (can’t speak in wolf form)|CR 3 (XP 700; PB +2)"},
	{"White Dragon Wyrmling", [2]int{339, 340}, "Medium Dragon (Chromatic), Chaotic Evil|AC 16 Initiative +2 (12)|HP 32 (5d8 + 10)|Speed 30 ft., Burrow 15 ft., Fly 60 ft., Swim 30 ft.|Str 14 +2 +2 Dex 10 +0 +2 Con 14 +2 +2|Int 5 −3 −3 Wis 10 +0 +2 Cha 11 +0 +0|Skills Perception +4, Stealth +2|Immunities Cold|Senses Blindsight 10 ft., Darkvision 60 ft.; Passive Perception 14|Languages Draconic|CR 2 (450 XP; PB +2)"},
}

// Typed values of every trait and action entry, reviewed against the same pages:
// limited uses, attack bonus/reach/range/damage, save ability and DC, target
// shape size/width/range, dice and dependency keys. A ranged attack printed with
// one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text, including each "Melee or Ranged Attack Roll". A lycanthrope heading's
// form qualifier is a requirement of its entry, and the Werebear's Handaxe
// prints "reach 5 ft" without its period.
var monsters176Entries = map[string][]string{
	"Warhorse": {
		"action|hooves|Hooves; melee +6 reach 5 range 0/0 9 (2d4 + 4) bludgeoning; effect/bludgeoning-damage 9 (2d4 + 4), requirement/target-size, requirement/straight-line-movement, effect/bludgeoning-damage 5 (2d4), condition/prone",
	},
	"Warhorse Skeleton": {
		"action|hooves|Hooves; melee +6 reach 5 range 0/0 7 (1d6 + 4) bludgeoning; effect/bludgeoning-damage 7 (1d6 + 4), requirement/target-size, requirement/straight-line-movement, condition/prone",
	},
	"Warrior Infantry": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage-with-nearby-ally, requirement/ally-not-incapacitated",
		"action|spear|Spear; effect/piercing-damage 4 (1d6 + 1)",
	},
	"Warrior Veteran": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack",
		"action|greatsword|Greatsword; melee +5 reach 5 range 0/0 10 (2d6 + 3) slashing; effect/slashing-damage 10 (2d6 + 3)",
		"action|heavy-crossbow|Heavy Crossbow; ranged +3 reach 0 range 100/400 12 (2d10 + 1) piercing; effect/piercing-damage 12 (2d10 + 1)",
		"reaction|parry|Parry; reaction/trigger, requirement/holding-weapon, effect/ac-bonus",
	},
	"Water Elemental": {
		"trait|trait/freeze|Freeze; requirement/takes-cold-damage, effect/speed-reduction",
		"trait|trait/water-form|Water Form; movement/enter-creature-space, movement/squeeze",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|slam|Slam; melee +7 reach 5 range 0/0 13 (2d8 + 4) bludgeoning; effect/bludgeoning-damage 13 (2d8 + 4), requirement/target-size, condition/prone",
		"action|whelm|Whelm (Recharge 4–6); str DC 15; save/str, target/creatures-in-space, limit/recharge, effect/bludgeoning-damage 22 (4d8 + 4), requirement/target-size, condition/grappled, condition/restrained, effect/suffocating, effect/start-of-turn-damage 9 (2d8), limit/grapple-capacity, effect/pull-creature-out, check/str-athletics, save/half-damage-only",
	},
	"Weasel": {
		"action|bite|Bite; melee +5 reach 5 range 0/0 1 piercing; effect/piercing-damage [damage/fixed]",
	},
	"Werebear": {
		"action|multiattack|Multiattack; action/melee-attack, action/melee-or-ranged-attack, multiattack/any-combination, multiattack/replace-attack",
		"action|bite|Bite; melee +7 reach 5 range 0/0 17 (2d12 + 4) piercing; con DC 14; effect/piercing-damage 17 (2d12 + 4), requirement/bear-or-hybrid-form, requirement/target-humanoid, save/con, effect/curse, effect/cursed-becomes-werebear, save/success-immunity",
		"action|handaxe|Handaxe; effect/slashing-damage 14 (3d6 + 4), requirement/humanoid-or-hybrid-form",
		"action|rend|Rend; melee +7 reach 5 range 0/0 13 (2d8 + 4) slashing; effect/slashing-damage 13 (2d8 + 4), requirement/bear-or-hybrid-form",
		"bonus|shape-shift|Shape-Shift; effect/shape-shift, effect/shapechange-form-limits, effect/equipment-not-transformed",
	},
	"Wereboar": {
		"action|multiattack|Multiattack; action/melee-attack, action/melee-or-ranged-attack, multiattack/any-combination, multiattack/replace-attack",
		"action|gore|Gore; melee +5 reach 5 range 0/0 12 (2d8 + 3) piercing; con DC 12; effect/piercing-damage 12 (2d8 + 3), requirement/boar-or-hybrid-form, requirement/target-humanoid, save/con, effect/curse, effect/cursed-becomes-wereboar, save/success-immunity",
		"action|javelin|Javelin; effect/piercing-damage 13 (3d6 + 3), requirement/humanoid-or-hybrid-form",
		"action|tusk|Tusk; melee +5 reach 5 range 0/0 10 (2d6 + 3) piercing; effect/piercing-damage 10 (2d6 + 3), requirement/boar-or-hybrid-form, requirement/target-size, requirement/straight-line-movement, effect/piercing-damage 7 (2d6), condition/prone",
		"bonus|shape-shift|Shape-Shift; effect/shape-shift, effect/shapechange-form-limits, effect/equipment-not-transformed",
	},
	"Wererat": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination, multiattack/replace-attack",
		"action|bite|Bite; melee +5 reach 5 range 0/0 8 (2d4 + 3) piercing; con DC 11; effect/piercing-damage 8 (2d4 + 3), requirement/rat-or-hybrid-form, requirement/target-humanoid, save/con, effect/curse, effect/cursed-becomes-wererat, save/success-immunity",
		"action|scratch|Scratch; melee +5 reach 5 range 0/0 6 (1d6 + 3) slashing; effect/slashing-damage 6 (1d6 + 3)",
		"action|hand-crossbow|Hand Crossbow; ranged +5 reach 0 range 30/120 6 (1d6 + 3) piercing; effect/piercing-damage 6 (1d6 + 3), requirement/humanoid-or-hybrid-form",
		"bonus|shape-shift|Shape-Shift; effect/shape-shift, effect/shapechange-form-limits, effect/equipment-not-transformed",
	},
	"Weretiger": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination, multiattack/replace-attack",
		"action|bite|Bite; melee +5 reach 5 range 0/0 12 (2d8 + 3) piercing; con DC 13; effect/piercing-damage 12 (2d8 + 3), requirement/tiger-or-hybrid-form, requirement/target-humanoid, save/con, effect/curse, effect/cursed-becomes-weretiger, save/success-immunity",
		"action|scratch|Scratch; melee +5 reach 5 range 0/0 10 (2d6 + 3) slashing; effect/slashing-damage 10 (2d6 + 3)",
		"action|longbow|Longbow; ranged +4 reach 0 range 150/600 11 (2d8 + 2) piercing; effect/piercing-damage 11 (2d8 + 2), requirement/humanoid-or-hybrid-form",
		"bonus|prowl|Prowl; requirement/tiger-or-hybrid-form, effect/move-without-opportunity-attacks, action/hide",
		"bonus|shape-shift|Shape-Shift; effect/shape-shift, effect/shapechange-form-limits, effect/equipment-not-transformed",
	},
	"Werewolf": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage-with-nearby-ally, requirement/ally-not-incapacitated",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination, multiattack/replace-attack",
		"action|bite|Bite; melee +5 reach 5 range 0/0 12 (2d8 + 3) piercing; con DC 12; effect/piercing-damage 12 (2d8 + 3), requirement/wolf-or-hybrid-form, requirement/target-humanoid, save/con, effect/curse, effect/cursed-becomes-werewolf, save/success-immunity",
		"action|scratch|Scratch; melee +5 reach 5 range 0/0 10 (2d6 + 3) slashing; effect/slashing-damage 10 (2d6 + 3)",
		"action|longbow|Longbow; ranged +4 reach 0 range 150/600 11 (2d8 + 2) piercing; effect/piercing-damage 11 (2d8 + 2), requirement/humanoid-or-hybrid-form",
		"bonus|shape-shift|Shape-Shift; effect/shape-shift, effect/shapechange-form-limits, effect/equipment-not-transformed",
	},
	"White Dragon Wyrmling": {
		"trait|trait/ice-walk|Ice Walk; movement/ice-walk, movement/ignore-difficult-terrain",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +4 reach 5 range 0/0 6 (1d8 + 2) slashing; effect/slashing-damage 6 (1d8 + 2), effect/cold-damage 2 (1d4)",
		"action|cold-breath|Cold Breath (Recharge 5–6); con DC 12; cone 15/0/0; save/con, area/cone, limit/recharge, effect/cold-damage 22 (5d8), save/half-damage",
	},
}

// printedSave176 holds a save the source prints without its sign: the Young
// White Dragon's Intelligence save reads "2" beside a −2 modifier (p. 340).
// The record keeps the derived −2.
var printedSave176 = map[string][2]string{"Young White Dragon": {"Int 6 −2 −2", "Int 6 −2 2"}}

// xpBeforeLabel176 prints "(450 XP; PB +2)" instead of "(XP 450; PB +2)"
// (p. 340).
var xpBeforeLabel176 = map[string]bool{"White Dragon Wyrmling": true, "Young White Dragon": true}

// monsterHeader176 extends monsterHeader170 with form-only speeds at the end
// of the speed list and the pinned printed exceptions of batches 26 to 28.
func monsterHeader176(m Monster) string {
	lines := strings.Split(monsterHeader170(m), "|")
	if m.ConditionalSpeedNote != "" {
		i := slices.IndexFunc(lines, func(l string) bool { return strings.HasPrefix(l, "Speed ") })
		lines[i] += ", " + m.ConditionalSpeedNote
	}
	header := strings.Join(lines, "|")
	if xpBeforeLabel176[m.Name] {
		header = strings.Replace(header, "(XP "+monsterThousands151(m.XP), "("+monsterThousands151(m.XP)+" XP", 1)
	}
	if save, ok := printedSave176[m.Name]; ok {
		header = strings.Replace(header, save[0], save[1], 1)
	}
	return header
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters176SourceValues(t *testing.T) {
	c := loadMonsters176(t)
	if len(c.Monsters) != 12 || len(monsters176Headers) != 12 || len(monsters176Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #176") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters176Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader176(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries164(m); !slices.Equal(got, monsters176Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters176Entries[want.name])
			}
			if m.VulnerabilityNote != "" || m.ResistanceNote != "" || m.TelepathyNote != "" || m.SwarmOf != "" || m.SpeedNote != "" || m.ConditionImmunityNote != "" || m.Telepathy != 0 || m.LairXP != 0 || m.Legendary != nil || m.LegendaryActions != nil || m.Senses.Notes != nil {
				t.Fatal("source lists no vulnerability note, resistance note, telepathy note, swarm line, speed note, condition immunity note, telepathy, lair XP, legendary actions, sense notes")
			}
		})
		for _, line := range monsters176Entries[want.name] {
			if !strings.HasPrefix(line, "legendary-uses|") {
				entries++
			}
		}
	}
	// Bold-italic entry headings counted on PDF pp. 326, 336–340 and 364.
	if entries != 45 {
		t.Fatalf("reviewed %d entries, want 45", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters176MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 176 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters176(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Animals appendix blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("animals monsters monsters monsters monsters animals monsters monsters monsters monsters monsters monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #176 cites the page each block starts on, in this order.
	names := []string{"Warhorse", "Warhorse Skeleton", "Warrior Infantry", "Warrior Veteran", "Water Elemental", "Weasel", "Werebear", "Wereboar", "Wererat", "Weretiger", "Werewolf", "White Dragon Wyrmling"}
	pages := []int{364, 326, 336, 337, 337, 364, 337, 338, 338, 339, 339, 339}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic for these twelve blocks; see monsterConsistency170.
func TestMonsters176DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"1/2": 2, "1/8": 2, "3": 2, "5": 3, "0": 2, "4": 2, "2": 2}
	attacks := map[string]attackAbility170{
		"Warrior Veteran/Heavy Crossbow": {"dex", true},
		"Weretiger/Longbow":              {"dex", true},
		"Werewolf/Longbow":               {"dex", true},
	}
	for _, m := range loadMonsters176(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			monsterConsistency170(t, m, pb, attacks)
		})
	}
}

// Every listed key is used and required, every effect quotes its own entry, and
// every dice expression is typed. Removing any one key invalidates the record.
func TestMonsters176DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters176(t).Monsters {
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
func TestMonsters176RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters176(t).Monsters {
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

func TestMonsters176MalformedFixtures(t *testing.T) {
	tests := []struct {
		name, record string
		mutate       func(*Monster)
	}{
		{"missing languages list", "Weasel", func(m *Monster) {
			m.Languages = nil
		}},
		{"animal or monster block in another family", "Warhorse", func(m *Monster) {
			m.ID = strings.Replace(strings.Replace(m.ID, "/animals/", "/spells/", 1), "/monsters/", "/spells/", 1)
		}},
		{"XP off CR", "Warhorse", func(m *Monster) {
			m.XP++
		}},
		{"save off modifier", "Warhorse", func(m *Monster) {
			m.Abilities["wis"] = MonsterAbility{m.Abilities["wis"].Score, m.Abilities["wis"].Modifier, m.Abilities["wis"].Modifier + 11}
		}},
		{"modifier", "Water Elemental", func(m *Monster) {
			m.Abilities["str"] = MonsterAbility{m.Abilities["str"].Score, m.Abilities["str"].Modifier + 1, m.Abilities["str"].Modifier + 1}
		}},
		{"HP average", "Water Elemental", func(m *Monster) {
			m.HP.Average++
		}},
		{"conditional speed note beside a typed mode", "Werebear", func(m *Monster) {
			m.Speeds["climb"] = 30
			m.RequiredMechanics = append(m.RequiredMechanics, "movement/climb")
		}},
		{"padded conditional speed note", "Wereboar", func(m *Monster) {
			m.ConditionalSpeedNote += " "
		}},
		{"fixed hit with a modifier", "Weasel", func(m *Monster) {
			m.Actions[0].Attack.Damage.Modifier = 1
		}},
		{"recharge", "Water Elemental", func(m *Monster) {
			m.Actions[2].Uses.Recharge = 7
		}},
		{"cone width", "White Dragon Wyrmling", func(m *Monster) {
			m.Actions[2].Target.Width = 5
		}},
		{"save DC off entry text", "Werebear", func(m *Monster) {
			m.Actions[1].Save.DC = 18
		}},
		{"empty bonus actions", "Warhorse", func(m *Monster) {
			m.BonusActions = []MonsterAction{}
		}},
		{"empty reactions", "Warhorse", func(m *Monster) {
			m.Reactions = []MonsterAction{}
		}},
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
		c := loadMonsters176(t)
		c.Monsters[1].ID = c.Monsters[2].ID
		rejected(t, c, true)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters176(t)
			tt.mutate(record(t, &c, tt.record))
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			rejected(t, c, !strings.HasPrefix(tt.name, "empty "))
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_176.json")
	if err != nil {
		t.Fatal(err)
	}
	// The fixtures are indented; compact them so each edit matches one field.
	var compact bytes.Buffer
	if err := json.Compact(&compact, raw); err != nil {
		t.Fatal(err)
	}
	for _, edit := range [][2]string{
		{`"long_range":0,`, `"long_range":null,`},
		{`"recharge":4`, `"recharge":null`},
		{`"per_day":0,`, ``},
		{`"dc":14`, `"dc":14,"dc":22`},
		{`"width":0,`, ``},
		{`"darkvision":60`, `"darkvision":null`},
		{`"shape":"cone",`, `"shape":"cone","height":5,`},
		{`"conditional_speed_note":`, `"conditional_speed_note":"40 ft. (bear form only)","conditional_speed_note":`},
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
// to 25 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters176ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters176(t)
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
	for _, load := range []func(*testing.T) MonsterCatalogue{loadMonsters158, loadMonsters159, loadMonsters160, loadMonsters161, loadMonsters162, loadMonsters163, loadMonsters164, loadMonsters165, loadMonsters166, loadMonsters167, loadMonsters168, loadMonsters169, loadMonsters170, loadMonsters171, loadMonsters172, loadMonsters173, loadMonsters174, loadMonsters175} {
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
	for _, id := range []string{"bite", "gore", "greatsword", "hand-crossbow", "heavy-crossbow", "hooves", "longbow", "rend", "scratch", "slam", "tusk"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
