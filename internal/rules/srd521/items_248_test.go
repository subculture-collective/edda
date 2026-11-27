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

func loadItems248(t *testing.T) ItemSourceBatch[ItemDetails248] {
	t.Helper()
	b, err := LoadItems248()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities and listed prices from the Weapons table on pinned PDF p. 91, in
// planning inventory order for issue #248.
func TestItems248SourceIdentities(t *testing.T) {
	expected := []struct {
		name, slug string
		cost       int
		currency   string
	}{
		{"Handaxe", "handaxe", 5, "gp"},
		{"Javelin", "javelin", 5, "sp"},
		{"Lance", "lance", 10, "gp"},
		{"Light Crossbow", "light-crossbow", 25, "gp"},
		{"Light Hammer", "light-hammer", 2, "gp"},
		{"Longbow", "longbow", 50, "gp"},
		{"Longsword", "longsword", 15, "gp"},
		{"Mace", "mace", 5, "gp"},
		{"Maul", "maul", 10, "gp"},
		{"Morningstar", "morningstar", 15, "gp"},
		{"Musket", "musket", 500, "gp"},
		{"Pike", "pike", 5, "gp"},
	}
	b := loadItems248(t)
	if len(b.Families) != 12 || b.Batch != "data-weapons-02" || b.CatalogueIssue != 248 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve weapons of issue #248")
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
		if r.Issue != 248 {
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
		t.Fatalf("inventory assigns %d weapons to issue #248", index)
	}
	b.Families[0].Variants[0].Details.Weapon.Mastery = "changed"
	if fresh := loadItems248(t); fresh.Families[0].Variants[0].Details.Weapon.Mastery != "vex" {
		t.Fatal("load must return independent data")
	}
}

// Literal table rows for all 12 weapons, transcribed from pinned PDF p. 91.
func TestItems248SourceValues(t *testing.T) {
	expected := []string{
		`{"weapon":{"proficiency_category":"simple","attack_kind":"melee","damage":{"dice":1,"sides":6,"flat":0,"type":"slashing"},"properties":["light","thrown"],"range":{"normal_feet":20,"long_feet":60},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"vex","weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"simple","attack_kind":"melee","damage":{"dice":1,"sides":6,"flat":0,"type":"piercing"},"properties":["thrown"],"range":{"normal_feet":30,"long_feet":120},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"slow","weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"martial","attack_kind":"melee","damage":{"dice":1,"sides":10,"flat":0,"type":"piercing"},"properties":["heavy","reach","two-handed"],"range":{"normal_feet":0,"long_feet":0},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"unless-mounted","mastery":"topple","weight":{"listed":true,"pounds_numerator":6,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"simple","attack_kind":"ranged","damage":{"dice":1,"sides":8,"flat":0,"type":"piercing"},"properties":["ammunition","loading","two-handed"],"range":{"normal_feet":80,"long_feet":320},"ammunition":"bolt","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"slow","weight":{"listed":true,"pounds_numerator":5,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"simple","attack_kind":"melee","damage":{"dice":1,"sides":4,"flat":0,"type":"bludgeoning"},"properties":["light","thrown"],"range":{"normal_feet":20,"long_feet":60},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"nick","weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"martial","attack_kind":"ranged","damage":{"dice":1,"sides":8,"flat":0,"type":"piercing"},"properties":["ammunition","heavy","two-handed"],"range":{"normal_feet":150,"long_feet":600},"ammunition":"arrow","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"slow","weight":{"listed":true,"pounds_numerator":2,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"martial","attack_kind":"melee","damage":{"dice":1,"sides":8,"flat":0,"type":"slashing"},"properties":["versatile"],"range":{"normal_feet":0,"long_feet":0},"ammunition":"none","versatile_damage":{"dice":1,"sides":10},"two_handed_exception":"none","mastery":"sap","weight":{"listed":true,"pounds_numerator":3,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"simple","attack_kind":"melee","damage":{"dice":1,"sides":6,"flat":0,"type":"bludgeoning"},"properties":[],"range":{"normal_feet":0,"long_feet":0},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"sap","weight":{"listed":true,"pounds_numerator":4,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"martial","attack_kind":"melee","damage":{"dice":2,"sides":6,"flat":0,"type":"bludgeoning"},"properties":["heavy","two-handed"],"range":{"normal_feet":0,"long_feet":0},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"topple","weight":{"listed":true,"pounds_numerator":10,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"martial","attack_kind":"melee","damage":{"dice":1,"sides":8,"flat":0,"type":"piercing"},"properties":[],"range":{"normal_feet":0,"long_feet":0},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"sap","weight":{"listed":true,"pounds_numerator":4,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"martial","attack_kind":"ranged","damage":{"dice":1,"sides":12,"flat":0,"type":"piercing"},"properties":["ammunition","loading","two-handed"],"range":{"normal_feet":40,"long_feet":120},"ammunition":"bullet","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"slow","weight":{"listed":true,"pounds_numerator":10,"pounds_denominator":1}}}`,
		`{"weapon":{"proficiency_category":"martial","attack_kind":"melee","damage":{"dice":1,"sides":10,"flat":0,"type":"piercing"},"properties":["heavy","reach","two-handed"],"range":{"normal_feet":0,"long_feet":0},"ammunition":"none","versatile_damage":{"dice":0,"sides":0},"two_handed_exception":"none","mastery":"push","weight":{"listed":true,"pounds_numerator":18,"pounds_denominator":1}}}`,
	}
	keys := []string{
		"item.wield weapon.attack weapon.simple weapon.melee damage.1d6-slashing property.light property.thrown range.20-60 mastery.vex",
		"item.wield weapon.attack weapon.simple weapon.melee damage.1d6-piercing property.thrown range.30-120 mastery.slow",
		"item.wield weapon.attack weapon.martial weapon.melee damage.1d10-piercing property.heavy property.reach property.two-handed property.two-handed-unless-mounted mastery.topple",
		"item.wield weapon.attack weapon.simple weapon.ranged damage.1d8-piercing property.ammunition property.loading property.two-handed range.80-320 ammunition.bolt mastery.slow",
		"item.wield weapon.attack weapon.simple weapon.melee damage.1d4-bludgeoning property.light property.thrown range.20-60 mastery.nick",
		"item.wield weapon.attack weapon.martial weapon.ranged damage.1d8-piercing property.ammunition property.heavy property.two-handed range.150-600 ammunition.arrow mastery.slow",
		"item.wield weapon.attack weapon.martial weapon.melee damage.1d8-slashing property.versatile damage.two-handed-1d10 mastery.sap",
		"item.wield weapon.attack weapon.simple weapon.melee damage.1d6-bludgeoning mastery.sap",
		"item.wield weapon.attack weapon.martial weapon.melee damage.2d6-bludgeoning property.heavy property.two-handed mastery.topple",
		"item.wield weapon.attack weapon.martial weapon.melee damage.1d8-piercing mastery.sap",
		"item.wield weapon.attack weapon.martial weapon.ranged damage.1d12-piercing property.ammunition property.loading property.two-handed range.40-120 ammunition.bullet mastery.slow",
		"item.wield weapon.attack weapon.martial weapon.melee damage.1d10-piercing property.heavy property.reach property.two-handed mastery.push",
	}
	b := loadItems248(t)
	if len(b.Families) != len(expected) || len(keys) != len(expected) {
		t.Fatal("fixture count disagrees with the batch")
	}
	for i, f := range b.Families {
		t.Run(f.Name, func(t *testing.T) {
			v := f.Variants[0]
			var want ItemDetails248
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

// The playable Longsword is a hard-coded fixture in attack.go with a coverage
// register entry, not an item record. This batch's row is the only record; the
// fixture's damage die and type and the register's ID must agree with it.
func TestItems248JoinLongswordFixture(t *testing.T) {
	var family ItemSourceFamily[ItemDetails248]
	for _, f := range loadItems248(t).Families {
		if f.Name == "Longsword" {
			family = f
		}
	}
	if family.ID == "" {
		t.Fatal("Longsword record missing")
	}
	row := family.Variants[0].Details.Weapon
	stats, err := os.ReadFile("testdata/fighter1.json")
	if err != nil {
		t.Fatal(err)
	}
	sheet, err := Project(stats, 1, 12, 12, 0)
	if err != nil {
		t.Fatal(err)
	}
	// A natural 11 against Armor Class 10 hits; the die then rolls its maximum.
	dice := &items248Dice{values: []int{10, row.Damage.Sides - 1}}
	receipt, err := ResolveLongsword(sheet, 10, dice)
	if err != nil {
		t.Fatal(err)
	}
	if row.Damage.Dice != 1 || !slices.Equal(dice.sides, []int{20, row.Damage.Sides}) || receipt.DamageType != row.Damage.Type {
		t.Fatalf("Longsword fixture disagrees with the table row: sides %v, %s", dice.sides, receipt.DamageType)
	}
	if !slices.Contains(receipt.SourceIDs, family.Presentations[0].ID) {
		t.Fatal("attack receipt no longer cites the catalogued presentation")
	}
	register, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	joined := 0
	for _, entry := range register.Entries {
		if entry.ID != family.Presentations[0].ID {
			continue
		}
		joined++
		if entry.Name != family.Name || entry.Status != "unsupported" || !slices.Contains(entry.SourcePages, family.Presentations[0].Pages[0]) {
			t.Fatalf("coverage entry disagrees with the catalogue: %+v", entry)
		}
	}
	if joined != 1 {
		t.Fatalf("expected one Longsword coverage entry, got %d", joined)
	}
}

type items248Dice struct {
	values, sides []int
}

func (d *items248Dice) Intn(n int) int {
	v := d.values[len(d.sides)]
	d.sides = append(d.sides, n)
	return v
}

// The Javelin appears by name in the creation fixture, and the Bandit's Light
// Crossbow action uses the weapon's damage die and range.
func TestItems248JoinRuntimeWeaponUses(t *testing.T) {
	build, err := NewSupportedMartialBuild("Human", "Fighter", "Soldier")
	if err != nil {
		t.Fatal(err)
	}
	rows := map[string]*ItemWeapon{}
	for _, f := range loadItems248(t).Families {
		rows[f.Name] = f.Variants[0].Details.Weapon
	}
	if rows["Javelin"] == nil || !slices.Contains(build.Equipment, "Javelin (8)") {
		t.Fatal("creation fixture Javelin does not join to the catalogue")
	}
	bandit, _, err := BanditRecord()
	if err != nil {
		t.Fatal(err)
	}
	joined := 0
	for _, a := range bandit.Actions {
		if a.ID != "light-crossbow" {
			continue
		}
		joined++
		row := rows["Light Crossbow"]
		if a.Attack == nil || a.Attack.Damage.Dice != row.Damage.Dice || a.Attack.Damage.Sides != row.Damage.Sides || a.Attack.DamageType != row.Damage.Type || a.Attack.Range != row.Range.NormalFeet || a.Attack.LongRange != row.Range.LongFeet {
			t.Fatalf("Bandit Light Crossbow disagrees with the table row: %+v", a.Attack)
		}
	}
	if joined != 1 {
		t.Fatalf("expected one Bandit Light Crossbow action, got %d", joined)
	}
}

func TestItems248CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems248(t)
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
		func() ([]ItemSourceIdentity, error) { b, err := LoadItems247(); return b.SourceIdentities(), err },
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
	if len(composed) != 282 || variants != 358 || composed[0].Name != "Potion of Healing" || composed[270].Name != "Handaxe" || composed[281].Name != "Pike" {
		t.Fatalf("expected Potion of Healing, batches #225 to #247 and the twelve weapons of #248: %d families, %d variants", len(composed), variants)
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
			if strings.Contains(fixture, "items_248") {
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
func TestItems248RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items248JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 600 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items248JSON, &fresh); err != nil {
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
		if _, err := DecodeItems248(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems248RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails248]
	// v and row address one weapon by index; see the source identity test.
	v := func(b *batch, family int) *ItemSourceVariant[ItemDetails248] { return &b.Families[family].Variants[0] }
	row := func(b *batch, family int) *ItemWeapon { return v(b, family).Details.Weapon }
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-weapons-04" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 247 }},
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
		{"magic presentation", func(b *batch) { b.Families[0].Presentations[0].ID = "srd-5.2.1:content/magic-items/handaxe" }},
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
		{"empty details", func(b *batch) { v(b, 0).Details = &ItemDetails248{} }},
		{"another weapon's row", func(b *batch) { v(b, 0).Details = v(b, 1).Details }},
		{"proficiency category", func(b *batch) { row(b, 0).ProficiencyCategory = "exotic" }},
		{"attack kind", func(b *batch) { row(b, 0).AttackKind = "thrown" }},
		{"damage die", func(b *batch) { row(b, 0).Damage.Sides += 2 }},
		{"two damage forms", func(b *batch) { row(b, 0).Damage.Dice, row(b, 0).Damage.Flat = 1, 1 }},
		{"damage type", func(b *batch) { row(b, 0).Damage.Type = "fire" }},
		{"unknown property", func(b *batch) { row(b, 0).Properties[0] = "keen" }},
		{"reordered properties", func(b *batch) { p := row(b, 0).Properties; p[0], p[1] = p[1], p[0] }},
		{"duplicate property", func(b *batch) { p := row(b, 0); p.Properties = append(p.Properties, p.Properties[len(p.Properties)-1]) }},
		{"dropped properties", func(b *batch) { row(b, 0).Properties = []string{} }},
		{"dropped range", func(b *batch) { row(b, 0).Range = ItemWeaponRange{} }},
		{"changed range", func(b *batch) { row(b, 0).Range.LongFeet += 5 }},
		{"dropped ammunition", func(b *batch) { row(b, 3).Ammunition = "none" }},
		{"changed ammunition", func(b *batch) { row(b, 3).Ammunition = "arrow" }},
		{"dropped two-handed damage", func(b *batch) { row(b, 6).VersatileDamage = ItemWeaponVersatility{} }},
		{"changed two-handed damage", func(b *batch) { row(b, 6).VersatileDamage.Sides += 2 }},
		{"invented mounted exception", func(b *batch) { row(b, 0).TwoHandedException = "while-prone" }},
		{"changed mastery", func(b *batch) { row(b, 0).Mastery = "sap" }},
		{"unlisted weight with value", func(b *batch) { row(b, 0).Weight.Listed = false }},
		{"zero weight denominator", func(b *batch) { row(b, 0).Weight.PoundsDenominator = 0 }},
		{"negative weight", func(b *batch) { row(b, 0).Weight.PoundsNumerator = -1 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems248(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems248(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems248RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/handaxe/variant/standard"`)
	if !bytes.Contains(items248JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"proficient":true,`, `"playable":true,`, `"available":true,`, `"properties":{"attack_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items248JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems248(data); err == nil {
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
			if !bytes.Contains(items248JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items248JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems248(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items248JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems248(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items248JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems247(items248JSON); err == nil {
		t.Fatal("#247 decoder accepted the #248 batch")
	}
	if _, err := DecodeItems248(items247JSON); err == nil {
		t.Fatal("#248 decoder accepted the #247 batch")
	}
}
