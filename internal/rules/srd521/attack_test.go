package srd521

import (
	"os"
	"slices"
	"testing"
)

type attackDice struct {
	values []int
	calls  int
}

func (d *attackDice) Intn(int) int { v := d.values[d.calls]; d.calls++; return v }
func TestLongswordFixture(t *testing.T) {
	stats, err := os.ReadFile("testdata/fighter1.json")
	if err != nil {
		t.Fatal(err)
	}
	sheet, err := Project(stats, 1, 12, 12, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name    string
		rolls   []int
		ac      int
		outcome string
		damage  int
	}{
		{"hit", []int{11, 4}, 10, "hit", 8},
		{"miss", []int{1}, 20, "miss", 0},
		{"natural one", []int{0}, 1, "miss", 0},
		{"critical ignores armor", []int{19, 4, 6}, 99, "critical-hit", 15},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dice := &attackDice{values: tc.rolls}
			r, err := ResolveLongsword(sheet, tc.ac, dice)
			if err != nil {
				t.Fatal(err)
			}
			if r.Outcome != tc.outcome || r.Damage != tc.damage || r.Modifier.Value != 5 || dice.calls != len(tc.rolls) {
				t.Fatalf("receipt=%+v calls=%d", r, dice.calls)
			}
		})
	}
	for _, rolls := range [][]int{{20}, {11, 8}} {
		if _, err := ResolveLongsword(sheet, 10, &attackDice{values: rolls}); err == nil {
			t.Fatal("invalid dice accepted")
		}
	}
}

func TestLongswordAdvantageDisadvantage(t *testing.T) {
	stats, err := os.ReadFile("testdata/fighter1.json")
	if err != nil {
		t.Fatal(err)
	}
	sheet, err := Project(stats, 1, 12, 12, 0)
	if err != nil {
		t.Fatal(err)
	}
	prone := []string{"attacker is Prone"}
	for _, tc := range []struct {
		name                    string
		advantage, disadvantage []string
		rolls, d20s             []int
		die                     int
		outcome, note           string
	}{
		// A natural 20 on one die is discarded: the lower die is kept.
		{name: "disadvantage keeps lower", disadvantage: prone, rolls: []int{19, 9, 3}, d20s: []int{20, 10}, die: 10, outcome: "hit", note: "Disadvantage (attacker is Prone): rolled 20 and 10, kept 10."},
		{name: "disadvantage both low", disadvantage: prone, rolls: []int{4, 2}, d20s: []int{5, 3}, die: 3, outcome: "miss", note: "Disadvantage (attacker is Prone): rolled 5 and 3, kept 3."},
		{name: "advantage keeps higher", advantage: []string{"x"}, rolls: []int{2, 9, 3}, d20s: []int{3, 10}, die: 10, outcome: "hit", note: "Advantage (x): rolled 3 and 10."},
		// Any Advantage cancels any Disadvantage: one d20, even two reasons against one.
		{name: "cancel", advantage: []string{"x"}, disadvantage: []string{"attacker is Prone", "y"}, rolls: []int{19, 3, 4}, die: 20, outcome: "critical-hit", note: "Advantage (x) and Disadvantage (attacker is Prone; y) cancel: one d20."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dice := &attackDice{values: tc.rolls}
			r, err := ResolveLongswordRoll(sheet, 15, tc.advantage, tc.disadvantage, dice)
			if err != nil {
				t.Fatal(err)
			}
			if r.Die != tc.die || r.Outcome != tc.outcome || !slices.Equal(r.D20s, tc.d20s) || !slices.Equal(r.Advantage, tc.advantage) || !slices.Equal(r.Disadvantage, tc.disadvantage) || dice.calls != len(tc.rolls) || r.RollNote() != tc.note {
				t.Fatalf("receipt=%+v calls=%d note=%q", r, dice.calls, r.RollNote())
			}
			if !slices.Contains(r.SourceIDs, prefix+"d20-tests/advantage-disadvantage") {
				t.Fatalf("sources=%v", r.SourceIDs)
			}
		})
	}
	if _, err := ResolveLongswordRoll(sheet, 15, nil, prone, &attackDice{values: []int{3, 20}}); err == nil {
		t.Fatal("invalid second d20 accepted")
	}
	if r, err := ResolveLongswordRoll(sheet, 15, nil, nil, &attackDice{values: []int{3}}); err != nil || r.RollNote() != "" || r.D20s != nil {
		t.Fatalf("ordinary roll=%+v %v", r, err)
	}
}
