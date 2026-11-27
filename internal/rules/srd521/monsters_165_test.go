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

func loadMonsters165(t *testing.T) MonsterCatalogue {
	t.Helper()
	b, err := os.ReadFile("testdata/monsters_165.json")
	if err != nil {
		t.Fatal(err)
	}
	c, err := DecodeMonsters(bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Stat-block header lines from the pinned PDF, in issue #165 order, with each
// record's inclusive page range. Wrapped lines are rejoined; nothing else is
// normalised. The Invisible Stalker prints "Initiative +7 (22)" (p. 301); the record keeps the
// score 17 that the modifier gives.
var monsters165Headers = []struct {
	name   string
	pages  [2]int
	header string
}{
	{"Hobgoblin Captain", [2]int{299, 299}, "Medium Fey (Goblinoid), Lawful Evil|AC 17 Initiative +4 (14)|HP 58 (9d8 + 18)|Speed 30 ft.|Str 15 +2 +2 Dex 14 +2 +2 Con 14 +2 +2|Int 12 +1 +1 Wis 10 +0 +0 Cha 13 +1 +1|Gear Greatsword, Half Plate Armor, Longbow|Senses Darkvision 60 ft.; Passive Perception 10|Languages Common, Goblin|CR 3 (XP 700; PB +2)"},
	{"Hobgoblin Warrior", [2]int{298, 299}, "Medium Fey (Goblinoid), Lawful Evil|AC 18 Initiative +3 (13)|HP 11 (2d8 + 2)|Speed 30 ft.|Str 13 +1 +1 Dex 12 +1 +1 Con 12 +1 +1|Int 10 +0 +0 Wis 10 +0 +0 Cha 9 −1 −1|Gear Half Plate Armor, Longbow, Longsword, Shield|Senses Darkvision 60 ft.; Passive Perception 10|Languages Common, Goblin|CR 1/2 (XP 100; PB +2)"},
	{"Homunculus", [2]int{299, 299}, "Tiny Construct, Neutral|AC 13 Initiative +2 (12)|HP 4 (1d4 + 2)|Speed 20 ft., Fly 40 ft.|Str 4 −3 −3 Dex 15 +2 +2 Con 14 +2 +2|Int 10 +0 +0 Wis 10 +0 +2 Cha 7 −2 +0|Immunities Poison; Charmed, Poisoned|Senses Darkvision 60 ft.; Passive Perception 10|Languages Understands Common plus one other language but can’t speak|CR 0 (XP 10; PB +2)"},
	{"Horned Devil", [2]int{299, 300}, "Large Fiend (Devil), Lawful Evil|AC 18 Initiative +7 (17)|HP 199 (19d10 + 95)|Speed 30 ft., Fly 60 ft.|Str 22 +6 +10 Dex 17 +3 +7 Con 21 +5 +5|Int 12 +1 +1 Wis 16 +3 +7 Cha 18 +4 +8|Resistances Cold|Immunities Fire, Poison; Poisoned|Senses Darkvision 150 ft. (unimpeded by magical Darkness); Passive Perception 13|Languages Infernal; telepathy 120 ft.|CR 11 (XP 7,200; PB +4)"},
	{"Hunter Shark", [2]int{356, 356}, "Large Beast, Unaligned|AC 12 Initiative +2 (12)|HP 45 (6d10 + 12)|Speed 5 ft., Swim 40 ft.|Str 18 +4 +4 Dex 14 +2 +2 Con 15 +2 +2|Int 1 −5 −5 Wis 10 +0 +0 Cha 4 −3 −3|Skills Perception +2|Senses Blindsight 60 ft.; Passive Perception 12|Languages None|CR 2 (XP 450; PB +2)"},
	{"Hydra", [2]int{300, 300}, "Huge Monstrosity, Unaligned|AC 15 Initiative +4 (14)|HP 184 (16d12 + 80)|Speed 40 ft., Swim 40 ft.|Str 20 +5 +5 Dex 12 +1 +1 Con 20 +5 +5|Int 2 −4 −4 Wis 10 +0 +0 Cha 7 −2 −2|Skills Perception +6|Immunities Blinded, Charmed, Deafened, Frightened, Stunned, Unconscious|Senses Darkvision 60 ft.; Passive Perception 16|Languages None|CR 8 (XP 3,900; PB +3)"},
	{"Hyena", [2]int{356, 356}, "Medium Beast, Unaligned|AC 11 Initiative +1 (11)|HP 5 (1d8 + 1)|Speed 50 ft.|Str 11 +0 +0 Dex 13 +1 +1 Con 12 +1 +1|Int 2 −4 −4 Wis 12 +1 +1 Cha 5 −3 −3|Skills Perception +3|Senses Darkvision 60 ft.; Passive Perception 13|Languages None|CR 0 (XP 10; PB +2)"},
	{"Ice Devil", [2]int{300, 300}, "Large Fiend (Devil), Lawful Evil|AC 18 Initiative +7 (17)|HP 228 (24d10 + 96)|Speed 40 ft.|Str 21 +5 +5 Dex 14 +2 +7 Con 18 +4 +9|Int 18 +4 +4 Wis 15 +2 +7 Cha 18 +4 +9|Skills Insight +7, Perception +7, Persuasion +9|Immunities Cold, Fire, Poison; Poisoned|Senses Blindsight 120 ft.; Passive Perception 17|Languages Infernal; telepathy 120 ft.|CR 14 (XP 11,500; PB +5)"},
	{"Ice Mephit", [2]int{307, 307}, "Small Elemental, Neutral Evil|AC 11 Initiative +1 (11)|HP 21 (6d6)|Speed 30 ft., Fly 30 ft.|Str 7 −2 −2 Dex 13 +1 +1 Con 10 +0 +0|Int 9 −1 −1 Wis 11 +0 +0 Cha 12 +1 +1|Skills Perception +2, Stealth +3|Vulnerabilities Fire|Immunities Cold, Poison; Exhaustion, Poisoned|Senses Darkvision 60 ft.; Passive Perception 12|Languages Primordial (Aquan, Auran)|CR 1/2 (XP 100; PB +2)"},
	{"Imp", [2]int{300, 301}, "Tiny Fiend (Devil), Lawful Evil|AC 13 Initiative +3 (13)|HP 21 (6d4 + 6)|Speed 20 ft., Fly 40 ft.|Str 6 −2 −2 Dex 17 +3 +3 Con 13 +1 +1|Int 11 +0 +0 Wis 12 +1 +1 Cha 14 +2 +2|Skills Deception +4, Insight +3, Stealth +5|Resistances Cold|Immunities Fire, Poison; Poisoned|Senses Darkvision 120 ft. (unimpeded by magical Darkness); Passive Perception 11|Languages Common, Infernal|CR 1 (XP 200; PB +2)"},
	{"Incubus", [2]int{301, 301}, "Medium Fiend, Neutral Evil|AC 15 Initiative +3 (13)|HP 66 (12d8 + 12)|Speed 30 ft., Fly 60 ft.|Str 8 −1 −1 Dex 17 +3 +3 Con 13 +1 +1|Int 15 +2 +2 Wis 12 +1 +1 Cha 20 +5 +5|Skills Deception +9, Insight +5, Perception +5, Persuasion +9, Stealth +7|Resistances Cold, Fire, Poison, Psychic|Senses Darkvision 60 ft.; Passive Perception 15|Languages Abyssal, Common, Infernal; telepathy 60 ft.|CR 4 (XP 1,100; PB +2)"},
	{"Invisible Stalker", [2]int{301, 302}, "Large Elemental, Neutral|AC 14 Initiative +7 (22)|HP 97 (13d10 + 26)|Speed 50 ft., Fly 50 ft. (hover)|Str 16 +3 +3 Dex 19 +4 +4 Con 14 +2 +2|Int 10 +0 +0 Wis 15 +2 +2 Cha 11 +0 +0|Skills Perception +8, Stealth +10|Resistances Bludgeoning, Piercing, Slashing|Immunities Poison; Exhaustion, Grappled, Paralyzed, Petrified, Poisoned, Prone, Restrained, Unconscious|Senses Darkvision 60 ft.; Passive Perception 18|Languages Common, Primordial (Auran)|CR 6 (XP 2,300; PB +3)"},
}

// Typed values of every trait and action entry, reviewed against the same pages:
// limited uses, attack bonus/reach/range/damage, save ability and DC, target
// shape size/width/range, dice and dependency keys. A ranged attack printed with
// one range repeats it; a fixed hit prints no dice.
//
// Printed values the contract cannot type stay in the entry's retained source
// text. The Ice Devil's Ice Spear is a "Melee or Ranged Attack Roll" (key
// action/melee-or-ranged-attack). The Homunculus's Bite types the save its hit imposes.
var monsters165Entries = map[string][]string{
	"Hobgoblin Captain": {
		"trait|trait/aura-of-authority|Aura of Authority; emanation 10/0/0; area/emanation, attack/advantage, save/advantage, target/self-and-allies, requirement/not-incapacitated",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination",
		"action|greatsword|Greatsword; melee +4 reach 5 range 0/0 9 (2d6 + 2) slashing; effect/slashing-damage 9 (2d6 + 2), effect/poison-damage 3 (1d6)",
		"action|longbow|Longbow; ranged +4 reach 0 range 150/600 6 (1d8 + 2) piercing; effect/piercing-damage 6 (1d8 + 2), effect/poison-damage 5 (2d4)",
	},
	"Hobgoblin Warrior": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage, requirement/ally-within-5-feet, requirement/ally-not-incapacitated",
		"action|longsword|Longsword; melee +3 reach 5 range 0/0 12 (2d10 + 1) slashing; effect/slashing-damage 12 (2d10 + 1)",
		"action|longbow|Longbow; ranged +3 reach 0 range 150/600 5 (1d8 + 1) piercing; effect/piercing-damage 5 (1d8 + 1), effect/poison-damage 7 (3d4)",
	},
	"Homunculus": {
		"trait|trait/telepathic-bond|Telepathic Bond; effect/telepathic-link, requirement/same-plane-as-master",
		"action|bite|Bite; melee +4 reach 5 range 0/0 1 piercing; con DC 12; effect/piercing-damage [damage/fixed], requirement/target-subjected-to-save, save/con, condition/poisoned, save/failure-by-five-or-more, condition/unconscious, effect/ends-on-damage",
	},
	"Horned Devil": {
		"trait|trait/diabolical-restoration|Diabolical Restoration; effect/revive",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-attack, action/ranged-attack, multiattack/any-combination, multiattack/replace-attack, action/saving-throw",
		"action|searing-fork|Searing Fork; melee +10 reach 10 range 0/0 15 (2d8 + 6) piercing; effect/piercing-damage 15 (2d8 + 6), effect/fire-damage 9 (2d8)",
		"action|hurl-flame|Hurl Flame; ranged +8 reach 0 range 150/150 26 (5d8 + 4) fire; effect/fire-damage 26 (5d8 + 4), effect/ignite-object",
		"action|infernal-tail|Infernal Tail; dex DC 17; creature 0/0/10; save/dex, target/creature, effect/necrotic-damage 10 (1d8 + 6), effect/infernal-wound, requirement/no-existing-infernal-wound, effect/hit-point-loss 10 (3d6), effect/wound-closes, check/wis-medicine",
	},
	"Hunter Shark": {
		"trait|trait/water-breathing|Water Breathing; effect/breathe-only-underwater",
		"action|bite|Bite; melee +6 reach 5 range 0/0 14 (3d6 + 4) piercing; effect/piercing-damage 14 (3d6 + 4), attack/advantage-against-target-missing-hit-points",
	},
	"Hydra": {
		"trait|trait/hold-breath|Hold Breath; effect/hold-breath",
		"trait|trait/multiple-heads|Multiple Heads; effect/five-heads, effect/lose-head-on-damage-threshold, effect/die-without-heads, effect/regrow-heads, effect/fire-damage-prevents-regrowth, effect/regain-hit-points",
		"trait|trait/reactive-heads|Reactive Heads; effect/extra-reaction-per-head, limit/opportunity-attacks-only",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/one-attack-per-head",
		"action|bite|Bite; melee +8 reach 10 range 0/0 10 (1d10 + 5) piercing; effect/piercing-damage 10 (1d10 + 5)",
	},
	"Hyena": {
		"trait|trait/pack-tactics|Pack Tactics; attack/advantage, requirement/ally-within-5-feet, requirement/ally-not-incapacitated",
		"action|bite|Bite; melee +2 reach 5 range 0/0 3 (1d6) piercing; effect/piercing-damage 3 (1d6)",
	},
	"Ice Devil": {
		"trait|trait/diabolical-restoration|Diabolical Restoration; effect/revive",
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|multiattack|Multiattack; action/melee-or-ranged-attack, multiattack/replace-attack, action/melee-attack",
		"action|ice-spear|Ice Spear; effect/piercing-damage 14 (2d8 + 5), effect/cold-damage 10 (3d6), effect/no-bonus-action-or-reaction, effect/speed-decrease, effect/move-or-action, effect/weapon-returns",
		"action|tail|Tail; melee +10 reach 10 range 0/0 15 (3d6 + 5) bludgeoning; effect/bludgeoning-damage 15 (3d6 + 5), effect/cold-damage 18 (4d8)",
		"action|ice-wall|Ice Wall (Recharge 6); limit/recharge, action/spellcasting, spell/wall-of-ice, spellcasting/no-components, spellcasting/fixed-save-dc",
	},
	"Ice Mephit": {
		"trait|trait/death-burst|Death Burst; con DC 10; emanation 5/0/0; 5 (2d4); trigger/death, save/con, area/emanation, effect/cold-damage, save/half-damage",
		"action|claw|Claw; melee +3 reach 5 range 0/0 3 (1d4 + 1) slashing; effect/slashing-damage 3 (1d4 + 1), effect/cold-damage 2 (1d4)",
		"action|fog-cloud|Fog Cloud (1/Day); limit/per-day, action/spellcasting, spell/fog-cloud, spellcasting/no-components, spellcasting/charisma-ability",
		"action|frost-breath|Frost Breath (Recharge 6); con DC 10; cone 15/0/0; limit/recharge, save/con, area/cone, effect/cold-damage 7 (3d4), save/half-damage",
	},
	"Imp": {
		"trait|trait/magic-resistance|Magic Resistance; save/advantage-against-magic",
		"action|sting|Sting; melee +5 reach 5 range 0/0 6 (1d6 + 3) piercing; effect/piercing-damage 6 (1d6 + 3), effect/poison-damage 7 (2d6)",
		"action|invisibility|Invisibility; action/spellcasting, spell/invisibility, spellcasting/no-components, spellcasting/charisma-ability",
		"action|shape-shift|Shape-Shift; effect/shape-shift, effect/shapechange-form-limits, effect/equipment-not-transformed",
	},
	"Incubus": {
		"trait|trait/succubus-form|Succubus Form; effect/shape-shift-to-succubus, requirement/finish-long-rest, effect/use-succubus-stat-block, effect/equipment-not-transformed",
		"action|multiattack|Multiattack; action/melee-attack",
		"action|restless-touch|Restless Touch; melee +7 reach 5 range 0/0 15 (3d6 + 5) psychic; effect/psychic-damage 15 (3d6 + 5), effect/curse, effect/no-short-rest-benefit",
		"action|spellcasting|Spellcasting; cha DC 15 +0; at will: Disguise Self, Etherealness; 1/day: Dream, Hypnotic Pattern; spellcasting/no-material-components, spell/disguise-self [spellcasting/at-will], spell/etherealness [spellcasting/at-will], spell/dream [limit/per-day], spell/hypnotic-pattern [limit/per-day]",
		"bonus|nightmare|Nightmare (Recharge 6); wis DC 15; creature 0/0/60; limit/recharge, save/wis, target/creature, requirement/target-at-most-20-hit-points, condition/unconscious, effect/ends-on-damage-or-wake, effect/psychic-damage 18 (4d8)",
	},
	"Invisible Stalker": {
		"trait|trait/air-form|Air Form; movement/enter-creature-space, movement/squeeze",
		"trait|trait/invisibility|Invisibility; condition/invisible",
		"action|multiattack|Multiattack; action/melee-attack, multiattack/replace-attack, action/saving-throw",
		"action|wind-swipe|Wind Swipe; melee +7 reach 5 range 0/0 11 (2d6 + 4) force; effect/force-damage 11 (2d6 + 4)",
		"action|vortex|Vortex; con DC 14; creature 0/0/0; save/con, target/creature, requirement/target-size, effect/thunder-damage 7 (1d8 + 3), condition/grappled, effect/no-verbal-spells, effect/start-of-turn-damage 7 (2d6)",
	},
}

// Literal PDF values for every record: header fields, page ranges and each
// trait/action payload, independent of the decoder and structural validation.
func TestMonsters165SourceValues(t *testing.T) {
	c := loadMonsters165(t)
	if len(c.Monsters) != 12 || len(monsters165Headers) != 12 || len(monsters165Entries) != 12 {
		t.Fatal("exact twelve stat blocks required")
	}
	if c.Source.SHA256 != artifactSHA256 || c.Source.Attribution != Attribution || c.Source.License != "CC-BY-4.0" || !strings.Contains(c.Source.Changes, "issue #165") || !strings.Contains(c.Source.Changes, "No playable NPC behavior") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	entries := 0
	for i, want := range monsters165Headers {
		m := c.Monsters[i]
		t.Run(want.name, func(t *testing.T) {
			if m.Name != want.name || !slices.Equal(m.SourcePages, want.pages[:]) {
				t.Fatalf("identity or pages differ: %s %v", m.Name, m.SourcePages)
			}
			if got := monsterHeader164(m); got != want.header {
				t.Fatalf("header differs from PDF p.%d:\n got %s\nwant %s", want.pages[0], got, want.header)
			}
			if got := monsterEntries164(m); !slices.Equal(got, monsters165Entries[want.name]) {
				t.Fatalf("entries differ from PDF pp.%d–%d:\n got %q\nwant %q", want.pages[0], want.pages[1], got, monsters165Entries[want.name])
			}
			if m.ResistanceNote != "" || m.LairXP != 0 || m.Legendary != nil || m.LegendaryActions != nil {
				t.Fatal("source lists no resistance note, lair XP, legendary actions")
			}
		})
		for _, line := range monsters165Entries[want.name] {
			if !strings.HasPrefix(line, "legendary-uses|") {
				entries++
			}
		}
	}
	// Bold-italic entry headings counted on PDF pp. 298–302, 307 and 356.
	if entries != 48 {
		t.Fatalf("reviewed %d entries, want 48", entries)
	}
}

// The decision inventory owns the exact family list and cited pages.
func TestMonsters165MatchesDecisionInventory(t *testing.T) {
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
		if r.CatalogueIssue == 165 {
			if r.Status != "unsupported" || !strings.HasPrefix(r.ID, "srd-5.2.1:content/"+r.Family+"/") {
				t.Fatalf("%s: inventory family or status changed", r.Name)
			}
			want = append(want, identity{r.ID, r.Name, r.SourcePage})
			families = append(families, r.Family)
		}
	}
	for _, m := range loadMonsters165(t).Monsters {
		got = append(got, identity{m.ID, m.Name, m.SourcePages[0]})
	}
	if len(want) != 12 || !reflect.DeepEqual(got, want) {
		t.Fatalf("batch differs from inventory:\n got %v\nwant %v", got, want)
	}
	// The inventory files the Animals appendix blocks under animals, the rest under monsters.
	if !slices.Equal(families, strings.Fields("monsters monsters monsters monsters animals monsters animals monsters monsters monsters monsters monsters")) {
		t.Fatalf("inventory families changed: %v", families)
	}
	// Issue #165 cites the page each block starts on, in this order.
	names := []string{"Hobgoblin Captain", "Hobgoblin Warrior", "Homunculus", "Horned Devil", "Hunter Shark", "Hydra", "Hyena", "Ice Devil", "Ice Mephit", "Imp", "Incubus", "Invisible Stalker"}
	pages := []int{299, 298, 299, 299, 356, 300, 356, 300, 307, 300, 301, 301}
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
func TestMonsters165DerivedSourceConsistency(t *testing.T) {
	pb := map[string]int{"3": 2, "1/2": 2, "0": 2, "11": 4, "2": 2, "8": 3, "14": 5, "1": 2, "4": 2, "6": 3}
	skillAbility := map[string]string{"deception": "cha", "insight": "wis", "perception": "wis", "persuasion": "cha", "stealth": "dex"}
	hitDie := map[string]int{"tiny": 4, "small": 6, "medium": 8, "large": 10, "huge": 12}
	// Attacks that do not use the better of Strength and Dexterity.
	attackAbility := map[string]string{"Horned Devil/hurl-flame": "cha", "Hyena/bite": "str", "Incubus/restless-touch": "cha"}
	for _, m := range loadMonsters165(t).Monsters {
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
func TestMonsters165DependencyKeys(t *testing.T) {
	for _, m := range loadMonsters165(t).Monsters {
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
func TestMonsters165RemainUnavailable(t *testing.T) {
	bandit := map[string]bool{"monster/stat-block": true, "action/melee-attack": true, "action/ranged-attack": true, "effect/slashing-damage": true, "effect/piercing-damage": true}
	for _, m := range loadMonsters165(t).Monsters {
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

func TestMonsters165MalformedFixtures(t *testing.T) {
	tests := []struct {
		name, record string
		mutate       func(*Monster)
	}{
		{"missing languages list", "Hobgoblin Warrior", func(m *Monster) { m.Languages = nil }},
		{"animal block in another family", "Hobgoblin Warrior", func(m *Monster) {
			m.ID = strings.Replace(strings.Replace(m.ID, "/animals/", "/spells/", 1), "/monsters/", "/spells/", 1)
		}},
		{"XP off CR", "Invisible Stalker", func(m *Monster) { m.XP++ }},
		{"save off modifier", "Invisible Stalker", func(m *Monster) {
			m.Abilities["wis"] = MonsterAbility{m.Abilities["wis"].Score, m.Abilities["wis"].Modifier, m.Abilities["wis"].Modifier + 7}
		}},
		{"modifier", "Hobgoblin Captain", func(m *Monster) {
			m.Abilities["str"] = MonsterAbility{m.Abilities["str"].Score, m.Abilities["str"].Modifier + 1, m.Abilities["str"].Modifier + 1}
		}},
		{"HP average", "Hobgoblin Captain", func(m *Monster) { m.HP.Average++ }},
		{"initiative score as printed", "Invisible Stalker", func(m *Monster) { m.Initiative.Score = 22 }},
		{"unlisted sense note", "Horned Devil", func(m *Monster) { m.Senses.Notes["darkvision"] = "impeded by magical Darkness" }},
		{"recharge", "Ice Devil", func(m *Monster) { m.Actions[3].Uses.Recharge = 7 }},
		{"daily uses with recharge", "Ice Mephit", func(m *Monster) { m.Actions[1].Uses.Recharge = 6 }},
		{"save DC off entry text", "Homunculus", func(m *Monster) { m.Actions[0].Save.DC = 13 }},
		{"creature size", "Horned Devil", func(m *Monster) { m.Actions[3].Target.Size = 5 }},
		{"cone width", "Ice Mephit", func(m *Monster) { m.Actions[2].Target.Width = 5 }},
		{"trait dice", "Ice Mephit", func(m *Monster) { m.Traits[0].Dice.Average = 6 }},
		{"unlisted spell", "Incubus", func(m *Monster) { m.Actions[2].Spellcasting.AtWill[0] = "Wish" }},
		{"untyped spellcasting", "Incubus", func(m *Monster) { m.Actions[2].Spellcasting = nil }},
		{"fixed damage with modifier", "Homunculus", func(m *Monster) { m.Actions[0].Attack.Damage.Modifier = 1 }},
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
		c := loadMonsters165(t)
		c.Monsters[1].ID = c.Monsters[2].ID
		rejected(t, c, true)
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := loadMonsters165(t)
			tt.mutate(record(t, &c, tt.record))
			// An empty optional list is omitted when encoded, so it cannot reach a decoder.
			rejected(t, c, !strings.HasPrefix(tt.name, "empty "))
		})
	}
	// Typed numbers cannot be omitted, nulled, duplicated or joined by unknown fields.
	raw, err := os.ReadFile("testdata/monsters_165.json")
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
		{`"per_day":0,`, ``},
		{`"dc":17`, `"dc":17,"dc":18`},
		{`"width":0,`, ``},
		{`"save_dc":15,`, ``},
		{`"attack_bonus":0`, `"attack_bonus":0,"attack_bonus":5`},
		{`"darkvision":150`, `"darkvision":null`},
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
// to 14 without entering the embedded runtime catalogue. Batch 04 repeats the
// Bandit revision #143 pinned, so composition keeps that one copy.
func TestMonsters165ComposesWithEarlierBatches(t *testing.T) {
	base, err := LoadMonsters()
	if err != nil {
		t.Fatal(err)
	}
	batch := loadMonsters165(t)
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
	for _, load := range []func(*testing.T) MonsterCatalogue{loadMonsters158, loadMonsters159, loadMonsters160, loadMonsters161, loadMonsters162, loadMonsters163, loadMonsters164} {
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
	for _, id := range []string{"bite", "claw", "greatsword", "hurl-flame", "longbow", "longsword", "restless-touch", "searing-fork", "sting", "tail", "wind-swipe"} {
		if _, err := ResolveBanditAttack(BanditProfile{Version: version, ActionID: id}, 12, nil); err == nil {
			t.Fatalf("catalogued action %s became executable", id)
		}
	}
}
