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

func loadMonsters175(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_175.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #175 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Vampire’s, the Vampire Spawn’s and the Vulture’s blocks run onto the next
// page.
var monsters175Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Triceratops", [2]int{363, 363}, "Huge Beast (Dinosaur), Unaligned|AC 14 Initiative −1 (9)|HP 114 (12d12 + 36)|Speed 50 ft.|Str 22 +6 +6 Dex 9 −1 −1 Con 17 +3 +3|Int 2 −4 −4 Wis 11 +0 +0 Cha 5 −3 −3|Senses Passive Perception 10|Languages None|CR 5 (XP 1,800; PB +3)"},
	{"Troll", [2]int{333, 333}, "Large Giant, Chaotic Evil|AC 15 Initiative +1 (11)|HP 94 (9d10 + 45)|Speed 30 ft.|Str 18 +4 +4 Dex 13 +1 +1 Con 20 +5 +5|Int 7 −2 −2 Wis 9 −1 −1 Cha 7 −2 −2|Skills Perception +5|Senses Darkvision 60 ft.; Passive Perception 15|Languages Giant|CR 5 (XP 1,800; PB +3)"},
	{"Troll Limb", [2]int{333, 333}, "Small Giant, Chaotic Evil|AC 13 Initiative +1 (11)|HP 14 (4d6)|Speed 20 ft.|Str 18 +4 +4 Dex 12 +1 +1 Con 10 +0 +0|Int 1 −5 −5 Wis 9 −1 −1 Cha 1 −5 −5|Senses Darkvision 60 ft.; Passive Perception 9|Languages None|CR 1/2 (XP 100; PB +2)"},
	{"Tyrannosaurus Rex", [2]int{363, 363}, "Huge Beast (Dinosaur), Unaligned|AC 13 Initiative +3 (13)|HP 136 (13d12 + 52)|Speed 50 ft.|Str 25 +7 +10 Dex 10 +0 +0 Con 19 +4 +4|Int 2 −4 −4 Wis 12 +1 +4 Cha 9 −1 −1|Skills Perception +4|Senses Passive Perception 14|Languages None|CR 8 (XP 3,900; PB +3)"},
	{"Unicorn", [2]int{334, 334}, "Large Celestial, Lawful Good|AC 12 Initiative +8 (18)|HP 97 (13d10 + 26)|Speed 50 ft.|Str 18 +4 +4 Dex 14 +2 +2 Con 15 +2 +2|Int 11 +0 +0 Wis 17 +3 +3 Cha 16 +3 +3|Immunities Poison; Charmed, Paralyzed, Poisoned|Senses Darkvision 60 ft.; Passive Perception 13|Languages Celestial, Elvish, Sylvan; telepathy 120 ft.|CR 5 (XP 1,800; PB +3)"},
	{"Vampire", [2]int{335, 336}, "Medium or Small Undead, Lawful Evil|AC 16 Initiative +14 (24)|HP 195 (23d8 + 92)|Speed 40 ft., Climb 40 ft.|Str 18 +4 +4 Dex 18 +4 +9 Con 18 +4 +9|Int 17 +3 +3 Wis 15 +2 +7 Cha 18 +4 +9|Skills Perception +7, Stealth +9|Resistances Necrotic|Senses Darkvision 120 ft.; Passive Perception 17|Languages Common plus two other languages|CR 13 (XP 10,000, or 11,500 in lair; PB +5)"},
	{"Vampire Familiar", [2]int{334, 334}, "Medium or Small Humanoid, Neutral Evil|AC 15 Initiative +5 (15)|HP 65 (10d8 + 20)|Speed 30 ft., Climb 30 ft.|Str 17 +3 +3 Dex 16 +3 +5 Con 15 +2 +2|Int 10 +0 +0 Wis 10 +0 +2 Cha 14 +2 +2|Skills Perception +4, Persuasion +4, Stealth +7|Resistances Necrotic|Immunities Charmed (except from its vampire master)|Gear Daggers (10)|Senses Darkvision 60 ft.; Passive Perception 14|Languages Common plus one other language|CR 3 (XP 700; PB +2)"},
	{"Vampire Spawn", [2]int{334, 335}, "Medium or Small Undead, Neutral Evil|AC 16 Initiative +3 (13)|HP 90 (12d8 + 36)|Speed 30 ft.|Str 16 +3 +3 Dex 16 +3 +6 Con 16 +3 +3|Int 11 +0 +0 Wis 10 +0 +3 Cha 12 +1 +1|Skills Perception +3, Stealth +6|Resistances Necrotic|Senses Darkvision 60 ft.; Passive Perception 13|Languages Common plus one other language|CR 5 (XP 1,800; PB +3)"},
	{"Venomous Snake", [2]int{363, 363}, "Tiny Beast, Unaligned|AC 12 Initiative +2 (12)|HP 5 (2d4)|Speed 30 ft., Swim 30 ft.|Str 2 −4 −4 Dex 15 +2 +2 Con 11 +0 +0|Int 1 −5 −5 Wis 10 +0 +0 Cha 3 −4 −4|Senses Blindsight 10 ft.; Passive Perception 10|Languages None|CR 1/8 (XP 25; PB +2)"},
	{"Violet Fungus", [2]int{286, 286}, "Medium Plant, Unaligned|AC 5 Initiative −5 (5)|HP 18 (4d8)|Speed 5 ft.|Str 3 −4 −4 Dex 1 −5 −5 Con 10 +0 +0|Int 1 −5 −5 Wis 3 −4 −4 Cha 1 −5 −5|Immunities Blinded, Charmed, Deafened, Frightened|Senses Blindsight 30 ft.; Passive Perception 6|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Vrock", [2]int{336, 336}, "Large Fiend (Demon), Chaotic Evil|AC 15 Initiative +2 (12)|HP 152 (16d10 + 64)|Speed 40 ft., Fly 60 ft.|Str 17 +3 +3 Dex 15 +2 +5 Con 18 +4 +4|Int 8 −1 −1 Wis 13 +1 +4 Cha 8 −1 +2|Resistances Cold, Fire, Lightning|Immunities Poison; Poisoned|Senses Darkvision 120 ft.; Passive Perception 11|Languages Abyssal; telepathy 120 ft.|CR 6 (XP 2,300; PB +3)"},
	{"Vulture", [2]int{363, 364}, "Medium Beast, Unaligned|AC 10 Initiative +0 (10)|HP 5 (1d8 + 1)|Speed 10 ft., Fly 50 ft.|Str 7 −2 −2 Dex 10 +0 +0 Con 13 +1 +1|Int 2 −4 −4 Wis 12 +1 +1 Cha 4 −3 −3|Skills Perception +3|Senses Passive Perception 13|Languages None|CR 0 (XP 10; PB +2)"},
}

// Typed values of every trait and action entry, reviewed against the same pages:
// limited uses, attack bonus/reach/range/damage, save ability and DC, target
// shape size/width/range, dice and dependency keys. A ranged attack printed with
// one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text. The Vampire Familiar's Umbral Dagger is a "Melee or Ranged Attack Roll".
// A Vampire heading's form qualifier is a requirement of its entry.
var monsters175Entries = map[string][]string{
	"Triceratops": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|gore|Gore; melee +9 reach 5 range 0/0 19 (2d12 + 6) piercing; effect/piercing-damage 19 (2d12 + 6), requirement/target-size, requirement/straight-line-movement, effect/piercing-damage 9 (2d8), condition/prone",
	},
	"Troll": {
		"trait|trait/loathsome-limbs|Loathsome Limbs (4/Day); limit/per-day, effect/sever-limb, requirement/bloodied-after-slashing-damage, effect/limb-acts, effect/exhaustion-per-missing-limb, effect/regrow-limbs",
		"trait|trait/regeneration|Regeneration; effect/regeneration, effect/regeneration-suppressed, effect/death-requires-no-regeneration",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +7 reach 10 range 0/0 11 (2d6 + 4) slashing; effect/slashing-damage 11 (2d6 + 4)",
		"bonus|charge|Charge; movement/half-speed-toward-enemy",
	},
	"Troll Limb": {
		"trait|trait/regeneration|Regeneration; effect/regeneration, effect/regeneration-suppressed, effect/death-requires-no-regeneration",
		"trait|trait/troll-spawn|Troll Spawn; effect/troll-senses, effect/grow-into-troll",
		"action|rend|Rend; melee +6 reach 5 range 0/0 9 (2d4 + 4) slashing; effect/slashing-damage 9 (2d4 + 4)",
	},
	"Tyrannosaurus Rex": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +10 reach 10 range 0/0 33 (4d12 + 7) piercing; effect/piercing-damage 33 (4d12 + 7), requirement/target-size, condition/grappled, condition/restrained, limit/no-tail-on-grappled-target",
		"action|tail|Tail; melee +10 reach 15 range 0/0 25 (4d8 + 7) bludgeoning; effect/bludgeoning-damage 25 (4d8 + 7), requirement/target-size, condition/prone",
	},
	"Unicorn": {
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day); limit/per-day, save/choose-success",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|hooves|Hooves; melee +7 reach 5 range 0/0 11 (2d6 + 4) bludgeoning; effect/bludgeoning-damage 11 (2d6 + 4)",
		"action|radiant-horn|Radiant Horn; melee +7 reach 5 range 0/0 9 (1d10 + 4) radiant; effect/radiant-damage 9 (1d10 + 4)",
		"action|spellcasting|Spellcasting; cha DC 14 +0; at will: Detect Evil and Good, Druidcraft; 1/day: Calm Emotions, Dispel Evil and Good, Entangle, Pass without Trace, Word of Recall; spellcasting/no-components, spell/detect-evil-and-good [spellcasting/at-will], spell/druidcraft [spellcasting/at-will], spell/calm-emotions [limit/per-day], spell/dispel-evil-and-good [limit/per-day], spell/entangle [limit/per-day], spell/pass-without-trace [limit/per-day], spell/word-of-recall [limit/per-day]",
		"bonus|unicorns-blessing|Unicorn’s Blessing (3/Day); limit/per-day, effect/touch-creature, action/spellcasting, spell/cure-wounds, spell/lesser-restoration",
		"legendary|charging-horn|Charging Horn; movement/half-speed, effect/no-opportunity-attacks, action/melee-attack",
		"legendary|shimmering-shield|Shimmering Shield; target/self-or-creature, effect/temporary-hit-points 10 (3d6), effect/ac-bonus, limit/not-again-until-next-turn",
		"legendary-uses|3|0",
	},
	"Vampire": {
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day, or 4/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"trait|trait/misty-escape|Misty Escape; effect/misty-escape, effect/destroyed-without-shape-shift, effect/mist-form-limits, effect/resting-place-deadline, condition/paralyzed, effect/resting-place-recovery",
		"trait|trait/spider-climb|Spider Climb; effect/climb-without-check",
		"trait|trait/vampire-weakness|Vampire Weakness; effect/no-uninvited-entry, effect/running-water-damage, effect/acid-damage, effect/stake-to-the-heart, effect/sunlight-damage, effect/radiant-damage, effect/sunlight-disadvantage",
		"action|multiattack|Multiattack; requirement/vampire-form, action/melee-attack, action/saving-throw",
		"action|grave-strike|Grave Strike; melee +9 reach 5 range 0/0 8 (1d8 + 4) bludgeoning; effect/bludgeoning-damage 8 (1d8 + 4), requirement/vampire-form, effect/necrotic-damage 7 (2d6), requirement/target-size, condition/grappled, limit/grapple-from-one-of-two-hands",
		"action|bite|Bite; con DC 17; creature 0/0/5; save/con, target/creature, requirement/bat-or-vampire-form, requirement/willing-or-held-target, effect/piercing-damage 6 (1d4 + 4), effect/necrotic-damage 13 (3d8), effect/hit-point-maximum-reduction, effect/vampire-healing, effect/rise-as-vampire-spawn",
		"bonus|charm|Charm (Recharge 5–6); limit/recharge, action/spellcasting, spell/charm-person, spellcasting/no-components, spellcasting/charisma-ability, spellcasting/fixed-save-dc, effect/spell-duration-24-hours, effect/charmed-target-willing-for-bite, effect/target-unaware-of-charm",
		"bonus|shape-shift|Shape-Shift; requirement/not-in-sunlight-or-running-water, effect/shape-shift, effect/equipment-transforms, effect/bat-form-limits, effect/shapechange-form-limits, effect/mist-form-limits, effect/mist-movement, effect/mist-damage-resistance",
		"legendary|beguile|Beguile; action/spellcasting, spell/command, spellcasting/no-components, spellcasting/charisma-ability, spellcasting/fixed-save-dc, limit/not-again-until-next-turn",
		"legendary|deathless-strike|Deathless Strike; movement/half-speed, action/melee-attack",
		"legendary-uses|3|4",
	},
	"Vampire Familiar": {
		"trait|trait/vampiric-connection|Vampiric Connection; effect/master-telepathy, effect/master-shared-senses, requirement/same-plane-as-master",
		"action|multiattack|Multiattack; action/melee-or-ranged-attack",
		"action|umbral-dagger|Umbral Dagger; effect/piercing-damage 5 (1d4 + 3), effect/necrotic-damage 7 (3d4), effect/stable-at-zero-hit-points, condition/poisoned, condition/paralyzed",
		"bonus|deathless-agility|Deathless Agility; action/dash, action/disengage",
	},
	"Vampire Spawn": {
		"trait|trait/spider-climb|Spider Climb; effect/climb-without-check",
		"trait|trait/vampire-weakness|Vampire Weakness; effect/no-uninvited-entry, effect/running-water-damage, effect/acid-damage, effect/stake-to-the-heart, effect/sunlight-damage, effect/radiant-damage, effect/sunlight-disadvantage",
		"action|multiattack|Multiattack; action/melee-attack, action/saving-throw",
		"action|claw|Claw; melee +6 reach 5 range 0/0 8 (2d4 + 3) slashing; effect/slashing-damage 8 (2d4 + 3), requirement/target-size, condition/grappled, limit/grapple-from-one-of-two-claws",
		"action|bite|Bite; con DC 14; creature 0/0/5; save/con, target/creature, requirement/willing-or-held-target, effect/piercing-damage 5 (1d4 + 3), effect/necrotic-damage 10 (3d6), effect/hit-point-maximum-reduction, effect/vampire-healing",
		"bonus|deathless-agility|Deathless Agility; action/dash, action/disengage",
	},
	"Venomous Snake": {
		"action|bite|Bite; melee +4 reach 5 range 0/0 4 (1d4 + 2) piercing; effect/piercing-damage 4 (1d4 + 2), effect/poison-damage 3 (1d6)",
	},
	"Violet Fungus": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rotting-touch|Rotting Touch; melee +2 reach 10 range 0/0 4 (1d8) necrotic; effect/necrotic-damage 4 (1d8)",
	},
	"Vrock": {
		"trait|trait/demonic-restoration|Demonic Restoration; effect/revive",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|shred|Shred; melee +6 reach 5 range 0/0 10 (2d6 + 3) piercing; effect/piercing-damage 10 (2d6 + 3), effect/poison-damage 10 (3d6)",
		"action|spores|Spores (Recharge 6); con DC 15; emanation 20/0/0; limit/recharge, save/con, area/emanation, condition/poisoned, save/repeat, effect/poison-damage 5 (1d10), effect/holy-water-ends-effect",
		"action|stunning-screech|Stunning Screech (1/Day); con DC 15; emanation 20/0/0; limit/per-day, save/con, area/emanation, save/automatic-success-for-demons, effect/thunder-damage 10 (3d6), condition/stunned",
	},
	"Vulture": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage-with-nearby-ally, requirement/ally-not-incapacitated",
		"action|beak|Beak; melee +2 reach 5 range 0/0 2 (1d4) piercing; effect/piercing-damage 2 (1d4)",
	},
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters175SourceValues(t *testing.T) {
	c := loadMonsters175(t)
	if len(c.Monsters) != 12 || len(monsters175Headers) != 12 || len(monsters175Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #175") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters175Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader170(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries164(m); !slices.Equal(got, monsters175Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters175Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.VulnerabilityNote != "" || m.ResistanceNote != "" || m.TelepathyNote != "" || m.Senses.Notes != nil || m.SwarmOf != "" || m.SpeedNote != "" {
				t.Fatal("source lists no vulnerabilities, vulnerability note, resistance note, telepathy note, sense notes, swarm line, speed note")
			}
		})
		for _, line := range monsters175Entries[want.name] {
			if !strings.HasPrefix(line, "legendary-uses|") {
				entries++
			}
		}
	}
	// Bold-italic entry headings counted on PDF pp. 286, 333–336 and 363–364.
	if entries != 54 {
		t.Fatalf("reviewed %d entries, want 54", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters175MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 175 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters175(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Animals appendix blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("animals monsters monsters animals monsters monsters monsters monsters animals monsters monsters animals")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #175 cites the page each block starts on, in this order.
	names := []string{"Triceratops", "Troll", "Troll Limb", "Tyrannosaurus Rex", "Unicorn", "Vampire", "Vampire Familiar", "Vampire Spawn", "Venomous Snake", "Violet Fungus", "Vrock", "Vulture"}
	pages := []int{363, 333, 333, 363, 334, 335, 334, 334, 363, 286, 336, 363}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic for these twelve blocks; see monsterConsistency170.
func TestMonsters175DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"5": 3, "1/2": 2, "8": 3, "13": 5, "3": 2, "1/8": 2, "1/4": 2, "6": 3, "0": 2}
	attacks := map[string]attackAbility170{
		"Violet Fungus/Rotting Touch": {"con", false},
	}
	for _, m := range loadMonsters175(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			monsterConsistency170(t, m, pb, attacks)
		})
	}
}

// Every listed key is used and required, every effect quotes its own entry, and
// every dice expression is typed. Removing any one key invalidates the record.
func TestMonsters175DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters175(t).Monsters {
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
func TestMonsters175RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters175(t).Monsters {
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

func TestMonsters175MalformedFixtures(t *testing.T) {
	tests := []struct {
		name, record string
		mutate       func(*Monster)
	}{
		{"missing languages list", "Vulture", func(m *Monster) { m.Languages = nil }},
		{"animal or monster block in another family", "Venomous Snake", func(m *Monster) {
			m.ID = strings.Replace(strings.Replace(m.ID, "/animals/", "/spells/", 1), "/monsters/", "/spells/", 1)
		}},
		{"XP off CR", "Triceratops", func(m *Monster) { m.XP++ }},
		{"save off modifier", "Triceratops", func(m *Monster) {
			m.Abilities["wis"] = MonsterAbility{m.Abilities["wis"].Score, m.Abilities["wis"].Modifier, m.Abilities["wis"].Modifier + 11}
		}},
		{"modifier", "Unicorn", func(m *Monster) {
			m.Abilities["str"] = MonsterAbility{m.Abilities["str"].Score, m.Abilities["str"].Modifier + 1, m.Abilities["str"].Modifier + 1}
		}},
		{"HP average", "Unicorn", func(m *Monster) { m.HP.Average++ }},
		{"condition immunity note beside a typed list", "Vampire Familiar", func(m *Monster) { m.ConditionImmunities = []string{"charmed"} }},
		{"padded condition immunity note", "Vampire Familiar", func(m *Monster) { m.ConditionImmunityNote += " " }},
		{"recharge", "Vrock", func(m *Monster) { m.Actions[2].Uses.Recharge = 7 }},
		{"emanation width", "Vrock", func(m *Monster) { m.Actions[2].Target.Width = 5 }},
		{"save DC off entry text", "Vampire", func(m *Monster) { m.Actions[2].Save.DC = 18 }},
		{"trait daily uses without count", "Troll", func(m *Monster) { m.Traits[0].Uses.PerDay = 0 }},
		{"unlisted spell", "Unicorn", func(m *Monster) { m.Actions[3].Spellcasting.AtWill[0] = "Wish" }},
		{"spell DC off entry text", "Unicorn", func(m *Monster) { m.Actions[3].Spellcasting.SaveDC = 15 }},
		{"lair XP not above XP", "Vampire", func(m *Monster) { m.LairXP = m.XP }},
		{"legendary lair uses not above uses", "Vampire", func(m *Monster) { m.Legendary.LairUses = 3 }},
		{"empty bonus actions", "Vulture", func(m *Monster) { m.BonusActions = []MonsterAction{} }},
		{"empty reactions", "Vulture", func(m *Monster) { m.Reactions = []MonsterAction{} }},
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
		c := loadMonsters175(t)
		c.Monsters[1].ID = c.Monsters[2].ID
		rejected(t, c, true)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters175(t)
			tt.mutate(record(t, &c, tt.record))
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			rejected(t, c, !strings.HasPrefix(tt.name, "empty "))
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_175.json")
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
		{`"dc":17`, `"dc":17,"dc":22`},
		{`"width":0,`, ``},
		{`"save_dc":14,`, ``},
		{`"attack_bonus":0`, `"attack_bonus":0,"attack_bonus":5`},
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
// to 24 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters175ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters175(t)
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
	for _, load := range []func(*testing.T) MonsterCatalogue{loadMonsters158, loadMonsters159, loadMonsters160, loadMonsters161, loadMonsters162, loadMonsters163, loadMonsters164, loadMonsters165, loadMonsters166, loadMonsters167, loadMonsters168, loadMonsters169, loadMonsters170, loadMonsters171, loadMonsters172, loadMonsters173, loadMonsters174} {
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
	for _, id := range []string{"beak", "bite", "claw", "gore", "grave-strike", "hooves", "radiant-horn", "rend", "rotting-touch", "shred", "tail"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
