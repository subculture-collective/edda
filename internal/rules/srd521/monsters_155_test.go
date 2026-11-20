package srd521

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func loadMonsters155(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_155.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #155 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The headers, the entry payloads and the printed-notation helpers
// are shared with the #154 test, which already reads this notation.
var monsters155Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Black Bear", [2]int{345, 345}, "Medium Beast, Unaligned|AC 11 Initiative +1 (11)|HP 19 (3d8 + 6)|Speed 30 ft., Climb 30 ft., Swim 30 ft.|Str 15 +2 +2 Dex 12 +1 +1 Con 14 +2 +2|Int 2 −4 −4 Wis 12 +1 +1 Cha 7 −2 −2|Skills Perception +5|Senses Darkvision 60 ft.; Passive Perception 15|Languages None|CR 1/2 (XP 100; PB +2)"},
	{"Black Dragon Wyrmling", [2]int{263, 264}, "Medium Dragon (Chromatic), Chaotic Evil|AC 17 Initiative +4 (14)|HP 33 (6d8 + 6)|Speed 30 ft., Fly 60 ft., Swim 30 ft.|Str 15 +2 +2 Dex 14 +2 +4 Con 13 +1 +1|Int 10 +0 +0 Wis 11 +0 +2 Cha 13 +1 +1|Skills Perception +4, Stealth +4|Immunities Acid|Senses Blindsight 10 ft., Darkvision 60 ft.; Passive Perception 14|Languages Draconic|CR 2 (XP 450; PB +2)"},
	{"Black Pudding", [2]int{265, 266}, "Large Ooze, Unaligned|AC 7 Initiative −3 (7)|HP 68 (8d10 + 24)|Speed 20 ft., Climb 20 ft.|Str 16 +3 +3 Dex 5 −3 −3 Con 16 +3 +3|Int 1 −5 −5 Wis 6 −2 −2 Cha 1 −5 −5|Immunities Acid, Cold, Lightning, Slashing; Charmed, Deafened, Exhaustion, Frightened, Grappled, Prone, Restrained|Senses Blindsight 60 ft.; Passive Perception 8|Languages None|CR 4 (XP 1,100; PB +2)"},
	{"Blink Dog", [2]int{266, 266}, "Medium Fey, Lawful Good|AC 13 Initiative +3 (13)|HP 22 (4d8 + 4)|Speed 40 ft.|Str 12 +1 +1 Dex 17 +3 +3 Con 12 +1 +1|Int 10 +0 +0 Wis 13 +1 +1 Cha 11 +0 +0|Skills Perception +5, Stealth +5|Senses Darkvision 60 ft.; Passive Perception 15|Languages Blink Dog; understands Elvish and Sylvan but can’t speak them|CR 1/4 (XP 50; PB +2)"},
	{"Blood Hawk", [2]int{345, 346}, "Small Beast, Unaligned|AC 12 Initiative +2 (12)|HP 7 (2d6)|Speed 10 ft., Fly 60 ft.|Str 6 −2 −2 Dex 14 +2 +2 Con 10 +0 +0|Int 3 −4 −4 Wis 14 +2 +2 Cha 5 −3 −3|Skills Perception +6|Senses Passive Perception 16|Languages None|CR 1/8 (XP 25; PB +2)"},
	{"Blue Dragon Wyrmling", [2]int{266, 266}, "Medium Dragon (Chromatic), Lawful Evil|AC 17 Initiative +2 (12)|HP 65 (10d8 + 20)|Speed 30 ft., Burrow 15 ft., Fly 60 ft.|Str 17 +3 +3 Dex 10 +0 +2 Con 15 +2 +2|Int 12 +1 +1 Wis 11 +0 +2 Cha 15 +2 +2|Skills Perception +4, Stealth +2|Immunities Lightning|Senses Blindsight 10 ft., Darkvision 60 ft.; Passive Perception 14|Languages Draconic|CR 3 (XP 700; PB +2)"},
	{"Boar", [2]int{346, 346}, "Medium Beast, Unaligned|AC 11 Initiative +0 (10)|HP 13 (2d8 + 4)|Speed 40 ft.|Str 13 +1 +1 Dex 11 +0 +0 Con 14 +2 +2|Int 2 −4 −4 Wis 9 −1 −1 Cha 5 −3 −3|Senses Passive Perception 9|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Bone Devil", [2]int{267, 268}, "Large Fiend (Devil), Lawful Evil|AC 16 Initiative +7 (17)|HP 161 (17d10 + 68)|Speed 40 ft., Fly 40 ft.|Str 18 +4 +8 Dex 16 +3 +3 Con 18 +4 +4|Int 13 +1 +5 Wis 14 +2 +6 Cha 16 +3 +7|Skills Deception +7, Insight +6|Resistances Cold|Immunities Fire, Poison; Poisoned|Senses Darkvision 120 ft. (unimpeded by magical Darkness); Passive Perception 12|Languages Infernal; telepathy 120 ft.|CR 9 (XP 5,000; PB +4)"},
	{"Brass Dragon Wyrmling", [2]int{268, 268}, "Medium Dragon (Metallic), Chaotic Good|AC 15 Initiative +2 (12)|HP 22 (4d8 + 4)|Speed 30 ft., Burrow 15 ft., Fly 60 ft.|Str 15 +2 +2 Dex 10 +0 +2 Con 13 +1 +1|Int 10 +0 +0 Wis 11 +0 +2 Cha 13 +1 +1|Skills Perception +4, Stealth +2|Immunities Fire|Senses Blindsight 10 ft., Darkvision 60 ft.; Passive Perception 14|Languages Draconic|CR 1 (XP 200; PB +2)"},
	{"Bronze Dragon Wyrmling", [2]int{269, 270}, "Medium Dragon (Metallic), Lawful Good|AC 15 Initiative +2 (12)|HP 39 (6d8 + 12)|Speed 30 ft., Fly 60 ft., Swim 30 ft.|Str 17 +3 +3 Dex 10 +0 +2 Con 15 +2 +2|Int 12 +1 +1 Wis 11 +0 +2 Cha 15 +2 +2|Skills Perception +4, Stealth +2|Immunities Lightning|Senses Blindsight 10 ft., Darkvision 60 ft.; Passive Perception 14|Languages Draconic|CR 2 (XP 450; PB +2)"},
	{"Brown Bear", [2]int{346, 346}, "Large Beast, Unaligned|AC 11 Initiative +1 (11)|HP 22 (3d10 + 6)|Speed 40 ft., Climb 30 ft.|Str 17 +3 +3 Dex 12 +1 +1 Con 15 +2 +2|Int 2 −4 −4 Wis 13 +1 +1 Cha 7 −2 −2|Skills Perception +3|Senses Darkvision 60 ft.; Passive Perception 13|Languages None|CR 1 (XP 200; PB +2)"},
	{"Bugbear Stalker", [2]int{271, 271}, "Medium Fey (Goblinoid), Chaotic Evil|AC 15 Initiative +2 (12)|HP 65 (10d8 + 20)|Speed 30 ft.|Str 17 +3 +3 Dex 14 +2 +2 Con 14 +2 +4|Int 11 +0 +0 Wis 12 +1 +3 Cha 11 +0 +0|Skills Stealth +6, Survival +3|Gear Chain Shirt, Javelins (6), Morningstar|Senses Darkvision 60 ft.; Passive Perception 11|Languages Common, Goblin|CR 3 (XP 700; PB +2)"},
}

// Typed values of every trait, action, bonus action and reaction, reviewed
// against the same pages: limited uses, attack bonus/reach/damage, save ability
// and DC, target shape size/width/range, dice and dependency keys. The Javelin
// is printed as melee or ranged, so it has no typed attack and its dice sit on
// its effect, as for the #153 Archmage.
var monsters155Entries = map[string][]string{
	"Black Bear": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +4 reach 5 range 0/0 5 (1d6 + 2) slashing; effect/slashing-damage 5 (1d6 + 2)",
	},
	"Black Dragon Wyrmling": {
		"trait|trait/amphibious|Amphibious; ",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +4 reach 5 range 0/0 5 (1d6 + 2) slashing; effect/slashing-damage 5 (1d6 + 2), effect/acid-damage 2 (1d4)",
		"action|acid-breath|Acid Breath (Recharge 5–6); dex DC 11; line 15/5/0; save/dex, area/line, limit/recharge, effect/acid-damage 22 (5d8), save/half-damage",
	},
	"Black Pudding": {
		"trait|trait/amorphous|Amorphous; ",
		"trait|trait/corrosive-form|Corrosive Form; 4 (1d8); effect/acid-damage, effect/destroy-nonmagical-ammunition, effect/weapon-attack-penalty, effect/destroy-weapon, effect/eat-through-material",
		"trait|trait/spider-climb|Spider Climb; ",
		"action|dissolving-pseudopod|Dissolving Pseudopod; melee +5 reach 10 range 0/0 17 (4d6 + 3) acid; effect/acid-damage 17 (4d6 + 3), effect/armor-ac-penalty, effect/destroy-armor",
		"reaction|split|Split; reaction/trigger, requirement/size-large-or-medium, requirement/minimum-hit-points, effect/split, effect/smaller-size, effect/divide-hit-points",
	},
	"Blink Dog": {
		"action|bite|Bite; melee +5 reach 5 range 0/0 5 (1d4 + 3) piercing; effect/piercing-damage 5 (1d4 + 3)",
		"bonus|teleport|Teleport (Recharge 4–6); limit/recharge, effect/teleport, requirement/unoccupied-visible-space",
	},
	"Blood Hawk": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage, requirement/ally-within-5-feet, requirement/ally-not-incapacitated",
		"action|beak|Beak; melee +4 reach 5 range 0/0 4 (1d4 + 2) piercing; effect/piercing-damage 4 (1d4 + 2), effect/bloodied-piercing-damage 6 (1d8 + 2), requirement/target-bloodied",
	},
	"Blue Dragon Wyrmling": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +5 reach 5 range 0/0 8 (1d10 + 3) slashing; effect/slashing-damage 8 (1d10 + 3), effect/lightning-damage 3 (1d6)",
		"action|lightning-breath|Lightning Breath (Recharge 5–6); dex DC 12; line 30/5/0; save/dex, area/line, limit/recharge, effect/lightning-damage 21 (6d6), save/half-damage",
	},
	"Boar": {
		"trait|trait/bloodied-fury|Bloodied Fury; requirement/bloodied, attack/advantage",
		"action|gore|Gore; melee +3 reach 5 range 0/0 4 (1d6 + 1) piercing; effect/piercing-damage 4 (1d6 + 1), requirement/target-size, requirement/charge-distance, effect/extra-piercing-damage 3 (1d6), condition/prone",
	},
	"Bone Devil": {
		"trait|trait/diabolical-restoration|Diabolical Restoration; effect/revive",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|claw|Claw; melee +8 reach 10 range 0/0 13 (2d8 + 4) slashing; effect/slashing-damage 13 (2d8 + 4)",
		"action|infernal-sting|Infernal Sting; melee +8 reach 10 range 0/0 15 (2d10 + 4) piercing; effect/piercing-damage 15 (2d10 + 4), effect/poison-damage 18 (4d8), condition/poisoned, effect/no-hit-point-recovery",
	},
	"Brass Dragon Wyrmling": {
		"action|rend|Rend; melee +4 reach 5 range 0/0 7 (1d10 + 2) slashing; effect/slashing-damage 7 (1d10 + 2)",
		"action|fire-breath|Fire Breath (Recharge 5–6); dex DC 11; line 20/5/0; save/dex, area/line, limit/recharge, effect/fire-damage 14 (4d6), save/half-damage",
		"action|sleep-breath|Sleep Breath; con DC 11; cone 15/0/0; save/con, area/cone, condition/incapacitated, save/repeat, condition/unconscious, save/second-failure, effect/end-on-damage-or-wake",
	},
	"Bronze Dragon Wyrmling": {
		"trait|trait/amphibious|Amphibious; ",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +5 reach 5 range 0/0 8 (1d10 + 3) slashing; effect/slashing-damage 8 (1d10 + 3)",
		"action|lightning-breath|Lightning Breath (Recharge 5–6); dex DC 12; line 40/5/0; save/dex, area/line, limit/recharge, effect/lightning-damage 16 (3d10), save/half-damage",
		"action|repulsion-breath|Repulsion Breath; str DC 12; cone 30/0/0; save/str, area/cone, effect/push, condition/prone",
	},
	"Brown Bear": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +5 reach 5 range 0/0 7 (1d8 + 3) piercing; effect/piercing-damage 7 (1d8 + 3)",
		"action|claw|Claw; melee +5 reach 5 range 0/0 5 (1d4 + 3) slashing; effect/slashing-damage 5 (1d4 + 3), condition/prone, requirement/target-size",
	},
	"Bugbear Stalker": {
		"trait|trait/abduct|Abduct; ",
		"action|multiattack|Multiattack; action/melee-or-ranged-attack, action/melee-attack",
		"action|javelin|Javelin; effect/piercing-damage 13 (3d6 + 3)",
		"action|morningstar|Morningstar; melee +5 reach 10 range 0/0 12 (2d8 + 3) piercing; effect/piercing-damage 12 (2d8 + 3), attack/advantage-against-grappled-target",
		"bonus|quick-grapple|Quick Grapple; dex DC 13; creature 0/0/10; save/dex, target/creature, requirement/target-size, condition/grappled",
	},
}

// Printed values the contract has no typed member for. Each stays in its
// entry's retained source text and is pinned here against the PDF.
var monsters155Printed = []struct{ name, entry, text string }{
	{"Black Pudding", "trait/corrosive-form", "Any nonmagical weapon takes a cumulative −1 penalty to attack rolls immediately after dealing damage to the pudding and coming into contact with it. The weapon is destroyed if the penalty reaches −5."},
	{"Black Pudding", "trait/corrosive-form", "In 1 minute, the pudding can eat through 2 feet of nonmagical wood or metal."},
	{"Black Pudding", "dissolving-pseudopod", "The armor is destroyed if the penalty reduces its AC to 10."},
	{"Black Pudding", "split", "The original pudding’s Hit Points are divided evenly between the new puddings (round down)."},
	{"Blood Hawk", "beak", "Hit: 4 (1d4 + 2) Piercing damage, or 6 (1d8 + 2) Piercing damage if the target is Bloodied."},
	{"Boar", "gore", "the boar moved 20+ feet straight toward it immediately before the hit, the target takes an extra 3 (1d6) Piercing damage and has the Prone condition."},
	{"Brass Dragon Wyrmling", "sleep-breath", "Second Failure: The target has the Unconscious condition for 1 minute."},
	{"Bronze Dragon Wyrmling", "repulsion-breath", "The target is pushed up to 30 feet straight away from the dragon and has the Prone condition."},
	{"Bugbear Stalker", "javelin", "Javelin. Melee or Ranged Attack Roll: +5, reach 10 ft. or range 30/120 ft. Hit: 13 (3d6 + 3) Piercing damage."},
	{"Bugbear Stalker", "morningstar", "Melee Attack Roll: +5 (with Advantage if the target is Grappled by the bugbear), reach 10 ft."},
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters155SourceValues(t *testing.T) {
	c := loadMonsters155(t)
	if len(c.Monsters) != 12 || len(monsters155Headers) != 12 || len(monsters155Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #155") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters155Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader154(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries154(m); !slices.Equal(got, monsters155Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters155Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.Hover || m.LairXP != 0 || m.LegendaryActions != nil || m.Legendary != nil || len(m.Gear) > 0 != (m.Name == "Bugbear Stalker") {
				t.Fatal("source lists no vulnerabilities, hover, lair XP or legendary actions, and gear only for the Bugbear Stalker")
			}
		})
		entries += len(monsters155Entries[want.name])
	}
	// Bold-italic entry headings counted on PDF pp. 263–271 and 345–346.
	if entries != 41 {
		t.Fatalf("reviewed %d entries, want 41", entries)
	}
	for _, p := range monsters155Printed {
		found := false
		for _, m := range c.Monsters {
			if m.Name != p.name {
				continue
			}
			for _, x := range m.Traits {
				found = found || x.Key == p.entry && strings.Contains(x.Text, p.text)
			}
			for _, a := range m.allActions() {
				found = found || a.ID == p.entry && strings.Contains(a.Text, p.text)
			}
		}
		if !found {
			t.Fatalf("%s %s no longer prints %q", p.name, p.entry, p.text)
		}
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters155MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 155 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters155(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the four p. 345–346 animals under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("animals monsters monsters monsters animals monsters animals monsters monsters monsters animals monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #155 cites the page each block starts on, in this order.
	names := strings.Split("Black Bear|Black Dragon Wyrmling|Black Pudding|Blink Dog|Blood Hawk|Blue Dragon Wyrmling|Boar|Bone Devil|Brass Dragon Wyrmling|Bronze Dragon Wyrmling|Brown Bear|Bugbear Stalker", "|")
	pages := []int{345, 263, 265, 266, 345, 266, 346, 267, 268, 269, 346, 271}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic that holds for these twelve blocks (PDF pp. 254–256):
// XP and PB follow CR, Hit Dice follow size and Constitution, saves and skills
// add PB to the ability modifier, and attack rolls add PB to the better of
// Strength and Dexterity (the Javelin is untyped and read from its text).
func TestMonsters155DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"1/8": 2, "1/4": 2, "1/2": 2, "1": 2, "2": 2, "3": 2, "4": 2, "9": 4}
	skillAbility := map[string]string{"deception": "cha", "insight": "wis", "perception": "wis", "stealth": "dex", "survival": "wis"}
	hitDie := map[string]int{"small": 6, "medium": 8, "large": 10}
	for _, m := range loadMonsters155(t).Monsters {
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
			mod := max(m.Abilities["str"].Modifier, m.Abilities["dex"].Modifier)
			for _, a := range m.allActions() {
				if a.Attack == nil {
					continue
				}
				if a.Attack.Bonus != mod+p || a.Attack.Damage.Modifier != mod {
					t.Fatalf("%s attack does not use the ability modifier and PB", a.Name)
				}
			}
		})
	}
}

// Every listed key is used and required, every effect quotes its own entry, and
// every dice expression is typed. Removing any one key invalidates the record.
func TestMonsters155DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters155(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			if m.RequiredMechanics[0] != "monster/stat-block" {
				t.Fatal("stat-block dependency must come first")
			}
			quoted := func(parent string, x MonsterMechanic) {
				if !strings.Contains(parent, strings.TrimSuffix(x.Text, ".")) {
					t.Fatalf("%q is not quoted from its entry", x.Text)
				}
				if x.Dice != nil && !strings.Contains(x.Text, monsterDice154(*x.Dice)) {
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
				if a.Attack == nil {
					continue
				}
				if first := a.Effects[0]; first.Dice == nil || *first.Dice != a.Attack.Damage {
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

// Catalogue presence grants nothing: no key is implemented and the register
// still reports both families unsupported.
func TestMonsters155RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters155(t).Monsters {
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
		all["monster/stat-block"] = true
		if !m.availableWith(all) {
			t.Fatalf("%s: the complete key set should be the only thing that unlocks it", m.Name)
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

func TestMonsters155MalformedFixtures(t *testing.T) {
	// Indexes: 0 Black Bear, 1 Black Dragon Wyrmling, 2 Black Pudding, 3 Blink
	// Dog, 4 Blood Hawk, 6 Boar, 7 Bone Devil, 8 Brass, 9 Bronze, 11 Bugbear.
	tests := []struct {
		name   string
		mutate func(*MonsterCatalogue)
	}{
		{"duplicate ID", func(c *MonsterCatalogue) { c.Monsters[1].ID = c.Monsters[2].ID }},
		{"animal block in another family", func(c *MonsterCatalogue) { c.Monsters[0].ID = "srd-5.2.1:content/spells/black-bear" }},
		{"missing languages list", func(c *MonsterCatalogue) { c.Monsters[0].Languages = nil }},
		{"XP off CR", func(c *MonsterCatalogue) { c.Monsters[1].XP = 451 }},
		{"save off modifier", func(c *MonsterCatalogue) { c.Monsters[1].Abilities["con"] = MonsterAbility{13, 1, 4} }},
		{"modifier", func(c *MonsterCatalogue) { c.Monsters[0].Abilities["str"] = MonsterAbility{15, 3, 3} }},
		{"HP average", func(c *MonsterCatalogue) { c.Monsters[0].HP.Average = 20 }},
		{"melee reach", func(c *MonsterCatalogue) { c.Monsters[0].Actions[1].Attack.Reach = 0 }},
		{"fixed damage without mechanic", func(c *MonsterCatalogue) { c.Monsters[0].Actions[1].Attack.Damage = MonsterHP{Average: 5} }},
		{"untyped weapon action", func(c *MonsterCatalogue) { c.Monsters[11].Actions[1].Key = "action/melee-attack" }},
		{"sense note without sense", func(c *MonsterCatalogue) {
			c.Monsters[7].Senses.Notes = map[string]string{"blindsight": "unimpeded by magical Darkness"}
		}},
		{"unlisted sense note", func(c *MonsterCatalogue) { c.Monsters[7].Senses.Notes["darkvision"] = "impeded by magical Darkness" }},
		{"padded sense note", func(c *MonsterCatalogue) { c.Monsters[7].Senses.Notes["darkvision"] += " " }},
		{"trait dice", func(c *MonsterCatalogue) { c.Monsters[2].Traits[1].Dice.Sides = 9 }},
		{"alternative damage dice", func(c *MonsterCatalogue) { c.Monsters[4].Actions[0].Effects[1].Dice.Average = 7 }},
		{"extra damage sides", func(c *MonsterCatalogue) { c.Monsters[6].Actions[0].Effects[3].Dice.Sides = 9 }},
		{"recharge", func(c *MonsterCatalogue) { c.Monsters[3].BonusActions[0].Uses.Recharge = 1 }},
		{"recharge and daily uses", func(c *MonsterCatalogue) { c.Monsters[3].BonusActions[0].Uses.PerDay = 1 }},
		{"save DC off entry text", func(c *MonsterCatalogue) { c.Monsters[11].BonusActions[0].Save.DC = 14 }},
		{"save ability off entry text", func(c *MonsterCatalogue) { c.Monsters[9].Actions[3].Save.Ability = "dex" }},
		{"creature size", func(c *MonsterCatalogue) { c.Monsters[11].BonusActions[0].Target.Size = 5 }},
		{"cone width", func(c *MonsterCatalogue) { c.Monsters[8].Actions[2].Target.Width = 5 }},
		{"empty reactions", func(c *MonsterCatalogue) { c.Monsters[2].Reactions = []MonsterAction{} }},
		{"empty bonus actions", func(c *MonsterCatalogue) { c.Monsters[3].BonusActions = []MonsterAction{} }},
		{"duplicate action across sections", func(c *MonsterCatalogue) { c.Monsters[3].BonusActions[0].ID = "bite" }},
		{"unlisted reaction effect", func(c *MonsterCatalogue) { c.Monsters[2].Reactions[0].Effects[0].Key = "reaction/future" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters155(t)
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
	raw, err := os.ReadFile("testdata/monsters_155.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, edit := range [][2]string{
		{`"sides": 8,`, `"sides": null,`},
		{`"long_range": 0,`, ``},
		{`"recharge": 5`, `"recharge": null`},
		{`"per_day": 0,`, ``},
		{`"dc": 11`, `"dc": 11, "dc": 12`},
		{`"width": 5,`, ``},
		{`"xp": 450,`, ``},
		{`"notes": {`, `"notes": {"truesight": "see invisible", `},
		{`"shape": "line",`, `"shape": "line", "height": 5,`},
	} {
		if !bytes.Contains(raw, []byte(edit[0])) {
			t.Fatalf("fixture lost %s", edit[0])
		}
		if _, err := DecodeMonsters(bytes.NewReader(bytes.Replace(raw, []byte(edit[0]), []byte(edit[1]), 1))); err == nil {
			t.Fatalf("accepted %s -> %s", edit[0], edit[1])
		}
	}
}

// The batch composes with every earlier stat-block catalogue on main (the #143
// Bandit and batches #151 to #154) without entering the embedded runtime
// catalogue. The #154 file repeats the pinned Bandit revision, so composition
// keeps the one copy from #143.
func TestMonsters155ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	if len(base.Monsters) != 1 || base.Monsters[0].ID != BanditID {
		t.Fatal("batch data must not enter the embedded runtime catalogue")
	}
	batch := loadMonsters155(t)
	earlier := []MonsterCatalogue{loadMonsters151(t), loadMonsters153(t), loadMonsters154(t)}
	b02, err := LoadMonsterBatch02()
	if err != nil {
		t.Fatal(err)
	}
	all := slices.Clone(base.Monsters)
	count := 1
	for i, c := range earlier {
		monsters := c.Monsters
		if c.Monsters[4].ID == BanditID { // the #154 repeat of the pinned #143 record
			monsters = slices.Concat(c.Monsters[:4], c.Monsters[5:])
		}
		all = slices.Concat(all, monsters)
		count += len(monsters)
		if i == 0 {
			for _, sb := range b02.StatBlocks {
				all = append(all, sb.Monster)
				count++
			}
		}
		a, b := c.Source, batch.Source
		a.Changes, b.Changes = "", ""
		if a != b {
			t.Fatal("batches must share one source pin")
		}
	}
	all = slices.Concat(all, batch.Monsters)
	combined := MonsterCatalogue{SchemaVersion: 1, Source: batch.Source, Monsters: all}
	if err := combined.Validate(); err != nil || len(combined.Monsters) != count+12 {
		t.Fatalf("batch must compose with the earlier catalogues: %v (%d records, want %d)", err, len(combined.Monsters), count+12)
	}
	combined.Monsters = append(combined.Monsters, batch.Monsters[3])
	if combined.Validate() == nil {
		t.Fatal("composition accepted a duplicate record")
	}
	_, version, err := BanditRecord()
	if err != nil || version != "bandit-v1:1169c77601cf83a9b96985fc761a6b8be9c901dc1a86e6a9758dcc399707c3d4" {
		t.Fatal("batch changed the pinned Bandit revision:", version, err)
	}
	for _, id := range []string{"rend", "javelin", "morningstar", "teleport", "split"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
