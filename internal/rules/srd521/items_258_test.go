package srd521

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func loadItems258(t *testing.T) ItemSourceBatch[ItemDetails258] {
	t.Helper()
	b, err := LoadItems258()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities, variants and listed prices from pinned PDF p. 97 for issue #258.
func TestItems258SourceIdentities(t *testing.T) {
	expected := []itemFamilyWant{
		{"Clothes, Fine", "clothes-fine", 97, []itemVariantWant{{"standard", "Clothes, Fine", 15, "gp"}}},
		{"Clothes, Traveler’s", "clothes-travelers", 97, []itemVariantWant{{"standard", "Clothes, Traveler’s", 2, "gp"}}},
		{"Component Pouch", "component-pouch", 97, []itemVariantWant{{"standard", "Component Pouch", 25, "gp"}}},
		{"Costume", "costume", 97, []itemVariantWant{{"standard", "Costume", 5, "gp"}}},
		{"Crowbar", "crowbar", 97, []itemVariantWant{{"standard", "Crowbar", 2, "gp"}}},
		{"Diplomat’s Pack", "diplomats-pack", 97, []itemVariantWant{{"standard", "Diplomat’s Pack", 39, "gp"}}},
		{"Druidic Focus", "druidic-focus", 97, []itemVariantWant{{"sprig-of-mistletoe", "Sprig of mistletoe", 1, "gp"}, {"wooden-staff", "Wooden staff", 5, "gp"}, {"yew-wand", "Yew wand", 10, "gp"}}},
		{"Dungeoneer’s Pack", "dungeoneers-pack", 97, []itemVariantWant{{"standard", "Dungeoneer’s Pack", 12, "gp"}}},
		{"Entertainer’s Pack", "entertainers-pack", 97, []itemVariantWant{{"standard", "Entertainer’s Pack", 40, "gp"}}},
		{"Explorer’s Pack", "explorers-pack", 97, []itemVariantWant{{"standard", "Explorer’s Pack", 10, "gp"}}},
		{"Flask", "flask", 97, []itemVariantWant{{"standard", "Flask", 2, "cp"}}},
		{"Grappling Hook", "grappling-hook", 97, []itemVariantWant{{"standard", "Grappling Hook", 2, "gp"}}},
	}
	b := loadItems258(t)
	if len(b.Families) != 12 || b.Batch != "data-adventuring-gear-03" || b.CatalogueIssue != 258 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the new gear families of issue #258")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "p. 97") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	checkMundaneIdentities(t, b.Families, "adventuring-gear", "adventuring-gear", expected)
	// Join every family to the decision inventory the issue cites, in its order.
	joinItemInventory(t, b.SourceIdentities(), b.Batch, 258, "adventuring-gear", []int{97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97, 97})
	b.Families[0].Variants[0].Details.Gear.SourcePages[0] = 1
	if fresh := loadItems258(t); fresh.Families[0].Variants[0].Details.Gear.SourcePages[0] != 97 {
		t.Fatal("load must return independent data")
	}
}

// Literal values, use modes and mechanic keys for all 14 records,
// transcribed from pinned PDF p. 97.
func TestItems258SourceValues(t *testing.T) {
	expected := []struct{ name, details, use, keys string }{
		{"Clothes, Fine/Clothes, Fine", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["expensive-fabrics-with-expertly-crafted-details","required-for-admission-to-some-events-and-locations"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":6,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear trait.expensive-fabrics-with-expertly-crafted-details trait.required-for-admission-to-some-events-and-locations"},
		{"Clothes, Traveler’s/Clothes, Traveler’s", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["resilient-garments-for-travel-in-various-environments"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":4,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear trait.resilient-garments-for-travel-in-various-environments"},
		{"Component Pouch/Component Pouch", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["watertight","compartments-hold-all-free-material-components-of-your-spells"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear trait.watertight trait.compartments-hold-all-free-material-components-of-your-spells"},
		{"Costume/Costume", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["advantage-on-checks-to-impersonate-represented-person-or-type-while-wearing"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":4,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear trait.advantage-on-checks-to-impersonate-represented-person-or-type-while-wearing"},
		{"Crowbar/Crowbar", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["advantage-on-strength-checks-where-leverage-applies"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":5,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear trait.advantage-on-strength-checks-where-leverage-applies"},
		{"Diplomat’s Pack/Diplomat’s Pack", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":[],"contents":[{"item":"chest","quantity":1,"unit":"item"},{"item":"clothes-fine","quantity":1,"unit":"item"},{"item":"ink","quantity":1,"unit":"item"},{"item":"ink-pen","quantity":5,"unit":"item"},{"item":"lamp","quantity":1,"unit":"item"},{"item":"case-map-or-scroll","quantity":2,"unit":"item"},{"item":"oil","quantity":4,"unit":"flask"},{"item":"paper","quantity":5,"unit":"sheet"},{"item":"parchment","quantity":5,"unit":"sheet"},{"item":"perfume","quantity":1,"unit":"item"},{"item":"tinderbox","quantity":1,"unit":"item"}],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":39,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear content.chest-1-item content.clothes-fine-1-item content.ink-1-item content.ink-pen-5-item content.lamp-1-item content.case-map-or-scroll-2-item content.oil-4-flask content.paper-5-sheet content.parchment-5-sheet content.perfume-1-item content.tinderbox-1-item"},
		{"Druidic Focus/Sprig of mistletoe", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["carved-ribboned-or-painted-to-channel-primal-magic"],"contents":[],"spellcasting_focus":["druid","ranger"],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear spellcasting-focus.druid spellcasting-focus.ranger trait.carved-ribboned-or-painted-to-channel-primal-magic"},
		{"Druidic Focus/Wooden staff", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["carved-ribboned-or-painted-to-channel-primal-magic","also-quarterstaff"],"contents":[],"spellcasting_focus":["druid","ranger"],"weight":{"listed":true,"pounds_numerator":4,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear spellcasting-focus.druid spellcasting-focus.ranger trait.carved-ribboned-or-painted-to-channel-primal-magic trait.also-quarterstaff"},
		{"Druidic Focus/Yew wand", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["carved-ribboned-or-painted-to-channel-primal-magic"],"contents":[],"spellcasting_focus":["druid","ranger"],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear spellcasting-focus.druid spellcasting-focus.ranger trait.carved-ribboned-or-painted-to-channel-primal-magic"},
		{"Dungeoneer’s Pack/Dungeoneer’s Pack", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":[],"contents":[{"item":"backpack","quantity":1,"unit":"item"},{"item":"caltrops","quantity":1,"unit":"item"},{"item":"crowbar","quantity":1,"unit":"item"},{"item":"oil","quantity":2,"unit":"flask"},{"item":"rations","quantity":10,"unit":"day"},{"item":"rope","quantity":1,"unit":"item"},{"item":"tinderbox","quantity":1,"unit":"item"},{"item":"torch","quantity":10,"unit":"item"},{"item":"waterskin","quantity":1,"unit":"item"}],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":55,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear content.backpack-1-item content.caltrops-1-item content.crowbar-1-item content.oil-2-flask content.rations-10-day content.rope-1-item content.tinderbox-1-item content.torch-10-item content.waterskin-1-item"},
		{"Entertainer’s Pack/Entertainer’s Pack", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":[],"contents":[{"item":"backpack","quantity":1,"unit":"item"},{"item":"bedroll","quantity":1,"unit":"item"},{"item":"bell","quantity":1,"unit":"item"},{"item":"lantern-bullseye","quantity":1,"unit":"item"},{"item":"costume","quantity":3,"unit":"item"},{"item":"mirror","quantity":1,"unit":"item"},{"item":"oil","quantity":8,"unit":"flask"},{"item":"rations","quantity":9,"unit":"day"},{"item":"tinderbox","quantity":1,"unit":"item"},{"item":"waterskin","quantity":1,"unit":"item"}],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":117,"pounds_denominator":2,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear content.backpack-1-item content.bedroll-1-item content.bell-1-item content.lantern-bullseye-1-item content.costume-3-item content.mirror-1-item content.oil-8-flask content.rations-9-day content.tinderbox-1-item content.waterskin-1-item"},
		{"Explorer’s Pack/Explorer’s Pack", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":[],"contents":[{"item":"backpack","quantity":1,"unit":"item"},{"item":"bedroll","quantity":1,"unit":"item"},{"item":"oil","quantity":2,"unit":"flask"},{"item":"rations","quantity":10,"unit":"day"},{"item":"rope","quantity":1,"unit":"item"},{"item":"tinderbox","quantity":1,"unit":"item"},{"item":"torch","quantity":10,"unit":"item"},{"item":"waterskin","quantity":1,"unit":"item"}],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":55,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear content.backpack-1-item content.bedroll-1-item content.oil-2-flask content.rations-10-day content.rope-1-item content.tinderbox-1-item content.torch-10-item content.waterskin-1-item"},
		{"Flask/Flask", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"capacity","numerator":1,"denominator":1,"unit":"pint"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear measure.capacity-1-pint"},
		{"Grappling Hook/Grappling Hook", `{"gear":{"actions":["utilize"],"checks":[{"purpose":"catch-hook","kind":"ability-check","ability":"dexterity","skill":"acrobatics","dc":13,"dc_rule":"fixed"}],"damage":[],"measures":[{"name":"throw-range","numerator":50,"denominator":1,"unit":"feet"}],"traits":["throw-at-railing-ledge-or-other-catch","tied-rope-can-then-be-climbed"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":4,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"utilize", "item.gear action.utilize ability-check.catch-hook-dexterity-acrobatics-dc-13 measure.throw-range-50-feet trait.throw-at-railing-ledge-or-other-catch trait.tied-rope-can-then-be-climbed"},
	}
	index := 0
	for _, f := range loadItems258(t).Families {
		for _, v := range f.Variants {
			if index >= len(expected) {
				t.Fatal("batch holds more records than the fixture")
			}
			want := expected[index]
			index++
			t.Run(want.name, func(t *testing.T) {
				var details ItemDetails258
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

func TestItems258CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems258(t)
	composed, err := ComposeItemSourceIdentities(append(itemIdentitiesThrough(t, 257), batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, f := range composed {
		total += len(f.VariantIDs)
	}
	if len(composed) != 370 || total != 468 || composed[0].Name != "Potion of Healing" || composed[369].Name != "Grappling Hook" {
		t.Fatalf("expected Potion of Healing, batches #225 to #257 and #258: %d families, %d variants", len(composed), total)
	}
	checkItemBatchComposition(t, batch.SourceIdentities(), "adventuring-gear", "items_258")
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems258RejectMissingFields(t *testing.T) {
	rejectMissingItemMembers(t, items258JSON, 810, func(data []byte) error {
		_, err := DecodeItems258(data)
		return err
	})
}

func TestItems258RejectInvalidFixtures(t *testing.T) {
	rejectInvalidGearFixtures(t, LoadItems258, DecodeItems258, itemGear258)
}

// The mechanic keys restate every payload value except weight and the last
// source page, so changing any other value of any record fails.
func TestItems258RejectChangedValues(t *testing.T) {
	rejectChangedGearValues(t, items258JSON, func(data []byte) error {
		_, err := DecodeItems258(data)
		return err
	})
}

func TestItems258RejectUnsupportedJSON(t *testing.T) {
	decode := func(data []byte) error {
		_, err := DecodeItems258(data)
		return err
	}
	rejectUnsupportedItemJSON(t, items258JSON, `"id": "srd-5.2.1:content/items/clothes-fine/variant/standard"`, decode, []itemJSONChange{
		{"unknown payload", `"gear": {`, `"gear_plus": {`},
		{"unknown detail", `"spellcasting_focus": `, `"proficiency": `},
		{"number for boolean", `"listed": true`, `"listed": 1`},
		{"string for number", `"pounds_numerator": `, `"pounds_numerator": "1", "x": `},
		{"null check list", `"checks": [`, `"checks": null, "x": [`},
	})
	// Neither the Healing contract decoder nor the previous batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items258JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems257(items258JSON); err == nil {
		t.Fatal("#257 decoder accepted the #258 batch")
	}
	if _, err := DecodeItems258(items257JSON); err == nil {
		t.Fatal("#258 decoder accepted the #257 batch")
	}
}

// The Fighter creation fixture names the Dungeoneer’s Pack, and the character
// source catalogues refer to these keys.
func TestItems258JoinRuntimeGear(t *testing.T) {
	families := loadItems258(t).Families
	joinCatalogueEquipment(t, families, []itemGearKeyWant{
		{"crowbar", "crowbar", "standard"},
		{"travelers-clothes", "clothes-travelers", "standard"},
		{"dungeoneers-pack", "dungeoneers-pack", "standard"},
		{"entertainers-pack", "entertainers-pack", "standard"},
		{"explorers-pack", "explorers-pack", "standard"},
		{"druidic-focus-quarterstaff", "druidic-focus", "wooden-staff"},
		{"druidic-focus-sprig-of-mistletoe", "druidic-focus", "sprig-of-mistletoe"},
	})
	build, err := SupportedBuildForConcept("human fighter", "soldier")
	if err != nil {
		t.Fatal(err)
	}
	requireCreationNames(t, families, build, map[string]string{"Dungeoneer's Pack": "Dungeoneer's Pack"})
}
