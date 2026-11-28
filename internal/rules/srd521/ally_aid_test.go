package srd521

import (
	"strings"
	"testing"
)

func TestResolveAllyStabilize(t *testing.T) {
	for name, tc := range map[string]struct {
		die, modifier int
		saves         DeathSaves
		outcome       string
		after         DeathSaves
	}{
		"exactly DC 10":           {die: 7, modifier: 3, saves: DeathSaves{Successes: 1, Failures: 2}, outcome: "stable", after: DeathSaves{Stable: true}},
		"one short":               {die: 6, modifier: 3, saves: DeathSaves{Failures: 1}, outcome: "failure", after: DeathSaves{Failures: 1}},
		"natural 20 is not auto":  {die: 20, modifier: -11, outcome: "failure"},
		"natural 1 is not auto":   {die: 1, modifier: 9, outcome: "stable", after: DeathSaves{Stable: true}},
		"negative modifier fails": {die: 12, modifier: -3, outcome: "failure"},
	} {
		faces := zeroFaces{tc.die}
		r, err := ResolveAllyStabilize(tc.saves, 12, tc.modifier, &faces)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if r.Event != ZeroHPAllyStabilize || r.Die != tc.die || r.Total != tc.die+tc.modifier || r.DC != 10 || r.Modifier == nil || *r.Modifier != tc.modifier || r.Outcome != tc.outcome || r.After != tc.after || r.Before != tc.saves || r.HPAfter != 0 {
			t.Fatalf("%s: %+v", name, r)
		}
		if !strings.Contains(r.Reason, "DC 10") || (tc.outcome == "stable") != strings.Contains(r.Reason, "Stable at 0 HP") {
			t.Fatalf("%s: reason %q", name, r.Reason)
		}
	}
	for name, saves := range map[string]DeathSaves{"stable": {Stable: true}, "dead": {Failures: 3, Dead: true}, "impossible": {Failures: 4}} {
		faces := zeroFaces{15}
		if _, err := ResolveAllyStabilize(saves, 12, 3, &faces); err == nil || len(faces) != 1 {
			t.Fatalf("%s: rolled or accepted: %v", name, err)
		}
	}
	if _, err := ResolveAllyStabilize(DeathSaves{}, 12, 3, nil); err == nil {
		t.Fatal("nil dice accepted")
	}
}

func TestResolveAdministeredPotion(t *testing.T) {
	faces := zeroFaces{3, 4}
	r, err := ResolveAdministeredPotion(0, 12, &faces)
	if err != nil || r.HPBefore != 0 || r.RolledHealing != 9 || r.HPAfter != 9 || r.Healing != 9 {
		t.Fatalf("from 0 HP: %+v %v", r, err)
	}
	faces = zeroFaces{4, 4}
	if r, err := ResolveAdministeredPotion(0, 5, &faces); err != nil || r.HPAfter != 5 || r.Healing != 5 {
		t.Fatalf("capped: %+v %v", r, err)
	}
	for name, hp := range map[string]int{"negative": -1, "over max": 13} {
		if _, err := ResolveAdministeredPotion(hp, 12, &zeroFaces{1, 1}); err == nil {
			t.Fatalf("%s accepted", name)
		}
	}
	// Drinking one yourself still needs a conscious drinker.
	if _, err := ResolveHealingPotion(0, 12, &zeroFaces{1, 1}); err == nil {
		t.Fatal("self-use at 0 HP accepted")
	}
}
