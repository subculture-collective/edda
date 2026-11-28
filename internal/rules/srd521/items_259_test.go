package srd521

import (
	"encoding/json"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func loadItems259(t *testing.T) ItemSourceBatch[ItemDetails259] {
	t.Helper()
	b, err := LoadItems259()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities, variants and listed prices from pinned PDF pp. 97–98 for issue #259.
func TestItems259SourceIdentities(t *testing.T) {
	expected := []itemFamilyWant{
		{"Healer’s Kit", "healers-kit", 97, []itemVariantWant{{"standard", "Healer’s Kit", 5, "gp"}}},
		{"Holy Symbol", "holy-symbol", 97, []itemVariantWant{{"amulet", "Amulet", 5, "gp"}, {"emblem", "Emblem", 5, "gp"}, {"reliquary", "Reliquary", 5, "gp"}}},
		{"Holy Water", "holy-water", 97, []itemVariantWant{{"standard", "Holy Water", 25, "gp"}}},
		{"Hunting Trap", "hunting-trap", 97, []itemVariantWant{{"standard", "Hunting Trap", 5, "gp"}}},
		{"Ink", "ink", 98, []itemVariantWant{{"standard", "Ink", 10, "gp"}}},
		{"Ink Pen", "ink-pen", 98, []itemVariantWant{{"standard", "Ink Pen", 2, "cp"}}},
		{"Jug", "jug", 98, []itemVariantWant{{"standard", "Jug", 2, "cp"}}},
		{"Ladder", "ladder", 98, []itemVariantWant{{"standard", "Ladder", 1, "sp"}}},
		{"Lamp", "lamp", 98, []itemVariantWant{{"standard", "Lamp", 5, "sp"}}},
		{"Lantern, Bullseye", "lantern-bullseye", 98, []itemVariantWant{{"standard", "Lantern, Bullseye", 10, "gp"}}},
		{"Lantern, Hooded", "lantern-hooded", 98, []itemVariantWant{{"standard", "Lantern, Hooded", 5, "gp"}}},
		{"Lock", "lock", 98, []itemVariantWant{{"standard", "Lock", 10, "gp"}}},
	}
	b := loadItems259(t)
	if len(b.Families) != 12 || b.Batch != "data-adventuring-gear-04" || b.CatalogueIssue != 259 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the new gear families of issue #259")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 97–98") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	checkMundaneIdentities(t, b.Families, "adventuring-gear", "adventuring-gear", expected)
	// Join every family to the decision inventory the issue cites, in its order.
	joinItemInventory(t, b.SourceIdentities(), b.Batch, 259, "adventuring-gear", []int{97, 97, 97, 97, 98, 98, 98, 98, 98, 98, 98, 98})
	b.Families[0].Variants[0].Details.Gear.SourcePages[0] = 1
	if fresh := loadItems259(t); fresh.Families[0].Variants[0].Details.Gear.SourcePages[0] != 97 {
		t.Fatal("load must return independent data")
	}
}

// Literal values, use modes and mechanic keys for all 14 records,
// transcribed from pinned PDF pp. 97–98.
func TestItems259SourceValues(t *testing.T) {
	expected := []struct{ name, details, use, keys string }{
		{"Healer’s Kit/Healer’s Kit", `{"gear":{"actions":["utilize"],"checks":[],"damage":[],"measures":[{"name":"uses","numerator":10,"denominator":1,"unit":"uses"},{"name":"stabilize-target-hit-points","numerator":0,"denominator":1,"unit":"hit-points"}],"traits":["expend-one-use-to-stabilize-unconscious-creature","no-wisdom-medicine-check-needed"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":3,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"utilize", "item.gear action.utilize measure.uses-10-uses measure.stabilize-target-hit-points-0-hit-points trait.expend-one-use-to-stabilize-unconscious-creature trait.no-wisdom-medicine-check-needed"},
		{"Holy Symbol/Amulet", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["bejeweled-or-painted-to-channel-divine-magic","table-states-how-each-symbol-is-carried","worn-or-held"],"contents":[],"spellcasting_focus":["cleric","paladin"],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear spellcasting-focus.cleric spellcasting-focus.paladin trait.bejeweled-or-painted-to-channel-divine-magic trait.table-states-how-each-symbol-is-carried trait.worn-or-held"},
		{"Holy Symbol/Emblem", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["bejeweled-or-painted-to-channel-divine-magic","table-states-how-each-symbol-is-carried","borne-on-fabric-or-shield"],"contents":[],"spellcasting_focus":["cleric","paladin"],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear spellcasting-focus.cleric spellcasting-focus.paladin trait.bejeweled-or-painted-to-channel-divine-magic trait.table-states-how-each-symbol-is-carried trait.borne-on-fabric-or-shield"},
		{"Holy Symbol/Reliquary", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["bejeweled-or-painted-to-channel-divine-magic","table-states-how-each-symbol-is-carried","held"],"contents":[],"spellcasting_focus":["cleric","paladin"],"weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"none", "item.gear spellcasting-focus.cleric spellcasting-focus.paladin trait.bejeweled-or-painted-to-channel-divine-magic trait.table-states-how-each-symbol-is-carried trait.held"},
		{"Holy Water/Holy Water", `{"gear":{"actions":["attack-replacement"],"checks":[{"purpose":"avoid-damage","kind":"saving-throw","ability":"dexterity","skill":"none","dc":8,"dc_rule":"plus-dexterity-modifier-and-proficiency-bonus"}],"damage":[{"purpose":"thrown-flask-against-fiend-or-undead","dice":2,"sides":8,"flat":0,"type":"radiant"}],"measures":[{"name":"throw-range","numerator":20,"denominator":1,"unit":"feet"}],"traits":["throw-flask-in-place-of-one-attack","target-one-creature-you-can-see","damages-only-fiend-or-undead"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[97,97]}}`,
			"attack-replacement", "item.gear action.attack-replacement saving-throw.avoid-damage-dexterity-dc-8-plus-dexterity-modifier-and-proficiency-bonus damage.thrown-flask-against-fiend-or-undead-2d8-radiant measure.throw-range-20-feet trait.throw-flask-in-place-of-one-attack trait.target-one-creature-you-can-see trait.damages-only-fiend-or-undead"},
		{"Hunting Trap/Hunting Trap", `{"gear":{"actions":["utilize","action"],"checks":[{"purpose":"avoid-damage-and-halt","kind":"saving-throw","ability":"dexterity","skill":"none","dc":13,"dc_rule":"fixed"},{"purpose":"free-creature","kind":"ability-check","ability":"strength","skill":"athletics","dc":13,"dc_rule":"fixed"}],"damage":[{"purpose":"failed-save","dice":1,"sides":4,"flat":0,"type":"piercing"},{"purpose":"each-failed-escape-check","dice":0,"sides":0,"flat":1,"type":"piercing"}],"measures":[{"name":"failed-save-speed","numerator":0,"denominator":1,"unit":"speed"},{"name":"typical-chain-length","numerator":3,"denominator":1,"unit":"feet"}],"traits":["sawtooth-steel-ring-snaps-shut-on-pressure-plate","chain-affixed-to-immobile-object","speed-penalty-lasts-until-start-of-next-turn","movement-limited-by-chain-until-free","check-frees-self-or-creature-within-reach"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":25,"pounds_denominator":1,"condition":"none"},"source_pages":[97,98]}}`,
			"multiple", "item.gear action.utilize action.action saving-throw.avoid-damage-and-halt-dexterity-dc-13 ability-check.free-creature-strength-athletics-dc-13 damage.failed-save-1d4-piercing damage.each-failed-escape-check-1-piercing measure.failed-save-speed-0-speed measure.typical-chain-length-3-feet trait.sawtooth-steel-ring-snaps-shut-on-pressure-plate trait.chain-affixed-to-immobile-object trait.speed-penalty-lasts-until-start-of-next-turn trait.movement-limited-by-chain-until-free trait.check-frees-self-or-creature-within-reach"},
		{"Ink/Ink", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"bottle-volume","numerator":1,"denominator":1,"unit":"ounce"},{"name":"writable-pages","numerator":500,"denominator":1,"unit":"pages"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[98,98]}}`,
			"none", "item.gear measure.bottle-volume-1-ounce measure.writable-pages-500-pages"},
		{"Ink Pen/Ink Pen", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[],"traits":["writes-or-draws-with-ink"],"contents":[],"spellcasting_focus":[],"weight":{"listed":false,"pounds_numerator":0,"pounds_denominator":1,"condition":"none"},"source_pages":[98,98]}}`,
			"none", "item.gear trait.writes-or-draws-with-ink"},
		{"Jug/Jug", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"capacity","numerator":1,"denominator":1,"unit":"gallon"}],"traits":[],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":4,"pounds_denominator":1,"condition":"none"},"source_pages":[98,98]}}`,
			"none", "item.gear measure.capacity-1-gallon"},
		{"Ladder/Ladder", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"height","numerator":10,"denominator":1,"unit":"feet"}],"traits":["must-climb-to-move-up-or-down"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":25,"pounds_denominator":1,"condition":"none"},"source_pages":[98,98]}}`,
			"none", "item.gear measure.height-10-feet trait.must-climb-to-move-up-or-down"},
		{"Lamp/Lamp", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"bright-light-radius","numerator":15,"denominator":1,"unit":"feet"},{"name":"additional-dim-light","numerator":30,"denominator":1,"unit":"feet"}],"traits":["burns-oil-as-fuel"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[98,98]}}`,
			"none", "item.gear measure.bright-light-radius-15-feet measure.additional-dim-light-30-feet trait.burns-oil-as-fuel"},
		{"Lantern, Bullseye/Lantern, Bullseye", `{"gear":{"actions":[],"checks":[],"damage":[],"measures":[{"name":"bright-light-cone","numerator":60,"denominator":1,"unit":"feet"},{"name":"additional-dim-light","numerator":60,"denominator":1,"unit":"feet"}],"traits":["burns-oil-as-fuel"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1,"condition":"none"},"source_pages":[98,98]}}`,
			"none", "item.gear measure.bright-light-cone-60-feet measure.additional-dim-light-60-feet trait.burns-oil-as-fuel"},
		{"Lantern, Hooded/Lantern, Hooded", `{"gear":{"actions":["bonus-action"],"checks":[],"damage":[],"measures":[{"name":"bright-light-radius","numerator":30,"denominator":1,"unit":"feet"},{"name":"additional-dim-light","numerator":30,"denominator":1,"unit":"feet"},{"name":"lowered-hood-dim-light-radius","numerator":5,"denominator":1,"unit":"feet"}],"traits":["burns-oil-as-fuel","hood-can-be-lowered-or-raised"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1,"condition":"none"},"source_pages":[98,98]}}`,
			"bonus-action", "item.gear action.bonus-action measure.bright-light-radius-30-feet measure.additional-dim-light-30-feet measure.lowered-hood-dim-light-radius-5-feet trait.burns-oil-as-fuel trait.hood-can-be-lowered-or-raised"},
		{"Lock/Lock", `{"gear":{"actions":[],"checks":[{"purpose":"pick-lock-without-key","kind":"ability-check","ability":"dexterity","skill":"sleight-of-hand","dc":15,"dc_rule":"fixed"}],"damage":[],"measures":[],"traits":["comes-with-key","thieves-tools-can-pick-without-key"],"contents":[],"spellcasting_focus":[],"weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1,"condition":"none"},"source_pages":[98,98]}}`,
			"none", "item.gear ability-check.pick-lock-without-key-dexterity-sleight-of-hand-dc-15 trait.comes-with-key trait.thieves-tools-can-pick-without-key"},
	}
	index := 0
	for _, f := range loadItems259(t).Families {
		for _, v := range f.Variants {
			if index >= len(expected) {
				t.Fatal("batch holds more records than the fixture")
			}
			want := expected[index]
			index++
			t.Run(want.name, func(t *testing.T) {
				var details ItemDetails259
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

func TestItems259CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems259(t)
	composed, err := ComposeItemSourceIdentities(append(itemIdentitiesThrough(t, 258), batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, f := range composed {
		total += len(f.VariantIDs)
	}
	if len(composed) != 382 || total != 482 || composed[0].Name != "Potion of Healing" || composed[381].Name != "Lock" {
		t.Fatalf("expected Potion of Healing, batches #225 to #258 and #259: %d families, %d variants", len(composed), total)
	}
	checkItemBatchComposition(t, batch.SourceIdentities(), "adventuring-gear", "items_259")
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems259RejectMissingFields(t *testing.T) {
	rejectMissingItemMembers(t, items259JSON, 785, func(data []byte) error {
		_, err := DecodeItems259(data)
		return err
	})
}

func TestItems259RejectInvalidFixtures(t *testing.T) {
	rejectInvalidGearFixtures(t, LoadItems259, DecodeItems259, itemGear259)
}

// The mechanic keys restate every payload value except weight and the last
// source page, so changing any other value of any record fails.
func TestItems259RejectChangedValues(t *testing.T) {
	rejectChangedGearValues(t, items259JSON, func(data []byte) error {
		_, err := DecodeItems259(data)
		return err
	})
}

func TestItems259RejectUnsupportedJSON(t *testing.T) {
	decode := func(data []byte) error {
		_, err := DecodeItems259(data)
		return err
	}
	rejectUnsupportedItemJSON(t, items259JSON, `"id": "srd-5.2.1:content/items/healers-kit/variant/standard"`, decode, []itemJSONChange{
		{"unknown payload", `"gear": {`, `"gear_plus": {`},
		{"unknown detail", `"spellcasting_focus": `, `"proficiency": `},
		{"number for boolean", `"listed": true`, `"listed": 1`},
		{"string for number", `"pounds_numerator": `, `"pounds_numerator": "1", "x": `},
		{"null check list", `"checks": [`, `"checks": null, "x": [`},
	})
	// Neither the Healing contract decoder nor the previous batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items259JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems258(items259JSON); err == nil {
		t.Fatal("#258 decoder accepted the #259 batch")
	}
	if _, err := DecodeItems259(items258JSON); err == nil {
		t.Fatal("#259 decoder accepted the #258 batch")
	}
}

// The Cleric creation fixture names a Holy Symbol, and the Shield of Faith
// fixture requires one as its focus.
func TestItems259JoinRuntimeGear(t *testing.T) {
	families := loadItems259(t).Families
	joinCatalogueEquipment(t, families, []itemGearKeyWant{
		{"healers-kit", "healers-kit", "standard"},
		{"holy-symbol", "holy-symbol", "amulet"},
		{"holy-symbol", "holy-symbol", "emblem"},
		{"holy-symbol", "holy-symbol", "reliquary"},
	})
	build, err := SupportedBuildForConcept("human cleric", "acolyte")
	if err != nil {
		t.Fatal(err)
	}
	requireCreationNames(t, families, build, map[string]string{"Holy Symbol": "Holy Symbol"})
	symbol := families[slices.IndexFunc(families, func(f ItemSourceFamily[ItemDetails259]) bool { return f.Name == "Holy Symbol" })]
	for _, v := range symbol.Variants {
		if !slices.Contains(v.Details.Gear.SpellcastingFocus, "cleric") {
			t.Fatalf("%s is not a Cleric focus", v.Name)
		}
	}
	if components := shieldSpell(&CastingBuild{Focus: true}).Components; !slices.Contains(components, "material (Holy Symbol focus)") {
		t.Fatalf("the Shield of Faith fixture no longer names the focus: %v", components)
	}
}
