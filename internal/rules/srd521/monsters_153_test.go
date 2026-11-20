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

func loadMonsters153(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_153.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #153 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Archelon header continues on p. 345, and the Archmage prints
// XP 8,000 for CR 12 although the table on p. 256 gives 8,400.
var monsters153Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Animated Armor", [2]int{259, 259}, "Medium Construct, Unaligned|AC 18 Initiative +2 (12)|HP 33 (6d8 + 6)|Speed 25 ft.|Str 14 +2 +2 Dex 11 +0 +0 Con 13 +1 +1|Int 1 −5 −5 Wis 3 −4 −4 Cha 1 −5 −5|Immunities Poison, Psychic; Charmed, Deafened, Exhaustion, Frightened, Paralyzed, Petrified, Poisoned|Senses Blindsight 60 ft.; Passive Perception 6|Languages None|CR 1 (XP 200; PB +2)"},
	{"Animated Flying Sword", [2]int{259, 259}, "Small Construct, Unaligned|AC 17 Initiative +4 (14)|HP 14 (4d6)|Speed 5 ft., Fly 50 ft. (hover)|Str 12 +1 +1 Dex 15 +2 +4 Con 11 +0 +0|Int 1 −5 −5 Wis 5 −3 −3 Cha 1 −5 −5|Immunities Poison, Psychic; Charmed, Deafened, Exhaustion, Frightened, Paralyzed, Petrified, Poisoned|Senses Blindsight 60 ft.; Passive Perception 7|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Animated Rug of Smothering", [2]int{259, 259}, "Large Construct, Unaligned|AC 12 Initiative +4 (14)|HP 27 (5d10)|Speed 10 ft.|Str 17 +3 +3 Dex 14 +2 +2 Con 10 +0 +0|Int 1 −5 −5 Wis 3 −4 −4 Cha 1 −5 −5|Immunities Poison, Psychic; Charmed, Deafened, Exhaustion, Frightened, Paralyzed, Petrified, Poisoned|Senses Blindsight 60 ft.; Passive Perception 6|Languages None|CR 2 (XP 450; PB +2)"},
	{"Ankheg", [2]int{259, 259}, "Large Monstrosity, Unaligned|AC 14 Initiative +0 (10)|HP 45 (6d10 + 12)|Speed 30 ft., Burrow 10 ft.|Str 17 +3 +3 Dex 11 +0 +0 Con 14 +2 +2|Int 1 −5 −5 Wis 13 +1 +1 Cha 6 −2 −2|Senses Darkvision 60 ft., Tremorsense 60 ft.; Passive Perception 11|Languages None|CR 2 (XP 450; PB +2)"},
	{"Ankylosaurus", [2]int{344, 344}, "Huge Beast (Dinosaur), Unaligned|AC 15 Initiative +0 (10)|HP 68 (8d12 + 16)|Speed 30 ft.|Str 19 +4 +6 Dex 11 +0 +0 Con 15 +2 +2|Int 2 −4 −4 Wis 12 +1 +1 Cha 5 −3 −3|Senses Passive Perception 11|Languages None|CR 3 (XP 700; PB +2)"},
	{"Ape", [2]int{344, 344}, "Medium Beast, Unaligned|AC 12 Initiative +2 (12)|HP 19 (3d8 + 6)|Speed 30 ft., Climb 30 ft.|Str 16 +3 +3 Dex 14 +2 +2 Con 14 +2 +2|Int 6 −2 −2 Wis 12 +1 +1 Cha 7 −2 −2|Skills Athletics +5, Perception +3|Senses Passive Perception 13|Languages None|CR 1/2 (XP 100; PB +2)"},
	{"Archelon", [2]int{344, 345}, "Huge Beast (Dinosaur), Unaligned|AC 17 Initiative +3 (13)|HP 90 (12d12 + 12)|Speed 20 ft., Swim 80 ft.|Str 18 +4 +4 Dex 16 +3 +3 Con 13 +1 +1|Int 4 −3 −3 Wis 14 +2 +2 Cha 6 −2 −2|Skills Stealth +5|Senses Passive Perception 12|Languages None|CR 4 (XP 1,100; PB +2)"},
	{"Archmage", [2]int{305, 305}, "Medium or Small Humanoid (Wizard), Neutral|AC 17 Initiative +7 (17)|HP 170 (31d8 + 31)|Speed 30 ft.|Str 10 +0 +0 Dex 14 +2 +2 Con 12 +1 +1|Int 20 +5 +9 Wis 15 +2 +6 Cha 16 +3 +3|Skills Arcana +13, History +9, Perception +6|Immunities Psychic; Charmed (with Mind Blank)|Gear Wand|Senses Passive Perception 16|Languages Common plus five other languages|CR 12 (XP 8,000; PB +4)"},
	{"Assassin", [2]int{260, 260}, "Medium or Small Humanoid, Neutral|AC 16 Initiative +10 (20)|HP 97 (15d8 + 30)|Speed 30 ft.|Str 11 +0 +0 Dex 18 +4 +7 Con 14 +2 +2|Int 16 +3 +6 Wis 11 +0 +0 Cha 10 +0 +0|Skills Acrobatics +7, Perception +6, Stealth +10|Resistances Poison|Gear Light Crossbow, Shortsword, Studded Leather Armor|Senses Passive Perception 16|Languages Common, Thieves’ Cant|CR 8 (XP 3,900; PB +3)"},
	{"Awakened Shrub", [2]int{260, 260}, "Small Plant, Neutral|AC 9 Initiative −1 (9)|HP 10 (3d6)|Speed 20 ft.|Str 3 −4 −4 Dex 8 −1 −1 Con 11 +0 +0|Int 10 +0 +0 Wis 10 +0 +0 Cha 6 −2 −2|Vulnerabilities Fire|Resistances Piercing|Senses Passive Perception 10|Languages Common plus one other language|CR 0 (XP 10; PB +2)"},
	{"Awakened Tree", [2]int{260, 260}, "Huge Plant, Neutral|AC 13 Initiative −2 (8)|HP 59 (7d12 + 14)|Speed 20 ft.|Str 19 +4 +4 Dex 6 −2 −2 Con 15 +2 +2|Int 10 +0 +0 Wis 10 +0 +0 Cha 7 −2 −2|Vulnerabilities Fire|Resistances Bludgeoning, Piercing|Senses Passive Perception 10|Languages Common plus one other language|CR 2 (XP 450; PB +2)"},
	{"Axe Beak", [2]int{260, 260}, "Large Monstrosity, Unaligned|AC 11 Initiative +1 (11)|HP 19 (3d10 + 3)|Speed 50 ft.|Str 14 +2 +2 Dex 12 +1 +1 Con 12 +1 +1|Int 2 −4 −4 Wis 10 +0 +0 Cha 5 −3 −3|Senses Passive Perception 10|Languages None|CR 1/4 (XP 50; PB +2)"},
}

// Typed values of every trait, action, bonus action and reaction, reviewed
// against the same pages: limited uses, attack bonus/reach/range/damage, save
// ability and DC, target shape size/width/range, dice, spell lists and
// dependency keys. Keys in brackets are an effect's own requirements.
var monsters153Entries = map[string][]string{
	"Animated Armor": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|slam|Slam; melee +4 reach 5 range 0/0 5 (1d6 + 2) bludgeoning; effect/bludgeoning-damage 5 (1d6 + 2)",
	},
	"Animated Flying Sword": {
		"action|slash|Slash; melee +4 reach 5 range 0/0 6 (1d8 + 2) slashing; effect/slashing-damage 6 (1d8 + 2)",
	},
	"Animated Rug of Smothering": {
		"action|smother|Smother; melee +5 reach 5 range 0/0 10 (2d6 + 3) bludgeoning; effect/bludgeoning-damage 10 (2d6 + 3), requirement/target-size, condition/grappled, effect/condition-instead-of-damage, condition/blinded, condition/restrained, effect/suffocating, effect/start-of-turn-damage 10 (2d6 + 3), limit/one-grappled-creature, requirement/not-grappling, effect/halve-damage-taken, effect/share-damage-taken",
	},
	"Ankheg": {
		"trait|trait/tunneler|Tunneler; movement/burrow-through-rock, movement/half-burrow-speed, effect/leave-tunnel",
		"action|bite|Bite; melee +5 reach 5 range 0/0 10 (2d6 + 3) slashing; effect/slashing-damage 10 (2d6 + 3), effect/acid-damage 3 (1d6), attack/advantage-against-grappled-target, condition/grappled, requirement/target-size",
		"action|acid-spray|Acid Spray (Recharge 6); dex DC 12; line 30/5/0; save/dex, area/line, limit/recharge, effect/acid-damage 14 (4d6), save/half-damage",
	},
	"Ankylosaurus": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|tail|Tail; melee +6 reach 10 range 0/0 9 (1d10 + 4) bludgeoning; effect/bludgeoning-damage 9 (1d10 + 4), condition/prone, requirement/target-size",
	},
	"Ape": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|fist|Fist; melee +5 reach 5 range 0/0 5 (1d4 + 3) bludgeoning; effect/bludgeoning-damage 5 (1d4 + 3)",
		"action|rock|Rock (Recharge 6); ranged +5 reach 0 range 25/50 10 (2d6 + 3) bludgeoning; effect/bludgeoning-damage 10 (2d6 + 3), limit/recharge",
	},
	"Archelon": {
		"trait|trait/amphibious|Amphibious; ",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +6 reach 5 range 0/0 14 (3d6 + 4) piercing; effect/piercing-damage 14 (3d6 + 4)",
	},
	"Archmage": {
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-or-ranged-attack",
		"action|arcane-burst|Arcane Burst; effect/force-damage 27 (4d10 + 5)",
		"action|spellcasting|Spellcasting; int DC 17 +0; at will: Detect Magic, Detect Thoughts, Disguise Self, Invisibility, Light, Mage Armor (included in AC), Mage Hand, Prestidigitation; 2/day: Fly, Lightning Bolt (level 7 version); 1/day: Cone of Cold (level 9 version), Mind Blank (cast before combat), Scrying, Teleport; " +
			"spell/detect-magic [spellcasting/at-will], spell/detect-thoughts [spellcasting/at-will], spell/disguise-self [spellcasting/at-will], spell/invisibility [spellcasting/at-will], spell/light [spellcasting/at-will], spell/mage-armor [spellcasting/at-will spellcasting/included-in-ac], spell/mage-hand [spellcasting/at-will], spell/prestidigitation [spellcasting/at-will], " +
			"spell/fly [limit/per-day], spell/lightning-bolt [limit/per-day spellcasting/higher-level-version], spell/cone-of-cold [limit/per-day spellcasting/higher-level-version], spell/mind-blank [limit/per-day spellcasting/cast-before-combat defense/immunity-with-mind-blank], spell/scrying [limit/per-day], spell/teleport [limit/per-day]",
		"bonus|misty-step|Misty Step (3/Day); action/spellcasting, spell/misty-step, limit/per-day",
		"reaction|protective-magic|Protective Magic (3/Day); action/spellcasting, spell/counterspell, spell/shield, reaction/spell-trigger, limit/per-day",
	},
	"Assassin": {
		"trait|trait/evasion|Evasion; save/dex, save/no-damage-on-success, save/half-damage-on-failure, requirement/not-incapacitated",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|shortsword|Shortsword; melee +7 reach 5 range 0/0 7 (1d6 + 4) piercing; effect/piercing-damage 7 (1d6 + 4), effect/poison-damage 17 (5d6), condition/poisoned",
		"action|light-crossbow|Light Crossbow; ranged +7 reach 0 range 80/320 8 (1d8 + 4) piercing; effect/piercing-damage 8 (1d8 + 4), effect/poison-damage 21 (6d6)",
		"bonus|cunning-action|Cunning Action; action/dash, action/disengage, action/hide",
	},
	"Awakened Shrub": {
		"action|rake|Rake; effect/slashing-damage",
	},
	"Awakened Tree": {
		"action|slam|Slam; melee +6 reach 10 range 0/0 14 (3d6 + 4) bludgeoning; effect/bludgeoning-damage 14 (3d6 + 4)",
	},
	"Axe Beak": {
		"action|beak|Beak; melee +4 reach 5 range 0/0 6 (1d8 + 2) slashing; effect/slashing-damage 6 (1d8 + 2)",
	},
}

// Printed values the #143 contract has no typed member for: a fixed damage
// number, an attack that is melee or ranged, escape DCs, size limits and an
// attack-roll condition. Each stays in its entry's retained source text.
var monsters153Printed = []struct{ name, entry, text string }{
	{"Animated Rug of Smothering", "smother", "If the target is a Medium or smaller creature, the rug can give it the Grappled condition (escape DC 13) instead of dealing damage."},
	{"Ankheg", "trait/tunneler", "at half its Burrow Speed and leaves a 10-foot-diameter tunnel"},
	{"Ankheg", "bite", "Melee Attack Roll: +5 (with Advantage if the target is Grappled by the ankheg), reach 5 ft."},
	{"Ankheg", "bite", "If the target is a Large or smaller creature, it has the Grappled condition (escape DC 13)."},
	{"Ankylosaurus", "tail", "If the target is a Huge or smaller creature, it has the Prone condition."},
	{"Archmage", "arcane-burst", "Arcane Burst. Melee or Ranged Attack Roll: +9, reach 5 ft. or range 150 ft. Hit: 27 (4d10 + 5) Force damage."},
	{"Archmage", "misty-step", "Misty Step (3/Day). The mage casts Misty Step, using the same spellcasting ability as Spellcasting."},
	{"Archmage", "protective-magic", "Protective Magic (3/Day). The archmage casts Counterspell or Shield in response to the spell’s trigger"},
	{"Assassin", "shortsword", "and the target has the Poisoned condition until the start of the assassin’s next turn."},
	{"Assassin", "cunning-action", "Cunning Action. The assassin takes the Dash, Disengage, or Hide action."},
	{"Awakened Shrub", "rake", "Rake. Melee Attack Roll: +1, reach 5 ft. Hit: 1 Slashing damage."},
}

func monsterSigned153(n int) string {
	if n < 0 {
		return fmt.Sprintf("−%d", -n)
	}
	return fmt.Sprintf("+%d", n)
}

func monsterThousands153(n int) string {
	if n < 1000 {
		return fmt.Sprint(n)
	}
	return fmt.Sprintf("%d,%03d", n/1000, n%1000)
}

func monsterDice153(d MonsterHP) string {
	if d.Modifier == 0 {
		return fmt.Sprintf("%d (%dd%d)", d.Average, d.Dice, d.Sides)
	}
	return fmt.Sprintf("%d (%dd%d + %d)", d.Average, d.Dice, d.Sides, d.Modifier)
}

func monsterTitle153(s string) string {
	words := strings.Split(strings.ReplaceAll(s, "-", " "), " ")
	for i, w := range words {
		words[i] = strings.ToUpper(w[:1]) + w[1:]
	}
	return strings.Join(words, " ")
}

func monsterTitles153(values []string, sep string) string {
	out := make([]string, len(values))
	for i, v := range values {
		out[i] = monsterTitle153(v)
	}
	return strings.Join(out, sep)
}

func monsterSorted153[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// monsterHeader153 prints typed fields in the source's stat-block notation.
func monsterHeader153(m Monster) string {
	kind := monsterTitle153(m.Type)
	if len(m.Tags) > 0 {
		kind += " (" + strings.Join(m.Tags, ", ") + ")"
	}
	lines := []string{
		monsterTitles153(m.Sizes, " or ") + " " + kind + ", " + m.Alignment,
		fmt.Sprintf("AC %d Initiative %s (%d)", m.AC, monsterSigned153(m.Initiative.Modifier), m.Initiative.Score),
		"HP " + monsterDice153(m.HP),
	}
	speed := fmt.Sprintf("Speed %d ft.", m.Speeds["walk"])
	for _, mode := range monsterSorted153(m.Speeds) {
		if mode != "walk" {
			speed += fmt.Sprintf(", %s %d ft.", monsterTitle153(mode), m.Speeds[mode])
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
			cells = append(cells, fmt.Sprintf("%s %d %s %s", monsterTitle153(key), a.Score, monsterSigned153(a.Modifier), monsterSigned153(a.Save)))
		}
		lines = append(lines, strings.Join(cells, " "))
	}
	if m.Skills != nil {
		skills := []string{}
		for _, skill := range monsterSorted153(m.Skills) {
			skills = append(skills, monsterTitle153(skill)+" "+monsterSigned153(m.Skills[skill]))
		}
		lines = append(lines, "Skills "+strings.Join(skills, ", "))
	}
	if m.Vulnerabilities != nil {
		lines = append(lines, "Vulnerabilities "+monsterTitles153(m.Vulnerabilities, ", "))
	}
	if m.Resistances != nil {
		lines = append(lines, "Resistances "+monsterTitles153(m.Resistances, ", "))
	}
	conditions := monsterTitles153(m.ConditionImmunities, ", ")
	if slices.Contains(m.RequiredMechanics, "defense/immunity-with-mind-blank") {
		conditions += " (with Mind Blank)"
	}
	if immunities := slices.DeleteFunc([]string{monsterTitles153(m.Immunities, ", "), conditions}, func(s string) bool { return s == "" }); len(immunities) > 0 {
		lines = append(lines, "Immunities "+strings.Join(immunities, "; "))
	}
	if len(m.Gear) > 0 {
		lines = append(lines, "Gear "+strings.Join(m.Gear, ", "))
	}
	senses := []string{}
	for _, sense := range monsterSorted153(m.Senses.Special) {
		senses = append(senses, fmt.Sprintf("%s %d ft.", monsterTitle153(sense), m.Senses.Special[sense]))
	}
	passive := fmt.Sprintf("Passive Perception %d", m.Senses.PassivePerception)
	if len(senses) > 0 {
		passive = strings.Join(senses, ", ") + "; " + passive
	}
	languages := "Languages None"
	if len(m.Languages) > 0 {
		languages = "Languages " + strings.Join(m.Languages, ", ")
	}
	lines = append(lines, "Senses "+passive, languages, fmt.Sprintf("CR %s (XP %s; PB +%d)", m.CR, monsterThousands153(m.XP), m.ProficiencyBonus))
	return strings.Join(lines, "|")
}

func monsterEntry153(section, id string, x MonsterMechanic, attack *MonsterAttack, sc *MonsterSpellcasting, effects []MonsterMechanic) string {
	head := x.Name
	if u := x.Uses; u != nil {
		switch {
		case u.Recharge == 6:
			head += " (Recharge 6)"
		case u.Recharge > 0:
			head += fmt.Sprintf(" (Recharge %d–6)", u.Recharge)
		default:
			head += fmt.Sprintf(" (%d/Day)", u.PerDay)
		}
	}
	parts := []string{head}
	if attack != nil {
		parts = append(parts, fmt.Sprintf("%s +%d reach %d range %d/%d %s %s", attack.Kind, attack.Bonus, attack.Reach, attack.Range, attack.LongRange, monsterDice153(attack.Damage), attack.DamageType))
	}
	if x.Save != nil {
		parts = append(parts, fmt.Sprintf("%s DC %d", x.Save.Ability, x.Save.DC))
	}
	if t := x.Target; t != nil {
		parts = append(parts, fmt.Sprintf("%s %d/%d/%d", t.Shape, t.Size, t.Width, t.Range))
	}
	if x.Dice != nil {
		parts = append(parts, monsterDice153(*x.Dice))
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
			key += " " + monsterDice153(*e.Dice)
		}
		if len(e.Requires) > 0 {
			key += " [" + strings.Join(e.Requires, " ") + "]"
		}
		keys = append(keys, key)
	}
	return section + "|" + id + "|" + strings.Join(append(parts, strings.Join(keys, ", ")), "; ")
}

func monsterEntries153(m Monster) []string {
	out := []string{}
	for _, trait := range m.Traits {
		out = append(out, monsterEntry153("trait", trait.Key, trait, nil, nil, nil))
	}
	for _, section := range []struct {
		label   string
		actions []MonsterAction
	}{{"action", m.Actions}, {"bonus", m.BonusActions}, {"reaction", m.Reactions}} {
		for _, a := range section.actions {
			out = append(out, monsterEntry153(section.label, a.ID, MonsterMechanic{Name: a.Name, Uses: a.Uses, Save: a.Save, Target: a.Target}, a.Attack, a.Spellcasting, a.Effects))
		}
	}
	return out
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters153SourceValues(t *testing.T) {
	c := loadMonsters153(t)
	if len(c.Monsters) != 12 || len(monsters153Headers) != 12 || len(monsters153Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #153") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	gear := map[string][]string{"Archmage": {"Wand"}, "Assassin": {"Light Crossbow", "Shortsword", "Studded Leather Armor"}}
	sections := map[string][2]int{"Archmage": {1, 1}, "Assassin": {1, 0}}
	for i, want := range monsters153Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader153(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries153(m); !slices.Equal(got, monsters153Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters153Entries[want.name])
			}
			if !slices.Equal(m.Gear, gear[m.Name]) && (len(m.Gear) != 0 || gear[m.Name] != nil) {
				t.Fatalf("gear differs: %v", m.Gear)
			}
			if len(m.BonusActions) != sections[m.Name][0] || len(m.Reactions) != sections[m.Name][1] {
				t.Fatal("bonus action or reaction sections differ from the source")
			}
			if m.Telepathy != 0 || m.LairXP != 0 || m.Legendary != nil || m.LegendaryActions != nil {
				t.Fatal("source lists no telepathy, lair XP or legendary actions")
			}
		})
	}
	byName := map[string]Monster{}
	for _, m := range c.Monsters {
		byName[m.Name] = m
	}
	for _, p := range monsters153Printed {
		text := ""
		for _, trait := range byName[p.name].Traits {
			if trait.Key == p.entry {
				text = trait.Text
			}
		}
		for _, a := range byName[p.name].allActions() {
			if a.ID == p.entry {
				text = a.Text
			}
		}
		if !strings.Contains(text, p.text) {
			t.Fatalf("%s %s lost printed value %q", p.name, p.entry, p.text)
		}
	}
	// A fixed damage number and a melee-or-ranged attack have no typed attack in
	// the #143 contract, so they must not masquerade as ordinary weapon attacks.
	for name, key := range map[string]string{"Awakened Shrub": "action/fixed-damage-melee-attack", "Archmage": "action/melee-or-ranged-attack"} {
		found := false
		for _, a := range byName[name].Actions {
			if a.Key == key {
				found = a.Attack == nil
			}
			if a.Attack != nil {
				t.Fatalf("%s has no ordinary weapon attack in the source", name)
			}
		}
		if !found {
			t.Fatalf("%s lost %s", name, key)
		}
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters153MatchesDecisionInventory(t *testing.T) {
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
	families := []string{}
	for _, r := range inventory.Records {
		if r.CatalogueIssue == 153 {
			if r.Status != "unsupported" || r.ID != "srd-5.2.1:content/"+r.Family+"/"+strings.ReplaceAll(strings.ToLower(r.Name), " ", "-") {
				t.Fatalf("%s: inventory identity or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters153(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the three p. 344 blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("monsters monsters monsters monsters animals animals animals monsters monsters monsters monsters monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #153 lists these names with the page each block starts on.
	names := []string{"Animated Armor", "Animated Flying Sword", "Animated Rug of Smothering", "Ankheg", "Ankylosaurus", "Ape", "Archelon", "Archmage", "Assassin", "Awakened Shrub", "Awakened Tree", "Axe Beak"}
	pages := []int{259, 259, 259, 259, 344, 344, 344, 305, 260, 260, 260, 260}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue lists %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic that holds for these twelve blocks (PDF pp. 254–256): PB
// follows CR, saves and skills add PB to the ability modifier, and typed attack
// rolls use the better of Strength and Dexterity. The Archmage is the only block
// whose printed XP and Initiative do not follow those rules.
func TestMonsters153DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"0": 2, "1/4": 2, "1/2": 2, "1": 2, "2": 2, "3": 2, "4": 2, "8": 3, "12": 4}
	skillAbility := map[string]string{"acrobatics": "dex", "arcana": "int", "athletics": "str", "history": "int", "perception": "wis", "stealth": "dex"}
	hitDie := map[string]int{"small": 6, "medium": 8, "large": 10, "huge": 12}
	for _, m := range loadMonsters153(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			p := m.ProficiencyBonus
			if pb[m.CR] != p {
				t.Fatal("PB does not follow CR")
			}
			if valid := encounterValidXP(m.CR, int64(m.XP)); valid == (m.Name == "Archmage") {
				t.Fatal("only the Archmage prints XP that differs from the CR table")
			}
			if m.HP.Sides != hitDie[m.Sizes[0]] || m.HP.Modifier != m.HP.Dice*m.Abilities["con"].Modifier {
				t.Fatal("Hit Dice do not follow size and Constitution")
			}
			for key, a := range m.Abilities {
				if a.Save != a.Modifier && a.Save != a.Modifier+p {
					t.Fatalf("%s save is not the modifier with optional PB", key)
				}
			}
			extra := m.Initiative.Modifier - m.Abilities["dex"].Modifier
			if follows := extra == 0 || extra == p || extra == 2*p; follows == (m.Name == "Archmage") || m.Name == "Archmage" && m.Initiative.Modifier != 7 {
				t.Fatal("only the Archmage's printed Initiative +7 is not Dexterity with optional PB or twice PB")
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
func TestMonsters153DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters153(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			if m.RequiredMechanics[0] != "monster/stat-block" {
				t.Fatal("stat-block dependency must come first")
			}
			quoted := func(parent string, x MonsterMechanic) {
				if !strings.Contains(parent, strings.TrimSuffix(x.Text, ".")) {
					t.Fatalf("%q is not quoted from its entry", x.Text)
				}
				if x.Dice != nil && !strings.Contains(x.Text, monsterDice153(*x.Dice)) {
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

// Catalogue presence grants nothing: Available consults no implemented key, and
// the register still reports both families unsupported. The Axe Beak is the one
// block whose keys are a subset of the Bandit fixture's keys (one melee attack
// dealing Slashing damage); the Bandit resolver is bound to the Bandit record, so
// those keys do not make the Axe Beak executable.
func TestMonsters153RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters153(t).Monsters {
		if m.Available() || m.availableWith(bandit) != (m.Name == "Axe Beak") {
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
	families := 0
	for _, e := range r.Entries {
		if e.Family == "monsters" || e.Family == "animals" {
			families++
			if e.Status != "unsupported" || len(e.Fixtures) != 0 {
				t.Fatalf("catalogue data promoted the %s family", e.Family)
			}
		}
	}
	if families < 2 {
		t.Fatal("register lost the monster or animal family")
	}
}

// The Archmage keeps its printed CR 12 and XP 8,000 (PDF p. 305). Only that
// exact pair is accepted, the table value is not substituted, no other record
// gains an exception, and the encounter budget still rejects the mismatch.
func TestMonsters153PrintedXPException(t *testing.T) {
	c := loadMonsters153(t)
	archmage, assassin := c.Monsters[7], c.Monsters[8]
	if archmage.Name != "Archmage" || archmage.CR != "12" || archmage.XP != 8000 || archmage.Validate() != nil {
		t.Fatal("Archmage must keep the printed CR 12 and XP 8,000")
	}
	if !reflect.DeepEqual(monsterPrintedXP, map[string][2]string{archmage.ID: {"12", "8000"}}) {
		t.Fatal("printed XP exceptions changed")
	}
	for _, mutate := range []func(*Monster){
		func(m *Monster) { m.XP = 8400 },
		func(m *Monster) { m.XP = 8001 },
		func(m *Monster) { m.CR, m.ProficiencyBonus = "11", 4 },
	} {
		m := archmage
		mutate(&m)
		if m.Validate() == nil {
			t.Fatal("Archmage accepted a CR or XP the source does not print")
		}
	}
	assassin.XP = 8000
	if assassin.Validate() == nil {
		t.Fatal("the printed exception leaked to another record")
	}
	if _, err := EncounterMonsterFromRecord(archmage, 1); err == nil {
		t.Fatal("encounter budget accepted XP that differs from the CR table")
	}
	if _, err := EncounterMonsterFromRecord(c.Monsters[8], 1); err != nil {
		t.Fatal("encounter budget must still read a consistent record:", err)
	}
}

func TestMonsters153MalformedFixtures(t *testing.T) {
	// Indexes: 0 Animated Armor, 1 Animated Flying Sword, 3 Ankheg, 4 Ankylosaurus,
	// 5 Ape, 7 Archmage, 8 Assassin, 9 Awakened Shrub.
	tests := []struct {
		name   string
		mutate func(*MonsterCatalogue)
	}{
		{"duplicate ID", func(c *MonsterCatalogue) { c.Monsters[1].ID = c.Monsters[0].ID }},
		{"unknown family", func(c *MonsterCatalogue) { c.Monsters[4].ID = "srd-5.2.1:content/beasts/ankylosaurus" }},
		{"empty tags", func(c *MonsterCatalogue) { c.Monsters[4].Tags = []string{} }},
		{"duplicate immunity", func(c *MonsterCatalogue) { c.Monsters[0].Immunities = []string{"poison", "poison"} }},
		{"unknown vulnerability", func(c *MonsterCatalogue) { c.Monsters[9].Vulnerabilities[0] = "flame" }},
		{"unknown resistance", func(c *MonsterCatalogue) { c.Monsters[9].Resistances[0] = "stabbing" }},
		{"unknown condition", func(c *MonsterCatalogue) { c.Monsters[0].ConditionImmunities[0] = "dazed" }},
		{"unknown skill", func(c *MonsterCatalogue) { c.Monsters[5].Skills["climbing"] = 5 }},
		{"empty skills", func(c *MonsterCatalogue) { c.Monsters[5].Skills = map[string]int{} }},
		{"hover without fly", func(c *MonsterCatalogue) { delete(c.Monsters[1].Speeds, "fly") }},
		{"modifier", func(c *MonsterCatalogue) { c.Monsters[1].Abilities["dex"] = MonsterAbility{15, 3, 4} }},
		{"save off modifier", func(c *MonsterCatalogue) { c.Monsters[1].Abilities["dex"] = MonsterAbility{15, 2, 5} }},
		{"HP average", func(c *MonsterCatalogue) { c.Monsters[0].HP.Average = 34 }},
		{"XP off CR", func(c *MonsterCatalogue) { c.Monsters[0].XP = 450 }},
		{"empty bonus actions", func(c *MonsterCatalogue) { c.Monsters[8].BonusActions = []MonsterAction{} }},
		{"empty reactions", func(c *MonsterCatalogue) { c.Monsters[7].Reactions = []MonsterAction{} }},
		{"duplicate action across sections", func(c *MonsterCatalogue) { c.Monsters[8].BonusActions[0].ID = "shortsword" }},
		{"no uses", func(c *MonsterCatalogue) { *c.Monsters[3].Actions[1].Uses = MonsterUses{} }},
		{"recharge", func(c *MonsterCatalogue) { c.Monsters[3].Actions[1].Uses.Recharge = 7 }},
		{"daily uses", func(c *MonsterCatalogue) { c.Monsters[7].Reactions[0].Uses.PerDay = -1 }},
		{"save ability", func(c *MonsterCatalogue) { c.Monsters[3].Actions[1].Save.Ability = "luck" }},
		{"save DC", func(c *MonsterCatalogue) { c.Monsters[3].Actions[1].Save.DC = 0 }},
		{"save DC off entry text", func(c *MonsterCatalogue) { c.Monsters[3].Actions[1].Save.DC = 13 }},
		{"target shape", func(c *MonsterCatalogue) { c.Monsters[3].Actions[1].Target.Shape = "cube" }},
		{"line width", func(c *MonsterCatalogue) { c.Monsters[3].Actions[1].Target.Width = 0 }},
		{"attack damage", func(c *MonsterCatalogue) { c.Monsters[0].Actions[1].Attack.Damage.Average = 6 }},
		{"attack reach", func(c *MonsterCatalogue) { c.Monsters[0].Actions[1].Attack.Reach = 0 }},
		{"ranged long range", func(c *MonsterCatalogue) { c.Monsters[5].Actions[2].Attack.LongRange = 20 }},
		{"effect dice", func(c *MonsterCatalogue) { c.Monsters[3].Actions[1].Effects[3].Dice.Average = 15 }},
		{"fixed damage typed as weapon attack", func(c *MonsterCatalogue) { c.Monsters[9].Actions[0].Key = "action/melee-attack" }},
		{"unlisted requirement", func(c *MonsterCatalogue) { c.Monsters[3].Traits[0].Requires = []string{"effect/future"} }},
		{"duplicate requirement", func(c *MonsterCatalogue) {
			c.Monsters[3].Traits[0].Requires = []string{"effect/leave-tunnel", "effect/leave-tunnel"}
		}},
		{"unlisted effect", func(c *MonsterCatalogue) { c.Monsters[3].Actions[1].Effects[4].Key = "save/future" }},
		{"untyped spellcasting", func(c *MonsterCatalogue) { c.Monsters[7].Actions[2].Spellcasting = nil }},
		{"misplaced spellcasting", func(c *MonsterCatalogue) {
			c.Monsters[7].Actions[0].Spellcasting = c.Monsters[7].Actions[2].Spellcasting
		}},
		{"spell ability", func(c *MonsterCatalogue) { c.Monsters[7].Actions[2].Spellcasting.Ability = "" }},
		{"spell DC", func(c *MonsterCatalogue) { c.Monsters[7].Actions[2].Spellcasting.SaveDC = 0 }},
		{"unlisted spell", func(c *MonsterCatalogue) { c.Monsters[7].Actions[2].Spellcasting.AtWill[0] = "Wish" }},
		{"duplicate spell", func(c *MonsterCatalogue) { c.Monsters[7].Actions[2].Spellcasting.PerDay[0].Spells[0] = "Light" }},
		{"daily spell uses", func(c *MonsterCatalogue) { c.Monsters[7].Actions[2].Spellcasting.PerDay[0].Uses = 0 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters153(t)
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
	raw, err := os.ReadFile("testdata/monsters_153.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, edit := range [][2]string{
		{`"recharge": 6`, `"recharge": null`},
		{`"per_day": 3,`, ``},
		{`"dc": 12`, `"dc": null`},
		{`"width": 5,`, ``},
		{`"attack_bonus": 0,`, ``},
		{`"save_dc": 17,`, `"save_dc": 17, "save_dc": 18,`},
		{`"uses": 2,`, ``},
		{`"long_range": 50,`, ``},
		{`"hover": true,`, `"hover": true, "swoop": true,`},
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

// The batch file composes with the #143 Bandit catalogue and the #151 batch
// without entering the embedded runtime catalogue or changing the Bandit
// fixture's pinned revision.
func TestMonsters153ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters153(t)
	earlier, err := os.ReadFile("testdata/monsters_151.json")
	if err != nil {
		t.Fatal(err)
	}
	first, err := DecodeMonsters(bytes.NewReader(earlier))
	if err != nil {
		t.Fatal(err)
	}
	firstSource, source := first.Source, batch.Source
	firstSource.Changes, source.Changes = "", ""
	b02, err := LoadMonsterBatch02()
	if err != nil {
		t.Fatal(err)
	}
	var second []Monster
	for _, sb := range b02.StatBlocks {
		second = append(second, sb.Monster)
	}
	all := MonsterCatalogue{SchemaVersion: 1, Source: batch.Source, Monsters: slices.Concat(base.Monsters, first.Monsters, second, batch.Monsters)}
	if err := all.Validate(); err != nil || len(all.Monsters) != 25+len(second) || firstSource != source {
		t.Fatal("batch must compose with the Bandit, #151 and #152 records under one source pin:", err)
	}
	all.Monsters = append(all.Monsters, first.Monsters[0])
	if all.Validate() == nil {
		t.Fatal("composition accepted a duplicate earlier record")
	}
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
	// The Assassin shares the Bandit's action names; neither they nor any other
	// catalogued action can be resolved through the Bandit fixture.
	for _, action := range []string{"shortsword", "slam", "arcane-burst"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: action}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", action)
		}
	}
}
