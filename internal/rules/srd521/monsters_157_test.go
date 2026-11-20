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

func loadMonsters157(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_157.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #157 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Constrictor Snake's block starts on p. 346 and the Wyrmling's,
// Couatl's and Darkmantle's run onto the next page.
var monsters157Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Commoner", [2]int{275, 275}, "Medium or Small Humanoid, Neutral|AC 10 Initiative +0 (10)|HP 4 (1d8)|Speed 30 ft.|Str 10 +0 +0 Dex 10 +0 +0 Con 10 +0 +0|Int 10 +0 +0 Wis 10 +0 +0 Cha 10 +0 +0|Gear Club|Senses Passive Perception 10|Languages Common|CR 0 (XP 10; PB +2)"},
	{"Constrictor Snake", [2]int{346, 347}, "Large Beast, Unaligned|AC 13 Initiative +2 (12)|HP 13 (2d10 + 2)|Speed 30 ft., Swim 30 ft.|Str 15 +2 +2 Dex 14 +2 +2 Con 12 +1 +1|Int 1 −5 −5 Wis 10 +0 +0 Cha 3 −4 −4|Skills Perception +2, Stealth +4|Senses Blindsight 10 ft.; Passive Perception 12|Languages None|CR 1/4 (XP 50; PB +2)"},
	{"Copper Dragon Wyrmling", [2]int{275, 276}, "Medium Dragon (Metallic), Chaotic Good|AC 16 Initiative +3 (13)|HP 22 (4d8 + 4)|Speed 30 ft., Climb 30 ft., Fly 60 ft.|Str 15 +2 +2 Dex 12 +1 +3 Con 13 +1 +1|Int 14 +2 +2 Wis 11 +0 +2 Cha 13 +1 +1|Skills Perception +4, Stealth +3|Immunities Acid|Senses Blindsight 10 ft., Darkvision 60 ft.; Passive Perception 14|Languages Draconic|CR 1 (XP 200; PB +2)"},
	{"Couatl", [2]int{277, 278}, "Medium Celestial, Lawful Good|AC 19 Initiative +5 (15)|HP 60 (8d8 + 24)|Speed 30 ft., Fly 90 ft.|Str 16 +3 +3 Dex 20 +5 +5 Con 17 +3 +5|Int 18 +4 +4 Wis 20 +5 +7 Cha 18 +4 +4|Resistances Bludgeoning, Piercing, Slashing|Immunities Psychic, Radiant|Senses Truesight 120 ft.; Passive Perception 15|Languages All; telepathy 120 ft.|CR 4 (XP 1,100; PB +2)"},
	{"Crab", [2]int{347, 347}, "Tiny Beast, Unaligned|AC 11 Initiative +0 (10)|HP 3 (1d4 + 1)|Speed 20 ft., Swim 20 ft.|Str 6 −2 −2 Dex 11 +0 +0 Con 12 +1 +1|Int 1 −5 −5 Wis 8 −1 −1 Cha 2 −4 −4|Skills Stealth +2|Senses Blindsight 30 ft.; Passive Perception 9|Languages None|CR 0 (XP 10; PB +2)"},
	{"Crocodile", [2]int{347, 347}, "Large Beast, Unaligned|AC 12 Initiative +0 (10)|HP 13 (2d10 + 2)|Speed 20 ft., Swim 30 ft.|Str 15 +2 +2 Dex 10 +0 +0 Con 13 +1 +3|Int 2 −4 −4 Wis 10 +0 +0 Cha 5 −3 −3|Skills Stealth +2|Senses Passive Perception 10|Languages None|CR 1/2 (XP 100; PB +2)"},
	{"Cultist", [2]int{278, 278}, "Medium or Small Humanoid, Neutral|AC 12 Initiative +1 (11)|HP 9 (2d8)|Speed 30 ft.|Str 11 +0 +0 Dex 12 +1 +1 Con 10 +0 +0|Int 10 +0 +0 Wis 11 +0 +2 Cha 10 +0 +0|Skills Deception +2, Religion +2|Gear Leather Armor, Sickle|Senses Passive Perception 10|Languages Common|CR 1/8 (XP 25; PB +2)"},
	{"Cultist Fanatic", [2]int{278, 278}, "Medium or Small Humanoid, Neutral|AC 13 Initiative +2 (12)|HP 44 (8d8 + 8)|Speed 30 ft.|Str 11 +0 +0 Dex 14 +2 +2 Con 12 +1 +1|Int 10 +0 +0 Wis 14 +2 +4 Cha 13 +1 +1|Skills Deception +3, Persuasion +3, Religion +2|Gear Holy Symbol, Leather Armor|Senses Passive Perception 12|Languages Common|CR 2 (XP 450; PB +2)"},
	{"Darkmantle", [2]int{278, 279}, "Small Aberration, Unaligned|AC 11 Initiative +3 (13)|HP 22 (5d6 + 5)|Speed 10 ft., Fly 30 ft.|Str 16 +3 +3 Dex 12 +1 +1 Con 13 +1 +1|Int 2 −4 −4 Wis 10 +0 +0 Cha 5 −3 −3|Skills Stealth +3|Senses Blindsight 60 ft.; Passive Perception 10|Languages None|CR 1/2 (XP 100; PB +2)"},
	{"Death Dog", [2]int{279, 279}, "Medium Monstrosity, Neutral Evil|AC 12 Initiative +2 (12)|HP 39 (6d8 + 12)|Speed 40 ft.|Str 15 +2 +2 Dex 14 +2 +2 Con 14 +2 +2|Int 3 −4 −4 Wis 13 +1 +1 Cha 6 −2 −2|Skills Perception +5, Stealth +4|Immunities Blinded, Charmed, Deafened, Frightened, Stunned, Unconscious|Senses Darkvision 120 ft.; Passive Perception 15|Languages None|CR 1 (XP 200; PB +2)"},
	{"Deer", [2]int{347, 347}, "Medium Beast, Unaligned|AC 13 Initiative +3 (13)|HP 4 (1d8)|Speed 50 ft.|Str 11 +0 +0 Dex 16 +3 +3 Con 11 +0 +0|Int 2 −4 −4 Wis 14 +2 +2 Cha 5 −3 −3|Skills Perception +4|Senses Darkvision 60 ft.; Passive Perception 14|Languages None|CR 0 (XP 10; PB +2)"},
	{"Deva", [2]int{279, 279}, "Medium Celestial (Angel), Lawful Good|AC 17 Initiative +4 (14)|HP 229 (27d8 + 108)|Speed 30 ft., Fly 90 ft. (hover)|Str 18 +4 +4 Dex 18 +4 +4 Con 18 +4 +4|Int 17 +3 +3 Wis 20 +5 +9 Cha 20 +5 +9|Skills Insight +9, Perception +9|Resistances Radiant|Immunities Charmed, Exhaustion, Frightened|Senses Darkvision 120 ft.; Passive Perception 19|Languages All; telepathy 120 ft.|CR 10 (XP 5,900; PB +4)"},
}

// Typed values of every trait, action, bonus action and reaction, reviewed
// against the same pages: limited uses, attack bonus/reach/range/damage, save
// ability and DC, target shape size/width/range, dice and dependency keys. A
// ranged attack printed with one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text: the Cultist's "plus 1 Necrotic damage" is a fixed second hit (effect
// requires damage/fixed), and the Cultist Fanatic's "+4 to hit with spell
// attacks" is typed as attack_bonus with the key spellcasting/spell-attack-bonus.
var monsters157Entries = map[string][]string{
	"Commoner": {
		"trait|trait/training|Training; skill/proficiency-of-gm-choice, check/advantage-with-chosen-skill",
		"action|club|Club; melee +2 reach 5 range 0/0 2 (1d4) bludgeoning; effect/bludgeoning-damage 2 (1d4)",
	},
	"Constrictor Snake": {
		"action|bite|Bite; melee +4 reach 5 range 0/0 6 (1d8 + 2) piercing; effect/piercing-damage 6 (1d8 + 2)",
		"action|constrict|Constrict; str DC 12; creature 0/0/5; save/str, target/creature, requirement/target-size, effect/bludgeoning-damage 7 (3d4), condition/grappled",
	},
	"Copper Dragon Wyrmling": {
		"action|rend|Rend; melee +4 reach 5 range 0/0 7 (1d10 + 2) slashing; effect/slashing-damage 7 (1d10 + 2)",
		"action|acid-breath|Acid Breath (Recharge 5–6); dex DC 11; line 20/5/0; save/dex, area/line, limit/recharge, effect/acid-damage 18 (4d8), save/half-damage",
		"action|slowing-breath|Slowing Breath; con DC 11; cone 15/0/0; save/con, area/cone, effect/speed-halved, effect/no-reactions, effect/action-or-bonus-action",
	},
	"Couatl": {
		"trait|trait/shielded-mind|Shielded Mind; effect/thoughts-unreadable, effect/telepathy-only-if-allowed",
		"action|bite|Bite; melee +7 reach 5 range 0/0 11 (1d12 + 5) piercing; effect/piercing-damage 11 (1d12 + 5), condition/poisoned",
		"action|constrict|Constrict; str DC 15; creature 0/0/5; save/str, target/creature, requirement/target-size, effect/bludgeoning-damage 8 (1d6 + 5), condition/grappled, condition/restrained",
		"action|spellcasting|Spellcasting; wis DC 15 +0; at will: Detect Evil and Good, Detect Magic, Detect Thoughts, Shapechange (Beast or Humanoid form only, no Temporary Hit Points gained from the spell, and no Concentration or Temporary Hit Points required to maintain the spell); 1/day: Create Food and Water, Dream, Greater Restoration, Scrying, Sleep; spellcasting/no-components, spell/detect-evil-and-good [spellcasting/at-will], spell/detect-magic [spellcasting/at-will], spell/detect-thoughts [spellcasting/at-will], spell/shapechange [spellcasting/at-will], spell/create-food-and-water [limit/per-day], spell/dream [limit/per-day], spell/greater-restoration [limit/per-day], spell/scrying [limit/per-day], spell/sleep [limit/per-day]",
		"bonus|divine-aid|Divine Aid (2/Day); limit/per-day, action/spellcasting, spell/bless, spell/lesser-restoration, spell/sanctuary, spellcasting/no-components",
	},
	"Crab": {
		"trait|trait/amphibious|Amphibious; effect/breathe-air-and-water",
		"action|claw|Claw; melee +2 reach 5 range 0/0 1 bludgeoning; effect/bludgeoning-damage [damage/fixed]",
	},
	"Crocodile": {
		"trait|trait/hold-breath|Hold Breath; effect/hold-breath",
		"action|bite|Bite; melee +4 reach 5 range 0/0 6 (1d8 + 2) piercing; effect/piercing-damage 6 (1d8 + 2), requirement/target-size, condition/grappled, condition/restrained",
	},
	"Cultist": {
		"action|ritual-sickle|Ritual Sickle; melee +3 reach 5 range 0/0 3 (1d4 + 1) slashing; effect/slashing-damage 3 (1d4 + 1), effect/necrotic-damage [damage/fixed]",
	},
	"Cultist Fanatic": {
		"action|pact-blade|Pact Blade; melee +4 reach 5 range 0/0 6 (1d8 + 2) slashing; effect/slashing-damage 6 (1d8 + 2), effect/necrotic-damage 7 (2d6)",
		"action|spellcasting|Spellcasting; wis DC 12 +4; at will: Light, Thaumaturgy; 2/day: Command; 1/day: Hold Person; spellcasting/spell-attack-bonus, spell/light [spellcasting/at-will], spell/thaumaturgy [spellcasting/at-will], spell/command [limit/per-day], spell/hold-person [limit/per-day]",
		"bonus|spiritual-weapon|Spiritual Weapon (2/Day); limit/per-day, action/spellcasting, spell/spiritual-weapon",
	},
	"Darkmantle": {
		"action|crush|Crush; melee +5 reach 5 range 0/0 6 (1d6 + 3) bludgeoning; effect/bludgeoning-damage 6 (1d6 + 3), effect/attach, requirement/target-size, condition/blinded, effect/suffocating, effect/attack-only-attached-target, attack/advantage-while-attached, effect/speed-zero-while-attached, effect/detach, check/str-athletics, effect/self-detach",
		"action|darkness-aura|Darkness Aura (1/Day); limit/per-day, area/emanation, effect/magical-darkness, effect/concentration, effect/blocks-darkvision-and-light",
	},
	"Death Dog": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|bite|Bite; melee +4 reach 5 range 0/0 4 (1d4 + 2) piercing; con DC 12; effect/piercing-damage 4 (1d4 + 2), requirement/target-creature, save/con, condition/poisoned, effect/hit-point-maximum-no-reset, save/repeat, save/subsequent-failure, effect/hit-point-maximum-reduction 5 (1d10)",
	},
	"Deer": {
		"trait|trait/agile|Agile; effect/no-opportunity-attack-leaving-reach",
		"action|ram|Ram; melee +2 reach 5 range 0/0 2 (1d4) bludgeoning; effect/bludgeoning-damage 2 (1d4)",
	},
	"Deva": {
		"trait|trait/exalted-restoration|Exalted Restoration; effect/revive",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|holy-mace|Holy Mace; melee +8 reach 5 range 0/0 7 (1d6 + 4) bludgeoning; effect/bludgeoning-damage 7 (1d6 + 4), effect/radiant-damage 18 (4d8)",
		"action|spellcasting|Spellcasting; cha DC 17 +0; at will: Detect Evil and Good, Shapechange (Beast or Humanoid form only, no Temporary Hit Points gained from the spell, and no Concentration or Temporary Hit Points required to maintain the spell); 1/day: Commune, Raise Dead; spellcasting/no-material-components, spell/detect-evil-and-good [spellcasting/at-will], spell/shapechange [spellcasting/at-will], spell/commune [limit/per-day], spell/raise-dead [limit/per-day]",
		"bonus|divine-aid|Divine Aid (2/Day); limit/per-day, action/spellcasting, spell/cure-wounds, spell/lesser-restoration, spell/remove-curse",
	},
}

// monsterDice157 prints a fixed hit, which the source prints as a bare number
// (the Cat's Scratch), and otherwise the source's "N (XdY + Z)" notation.
func monsterDice157(d MonsterHP) string {
	switch {
	case d.Dice == 0:
		return fmt.Sprint(d.Average)
	case d.Modifier < 0:
		return fmt.Sprintf("%d (%dd%d − %d)", d.Average, d.Dice, d.Sides, -d.Modifier)
	}
	return monsterDice151(d)
}

// monsterHeader157 prints typed fields in the source's stat-block notation.
func monsterHeader157(m Monster) string {
	kind := monsterTitle151(m.Type)
	if len(m.Tags) > 0 {
		kind += " (" + strings.Join(m.Tags, ", ") + ")"
	}
	sizes := make([]string, len(m.Sizes))
	for i, size := range m.Sizes {
		sizes[i] = monsterTitle151(size)
	}
	lines := []string{
		strings.Join(sizes, " or ") + " " + kind + ", " + m.Alignment,
		fmt.Sprintf("AC %d Initiative %s (%d)", m.AC, monsterSigned151(m.Initiative.Modifier), m.Initiative.Score),
		"HP " + monsterDice157(m.HP),
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
	if m.Resistances != nil {
		lines = append(lines, "Resistances "+monsterTitles151(m.Resistances))
	}
	if immunities := slices.DeleteFunc([]string{monsterTitles151(m.Immunities), monsterTitles151(m.ConditionImmunities)}, func(s string) bool { return s == "" }); len(immunities) > 0 {
		lines = append(lines, "Immunities "+strings.Join(immunities, "; "))
	}
	if len(m.Gear) > 0 {
		lines = append(lines, "Gear "+strings.Join(m.Gear, ", "))
	}
	senses := []string{}
	for _, sense := range monsterSorted151(m.Senses.Special) {
		printed := fmt.Sprintf("%s %d ft.", monsterTitle151(sense), m.Senses.Special[sense])
		if note, ok := m.Senses.Notes[sense]; ok {
			printed += " (" + note + ")"
		}
		senses = append(senses, printed)
	}
	passive := fmt.Sprintf("Passive Perception %d", m.Senses.PassivePerception)
	if len(senses) > 0 {
		passive = strings.Join(senses, ", ") + "; " + passive
	}
	languages := "Languages None"
	if len(m.Languages) > 0 {
		languages = "Languages " + strings.Join(m.Languages, ", ")
	}
	if m.Telepathy > 0 {
		languages += fmt.Sprintf("; telepathy %d ft.", m.Telepathy)
	}
	lines = append(lines, "Senses "+passive, languages, fmt.Sprintf("CR %s (XP %s; PB +%d)", m.CR, monsterThousands151(m.XP), m.ProficiencyBonus))
	return strings.Join(lines, "|")
}

func monsterEntry157(section, id string, x MonsterMechanic, attack *MonsterAttack, sc *MonsterSpellcasting, effects []MonsterMechanic) string {
	head := x.Name
	if u := x.Uses; u != nil {
		if u.Recharge > 0 {
			head += fmt.Sprintf(" (Recharge %d–6)", u.Recharge)
		} else {
			head += fmt.Sprintf(" (%d/Day)", u.PerDay)
		}
	}
	parts := []string{head}
	if attack != nil {
		parts = append(parts, fmt.Sprintf("%s +%d reach %d range %d/%d %s %s", attack.Kind, attack.Bonus, attack.Reach, attack.Range, attack.LongRange, monsterDice157(attack.Damage), attack.DamageType))
	}
	if x.Save != nil {
		parts = append(parts, fmt.Sprintf("%s DC %d", x.Save.Ability, x.Save.DC))
	}
	if t := x.Target; t != nil {
		parts = append(parts, fmt.Sprintf("%s %d/%d/%d", t.Shape, t.Size, t.Width, t.Range))
	}
	if x.Dice != nil {
		parts = append(parts, monsterDice157(*x.Dice))
	}
	for _, d := range x.MoreDice {
		parts = append(parts, d.Key+" "+monsterDice157(d.Dice))
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
			key += " " + monsterDice157(*e.Dice)
		}
		if len(e.Requires) > 0 {
			key += " [" + strings.Join(e.Requires, ", ") + "]"
		}
		keys = append(keys, key)
	}
	return section + "|" + id + "|" + strings.Join(append(parts, strings.Join(keys, ", ")), "; ")
}

func monsterEntries157(m Monster) []string {
	out := []string{}
	for _, trait := range m.Traits {
		out = append(out, monsterEntry157("trait", trait.Key, trait, nil, nil, nil))
	}
	for _, section := range []struct {
		label   string
		actions []MonsterAction
	}{{"action", m.Actions}, {"bonus", m.BonusActions}, {"reaction", m.Reactions}} {
		for _, a := range section.actions {
			out = append(out, monsterEntry157(section.label, a.ID, MonsterMechanic{Name: a.Name, Uses: a.Uses, Save: a.Save, Target: a.Target}, a.Attack, a.Spellcasting, a.Effects))
		}
	}
	return out
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters157SourceValues(t *testing.T) {
	c := loadMonsters157(t)
	if len(c.Monsters) != 12 || len(monsters157Headers) != 12 || len(monsters157Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #157") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters157Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader157(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries157(m); !slices.Equal(got, monsters157Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters157Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.Hover != (m.Name == "Deva") || m.LairXP != 0 || m.LegendaryActions != nil || m.Legendary != nil {
				t.Fatal("source lists no vulnerabilities, lair XP or legendary actions, and only the Deva hovers")
			}
		})
		entries += len(monsters157Entries[want.name])
	}
	// Bold-italic entry headings counted on PDF pp. 275–279 and 346–347.
	if entries != 32 {
		t.Fatalf("reviewed %d entries, want 32", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters157MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 157 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters157(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the four p. 346-347 animals under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("monsters animals monsters monsters animals animals monsters monsters monsters monsters animals monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #157 cites the page each block starts on, in this order.
	names := []string{"Commoner", "Constrictor Snake", "Copper Dragon Wyrmling", "Couatl", "Crab", "Crocodile", "Cultist", "Cultist Fanatic", "Darkmantle", "Death Dog", "Deer", "Deva"}
	pages := []int{275, 346, 275, 277, 347, 347, 278, 278, 278, 279, 347, 279}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic that holds for these twelve blocks (PDF pp. 254–256):
// XP and PB follow CR, Hit Dice follow size and Constitution, saves and skills
// add PB to the ability modifier, and attack rolls add PB to the listed ability.
// Initiative adds PB (Copper Dragon Wyrmling) or nothing.
func TestMonsters157DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"0": 2, "1/8": 2, "1/4": 2, "1/2": 2, "1": 2, "2": 2, "4": 2, "10": 4}
	skillAbility := map[string]string{"deception": "cha", "insight": "wis", "perception": "wis", "persuasion": "cha", "religion": "int", "stealth": "dex"}
	hitDie := map[string]int{"tiny": 4, "small": 6, "medium": 8, "large": 10}
	// Attacks that do not use the better of Strength and Dexterity.
	attackAbility := map[string]string{"Deer/ram": "str"}
	for _, m := range loadMonsters157(t).Monsters {
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
				ability, special := attackAbility[m.Name+"/"+a.ID]
				if special {
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
func TestMonsters157DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters157(t).Monsters {
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
				if a.Attack == nil || m.ID == BanditID { // the #143 Bandit revision types damage on the attack only
					continue
				}
				first := a.Effects[0]
				if fixedMonsterDamage(a.Attack.Damage) {
					if first.Dice != nil || !slices.Equal(first.Requires, []string{"damage/fixed"}) || !strings.Contains(a.Text, fmt.Sprintf("Hit: %d %s damage.", a.Attack.Damage.Average, monsterTitle151(a.Attack.DamageType))) {
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
func TestMonsters157RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters157(t).Monsters {
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

func TestMonsters157MalformedFixtures(t *testing.T) {
	// Indexes: 0 Commoner, 1 Constrictor Snake, 2 Copper Dragon Wyrmling, 3 Couatl,
	// 4 Crab, 5 Crocodile, 6 Cultist, 7 Cultist Fanatic, 8 Darkmantle, 9 Death Dog,
	// 10 Deer, 11 Deva.
	tests := []struct {
		name   string
		mutate func(*MonsterCatalogue)
	}{
		{"duplicate ID", func(c *MonsterCatalogue) { c.Monsters[1].ID = c.Monsters[2].ID }},
		{"animal block in another family", func(c *MonsterCatalogue) { c.Monsters[1].ID = "srd-5.2.1:content/spells/constrictor-snake" }},
		{"missing languages list", func(c *MonsterCatalogue) { c.Monsters[1].Languages = nil }},
		{"fixed damage with sides", func(c *MonsterCatalogue) { c.Monsters[4].Actions[0].Attack.Damage.Sides = 4 }},
		{"fixed damage with modifier", func(c *MonsterCatalogue) { c.Monsters[4].Actions[0].Attack.Damage.Modifier = 1 }},
		{"fixed damage of zero", func(c *MonsterCatalogue) { c.Monsters[4].Actions[0].Attack.Damage.Average = 0 }},
		{"fixed damage without mechanic", func(c *MonsterCatalogue) {
			m := &c.Monsters[4]
			m.Actions[0].Effects[0].Requires = nil
			m.RequiredMechanics = slices.DeleteFunc(m.RequiredMechanics, func(key string) bool { return key == "damage/fixed" })
		}},
		{"fixed effect dice", func(c *MonsterCatalogue) { c.Monsters[4].Actions[0].Effects[0].Dice = &MonsterHP{Average: 1} }},
		{"sense note without sense", func(c *MonsterCatalogue) {
			c.Monsters[0].Senses.Notes = map[string]string{"blindsight": "unimpeded by magical Darkness"}
		}},
		{"blank sense note", func(c *MonsterCatalogue) { c.Monsters[2].Senses.Notes = map[string]string{"darkvision": " "} }},
		{"unlisted sense note", func(c *MonsterCatalogue) {
			c.Monsters[2].Senses.Notes = map[string]string{"darkvision": "impeded by magical Darkness"}
		}},
		{"empty sense notes", func(c *MonsterCatalogue) { c.Monsters[2].Senses.Notes = map[string]string{} }},
		{"melee range", func(c *MonsterCatalogue) { c.Monsters[0].Actions[0].Attack.Range = 5 }},
		{"further dice without first", func(c *MonsterCatalogue) {
			c.Monsters[0].Traits[0].MoreDice = []MonsterKeyedDice{{Key: "effect/future-damage", Dice: MonsterHP{Average: 3, Dice: 1, Sides: 4, Modifier: 1}}}
		}},
		{"XP off CR", func(c *MonsterCatalogue) { c.Monsters[11].XP = 5901 }},
		{"save off modifier", func(c *MonsterCatalogue) { c.Monsters[11].Abilities["wis"] = MonsterAbility{20, 5, 10} }},
		{"modifier", func(c *MonsterCatalogue) { c.Monsters[0].Abilities["str"] = MonsterAbility{10, 1, 1} }},
		{"HP average", func(c *MonsterCatalogue) { c.Monsters[4].HP.Average = 4 }},
		{"hover without flight", func(c *MonsterCatalogue) { c.Monsters[11].Speeds["fly"] = 0 }},
		{"empty bonus actions", func(c *MonsterCatalogue) { c.Monsters[0].BonusActions = []MonsterAction{} }},
		{"empty reactions", func(c *MonsterCatalogue) { c.Monsters[0].Reactions = []MonsterAction{} }},
		{"duplicate action across sections", func(c *MonsterCatalogue) { c.Monsters[3].BonusActions[0].ID = "bite" }},
		{"recharge", func(c *MonsterCatalogue) { c.Monsters[2].Actions[1].Uses.Recharge = 1 }},
		{"daily uses with recharge", func(c *MonsterCatalogue) { c.Monsters[2].Actions[1].Uses.PerDay = 1 }},
		{"save DC off entry text", func(c *MonsterCatalogue) { c.Monsters[1].Actions[1].Save.DC = 16 }},
		{"save ability off entry text", func(c *MonsterCatalogue) { c.Monsters[1].Actions[1].Save.Ability = "dex" }},
		{"creature size", func(c *MonsterCatalogue) { c.Monsters[1].Actions[1].Target.Size = 5 }},
		{"cone width", func(c *MonsterCatalogue) { c.Monsters[2].Actions[2].Target.Width = 5 }},
		{"effect dice", func(c *MonsterCatalogue) { c.Monsters[3].Actions[1].Effects[3].Dice.Average = 9 }},
		{"unlisted action effect", func(c *MonsterCatalogue) { c.Monsters[9].Actions[1].Effects[1].Key = "effect/future" }},
		{"untyped spellcasting", func(c *MonsterCatalogue) { c.Monsters[11].Actions[2].Spellcasting = nil }},
		{"misplaced spellcasting", func(c *MonsterCatalogue) {
			c.Monsters[11].Actions[1].Spellcasting = c.Monsters[11].Actions[2].Spellcasting
		}},
		{"spell ability", func(c *MonsterCatalogue) { c.Monsters[11].Actions[2].Spellcasting.Ability = "" }},
		{"spell DC", func(c *MonsterCatalogue) { c.Monsters[11].Actions[2].Spellcasting.SaveDC = 0 }},
		{"negative spell attack bonus", func(c *MonsterCatalogue) { c.Monsters[7].Actions[1].Spellcasting.AttackBonus = -1 }},
		{"unlisted spell", func(c *MonsterCatalogue) { c.Monsters[11].Actions[2].Spellcasting.AtWill[0] = "Wish" }},
		{"duplicate spell", func(c *MonsterCatalogue) { c.Monsters[11].Actions[2].Spellcasting.AtWill[0] = "Commune" }},
		{"daily spell uses", func(c *MonsterCatalogue) { c.Monsters[7].Actions[1].Spellcasting.PerDay[0].Uses = 0 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters157(t)
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
	raw, err := os.ReadFile("testdata/monsters_157.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, edit := range [][2]string{
		{`"dice": 0,`, ``},
		{`"sides": 0,`, `"sides": null,`},
		{`"long_range": 0,`, `"long_range": null,`},
		{`"recharge": 5`, `"recharge": null`},
		{`"per_day": 0,`, ``},
		{`"dc": 15`, `"dc": 15, "dc": 16`},
		{`"width": 0,`, ``},
		{`"special": {}`, `"special": {"darkvision": null}`},
		{`"shape": "cone",`, `"shape": "cone", "height": 5,`},
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
// to 06 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters157ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters157(t)
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
	fifth, sixth := loadMonsters155(t), loadMonsters156(t)
	for _, c := range []MonsterCatalogue{base, first, third, fourth, fifth, sixth} {
		a, b := c.Source, batch.Source
		a.Changes, b.Changes = "", ""
		if a != b {
			t.Fatal("batches must share one source pin")
		}
	}
	earlier := slices.Concat(base.Monsters, first.Monsters, second, third.Monsters, slices.Concat(fourth.Monsters[:4], fourth.Monsters[5:]), fifth.Monsters, sixth.Monsters)
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
	for _, id := range []string{"bite", "claw", "club", "crush", "holy-mace", "pact-blade", "ram", "rend", "ritual-sickle"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
