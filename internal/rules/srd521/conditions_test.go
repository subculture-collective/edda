package srd521

import (
	"reflect"
	"strings"
	"testing"
)

func TestConditionChangeUnconsciousIncludesProneAndEndsConcentration(t *testing.T) {
	r, events, err := ResolveConditionChange(Conditions{}, ConditionChange{Condition: ConditionUnconscious, Apply: true})
	if err != nil {
		t.Fatal(err)
	}
	if want := (Conditions{Prone: true, Unconscious: true}); r.After != want || !r.Changed {
		t.Fatalf("after = %+v, changed %v", r.After, r.Changed)
	}
	if !r.After.IsIncapacitated() || r.After.Incapacitated {
		t.Fatalf("Unconscious must imply Incapacitated without a second flag: %+v", r.After)
	}
	if len(events) != 1 || events[0].Kind != ConcentrationUnconscious {
		t.Fatalf("events = %+v", events)
	}
	// Ending Unconscious leaves Prone.
	end, events, err := ResolveConditionChange(r.After, ConditionChange{Condition: ConditionUnconscious})
	if err != nil {
		t.Fatal(err)
	}
	if end.After != (Conditions{Prone: true}) || len(events) != 0 || end.After.IsIncapacitated() {
		t.Fatalf("ended = %+v events %+v", end.After, events)
	}
}

func TestConditionChangeIncapacitatedSources(t *testing.T) {
	r, events, err := ResolveConditionChange(Conditions{}, ConditionChange{Condition: ConditionIncapacitated, Apply: true})
	if err != nil || len(events) != 1 || events[0].Kind != ConcentrationIncapacitated || r.After != (Conditions{Incapacitated: true}) {
		t.Fatalf("receipt %+v events %+v err %v", r, events, err)
	}
	// Becoming Unconscious while already Incapacitated raises no second ending.
	r, events, err = ResolveConditionChange(r.After, ConditionChange{Condition: ConditionUnconscious, Apply: true})
	if err != nil || len(events) != 0 {
		t.Fatalf("events %+v err %v", events, err)
	}
	// Waking keeps the other Incapacitated source.
	r, _, err = ResolveConditionChange(r.After, ConditionChange{Condition: ConditionUnconscious})
	if err != nil || r.After != (Conditions{Prone: true, Incapacitated: true}) || !r.After.IsIncapacitated() {
		t.Fatalf("after %+v err %v", r.After, err)
	}
	// Re-applying is a no-op receipt with no events.
	r, events, err = ResolveConditionChange(r.After, ConditionChange{Condition: ConditionIncapacitated, Apply: true})
	if err != nil || r.Changed || len(events) != 0 {
		t.Fatalf("repeat %+v events %+v err %v", r, events, err)
	}
}

func TestConditionChangeFailsClosed(t *testing.T) {
	for name, tc := range map[string]struct {
		before Conditions
		change ConditionChange
	}{
		"unknown condition":       {Conditions{}, ConditionChange{Condition: "paralyzed", Apply: true}},
		"empty condition":         {Conditions{}, ConditionChange{Apply: true}},
		"unconscious, not prone":  {Conditions{Unconscious: true}, ConditionChange{Condition: ConditionIncapacitated, Apply: true}},
		"forged case of the name": {Conditions{}, ConditionChange{Condition: "Unconscious", Apply: true}},
	} {
		if r, events, err := ResolveConditionChange(tc.before, tc.change); err == nil || r != nil || events != nil {
			t.Errorf("%s: receipt %+v events %+v err %v", name, r, events, err)
		}
	}
}

func TestDecodeConditions(t *testing.T) {
	for stats, want := range map[string]Conditions{
		``:                                     {},
		`{}`:                                   {},
		`{"srd521_conditions":null}`:           {},
		`{"srd521_conditions":{"prone":true}}`: {Prone: true},
		`{"srd521_conditions":{"prone":true,"unconscious":true},"custom":1}`: {Prone: true, Unconscious: true},
	} {
		got, err := DecodeConditions([]byte(stats))
		if err != nil || got != want {
			t.Errorf("%s: %+v %v", stats, got, err)
		}
	}
	for _, stats := range []string{`{"srd521_conditions":{"stunned":true}}`, `{"srd521_conditions":{"unconscious":true}}`, `{"srd521_conditions":{"prone":"yes"}}`, `not json`} {
		if _, err := DecodeConditions([]byte(stats)); err == nil {
			t.Errorf("%s: want error", stats)
		}
	}
}

func TestConditionEffects(t *testing.T) {
	if e := (Conditions{Prone: true}).Effects(); !reflect.DeepEqual(e, ConditionEffects{}) {
		t.Fatalf("Prone alone is not incapacitating: %+v", e)
	}
	inc := Conditions{Incapacitated: true}.Effects()
	if !inc.NoActions || !inc.NoSpeech || !inc.InitiativeDisadvantage || inc.SpeedZero || inc.AttackedAdvantage || len(inc.AutoFailSaves) != 0 {
		t.Fatalf("incapacitated effects %+v", inc)
	}
	unc := Conditions{Prone: true, Unconscious: true}
	e := unc.Effects()
	if !e.NoActions || !e.NoSpeech || !e.InitiativeDisadvantage || !e.DropsHeldObjects || !e.SpeedZero || !e.AttackedAdvantage || e.CriticalWithinFeet != 5 || !e.Unaware {
		t.Fatalf("unconscious effects %+v", e)
	}
	if !unc.AutoFailsSave("strength") || !unc.AutoFailsSave("dexterity") || unc.AutoFailsSave("constitution") || (Conditions{Incapacitated: true}).AutoFailsSave("strength") {
		t.Fatal("only Unconscious fails Strength and Dexterity saves")
	}
	if !strings.Contains(strings.Join(e.SourceIDs, " "), "rules-glossary/unconscious") {
		t.Fatalf("sources %v", e.SourceIDs)
	}
}

func TestConditionUnavailableReason(t *testing.T) {
	if (Conditions{Prone: true}).UnavailableReason(false) != "" {
		t.Fatal("Prone alone restricts no choice here")
	}
	inc := Conditions{Incapacitated: true}
	if inc.UnavailableReason(false) == "" || inc.UnavailableReason(true) != "" {
		t.Fatal("Incapacitated blocks actions but not movement")
	}
	unc := Conditions{Prone: true, Unconscious: true}
	if unc.UnavailableReason(false) == "" || unc.UnavailableReason(true) == "" {
		t.Fatal("Unconscious blocks actions and movement")
	}
}

func TestConcentrationEndsOnConditionEvents(t *testing.T) {
	// The condition events this resolver raises are accepted by #353's adapter.
	for _, change := range []ConditionChange{{Condition: ConditionIncapacitated, Apply: true}, {Condition: ConditionUnconscious, Apply: true}} {
		_, events, err := ResolveConditionChange(Conditions{}, change)
		if err != nil {
			t.Fatal(err)
		}
		b := &Build{}
		if _, err := ResolveConcentrationEvents(b, &CharacterSheet{}, events, 0, nil); err != nil {
			t.Fatalf("%s: %v", change.Condition, err)
		}
	}
}
