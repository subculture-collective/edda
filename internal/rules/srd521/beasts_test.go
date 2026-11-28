package srd521

import (
	"bytes"
	"slices"
	"strings"
	"testing"
)

func beastProfile(t *testing.T, profile, size string) BeastProfile {
	t.Helper()
	_, version, err := BeastRecord(profile)
	if err != nil {
		t.Fatal(err)
	}
	return BeastProfile{Version: version, ActionID: beastFixtures[profile].action, TargetSize: size}
}

func TestBeastRecordsStayOutsideCatalogueAvailability(t *testing.T) {
	for profile, id := range map[string]string{"wolf-v1": WolfID, "brown-bear-v1": BrownBearID} {
		m, version, err := BeastRecord(profile)
		if err != nil || m.ID != id || !strings.HasPrefix(version, profile+":") {
			t.Fatalf("%s: %v %s", profile, err, version)
		}
		if m.Available() {
			t.Fatal("bounded runtime must not promote whole catalogue support")
		}
	}
	if _, _, err := BeastRecord("bandit-v1"); err == nil {
		t.Fatal("unknown profile loaded")
	}
	// Every other record in both batch files stays refused by its capability set.
	for _, data := range [][]byte{monsters155JSON, monsters177JSON} {
		c, err := DecodeMonsters(bytes.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range c.Monsters {
			if m.Available() {
				t.Fatalf("%s became available", m.ID)
			}
		}
	}
}

// Any unknown key, extra action, trait or changed size gate refuses the block.
func TestBeastRecordsFailClosed(t *testing.T) {
	gate := func(id, text string) func(*Monster) {
		return func(m *Monster) {
			for i := range m.Actions {
				for j := range m.Actions[i].Effects {
					if m.Actions[i].ID == id && m.Actions[i].Effects[j].Key == "requirement/target-size" {
						m.Actions[i].Effects[j].Text = text
					}
				}
			}
		}
	}
	for profile, edits := range map[string]map[string]func(*Monster){
		"wolf-v1": {
			"unknown key":   func(m *Monster) { m.RequiredMechanics = append(m.RequiredMechanics, "trait/unknown") },
			"size gate":     gate("bite", "If the target is a Small or smaller creature"),
			"trait require": func(m *Monster) { m.Traits[0].Requires = m.Traits[0].Requires[:1] },
			"extra trait":   func(m *Monster) { m.Traits = append(m.Traits, m.Traits[0]) },
		},
		"brown-bear-v1": {
			"unknown key":  func(m *Monster) { m.RequiredMechanics = append(m.RequiredMechanics, "action/unknown") },
			"size gate":    gate("claw", "If the target is a Huge or smaller creature, it has the Prone condition."),
			"reordered":    func(m *Monster) { m.Actions[1], m.Actions[2] = m.Actions[2], m.Actions[1] },
			"bonus action": func(m *Monster) { m.BonusActions = []MonsterAction{m.Actions[1]} },
		},
	} {
		f := beastFixtures[profile]
		implemented := map[string]bool{}
		for _, key := range f.implemented {
			implemented[key] = true
		}
		for name, edit := range edits {
			c, err := DecodeMonsters(bytes.NewReader(f.data))
			if err != nil {
				t.Fatal(err)
			}
			for _, m := range c.Monsters {
				if m.ID != f.id {
					continue
				}
				if !m.availableWith(implemented) || !f.shape(m) {
					t.Fatal("pinned record refused")
				}
				edit(&m)
				if m.availableWith(implemented) && f.shape(m) {
					t.Fatalf("%s accepted %s", profile, name)
				}
			}
		}
	}
}

func TestWolfBite(t *testing.T) {
	pack := []string{"Pack Tactics"}
	for _, tc := range []struct {
		name, size, outcome, condition string
		rolls, d20s                    []int
		advantage                      []string
		prone                          bool
		damage                         int
	}{
		{name: "hit", size: "medium", rolls: []int{15, 2}, outcome: "hit", damage: 5, condition: "prone-applied"},
		{name: "miss", size: "medium", rolls: []int{1}, outcome: "miss"},
		{name: "advantage keeps higher", size: "small", rolls: []int{1, 15, 2}, advantage: pack, d20s: []int{2, 16}, outcome: "hit", damage: 5, condition: "prone-applied"},
		{name: "advantage both low", size: "medium", rolls: []int{1, 2}, advantage: pack, d20s: []int{2, 3}, outcome: "miss"},
		{name: "critical", size: "medium", rolls: []int{19, 0, 5}, outcome: "critical-hit", damage: 9, condition: "prone-applied"},
		{name: "already prone", size: "medium", rolls: []int{15, 2}, prone: true, outcome: "hit", damage: 5, condition: "already-prone"},
		{name: "large target exempt", size: "large", rolls: []int{15, 2}, outcome: "hit", damage: 5, condition: "size-exempt"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := &banditTestDice{values: tc.rolls}
			r, err := ResolveBeastAttack(beastProfile(t, "wolf-v1", tc.size), BeastAttack{ActionID: "bite", ArmorClass: 12, Advantage: tc.advantage, TargetProne: tc.prone}, d)
			if err != nil || r.Outcome != tc.outcome || r.Damage != tc.damage || r.Condition != tc.condition || !slices.Equal(r.D20s, tc.d20s) || r.Modifier.Value != 4 || r.DamageType != "piercing" || r.Weapon != "Bite" {
				t.Fatalf("%+v %v", r, err)
			}
			if d.calls != len(tc.rolls) || r.SourceIDs[0] != WolfID {
				t.Fatalf("calls=%d sources=%v", d.calls, r.SourceIDs)
			}
		})
	}
	valid := beastProfile(t, "wolf-v1", "medium")
	for name, p := range map[string]BeastProfile{
		"stale":        {Version: "wolf-v1:stale", ActionID: "bite", TargetSize: "medium"},
		"wrong action": {Version: valid.Version, ActionID: "claw", TargetSize: "medium"},
		"no size":      {Version: valid.Version, ActionID: "bite"},
		"bad size":     {Version: valid.Version, ActionID: "bite", TargetSize: "colossal"},
		"bear version": {Version: beastProfile(t, "brown-bear-v1", "medium").Version, ActionID: "bite", TargetSize: "medium"},
	} {
		d := &banditTestDice{values: []int{15, 2}}
		if _, err := ResolveBeastAttack(p, BeastAttack{ActionID: "bite", ArmorClass: 12}, d); err == nil || d.calls != 0 {
			t.Fatalf("%s rolled", name)
		}
	}
	for name, rolls := range map[string][]int{"bad d20": {20}, "bad second d20": {3, -1}, "bad damage": {15, 6}} {
		if _, err := ResolveBeastAttack(valid, BeastAttack{ActionID: "bite", ArmorClass: 12, Advantage: pack}, &banditTestDice{values: rolls}); err == nil {
			t.Fatalf("%s accepted", name)
		}
	}
}

func TestBrownBearAttacks(t *testing.T) {
	p := beastProfile(t, "brown-bear-v1", "large")
	m, attacks, err := p.Record()
	if err != nil || len(attacks) != 2 || attacks[0].ID != "bite" || attacks[1].ID != "claw" || m.AC != 11 {
		t.Fatalf("%v %+v", err, attacks)
	}
	bite, err := ResolveBeastAttack(p, BeastAttack{ActionID: "bite", ArmorClass: 12}, &banditTestDice{values: []int{15, 3}})
	if err != nil || bite.Damage != 7 || bite.DamageType != "piercing" || bite.Condition != "" || bite.Modifier.Value != 5 {
		t.Fatalf("%+v %v", bite, err)
	}
	for size, condition := range map[string]string{"large": "prone-applied", "medium": "prone-applied", "huge": "size-exempt"} {
		claw, err := ResolveBeastAttack(beastProfile(t, "brown-bear-v1", size), BeastAttack{ActionID: "claw", ArmorClass: 12}, &banditTestDice{values: []int{15, 1}})
		if err != nil || claw.Damage != 5 || claw.DamageType != "slashing" || claw.Condition != condition {
			t.Fatalf("%s: %+v %v", size, claw, err)
		}
	}
	if _, err := ResolveBeastAttack(p, BeastAttack{ActionID: "multiattack", ArmorClass: 12}, &banditTestDice{values: []int{15}}); err == nil {
		t.Fatal("multiattack resolved as one attack roll")
	}
}
