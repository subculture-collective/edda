package srd521

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func loadItems257(t *testing.T) ItemSourceBatch[ItemDetails257] {
	t.Helper()
	b, err := LoadItems257()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities, variants and listed prices from pinned PDF pp. 96–97 for issue #257.
func TestItems257SourceIdentities(t *testing.T) {
	expected := []itemFamilyWant{
		{"Block and Tackle", "block-and-tackle", 96, []itemVariantWant{{"standard", "Block and Tackle", 1, "gp"}}},
		{"Book", "book", 96, []itemVariantWant{{"standard", "Book", 25, "gp"}}},
		{"Bottle, Glass", "bottle-glass", 96, []itemVariantWant{{"standard", "Bottle, Glass", 2, "gp"}}},
		{"Bucket", "bucket", 96, []itemVariantWant{{"standard", "Bucket", 5, "cp"}}},
		{"Burglar’s Pack", "burglars-pack", 96, []itemVariantWant{{"standard", "Burglar’s Pack", 16, "gp"}}},
		{"Caltrops", "caltrops", 96, []itemVariantWant{{"standard", "Caltrops", 1, "gp"}}},
		{"Candle", "candle", 96, []itemVariantWant{{"standard", "Candle", 1, "cp"}}},
		{"Case, Crossbow Bolt", "case-crossbow-bolt", 96, []itemVariantWant{{"standard", "Case, Crossbow Bolt", 1, "gp"}}},
		{"Case, Map or Scroll", "case-map-or-scroll", 96, []itemVariantWant{{"standard", "Case, Map or Scroll", 1, "gp"}}},
		{"Chain", "chain", 96, []itemVariantWant{{"standard", "Chain", 5, "gp"}}},
		{"Chest", "chest", 97, []itemVariantWant{{"standard", "Chest", 5, "gp"}}},
		{"Climber’s Kit", "climbers-kit", 97, []itemVariantWant{{"standard", "Climber’s Kit", 25, "gp"}}},
	}
	b := loadItems257(t)
	if len(b.Families) != 12 || b.Batch != "data-adventuring-gear-02" || b.CatalogueIssue != 257 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the new gear families of issue #257")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 96–97") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	checkMundaneIdentities(t, b.Families, "adventuring-gear", "adventuring-gear", expected)
	// Join every family to the decision inventory the issue cites, in its order.
	joinItemInventory(t, b.SourceIdentities(), b.Batch, 257, "adventuring-gear", []int{96, 96, 96, 96, 96, 96, 96, 96, 96, 96, 97, 97})
	b.Families[0].Variants[0].Details.Gear.SourcePages[0] = 1
	if fresh := loadItems257(t); fresh.Families[0].Variants[0].Details.Gear.SourcePages[0] != 96 {
		t.Fatal("load must return independent data")
	}
}

// Literal values, use modes and mechanic keys for all 12 records,
// transcribed from pinned PDF pp. 96–97.
func TestItems257SourceValues(t *testing.T) {
	expected := []struct{ name, details, use, keys string }{
		{"Block and Tackle/Block and Tackle", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"hoist-multiplier","numerator":4,"denominator":1,"unit":"times-normal-lift"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":5,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear measure.hoist-multiplier-4-times-normal-lift"},
		{"Book/Book", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"knowledge-check-bonus","numerator":5,"denominator":1,"unit":"bonus"}],"traits":["contains-fiction-or-nonfiction","bonus-requires-consulting-accurate-nonfiction-about-topic","bonus-to-intelligence-arcana-history-nature-or-religion-checks"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":5,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear measure.knowledge-check-bonus-5-bonus trait.contains-fiction-or-nonfiction trait.bonus-requires-consulting-accurate-nonfiction-about-topic trait.bonus-to-intelligence-arcana-history-nature-or-religion-checks"},
		{"Bottle, Glass/Bottle, Glass", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"capacity","numerator":3,"denominator":2,"unit":"pints"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear measure.capacity-3-over-2-pints"},
		{"Bucket/Bucket", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"capacity-volume","numerator":1,"denominator":2,"unit":"cubic-foot"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear measure.capacity-volume-1-over-2-cubic-foot"},
		{"Burglar’s Pack/Burglar’s Pack", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":[],"contents":[{"item":"backpack","quantity":1,"unit":"item"},{"item":"ball-bearings","quantity":1,"unit":"item"},{"item":"bell","quantity":1,"unit":"item"},{"item":"candle","quantity":10,"unit":"item"},{"item":"crowbar","quantity":1,"unit":"item"},{"item":"lantern-hooded","quantity":1,"unit":"item"},{"item":"oil","quantity":7,"unit":"flask"},{"item":"rations","quantity":5,"unit":"day"},{"item":"rope","quantity":1,"unit":"item"},{"item":"tinderbox","quantity":1,"unit":"item"},{"item":"waterskin","quantity":1,"unit":"item"}],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":42,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear content.backpack-1-item content.ball-bearings-1-item content.bell-1-item content.candle-10-item content.crowbar-1-item content.lantern-hooded-1-item content.oil-7-flask content.rations-5-day content.rope-1-item content.tinderbox-1-item content.waterskin-1-item"},
		{"Caltrops/Caltrops", `{"gear":{"actions":["utilize"],"checks":[{"purpose":"avoid-damage-and-halt","kind":"saving-throw","ability":"dexterity","skill":"none","dc":15,"dc_rule":"fixed"}],"damage":[{"purpose":"failed-save","dice":0,"sides":0,"flat":1,"type":"piercing"}],"measures":[{"name":"covered-square-side","numerator":5,"denominator":1,"unit":"feet"},{"name":"placement-range","numerator":5,"denominator":1,"unit":"feet"},{"name":"failed-save-speed","numerator":0,"denominator":1,"unit":"speed"},{"name":"recovery-time","numerator":10,"denominator":1,"unit":"minutes"}],"traits":["spread-from-bag","triggers-on-first-entry-on-a-turn","speed-penalty-lasts-until-start-of-next-turn"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"utilize", "item.gear action.utilize saving-throw.avoid-damage-and-halt-dexterity-dc-15 damage.failed-save-1-piercing measure.covered-square-side-5-feet measure.placement-range-5-feet measure.failed-save-speed-0-speed measure.recovery-time-10-minutes trait.spread-from-bag trait.triggers-on-first-entry-on-a-turn trait.speed-penalty-lasts-until-start-of-next-turn"},
		{"Candle/Candle", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"burn-duration","numerator":1,"denominator":1,"unit":"hour"},{"name":"bright-light-radius","numerator":5,"denominator":1,"unit":"feet"},{"name":"additional-dim-light","numerator":5,"denominator":1,"unit":"feet"}],"traits":["sheds-light-while-lit"],"contents":[],"spellcasting_focus":[],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear measure.burn-duration-1-hour measure.bright-light-radius-5-feet measure.additional-dim-light-5-feet trait.sheds-light-while-lit"},
		{"Case, Crossbow Bolt/Case, Crossbow Bolt", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"capacity","numerator":20,"denominator":1,"unit":"bolts"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear measure.capacity-20-bolts"},
		{"Case, Map or Scroll/Case, Map or Scroll", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"paper-capacity","numerator":10,"denominator":1,"unit":"sheets"},{"name":"parchment-capacity","numerator":5,"denominator":1,"unit":"sheets"}],"traits":["holds-paper-or-parchment"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"none", "item.gear measure.paper-capacity-10-sheets measure.parchment-capacity-5-sheets trait.holds-paper-or-parchment"},
		{"Chain/Chain", `{"gear":{"actions":["utilize","action"],"checks":[{"purpose":"wrap-creature","kind":"ability-check","ability":"strength","skill":"athletics","dc":13,"dc_rule":"fixed"},{"purpose":"escape","kind":"ability-check","ability":"dexterity","skill":"acrobatics","dc":18,"dc_rule":"fixed"},{"purpose":"burst","kind":"ability-check","ability":"strength","skill":"athletics","dc":20,"dc_rule":"fixed"}],"damage":[],"measures":[{"name":"wrap-range","numerator":5,"denominator":1,"unit":"feet"}],"traits":["target-unwilling-grappled-incapacitated-or-restrained-creature","bound-legs-restrain-until-escape"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":10,"pounds_denominator":1,"condition":"none"},"source_pages":[96,96]}}`,
			"multiple", "item.gear action.utilize action.action ability-check.wrap-creature-strength-athletics-dc-13 ability-check.escape-dexterity-acrobatics-dc-18 ability-check.burst-strength-athletics-dc-20 measure.wrap-range-5-feet trait.target-unwilling-grappled-incapacitated-or-restrained-creature trait.bound-legs-restrain-until-escape"},
		{"Chest/Chest", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"capacity-volume","numerator":12,"denominator":1,"unit":"cubic-feet"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":25,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear measure.capacity-volume-12-cubic-feet"},
		{"Climber’s Kit/Climber’s Kit", `{"gear":{"actions":["utilize","bonus-action"],"checks":[],"damage":[],"measures":[{"name":"maximum-fall-from-anchor","numerator":25,"denominator":1,"unit":"feet"},{"name":"maximum-movement-from-anchor","numerator":25,"denominator":1,"unit":"feet"}],"traits":["includes-boot-tips-gloves-pitons-and-harness","anchor-yourself","moving-farther-requires-undoing-anchor"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":12,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"multiple", "item.gear action.utilize action.bonus-action measure.maximum-fall-from-anchor-25-feet measure.maximum-movement-from-anchor-25-feet trait.includes-boot-tips-gloves-pitons-and-harness trait.anchor-yourself trait.moving-farther-requires-undoing-anchor"},
	}
	index := 0
	for _, f := range loadItems257(t).Families {
		for _, v := range f.Variants {
			if index >= len(expected) {
				t.Fatal("batch holds more records than the fixture")
			}
			want := expected[index]
			index++
			t.Run(want.name, func(t *testing.T) {
				var details ItemDetails257
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

func TestItems257CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems257(t)
	composed, err := ComposeItemSourceIdentities(append(itemIdentitiesThrough(t, 256), batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, f := range composed {
		total += len(f.VariantIDs)
	}
	if len(composed) != 358 || total != 454 || composed[0].Name != "Potion of Healing" || composed[357].Name != "Climber’s Kit" {
		t.Fatalf("expected Potion of Healing, batches #225 to #256 and #257: %d families, %d variants", len(composed), total)
	}
	checkItemBatchComposition(t, batch.SourceIdentities(), "adventuring-gear", "items_257")
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems257RejectMissingFields(t *testing.T) {
	rejectMissingItemMembers(t, items257JSON, 740, func(data []byte) error {
		_, err := DecodeItems257(data)
		return err
	})
}

func TestItems257RejectInvalidFixtures(t *testing.T) {
	rejectInvalidGearFixtures(t, LoadItems257, DecodeItems257, itemGear257)
}

// The mechanic keys restate every payload value except weight and the last
// source page, so changing any other value of any record fails.
func TestItems257RejectChangedValues(t *testing.T) {
	rejectChangedGearValues(t, items257JSON, func(data []byte) error {
		_, err := DecodeItems257(data)
		return err
	})
}

func TestItems257RejectUnsupportedJSON(t *testing.T) {
	decode := func(data []byte) error {
		_, err := DecodeItems257(data)
		return err
	}
	rejectUnsupportedItemJSON(t, items257JSON, `"id": "srd-5.2.1:content/items/block-and-tackle/variant/standard"`, decode, []itemJSONChange{
		{"unknown payload", `"gear": {`, `"gear_plus": {`},
		{"unknown detail", `"spellcasting_focus": `, `"proficiency": `},
		{"number for boolean", `"listed": true`, `"listed": 1`},
		{"string for number", `"pounds_numerator": `, `"pounds_numerator": "1", "x": `},
		{"null check list", `"checks": [`, `"checks": null, "x": [`},
	})
	// Neither the Healing contract decoder nor the previous batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items257JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems256(items257JSON); err == nil {
		t.Fatal("#256 decoder accepted the #257 batch")
	}
	if _, err := DecodeItems257(items256JSON); err == nil {
		t.Fatal("#257 decoder accepted the #256 batch")
	}
}

// The character source catalogues refer to these keys. Their three books
// have different contents but are all the Book entry.
func TestItems257JoinCatalogueEquipment(t *testing.T) {
	families := loadItems257(t).Families
	joinCatalogueEquipment(t, families, []itemGearKeyWant{
		{"burglars-pack", "burglars-pack", "standard"},
		{"book/prayers", "book", "standard"},
		{"book/history", "book", "standard"},
		{"book-occult-lore", "book", "standard"},
	})
}
