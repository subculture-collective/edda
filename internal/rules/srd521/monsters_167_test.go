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

func loadMonsters167(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_167.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #167 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Magmin's, Medusa's and Mimic's blocks run onto the next page.
var monsters167Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Magma Mephit", [2]int{307, 307}, "Small Elemental, Neutral Evil|AC 11 Initiative +1 (11)|HP 18 (4d6 + 4)|Speed 30 ft., Fly 30 ft.|Str 8 −1 −1 Dex 12 +1 +1 Con 12 +1 +1|Int 7 −2 −2 Wis 10 +0 +0 Cha 10 +0 +0|Skills Stealth +3|Vulnerabilities Cold|Immunities Fire, Poison; Exhaustion, Poisoned|Senses Darkvision 60 ft.; Passive Perception 10|Languages Primordial (Ignan, Terran)|CR 1/2 (XP 100; PB +2)"},
	{"Magmin", [2]int{305, 306}, "Small Elemental, Chaotic Neutral|AC 14 Initiative +2 (12)|HP 13 (3d6 + 3)|Speed 30 ft.|Str 7 −2 −2 Dex 15 +2 +2 Con 12 +1 +1|Int 8 −1 −1 Wis 11 +0 +0 Cha 10 +0 +0|Immunities Fire|Senses Darkvision 60 ft.; Passive Perception 10|Languages Primordial (Ignan)|CR 1/2 (XP 100; PB +2)"},
	{"Mammoth", [2]int{357, 357}, "Huge Beast, Unaligned|AC 13 Initiative +2 (12)|HP 126 (11d12 + 55)|Speed 50 ft.|Str 24 +7 +10 Dex 9 −1 −1 Con 21 +5 +8|Int 3 −4 −4 Wis 11 +0 +0 Cha 6 −2 −2|Senses Passive Perception 10|Languages None|CR 6 (XP 2,300; PB +3)"},
	{"Manticore", [2]int{306, 306}, "Large Monstrosity, Lawful Evil|AC 14 Initiative +3 (13)|HP 68 (8d10 + 24)|Speed 30 ft., Fly 50 ft.|Str 17 +3 +3 Dex 16 +3 +3 Con 17 +3 +3|Int 7 −2 −2 Wis 12 +1 +1 Cha 8 −1 −1|Senses Darkvision 60 ft.; Passive Perception 11|Languages Common|CR 3 (XP 700; PB +2)"},
	{"Marilith", [2]int{306, 306}, "Large Fiend (Demon), Chaotic Evil|AC 16 Initiative +10 (20)|HP 220 (21d10 + 105)|Speed 40 ft., Climb 40 ft.|Str 18 +4 +9 Dex 20 +5 +5 Con 20 +5 +10|Int 18 +4 +4 Wis 16 +3 +8 Cha 20 +5 +10|Skills Perception +8|Resistances Cold, Fire, Lightning|Immunities Poison; Poisoned|Senses Truesight 120 ft.; Passive Perception 18|Languages Abyssal; telepathy 120 ft.|CR 16 (XP 15,000; PB +5)"},
	{"Mastiff", [2]int{357, 357}, "Medium Beast, Unaligned|AC 12 Initiative +2 (12)|HP 5 (1d8 + 1)|Speed 40 ft.|Str 13 +1 +1 Dex 14 +2 +2 Con 12 +1 +1|Int 3 −4 −4 Wis 12 +1 +3 Cha 7 −2 −2|Skills Perception +5|Senses Darkvision 60 ft.; Passive Perception 15|Languages None|CR 1/8 (XP 25; PB +2)"},
	{"Medusa", [2]int{306, 307}, "Medium Monstrosity, Lawful Evil|AC 15 Initiative +6 (16)|HP 127 (17d8 + 51)|Speed 30 ft.|Str 10 +0 +0 Dex 17 +3 +3 Con 16 +3 +3|Int 12 +1 +1 Wis 13 +1 +4 Cha 15 +2 +2|Skills Deception +5, Perception +4, Stealth +6|Senses Darkvision 150 ft.; Passive Perception 14|Languages Common plus one other language|CR 6 (XP 2,300; PB +3)"},
	{"Merfolk Skirmisher", [2]int{308, 308}, "Medium Elemental, Neutral|AC 11 Initiative +1 (11)|HP 11 (2d8 + 2)|Speed 10 ft., Swim 40 ft.|Str 10 +0 +0 Dex 13 +1 +1 Con 12 +1 +1|Int 11 +0 +0 Wis 14 +2 +2 Cha 12 +1 +1|Senses Passive Perception 12|Languages Common, Primordial (Aquan)|CR 1/8 (XP 25; PB +2)"},
	{"Merrow", [2]int{308, 308}, "Large Monstrosity, Chaotic Evil|AC 13 Initiative +2 (12)|HP 45 (6d10 + 12)|Speed 10 ft., Swim 40 ft.|Str 18 +4 +4 Dex 15 +2 +2 Con 15 +2 +2|Int 8 −1 −1 Wis 10 +0 +0 Cha 9 −1 −1|Senses Darkvision 60 ft.; Passive Perception 10|Languages Abyssal, Primordial (Aquan)|CR 2 (XP 450; PB +2)"},
	{"Mimic", [2]int{308, 309}, "Medium Monstrosity, Neutral|AC 12 Initiative +3 (13)|HP 58 (9d8 + 18)|Speed 20 ft.|Str 17 +3 +3 Dex 12 +1 +1 Con 15 +2 +2|Int 5 −3 −3 Wis 13 +1 +1 Cha 8 −1 −1|Skills Stealth +5|Immunities Acid; Prone|Senses Darkvision 60 ft.; Passive Perception 11|Languages None|CR 2 (XP 450; PB +2)"},
	{"Minotaur Skeleton", [2]int{326, 326}, "Large Undead, Lawful Evil|AC 12 Initiative +0 (10)|HP 45 (6d10 + 12)|Speed 40 ft.|Str 18 +4 +4 Dex 11 +0 +0 Con 15 +2 +2|Int 6 −2 −2 Wis 8 −1 −1 Cha 5 −3 −3|Vulnerabilities Bludgeoning|Immunities Poison; Exhaustion, Poisoned|Senses Darkvision 60 ft.; Passive Perception 9|Languages Understands Abyssal but can’t speak|CR 2 (XP 450; PB +2)"},
	{"Minotaur of Baphomet", [2]int{309, 309}, "Large Monstrosity, Chaotic Evil|AC 14 Initiative +0 (10)|HP 85 (10d10 + 30)|Speed 40 ft.|Str 18 +4 +4 Dex 11 +0 +0 Con 16 +3 +3|Int 6 −2 −2 Wis 16 +3 +3 Cha 9 −1 −1|Skills Perception +7, Survival +7|Senses Darkvision 60 ft.; Passive Perception 17|Languages Abyssal|CR 3 (XP 700; PB +2)"},
}

// Typed values of every trait and action entry, reviewed against the same pages:
// limited uses, attack bonus/reach/range/damage, save ability and DC, target
// shape size/width/range, dice and dependency keys. A ranged attack printed with
// one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text. The Merfolk Skirmisher's Ocean Spear and the Merrow's Harpoon are "Melee or
// Ranged Attack Roll" entries (key action/melee-or-ranged-attack). The Mimic's larger hit
// against a Grappled target stays a second typed effect.
var monsters167Entries = map[string][]string{
	"Magma Mephit": {
		"trait|trait/death-burst|Death Burst; dex DC 11; emanation 5/0/0; 7 (2d6); trigger/death, save/dex, area/emanation, effect/fire-damage, save/half-damage",
		"action|claw|Claw; melee +3 reach 5 range 0/0 3 (1d4 + 1) slashing; effect/slashing-damage 3 (1d4 + 1), effect/fire-damage 3 (1d6)",
		"action|fire-breath|Fire Breath (Recharge 6); dex DC 11; cone 15/0/0; limit/recharge, save/dex, area/cone, effect/fire-damage 7 (2d6), save/half-damage",
	},
	"Magmin": {
		"trait|trait/death-burst|Death Burst; dex DC 11; emanation 10/0/0; 7 (2d6); trigger/death, save/dex, area/emanation, effect/fire-damage, save/half-damage",
		"action|touch|Touch; melee +4 reach 5 range 0/0 7 (2d4 + 2) fire; effect/fire-damage 7 (2d4 + 2), requirement/creature-or-unattended-flammable-object, effect/start-burning",
		"bonus|ignited-illumination|Ignited Illumination; effect/toggle-flames, effect/shed-light",
	},
	"Mammoth": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|gore|Gore; melee +10 reach 10 range 0/0 18 (2d10 + 7) piercing; effect/piercing-damage 18 (2d10 + 7), requirement/target-size, requirement/charge-distance, condition/prone",
		"bonus|trample|Trample; dex DC 18; creature 0/0/5; save/dex, target/creature, requirement/target-prone, effect/bludgeoning-damage 29 (4d10 + 7), save/half-damage",
	},
	"Manticore": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|rend|Rend; melee +5 reach 5 range 0/0 7 (1d8 + 3) slashing; effect/slashing-damage 7 (1d8 + 3)",
		"action|tail-spike|Tail Spike; ranged +5 reach 0 range 100/200 7 (1d8 + 3) piercing; effect/piercing-damage 7 (1d8 + 3)",
	},
	"Marilith": {
		"trait|trait/demonic-restoration|Demonic Restoration; effect/revive",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"trait|trait/reactive|Reactive; effect/reaction-every-turn",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/additional-action, action/saving-throw",
		"action|pact-blade|Pact Blade; melee +10 reach 5 range 0/0 10 (1d10 + 5) slashing; effect/slashing-damage 10 (1d10 + 5), effect/necrotic-damage 7 (2d6)",
		"action|constrict|Constrict; str DC 17; creature 0/0/5; save/str, target/creature, requirement/target-size, effect/bludgeoning-damage 15 (2d10 + 4), condition/grappled, condition/restrained",
		"bonus|teleport|Teleport (Recharge 5–6); limit/recharge, effect/teleport, requirement/unoccupied-visible-space",
		"reaction|parry|Parry; reaction/trigger, requirement/holding-weapon, effect/ac-bonus",
	},
	"Mastiff": {
		"action|bite|Bite; melee +3 reach 5 range 0/0 4 (1d6 + 1) piercing; effect/piercing-damage 4 (1d6 + 1), requirement/target-size, condition/prone",
	},
	"Medusa": {
		"action|multiattack|Multiattack; action/melee-attack, multiattack/alternative-routine, action/ranged-attack",
		"action|claw|Claw; melee +6 reach 5 range 0/0 10 (2d6 + 3) slashing; effect/slashing-damage 10 (2d6 + 3)",
		"action|snake-hair|Snake Hair; melee +6 reach 5 range 0/0 5 (1d4 + 3) piercing; effect/piercing-damage 5 (1d4 + 3), effect/poison-damage 14 (4d6)",
		"action|poison-ray|Poison Ray; ranged +5 reach 0 range 150/150 11 (2d8 + 2) poison; effect/poison-damage 11 (2d8 + 2)",
		"bonus|petrifying-gaze|Petrifying Gaze (Recharge 5–6); con DC 13; cone 30/0/0; limit/recharge, save/con, area/cone, effect/reflected-gaze, save/first-failure, condition/restrained, save/repeat, save/second-failure, condition/petrified",
	},
	"Merfolk Skirmisher": {
		"trait|trait/amphibious|Amphibious; effect/breathe-air-and-water",
		"action|ocean-spear|Ocean Spear; effect/piercing-damage 3 (1d6), effect/cold-damage 2 (1d4), requirement/target-creature, effect/speed-decrease, effect/weapon-returns",
	},
	"Merrow": {
		"trait|trait/amphibious|Amphibious; effect/breathe-air-and-water",
		"action|multiattack|Multiattack; action/melee-attack, action/melee-or-ranged-attack, multiattack/any-combination",
		"action|bite|Bite; melee +6 reach 5 range 0/0 6 (1d4 + 4) piercing; effect/piercing-damage 6 (1d4 + 4), condition/poisoned",
		"action|claw|Claw; melee +6 reach 5 range 0/0 9 (2d4 + 4) slashing; effect/slashing-damage 9 (2d4 + 4)",
		"action|harpoon|Harpoon; effect/piercing-damage 11 (2d6 + 4), requirement/target-size, effect/pull",
	},
	"Mimic": {
		"trait|trait/adhesive|Adhesive; requirement/object-form, effect/adhere, condition/grappled, requirement/target-size, check/disadvantage",
		"action|bite|Bite; melee +5 reach 5 range 0/0 7 (1d8 + 3) piercing; effect/piercing-damage 7 (1d8 + 3), attack/advantage-against-grappled-target, effect/piercing-damage-against-grappled-target 12 (2d8 + 3), effect/acid-damage 4 (1d8)",
		"action|pseudopod|Pseudopod; melee +5 reach 5 range 0/0 7 (1d8 + 3) bludgeoning; effect/bludgeoning-damage 7 (1d8 + 3), effect/acid-damage 4 (1d8), requirement/target-size, condition/grappled, check/disadvantage",
		"bonus|shape-shift|Shape-Shift; effect/shape-shift, effect/shapechange-form-limits, effect/equipment-not-transformed",
	},
	"Minotaur Skeleton": {
		"action|gore|Gore; melee +6 reach 5 range 0/0 11 (2d6 + 4) piercing; effect/piercing-damage 11 (2d6 + 4), requirement/target-size, requirement/charge-distance, effect/extra-piercing-damage 9 (2d8), condition/prone",
		"action|slam|Slam; melee +6 reach 5 range 0/0 15 (2d10 + 4) bludgeoning; effect/bludgeoning-damage 15 (2d10 + 4)",
	},
	"Minotaur of Baphomet": {
		"action|abyssal-glaive|Abyssal Glaive; melee +6 reach 10 range 0/0 10 (1d12 + 4) slashing; effect/slashing-damage 10 (1d12 + 4), effect/necrotic-damage 10 (3d6)",
		"action|gore|Gore (Recharge 5–6); melee +6 reach 5 range 0/0 18 (4d6 + 4) piercing; effect/piercing-damage 18 (4d6 + 4), limit/recharge, requirement/target-size, requirement/charge-distance, effect/extra-piercing-damage 10 (3d6), condition/prone",
	},
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters167SourceValues(t *testing.T) {
	c := loadMonsters167(t)
	if len(c.Monsters) != 12 || len(monsters167Headers) != 12 || len(monsters167Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #167") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters167Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader164(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries164(m); !slices.Equal(got, monsters167Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters167Entries[want.name])
			}
			if m.ResistanceNote != "" || m.Hover || m.LairXP != 0 || m.Legendary != nil || m.LegendaryActions != nil {
				t.Fatal("source lists no resistance note, hover, lair XP, legendary actions")
			}
		})
		for _, line := range monsters167Entries[want.name] {
			if !strings.HasPrefix(line, "legendary-uses|") {
				entries++
			}
		}
	}
	// Bold-italic entry headings counted on PDF pp. 305–309, 326 and 357.
	if entries != 41 {
		t.Fatalf("reviewed %d entries, want 41", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters167MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 167 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters167(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Animals appendix blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("monsters monsters animals monsters monsters animals monsters monsters monsters monsters monsters monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #167 cites the page each block starts on, in this order.
	names := []string{"Magma Mephit", "Magmin", "Mammoth", "Manticore", "Marilith", "Mastiff", "Medusa", "Merfolk Skirmisher", "Merrow", "Mimic", "Minotaur Skeleton", "Minotaur of Baphomet"}
	pages := []int{307, 305, 357, 306, 306, 357, 306, 308, 308, 308, 326, 309}
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
func TestMonsters167DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"1/2": 2, "6": 3, "3": 2, "16": 5, "1/8": 2, "2": 2}
	skillAbility := map[string]string{"deception": "cha", "perception": "wis", "stealth": "dex", "survival": "wis"}
	hitDie := map[string]int{"small": 6, "medium": 8, "large": 10, "huge": 12}
	// Attacks that do not use the better of Strength and Dexterity.
	attackAbility := map[string]string{"Mastiff/bite": "str", "Medusa/poison-ray": "cha"}
	for _, m := range loadMonsters167(t).Monsters {
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
				if ability, special := attackAbility[m.Name+"/"+a.ID]; special {
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
func TestMonsters167DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters167(t).Monsters {
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
func TestMonsters167RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters167(t).Monsters {
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

func TestMonsters167MalformedFixtures(t *testing.T) {
	tests := []struct {
		name, record string
		mutate       func(*Monster)
	}{
		{"missing languages list", "Magmin", func(m *Monster) { m.Languages = nil }},
		{"animal block in another family", "Magmin", func(m *Monster) {
			m.ID = strings.Replace(strings.Replace(m.ID, "/animals/", "/spells/", 1), "/monsters/", "/spells/", 1)
		}},
		{"XP off CR", "Minotaur of Baphomet", func(m *Monster) { m.XP++ }},
		{"save off modifier", "Minotaur of Baphomet", func(m *Monster) {
			m.Abilities["wis"] = MonsterAbility{m.Abilities["wis"].Score, m.Abilities["wis"].Modifier, m.Abilities["wis"].Modifier + 7}
		}},
		{"modifier", "Magma Mephit", func(m *Monster) {
			m.Abilities["str"] = MonsterAbility{m.Abilities["str"].Score, m.Abilities["str"].Modifier + 1, m.Abilities["str"].Modifier + 1}
		}},
		{"HP average", "Magma Mephit", func(m *Monster) { m.HP.Average++ }},
		{"vulnerability value", "Magma Mephit", func(m *Monster) { m.Vulnerabilities = []string{"frost"} }},
		{"recharge", "Marilith", func(m *Monster) { m.BonusActions[0].Uses.Recharge = 1 }},
		{"trait save DC off entry text", "Magmin", func(m *Monster) { m.Traits[0].Save.DC = 12 }},
		{"emanation width", "Magmin", func(m *Monster) { m.Traits[0].Target.Width = 5 }},
		{"creature size", "Mammoth", func(m *Monster) { m.BonusActions[0].Target.Size = 5 }},
		{"trait dice", "Magma Mephit", func(m *Monster) { m.Traits[0].Dice.Average = 8 }},
		{"effect dice", "Medusa", func(m *Monster) { m.Actions[2].Effects[1].Dice.Average = 15 }},
		{"unlisted action effect", "Mimic", func(m *Monster) { m.Actions[0].Effects[2].Key = "effect/future" }},
		{"empty bonus actions", "Mastiff", func(m *Monster) { m.BonusActions = []MonsterAction{} }},
		{"duplicate action across sections", "Marilith", func(m *Monster) { m.Reactions[0].ID = "constrict" }},
		{"negative attack reach", "Manticore", func(m *Monster) { m.Actions[1].Attack.Reach = -5 }},
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
		c := loadMonsters167(t)
		c.Monsters[1].ID = c.Monsters[2].ID
		rejected(t, c, true)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters167(t)
			tt.mutate(record(t, &c, tt.record))
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			rejected(t, c, !strings.HasPrefix(tt.name, "empty "))
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_167.json")
	if err != nil {
		t.Fatal(err)
	}
	// The fixtures are indented; compact them so each edit matches one field.
	var compact bytes.Buffer
	if err := json.Compact(&compact, raw); err != nil {
		t.Fatal(err)
	}
	for _, edit := range [][2]string{
		{`"range":0,`, ``},
		{`"sides":8,`, `"sides":null,`},
		{`"long_range":0,`, `"long_range":null,`},
		{`"recharge":6`, `"recharge":null`},
		{`"per_day":0,`, ``},
		{`"dc":18`, `"dc":18,"dc":19`},
		{`"width":0,`, ``},
		{`"modifier":7`, `"modifier":null`},
		{`"bonus":10,`, ``},
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
// to 16 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters167ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters167(t)
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
	for _, load := range []func(*testing.T) MonsterCatalogue{loadMonsters158, loadMonsters159, loadMonsters160, loadMonsters161, loadMonsters162, loadMonsters163, loadMonsters164, loadMonsters165, loadMonsters166} {
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
	for _, id := range []string{"abyssal-glaive", "bite", "claw", "gore", "pact-blade", "poison-ray", "pseudopod", "rend", "slam", "snake-hair", "tail-spike", "touch"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
