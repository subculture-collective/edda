package srd521

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func loadItems256(t *testing.T) ItemSourceBatch[ItemDetails256] {
	t.Helper()
	b, err := LoadItems256()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities, variants and listed prices from pinned PDF pp. 94–96 for issue #256.
func TestItems256SourceIdentities(t *testing.T) {
	expected := []itemFamilyWant{
		{"Acid", "acid", 94, []itemVariantWant{{"standard", "Acid", 25, "gp"}}},
		{"Alchemist’s Fire", "alchemists-fire", 94, []itemVariantWant{{"standard", "Alchemist’s Fire", 50, "gp"}}},
		{"Ammunition", "ammunition", 94, []itemVariantWant{{"arrows", "Arrows", 1, "gp"}, {"bolts", "Bolts", 1, "gp"}, {"bullets-firearm", "Bullets, Firearm", 3, "gp"}, {"bullets-sling", "Bullets, Sling", 4, "cp"}, {"needles", "Needles", 1, "gp"}}},
		{"Antitoxin", "antitoxin", 96, []itemVariantWant{{"standard", "Antitoxin", 50, "gp"}}},
		{"Arcane Focus", "arcane-focus", 96, []itemVariantWant{{"crystal", "Crystal", 10, "gp"}, {"orb", "Orb", 20, "gp"}, {"rod", "Rod", 10, "gp"}, {"staff", "Staff", 5, "gp"}, {"wand", "Wand", 10, "gp"}}},
		{"Backpack", "backpack", 96, []itemVariantWant{{"standard", "Backpack", 2, "gp"}}},
		{"Ball Bearings", "ball-bearings", 96, []itemVariantWant{{"standard", "Ball Bearings", 1, "gp"}}},
		{"Barrel", "barrel", 96, []itemVariantWant{{"standard", "Barrel", 2, "gp"}}},
		{"Basket", "basket", 96, []itemVariantWant{{"standard", "Basket", 4, "sp"}}},
		{"Bedroll", "bedroll", 96, []itemVariantWant{{"standard", "Bedroll", 1, "gp"}}},
		{"Bell", "bell", 96, []itemVariantWant{{"standard", "Bell", 1, "gp"}}},
		{"Blanket", "blanket", 96, []itemVariantWant{{"standard", "Blanket", 5, "sp"}}},
	}
	b := loadItems256(t)
	if len(b.Families) != 12 || b.Batch != "data-adventuring-gear-01" || b.CatalogueIssue != 256 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the new gear families of issue #256")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 94–96") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	checkMundaneIdentities(t, b.Families, "adventuring-gear", "adventuring-gear", expected, map[string]int{"ammunition": 96})
	// Join every family to the decision inventory the issue cites, in its order.
	joinItemInventory(t, b.SourceIdentities(), b.Batch, 256, "adventuring-gear", []int{94, 94, 94, 96, 96, 96, 96, 96, 96, 96, 96, 96})
	b.Families[0].Variants[0].Details.Gear.SourcePages[0] = 1
	if fresh := loadItems256(t); fresh.Families[0].Variants[0].Details.Gear.SourcePages[0] != 94 {
		t.Fatal("load must return independent data")
	}
}

// Literal values, use modes and mechanic keys for all 20 records,
// transcribed from pinned PDF pp. 94–96.
func TestItems256SourceValues(t *testing.T) {
	expected := []struct{ name, details, use, keys string }{
		{"Acid/Acid", `{"gear":{"actions":["attack-replacement"],"checks":[{"purpose":"avoid-damage","kind":"saving-throw","ability":"dexterity","skill":"none","dc":8,"dc_rule":"plus-dexterity-modifier-and-proficiency-bonus"}],"damage":[{"purpose":"thrown-vial","dice":2,"sides":6,"flat":0,"type":"acid"}],"measures":[{"name":"throw-range","numerator":20,"denominator":1,"unit":"feet"}],"traits":["throw-vial-in-place-of-one-attack","target-one-creature-or-object-you-can-see"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[94,94]}}`,
			"attack-replacement", "item.gear action.attack-replacement saving-throw.avoid-damage-dexterity-dc-8-plus-dexterity-modifier-and-proficiency-bonus damage.thrown-vial-2d6-acid measure.throw-range-20-feet trait.throw-vial-in-place-of-one-attack trait.target-one-creature-or-object-you-can-see"},
		{"Alchemist’s Fire/Alchemist’s Fire", `{"gear":{"actions":["attack-replacement"],"checks":[{"purpose":"avoid-damage-and-burning","kind":"saving-throw","ability":"dexterity","skill":"none","dc":8,"dc_rule":"plus-dexterity-modifier-and-proficiency-bonus"}],"damage":[{"purpose":"thrown-flask","dice":1,"sides":4,"flat":0,"type":"fire"}],"measures":[{"name":"throw-range","numerator":20,"denominator":1,"unit":"feet"}],"traits":["throw-flask-in-place-of-one-attack","target-one-creature-or-object-you-can-see","failed-save-starts-burning"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[94,94]}}`,
			"attack-replacement", "item.gear action.attack-replacement saving-throw.avoid-damage-and-burning-dexterity-dc-8-plus-dexterity-modifier-and-proficiency-bonus damage.thrown-flask-1d4-fire measure.throw-range-20-feet trait.throw-flask-in-place-of-one-attack trait.target-one-creature-or-object-you-can-see trait.failed-save-starts-burning"},
		{"Ammunition/Arrows", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"bundle-amount","numerator":20,"denominator":1,"unit":"pieces"}],"traits":["required-by-weapon-with-ammunition-property","weapon-description-specifies-ammunition-type","bought-in-listed-amounts","storage-bought-separately","stored-in-quiver"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[94,96]}}`,
			"none", "item.gear measure.bundle-amount-20-pieces trait.required-by-weapon-with-ammunition-property trait.weapon-description-specifies-ammunition-type trait.bought-in-listed-amounts trait.storage-bought-separately trait.stored-in-quiver"},
		{"Ammunition/Bolts", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"bundle-amount","numerator":20,"denominator":1,"unit":"pieces"}],"traits":["required-by-weapon-with-ammunition-property","weapon-description-specifies-ammunition-type","bought-in-listed-amounts","storage-bought-separately","stored-in-case"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":3,"pounds_denominator":2,"condition":"none"},"source_pages":[94,96]}}`,
			"none", "item.gear measure.bundle-amount-20-pieces trait.required-by-weapon-with-ammunition-property trait.weapon-description-specifies-ammunition-type trait.bought-in-listed-amounts trait.storage-bought-separately trait.stored-in-case"},
		{"Ammunition/Bullets, Firearm", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"bundle-amount","numerator":10,"denominator":1,"unit":"pieces"}],"traits":["required-by-weapon-with-ammunition-property","weapon-description-specifies-ammunition-type","bought-in-listed-amounts","storage-bought-separately","stored-in-pouch"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1,"condition":"none"},"source_pages":[94,96]}}`,
			"none", "item.gear measure.bundle-amount-10-pieces trait.required-by-weapon-with-ammunition-property trait.weapon-description-specifies-ammunition-type trait.bought-in-listed-amounts trait.storage-bought-separately trait.stored-in-pouch"},
		{"Ammunition/Bullets, Sling", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"bundle-amount","numerator":20,"denominator":1,"unit":"pieces"}],"traits":["required-by-weapon-with-ammunition-property","weapon-description-specifies-ammunition-type","bought-in-listed-amounts","storage-bought-separately","stored-in-pouch"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":3,"pounds_denominator":2,"condition":"none"},"source_pages":[94,96]}}`,
			"none", "item.gear measure.bundle-amount-20-pieces trait.required-by-weapon-with-ammunition-property trait.weapon-description-specifies-ammunition-type trait.bought-in-listed-amounts trait.storage-bought-separately trait.stored-in-pouch"},
		{"Ammunition/Needles", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"bundle-amount","numerator":50,"denominator":1,"unit":"pieces"}],"traits":["required-by-weapon-with-ammunition-property","weapon-description-specifies-ammunition-type","bought-in-listed-amounts","storage-bought-separately","stored-in-pouch"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[94,96]}}`,
			"none", "item.gear measure.bundle-amount-50-pieces trait.required-by-weapon-with-ammunition-property trait.weapon-description-specifies-ammunition-type trait.bought-in-listed-amounts trait.storage-bought-separately trait.stored-in-pouch"},
		{"Antitoxin/Antitoxin", `{"gear":{"actions":["bonus-action"],"checks":[],"damage":[],"measures":[{"name":"benefit-duration","numerator":1,"denominator":1,"unit":"hour"}],"traits":["drink-one-vial","advantage-on-saves-to-avoid-or-end-poisoned"],"contents":[],"spellcasting_focus":[],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"bonus-action", "item.gear action.bonus-action measure.benefit-duration-1-hour trait.drink-one-vial trait.advantage-on-saves-to-avoid-or-end-poisoned"},
		{"Arcane Focus/Crystal", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["bejeweled-or-carved-to-channel-arcane-magic"],"contents":[],"spellcasting_focus":["sorcerer","warlock","wizard"],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear spellcasting-focus.sorcerer spellcasting-focus.warlock spellcasting-focus.wizard trait.bejeweled-or-carved-to-channel-arcane-magic"},
		{"Arcane Focus/Orb", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["bejeweled-or-carved-to-channel-arcane-magic"],"contents":[],"spellcasting_focus":["sorcerer","warlock","wizard"],"weight":{"listed":true,"pounds_numerator":3,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear spellcasting-focus.sorcerer spellcasting-focus.warlock spellcasting-focus.wizard trait.bejeweled-or-carved-to-channel-arcane-magic"},
		{"Arcane Focus/Rod", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["bejeweled-or-carved-to-channel-arcane-magic"],"contents":[],"spellcasting_focus":["sorcerer","warlock","wizard"],"weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear spellcasting-focus.sorcerer spellcasting-focus.warlock spellcasting-focus.wizard trait.bejeweled-or-carved-to-channel-arcane-magic"},
		{"Arcane Focus/Staff", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["bejeweled-or-carved-to-channel-arcane-magic","also-quarterstaff"],"contents":[],"spellcasting_focus":["sorcerer","warlock","wizard"],"weight":{"listed":true,"pounds_numerator":4,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear spellcasting-focus.sorcerer spellcasting-focus.warlock spellcasting-focus.wizard trait.bejeweled-or-carved-to-channel-arcane-magic trait.also-quarterstaff"},
		{"Arcane Focus/Wand", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["bejeweled-or-carved-to-channel-arcane-magic"],"contents":[],"spellcasting_focus":["sorcerer","warlock","wizard"],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear spellcasting-focus.sorcerer spellcasting-focus.warlock spellcasting-focus.wizard trait.bejeweled-or-carved-to-channel-arcane-magic"},
		{"Backpack/Backpack", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"capacity-weight","numerator":30,"denominator":1,"unit":"pounds"},{"name":"capacity-volume","numerator":1,"denominator":1,"unit":"cubic-foot"}],"traits":["can-serve-as-saddlebag"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":5,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear measure.capacity-weight-30-pounds measure.capacity-volume-1-cubic-foot trait.can-serve-as-saddlebag"},
		{"Ball Bearings/Ball Bearings", `{"gear":{"actions":["utilize"],"checks":[{"purpose":"avoid-prone","kind":"saving-throw","ability":"dexterity","skill":"none","dc":10,"dc_rule":"fixed"}],"damage":[],"measures":[{"name":"covered-square-side","numerator":10,"denominator":1,"unit":"feet"},{"name":"placement-range","numerator":10,"denominator":1,"unit":"feet"},{"name":"recovery-time","numerator":10,"denominator":1,"unit":"minutes"}],"traits":["spill-from-pouch-over-level-area","triggers-on-first-entry-on-a-turn","failed-save-causes-prone"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"utilize", "item.gear action.utilize saving-throw.avoid-prone-dexterity-dc-10 measure.covered-square-side-10-feet measure.placement-range-10-feet measure.recovery-time-10-minutes trait.spill-from-pouch-over-level-area trait.triggers-on-first-entry-on-a-turn trait.failed-save-causes-prone"},
		{"Barrel/Barrel", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"liquid-capacity","numerator":40,"denominator":1,"unit":"gallons"},{"name":"dry-goods-capacity","numerator":4,"denominator":1,"unit":"cubic-feet"}],"traits":["holds-liquid-or-dry-goods"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":70,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear measure.liquid-capacity-40-gallons measure.dry-goods-capacity-4-cubic-feet trait.holds-liquid-or-dry-goods"},
		{"Basket/Basket", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"capacity-weight","numerator":40,"denominator":1,"unit":"pounds"},{"name":"capacity-volume","numerator":2,"denominator":1,"unit":"cubic-feet"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear measure.capacity-weight-40-pounds measure.capacity-volume-2-cubic-feet"},
		{"Bedroll/Bedroll", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"occupants","numerator":1,"denominator":1,"unit":"creature"}],"traits":["occupant-is-small-or-medium","automatic-success-on-saves-against-extreme-cold-while-inside"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":7,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear measure.occupants-1-creature trait.occupant-is-small-or-medium trait.automatic-success-on-saves-against-extreme-cold-while-inside"},
		{"Bell/Bell", `{"gear":{"actions":["utilize"],"checks":[],"damage":[],"measures":[{"name":"audible-distance","numerator":60,"denominator":1,"unit":"feet"}],"traits":["rung-to-produce-sound"],"contents":[],"spellcasting_focus":[],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"utilize", "item.gear action.utilize measure.audible-distance-60-feet trait.rung-to-produce-sound"},
		{"Blanket/Blanket", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["advantage-on-saves-against-extreme-cold-while-wrapped"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":3,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear trait.advantage-on-saves-against-extreme-cold-while-wrapped"},
	}
	index := 0
	for _, f := range loadItems256(t).Families {
		for _, v := range f.Variants {
			if index >= len(expected) {
				t.Fatal("batch holds more records than the fixture")
			}
			want := expected[index]
			index++
			t.Run(want.name, func(t *testing.T) {
				var details ItemDetails256
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

func TestItems256CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems256(t)
	composed, err := ComposeItemSourceIdentities(append(itemIdentitiesThrough(t, 255), batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, f := range composed {
		total += len(f.VariantIDs)
	}
	if len(composed) != 346 || total != 442 || composed[0].Name != "Potion of Healing" || composed[345].Name != "Blanket" {
		t.Fatalf("expected Potion of Healing, batches #225 to #255 and #256: %d families, %d variants", len(composed), total)
	}
	checkItemBatchComposition(t, batch.SourceIdentities(), "adventuring-gear", "items_256")
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems256RejectMissingFields(t *testing.T) {
	rejectMissingItemMembers(t, items256JSON, 1014, func(data []byte) error {
		_, err := DecodeItems256(data)
		return err
	})
}

func TestItems256RejectInvalidFixtures(t *testing.T) {
	rejectInvalidGearFixtures(t, LoadItems256, DecodeItems256, itemGear256)
}

// The mechanic keys restate every payload value except weight and the last
// source page, so changing any other value of any record fails.
func TestItems256RejectChangedValues(t *testing.T) {
	rejectChangedGearValues(t, items256JSON, func(data []byte) error {
		_, err := DecodeItems256(data)
		return err
	})
}

func TestItems256RejectUnsupportedJSON(t *testing.T) {
	decode := func(data []byte) error {
		_, err := DecodeItems256(data)
		return err
	}
	rejectUnsupportedItemJSON(t, items256JSON, `"id": "srd-5.2.1:content/items/acid/variant/standard"`, decode, []itemJSONChange{
		{"unknown payload", `"gear": {`, `"gear_plus": {`},
		{"unknown detail", `"spellcasting_focus": `, `"proficiency": `},
		{"number for boolean", `"listed": true`, `"listed": 1`},
		{"string for number", `"pounds_numerator": `, `"pounds_numerator": "1", "x": `},
		{"null check list", `"checks": [`, `"checks": null, "x": [`},
	})
	// Neither the Healing contract decoder nor the previous batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items256JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems255(items256JSON); err == nil {
		t.Fatal("#255 decoder accepted the #256 batch")
	}
	if _, err := DecodeItems256(items255JSON); err == nil {
		t.Fatal("#256 decoder accepted the #255 batch")
	}
}

// The character source catalogues refer to these ammunition and focus keys.
func TestItems256JoinCatalogueEquipment(t *testing.T) {
	families := loadItems256(t).Families
	joinCatalogueEquipment(t, families, []itemGearKeyWant{
		{"arrow", "ammunition", "arrows"},
		{"arrows", "ammunition", "arrows"},
		{"arcane-focus-crystal", "arcane-focus", "crystal"},
		{"arcane-focus-orb", "arcane-focus", "orb"},
		{"arcane-focus-quarterstaff", "arcane-focus", "staff"},
	})
}
