package srd521

import (
	"strings"
	"testing"
)

// faces returns 1-based die faces in order.
type zeroFaces []int

func (f *zeroFaces) Intn(int) int { v := (*f)[0]; *f = (*f)[1:]; return v - 1 }

func TestResolvePlayerDamage(t *testing.T) {
	for name, tc := range map[string]struct {
		hp, damage int
		critical   bool
		saves      DeathSaves
		hpAfter    int
		event      string
		after      DeathSaves
	}{
		"above 0":                       {hp: 10, damage: 9, hpAfter: 1},
		"no damage at 0":                {hp: 0, damage: 0},
		"drops to 0":                    {hp: 10, damage: 21, event: ZeroHPDropped},
		"massive damage at max":         {hp: 6, damage: 18, event: ZeroHPMassiveDamage, after: DeathSaves{Dead: true}},
		"one below massive":             {hp: 6, damage: 17, event: ZeroHPDropped},
		"hit at 0":                      {hp: 0, damage: 3, saves: DeathSaves{Successes: 2}, event: ZeroHPDamageAtZero, after: DeathSaves{Successes: 2, Failures: 1}},
		"critical at 0":                 {hp: 0, damage: 3, critical: true, event: ZeroHPDamageAtZero, after: DeathSaves{Failures: 2}},
		"third failure":                 {hp: 0, damage: 1, saves: DeathSaves{Failures: 2}, event: ZeroHPDamageAtZero, after: DeathSaves{Failures: 3, Dead: true}},
		"critical past three":           {hp: 0, damage: 1, critical: true, saves: DeathSaves{Failures: 2}, event: ZeroHPDamageAtZero, after: DeathSaves{Failures: 3, Dead: true}},
		"max damage at 0":               {hp: 0, damage: 12, event: ZeroHPDamageAtZero, after: DeathSaves{Dead: true}},
		"stable hit stops being stable": {hp: 0, damage: 2, saves: DeathSaves{Stable: true}, event: ZeroHPDamageAtZero, after: DeathSaves{Failures: 1}},
	} {
		hp, r, err := ResolvePlayerDamage(tc.hp, 12, tc.damage, tc.critical, tc.saves)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if hp != tc.hpAfter || (r == nil) != (tc.event == "") {
			t.Fatalf("%s: hp %d receipt %+v", name, hp, r)
		}
		if r != nil && (r.Event != tc.event || r.After != tc.after || r.Reason == "" || len(r.SourceIDs) == 0) {
			t.Errorf("%s: receipt %+v", name, r)
		}
	}
	if _, _, err := ResolvePlayerDamage(0, 12, 1, false, DeathSaves{Dead: true}); err == nil {
		t.Error("damage to a dead character accepted")
	}
	if _, _, err := ResolvePlayerDamage(13, 12, 1, false, DeathSaves{}); err == nil {
		t.Error("HP above maximum accepted")
	}
}

func TestResolveDeathSave(t *testing.T) {
	for name, tc := range map[string]struct {
		die     int
		saves   DeathSaves
		outcome string
		after   DeathSaves
		hp      int
	}{
		"10 succeeds":     {die: 10, outcome: "success", after: DeathSaves{Successes: 1}},
		"9 fails":         {die: 9, saves: DeathSaves{Successes: 1}, outcome: "failure", after: DeathSaves{Successes: 1, Failures: 1}},
		"1 is two":        {die: 1, outcome: "two-failures", after: DeathSaves{Failures: 2}},
		"1 kills at one":  {die: 1, saves: DeathSaves{Failures: 1}, outcome: "dead", after: DeathSaves{Failures: 3, Dead: true}},
		"third failure":   {die: 2, saves: DeathSaves{Successes: 2, Failures: 2}, outcome: "dead", after: DeathSaves{Successes: 2, Failures: 3, Dead: true}},
		"third success":   {die: 19, saves: DeathSaves{Successes: 2, Failures: 2}, outcome: "stable", after: DeathSaves{Stable: true}},
		"20 regains 1 HP": {die: 20, saves: DeathSaves{Successes: 1, Failures: 2}, outcome: "regain-1-hp", hp: 1},
	} {
		dice := zeroFaces{tc.die}
		r, err := ResolveDeathSave(tc.saves, 12, &dice)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if r.Die != tc.die || r.Outcome != tc.outcome || r.After != tc.after || r.HPAfter != tc.hp || r.Regained() != (tc.hp == 1) || !strings.Contains(r.Reason, "Death Saving Throw") {
			t.Errorf("%s: %+v", name, r)
		}
	}
	for _, s := range []DeathSaves{{Stable: true}, {Dead: true}, {Failures: 3}} {
		dice := zeroFaces{10}
		if _, err := ResolveDeathSave(s, 12, &dice); err == nil {
			t.Errorf("save rolled for %+v", s)
		}
	}
	bad := zeroFaces{21}
	if _, err := ResolveDeathSave(DeathSaves{}, 12, &bad); err == nil {
		t.Error("invalid d20 accepted")
	}
}

func TestResolveStableRecovery(t *testing.T) {
	dice := zeroFaces{3}
	r, err := ResolveStableRecovery(DeathSaves{Stable: true}, 12, &dice)
	if err != nil || r.RecoveryHours != 3 || r.HPAfter != 1 || r.After != (DeathSaves{}) || !r.Regained() {
		t.Fatalf("recovery %+v %v", r, err)
	}
	dice = zeroFaces{3}
	if _, err := ResolveStableRecovery(DeathSaves{}, 12, &dice); err == nil {
		t.Error("recovery without Stable accepted")
	}
}

func TestDecodeDeathSaves(t *testing.T) {
	for stats, ok := range map[string]bool{
		`{}`:                          true,
		`{"srd521_death_saves":null}`: true,
		`{"srd521_death_saves":{"successes":2,"failures":2}}`: true,
		`{"srd521_death_saves":{"stable":true}}`:              true,
		`{"srd521_death_saves":{"failures":3,"dead":true}}`:   true,
		`{"srd521_death_saves":{"failures":3}}`:               false,
		`{"srd521_death_saves":{"successes":3}}`:              false,
		`{"srd521_death_saves":{"stable":true,"failures":1}}`: false,
		`{"srd521_death_saves":{"dead":true,"stable":true}}`:  false,
		`{"srd521_death_saves":{"failures":-1}}`:              false,
		`{"srd521_death_saves":{"revived":true}}`:             false,
		`{"srd521_death_saves":"dead"}`:                       false,
	} {
		if _, err := DecodeDeathSaves([]byte(stats)); (err == nil) != ok {
			t.Errorf("%s: err %v", stats, err)
		}
	}
}

func TestAttackAgainstUnconsciousTarget(t *testing.T) {
	target := IncomingAttack{Advantage: []string{"target is Unconscious", "target is Prone within 5 feet"}, CriticalOnHit: "target is Unconscious within 5 feet"}
	// Bandit Scimitar: d20s 3 and 8 keep 8, +3 = 11 hits AC 11; the automatic
	// Critical Hit doubles the d6: 4 + 2, +1.
	dice := zeroFaces{3, 8, 4, 2}
	r, err := ResolveBanditAttackAgainst(BanditProfile{Version: banditTestVersion(t), ActionID: "scimitar"}, 11, target, &dice)
	if err != nil {
		t.Fatal(err)
	}
	if r.Die != 8 || r.Outcome != "critical-hit" || r.AutoCritical == "" || len(r.DamageDice) != 2 || r.Damage != 7 || !strings.Contains(strings.Join(r.SourceIDs, " "), "rules-glossary/unconscious") {
		t.Fatalf("receipt %+v", r)
	}
	// A miss stays a miss.
	dice = zeroFaces{1, 2}
	if r, err := ResolveBanditAttackAgainst(BanditProfile{Version: banditTestVersion(t), ActionID: "scimitar"}, 20, target, &dice); err != nil || r.Outcome != "miss" || r.AutoCritical != "" || r.Damage != 0 {
		t.Fatalf("miss %+v %v", r, err)
	}
	// Unconscious Advantage and Prone-at-range Disadvantage cancel: one d20,
	// and no automatic Critical Hit beyond 5 feet.
	ranged := IncomingAttack{Advantage: []string{"target is Unconscious"}, Disadvantage: []string{"target is Prone beyond 5 feet"}}
	dice = zeroFaces{15, 5}
	r, err = ResolveBanditAttackAgainst(BanditProfile{Version: banditTestVersion(t), ActionID: "light-crossbow"}, 11, ranged, &dice)
	if err != nil || len(r.D20s) != 0 || r.Outcome != "hit" || r.AutoCritical != "" || len(r.DamageDice) != 1 {
		t.Fatalf("ranged %+v %v", r, err)
	}
	// The fixture Unarmed Strike has no dice to double.
	dice = zeroFaces{12}
	if r, err := ResolveFixtureUnarmedAgainst(10, IncomingAttack{CriticalOnHit: "target is Unconscious within 5 feet"}, &dice); err != nil || r.Outcome != "critical-hit" || r.Damage != 1 {
		t.Fatalf("unarmed %+v %v", r, err)
	}
}

func banditTestVersion(t *testing.T) string {
	t.Helper()
	_, version, err := BanditRecord()
	if err != nil {
		t.Fatal(err)
	}
	return version
}
