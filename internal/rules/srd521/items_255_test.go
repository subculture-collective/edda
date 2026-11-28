package srd521

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func loadItems255(t *testing.T) ItemSourceBatch[ItemDetails255] {
	t.Helper()
	b, err := LoadItems255()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities, variants and listed prices from pinned PDF p. 94 for issue #255.
func TestItems255SourceIdentities(t *testing.T) {
	expected := []itemFamilyWant{
		{"Woodcarver’s Tools", "woodcarvers-tools", 94, []itemVariantWant{{"standard", "Woodcarver’s Tools", 1, "gp"}}},
	}
	b := loadItems255(t)
	if len(b.Families) != 1 || b.Batch != "data-tools-03" || b.CatalogueIssue != 255 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the tool entry of issue #255")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "p. 94") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	checkMundaneIdentities(t, b.Families, "tool", "tools", expected)
	// Join every family to the decision inventory the issue cites, in its order.
	joinItemInventory(t, b.SourceIdentities(), b.Batch, 255, "tools", []int{94})
	b.Families[0].Variants[0].Details.Tool.Ability = "changed"
	if fresh := loadItems255(t); fresh.Families[0].Variants[0].Details.Tool.Ability != "dexterity" {
		t.Fatal("load must return independent data")
	}
}

// Literal values, equip modes and mechanic keys for all 1 records,
// transcribed from pinned PDF p. 94.
func TestItems255SourceValues(t *testing.T) {
	expected := []struct{ name, details, equip, keys string }{
		{"Woodcarver’s Tools/Woodcarver’s Tools", `{"tool":{"tool_group":"artisans-tools","ability":"dexterity","utilize":[{"task":"carve-pattern-in-wood","dc":10}],"craft":[{"item":"club","except":[]},{"item":"greatclub","except":[]},{"item":"quarterstaff","except":[]},{"item":"ranged-weapons","except":["pistol","musket","sling"]},{"item":"arcane-focus","except":[]},{"item":"arrows","except":[]},{"item":"bolts","except":[]},{"item":"druidic-focus","except":[]},{"item":"ink-pen","except":[]},{"item":"needles","except":[]}],"weight":{"listed":true,"pounds":5},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.dexterity utilize.carve-pattern-in-wood-dc-10 craft.club craft.greatclub craft.quarterstaff craft.ranged-weapons-except-pistol-musket-sling craft.arcane-focus craft.arrows craft.bolts craft.druidic-focus craft.ink-pen craft.needles"},
	}
	index := 0
	for _, f := range loadItems255(t).Families {
		for _, v := range f.Variants {
			if index >= len(expected) {
				t.Fatal("batch holds more records than the fixture")
			}
			want := expected[index]
			index++
			t.Run(want.name, func(t *testing.T) {
				var details ItemDetails255
				d := json.NewDecoder(strings.NewReader(want.details))
				d.DisallowUnknownFields()
				if err := d.Decode(&details); err != nil {
					t.Fatal(err)
				}
				if f.Name+"/"+v.Name != want.name || !reflect.DeepEqual(v.Details, &details) {
					t.Fatalf("source value mismatch for %s/%s: %+v", f.Name, v.Name, v.Details.Tool)
				}
				if v.Equip.Mode != want.equip || v.Use.Mode != "utilize" || len(v.Equip.Restrictions) != 0 || len(v.Use.Restrictions) != 0 || strings.Join(v.RequiredMechanicKeys, " ") != want.keys {
					t.Fatalf("source requirements mismatch: %v", v.RequiredMechanicKeys)
				}
			})
		}
	}
	if index != len(expected) {
		t.Fatalf("fixture lists %d records, batch holds %d", len(expected), index)
	}
}

func TestItems255CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems255(t)
	composed, err := ComposeItemSourceIdentities(append(itemIdentitiesThrough(t, 254), batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, f := range composed {
		total += len(f.VariantIDs)
	}
	if len(composed) != 334 || total != 422 || composed[0].Name != "Potion of Healing" || composed[333].Name != "Woodcarver’s Tools" {
		t.Fatalf("expected Potion of Healing, batches #225 to #254 and #255: %d families, %d variants", len(composed), total)
	}
	checkItemBatchComposition(t, batch.SourceIdentities(), "tools", "items_255")
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems255RejectMissingFields(t *testing.T) {
	rejectMissingItemMembers(t, items255JSON, 84, func(data []byte) error {
		_, err := DecodeItems255(data)
		return err
	})
}

func TestItems255RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails255]
	// v and row address a family's first variant by family index.
	v := func(b *batch, family int) *ItemSourceVariant[ItemDetails255] { return &b.Families[family].Variants[0] }
	row := func(b *batch, family int) *ItemTool { return v(b, family).Details.Tool }
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-other" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 254 }},
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
		{"magic presentation", func(b *batch) { b.Families[0].Presentations[0].ID = "srd-5.2.1:content/magic-items/woodcarvers-tools" }},
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
		{"empty details", func(b *batch) { v(b, 0).Details = &ItemDetails255{} }},
		{"tool group", func(b *batch) { row(b, 0).ToolGroup = "magic-tools" }},
		{"unknown ability", func(b *batch) { row(b, 0).Ability = "luck" }},
		{"changed ability", func(b *batch) { row(b, 0).Ability = "charisma" }},
		{"changed DC", func(b *batch) { row(b, 0).Utilize[0].DC++ }},
		{"zero DC", func(b *batch) { row(b, 0).Utilize[0].DC = 0 }},
		{"changed task", func(b *batch) { row(b, 0).Utilize[0].Task = "other-task" }},
		{"prose task", func(b *batch) { row(b, 0).Utilize[0].Task = "Do the thing." }},
		{"dropped utilize", func(b *batch) { row(b, 0).Utilize = []ItemToolUtilize{} }},
		{"duplicate task", func(b *batch) {
			p := row(b, 0)
			p.Utilize = append(p.Utilize, p.Utilize[0])
		}},
		{"invented weight", func(b *batch) { row(b, 0).Weight.Pounds = 0 }},
		{"unlisted weight with value", func(b *batch) { row(b, 0).Weight.Listed = false }},
		{"first source page", func(b *batch) { row(b, 0).SourcePages = []int{92, 93} }},
		{"overlong source pages", func(b *batch) {
			p := row(b, 0)
			p.SourcePages = []int{p.SourcePages[0], p.SourcePages[0] + 2}
		}},
		{"short source pages", func(b *batch) { row(b, 0).SourcePages = row(b, 0).SourcePages[:1] }},
		{"dropped craft", func(b *batch) { row(b, 0).Craft = []ItemToolCraft{} }},
		{"changed craft", func(b *batch) { row(b, 0).Craft[0].Item = "other-item" }},
		{"duplicate craft", func(b *batch) {
			p := row(b, 0)
			p.Craft = append(p.Craft, p.Craft[0])
		}},
		{"invented exception", func(b *batch) { row(b, 0).Craft[0].Except = []string{"other-item"} }},
		{"dropped exception", func(b *batch) { row(b, 0).Craft[3].Except = []string{} }},
		{"duplicate exception", func(b *batch) { row(b, 0).Craft[3].Except[1] = "pistol" }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems255(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems255(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems255RejectUnsupportedJSON(t *testing.T) {
	decode := func(data []byte) error {
		_, err := DecodeItems255(data)
		return err
	}
	rejectUnsupportedItemJSON(t, items255JSON, `"id": "srd-5.2.1:content/items/woodcarvers-tools/variant/standard"`, decode, []itemJSONChange{
		{"unknown payload", `"tool": {`, `"tool_plus": {`},
		{"unknown detail", `"tool_group": `, `"tool_family": `},
		{"string for number", `"dc": `, `"dc": "1", "dc2": `},
		{"number for boolean", `"listed": true`, `"listed": 1`},
		{"null exception list", `"except": []`, `"except": null`},
	})
	// Neither the Healing contract decoder nor the previous batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items255JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems254(items255JSON); err == nil {
		t.Fatal("#254 decoder accepted the #255 batch")
	}
	if _, err := DecodeItems255(items254JSON); err == nil {
		t.Fatal("#255 decoder accepted the #254 batch")
	}
}
