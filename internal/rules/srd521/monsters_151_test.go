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

func loadMonsters151(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_151.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #151 order, with each
// record's inclusive page range. Wrapped lines are rejoined and the Adult White
// Dragon's printed "Con22" is spaced; nothing else is normalised.
var monsters151Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Aboleth", [2]int{258, 258}, "Large Aberration, Lawful Evil|AC 17 Initiative +7 (17)|HP 150 (20d10 + 40)|Speed 10 ft., Swim 40 ft.|Str 21 +5 +5 Dex 9 −1 +3 Con 15 +2 +6|Int 18 +4 +8 Wis 15 +2 +6 Cha 18 +4 +4|Skills History +12, Perception +10|Senses Darkvision 120 ft.; Passive Perception 20|Languages Deep Speech; telepathy 120 ft.|CR 10 (XP 5,900, or 7,200 in lair; PB +4)"},
	{"Adult Black Dragon", [2]int{264, 264}, "Huge Dragon (Chromatic), Chaotic Evil|AC 19 Initiative +12 (22)|HP 195 (17d12 + 85)|Speed 40 ft., Fly 80 ft., Swim 40 ft.|Str 23 +6 +6 Dex 14 +2 +7 Con 21 +5 +5|Int 14 +2 +2 Wis 13 +1 +6 Cha 19 +4 +4|Skills Perception +11, Stealth +7|Immunities Acid|Senses Blindsight 60 ft., Darkvision 120 ft.; Passive Perception 21|Languages Common, Draconic|CR 14 (XP 11,500, or 13,000 in lair; PB +5)"},
	{"Adult Blue Dragon", [2]int{266, 267}, "Huge Dragon (Chromatic), Lawful Evil|AC 19 Initiative +10 (20)|HP 212 (17d12 + 102)|Speed 40 ft., Burrow 30 ft., Fly 80 ft.|Str 25 +7 +7 Dex 10 +0 +5 Con 23 +6 +6|Int 16 +3 +3 Wis 15 +2 +7 Cha 20 +5 +5|Skills Perception +12, Stealth +5|Immunities Lightning|Senses Blindsight 60 ft., Darkvision 120 ft.; Passive Perception 22|Languages Common, Draconic|CR 16 (XP 15,000, or 18,000 in lair; PB +5)"},
	{"Adult Brass Dragon", [2]int{268, 269}, "Huge Dragon (Metallic), Chaotic Good|AC 18 Initiative +10 (20)|HP 172 (15d12 + 75)|Speed 40 ft., Burrow 30 ft., Fly 80 ft.|Str 23 +6 +6 Dex 10 +0 +5 Con 21 +5 +5|Int 14 +2 +2 Wis 13 +1 +6 Cha 17 +3 +3|Skills History +7, Perception +11, Persuasion +8, Stealth +5|Immunities Fire|Senses Blindsight 60 ft., Darkvision 120 ft.; Passive Perception 21|Languages Common, Draconic|CR 13 (XP 10,000, or 11,500 in lair; PB +5)"},
	{"Adult Bronze Dragon", [2]int{270, 271}, "Huge Dragon (Metallic), Lawful Good|AC 18 Initiative +10 (20)|HP 212 (17d12 + 102)|Speed 40 ft., Fly 80 ft., Swim 40 ft.|Str 25 +7 +7 Dex 10 +0 +5 Con 23 +6 +6|Int 16 +3 +3 Wis 15 +2 +7 Cha 20 +5 +5|Skills Insight +7, Perception +12, Stealth +5|Immunities Lightning|Senses Blindsight 60 ft., Darkvision 120 ft.; Passive Perception 22|Languages Common, Draconic|CR 15 (XP 13,000, or 15,000 in lair; PB +5)"},
	{"Adult Copper Dragon", [2]int{276, 277}, "Huge Dragon (Metallic), Chaotic Good|AC 18 Initiative +11 (21)|HP 184 (16d12 + 80)|Speed 40 ft., Climb 40 ft., Fly 80 ft.|Str 23 +6 +6 Dex 12 +1 +6 Con 21 +5 +5|Int 18 +4 +4 Wis 15 +2 +7 Cha 18 +4 +4|Skills Deception +9, Perception +12, Stealth +6|Immunities Acid|Senses Blindsight 60 ft., Darkvision 120 ft.; Passive Perception 22|Languages Common, Draconic|CR 14 (XP 11,500, or 13,000 in lair; PB +5)"},
	{"Adult Gold Dragon", [2]int{291, 292}, "Huge Dragon (Metallic), Lawful Good|AC 19 Initiative +14 (24)|HP 243 (18d12 + 126)|Speed 40 ft., Fly 80 ft., Swim 40 ft.|Str 27 +8 +8 Dex 14 +2 +8 Con 25 +7 +7|Int 16 +3 +3 Wis 15 +2 +8 Cha 24 +7 +7|Skills Insight +8, Perception +14, Persuasion +13, Stealth +8|Immunities Fire|Senses Blindsight 60 ft., Darkvision 120 ft.; Passive Perception 24|Languages Common, Draconic|CR 17 (XP 18,000, or 20,000 in lair; PB +6)"},
	{"Adult Green Dragon", [2]int{294, 294}, "Huge Dragon (Chromatic), Lawful Evil|AC 19 Initiative +11 (21)|HP 207 (18d12 + 90)|Speed 40 ft., Fly 80 ft., Swim 40 ft.|Str 23 +6 +6 Dex 12 +1 +6 Con 21 +5 +5|Int 18 +4 +4 Wis 15 +2 +7 Cha 18 +4 +4|Skills Deception +9, Perception +12, Persuasion +9, Stealth +6|Immunities Poison; Poisoned|Senses Blindsight 60 ft., Darkvision 120 ft.; Passive Perception 22|Languages Common, Draconic|CR 15 (XP 13,000, or 15,000 in lair; PB +5)"},
	{"Adult Red Dragon", [2]int{318, 319}, "Huge Dragon (Chromatic), Chaotic Evil|AC 19 Initiative +12 (22)|HP 256 (19d12 + 133)|Speed 40 ft., Climb 40 ft., Fly 80 ft.|Str 27 +8 +8 Dex 10 +0 +6 Con 25 +7 +7|Int 16 +3 +3 Wis 13 +1 +7 Cha 23 +6 +6|Skills Perception +13, Stealth +6|Immunities Fire|Senses Blindsight 60 ft., Darkvision 120 ft.; Passive Perception 23|Languages Common, Draconic|CR 17 (XP 18,000, or 20,000 in lair; PB +6)"},
	{"Adult Silver Dragon", [2]int{324, 325}, "Huge Dragon (Metallic), Lawful Good|AC 19 Initiative +10 (20)|HP 216 (16d12 + 112)|Speed 40 ft., Fly 80 ft.|Str 27 +8 +8 Dex 10 +0 +5 Con 25 +7 +7|Int 16 +3 +3 Wis 13 +1 +6 Cha 22 +6 +6|Skills History +8, Perception +11, Stealth +5|Immunities Cold|Senses Blindsight 60 ft., Darkvision 120 ft.; Passive Perception 21|Languages Common, Draconic|CR 16 (XP 15,000, or 18,000 in lair; PB +5)"},
	{"Adult White Dragon", [2]int{340, 340}, "Huge Dragon (Chromatic), Chaotic Evil|AC 18 Initiative +10 (20)|HP 200 (16d12 + 96)|Speed 40 ft., Burrow 30 ft., Fly 80 ft., Swim 40 ft.|Str 22 +6 +6 Dex 10 +0 +5 Con 22 +6 +6|Int 8 −1 −1 Wis 12 +1 +6 Cha 12 +1 +1|Skills Perception +11, Stealth +5|Immunities Cold|Senses Blindsight 60 ft., Darkvision 120 ft.; Passive Perception 21|Languages Common, Draconic|CR 13 (XP 10,000, or 11,500 in lair; PB +5)"},
	{"Air Elemental", [2]int{258, 259}, "Large Elemental, Neutral|AC 15 Initiative +5 (15)|HP 90 (12d10 + 24)|Speed 10 ft., Fly 90 ft. (hover)|Str 14 +2 +2 Dex 20 +5 +5 Con 14 +2 +2|Int 6 −2 −2 Wis 10 +0 +0 Cha 6 −2 −2|Resistances Bludgeoning, Lightning, Piercing, Slashing|Immunities Poison, Thunder; Exhaustion, Grappled, Paralyzed, Petrified, Poisoned, Prone, Restrained, Unconscious|Senses Darkvision 60 ft.; Passive Perception 10|Languages Primordial (Auran)|CR 5 (XP 1,800; PB +3)"},
}

// Typed values of every trait, action and legendary action, reviewed against
// the same pages: limited uses, attack bonus/reach/damage, save ability and DC,
// target shape size/width/range, dice, spell lists and dependency keys.
var monsters151Entries = map[string][]string{
	"Aboleth": {
		"trait|trait/amphibious|Amphibious; ",
		"trait|trait/eldritch-restoration|Eldritch Restoration; effect/revive",
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day, or 4/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"trait|trait/mucus-cloud|Mucus Cloud; con DC 14; emanation 5/0/0; 6 (1d12); save/con, area/emanation, effect/acid-damage, effect/curse, effect/no-hit-point-recovery",
		"trait|trait/probing-telepathy|Probing Telepathy; effect/learn-desires",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/additional-action",
		"action|tentacle|Tentacle; melee +9 reach 15 range 0/0 12 (2d6 + 5) bludgeoning; effect/bludgeoning-damage 12 (2d6 + 5), condition/grappled, requirement/target-size",
		"action|consume-memories|Consume Memories; int DC 16; creature 0/0/30; save/int, target/creature, requirement/charmed-or-grappled-creature, effect/psychic-damage 10 (3d6), save/half-damage, effect/gain-memories",
		"action|dominate-mind|Dominate Mind (2/Day); wis DC 16; creature 0/0/30; save/wis, target/creature, limit/per-day, condition/charmed, effect/control-charmed-creature, effect/telepathic-link, save/repeat",
		"legendary|lash|Lash; action/melee-attack",
		"legendary|psychic-drain|Psychic Drain; effect/regain-hit-points 5 (1d10), action/consume-memories, requirement/charmed-or-grappled-creature",
		"legendary-uses|3|4",
	},
	"Adult Black Dragon": {
		"trait|trait/amphibious|Amphibious; ",
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day, or 4/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, spell/acid-arrow",
		"action|rend|Rend; melee +11 reach 10 range 0/0 13 (2d6 + 6) slashing; effect/slashing-damage 13 (2d6 + 6), effect/acid-damage 4 (1d8)",
		"action|acid-breath|Acid Breath (Recharge 5–6); dex DC 18; line 60/5/0; save/dex, area/line, limit/recharge, effect/acid-damage 54 (12d8), save/half-damage",
		"action|spellcasting|Spellcasting; cha DC 17 +9; at will: Acid Arrow (level 3 version), Detect Magic, Fear; 1/day: Speak with Dead, Vitriolic Sphere; spellcasting/no-material-components, spell/acid-arrow, spell/detect-magic, spell/fear, spell/speak-with-dead, spell/vitriolic-sphere",
		"legendary|cloud-of-insects|Cloud of Insects; dex DC 17; creature 0/0/120; save/dex, target/creature, effect/poison-damage 22 (4d10), effect/concentration-save-disadvantage, limit/once-per-round",
		"legendary|frightful-presence|Frightful Presence; action/spellcasting, spell/fear, limit/once-per-round",
		"legendary|pounce|Pounce; movement/half-speed, action/melee-attack",
		"legendary-uses|3|4",
	},
	"Adult Blue Dragon": {
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day, or 4/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, spell/shatter",
		"action|rend|Rend; melee +12 reach 10 range 0/0 16 (2d8 + 7) slashing; effect/slashing-damage 16 (2d8 + 7), effect/lightning-damage 5 (1d10)",
		"action|lightning-breath|Lightning Breath (Recharge 5–6); dex DC 19; line 90/5/0; save/dex, area/line, limit/recharge, effect/lightning-damage 60 (11d10), save/half-damage",
		"action|spellcasting|Spellcasting; cha DC 18 +0; at will: Detect Magic, Invisibility, Mage Hand, Shatter; 1/day: Scrying, Sending; spellcasting/no-material-components, spell/detect-magic, spell/invisibility, spell/mage-hand, spell/shatter, spell/scrying, spell/sending",
		"legendary|cloaked-flight|Cloaked Flight; movement/half-fly-speed, action/spellcasting, spell/invisibility, limit/once-per-round",
		"legendary|sonic-boom|Sonic Boom; action/spellcasting, spell/shatter, limit/once-per-round",
		"legendary|tail-swipe|Tail Swipe; action/melee-attack",
		"legendary-uses|3|4",
	},
	"Adult Brass Dragon": {
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day, or 4/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, spell/scorching-ray",
		"action|rend|Rend; melee +11 reach 10 range 0/0 17 (2d10 + 6) slashing; effect/slashing-damage 17 (2d10 + 6), effect/fire-damage 4 (1d8)",
		"action|fire-breath|Fire Breath (Recharge 5–6); dex DC 18; line 60/5/0; save/dex, area/line, limit/recharge, effect/fire-damage 45 (10d8), save/half-damage",
		"action|sleep-breath|Sleep Breath; con DC 18; cone 60/0/0; save/con, area/cone, condition/incapacitated, save/repeat, condition/unconscious, save/second-failure, effect/ends-on-damage-or-wake",
		"action|spellcasting|Spellcasting; cha DC 16 +0; at will: Detect Magic, Minor Illusion, Scorching Ray, Shapechange (Beast or Humanoid form only, no Temporary Hit Points gained from the spell, and no Concentration or Temporary Hit Points required to maintain the spell), Speak with Animals; 1/day: Detect Thoughts, Control Weather; spellcasting/no-material-components, spell/detect-magic, spell/minor-illusion, spell/scorching-ray, spell/shapechange, spell/speak-with-animals, spell/detect-thoughts, spell/control-weather",
		"legendary|blazing-light|Blazing Light; action/spellcasting, spell/scorching-ray",
		"legendary|pounce|Pounce; movement/half-speed, action/melee-attack",
		"legendary|scorching-sands|Scorching Sands; dex DC 16; creature 0/0/120; save/dex, target/creature, effect/fire-damage 27 (6d8), effect/speed-halved, limit/once-per-round",
		"legendary-uses|3|4",
	},
	"Adult Bronze Dragon": {
		"trait|trait/amphibious|Amphibious; ",
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day, or 4/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, spell/guiding-bolt",
		"action|rend|Rend; melee +12 reach 10 range 0/0 16 (2d8 + 7) slashing; effect/slashing-damage 16 (2d8 + 7), effect/lightning-damage 5 (1d10)",
		"action|lightning-breath|Lightning Breath (Recharge 5–6); dex DC 19; line 90/5/0; save/dex, area/line, limit/recharge, effect/lightning-damage 55 (10d10), save/half-damage",
		"action|repulsion-breath|Repulsion Breath; str DC 19; cone 30/0/0; save/str, area/cone, condition/prone, effect/push",
		"action|spellcasting|Spellcasting; cha DC 17 +10; at will: Detect Magic, Guiding Bolt (level 2 version), Shapechange (Beast or Humanoid form only, no Temporary Hit Points gained from the spell, and no Concentration or Temporary Hit Points required to maintain the spell), Speak with Animals, Thaumaturgy; 1/day: Detect Thoughts, Water Breathing; spellcasting/no-material-components, spell/detect-magic, spell/guiding-bolt, spell/shapechange, spell/speak-with-animals, spell/thaumaturgy, spell/detect-thoughts, spell/water-breathing",
		"legendary|guiding-light|Guiding Light; action/spellcasting, spell/guiding-bolt",
		"legendary|pounce|Pounce; movement/half-speed, action/melee-attack",
		"legendary|thunderclap|Thunderclap; con DC 17; sphere 20/0/90; save/con, area/sphere, effect/thunder-damage 10 (3d6), condition/deafened",
		"legendary-uses|3|4",
	},
	"Adult Copper Dragon": {
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day, or 4/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, spell/mind-spike",
		"action|rend|Rend; melee +11 reach 10 range 0/0 17 (2d10 + 6) slashing; effect/slashing-damage 17 (2d10 + 6), effect/acid-damage 4 (1d8)",
		"action|acid-breath|Acid Breath (Recharge 5–6); dex DC 18; line 60/5/0; save/dex, area/line, limit/recharge, effect/acid-damage 54 (12d8), save/half-damage",
		"action|slowing-breath|Slowing Breath; con DC 18; cone 60/0/0; save/con, area/cone, effect/speed-halved, effect/no-reactions, effect/action-or-bonus-action",
		"action|spellcasting|Spellcasting; cha DC 17 +0; at will: Detect Magic, Mind Spike (level 4 version), Minor Illusion, Shapechange (Beast or Humanoid form only, no Temporary Hit Points gained from the spell, and no Concentration or Temporary Hit Points required to maintain the spell); 1/day: Greater Restoration, Major Image; spellcasting/no-material-components, spell/detect-magic, spell/mind-spike, spell/minor-illusion, spell/shapechange, spell/greater-restoration, spell/major-image",
		"legendary|giggling-magic|Giggling Magic; cha DC 17; creature 0/0/90; save/cha, target/creature, effect/psychic-damage 24 (7d6), effect/d20-test-penalty, limit/once-per-round",
		"legendary|mind-jolt|Mind Jolt; action/spellcasting, spell/mind-spike, limit/once-per-round",
		"legendary|pounce|Pounce; movement/half-speed, action/melee-attack",
		"legendary-uses|3|4",
	},
	"Adult Gold Dragon": {
		"trait|trait/amphibious|Amphibious; ",
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day, or 4/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, spell/guiding-bolt",
		"action|rend|Rend; melee +14 reach 10 range 0/0 17 (2d8 + 8) slashing; effect/slashing-damage 17 (2d8 + 8), effect/fire-damage 4 (1d8)",
		"action|fire-breath|Fire Breath (Recharge 5–6); dex DC 21; cone 60/0/0; save/dex, area/cone, limit/recharge, effect/fire-damage 66 (12d10), save/half-damage",
		"action|spellcasting|Spellcasting; cha DC 21 +13; at will: Detect Magic, Guiding Bolt (level 2 version), Shapechange (Beast or Humanoid form only, no Temporary Hit Points gained from the spell, and no Concentration or Temporary Hit Points required to maintain the spell); 1/day: Flame Strike, Zone of Truth; spellcasting/no-material-components, spell/detect-magic, spell/guiding-bolt, spell/shapechange, spell/flame-strike, spell/zone-of-truth",
		"action|weakening-breath|Weakening Breath; str DC 21; cone 60/0/0; save/str, area/cone, requirement/unaffected-target, effect/damage-roll-penalty 3 (1d6), effect/strength-test-disadvantage, save/repeat, save/automatic-success",
		"legendary|banish|Banish; cha DC 21; creature 0/0/120; save/cha, target/creature, effect/force-damage 10 (3d6), condition/incapacitated, effect/banishment, limit/once-per-round",
		"legendary|guiding-light|Guiding Light; action/spellcasting, spell/guiding-bolt",
		"legendary|pounce|Pounce; movement/half-speed, action/melee-attack",
		"legendary-uses|3|4",
	},
	"Adult Green Dragon": {
		"trait|trait/amphibious|Amphibious; ",
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day, or 4/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, spell/mind-spike",
		"action|rend|Rend; melee +11 reach 10 range 0/0 15 (2d8 + 6) slashing; effect/slashing-damage 15 (2d8 + 6), effect/poison-damage 7 (2d6)",
		"action|poison-breath|Poison Breath (Recharge 5–6); con DC 18; cone 60/0/0; save/con, area/cone, limit/recharge, effect/poison-damage 56 (16d6), save/half-damage",
		"action|spellcasting|Spellcasting; cha DC 17 +0; at will: Detect Magic, Mind Spike (level 3 version); 1/day: Geas; spellcasting/no-material-components, spell/detect-magic, spell/mind-spike, spell/geas",
		"legendary|mind-invasion|Mind Invasion; action/spellcasting, spell/mind-spike",
		"legendary|noxious-miasma|Noxious Miasma; con DC 17; sphere 20/0/90; save/con, area/sphere, effect/poison-damage 7 (2d6), effect/ac-penalty, limit/once-per-round",
		"legendary|pounce|Pounce; movement/half-speed, action/melee-attack",
		"legendary-uses|3|4",
	},
	"Adult Red Dragon": {
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day, or 4/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, spell/scorching-ray",
		"action|rend|Rend; melee +14 reach 10 range 0/0 13 (1d10 + 8) slashing; effect/slashing-damage 13 (1d10 + 8), effect/fire-damage 5 (2d4)",
		"action|fire-breath|Fire Breath (Recharge 5–6); dex DC 21; cone 60/0/0; save/dex, area/cone, limit/recharge, effect/fire-damage 59 (17d6), save/half-damage",
		"action|spellcasting|Spellcasting; cha DC 20 +12; at will: Command (level 2 version), Detect Magic, Scorching Ray; 1/day: Fireball; spellcasting/no-material-components, spell/command, spell/detect-magic, spell/scorching-ray, spell/fireball",
		"legendary|commanding-presence|Commanding Presence; action/spellcasting, spell/command, limit/once-per-round",
		"legendary|fiery-rays|Fiery Rays; action/spellcasting, spell/scorching-ray, limit/once-per-round",
		"legendary|pounce|Pounce; movement/half-speed, action/melee-attack",
		"legendary-uses|3|4",
	},
	"Adult Silver Dragon": {
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day, or 4/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, spell/ice-knife",
		"action|rend|Rend; melee +13 reach 10 range 0/0 17 (2d8 + 8) slashing; effect/slashing-damage 17 (2d8 + 8), effect/cold-damage 4 (1d8)",
		"action|cold-breath|Cold Breath (Recharge 5–6); con DC 20; cone 60/0/0; save/con, area/cone, limit/recharge, effect/cold-damage 54 (12d8), save/half-damage",
		"action|paralyzing-breath|Paralyzing Breath; con DC 20; cone 60/0/0; save/con, area/cone, condition/incapacitated, save/repeat, condition/paralyzed, save/repeat, save/second-failure, save/automatic-success",
		"action|spellcasting|Spellcasting; cha DC 19 +11; at will: Detect Magic, Hold Monster, Ice Knife, Shapechange (Beast or Humanoid form only, no Temporary Hit Points gained from the spell, and no Concentration or Temporary Hit Points required to maintain the spell); 1/day: Ice Storm (level 5 version), Zone of Truth; spellcasting/no-material-components, spell/detect-magic, spell/hold-monster, spell/ice-knife, spell/shapechange, spell/ice-storm, spell/zone-of-truth",
		"legendary|chill|Chill; action/spellcasting, spell/hold-monster, limit/once-per-round",
		"legendary|cold-gale|Cold Gale; dex DC 19; line 60/10/0; save/dex, area/line, effect/cold-damage 14 (4d6), effect/push, save/half-damage-only, limit/once-per-round",
		"legendary|pounce|Pounce; movement/half-speed, action/melee-attack",
		"legendary-uses|3|4",
	},
	"Adult White Dragon": {
		"trait|trait/ice-walk|Ice Walk; movement/ice-walk, movement/ignore-difficult-terrain",
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day, or 4/Day in Lair); limit/per-day, lair/uses, save/choose-success",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +11 reach 10 range 0/0 13 (2d6 + 6) slashing; effect/slashing-damage 13 (2d6 + 6), effect/cold-damage 4 (1d8)",
		"action|cold-breath|Cold Breath (Recharge 5–6); con DC 19; cone 60/0/0; save/con, area/cone, limit/recharge, effect/cold-damage 54 (12d8), save/half-damage",
		"legendary|freezing-burst|Freezing Burst; con DC 14; sphere 30/0/120; save/con, area/sphere, effect/cold-damage 7 (2d6), effect/speed-zero, limit/once-per-round",
		"legendary|frightful-presence|Frightful Presence; spellcasting/fixed-save-dc, spellcasting/no-material-components, spell/fear, limit/once-per-round",
		"legendary|pounce|Pounce; movement/half-speed, action/melee-attack",
		"legendary-uses|3|4",
	},
	"Air Elemental": {
		"trait|trait/air-form|Air Form; movement/enter-creature-space, movement/squeeze",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|thunderous-slam|Thunderous Slam; melee +8 reach 10 range 0/0 14 (2d8 + 5) thunder; effect/thunder-damage 14 (2d8 + 5)",
		"action|whirlwind|Whirlwind (Recharge 4–6); str DC 13; creature 0/0/0; save/str, target/creature, limit/recharge, requirement/target-size, effect/thunder-damage 24 (4d10 + 2), condition/prone, effect/push, save/half-damage-only",
	},
}

func monsterSigned151(n int) string {
	if n < 0 {
		return fmt.Sprintf("−%d", -n)
	}
	return fmt.Sprintf("+%d", n)
}

func monsterThousands151(n int) string {
	if n < 1000 {
		return fmt.Sprint(n)
	}
	return fmt.Sprintf("%d,%03d", n/1000, n%1000)
}

func monsterDice151(d MonsterHP) string {
	if d.Modifier == 0 {
		return fmt.Sprintf("%d (%dd%d)", d.Average, d.Dice, d.Sides)
	}
	return fmt.Sprintf("%d (%dd%d + %d)", d.Average, d.Dice, d.Sides, d.Modifier)
}

func monsterTitle151(s string) string {
	words := strings.Split(strings.ReplaceAll(s, "-", " "), " ")
	for i, w := range words {
		words[i] = strings.ToUpper(w[:1]) + w[1:]
	}
	return strings.Join(words, " ")
}

func monsterTitles151(values []string) string {
	out := make([]string, len(values))
	for i, v := range values {
		out[i] = monsterTitle151(v)
	}
	return strings.Join(out, ", ")
}

func monsterSorted151[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// monsterHeader151 prints typed fields in the source's stat-block notation.
func monsterHeader151(m Monster) string {
	kind := monsterTitle151(m.Type)
	if len(m.Tags) > 0 {
		kind += " (" + strings.Join(m.Tags, ", ") + ")"
	}
	lines := []string{
		monsterTitles151(m.Sizes) + " " + kind + ", " + m.Alignment,
		fmt.Sprintf("AC %d Initiative %s (%d)", m.AC, monsterSigned151(m.Initiative.Modifier), m.Initiative.Score),
		"HP " + monsterDice151(m.HP),
	}
	speed := fmt.Sprintf("Speed %d ft.", m.Speeds["walk"])
	for _, mode := range monsterSorted151(m.Speeds) {
		if mode != "walk" {
			speed += fmt.Sprintf(", %s %d ft.", monsterTitle151(mode), m.Speeds[mode])
		}
	}
	if m.Hover {
		speed += " (hover)"
	}
	lines = append(lines, speed)
	for _, row := range [][]string{{"str", "dex", "con"}, {"int", "wis", "cha"}} {
		cells := []string{}
		for _, key := range row {
			a := m.Abilities[key]
			cells = append(cells, fmt.Sprintf("%s %d %s %s", monsterTitle151(key), a.Score, monsterSigned151(a.Modifier), monsterSigned151(a.Save)))
		}
		lines = append(lines, strings.Join(cells, " "))
	}
	if m.Skills != nil {
		skills := []string{}
		for _, skill := range monsterSorted151(m.Skills) {
			skills = append(skills, monsterTitle151(skill)+" "+monsterSigned151(m.Skills[skill]))
		}
		lines = append(lines, "Skills "+strings.Join(skills, ", "))
	}
	if m.Vulnerabilities != nil {
		lines = append(lines, "Vulnerabilities "+monsterTitles151(m.Vulnerabilities))
	}
	if m.Resistances != nil {
		lines = append(lines, "Resistances "+monsterTitles151(m.Resistances))
	}
	if immunities := slices.DeleteFunc([]string{monsterTitles151(m.Immunities), monsterTitles151(m.ConditionImmunities)}, func(s string) bool { return s == "" }); len(immunities) > 0 {
		lines = append(lines, "Immunities "+strings.Join(immunities, "; "))
	}
	senses := []string{}
	for _, sense := range monsterSorted151(m.Senses.Special) {
		senses = append(senses, fmt.Sprintf("%s %d ft.", monsterTitle151(sense), m.Senses.Special[sense]))
	}
	passive := fmt.Sprintf("Passive Perception %d", m.Senses.PassivePerception)
	if len(senses) > 0 {
		passive = strings.Join(senses, ", ") + "; " + passive
	}
	languages := "Languages " + strings.Join(m.Languages, ", ")
	if m.Telepathy > 0 {
		languages += fmt.Sprintf("; telepathy %d ft.", m.Telepathy)
	}
	xp := "XP " + monsterThousands151(m.XP)
	if m.LairXP > 0 {
		xp += ", or " + monsterThousands151(m.LairXP) + " in lair"
	}
	lines = append(lines, "Senses "+passive, languages, fmt.Sprintf("CR %s (%s; PB +%d)", m.CR, xp, m.ProficiencyBonus))
	return strings.Join(lines, "|")
}

func monsterEntry151(section, id string, x MonsterMechanic, attack *MonsterAttack, sc *MonsterSpellcasting, effects []MonsterMechanic) string {
	head := x.Name
	if u := x.Uses; u != nil {
		switch {
		case u.Recharge > 0:
			head += fmt.Sprintf(" (Recharge %d–6)", u.Recharge)
		case u.LairPerDay > 0:
			head += fmt.Sprintf(" (%d/Day, or %d/Day in Lair)", u.PerDay, u.LairPerDay)
		default:
			head += fmt.Sprintf(" (%d/Day)", u.PerDay)
		}
	}
	parts := []string{head}
	if attack != nil {
		parts = append(parts, fmt.Sprintf("%s +%d reach %d range %d/%d %s %s", attack.Kind, attack.Bonus, attack.Reach, attack.Range, attack.LongRange, monsterDice151(attack.Damage), attack.DamageType))
	}
	if x.Save != nil {
		parts = append(parts, fmt.Sprintf("%s DC %d", x.Save.Ability, x.Save.DC))
	}
	if t := x.Target; t != nil {
		parts = append(parts, fmt.Sprintf("%s %d/%d/%d", t.Shape, t.Size, t.Width, t.Range))
	}
	if x.Dice != nil {
		parts = append(parts, monsterDice151(*x.Dice))
	}
	if sc != nil {
		parts = append(parts, fmt.Sprintf("%s DC %d +%d", sc.Ability, sc.SaveDC, sc.AttackBonus), "at will: "+strings.Join(sc.AtWill, ", "))
		for _, group := range sc.PerDay {
			parts = append(parts, fmt.Sprintf("%d/day: %s", group.Uses, strings.Join(group.Spells, ", ")))
		}
	}
	keys := slices.Clone(x.Requires)
	for _, e := range effects {
		key := e.Key
		if e.Dice != nil {
			key += " " + monsterDice151(*e.Dice)
		}
		keys = append(keys, key)
	}
	return section + "|" + id + "|" + strings.Join(append(parts, strings.Join(keys, ", ")), "; ")
}

func monsterEntries151(m Monster) []string {
	out := []string{}
	for _, trait := range m.Traits {
		out = append(out, monsterEntry151("trait", trait.Key, trait, nil, nil, nil))
	}
	for _, section := range []struct {
		label   string
		actions []MonsterAction
	}{{"action", m.Actions}, {"bonus", m.BonusActions}, {"reaction", m.Reactions}, {"legendary", m.LegendaryActions}} {
		for _, a := range section.actions {
			out = append(out, monsterEntry151(section.label, a.ID, MonsterMechanic{Name: a.Name, Uses: a.Uses, Save: a.Save, Target: a.Target}, a.Attack, a.Spellcasting, a.Effects))
		}
	}
	if m.Legendary != nil {
		out = append(out, fmt.Sprintf("legendary-uses|%d|%d", m.Legendary.Uses, m.Legendary.LairUses))
	}
	return out
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters151SourceValues(t *testing.T) {
	c := loadMonsters151(t)
	if len(c.Monsters) != 12 || len(monsters151Headers) != 12 || len(monsters151Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #151") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	for i, want := range monsters151Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader151(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries151(m); !slices.Equal(got, monsters151Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters151Entries[want.name])
			}
			if len(m.Gear) != 0 || m.BonusActions != nil || m.Reactions != nil {
				t.Fatal("source lists no gear, bonus actions or reactions")
			}
		})
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters151MatchesDecisionInventory(t *testing.T) {
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
	want, got := []identity{}, []identity{}
	for _, r := range inventory.Records {
		if r.CatalogueIssue == 151 {
			if r.Family != "monsters" || r.Status != "unsupported" {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
		}
	}
	for _, m := range loadMonsters151(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// Issue #151 cites the page each block starts on.
	pages := []int{258, 264, 266, 268, 270, 276, 291, 294, 318, 324, 340, 258}
	for i, r := range want {
		if r.page != pages[i] {
			t.Fatalf("%s: issue cites p.%d", r.name, pages[i])
		}
	}
}

// Source arithmetic that holds for these twelve blocks (PDF pp. 254–256):
// XP and PB follow CR, saves and skills add PB to the ability modifier, and
// attack rolls use the better of Strength and Dexterity.
func TestMonsters151DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"5": 3, "10": 4, "13": 5, "14": 5, "15": 5, "16": 5, "17": 6}
	skillAbility := map[string]string{"deception": "cha", "history": "int", "insight": "wis", "perception": "wis", "persuasion": "cha", "stealth": "dex"}
	hitDie := map[string]int{"large": 10, "huge": 12}
	for _, m := range loadMonsters151(t).Monsters {
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
			best := max(m.Abilities["str"].Modifier, m.Abilities["dex"].Modifier)
			for _, a := range m.allActions() {
				if a.Attack != nil && (a.Attack.Bonus != best+p || a.Attack.Damage.Modifier != best) {
					t.Fatalf("%s attack does not use the ability modifier and PB", a.Name)
				}
			}
		})
	}
}

// Every listed key is used and required, every effect quotes its own entry, and
// every dice expression is typed. Removing any one key invalidates the record.
func TestMonsters151DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters151(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			if m.RequiredMechanics[0] != "monster/stat-block" {
				t.Fatal("stat-block dependency must come first")
			}
			quoted := func(parent string, x MonsterMechanic) {
				if !strings.Contains(parent, strings.TrimSuffix(x.Text, ".")) {
					t.Fatalf("%q is not quoted from its entry", x.Text)
				}
				if x.Dice != nil && !strings.Contains(x.Text, monsterDice151(*x.Dice)) {
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
				if a.Attack != nil && (a.Effects[0].Dice == nil || *a.Effects[0].Dice != a.Attack.Damage) {
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
// keys do not cover any block, and the register still reports the family unsupported.
func TestMonsters151RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters151(t).Monsters {
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
		if e.Family == "monsters" && (e.Status != "unsupported" || len(e.Fixtures) != 0) {
			t.Fatal("catalogue data promoted the monster family")
		}
	}
}

func TestMonsters151MalformedFixtures(t *testing.T) {
	// Indexes: 0 Aboleth, 1 Adult Black Dragon, 11 Air Elemental.
	tests := []struct {
		name   string
		mutate func(*MonsterCatalogue)
	}{
		{"duplicate ID", func(c *MonsterCatalogue) { c.Monsters[1].ID = c.Monsters[0].ID }},
		{"empty tags", func(c *MonsterCatalogue) { c.Monsters[1].Tags = []string{} }},
		{"duplicate immunity", func(c *MonsterCatalogue) { c.Monsters[1].Immunities = []string{"acid", "acid"} }},
		{"unknown resistance", func(c *MonsterCatalogue) { c.Monsters[11].Resistances[0] = "sonic" }},
		{"unknown condition", func(c *MonsterCatalogue) { c.Monsters[11].ConditionImmunities[0] = "dazed" }},
		{"unknown skill", func(c *MonsterCatalogue) { c.Monsters[0].Skills["lore"] = 12 }},
		{"empty skills", func(c *MonsterCatalogue) { c.Monsters[0].Skills = map[string]int{} }},
		{"hover without fly", func(c *MonsterCatalogue) { delete(c.Monsters[11].Speeds, "fly") }},
		{"telepathy", func(c *MonsterCatalogue) { c.Monsters[0].Telepathy = -1 }},
		{"lair XP", func(c *MonsterCatalogue) { c.Monsters[0].LairXP = c.Monsters[0].XP }},
		{"modifier", func(c *MonsterCatalogue) { c.Monsters[0].Abilities["dex"] = MonsterAbility{9, 0, 3} }},
		{"HP average", func(c *MonsterCatalogue) { c.Monsters[0].HP.Average = 151 }},
		{"legendary uses without actions", func(c *MonsterCatalogue) { c.Monsters[1].LegendaryActions = nil }},
		{"legendary actions without uses", func(c *MonsterCatalogue) { c.Monsters[1].Legendary = nil }},
		{"lair uses", func(c *MonsterCatalogue) { c.Monsters[1].Legendary.LairUses = 3 }},
		{"legendary text", func(c *MonsterCatalogue) { c.Monsters[1].Legendary.Text = " " }},
		{"empty reactions", func(c *MonsterCatalogue) { c.Monsters[1].Reactions = []MonsterAction{} }},
		{"duplicate action across sections", func(c *MonsterCatalogue) { c.Monsters[1].LegendaryActions[0].ID = "rend" }},
		{"no uses", func(c *MonsterCatalogue) { *c.Monsters[1].Actions[2].Uses = MonsterUses{} }},
		{"recharge", func(c *MonsterCatalogue) { c.Monsters[1].Actions[2].Uses.Recharge = 7 }},
		{"lair daily uses", func(c *MonsterCatalogue) { c.Monsters[1].Traits[1].Uses.LairPerDay = 3 }},
		{"save ability", func(c *MonsterCatalogue) { c.Monsters[1].Actions[2].Save.Ability = "luck" }},
		{"save DC", func(c *MonsterCatalogue) { c.Monsters[1].Actions[2].Save.DC = 0 }},
		{"target shape", func(c *MonsterCatalogue) { c.Monsters[1].Actions[2].Target.Shape = "cube" }},
		{"line width", func(c *MonsterCatalogue) { c.Monsters[1].Actions[2].Target.Width = 0 }},
		{"creature size", func(c *MonsterCatalogue) { c.Monsters[1].LegendaryActions[0].Target.Size = 5 }},
		{"attack damage", func(c *MonsterCatalogue) { c.Monsters[1].Actions[1].Attack.Damage.Average = 14 }},
		{"attack reach", func(c *MonsterCatalogue) { c.Monsters[1].Actions[1].Attack.Reach = 0 }},
		{"effect dice", func(c *MonsterCatalogue) { c.Monsters[1].Actions[2].Effects[3].Dice.Average = 55 }},
		{"trait dice", func(c *MonsterCatalogue) { c.Monsters[0].Traits[3].Dice.Sides = 13 }},
		{"unlisted requirement", func(c *MonsterCatalogue) { c.Monsters[0].Traits[1].Requires = []string{"effect/future"} }},
		{"duplicate requirement", func(c *MonsterCatalogue) {
			c.Monsters[0].Traits[1].Requires = []string{"effect/revive", "effect/revive"}
		}},
		{"unlisted effect", func(c *MonsterCatalogue) { c.Monsters[1].Actions[2].Effects[4].Key = "save/future" }},
		{"untyped spellcasting", func(c *MonsterCatalogue) { c.Monsters[1].Actions[3].Spellcasting = nil }},
		{"misplaced spellcasting", func(c *MonsterCatalogue) {
			c.Monsters[1].Actions[0].Spellcasting = c.Monsters[1].Actions[3].Spellcasting
		}},
		{"spell ability", func(c *MonsterCatalogue) { c.Monsters[1].Actions[3].Spellcasting.Ability = "" }},
		{"spell DC", func(c *MonsterCatalogue) { c.Monsters[1].Actions[3].Spellcasting.SaveDC = 0 }},
		{"unlisted spell", func(c *MonsterCatalogue) { c.Monsters[1].Actions[3].Spellcasting.AtWill[0] = "Wish" }},
		{"duplicate spell", func(c *MonsterCatalogue) { c.Monsters[1].Actions[3].Spellcasting.PerDay[0].Spells[0] = "Fear" }},
		{"daily spell uses", func(c *MonsterCatalogue) { c.Monsters[1].Actions[3].Spellcasting.PerDay[0].Uses = 0 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters151(t)
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
	raw, err := os.ReadFile("testdata/monsters_151.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, edit := range [][2]string{
		{`"recharge": 5`, `"recharge": null`},
		{`"per_day": 3,`, ``},
		{`"lair_uses": 4,`, ``},
		{`"dc": 18`, `"dc": null`},
		{`"width": 5,`, ``},
		{`"attack_bonus": 9,`, ``},
		{`"save_dc": 17,`, `"save_dc": 17, "save_dc": 18,`},
		{`"uses": 1,`, ``},
		{`"lair_xp": 7200,`, `"lair_xp": 7200, "mythic_xp": 1,`},
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

// The batch file composes with the #143 Bandit catalogue without entering the
// embedded runtime catalogue or changing the Bandit fixture's pinned revision.
func TestMonsters151ComposesWithBandit(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters151(t)
	if len(base.Monsters) != 1 || base.Monsters[0].ID != BanditID {
		t.Fatal("batch data must not enter the embedded runtime catalogue")
	}
	baseSource, batchSource := base.Source, batch.Source
	baseSource.Changes, batchSource.Changes = "", ""
	if baseSource != batchSource {
		t.Fatal("batches must share one source pin")
	}
	combined := MonsterCatalogue{SchemaVersion: 1, Source: batch.Source, Monsters: slices.Concat(base.Monsters, batch.Monsters)}
	if err := combined.Validate(); err != nil || len(combined.Monsters) != 13 {
		t.Fatal("batch must compose with the earlier catalogue:", err)
	}
	combined.Monsters = append(combined.Monsters, batch.Monsters[11])
	if combined.Validate() == nil {
		t.Fatal("composition accepted a duplicate record")
	}
	_, version, err := BanditRecord()
	if err != nil || version != "bandit-v1:1169c77601cf83a9b96985fc761a6b8be9c901dc1a86e6a9758dcc399707c3d4" {
		t.Fatal("contract extension changed the pinned Bandit revision:", version, err)
	}
	if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: "rend"}, 12, nil); err == nil {
		t.Fatal("a catalogued action became executable")
	}
}
