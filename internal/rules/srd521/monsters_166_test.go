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

func loadMonsters166(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_166.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #166 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Lamia's and Lion's blocks run onto the next page.
var monsters166Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Iron Golem", [2]int{302, 302}, "Large Construct, Unaligned|AC 20 Initiative +9 (19)|HP 252 (24d10 + 120)|Speed 30 ft.|Str 24 +7 +7 Dex 9 −1 −1 Con 20 +5 +5|Int 3 −4 −4 Wis 11 +0 +0 Cha 1 −5 −5|Immunities Fire, Poison, Psychic; Charmed, Exhaustion, Frightened, Paralyzed, Petrified, Poisoned|Senses Darkvision 120 ft.; Passive Perception 10|Languages Understands Common plus two other languages but can’t speak|CR 16 (XP 15,000; PB +5)"},
	{"Jackal", [2]int{356, 356}, "Small Beast, Unaligned|AC 12 Initiative +2 (12)|HP 3 (1d6)|Speed 40 ft.|Str 8 −1 −1 Dex 15 +2 +2 Con 11 +0 +0|Int 3 −4 −4 Wis 12 +1 +1 Cha 6 −2 −2|Skills Perception +5, Stealth +4|Senses Darkvision 90 ft.; Passive Perception 15|Languages None|CR 0 (XP 10; PB +2)"},
	{"Killer Whale", [2]int{356, 356}, "Huge Beast, Unaligned|AC 12 Initiative +2 (12)|HP 90 (12d12 + 12)|Speed 5 ft., Swim 60 ft.|Str 19 +4 +4 Dex 14 +2 +2 Con 13 +1 +1|Int 3 −4 −4 Wis 12 +1 +1 Cha 7 −2 −2|Skills Perception +3, Stealth +4|Senses Blindsight 120 ft.; Passive Perception 13|Languages None|CR 3 (XP 700; PB +2)"},
	{"Knight", [2]int{302, 302}, "Medium or Small Humanoid, Neutral|AC 18 Initiative +0 (10)|HP 52 (8d8 + 16)|Speed 30 ft.|Str 16 +3 +3 Dex 11 +0 +0 Con 14 +2 +4|Int 11 +0 +0 Wis 11 +0 +2 Cha 15 +2 +2|Immunities Frightened|Gear Greatsword, Heavy Crossbow, Plate Armor|Senses Passive Perception 10|Languages Common plus one other language|CR 3 (XP 700; PB +2)"},
	{"Kobold Warrior", [2]int{302, 302}, "Small Dragon, Neutral|AC 14 Initiative +2 (12)|HP 7 (3d6 − 3)|Speed 30 ft.|Str 7 −2 −2 Dex 15 +2 +2 Con 9 −1 −1|Int 8 −1 −1 Wis 7 −2 −2 Cha 8 −1 −1|Gear Daggers (3)|Senses Darkvision 60 ft.; Passive Perception 8|Languages Common, Draconic|CR 1/8 (XP 25; PB +2)"},
	{"Kraken", [2]int{303, 303}, "Gargantuan Monstrosity (Titan), Chaotic Evil|AC 18 Initiative +14 (24)|HP 481 (26d20 + 208)|Speed 30 ft., Swim 120 ft.|Str 30 +10 +17 Dex 11 +0 +7 Con 26 +8 +15|Int 22 +6 +6 Wis 18 +4 +11 Cha 20 +5 +5|Skills History +13, Perception +11|Immunities Cold, Lightning; Frightened, Grappled, Paralyzed, Restrained|Senses Truesight 120 ft.; Passive Perception 21|Languages Understands Abyssal, Celestial, Infernal, and Primordial but can’t speak; telepathy 120 ft.|CR 23 (XP 50,000, or 62,000 in lair; PB +7)"},
	{"Lamia", [2]int{303, 304}, "Large Fiend, Chaotic Evil|AC 13 Initiative +1 (11)|HP 97 (13d10 + 26)|Speed 40 ft.|Str 16 +3 +3 Dex 13 +1 +1 Con 15 +2 +2|Int 14 +2 +2 Wis 15 +2 +2 Cha 16 +3 +3|Skills Deception +7, Insight +4, Stealth +5|Senses Darkvision 60 ft.; Passive Perception 12|Languages Abyssal, Common|CR 4 (XP 1,100; PB +2)"},
	{"Lemure", [2]int{304, 304}, "Medium Fiend (Devil), Lawful Evil|AC 9 Initiative −3 (7)|HP 9 (2d8)|Speed 20 ft.|Str 10 +0 +0 Dex 5 −3 −3 Con 11 +0 +0|Int 1 −5 −5 Wis 11 +0 +0 Cha 3 −4 −4|Resistances Cold|Immunities Fire, Poison; Charmed, Frightened, Poisoned|Senses Darkvision 120 ft. (unimpeded by magical Darkness); Passive Perception 10|Languages Understands Infernal but can’t speak|CR 0 (XP 10; PB +2)"},
	{"Lich", [2]int{304, 304}, "Medium Undead (Wizard), Neutral Evil|AC 20 Initiative +17 (27)|HP 315 (42d8 + 126)|Speed 30 ft.|Str 11 +0 +0 Dex 16 +3 +10 Con 16 +3 +10|Int 21 +5 +12 Wis 14 +2 +9 Cha 16 +3 +3|Skills Arcana +19, History +12, Insight +9, Perception +9|Resistances Cold, Lightning|Immunities Necrotic, Poison; Charmed, Exhaustion, Frightened, Paralyzed, Poisoned|Gear Component Pouch|Senses Truesight 120 ft.; Passive Perception 19|Languages All|CR 21 (XP 33,000, or 41,000 in lair; PB +7)"},
	{"Lion", [2]int{356, 357}, "Large Beast, Unaligned|AC 12 Initiative +2 (12)|HP 22 (4d10)|Speed 50 ft.|Str 17 +3 +3 Dex 15 +2 +2 Con 11 +0 +0|Int 3 −4 −4 Wis 12 +1 +1 Cha 8 −1 −1|Skills Perception +3, Stealth +4|Senses Darkvision 60 ft.; Passive Perception 13|Languages None|CR 1 (XP 200; PB +2)"},
	{"Lizard", [2]int{357, 357}, "Tiny Beast, Unaligned|AC 10 Initiative +0 (10)|HP 2 (1d4)|Speed 20 ft., Climb 20 ft.|Str 2 −4 −4 Dex 11 +0 +0 Con 10 +0 +0|Int 1 −5 −5 Wis 8 −1 −1 Cha 3 −4 −4|Senses Darkvision 30 ft.; Passive Perception 9|Languages None|CR 0 (XP 10; PB +2)"},
	{"Mage", [2]int{305, 305}, "Medium or Small Humanoid (Wizard), Neutral|AC 15 Initiative +2 (12)|HP 81 (18d8)|Speed 30 ft.|Str 9 −1 −1 Dex 14 +2 +2 Con 11 +0 +0|Int 17 +3 +6 Wis 12 +1 +4 Cha 11 +0 +0|Skills Arcana +6, History +6, Perception +4|Gear Wand|Senses Passive Perception 14|Languages Common plus three other languages|CR 6 (XP 2,300; PB +3)"},
}

// Typed values of every trait and action entry, reviewed against the same pages:
// limited uses, attack bonus/reach/range/damage, save ability and DC, target
// shape size/width/range, dice and dependency keys. A ranged attack printed with
// one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text. The Kobold Warrior's Dagger, Lich's Eldritch Burst and Mage's Arcane Burst are
// "Melee or Ranged Attack Roll" entries (key action/melee-or-ranged-attack). The Iron Golem's
// Fiery Bolt prints +10 with no damage modifier, which no ability score plus PB explains.
var monsters166Entries = map[string][]string{
	"Iron Golem": {
		"trait|trait/fire-absorption|Fire Absorption; effect/absorb-fire-damage",
		"trait|trait/immutable-form|Immutable Form; effect/no-shape-shift",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|bladed-arm|Bladed Arm; melee +12 reach 10 range 0/0 20 (3d8 + 7) slashing; effect/slashing-damage 20 (3d8 + 7), effect/fire-damage 10 (3d6)",
		"action|fiery-bolt|Fiery Bolt; ranged +10 reach 0 range 120/120 36 (8d8) fire; effect/fire-damage 36 (8d8)",
		"action|poison-breath|Poison Breath (Recharge 6); con DC 18; cone 60/0/0; limit/recharge, save/con, area/cone, effect/poison-damage 55 (10d10), save/half-damage",
	},
	"Jackal": {
		"action|bite|Bite; melee +1 reach 5 range 0/0 1 (1d4 − 1) piercing; effect/piercing-damage 1 (1d4 − 1)",
	},
	"Killer Whale": {
		"trait|trait/hold-breath|Hold Breath; effect/hold-breath",
		"action|bite|Bite; melee +6 reach 5 range 0/0 21 (5d6 + 4) piercing; effect/piercing-damage 21 (5d6 + 4)",
	},
	"Knight": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|greatsword|Greatsword; melee +5 reach 5 range 0/0 10 (2d6 + 3) slashing; effect/slashing-damage 10 (2d6 + 3), effect/radiant-damage 4 (1d8)",
		"action|heavy-crossbow|Heavy Crossbow; ranged +2 reach 0 range 100/400 11 (2d10) piercing; effect/piercing-damage 11 (2d10), effect/radiant-damage 4 (1d8)",
		"reaction|parry|Parry; reaction/trigger, requirement/holding-weapon, effect/ac-bonus",
	},
	"Kobold Warrior": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage, requirement/ally-within-5-feet, requirement/ally-not-incapacitated",
		"trait|trait/sunlight-sensitivity|Sunlight Sensitivity; requirement/in-sunlight, check/disadvantage, attack/disadvantage",
		"action|dagger|Dagger; effect/piercing-damage 4 (1d4 + 2)",
	},
	"Kraken": {
		"trait|trait/amphibious|Amphibious; effect/breathe-air-and-water",
		"trait|trait/legendary-resistance|Legendary Resistance (4/Day, or 5/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"trait|trait/siege-monster|Siege Monster; effect/double-damage-to-objects",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/additional-action, action/saving-throw",
		"action|tentacle|Tentacle; melee +17 reach 30 range 0/0 24 (4d6 + 10) bludgeoning; effect/bludgeoning-damage 24 (4d6 + 10), condition/grappled, limit/grapple-from-one-of-ten-tentacles, condition/restrained",
		"action|fling|Fling; dex DC 25; requirement/grappled-creature, effect/throw-creature, save/dex, target/thrown-creature-and-destination-space, effect/bludgeoning-damage 18 (4d8), condition/prone, save/half-damage",
		"action|lightning-strike|Lightning Strike; dex DC 23; creature 0/0/120; save/dex, target/creature, effect/lightning-damage 33 (6d10), save/half-damage",
		"action|swallow|Swallow; dex DC 25; creature 0/0/0; save/dex, target/creature, requirement/grappled-creature, limit/four-swallowed-creatures, effect/piercing-damage 23 (3d8 + 10), requirement/target-size, effect/swallow, condition/restrained, effect/total-cover, effect/acid-damage 24 (7d6), save/con-or-regurgitate, condition/prone, effect/escape-corpse",
		"legendary|storm-bolt|Storm Bolt; action/saving-throw",
		"legendary|toxic-ink|Toxic Ink; con DC 23; emanation 15/0/0; save/con, area/emanation, requirement/underwater, condition/blinded, condition/poisoned, effect/move-speed, limit/once-per-round",
		"legendary-uses|3|4",
	},
	"Lamia": {
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, action/saving-throw",
		"action|claw|Claw; melee +5 reach 5 range 0/0 7 (1d8 + 3) slashing; effect/slashing-damage 7 (1d8 + 3), effect/psychic-damage 7 (2d6)",
		"action|corrupting-touch|Corrupting Touch; wis DC 13; creature 0/0/5; save/wis, target/creature, effect/psychic-damage 13 (3d8), effect/curse, condition/charmed, condition/poisoned",
		"action|spellcasting|Spellcasting; cha DC 13 +0; at will: Disguise Self (can appear as a Large or Medium biped), Minor Illusion; 1/day: Geas, Major Image, Scrying; spellcasting/no-material-components, spell/disguise-self [spellcasting/at-will], spell/minor-illusion [spellcasting/at-will], spell/geas [limit/per-day], spell/major-image [limit/per-day], spell/scrying [limit/per-day]",
		"bonus|leap|Leap; effect/jump",
	},
	"Lemure": {
		"trait|trait/hellish-restoration|Hellish Restoration; effect/revive, effect/revival-delay, requirement/not-killed-under-bless, requirement/remains-not-sprinkled-with-holy-water, spell/bless",
		"action|vile-slime|Vile Slime; melee +2 reach 5 range 0/0 2 (1d4) poison; effect/poison-damage 2 (1d4)",
	},
	"Lich": {
		"trait|trait/legendary-resistance|Legendary Resistance (4/Day, or 5/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"trait|trait/spirit-jar|Spirit Jar; effect/revive, effect/revival-delay, requirement/spirit-jar, effect/reform-in-lair",
		"action|multiattack|Multiattack; action/melee-or-ranged-attack, action/melee-attack, multiattack/any-combination",
		"action|eldritch-burst|Eldritch Burst; effect/force-damage 31 (4d12 + 5)",
		"action|paralyzing-touch|Paralyzing Touch; melee +12 reach 5 range 0/0 15 (3d6 + 5) cold; effect/cold-damage 15 (3d6 + 5), condition/paralyzed",
		"action|spellcasting|Spellcasting; int DC 20 +0; at will: Detect Magic, Detect Thoughts, Dispel Magic, Fireball (level 5 version), Invisibility, Lightning Bolt (level 5 version), Mage Hand, Prestidigitation; 2/day: Animate Dead, Dimension Door, Plane Shift; 1/day: Chain Lightning, Finger of Death, Power Word Kill, Scrying; spell/detect-magic [spellcasting/at-will], spell/detect-thoughts [spellcasting/at-will], spell/dispel-magic [spellcasting/at-will], spell/fireball [spellcasting/at-will], spell/invisibility [spellcasting/at-will], spell/lightning-bolt [spellcasting/at-will], spell/mage-hand [spellcasting/at-will], spell/prestidigitation [spellcasting/at-will], spell/animate-dead [limit/per-day], spell/dimension-door [limit/per-day], spell/plane-shift [limit/per-day], spell/chain-lightning [limit/per-day], spell/finger-of-death [limit/per-day], spell/power-word-kill [limit/per-day], spell/scrying [limit/per-day]",
		"reaction|protective-magic|Protective Magic; action/spellcasting, spell/counterspell, spell/shield, reaction/spell-trigger",
		"legendary|deathly-teleport|Deathly Teleport; effect/teleport, requirement/unoccupied-visible-space, target/creatures-near-departure-space, effect/necrotic-damage 11 (2d10)",
		"legendary|disrupt-life|Disrupt Life; con DC 20; emanation 20/0/0; save/con, area/emanation, requirement/not-undead, effect/necrotic-damage 31 (9d6), save/half-damage, limit/once-per-round",
		"legendary|frightening-gaze|Frightening Gaze; action/spellcasting, spell/fear, limit/once-per-round",
		"legendary-uses|3|4",
	},
	"Lion": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage, requirement/ally-within-5-feet, requirement/ally-not-incapacitated",
		"trait|trait/running-leap|Running Leap; effect/long-jump-distance, requirement/running-start",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, action/saving-throw",
		"action|rend|Rend; melee +5 reach 5 range 0/0 7 (1d8 + 3) slashing; effect/slashing-damage 7 (1d8 + 3)",
		"action|roar|Roar; wis DC 11; creature 0/0/15; save/wis, target/creature, condition/frightened",
	},
	"Lizard": {
		"trait|trait/spider-climb|Spider Climb; effect/climb-without-check",
		"action|bite|Bite; melee +2 reach 5 range 0/0 1 piercing; effect/piercing-damage [damage/fixed]",
	},
	"Mage": {
		"action|multiattack|Multiattack; action/melee-or-ranged-attack",
		"action|arcane-burst|Arcane Burst; effect/force-damage 16 (3d8 + 3)",
		"action|spellcasting|Spellcasting; int DC 14 +0; at will: Detect Magic, Light, Mage Armor (included in AC), Mage Hand, Prestidigitation; 2/day: Fireball (level 4 version), Invisibility; 1/day: Cone of Cold, Fly; spell/detect-magic [spellcasting/at-will], spell/light [spellcasting/at-will], spell/mage-armor [spellcasting/at-will], spell/mage-hand [spellcasting/at-will], spell/prestidigitation [spellcasting/at-will], spell/fireball [limit/per-day], spell/invisibility [limit/per-day], spell/cone-of-cold [limit/per-day], spell/fly [limit/per-day]",
		"bonus|misty-step|Misty Step (3/Day); limit/per-day, action/spellcasting, spell/misty-step",
		"reaction|protective-magic|Protective Magic (3/Day); limit/per-day, action/spellcasting, spell/counterspell, spell/shield, reaction/spell-trigger",
	},
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters166SourceValues(t *testing.T) {
	c := loadMonsters166(t)
	if len(c.Monsters) != 12 || len(monsters166Headers) != 12 || len(monsters166Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #166") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters166Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader164(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries164(m); !slices.Equal(got, monsters166Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters166Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.ResistanceNote != "" || m.Hover {
				t.Fatal("source lists no vulnerabilities, resistance note, hover")
			}
		})
		for _, line := range monsters166Entries[want.name] {
			if !strings.HasPrefix(line, "legendary-uses|") {
				entries++
			}
		}
	}
	// Bold-italic entry headings counted on PDF pp. 302–305 and 356–357.
	if entries != 56 {
		t.Fatalf("reviewed %d entries, want 56", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters166MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 166 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters166(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Animals appendix blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("monsters animals animals monsters monsters monsters monsters monsters monsters animals animals monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #166 cites the page each block starts on, in this order.
	names := []string{"Iron Golem", "Jackal", "Killer Whale", "Knight", "Kobold Warrior", "Kraken", "Lamia", "Lemure", "Lich", "Lion", "Lizard", "Mage"}
	pages := []int{302, 356, 356, 302, 302, 303, 303, 304, 304, 356, 357, 305}
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
func TestMonsters166DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"16": 5, "0": 2, "3": 2, "1/8": 2, "23": 7, "4": 2, "21": 7, "1": 2, "6": 3}
	skillAbility := map[string]string{"arcana": "int", "deception": "cha", "history": "int", "insight": "wis", "perception": "wis", "stealth": "dex"}
	hitDie := map[string]int{"tiny": 4, "small": 6, "medium": 8, "large": 10, "huge": 12, "gargantuan": 20}
	// Attacks that do not use the better of Strength and Dexterity.
	attackAbility := map[string]string{"Jackal/bite": "str", "Knight/heavy-crossbow": "dex", "Lich/paralyzing-touch": "int"}
	// Printed attack bonus and damage modifier that no ability explains.
	printed := map[string][2]int{"Iron Golem/fiery-bolt": {10, 0}}
	for _, m := range loadMonsters166(t).Monsters {
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
				if want, ok := printed[m.Name+"/"+a.ID]; ok {
					if a.Attack.Bonus != want[0] || a.Attack.Damage.Modifier != want[1] {
						t.Fatalf("%s differs from its printed values", a.Name)
					}
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
func TestMonsters166DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters166(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			if m.RequiredMechanics[0] != "monster/stat-block" {
				t.Fatal("stat-block dependency must come first")
			}
			quoted := func(parent string, x MonsterMechanic) {
				if !strings.Contains(parent, strings.TrimSuffix(x.Text, ".")) {
					t.Fatalf("%q is not quoted from its entry", x.Text)
				}
				// The Jackal's Bite prints its negative modifier with an en dash (p. 356).
				if x.Dice != nil && !strings.Contains(strings.ReplaceAll(x.Text, " – ", " − "), monsterDice163(*x.Dice)) {
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
func TestMonsters166RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters166(t).Monsters {
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

func TestMonsters166MalformedFixtures(t *testing.T) {
	tests := []struct {
		name, record string
		mutate       func(*Monster)
	}{
		{"missing languages list", "Jackal", func(m *Monster) { m.Languages = nil }},
		{"animal block in another family", "Jackal", func(m *Monster) {
			m.ID = strings.Replace(strings.Replace(m.ID, "/animals/", "/spells/", 1), "/monsters/", "/spells/", 1)
		}},
		{"XP off CR", "Mage", func(m *Monster) { m.XP++ }},
		{"save off modifier", "Mage", func(m *Monster) {
			m.Abilities["wis"] = MonsterAbility{m.Abilities["wis"].Score, m.Abilities["wis"].Modifier, m.Abilities["wis"].Modifier + 7}
		}},
		{"modifier", "Iron Golem", func(m *Monster) {
			m.Abilities["str"] = MonsterAbility{m.Abilities["str"].Score, m.Abilities["str"].Modifier + 1, m.Abilities["str"].Modifier + 1}
		}},
		{"HP average", "Iron Golem", func(m *Monster) { m.HP.Average++ }},
		{"legendary uses without actions", "Kraken", func(m *Monster) { m.LegendaryActions = nil }},
		{"legendary actions without uses", "Lich", func(m *Monster) { m.Legendary = nil }},
		{"lair uses not above uses", "Kraken", func(m *Monster) { m.Legendary.LairUses = 3 }},
		{"lair XP not above XP", "Lich", func(m *Monster) { m.LairXP = m.XP }},
		{"lair daily uses not above daily uses", "Lich", func(m *Monster) { m.Traits[0].Uses.LairPerDay = 4 }},
		{"blank legendary text", "Kraken", func(m *Monster) { m.Legendary.Text = " " }},
		{"duplicate action across sections", "Kraken", func(m *Monster) { m.LegendaryActions[0].ID = "tentacle" }},
		{"save ability off entry text", "Kraken", func(m *Monster) { m.LegendaryActions[1].Save.Ability = "wis" }},
		{"unlisted spell", "Lich", func(m *Monster) { m.Actions[3].Spellcasting.AtWill[0] = "Wish" }},
		{"duplicate spell", "Mage", func(m *Monster) { m.Actions[2].Spellcasting.AtWill[0] = "Light" }},
		{"fixed damage with sides", "Lizard", func(m *Monster) { m.Actions[0].Attack.Damage.Sides = 4 }},
		{"empty reactions", "Knight", func(m *Monster) { m.Reactions = []MonsterAction{} }},
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
		c := loadMonsters166(t)
		c.Monsters[1].ID = c.Monsters[2].ID
		rejected(t, c, true)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters166(t)
			tt.mutate(record(t, &c, tt.record))
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			rejected(t, c, !strings.HasPrefix(tt.name, "empty "))
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_166.json")
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
		{`"lair_uses":4`, ``},
		{`"dc":25`, `"dc":25,"dc":26`},
		{`"width":0,`, ``},
		{`"save_dc":20,`, ``},
		{`"attack_bonus":0`, `"attack_bonus":null`},
		{`"truesight":120`, `"truesight":null`},
		{`"lair_per_day":5,`, ``},
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
// to 15 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters166ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters166(t)
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
	for _, load := range []func(*testing.T) MonsterCatalogue{loadMonsters158, loadMonsters159, loadMonsters160, loadMonsters161, loadMonsters162, loadMonsters163, loadMonsters164, loadMonsters165} {
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
	for _, id := range []string{"bite", "bladed-arm", "claw", "fiery-bolt", "greatsword", "heavy-crossbow", "paralyzing-touch", "rend", "tentacle", "vile-slime"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
