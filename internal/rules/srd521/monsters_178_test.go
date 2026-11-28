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

func loadMonsters178(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_178.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #178 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Young Green Dragon’s and the Zombie’s blocks run onto the next page.
var monsters178Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Young Copper Dragon", [2]int{276, 276}, "Large Dragon (Metallic), Chaotic Good|AC 17 Initiative +4 (14)|HP 119 (14d10 + 42)|Speed 40 ft., Climb 40 ft., Fly 80 ft.|Str 19 +4 +4 Dex 12 +1 +4 Con 17 +3 +3|Int 16 +3 +3 Wis 13 +1 +4 Cha 15 +2 +2|Skills Deception +5, Perception +7, Stealth +4|Immunities Acid|Senses Blindsight 30 ft., Darkvision 120 ft.; Passive Perception 17|Languages Common, Draconic|CR 7 (XP 2,900; PB +3)"},
	{"Young Gold Dragon", [2]int{291, 291}, "Large Dragon (Metallic), Lawful Good|AC 18 Initiative +6 (16)|HP 178 (17d10 + 85)|Speed 40 ft., Fly 80 ft., Swim 40 ft.|Str 23 +6 +6 Dex 14 +2 +6 Con 21 +5 +5|Int 16 +3 +3 Wis 13 +1 +5 Cha 20 +5 +5|Skills Insight +5, Perception +9, Persuasion +9, Stealth +6|Immunities Fire|Senses Blindsight 30 ft., Darkvision 120 ft.; Passive Perception 19|Languages Common, Draconic|CR 10 (XP 5,900; PB +4)"},
	{"Young Green Dragon", [2]int{293, 294}, "Large Dragon (Chromatic), Lawful Evil|AC 18 Initiative +4 (14)|HP 136 (16d10 + 48)|Speed 40 ft., Fly 80 ft., Swim 40 ft.|Str 19 +4 +4 Dex 12 +1 +4 Con 17 +3 +3|Int 16 +3 +3 Wis 13 +1 +4 Cha 15 +2 +2|Skills Deception +5, Perception +7, Stealth +4|Immunities Poison; Poisoned|Senses Blindsight 30 ft., Darkvision 120 ft.; Passive Perception 17|Languages Common, Draconic|CR 8 (XP 3,900; PB +3)"},
	{"Young Red Dragon", [2]int{318, 318}, "Large Dragon (Chromatic), Chaotic Evil|AC 18 Initiative +4 (14)|HP 178 (17d10 + 85)|Speed 40 ft., Climb 40 ft., Fly 80 ft.|Str 23 +6 +6 Dex 10 +0 +4 Con 21 +5 +5|Int 14 +2 +2 Wis 11 +0 +4 Cha 19 +4 +4|Skills Perception +8, Stealth +4|Immunities Fire|Senses Blindsight 30 ft., Darkvision 120 ft.; Passive Perception 18|Languages Common, Draconic|CR 10 (XP 5,900; PB +4)"},
	{"Young Silver Dragon", [2]int{324, 324}, "Large Dragon (Metallic), Lawful Good|AC 18 Initiative +4 (14)|HP 168 (16d10 + 80)|Speed 40 ft., Fly 80 ft.|Str 23 +6 +6 Dex 10 +0 +4 Con 21 +5 +5|Int 14 +2 +2 Wis 11 +0 +4 Cha 19 +4 +4|Skills History +6, Perception +8, Stealth +4|Immunities Cold|Senses Blindsight 30 ft., Darkvision 120 ft.; Passive Perception 18|Languages Common, Draconic|CR 9 (XP 5,000; PB +4)"},
	{"Young White Dragon", [2]int{340, 340}, "Large Dragon (Chromatic), Chaotic Evil|AC 17 Initiative +3 (13)|HP 123 (13d10 + 52)|Speed 40 ft., Burrow 20 ft., Fly 80 ft., Swim 40 ft.|Str 18 +4 +4 Dex 10 +0 +3 Con 18 +4 +4|Int 6 −2 2 Wis 11 +0 +3 Cha 12 +1 +1|Skills Perception +6, Stealth +3|Immunities Cold|Senses Blindsight 30 ft., Darkvision 120 ft.; Passive Perception 16|Languages Common, Draconic|CR 6 (2,300 XP; PB +3)"},
	{"Zombie", [2]int{343, 344}, "Medium Undead, Neutral Evil|AC 8 Initiative −2 (8)|HP 15 (2d8 + 6)|Speed 20 ft.|Str 13 +1 +1 Dex 6 −2 −2 Con 16 +3 +3|Int 3 −4 −4 Wis 6 −2 +0 Cha 5 −3 −3|Immunities Poison; Exhaustion, Poisoned|Senses Darkvision 60 ft.; Passive Perception 8|Languages Understands Common plus one other language but can’t speak|CR 1/4 (XP 50; PB +2)"},
}

// Typed values of every trait and action entry, reviewed against the same pages:
// limited uses, attack bonus/reach/range/damage, save ability and DC, target
// shape size/width/range, dice and dependency keys. A ranged attack printed with
// one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text.
var monsters178Entries = map[string][]string{
	"Young Copper Dragon": {
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, action/saving-throw",
		"action|rend|Rend; melee +7 reach 10 range 0/0 15 (2d10 + 4) slashing; effect/slashing-damage 15 (2d10 + 4)",
		"action|acid-breath|Acid Breath (Recharge 5–6); dex DC 14; line 40/5/0; save/dex, area/line, limit/recharge, effect/acid-damage 40 (9d8), save/half-damage",
		"action|slowing-breath|Slowing Breath; con DC 14; cone 30/0/0; save/con, area/cone, effect/no-reactions, effect/speed-halved, effect/action-or-bonus-action, effect/lasts-until-end-of-next-turn",
	},
	"Young Gold Dragon": {
		"trait|trait/amphibious|Amphibious; ",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, action/saving-throw",
		"action|rend|Rend; melee +10 reach 10 range 0/0 17 (2d10 + 6) slashing; effect/slashing-damage 17 (2d10 + 6)",
		"action|fire-breath|Fire Breath (Recharge 5–6); dex DC 17; cone 30/0/0; save/dex, area/cone, limit/recharge, effect/fire-damage 55 (10d10), save/half-damage",
		"action|weakening-breath|Weakening Breath; str DC 17; cone 30/0/0; save/str, area/cone, requirement/unaffected-target, effect/strength-test-disadvantage, effect/damage-roll-penalty 3 (1d6), save/repeat, save/automatic-success",
	},
	"Young Green Dragon": {
		"trait|trait/amphibious|Amphibious; ",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +7 reach 10 range 0/0 11 (2d6 + 4) slashing; effect/slashing-damage 11 (2d6 + 4), effect/poison-damage 7 (2d6)",
		"action|poison-breath|Poison Breath (Recharge 5–6); con DC 14; cone 30/0/0; save/con, area/cone, limit/recharge, effect/poison-damage 42 (12d6), save/half-damage",
	},
	"Young Red Dragon": {
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +10 reach 10 range 0/0 13 (2d6 + 6) slashing; effect/slashing-damage 13 (2d6 + 6), effect/fire-damage 3 (1d6)",
		"action|fire-breath|Fire Breath (Recharge 5–6); dex DC 17; cone 30/0/0; save/dex, area/cone, limit/recharge, effect/fire-damage 56 (16d6), save/half-damage",
	},
	"Young Silver Dragon": {
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, action/saving-throw",
		"action|rend|Rend; melee +10 reach 10 range 0/0 15 (2d8 + 6) slashing; effect/slashing-damage 15 (2d8 + 6)",
		"action|cold-breath|Cold Breath (Recharge 5–6); con DC 17; cone 30/0/0; save/con, area/cone, limit/recharge, effect/cold-damage 49 (11d8), save/half-damage",
		"action|paralyzing-breath|Paralyzing Breath; con DC 17; cone 30/0/0; save/con, area/cone, condition/incapacitated, save/repeat, condition/paralyzed, save/second-failure, save/repeat, save/automatic-success",
	},
	"Young White Dragon": {
		"trait|trait/ice-walk|Ice Walk; movement/ice-walk, movement/ignore-difficult-terrain",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|rend|Rend; melee +7 reach 10 range 0/0 9 (2d4 + 4) slashing; effect/slashing-damage 9 (2d4 + 4), effect/cold-damage 2 (1d4)",
		"action|cold-breath|Cold Breath (Recharge 5–6); con DC 15; cone 30/0/0; save/con, area/cone, limit/recharge, effect/cold-damage 40 (9d8), save/half-damage",
	},
	"Zombie": {
		"trait|trait/undead-fortitude|Undead Fortitude; requirement/reduced-to-zero-hit-points, save/con, effect/damage-based-save-dc, requirement/not-radiant-or-critical-hit, effect/drop-to-one-hit-point",
		"action|slam|Slam; melee +3 reach 5 range 0/0 5 (1d8 + 1) bludgeoning; effect/bludgeoning-damage 5 (1d8 + 1)",
	},
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters178SourceValues(t *testing.T) {
	c := loadMonsters178(t)
	if len(c.Monsters) != 7 || len(monsters178Headers) != 7 || len(monsters178Entries) != 7 {
		t.Fatal("exact seven stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #178") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters178Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader176(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if m.Name == "Young White Dragon" && m.Abilities["int"] != (MonsterAbility{6, -2, -2}) {
				t.Fatalf("Intelligence must keep the derived save: %+v", m.Abilities["int"])
			}
			if got := monsterEntries164(m); !slices.Equal(got, monsters178Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters178Entries[want.name])
			}
			if m.Vulnerabilities != nil || m.VulnerabilityNote != "" || m.ResistanceNote != "" || m.TelepathyNote != "" || m.SwarmOf != "" || m.SpeedNote != "" || m.ConditionalSpeedNote != "" || m.ConditionImmunityNote != "" || m.Telepathy != 0 || m.LairXP != 0 || m.Legendary != nil || m.LegendaryActions != nil || m.Senses.Notes != nil {
				t.Fatal("source lists no vulnerabilities, vulnerability note, resistance note, telepathy note, swarm line, speed note, conditional speed note, condition immunity note, telepathy, lair XP, legendary actions, sense notes")
			}
		})
		for _, line := range monsters178Entries[want.name] {
			if !strings.HasPrefix(line, "legendary-uses|") {
				entries++
			}
		}
	}
	// Bold-italic entry headings counted on PDF pp. 276, 291, 293–294, 318, 324, 340 and 343–344.
	if entries != 26 {
		t.Fatalf("reviewed %d entries, want 26", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters178MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 178 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters178(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 7 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Animals appendix blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("monsters monsters monsters monsters monsters monsters monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #178 cites the page each block starts on, in this order.
	names := []string{"Young Copper Dragon", "Young Gold Dragon", "Young Green Dragon", "Young Red Dragon", "Young Silver Dragon", "Young White Dragon", "Zombie"}
	pages := []int{276, 291, 293, 318, 324, 340, 343}
	for i, r := range want {
		if r.name != names[i] || r.page != pages[i] {
			t.Fatalf("%s: issue cites %s on p.%d", r.name, names[i], pages[i])
		}
	}
}

// Source arithmetic for these seven blocks; see monsterConsistency170.
func TestMonsters178DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"7": 3, "10": 4, "8": 3, "9": 4, "6": 3, "1/4": 2}
	attacks := map[string]attackAbility170{}
	for _, m := range loadMonsters178(t).Monsters {
		t.Run(m.Name, func(t *testing.T) {
			monsterConsistency170(t, m, pb, attacks)
		})
	}
}

// Every listed key is used and required, every effect quotes its own entry, and
// every dice expression is typed. Removing any one key invalidates the record.
func TestMonsters178DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters178(t).Monsters {
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
func TestMonsters178RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters178(t).Monsters {
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

func TestMonsters178MalformedFixtures(t *testing.T) {
	tests := []struct {
		name, record string
		mutate       func(*Monster)
	}{
		{"missing languages list", "Zombie", func(m *Monster) {
			m.Languages = nil
		}},
		{"animal or monster block in another family", "Zombie", func(m *Monster) {
			m.ID = strings.Replace(m.ID, "/monsters/", "/spells/", 1)
		}},
		{"XP off CR", "Zombie", func(m *Monster) {
			m.XP++
		}},
		{"printed Intelligence save without its sign", "Young White Dragon", func(m *Monster) {
			m.Abilities["int"] = MonsterAbility{m.Abilities["int"].Score, m.Abilities["int"].Modifier, 2}
		}},
		{"modifier", "Young Gold Dragon", func(m *Monster) {
			m.Abilities["str"] = MonsterAbility{m.Abilities["str"].Score, m.Abilities["str"].Modifier + 1, m.Abilities["str"].Modifier + 1}
		}},
		{"HP average", "Young Gold Dragon", func(m *Monster) {
			m.HP.Average++
		}},
		{"recharge", "Young Red Dragon", func(m *Monster) {
			m.Actions[2].Uses.Recharge = 7
		}},
		{"cone width", "Young Silver Dragon", func(m *Monster) {
			m.Actions[2].Target.Width = 5
		}},
		{"line without width", "Young Copper Dragon", func(m *Monster) {
			m.Actions[2].Target.Width = 0
		}},
		{"save DC off entry text", "Young Gold Dragon", func(m *Monster) {
			m.Actions[3].Save.DC = 18
		}},
		{"empty bonus actions", "Zombie", func(m *Monster) {
			m.BonusActions = []MonsterAction{}
		}},
		{"empty reactions", "Zombie", func(m *Monster) {
			m.Reactions = []MonsterAction{}
		}},
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
		c := loadMonsters178(t)
		c.Monsters[1].ID = c.Monsters[2].ID
		rejected(t, c, true)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters178(t)
			tt.mutate(record(t, &c, tt.record))
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			rejected(t, c, !strings.HasPrefix(tt.name, "empty "))
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_178.json")
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
		{`"dc":14`, `"dc":14,"dc":22`},
		{`"width":5,`, ``},
		{`"range":0}`, `"range":null}`},
		{`"darkvision":120`, `"darkvision":null`},
		{`"shape":"cone",`, `"shape":"cone","height":5,`},
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
// to 27 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters178ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters178(t)
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
	for _, load := range []func(*testing.T) MonsterCatalogue{loadMonsters158, loadMonsters159, loadMonsters160, loadMonsters161, loadMonsters162, loadMonsters163, loadMonsters164, loadMonsters165, loadMonsters166, loadMonsters167, loadMonsters168, loadMonsters169, loadMonsters170, loadMonsters171, loadMonsters172, loadMonsters173, loadMonsters174, loadMonsters175, loadMonsters176, loadMonsters177} {
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
	if err := combined.Validate(); err != nil || len(combined.Monsters) != len(earlier)+7 {
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
	for _, id := range []string{"rend", "slam"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
