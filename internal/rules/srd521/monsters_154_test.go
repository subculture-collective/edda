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

func loadMonsters154(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_154.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #154 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised.
var monsters154Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Azer Sentinel", [2]int{261, 261}, "Medium Elemental, Lawful Neutral|AC 17 Initiative +1 (11)|HP 39 (6d8 + 12)|Speed 30 ft.|Str 17 +3 +3 Dex 12 +1 +1 Con 15 +2 +4|Int 12 +1 +1 Wis 13 +1 +1 Cha 10 +0 +0|Immunities Fire, Poison; Poisoned|Senses Passive Perception 11|Languages Primordial (Ignan)|CR 2 (XP 450; PB +2)"},
	{"Baboon", [2]int{345, 345}, "Small Beast, Unaligned|AC 12 Initiative +2 (12)|HP 3 (1d6)|Speed 30 ft., Climb 30 ft.|Str 8 −1 −1 Dex 14 +2 +2 Con 11 +0 +0|Int 4 −3 −3 Wis 12 +1 +1 Cha 6 −2 −2|Senses Passive Perception 11|Languages None|CR 0 (XP 10; PB +2)"},
	{"Badger", [2]int{345, 345}, "Tiny Beast, Unaligned|AC 11 Initiative +0 (10)|HP 5 (1d4 + 3)|Speed 20 ft., Burrow 5 ft.|Str 10 +0 +0 Dex 11 +0 +0 Con 16 +3 +3|Int 2 −4 −4 Wis 12 +1 +1 Cha 5 −3 −3|Skills Perception +3|Resistances Poison|Senses Darkvision 30 ft.; Passive Perception 13|Languages None|CR 0 (XP 10; PB +2)"},
	{"Balor", [2]int{261, 261}, "Huge Fiend (Demon), Chaotic Evil|AC 19 Initiative +14 (24)|HP 287 (23d12 + 138)|Speed 40 ft., Fly 80 ft.|Str 26 +8 +8 Dex 15 +2 +2 Con 22 +6 +12|Int 20 +5 +5 Wis 16 +3 +9 Cha 22 +6 +6|Skills Perception +9|Resistances Cold, Lightning|Immunities Fire, Poison; Charmed, Frightened, Poisoned|Senses Truesight 120 ft.; Passive Perception 19|Languages Abyssal; telepathy 120 ft.|CR 19 (XP 22,000; PB +6)"},
	{"Bandit", [2]int{261, 261}, "Medium or Small Humanoid, Neutral|AC 12 Initiative +1 (11)|HP 11 (2d8 + 2)|Speed 30 ft.|Str 11 +0 +0 Dex 12 +1 +1 Con 12 +1 +1|Int 10 +0 +0 Wis 10 +0 +0 Cha 10 +0 +0|Gear Leather Armor, Light Crossbow, Scimitar|Senses Passive Perception 10|Languages Common, Thieves’ Cant|CR 1/8 (XP 25; PB +2)"},
	{"Bandit Captain", [2]int{261, 262}, "Medium or Small Humanoid, Neutral|AC 15 Initiative +3 (13)|HP 52 (8d8 + 16)|Speed 30 ft.|Str 15 +2 +4 Dex 16 +3 +5 Con 14 +2 +2|Int 14 +2 +2 Wis 11 +0 +2 Cha 14 +2 +2|Skills Athletics +4, Deception +4|Gear Pistol, Scimitar, Studded Leather Armor|Senses Passive Perception 10|Languages Common, Thieves’ Cant|CR 2 (XP 450; PB +2)"},
	{"Barbed Devil", [2]int{262, 262}, "Medium Fiend (Devil), Lawful Evil|AC 15 Initiative +3 (13)|HP 110 (13d8 + 52)|Speed 30 ft., Climb 30 ft.|Str 16 +3 +6 Dex 17 +3 +3 Con 18 +4 +7|Int 12 +1 +1 Wis 14 +2 +5 Cha 14 +2 +5|Skills Deception +5, Insight +5, Perception +8|Resistances Cold|Immunities Fire, Poison; Poisoned|Senses Darkvision 120 ft. (unimpeded by magical Darkness); Passive Perception 18|Languages Infernal; telepathy 120 ft.|CR 5 (XP 1,800; PB +3)"},
	{"Basilisk", [2]int{262, 262}, "Medium Monstrosity, Unaligned|AC 15 Initiative −1 (9)|HP 52 (8d8 + 16)|Speed 20 ft.|Str 16 +3 +3 Dex 8 −1 −1 Con 15 +2 +2|Int 2 −4 −4 Wis 8 −1 −1 Cha 7 −2 −2|Senses Darkvision 60 ft.; Passive Perception 9|Languages None|CR 3 (XP 700; PB +2)"},
	{"Bat", [2]int{345, 345}, "Tiny Beast, Unaligned|AC 12 Initiative +2 (12)|HP 1 (1d4 − 1)|Speed 5 ft., Fly 30 ft.|Str 2 −4 −4 Dex 15 +2 +2 Con 8 −1 −1|Int 2 −4 −4 Wis 12 +1 +1 Cha 4 −3 −3|Senses Blindsight 60 ft.; Passive Perception 11|Languages None|CR 0 (XP 10; PB +2)"},
	{"Bearded Devil", [2]int{262, 263}, "Medium Fiend (Devil), Lawful Evil|AC 13 Initiative +2 (12)|HP 58 (9d8 + 18)|Speed 30 ft.|Str 16 +3 +5 Dex 15 +2 +2 Con 15 +2 +4|Int 9 −1 −1 Wis 11 +0 +0 Cha 14 +2 +4|Resistances Cold|Immunities Fire, Poison; Frightened, Poisoned|Senses Darkvision 120 ft. (unimpeded by magical Darkness); Passive Perception 10|Languages Infernal; telepathy 120 ft.|CR 3 (XP 700; PB +2)"},
	{"Behir", [2]int{263, 263}, "Huge Monstrosity, Neutral Evil|AC 17 Initiative +3 (13)|HP 168 (16d12 + 64)|Speed 50 ft., Climb 50 ft.|Str 23 +6 +6 Dex 16 +3 +3 Con 18 +4 +4|Int 7 −2 −2 Wis 14 +2 +2 Cha 12 +1 +1|Skills Perception +6, Stealth +7|Immunities Lightning|Senses Darkvision 90 ft.; Passive Perception 16|Languages Draconic|CR 11 (XP 7,200; PB +4)"},
	{"Berserker", [2]int{263, 263}, "Medium or Small Humanoid, Neutral|AC 13 Initiative +1 (11)|HP 67 (9d8 + 27)|Speed 30 ft.|Str 16 +3 +3 Dex 12 +1 +1 Con 17 +3 +3|Int 9 −1 −1 Wis 11 +0 +0 Cha 9 −1 −1|Gear Greataxe, Hide Armor|Senses Passive Perception 10|Languages Common|CR 2 (XP 450; PB +2)"},
}

// Typed values of every trait, action, bonus action and reaction, reviewed
// against the same pages: limited uses, attack bonus/reach/range/damage, save
// ability and DC, target shape size/width/range, dice and dependency keys. A
// ranged attack printed with one range repeats it; a fixed hit prints no dice.
var monsters154Entries = map[string][]string{
	"Azer Sentinel": {
		"trait|trait/fire-aura|Fire Aura; emanation 5/0/0; 5 (1d10); area/emanation, effect/fire-damage, target/creatures-of-choice, requirement/not-incapacitated",
		"trait|trait/illumination|Illumination; effect/shed-light",
		"action|burning-hammer|Burning Hammer; melee +5 reach 5 range 0/0 8 (1d10 + 3) bludgeoning; effect/bludgeoning-damage 8 (1d10 + 3), effect/fire-damage 3 (1d6)",
	},
	"Baboon": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage, requirement/ally-within-5-feet, requirement/ally-not-incapacitated",
		"action|bite|Bite; melee +1 reach 5 range 0/0 1 (1d4 − 1) piercing; effect/piercing-damage 1 (1d4 − 1)",
	},
	"Badger": {
		"action|bite|Bite; melee +2 reach 5 range 0/0 1 piercing; effect/piercing-damage [damage/fixed]",
	},
	"Balor": {
		"trait|trait/death-throes|Death Throes; dex DC 20; emanation 30/0/0; 31 (9d6); effect/force-damage 31 (9d6); trigger/death, save/dex, area/emanation, effect/fire-damage, save/half-damage, effect/revive",
		"trait|trait/fire-aura|Fire Aura; emanation 5/0/0; 13 (3d8); area/emanation, effect/fire-damage",
		"trait|trait/legendary-resistance|Legendary Resistance (3/Day); limit/per-day, save/choose-success",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|flame-whip|Flame Whip; melee +14 reach 30 range 0/0 18 (3d6 + 8) force; effect/force-damage 18 (3d6 + 8), effect/fire-damage 17 (5d6), effect/pull, condition/prone, requirement/target-size",
		"action|lightning-blade|Lightning Blade; melee +14 reach 10 range 0/0 21 (3d8 + 8) force; effect/force-damage 21 (3d8 + 8), effect/lightning-damage 22 (4d10), effect/no-reactions",
		"bonus|teleport|Teleport; effect/teleport, target/self-or-willing-demon",
	},
	"Bandit": {
		"action|scimitar|Scimitar; melee +3 reach 5 range 0/0 4 (1d6 + 1) slashing; effect/slashing-damage",
		"action|light-crossbow|Light Crossbow; ranged +3 reach 0 range 80/320 5 (1d8 + 1) piercing; effect/piercing-damage",
	},
	"Bandit Captain": {
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack",
		"action|scimitar|Scimitar; melee +5 reach 5 range 0/0 6 (1d6 + 3) slashing; effect/slashing-damage 6 (1d6 + 3)",
		"action|pistol|Pistol; ranged +5 reach 0 range 30/90 8 (1d10 + 3) piercing; effect/piercing-damage 8 (1d10 + 3)",
		"reaction|parry|Parry; reaction/trigger, requirement/holding-weapon, effect/ac-bonus",
	},
	"Barbed Devil": {
		"trait|trait/barbed-hide|Barbed Hide; 5 (1d10); effect/piercing-damage, requirement/grappling-or-grappled",
		"trait|trait/diabolical-restoration|Diabolical Restoration; effect/revive",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack",
		"action|claws|Claws; melee +6 reach 5 range 0/0 10 (2d6 + 3) piercing; effect/piercing-damage 10 (2d6 + 3), condition/grappled, requirement/target-size",
		"action|tail|Tail; melee +6 reach 10 range 0/0 14 (2d10 + 3) slashing; effect/slashing-damage 14 (2d10 + 3)",
		"action|hurl-flame|Hurl Flame; ranged +5 reach 0 range 150/150 17 (5d6) fire; effect/fire-damage 17 (5d6), effect/ignite-object",
	},
	"Basilisk": {
		"action|bite|Bite; melee +5 reach 5 range 0/0 10 (2d6 + 3) piercing; effect/piercing-damage 10 (2d6 + 3), effect/poison-damage 7 (2d6)",
		"bonus|petrifying-gaze|Petrifying Gaze (Recharge 4–6); con DC 12; cone 30/0/0; save/con, area/cone, limit/recharge, effect/reflected-gaze, condition/restrained, save/repeat, condition/petrified, save/second-failure",
	},
	"Bat": {
		"action|bite|Bite; melee +4 reach 5 range 0/0 1 piercing; effect/piercing-damage [damage/fixed]",
	},
	"Bearded Devil": {
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|beard|Beard; melee +5 reach 5 range 0/0 7 (1d8 + 3) piercing; effect/piercing-damage 7 (1d8 + 3), condition/poisoned, effect/no-hit-point-recovery",
		"action|infernal-glaive|Infernal Glaive; melee +5 reach 10 range 0/0 8 (1d10 + 3) slashing; con DC 12; effect/slashing-damage 8 (1d10 + 3), requirement/no-existing-infernal-wound, save/con, effect/infernal-wound, effect/hit-point-loss 5 (1d10), effect/wound-closes, check/wis-medicine",
	},
	"Behir": {
		"action|multiattack|Multiattack; action/melee-attack, multiattack/additional-action",
		"action|bite|Bite; melee +10 reach 10 range 0/0 19 (2d12 + 6) piercing; effect/piercing-damage 19 (2d12 + 6), effect/lightning-damage 11 (2d10)",
		"action|constrict|Constrict; str DC 18; creature 0/0/5; save/str, target/creature, requirement/target-size, effect/bludgeoning-damage 28 (5d8 + 6), condition/grappled, condition/restrained",
		"action|lightning-breath|Lightning Breath (Recharge 5–6); dex DC 16; line 90/5/0; save/dex, area/line, limit/recharge, effect/lightning-damage 66 (12d10), save/half-damage",
		"bonus|swallow|Swallow; dex DC 18; creature 0/0/0; save/dex, target/creature, requirement/target-size, requirement/grappled-creature, limit/one-swallowed-creature, effect/swallow, condition/blinded, condition/restrained, effect/total-cover, effect/acid-damage 21 (6d6), effect/regurgitate, save/con-or-regurgitate, condition/prone, effect/escape-corpse",
	},
	"Berserker": {
		"trait|trait/bloodied-frenzy|Bloodied Frenzy; requirement/bloodied, attack/advantage, save/advantage",
		"action|greataxe|Greataxe; melee +5 reach 5 range 0/0 9 (1d12 + 3) slashing; effect/slashing-damage 9 (1d12 + 3)",
	},
}

// monsterDice154 also prints the subtracted modifier of p. 345's small beasts
// and a fixed hit, which the source prints as a bare number.
func monsterDice154(d MonsterHP) string {
	switch {
	case d.Dice == 0:
		return fmt.Sprint(d.Average)
	case d.Modifier < 0:
		return fmt.Sprintf("%d (%dd%d − %d)", d.Average, d.Dice, d.Sides, -d.Modifier)
	}
	return monsterDice151(d)
}

// monsterHeader154 prints typed fields in the source's stat-block notation.
func monsterHeader154(m Monster) string {
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
		"HP " + monsterDice154(m.HP),
	}
	speed := fmt.Sprintf("Speed %d ft.", m.Speeds["walk"])
	for _, mode := range monsterSorted151(m.Speeds) {
		if mode != "walk" {
			speed += fmt.Sprintf(", %s %d ft.", monsterTitle151(mode), m.Speeds[mode])
		}
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

func monsterEntry154(section, id string, x MonsterMechanic, attack *MonsterAttack, effects []MonsterMechanic) string {
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
		parts = append(parts, fmt.Sprintf("%s +%d reach %d range %d/%d %s %s", attack.Kind, attack.Bonus, attack.Reach, attack.Range, attack.LongRange, monsterDice154(attack.Damage), attack.DamageType))
	}
	if x.Save != nil {
		parts = append(parts, fmt.Sprintf("%s DC %d", x.Save.Ability, x.Save.DC))
	}
	if t := x.Target; t != nil {
		parts = append(parts, fmt.Sprintf("%s %d/%d/%d", t.Shape, t.Size, t.Width, t.Range))
	}
	if x.Dice != nil {
		parts = append(parts, monsterDice154(*x.Dice))
	}
	for _, d := range x.MoreDice {
		parts = append(parts, d.Key+" "+monsterDice154(d.Dice))
	}
	keys := slices.Clone(x.Requires)
	for _, e := range effects {
		key := e.Key
		if e.Dice != nil {
			key += " " + monsterDice154(*e.Dice)
		}
		if len(e.Requires) > 0 {
			key += " [" + strings.Join(e.Requires, ", ") + "]"
		}
		keys = append(keys, key)
	}
	return section + "|" + id + "|" + strings.Join(append(parts, strings.Join(keys, ", ")), "; ")
}

func monsterEntries154(m Monster) []string {
	out := []string{}
	for _, trait := range m.Traits {
		out = append(out, monsterEntry154("trait", trait.Key, trait, nil, nil))
	}
	for _, section := range []struct {
		label   string
		actions []MonsterAction
	}{{"action", m.Actions}, {"bonus", m.BonusActions}, {"reaction", m.Reactions}} {
		for _, a := range section.actions {
			out = append(out, monsterEntry154(section.label, a.ID, MonsterMechanic{Name: a.Name, Uses: a.Uses, Save: a.Save, Target: a.Target}, a.Attack, a.Effects))
		}
	}
	return out
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters154SourceValues(t *testing.T) {
	c := loadMonsters154(t)
	if len(c.Monsters) != 12 || len(monsters154Headers) != 12 || len(monsters154Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #154") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters154Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader154(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries154(m); !slices.Equal(got, monsters154Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters154Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.Hover || m.LairXP != 0 || m.LegendaryActions != nil || m.Legendary != nil {
				t.Fatal("source lists no vulnerabilities, hover, lair XP or legendary actions")
			}
		})
		entries += len(monsters154Entries[want.name])
	}
	// Bold-italic entry headings counted on PDF pp. 261–263 and 345.
	if entries != 41 {
		t.Fatalf("reviewed %d entries, want 41", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters154MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 154 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters154(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the three p. 345 blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("monsters animals animals monsters monsters monsters monsters monsters animals monsters monsters monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #154 cites the page each block starts on, in this order.
	names := []string{"Azer Sentinel", "Baboon", "Badger", "Balor", "Bandit", "Bandit Captain", "Barbed Devil", "Basilisk", "Bat", "Bearded Devil", "Behir", "Berserker"}
	pages := []int{261, 345, 345, 261, 261, 261, 262, 262, 345, 262, 263, 263}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic that holds for these twelve blocks (PDF pp. 254–256):
// XP and PB follow CR, Hit Dice follow size and Constitution, saves and skills
// add PB to the ability modifier, and attack rolls add PB to the listed ability.
func TestMonsters154DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"0": 2, "1/8": 2, "2": 2, "3": 2, "5": 3, "11": 4, "19": 6}
	skillAbility := map[string]string{"athletics": "str", "deception": "cha", "insight": "wis", "perception": "wis", "stealth": "dex"}
	hitDie := map[string]int{"tiny": 4, "small": 6, "medium": 8, "huge": 12}
	// Attacks that do not use the better of Strength and Dexterity.
	attackAbility := map[string]string{"Baboon/bite": "str", "Barbed Devil/hurl-flame": "cha"}
	for _, m := range loadMonsters154(t).Monsters {
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
			if extra := m.Initiative.Modifier - m.Abilities["dex"].Modifier; extra != 0 && extra != 2*p {
				t.Fatal("initiative is not Dexterity with optional twice PB")
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
				// A fixed hit and Hurl Flame's spell-like dice add no ability modifier.
				damage := mod
				if fixedMonsterDamage(a.Attack.Damage) || ability == "cha" {
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
func TestMonsters154DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters154(t).Monsters {
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
func TestMonsters154RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters154(t).Monsters {
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

func TestMonsters154MalformedFixtures(t *testing.T) {
	// Indexes: 1 Baboon, 2 Badger, 3 Balor, 6 Barbed Devil, 7 Basilisk,
	// 9 Bearded Devil, 10 Behir.
	tests := []struct {
		name   string
		mutate func(*MonsterCatalogue)
	}{
		{"duplicate ID", func(c *MonsterCatalogue) { c.Monsters[1].ID = c.Monsters[2].ID }},
		{"animal block in another family", func(c *MonsterCatalogue) { c.Monsters[1].ID = "srd-5.2.1:content/spells/baboon" }},
		{"missing languages list", func(c *MonsterCatalogue) { c.Monsters[1].Languages = nil }},
		{"subtracted damage average", func(c *MonsterCatalogue) { c.Monsters[1].Actions[0].Attack.Damage.Average = 2 }},
		{"fixed damage with sides", func(c *MonsterCatalogue) { c.Monsters[2].Actions[0].Attack.Damage.Sides = 4 }},
		{"fixed damage with modifier", func(c *MonsterCatalogue) { c.Monsters[2].Actions[0].Attack.Damage.Modifier = 1 }},
		{"fixed damage of zero", func(c *MonsterCatalogue) { c.Monsters[2].Actions[0].Attack.Damage.Average = 0 }},
		{"fixed damage without mechanic", func(c *MonsterCatalogue) {
			m := &c.Monsters[2]
			m.Actions[0].Effects[0].Requires = nil
			m.RequiredMechanics = slices.DeleteFunc(m.RequiredMechanics, func(key string) bool { return key == "damage/fixed" })
		}},
		{"fixed effect dice", func(c *MonsterCatalogue) { c.Monsters[2].Actions[0].Effects[0].Dice = &MonsterHP{Average: 1} }},
		{"sense note without sense", func(c *MonsterCatalogue) {
			c.Monsters[6].Senses.Notes = map[string]string{"blindsight": "unimpeded by magical Darkness"}
		}},
		{"blank sense note", func(c *MonsterCatalogue) { c.Monsters[6].Senses.Notes["darkvision"] = " " }},
		{"padded sense note", func(c *MonsterCatalogue) { c.Monsters[6].Senses.Notes["darkvision"] += " " }},
		{"unlisted sense note", func(c *MonsterCatalogue) { c.Monsters[6].Senses.Notes["darkvision"] = "impeded by magical Darkness" }},
		{"empty sense notes", func(c *MonsterCatalogue) { c.Monsters[6].Senses.Notes = map[string]string{} }},
		{"single range below range", func(c *MonsterCatalogue) { c.Monsters[6].Actions[3].Attack.LongRange = 90 }},
		{"ranged reach", func(c *MonsterCatalogue) { c.Monsters[6].Actions[3].Attack.Reach = 5 }},
		{"further dice without first", func(c *MonsterCatalogue) { c.Monsters[3].Traits[0].Dice = nil }},
		{"further dice average", func(c *MonsterCatalogue) { c.Monsters[3].Traits[0].MoreDice[0].Dice.Average = 32 }},
		{"further dice off entry text", func(c *MonsterCatalogue) {
			c.Monsters[3].Traits[0].MoreDice[0].Dice = MonsterHP{Average: 35, Dice: 10, Sides: 6}
		}},
		{"further dice unlisted key", func(c *MonsterCatalogue) { c.Monsters[3].Traits[0].MoreDice[0].Key = "effect/future-damage" }},
		{"further dice invalid key", func(c *MonsterCatalogue) { c.Monsters[3].Traits[0].MoreDice[0].Key = "Force" }},
		{"duplicate further dice key", func(c *MonsterCatalogue) {
			t := &c.Monsters[3].Traits[0]
			t.MoreDice = append(t.MoreDice, t.MoreDice[0])
		}},
		{"empty further dice", func(c *MonsterCatalogue) { c.Monsters[3].Traits[1].MoreDice = []MonsterKeyedDice{} }},
		{"lair daily uses", func(c *MonsterCatalogue) { c.Monsters[3].Traits[2].Uses.LairPerDay = 2 }},
		{"XP off CR", func(c *MonsterCatalogue) { c.Monsters[3].XP = 22001 }},
		{"save off modifier", func(c *MonsterCatalogue) { c.Monsters[3].Abilities["con"] = MonsterAbility{22, 6, 11} }},
		{"modifier", func(c *MonsterCatalogue) { c.Monsters[1].Abilities["str"] = MonsterAbility{8, 0, 0} }},
		{"HP average", func(c *MonsterCatalogue) { c.Monsters[8].HP.Average = 2 }},
		{"empty bonus actions", func(c *MonsterCatalogue) { c.Monsters[1].BonusActions = []MonsterAction{} }},
		{"empty reactions", func(c *MonsterCatalogue) { c.Monsters[1].Reactions = []MonsterAction{} }},
		{"duplicate action across sections", func(c *MonsterCatalogue) { c.Monsters[10].BonusActions[0].ID = "bite" }},
		{"recharge", func(c *MonsterCatalogue) { c.Monsters[7].BonusActions[0].Uses.Recharge = 1 }},
		{"save DC off entry text", func(c *MonsterCatalogue) { c.Monsters[9].Actions[2].Save.DC = 13 }},
		{"save ability off entry text", func(c *MonsterCatalogue) { c.Monsters[10].BonusActions[0].Save.Ability = "str" }},
		{"creature size", func(c *MonsterCatalogue) { c.Monsters[10].Actions[2].Target.Size = 5 }},
		{"cone width", func(c *MonsterCatalogue) { c.Monsters[7].BonusActions[0].Target.Width = 5 }},
		{"trait dice", func(c *MonsterCatalogue) { c.Monsters[6].Traits[0].Dice.Sides = 9 }},
		{"effect dice", func(c *MonsterCatalogue) { c.Monsters[9].Actions[2].Effects[4].Dice.Average = 6 }},
		{"unlisted reaction effect", func(c *MonsterCatalogue) { c.Monsters[5].Reactions[0].Effects[0].Key = "reaction/future" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters154(t)
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
	raw, err := os.ReadFile("testdata/monsters_154.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, edit := range [][2]string{
		{`"dice": 0,`, ``},
		{`"sides": 0,`, `"sides": null,`},
		{`"long_range": 150,`, ``},
		{`"long_range": 90,`, `"long_range": null,`},
		{`"recharge": 4`, `"recharge": null`},
		{`"per_day": 3,`, ``},
		{`"dc": 20`, `"dc": 20, "dc": 21`},
		{`"width": 5,`, ``},
		{`"notes": {`, `"notes": {"truesight": "see invisible", `},
		{`"more_dice": [`, `"extra_dice": [`},
		{`"shape": "cone",`, `"shape": "cone", "height": 5,`},
	} {
		if !bytes.Contains(raw, []byte(edit[0])) {
			t.Fatalf("fixture lost %s", edit[0])
		}
		if _, err := DecodeMonsters(bytes.NewReader(bytes.Replace(raw, []byte(edit[0]), []byte(edit[1]), 1))); err == nil {
			t.Fatalf("accepted %s -> %s", edit[0], edit[1])
		}
	}
}

// The batch composes with the #143 Bandit catalogue and stat-block batch 01
// without entering the embedded runtime catalogue. Its Bandit record is the
// revision #143 pinned, so composition keeps one copy and the revision stays.
func TestMonsters154ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	first, batch := loadMonsters151(t), loadMonsters154(t)
	if len(base.Monsters) != 1 || base.Monsters[0].ID != BanditID {
		t.Fatal("batch data must not enter the embedded runtime catalogue")
	}
	if !reflect.DeepEqual(batch.Monsters[4], base.Monsters[0]) {
		t.Fatal("the batch must repeat the pinned Bandit revision unchanged")
	}
	for _, c := range []MonsterCatalogue{base, first} {
		a, b := c.Source, batch.Source
		a.Changes, b.Changes = "", ""
		if a != b {
			t.Fatal("batches must share one source pin")
		}
	}
	rest := slices.Concat(batch.Monsters[:4], batch.Monsters[5:])
	combined := MonsterCatalogue{SchemaVersion: 1, Source: batch.Source, Monsters: slices.Concat(base.Monsters, first.Monsters, rest)}
	if err := combined.Validate(); err != nil || len(combined.Monsters) != 24 {
		t.Fatal("batch must compose with the earlier catalogues:", err)
	}
	combined.Monsters = append(combined.Monsters, batch.Monsters[4])
	if combined.Validate() == nil {
		t.Fatal("composition accepted a second Bandit record")
	}
	_, version, err := BanditRecord()
	if err != nil || version != "bandit-v1:1169c77601cf83a9b96985fc761a6b8be9c901dc1a86e6a9758dcc399707c3d4" {
		t.Fatal("contract extension changed the pinned Bandit revision:", version, err)
	}
	for _, id := range []string{"burning-hammer", "pistol", "parry", "teleport"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
