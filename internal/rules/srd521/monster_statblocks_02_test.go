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

func monsterBatch02(t *testing.T) MonsterBatch02 {
	t.Helper()
	b, err := LoadMonsterBatch02()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// batch02Want is typed by hand from the cited PDF page, not from the fixture.
// abilities are score and save in Str, Dex, Con, Int, Wis, Cha order; actions
// are rendered by batch02Line so every typed number has a literal.
type batch02Want struct {
	id, name       string
	pages          []int
	size, typ, tag string
	alignment      string
	ac             int
	hp             MonsterHP
	initiative     MonsterInitiative
	speeds         map[string]int
	abilities      [6][2]int
	skills         map[string]int
	immunities     string
	special        map[string]int
	passive        int
	languages      []string
	cr             string
	xp, lairXP, pb int
	legendary      string
	traits         []string
	actions        []string
}

var (
	dragonSenses    = map[string]int{"blindsight": 60, "darkvision": 120}
	dragonLanguages = []string{"Common", "Draconic"}
	dragonLegendary = "3/4 LR 4/5"
)

var monsterBatch02Want = []batch02Want{
	{id: "animals/allosaurus", name: "Allosaurus", pages: []int{344, 344}, size: "large", typ: "beast", tag: "dinosaur", alignment: "Unaligned",
		ac: 13, hp: MonsterHP{51, 6, 10, 18}, initiative: MonsterInitiative{1, 11}, speeds: map[string]int{"walk": 60},
		abilities: [6][2]int{{19, 4}, {13, 1}, {17, 3}, {2, -4}, {12, 1}, {5, -3}}, skills: map[string]int{"perception": 5},
		special: map[string]int{}, passive: 15, languages: []string{}, cr: "2", xp: 450, pb: 2, traits: []string{},
		actions: []string{
			"bite/action/melee-attack atk+6 r5 15(2d10+4)piercing@hit",
			"claws/action/melee-attack atk+6 r5 8(1d8+4)slashing@hit c:prone x1bite",
		}},
	{id: "monsters/ancient-black-dragon", name: "Ancient Black Dragon", pages: []int{265, 265}, size: "gargantuan", typ: "dragon", tag: "chromatic", alignment: "Chaotic Evil",
		ac: 22, hp: MonsterHP{367, 21, 20, 147}, initiative: MonsterInitiative{16, 26}, speeds: map[string]int{"walk": 40, "fly": 80, "swim": 40},
		abilities: [6][2]int{{27, 8}, {14, 9}, {25, 7}, {16, 3}, {15, 9}, {22, 6}}, skills: map[string]int{"perception": 16, "stealth": 9},
		immunities: "acid;", special: dragonSenses, passive: 26, languages: dragonLanguages, cr: "21", xp: 33000, lairXP: 41000, pb: 7,
		legendary: dragonLegendary, traits: []string{"Amphibious", "Legendary Resistance"},
		actions: []string{
			"multiattack/action/multiattack x3rend rep1:spellcasting/Acid Arrow@4",
			"rend/action/melee-attack atk+15 r15 17(2d8+8)slashing@hit 9(2d8+0)acid@hit",
			"acid-breath/action/acid-breath rch5-6 dex22½ line90x10 67(15d8+0)acid@failure",
			"spellcasting/action/spellcasting sc:cha21+13[Acid Arrow@4,Detect Magic@0,Fear@0;Create Undead@0,Speak with Dead@0,Vitriolic Sphere@5]",
			"cloud-of-insects/legendary-action/cloud-of-insects dex21 tgt120 33(6d10+0)poison@failure 1/cycle",
			"frightful-presence/legendary-action/frightful-presence 1/cycle cast:Fear@0",
			"pounce/legendary-action/pounce mv:half-speed x1rend",
		}},
	{id: "monsters/ancient-blue-dragon", name: "Ancient Blue Dragon", pages: []int{267, 267}, size: "gargantuan", typ: "dragon", tag: "chromatic", alignment: "Lawful Evil",
		ac: 22, hp: MonsterHP{481, 26, 20, 208}, initiative: MonsterInitiative{14, 24}, speeds: map[string]int{"walk": 40, "burrow": 40, "fly": 80},
		abilities: [6][2]int{{29, 9}, {10, 7}, {27, 8}, {18, 4}, {17, 10}, {25, 7}}, skills: map[string]int{"perception": 17, "stealth": 7},
		immunities: "lightning;", special: dragonSenses, passive: 27, languages: dragonLanguages, cr: "23", xp: 50000, lairXP: 62000, pb: 7,
		legendary: dragonLegendary, traits: []string{"Legendary Resistance"},
		actions: []string{
			"multiattack/action/multiattack x3rend rep1:spellcasting/Shatter@3",
			"rend/action/melee-attack atk+16 r15 18(2d8+9)slashing@hit 11(2d10+0)lightning@hit",
			"lightning-breath/action/lightning-breath rch5-6 dex23½ line120x10 88(16d10+0)lightning@failure",
			"spellcasting/action/spellcasting sc:cha22[Detect Magic@0,Invisibility@0,Mage Hand@0,Shatter@3;Scrying@0,Sending@0]",
			"cloaked-flight/legendary-action/cloaked-flight 1/cycle mv:half-fly-speed cast:Invisibility@0",
			"sonic-boom/legendary-action/sonic-boom 1/cycle cast:Shatter@3",
			"tail-swipe/legendary-action/tail-swipe x1rend",
		}},
	{id: "monsters/ancient-brass-dragon", name: "Ancient Brass Dragon", pages: []int{269, 269}, size: "gargantuan", typ: "dragon", tag: "metallic", alignment: "Chaotic Good",
		ac: 20, hp: MonsterHP{332, 19, 20, 133}, initiative: MonsterInitiative{12, 22}, speeds: map[string]int{"walk": 40, "burrow": 40, "fly": 80},
		abilities:  [6][2]int{{27, 8}, {10, 6}, {25, 7}, {16, 3}, {15, 8}, {22, 6}},
		skills:     map[string]int{"history": 9, "perception": 14, "persuasion": 12, "stealth": 6},
		immunities: "fire;", special: dragonSenses, passive: 24, languages: dragonLanguages, cr: "20", xp: 25000, lairXP: 33000, pb: 6,
		legendary: dragonLegendary, traits: []string{"Legendary Resistance"},
		actions: []string{
			"multiattack/action/multiattack x3rend rep1:sleep-breath|spellcasting/Scorching Ray@3",
			"rend/action/melee-attack atk+14 r15 19(2d10+8)slashing@hit 7(2d6+0)fire@hit",
			"fire-breath/action/fire-breath rch5-6 dex21½ line90x5 58(13d8+0)fire@failure",
			"sleep-breath/action/sleep-breath con21 cone90 c:incapacitated,unconscious min10",
			"spellcasting/action/spellcasting sc:cha20[Detect Magic@0,Minor Illusion@0,Scorching Ray@3,Shapechange@0*,Speak with Animals@0;Control Weather@0,Detect Thoughts@0]",
			"blazing-light/legendary-action/blazing-light cast:Scorching Ray@3",
			"pounce/legendary-action/pounce mv:half-speed x1rend",
			"scorching-sands/legendary-action/scorching-sands dex20 tgt120 36(8d8+0)fire@failure 1/cycle",
		}},
	{id: "monsters/ancient-bronze-dragon", name: "Ancient Bronze Dragon", pages: []int{271, 271}, size: "gargantuan", typ: "dragon", tag: "metallic", alignment: "Lawful Good",
		ac: 22, hp: MonsterHP{444, 24, 20, 192}, initiative: MonsterInitiative{14, 24}, speeds: map[string]int{"walk": 40, "fly": 80, "swim": 40},
		abilities: [6][2]int{{29, 9}, {10, 7}, {27, 8}, {18, 4}, {17, 10}, {25, 7}}, skills: map[string]int{"insight": 10, "perception": 17, "stealth": 7},
		immunities: "lightning;", special: dragonSenses, passive: 27, languages: dragonLanguages, cr: "22", xp: 41000, lairXP: 50000, pb: 7,
		legendary: dragonLegendary, traits: []string{"Amphibious", "Legendary Resistance"},
		actions: []string{
			"multiattack/action/multiattack x3rend rep1:repulsion-breath|spellcasting/Guiding Bolt@2",
			"rend/action/melee-attack atk+16 r15 18(2d8+9)slashing@hit 9(2d8+0)lightning@hit",
			"lightning-breath/action/lightning-breath rch5-6 dex23½ line120x10 82(15d10+0)lightning@failure",
			"repulsion-breath/action/repulsion-breath str23 cone30 c:prone push60",
			"spellcasting/action/spellcasting sc:cha22+14[Detect Magic@0,Guiding Bolt@2,Shapechange@0*,Speak with Animals@0,Thaumaturgy@0;Detect Thoughts@0,Control Water@0,Scrying@0,Water Breathing@0]",
			"guiding-light/legendary-action/guiding-light cast:Guiding Bolt@2",
			"pounce/legendary-action/pounce mv:half-speed x1rend",
			"thunderclap/legendary-action/thunderclap con22 sphere20@120 13(3d8+0)thunder@failure c:deafened",
		}},
	{id: "monsters/ancient-copper-dragon", name: "Ancient Copper Dragon", pages: []int{277, 277}, size: "gargantuan", typ: "dragon", tag: "metallic", alignment: "Chaotic Good",
		ac: 21, hp: MonsterHP{367, 21, 20, 147}, initiative: MonsterInitiative{15, 25}, speeds: map[string]int{"walk": 40, "climb": 40, "fly": 80},
		abilities: [6][2]int{{27, 8}, {12, 8}, {25, 7}, {20, 5}, {17, 10}, {22, 6}}, skills: map[string]int{"deception": 13, "perception": 17, "stealth": 8},
		immunities: "acid;", special: dragonSenses, passive: 27, languages: dragonLanguages, cr: "21", xp: 33000, lairXP: 41000, pb: 7,
		legendary: dragonLegendary, traits: []string{"Legendary Resistance"},
		actions: []string{
			"multiattack/action/multiattack x3rend rep1:slowing-breath|spellcasting/Mind Spike@5",
			"rend/action/melee-attack atk+15 r15 19(2d10+8)slashing@hit 9(2d8+0)acid@hit",
			"acid-breath/action/acid-breath rch5-6 dex22½ line90x10 63(14d8+0)acid@failure",
			"slowing-breath/action/slowing-breath con22 cone90",
			"spellcasting/action/spellcasting sc:cha21[Detect Magic@0,Mind Spike@5,Minor Illusion@0,Shapechange@0*;Greater Restoration@0,Major Image@0,Project Image@0]",
			"giggling-magic/legendary-action/giggling-magic cha21 tgt120 31(9d6+0)psychic@failure pen:d20-test:1d8 1/cycle",
			"mind-jolt/legendary-action/mind-jolt 1/cycle cast:Mind Spike@5",
			"pounce/legendary-action/pounce mv:half-speed x1rend",
		}},
	{id: "monsters/ancient-gold-dragon", name: "Ancient Gold Dragon", pages: []int{292, 292}, size: "gargantuan", typ: "dragon", tag: "metallic", alignment: "Lawful Good",
		ac: 22, hp: MonsterHP{546, 28, 20, 252}, initiative: MonsterInitiative{16, 26}, speeds: map[string]int{"walk": 40, "fly": 80, "swim": 40},
		abilities:  [6][2]int{{30, 10}, {14, 9}, {29, 9}, {18, 4}, {17, 10}, {28, 9}},
		skills:     map[string]int{"insight": 10, "perception": 17, "persuasion": 16, "stealth": 9},
		immunities: "fire;", special: dragonSenses, passive: 27, languages: dragonLanguages, cr: "24", xp: 62000, lairXP: 75000, pb: 7,
		legendary: dragonLegendary, traits: []string{"Amphibious", "Legendary Resistance"},
		actions: []string{
			"multiattack/action/multiattack x3rend rep1:spellcasting/Guiding Bolt@4|weakening-breath",
			"rend/action/melee-attack atk+17 r15 19(2d8+10)slashing@hit 9(2d8+0)fire@hit",
			"fire-breath/action/fire-breath rch5-6 dex24½ cone90 71(13d10+0)fire@failure",
			"spellcasting/action/spellcasting sc:cha24+16[Detect Magic@0,Guiding Bolt@4,Shapechange@0*;Flame Strike@6,Word of Recall@0,Zone of Truth@0]",
			"weakening-breath/action/weakening-breath str24 cone90 min1 pen:damage-roll:5(1d10)",
			"banish/legendary-action/banish cha24 tgt120 24(7d6+0)force@failure c:incapacitated 1/cycle",
			"guiding-light/legendary-action/guiding-light cast:Guiding Bolt@4",
			"pounce/legendary-action/pounce mv:half-speed x1rend",
		}},
	{id: "monsters/ancient-green-dragon", name: "Ancient Green Dragon", pages: []int{294, 295}, size: "gargantuan", typ: "dragon", tag: "chromatic", alignment: "Lawful Evil",
		ac: 21, hp: MonsterHP{402, 23, 20, 161}, initiative: MonsterInitiative{15, 25}, speeds: map[string]int{"walk": 40, "fly": 80, "swim": 40},
		abilities:  [6][2]int{{27, 8}, {12, 8}, {25, 7}, {20, 5}, {17, 10}, {22, 6}},
		skills:     map[string]int{"deception": 13, "perception": 17, "persuasion": 13, "stealth": 8},
		immunities: "poison;poisoned", special: dragonSenses, passive: 27, languages: dragonLanguages, cr: "22", xp: 41000, lairXP: 50000, pb: 7,
		legendary: dragonLegendary, traits: []string{"Amphibious", "Legendary Resistance"},
		actions: []string{
			"multiattack/action/multiattack x3rend rep1:spellcasting/Mind Spike@5",
			"rend/action/melee-attack atk+15 r15 17(2d8+8)slashing@hit 10(3d6+0)poison@hit",
			"poison-breath/action/poison-breath rch5-6 con22½ cone90 77(22d6+0)poison@failure",
			"spellcasting/action/spellcasting sc:cha21[Detect Magic@0,Mind Spike@5;Geas@0,Modify Memory@0]",
			"mind-invasion/legendary-action/mind-invasion cast:Mind Spike@5",
			"noxious-miasma/legendary-action/noxious-miasma con21 sphere30@90 17(5d6+0)poison@failure pen:armor-class:2 1/cycle",
			"pounce/legendary-action/pounce mv:half-speed x1rend",
		}},
	{id: "monsters/ancient-red-dragon", name: "Ancient Red Dragon", pages: []int{319, 319}, size: "gargantuan", typ: "dragon", tag: "chromatic", alignment: "Chaotic Evil",
		ac: 22, hp: MonsterHP{507, 26, 20, 234}, initiative: MonsterInitiative{14, 24}, speeds: map[string]int{"walk": 40, "climb": 40, "fly": 80},
		abilities: [6][2]int{{30, 10}, {10, 7}, {29, 9}, {18, 4}, {15, 9}, {27, 8}}, skills: map[string]int{"perception": 16, "stealth": 7},
		immunities: "fire;", special: dragonSenses, passive: 26, languages: dragonLanguages, cr: "24", xp: 62000, lairXP: 75000, pb: 7,
		legendary: dragonLegendary, traits: []string{"Legendary Resistance"},
		actions: []string{
			"multiattack/action/multiattack x3rend rep1:spellcasting/Scorching Ray@3",
			"rend/action/melee-attack atk+17 r15 19(2d8+10)slashing@hit 10(3d6+0)fire@hit",
			"fire-breath/action/fire-breath rch5-6 dex24½ cone90 91(26d6+0)fire@failure",
			"spellcasting/action/spellcasting sc:cha23+15[Command@2,Detect Magic@0,Scorching Ray@3;Fireball@6,Scrying@0]",
			"commanding-presence/legendary-action/commanding-presence 1/cycle cast:Command@2",
			"fiery-rays/legendary-action/fiery-rays 1/cycle cast:Scorching Ray@3",
			"pounce/legendary-action/pounce mv:half-speed x1rend",
		}},
	{id: "monsters/ancient-silver-dragon", name: "Ancient Silver Dragon", pages: []int{325, 325}, size: "gargantuan", typ: "dragon", tag: "metallic", alignment: "Lawful Good",
		ac: 22, hp: MonsterHP{468, 24, 20, 216}, initiative: MonsterInitiative{14, 24}, speeds: map[string]int{"walk": 40, "fly": 80},
		abilities: [6][2]int{{30, 10}, {10, 7}, {29, 9}, {18, 4}, {15, 9}, {26, 8}}, skills: map[string]int{"history": 11, "perception": 16, "stealth": 7},
		immunities: "cold;", special: dragonSenses, passive: 26, languages: dragonLanguages, cr: "23", xp: 50000, lairXP: 62000, pb: 7,
		legendary: dragonLegendary, traits: []string{"Legendary Resistance"},
		actions: []string{
			"multiattack/action/multiattack x3rend rep1:paralyzing-breath|spellcasting/Ice Knife@2",
			"rend/action/melee-attack atk+17 r15 19(2d8+10)slashing@hit 9(2d8+0)cold@hit",
			"cold-breath/action/cold-breath rch5-6 con24½ cone90 67(15d8+0)cold@failure",
			"paralyzing-breath/action/paralyzing-breath con24 cone90 c:incapacitated,paralyzed min1",
			"spellcasting/action/spellcasting sc:cha23+15[Detect Magic@0,Hold Monster@0,Ice Knife@2,Shapechange@0*;Control Weather@0,Ice Storm@7,Teleport@0,Zone of Truth@0]",
			"chill/legendary-action/chill 1/cycle cast:Hold Monster@0",
			"cold-gale/legendary-action/cold-gale dex23½ line60x10 14(4d6+0)cold@failure push30 1/cycle",
			"pounce/legendary-action/pounce mv:half-speed x1rend",
		}},
	{id: "monsters/ancient-white-dragon", name: "Ancient White Dragon", pages: []int{341, 341}, size: "gargantuan", typ: "dragon", tag: "chromatic", alignment: "Chaotic Evil",
		ac: 20, hp: MonsterHP{333, 18, 20, 144}, initiative: MonsterInitiative{12, 22}, speeds: map[string]int{"walk": 40, "burrow": 40, "fly": 80, "swim": 40},
		abilities: [6][2]int{{26, 8}, {10, 6}, {26, 8}, {10, 0}, {13, 7}, {18, 4}}, skills: map[string]int{"perception": 13, "stealth": 6},
		immunities: "cold;", special: dragonSenses, passive: 23, languages: dragonLanguages, cr: "20", xp: 25000, lairXP: 33000, pb: 6,
		legendary: dragonLegendary, traits: []string{"Ice Walk", "Legendary Resistance"},
		actions: []string{
			"multiattack/action/multiattack x3rend",
			"rend/action/melee-attack atk+14 r15 17(2d8+8)slashing@hit 7(2d6+0)cold@hit",
			"cold-breath/action/cold-breath rch5-6 con22½ cone90 63(14d8+0)cold@failure",
			"freezing-burst/legendary-action/freezing-burst con20 sphere30@120 14(4d6+0)cold@failure 1/cycle",
			"frightful-presence/legendary-action/frightful-presence 1/cycle cast:Fear@0! sc:cha18[;]",
			"pounce/legendary-action/pounce mv:half-speed x1rend",
		}},
}

// batch02Line renders every typed member of one action so a literal can pin it.
func batch02Line(a MonsterAction, d MonsterBatch02ActionDetail) string {
	out := []string{d.ActionID + "/" + d.Section + "/" + strings.SplitN(a.Key, "/", 2)[1]}
	if at := a.Attack; at != nil {
		out = append(out, fmt.Sprintf("atk+%d r%d", at.Bonus, at.Reach))
		if at.Kind != "melee" || at.Range != 0 || at.LongRange != 0 || len(d.Damage) == 0 ||
			at.Damage != (MonsterHP{d.Damage[0].Average, d.Damage[0].Dice, d.Damage[0].Sides, d.Damage[0].Modifier}) || at.DamageType != d.Damage[0].Type {
			out = append(out, "ATTACK-MISMATCH")
		}
	}
	if d.Recharge != nil {
		out = append(out, fmt.Sprintf("rch%d-%d", d.Recharge.Min, d.Recharge.Max))
	}
	if d.Save != nil {
		half := ""
		if d.Save.HalfOnSuccess {
			half = "½"
		}
		out = append(out, fmt.Sprintf("%s%d%s", d.Save.Ability, d.Save.DC, half))
	}
	if ar := d.Area; ar != nil {
		switch ar.Shape {
		case "line":
			out = append(out, fmt.Sprintf("line%dx%d", ar.Feet, ar.WidthFeet))
		case "sphere":
			out = append(out, fmt.Sprintf("sphere%d@%d", ar.Feet, ar.RangeFeet))
		default:
			out = append(out, fmt.Sprintf("%s%d", ar.Shape, ar.Feet))
		}
	}
	if d.TargetRangeFeet != 0 {
		out = append(out, fmt.Sprintf("tgt%d", d.TargetRangeFeet))
	}
	for _, r := range d.Damage {
		out = append(out, fmt.Sprintf("%d(%dd%d+%d)%s@%s", r.Average, r.Dice, r.Sides, r.Modifier, r.Type, r.Trigger))
	}
	if len(d.Conditions) > 0 {
		out = append(out, "c:"+strings.Join(d.Conditions, ","))
	}
	if d.PushFeet != 0 {
		out = append(out, fmt.Sprintf("push%d", d.PushFeet))
	}
	if d.DurationMinutes != 0 {
		out = append(out, fmt.Sprintf("min%d", d.DurationMinutes))
	}
	for _, p := range d.Penalties {
		switch {
		case p.Flat != 0:
			out = append(out, fmt.Sprintf("pen:%s:%d", p.AppliesTo, p.Flat))
		case p.Average != 0:
			out = append(out, fmt.Sprintf("pen:%s:%d(%dd%d)", p.AppliesTo, p.Average, p.Dice, p.Sides))
		default:
			out = append(out, fmt.Sprintf("pen:%s:%dd%d", p.AppliesTo, p.Dice, p.Sides))
		}
	}
	if d.OncePerTurnCycle {
		out = append(out, "1/cycle")
	}
	if d.Movement != "" {
		out = append(out, "mv:"+d.Movement)
	}
	for _, u := range d.Attacks {
		out = append(out, fmt.Sprintf("x%d%s", u.Count, u.ActionID))
	}
	if r := d.Replacement; r != nil {
		var options []string
		for _, o := range r.Options {
			if o.Spell != "" {
				options = append(options, fmt.Sprintf("%s/%s@%d", o.ActionID, o.Spell, o.Level))
			} else {
				options = append(options, o.ActionID)
			}
		}
		out = append(out, fmt.Sprintf("rep%d:%s", r.Attacks, strings.Join(options, "|")))
	}
	for _, c := range d.Casts {
		direct := ""
		if !c.ViaSpellcasting {
			direct = "!"
		}
		out = append(out, fmt.Sprintf("cast:%s@%d%s", c.Spell, c.Level, direct))
	}
	if sc := d.Spellcasting; sc != nil {
		lists := map[string][]string{}
		for _, sp := range sc.Spells {
			mark := ""
			if sp.Restriction != "" {
				mark = "*"
			}
			lists[sp.Frequency] = append(lists[sp.Frequency], fmt.Sprintf("%s@%d%s", sp.Name, sp.Level, mark))
		}
		bonus := ""
		if sc.AttackBonus != nil {
			bonus = fmt.Sprintf("+%d", *sc.AttackBonus)
		}
		out = append(out, fmt.Sprintf("sc:%s%d%s[%s;%s]", sc.Ability, sc.SaveDC, bonus, strings.Join(lists["at-will"], ","), strings.Join(lists["1/day"], ",")))
	}
	return strings.Join(out, " ")
}

func TestMonsterBatch02SourceFixtures(t *testing.T) {
	b := monsterBatch02(t)
	if len(b.StatBlocks) != len(monsterBatch02Want) || len(b.StatBlocks) != 11 {
		t.Fatalf("stat blocks = %d, want 11", len(b.StatBlocks))
	}
	const restriction = "Beast or Humanoid form only, no Temporary Hit Points gained from the spell, and no Concentration or Temporary Hit Points required to maintain the spell"
	for i, want := range monsterBatch02Want {
		t.Run(want.name, func(t *testing.T) {
			s := b.StatBlocks[i]
			m := s.Monster
			if m.ID != "srd-5.2.1:content/"+want.id || m.Name != want.name || !reflect.DeepEqual(m.SourcePages, want.pages) {
				t.Fatalf("identity %s %v", m.ID, m.SourcePages)
			}
			if !reflect.DeepEqual(m.Sizes, []string{want.size}) || m.Type != want.typ || !reflect.DeepEqual(s.Tags, []string{want.tag}) || m.Alignment != want.alignment {
				t.Fatal("size, type, tag or alignment differs")
			}
			if m.AC != want.ac || m.HP != want.hp || m.Initiative != want.initiative || !reflect.DeepEqual(m.Speeds, want.speeds) {
				t.Fatal("AC, HP, initiative or speeds differ")
			}
			for j, key := range monsterBatch02Abilities {
				score, save := want.abilities[j][0], want.abilities[j][1]
				if m.Abilities[key] != (MonsterAbility{score, score/2 - 5, save}) {
					t.Fatalf("%s = %+v", key, m.Abilities[key])
				}
			}
			immunities := strings.Join(s.DamageImmunities, ",") + ";" + strings.Join(s.ConditionImmunities, ",")
			if want.immunities == "" {
				want.immunities = ";"
			}
			if !reflect.DeepEqual(s.Skills, want.skills) || immunities != want.immunities || len(s.DamageResistances) != 0 || len(s.DamageVulnerabilities) != 0 {
				t.Fatalf("skills %v or immunities %q differ", s.Skills, immunities)
			}
			if !reflect.DeepEqual(m.Senses.Special, want.special) || m.Senses.PassivePerception != want.passive || !reflect.DeepEqual(m.Languages, want.languages) || len(m.Gear) != 0 {
				t.Fatal("senses, languages or gear differ")
			}
			if m.CR != want.cr || m.XP != want.xp || s.LairXP != want.lairXP || m.ProficiencyBonus != want.pb {
				t.Fatal("CR, XP, lair XP or proficiency differs")
			}
			legendary := ""
			if l := s.Legendary; l != nil {
				legendary = fmt.Sprintf("%d/%d LR %d/%d", l.Uses, l.LairUses, l.ResistanceUses, l.ResistanceLairUses)
				if !strings.HasPrefix(l.Text, fmt.Sprintf("Legendary Action Uses: %d (%d in Lair). Immediately after another creature’s turn", l.Uses, l.LairUses)) {
					t.Fatal("legendary rule text lost")
				}
			}
			if legendary != want.legendary {
				t.Fatalf("legendary = %q", legendary)
			}
			names := []string{}
			for _, tr := range m.Traits {
				names = append(names, tr.Name)
				if tr.Key != "trait/"+strings.ReplaceAll(strings.ToLower(tr.Name), " ", "-") || !strings.HasPrefix(tr.Text, tr.Name) {
					t.Fatalf("trait %q", tr.Name)
				}
			}
			if !reflect.DeepEqual(names, want.traits) {
				t.Fatalf("traits = %v", names)
			}
			lines := []string{}
			for j, a := range m.Actions {
				lines = append(lines, batch02Line(a, s.ActionDetails[j]))
				if !strings.HasPrefix(a.Text, a.Name) || a.ID != strings.ReplaceAll(strings.ToLower(a.Name), " ", "-") {
					t.Fatalf("action %q does not carry its source name", a.ID)
				}
				if sc := s.ActionDetails[j].Spellcasting; sc != nil {
					for _, sp := range sc.Spells {
						if sp.Restriction != "" && (sp.Name != "Shapechange" || sp.Restriction != restriction) {
							t.Fatalf("restriction on %s", sp.Name)
						}
					}
				}
			}
			if !reflect.DeepEqual(lines, want.actions) {
				t.Fatalf("actions differ:\n got %s\nwant %s", strings.Join(lines, "\n     "), strings.Join(want.actions, "\n     "))
			}
			if s.Available() || m.Available() {
				t.Fatal("catalogue presence cannot grant NPC behavior")
			}
		})
	}
	// Source oddities kept verbatim rather than normalized.
	gold := b.StatBlocks[6].Monster.Actions[1].Text
	if gold != "Rend. Melee Attack Roll: +17 to hit, reach 15 ft. Hit: 19 (2d8 + 10) Slashing damage plus 9 (2d8) Fire damage." {
		t.Fatalf("gold Rend text = %q", gold)
	}
	if text := b.StatBlocks[5].Monster.Actions[2].Text; !strings.Contains(text, "each creature in an 90-foot-long, 10-foot-wide Line") {
		t.Fatalf("copper breath text = %q", text)
	}
	// The green dragon's legendary actions continue on the next page.
	if last := b.StatBlocks[7].Monster.Actions; last[len(last)-1].Text != "Pounce. The dragon moves up to half its Speed, and it makes one Rend attack." {
		t.Fatal("green dragon continuation lost")
	}
	if c := b.StatBlocks[0].Monster.Actions[1]; c.Text != "Claws. Melee Attack Roll: +6, reach 5 ft. Hit: 8 (1d8 + 4) Slashing damage. If the target is a Large or smaller creature and the allosaurus moved 30+ feet straight toward it immediately before the hit, the target has the Prone condition, and the allosaurus can make one Bite attack against it." {
		t.Fatal("allosaurus column continuation lost")
	}
}

func TestMonsterBatch02MatchesInventory(t *testing.T) {
	data, err := os.ReadFile("../../../docs/investigations/issue-82-world-inventory.json")
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Records []struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			SourcePage     int    `json:"source_page"`
			Family         string `json:"family"`
			Status         string `json:"status"`
			CatalogueIssue int    `json:"catalogue_issue"`
		} `json:"records"`
	}
	if err := json.Unmarshal(data, &inventory); err != nil {
		t.Fatal(err)
	}
	type entry struct {
		name string
		page int
	}
	want, got := map[string]entry{}, map[string]entry{}
	for _, r := range inventory.Records {
		if r.CatalogueIssue == 152 {
			want[r.ID] = entry{r.Name, r.SourcePage}
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("inventory record %s", r.ID)
			}
		}
	}
	b := monsterBatch02(t)
	var order []string
	for _, s := range b.StatBlocks {
		got[s.Monster.ID] = entry{s.Monster.Name, s.Monster.SourcePages[0]}
		order = append(order, s.Monster.Name)
	}
	for _, h := range b.Headings {
		got[h.ID] = entry{h.Name, h.SourcePages[0]}
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch = %v\ninventory = %v", got, want)
	}
	if !slices.IsSorted(order) {
		t.Fatal("stat blocks must stay in index order")
	}
	if !reflect.DeepEqual(b.Headings, []MonsterBatch02Heading{{ID: "srd-5.2.1:content/animals/animals", Name: "Animals", Kind: "section-heading", SourcePages: []int{344, 344}}}) {
		t.Fatalf("headings = %+v", b.Headings)
	}
	if b.Source.Attribution != Attribution || b.Source.SHA256 != artifactSHA256 || !strings.Contains(b.Source.Changes, "No playable NPC behavior is granted.") {
		t.Fatal("provenance or adaptation notice lost")
	}
}

func TestMonsterBatch02MechanicKeysStayUnavailable(t *testing.T) {
	b := monsterBatch02(t)
	bandit, _, err := BanditRecord()
	if err != nil {
		t.Fatal(err)
	}
	banditKeys := map[string]bool{}
	for _, key := range bandit.RequiredMechanics {
		banditKeys[key] = true
	}
	union := map[string]bool{}
	counts := []int{}
	for _, s := range b.StatBlocks {
		m := s.Monster
		listed := map[string]bool{}
		for _, key := range m.RequiredMechanics {
			listed[key] = true
			union[key] = true
		}
		used := []string{"monster/stat-block"}
		for _, tr := range m.Traits {
			used = append(used, tr.Key)
		}
		for _, a := range m.Actions {
			used = append(used, a.Key)
			if len(a.Effects) == 0 {
				t.Fatalf("%s %s lists no effect", m.Name, a.ID)
			}
			for _, e := range a.Effects {
				used = append(used, e.Key)
			}
		}
		for _, key := range used {
			if !listed[key] {
				t.Fatalf("%s uses unlisted key %s", m.Name, key)
			}
		}
		for key := range listed {
			family := strings.SplitN(key, "/", 2)[0]
			if !slices.Contains(used, key) && !slices.Contains([]string{"speed", "movement", "spell", "skill", "immunity", "sense", "monster"}, family) {
				t.Fatalf("%s lists unused key %s", m.Name, key)
			}
		}
		// The Bandit fixture's capability set must not unlock any imported block.
		if m.availableWith(banditKeys) || s.Available() {
			t.Fatalf("%s became available", m.Name)
		}
		counts = append(counts, len(m.RequiredMechanics))
	}
	if !reflect.DeepEqual(counts, []int{9, 55, 50, 64, 66, 62, 67, 52, 48, 64, 39}) || len(union) != 192 {
		t.Fatalf("required key counts = %v, union = %d", counts, len(union))
	}
	// Composition with the earlier Bandit record and the #151 batch: disjoint
	// IDs under one pinned source, Bandit untouched.
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch01 := loadMonsters151(t)
	ids := map[string]bool{}
	for _, m := range slices.Concat(base.Monsters, batch01.Monsters) {
		ids[m.ID] = true
	}
	combined := MonsterCatalogue{SchemaVersion: 1, Source: batch01.Source, Monsters: slices.Concat(base.Monsters, batch01.Monsters)}
	for _, s := range b.StatBlocks {
		combined.Monsters = append(combined.Monsters, s.Monster)
	}
	if err := combined.Validate(); err != nil || len(combined.Monsters) != 24 {
		t.Fatal("batch must compose with the earlier catalogues:", err)
	}
	samePin := func(a, b Source) bool { a.Changes, b.Changes = "", ""; return a == b }
	if len(batch01.Monsters) != 12 || !samePin(batch01.Source, b.Source) {
		t.Fatal("batch must share the #151 batch's source pin")
	}
	for _, s := range b.StatBlocks {
		if ids[s.Monster.ID] {
			t.Fatalf("%s duplicates an earlier record", s.Monster.ID)
		}
		ids[s.Monster.ID] = true
	}
	if len(base.Monsters) != 1 || len(ids) != 24 || base.Source != func() Source { s := b.Source; s.Changes = base.Source.Changes; return s }() {
		t.Fatal("batch must compose with the Bandit catalogue and the #151 batch under one pinned source")
	}
}

func TestMonsterBatch02Malformed(t *testing.T) {
	tests := map[string]func(*MonsterBatch02){
		"schema":            func(b *MonsterBatch02) { b.SchemaVersion = 2 },
		"pin":               func(b *MonsterBatch02) { b.Source.SHA256 = "bad" },
		"attribution":       func(b *MonsterBatch02) { b.Source.Attribution = "" },
		"adaptation":        func(b *MonsterBatch02) { b.Source.Changes = " " },
		"no headings":       func(b *MonsterBatch02) { b.Headings = nil },
		"heading kind":      func(b *MonsterBatch02) { b.Headings[0].Kind = "stat-block" },
		"heading pages":     func(b *MonsterBatch02) { b.Headings[0].SourcePages = []int{344, 365} },
		"duplicate heading": func(b *MonsterBatch02) { b.Headings = append(b.Headings, b.Headings[0]) },
		"heading as block":  func(b *MonsterBatch02) { b.Headings[0].ID = b.StatBlocks[0].Monster.ID },
		"duplicate block":   func(b *MonsterBatch02) { b.StatBlocks = append(b.StatBlocks, b.StatBlocks[1]) },
		"wrong family":      func(b *MonsterBatch02) { b.StatBlocks[0].Monster.ID = "srd-5.2.1:content/spells/allosaurus" },
		"ability modifier":  func(b *MonsterBatch02) { b.StatBlocks[1].Monster.Abilities["str"] = MonsterAbility{27, 9, 9} },
		"save":              func(b *MonsterBatch02) { b.StatBlocks[1].Monster.Abilities["dex"] = MonsterAbility{14, 2, 8} },
		"XP":                func(b *MonsterBatch02) { b.StatBlocks[1].Monster.XP = 33001 },
		"lair XP":           func(b *MonsterBatch02) { b.StatBlocks[1].LairXP = 40000 },
		"lair XP key":       func(b *MonsterBatch02) { dropKey(&b.StatBlocks[1].Monster, "monster/lair-experience") },
		"skill bonus":       func(b *MonsterBatch02) { b.StatBlocks[1].Skills["stealth"] = 10 },
		"unknown skill":     func(b *MonsterBatch02) { b.StatBlocks[1].Skills["cooking"] = 9 },
		"skill key":         func(b *MonsterBatch02) { dropKey(&b.StatBlocks[0].Monster, "skill/perception") },
		"passive":           func(b *MonsterBatch02) { b.StatBlocks[0].Monster.Senses.PassivePerception = 11 },
		"sense key":         func(b *MonsterBatch02) { dropKey(&b.StatBlocks[1].Monster, "sense/blindsight") },
		"speed key":         func(b *MonsterBatch02) { dropKey(&b.StatBlocks[1].Monster, "speed/fly") },
		"immunity type":     func(b *MonsterBatch02) { b.StatBlocks[1].DamageImmunities = []string{"sonic"} },
		"immunity key":      func(b *MonsterBatch02) { b.StatBlocks[1].DamageImmunities = []string{"acid", "fire"} },
		"condition key":     func(b *MonsterBatch02) { b.StatBlocks[7].ConditionImmunities = []string{"poisoned", "charmed"} },
		"no resistances":    func(b *MonsterBatch02) { b.StatBlocks[1].DamageResistances = nil },
		"no tags":           func(b *MonsterBatch02) { b.StatBlocks[1].Tags = nil },
		"tag":               func(b *MonsterBatch02) { b.StatBlocks[1].Tags = []string{"Chromatic"} },
		"no legendary":      func(b *MonsterBatch02) { b.StatBlocks[1].Legendary = nil },
		"extra legendary":   func(b *MonsterBatch02) { b.StatBlocks[0].Legendary = b.StatBlocks[1].Legendary },
		"legendary uses":    func(b *MonsterBatch02) { b.StatBlocks[1].Legendary.LairUses = 2 },
		"resistance uses":   func(b *MonsterBatch02) { b.StatBlocks[1].Legendary.ResistanceUses = 0 },
		"legendary key":     func(b *MonsterBatch02) { dropKey(&b.StatBlocks[1].Monster, "monster/legendary-actions") },
		"detail count":      func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails = b.StatBlocks[1].ActionDetails[1:] },
		"detail order": func(b *MonsterBatch02) {
			d := b.StatBlocks[1].ActionDetails
			d[0], d[1] = d[1], d[0]
		},
		"section":           func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[6].Section = "action" },
		"section value":     func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[0].Section = "reaction" },
		"attack bonus":      func(b *MonsterBatch02) { b.StatBlocks[1].Monster.Actions[1].Attack.Bonus = 16 },
		"attack damage":     func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[1].Damage[0].Dice = 3 },
		"extra damage":      func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[1].Damage[1].Average = 10 },
		"dropped damage":    func(b *MonsterBatch02) { d := &b.StatBlocks[1].ActionDetails[1]; d.Damage = d.Damage[:1] },
		"damage type":       func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Damage[0].Type = "fire" },
		"damage trigger":    func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Damage[0].Trigger = "hit" },
		"recharge":          func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Recharge = nil },
		"recharge range":    func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Recharge.Min = 7 },
		"dropped save":      func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Save = nil },
		"save ability":      func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Save.Ability = "luck" },
		"save DC":           func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Save.DC = 0 },
		"half damage":       func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Save.HalfOnSuccess = false },
		"area":              func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Area = nil },
		"area shape":        func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Area.Shape = "cone" },
		"line width":        func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Area.WidthFeet = 0 },
		"sphere range":      func(b *MonsterBatch02) { b.StatBlocks[4].ActionDetails[7].Area.RangeFeet = 0 },
		"target range":      func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[4].TargetRangeFeet = 0 },
		"condition":         func(b *MonsterBatch02) { b.StatBlocks[0].ActionDetails[1].Conditions = []string{"stunned"} },
		"push":              func(b *MonsterBatch02) { b.StatBlocks[4].ActionDetails[3].PushFeet = 0 },
		"penalty":           func(b *MonsterBatch02) { b.StatBlocks[7].ActionDetails[5].Penalties[0].Flat = 0 },
		"penalty target":    func(b *MonsterBatch02) { b.StatBlocks[7].ActionDetails[5].Penalties[0].AppliesTo = "speed" },
		"cycle limit":       func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[4].OncePerTurnCycle = false },
		"movement":          func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[6].Movement = "half-fly-speed" },
		"attack use":        func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[0].Attacks[0].ActionID = "bite" },
		"attack count":      func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[0].Attacks[0].Count = 0 },
		"no attack list":    func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[0].Attacks = []MonsterBatch02AttackUse{} },
		"replacement":       func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[0].Replacement = nil },
		"replacement spell": func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[0].Replacement.Options[0].Level = 3 },
		"replace action":    func(b *MonsterBatch02) { b.StatBlocks[3].ActionDetails[0].Replacement.Options[0].ActionID = "tail" },
		"cast level":        func(b *MonsterBatch02) { b.StatBlocks[2].ActionDetails[5].Casts[0].Level = 0 },
		"cast spell":        func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[5].Casts[0].Spell = "Wish" },
		"spell DC":          func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[3].Spellcasting.SaveDC = 22 },
		"spell attack":      func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[3].Spellcasting.AttackBonus = nil },
		"invented attack":   func(b *MonsterBatch02) { n := 14; b.StatBlocks[2].ActionDetails[3].Spellcasting.AttackBonus = &n },
		"spell list":        func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[3].Spellcasting = nil },
		"spell frequency":   func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[3].Spellcasting.Spells[0].Frequency = "3/day" },
		"unlisted spell":    func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[3].Spellcasting.Spells[0].Name = "Wish" },
		"restriction": func(b *MonsterBatch02) {
			b.StatBlocks[3].ActionDetails[4].Spellcasting.Spells[3].Restriction = "Dragon form only"
		},
		"effect text":     func(b *MonsterBatch02) { b.StatBlocks[1].Monster.Actions[2].Effects[0].Text = "Recharge 4–6" },
		"unlisted effect": func(b *MonsterBatch02) { b.StatBlocks[1].Monster.Actions[2].Effects[0].Key = "effect/new-effect" },
		"missing keys":    func(b *MonsterBatch02) { b.StatBlocks[1].Monster.RequiredMechanics = []string{"monster/stat-block"} },
	}
	// Each case gets a fresh unvalidated copy and one Validate call. The
	// decoder runs the same Validate, so its round trip is checked for a
	// sample of cases instead of re-decoding the 223 KB document every time.
	roundTrip := []string{"schema", "XP", "section", "spell DC", "effect text"}
	fresh := func(t *testing.T) MonsterBatch02 {
		t.Helper()
		var b MonsterBatch02
		if err := json.Unmarshal(monsterBatch02JSON, &b); err != nil {
			t.Fatal(err)
		}
		return b
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			b := fresh(t)
			mutate(&b)
			if b.Validate() == nil {
				t.Fatal("accepted malformed batch")
			}
			if !slices.Contains(roundTrip, name) {
				return
			}
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := DecodeMonsterBatch02(bytes.NewReader(data))
			if err == nil || len(decoded.StatBlocks) != 0 {
				t.Fatal("decoder must return no batch on invalid data")
			}
		})
	}
	for _, name := range roundTrip {
		if tests[name] == nil {
			t.Fatalf("round-trip case %q is not in the table", name)
		}
	}
}

// Typed source numbers must be printed in the retained text (review of #152):
// a self-consistent wrong number is rejected even without the literal table.
func TestMonsterBatch02TypedNumbersFollowText(t *testing.T) {
	for name, mutate := range map[string]func(*MonsterBatch02){
		"breath DC":       func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Save.DC = 23 },
		"DC prefix":       func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Save.DC = 2 },
		"line length":     func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Area.Feet = 60 },
		"line width":      func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Area.WidthFeet = 5 },
		"cone length":     func(b *MonsterBatch02) { b.StatBlocks[3].ActionDetails[3].Area.Feet = 60 },
		"sphere radius":   func(b *MonsterBatch02) { b.StatBlocks[4].ActionDetails[7].Area.Feet = 30 },
		"sphere range":    func(b *MonsterBatch02) { b.StatBlocks[4].ActionDetails[7].Area.RangeFeet = 90 },
		"target range":    func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[4].TargetRangeFeet = 60 },
		"rend reach":      func(b *MonsterBatch02) { b.StatBlocks[1].Monster.Actions[1].Attack.Reach = 10 },
		"recharge":        func(b *MonsterBatch02) { b.StatBlocks[1].ActionDetails[2].Recharge.Min = 4 },
		"push":            func(b *MonsterBatch02) { b.StatBlocks[4].ActionDetails[3].PushFeet = 30 },
		"duration":        func(b *MonsterBatch02) { b.StatBlocks[3].ActionDetails[3].DurationMinutes = 1 },
		"flat penalty":    func(b *MonsterBatch02) { b.StatBlocks[7].ActionDetails[5].Penalties[0].Flat = 4 },
		"rolled penalty":  func(b *MonsterBatch02) { b.StatBlocks[5].ActionDetails[5].Penalties[0].Sides = 6 },
		"resistance uses": func(b *MonsterBatch02) { b.StatBlocks[1].Legendary.ResistanceUses = 3 },
		"lair resistance": func(b *MonsterBatch02) { b.StatBlocks[1].Legendary.ResistanceLairUses = 6 },
		"legendary uses":  func(b *MonsterBatch02) { b.StatBlocks[1].Legendary.Uses = 2 },
		"save damage": func(b *MonsterBatch02) {
			b.StatBlocks[1].ActionDetails[2].Damage[0] = MonsterBatch02Damage{Average: 72, Dice: 16, Sides: 8, Type: "acid", Trigger: "failure"}
		},
		"extra hit damage": func(b *MonsterBatch02) {
			r := &b.StatBlocks[1].ActionDetails[1].Damage[1]
			r.Average, r.Dice, r.Sides = 7, 2, 6
		},
	} {
		t.Run(name, func(t *testing.T) {
			var b MonsterBatch02
			if err := json.Unmarshal(monsterBatch02JSON, &b); err != nil {
				t.Fatal(err)
			}
			mutate(&b)
			if b.Validate() == nil {
				t.Fatal("accepted a typed number the source text does not print")
			}
		})
	}
	// The source prints "+17 to hit" for the Ancient Gold Dragon's Rend (p. 292)
	// and "+N" elsewhere; both wordings satisfy the check.
	gold := monsterBatch02(t).StatBlocks[6].Monster.Actions[1]
	if !strings.HasPrefix(gold.Text, "Rend. Melee Attack Roll: +17 to hit, reach 15 ft.") {
		t.Fatalf("unexpected text %q", gold.Text)
	}
	for text, phrase := range map[string]string{"DC 22,": "DC 2", "a 190-foot Cone": "90-foot Cone"} {
		if monsterBatch02Prints(text, phrase) {
			t.Fatalf("%q matched inside %q", phrase, text)
		}
	}
}

func dropKey(m *Monster, key string) {
	m.RequiredMechanics = slices.DeleteFunc(slices.Clone(m.RequiredMechanics), func(k string) bool { return k == key })
}

func TestMonsterBatch02StrictDecoding(t *testing.T) {
	text := string(monsterBatch02JSON)
	for name, data := range map[string]string{
		"trailing object":  text + " {}",
		"unknown member":   strings.Replace(text, `"tags": [`, `"lair_actions": [], "tags": [`, 1),
		"unknown detail":   strings.Replace(text, `"push_feet": 0,`, `"push_feet": 0, "pull_feet": 0,`, 1),
		"duplicate member": strings.Replace(text, `"lair_xp": 41000,`, `"lair_xp": 1, "lair_xp": 41000,`, 1),
		"omitted lair XP":  strings.Replace(text, `"lair_xp": 41000,`, ``, 1),
		"omitted save":     strings.Replace(text, `"save": null,`, ``, 1),
		"omitted spells":   strings.Replace(text, `"spellcasting": null`, `"casts_per_day": 0`, 1),
		"malformed":        `{`,
	} {
		t.Run(name, func(t *testing.T) {
			if data == text {
				t.Fatal("mutation did not apply")
			}
			if _, err := DecodeMonsterBatch02(strings.NewReader(data)); err == nil {
				t.Fatal("accepted non-strict document")
			}
		})
	}
	// Every load owns its data; no caller can mutate the embedded batch.
	b := monsterBatch02(t)
	b.StatBlocks[0].Skills["perception"] = 0
	if monsterBatch02(t).StatBlocks[0].Skills["perception"] != 5 {
		t.Fatal("shared mutable batch")
	}
}

func TestMonsterAnimalIdentitiesShareTheContract(t *testing.T) {
	m := monsterBatch02(t).StatBlocks[0].Monster
	if !strings.HasPrefix(m.ID, "srd-5.2.1:content/animals/") || m.Validate() != nil {
		t.Fatal("inventory animal identity must validate")
	}
	for _, id := range []string{"srd-5.2.1:content/spells/allosaurus", "srd-5.2.1:content/animals/Bad ID", "srd-5.2.1:family/animals"} {
		m.ID = id
		if m.Validate() == nil {
			t.Fatalf("accepted %s", id)
		}
	}
}
