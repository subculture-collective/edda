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

func loadMonsters177(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_177.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #177 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Will-o’-Wisp’s and the Wraith’s blocks run onto the next page.
var monsters177Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Wight", [2]int{341, 341}, "Medium Undead, Neutral Evil|AC 14 Initiative +4 (14)|HP 82 (11d8 + 33)|Speed 30 ft.|Str 15 +2 +2 Dex 14 +2 +2 Con 16 +3 +3|Int 10 +0 +0 Wis 13 +1 +1 Cha 15 +2 +2|Skills Perception +3, Stealth +4|Resistances Necrotic|Immunities Poison; Exhaustion, Poisoned|Gear Studded Leather Armor|Senses Darkvision 60 ft.; Passive Perception 13|Languages Common plus one other language|CR 3 (XP 700; PB +2)"},
	{"Will-o’-Wisp", [2]int{341, 342}, "Tiny Undead, Chaotic Evil|AC 19 Initiative +9 (19)|HP 27 (11d4)|Speed 5 ft., Fly 50 ft. (hover)|Str 1 −5 −5 Dex 28 +9 +9 Con 10 +0 +0|Int 13 +1 +1 Wis 14 +2 +2 Cha 11 +0 +0|Resistances Acid, Bludgeoning, Cold, Fire, Necrotic, Piercing, Slashing|Immunities Lightning, Poison; Exhaustion, Grappled, Paralyzed, Petrified, Poisoned, Prone, Restrained, Unconscious|Senses Darkvision 120 ft.; Passive Perception 12|Languages Common plus one other language|CR 2 (XP 450; PB +2)"},
	{"Winter Wolf", [2]int{342, 342}, "Large Monstrosity, Neutral Evil|AC 13 Initiative +1 (11)|HP 75 (10d10 + 20)|Speed 50 ft.|Str 18 +4 +4 Dex 13 +1 +1 Con 14 +2 +2|Int 7 −2 −2 Wis 12 +1 +1 Cha 8 −1 −1|Skills Perception +5, Stealth +5|Immunities Cold|Senses Passive Perception 15|Languages Common, Giant|CR 3 (XP 700; PB +2)"},
	{"Wolf", [2]int{364, 364}, "Medium Beast, Unaligned|AC 12 Initiative +2 (12)|HP 11 (2d8 + 2)|Speed 40 ft.|Str 14 +2 +2 Dex 15 +2 +2 Con 12 +1 +1|Int 3 −4 −4 Wis 12 +1 +1 Cha 6 −2 −2|Skills Perception +5, Stealth +4|Senses Darkvision 60 ft.; Passive Perception 15|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Worg", [2]int{342, 342}, "Large Fey, Neutral Evil|AC 13 Initiative +1 (11)|HP 26 (4d10 + 4)|Speed 50 ft.|Str 16 +3 +3 Dex 13 +1 +1 Con 13 +1 +1|Int 7 −2 −2 Wis 11 +0 +0 Cha 8 −1 −1|Skills Perception +4|Senses Darkvision 60 ft.; Passive Perception 14|Languages Goblin, Worg|CR 1/2 (XP 100; PB +2)"},
	{"Wraith", [2]int{342, 343}, "Medium or Small Undead, Neutral Evil|AC 13 Initiative +3 (13)|HP 67 (9d8 + 27)|Speed 5 ft., Fly 60 ft. (hover)|Str 6 −2 −2 Dex 16 +3 +3 Con 16 +3 +3|Int 12 +1 +1 Wis 14 +2 +2 Cha 15 +2 +2|Resistances Acid, Bludgeoning, Cold, Fire, Piercing, Slashing|Immunities Necrotic, Poison; Charmed, Exhaustion, Grappled, Paralyzed, Petrified, Poisoned, Prone, Restrained, Unconscious|Senses Darkvision 60 ft.; Passive Perception 12|Languages Common plus two other languages|CR 5 (XP 1,800; PB +3)"},
	{"Wyvern", [2]int{343, 343}, "Large Dragon, Unaligned|AC 14 Initiative +0 (10)|HP 127 (15d10 + 45)|Speed 30 ft., Fly 80 ft.|Str 19 +4 +4 Dex 10 +0 +0 Con 16 +3 +3|Int 5 −3 −3 Wis 12 +1 +1 Cha 6 −2 −2|Skills Perception +4|Senses Darkvision 120 ft.; Passive Perception 14|Languages None|CR 6 (XP 2,300; PB +3)"},
	{"Xorn", [2]int{343, 343}, "Medium Elemental, Neutral|AC 19 Initiative +0 (10)|HP 84 (8d8 + 48)|Speed 20 ft., Burrow 20 ft.|Str 17 +3 +3 Dex 10 +0 +0 Con 22 +6 +6|Int 11 +0 +0 Wis 10 +0 +0 Cha 11 +0 +0|Skills Perception +6, Stealth +6|Immunities Poison; Paralyzed, Petrified, Poisoned|Senses Darkvision 60 ft., Tremorsense 60 ft.; Passive Perception 16|Languages Primordial (Terran)|CR 5 (XP 1,800; PB +3)"},
	{"Young Black Dragon", [2]int{264, 264}, "Large Dragon (Chromatic), Chaotic Evil|AC 18 Initiative +5 (15)|HP 127 (15d10 + 45)|Speed 40 ft., Fly 80 ft., Swim 40 ft.|Str 19 +4 +4 Dex 14 +2 +5 Con 17 +3 +3|Int 12 +1 +1 Wis 11 +0 +3 Cha 15 +2 +2|Skills Perception +6, Stealth +5|Immunities Acid|Senses Blindsight 30 ft., Darkvision 120 ft.; Passive Perception 16|Languages Common, Draconic|CR 7 (XP 2,900; PB +3)"},
	{"Young Blue Dragon", [2]int{266, 266}, "Large Dragon (Chromatic), Lawful Evil|AC 18 Initiative +4 (14)|HP 152 (16d10 + 64)|Speed 40 ft., Burrow 20 ft., Fly 80 ft.|Str 21 +5 +5 Dex 10 +0 +4 Con 19 +4 +4|Int 14 +2 +2 Wis 13 +1 +5 Cha 17 +3 +3|Skills Perception +9, Stealth +4|Immunities Lightning|Senses Blindsight 30 ft., Darkvision 120 ft.; Passive Perception 19|Languages Common, Draconic|CR 9 (XP 5,000; PB +4)"},
	{"Young Brass Dragon", [2]int{268, 268}, "Large Dragon (Metallic), Chaotic Good|AC 17 Initiative +3 (13)|HP 110 (13d10 + 39)|Speed 40 ft., Burrow 20 ft., Fly 80 ft.|Str 19 +4 +4 Dex 10 +0 +3 Con 17 +3 +3|Int 12 +1 +1 Wis 11 +0 +3 Cha 15 +2 +2|Skills Perception +6, Persuasion +5, Stealth +3|Immunities Fire|Senses Blindsight 30 ft., Darkvision 120 ft.; Passive Perception 16|Languages Common, Draconic|CR 6 (XP 2,300; PB +3)"},
	{"Young Bronze Dragon", [2]int{270, 270}, "Large Dragon (Metallic), Lawful Good|AC 17 Initiative +3 (13)|HP 142 (15d10 + 60)|Speed 40 ft., Fly 80 ft., Swim 40 ft.|Str 21 +5 +5 Dex 10 +0 +3 Con 19 +4 +4|Int 14 +2 +2 Wis 13 +1 +4 Cha 17 +3 +3|Skills Insight +4, Perception +7, Stealth +3|Immunities Lightning|Senses Blindsight 30 ft., Darkvision 120 ft.; Passive Perception 17|Languages Common, Draconic|CR 8 (XP 3,900; PB +3)"},
}

// Typed values of every trait and action entry, reviewed against the same pages:
// limited uses, attack bonus/reach/range/damage, save ability and DC, target
// shape size/width/range, dice and dependency keys. A ranged attack printed with
// one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text.
var monsters177Entries = map[string][]string{
	"Wight": {
		"trait|trait/sunlight-sensitivity|Sunlight Sensitivity; requirement/in-sunlight, check/disadvantage, attack/disadvantage",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination, multiattack/replace-attack, action/saving-throw",
		"action|necrotic-sword|Necrotic Sword; melee +4 reach 5 range 0/0 6 (1d8 + 2) slashing; effect/slashing-damage 6 (1d8 + 2), effect/necrotic-damage 4 (1d8)",
		"action|necrotic-bow|Necrotic Bow; ranged +4 reach 0 range 150/600 6 (1d8 + 2) piercing; effect/piercing-damage 6 (1d8 + 2), effect/necrotic-damage 4 (1d8)",
		"action|life-drain|Life Drain; con DC 13; creature 0/0/5; save/con, target/creature, effect/necrotic-damage 6 (1d8 + 2), effect/hit-point-maximum-reduction, effect/rise-as-zombie, limit/controlled-zombies",
	},
	"Will-o’-Wisp": {
		"trait|trait/ephemeral|Ephemeral; effect/cannot-wear-or-carry",
		"trait|trait/illumination|Illumination; effect/shed-light",
		"trait|trait/incorporeal-movement|Incorporeal Movement; 5 (1d10); movement/through-creatures-and-objects, effect/end-turn-in-object, effect/force-damage",
		"action|shock|Shock; melee +4 reach 5 range 0/0 11 (2d8 + 2) lightning; effect/lightning-damage 11 (2d8 + 2)",
		"bonus|consume-life|Consume Life; con DC 10; creature 0/0/5; save/con, target/creature, requirement/zero-hit-points, effect/target-dies, effect/regain-hit-points 10 (3d6)",
		"bonus|vanish|Vanish; condition/invisible, requirement/concentration, effect/concentration-ends-early",
	},
	"Winter Wolf": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage-with-nearby-ally, requirement/ally-not-incapacitated",
		"action|bite|Bite; melee +6 reach 5 range 0/0 11 (2d6 + 4) piercing; effect/piercing-damage 11 (2d6 + 4), requirement/target-size, condition/prone",
		"action|cold-breath|Cold Breath (Recharge 5–6); con DC 12; cone 15/0/0; save/con, area/cone, limit/recharge, effect/cold-damage 18 (4d8), save/half-damage",
	},
	"Wolf": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage-with-nearby-ally, requirement/ally-not-incapacitated",
		"action|bite|Bite; melee +4 reach 5 range 0/0 5 (1d6 + 2) piercing; effect/piercing-damage 5 (1d6 + 2), requirement/target-size, condition/prone",
	},
	"Worg": {
		"action|bite|Bite; melee +5 reach 5 range 0/0 7 (1d8 + 3) piercing; effect/piercing-damage 7 (1d8 + 3), attack/advantage-on-next-attack",
	},
	"Wraith": {
		"trait|trait/incorporeal-movement|Incorporeal Movement; 5 (1d10); movement/through-creatures-and-objects, effect/end-turn-in-object, effect/force-damage",
		"trait|trait/sunlight-sensitivity|Sunlight Sensitivity; requirement/in-sunlight, check/disadvantage, attack/disadvantage",
		"action|life-drain|Life Drain; melee +6 reach 5 range 0/0 21 (4d8 + 3) necrotic; effect/necrotic-damage 21 (4d8 + 3), requirement/target-creature, effect/hit-point-maximum-reduction",
		"action|create-specter|Create Specter; target/humanoid-corpse, requirement/recent-corpse, effect/raise-specter, effect/specter-control, limit/controlled-specters",
	},
	"Wyvern": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +7 reach 5 range 0/0 13 (2d8 + 4) piercing; effect/piercing-damage 13 (2d8 + 4)",
		"action|sting|Sting; melee +7 reach 10 range 0/0 11 (2d6 + 4) piercing; effect/piercing-damage 11 (2d6 + 4), effect/poison-damage 24 (7d6), condition/poisoned",
	},
	"Xorn": {
		"trait|trait/earth-glide|Earth Glide; movement/burrow-through-earth-and-stone, effect/no-disturbance",
		"trait|trait/treasure-sense|Treasure Sense; sense/locate-precious-metals-and-stones",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +6 reach 5 range 0/0 17 (4d6 + 3) piercing; effect/piercing-damage 17 (4d6 + 3)",
		"action|claw|Claw; melee +6 reach 5 range 0/0 8 (1d10 + 3) slashing; effect/slashing-damage 8 (1d10 + 3)",
		"bonus|charge|Charge; movement/speed-toward-enemy",
	},
	"Young Black Dragon": {
		"trait|trait/amphibious|Amphibious; ",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +7 reach 10 range 0/0 9 (2d4 + 4) slashing; effect/slashing-damage 9 (2d4 + 4), effect/acid-damage 3 (1d6)",
		"action|acid-breath|Acid Breath (Recharge 5–6); dex DC 14; line 30/5/0; save/dex, area/line, limit/recharge, effect/acid-damage 49 (14d6), save/half-damage",
	},
	"Young Blue Dragon": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +9 reach 10 range 0/0 12 (2d6 + 5) slashing; effect/slashing-damage 12 (2d6 + 5), effect/lightning-damage 5 (1d10)",
		"action|lightning-breath|Lightning Breath (Recharge 5–6); dex DC 16; line 60/5/0; save/dex, area/line, limit/recharge, effect/lightning-damage 55 (10d10), save/half-damage",
	},
	"Young Brass Dragon": {
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, action/saving-throw",
		"action|rend|Rend; melee +7 reach 10 range 0/0 15 (2d10 + 4) slashing; effect/slashing-damage 15 (2d10 + 4)",
		"action|fire-breath|Fire Breath (Recharge 5–6); dex DC 14; line 40/5/0; save/dex, area/line, limit/recharge, effect/fire-damage 38 (11d6), save/half-damage",
		"action|sleep-breath|Sleep Breath; con DC 14; cone 30/0/0; save/con, area/cone, condition/incapacitated, save/repeat, condition/unconscious, save/second-failure, effect/end-on-damage-or-wake",
	},
	"Young Bronze Dragon": {
		"trait|trait/amphibious|Amphibious; ",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, action/saving-throw",
		"action|rend|Rend; melee +8 reach 10 range 0/0 16 (2d10 + 5) slashing; effect/slashing-damage 16 (2d10 + 5)",
		"action|lightning-breath|Lightning Breath (Recharge 5–6); dex DC 15; line 60/5/0; save/dex, area/line, limit/recharge, effect/lightning-damage 49 (9d10), save/half-damage",
		"action|repulsion-breath|Repulsion Breath; str DC 15; cone 30/0/0; save/str, area/cone, effect/push, condition/prone",
	},
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters177SourceValues(t *testing.T) {
	c := loadMonsters177(t)
	if len(c.Monsters) != 12 || len(monsters177Headers) != 12 || len(monsters177Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #177") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters177Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader176(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries164(m); !slices.Equal(got, monsters177Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters177Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.VulnerabilityNote != "" || m.ResistanceNote != "" || m.TelepathyNote != "" || m.SwarmOf != "" || m.SpeedNote != "" || m.ConditionalSpeedNote != "" || m.ConditionImmunityNote != "" || m.Telepathy != 0 || m.LairXP != 0 || m.Legendary != nil || m.LegendaryActions != nil || m.Senses.Notes != nil {
				t.Fatal("source lists no vulnerabilities, vulnerability note, resistance note, telepathy note, swarm line, speed note, conditional speed note, condition immunity note, telepathy, lair XP, legendary actions, sense notes")
			}
		})
		for _, line := range monsters177Entries[want.name] {
			if !strings.HasPrefix(line, "legendary-uses|") {
				entries++
			}
		}
	}
	// Bold-italic entry headings counted on PDF pp. 264, 266, 268, 270, 341–343 and 364.
	if entries != 46 {
		t.Fatalf("reviewed %d entries, want 46", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters177MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 177 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters177(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Animals appendix blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("monsters monsters monsters animals monsters monsters monsters monsters monsters monsters monsters monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #177 cites the page each block starts on, in this order.
	names := []string{"Wight", "Will-o’-Wisp", "Winter Wolf", "Wolf", "Worg", "Wraith", "Wyvern", "Xorn", "Young Black Dragon", "Young Blue Dragon", "Young Brass Dragon", "Young Bronze Dragon"}
	pages := []int{341, 341, 342, 364, 342, 342, 343, 343, 264, 266, 268, 270}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic for these twelve blocks; see monsterConsistency170.
func TestMonsters177DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"3": 2, "2": 2, "1/4": 2, "1/2": 2, "5": 3, "6": 3, "7": 3, "9": 4, "8": 3}
	attacks := map[string]attackAbility170{
		"Will-o’-Wisp/Shock": {"wis", true},
	}
	for _, m := range loadMonsters177(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			monsterConsistency170(t, m, pb, attacks)
		})
	}
}

// Every listed key is used and required, every effect quotes its own entry, and
// every dice expression is typed. Removing any one key invalidates the record.
func TestMonsters177DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters177(t).Monsters {
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
func TestMonsters177RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters177(t).Monsters {
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

func TestMonsters177MalformedFixtures(t *testing.T) {
	tests := []struct {
		name, record string
		mutate       func(*Monster)
	}{
		{"missing languages list", "Wolf", func(m *Monster) {
			m.Languages = nil
		}},
		{"animal or monster block in another family", "Wolf", func(m *Monster) {
			m.ID = strings.Replace(strings.Replace(m.ID, "/animals/", "/spells/", 1), "/monsters/", "/spells/", 1)
		}},
		{"XP off CR", "Wolf", func(m *Monster) {
			m.XP++
		}},
		{"save off modifier", "Wolf", func(m *Monster) {
			m.Abilities["wis"] = MonsterAbility{m.Abilities["wis"].Score, m.Abilities["wis"].Modifier, m.Abilities["wis"].Modifier + 11}
		}},
		{"modifier", "Wight", func(m *Monster) {
			m.Abilities["str"] = MonsterAbility{m.Abilities["str"].Score, m.Abilities["str"].Modifier + 1, m.Abilities["str"].Modifier + 1}
		}},
		{"HP average", "Wight", func(m *Monster) {
			m.HP.Average++
		}},
		{"hover without fly speed", "Wraith", func(m *Monster) {
			delete(m.Speeds, "fly")
		}},
		{"recharge", "Winter Wolf", func(m *Monster) {
			m.Actions[1].Uses.Recharge = 7
		}},
		{"line without width", "Young Black Dragon", func(m *Monster) {
			m.Actions[2].Target.Width = 0
		}},
		{"save DC off entry text", "Wight", func(m *Monster) {
			m.Actions[3].Save.DC = 18
		}},
		{"trait dice off size", "Wraith", func(m *Monster) {
			m.Traits[0].Dice.Average++
		}},
		{"empty bonus actions", "Wolf", func(m *Monster) {
			m.BonusActions = []MonsterAction{}
		}},
		{"empty reactions", "Wolf", func(m *Monster) {
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
		c := loadMonsters177(t)
		c.Monsters[1].ID = c.Monsters[2].ID
		rejected(t, c, true)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters177(t)
			tt.mutate(record(t, &c, tt.record))
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			rejected(t, c, !strings.HasPrefix(tt.name, "empty "))
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_177.json")
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
		{`"recharge":5`, `"recharge":null`},
		{`"per_day":0,`, ``},
		{`"dc":13`, `"dc":13,"dc":22`},
		{`"width":0,`, ``},
		{`"darkvision":60`, `"darkvision":null`},
		{`"shape":"creature",`, `"shape":"creature","height":5,`},
		{`"hover":true,`, `"hover":true,"hover":false,`},
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
// to 26 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters177ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters177(t)
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
	for _, load := range []func(*testing.T) MonsterCatalogue{loadMonsters158, loadMonsters159, loadMonsters160, loadMonsters161, loadMonsters162, loadMonsters163, loadMonsters164, loadMonsters165, loadMonsters166, loadMonsters167, loadMonsters168, loadMonsters169, loadMonsters170, loadMonsters171, loadMonsters172, loadMonsters173, loadMonsters174, loadMonsters175, loadMonsters176} {
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
	for _, id := range []string{"bite", "claw", "life-drain", "necrotic-bow", "necrotic-sword", "rend", "shock", "sting"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
