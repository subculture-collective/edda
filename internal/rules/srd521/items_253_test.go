package srd521

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func loadItems253(t *testing.T) ItemSourceBatch[ItemDetails253] {
	t.Helper()
	b, err := LoadItems253()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities, variants and listed prices from pinned PDF pp. 93-94 for issue #253.
func TestItems253SourceIdentities(t *testing.T) {
	expected := []itemFamilyWant{
		{"Alchemist’s Supplies", "alchemists-supplies", 93, []itemVariantWant{{"standard", "Alchemist’s Supplies", 50, "gp"}}},
		{"Brewer’s Supplies", "brewers-supplies", 93, []itemVariantWant{{"standard", "Brewer’s Supplies", 20, "gp"}}},
		{"Calligrapher’s Supplies", "calligraphers-supplies", 93, []itemVariantWant{{"standard", "Calligrapher’s Supplies", 10, "gp"}}},
		{"Carpenter’s Tools", "carpenters-tools", 93, []itemVariantWant{{"standard", "Carpenter’s Tools", 8, "gp"}}},
		{"Cartographer’s Tools", "cartographers-tools", 93, []itemVariantWant{{"standard", "Cartographer’s Tools", 15, "gp"}}},
		{"Cobbler’s Tools", "cobblers-tools", 93, []itemVariantWant{{"standard", "Cobbler’s Tools", 5, "gp"}}},
		{"Cook’s Utensils", "cooks-utensils", 93, []itemVariantWant{{"standard", "Cook’s Utensils", 1, "gp"}}},
		{"Disguise Kit", "disguise-kit", 94, []itemVariantWant{{"standard", "Disguise Kit", 25, "gp"}}},
		{"Forgery Kit", "forgery-kit", 94, []itemVariantWant{{"standard", "Forgery Kit", 15, "gp"}}},
		{"Gaming Set", "gaming-set", 94, []itemVariantWant{{"dice", "Dice", 1, "sp"}, {"dragonchess", "Dragonchess", 1, "gp"}, {"playing-cards", "Playing cards", 5, "sp"}, {"three-dragon-ante", "Three-dragon ante", 1, "gp"}}},
		{"Glassblower’s Tools", "glassblowers-tools", 93, []itemVariantWant{{"standard", "Glassblower’s Tools", 30, "gp"}}},
		{"Herbalism Kit", "herbalism-kit", 94, []itemVariantWant{{"standard", "Herbalism Kit", 5, "gp"}}},
	}
	b := loadItems253(t)
	if len(b.Families) != 12 || b.Batch != "data-tools-01" || b.CatalogueIssue != 253 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the tool entries of issue #253")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 93–94") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	checkMundaneIdentities(t, b.Families, "tool", "tools", expected)
	// Join every family to the decision inventory the issue cites, in its order.
	joinItemInventory(t, b.SourceIdentities(), b.Batch, 253, "tools", []int{93, 93, 93, 93, 93, 93, 93, 94, 94, 94, 93, 94})
	b.Families[0].Variants[0].Details.Tool.Ability = "changed"
	if fresh := loadItems253(t); fresh.Families[0].Variants[0].Details.Tool.Ability != "intelligence" {
		t.Fatal("load must return independent data")
	}
}

// Literal values, equip modes and mechanic keys for all 15 records,
// transcribed from pinned PDF pp. 93-94.
func TestItems253SourceValues(t *testing.T) {
	expected := []struct{ name, details, equip, keys string }{
		{"Alchemist’s Supplies/Alchemist’s Supplies", `{"tool":{"tool_group":"artisans-tools","ability":"intelligence","utilize":[{"task":"identify-substance","dc":15},{"task":"start-fire","dc":15}],"craft":[{"item":"acid","except":[]},{"item":"alchemists-fire","except":[]},{"item":"component-pouch","except":[]},{"item":"oil","except":[]},{"item":"paper","except":[]},{"item":"perfume","except":[]}],"weight":{"listed":true,"pounds":8},"source_pages":[93,93]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.intelligence utilize.identify-substance-dc-15 utilize.start-fire-dc-15 craft.acid craft.alchemists-fire craft.component-pouch craft.oil craft.paper craft.perfume"},
		{"Brewer’s Supplies/Brewer’s Supplies", `{"tool":{"tool_group":"artisans-tools","ability":"intelligence","utilize":[{"task":"detect-poisoned-drink","dc":15},{"task":"identify-alcohol","dc":10}],"craft":[{"item":"antitoxin","except":[]}],"weight":{"listed":true,"pounds":9},"source_pages":[93,93]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.intelligence utilize.detect-poisoned-drink-dc-15 utilize.identify-alcohol-dc-10 craft.antitoxin"},
		{"Calligrapher’s Supplies/Calligrapher’s Supplies", `{"tool":{"tool_group":"artisans-tools","ability":"dexterity","utilize":[{"task":"write-text-with-impressive-flourishes-that-guard-against-forgery","dc":15}],"craft":[{"item":"ink","except":[]},{"item":"spell-scroll","except":[]}],"weight":{"listed":true,"pounds":5},"source_pages":[93,93]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.dexterity utilize.write-text-with-impressive-flourishes-that-guard-against-forgery-dc-15 craft.ink craft.spell-scroll"},
		{"Carpenter’s Tools/Carpenter’s Tools", `{"tool":{"tool_group":"artisans-tools","ability":"strength","utilize":[{"task":"seal-or-pry-open-door-or-container","dc":20}],"craft":[{"item":"club","except":[]},{"item":"greatclub","except":[]},{"item":"quarterstaff","except":[]},{"item":"barrel","except":[]},{"item":"chest","except":[]},{"item":"ladder","except":[]},{"item":"pole","except":[]},{"item":"portable-ram","except":[]},{"item":"torch","except":[]}],"weight":{"listed":true,"pounds":6},"source_pages":[93,93]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.strength utilize.seal-or-pry-open-door-or-container-dc-20 craft.club craft.greatclub craft.quarterstaff craft.barrel craft.chest craft.ladder craft.pole craft.portable-ram craft.torch"},
		{"Cartographer’s Tools/Cartographer’s Tools", `{"tool":{"tool_group":"artisans-tools","ability":"wisdom","utilize":[{"task":"draft-map-of-small-area","dc":15}],"craft":[{"item":"map","except":[]}],"weight":{"listed":true,"pounds":6},"source_pages":[93,93]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.wisdom utilize.draft-map-of-small-area-dc-15 craft.map"},
		{"Cobbler’s Tools/Cobbler’s Tools", `{"tool":{"tool_group":"artisans-tools","ability":"dexterity","utilize":[{"task":"modify-footwear-to-give-advantage-on-wearers-next-dexterity-acrobatics-check","dc":10}],"craft":[{"item":"climbers-kit","except":[]}],"weight":{"listed":true,"pounds":5},"source_pages":[93,93]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.dexterity utilize.modify-footwear-to-give-advantage-on-wearers-next-dexterity-acrobatics-check-dc-10 craft.climbers-kit"},
		{"Cook’s Utensils/Cook’s Utensils", `{"tool":{"tool_group":"artisans-tools","ability":"wisdom","utilize":[{"task":"improve-foods-flavor","dc":10},{"task":"detect-spoiled-or-poisoned-food","dc":15}],"craft":[{"item":"rations","except":[]}],"weight":{"listed":true,"pounds":8},"source_pages":[93,93]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.wisdom utilize.improve-foods-flavor-dc-10 utilize.detect-spoiled-or-poisoned-food-dc-15 craft.rations"},
		{"Disguise Kit/Disguise Kit", `{"tool":{"tool_group":"other-tools","ability":"charisma","utilize":[{"task":"apply-makeup","dc":10}],"craft":[{"item":"costume","except":[]}],"weight":{"listed":true,"pounds":3},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.apply-makeup-dc-10 craft.costume"},
		{"Forgery Kit/Forgery Kit", `{"tool":{"tool_group":"other-tools","ability":"dexterity","utilize":[{"task":"mimic-10-or-fewer-words-of-someone-elses-handwriting","dc":15},{"task":"duplicate-wax-seal","dc":20}],"craft":[],"weight":{"listed":true,"pounds":5},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.dexterity utilize.mimic-10-or-fewer-words-of-someone-elses-handwriting-dc-15 utilize.duplicate-wax-seal-dc-20"},
		{"Gaming Set/Dice", `{"tool":{"tool_group":"other-tools","ability":"wisdom","utilize":[{"task":"discern-whether-someone-is-cheating","dc":10},{"task":"win-game","dc":20}],"craft":[],"weight":{"listed":false,"pounds":0},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.wisdom utilize.discern-whether-someone-is-cheating-dc-10 utilize.win-game-dc-20"},
		{"Gaming Set/Dragonchess", `{"tool":{"tool_group":"other-tools","ability":"wisdom","utilize":[{"task":"discern-whether-someone-is-cheating","dc":10},{"task":"win-game","dc":20}],"craft":[],"weight":{"listed":false,"pounds":0},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.wisdom utilize.discern-whether-someone-is-cheating-dc-10 utilize.win-game-dc-20"},
		{"Gaming Set/Playing cards", `{"tool":{"tool_group":"other-tools","ability":"wisdom","utilize":[{"task":"discern-whether-someone-is-cheating","dc":10},{"task":"win-game","dc":20}],"craft":[],"weight":{"listed":false,"pounds":0},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.wisdom utilize.discern-whether-someone-is-cheating-dc-10 utilize.win-game-dc-20"},
		{"Gaming Set/Three-dragon ante", `{"tool":{"tool_group":"other-tools","ability":"wisdom","utilize":[{"task":"discern-whether-someone-is-cheating","dc":10},{"task":"win-game","dc":20}],"craft":[],"weight":{"listed":false,"pounds":0},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.wisdom utilize.discern-whether-someone-is-cheating-dc-10 utilize.win-game-dc-20"},
		{"Glassblower’s Tools/Glassblower’s Tools", `{"tool":{"tool_group":"artisans-tools","ability":"intelligence","utilize":[{"task":"discern-what-glass-object-held-in-past-24-hours","dc":15}],"craft":[{"item":"glass-bottle","except":[]},{"item":"magnifying-glass","except":[]},{"item":"spyglass","except":[]},{"item":"vial","except":[]}],"weight":{"listed":true,"pounds":5},"source_pages":[93,93]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.intelligence utilize.discern-what-glass-object-held-in-past-24-hours-dc-15 craft.glass-bottle craft.magnifying-glass craft.spyglass craft.vial"},
		{"Herbalism Kit/Herbalism Kit", `{"tool":{"tool_group":"other-tools","ability":"intelligence","utilize":[{"task":"identify-plant","dc":10}],"craft":[{"item":"antitoxin","except":[]},{"item":"candle","except":[]},{"item":"healers-kit","except":[]},{"item":"potion-of-healing","except":[]}],"weight":{"listed":true,"pounds":3},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.intelligence utilize.identify-plant-dc-10 craft.antitoxin craft.candle craft.healers-kit craft.potion-of-healing"},
	}
	index := 0
	for _, f := range loadItems253(t).Families {
		for _, v := range f.Variants {
			if index >= len(expected) {
				t.Fatal("batch holds more records than the fixture")
			}
			want := expected[index]
			index++
			t.Run(want.name, func(t *testing.T) {
				var details ItemDetails253
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

// No tool had an item record before this batch. The background source catalogue
// refers to tools by key, and the Fighter creation fixture names a "Dice Set".
// This batch's records are the only item records; the references stay as they are.
func TestItems253JoinToolReferences(t *testing.T) {
	families := map[string]ItemSourceFamily[ItemDetails253]{}
	for _, f := range loadItems253(t).Families {
		families[f.Name] = f
	}
	data, err := os.ReadFile("testdata/character-backgrounds.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"Calligrapher’s Supplies", "Gaming Set"} {
		f := families[name]
		if len(f.Presentations) != 1 {
			t.Fatalf("%s is missing from the batch", name)
		}
		reference := strings.Replace(f.Presentations[0].ID, "content/tools/", "mechanic/tools/", 1)
		if !strings.Contains(string(data), `"`+reference+`"`) {
			t.Fatalf("background catalogue does not refer to %s", reference)
		}
	}
	fixture, err := NewSupportedMartialBuild("human", "fighter", "soldier")
	if err != nil {
		t.Fatal(err)
	}
	set := families["Gaming Set"]
	if fixture.Tool != "Dice Set" || len(set.Variants) != 4 || set.Variants[0].Name != "Dice" {
		t.Fatalf("creation tool %q or the Gaming Set dice variant changed", fixture.Tool)
	}
}

func TestItems253CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems253(t)
	composed, err := ComposeItemSourceIdentities(append(itemIdentitiesThrough(t, 252), batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, f := range composed {
		total += len(f.VariantIDs)
	}
	if len(composed) != 321 || total != 400 || composed[0].Name != "Potion of Healing" || composed[309].Name != "Alchemist’s Supplies" || composed[320].Name != "Herbalism Kit" {
		t.Fatalf("expected Potion of Healing, batches #225 to #252 and #253: %d families, %d variants", len(composed), total)
	}
	checkItemBatchComposition(t, batch.SourceIdentities(), "tools", "items_253")
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems253RejectMissingFields(t *testing.T) {
	rejectMissingItemMembers(t, items253JSON, 751, func(data []byte) error {
		_, err := DecodeItems253(data)
		return err
	})
}

func TestItems253RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails253]
	// v and row address a family's first variant by family index.
	v := func(b *batch, family int) *ItemSourceVariant[ItemDetails253] { return &b.Families[family].Variants[0] }
	row := func(b *batch, family int) *ItemTool { return v(b, family).Details.Tool }
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-other" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 252 }},
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
		{"magic presentation", func(b *batch) {
			b.Families[0].Presentations[0].ID = "srd-5.2.1:content/magic-items/alchemists-supplies"
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
		{"empty details", func(b *batch) { v(b, 0).Details = &ItemDetails253{} }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"another family's row", func(b *batch) { v(b, 0).Details = v(b, 1).Details }},
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
		{"missing set variant", func(b *batch) { b.Families[9].Variants = b.Families[9].Variants[:3] }},
		{"reordered set variants", func(b *batch) {
			p := b.Families[9].Variants
			p[0], p[1] = p[1], p[0]
		}},
		{"changed set price", func(b *batch) { b.Families[9].Variants[2].Cost.Amount = 1 }},
		{"collapsed set", func(b *batch) { b.Families[9].Variants = b.Families[9].Variants[:1] }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems253(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems253(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems253RejectUnsupportedJSON(t *testing.T) {
	decode := func(data []byte) error {
		_, err := DecodeItems253(data)
		return err
	}
	rejectUnsupportedItemJSON(t, items253JSON, `"id": "srd-5.2.1:content/items/alchemists-supplies/variant/standard"`, decode, []itemJSONChange{
		{"unknown payload", `"tool": {`, `"tool_plus": {`},
		{"unknown detail", `"tool_group": `, `"tool_family": `},
		{"string for number", `"dc": `, `"dc": "1", "dc2": `},
		{"number for boolean", `"listed": true`, `"listed": 1`},
		{"null exception list", `"except": []`, `"except": null`},
	})
	// Neither the Healing contract decoder nor the previous batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items253JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems252(items253JSON); err == nil {
		t.Fatal("#252 decoder accepted the #253 batch")
	}
	if _, err := DecodeItems253(items252JSON); err == nil {
		t.Fatal("#253 decoder accepted the #252 batch")
	}
}
