package srd521

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func loadLightFixture(t *testing.T) SpellRecord {
	t.Helper()
	data, err := os.ReadFile("testdata/spell_light_fixture.json")
	if err != nil {
		t.Fatal(err)
	}
	r, err := DecodeSpellRecords(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Records) != 1 || r.Source.Attribution != Attribution {
		t.Fatal("one attributed Light fixture required")
	}
	return r.Records[0]
}

// The executable duration and radii must equal the pinned PDF p.144 record.
func TestLightConstantsMatchSource(t *testing.T) {
	var source *SpellRecord
	for _, r := range loadSpells211(t).Records {
		if r.ID == prefix+"spells/light" {
			source = &r
		}
	}
	if source == nil || source.SourceDetails == nil || source.SourceDetails.Light == nil {
		t.Fatal("Light source record missing")
	}
	d := source.SourceDetails.Light
	if source.DurationSeconds != LightDurationSeconds || d.BrightLightRadiusFeet != LightBrightFeet || d.DimLightAdditionalFeet != LightDimFeet || *source.Concentration || source.Level != 0 || source.CastingTime != "action" {
		t.Fatal("Light duration, radii or timing differ from source")
	}
	if d.ObjectRequirement != "not-worn-or-carried-by-someone-else" || d.Recast != "spell-ends-if-cast-again" || d.Color != "colored-as-caster-likes" {
		t.Fatal("Light object, recast or color rule differs from source")
	}
	if source.MechanicsAvailable() {
		t.Fatal("complete source record became available")
	}
	fixture := loadLightFixture(t)
	if fixture.DurationSeconds != source.DurationSeconds || fixture.Components.Material.Description != source.Components.Material.Description || fixture.Range != source.Range || fixture.SourcePage != source.SourcePage {
		t.Fatal("bounded fixture disagrees with source")
	}
}

func TestLightFixtureFailsClosed(t *testing.T) {
	fixture := loadLightFixture(t)
	if !fixture.MechanicsAvailable() {
		t.Fatal("bounded Light fixture unavailable")
	}
	mutations := map[string]func(*SpellRecord){
		"unknown key": func(s *SpellRecord) {
			s.RequiredMechanics = append(append([]string{}, s.RequiredMechanics...), "light.arbitrary-illumination")
		},
		"no keys":        func(s *SpellRecord) { s.RequiredMechanics = nil },
		"hostile holder": func(s *SpellRecord) { s.Target.Eligibility = "creature-held-object" },
		"source target":  func(s *SpellRecord) { s.Target.Eligibility = "unattended-large-or-smaller-object-touch" },
		"duration":       func(s *SpellRecord) { s.DurationSeconds = 7200 },
		"concentration":  func(s *SpellRecord) { v := true; s.Concentration = &v },
		"mechanic": func(s *SpellRecord) {
			s.Effects = []SpellEffect{s.Effects[0]}
			s.Effects[0].Mechanic = "light.bright-40"
		},
		"source details": func(s *SpellRecord) { s.SourceDetails = &SpellSourceDetails{Light: &SpellLight{}} },
		"slot":           func(s *SpellRecord) { s.SlotLevels = []int{1} },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			data, _ := json.Marshal(fixture)
			var changed SpellRecord
			if err := json.Unmarshal(data, &changed); err != nil {
				t.Fatal(err)
			}
			mutate(&changed)
			if changed.MechanicsAvailable() {
				t.Fatal("changed fixture available")
			}
		})
	}
}

func TestCastLightRecastDismissAndExpiry(t *testing.T) {
	b, sheet := casterBuild(t)
	var light *Spell
	for i := range sheet.Spellcasting.Spells {
		if sheet.Spellcasting.Spells[i].ID == LightSpell {
			light = &sheet.Spellcasting.Spells[i]
		}
	}
	if light == nil || !light.ComponentsAvailable || len(light.SlotLevels) != 0 || light.SourceID != prefix+"spells/light" {
		t.Fatalf("Light not projected as a castable cantrip: %+v", light)
	}
	if _, err := Lifecycle(b, sheet, SpellCast, "shield-of-faith", "faith", 100, &spellDice{}); err != nil {
		t.Fatal(err)
	}
	faith := *b.Casting.Concentration
	slots := b.Casting.Slots
	for _, color := range []string{"", "ultraviolet", "White"} {
		if _, err := CastLight(b, sheet, "item-a", "Torch", color, "cast-1", 100); err == nil {
			t.Fatalf("color %q accepted", color)
		}
	}
	first, err := CastLight(b, sheet, "item-a", "Torch", "blue", "cast-1", 100)
	if err != nil {
		t.Fatal(err)
	}
	want := LightEffect{ID: "cast-1", ItemID: "item-a", Color: "blue", StartedAt: 100, ExpiresAt: 3700, BrightFeet: 20, DimFeet: 20}
	if first.Before != nil || *first.After != want || b.Casting.Slots != slots || *b.Casting.Concentration != faith {
		t.Fatalf("first cast wrong or touched other effects: %+v", first)
	}
	second, err := CastLight(b, sheet, "item-b", "Mace", "red", "cast-2", 200)
	if err != nil {
		t.Fatal(err)
	}
	if second.Before == nil || *second.Before != want || second.After.ItemID != "item-b" || second.After.ExpiresAt != 3800 || *b.Casting.Concentration != faith {
		t.Fatal("recast did not replace exactly the previous Light")
	}
	// Save and reload keeps the exact absolute expiry.
	stats, err := json.Marshal(map[string]any{"srd521": b})
	if err != nil {
		t.Fatal(err)
	}
	reloaded, err := Project(stats, 1, sheet.HP.Value, sheet.MaxHP.Value, 0)
	if err != nil {
		t.Fatal(err)
	}
	var envelope struct {
		Build Build `json:"srd521"`
	}
	if err := json.Unmarshal(stats, &envelope); err != nil || reloaded == nil {
		t.Fatal(err)
	}
	saved := envelope.Build.Casting.Light
	if *saved != *second.After || !saved.Active(3799) || saved.Active(3800) || saved.Active(199) {
		t.Fatal("reload changed Light expiry")
	}
	if _, err := EndLight(&envelope.Build, "Mace", 3800); err == nil {
		t.Fatal("expired Light dismissed")
	}
	if _, err := CastLight(&envelope.Build, sheet, "item-c", "Torch", "white", "cast-3", 3800); err != nil {
		t.Fatal(err)
	}
	if envelope.Build.Casting.Light.ID != "cast-3" {
		t.Fatal("recast after expiry failed")
	}
	ended, err := EndLight(&envelope.Build, "Torch", 3801)
	if err != nil || ended.Before.ID != "cast-3" || envelope.Build.Casting.Light != nil || envelope.Build.Casting.Concentration.ID != faith.ID {
		t.Fatal("dismissal removed the wrong effect")
	}
	if _, err := EndLight(&envelope.Build, "Torch", 3801); err == nil {
		t.Fatal("second dismissal accepted")
	}
}

func TestLightUnavailableAndTamperedSaves(t *testing.T) {
	b, sheet := casterBuild(t)
	b.Casting.Focus = false
	if _, err := CastLight(b, sheet, "item", "Torch", "white", "id", 0); err == nil {
		t.Fatal("Light cast without material substitute")
	}
	b, sheet = casterBuild(t)
	down := *sheet
	down.HP.Value = 0
	if _, err := CastLight(b, &down, "item", "Torch", "white", "id", 0); err == nil {
		t.Fatal("unconscious caster cast Light")
	}
	for name, mutate := range map[string]func(*LightEffect){
		"bright":   func(l *LightEffect) { l.BrightFeet = 40 },
		"dim":      func(l *LightEffect) { l.DimFeet = 0 },
		"duration": func(l *LightEffect) { l.ExpiresAt += 60 },
		"color":    func(l *LightEffect) { l.Color = "black" },
		"item":     func(l *LightEffect) { l.ItemID = "" },
	} {
		t.Run(name, func(t *testing.T) {
			b, sheet := casterBuild(t)
			if _, err := CastLight(b, sheet, "item", "Torch", "white", "id", 10); err != nil {
				t.Fatal(err)
			}
			mutate(b.Casting.Light)
			stats, _ := json.Marshal(map[string]any{"srd521": b})
			if _, err := Project(stats, 1, sheet.HP.Value, sheet.MaxHP.Value, 0); err == nil {
				t.Fatal("tampered Light accepted")
			}
		})
	}
}
