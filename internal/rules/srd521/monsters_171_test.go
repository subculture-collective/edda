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

func loadMonsters171(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_171.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #171 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Remorhaz’s, the Rust Monster’s and the Satyr’s blocks run onto the next
// page.
var monsters171Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Reef Shark", [2]int{360, 360}, "Medium Beast, Unaligned|AC 12 Initiative +2 (12)|HP 22 (4d8 + 4)|Speed 5 ft., Swim 30 ft.|Str 14 +2 +2 Dex 15 +2 +2 Con 13 +1 +1|Int 1 −5 −5 Wis 10 +0 +0 Cha 4 −3 −3|Skills Perception +2|Senses Blindsight 30 ft.; Passive Perception 12|Languages None|CR 1/2 (XP 100; PB +2)"},
	{"Remorhaz", [2]int{319, 320}, "Huge Monstrosity, Unaligned|AC 17 Initiative +5 (15)|HP 195 (17d12 + 85)|Speed 40 ft., Burrow 30 ft.|Str 24 +7 +7 Dex 13 +1 +1 Con 21 +5 +5|Int 4 −3 −3 Wis 10 +0 +0 Cha 5 −3 −3|Immunities Cold, Fire|Senses Darkvision 60 ft., Tremorsense 60 ft.; Passive Perception 10|Languages None|CR 11 (XP 7,200; PB +4)"},
	{"Rhinoceros", [2]int{360, 360}, "Large Beast, Unaligned|AC 13 Initiative −1 (9)|HP 45 (6d10 + 12)|Speed 40 ft.|Str 21 +5 +5 Dex 8 −1 −1 Con 15 +2 +2|Int 2 −4 −4 Wis 12 +1 +1 Cha 6 −2 −2|Senses Passive Perception 11|Languages None|CR 2 (XP 450; PB +2)"},
	{"Riding Horse", [2]int{360, 360}, "Large Beast, Unaligned|AC 11 Initiative +1 (11)|HP 13 (2d10 + 2)|Speed 60 ft.|Str 16 +3 +3 Dex 13 +1 +1 Con 12 +1 +1|Int 2 −4 −4 Wis 11 +0 +0 Cha 7 −2 −2|Senses Passive Perception 10|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Roc", [2]int{320, 320}, "Gargantuan Monstrosity, Unaligned|AC 15 Initiative +8 (18)|HP 248 (16d20 + 80)|Speed 20 ft., Fly 120 ft.|Str 28 +9 +9 Dex 10 +0 +4 Con 20 +5 +5|Int 3 −4 −4 Wis 10 +0 +4 Cha 9 −1 −1|Skills Perception +8|Senses Passive Perception 18|Languages None|CR 11 (XP 7,200; PB +4)"},
	{"Roper", [2]int{320, 320}, "Large Aberration, Neutral Evil|AC 20 Initiative +5 (15)|HP 93 (11d10 + 33)|Speed 10 ft., Climb 20 ft.|Str 18 +4 +4 Dex 8 −1 −1 Con 17 +3 +3|Int 7 −2 −2 Wis 16 +3 +3 Cha 6 −2 −2|Skills Perception +6, Stealth +5|Senses Darkvision 60 ft.; Passive Perception 16|Languages None|CR 5 (XP 1,800; PB +3)"},
	{"Rust Monster", [2]int{320, 321}, "Medium Monstrosity, Unaligned|AC 14 Initiative +1 (11)|HP 33 (6d8 + 6)|Speed 40 ft.|Str 13 +1 +1 Dex 12 +1 +1 Con 13 +1 +1|Int 2 −4 −4 Wis 13 +1 +1 Cha 6 −2 −2|Senses Darkvision 60 ft.; Passive Perception 11|Languages None|CR 1/2 (XP 100; PB +2)"},
	{"Saber-Toothed Tiger", [2]int{360, 360}, "Large Beast, Unaligned|AC 13 Initiative +3 (13)|HP 52 (7d10 + 14)|Speed 40 ft.|Str 18 +4 +6 Dex 17 +3 +5 Con 15 +2 +2|Int 3 −4 −4 Wis 12 +1 +1 Cha 8 −1 −1|Skills Perception +5, Stealth +7|Senses Darkvision 60 ft.; Passive Perception 15|Languages None|CR 2 (XP 450; PB +2)"},
	{"Sahuagin Warrior", [2]int{321, 321}, "Medium Fiend, Lawful Evil|AC 12 Initiative +0 (10)|HP 22 (4d8 + 4)|Speed 30 ft., Swim 40 ft.|Str 13 +1 +1 Dex 11 +0 +0 Con 12 +1 +1|Int 12 +1 +1 Wis 13 +1 +1 Cha 9 −1 −1|Skills Perception +5|Resistances Acid, Cold|Senses Darkvision 120 ft.; Passive Perception 15|Languages Sahuagin|CR 1/2 (XP 100; PB +2)"},
	{"Salamander", [2]int{321, 321}, "Large Elemental, Neutral Evil|AC 15 Initiative +2 (12)|HP 90 (12d10 + 24)|Speed 30 ft., Climb 30 ft.|Str 18 +4 +4 Dex 14 +2 +2 Con 15 +2 +2|Int 11 +0 +0 Wis 10 +0 +0 Cha 12 +1 +1|Vulnerabilities Cold|Immunities Fire|Senses Darkvision 60 ft.; Passive Perception 10|Languages Primordial (Ignan)|CR 5 (XP 1,800; PB +3)"},
	{"Satyr", [2]int{321, 322}, "Medium Fey, Chaotic Neutral|AC 13 Initiative +3 (13)|HP 31 (7d8)|Speed 40 ft.|Str 12 +1 +1 Dex 16 +3 +3 Con 11 +0 +0|Int 12 +1 +1 Wis 10 +0 +0 Cha 14 +2 +2|Skills Perception +2, Performance +6, Stealth +5|Senses Passive Perception 12|Languages Common, Elvish, Sylvan|CR 1/2 (XP 100; PB +2)"},
	{"Scorpion", [2]int{360, 360}, "Tiny Beast, Unaligned|AC 11 Initiative +0 (10)|HP 1 (1d4 − 1)|Speed 10 ft.|Str 2 −4 −4 Dex 11 +0 +0 Con 8 −1 −1|Int 1 −5 −5 Wis 8 −1 −1 Cha 2 −4 −4|Senses Blindsight 10 ft.; Passive Perception 9|Languages None|CR 0 (XP 10; PB +2)"},
}

// Typed values of every trait and action entry, reviewed against the same pages:
// limited uses, attack bonus/reach/range/damage, save ability and DC, target
// shape size/width/range, dice and dependency keys. A ranged attack printed with
// one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text. The Roper's Tentacle prints an attack roll whose hit deals no damage (key
// action/melee-attack-without-damage), the Salamander's Flame Spear is a "Melee
// or Ranged Attack Roll", and the Rhinoceros's Gore prints "If target is" without
// "the".
var monsters171Entries = map[string][]string{
	"Reef Shark": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage-with-nearby-ally, requirement/ally-not-incapacitated",
		"trait|trait/water-breathing|Water Breathing; effect/breathe-only-underwater",
		"action|bite|Bite; melee +4 reach 5 range 0/0 7 (2d4 + 2) piercing; effect/piercing-damage 7 (2d4 + 2)",
	},
	"Remorhaz": {
		"trait|trait/heat-aura|Heat Aura; 16 (3d10); area/emanation, target/creatures-in-area, effect/fire-damage, effect/end-of-turn-damage",
		"action|bite|Bite; melee +11 reach 10 range 0/0 18 (2d10 + 7) piercing; effect/piercing-damage 18 (2d10 + 7), effect/fire-damage 14 (4d6), requirement/target-size, condition/grappled, condition/restrained",
		"bonus|swallow|Swallow; str DC 19; creature 0/0/0; save/str, target/creature, requirement/target-size, requirement/grappled-creature, limit/two-swallowed-creatures, effect/swallow, condition/blinded, condition/restrained, effect/total-cover, effect/acid-damage 10 (3d6), effect/fire-damage 10 (3d6), effect/regurgitation-save, condition/prone, effect/escape-corpse",
	},
	"Rhinoceros": {
		"action|gore|Gore; melee +7 reach 5 range 0/0 14 (2d8 + 5) piercing; effect/piercing-damage 14 (2d8 + 5), requirement/target-size, requirement/straight-line-movement, effect/piercing-damage 9 (2d8), condition/prone",
	},
	"Riding Horse": {
		"action|hooves|Hooves; melee +5 reach 5 range 0/0 7 (1d8 + 3) bludgeoning; effect/bludgeoning-damage 7 (1d8 + 3)",
	},
	"Roc": {
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack",
		"action|beak|Beak; melee +13 reach 10 range 0/0 28 (3d12 + 9) piercing; effect/piercing-damage 28 (3d12 + 9)",
		"action|talons|Talons; melee +13 reach 5 range 0/0 23 (4d6 + 9) slashing; effect/slashing-damage 23 (4d6 + 9), requirement/target-size, condition/grappled, condition/restrained",
		"bonus|swoop|Swoop (Recharge 5–6); limit/recharge, requirement/grappled-creature, movement/half-fly-speed, effect/no-opportunity-attacks, effect/drop-creature",
	},
	"Roper": {
		"trait|trait/spider-climb|Spider Climb; effect/climb-without-check",
		"action|multiattack|Multiattack; action/melee-attack, action/special",
		"action|bite|Bite; melee +7 reach 5 range 0/0 17 (3d8 + 4) piercing; effect/piercing-damage 17 (3d8 + 4)",
		"action|tentacle|Tentacle; attack/roll-without-damage, condition/grappled, limit/grapple-from-one-of-six-tentacles, condition/poisoned, effect/destructible-tentacle, effect/tentacle-damage-not-transferred, effect/tentacle-regrows",
		"action|reel|Reel; target/grappled-creatures, effect/pull",
	},
	"Rust Monster": {
		"trait|trait/iron-scent|Iron Scent; sense/locate-ferrous-metal",
		"action|multiattack|Multiattack; action/melee-attack, action/saving-throw",
		"action|bite|Bite; melee +3 reach 5 range 0/0 5 (1d8 + 1) piercing; effect/piercing-damage 5 (1d8 + 1)",
		"action|antennae|Antennae; dex DC 11; save/dex, target/worn-or-carried-metal-object, target/object-bearer, effect/object-penalty, effect/armor-destruction, effect/weapon-destruction, effect/penalty-removed-by-mending",
		"action|destroy-metal|Destroy Metal; target/unattended-metal-object, effect/destroy-metal",
		"reaction|reflexive-antennae|Reflexive Antennae; reaction/trigger, action/saving-throw",
	},
	"Saber-Toothed Tiger": {
		"trait|trait/running-leap|Running Leap; movement/running-long-jump",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +6 reach 5 range 0/0 11 (2d6 + 4) slashing; effect/slashing-damage 11 (2d6 + 4)",
		"bonus|nimble-escape|Nimble Escape; action/disengage, action/hide",
	},
	"Sahuagin Warrior": {
		"trait|trait/blood-frenzy|Blood Frenzy; attack/advantage-against-target-missing-hit-points",
		"trait|trait/limited-amphibiousness|Limited Amphibiousness; effect/breathe-air-and-water, requirement/submerge-every-four-hours",
		"trait|trait/shark-telepathy|Shark Telepathy; effect/control-sharks",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|claw|Claw; melee +3 reach 5 range 0/0 4 (1d6 + 1) slashing; effect/slashing-damage 4 (1d6 + 1)",
		"bonus|aquatic-charge|Aquatic Charge; movement/swim-toward-enemy",
	},
	"Salamander": {
		"trait|trait/fire-aura|Fire Aura; 7 (2d6); area/emanation, target/chosen-creatures-in-area, effect/fire-damage, effect/end-of-turn-damage",
		"action|multiattack|Multiattack; action/melee-or-ranged-attack, multiattack/replace-attack, action/saving-throw",
		"action|flame-spear|Flame Spear; effect/piercing-damage 13 (2d8 + 4), effect/fire-damage 7 (2d6), effect/return-after-ranged-attack",
		"action|constrict|Constrict; str DC 15; creature 0/0/10; save/str, target/creature, requirement/target-size, effect/bludgeoning-damage 11 (2d6 + 4), effect/fire-damage 7 (2d6), condition/grappled, condition/restrained",
	},
	"Satyr": {
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|hooves|Hooves; melee +5 reach 5 range 0/0 5 (1d4 + 3) bludgeoning; effect/bludgeoning-damage 5 (1d4 + 3), requirement/target-size, effect/push",
		"action|mockery|Mockery; wis DC 12; creature 0/0/90; save/wis, target/creature, effect/psychic-damage 5 (1d6 + 2)",
	},
	"Scorpion": {
		"action|sting|Sting; melee +2 reach 5 range 0/0 1 piercing; effect/piercing-damage [damage/fixed], effect/poison-damage 3 (1d6)",
	},
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters171SourceValues(t *testing.T) {
	c := loadMonsters171(t)
	if len(c.Monsters) != 12 || len(monsters171Headers) != 12 || len(monsters171Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #171") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters171Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader170(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries164(m); !slices.Equal(got, monsters171Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters171Entries[want.name])
			}
			if m.VulnerabilityNote != "" || m.ResistanceNote != "" || m.ConditionImmunityNote != "" || m.TelepathyNote != "" || m.Senses.Notes != nil || m.SwarmOf != "" || m.SpeedNote != "" || m.LairXP != 0 || m.Legendary != nil || m.LegendaryActions != nil {
				t.Fatal("source lists no vulnerability note, resistance note, condition immunity note, telepathy note, sense notes, swarm line, speed note, lair XP, legendary actions")
			}
		})
		for _, line := range monsters171Entries[want.name] {
			if !strings.HasPrefix(line, "legendary-uses|") {
				entries++
			}
		}
	}
	// Bold-italic entry headings counted on PDF pp. 319–322 and 360.
	if entries != 41 {
		t.Fatalf("reviewed %d entries, want 41", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters171MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 171 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters171(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Animals appendix blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("animals monsters animals animals monsters monsters monsters animals monsters monsters monsters animals")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #171 cites the page each block starts on, in this order.
	names := []string{"Reef Shark", "Remorhaz", "Rhinoceros", "Riding Horse", "Roc", "Roper", "Rust Monster", "Saber-Toothed Tiger", "Sahuagin Warrior", "Salamander", "Satyr", "Scorpion"}
	pages := []int{360, 319, 360, 360, 320, 320, 320, 360, 321, 321, 321, 360}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic for these twelve blocks; see monsterConsistency170.
func TestMonsters171DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"1/2": 2, "11": 4, "2": 2, "1/4": 2, "5": 3, "0": 2}
	for _, m := range loadMonsters171(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			monsterConsistency170(t, m, pb, nil)
		})
	}
}

// Every listed key is used and required, every effect quotes its own entry, and
// every dice expression is typed. Removing any one key invalidates the record.
func TestMonsters171DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters171(t).Monsters {
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
					if first.Dice != nil || !slices.Equal(first.Requires, []string{"damage/fixed"}) || !regexp.MustCompile(fmt.Sprintf(`Hit: %d %s damage(?:[.,]| plus )`, a.Attack.Damage.Average, monsterTitle151(a.Attack.DamageType))).MatchString(a.Text) {
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
func TestMonsters171RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters171(t).Monsters {
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

func TestMonsters171MalformedFixtures(t *testing.T) {
	tests := []struct {
		name, record string
		mutate       func(*Monster)
	}{
		{"missing languages list", "Scorpion", func(m *Monster) { m.Languages = nil }},
		{"animal or monster block in another family", "Reef Shark", func(m *Monster) {
			m.ID = strings.Replace(strings.Replace(m.ID, "/animals/", "/spells/", 1), "/monsters/", "/spells/", 1)
		}},
		{"XP off CR", "Roc", func(m *Monster) { m.XP++ }},
		{"save off modifier", "Roc", func(m *Monster) {
			m.Abilities["wis"] = MonsterAbility{m.Abilities["wis"].Score, m.Abilities["wis"].Modifier, m.Abilities["wis"].Modifier + 11}
		}},
		{"modifier", "Roc", func(m *Monster) {
			m.Abilities["str"] = MonsterAbility{m.Abilities["str"].Score, m.Abilities["str"].Modifier + 1, m.Abilities["str"].Modifier + 1}
		}},
		{"HP average", "Roc", func(m *Monster) { m.HP.Average++ }},
		{"recharge", "Roc", func(m *Monster) { m.BonusActions[0].Uses.Recharge = 7 }},
		{"save DC off entry text", "Satyr", func(m *Monster) { m.Actions[1].Save.DC = 13 }},
		{"creature target with a size", "Satyr", func(m *Monster) { m.Actions[1].Target.Size = 5 }},
		{"trait dice average", "Remorhaz", func(m *Monster) { m.Traits[0].Dice.Average++ }},
		{"melee attack without reach", "Roc", func(m *Monster) { m.Actions[1].Attack.Reach = 0 }},
		{"untyped attack under the weapon key", "Roper", func(m *Monster) { m.Actions[2].Key = "action/melee-attack" }},
		{"fixed damage of zero", "Scorpion", func(m *Monster) { m.Actions[0].Attack.Damage.Average = 0 }},
		{"empty bonus actions", "Rhinoceros", func(m *Monster) { m.BonusActions = []MonsterAction{} }},
		{"empty reactions", "Rhinoceros", func(m *Monster) { m.Reactions = []MonsterAction{} }},
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
		c := loadMonsters171(t)
		c.Monsters[1].ID = c.Monsters[2].ID
		rejected(t, c, true)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters171(t)
			tt.mutate(record(t, &c, tt.record))
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			rejected(t, c, !strings.HasPrefix(tt.name, "empty "))
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_171.json")
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
		{`"dc":19`, `"dc":19,"dc":22`},
		{`"width":0,`, ``},
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
// to 20 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters171ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters171(t)
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
	for _, load := range []func(*testing.T) MonsterCatalogue{loadMonsters158, loadMonsters159, loadMonsters160, loadMonsters161, loadMonsters162, loadMonsters163, loadMonsters164, loadMonsters165, loadMonsters166, loadMonsters167, loadMonsters168, loadMonsters169, loadMonsters170} {
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
	for _, id := range []string{"beak", "bite", "claw", "gore", "hooves", "rend", "sting", "talons"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
