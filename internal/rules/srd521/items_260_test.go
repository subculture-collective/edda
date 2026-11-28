package srd521

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func loadItems260(t *testing.T) ItemSourceBatch[ItemDetails260] {
	t.Helper()
	b, err := LoadItems260()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities, variants and listed prices from pinned PDF pp. 98–99 for issue #260.
func TestItems260SourceIdentities(t *testing.T) {
	expected := []itemFamilyWant{
		{"Magnifying Glass", "magnifying-glass", 98, []itemVariantWant{{"standard", "Magnifying Glass", 100, "gp"}}},
		{"Manacles", "manacles", 98, []itemVariantWant{{"standard", "Manacles", 2, "gp"}}},
		{"Map", "map", 98, []itemVariantWant{{"standard", "Map", 1, "gp"}}},
		{"Mirror", "mirror", 98, []itemVariantWant{{"standard", "Mirror", 5, "gp"}}},
		{"Net", "net", 98, []itemVariantWant{{"standard", "Net", 1, "gp"}}},
		{"Oil", "oil", 98, []itemVariantWant{{"standard", "Oil", 1, "sp"}}},
		{"Paper", "paper", 99, []itemVariantWant{{"standard", "Paper", 2, "sp"}}},
		{"Parchment", "parchment", 99, []itemVariantWant{{"standard", "Parchment", 1, "sp"}}},
		{"Perfume", "perfume", 99, []itemVariantWant{{"standard", "Perfume", 5, "gp"}}},
		{"Poison, Basic", "poison-basic", 99, []itemVariantWant{{"standard", "Poison, Basic", 100, "gp"}}},
		{"Pole", "pole", 99, []itemVariantWant{{"standard", "Pole", 5, "cp"}}},
		{"Pot, Iron", "pot-iron", 99, []itemVariantWant{{"standard", "Pot, Iron", 2, "gp"}}},
	}
	b := loadItems260(t)
	if len(b.Families) != 12 || b.Batch != "data-adventuring-gear-05" || b.CatalogueIssue != 260 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the new gear families of issue #260")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 98–99") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	checkMundaneIdentities(t, b.Families, "adventuring-gear", "adventuring-gear", expected)
	// Join every family to the decision inventory the issue cites, in its order.
	joinItemInventory(t, b.SourceIdentities(), b.Batch, 260, "adventuring-gear", []int{98, 98, 98, 98, 98, 98, 99, 99, 99, 99, 99, 99})
	b.Families[0].Variants[0].Details.Gear.SourcePages[0] = 1
	if fresh := loadItems260(t); fresh.Families[0].Variants[0].Details.Gear.SourcePages[0] != 98 {
		t.Fatal("load must return independent data")
	}
}

// Literal values, use modes and mechanic keys for all 12 records,
// transcribed from pinned PDF pp. 98–99.
func TestItems260SourceValues(t *testing.T) {
	expected := []struct{ name, details, use, keys string }{
		{"Magnifying Glass/Magnifying Glass", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"fire-ignition-time","numerator":5,"denominator":1,"unit":"minutes"}],"traits":["advantage-on-checks-to-appraise-or-inspect-highly-detailed-item","lighting-fire-requires-sunlight-bright-light-and-tinder"],"contents":[],"spellcasting_focus":[],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[98,98]}}`,
			"none", "item.gear measure.fire-ignition-time-5-minutes trait.advantage-on-checks-to-appraise-or-inspect-highly-detailed-item trait.lighting-fire-requires-sunlight-bright-light-and-tinder"},
		{"Manacles/Manacles", `{"gear":{"actions":["utilize","action"],"checks":[{"purpose":"bind-creature","kind":"ability-check","ability":"dexterity","skill":"sleight-of-hand","dc":13,"dc_rule":"fixed"},{"purpose":"escape","kind":"ability-check","ability":"dexterity","skill":"sleight-of-hand","dc":20,"dc_rule":"fixed"},{"purpose":"burst","kind":"ability-check","ability":"strength","skill":"athletics","dc":25,"dc_rule":"fixed"},{"purpose":"pick-lock-without-key","kind":"ability-check","ability":"dexterity","skill":"sleight-of-hand","dc":15,"dc_rule":"fixed"}],"damage":[],"measures":[{"name":"bind-range","numerator":5,"denominator":1,"unit":"feet"}],"traits":["target-unwilling-small-or-medium-grappled-incapacitated-or-restrained-creature","bound-creature-has-disadvantage-on-attack-rolls","restrained-if-attached-to-fixed-chain-or-hook","comes-with-key","thieves-tools-can-pick-without-key"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":6,"pounds_denominator":1,"condition":"none"},"source_pages":[98,98]}}`,
			"multiple", "item.gear action.utilize action.action ability-check.bind-creature-dexterity-sleight-of-hand-dc-13 ability-check.escape-dexterity-sleight-of-hand-dc-20 ability-check.burst-strength-athletics-dc-25 ability-check.pick-lock-without-key-dexterity-sleight-of-hand-dc-15 measure.bind-range-5-feet trait.target-unwilling-small-or-medium-grappled-incapacitated-or-restrained-creature trait.bound-creature-has-disadvantage-on-attack-rolls trait.restrained-if-attached-to-fixed-chain-or-hook trait.comes-with-key trait.thieves-tools-can-pick-without-key"},
		{"Map/Map", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"wayfinding-check-bonus","numerator":5,"denominator":1,"unit":"bonus"}],"traits":["bonus-requires-consulting-accurate-map","bonus-to-wisdom-survival-checks-to-find-way-in-mapped-place"],"contents":[],"spellcasting_focus":[],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[98,98]}}`,
			"none", "item.gear measure.wayfinding-check-bonus-5-bonus trait.bonus-requires-consulting-accurate-map trait.bonus-to-wisdom-survival-checks-to-find-way-in-mapped-place"},
		{"Mirror/Mirror", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["handheld-steel","useful-for-cosmetics-peeking-around-corners-and-signals"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":2,"condition":"none"},"source_pages":[98,98]}}`,
			"none", "item.gear trait.handheld-steel trait.useful-for-cosmetics-peeking-around-corners-and-signals"},
		{"Net/Net", `{"gear":{"actions":["attack-replacement","action"],"checks":[{"purpose":"avoid-restrained","kind":"saving-throw","ability":"dexterity","skill":"none","dc":8,"dc_rule":"plus-dexterity-modifier-and-proficiency-bonus"},{"purpose":"escape","kind":"ability-check","ability":"strength","skill":"athletics","dc":10,"dc_rule":"fixed"}],"damage":[],"measures":[{"name":"throw-range","numerator":15,"denominator":1,"unit":"feet"},{"name":"escape-helper-range","numerator":5,"denominator":1,"unit":"feet"},{"name":"armor-class","numerator":10,"denominator":1,"unit":"armor-class"},{"name":"hit-points","numerator":5,"denominator":1,"unit":"hit-points"}],"traits":["throw-net-in-place-of-one-attack","target-one-creature-you-can-see","failed-save-restrains-until-escape","huge-or-larger-target-succeeds-automatically","successful-check-frees-restrained-creature","immune-to-bludgeoning-poison-and-psychic-damage","destroying-net-frees-target-and-ends-effect"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":3,"pounds_denominator":1,"condition":"none"},"source_pages":[98,98]}}`,
			"multiple", "item.gear action.attack-replacement action.action saving-throw.avoid-restrained-dexterity-dc-8-plus-dexterity-modifier-and-proficiency-bonus ability-check.escape-strength-athletics-dc-10 measure.throw-range-15-feet measure.escape-helper-range-5-feet measure.armor-class-10-armor-class measure.hit-points-5-hit-points trait.throw-net-in-place-of-one-attack trait.target-one-creature-you-can-see trait.failed-save-restrains-until-escape trait.huge-or-larger-target-succeeds-automatically trait.successful-check-frees-restrained-creature trait.immune-to-bludgeoning-poison-and-psychic-damage trait.destroying-net-frees-target-and-ends-effect"},
		{"Oil/Oil", `{"gear":{"actions":["attack-replacement","utilize"],"checks":[{"purpose":"avoid-being-covered-in-oil","kind":"saving-throw","ability":"dexterity","skill":"none","dc":8,"dc_rule":"plus-dexterity-modifier-and-proficiency-bonus"}],"damage":[{"purpose":"extra-from-burning-oil","dice":0,"sides":0,"flat":5,"type":"fire"},{"purpose":"burning-area","dice":0,"sides":0,"flat":5,"type":"fire"}],"measures":[{"name":"throw-range","numerator":20,"denominator":1,"unit":"feet"},{"name":"drying-time","numerator":1,"denominator":1,"unit":"minute"},{"name":"covered-square-side","numerator":5,"denominator":1,"unit":"feet"},{"name":"pour-range","numerator":5,"denominator":1,"unit":"feet"},{"name":"area-burn-rounds","numerator":2,"denominator":1,"unit":"rounds"},{"name":"area-burn-seconds","numerator":12,"denominator":1,"unit":"seconds"},{"name":"fuel-burn-duration","numerator":6,"denominator":1,"unit":"hours"}],"traits":["douse-creature-object-or-space-or-use-as-fuel","throw-flask-in-place-of-one-attack","target-one-creature-or-object","extra-damage-if-fire-damage-taken-before-oil-dries","pour-flask-on-level-ground","area-damages-on-entry-or-turn-end-once-per-turn","fuel-for-lamps-and-lanterns","fuel-duration-need-not-be-consecutive","extinguish-and-rekindle-until-total-duration-burned"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[98,99]}}`,
			"multiple", "item.gear action.attack-replacement action.utilize saving-throw.avoid-being-covered-in-oil-dexterity-dc-8-plus-dexterity-modifier-and-proficiency-bonus damage.extra-from-burning-oil-5-fire damage.burning-area-5-fire measure.throw-range-20-feet measure.drying-time-1-minute measure.covered-square-side-5-feet measure.pour-range-5-feet measure.area-burn-rounds-2-rounds measure.area-burn-seconds-12-seconds measure.fuel-burn-duration-6-hours trait.douse-creature-object-or-space-or-use-as-fuel trait.throw-flask-in-place-of-one-attack trait.target-one-creature-or-object trait.extra-damage-if-fire-damage-taken-before-oil-dries trait.pour-flask-on-level-ground trait.area-damages-on-entry-or-turn-end-once-per-turn trait.fuel-for-lamps-and-lanterns trait.fuel-duration-need-not-be-consecutive trait.extinguish-and-rekindle-until-total-duration-burned"},
		{"Paper/Paper", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"handwritten-words-per-sheet","numerator":250,"denominator":1,"unit":"words"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"none", "item.gear measure.handwritten-words-per-sheet-250-words"},
		{"Parchment/Parchment", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"handwritten-words-per-sheet","numerator":250,"denominator":1,"unit":"words"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"none", "item.gear measure.handwritten-words-per-sheet-250-words"},
		{"Perfume/Perfume", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"vial-volume","numerator":4,"denominator":1,"unit":"ounces"},{"name":"benefit-duration","numerator":1,"denominator":1,"unit":"hour"},{"name":"influence-range","numerator":5,"denominator":1,"unit":"feet"}],"traits":["advantage-on-charisma-persuasion-to-influence-indifferent-humanoid"],"contents":[],"spellcasting_focus":[],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"none", "item.gear measure.vial-volume-4-ounces measure.benefit-duration-1-hour measure.influence-range-5-feet trait.advantage-on-charisma-persuasion-to-influence-indifferent-humanoid"},
		{"Poison, Basic/Poison, Basic", `{"gear":{"actions":["bonus-action"],"checks":[],"damage":[{"purpose":"extra-from-poisoned-weapon-or-ammunition","dice":1,"sides":4,"flat":0,"type":"poison"}],"measures":[{"name":"coated-weapons","numerator":1,"denominator":1,"unit":"weapon"},{"name":"coated-ammunition","numerator":3,"denominator":1,"unit":"pieces"},{"name":"potency-duration","numerator":1,"denominator":1,"unit":"minute"}],"traits":["coat-one-weapon-or-ammunition-from-vial","applies-on-piercing-or-slashing-damage","potency-ends-early-when-damage-is-dealt"],"contents":[],"spellcasting_focus":[],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"bonus-action", "item.gear action.bonus-action damage.extra-from-poisoned-weapon-or-ammunition-1d4-poison measure.coated-weapons-1-weapon measure.coated-ammunition-3-pieces measure.potency-duration-1-minute trait.coat-one-weapon-or-ammunition-from-vial trait.applies-on-piercing-or-slashing-damage trait.potency-ends-early-when-damage-is-dealt"},
		{"Pole/Pole", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"length","numerator":10,"denominator":1,"unit":"feet"},{"name":"touch-range","numerator":10,"denominator":1,"unit":"feet"}],"traits":["advantage-when-vaulting-in-high-or-long-jump-athletics-check"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":7,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"none", "item.gear measure.length-10-feet measure.touch-range-10-feet trait.advantage-when-vaulting-in-high-or-long-jump-athletics-check"},
		{"Pot, Iron/Pot, Iron", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"capacity","numerator":1,"denominator":1,"unit":"gallon"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":10,"pounds_denominator":1,"condition":"none"},"source_pages":[99,99]}}`,
			"none", "item.gear measure.capacity-1-gallon"},
	}
	index := 0
	for _, f := range loadItems260(t).Families {
		for _, v := range f.Variants {
			if index >= len(expected) {
				t.Fatal("batch holds more records than the fixture")
			}
			want := expected[index]
			index++
			t.Run(want.name, func(t *testing.T) {
				var details ItemDetails260
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

func TestItems260CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems260(t)
	composed, err := ComposeItemSourceIdentities(append(itemIdentitiesThrough(t, 259), batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, f := range composed {
		total += len(f.VariantIDs)
	}
	if len(composed) != 394 || total != 494 || composed[0].Name != "Potion of Healing" || composed[393].Name != "Pot, Iron" {
		t.Fatalf("expected Potion of Healing, batches #225 to #259 and #260: %d families, %d variants", len(composed), total)
	}
	checkItemBatchComposition(t, batch.SourceIdentities(), "adventuring-gear", "items_260")
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems260RejectMissingFields(t *testing.T) {
	rejectMissingItemMembers(t, items260JSON, 763, func(data []byte) error {
		_, err := DecodeItems260(data)
		return err
	})
}

func TestItems260RejectInvalidFixtures(t *testing.T) {
	rejectInvalidGearFixtures(t, LoadItems260, DecodeItems260, itemGear260)
}

// The mechanic keys restate every payload value except weight and the last
// source page, so changing any other value of any record fails.
func TestItems260RejectChangedValues(t *testing.T) {
	rejectChangedGearValues(t, items260JSON, func(data []byte) error {
		_, err := DecodeItems260(data)
		return err
	})
}

func TestItems260RejectUnsupportedJSON(t *testing.T) {
	decode := func(data []byte) error {
		_, err := DecodeItems260(data)
		return err
	}
	rejectUnsupportedItemJSON(t, items260JSON, `"id": "srd-5.2.1:content/items/magnifying-glass/variant/standard"`, decode, []itemJSONChange{
		{"unknown payload", `"gear": {`, `"gear_plus": {`},
		{"unknown detail", `"spellcasting_focus": `, `"proficiency": `},
		{"number for boolean", `"listed": true`, `"listed": 1`},
		{"string for number", `"pounds_numerator": `, `"pounds_numerator": "1", "x": `},
		{"null check list", `"checks": [`, `"checks": null, "x": [`},
	})
	// Neither the Healing contract decoder nor the previous batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items260JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems259(items260JSON); err == nil {
		t.Fatal("#259 decoder accepted the #260 batch")
	}
	if _, err := DecodeItems260(items259JSON); err == nil {
		t.Fatal("#260 decoder accepted the #259 batch")
	}
}

// The Cleric creation fixture names ten sheets of Parchment, and the
// background catalogue counts sheets.
func TestItems260JoinRuntimeGear(t *testing.T) {
	families := loadItems260(t).Families
	joinCatalogueEquipment(t, families, []itemGearKeyWant{
		{"parchment/sheets", "parchment", "standard"},
	})
	build, err := SupportedBuildForConcept("human cleric", "acolyte")
	if err != nil {
		t.Fatal(err)
	}
	requireCreationNames(t, families, build, map[string]string{"Parchment (10)": "Parchment"})
}
