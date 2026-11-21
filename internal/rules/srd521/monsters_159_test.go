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

func loadMonsters159(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_159.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #159 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Fire Elemental's block starts on p. 284, the Frost Giant's on
// p. 285 and the Frog's on p. 348; each runs onto the next page.
var monsters159Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Efreeti", [2]int{283, 283}, "Large Elemental (Genie), Neutral|AC 17 Initiative +1 (11)|HP 212 (17d10 + 119)|Speed 40 ft., Fly 60 ft. (hover)|Str 22 +6 +6 Dex 12 +1 +1 Con 24 +7 +7|Int 16 +3 +3 Wis 15 +2 +6 Cha 19 +4 +8|Immunities Fire|Senses Darkvision 120 ft.; Passive Perception 12|Languages Primordial (Ignan)|CR 11 (XP 7,200; PB +4)"},
	{"Elephant", [2]int{348, 348}, "Huge Beast, Unaligned|AC 12 Initiative −1 (9)|HP 76 (8d12 + 24)|Speed 40 ft.|Str 22 +6 +6 Dex 9 −1 −1 Con 17 +3 +3|Int 3 −4 −4 Wis 11 +0 +0 Cha 6 −2 −2|Senses Passive Perception 10|Languages None|CR 4 (XP 1,100; PB +2)"},
	{"Elk", [2]int{348, 348}, "Large Beast, Unaligned|AC 10 Initiative +0 (10)|HP 11 (2d10)|Speed 50 ft.|Str 16 +3 +3 Dex 10 +0 +0 Con 11 +0 +0|Int 2 −4 −4 Wis 10 +0 +0 Cha 6 −2 −2|Skills Perception +2|Senses Darkvision 60 ft.; Passive Perception 12|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Erinyes", [2]int{283, 283}, "Medium Fiend (Devil), Lawful Evil|AC 18 Initiative +7 (17)|HP 178 (21d8 + 84)|Speed 30 ft., Fly 60 ft.|Str 18 +4 +4 Dex 16 +3 +7 Con 18 +4 +8|Int 14 +2 +2 Wis 14 +2 +2 Cha 18 +4 +8|Skills Perception +6, Persuasion +8|Resistances Cold|Immunities Fire, Poison; Poisoned|Senses Truesight 120 ft.; Passive Perception 16|Languages Infernal; telepathy 120 ft.|CR 12 (XP 8,400; PB +4)"},
	{"Ettercap", [2]int{284, 284}, "Medium Monstrosity, Neutral Evil|AC 13 Initiative +2 (12)|HP 44 (8d8 + 8)|Speed 30 ft., Climb 30 ft.|Str 14 +2 +2 Dex 15 +2 +2 Con 13 +1 +1|Int 7 −2 −2 Wis 12 +1 +1 Cha 8 −1 −1|Skills Perception +3, Stealth +4, Survival +3|Senses Darkvision 60 ft.; Passive Perception 13|Languages None|CR 2 (XP 450; PB +2)"},
	{"Ettin", [2]int{284, 284}, "Large Giant, Chaotic Evil|AC 12 Initiative −1 (9)|HP 85 (10d10 + 30)|Speed 40 ft.|Str 21 +5 +5 Dex 8 −1 −1 Con 17 +3 +3|Int 6 −2 −2 Wis 10 +0 +0 Cha 8 −1 −1|Skills Perception +4|Immunities Blinded, Charmed, Deafened, Frightened, Stunned, Unconscious|Gear Battleaxe, Morningstar|Senses Darkvision 60 ft.; Passive Perception 14|Languages Giant|CR 4 (XP 1,100; PB +2)"},
	{"Fire Elemental", [2]int{284, 285}, "Large Elemental, Neutral|AC 13 Initiative +3 (13)|HP 93 (11d10 + 33)|Speed 50 ft.|Str 10 +0 +0 Dex 17 +3 +3 Con 16 +3 +3|Int 6 −2 −2 Wis 10 +0 +0 Cha 7 −2 −2|Resistances Bludgeoning, Piercing, Slashing|Immunities Fire, Poison; Exhaustion, Grappled, Paralyzed, Petrified, Poisoned, Prone, Restrained, Unconscious|Senses Darkvision 60 ft.; Passive Perception 10|Languages Primordial (Ignan)|CR 5 (XP 1,800; PB +3)"},
	{"Fire Giant", [2]int{285, 285}, "Huge Giant, Lawful Evil|AC 18 Initiative +3 (13)|HP 162 (13d12 + 78)|Speed 30 ft.|Str 25 +7 +7 Dex 9 −1 +3 Con 23 +6 +10|Int 10 +0 +0 Wis 14 +2 +2 Cha 13 +1 +5|Skills Athletics +11, Perception +6|Immunities Fire|Senses Passive Perception 16|Languages Giant|CR 9 (XP 5,000; PB +4)"},
	{"Flesh Golem", [2]int{285, 285}, "Medium Construct, Neutral|AC 9 Initiative −1 (9)|HP 127 (15d8 + 60)|Speed 30 ft.|Str 19 +4 +4 Dex 9 −1 −1 Con 18 +4 +4|Int 6 −2 −2 Wis 10 +0 +0 Cha 5 −3 −3|Immunities Lightning, Poison; Charmed, Exhaustion, Frightened, Paralyzed, Petrified, Poisoned|Senses Darkvision 60 ft.; Passive Perception 10|Languages Understands Common plus one other language but can’t speak|CR 5 (XP 1,800; PB +3)"},
	{"Flying Snake", [2]int{348, 348}, "Tiny Monstrosity, Unaligned|AC 14 Initiative +2 (12)|HP 5 (2d4)|Speed 30 ft., Fly 60 ft., Swim 30 ft.|Str 4 −3 −3 Dex 15 +2 +2 Con 11 +0 +0|Int 2 −4 −4 Wis 12 +1 +1 Cha 5 −3 −3|Senses Blindsight 10 ft.; Passive Perception 11|Languages None|CR 1/8 (XP 25; PB +2)"},
	{"Frog", [2]int{348, 349}, "Tiny Beast, Unaligned|AC 11 Initiative +1 (11)|HP 1 (1d4 − 1)|Speed 20 ft., Swim 20 ft.|Str 1 −5 −5 Dex 13 +1 +1 Con 8 −1 −1|Int 1 −5 −5 Wis 8 −1 −1 Cha 3 −4 −4|Skills Perception +1, Stealth +3|Senses Darkvision 30 ft.; Passive Perception 11|Languages None|CR 0 (XP 10; PB +2)"},
	{"Frost Giant", [2]int{285, 286}, "Huge Giant, Neutral Evil|AC 15 Initiative +2 (12)|HP 149 (13d12 + 65)|Speed 40 ft.|Str 23 +6 +6 Dex 9 −1 −1 Con 21 +5 +8|Int 9 −1 −1 Wis 10 +0 +3 Cha 12 +1 +4|Skills Athletics +9, Perception +3|Immunities Cold|Senses Passive Perception 13|Languages Giant|CR 8 (XP 3,900; PB +3)"},
}

// Typed values of every trait, action, bonus action and reaction, reviewed
// against the same pages: limited uses, attack bonus/reach/range/damage, save
// ability and DC, target shape size/width/range, dice and dependency keys. A
// ranged attack printed with one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text: the Efreeti's "Wall of Fire (level 7 version)" keeps its parenthetical
// and requires spellcasting/higher-level-version, the Erinyes' "(Requires Magic
// Rope)" heading qualifier is the requirement/magic-rope effect (the verifier
// compares the full heading), and the Flying Snake's
// bite is a fixed 1 Piercing damage hit followed by typed Poison damage.
var monsters159Entries = map[string][]string{
	"Efreeti": {
		"trait|trait/elemental-restoration|Elemental Restoration; effect/revive",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"trait|trait/wishes|Wishes; spell/wish, effect/wish-chance, effect/wish-non-genie-beneficiary, effect/wish-no-stress, effect/wish-casting-limit",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|heated-blade|Heated Blade; melee +10 reach 5 range 0/0 13 (2d6 + 6) slashing; effect/slashing-damage 13 (2d6 + 6), effect/fire-damage 13 (2d12)",
		"action|hurl-flame|Hurl Flame; ranged +8 reach 0 range 120/120 24 (7d6) fire; effect/fire-damage 24 (7d6)",
		"action|spellcasting|Spellcasting; cha DC 16 +0; at will: Detect Magic, Elementalism; 1/day: Gaseous Form, Invisibility, Major Image, Plane Shift, Tongues, Wall of Fire (level 7 version); spellcasting/no-material-components, spell/detect-magic [spellcasting/at-will], spell/elementalism [spellcasting/at-will], spell/gaseous-form [limit/per-day], spell/invisibility [limit/per-day], spell/major-image [limit/per-day], spell/plane-shift [limit/per-day], spell/tongues [limit/per-day], spell/wall-of-fire [limit/per-day, spellcasting/higher-level-version]",
	},
	"Elephant": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|gore|Gore; melee +8 reach 5 range 0/0 15 (2d8 + 6) piercing; effect/piercing-damage 15 (2d8 + 6), requirement/target-size, requirement/charge-distance, condition/prone",
		"bonus|trample|Trample; dex DC 16; creature 0/0/5; save/dex, target/creature, requirement/target-prone, effect/bludgeoning-damage 17 (2d10 + 6), save/half-damage",
	},
	"Elk": {
		"action|ram|Ram; melee +5 reach 5 range 0/0 6 (1d6 + 3) bludgeoning; effect/bludgeoning-damage 6 (1d6 + 3), requirement/target-size, requirement/charge-distance, effect/extra-bludgeoning-damage 3 (1d6), condition/prone",
	},
	"Erinyes": {
		"trait|trait/diabolical-restoration|Diabolical Restoration; effect/revive",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"trait|trait/magic-rope|Magic Rope; effect/magic-rope, effect/magic-rope-statistics, effect/magic-rope-destruction, effect/magic-rope-restoration",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/additional-action",
		"action|withering-sword|Withering Sword; melee +8 reach 5 range 0/0 13 (2d8 + 4) slashing; effect/slashing-damage 13 (2d8 + 4), effect/necrotic-damage 11 (2d10)",
		"action|entangling-rope|Entangling Rope; str DC 16; creature 0/0/120; requirement/magic-rope, save/str, target/creature, effect/force-damage 14 (4d6), condition/restrained, effect/restraint-duration",
		"reaction|parry|Parry; reaction/trigger, requirement/holding-weapon, effect/ac-bonus",
	},
	"Ettercap": {
		"trait|trait/spider-climb|Spider Climb; ",
		"trait|trait/web-walker|Web Walker; effect/ignore-web-movement-restrictions, effect/know-web-contact-locations",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +4 reach 5 range 0/0 5 (1d6 + 2) piercing; effect/piercing-damage 5 (1d6 + 2), effect/poison-damage 2 (1d4), condition/poisoned",
		"action|claw|Claw; melee +4 reach 5 range 0/0 7 (2d4 + 2) slashing; effect/slashing-damage 7 (2d4 + 2)",
		"action|web-strand|Web Strand (Recharge 5–6); dex DC 12; creature 0/0/30; limit/recharge, save/dex, target/creature, requirement/target-size, condition/restrained, effect/web-statistics",
		"bonus|reel|Reel; effect/pull, requirement/restrained-by-web-strand",
	},
	"Ettin": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|battleaxe|Battleaxe; melee +7 reach 5 range 0/0 14 (2d8 + 5) slashing; effect/slashing-damage 14 (2d8 + 5), condition/prone, requirement/target-size",
		"action|morningstar|Morningstar; melee +7 reach 5 range 0/0 14 (2d8 + 5) piercing; effect/piercing-damage 14 (2d8 + 5), attack/disadvantage",
	},
	"Fire Elemental": {
		"trait|trait/fire-aura|Fire Aura; emanation 10/0/0; 5 (1d10); area/emanation, effect/fire-damage, effect/start-burning",
		"trait|trait/fire-form|Fire Form; 5 (1d10); movement/squeeze, movement/enter-creature-space, effect/fire-damage",
		"trait|trait/illumination|Illumination; effect/shed-light",
		"trait|trait/water-susceptibility|Water Susceptibility; 3 (1d6); effect/cold-damage, effect/water-movement-damage, effect/water-splash-damage",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|burn|Burn; melee +6 reach 5 range 0/0 10 (2d6 + 3) fire; effect/fire-damage 10 (2d6 + 3), effect/start-burning",
	},
	"Fire Giant": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|flame-sword|Flame Sword; melee +11 reach 10 range 0/0 21 (4d6 + 7) slashing; effect/slashing-damage 21 (4d6 + 7), effect/fire-damage 10 (3d6)",
		"action|hammer-throw|Hammer Throw; ranged +11 reach 0 range 60/240 23 (3d10 + 7) bludgeoning; effect/bludgeoning-damage 23 (3d10 + 7), effect/fire-damage 4 (1d8), effect/push, attack/disadvantage",
	},
	"Flesh Golem": {
		"trait|trait/aversion-to-fire|Aversion to Fire; effect/fire-damage-trigger, attack/disadvantage, check/disadvantage",
		"trait|trait/berserk|Berserk; requirement/bloodied, effect/berserk-roll, effect/berserk, effect/berserk-calm, check/cha-persuasion",
		"trait|trait/immutable-form|Immutable Form; effect/no-shape-shift",
		"trait|trait/lightning-absorption|Lightning Absorption; effect/absorb-lightning-damage",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|slam|Slam; melee +7 reach 5 range 0/0 13 (2d8 + 4) bludgeoning; effect/bludgeoning-damage 13 (2d8 + 4), effect/lightning-damage 4 (1d8)",
	},
	"Flying Snake": {
		"trait|trait/flyby|Flyby; effect/no-opportunity-attack-flying-out",
		"action|bite|Bite; melee +4 reach 5 range 0/0 1 piercing; effect/piercing-damage [damage/fixed], effect/poison-damage 5 (2d4)",
	},
	"Frog": {
		"trait|trait/amphibious|Amphibious; effect/breathe-air-and-water",
		"trait|trait/standing-leap|Standing Leap; effect/jump-distance-limits, effect/jump-without-running-start",
		"action|bite|Bite; melee +3 reach 5 range 0/0 1 piercing; effect/piercing-damage [damage/fixed]",
	},
	"Frost Giant": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|frost-axe|Frost Axe; melee +9 reach 10 range 0/0 19 (2d12 + 6) slashing; effect/slashing-damage 19 (2d12 + 6), effect/cold-damage 9 (2d8)",
		"action|great-bow|Great Bow; ranged +9 reach 0 range 150/600 17 (2d10 + 6) piercing; effect/piercing-damage 17 (2d10 + 6), effect/cold-damage 7 (2d6), effect/speed-decrease",
		"bonus|war-cry|War Cry (Recharge 5–6); limit/recharge, target/creature-of-choice, effect/temporary-hit-points 16 (2d10 + 5), attack/advantage",
	},
}

// monsterHeader157, monsterEntries157 and monsterDice157 print typed fields in
// the source's notation; batch #157 introduced them for the same record shapes.

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters159SourceValues(t *testing.T) {
	c := loadMonsters159(t)
	if len(c.Monsters) != 12 || len(monsters159Headers) != 12 || len(monsters159Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #159") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters159Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader157(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries157(m); !slices.Equal(got, monsters159Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters159Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.Hover != (m.Name == "Efreeti") || m.LairXP != 0 || m.LegendaryActions != nil || m.Legendary != nil {
				t.Fatal("source lists no vulnerabilities, lair XP or legendary actions, and only the Efreeti hovers")
			}
			if (m.Telepathy != 0) != (m.Name == "Erinyes") || (m.Reactions != nil) != (m.Name == "Erinyes") {
				t.Fatal("only the Erinyes has telepathy and a reaction")
			}
		})
		entries += len(monsters159Entries[want.name])
	}
	// Bold-italic entry headings counted on PDF pp. 283–286 and 348–349.
	if entries != 53 {
		t.Fatalf("reviewed %d entries, want 53", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters159MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 159 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters159(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the four p. 348 animals under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("monsters animals animals monsters monsters monsters monsters monsters monsters animals animals monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #159 cites the page each block starts on, in this order.
	names := []string{"Efreeti", "Elephant", "Elk", "Erinyes", "Ettercap", "Ettin", "Fire Elemental", "Fire Giant", "Flesh Golem", "Flying Snake", "Frog", "Frost Giant"}
	pages := []int{283, 348, 348, 283, 284, 284, 284, 285, 285, 348, 348, 285}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic that holds for these twelve blocks (PDF pp. 254–256):
// XP and PB follow CR, Hit Dice follow size and Constitution, saves and skills
// add PB to the ability modifier (twice PB for the Ettin's Perception), and
// attack rolls add PB to the listed ability. Initiative adds PB or nothing.
func TestMonsters159DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"0": 2, "1/8": 2, "1/4": 2, "2": 2, "4": 2, "5": 3, "8": 3, "9": 4, "11": 4, "12": 4}
	skillAbility := map[string]string{"athletics": "str", "perception": "wis", "persuasion": "cha", "stealth": "dex", "survival": "wis"}
	hitDie := map[string]int{"tiny": 4, "medium": 8, "large": 10, "huge": 12}
	// Attacks that do not use the better of Strength and Dexterity: the Efreeti's
	// Hurl Flame is a Charisma attack whose 7d6 hit adds no ability modifier.
	attackAbility := map[string]string{"Efreeti/hurl-flame": "cha"}
	plainDamage := map[string]bool{"Efreeti/hurl-flame": true}
	for _, m := range loadMonsters159(t).Monsters {
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
				key := m.Name + "/" + a.ID
				if ability, special := attackAbility[key]; special {
					mod = m.Abilities[ability].Modifier
				}
				// A fixed hit adds no ability modifier.
				damage := mod
				if fixedMonsterDamage(a.Attack.Damage) || plainDamage[key] {
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
func TestMonsters159DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters159(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			if m.RequiredMechanics[0] != "monster/stat-block" {
				t.Fatal("stat-block dependency must come first")
			}
			quoted := func(parent string, x MonsterMechanic) {
				if !strings.Contains(parent, strings.TrimSuffix(x.Text, ".")) {
					t.Fatalf("%q is not quoted from its entry", x.Text)
				}
				if x.Dice != nil && !strings.Contains(x.Text, monsterDice157(*x.Dice)) {
					t.Fatalf("%q does not print its typed dice", x.Text)
				}
			}
			for _, trait := range m.Traits {
				if !strings.HasPrefix(trait.Text, trait.Name+". ") || trait.Key != "trait/"+strings.ReplaceAll(strings.ToLower(trait.Name), " ", "-") {
					t.Fatalf("trait %q key or text differs", trait.Name)
				}
				quoted(trait.Text, trait)
			}
			for _, a := range m.allActions() {
				if !strings.HasPrefix(a.Text, a.Name+" ") && !strings.HasPrefix(a.Text, a.Name+".") || a.ID != strings.ReplaceAll(strings.ToLower(a.Name), " ", "-") {
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
					if first.Dice != nil || !slices.Equal(first.Requires, []string{"damage/fixed"}) || !strings.Contains(a.Text, fmt.Sprintf("Hit: %d %s damage", a.Attack.Damage.Average, monsterTitle151(a.Attack.DamageType))) {
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
func TestMonsters159RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters159(t).Monsters {
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

func TestMonsters159MalformedFixtures(t *testing.T) {
	// Indexes: 0 Efreeti, 1 Elephant, 2 Elk, 3 Erinyes, 4 Ettercap, 5 Ettin,
	// 6 Fire Elemental, 7 Fire Giant, 8 Flesh Golem, 9 Flying Snake, 10 Frog,
	// 11 Frost Giant.
	tests := []struct {
		name   string
		mutate func(*MonsterCatalogue)
	}{
		{"duplicate ID", func(c *MonsterCatalogue) { c.Monsters[1].ID = c.Monsters[2].ID }},
		{"animal block in another family", func(c *MonsterCatalogue) { c.Monsters[1].ID = "srd-5.2.1:content/spells/elephant" }},
		{"missing languages list", func(c *MonsterCatalogue) { c.Monsters[1].Languages = nil }},
		{"fixed damage with sides", func(c *MonsterCatalogue) { c.Monsters[10].Actions[0].Attack.Damage.Sides = 4 }},
		{"fixed damage with modifier", func(c *MonsterCatalogue) { c.Monsters[10].Actions[0].Attack.Damage.Modifier = 1 }},
		{"fixed damage of zero", func(c *MonsterCatalogue) { c.Monsters[10].Actions[0].Attack.Damage.Average = 0 }},
		{"fixed damage without mechanic", func(c *MonsterCatalogue) {
			m := &c.Monsters[10]
			m.Actions[0].Effects[0].Requires = nil
			m.RequiredMechanics = slices.DeleteFunc(m.RequiredMechanics, func(key string) bool { return key == "damage/fixed" })
		}},
		{"fixed effect dice", func(c *MonsterCatalogue) { c.Monsters[10].Actions[0].Effects[0].Dice = &MonsterHP{Average: 1} }},
		{"sense note without sense", func(c *MonsterCatalogue) {
			c.Monsters[0].Senses.Notes = map[string]string{"blindsight": "unimpeded by magical Darkness"}
		}},
		{"blank sense note", func(c *MonsterCatalogue) { c.Monsters[0].Senses.Notes = map[string]string{"darkvision": " "} }},
		{"unlisted sense note", func(c *MonsterCatalogue) {
			c.Monsters[0].Senses.Notes = map[string]string{"darkvision": "impeded by magical Darkness"}
		}},
		{"empty sense notes", func(c *MonsterCatalogue) { c.Monsters[0].Senses.Notes = map[string]string{} }},
		{"melee range", func(c *MonsterCatalogue) { c.Monsters[1].Actions[1].Attack.Range = 5 }},
		{"further dice without first", func(c *MonsterCatalogue) {
			c.Monsters[0].Traits[0].MoreDice = []MonsterKeyedDice{{Key: "effect/future-damage", Dice: MonsterHP{Average: 3, Dice: 1, Sides: 4, Modifier: 1}}}
		}},
		{"XP off CR", func(c *MonsterCatalogue) { c.Monsters[11].XP = 3901 }},
		{"save off modifier", func(c *MonsterCatalogue) { c.Monsters[11].Abilities["wis"] = MonsterAbility{10, 0, 5} }},
		{"modifier", func(c *MonsterCatalogue) { c.Monsters[0].Abilities["str"] = MonsterAbility{22, 5, 5} }},
		{"HP average", func(c *MonsterCatalogue) { c.Monsters[10].HP.Average = 2 }},
		{"hover without flight", func(c *MonsterCatalogue) { c.Monsters[0].Speeds["fly"] = 0 }},
		{"empty bonus actions", func(c *MonsterCatalogue) { c.Monsters[0].BonusActions = []MonsterAction{} }},
		{"empty reactions", func(c *MonsterCatalogue) { c.Monsters[0].Reactions = []MonsterAction{} }},
		{"duplicate action across sections", func(c *MonsterCatalogue) { c.Monsters[4].BonusActions[0].ID = "bite" }},
		{"recharge", func(c *MonsterCatalogue) { c.Monsters[4].Actions[3].Uses.Recharge = 1 }},
		{"daily uses with recharge", func(c *MonsterCatalogue) { c.Monsters[4].Actions[3].Uses.PerDay = 1 }},
		{"save DC off entry text", func(c *MonsterCatalogue) { c.Monsters[4].Actions[3].Save.DC = 16 }},
		{"save ability off entry text", func(c *MonsterCatalogue) { c.Monsters[4].Actions[3].Save.Ability = "str" }},
		{"creature size", func(c *MonsterCatalogue) { c.Monsters[4].Actions[3].Target.Size = 5 }},
		{"emanation width", func(c *MonsterCatalogue) { c.Monsters[6].Traits[0].Target.Width = 5 }},
		{"effect dice", func(c *MonsterCatalogue) { c.Monsters[3].Actions[2].Effects[3].Dice.Average = 15 }},
		{"unlisted action effect", func(c *MonsterCatalogue) { c.Monsters[3].Actions[1].Effects[1].Key = "effect/future" }},
		{"untyped spellcasting", func(c *MonsterCatalogue) { c.Monsters[0].Actions[3].Spellcasting = nil }},
		{"misplaced spellcasting", func(c *MonsterCatalogue) {
			c.Monsters[0].Actions[2].Spellcasting = c.Monsters[0].Actions[3].Spellcasting
		}},
		{"spell ability", func(c *MonsterCatalogue) { c.Monsters[0].Actions[3].Spellcasting.Ability = "" }},
		{"spell DC", func(c *MonsterCatalogue) { c.Monsters[0].Actions[3].Spellcasting.SaveDC = 0 }},
		{"negative spell attack bonus", func(c *MonsterCatalogue) { c.Monsters[0].Actions[3].Spellcasting.AttackBonus = -1 }},
		{"unlisted spell", func(c *MonsterCatalogue) { c.Monsters[0].Actions[3].Spellcasting.AtWill[0] = "Fireball" }},
		{"duplicate spell", func(c *MonsterCatalogue) { c.Monsters[0].Actions[3].Spellcasting.AtWill[0] = "Tongues" }},
		{"daily spell uses", func(c *MonsterCatalogue) { c.Monsters[0].Actions[3].Spellcasting.PerDay[0].Uses = 0 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters159(t)
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
	raw, err := os.ReadFile("testdata/monsters_159.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, edit := range [][2]string{
		{`"dice": 0,`, ``},
		{`"sides": 0,`, `"sides": null,`},
		{`"long_range": 0,`, `"long_range": null,`},
		{`"recharge": 5`, `"recharge": null`},
		{`"per_day": 0,`, ``},
		{`"dc": 16`, `"dc": 16, "dc": 17`},
		{`"width": 0,`, ``},
		{`"special": {}`, `"special": {"darkvision": null}`},
		{`"shape": "creature",`, `"shape": "creature", "height": 5,`},
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
func TestMonsters159ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters159(t)
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
	for _, id := range []string{"burn", "claw", "flame-sword", "frost-axe", "gore", "great-bow", "hammer-throw", "heated-blade", "hurl-flame", "morningstar", "ram", "slam", "withering-sword"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
