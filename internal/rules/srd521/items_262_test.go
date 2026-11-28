package srd521

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func loadItems262(t *testing.T) ItemSourceBatch[ItemDetails262] {
	t.Helper()
	b, err := LoadItems262()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities, variants and listed prices from pinned PDF p. 100 for issue #262.
func TestItems262SourceIdentities(t *testing.T) {
	expected := []itemFamilyWant{
		{"Spikes, Iron", "spikes-iron", 100, []itemVariantWant{{"standard", "Spikes, Iron", 1, "gp"}}},
		{"Spyglass", "spyglass", 100, []itemVariantWant{{"standard", "Spyglass", 1000, "gp"}}},
		{"String", "string", 100, []itemVariantWant{{"standard", "String", 1, "sp"}}},
		{"Tent", "tent", 100, []itemVariantWant{{"standard", "Tent", 2, "gp"}}},
		{"Tinderbox", "tinderbox", 100, []itemVariantWant{{"standard", "Tinderbox", 5, "sp"}}},
		{"Torch", "torch", 100, []itemVariantWant{{"standard", "Torch", 1, "cp"}}},
		{"Vial", "vial", 100, []itemVariantWant{{"standard", "Vial", 1, "gp"}}},
		{"Waterskin", "waterskin", 100, []itemVariantWant{{"standard", "Waterskin", 2, "sp"}}},
	}
	b := loadItems262(t)
	if len(b.Families) != 8 || b.Batch != "data-adventuring-gear-07" || b.CatalogueIssue != 262 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the new gear families of issue #262")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "p. 100") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	checkMundaneIdentities(t, b.Families, "adventuring-gear", "adventuring-gear", expected)
	// Join every family to the decision inventory the issue cites, in its order.
	// Spell Scroll is first in the inventory for this issue. Batch #242 owns it,
	// so it joins there and the batch holds the other eight.
	owner := requireSpellScrollOwner(t)
	joinItemInventory(t, append([]ItemSourceIdentity{owner}, b.SourceIdentities()...), b.Batch, 262, "adventuring-gear", []int{99, 100, 100, 100, 100, 100, 100, 100, 100})
	b.Families[0].Variants[0].Details.Gear.SourcePages[0] = 1
	if fresh := loadItems262(t); fresh.Families[0].Variants[0].Details.Gear.SourcePages[0] != 100 {
		t.Fatal("load must return independent data")
	}
}

// Literal values, use modes and mechanic keys for all 8 records,
// transcribed from pinned PDF p. 100.
func TestItems262SourceValues(t *testing.T) {
	expected := []struct{ name, details, use, keys string }{
		{"Spikes, Iron/Spikes, Iron", `{"gear":{"actions":["utilize"],"checks":[],"damage":[],"measures":[{"name":"bundle-amount","numerator":10,"denominator":1,"unit":"spikes"}],"traits":["hammer-with-blunt-object-into-wood-earth-or-similar-material","jam-door-shut-or-anchor-rope-or-chain"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":5,"pounds_denominator":1,"condition":"none"},"source_pages":[100,100]}}`,
			"utilize", "item.gear action.utilize measure.bundle-amount-10-spikes trait.hammer-with-blunt-object-into-wood-earth-or-similar-material trait.jam-door-shut-or-anchor-rope-or-chain"},
		{"Spyglass/Spyglass", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"magnification","numerator":2,"denominator":1,"unit":"times-size"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[100,100]}}`,
			"none", "item.gear measure.magnification-2-times-size"},
		{"String/String", `{"gear":{"actions":["utilize"],"checks":[],"damage":[],"measures":[{"name":"length","numerator":10,"denominator":1,"unit":"feet"}],"traits":["can-be-knotted"],"contents":[],"spellcasting_focus":[],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[100,100]}}`,
			"utilize", "item.gear action.utilize measure.length-10-feet trait.can-be-knotted"},
		{"Tent/Tent", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"occupants","numerator":2,"denominator":1,"unit":"creatures"}],"traits":["occupants-are-small-or-medium"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":20,"pounds_denominator":1,"condition":"none"},"source_pages":[100,100]}}`,
			"none", "item.gear measure.occupants-2-creatures trait.occupants-are-small-or-medium"},
		{"Tinderbox/Tinderbox", `{"gear":{"actions":["bonus-action"],"checks":[],"damage":[],"measures":[{"name":"other-fire-lighting-time","numerator":1,"denominator":1,"unit":"minute"}],"traits":["small-container-of-flint-fire-steel-and-tinder","lights-candle-lamp-lantern-torch-or-exposed-fuel"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[100,100]}}`,
			"bonus-action", "item.gear action.bonus-action measure.other-fire-lighting-time-1-minute trait.small-container-of-flint-fire-steel-and-tinder trait.lights-candle-lamp-lantern-torch-or-exposed-fuel"},
		{"Torch/Torch", `{"gear":{"actions":["attack-action"],"checks":[],"damage":[{"purpose":"hit","dice":0,"sides":0,"flat":1,"type":"fire"}],"measures":[{"name":"burn-duration","numerator":1,"denominator":1,"unit":"hour"},{"name":"bright-light-radius","numerator":20,"denominator":1,"unit":"feet"},{"name":"additional-dim-light","numerator":20,"denominator":1,"unit":"feet"}],"traits":["attack-with-torch-as-simple-melee-weapon"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[100,100]}}`,
			"attack-action", "item.gear action.attack-action damage.hit-1-fire measure.burn-duration-1-hour measure.bright-light-radius-20-feet measure.additional-dim-light-20-feet trait.attack-with-torch-as-simple-melee-weapon"},
		{"Vial/Vial", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"capacity","numerator":4,"denominator":1,"unit":"ounces"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[100,100]}}`,
			"none", "item.gear measure.capacity-4-ounces"},
		{"Waterskin/Waterskin", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"capacity","numerator":4,"denominator":1,"unit":"pints"}],"traits":["insufficient-water-risks-dehydration"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":5,"pounds_denominator":1,"condition":"full"},"source_pages":[100,100]}}`,
			"none", "item.gear measure.capacity-4-pints trait.insufficient-water-risks-dehydration"},
	}
	index := 0
	for _, f := range loadItems262(t).Families {
		for _, v := range f.Variants {
			if index >= len(expected) {
				t.Fatal("batch holds more records than the fixture")
			}
			want := expected[index]
			index++
			t.Run(want.name, func(t *testing.T) {
				var details ItemDetails262
				d := json.NewDecoder(strings.NewReader(want.details))
				d.DisallowUnknownFields()
				if err := d.Decode(&details); err != nil {
					t.Fatal(err)
				}
				if f.Name+"/"+v.Name != want.name || !reflect.DeepEqual(v.Details, &details) {
					t.Fatalf("source value mismatch for %s/%s: %+v", f.Name, v.Name, v.Details.Gear)
				}
				if v.Equip.Mode != "none" || v.Use.Mode != want.use || len(v.Equip.Restrictions) != 0 || len(v.Use.Restrictions) != 0 || strings.Join(v.RequiredMechanicKeys, " ") != want.keys {
					t.Fatalf("source requirements mismatch: %s %v", v.Use.Mode, v.RequiredMechanicKeys)
				}
			})
		}
	}
	if index != len(expected) {
		t.Fatalf("fixture lists %d records, batch holds %d", len(expected), index)
	}
}

func TestItems262CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems262(t)
	composed, err := ComposeItemSourceIdentities(append(itemIdentitiesThrough(t, 261), batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, f := range composed {
		total += len(f.VariantIDs)
	}
	if len(composed) != 413 || total != 513 || composed[0].Name != "Potion of Healing" || composed[412].Name != "Waterskin" {
		t.Fatalf("expected Potion of Healing, batches #225 to #261 and #262: %d families, %d variants", len(composed), total)
	}
	checkItemBatchComposition(t, batch.SourceIdentities(), "adventuring-gear", "items_262")
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems262RejectMissingFields(t *testing.T) {
	rejectMissingItemMembers(t, items262JSON, 455, func(data []byte) error {
		_, err := DecodeItems262(data)
		return err
	})
}

func TestItems262RejectInvalidFixtures(t *testing.T) {
	rejectInvalidGearFixtures(t, LoadItems262, DecodeItems262, itemGear262)
}

// The mechanic keys restate every payload value except weight and the last
// source page, so changing any other value of any record fails.
func TestItems262RejectChangedValues(t *testing.T) {
	rejectChangedGearValues(t, items262JSON, func(data []byte) error {
		_, err := DecodeItems262(data)
		return err
	})
}

func TestItems262RejectUnsupportedJSON(t *testing.T) {
	decode := func(data []byte) error {
		_, err := DecodeItems262(data)
		return err
	}
	rejectUnsupportedItemJSON(t, items262JSON, `"id": "srd-5.2.1:content/items/spikes-iron/variant/standard"`, decode, []itemJSONChange{
		{"unknown payload", `"gear": {`, `"gear_plus": {`},
		{"unknown detail", `"spellcasting_focus": `, `"proficiency": `},
		{"number for boolean", `"listed": true`, `"listed": 1`},
		{"string for number", `"pounds_numerator": `, `"pounds_numerator": "1", "x": `},
		{"null check list", `"checks": [`, `"checks": null, "x": [`},
	})
	// Neither the Healing contract decoder nor the previous batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items262JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems261(items262JSON); err == nil {
		t.Fatal("#261 decoder accepted the #262 batch")
	}
	if _, err := DecodeItems262(items261JSON); err == nil {
		t.Fatal("#262 decoder accepted the #261 batch")
	}
}
