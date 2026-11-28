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

func loadMonsters173(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_173.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #173 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Stirge’s and the Storm Giant’s blocks run onto the next
// page.
var monsters173Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Sphinx of Valor", [2]int{328, 328}, "Large Celestial, Lawful Neutral|AC 17 Initiative +12 (22)|HP 199 (19d10 + 95)|Speed 40 ft., Fly 60 ft.|Str 22 +6 +6 Dex 10 +0 +6 Con 20 +5 +11|Int 16 +3 +9 Wis 23 +6 +12 Cha 18 +4 +4|Skills Arcana +9, Perception +12, Religion +15|Resistances Necrotic, Radiant|Immunities Psychic; Charmed, Frightened|Senses Truesight 120 ft.; Passive Perception 22|Languages Celestial, Common|CR 17 (XP 18,000, or 20,000 in lair; PB +6)"},
	{"Sphinx of Wonder", [2]int{327, 327}, "Tiny Celestial, Lawful Good|AC 13 Initiative +3 (13)|HP 24 (7d4 + 7)|Speed 20 ft., Fly 40 ft.|Str 6 −2 −2 Dex 17 +3 +3 Con 13 +1 +1|Int 15 +2 +2 Wis 12 +1 +1 Cha 11 +0 +0|Skills Arcana +4, Religion +4, Stealth +5|Resistances Necrotic, Psychic, Radiant|Senses Darkvision 60 ft.; Passive Perception 11|Languages Celestial, Common|CR 1 (XP 200; PB +2)"},
	{"Spider", [2]int{361, 361}, "Tiny Beast, Unaligned|AC 12 Initiative +2 (12)|HP 1 (1d4 − 1)|Speed 20 ft., Climb 20 ft.|Str 2 −4 −4 Dex 14 +2 +2 Con 8 −1 −1|Int 1 −5 −5 Wis 10 +0 +0 Cha 2 −4 −4|Skills Stealth +4|Senses Darkvision 30 ft.; Passive Perception 10|Languages None|CR 0 (XP 10; PB +2)"},
	{"Spirit Naga", [2]int{329, 329}, "Large Fiend, Chaotic Evil|AC 17 Initiative +3 (13)|HP 135 (18d10 + 36)|Speed 40 ft.|Str 18 +4 +4 Dex 17 +3 +6 Con 14 +2 +5|Int 16 +3 +3 Wis 15 +2 +5 Cha 16 +3 +6|Immunities Poison; Charmed, Poisoned|Senses Darkvision 60 ft.; Passive Perception 12|Languages Abyssal, Common|CR 8 (XP 3,900; PB +3)"},
	{"Sprite", [2]int{329, 329}, "Tiny Fey, Neutral Good|AC 15 Initiative +4 (14)|HP 10 (4d4)|Speed 10 ft., Fly 40 ft.|Str 3 −4 −4 Dex 18 +4 +4 Con 10 +0 +0|Int 14 +2 +2 Wis 13 +1 +1 Cha 11 +0 +0|Skills Perception +3, Stealth +8|Senses Passive Perception 13|Languages Common, Elvish, Sylvan|CR 1/4 (XP 50; PB +2)"},
	{"Spy", [2]int{329, 329}, "Medium or Small Humanoid, Neutral|AC 12 Initiative +4 (14)|HP 27 (6d8)|Speed 30 ft., Climb 30 ft.|Str 10 +0 +0 Dex 15 +2 +2 Con 10 +0 +0|Int 12 +1 +1 Wis 14 +2 +2 Cha 16 +3 +3|Skills Deception +5, Insight +4, Investigation +5, Perception +6, Sleight of Hand +4, Stealth +6|Gear Hand Crossbow, Shortsword, Thieves’ Tools|Senses Passive Perception 16|Languages Common plus one other language|CR 1 (XP 200; PB +2)"},
	{"Steam Mephit", [2]int{308, 308}, "Small Elemental, Neutral Evil|AC 10 Initiative +0 (10)|HP 17 (5d6)|Speed 30 ft., Fly 30 ft.|Str 5 −3 −3 Dex 11 +0 +0 Con 10 +0 +0|Int 11 +0 +0 Wis 10 +0 +0 Cha 12 +1 +1|Skills Stealth +2|Immunities Fire, Poison; Exhaustion, Poisoned|Senses Darkvision 60 ft.; Passive Perception 10|Languages Primordial (Aquan, Ignan)|CR 1/4 (XP 50; PB +2)"},
	{"Stirge", [2]int{329, 330}, "Tiny Monstrosity, Unaligned|AC 13 Initiative +3 (13)|HP 5 (2d4)|Speed 10 ft., Fly 40 ft.|Str 4 −3 −3 Dex 16 +3 +3 Con 11 +0 +0|Int 2 −4 −4 Wis 8 −1 −1 Cha 6 −2 −2|Senses Darkvision 60 ft.; Passive Perception 9|Languages None|CR 1/8 (XP 25; PB +2)"},
	{"Stone Giant", [2]int{330, 330}, "Huge Giant, Neutral|AC 17 Initiative +5 (15)|HP 126 (11d12 + 55)|Speed 40 ft.|Str 23 +6 +6 Dex 15 +2 +5 Con 20 +5 +8|Int 10 +0 +0 Wis 12 +1 +4 Cha 9 −1 −1|Skills Athletics +12, Perception +4, Stealth +5|Senses Darkvision 60 ft.; Passive Perception 14|Languages Giant|CR 7 (XP 2,900; PB +3)"},
	{"Stone Golem", [2]int{330, 330}, "Large Construct, Unaligned|AC 18 Initiative +3 (13)|HP 220 (21d10 + 105)|Speed 30 ft.|Str 22 +6 +6 Dex 9 −1 −1 Con 20 +5 +5|Int 3 −4 −4 Wis 11 +0 +0 Cha 1 −5 −5|Immunities Poison, Psychic; Charmed, Exhaustion, Frightened, Paralyzed, Petrified, Poisoned|Senses Darkvision 120 ft.; Passive Perception 10|Languages Understands Common plus two other languages but can’t speak|CR 10 (XP 5,900; PB +4)"},
	{"Storm Giant", [2]int{330, 331}, "Huge Giant, Chaotic Good|AC 16 Initiative +7 (17)|HP 230 (20d12 + 100)|Speed 50 ft., Fly 25 ft. (hover), Swim 50 ft.|Str 29 +9 +14 Dex 14 +2 +2 Con 20 +5 +10|Int 16 +3 +3 Wis 20 +5 +10 Cha 18 +4 +9|Skills Arcana +8, Athletics +14, History +8, Perception +10|Resistances Cold|Immunities Lightning, Thunder|Senses Darkvision 120 ft., Truesight 30 ft.; Passive Perception 20|Languages Common, Giant|CR 13 (XP 10,000; PB +5)"},
	{"Succubus", [2]int{331, 331}, "Medium Fiend, Neutral Evil|AC 15 Initiative +3 (13)|HP 71 (13d8 + 13)|Speed 30 ft., Fly 60 ft.|Str 8 −1 −1 Dex 17 +3 +3 Con 13 +1 +1|Int 15 +2 +2 Wis 12 +1 +1 Cha 20 +5 +5|Skills Deception +9, Insight +5, Perception +5, Persuasion +9, Stealth +7|Resistances Cold, Fire, Poison, Psychic|Senses Darkvision 60 ft.; Passive Perception 15|Languages Abyssal, Common, Infernal; telepathy 60 ft.|CR 4 (XP 1,100; PB +2)"},
}

// Typed values of every trait and action entry, reviewed against the same pages:
// limited uses, attack bonus/reach/range/damage, save ability and DC, target
// shape size/width/range, dice and dependency keys. A ranged attack printed with
// one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text. The Storm Giant's Lightning Storm covers a Cylinder, and the Sphinx of
// Valor's Roar prints three saves. The Sphinx of Valor's Weight of Years prints
// DC 16, the Sphinx of Lore's value, beside its other DC 20 entries (p. 328).
var monsters173Entries = map[string][]string{
	"Sphinx of Valor": {
		"trait|trait/inscrutable|Inscrutable; effect/scrying-immunity, check/disadvantage-to-read-intentions",
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day, or 4/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"action|multiattack|Multiattack; action/melee-attack, action/special",
		"action|claw|Claw; melee +12 reach 5 range 0/0 20 (4d6 + 6) slashing; effect/slashing-damage 20 (4d6 + 6)",
		"action|roar|Roar (3/Day); limit/per-day, effect/sequential-roars, limit/long-rest-sequence-reset, target/enemies-only, save/wis, area/emanation, condition/frightened, condition/paralyzed, save/repeat, save/automatic-success, save/con, effect/thunder-damage 44 (8d10), condition/prone, save/half-damage-only",
		"action|spellcasting|Spellcasting; wis DC 20 +0; at will: Detect Evil and Good, Thaumaturgy; 1/day: Detect Magic, Dispel Magic, Greater Restoration, Heroes’ Feast, Zone of Truth; spellcasting/no-material-components, spell/detect-evil-and-good [spellcasting/at-will], spell/thaumaturgy [spellcasting/at-will], spell/detect-magic [limit/per-day], spell/dispel-magic [limit/per-day], spell/greater-restoration [limit/per-day], spell/heroes-feast [limit/per-day], spell/zone-of-truth [limit/per-day]",
		"legendary|arcane-prowl|Arcane Prowl; movement/teleport, action/melee-attack",
		"legendary|weight-of-years|Weight of Years; con DC 16; creature 0/0/120; save/con, target/creature, effect/exhaustion-level, effect/appear-older, limit/not-again-until-next-turn",
		"legendary-uses|3|4",
	},
	"Sphinx of Wonder": {
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|rend|Rend; melee +5 reach 5 range 0/0 5 (1d4 + 3) slashing; effect/slashing-damage 5 (1d4 + 3), effect/radiant-damage 7 (2d6)",
		"reaction|burst-of-ingenuity|Burst of Ingenuity (2/Day); limit/per-day, reaction/trigger, effect/add-to-roll",
	},
	"Spider": {
		"trait|trait/spider-climb|Spider Climb; effect/climb-without-check",
		"trait|trait/web-walker|Web Walker; effect/ignore-web-restrictions, effect/know-web-contacts",
		"action|bite|Bite; melee +4 reach 5 range 0/0 1 piercing; effect/piercing-damage [damage/fixed], effect/poison-damage 2 (1d4)",
	},
	"Spirit Naga": {
		"trait|trait/fiendish-restoration|Fiendish Restoration; effect/revive, effect/wish-prevents-revival",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|bite|Bite; melee +7 reach 10 range 0/0 7 (1d6 + 4) piercing; effect/piercing-damage 7 (1d6 + 4), effect/poison-damage 14 (4d6)",
		"action|necrotic-ray|Necrotic Ray; ranged +6 reach 0 range 60/60 21 (6d6) necrotic; effect/necrotic-damage 21 (6d6)",
		"action|spellcasting|Spellcasting; int DC 14 +0; at will: Detect Magic, Mage Hand, Minor Illusion, Water Breathing; 2/day: Detect Thoughts, Dimension Door, Hold Person (level 3 version), Lightning Bolt (level 4 version); spellcasting/no-somatic-or-material-components, spell/detect-magic [spellcasting/at-will], spell/mage-hand [spellcasting/at-will], spell/minor-illusion [spellcasting/at-will], spell/water-breathing [spellcasting/at-will], spell/detect-thoughts [limit/per-day], spell/dimension-door [limit/per-day], spell/hold-person [limit/per-day], spell/lightning-bolt [limit/per-day]",
	},
	"Sprite": {
		"action|needle-sword|Needle Sword; melee +6 reach 5 range 0/0 6 (1d4 + 4) piercing; effect/piercing-damage 6 (1d4 + 4)",
		"action|enchanting-bow|Enchanting Bow; ranged +6 reach 0 range 40/160 1 piercing; effect/piercing-damage [damage/fixed], condition/charmed",
		"action|heart-sight|Heart Sight; cha DC 10; creature 0/0/5; save/cha, target/creature, save/automatic-failure, effect/know-emotions-and-alignment",
		"action|invisibility|Invisibility; action/spellcasting, spell/invisibility, spellcasting/no-components, spellcasting/charisma-ability",
	},
	"Spy": {
		"action|shortsword|Shortsword; melee +4 reach 5 range 0/0 5 (1d6 + 2) piercing; effect/piercing-damage 5 (1d6 + 2), effect/poison-damage 7 (2d6)",
		"action|hand-crossbow|Hand Crossbow; ranged +4 reach 0 range 30/120 5 (1d6 + 2) piercing; effect/piercing-damage 5 (1d6 + 2), effect/poison-damage 7 (2d6)",
		"bonus|cunning-action|Cunning Action; action/dash, action/disengage, action/hide",
	},
	"Steam Mephit": {
		"trait|trait/blurred-form|Blurred Form; attack/disadvantage-against-self, requirement/not-incapacitated",
		"trait|trait/death-burst|Death Burst; dex DC 10; emanation 5/0/0; 5 (2d4); save/dex, area/emanation, effect/death-explosion, effect/fire-damage, save/half-damage",
		"action|claw|Claw; melee +2 reach 5 range 0/0 2 (1d4) slashing; effect/slashing-damage 2 (1d4), effect/fire-damage 2 (1d4)",
		"action|steam-breath|Steam Breath (Recharge 6); con DC 10; cone 15/0/0; limit/recharge, save/con, area/cone, effect/fire-damage 5 (2d4), effect/speed-reduction, save/half-damage-only, effect/no-underwater-fire-resistance",
	},
	"Stirge": {
		"action|proboscis|Proboscis; melee +5 reach 5 range 0/0 6 (1d6 + 3) piercing; effect/piercing-damage 6 (1d6 + 3), effect/attach, limit/no-proboscis-while-attached, effect/necrotic-damage 5 (2d4), effect/detach-with-movement, effect/detach-as-action",
	},
	"Stone Giant": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|stone-club|Stone Club; melee +9 reach 15 range 0/0 22 (3d10 + 6) bludgeoning; effect/bludgeoning-damage 22 (3d10 + 6)",
		"action|boulder|Boulder; ranged +9 reach 0 range 60/240 15 (2d8 + 6) bludgeoning; effect/bludgeoning-damage 15 (2d8 + 6), requirement/target-size, condition/prone",
		"reaction|deflect-missile|Deflect Missile (Recharge 5–6); dex DC 17; creature 0/0/60; limit/recharge, save/dex, target/creature, reaction/trigger, effect/damage-reduction 11 (1d10 + 6), effect/redirect-force, effect/force-damage 11 (1d10 + 6)",
	},
	"Stone Golem": {
		"trait|trait/immutable-form|Immutable Form; effect/no-shape-shift",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|slam|Slam; melee +10 reach 5 range 0/0 15 (2d8 + 6) bludgeoning; effect/bludgeoning-damage 15 (2d8 + 6), effect/force-damage 9 (2d8)",
		"action|force-bolt|Force Bolt; ranged +9 reach 0 range 120/120 22 (4d10) force; effect/force-damage 22 (4d10)",
		"bonus|slow|Slow (Recharge 5–6); limit/recharge, action/spellcasting, spell/slow, spellcasting/no-components, spellcasting/constitution-ability, spellcasting/fixed-save-dc",
	},
	"Storm Giant": {
		"trait|trait/amphibious|Amphibious; effect/breathe-air-and-water",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|storm-sword|Storm Sword; melee +14 reach 10 range 0/0 23 (4d6 + 9) slashing; effect/slashing-damage 23 (4d6 + 9), effect/lightning-damage 13 (3d8)",
		"action|thunderbolt|Thunderbolt; ranged +14 reach 0 range 500/500 22 (2d12 + 9) lightning; effect/lightning-damage 22 (2d12 + 9), condition/blinded, condition/deafened",
		"action|lightning-storm|Lightning Storm (Recharge 5–6); dex DC 18; limit/recharge, save/dex, area/cylinder, effect/lightning-damage 55 (10d10), save/half-damage",
		"action|spellcasting|Spellcasting; wis DC 18 +0; at will: Detect Magic, Light; 1/day: Control Weather; spellcasting/no-material-components, spell/detect-magic [spellcasting/at-will], spell/light [spellcasting/at-will], spell/control-weather [limit/per-day]",
	},
	"Succubus": {
		"trait|trait/incubus-form|Incubus Form; effect/change-to-incubus, requirement/long-rest-finished",
		"action|multiattack|Multiattack; action/melee-attack, action/special, action/saving-throw",
		"action|fiendish-touch|Fiendish Touch; melee +7 reach 5 range 0/0 16 (2d10 + 5) psychic; effect/psychic-damage 16 (2d10 + 5)",
		"action|charm|Charm; action/spellcasting, spell/dominate-person, spellcasting/no-components, spellcasting/charisma-ability, spellcasting/fixed-save-dc",
		"action|draining-kiss|Draining Kiss; con DC 15; creature 0/0/5; save/con, target/creature, requirement/charmed-target, effect/psychic-damage 13 (3d8), save/half-damage, effect/hit-point-maximum-reduction",
		"bonus|shape-shift|Shape-Shift; effect/shape-shift, effect/shapechange-form-limits, effect/equipment-not-transformed",
	},
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters173SourceValues(t *testing.T) {
	c := loadMonsters173(t)
	if len(c.Monsters) != 12 || len(monsters173Headers) != 12 || len(monsters173Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #173") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters173Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader170(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries164(m); !slices.Equal(got, monsters173Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters173Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.VulnerabilityNote != "" || m.ResistanceNote != "" || m.ConditionImmunityNote != "" || m.TelepathyNote != "" || m.Senses.Notes != nil || m.SwarmOf != "" || m.SpeedNote != "" {
				t.Fatal("source lists no vulnerabilities, vulnerability note, resistance note, condition immunity note, telepathy note, sense notes, swarm line, speed note")
			}
		})
		for _, line := range monsters173Entries[want.name] {
			if !strings.HasPrefix(line, "legendary-uses|") {
				entries++
			}
		}
	}
	// Bold-italic entry headings counted on PDF pp. 308, 327–331 and 361.
	if entries != 53 {
		t.Fatalf("reviewed %d entries, want 53", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters173MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 173 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters173(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Animals appendix blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("monsters monsters animals monsters monsters monsters monsters monsters monsters monsters monsters monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #173 cites the page each block starts on, in this order.
	names := []string{"Sphinx of Valor", "Sphinx of Wonder", "Spider", "Spirit Naga", "Sprite", "Spy", "Steam Mephit", "Stirge", "Stone Giant", "Stone Golem", "Storm Giant", "Succubus"}
	pages := []int{328, 327, 361, 329, 329, 329, 308, 329, 330, 330, 330, 331}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic for these twelve blocks; see monsterConsistency170.
func TestMonsters173DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"17": 6, "1": 2, "0": 2, "8": 3, "1/4": 2, "1/8": 2, "7": 3, "10": 4, "13": 5, "4": 2}
	attacks := map[string]attackAbility170{
		"Spirit Naga/Necrotic Ray": {"cha", false},
		"Stone Golem/Force Bolt":   {"con", false},
		"Succubus/Fiendish Touch":  {"cha", true},
	}
	for _, m := range loadMonsters173(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			monsterConsistency170(t, m, pb, attacks)
		})
	}
}

// Every listed key is used and required, every effect quotes its own entry, and
// every dice expression is typed. Removing any one key invalidates the record.
func TestMonsters173DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters173(t).Monsters {
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
func TestMonsters173RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters173(t).Monsters {
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

func TestMonsters173MalformedFixtures(t *testing.T) {
	tests := []struct {
		name, record string
		mutate       func(*Monster)
	}{
		{"missing languages list", "Sprite", func(m *Monster) { m.Languages = nil }},
		{"animal or monster block in another family", "Spider", func(m *Monster) {
			m.ID = strings.Replace(strings.Replace(m.ID, "/animals/", "/spells/", 1), "/monsters/", "/spells/", 1)
		}},
		{"XP off CR", "Storm Giant", func(m *Monster) { m.XP++ }},
		{"save off modifier", "Storm Giant", func(m *Monster) {
			m.Abilities["wis"] = MonsterAbility{m.Abilities["wis"].Score, m.Abilities["wis"].Modifier, m.Abilities["wis"].Modifier + 11}
		}},
		{"modifier", "Storm Giant", func(m *Monster) {
			m.Abilities["str"] = MonsterAbility{m.Abilities["str"].Score, m.Abilities["str"].Modifier + 1, m.Abilities["str"].Modifier + 1}
		}},
		{"HP average", "Storm Giant", func(m *Monster) { m.HP.Average++ }},
		{"recharge", "Steam Mephit", func(m *Monster) { m.Actions[1].Uses.Recharge = 7 }},
		{"trait save DC off entry text", "Steam Mephit", func(m *Monster) { m.Traits[1].Save.DC = 11 }},
		{"trait dice average", "Steam Mephit", func(m *Monster) { m.Traits[1].Dice.Average++ }},
		{"untypable cylinder shape", "Storm Giant", func(m *Monster) { m.Actions[3].Target = &MonsterTarget{Shape: "cylinder", Size: 10, Range: 500} }},
		{"unlisted spell", "Spirit Naga", func(m *Monster) { m.Actions[3].Spellcasting.AtWill[0] = "Wish" }},
		{"spell DC off entry text", "Spirit Naga", func(m *Monster) { m.Actions[3].Spellcasting.SaveDC = 15 }},
		{"daily uses without count", "Sphinx of Wonder", func(m *Monster) { m.Reactions[0].Uses.PerDay = 0 }},
		{"legendary lair uses not above uses", "Sphinx of Valor", func(m *Monster) { m.Legendary.LairUses = 3 }},
		{"fixed damage of zero", "Spider", func(m *Monster) { m.Actions[0].Attack.Damage.Average = 0 }},
		{"empty bonus actions", "Spy", func(m *Monster) { m.BonusActions = []MonsterAction{} }},
		{"empty reactions", "Spy", func(m *Monster) { m.Reactions = []MonsterAction{} }},
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
		c := loadMonsters173(t)
		c.Monsters[1].ID = c.Monsters[2].ID
		rejected(t, c, true)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters173(t)
			tt.mutate(record(t, &c, tt.record))
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			rejected(t, c, !strings.HasPrefix(tt.name, "empty "))
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_173.json")
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
		{`"recharge":6`, `"recharge":null`},
		{`"per_day":0,`, ``},
		{`"dc":16`, `"dc":16,"dc":22`},
		{`"width":0,`, ``},
		{`"save_dc":20,`, ``},
		{`"attack_bonus":0`, `"attack_bonus":0,"attack_bonus":5`},
		{`"darkvision":60`, `"darkvision":null`},
		{`"shape":"creature",`, `"shape":"creature","height":5,`},
		{`"lair_uses":4,`, ``},
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
// to 22 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters173ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters173(t)
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
	for _, load := range []func(*testing.T) MonsterCatalogue{loadMonsters158, loadMonsters159, loadMonsters160, loadMonsters161, loadMonsters162, loadMonsters163, loadMonsters164, loadMonsters165, loadMonsters166, loadMonsters167, loadMonsters168, loadMonsters169, loadMonsters170, loadMonsters171, loadMonsters172} {
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
	for _, id := range []string{"bite", "boulder", "claw", "enchanting-bow", "fiendish-touch", "force-bolt", "hand-crossbow", "necrotic-ray", "needle-sword", "proboscis", "rend", "shortsword", "slam", "stone-club", "storm-sword", "thunderbolt"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
