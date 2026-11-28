package srd521

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func loadItems254(t *testing.T) ItemSourceBatch[ItemDetails254] {
	t.Helper()
	b, err := LoadItems254()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities, variants and listed prices from pinned PDF pp. 93-94 for issue #254.
func TestItems254SourceIdentities(t *testing.T) {
	expected := []itemFamilyWant{
		{"Jeweler’s Tools", "jewelers-tools", 93, []itemVariantWant{{"standard", "Jeweler’s Tools", 25, "gp"}}},
		{"Leatherworker’s Tools", "leatherworkers-tools", 93, []itemVariantWant{{"standard", "Leatherworker’s Tools", 5, "gp"}}},
		{"Mason’s Tools", "masons-tools", 93, []itemVariantWant{{"standard", "Mason’s Tools", 10, "gp"}}},
		{"Musical Instrument", "musical-instrument", 94, []itemVariantWant{{"bagpipes", "Bagpipes", 30, "gp"}, {"drum", "Drum", 6, "gp"}, {"dulcimer", "Dulcimer", 25, "gp"}, {"flute", "Flute", 2, "gp"}, {"horn", "Horn", 3, "gp"}, {"lute", "Lute", 35, "gp"}, {"lyre", "Lyre", 30, "gp"}, {"pan-flute", "Pan flute", 12, "gp"}, {"shawm", "Shawm", 2, "gp"}, {"viol", "Viol", 30, "gp"}}},
		{"Navigator’s Tools", "navigators-tools", 94, []itemVariantWant{{"standard", "Navigator’s Tools", 25, "gp"}}},
		{"Painter’s Supplies", "painters-supplies", 93, []itemVariantWant{{"standard", "Painter’s Supplies", 10, "gp"}}},
		{"Poisoner’s Kit", "poisoners-kit", 94, []itemVariantWant{{"standard", "Poisoner’s Kit", 50, "gp"}}},
		{"Potter’s Tools", "potters-tools", 93, []itemVariantWant{{"standard", "Potter’s Tools", 10, "gp"}}},
		{"Smith’s Tools", "smiths-tools", 93, []itemVariantWant{{"standard", "Smith’s Tools", 20, "gp"}}},
		{"Thieves’ Tools", "thieves-tools", 94, []itemVariantWant{{"standard", "Thieves’ Tools", 25, "gp"}}},
		{"Tinker’s Tools", "tinkers-tools", 94, []itemVariantWant{{"standard", "Tinker’s Tools", 50, "gp"}}},
		{"Weaver’s Tools", "weavers-tools", 94, []itemVariantWant{{"standard", "Weaver’s Tools", 1, "gp"}}},
	}
	b := loadItems254(t)
	if len(b.Families) != 12 || b.Batch != "data-tools-02" || b.CatalogueIssue != 254 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the tool entries of issue #254")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 93–94") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	checkMundaneIdentities(t, b.Families, "tool", "tools", expected)
	// Join every family to the decision inventory the issue cites, in its order.
	joinItemInventory(t, b.SourceIdentities(), b.Batch, 254, "tools", []int{93, 93, 93, 94, 94, 93, 94, 93, 93, 94, 94, 94})
	b.Families[0].Variants[0].Details.Tool.Ability = "changed"
	if fresh := loadItems254(t); fresh.Families[0].Variants[0].Details.Tool.Ability != "intelligence" {
		t.Fatal("load must return independent data")
	}
}

// Literal values, equip modes and mechanic keys for all 21 records,
// transcribed from pinned PDF pp. 93-94.
func TestItems254SourceValues(t *testing.T) {
	expected := []struct{ name, details, equip, keys string }{
		{"Jeweler’s Tools/Jeweler’s Tools", `{"tool":{"tool_group":"artisans-tools","ability":"intelligence","utilize":[{"task":"discern-gems-value","dc":15}],"craft":[{"item":"arcane-focus","except":[]},{"item":"holy-symbol","except":[]}],"weight":{"listed":true,"pounds":2},"source_pages":[93,93]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.intelligence utilize.discern-gems-value-dc-15 craft.arcane-focus craft.holy-symbol"},
		{"Leatherworker’s Tools/Leatherworker’s Tools", `{"tool":{"tool_group":"artisans-tools","ability":"dexterity","utilize":[{"task":"add-design-to-leather-item","dc":10}],"craft":[{"item":"sling","except":[]},{"item":"whip","except":[]},{"item":"hide-armor","except":[]},{"item":"leather-armor","except":[]},{"item":"studded-leather-armor","except":[]},{"item":"backpack","except":[]},{"item":"crossbow-bolt-case","except":[]},{"item":"map-or-scroll-case","except":[]},{"item":"parchment","except":[]},{"item":"pouch","except":[]},{"item":"quiver","except":[]},{"item":"waterskin","except":[]}],"weight":{"listed":true,"pounds":5},"source_pages":[93,93]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.dexterity utilize.add-design-to-leather-item-dc-10 craft.sling craft.whip craft.hide-armor craft.leather-armor craft.studded-leather-armor craft.backpack craft.crossbow-bolt-case craft.map-or-scroll-case craft.parchment craft.pouch craft.quiver craft.waterskin"},
		{"Mason’s Tools/Mason’s Tools", `{"tool":{"tool_group":"artisans-tools","ability":"strength","utilize":[{"task":"chisel-symbol-or-hole-in-stone","dc":10}],"craft":[{"item":"block-and-tackle","except":[]}],"weight":{"listed":true,"pounds":8},"source_pages":[93,93]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.strength utilize.chisel-symbol-or-hole-in-stone-dc-10 craft.block-and-tackle"},
		{"Musical Instrument/Bagpipes", `{"tool":{"tool_group":"other-tools","ability":"charisma","utilize":[{"task":"play-known-tune","dc":10},{"task":"improvise-song","dc":15}],"craft":[],"weight":{"listed":true,"pounds":6},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{"Musical Instrument/Drum", `{"tool":{"tool_group":"other-tools","ability":"charisma","utilize":[{"task":"play-known-tune","dc":10},{"task":"improvise-song","dc":15}],"craft":[],"weight":{"listed":true,"pounds":3},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{"Musical Instrument/Dulcimer", `{"tool":{"tool_group":"other-tools","ability":"charisma","utilize":[{"task":"play-known-tune","dc":10},{"task":"improvise-song","dc":15}],"craft":[],"weight":{"listed":true,"pounds":10},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{"Musical Instrument/Flute", `{"tool":{"tool_group":"other-tools","ability":"charisma","utilize":[{"task":"play-known-tune","dc":10},{"task":"improvise-song","dc":15}],"craft":[],"weight":{"listed":true,"pounds":1},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{"Musical Instrument/Horn", `{"tool":{"tool_group":"other-tools","ability":"charisma","utilize":[{"task":"play-known-tune","dc":10},{"task":"improvise-song","dc":15}],"craft":[],"weight":{"listed":true,"pounds":2},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{"Musical Instrument/Lute", `{"tool":{"tool_group":"other-tools","ability":"charisma","utilize":[{"task":"play-known-tune","dc":10},{"task":"improvise-song","dc":15}],"craft":[],"weight":{"listed":true,"pounds":2},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{"Musical Instrument/Lyre", `{"tool":{"tool_group":"other-tools","ability":"charisma","utilize":[{"task":"play-known-tune","dc":10},{"task":"improvise-song","dc":15}],"craft":[],"weight":{"listed":true,"pounds":2},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{"Musical Instrument/Pan flute", `{"tool":{"tool_group":"other-tools","ability":"charisma","utilize":[{"task":"play-known-tune","dc":10},{"task":"improvise-song","dc":15}],"craft":[],"weight":{"listed":true,"pounds":2},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{"Musical Instrument/Shawm", `{"tool":{"tool_group":"other-tools","ability":"charisma","utilize":[{"task":"play-known-tune","dc":10},{"task":"improvise-song","dc":15}],"craft":[],"weight":{"listed":true,"pounds":1},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{"Musical Instrument/Viol", `{"tool":{"tool_group":"other-tools","ability":"charisma","utilize":[{"task":"play-known-tune","dc":10},{"task":"improvise-song","dc":15}],"craft":[],"weight":{"listed":true,"pounds":1},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.charisma utilize.play-known-tune-dc-10 utilize.improvise-song-dc-15"},
		{"Navigator’s Tools/Navigator’s Tools", `{"tool":{"tool_group":"other-tools","ability":"wisdom","utilize":[{"task":"plot-course","dc":10},{"task":"determine-position-by-stargazing","dc":15}],"craft":[],"weight":{"listed":true,"pounds":2},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.wisdom utilize.plot-course-dc-10 utilize.determine-position-by-stargazing-dc-15"},
		{"Painter’s Supplies/Painter’s Supplies", `{"tool":{"tool_group":"artisans-tools","ability":"wisdom","utilize":[{"task":"paint-recognizable-image-of-something-youve-seen","dc":10}],"craft":[{"item":"druidic-focus","except":[]},{"item":"holy-symbol","except":[]}],"weight":{"listed":true,"pounds":5},"source_pages":[93,93]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.wisdom utilize.paint-recognizable-image-of-something-youve-seen-dc-10 craft.druidic-focus craft.holy-symbol"},
		{"Poisoner’s Kit/Poisoner’s Kit", `{"tool":{"tool_group":"other-tools","ability":"intelligence","utilize":[{"task":"detect-poisoned-object","dc":10}],"craft":[{"item":"basic-poison","except":[]}],"weight":{"listed":true,"pounds":2},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.intelligence utilize.detect-poisoned-object-dc-10 craft.basic-poison"},
		{"Potter’s Tools/Potter’s Tools", `{"tool":{"tool_group":"artisans-tools","ability":"intelligence","utilize":[{"task":"discern-what-ceramic-object-held-in-past-24-hours","dc":15}],"craft":[{"item":"jug","except":[]},{"item":"lamp","except":[]}],"weight":{"listed":true,"pounds":3},"source_pages":[93,93]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.intelligence utilize.discern-what-ceramic-object-held-in-past-24-hours-dc-15 craft.jug craft.lamp"},
		{"Smith’s Tools/Smith’s Tools", `{"tool":{"tool_group":"artisans-tools","ability":"strength","utilize":[{"task":"pry-open-door-or-container","dc":20}],"craft":[{"item":"any-melee-weapon","except":["club","greatclub","quarterstaff","whip"]},{"item":"medium-armor","except":["hide"]},{"item":"heavy-armor","except":[]},{"item":"ball-bearings","except":[]},{"item":"bucket","except":[]},{"item":"caltrops","except":[]},{"item":"chain","except":[]},{"item":"crowbar","except":[]},{"item":"firearm-bullets","except":[]},{"item":"grappling-hook","except":[]},{"item":"iron-pot","except":[]},{"item":"iron-spikes","except":[]},{"item":"sling-bullets","except":[]}],"weight":{"listed":true,"pounds":8},"source_pages":[93,94]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.strength utilize.pry-open-door-or-container-dc-20 craft.any-melee-weapon-except-club-greatclub-quarterstaff-whip craft.medium-armor-except-hide craft.heavy-armor craft.ball-bearings craft.bucket craft.caltrops craft.chain craft.crowbar craft.firearm-bullets craft.grappling-hook craft.iron-pot craft.iron-spikes craft.sling-bullets"},
		{"Thieves’ Tools/Thieves’ Tools", `{"tool":{"tool_group":"other-tools","ability":"dexterity","utilize":[{"task":"pick-lock","dc":15},{"task":"disarm-trap","dc":15}],"craft":[],"weight":{"listed":true,"pounds":1},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.other-tools ability.dexterity utilize.pick-lock-dc-15 utilize.disarm-trap-dc-15"},
		{"Tinker’s Tools/Tinker’s Tools", `{"tool":{"tool_group":"artisans-tools","ability":"dexterity","utilize":[{"task":"assemble-tiny-item-composed-of-scrap-which-falls-apart-in-1-minute","dc":20}],"craft":[{"item":"musket","except":[]},{"item":"pistol","except":[]},{"item":"bell","except":[]},{"item":"bullseye-lantern","except":[]},{"item":"flask","except":[]},{"item":"hooded-lantern","except":[]},{"item":"hunting-trap","except":[]},{"item":"lock","except":[]},{"item":"manacles","except":[]},{"item":"mirror","except":[]},{"item":"shovel","except":[]},{"item":"signal-whistle","except":[]},{"item":"tinderbox","except":[]}],"weight":{"listed":true,"pounds":10},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.dexterity utilize.assemble-tiny-item-composed-of-scrap-which-falls-apart-in-1-minute-dc-20 craft.musket craft.pistol craft.bell craft.bullseye-lantern craft.flask craft.hooded-lantern craft.hunting-trap craft.lock craft.manacles craft.mirror craft.shovel craft.signal-whistle craft.tinderbox"},
		{"Weaver’s Tools/Weaver’s Tools", `{"tool":{"tool_group":"artisans-tools","ability":"dexterity","utilize":[{"task":"mend-tear-in-clothing","dc":10},{"task":"sew-tiny-design","dc":10}],"craft":[{"item":"padded-armor","except":[]},{"item":"basket","except":[]},{"item":"bedroll","except":[]},{"item":"blanket","except":[]},{"item":"fine-clothes","except":[]},{"item":"net","except":[]},{"item":"robe","except":[]},{"item":"rope","except":[]},{"item":"sack","except":[]},{"item":"string","except":[]},{"item":"tent","except":[]},{"item":"travelers-clothes","except":[]}],"weight":{"listed":true,"pounds":5},"source_pages":[94,94]}}`,
			"none", "tool.utilize tool.proficiency tool.artisans-tools ability.dexterity utilize.mend-tear-in-clothing-dc-10 utilize.sew-tiny-design-dc-10 craft.padded-armor craft.basket craft.bedroll craft.blanket craft.fine-clothes craft.net craft.robe craft.rope craft.sack craft.string craft.tent craft.travelers-clothes"},
	}
	index := 0
	for _, f := range loadItems254(t).Families {
		for _, v := range f.Variants {
			if index >= len(expected) {
				t.Fatal("batch holds more records than the fixture")
			}
			want := expected[index]
			index++
			t.Run(want.name, func(t *testing.T) {
				var details ItemDetails254
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

// The background and Rogue source catalogues refer to Thieves' Tools by key.
// This batch's record is the only item record.
func TestItems254JoinThievesToolsReferences(t *testing.T) {
	var presentation string
	for _, f := range loadItems254(t).Families {
		if f.Name == "Thieves’ Tools" {
			presentation = f.Presentations[0].ID
		}
	}
	if presentation != "srd-5.2.1:content/tools/thieves-tools" {
		t.Fatalf("unexpected Thieves' Tools presentation %q", presentation)
	}
	for file, reference := range map[string]string{
		"testdata/character-backgrounds.json": "srd-5.2.1:mechanic/tools/thieves-tools",
		"testdata/character-rogue.json":       "srd-5.2.1:mechanic/rogue/tools/thieves-tools",
	} {
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), `"`+reference+`"`) {
			t.Fatalf("%s does not refer to %s", file, reference)
		}
	}
}

func TestItems254CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems254(t)
	composed, err := ComposeItemSourceIdentities(append(itemIdentitiesThrough(t, 253), batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, f := range composed {
		total += len(f.VariantIDs)
	}
	if len(composed) != 333 || total != 421 || composed[0].Name != "Potion of Healing" || composed[321].Name != "Jeweler’s Tools" || composed[332].Name != "Weaver’s Tools" {
		t.Fatalf("expected Potion of Healing, batches #225 to #253 and #254: %d families, %d variants", len(composed), total)
	}
	checkItemBatchComposition(t, batch.SourceIdentities(), "tools", "items_254")
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems254RejectMissingFields(t *testing.T) {
	rejectMissingItemMembers(t, items254JSON, 1027, func(data []byte) error {
		_, err := DecodeItems254(data)
		return err
	})
}

func TestItems254RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails254]
	// v and row address a family's first variant by family index.
	v := func(b *batch, family int) *ItemSourceVariant[ItemDetails254] { return &b.Families[family].Variants[0] }
	row := func(b *batch, family int) *ItemTool { return v(b, family).Details.Tool }
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-other" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 253 }},
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
		{"magic presentation", func(b *batch) { b.Families[0].Presentations[0].ID = "srd-5.2.1:content/magic-items/jewelers-tools" }},
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
		{"empty details", func(b *batch) { v(b, 0).Details = &ItemDetails254{} }},
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
		{"missing instrument", func(b *batch) { b.Families[3].Variants = b.Families[3].Variants[:9] }},
		{"reordered instruments", func(b *batch) {
			p := b.Families[3].Variants
			p[0], p[1] = p[1], p[0]
		}},
		{"changed instrument price", func(b *batch) { b.Families[3].Variants[5].Cost.Amount++ }},
		{"dropped exception", func(b *batch) { row(b, 8).Craft[0].Except = []string{} }},
		{"changed exception", func(b *batch) { row(b, 8).Craft[0].Except[3] = "net" }},
		{"duplicate exception", func(b *batch) { row(b, 8).Craft[0].Except[3] = "club" }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems254(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems254(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems254RejectUnsupportedJSON(t *testing.T) {
	decode := func(data []byte) error {
		_, err := DecodeItems254(data)
		return err
	}
	rejectUnsupportedItemJSON(t, items254JSON, `"id": "srd-5.2.1:content/items/jewelers-tools/variant/standard"`, decode, []itemJSONChange{
		{"unknown payload", `"tool": {`, `"tool_plus": {`},
		{"unknown detail", `"tool_group": `, `"tool_family": `},
		{"string for number", `"dc": `, `"dc": "1", "dc2": `},
		{"number for boolean", `"listed": true`, `"listed": 1`},
		{"null exception list", `"except": []`, `"except": null`},
	})
	// Neither the Healing contract decoder nor the previous batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items254JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems253(items254JSON); err == nil {
		t.Fatal("#253 decoder accepted the #254 batch")
	}
	if _, err := DecodeItems254(items253JSON); err == nil {
		t.Fatal("#254 decoder accepted the #253 batch")
	}
}
