package srd521

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func loadItems252(t *testing.T) ItemSourceBatch[ItemDetails252] {
	t.Helper()
	b, err := LoadItems252()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities, variants and listed prices from pinned PDF p. 92 for issue #252.
func TestItems252SourceIdentities(t *testing.T) {
	expected := []itemFamilyWant{
		{"Studded Leather Armor", "studded-leather-armor", 92, []itemVariantWant{{"standard", "Studded Leather Armor", 45, "gp"}}},
	}
	b := loadItems252(t)
	if len(b.Families) != 1 || b.Batch != "data-armor-02" || b.CatalogueIssue != 252 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the armor row of issue #252")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "p. 92") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	checkMundaneIdentities(t, b.Families, "armor", "armor", expected)
	// Join every family to the decision inventory the issue cites, in its order.
	joinItemInventory(t, b.SourceIdentities(), b.Batch, 252, "armor", []int{92})
	b.Families[0].Variants[0].Details.Armor.ArmorCategory = "changed"
	if fresh := loadItems252(t); fresh.Families[0].Variants[0].Details.Armor.ArmorCategory != "light" {
		t.Fatal("load must return independent data")
	}
}

// Literal values, equip modes and mechanic keys for all 1 records,
// transcribed from pinned PDF p. 92.
func TestItems252SourceValues(t *testing.T) {
	expected := []struct{ name, details, equip, keys string }{
		{"Studded Leather Armor/Studded Leather Armor", `{"armor":{"armor_category":"light","armor_class":{"base":12,"bonus":0,"dexterity":"full","dexterity_maximum":0},"strength_requirement":0,"stealth_disadvantage":false,"don_doff":{"action":"none","don_minutes":1,"doff_minutes":1},"weight_pounds":13}}`,
			"wear", "item.wear armor.light armor.training armor-class.12-plus-dexterity don.minutes-1 doff.minutes-1"},
	}
	index := 0
	for _, f := range loadItems252(t).Families {
		for _, v := range f.Variants {
			if index >= len(expected) {
				t.Fatal("batch holds more records than the fixture")
			}
			want := expected[index]
			index++
			t.Run(want.name, func(t *testing.T) {
				var details ItemDetails252
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

func TestItems252CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems252(t)
	composed, err := ComposeItemSourceIdentities(append(itemIdentitiesThrough(t, 251), batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, f := range composed {
		total += len(f.VariantIDs)
	}
	if len(composed) != 309 || total != 385 || composed[0].Name != "Potion of Healing" || composed[308].Name != "Studded Leather Armor" {
		t.Fatalf("expected Potion of Healing, batches #225 to #251 and #252: %d families, %d variants", len(composed), total)
	}
	checkItemBatchComposition(t, batch.SourceIdentities(), "armor", "items_252")
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems252RejectMissingFields(t *testing.T) {
	rejectMissingItemMembers(t, items252JSON, 67, func(data []byte) error {
		_, err := DecodeItems252(data)
		return err
	})
}

func TestItems252RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails252]
	// v and row address a family's first variant by family index.
	v := func(b *batch, family int) *ItemSourceVariant[ItemDetails252] { return &b.Families[family].Variants[0] }
	row := func(b *batch, family int) *ItemArmor { return v(b, family).Details.Armor }
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-other" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 251 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[1].ID = "srd-5.2.1:content/items/other"
			b.Families[1].Name = "Other"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/unreviewed" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Unreviewed Name" }},
		{"magic kind", func(b *batch) { b.Families[0].Kind = "magic" }},
		{"category", func(b *batch) { b.Families[0].Category = "wondrous-item" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"any-simple-weapon"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{90, 90} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"magic presentation", func(b *batch) {
			b.Families[0].Presentations[0].ID = "srd-5.2.1:content/magic-items/studded-leather-armor"
		}},
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
		{"empty details", func(b *batch) { v(b, 0).Details = &ItemDetails252{} }},
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
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems252(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems252(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems252RejectUnsupportedJSON(t *testing.T) {
	decode := func(data []byte) error {
		_, err := DecodeItems252(data)
		return err
	}
	rejectUnsupportedItemJSON(t, items252JSON, `"id": "srd-5.2.1:content/items/studded-leather-armor/variant/standard"`, decode, []itemJSONChange{
		{"unknown payload", `"armor": {`, `"armor_plus": {`},
		{"unknown detail", `"weight_pounds": `, `"weight_ounces": `},
		{"string for number", `"strength_requirement": 0`, `"strength_requirement": "0"`},
		{"fractional number", `"strength_requirement": 0`, `"strength_requirement": 0.5`},
		{"number for boolean", `"stealth_disadvantage": false`, `"stealth_disadvantage": 0`},
	})
	// Neither the Healing contract decoder nor the previous batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items252JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems251(items252JSON); err == nil {
		t.Fatal("#251 decoder accepted the #252 batch")
	}
	if _, err := DecodeItems252(items251JSON); err == nil {
		t.Fatal("#252 decoder accepted the #251 batch")
	}
}
