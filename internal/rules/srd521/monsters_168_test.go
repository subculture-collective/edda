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

func loadMonsters168(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_168.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #168 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Mummy Lord's, Nalfeshnee's, Octopus's and Oni's blocks run onto the next
// page. The Oni's Spellcasting prints no At Will list, so that typed list is empty.
var monsters168Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Mule", [2]int{357, 357}, "Medium Beast, Unaligned|AC 10 Initiative +0 (10)|HP 11 (2d8 + 2)|Speed 40 ft.|Str 14 +2 +4 Dex 10 +0 +0 Con 13 +1 +1|Int 2 −4 −4 Wis 10 +0 +0 Cha 5 −3 −3|Senses Passive Perception 10|Languages None|CR 1/8 (XP 25; PB +2)"},
	{"Mummy", [2]int{309, 309}, "Medium or Small Undead, Lawful Evil|AC 11 Initiative −1 (9)|HP 58 (9d8 + 18)|Speed 20 ft.|Str 16 +3 +3 Dex 8 −1 −1 Con 15 +2 +2|Int 6 −2 −2 Wis 12 +1 +3 Cha 12 +1 +1|Vulnerabilities Fire|Immunities Necrotic, Poison; Charmed, Exhaustion, Frightened, Paralyzed, Poisoned|Senses Darkvision 60 ft.; Passive Perception 11|Languages Common plus two other languages|CR 3 (XP 700; PB +2)"},
	{"Mummy Lord", [2]int{309, 310}, "Medium or Small Undead (Cleric), Lawful Evil|AC 17 Initiative +10 (20)|HP 187 (25d8 + 75)|Speed 30 ft.|Str 18 +4 +4 Dex 10 +0 +0 Con 17 +3 +3|Int 11 +0 +5 Wis 19 +4 +9 Cha 16 +3 +3|Skills History +5, Perception +9, Religion +5|Vulnerabilities Fire|Immunities Necrotic, Poison; Charmed, Exhaustion, Frightened, Paralyzed, Poisoned|Senses Truesight 60 ft.; Passive Perception 19|Languages Common plus three other languages|CR 15 (XP 13,000, or 15,000 in lair; PB +5)"},
	{"Nalfeshnee", [2]int{310, 311}, "Large Fiend (Demon), Chaotic Evil|AC 18 Initiative +5 (15)|HP 184 (16d10 + 96)|Speed 20 ft., Fly 30 ft.|Str 21 +5 +5 Dex 10 +0 +0 Con 22 +6 +11|Int 19 +4 +9 Wis 12 +1 +6 Cha 15 +2 +7|Resistances Cold, Fire, Lightning|Immunities Poison; Frightened, Poisoned|Senses Truesight 120 ft.; Passive Perception 11|Languages Abyssal; telepathy 120 ft.|CR 13 (XP 10,000; PB +5)"},
	{"Night Hag", [2]int{311, 311}, "Medium Fiend, Neutral Evil|AC 17 Initiative +5 (15)|HP 112 (15d8 + 45)|Speed 30 ft.|Str 18 +4 +4 Dex 15 +2 +2 Con 16 +3 +3|Int 16 +3 +3 Wis 14 +2 +2 Cha 16 +3 +3|Skills Deception +6, Insight +5, Perception +5, Stealth +5|Resistances Cold, Fire|Immunities Charmed|Senses Darkvision 120 ft.; Passive Perception 15|Languages Abyssal, Common, Infernal, Primordial|CR 5 (XP 1,800; PB +3)"},
	{"Nightmare", [2]int{311, 311}, "Large Fiend, Neutral Evil|AC 13 Initiative +2 (12)|HP 68 (8d10 + 24)|Speed 60 ft., Fly 90 ft. (hover)|Str 18 +4 +4 Dex 15 +2 +2 Con 16 +3 +3|Int 10 +0 +0 Wis 13 +1 +1 Cha 15 +2 +2|Immunities Fire|Senses Passive Perception 11|Languages Understands Abyssal, Common, and Infernal but can’t speak|CR 3 (XP 700; PB +2)"},
	{"Noble", [2]int{312, 312}, "Medium or Small Humanoid, Neutral|AC 15 Initiative +1 (11)|HP 9 (2d8)|Speed 30 ft.|Str 11 +0 +0 Dex 12 +1 +1 Con 11 +0 +0|Int 12 +1 +1 Wis 14 +2 +2 Cha 16 +3 +3|Skills Deception +5, Insight +4, Persuasion +5|Gear Breastplate, Rapier|Senses Passive Perception 12|Languages Common plus two other languages|CR 1/8 (XP 25; PB +2)"},
	{"Ochre Jelly", [2]int{312, 312}, "Large Ooze, Unaligned|AC 8 Initiative −2 (8)|HP 52 (7d10 + 14)|Speed 20 ft., Climb 20 ft.|Str 15 +2 +2 Dex 6 −2 −2 Con 14 +2 +2|Int 2 −4 −4 Wis 6 −2 −2 Cha 1 −5 −5|Resistances Acid|Immunities Lightning, Slashing; Charmed, Deafened, Exhaustion, Frightened, Grappled, Prone, Restrained|Senses Blindsight 60 ft.; Passive Perception 8|Languages None|CR 2 (XP 450; PB +2)"},
	{"Octopus", [2]int{357, 358}, "Small Beast, Unaligned|AC 12 Initiative +2 (12)|HP 3 (1d6)|Speed 5 ft., Swim 30 ft.|Str 4 −3 −3 Dex 15 +2 +2 Con 11 +0 +0|Int 3 −4 −4 Wis 10 +0 +0 Cha 4 −3 −3|Skills Perception +2, Stealth +6|Senses Darkvision 30 ft.; Passive Perception 12|Languages None|CR 0 (XP 10; PB +2)"},
	{"Ogre", [2]int{312, 312}, "Large Giant, Chaotic Evil|AC 11 Initiative −1 (9)|HP 68 (8d10 + 24)|Speed 40 ft.|Str 19 +4 +4 Dex 8 −1 −1 Con 16 +3 +3|Int 5 −3 −3 Wis 7 −2 −2 Cha 7 −2 −2|Gear Greatclub, Javelins (3)|Senses Darkvision 60 ft.; Passive Perception 8|Languages Common, Giant|CR 2 (XP 450; PB +2)"},
	{"Ogre Zombie", [2]int{344, 344}, "Large Undead, Neutral Evil|AC 8 Initiative −2 (8)|HP 85 (9d10 + 36)|Speed 30 ft.|Str 19 +4 +4 Dex 6 −2 −2 Con 18 +4 +4|Int 3 −4 −4 Wis 6 −2 +0 Cha 5 −3 −3|Immunities Poison; Exhaustion, Poisoned|Senses Darkvision 60 ft.; Passive Perception 8|Languages Understands Common and Giant but can’t speak|CR 2 (XP 450; PB +2)"},
	{"Oni", [2]int{312, 313}, "Large Fiend, Lawful Evil|AC 17 Initiative +0 (10)|HP 119 (14d10 + 42)|Speed 30 ft., Fly 30 ft. (hover)|Str 19 +4 +4 Dex 11 +0 +3 Con 16 +3 +6|Int 14 +2 +2 Wis 12 +1 +4 Cha 15 +2 +5|Skills Arcana +5, Deception +8, Perception +4|Resistances Cold|Senses Darkvision 60 ft.; Passive Perception 14|Languages Common, Giant|CR 7 (XP 2,900; PB +3)"},
}

// Typed values of every trait and action entry, reviewed against the same pages:
// limited uses, attack bonus/reach/range/damage, save ability and DC, target
// shape size/width/range, dice and dependency keys. A ranged attack printed with
// one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text. The Ogre's Javelin is a "Melee or Ranged Attack Roll" (key
// action/melee-or-ranged-attack).
var monsters168Entries = map[string][]string{
	"Mule": {
		"trait|trait/beast-of-burden|Beast of Burden; effect/larger-carrying-capacity",
		"action|hooves|Hooves; melee +4 reach 5 range 0/0 4 (1d4 + 2) bludgeoning; effect/bludgeoning-damage 4 (1d4 + 2)",
	},
	"Mummy": {
		"action|multiattack|Multiattack; action/melee-attack, multiattack/additional-action, action/saving-throw",
		"action|rotting-fist|Rotting Fist; melee +5 reach 5 range 0/0 8 (1d10 + 3) bludgeoning; effect/bludgeoning-damage 8 (1d10 + 3), effect/necrotic-damage 10 (3d6), requirement/target-creature, effect/curse, effect/no-hit-point-recovery, effect/hit-point-maximum-no-reset, effect/hit-point-maximum-reduction 10 (3d6), effect/die-at-zero-hit-points",
		"action|dreadful-glare|Dreadful Glare; wis DC 11; creature 0/0/60; save/wis, target/creature, condition/frightened, save/success-immunity",
	},
	"Mummy Lord": {
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day, or 4/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"trait|trait/undead-restoration|Undead Restoration; effect/revive, effect/revival-delay, requirement/intact-heart, effect/reform-in-lair, effect/heart-statistics",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/additional-action, action/saving-throw",
		"action|rotting-fist|Rotting Fist; melee +9 reach 5 range 0/0 15 (2d10 + 4) bludgeoning; effect/bludgeoning-damage 15 (2d10 + 4), effect/necrotic-damage 10 (3d6), requirement/target-creature, effect/curse, effect/no-hit-point-recovery, effect/no-long-rest-benefit, effect/hit-point-maximum-reduction 10 (3d6), effect/die-at-zero-hit-points",
		"action|channel-negative-energy|Channel Negative Energy; ranged +9 reach 0 range 60/60 25 (6d6 + 4) necrotic; effect/necrotic-damage 25 (6d6 + 4)",
		"action|dreadful-glare|Dreadful Glare; wis DC 17; creature 0/0/60; save/wis, target/creature, effect/psychic-damage 25 (6d6 + 4), condition/paralyzed",
		"action|spellcasting|Spellcasting; wis DC 17 +9; at will: Dispel Magic, Thaumaturgy; 1/day: Animate Dead, Harm, Insect Plague (level 7 version); spellcasting/no-material-components, spellcasting/spell-attack-bonus, spell/dispel-magic [spellcasting/at-will], spell/thaumaturgy [spellcasting/at-will], spell/animate-dead [limit/per-day], spell/harm [limit/per-day], spell/insect-plague [limit/per-day]",
		"reaction|whirlwind-of-sand|Whirlwind of Sand; reaction/trigger, effect/ac-bonus, effect/teleport, requirement/unoccupied-visible-space, target/creatures-of-choice, condition/blinded",
		"legendary|dread-command|Dread Command; action/spellcasting, spell/command, limit/once-per-round",
		"legendary|glare|Glare; action/saving-throw, limit/once-per-round",
		"legendary|necrotic-strike|Necrotic Strike; action/melee-attack, action/ranged-attack",
		"legendary-uses|3|4",
	},
	"Nalfeshnee": {
		"trait|trait/demonic-restoration|Demonic Restoration; effect/revive",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +10 reach 10 range 0/0 16 (2d10 + 5) slashing; effect/slashing-damage 16 (2d10 + 5), effect/force-damage 11 (2d10)",
		"action|teleport|Teleport; effect/teleport, requirement/unoccupied-visible-space",
		"bonus|horror-nimbus|Horror Nimbus (Recharge 5–6); wis DC 15; emanation 15/0/0; limit/recharge, save/wis, area/emanation, effect/psychic-damage 28 (8d6), condition/frightened, effect/ends-on-damage-or-lost-line-of-sight, save/success-immunity",
		"reaction|pursuit|Pursuit; reaction/trigger, effect/teleport, requirement/destination-near-triggering-creature",
	},
	"Night Hag": {
		"trait|trait/coven-magic|Coven Magic; effect/coven-casting, requirement/two-hag-allies-within-30-feet, spellcasting/no-material-components, spellcasting/normal-casting-time, limit/long-rest-per-spell, spell/augury, spell/find-familiar, spell/identify, spell/locate-object, spell/scrying, spell/unseen-servant",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"trait|trait/soul-bag|Soul Bag; effect/soul-bag, requirement/holding-or-carrying-soul-bag, effect/soul-bag-statistics, effect/soul-bag-destruction, effect/soul-bag-replacement",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|claw|Claw; melee +7 reach 5 range 0/0 13 (2d8 + 4) slashing; effect/slashing-damage 13 (2d8 + 4)",
		"action|nightmare-haunting|Nightmare Haunting (1/Day); limit/per-day, requirement/soul-bag, requirement/on-ethereal-plane, action/spellcasting, spell/dream, effect/hag-must-be-messenger, requirement/visible-target-on-material-plane, spell/protection-from-evil-and-good, spell/magic-circle, effect/hit-point-maximum-reduction, effect/trap-soul, effect/no-revival-until-soul-released",
		"action|spellcasting|Spellcasting; int DC 14 +0; at will: Detect Magic, Etherealness, Magic Missile (level 4 version); 2/day: Phantasmal Killer, Plane Shift (self only); spellcasting/no-material-components, spell/detect-magic [spellcasting/at-will], spell/etherealness [spellcasting/at-will], spell/magic-missile [spellcasting/at-will], spell/phantasmal-killer [limit/per-day], spell/plane-shift [limit/per-day]",
		"bonus|shape-shift|Shape-Shift; effect/shape-shift, effect/shapechange-form-limits, effect/equipment-not-transformed",
	},
	"Nightmare": {
		"trait|trait/confer-fire-resistance|Confer Fire Resistance; effect/grant-fire-resistance, requirement/rider",
		"trait|trait/illumination|Illumination; effect/shed-light",
		"action|hooves|Hooves; melee +6 reach 5 range 0/0 13 (2d8 + 4) bludgeoning; effect/bludgeoning-damage 13 (2d8 + 4), effect/fire-damage 10 (3d6)",
		"action|ethereal-stride|Ethereal Stride; target/self-and-willing-creatures, effect/plane-shift-to-or-from-ethereal",
	},
	"Noble": {
		"action|rapier|Rapier; melee +3 reach 5 range 0/0 5 (1d8 + 1) piercing; effect/piercing-damage 5 (1d8 + 1)",
		"reaction|parry|Parry; reaction/trigger, requirement/holding-weapon, effect/ac-bonus",
	},
	"Ochre Jelly": {
		"trait|trait/amorphous|Amorphous; movement/squeeze",
		"trait|trait/spider-climb|Spider Climb; effect/climb-without-check",
		"action|pseudopod|Pseudopod; melee +4 reach 5 range 0/0 12 (3d6 + 2) acid; effect/acid-damage 12 (3d6 + 2)",
		"reaction|split|Split; reaction/trigger, requirement/size-large-or-medium, requirement/minimum-hit-points, effect/split, effect/smaller-size, effect/divide-hit-points",
	},
	"Octopus": {
		"trait|trait/compression|Compression; movement/squeeze",
		"trait|trait/water-breathing|Water Breathing; effect/breathe-only-underwater",
		"action|tentacles|Tentacles; melee +4 reach 5 range 0/0 1 bludgeoning; effect/bludgeoning-damage [damage/fixed]",
		"reaction|ink-cloud|Ink Cloud (1/Day); limit/per-day, reaction/trigger, requirement/underwater, effect/ink-cloud, effect/move-swim-speed, effect/heavily-obscured",
	},
	"Ogre": {
		"action|greatclub|Greatclub; melee +6 reach 5 range 0/0 13 (2d8 + 4) bludgeoning; effect/bludgeoning-damage 13 (2d8 + 4)",
		"action|javelin|Javelin; effect/piercing-damage 11 (2d6 + 4)",
	},
	"Ogre Zombie": {
		"trait|trait/undead-fortitude|Undead Fortitude; requirement/reduced-to-zero-hit-points, save/con, effect/damage-based-save-dc, requirement/not-radiant-or-critical-hit, effect/drop-to-one-hit-point",
		"action|slam|Slam; melee +6 reach 5 range 0/0 13 (2d8 + 4) bludgeoning; effect/bludgeoning-damage 13 (2d8 + 4)",
	},
	"Oni": {
		"trait|trait/regeneration|Regeneration; effect/regain-hit-points, requirement/at-least-one-hit-point",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/replace-attack, action/spellcasting",
		"action|claw|Claw; melee +7 reach 10 range 0/0 10 (1d12 + 4) slashing; effect/slashing-damage 10 (1d12 + 4), effect/necrotic-damage 9 (2d8)",
		"action|nightmare-ray|Nightmare Ray; ranged +5 reach 0 range 60/60 9 (2d6 + 2) psychic; effect/psychic-damage 9 (2d6 + 2), condition/frightened",
		"action|shape-shift|Shape-Shift; effect/shape-shift, effect/shapechange-form-limits, effect/equipment-not-transformed",
		"action|spellcasting|Spellcasting; cha DC 13 +0; at will: ; 1/day: Charm Person (level 2 version), Darkness, Gaseous Form, Sleep; spellcasting/no-material-components, spell/charm-person [limit/per-day], spell/darkness [limit/per-day], spell/gaseous-form [limit/per-day], spell/sleep [limit/per-day]",
		"bonus|invisibility|Invisibility; action/spellcasting, spell/invisibility, spellcasting/no-components",
	},
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters168SourceValues(t *testing.T) {
	c := loadMonsters168(t)
	if len(c.Monsters) != 12 || len(monsters168Headers) != 12 || len(monsters168Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #168") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters168Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader164(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries164(m); !slices.Equal(got, monsters168Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters168Entries[want.name])
			}
			if m.ResistanceNote != "" {
				t.Fatal("source lists no resistance note")
			}
		})
		for _, line := range monsters168Entries[want.name] {
			if !strings.HasPrefix(line, "legendary-uses|") {
				entries++
			}
		}
	}
	// Bold-italic entry headings counted on PDF pp. 309–313, 344 and 357–358.
	if entries != 57 {
		t.Fatalf("reviewed %d entries, want 57", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters168MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 168 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters168(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Animals appendix blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("animals monsters monsters monsters monsters monsters monsters monsters animals monsters monsters monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #168 cites the page each block starts on, in this order.
	names := []string{"Mule", "Mummy", "Mummy Lord", "Nalfeshnee", "Night Hag", "Nightmare", "Noble", "Ochre Jelly", "Octopus", "Ogre", "Ogre Zombie", "Oni"}
	pages := []int{357, 309, 309, 310, 311, 311, 312, 312, 357, 312, 344, 312}
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
func TestMonsters168DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"1/8": 2, "3": 2, "15": 5, "13": 5, "5": 3, "2": 2, "0": 2, "7": 3}
	skillAbility := map[string]string{"arcana": "int", "deception": "cha", "history": "int", "insight": "wis", "perception": "wis", "persuasion": "cha", "religion": "int", "stealth": "dex"}
	hitDie := map[string]int{"small": 6, "medium": 8, "large": 10}
	// Attacks that do not use the better of Strength and Dexterity.
	attackAbility := map[string]string{"Oni/nightmare-ray": "cha"}
	for _, m := range loadMonsters168(t).Monsters {
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
func TestMonsters168DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters168(t).Monsters {
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
func TestMonsters168RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters168(t).Monsters {
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

func TestMonsters168MalformedFixtures(t *testing.T) {
	tests := []struct {
		name, record string
		mutate       func(*Monster)
	}{
		{"missing languages list", "Mummy", func(m *Monster) { m.Languages = nil }},
		{"animal block in another family", "Mummy", func(m *Monster) {
			m.ID = strings.Replace(strings.Replace(m.ID, "/animals/", "/spells/", 1), "/monsters/", "/spells/", 1)
		}},
		{"XP off CR", "Oni", func(m *Monster) { m.XP++ }},
		{"save off modifier", "Oni", func(m *Monster) {
			m.Abilities["wis"] = MonsterAbility{m.Abilities["wis"].Score, m.Abilities["wis"].Modifier, m.Abilities["wis"].Modifier + 7}
		}},
		{"modifier", "Mule", func(m *Monster) {
			m.Abilities["str"] = MonsterAbility{m.Abilities["str"].Score, m.Abilities["str"].Modifier + 1, m.Abilities["str"].Modifier + 1}
		}},
		{"HP average", "Mule", func(m *Monster) { m.HP.Average++ }},
		{"legendary uses without actions", "Mummy Lord", func(m *Monster) { m.LegendaryActions = nil }},
		{"lair uses not above uses", "Mummy Lord", func(m *Monster) { m.Legendary.LairUses = 3 }},
		{"lair XP not above XP", "Mummy Lord", func(m *Monster) { m.LairXP = m.XP }},
		{"hover without flight", "Mule", func(m *Monster) { m.Hover = true }},
		{"recharge", "Nalfeshnee", func(m *Monster) { m.BonusActions[0].Uses.Recharge = 1 }},
		{"save DC off entry text", "Mummy", func(m *Monster) { m.Actions[2].Save.DC = 12 }},
		{"spell DC off entry text", "Mummy Lord", func(m *Monster) { m.Actions[4].Spellcasting.SaveDC = 16 }},
		{"missing at-will list", "Oni", func(m *Monster) { m.Actions[4].Spellcasting.AtWill = nil }},
		{"missing daily spell groups", "Oni", func(m *Monster) { m.Actions[4].Spellcasting.PerDay = nil }},
		{"negative spell attack bonus", "Mummy Lord", func(m *Monster) { m.Actions[4].Spellcasting.AttackBonus = -1 }},
		{"fixed damage with sides", "Octopus", func(m *Monster) { m.Actions[0].Attack.Damage.Sides = 4 }},
		{"empty reactions", "Noble", func(m *Monster) { m.Reactions = []MonsterAction{} }},
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
		c := loadMonsters168(t)
		c.Monsters[1].ID = c.Monsters[2].ID
		rejected(t, c, true)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters168(t)
			tt.mutate(record(t, &c, tt.record))
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			rejected(t, c, !strings.HasPrefix(tt.name, "empty "))
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_168.json")
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
		{`"uses":3,`, ``},
		{`"dc":17`, `"dc":17,"dc":18`},
		{`"width":0,`, ``},
		{`"save_dc":17,`, ``},
		{`"attack_bonus":9`, `"attack_bonus":9,"attack_bonus":10`},
		{`"truesight":60`, `"truesight":null`},
		{`"shape":"emanation",`, `"shape":"emanation","height":5,`},
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
// to 17 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters168ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters168(t)
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
	for _, load := range []func(*testing.T) MonsterCatalogue{loadMonsters158, loadMonsters159, loadMonsters160, loadMonsters161, loadMonsters162, loadMonsters163, loadMonsters164, loadMonsters165, loadMonsters166, loadMonsters167} {
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
	for _, id := range []string{"channel-negative-energy", "claw", "greatclub", "hooves", "nightmare-ray", "pseudopod", "rapier", "rend", "rotting-fist", "slam", "tentacles"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
