package srd521

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
)

func magicFacts(t *testing.T, change func(map[string]any)) []byte {
	t.Helper()
	m := map[string]any{"profile": DetectMagicFactsProfile, "bearer": "object", "actor_ids": []string{"actor"}, "location_id": "loc", "hidden": false, "distance_feet": 20, "barrier": "none", "visible": true, "origin": "spell", "school": "abjuration"}
	if change != nil {
		change(m)
	}
	raw, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func magicSubject(t *testing.T, id string, change func(map[string]any)) DetectMagicSubject {
	t.Helper()
	facts, ok := ParseDetectMagicFacts(magicFacts(t, change))
	if !ok {
		t.Fatalf("%s: facts rejected", id)
	}
	return DetectMagicSubject{ID: id, Label: id, Facts: facts}
}

// The resolver's constants must agree with the catalogued complete source
// record, which itself stays unavailable as a general spell.
func TestDetectMagicAgreesWithSourceRecord(t *testing.T) {
	var record *SpellRecord
	for _, s := range loadSpells202(t).Records {
		if s.ID == prefix+"spells/detect-magic" {
			record = &s
		}
	}
	if record == nil || record.MechanicsAvailable() || record.Scope != "complete-source-record" {
		t.Fatal("Detect Magic source record missing or promoted")
	}
	d := record.SourceDetails.DetectMagic
	if record.SourcePage != 123 || record.Level != 1 || record.School != "divination" || record.CastingTime != "action" || !*record.Components.Verbal || !*record.Components.Somatic || record.Components.Material != nil || !*record.Ritual || !*record.Concentration || record.DurationSeconds != detectMagicSeconds || d == nil || d.SenseRadiusFeet != DetectMagicRadiusFeet || d.FollowUpAction != "magic-action" || !slices.Contains(record.Lists, "cleric") {
		t.Fatalf("resolver disagrees with source record: %+v", record)
	}
	got := []string{}
	for _, b := range d.Blockers {
		got = append(got, b.Material+"-"+b.Thickness)
	}
	if !slices.Equal(got, detectMagicBarriers) {
		t.Fatalf("barriers %v, want %v", got, detectMagicBarriers)
	}
}

func TestDetectMagicProjectsOnlyTheFreeUse(t *testing.T) {
	_, sheet := casterBuild(t)
	var spell *Spell
	for i, s := range sheet.Spellcasting.Spells {
		if s.ID == DetectMagicID {
			spell = &sheet.Spellcasting.Spells[i]
		}
		if s.ID == "bless" && s.Playable {
			t.Fatal("unrelated recorded spell became playable")
		}
	}
	if spell == nil || !spell.Playable || len(spell.SlotLevels) != 0 || !spell.ComponentsAvailable {
		t.Fatalf("Detect Magic free use not projected: %+v", spell)
	}
	r := sheet.Resources[len(sheet.Resources)-1]
	if r.ID != "magic-initiate" || r.Current != 1 || r.Maximum != 1 || strings.Contains(r.Label, "unavailable") {
		t.Fatalf("feat resource=%+v", r)
	}
}

func TestDetectMagicFactsFailClosed(t *testing.T) {
	for name, change := range map[string]func(map[string]any){
		"missing bearer":  func(m map[string]any) { delete(m, "bearer") },
		"null bearer":     func(m map[string]any) { m["bearer"] = nil },
		"empty bearer":    func(m map[string]any) { m["bearer"] = "" },
		"unknown bearer":  func(m map[string]any) { m["bearer"] = "construct" },
		"hidden lore":     func(m map[string]any) { m["hidden"] = true },
		"missing field":   func(m map[string]any) { delete(m, "visible") },
		"null field":      func(m map[string]any) { m["barrier"] = nil },
		"extra field":     func(m map[string]any) { m["aura"] = "bright" },
		"unknown barrier": func(m map[string]any) { m["barrier"] = "wood-3-inches" },
		"unknown school":  func(m map[string]any) { m["school"] = "chronomancy" },
		"school without spell": func(m map[string]any) {
			m["origin"] = "other"
		},
		"spell without school": func(m map[string]any) { m["school"] = "" },
		"unknown origin":       func(m map[string]any) { m["origin"] = "rumour" },
		"wrong profile":        func(m map[string]any) { m["profile"] = "detect-magic-v2" },
		"no actors":            func(m map[string]any) { m["actor_ids"] = []string{} },
		"negative distance":    func(m map[string]any) { m["distance_feet"] = -5 },
	} {
		if _, ok := ParseDetectMagicFacts(magicFacts(t, change)); ok {
			t.Errorf("%s accepted", name)
		}
	}
	if _, ok := ParseDetectMagicFacts(magicFacts(t, func(m map[string]any) { m["origin"], m["school"] = "other", "" })); !ok {
		t.Fatal("non-spell magic rejected")
	}
	if _, ok := ParseDetectMagicFacts(magicFacts(t, func(m map[string]any) { m["bearer"] = "creature" })); !ok {
		t.Fatal("creature bearer rejected")
	}
}

func TestDetectMagicCastSpendsFeatUseAndReplacesConcentration(t *testing.T) {
	b, sheet := casterBuild(t)
	if _, err := Lifecycle(b, sheet, SpellCast, "shield-of-faith", "faith", 100, &spellDice{}); err != nil {
		t.Fatal(err)
	}
	slots := b.Casting.Slots
	subjects := []DetectMagicSubject{
		magicSubject(t, "far", func(m map[string]any) { m["distance_feet"] = 35 }),
		magicSubject(t, "leaded", func(m map[string]any) { m["barrier"] = "lead-thin-sheet" }),
	}
	r, err := CastDetectMagic(b, sheet, subjects, "cast-1", 160)
	if err != nil {
		t.Fatal(err)
	}
	if r.Presence || r.FeatUsesBefore != 1 || r.FeatUsesAfter != 0 || b.Casting.Choices.FeatUses != 0 || r.SlotsAfter != slots || b.Casting.Slots != slots || r.Before == nil || r.Before.ID != "faith" || r.After.ID != "cast-1" || r.After.Spell != DetectMagicID || r.After.ExpiresAt != 760 || !strings.Contains(r.Reason, "does not establish") {
		t.Fatalf("cast receipt=%+v", r)
	}
	if _, err := CastDetectMagic(b, sheet, nil, "cast-2", 170); err == nil || !strings.Contains(err.Error(), "spent") {
		t.Fatalf("spent free use cast again: %v", err)
	}
	// Reload keeps the effect, without Shield of Faith's AC term.
	data, err := json.Marshal(map[string]any{"srd521": b})
	if err != nil {
		t.Fatal(err)
	}
	reloaded, err := Project(data, 1, 2, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Concentration == nil || reloaded.Concentration.ID != "cast-1" || reloaded.Defenses["armor-class"].Value != 11 || reloaded.Resources[len(reloaded.Resources)-1].Current != 0 {
		t.Fatalf("reload=%+v", reloaded)
	}
	b.Casting.Concentration.RecipientID = "someone"
	data, _ = json.Marshal(map[string]any{"srd521": b})
	if _, err := Project(data, 1, 2, 10, 0); err == nil {
		t.Fatal("Detect Magic recipient accepted")
	}
}

func TestDetectMagicCastRequiresRecordedChoiceAndComponents(t *testing.T) {
	b, sheet := casterBuild(t)
	b.Casting.Choices = nil
	if _, err := CastDetectMagic(b, sheet, nil, "x", 0); err == nil {
		t.Fatal("missing feat choice cast")
	}
	b, sheet = casterBuild(t)
	b.Casting.Verbal = false
	if _, err := CastDetectMagic(b, sheet, nil, "x", 0); err == nil || b.Casting.Choices.FeatUses != 1 {
		t.Fatal("cast without components or spent on failure")
	}
	b, sheet = casterBuild(t)
	sheet.HP.Value = 0
	if _, err := CastDetectMagic(b, sheet, nil, "x", 0); err == nil {
		t.Fatal("unconscious cast")
	}
}

func TestDetectMagicStudyDisclosesOnlyVisibleAdjudicatedAuras(t *testing.T) {
	b, sheet := casterBuild(t)
	subjects := []DetectMagicSubject{
		magicSubject(t, "ward", nil),
		magicSubject(t, "relic", func(m map[string]any) { m["bearer"], m["origin"], m["school"] = "creature", "other", "" }),
		magicSubject(t, "unseen", func(m map[string]any) { m["visible"] = false; m["school"] = "illusion" }),
		magicSubject(t, "walled", func(m map[string]any) { m["barrier"] = "stone-dirt-or-wood-1-foot"; m["school"] = "necromancy" }),
	}
	if _, err := StudyDetectMagic(b, sheet, subjects, 100); err == nil {
		t.Fatal("study without concentration")
	}
	cast, err := CastDetectMagic(b, sheet, subjects, "cast", 100)
	if err != nil || !cast.Presence {
		t.Fatalf("cast=%+v err=%v", cast, err)
	}
	if _, err := StudyDetectMagic(b, sheet, subjects[3:], 110); err == nil {
		t.Fatal("study with nothing sensed")
	}
	r, err := StudyDetectMagic(b, sheet, subjects, 110)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Auras) != 2 || r.Auras[0].SubjectID != "ward" || r.Auras[0].School != "abjuration" || r.Auras[1].SubjectID != "relic" || r.Auras[1].School != "" || r.FeatUsesAfter != 0 || r.SlotsAfter != r.SlotsBefore {
		t.Fatalf("study=%+v", r)
	}
	for _, hidden := range []string{"unseen", "walled", "illusion", "necromancy"} {
		if strings.Contains(r.Reason, hidden) {
			t.Fatalf("study disclosed %s: %s", hidden, r.Reason)
		}
	}
	if _, err := StudyDetectMagic(b, sheet, subjects, 700); err == nil {
		t.Fatal("study after expiry")
	}
}

func TestDetectMagicLifecycleEndsAndRecovers(t *testing.T) {
	b, sheet := casterBuild(t)
	if _, err := CastDetectMagic(b, sheet, nil, "cast", 0); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		r, err := Lifecycle(b, sheet, WaitMinute, "", "", i*60, nil)
		if err != nil {
			t.Fatal(err)
		}
		if i == 9 && (b.Casting.Concentration != nil || !strings.Contains(r.Reason, "Detect Magic expired")) {
			t.Fatalf("expiry=%+v", r)
		}
	}
	if _, err := CastDetectMagic(b, sheet, nil, "again", 600); err == nil {
		t.Fatal("cast without Long Rest")
	}
	if _, err := Lifecycle(b, sheet, LongRest, "", "", 600, nil); err != nil || b.Casting.Choices.FeatUses != 1 {
		t.Fatalf("Long Rest did not recover free use: %v", err)
	}
	if _, err := CastDetectMagic(b, sheet, nil, "again", 30000); err != nil {
		t.Fatal(err)
	}
	r, err := ConcentrationDamage(b, sheet, 4, 1, 30010, &spellDice{value: 0})
	if err != nil || r == nil || r.Check == nil || b.Casting.Concentration != nil {
		t.Fatalf("failed save kept concentration: %+v %v", r, err)
	}
}
