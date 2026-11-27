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

func loadMonsters169(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_169.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #169 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Phase Spider's, Pit Fiend's and Plesiosaurus's blocks run onto the next
// page.
var monsters169Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Otyugh", [2]int{313, 313}, "Large Aberration, Neutral|AC 14 Initiative +0 (10)|HP 104 (11d10 + 44)|Speed 30 ft.|Str 16 +3 +3 Dex 11 +0 +0 Con 19 +4 +7|Int 6 −2 −2 Wis 13 +1 +1 Cha 6 −2 −2|Senses Darkvision 120 ft.; Passive Perception 11|Languages Otyugh; telepathy 120 ft. (doesn’t allow the receiving creature to respond telepathically)|CR 5 (XP 1,800; PB +3)"},
	{"Owl", [2]int{358, 358}, "Tiny Beast, Unaligned|AC 11 Initiative +1 (11)|HP 1 (1d4 − 1)|Speed 5 ft., Fly 60 ft.|Str 3 −4 −4 Dex 13 +1 +1 Con 8 −1 −1|Int 2 −4 −4 Wis 12 +1 +1 Cha 7 −2 −2|Skills Perception +5, Stealth +5|Senses Darkvision 120 ft.; Passive Perception 15|Languages None|CR 0 (XP 10; PB +2)"},
	{"Owlbear", [2]int{313, 313}, "Large Monstrosity, Unaligned|AC 13 Initiative +1 (11)|HP 59 (7d10 + 21)|Speed 40 ft., Climb 40 ft.|Str 20 +5 +5 Dex 12 +1 +1 Con 17 +3 +3|Int 3 −4 −4 Wis 12 +1 +1 Cha 7 −2 −2|Skills Perception +5|Senses Darkvision 60 ft.; Passive Perception 15|Languages None|CR 3 (XP 700; PB +2)"},
	{"Panther", [2]int{358, 358}, "Medium Beast, Unaligned|AC 13 Initiative +3 (13)|HP 13 (3d8)|Speed 50 ft., Climb 40 ft.|Str 14 +2 +2 Dex 16 +3 +3 Con 10 +0 +0|Int 3 −4 −4 Wis 14 +2 +2 Cha 7 −2 −2|Skills Perception +4, Stealth +7|Senses Darkvision 60 ft.; Passive Perception 14|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Pegasus", [2]int{313, 313}, "Large Celestial, Chaotic Good|AC 12 Initiative +2 (12)|HP 59 (7d10 + 21)|Speed 60 ft., Fly 90 ft.|Str 18 +4 +4 Dex 15 +2 +4 Con 16 +3 +5|Int 10 +0 +0 Wis 15 +2 +4 Cha 13 +1 +3|Skills Perception +6|Senses Passive Perception 16|Languages Understands Celestial, Common, Elvish, and Sylvan but can’t speak|CR 2 (XP 450; PB +2)"},
	{"Phase Spider", [2]int{313, 314}, "Large Monstrosity, Unaligned|AC 14 Initiative +3 (13)|HP 45 (7d10 + 7)|Speed 30 ft., Climb 30 ft.|Str 15 +2 +2 Dex 16 +3 +3 Con 12 +1 +1|Int 6 −2 −2 Wis 10 +0 +0 Cha 6 −2 −2|Skills Stealth +7|Senses Darkvision 60 ft.; Passive Perception 10|Languages None|CR 3 (XP 700; PB +2)"},
	{"Piranha", [2]int{358, 358}, "Tiny Beast, Unaligned|AC 13 Initiative +3 (13)|HP 1 (1d4 − 1)|Speed 5 ft., Swim 40 ft.|Str 2 −4 −4 Dex 16 +3 +3 Con 9 −1 −1|Int 1 −5 −5 Wis 7 −2 −2 Cha 2 −4 −4|Senses Darkvision 60 ft.; Passive Perception 8|Languages None|CR 0 (XP 10; PB +2)"},
	{"Pirate", [2]int{314, 314}, "Medium or Small Humanoid, Neutral|AC 14 Initiative +5 (15)|HP 33 (6d8 + 6)|Speed 30 ft.|Str 10 +0 +0 Dex 16 +3 +5 Con 12 +1 +1|Int 8 −1 −1 Wis 12 +1 +1 Cha 14 +2 +4|Gear Daggers (6), Leather Armor|Senses Passive Perception 11|Languages Common plus one other language|CR 1 (XP 200; PB +2)"},
	{"Pirate Captain", [2]int{314, 314}, "Medium or Small Humanoid, Neutral|AC 17 Initiative +7 (17)|HP 84 (13d8 + 26)|Speed 30 ft.|Str 10 +0 +3 Dex 18 +4 +7 Con 14 +2 +2|Int 10 +0 +0 Wis 14 +2 +5 Cha 17 +3 +6|Skills Acrobatics +7, Perception +5|Gear Pistol, Rapier|Senses Passive Perception 15|Languages Common plus one other language|CR 6 (XP 2,300; PB +3)"},
	{"Pit Fiend", [2]int{314, 315}, "Large Fiend (Devil), Lawful Evil|AC 21 Initiative +14 (24)|HP 337 (27d10 + 189)|Speed 30 ft., Fly 60 ft.|Str 26 +8 +8 Dex 14 +2 +8 Con 24 +7 +7|Int 22 +6 +6 Wis 18 +4 +10 Cha 24 +7 +7|Skills Perception +10, Persuasion +19|Resistances Cold|Immunities Fire, Poison; Poisoned|Senses Truesight 120 ft.; Passive Perception 20|Languages Infernal; telepathy 120 ft.|CR 20 (XP 25,000; PB +6)"},
	{"Planetar", [2]int{315, 315}, "Large Celestial (Angel), Lawful Good|AC 19 Initiative +10 (20)|HP 262 (21d10 + 147)|Speed 40 ft., Fly 120 ft. (hover)|Str 24 +7 +12 Dex 20 +5 +5 Con 24 +7 +12|Int 19 +4 +4 Wis 22 +6 +11 Cha 25 +7 +12|Skills Perception +11|Resistances Radiant|Immunities Charmed, Exhaustion, Frightened|Senses Truesight 120 ft.; Passive Perception 21|Languages All; telepathy 120 ft.|CR 16 (XP 15,000; PB +5)"},
	{"Plesiosaurus", [2]int{358, 359}, "Large Beast (Dinosaur), Unaligned|AC 13 Initiative +2 (12)|HP 68 (8d10 + 24)|Speed 20 ft., Swim 40 ft.|Str 18 +4 +4 Dex 15 +2 +2 Con 16 +3 +3|Int 2 −4 −4 Wis 12 +1 +1 Cha 5 −3 −3|Skills Perception +3, Stealth +4|Senses Passive Perception 13|Languages None|CR 2 (XP 450; PB +2)"},
}

// Typed values of every trait and action entry, reviewed against the same pages:
// limited uses, attack bonus/reach/range/damage, save ability and DC, target
// shape size/width/range, dice and dependency keys. A ranged attack printed with
// one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text. The Pirate's Dagger is a "Melee or Ranged Attack Roll" (key
// action/melee-or-ranged-attack). The Planetar prints "using Charisma as spellcasting ability".
var monsters169Entries = map[string][]string{
	"Otyugh": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +6 reach 5 range 0/0 12 (2d8 + 3) piercing; con DC 15; effect/piercing-damage 12 (2d8 + 3), condition/poisoned, requirement/poisoned-target-finishes-long-rest, save/con, effect/hit-point-maximum-reduction 5 (1d10), effect/hit-point-maximum-no-reset, effect/success-ends-condition",
		"action|tentacle|Tentacle; melee +6 reach 10 range 0/0 12 (2d8 + 3) piercing; effect/piercing-damage 12 (2d8 + 3), requirement/target-size, condition/grappled, limit/grapple-from-one-of-two-tentacles",
		"action|tentacle-slam|Tentacle Slam; con DC 14; save/con, target/grappled-creatures, effect/bludgeoning-damage 16 (3d8 + 3), condition/stunned, save/half-damage",
	},
	"Owl": {
		"trait|trait/flyby|Flyby; effect/no-opportunity-attack-flying-out",
		"action|talons|Talons; melee +3 reach 5 range 0/0 1 slashing; effect/slashing-damage [damage/fixed]",
	},
	"Owlbear": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +7 reach 5 range 0/0 14 (2d8 + 5) slashing; effect/slashing-damage 14 (2d8 + 5)",
	},
	"Panther": {
		"action|rend|Rend; melee +5 reach 5 range 0/0 6 (1d6 + 3) slashing; effect/slashing-damage 6 (1d6 + 3)",
		"bonus|nimble-escape|Nimble Escape; action/disengage, action/hide",
	},
	"Pegasus": {
		"action|hooves|Hooves; melee +6 reach 5 range 0/0 7 (1d6 + 4) bludgeoning; effect/bludgeoning-damage 7 (1d6 + 4), effect/radiant-damage 5 (2d4)",
	},
	"Phase Spider": {
		"trait|trait/ethereal-sight|Ethereal Sight; effect/see-into-ethereal-plane",
		"trait|trait/spider-climb|Spider Climb; effect/climb-without-check",
		"trait|trait/web-walker|Web Walker; effect/ignore-web-restrictions, effect/know-web-contacts",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +5 reach 5 range 0/0 8 (1d10 + 3) piercing; effect/piercing-damage 8 (1d10 + 3), effect/poison-damage 9 (2d8), requirement/damage-reduces-to-zero-hit-points, effect/becomes-stable, condition/poisoned, condition/paralyzed",
		"bonus|ethereal-jaunt|Ethereal Jaunt; effect/plane-shift-to-or-from-ethereal",
	},
	"Piranha": {
		"trait|trait/water-breathing|Water Breathing; effect/breathe-only-underwater",
		"action|bite|Bite; melee +5 reach 5 range 0/0 1 piercing; effect/piercing-damage [damage/fixed], attack/advantage-against-target-missing-hit-points",
	},
	"Pirate": {
		"action|multiattack|Multiattack; action/melee-or-ranged-attack, multiattack/replace-attack, action/saving-throw",
		"action|dagger|Dagger; effect/piercing-damage 5 (1d4 + 3)",
		"action|enthralling-panache|Enthralling Panache; wis DC 12; creature 0/0/30; save/wis, target/creature, condition/charmed",
	},
	"Pirate Captain": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|rapier|Rapier; melee +7 reach 5 range 0/0 13 (2d8 + 4) piercing; effect/piercing-damage 13 (2d8 + 4), attack/advantage-on-next-attack",
		"action|pistol|Pistol; ranged +7 reach 0 range 30/90 15 (2d10 + 4) piercing; effect/piercing-damage 15 (2d10 + 4)",
		"bonus|captains-charm|Captain’s Charm; wis DC 14; creature 0/0/30; save/wis, target/creature, condition/charmed",
		"reaction|riposte|Riposte; reaction/trigger, requirement/holding-weapon, effect/ac-bonus, action/melee-attack, target/triggering-creature",
	},
	"Pit Fiend": {
		"trait|trait/diabolical-restoration|Diabolical Restoration; effect/revive",
		"trait|trait/fear-aura|Fear Aura; wis DC 21; emanation 20/0/0; area/emanation, requirement/not-incapacitated, save/wis, target/enemies-starting-turn-in-area, condition/frightened, save/success-immunity",
		"trait|trait/legendary-resistance|Legendary Resistance (4/Day); limit/per-day, save/choose-success",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +14 reach 10 range 0/0 18 (3d6 + 8) piercing; con DC 21; effect/piercing-damage 18 (3d6 + 8), requirement/target-creature, save/con, condition/poisoned, effect/no-hit-point-recovery, effect/start-of-turn-damage 21 (6d6), save/repeat, save/automatic-success",
		"action|devilish-claw|Devilish Claw; melee +14 reach 10 range 0/0 26 (4d8 + 8) necrotic; effect/necrotic-damage 26 (4d8 + 8)",
		"action|fiery-mace|Fiery Mace; melee +14 reach 10 range 0/0 22 (4d6 + 8) force; effect/force-damage 22 (4d6 + 8), effect/fire-damage 21 (6d6)",
		"action|hellfire-spellcasting|Hellfire Spellcasting (Recharge 4–6); limit/recharge, action/spellcasting, spell/fireball, spellcasting/no-material-components, spellcasting/charisma-ability, spellcasting/fixed-save-dc, spellcasting/replace-one-casting, spell/hold-monster, spell/wall-of-fire",
	},
	"Planetar": {
		"trait|trait/divine-awareness|Divine Awareness; effect/know-lies",
		"trait|trait/exalted-restoration|Exalted Restoration; effect/revive",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/alternative-routine, action/saving-throw",
		"action|radiant-sword|Radiant Sword; melee +12 reach 10 range 0/0 14 (2d6 + 7) slashing; effect/slashing-damage 14 (2d6 + 7), effect/radiant-damage 18 (4d8)",
		"action|holy-burst|Holy Burst; dex DC 20; sphere 20/0/120; save/dex, area/sphere, target/enemies-only, effect/radiant-damage 24 (7d6), save/half-damage",
		"action|spellcasting|Spellcasting; cha DC 20 +0; at will: Detect Evil and Good; 1/day: Commune, Control Weather, Dispel Evil and Good, Raise Dead; spellcasting/no-material-components, spell/detect-evil-and-good [spellcasting/at-will], spell/commune [limit/per-day], spell/control-weather [limit/per-day], spell/dispel-evil-and-good [limit/per-day], spell/raise-dead [limit/per-day]",
		"bonus|divine-aid|Divine Aid (2/Day); limit/per-day, action/spellcasting, spell/cure-wounds, spell/invisibility, spell/lesser-restoration, spell/remove-curse",
	},
	"Plesiosaurus": {
		"trait|trait/hold-breath|Hold Breath; effect/hold-breath",
		"action|bite|Bite; melee +6 reach 10 range 0/0 11 (2d6 + 4) piercing; effect/piercing-damage 11 (2d6 + 4)",
	},
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters169SourceValues(t *testing.T) {
	c := loadMonsters169(t)
	if len(c.Monsters) != 12 || len(monsters169Headers) != 12 || len(monsters169Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #169") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters169Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader164(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries164(m); !slices.Equal(got, monsters169Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters169Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.ResistanceNote != "" || m.LairXP != 0 || m.Legendary != nil || m.LegendaryActions != nil {
				t.Fatal("source lists no vulnerabilities, resistance note, lair XP, legendary actions")
			}
		})
		for _, line := range monsters169Entries[want.name] {
			if !strings.HasPrefix(line, "legendary-uses|") {
				entries++
			}
		}
	}
	// Bold-italic entry headings counted on PDF pp. 313–315 and 358–359.
	if entries != 46 {
		t.Fatalf("reviewed %d entries, want 46", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters169MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 169 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters169(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Animals appendix blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("monsters animals monsters animals monsters monsters animals monsters monsters monsters monsters animals")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #169 cites the page each block starts on, in this order.
	names := []string{"Otyugh", "Owl", "Owlbear", "Panther", "Pegasus", "Phase Spider", "Piranha", "Pirate", "Pirate Captain", "Pit Fiend", "Planetar", "Plesiosaurus"}
	pages := []int{313, 358, 313, 358, 313, 313, 358, 314, 314, 314, 315, 358}
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
func TestMonsters169DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"5": 3, "0": 2, "3": 2, "1/4": 2, "2": 2, "1": 2, "6": 3, "20": 6, "16": 5}
	skillAbility := map[string]string{"acrobatics": "dex", "perception": "wis", "persuasion": "cha", "stealth": "dex"}
	hitDie := map[string]int{"tiny": 4, "medium": 8, "large": 10}
	for _, m := range loadMonsters169(t).Monsters {
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
func TestMonsters169DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters169(t).Monsters {
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
func TestMonsters169RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters169(t).Monsters {
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

func TestMonsters169MalformedFixtures(t *testing.T) {
	tests := []struct {
		name, record string
		mutate       func(*Monster)
	}{
		{"missing languages list", "Owl", func(m *Monster) { m.Languages = nil }},
		{"animal block in another family", "Owl", func(m *Monster) {
			m.ID = strings.Replace(strings.Replace(m.ID, "/animals/", "/spells/", 1), "/monsters/", "/spells/", 1)
		}},
		{"XP off CR", "Plesiosaurus", func(m *Monster) { m.XP++ }},
		{"save off modifier", "Plesiosaurus", func(m *Monster) {
			m.Abilities["wis"] = MonsterAbility{m.Abilities["wis"].Score, m.Abilities["wis"].Modifier, m.Abilities["wis"].Modifier + 7}
		}},
		{"modifier", "Otyugh", func(m *Monster) {
			m.Abilities["str"] = MonsterAbility{m.Abilities["str"].Score, m.Abilities["str"].Modifier + 1, m.Abilities["str"].Modifier + 1}
		}},
		{"HP average", "Otyugh", func(m *Monster) { m.HP.Average++ }},
		{"telepathy note without telepathy", "Otyugh", func(m *Monster) { m.Telepathy = 0 }},
		{"padded telepathy note", "Otyugh", func(m *Monster) { m.TelepathyNote += " " }},
		{"recharge", "Pit Fiend", func(m *Monster) { m.Actions[4].Uses.Recharge = 7 }},
		{"trait save DC off entry text", "Pit Fiend", func(m *Monster) { m.Traits[1].Save.DC = 22 }},
		{"sphere width", "Planetar", func(m *Monster) { m.Actions[2].Target.Width = 5 }},
		{"save DC off entry text", "Pirate", func(m *Monster) { m.Actions[2].Save.DC = 13 }},
		{"unlisted spell", "Planetar", func(m *Monster) { m.Actions[3].Spellcasting.AtWill[0] = "Wish" }},
		{"spell DC off entry text", "Planetar", func(m *Monster) { m.Actions[3].Spellcasting.SaveDC = 19 }},
		{"daily uses without count", "Planetar", func(m *Monster) { m.BonusActions[0].Uses.PerDay = 0 }},
		{"fixed damage of zero", "Owl", func(m *Monster) { m.Actions[0].Attack.Damage.Average = 0 }},
		{"empty bonus actions", "Panther", func(m *Monster) { m.BonusActions = []MonsterAction{} }},
		{"empty reactions", "Pirate Captain", func(m *Monster) { m.Reactions = []MonsterAction{} }},
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
		c := loadMonsters169(t)
		c.Monsters[1].ID = c.Monsters[2].ID
		rejected(t, c, true)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters169(t)
			tt.mutate(record(t, &c, tt.record))
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			rejected(t, c, !strings.HasPrefix(tt.name, "empty "))
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_169.json")
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
		{`"recharge":4`, `"recharge":null`},
		{`"per_day":0,`, ``},
		{`"dc":21`, `"dc":21,"dc":22`},
		{`"width":0,`, ``},
		{`"save_dc":20,`, ``},
		{`"attack_bonus":0`, `"attack_bonus":0,"attack_bonus":5`},
		{`"darkvision":120`, `"darkvision":null`},
		{`"shape":"sphere",`, `"shape":"sphere","height":5,`},
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
// to 18 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters169ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters169(t)
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
	for _, load := range []func(*testing.T) MonsterCatalogue{loadMonsters158, loadMonsters159, loadMonsters160, loadMonsters161, loadMonsters162, loadMonsters163, loadMonsters164, loadMonsters165, loadMonsters166, loadMonsters167, loadMonsters168} {
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
	for _, id := range []string{"bite", "devilish-claw", "fiery-mace", "hooves", "pistol", "radiant-sword", "rapier", "rend", "talons", "tentacle"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
