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

func loadMonsters174(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_174.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #174 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Swarm of Venomous Snakes’ and the Tarrasque’s blocks run onto the next
// page.
var monsters174Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Swarm of Bats", [2]int{361, 361}, "Large Swarm of Tiny Beasts, Unaligned|AC 12 Initiative +2 (12)|HP 11 (2d10)|Speed 5 ft., Fly 30 ft.|Str 5 −3 −3 Dex 15 +2 +2 Con 10 +0 +0|Int 2 −4 −4 Wis 12 +1 +1 Cha 4 −3 −3|Resistances Bludgeoning, Piercing, Slashing|Immunities Charmed, Frightened, Grappled, Paralyzed, Petrified, Prone, Restrained, Stunned|Senses Blindsight 60 ft.; Passive Perception 11|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Swarm of Crawling Claws", [2]int{278, 278}, "Medium Swarm of Tiny Undead, Neutral Evil|AC 12 Initiative +2 (12)|HP 49 (11d8)|Speed 30 ft., Climb 30 ft.|Str 14 +2 +2 Dex 14 +2 +2 Con 11 +0 +0|Int 5 −3 −3 Wis 10 +0 +0 Cha 4 −3 −3|Resistances Bludgeoning, Piercing, Slashing|Immunities Necrotic, Poison; Charmed, Exhaustion, Frightened, Grappled, Incapacitated, Paralyzed, Petrified, Poisoned, Prone, Restrained, Stunned|Senses Blindsight 30 ft.; Passive Perception 10|Languages Understands Common but can’t speak|CR 3 (XP 700; PB +2)"},
	{"Swarm of Insects", [2]int{361, 361}, "Medium Swarm of Tiny Beasts, Unaligned|AC 11 Initiative +1 (11)|HP 19 (3d8 + 6)|Speed 20 ft., Climb or Fly 20 ft. (GM’s choice)|Str 3 −4 −4 Dex 13 +1 +1 Con 14 +2 +2|Int 1 −5 −5 Wis 7 −2 −2 Cha 1 −5 −5|Resistances Bludgeoning, Piercing, Slashing|Immunities Charmed, Frightened, Grappled, Paralyzed, Petrified, Prone, Restrained, Stunned|Senses Blindsight 30 ft.; Passive Perception 8|Languages None|CR 1/2 (XP 100; PB +2)"},
	{"Swarm of Piranhas", [2]int{362, 362}, "Medium Swarm of Tiny Beasts, Unaligned|AC 13 Initiative +3 (13)|HP 28 (8d8 − 8)|Speed 5 ft., Swim 40 ft.|Str 13 +1 +1 Dex 16 +3 +3 Con 9 −1 −1|Int 1 −5 −5 Wis 7 −2 −2 Cha 2 −4 −4|Resistances Bludgeoning, Piercing, Slashing|Immunities Charmed, Frightened, Grappled, Paralyzed, Petrified, Prone, Restrained, Stunned|Senses Darkvision 60 ft.; Passive Perception 8|Languages None|CR 1 (XP 200; PB +2)"},
	{"Swarm of Rats", [2]int{362, 362}, "Medium Swarm of Tiny Beasts, Unaligned|AC 10 Initiative +0 (10)|HP 14 (4d8 − 4)|Speed 30 ft., Climb 30 ft.|Str 9 −1 −1 Dex 11 +0 +2 Con 9 −1 −1|Int 2 −4 −4 Wis 10 +0 +0 Cha 3 −4 −4|Resistances Bludgeoning, Piercing, Slashing|Immunities Charmed, Frightened, Grappled, Paralyzed, Petrified, Prone, Restrained, Stunned|Senses Darkvision 30 ft.; Passive Perception 10|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Swarm of Ravens", [2]int{362, 362}, "Medium Swarm of Tiny Beasts, Unaligned|AC 12 Initiative +2 (12)|HP 11 (2d8 + 2)|Speed 10 ft., Fly 50 ft.|Str 6 −2 −2 Dex 14 +2 +2 Con 12 +1 +1|Int 5 −3 −3 Wis 12 +1 +1 Cha 6 −2 −2|Skills Perception +5|Resistances Bludgeoning, Piercing, Slashing|Immunities Charmed, Frightened, Grappled, Paralyzed, Petrified, Prone, Restrained, Stunned|Senses Passive Perception 15|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Swarm of Venomous Snakes", [2]int{362, 363}, "Medium Swarm of Tiny Beasts, Unaligned|AC 14 Initiative +4 (14)|HP 36 (8d8)|Speed 30 ft., Swim 30 ft.|Str 8 −1 −1 Dex 18 +4 +4 Con 11 +0 +0|Int 1 −5 −5 Wis 10 +0 +0 Cha 3 −4 −4|Resistances Bludgeoning, Piercing, Slashing|Immunities Charmed, Frightened, Grappled, Paralyzed, Petrified, Prone, Restrained, Stunned|Senses Blindsight 10 ft.; Passive Perception 10|Languages None|CR 2 (XP 450; PB +2)"},
	{"Tarrasque", [2]int{331, 332}, "Gargantuan Monstrosity (Titan), Unaligned|AC 25 Initiative +18 (28)|HP 697 (34d20 + 340)|Speed 60 ft., Burrow 40 ft., Climb 60 ft.|Str 30 +10 +10 Dex 11 +0 +9 Con 30 +10 +10|Int 3 −4 +5 Wis 11 +0 +9 Cha 11 +0 +9|Skills Perception +9|Resistances Bludgeoning, Piercing, Slashing|Immunities Fire, Poison; Charmed, Deafened, Frightened, Paralyzed, Poisoned|Senses Blindsight 120 ft.; Passive Perception 19|Languages None|CR 30 (XP 155,000; PB +9)"},
	{"Tiger", [2]int{363, 363}, "Large Beast, Unaligned|AC 13 Initiative +3 (13)|HP 30 (4d10 + 8)|Speed 40 ft.|Str 17 +3 +3 Dex 16 +3 +3 Con 14 +2 +2|Int 3 −4 −4 Wis 12 +1 +1 Cha 8 −1 −1|Skills Perception +3, Stealth +7|Senses Darkvision 60 ft.; Passive Perception 13|Languages None|CR 1 (XP 200; PB +2)"},
	{"Tough", [2]int{332, 332}, "Medium or Small Humanoid, Neutral|AC 12 Initiative +1 (11)|HP 32 (5d8 + 10)|Speed 30 ft.|Str 15 +2 +2 Dex 12 +1 +1 Con 14 +2 +2|Int 10 +0 +0 Wis 10 +0 +0 Cha 11 +0 +0|Gear Heavy Crossbow, Leather Armor, Mace|Senses Passive Perception 10|Languages Common|CR 1/2 (XP 100; PB +2)"},
	{"Tough Boss", [2]int{332, 332}, "Medium or Small Humanoid, Neutral|AC 16 Initiative +2 (12)|HP 82 (11d8 + 33)|Speed 30 ft.|Str 17 +3 +5 Dex 14 +2 +2 Con 16 +3 +5|Int 11 +0 +0 Wis 10 +0 +0 Cha 11 +0 +2|Gear Chain Mail, Heavy Crossbow, Warhammer|Senses Passive Perception 10|Languages Common plus one other language|CR 4 (XP 1,100; PB +2)"},
	{"Treant", [2]int{333, 333}, "Huge Plant, Chaotic Good|AC 16 Initiative +3 (13)|HP 138 (12d12 + 60)|Speed 30 ft.|Str 23 +6 +6 Dex 8 −1 −1 Con 21 +5 +5|Int 12 +1 +1 Wis 16 +3 +3 Cha 12 +1 +1|Vulnerabilities Fire|Resistances Bludgeoning, Piercing|Senses Passive Perception 13|Languages Common, Druidic, Elvish, Sylvan|CR 9 (XP 5,000; PB +4)"},
}

// Typed values of every trait and action entry, reviewed against the same pages:
// limited uses, attack bonus/reach/range/damage, save ability and DC, target
// shape size/width/range, dice and dependency keys. A ranged attack printed with
// one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text. A swarm's lower Bloodied damage requires requirement/swarm-bloodied.
var monsters174Entries = map[string][]string{
	"Swarm of Bats": {
		"trait|trait/swarm|Swarm; effect/occupy-creature-space, movement/squeeze, effect/no-hit-point-recovery",
		"action|bites|Bites; melee +4 reach 5 range 0/0 5 (2d4) piercing; effect/piercing-damage 5 (2d4), effect/piercing-damage 2 (1d4) [requirement/swarm-bloodied]",
	},
	"Swarm of Crawling Claws": {
		"trait|trait/swarm|Swarm; effect/occupy-creature-space, movement/squeeze, effect/no-hit-point-recovery",
		"action|swarm-of-grasping-hands|Swarm of Grasping Hands; melee +4 reach 5 range 0/0 20 (4d8 + 2) necrotic; effect/necrotic-damage 20 (4d8 + 2), effect/necrotic-damage 11 (2d8 + 2) [requirement/swarm-bloodied], requirement/target-size, condition/prone",
	},
	"Swarm of Insects": {
		"trait|trait/spider-climb|Spider Climb; requirement/climb-speed, effect/climb-without-check",
		"trait|trait/swarm|Swarm; effect/occupy-creature-space, movement/squeeze, effect/no-hit-point-recovery",
		"action|bites|Bites; melee +3 reach 5 range 0/0 6 (2d4 + 1) poison; effect/poison-damage 6 (2d4 + 1), effect/poison-damage 3 (1d4 + 1) [requirement/swarm-bloodied]",
	},
	"Swarm of Piranhas": {
		"trait|trait/swarm|Swarm; effect/occupy-creature-space, movement/squeeze, effect/no-hit-point-recovery",
		"trait|trait/water-breathing|Water Breathing; effect/breathe-only-underwater",
		"action|bites|Bites; melee +5 reach 5 range 0/0 8 (2d4 + 3) piercing; effect/piercing-damage 8 (2d4 + 3), attack/advantage-against-target-missing-hit-points, effect/piercing-damage 5 (1d4 + 3) [requirement/swarm-bloodied]",
	},
	"Swarm of Rats": {
		"trait|trait/swarm|Swarm; effect/occupy-creature-space, movement/squeeze, effect/no-hit-point-recovery",
		"action|bites|Bites; melee +2 reach 5 range 0/0 5 (2d4) piercing; effect/piercing-damage 5 (2d4), effect/piercing-damage 2 (1d4) [requirement/swarm-bloodied]",
	},
	"Swarm of Ravens": {
		"trait|trait/swarm|Swarm; effect/occupy-creature-space, movement/squeeze, effect/no-hit-point-recovery",
		"action|beaks|Beaks; melee +4 reach 5 range 0/0 5 (1d6 + 2) piercing; effect/piercing-damage 5 (1d6 + 2), effect/piercing-damage 2 (1d4) [requirement/swarm-bloodied]",
		"action|cacophony|Cacophony (Recharge 6); wis DC 10; creature 0/0/0; limit/recharge, save/wis, target/creature, condition/deafened, effect/deafened-disadvantage",
	},
	"Swarm of Venomous Snakes": {
		"trait|trait/swarm|Swarm; effect/occupy-creature-space, movement/squeeze, effect/no-hit-point-recovery",
		"action|bites|Bites; melee +6 reach 5 range 0/0 8 (1d8 + 4) piercing; effect/piercing-damage 8 (1d8 + 4), effect/piercing-damage 6 (1d4 + 4) [requirement/swarm-bloodied], effect/poison-damage 10 (3d6)",
	},
	"Tarrasque": {
		"trait|trait/legendary-resistance|Legendary Resistance (6/Day); limit/per-day, save/choose-success",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"trait|trait/reflective-carapace|Reflective Carapace; effect/spell-reflection, requirement/targeted-by-magic-missile-or-ranged-spell-attack",
		"trait|trait/siege-monster|Siege Monster; effect/double-damage-to-objects-and-structures",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/any-combination",
		"action|bite|Bite; melee +19 reach 15 range 0/0 36 (4d12 + 10) piercing; effect/piercing-damage 36 (4d12 + 10), condition/grappled, condition/restrained, effect/no-teleport",
		"action|claw|Claw; melee +19 reach 15 range 0/0 28 (4d8 + 10) slashing; effect/slashing-damage 28 (4d8 + 10)",
		"action|tail|Tail; melee +19 reach 30 range 0/0 23 (3d8 + 10) bludgeoning; effect/bludgeoning-damage 23 (3d8 + 10), requirement/target-size, condition/prone",
		"action|thunderous-bellow|Thunderous Bellow (Recharge 5–6); con DC 27; cone 150/0/0; limit/recharge, save/con, area/cone, target/unattended-objects, effect/thunder-damage 78 (12d12), condition/deafened, condition/frightened, save/half-damage-only",
		"bonus|swallow|Swallow; str DC 27; creature 0/0/0; save/str, target/creature, requirement/target-size, requirement/grappled-creature, limit/six-swallowed-creatures, effect/swallow, condition/blinded, condition/restrained, effect/no-teleport, effect/total-cover, effect/acid-damage 56 (16d6), effect/regurgitation-save, condition/prone, effect/escape-corpse",
		"legendary|onslaught|Onslaught; movement/half-speed, action/melee-attack",
		"legendary|world-shaking-movement|World-Shaking Movement; movement/full-speed, area/emanation, effect/lose-concentration, condition/prone, limit/not-again-until-next-turn",
		"legendary-uses|3|0",
	},
	"Tiger": {
		"action|rend|Rend; melee +5 reach 5 range 0/0 10 (2d6 + 3) slashing; effect/slashing-damage 10 (2d6 + 3), requirement/target-size, condition/prone",
		"bonus|nimble-escape|Nimble Escape; action/disengage, action/hide",
	},
	"Tough": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage-with-nearby-ally, requirement/ally-not-incapacitated",
		"action|mace|Mace; melee +4 reach 5 range 0/0 5 (1d6 + 2) bludgeoning; effect/bludgeoning-damage 5 (1d6 + 2)",
		"action|heavy-crossbow|Heavy Crossbow; ranged +3 reach 0 range 100/400 6 (1d10 + 1) piercing; effect/piercing-damage 6 (1d10 + 1)",
	},
	"Tough Boss": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage-with-nearby-ally, requirement/ally-not-incapacitated",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|warhammer|Warhammer; melee +5 reach 5 range 0/0 12 (2d8 + 3) bludgeoning; effect/bludgeoning-damage 12 (2d8 + 3), requirement/target-size, effect/push",
		"action|heavy-crossbow|Heavy Crossbow; ranged +4 reach 0 range 100/400 13 (2d10 + 2) piercing; effect/piercing-damage 13 (2d10 + 2)",
	},
	"Treant": {
		"trait|trait/siege-monster|Siege Monster; effect/double-damage-to-objects-and-structures",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|slam|Slam; melee +10 reach 5 range 0/0 16 (3d6 + 6) bludgeoning; effect/bludgeoning-damage 16 (3d6 + 6)",
		"action|hail-of-bark|Hail of Bark; ranged +10 reach 0 range 180/180 28 (4d10 + 6) piercing; effect/piercing-damage 28 (4d10 + 6)",
		"action|animate-trees|Animate Trees (1/Day); limit/per-day, effect/animate-trees, effect/animated-tree-statistics, effect/animated-tree-turn-order, effect/animation-duration, effect/tree-takes-root",
	},
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters174SourceValues(t *testing.T) {
	c := loadMonsters174(t)
	if len(c.Monsters) != 12 || len(monsters174Headers) != 12 || len(monsters174Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #174") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters174Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader170(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries164(m); !slices.Equal(got, monsters174Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters174Entries[want.name])
			}
			if m.VulnerabilityNote != "" || m.ResistanceNote != "" || m.ConditionImmunityNote != "" || m.TelepathyNote != "" || m.Senses.Notes != nil || m.LairXP != 0 {
				t.Fatal("source lists no vulnerability note, resistance note, condition immunity note, telepathy note, sense notes, lair XP")
			}
		})
		for _, line := range monsters174Entries[want.name] {
			if !strings.HasPrefix(line, "legendary-uses|") {
				entries++
			}
		}
	}
	// Bold-italic entry headings counted on PDF pp. 278, 331–333 and 361–363.
	if entries != 43 {
		t.Fatalf("reviewed %d entries, want 43", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters174MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 174 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters174(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Animals appendix blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("animals monsters animals animals animals animals animals monsters animals monsters monsters monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #174 cites the page each block starts on, in this order.
	names := []string{"Swarm of Bats", "Swarm of Crawling Claws", "Swarm of Insects", "Swarm of Piranhas", "Swarm of Rats", "Swarm of Ravens", "Swarm of Venomous Snakes", "Tarrasque", "Tiger", "Tough", "Tough Boss", "Treant"}
	pages := []int{361, 278, 361, 362, 362, 362, 362, 331, 363, 332, 332, 333}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic for these twelve blocks; see monsterConsistency170.
func TestMonsters174DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"1/4": 2, "3": 2, "1/2": 2, "1": 2, "2": 2, "30": 9, "4": 2, "9": 4}
	attacks := map[string]attackAbility170{
		"Swarm of Bats/Bites":       {"dex", false},
		"Tough/Heavy Crossbow":      {"dex", true},
		"Tough Boss/Heavy Crossbow": {"dex", true},
	}
	for _, m := range loadMonsters174(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			monsterConsistency170(t, m, pb, attacks)
		})
	}
}

// Every listed key is used and required, every effect quotes its own entry, and
// every dice expression is typed. Removing any one key invalidates the record.
func TestMonsters174DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters174(t).Monsters {
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
func TestMonsters174RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters174(t).Monsters {
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

func TestMonsters174MalformedFixtures(t *testing.T) {
	tests := []struct {
		name, record string
		mutate       func(*Monster)
	}{
		{"missing languages list", "Tiger", func(m *Monster) { m.Languages = nil }},
		{"animal or monster block in another family", "Swarm of Rats", func(m *Monster) {
			m.ID = strings.Replace(strings.Replace(m.ID, "/animals/", "/spells/", 1), "/monsters/", "/spells/", 1)
		}},
		{"XP off CR", "Tarrasque", func(m *Monster) { m.XP++ }},
		{"save off modifier", "Tarrasque", func(m *Monster) {
			m.Abilities["wis"] = MonsterAbility{m.Abilities["wis"].Score, m.Abilities["wis"].Modifier, m.Abilities["wis"].Modifier + 11}
		}},
		{"modifier", "Treant", func(m *Monster) {
			m.Abilities["str"] = MonsterAbility{m.Abilities["str"].Score, m.Abilities["str"].Modifier + 1, m.Abilities["str"].Modifier + 1}
		}},
		{"HP average", "Treant", func(m *Monster) { m.HP.Average++ }},
		{"unknown swarm member size", "Swarm of Bats", func(m *Monster) { m.SwarmOf = "small" }},
		{"swarm of its own size", "Swarm of Bats", func(m *Monster) { m.Sizes = []string{"tiny"} }},
		{"padded speed note", "Swarm of Insects", func(m *Monster) { m.SpeedNote += " " }},
		{"Bloodied dice average", "Swarm of Bats", func(m *Monster) { m.Actions[0].Effects[1].Dice.Average++ }},
		{"recharge", "Tarrasque", func(m *Monster) { m.Actions[4].Uses.Recharge = 7 }},
		{"save DC off entry text", "Tarrasque", func(m *Monster) { m.Actions[4].Save.DC = 26 }},
		{"cone width", "Tarrasque", func(m *Monster) { m.Actions[4].Target.Width = 5 }},
		{"daily uses without count", "Treant", func(m *Monster) { m.Actions[3].Uses.PerDay = 0 }},
		{"legendary actions without uses", "Tarrasque", func(m *Monster) { m.Legendary = nil }},
		{"empty bonus actions", "Tough", func(m *Monster) { m.BonusActions = []MonsterAction{} }},
		{"empty reactions", "Tough", func(m *Monster) { m.Reactions = []MonsterAction{} }},
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
		c := loadMonsters174(t)
		c.Monsters[1].ID = c.Monsters[2].ID
		rejected(t, c, true)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters174(t)
			tt.mutate(record(t, &c, tt.record))
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			rejected(t, c, !strings.HasPrefix(tt.name, "empty "))
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_174.json")
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
		{`"recharge":6`, `"recharge":null`},
		{`"per_day":0,`, ``},
		{`"dc":10`, `"dc":10,"dc":22`},
		{`"width":0,`, ``},
		{`"darkvision":60`, `"darkvision":null`},
		{`"shape":"creature",`, `"shape":"creature","height":5,`},
		{`"lair_uses":0,`, ``},
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
// to 23 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters174ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters174(t)
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
	for _, load := range []func(*testing.T) MonsterCatalogue{loadMonsters158, loadMonsters159, loadMonsters160, loadMonsters161, loadMonsters162, loadMonsters163, loadMonsters164, loadMonsters165, loadMonsters166, loadMonsters167, loadMonsters168, loadMonsters169, loadMonsters170, loadMonsters171, loadMonsters172, loadMonsters173} {
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
	for _, id := range []string{"beaks", "bite", "bites", "claw", "hail-of-bark", "heavy-crossbow", "mace", "rend", "slam", "swarm-of-grasping-hands", "tail", "warhammer"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
