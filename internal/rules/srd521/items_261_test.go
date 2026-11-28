package srd521

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func loadItems261(t *testing.T) ItemSourceBatch[ItemDetails261] {
	t.Helper()
	b, err := LoadItems261()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities, variants and listed prices from pinned PDF p. 99 for issue #261.
func TestItems261SourceIdentities(t *testing.T) {
	expected := []itemFamilyWant{
		{"Pouch", "pouch", 99, []itemVariantWant{{"standard", "Pouch", 5, "sp"}}},
		{"Priest’s Pack", "priests-pack", 99, []itemVariantWant{{"standard", "Priest’s Pack", 33, "gp"}}},
		{"Quiver", "quiver", 99, []itemVariantWant{{"standard", "Quiver", 1, "gp"}}},
		{"Ram, Portable", "ram-portable", 99, []itemVariantWant{{"standard", "Ram, Portable", 4, "gp"}}},
		{"Rations", "rations", 99, []itemVariantWant{{"standard", "Rations", 5, "sp"}}},
		{"Robe", "robe", 99, []itemVariantWant{{"standard", "Robe", 1, "gp"}}},
		{"Rope", "rope", 99, []itemVariantWant{{"standard", "Rope", 1, "gp"}}},
		{"Sack", "sack", 99, []itemVariantWant{{"standard", "Sack", 1, "cp"}}},
		{"Scholar’s Pack", "scholars-pack", 99, []itemVariantWant{{"standard", "Scholar’s Pack", 40, "gp"}}},
		{"Shovel", "shovel", 99, []itemVariantWant{{"standard", "Shovel", 2, "gp"}}},
		{"Signal Whistle", "signal-whistle", 99, []itemVariantWant{{"standard", "Signal Whistle", 5, "cp"}}},
	}
	b := loadItems261(t)
	if len(b.Families) != 11 || b.Batch != "data-adventuring-gear-06" || b.CatalogueIssue != 261 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the new gear families of issue #261")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "p. 99") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	checkMundaneIdentities(t, b.Families, "adventuring-gear", "adventuring-gear", expected)
	// Join every family to the decision inventory the issue cites, in its order.
	// Potion of Healing is first in the inventory for this issue. The Healing
	// contract owns it, so it joins there and the batch holds the other eleven.
	owner := requirePotionOfHealingOwner(t)
	joinItemInventory(t, append([]ItemSourceIdentity{owner}, b.SourceIdentities()...), b.Batch, 261, "adventuring-gear", []int{99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99, 99})
	b.Families[0].Variants[0].Details.Gear.SourcePages[0] = 1
	if fresh := loadItems261(t); fresh.Families[0].Variants[0].Details.Gear.SourcePages[0] != 99 {
		t.Fatal("load must return independent data")
	}
}

// Literal values, use modes and mechanic keys for all 11 records,
// transcribed from pinned PDF p. 99.
func TestItems261SourceValues(t *testing.T) {
	expected := []struct{ name, details, use, keys string }{
		{"Pouch/Pouch", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"capacity-weight","numerator":6,"denominator":1,"unit":"pounds"},{"name":"capacity-volume","numerator":1,"denominator":5,"unit":"cubic-foot"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"none", "item.gear measure.capacity-weight-6-pounds measure.capacity-volume-1-over-5-cubic-foot"},
		{"Priest’s Pack/Priest’s Pack", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":[],"contents":[{"item":"backpack","quantity":1,"unit":"item"},{"item":"blanket","quantity":1,"unit":"item"},{"item":"holy-water","quantity":1,"unit":"item"},{"item":"lamp","quantity":1,"unit":"item"},{"item":"rations","quantity":7,"unit":"day"},{"item":"robe","quantity":1,"unit":"item"},{"item":"tinderbox","quantity":1,"unit":"item"}],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":29,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"none", "item.gear content.backpack-1-item content.blanket-1-item content.holy-water-1-item content.lamp-1-item content.rations-7-day content.robe-1-item content.tinderbox-1-item"},
		{"Quiver/Quiver", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"capacity","numerator":20,"denominator":1,"unit":"arrows"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"none", "item.gear measure.capacity-20-arrows"},
		{"Ram, Portable/Ram, Portable", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"strength-check-bonus","numerator":4,"denominator":1,"unit":"bonus"},{"name":"helpers","numerator":1,"denominator":1,"unit":"character"}],"traits":["used-to-break-down-doors","helper-grants-advantage-on-check"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":35,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"none", "item.gear measure.strength-check-bonus-4-bonus measure.helpers-1-character trait.used-to-break-down-doors trait.helper-grants-advantage-on-check"},
		{"Rations/Rations", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["travel-ready-food-including-jerky-dried-fruit-hardtack-and-nuts","not-eating-risks-malnutrition"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"none", "item.gear trait.travel-ready-food-including-jerky-dried-fruit-hardtack-and-nuts trait.not-eating-risks-malnutrition"},
		{"Robe/Robe", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["vocational-or-ceremonial-significance","required-colors-or-symbols-for-admission-to-some-events-and-locations"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":4,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"none", "item.gear trait.vocational-or-ceremonial-significance trait.required-colors-or-symbols-for-admission-to-some-events-and-locations"},
		{"Rope/Rope", `{"gear":{"actions":["utilize","action"],"checks":[{"purpose":"tie-knot","kind":"ability-check","ability":"dexterity","skill":"sleight-of-hand","dc":10,"dc_rule":"fixed"},{"purpose":"burst","kind":"ability-check","ability":"strength","skill":"athletics","dc":20,"dc_rule":"fixed"},{"purpose":"escape","kind":"ability-check","ability":"dexterity","skill":"acrobatics","dc":15,"dc_rule":"fixed"}],"damage":[],"measures":[],"traits":["bind-only-unwilling-grappled-incapacitated-or-restrained-creature","bound-legs-restrain-until-escape"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":5,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"multiple", "item.gear action.utilize action.action ability-check.tie-knot-dexterity-sleight-of-hand-dc-10 ability-check.burst-strength-athletics-dc-20 ability-check.escape-dexterity-acrobatics-dc-15 trait.bind-only-unwilling-grappled-incapacitated-or-restrained-creature trait.bound-legs-restrain-until-escape"},
		{"Sack/Sack", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"capacity-weight","numerator":30,"denominator":1,"unit":"pounds"},{"name":"capacity-volume","numerator":1,"denominator":1,"unit":"cubic-foot"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":2,"condition":"none"},"source_pages":[99,99]}}`,
			"none", "item.gear measure.capacity-weight-30-pounds measure.capacity-volume-1-cubic-foot"},
		{"Scholar’s Pack/Scholar’s Pack", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":[],"contents":[{"item":"backpack","quantity":1,"unit":"item"},{"item":"book","quantity":1,"unit":"item"},{"item":"ink","quantity":1,"unit":"item"},{"item":"ink-pen","quantity":1,"unit":"item"},{"item":"lamp","quantity":1,"unit":"item"},{"item":"oil","quantity":10,"unit":"flask"},{"item":"parchment","quantity":10,"unit":"sheet"},{"item":"tinderbox","quantity":1,"unit":"item"}],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":22,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"none", "item.gear content.backpack-1-item content.book-1-item content.ink-1-item content.ink-pen-1-item content.lamp-1-item content.oil-10-flask content.parchment-10-sheet content.tinderbox-1-item"},
		{"Shovel/Shovel", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"digging-time","numerator":1,"denominator":1,"unit":"hour"},{"name":"hole-side","numerator":5,"denominator":1,"unit":"feet"}],"traits":["digs-hole-in-soil-or-similar-material"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":5,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"none", "item.gear measure.digging-time-1-hour measure.hole-side-5-feet trait.digs-hole-in-soil-or-similar-material"},
		{"Signal Whistle/Signal Whistle", `{"gear":{"actions":["utilize"],"checks":[],"damage":[],"measures":[{"name":"audible-distance","numerator":600,"denominator":1,"unit":"feet"}],"traits":["blown-to-produce-sound"],"contents":[],"spellcasting_focus":[],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"utilize", "item.gear action.utilize measure.audible-distance-600-feet trait.blown-to-produce-sound"},
	}
	index := 0
	for _, f := range loadItems261(t).Families {
		for _, v := range f.Variants {
			if index >= len(expected) {
				t.Fatal("batch holds more records than the fixture")
			}
			want := expected[index]
			index++
			t.Run(want.name, func(t *testing.T) {
				var details ItemDetails261
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

func TestItems261CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems261(t)
	composed, err := ComposeItemSourceIdentities(append(itemIdentitiesThrough(t, 260), batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, f := range composed {
		total += len(f.VariantIDs)
	}
	if len(composed) != 405 || total != 505 || composed[0].Name != "Potion of Healing" || composed[404].Name != "Signal Whistle" {
		t.Fatalf("expected Potion of Healing, batches #225 to #260 and #261: %d families, %d variants", len(composed), total)
	}
	checkItemBatchComposition(t, batch.SourceIdentities(), "adventuring-gear", "items_261")
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems261RejectMissingFields(t *testing.T) {
	rejectMissingItemMembers(t, items261JSON, 660, func(data []byte) error {
		_, err := DecodeItems261(data)
		return err
	})
}

func TestItems261RejectInvalidFixtures(t *testing.T) {
	rejectInvalidGearFixtures(t, LoadItems261, DecodeItems261, itemGear261)
}

// The mechanic keys restate every payload value except weight and the last
// source page, so changing any other value of any record fails.
func TestItems261RejectChangedValues(t *testing.T) {
	rejectChangedGearValues(t, items261JSON, func(data []byte) error {
		_, err := DecodeItems261(data)
		return err
	})
}

func TestItems261RejectUnsupportedJSON(t *testing.T) {
	decode := func(data []byte) error {
		_, err := DecodeItems261(data)
		return err
	}
	rejectUnsupportedItemJSON(t, items261JSON, `"id": "srd-5.2.1:content/items/pouch/variant/standard"`, decode, []itemJSONChange{
		{"unknown payload", `"gear": {`, `"gear_plus": {`},
		{"unknown detail", `"spellcasting_focus": `, `"proficiency": `},
		{"number for boolean", `"listed": true`, `"listed": 1`},
		{"string for number", `"pounds_numerator": `, `"pounds_numerator": "1", "x": `},
		{"null check list", `"checks": [`, `"checks": null, "x": [`},
	})
	// Neither the Healing contract decoder nor the previous batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items261JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems260(items261JSON); err == nil {
		t.Fatal("#260 decoder accepted the #261 batch")
	}
	if _, err := DecodeItems261(items260JSON); err == nil {
		t.Fatal("#261 decoder accepted the #260 batch")
	}
}

// The Cleric creation fixture names a Robe, and the character source
// catalogues refer to these keys.
func TestItems261JoinRuntimeGear(t *testing.T) {
	families := loadItems261(t).Families
	joinCatalogueEquipment(t, families, []itemGearKeyWant{
		{"pouch", "pouch", "standard"},
		{"priests-pack", "priests-pack", "standard"},
		{"quiver", "quiver", "standard"},
		{"robe", "robe", "standard"},
		{"scholars-pack", "scholars-pack", "standard"},
	})
	build, err := SupportedBuildForConcept("human cleric", "acolyte")
	if err != nil {
		t.Fatal(err)
	}
	requireCreationNames(t, families, build, map[string]string{"Robe": "Robe"})
}
