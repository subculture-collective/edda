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

func loadMonsters172(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_172.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #172 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Shadow’s, the Skeleton’s, the Solar’s and the Sphinx of Lore’s blocks run onto the next
// page.
var monsters172Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Scout", [2]int{322, 322}, "Medium or Small Humanoid, Neutral|AC 13 Initiative +2 (12)|HP 16 (3d8 + 3)|Speed 30 ft.|Str 11 +0 +0 Dex 14 +2 +2 Con 12 +1 +1|Int 11 +0 +0 Wis 13 +1 +1 Cha 11 +0 +0|Skills Nature +4, Perception +5, Stealth +6, Survival +5|Gear Leather Armor, Longbow, Shortsword|Senses Passive Perception 15|Languages Common plus one other language|CR 1/2 (XP 100; PB +2)"},
	{"Sea Hag", [2]int{322, 322}, "Medium Fey, Chaotic Evil|AC 14 Initiative +1 (11)|HP 52 (7d8 + 21)|Speed 30 ft., Swim 40 ft.|Str 16 +3 +3 Dex 13 +1 +1 Con 16 +3 +3|Int 12 +1 +1 Wis 12 +1 +1 Cha 13 +1 +1|Senses Darkvision 60 ft.; Passive Perception 11|Languages Common, Giant, Primordial (Aquan)|CR 2 (XP 450; PB +2)"},
	{"Seahorse", [2]int{361, 361}, "Tiny Beast, Unaligned|AC 12 Initiative +1 (11)|HP 1 (1d4 − 1)|Speed 5 ft., Swim 20 ft.|Str 1 −5 −5 Dex 12 +1 +1 Con 8 −1 −1|Int 1 −5 −5 Wis 10 +0 +0 Cha 2 −4 −4|Skills Perception +2, Stealth +5|Senses Passive Perception 12|Languages None|CR 0 (XP 0; PB +2)"},
	{"Shadow", [2]int{322, 323}, "Medium Undead, Chaotic Evil|AC 12 Initiative +2 (12)|HP 27 (5d8 + 5)|Speed 40 ft.|Str 6 −2 −2 Dex 14 +2 +2 Con 13 +1 +1|Int 6 −2 −2 Wis 10 +0 +0 Cha 8 −1 −1|Skills Stealth +6|Vulnerabilities Radiant|Resistances Acid, Cold, Fire, Lightning, Thunder|Immunities Necrotic, Poison; Exhaustion, Frightened, Grappled, Paralyzed, Petrified, Poisoned, Prone, Restrained, Unconscious|Senses Darkvision 60 ft.; Passive Perception 10|Languages None|CR 1/2 (XP 100; PB +2)"},
	{"Shambling Mound", [2]int{323, 323}, "Large Plant, Unaligned|AC 15 Initiative −1 (9)|HP 110 (13d10 + 39)|Speed 30 ft., Swim 20 ft.|Str 18 +4 +4 Dex 8 −1 −1 Con 16 +3 +3|Int 5 −3 −3 Wis 10 +0 +0 Cha 5 −3 −3|Skills Stealth +3|Resistances Cold, Fire|Immunities Lightning; Deafened, Exhaustion|Senses Blindsight 60 ft.; Passive Perception 10|Languages None|CR 5 (XP 1,800; PB +3)"},
	{"Shield Guardian", [2]int{323, 323}, "Large Construct, Unaligned|AC 17 Initiative −1 (9)|HP 142 (15d10 + 60)|Speed 30 ft.|Str 18 +4 +4 Dex 8 −1 −1 Con 18 +4 +4|Int 7 −2 −2 Wis 10 +0 +0 Cha 3 −4 −4|Immunities Poison; Charmed, Exhaustion, Frightened, Paralyzed, Petrified, Poisoned|Senses Blindsight 10 ft., Darkvision 60 ft.; Passive Perception 10|Languages Understands commands given in any language but can’t speak|CR 7 (XP 2,900; PB +3)"},
	{"Shrieker Fungus", [2]int{286, 286}, "Medium Plant, Unaligned|AC 5 Initiative −5 (5)|HP 13 (3d8)|Speed 5 ft.|Str 1 −5 −5 Dex 1 −5 −5 Con 10 +0 +0|Int 1 −5 −5 Wis 3 −4 −4 Cha 1 −5 −5|Immunities Blinded, Charmed, Deafened, Frightened|Senses Blindsight 30 ft.; Passive Perception 6|Languages None|CR 0 (XP 0; PB +2)"},
	{"Silver Dragon Wyrmling", [2]int{324, 324}, "Medium Dragon (Metallic), Lawful Good|AC 17 Initiative +2 (12)|HP 45 (6d8 + 18)|Speed 30 ft., Fly 60 ft.|Str 19 +4 +4 Dex 10 +0 +2 Con 17 +3 +3|Int 12 +1 +1 Wis 11 +0 +2 Cha 15 +2 +2|Skills Perception +4, Stealth +2|Immunities Cold|Senses Blindsight 10 ft., Darkvision 60 ft.; Passive Perception 14|Languages Draconic|CR 2 (450 XP; PB +2)"},
	{"Skeleton", [2]int{325, 326}, "Medium Undead, Lawful Evil|AC 14 Initiative +3 (13)|HP 13 (2d8 + 4)|Speed 30 ft.|Str 10 +0 +0 Dex 16 +3 +3 Con 15 +2 +2|Int 6 −2 −2 Wis 8 −1 −1 Cha 5 −3 −3|Vulnerabilities Bludgeoning|Immunities Poison; Exhaustion, Poisoned|Gear Shortbow, Shortsword|Senses Darkvision 60 ft.; Passive Perception 9|Languages Understands Common plus one other language but can’t speak|CR 1/4 (XP 50; PB +2)"},
	{"Solar", [2]int{326, 327}, "Large Celestial (Angel), Lawful Good|AC 21 Initiative +20 (30)|HP 297 (22d10 + 176)|Speed 50 ft., Fly 150 ft. (hover)|Str 26 +8 +8 Dex 22 +6 +6 Con 26 +8 +8|Int 25 +7 +7 Wis 25 +7 +7 Cha 30 +10 +10|Skills Perception +14|Immunities Poison, Radiant; Charmed, Exhaustion, Frightened, Poisoned|Senses Truesight 120 ft.; Passive Perception 24|Languages All; telepathy 120 ft.|CR 21 (XP 33,000; PB +7)"},
	{"Specter", [2]int{327, 327}, "Medium Undead, Chaotic Evil|AC 12 Initiative +2 (12)|HP 22 (5d8)|Speed 30 ft., Fly 50 ft. (hover)|Str 1 −5 −5 Dex 14 +2 +2 Con 11 +0 +0|Int 10 +0 +0 Wis 10 +0 +0 Cha 11 +0 +0|Resistances Acid, Bludgeoning, Cold, Fire, Lightning, Piercing, Slashing, Thunder|Immunities Necrotic, Poison; Charmed, Exhaustion, Grappled, Paralyzed, Petrified, Poisoned, Prone, Restrained, Unconscious|Senses Darkvision 60 ft.; Passive Perception 10|Languages Understands Common plus one other language but can’t speak|CR 1 (XP 200; PB +2)"},
	{"Sphinx of Lore", [2]int{327, 328}, "Large Celestial, Lawful Neutral|AC 17 Initiative +10 (20)|HP 170 (20d10 + 60)|Speed 40 ft., Fly 60 ft.|Str 18 +4 +4 Dex 15 +2 +2 Con 16 +3 +3|Int 18 +4 +4 Wis 18 +4 +4 Cha 18 +4 +4|Skills Arcana +12, History +12, Perception +8, Religion +12|Resistances Necrotic, Radiant|Immunities Psychic; Charmed, Frightened|Senses Truesight 120 ft.; Passive Perception 18|Languages Celestial, Common|CR 11 (XP 7,200, or 8,400 in lair; PB +4)"},
}

// Typed values of every trait and action entry, reviewed against the same pages:
// limited uses, attack bonus/reach/range/damage, save ability and DC, target
// shape size/width/range, dice and dependency keys. A ranged attack printed with
// one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text. The Solar's Flying Sword is a "Melee or Ranged Attack Roll".
var monsters172Entries = map[string][]string{
	"Scout": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|shortsword|Shortsword; melee +4 reach 5 range 0/0 5 (1d6 + 2) piercing; effect/piercing-damage 5 (1d6 + 2)",
		"action|longbow|Longbow; ranged +4 reach 0 range 150/600 6 (1d8 + 2) piercing; effect/piercing-damage 6 (1d8 + 2)",
	},
	"Sea Hag": {
		"trait|trait/amphibious|Amphibious; effect/breathe-air-and-water",
		"trait|trait/coven-magic|Coven Magic; effect/coven-casting, requirement/two-hag-allies-within-30-feet, spellcasting/no-material-components, spellcasting/normal-casting-time, limit/long-rest-per-spell, spell/augury, spell/find-familiar, spell/identify, spell/locate-object, spell/scrying, spell/unseen-servant",
		"trait|trait/vile-appearance|Vile Appearance; wis DC 11; save/wis, target/beasts-and-humanoids-starting-turn-nearby, requirement/can-see-true-form, condition/frightened, save/success-immunity",
		"action|claw|Claw; melee +5 reach 5 range 0/0 10 (2d6 + 3) slashing; effect/slashing-damage 10 (2d6 + 3)",
		"action|death-glare|Death Glare (Recharge 5–6); wis DC 11; creature 0/0/30; limit/recharge, save/wis, target/creature, requirement/frightened-target, effect/drop-to-zero-hit-points, effect/psychic-damage 13 (3d8)",
		"action|illusory-appearance|Illusory Appearance; action/spellcasting, spell/disguise-self, spellcasting/constitution-ability, spellcasting/fixed-save-dc, effect/spell-duration-24-hours",
	},
	"Seahorse": {
		"trait|trait/water-breathing|Water Breathing; effect/breathe-only-underwater",
		"action|bubble-dash|Bubble Dash; requirement/underwater, movement/swim-speed-move, effect/no-opportunity-attacks",
	},
	"Shadow": {
		"trait|trait/amorphous|Amorphous; movement/squeeze",
		"trait|trait/sunlight-weakness|Sunlight Weakness; effect/sunlight-disadvantage",
		"action|draining-swipe|Draining Swipe; melee +4 reach 5 range 0/0 5 (1d6 + 2) necrotic; effect/necrotic-damage 5 (1d6 + 2), effect/strength-score-reduction, effect/death-at-zero-strength, effect/rising-shadow",
		"bonus|shadow-stealth|Shadow Stealth; requirement/dim-light-or-darkness, action/hide",
	},
	"Shambling Mound": {
		"trait|trait/lightning-absorption|Lightning Absorption; effect/lightning-absorption",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, action/saving-throw",
		"action|charged-tendril|Charged Tendril; melee +7 reach 10 range 0/0 7 (1d6 + 4) bludgeoning; effect/bludgeoning-damage 7 (1d6 + 4), effect/lightning-damage 5 (2d4), requirement/target-size, effect/pull",
		"action|engulf|Engulf; str DC 15; creature 0/0/5; save/str, target/creature, requirement/target-size, effect/engulf, condition/grappled, condition/blinded, condition/restrained, effect/lightning-damage 10 (3d6), effect/engulfed-target-moves-along, limit/one-grappled-creature",
	},
	"Shield Guardian": {
		"trait|trait/bound|Bound; effect/amulet-bond, effect/telepathic-summons, effect/know-amulet-location, effect/damage-transfer",
		"trait|trait/regeneration|Regeneration; effect/regeneration",
		"trait|trait/spell-storing|Spell Storing; effect/store-spell, effect/stored-spell-replaced, effect/cast-stored-spell, spellcasting/no-components, spellcasting/caster-ability",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|fist|Fist; melee +7 reach 10 range 0/0 11 (2d6 + 4) bludgeoning; effect/bludgeoning-damage 11 (2d6 + 4), effect/force-damage 7 (2d6)",
		"reaction|protection|Protection; reaction/trigger, effect/ac-bonus",
	},
	"Shrieker Fungus": {
		"reaction|shriek|Shriek; reaction/trigger, effect/shriek",
	},
	"Silver Dragon Wyrmling": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +6 reach 5 range 0/0 9 (1d10 + 4) piercing; effect/piercing-damage 9 (1d10 + 4)",
		"action|cold-breath|Cold Breath (Recharge 5–6); con DC 13; cone 15/0/0; limit/recharge, save/con, area/cone, effect/cold-damage 18 (4d8), save/half-damage",
		"action|paralyzing-breath|Paralyzing Breath; con DC 13; cone 15/0/0; save/con, area/cone, condition/incapacitated, save/repeat, condition/paralyzed, save/repeat-until-success, save/automatic-success",
	},
	"Skeleton": {
		"action|shortsword|Shortsword; melee +5 reach 5 range 0/0 6 (1d6 + 3) piercing; effect/piercing-damage 6 (1d6 + 3)",
		"action|shortbow|Shortbow; ranged +5 reach 0 range 80/320 6 (1d6 + 3) piercing; effect/piercing-damage 6 (1d6 + 3)",
	},
	"Solar": {
		"trait|trait/divine-awareness|Divine Awareness; effect/know-lies",
		"trait|trait/exalted-restoration|Exalted Restoration; effect/revive",
		"trait|trait/legendary-resistance|Legendary Resistance (4/Day); limit/per-day, save/choose-success",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-or-ranged-attack, multiattack/replace-attack, action/saving-throw",
		"action|flying-sword|Flying Sword; effect/slashing-damage 22 (4d6 + 8), effect/radiant-damage 36 (8d8), effect/return-after-ranged-attack",
		"action|slaying-bow|Slaying Bow; dex DC 21; creature 0/0/600; save/dex, target/creature, effect/slay, effect/piercing-damage 24 (4d8 + 6), effect/radiant-damage 36 (8d8)",
		"action|spellcasting|Spellcasting; cha DC 25 +0; at will: Detect Evil and Good; 1/day: Commune, Control Weather, Dispel Evil and Good, Resurrection; spellcasting/no-material-components, spell/detect-evil-and-good [spellcasting/at-will], spell/commune [limit/per-day], spell/control-weather [limit/per-day], spell/dispel-evil-and-good [limit/per-day], spell/resurrection [limit/per-day]",
		"bonus|divine-aid|Divine Aid (3/Day); limit/per-day, action/spellcasting, spell/cure-wounds, spell/lesser-restoration, spell/remove-curse",
		"legendary|blinding-gaze|Blinding Gaze; con DC 25; creature 0/0/120; save/con, target/creature, condition/blinded, limit/not-again-until-next-turn",
		"legendary|radiant-teleport|Radiant Teleport; dex DC 25; emanation 10/0/0; save/dex, area/emanation, movement/teleport, effect/radiant-damage 11 (2d10), save/half-damage",
		"legendary-uses|3|0",
	},
	"Specter": {
		"trait|trait/incorporeal-movement|Incorporeal Movement; 5 (1d10); movement/through-creatures-and-objects, effect/end-turn-in-object, effect/force-damage",
		"trait|trait/sunlight-sensitivity|Sunlight Sensitivity; effect/sunlight-disadvantage",
		"action|life-drain|Life Drain; melee +4 reach 5 range 0/0 7 (2d6) necrotic; effect/necrotic-damage 7 (2d6), requirement/target-creature, effect/hit-point-maximum-reduction",
	},
	"Sphinx of Lore": {
		"trait|trait/inscrutable|Inscrutable; effect/scrying-immunity, check/disadvantage-to-read-intentions",
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day, or 4/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|claw|Claw; melee +8 reach 5 range 0/0 14 (3d6 + 4) slashing; effect/slashing-damage 14 (3d6 + 4)",
		"action|mind-rending-roar|Mind-Rending Roar (Recharge 5–6); wis DC 16; emanation 300/0/0; limit/recharge, save/wis, area/emanation, target/enemies-only, effect/psychic-damage 35 (10d6), condition/incapacitated",
		"action|spellcasting|Spellcasting; int DC 16 +0; at will: Detect Magic, Identify, Mage Hand, Minor Illusion, Prestidigitation; 1/day: Dispel Magic, Legend Lore, Locate Object, Plane Shift, Remove Curse, Tongues; spellcasting/no-material-components, spell/detect-magic [spellcasting/at-will], spell/identify [spellcasting/at-will], spell/mage-hand [spellcasting/at-will], spell/minor-illusion [spellcasting/at-will], spell/prestidigitation [spellcasting/at-will], spell/dispel-magic [limit/per-day], spell/legend-lore [limit/per-day], spell/locate-object [limit/per-day], spell/plane-shift [limit/per-day], spell/remove-curse [limit/per-day], spell/tongues [limit/per-day]",
		"legendary|arcane-prowl|Arcane Prowl; movement/teleport, action/melee-attack",
		"legendary|weight-of-years|Weight of Years; con DC 16; creature 0/0/120; save/con, target/creature, effect/exhaustion-level, effect/appear-older, limit/not-again-until-next-turn",
		"legendary-uses|3|4",
	},
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters172SourceValues(t *testing.T) {
	c := loadMonsters172(t)
	if len(c.Monsters) != 12 || len(monsters172Headers) != 12 || len(monsters172Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #172") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters172Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader170(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries164(m); !slices.Equal(got, monsters172Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters172Entries[want.name])
			}
			if m.VulnerabilityNote != "" || m.ResistanceNote != "" || m.ConditionImmunityNote != "" || m.TelepathyNote != "" || m.Senses.Notes != nil || m.SwarmOf != "" || m.SpeedNote != "" {
				t.Fatal("source lists no vulnerability note, resistance note, condition immunity note, telepathy note, sense notes, swarm line, speed note")
			}
		})
		for _, line := range monsters172Entries[want.name] {
			if !strings.HasPrefix(line, "legendary-uses|") {
				entries++
			}
		}
	}
	// Bold-italic entry headings counted on PDF pp. 286, 322–328 and 361.
	if entries != 54 {
		t.Fatalf("reviewed %d entries, want 54", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters172MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 172 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters172(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Animals appendix blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("monsters monsters animals monsters monsters monsters monsters monsters monsters monsters monsters monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #172 cites the page each block starts on, in this order.
	names := []string{"Scout", "Sea Hag", "Seahorse", "Shadow", "Shambling Mound", "Shield Guardian", "Shrieker Fungus", "Silver Dragon Wyrmling", "Skeleton", "Solar", "Specter", "Sphinx of Lore"}
	pages := []int{322, 322, 361, 322, 323, 323, 286, 324, 325, 326, 327, 327}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic for these twelve blocks; see monsterConsistency170.
// The Shambling Mound's printed Stealth bonus is pinned in printedSkill170.
func TestMonsters172DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"1/2": 2, "2": 2, "0": 2, "5": 3, "7": 3, "1/4": 2, "21": 7, "1": 2, "11": 4}
	attacks := map[string]attackAbility170{
		"Specter/Life Drain": {"dex", false},
	}
	for _, m := range loadMonsters172(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			monsterConsistency170(t, m, pb, attacks)
		})
	}
}

// Every listed key is used and required, every effect quotes its own entry, and
// every dice expression is typed. Removing any one key invalidates the record.
func TestMonsters172DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters172(t).Monsters {
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
func TestMonsters172RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters172(t).Monsters {
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

func TestMonsters172MalformedFixtures(t *testing.T) {
	tests := []struct {
		name, record string
		mutate       func(*Monster)
	}{
		{"missing languages list", "Scout", func(m *Monster) { m.Languages = nil }},
		{"animal or monster block in another family", "Seahorse", func(m *Monster) {
			m.ID = strings.Replace(strings.Replace(m.ID, "/animals/", "/spells/", 1), "/monsters/", "/spells/", 1)
		}},
		{"XP off CR", "Solar", func(m *Monster) { m.XP++ }},
		{"save off modifier", "Solar", func(m *Monster) {
			m.Abilities["wis"] = MonsterAbility{m.Abilities["wis"].Score, m.Abilities["wis"].Modifier, m.Abilities["wis"].Modifier + 11}
		}},
		{"modifier", "Solar", func(m *Monster) {
			m.Abilities["str"] = MonsterAbility{m.Abilities["str"].Score, m.Abilities["str"].Modifier + 1, m.Abilities["str"].Modifier + 1}
		}},
		{"HP average", "Solar", func(m *Monster) { m.HP.Average++ }},
		{"reaction-only block without its reaction", "Shrieker Fungus", func(m *Monster) { m.Reactions = nil }},
		{"null action list", "Shrieker Fungus", func(m *Monster) { m.Actions = nil }},
		{"recharge", "Sea Hag", func(m *Monster) { m.Actions[1].Uses.Recharge = 7 }},
		{"trait save DC off entry text", "Sea Hag", func(m *Monster) { m.Traits[2].Save.DC = 12 }},
		{"emanation width", "Sphinx of Lore", func(m *Monster) { m.Actions[2].Target.Width = 5 }},
		{"unlisted spell", "Solar", func(m *Monster) { m.Actions[3].Spellcasting.AtWill[0] = "Wish" }},
		{"spell DC off entry text", "Solar", func(m *Monster) { m.Actions[3].Spellcasting.SaveDC = 24 }},
		{"daily uses without count", "Solar", func(m *Monster) { m.BonusActions[0].Uses.PerDay = 0 }},
		{"lair XP not above XP", "Sphinx of Lore", func(m *Monster) { m.LairXP = m.XP }},
		{"trait lair uses not above daily uses", "Sphinx of Lore", func(m *Monster) { m.Traits[1].Uses.LairPerDay = 3 }},
		{"legendary lair uses not above uses", "Sphinx of Lore", func(m *Monster) { m.Legendary.LairUses = 3 }},
		{"legendary uses without actions", "Solar", func(m *Monster) { m.LegendaryActions = nil }},
		{"empty bonus actions", "Scout", func(m *Monster) { m.BonusActions = []MonsterAction{} }},
		{"empty reactions", "Scout", func(m *Monster) { m.Reactions = []MonsterAction{} }},
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
		c := loadMonsters172(t)
		c.Monsters[1].ID = c.Monsters[2].ID
		rejected(t, c, true)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters172(t)
			tt.mutate(record(t, &c, tt.record))
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			rejected(t, c, !strings.HasPrefix(tt.name, "empty "))
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_172.json")
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
		{`"dc":11`, `"dc":11,"dc":22`},
		{`"width":0,`, ``},
		{`"save_dc":25,`, ``},
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
// to 21 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters172ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters172(t)
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
	for _, load := range []func(*testing.T) MonsterCatalogue{loadMonsters158, loadMonsters159, loadMonsters160, loadMonsters161, loadMonsters162, loadMonsters163, loadMonsters164, loadMonsters165, loadMonsters166, loadMonsters167, loadMonsters168, loadMonsters169, loadMonsters170, loadMonsters171} {
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
	for _, id := range []string{"charged-tendril", "claw", "draining-swipe", "fist", "life-drain", "longbow", "rend", "shortbow", "shortsword"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
