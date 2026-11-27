package srd521

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func loadItems247(t *testing.T) ItemSourceBatch[ItemDetails247] {
	t.Helper()
	b, err := LoadItems247()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities and listed prices from the Weapons table on pinned PDF p. 91, in
// planning inventory order for issue #247.
func TestItems247SourceIdentities(t *testing.T) {
	expected := []struct {
		name, slug string
		cost       int
		currency   string
	}{
		{"Battleaxe", "battleaxe", 10, "gp"},
		{"Blowgun", "blowgun", 10, "gp"},
		{"Club", "club", 1, "sp"},
		{"Dagger", "dagger", 2, "gp"},
		{"Dart", "dart", 5, "cp"},
		{"Flail", "flail", 10, "gp"},
		{"Glaive", "glaive", 20, "gp"},
		{"Greataxe", "greataxe", 30, "gp"},
		{"Greatclub", "greatclub", 2, "sp"},
		{"Greatsword", "greatsword", 50, "gp"},
		{"Halberd", "halberd", 20, "gp"},
		{"Hand Crossbow", "hand-crossbow", 75, "gp"},
	}
	b := loadItems247(t)
	if len(b.Families) != 12 || b.Batch != "data-weapons-01" || b.CatalogueIssue != 247 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve weapons of issue #247")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "p. 91") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	for i, want := range expected {
		t.Run(want.name, func(t *testing.T) {
			f := b.Families[i]
			id := "srd-5.2.1:content/items/" + want.slug
			if f.ID != id || f.Name != want.name || f.Kind != "mundane" || f.Category != "weapon" || len(f.Base.Options) != 0 || len(f.Base.Excluded) != 0 {
				t.Fatalf("source identity mismatch: %+v", f)
			}
			if len(f.Presentations) != 1 || f.Presentations[0].ID != "srd-5.2.1:content/weapons/"+want.slug || !slices.Equal(f.Presentations[0].Pages, []int{91, 91}) {
				t.Fatal("source presentation mismatch")
			}
			if len(f.Variants) != 1 {
				t.Fatalf("expected one table row, got %d variants", len(f.Variants))
			}
			v := f.Variants[0]
			if v.ID != id+"/variant/standard" || v.Name != want.name || v.Rarity != "none" || v.Attunement.Mode != "none" || *v.Charges != (ItemCharges{Mode: "none", Recharge: "none"}) {
				t.Fatalf("variant mismatch: %+v", v)
			}
			if v.Cost.Kind != "listed" || v.Cost.Amount != want.cost || v.Cost.Currency != want.currency || !slices.Equal(v.Cost.SourcePages, []int{91, 91}) {
				t.Fatalf("listed price mismatch: %+v", v.Cost)
			}
			if v.MechanicsAvailable() {
				t.Fatal("catalogue record became available")
			}
		})
	}
	// Join every weapon to the decision inventory the issue cites, in its order.
	data, err := os.ReadFile("../../../docs/investigations/issue-80-spell-item-inventory.json")
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Records []struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Family string `json:"family"`
			Page   int    `json:"source_page"`
			Batch  string `json:"batch_key"`
			Issue  int    `json:"catalogue_issue"`
			Status string `json:"status"`
		} `json:"records"`
	}
	if err := json.Unmarshal(data, &inventory); err != nil {
		t.Fatal(err)
	}
	index := 0
	for _, r := range inventory.Records {
		if r.Issue != 247 {
			continue
		}
		if r.Family != "weapons" || r.Batch != b.Batch || index >= len(b.Families) {
			t.Fatalf("inventory record %s changed family or batch", r.ID)
		}
		f := b.Families[index]
		if r.Status != "unsupported" || r.ID != f.Presentations[0].ID || r.Name != f.Name || r.Page != 91 {
			t.Fatalf("inventory record %s disagrees with batch weapon %s", r.ID, f.Name)
		}
		index++
	}
	if index != 12 {
		t.Fatalf("inventory assigns %d weapons to issue #247", index)
	}
	b.Families[0].Variants[0].Details.Weapon.Mastery = "changed"
	if fresh := loadItems247(t); fresh.Families[0].Variants[0].Details.Weapon.Mastery != "topple" {
		t.Fatal("load must return independent data")
	}
}

// Literal table rows for all 12 weapons, transcribed from pinned PDF p. 91.
func TestItems247SourceValues(t *testing.T) {
	expected := []string{
		`{"weapon":{"proficiency_category":"martial","attack_kind":"melee","damage":{"dice":1,"sides":8,"flat":0,"type":"slashing"},"properties":["versatile"],"range":{"normal_feet":0,"long_feet":0},"ammunition":"none","versatile_damage":{"dice":1,"sides":10},"two_handed_exception":"none","mastery":"topple","weight":{"listed":true,"pounds_numerator":4,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"martial","attack_kind":"ranged","damage":{"dice":0,"sides":0,"flat":1,"type":"piercing"},"properties":["ammunition","loading"],"range":{"normal_feet":25,"long_feet":100},"ammunition":"needle","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"vex","weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"simple","attack_kind":"melee","damage":{"dice":1,"sides":4,"flat":0,"type":"bludgeoning"},"properties":["light"],"range":{"normal_feet":0,"long_feet":0},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"slow","weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"simple","attack_kind":"melee","damage":{"dice":1,"sides":4,"flat":0,"type":"piercing"},"properties":["finesse","light","thrown"],"range":{"normal_feet":20,"long_feet":60},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"nick","weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"simple","attack_kind":"ranged","damage":{"dice":1,"sides":4,"flat":0,"type":"piercing"},"properties":["finesse","thrown"],"range":{"normal_feet":20,"long_feet":60},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"vex","weight":{"listed":true,"pounds_numerator":1,"pounds_denominator":4}}}`,
		`{"weapon":{"proficiency_category":"martial","attack_kind":"melee","damage":{"dice":1,"sides":8,"flat":0,"type":"bludgeoning"},"properties":[],"range":{"normal_feet":0,"long_feet":0},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"sap","weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"martial","attack_kind":"melee","damage":{"dice":1,"sides":10,"flat":0,"type":"slashing"},"properties":["heavy","reach","two-handed"],"range":{"normal_feet":0,"long_feet":0},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"graze","weight":{"listed":true,"pounds_numerator":6,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"martial","attack_kind":"melee","damage":{"dice":1,"sides":12,"flat":0,"type":"slashing"},"properties":["heavy","two-handed"],"range":{"normal_feet":0,"long_feet":0},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"cleave","weight":{"listed":true,"pounds_numerator":7,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"simple","attack_kind":"melee","damage":{"dice":1,"sides":8,"flat":0,"type":"bludgeoning"},"properties":["two-handed"],"range":{"normal_feet":0,"long_feet":0},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"push","weight":{"listed":true,"pounds_numerator":10,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"martial","attack_kind":"melee","damage":{"dice":2,"sides":6,"flat":0,"type":"slashing"},"properties":["heavy","two-handed"],"range":{"normal_feet":0,"long_feet":0},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"graze","weight":{"listed":true,"pounds_numerator":6,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"martial","attack_kind":"melee","damage":{"dice":1,"sides":10,"flat":0,"type":"slashing"},"properties":["heavy","reach","two-handed"],"range":{"normal_feet":0,"long_feet":0},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"cleave","weight":{"listed":true,"pounds_numerator":6,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"martial","attack_kind":"ranged","damage":{"dice":1,"sides":6,"flat":0,"type":"piercing"},"properties":["ammunition","light","loading"],"range":{"normal_feet":30,"long_feet":120},"ammunition":"bolt","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"vex","weight":{"listed":true,"pounds_numerator":3,"pounds_denominator":1}}}`,
	}
	keys := []string{
		"item.wield weapon.attack weapon.martial weapon.melee damage.1d8-slashing property.versatile damage.two-handed-1d10 mastery.topple",
		"item.wield weapon.attack weapon.martial weapon.ranged damage.1-piercing property.ammunition property.loading range.25-100 ammunition.needle mastery.vex",
		"item.wield weapon.attack weapon.simple weapon.melee damage.1d4-bludgeoning property.light mastery.slow",
		"item.wield weapon.attack weapon.simple weapon.melee damage.1d4-piercing property.finesse property.light property.thrown range.20-60 mastery.nick",
		"item.wield weapon.attack weapon.simple weapon.ranged damage.1d4-piercing property.finesse property.thrown range.20-60 mastery.vex",
		"item.wield weapon.attack weapon.martial weapon.melee damage.1d8-bludgeoning mastery.sap",
		"item.wield weapon.attack weapon.martial weapon.melee damage.1d10-slashing property.heavy property.reach property.two-handed mastery.graze",
		"item.wield weapon.attack weapon.martial weapon.melee damage.1d12-slashing property.heavy property.two-handed mastery.cleave",
		"item.wield weapon.attack weapon.simple weapon.melee damage.1d8-bludgeoning property.two-handed mastery.push",
		"item.wield weapon.attack weapon.martial weapon.melee damage.2d6-slashing property.heavy property.two-handed mastery.graze",
		"item.wield weapon.attack weapon.martial weapon.melee damage.1d10-slashing property.heavy property.reach property.two-handed mastery.cleave",
		"item.wield weapon.attack weapon.martial weapon.ranged damage.1d6-piercing property.ammunition property.light property.loading range.30-120 ammunition.bolt mastery.vex",
	}
	b := loadItems247(t)
	if len(b.Families) != len(expected) || len(keys) != len(expected) {
		t.Fatal("fixture count disagrees with the batch")
	}
	for i, f := range b.Families {
		t.Run(f.Name, func(t *testing.T) {
			v := f.Variants[0]
			var want ItemDetails247
			d := json.NewDecoder(strings.NewReader(expected[i]))
			d.DisallowUnknownFields()
			if err := d.Decode(&want); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(v.Details, &want) {
				t.Fatalf("PDF p.91 row mismatch: %+v", v.Details.Weapon)
			}
			if v.Equip.Mode != "wield" || v.Use.Mode != "weapon-attack" || len(v.Equip.Restrictions) != 0 || len(v.Use.Restrictions) != 0 || strings.Join(v.RequiredMechanicKeys, " ") != keys[i] {
				t.Fatalf("source requirements mismatch: %v", v.RequiredMechanicKeys)
			}
		})
	}
}

// The fixed Fighter creation fixture names weapons but holds no weapon record.
// Its names must match this catalogue, which stays the only record of each row.
func TestItems247JoinRuntimeWeaponNames(t *testing.T) {
	build, err := NewSupportedMartialBuild("Human", "Fighter", "Soldier")
	if err != nil {
		t.Fatal(err)
	}
	rows := map[string]*ItemWeapon{}
	for _, f := range loadItems247(t).Families {
		rows[f.Name] = f.Variants[0].Details.Weapon
	}
	for _, name := range []string{"Greatsword", "Flail"} {
		if rows[name] == nil || !slices.Contains(build.Equipment, name) {
			t.Fatalf("creation fixture weapon %s does not join to the catalogue", name)
		}
	}
	// The fixture records Greatsword mastery; the row names which property.
	if !slices.Contains(build.WeaponMastery, "Greatsword") || rows["Greatsword"].Mastery != "graze" {
		t.Fatal("Greatsword mastery join changed")
	}
}

func TestItems247CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems247(t)
	contract, err := LoadItemCatalogue()
	if err != nil {
		t.Fatal(err)
	}
	identities := [][]ItemSourceIdentity{contract.SourceIdentities()}
	for _, load := range []func() ([]ItemSourceIdentity, error){
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems225(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems226(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems227(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems228(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems229(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems230(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems231(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems232(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems233(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems234(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems235(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems236(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems237(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems238(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems239(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems240(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems241(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems242(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems243(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems244(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems245(); return b.SourceIdentities(), err },
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems246(); return b.SourceIdentities(), err },
	} {
		ids, err := load()
		if err != nil {
			t.Fatal(err)
		}
		identities = append(identities, ids)
	}
	composed, err := ComposeItemSourceIdentities(append(identities, batch.SourceIdentities())...)
	if err != nil {
		t.Fatal(err)
	}
	variants := 0
	for _, f := range composed {
		variants += len(f.VariantIDs)
	}
	if len(composed) != 270 || variants != 346 || composed[0].Name != "Potion of Healing" || composed[258].Name != "Battleaxe" || composed[269].Name != "Hand Crossbow" {
		t.Fatalf("expected Potion of Healing, batches #225 to #246 and the twelve weapons of #247: %d families, %d variants", len(composed), variants)
	}
	if _, err := ComposeItemSourceIdentities(batch.SourceIdentities(), batch.SourceIdentities()); err == nil {
		t.Fatal("duplicate batch accepted")
	}
	clash := batch.SourceIdentities()[:1]
	clash[0].FamilyID, clash[0].Name, clash[0].PresentationIDs = "srd-5.2.1:content/items/other", "Other", []string{"srd-5.2.1:content/weapons/other"}
	if _, err := ComposeItemSourceIdentities(batch.SourceIdentities(), clash); err == nil {
		t.Fatal("duplicate variant identity accepted")
	}
	// The catalogue is descriptive: it must not promote coverage or add entries.
	register, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range register.Entries {
		for _, fixture := range entry.Fixtures {
			if strings.Contains(fixture, "items_247") {
				t.Fatal("catalogue batch must not be a coverage fixture")
			}
		}
	}
	if report := register.Coverage(); report.PlannedSpellItemScopes != 753 || report.FullSupport {
		t.Fatal("catalogue data changed the coverage boundary")
	}
}

// Every field of every record is required: deleting any single member anywhere
// in the document, including zero, false and empty-list members, fails.
func TestItems247RejectMissingFields(t *testing.T) {
	var paths [][]any
	var walk func(node any, path []any)
	walk = func(node any, path []any) {
		switch n := node.(type) {
		case map[string]any:
			for key, child := range n {
				p := append(slices.Clone(path), key)
				paths = append(paths, p)
				walk(child, p)
			}
		case []any:
			for i, child := range n {
				walk(child, append(slices.Clone(path), i))
			}
		}
	}
	var root any
	if err := json.Unmarshal(items247JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 600 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items247JSON, &fresh); err != nil {
			t.Fatal(err)
		}
		node := fresh
		for _, step := range path[:len(path)-1] {
			switch s := step.(type) {
			case string:
				node = node.(map[string]any)[s]
			case int:
				node = node.([]any)[s]
			}
		}
		delete(node.(map[string]any), path[len(path)-1].(string))
		data, err := json.Marshal(fresh)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := DecodeItems247(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems247RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails247]
	// v and row address one weapon by index; see the source identity test.
	v := func(b *batch, family int) *ItemSourceVariant[ItemDetails247] { return &b.Families[family].Variants[0] }
	row := func(b *batch, family int) *ItemWeapon { return v(b, family).Details.Weapon }
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-weapons-04" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 246 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing weapon", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered weapon", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate weapon", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent weapon", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/whip"
			b.Families[12].Name = "Whip"
		}},
		{"unknown weapon", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/unreviewed" }},
		{"weapon name", func(b *batch) { b.Families[0].Name = "Unreviewed Name" }},
		{"magic kind", func(b *batch) { b.Families[0].Kind = "magic" }},
		{"category", func(b *batch) { b.Families[0].Category = "wondrous-item" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"any-simple-weapon"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{90, 90} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"magic presentation", func(b *batch) { b.Families[0].Presentations[0].ID = "srd-5.2.1:content/magic-items/battleaxe" }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"variant id", func(b *batch) { v(b, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0).Name = "Unreviewed" }},
		{"invented rarity", func(b *batch) { v(b, 0).Rarity = "common" }},
		{"attunement", func(b *batch) { v(b, 0).Attunement.Mode = "required" }},
		{"invented charges", func(b *batch) { v(b, 0).Charges.Maximum = 1 }},
		{"equip mode", func(b *batch) { v(b, 0).Equip.Mode = "carry" }},
		{"use mode", func(b *batch) { v(b, 0).Use.Mode = "none" }},
		{"invented use restriction", func(b *batch) { v(b, 0).Use.Restrictions = []string{"proficiency"} }},
		{"unlisted price", func(b *batch) { v(b, 0).Cost.Kind = "unspecified" }},
		{"changed price", func(b *batch) { v(b, 0).Cost.Amount++ }},
		{"changed currency", func(b *batch) { v(b, 0).Cost.Currency = "ep" }},
		{"cost page", func(b *batch) { v(b, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys = v(b, 0).RequiredMechanicKeys[:1] }},
		{"unknown mechanic", func(b *batch) { v(b, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0).Details = &ItemDetails247{} }},
		{"another weapon's row", func(b *batch) { v(b, 0).Details = v(b, 1).Details }},
		{"proficiency category", func(b *batch) { row(b, 0).ProficiencyCategory = "exotic" }},
		{"attack kind", func(b *batch) { row(b, 0).AttackKind = "thrown" }},
		{"damage die", func(b *batch) { row(b, 0).Damage.Sides += 2 }},
		{"two damage forms", func(b *batch) { row(b, 0).Damage.Dice, row(b, 0).Damage.Flat = 1, 1 }},
		{"damage type", func(b *batch) { row(b, 0).Damage.Type = "fire" }},
		{"unknown property", func(b *batch) { row(b, 1).Properties[0] = "keen" }},
		{"reordered properties", func(b *batch) { p := row(b, 1).Properties; p[0], p[1] = p[1], p[0] }},
		{"duplicate property", func(b *batch) { p := row(b, 1); p.Properties = append(p.Properties, p.Properties[len(p.Properties)-1]) }},
		{"dropped properties", func(b *batch) { row(b, 1).Properties = []string{} }},
		{"dropped range", func(b *batch) { row(b, 1).Range = ItemWeaponRange{} }},
		{"changed range", func(b *batch) { row(b, 1).Range.LongFeet += 5 }},
		{"dropped ammunition", func(b *batch) { row(b, 1).Ammunition = "none" }},
		{"changed ammunition", func(b *batch) { row(b, 1).Ammunition = "arrow" }},
		{"dropped two-handed damage", func(b *batch) { row(b, 0).VersatileDamage = ItemWeaponVersatility{} }},
		{"changed two-handed damage", func(b *batch) { row(b, 0).VersatileDamage.Sides += 2 }},
		{"invented mounted exception", func(b *batch) { row(b, 0).TwoHandedException = "while-prone" }},
		{"changed mastery", func(b *batch) { row(b, 0).Mastery = "vex" }},
		{"unlisted weight with value", func(b *batch) { row(b, 0).Weight.Listed = false }},
		{"zero weight denominator", func(b *batch) { row(b, 0).Weight.PoundsDenominator = 0 }},
		{"negative weight", func(b *batch) { row(b, 0).Weight.PoundsNumerator = -1 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems247(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems247(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems247RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/battleaxe/variant/standard"`)
	if !bytes.Contains(items247JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"proficient":true,`, `"playable":true,`, `"available":true,`, `"properties":{"attack_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items247JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems247(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"weapon": {`, `"weapon_plus": {`},
		"unknown detail":     {`"normal_feet": `, `"normal_feet_all": `},
		"string for number":  {`"normal_feet": 0`, `"normal_feet": "0"`},
		"fractional number":  {`"normal_feet": 0`, `"normal_feet": 0.5`},
		"number for boolean": {`"listed": true`, `"listed": 1`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items247JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items247JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems247(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items247JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems247(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items247JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems246(items247JSON); err == nil {
		t.Fatal("#246 decoder accepted the #247 batch")
	}
	if _, err := DecodeItems247(items246JSON); err == nil {
		t.Fatal("#247 decoder accepted the #246 batch")
	}
}
