package srd521

import (
	"encoding/json"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func loadItems251(t *testing.T) ItemSourceBatch[ItemDetails251] {
	t.Helper()
	b, err := LoadItems251()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities, variants and listed prices from pinned PDF p. 92 for issue #251.
func TestItems251SourceIdentities(t *testing.T) {
	expected := []itemFamilyWant{
		{"Breastplate", "breastplate", 92, []itemVariantWant{{"standard", "Breastplate", 400, "gp"}}},
		{"Chain Mail", "chain-mail", 92, []itemVariantWant{{"standard", "Chain Mail", 75, "gp"}}},
		{"Chain Shirt", "chain-shirt", 92, []itemVariantWant{{"standard", "Chain Shirt", 50, "gp"}}},
		{"Half Plate Armor", "half-plate-armor", 92, []itemVariantWant{{"standard", "Half Plate Armor", 750, "gp"}}},
		{"Hide Armor", "hide-armor", 92, []itemVariantWant{{"standard", "Hide Armor", 10, "gp"}}},
		{"Leather Armor", "leather-armor", 92, []itemVariantWant{{"standard", "Leather Armor", 10, "gp"}}},
		{"Padded Armor", "padded-armor", 92, []itemVariantWant{{"standard", "Padded Armor", 5, "gp"}}},
		{"Plate Armor", "plate-armor", 92, []itemVariantWant{{"standard", "Plate Armor", 1500, "gp"}}},
		{"Ring Mail", "ring-mail", 92, []itemVariantWant{{"standard", "Ring Mail", 30, "gp"}}},
		{"Scale Mail", "scale-mail", 92, []itemVariantWant{{"standard", "Scale Mail", 50, "gp"}}},
		{"Shield", "shield", 92, []itemVariantWant{{"standard", "Shield", 10, "gp"}}},
		{"Splint Armor", "splint-armor", 92, []itemVariantWant{{"standard", "Splint Armor", 200, "gp"}}},
	}
	b := loadItems251(t)
	if len(b.Families) != 12 || b.Batch != "data-armor-01" || b.CatalogueIssue != 251 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the armor rows of issue #251")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "p. 92") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	checkMundaneIdentities(t, b.Families, "armor", "armor", expected)
	// Join every family to the decision inventory the issue cites, in its order.
	joinItemInventory(t, b.SourceIdentities(), b.Batch, 251, "armor", []int{92, 92, 92, 92, 92, 92, 92, 92, 92, 92, 92, 92})
	b.Families[0].Variants[0].Details.Armor.ArmorCategory = "changed"
	if fresh := loadItems251(t); fresh.Families[0].Variants[0].Details.Armor.ArmorCategory != "medium" {
		t.Fatal("load must return independent data")
	}
}

// Literal values, equip modes and mechanic keys for all 12 records,
// transcribed from pinned PDF p. 92.
func TestItems251SourceValues(t *testing.T) {
	expected := []struct{ name, details, equip, keys string }{
		{"Breastplate/Breastplate", `{"armor":{"armor_category":"medium","armor_class":{"base":14,"bonus":0,"dexterity":"capped","dexterity_maximum":2},"strength_requirement":0,"stealth_disadvantage":false,"don_doff":{"action":"none","don_minutes":5,"doff_minutes":1},"weight_pounds":20}}`,
			"wear", "item.wear armor.medium armor.training armor-class.14-plus-dexterity-max-2 don.minutes-5 doff.minutes-1"},
		{"Chain Mail/Chain Mail", `{"armor":{"armor_category":"heavy","armor_class":{"base":16,"bonus":0,"dexterity":"none","dexterity_maximum":0},"strength_requirement":13,"stealth_disadvantage":true,"don_doff":{"action":"none","don_minutes":10,"doff_minutes":5},"weight_pounds":55}}`,
			"wear", "item.wear armor.heavy armor.training armor-class.16 strength.13 stealth.disadvantage don.minutes-10 doff.minutes-5"},
		{"Chain Shirt/Chain Shirt", `{"armor":{"armor_category":"medium","armor_class":{"base":13,"bonus":0,"dexterity":"capped","dexterity_maximum":2},"strength_requirement":0,"stealth_disadvantage":false,"don_doff":{"action":"none","don_minutes":5,"doff_minutes":1},"weight_pounds":20}}`,
			"wear", "item.wear armor.medium armor.training armor-class.13-plus-dexterity-max-2 don.minutes-5 doff.minutes-1"},
		{"Half Plate Armor/Half Plate Armor", `{"armor":{"armor_category":"medium","armor_class":{"base":15,"bonus":0,"dexterity":"capped","dexterity_maximum":2},"strength_requirement":0,"stealth_disadvantage":true,"don_doff":{"action":"none","don_minutes":5,"doff_minutes":1},"weight_pounds":40}}`,
			"wear", "item.wear armor.medium armor.training armor-class.15-plus-dexterity-max-2 stealth.disadvantage don.minutes-5 doff.minutes-1"},
		{"Hide Armor/Hide Armor", `{"armor":{"armor_category":"medium","armor_class":{"base":12,"bonus":0,"dexterity":"capped","dexterity_maximum":2},"strength_requirement":0,"stealth_disadvantage":false,"don_doff":{"action":"none","don_minutes":5,"doff_minutes":1},"weight_pounds":12}}`,
			"wear", "item.wear armor.medium armor.training armor-class.12-plus-dexterity-max-2 don.minutes-5 doff.minutes-1"},
		{"Leather Armor/Leather Armor", `{"armor":{"armor_category":"light","armor_class":{"base":11,"bonus":0,"dexterity":"full","dexterity_maximum":0},"strength_requirement":0,"stealth_disadvantage":false,"don_doff":{"action":"none","don_minutes":1,"doff_minutes":1},"weight_pounds":10}}`,
			"wear", "item.wear armor.light armor.training armor-class.11-plus-dexterity don.minutes-1 doff.minutes-1"},
		{"Padded Armor/Padded Armor", `{"armor":{"armor_category":"light","armor_class":{"base":11,"bonus":0,"dexterity":"full","dexterity_maximum":0},"strength_requirement":0,"stealth_disadvantage":true,"don_doff":{"action":"none","don_minutes":1,"doff_minutes":1},"weight_pounds":8}}`,
			"wear", "item.wear armor.light armor.training armor-class.11-plus-dexterity stealth.disadvantage don.minutes-1 doff.minutes-1"},
		{"Plate Armor/Plate Armor", `{"armor":{"armor_category":"heavy","armor_class":{"base":18,"bonus":0,"dexterity":"none","dexterity_maximum":0},"strength_requirement":15,"stealth_disadvantage":true,"don_doff":{"action":"none","don_minutes":10,"doff_minutes":5},"weight_pounds":65}}`,
			"wear", "item.wear armor.heavy armor.training armor-class.18 strength.15 stealth.disadvantage don.minutes-10 doff.minutes-5"},
		{"Ring Mail/Ring Mail", `{"armor":{"armor_category":"heavy","armor_class":{"base":14,"bonus":0,"dexterity":"none","dexterity_maximum":0},"strength_requirement":0,"stealth_disadvantage":true,"don_doff":{"action":"none","don_minutes":10,"doff_minutes":5},"weight_pounds":40}}`,
			"wear", "item.wear armor.heavy armor.training armor-class.14 stealth.disadvantage don.minutes-10 doff.minutes-5"},
		{"Scale Mail/Scale Mail", `{"armor":{"armor_category":"medium","armor_class":{"base":14,"bonus":0,"dexterity":"capped","dexterity_maximum":2},"strength_requirement":0,"stealth_disadvantage":true,"don_doff":{"action":"none","don_minutes":5,"doff_minutes":1},"weight_pounds":45}}`,
			"wear", "item.wear armor.medium armor.training armor-class.14-plus-dexterity-max-2 stealth.disadvantage don.minutes-5 doff.minutes-1"},
		{"Shield/Shield", `{"armor":{"armor_category":"shield","armor_class":{"base":0,"bonus":2,"dexterity":"none","dexterity_maximum":0},"strength_requirement":0,"stealth_disadvantage":false,"don_doff":{"action":"utilize","don_minutes":0,"doff_minutes":0},"weight_pounds":6}}`,
			"wield", "item.wield armor.shield armor.training armor-class.bonus-2 don.utilize-action doff.utilize-action"},
		{"Splint Armor/Splint Armor", `{"armor":{"armor_category":"heavy","armor_class":{"base":17,"bonus":0,"dexterity":"none","dexterity_maximum":0},"strength_requirement":15,"stealth_disadvantage":true,"don_doff":{"action":"none","don_minutes":10,"doff_minutes":5},"weight_pounds":60}}`,
			"wear", "item.wear armor.heavy armor.training armor-class.17 strength.15 stealth.disadvantage don.minutes-10 doff.minutes-5"},
	}
	index := 0
	for _, f := range loadItems251(t).Families {
		for _, v := range f.Variants {
			if index >= len(expected) {
				t.Fatal("batch holds more records than the fixture")
			}
			want := expected[index]
			index++
			t.Run(want.name, func(t *testing.T) {
				var details ItemDetails251
				d := json.NewDecoder(strings.NewReader(want.details))
				d.DisallowUnknownFields()
				if err := d.Decode(&details); err != nil {
					t.Fatal(err)
				}
				if f.Name+"/"+v.Name != want.name || !reflect.DeepEqual(v.Details, &details) {
					t.Fatalf("source value mismatch for %s/%s: %+v", f.Name, v.Name, v.Details.Armor)
				}
				if v.Equip.Mode != want.equip || v.Use.Mode != "none" || len(v.Equip.Restrictions) != 0 || len(v.Use.Restrictions) != 0 || strings.Join(v.RequiredMechanicKeys, " ") != want.keys {
					t.Fatalf("source requirements mismatch: %v", v.RequiredMechanicKeys)
				}
			})
		}
	}
	if index != len(expected) {
		t.Fatalf("fixture lists %d records, batch holds %d", len(expected), index)
	}
}

// No item record existed for Chain Mail or Shield before this batch. The runtime
// holds fixed behaviour for them: the Armor Class terms in the character
// projection, and the Strength requirement, don and doff times and source IDs
// in ChangeEquipment. This batch's rows are the only item records, and this test
// holds the runtime values to them. No runtime code reads the batch.
func TestItems251JoinRuntimeChainMailAndShield(t *testing.T) {
	rows, ids := map[string]*ItemArmor{}, map[string]string{}
	for _, f := range loadItems251(t).Families {
		rows[f.Name], ids[f.Name] = f.Variants[0].Details.Armor, f.Presentations[0].ID
	}
	mail, shield := rows["Chain Mail"], rows["Shield"]
	if mail == nil || shield == nil {
		t.Fatal("runtime armor is missing from the batch")
	}
	build := func(armor string, hasShield bool, strength int) *Build {
		return &Build{Class: "fighter", Armor: armor, Shield: hasShield, Abilities: map[string]int{"strength": strength}}
	}

	for _, tc := range []struct {
		item, name, kind string
		equipped         bool
		seconds          int
	}{
		{"chain-mail", "Chain Mail", EquipmentEquip, false, mail.DonDoff.DonMinutes * 60},
		{"chain-mail", "Chain Mail", EquipmentRemove, true, mail.DonDoff.DoffMinutes * 60},
		// The runtime accounts a Utilize action as one six-second interval.
		{"shield", "Shield", EquipmentEquip, false, 6},
		{"shield", "Shield", EquipmentRemove, true, 6},
	} {
		b := build("none", false, mail.StrengthRequirement)
		if tc.equipped {
			b = build("chain-mail", true, mail.StrengthRequirement)
		}
		r, err := ChangeEquipment(b, tc.kind, tc.item, tc.equipped, 0)
		if err != nil {
			t.Fatal(err)
		}
		if r.TimeAfter != tc.seconds || !slices.Equal(r.SourceIDs, []string{ids[tc.name]}) {
			t.Fatalf("%s %s disagrees with the table: %+v", tc.kind, tc.item, r)
		}
	}
	if shield.DonDoff.Action != "utilize" {
		t.Fatal("the Shield row must take the Utilize action")
	}

	if _, err := ChangeEquipment(build("none", false, mail.StrengthRequirement-1), EquipmentEquip, "chain-mail", false, 0); err == nil {
		t.Fatal("runtime accepted Strength below the printed requirement")
	}

	data, err := os.ReadFile("testdata/fighter1.json")
	if err != nil {
		t.Fatal(err)
	}
	sheet, err := Project(data, 1, 12, 12, 0)
	if err != nil {
		t.Fatal(err)
	}
	joined := map[string]int{}
	for _, term := range sheet.Defenses["armor-class"].Derivation {
		switch term.SourceID {
		case ids["Chain Mail"]:
			joined[term.SourceID] = term.Value
		case ids["Shield"]:
			joined[term.SourceID] = term.Value
		}
	}
	want := map[string]int{ids["Chain Mail"]: mail.ArmorClass.Base, ids["Shield"]: shield.ArmorClass.Bonus}
	if mail.ArmorClass.Dexterity != "none" || !reflect.DeepEqual(joined, want) {
		t.Fatalf("projected Armor Class terms disagree with the table: %v", joined)
	}

	// The creation fixture names both items.
	fixture, err := NewSupportedMartialBuild("human", "fighter", "soldier")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"Chain Mail", "Shield"} {
		if !slices.Contains(fixture.Equipment, name) {
			t.Fatalf("creation fixture no longer names %s", name)
		}
	}
}

func TestItems251CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems251(t)
	composed, err := ComposeItemSourceIdentities(append(itemIdentitiesThrough(t, 250), batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, f := range composed {
		total += len(f.VariantIDs)
	}
	if len(composed) != 308 || total != 384 || composed[0].Name != "Potion of Healing" || composed[296].Name != "Breastplate" || composed[307].Name != "Splint Armor" {
		t.Fatalf("expected Potion of Healing, batches #225 to #250 and #251: %d families, %d variants", len(composed), total)
	}
	checkItemBatchComposition(t, batch.SourceIdentities(), "armor", "items_251")
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems251RejectMissingFields(t *testing.T) {
	rejectMissingItemMembers(t, items251JSON, 606, func(data []byte) error {
		_, err := DecodeItems251(data)
		return err
	})
}

func TestItems251RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails251]
	// v and row address a family's first variant by family index.
	v := func(b *batch, family int) *ItemSourceVariant[ItemDetails251] { return &b.Families[family].Variants[0] }
	row := func(b *batch, family int) *ItemArmor { return v(b, family).Details.Armor }
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-other" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 250 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/other"
			b.Families[12].Name = "Other"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/unreviewed" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Unreviewed Name" }},
		{"magic kind", func(b *batch) { b.Families[0].Kind = "magic" }},
		{"category", func(b *batch) { b.Families[0].Category = "wondrous-item" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"any-simple-weapon"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{90, 90} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"magic presentation", func(b *batch) { b.Families[0].Presentations[0].ID = "srd-5.2.1:content/magic-items/breastplate" }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"variant id", func(b *batch) { v(b, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0).Name = "Unreviewed" }},
		{"invented rarity", func(b *batch) { v(b, 0).Rarity = "common" }},
		{"attunement", func(b *batch) { v(b, 0).Attunement.Mode = "required" }},
		{"invented charges", func(b *batch) { v(b, 0).Charges.Maximum = 1 }},
		{"equip mode", func(b *batch) { v(b, 0).Equip.Mode = "carry" }},
		{"use mode", func(b *batch) { v(b, 0).Use.Mode = "drink" }},
		{"invented use restriction", func(b *batch) { v(b, 0).Use.Restrictions = []string{"proficiency"} }},
		{"unlisted price", func(b *batch) { v(b, 0).Cost.Kind = "unspecified" }},
		{"changed price", func(b *batch) { v(b, 0).Cost.Amount++ }},
		{"changed currency", func(b *batch) { v(b, 0).Cost.Currency = "ep" }},
		{"cost page", func(b *batch) { v(b, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys = v(b, 0).RequiredMechanicKeys[:1] }},
		{"unknown mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0).Details = &ItemDetails251{} }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"another family's row", func(b *batch) { v(b, 0).Details = v(b, 1).Details }},
		{"armor category", func(b *batch) { row(b, 0).ArmorCategory = "exotic" }},
		{"base armor class", func(b *batch) { row(b, 0).ArmorClass.Base++ }},
		{"base and bonus", func(b *batch) { row(b, 0).ArmorClass.Bonus = 2 }},
		{"no armor class", func(b *batch) { row(b, 0).ArmorClass.Base = 0 }},
		{"unknown dexterity rule", func(b *batch) { row(b, 0).ArmorClass.Dexterity = "half" }},
		{"changed first dexterity maximum", func(b *batch) { row(b, 0).ArmorClass.DexterityMaximum++ }},
		{"invented strength", func(b *batch) { row(b, 0).StrengthRequirement = 11 }},
		{"changed stealth", func(b *batch) { row(b, 0).StealthDisadvantage = !row(b, 0).StealthDisadvantage }},
		{"changed don time", func(b *batch) { row(b, 0).DonDoff.DonMinutes++ }},
		{"changed doff time", func(b *batch) { row(b, 0).DonDoff.DoffMinutes++ }},
		{"no doff time", func(b *batch) { row(b, 0).DonDoff.DoffMinutes = 0 }},
		{"utilize armor", func(b *batch) { row(b, 0).DonDoff.Action = "utilize" }},
		{"zero weight", func(b *batch) { row(b, 0).WeightPounds = 0 }},
		{"dropped Chain Mail strength", func(b *batch) { row(b, 1).StrengthRequirement = 0 }},
		{"changed Chain Mail strength", func(b *batch) { row(b, 1).StrengthRequirement = 15 }},
		{"changed dexterity maximum", func(b *batch) { row(b, 2).ArmorClass.DexterityMaximum = 3 }},
		{"changed Shield bonus", func(b *batch) { row(b, 10).ArmorClass.Bonus = 3 }},
		{"Shield base armor class", func(b *batch) { row(b, 10).ArmorClass.Base = 10 }},
		{"Shield minutes", func(b *batch) { row(b, 10).DonDoff.DonMinutes = 1 }},
		{"worn Shield", func(b *batch) { v(b, 10).Equip.Mode = "wear" }},
		{"wielded armor", func(b *batch) { v(b, 1).Equip.Mode = "wield" }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems251(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems251(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems251RejectUnsupportedJSON(t *testing.T) {
	decode := func(data []byte) error {
		_, err := DecodeItems251(data)
		return err
	}
	rejectUnsupportedItemJSON(t, items251JSON, `"id": "srd-5.2.1:content/items/breastplate/variant/standard"`, decode, []itemJSONChange{
		{"unknown payload", `"armor": {`, `"armor_plus": {`},
		{"unknown detail", `"weight_pounds": `, `"weight_ounces": `},
		{"string for number", `"strength_requirement": 0`, `"strength_requirement": "0"`},
		{"fractional number", `"strength_requirement": 0`, `"strength_requirement": 0.5`},
		{"number for boolean", `"stealth_disadvantage": false`, `"stealth_disadvantage": 0`},
	})
	// Neither the Healing contract decoder nor the previous batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items251JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems250(items251JSON); err == nil {
		t.Fatal("#250 decoder accepted the #251 batch")
	}
	if _, err := DecodeItems251(items250JSON); err == nil {
		t.Fatal("#251 decoder accepted the #250 batch")
	}
}
