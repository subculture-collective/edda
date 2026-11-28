package srd521

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func loadItems250(t *testing.T) ItemSourceBatch[ItemDetails250] {
	t.Helper()
	b, err := LoadItems250()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities, variants and listed prices from pinned PDF p. 91 for issue #250.
func TestItems250SourceIdentities(t *testing.T) {
	expected := []itemFamilyWant{
		{"Whip", "whip", 91, []itemVariantWant{{"standard", "Whip", 2, "gp"}}},
		{"Heavy Crossbow", "heavy-crossbow", 91, []itemVariantWant{{"standard", "Heavy Crossbow", 50, "gp"}}},
	}
	b := loadItems250(t)
	if len(b.Families) != 2 || b.Batch != "data-weapons-04" || b.CatalogueIssue != 250 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the weapons of issue #250")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "p. 91") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	checkMundaneIdentities(t, b.Families, "weapon", "weapons", expected)
	// The inventory assigns only the Whip to #250 and omits Heavy Crossbow.
	joinItemInventory(t, b.SourceIdentities()[:1], b.Batch, 250, "weapons", []int{91})
	requireInventoryOmits(t, "Heavy Crossbow")
	b.Families[0].Variants[0].Details.Weapon.Mastery = "changed"
	if fresh := loadItems250(t); fresh.Families[0].Variants[0].Details.Weapon.Mastery != "slow" {
		t.Fatal("load must return independent data")
	}
}

// Literal values, equip modes and mechanic keys for all 2 records,
// transcribed from pinned PDF p. 91.
func TestItems250SourceValues(t *testing.T) {
	expected := []struct{ name, details, equip, keys string }{
		{"Whip/Whip", `{"weapon":{"proficiency_category":"martial","attack_kind":"melee","damage":{"dice":1,"sides":4,"flat":0,"type":"slashing"},"properties":["finesse","reach"],"range":{"normal_feet":0,"long_feet":0},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"slow","weight":{"listed":true,"pounds_numerator":3,"pounds_denominator":1}}}`,
			"wield", "item.wield weapon.attack weapon.martial weapon.melee damage.1d4-slashing property.finesse property.reach mastery.slow"},
		{"Heavy Crossbow/Heavy Crossbow", `{"weapon":{"proficiency_category":"martial","attack_kind":"ranged","damage":{"dice":1,"sides":10,"flat":0,"type":"piercing"},"properties":["ammunition","heavy","loading","two-handed"],"range":{"normal_feet":100,"long_feet":400},"ammunition":"bolt","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"push","weight":{"listed":true,"pounds_numerator":18,"pounds_denominator":1}}}`,
			"wield", "item.wield weapon.attack weapon.martial weapon.ranged damage.1d10-piercing property.ammunition property.heavy property.loading property.two-handed range.100-400 ammunition.bolt mastery.push"},
	}
	index := 0
	for _, f := range loadItems250(t).Families {
		for _, v := range f.Variants {
			if index >= len(expected) {
				t.Fatal("batch holds more records than the fixture")
			}
			want := expected[index]
			index++
			t.Run(want.name, func(t *testing.T) {
				var details ItemDetails250
				d := json.NewDecoder(strings.NewReader(want.details))
				d.DisallowUnknownFields()
				if err := d.Decode(&details); err != nil {
					t.Fatal(err)
				}
				if f.Name+"/"+v.Name != want.name || !reflect.DeepEqual(v.Details, &details) {
					t.Fatalf("source value mismatch for %s/%s: %+v", f.Name, v.Name, v.Details.Weapon)
				}
				if v.Equip.Mode != want.equip || v.Use.Mode != "weapon-attack" || len(v.Equip.Restrictions) != 0 || len(v.Use.Restrictions) != 0 || strings.Join(v.RequiredMechanicKeys, " ") != want.keys {
					t.Fatalf("source requirements mismatch: %v", v.RequiredMechanicKeys)
				}
			})
		}
	}
	if index != len(expected) {
		t.Fatalf("fixture lists %d records, batch holds %d", len(expected), index)
	}
}

func TestItems250CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems250(t)
	composed, err := ComposeItemSourceIdentities(append(itemIdentitiesThrough(t, 249), batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	for _, f := range composed {
		total += len(f.VariantIDs)
	}
	if len(composed) != 296 || total != 372 || composed[0].Name != "Potion of Healing" || composed[294].Name != "Whip" || composed[295].Name != "Heavy Crossbow" {
		t.Fatalf("expected Potion of Healing, batches #225 to #249 and #250: %d families, %d variants", len(composed), total)
	}
	checkItemBatchComposition(t, batch.SourceIdentities(), "weapons", "items_250")
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems250RejectMissingFields(t *testing.T) {
	rejectMissingItemMembers(t, items250JSON, 132, func(data []byte) error {
		_, err := DecodeItems250(data)
		return err
	})
}

func TestItems250RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails250]
	// v and row address a family's first variant by family index.
	v := func(b *batch, family int) *ItemSourceVariant[ItemDetails250] { return &b.Families[family].Variants[0] }
	row := func(b *batch, family int) *ItemWeapon { return v(b, family).Details.Weapon }
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-other" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 249 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:1] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[2].ID = "srd-5.2.1:content/items/other"
			b.Families[2].Name = "Other"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/unreviewed" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Unreviewed Name" }},
		{"magic kind", func(b *batch) { b.Families[0].Kind = "magic" }},
		{"category", func(b *batch) { b.Families[0].Category = "wondrous-item" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"any-simple-weapon"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{90, 90} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"magic presentation", func(b *batch) { b.Families[0].Presentations[0].ID = "srd-5.2.1:content/magic-items/whip" }},
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
		{"empty details", func(b *batch) { v(b, 0).Details = &ItemDetails250{} }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"another family's row", func(b *batch) { v(b, 0).Details = v(b, 1).Details }},
		{"proficiency category", func(b *batch) { row(b, 0).ProficiencyCategory = "exotic" }},
		{"attack kind", func(b *batch) { row(b, 0).AttackKind = "thrown" }},
		{"damage die", func(b *batch) { row(b, 0).Damage.Sides += 2 }},
		{"two damage forms", func(b *batch) { row(b, 0).Damage.Dice, row(b, 0).Damage.Flat = 1, 1 }},
		{"damage type", func(b *batch) { row(b, 0).Damage.Type = "fire" }},
		{"unknown property", func(b *batch) { row(b, 0).Properties[0] = "keen" }},
		{"reordered properties", func(b *batch) {
			p := row(b, 0).Properties
			p[0], p[1] = p[1], p[0]
		}},
		{"duplicate property", func(b *batch) {
			p := row(b, 0)
			p.Properties = append(p.Properties, p.Properties[len(p.Properties)-1])
		}},
		{"dropped properties", func(b *batch) { row(b, 0).Properties = []string{} }},
		{"invented range", func(b *batch) { row(b, 0).Range = ItemWeaponRange{NormalFeet: 10, LongFeet: 20} }},
		{"dropped range", func(b *batch) { row(b, 1).Range = ItemWeaponRange{} }},
		{"changed range", func(b *batch) { row(b, 1).Range.LongFeet += 5 }},
		{"dropped ammunition", func(b *batch) { row(b, 1).Ammunition = "none" }},
		{"changed ammunition", func(b *batch) { row(b, 1).Ammunition = "arrow" }},
		{"invented two-handed damage", func(b *batch) { row(b, 0).VersatileDamage = ItemWeaponVersatility{Dice: 1, Sides: 6} }},
		{"invented mounted exception", func(b *batch) { row(b, 1).TwoHandedException = "unless-mounted" }},
		{"changed mastery", func(b *batch) { row(b, 0).Mastery = "sap" }},
		{"unlisted weight with value", func(b *batch) { row(b, 0).Weight.Listed = false }},
		{"zero weight denominator", func(b *batch) { row(b, 0).Weight.PoundsDenominator = 0 }},
		{"negative weight", func(b *batch) { row(b, 0).Weight.PoundsNumerator = -1 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems250(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems250(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems250RejectUnsupportedJSON(t *testing.T) {
	decode := func(data []byte) error {
		_, err := DecodeItems250(data)
		return err
	}
	rejectUnsupportedItemJSON(t, items250JSON, `"id": "srd-5.2.1:content/items/whip/variant/standard"`, decode, []itemJSONChange{
		{"unknown payload", `"weapon": {`, `"weapon_plus": {`},
		{"unknown detail", `"normal_feet": `, `"normal_feet_all": `},
		{"string for number", `"normal_feet": 0`, `"normal_feet": "0"`},
		{"fractional number", `"normal_feet": 0`, `"normal_feet": 0.5`},
		{"number for boolean", `"listed": true`, `"listed": 1`},
	})
	// Neither the Healing contract decoder nor the previous batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items250JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems249(items250JSON); err == nil {
		t.Fatal("#249 decoder accepted the #250 batch")
	}
	if _, err := DecodeItems250(items249JSON); err == nil {
		t.Fatal("#250 decoder accepted the #249 batch")
	}
}
