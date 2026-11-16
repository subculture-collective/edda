package srd521

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func loadSpellRecords(t *testing.T) SpellRecords {
	t.Helper()
	data, err := os.ReadFile("testdata/spell_records.json")
	if err != nil {
		t.Fatal(err)
	}
	records, err := DecodeSpellRecords(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	return records
}

func TestSpellRecordsSourceAndAcidArrow(t *testing.T) {
	r := loadSpellRecords(t)
	if len(r.Records) != 3 || r.SchemaVersion != 1 || r.Source.Attribution != Attribution {
		t.Fatal("unexpected catalogue or provenance")
	}
	s := r.Records[0]
	if s.ID != prefix+"spells/acid-arrow" || s.SourcePage != 107 || s.SourceID != s.ID || s.Name != "Acid Arrow" || s.Level != 2 || s.School != "evocation" || !reflect.DeepEqual(s.Lists, []string{"wizard"}) || s.CastingTime != "action" || s.Range != (SpellRange{"feet", 90}) || s.Target != (SpellTarget{"single", 1, "target-in-range"}) {
		t.Fatalf("wrong Acid Arrow identity/target: %+v", s)
	}
	m := s.Components.Material
	if !*s.Components.Verbal || !*s.Components.Somatic || m == nil || m.Description != "powdered rhubarb leaf" || m.CostCP != 0 || *m.Consumed || *s.Ritual || *s.Concentration || s.Duration != "instantaneous" || s.DurationSeconds != 0 {
		t.Fatal("wrong Acid Arrow casting requirements")
	}
	if len(s.Effects) != 2 || s.Effects[0].DiceCount != 4 || s.Effects[0].DieSides != 4 || s.Effects[0].DamageType != "acid" || !*s.Effects[0].HalfInitialOnly || *s.Effects[0].AddCastingModifier || s.Effects[1].DiceCount != 2 || s.Effects[1].DieSides != 4 || s.Effects[1].DamageType != "acid" || s.Effects[1].Timing != "on-hit-end-next-target-turn" || *s.Effects[1].HalfInitialOnly {
		t.Fatal("wrong hit, miss or delayed damage")
	}
	if !reflect.DeepEqual(s.Scaling, SpellScaling{"slot", 2, 1, 4, []string{"initial", "later"}}) || !reflect.DeepEqual(s.SlotLevels, []int{2, 3, 4, 5, 6, 7, 8, 9}) || s.MechanicsAvailable() {
		t.Fatal("wrong scaling or unsupported record became available")
	}
	if !reflect.DeepEqual(s.RequiredMechanics, []string{"attack.ranged-spell", "damage.acid", "damage.end-next-target-turn", "damage.miss-half-initial-only", "scaling.per-slot", "components.material"}) {
		t.Fatal("missing dependency")
	}
}

func TestSpellRecordsBoundedFixtures(t *testing.T) {
	r := loadSpellRecords(t)
	data, err := os.ReadFile("testdata/cleric1.json")
	if err != nil {
		t.Fatal(err)
	}
	sheet, err := Project(data, 1, 1, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range r.Records[1:] {
		t.Run(s.Name, func(t *testing.T) {
			if !s.MechanicsAvailable() || s.Target.Eligibility != "conscious-self" || !reflect.DeepEqual(s.SlotLevels, []int{1}) || s.Scaling.Kind != "none" {
				t.Fatal("bounded fixture changed")
			}
			for _, key := range []string{"unknown.future-key", "attack.ranged-spell", "scaling.per-slot"} {
				changed := s
				changed.RequiredMechanics = append(append([]string{}, s.RequiredMechanics...), key)
				if changed.MechanicsAvailable() {
					t.Fatalf("unsupported key %s available", key)
				}
			}
			changed := s
			changed.RequiredMechanics = nil
			if changed.MechanicsAvailable() {
				t.Fatal("missing dependencies available")
			}
		})
	}
	cure := r.Records[1]
	receipt, err := ResolveCureWounds(sheet, &spellDice{value: 0})
	if err != nil {
		t.Fatal(err)
	}
	if receipt.SlotLevel != cure.SlotLevels[0] || len(receipt.Dice) != cure.Effects[0].DiceCount || receipt.RolledHealing != 2+receipt.Modifier.Value || cure.SourcePage != 121 || cure.Range.Kind != "touch" {
		t.Fatal("Cure Wounds disagrees with existing fixture")
	}
	shield := r.Records[2]
	if shield.SourcePage != 162 || shield.Effects[0].Bonus != 2 || shield.DurationSeconds != 600 || !*shield.Concentration || shield.CastingTime != "bonus-action" || shield.Components.Material.Description != "a prayer scroll" {
		t.Fatal("Shield of Faith source mismatch")
	}
	b, shieldSheet := casterBuild(t)
	shieldReceipt, err := Lifecycle(b, shieldSheet, SpellCast, "shield-of-faith", "contract-test", 100, &spellDice{})
	if err != nil {
		t.Fatal(err)
	}
	if shieldReceipt.After.ExpiresAt-shieldReceipt.After.StartedAt != shield.DurationSeconds || shieldReceipt.SlotsBefore-shieldReceipt.SlotsAfter != 1 {
		t.Fatal("Shield of Faith disagrees with existing lifecycle")
	}
	// Existing lifecycle tests verify the focus substitution, +2 AC, exact expiry,
	// slot cost, replacement and damage concentration receipts. No resolver changes.
	before, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if before.Coverage().FullSupport || before.Coverage().PlannedSpellItemScopes != 753 {
		t.Fatal("catalogue promoted coverage")
	}
}

func TestSpellRecordsUnknownMechanicsDecodeUnavailable(t *testing.T) {
	r := loadSpellRecords(t)
	for i := range r.Records {
		r.Records[i].RequiredMechanics = append(r.Records[i].RequiredMechanics, "unknown.future-mechanic")
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeSpellRecords(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range decoded.Records {
		if s.MechanicsAvailable() {
			t.Fatalf("%s: unknown mechanic available", s.Name)
		}
	}
}

func TestSpellRecordsNegativeFixtures(t *testing.T) {
	data, err := os.ReadFile("testdata/spell_record_invalid.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixtures []struct {
		Name   string          `json:"name"`
		Record int             `json:"record"`
		Field  string          `json:"field"`
		Value  json.RawMessage `json:"value"`
	}
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	for _, fixture := range fixtures {
		t.Run(fixture.Name, func(t *testing.T) {
			original, err := os.ReadFile("testdata/spell_records.json")
			if err != nil {
				t.Fatal(err)
			}
			var root map[string]any
			if err := json.Unmarshal(original, &root); err != nil {
				t.Fatal(err)
			}
			var value any
			if err := json.Unmarshal(fixture.Value, &value); err != nil {
				t.Fatal(err)
			}
			target := root
			if fixture.Record >= 0 {
				target = root["records"].([]any)[fixture.Record].(map[string]any)
			}
			fields := strings.Split(fixture.Field, ".")
			for _, field := range fields[:len(fields)-1] {
				target = target[field].(map[string]any)
			}
			target[fields[len(fields)-1]] = value
			changed, err := json.Marshal(root)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeSpellRecords(bytes.NewReader(changed)); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
	for _, bad := range []string{"{", "null", "{}", "[]"} {
		if _, err := DecodeSpellRecords(strings.NewReader(bad)); err == nil {
			t.Fatalf("accepted %q", bad)
		}
	}
	original, err := os.ReadFile("testdata/spell_records.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeSpellRecords(bytes.NewReader(append(original, []byte(" {}")...))); err == nil {
		t.Fatal("trailing JSON accepted")
	}
	r := loadSpellRecords(t)
	r.Records = append(r.Records, r.Records[0])
	if r.Validate() == nil {
		t.Fatal("duplicate record accepted")
	}
}
