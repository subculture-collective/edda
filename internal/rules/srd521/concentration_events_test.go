package srd521

import (
	"strings"
	"testing"
)

// queueDice returns queued Intn values; the Cleric fixture's Constitution
// save total is value+1+2.
type queueDice struct {
	values []int
	calls  int
}

func (d *queueDice) Intn(int) int {
	d.calls++
	v := d.values[0]
	d.values = d.values[1:]
	return v
}

func concentratingCleric(t *testing.T, hp int) (*Build, *CharacterSheet) {
	t.Helper()
	b, sheet := casterBuild(t)
	sheet.HP.Value = hp
	b.Casting.Concentration = &Concentration{ID: "cast", Spell: "shield-of-faith", RecipientID: "ally", StartedAt: 28800, ExpiresAt: 29400}
	return b, sheet
}

func damage(source string, amount, hpAfter int) ConcentrationEvent {
	return ConcentrationEvent{Kind: ConcentrationDamageSource, SourceID: source, Damage: amount, HPAfter: hpAfter}
}

// The exact #353 fixture: two independent damage sources, then a trusted
// Unconscious transition at positive HP.
func TestConcentrationEventsTwoSourcesThenUnconscious(t *testing.T) {
	b, sheet := concentratingCleric(t, 10)
	dice := &queueDice{values: []int{7, 7}}
	events := []ConcentrationEvent{damage("attack-1", 3, 7), damage("attack-2", 2, 5), {Kind: ConcentrationUnconscious}}
	receipts, err := ResolveConcentrationEvents(b, sheet, events, 28806, dice)
	if err != nil {
		t.Fatal(err)
	}
	if dice.calls != 2 || len(receipts) != 3 || b.Casting.Concentration != nil {
		t.Fatalf("calls=%d receipts=%+v", dice.calls, receipts)
	}
	for i, want := range []struct{ before, after int }{{10, 7}, {7, 5}} {
		r := receipts[i]
		if r.Check == nil || r.Check.Outcome != "success" || *r.Check.DC != 10 || r.HPBefore != want.before || r.HPAfter != want.after || r.After == nil || r.Mechanic != ConcentrationDamageSource {
			t.Fatalf("receipt %d: %+v", i, r)
		}
	}
	end := receipts[2]
	if end.Check != nil || end.After != nil || end.Before.ID != "cast" || end.Kind != "concentration-ended" || end.Mechanic != ConcentrationUnconscious || end.HPAfter != 5 || !strings.Contains(end.Reason, "Unconscious") {
		t.Fatalf("unconscious receipt: %+v", end)
	}
}

func TestConcentrationEventsSavesAndEndings(t *testing.T) {
	for _, tc := range []struct {
		name    string
		events  []ConcentrationEvent
		rolls   []int
		dcs     []int
		kept    bool
		reasons int
	}{
		{"success", []ConcentrationEvent{damage("a", 4, 6)}, []int{7}, []int{10}, true, 1},
		{"failure", []ConcentrationEvent{damage("a", 4, 6)}, []int{6}, []int{10}, false, 1},
		{"odd damage rounds down", []ConcentrationEvent{damage("a", 23, 1)}, []int{7}, []int{11}, false, 1},
		{"cap 30", []ConcentrationEvent{damage("a", 61, 1)}, []int{19}, []int{30}, false, 1},
		{"one instance with two parts saves once", []ConcentrationEvent{damage("a", 11, 9), damage("a", 11, 1)}, []int{8}, []int{11}, true, 1},
		{"two sources save twice", []ConcentrationEvent{damage("a", 4, 6), damage("b", 4, 2)}, []int{7, 7}, []int{10, 10}, true, 2},
		{"zero damage", []ConcentrationEvent{damage("a", 0, 10)}, nil, nil, true, 0},
		{"zero HP ends without a save", []ConcentrationEvent{damage("a", 12, 0)}, nil, nil, false, 1},
		{"incapacitated", []ConcentrationEvent{{Kind: ConcentrationIncapacitated}}, nil, nil, false, 1},
		{"dead", []ConcentrationEvent{{Kind: ConcentrationDead}}, nil, nil, false, 1},
		{"no save after a failed save", []ConcentrationEvent{damage("a", 4, 6), damage("b", 4, 2), {Kind: ConcentrationDead}}, []int{0}, []int{10}, false, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b, sheet := concentratingCleric(t, 10)
			dice := &queueDice{values: tc.rolls}
			receipts, err := ResolveConcentrationEvents(b, sheet, tc.events, 28806, dice)
			if err != nil {
				t.Fatal(err)
			}
			if dice.calls != len(tc.rolls) || len(receipts) != tc.reasons || (b.Casting.Concentration != nil) != tc.kept {
				t.Fatalf("calls=%d receipts=%d concentration=%+v", dice.calls, len(receipts), b.Casting.Concentration)
			}
			saves := 0
			for _, r := range receipts {
				if r.Check != nil {
					if *r.Check.DC != tc.dcs[saves] || r.Check.Ability != "constitution" {
						t.Fatalf("check %d: %+v", saves, r.Check)
					}
					saves++
				}
			}
			if saves != len(tc.dcs) {
				t.Fatalf("saves=%d want %d", saves, len(tc.dcs))
			}
		})
	}
}

// Ended Concentration stays ended: later damage, waking or healing never
// restore it, and an event with no Concentration rolls nothing.
func TestConcentrationEventsNoResurrection(t *testing.T) {
	b, sheet := concentratingCleric(t, 10)
	dice := &queueDice{}
	if _, err := ResolveConcentrationEvents(b, sheet, []ConcentrationEvent{{Kind: ConcentrationUnconscious}}, 28806, dice); err != nil {
		t.Fatal(err)
	}
	sheet.HP.Value = 10 // healed and awake
	receipts, err := ResolveConcentrationEvents(b, sheet, []ConcentrationEvent{damage("a", 4, 6), {Kind: ConcentrationIncapacitated}}, 28812, dice)
	if err != nil || len(receipts) != 0 || dice.calls != 0 || b.Casting.Concentration != nil {
		t.Fatalf("receipts=%+v calls=%d err=%v", receipts, dice.calls, err)
	}
}

// Unknown or malformed events fail closed before anything changes.
func TestConcentrationEventsFailClosed(t *testing.T) {
	for name, events := range map[string][]ConcentrationEvent{
		"unknown key":          {damage("a", 4, 6), {Kind: "concentration.end.asleep"}},
		"missing source":       {damage("", 4, 6)},
		"negative damage":      {damage("a", -1, 10)},
		"repeated source":      {damage("a", 2, 8), damage("b", 2, 6), damage("a", 2, 4)},
		"damage on condition":  {{Kind: ConcentrationUnconscious, Damage: 3}},
		"damage that heals":    {damage("a", 4, 6), damage("b", 2, 12)},
		"client sleep request": {{Kind: "sleep"}},
	} {
		t.Run(name, func(t *testing.T) {
			b, sheet := concentratingCleric(t, 10)
			dice := &queueDice{values: []int{0, 0}}
			if _, err := ResolveConcentrationEvents(b, sheet, events, 28806, dice); err == nil {
				t.Fatal("accepted")
			}
			if dice.calls != 0 {
				t.Fatalf("rolled %d", dice.calls)
			}
			if b.Casting.Concentration == nil {
				t.Fatal("changed concentration before validation finished")
			}
		})
	}
}
