package srd521

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// Shared helpers for the adventuring gear batches #256 to #262, and the checks
// that need the whole item catalogue.

const (
	spellScrollGearID = prefix + "adventuring-gear/spell-scroll"
	heavyCrossbowID   = prefix + "weapons/heavy-crossbow"
)

// requirePotionOfHealingOwner checks that the Healing contract still owns the
// adventuring gear entry of p. 99 with its printed values, and returns that
// presentation as an identity for the batch #261 inventory join.
func requirePotionOfHealingOwner(t *testing.T) ItemSourceIdentity {
	t.Helper()
	contract, err := LoadItemCatalogue()
	if err != nil {
		t.Fatal(err)
	}
	family := contract.Families[0]
	gear := family.Presentations[0]
	if gear.ID != prefix+"adventuring-gear/potion-of-healing" || !slices.Equal(gear.Pages, []int{99, 99}) {
		t.Fatalf("the contract no longer owns the gear entry: %+v", gear)
	}
	basic := family.Variants[0]
	if basic.Name != "Potion of Healing" || basic.Cost.Kind != "listed" || basic.Cost.Amount != 50 || basic.Cost.Currency != "gp" || !slices.Equal(basic.Cost.SourcePages, []int{99, 99}) {
		t.Fatalf("the p. 99 price disagrees with the contract: %+v", basic.Cost)
	}
	if *basic.Healing != (ItemHealing{Dice: 2, Sides: 4, Bonus: 2}) || !slices.Equal(basic.Use.Restrictions, []string{"bonus-action", "consume-one", "administer-within-5-feet"}) {
		t.Fatalf("the p. 99 use or healing disagrees with the contract: %+v %v", basic.Healing, basic.Use.Restrictions)
	}
	if basic.Weight == nil || *basic.Weight != (ItemGearWeight{Listed: true, PoundsNumerator: 1, PoundsDenominator: 2, Condition: "none"}) {
		t.Fatalf("the gear table weight of 1/2 lb. disagrees with the contract: %+v", basic.Weight)
	}
	owner := ItemSourceIdentity{FamilyID: family.ID, Name: family.Name, PresentationIDs: []string{gear.ID}, VariantIDs: []string{basic.ID}}
	second := []ItemSourceIdentity{{FamilyID: family.ID, Name: family.Name, PresentationIDs: []string{gear.ID + "-copy"}, VariantIDs: []string{family.ID + "/variant/standard"}}}
	if _, err := ComposeItemSourceIdentities(contract.SourceIdentities(), second); err == nil {
		t.Fatal("a second Potion of Healing record was accepted")
	}
	return owner
}

// requireSpellScrollOwner checks that batch #242 owns Spell Scroll, holds the
// p. 99 gear listing with its two prices and agrees with that entry (cantrip and
// level 1, save DC 13, attack +5, no Material components). It returns the gear
// entry's identity for the batch #262 inventory join.
func requireSpellScrollOwner(t *testing.T) ItemSourceIdentity {
	t.Helper()
	batch := loadItems242(t)
	index := slices.IndexFunc(batch.Families, func(f ItemSourceFamily[ItemDetails242]) bool { return f.Name == "Spell Scroll" })
	if index < 0 {
		t.Fatal("batch #242 no longer owns Spell Scroll")
	}
	family := batch.Families[index]
	if len(family.Presentations) != 2 || family.Presentations[1].ID != spellScrollGearID || !slices.Equal(family.Presentations[1].Pages, []int{99, 99}) {
		t.Fatalf("batch #242 no longer holds the gear entry: %+v", family.Presentations)
	}
	for level, v := range family.Variants[:2] {
		scroll := v.Details.Scroll
		if scroll.SpellLevel != level || scroll.SaveDC != 13 || scroll.AttackBonus != 5 || scroll.MaterialComponentsRequired ||
			scroll.ReadableIf != "spell-on-your-spell-list" || scroll.CastingTime != "spells-normal-casting-time" {
			t.Fatalf("%s disagrees with the p. 99 entry: %+v", v.Name, scroll)
		}
		if v.Cost.Kind != "listed" || v.Cost.Amount != []int{30, 50}[level] || v.Cost.Currency != "gp" || !slices.Equal(v.Cost.SourcePages, []int{99, 99}) {
			t.Fatalf("%s disagrees with the p. 99 price: %+v", v.Name, v.Cost)
		}
	}
	owner := ItemSourceIdentity{FamilyID: family.ID, Name: family.Name, PresentationIDs: []string{spellScrollGearID}, VariantIDs: []string{family.Variants[0].ID}}
	second := []ItemSourceIdentity{{FamilyID: family.ID, Name: family.Name, PresentationIDs: []string{spellScrollGearID}, VariantIDs: []string{family.ID + "/variant/standard"}}}
	if _, err := ComposeItemSourceIdentities(batch.SourceIdentities(), second); err == nil {
		t.Fatal("a second Spell Scroll record was accepted")
	}
	return owner
}

type itemGearKeyWant struct{ key, family, variant string }

// joinCatalogueEquipment joins equipment keys of the character source
// catalogues to the batch's records. The catalogues are not changed.
func joinCatalogueEquipment[D any](t *testing.T, families []ItemSourceFamily[D], joins []itemGearKeyWant) {
	t.Helper()
	paths, err := filepath.Glob("testdata/character-*.json")
	if err != nil || len(paths) == 0 {
		t.Fatalf("character catalogues not found: %v", err)
	}
	var catalogues strings.Builder
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		catalogues.Write(data)
	}
	for _, join := range joins {
		if !strings.Contains(catalogues.String(), `"srd-5.2.1:mechanic/equipment/`+join.key+`"`) {
			t.Fatalf("no character catalogue refers to equipment key %s", join.key)
		}
		id := prefix + "items/" + join.family + "/variant/" + join.variant
		found := slices.ContainsFunc(families, func(f ItemSourceFamily[D]) bool {
			return slices.ContainsFunc(f.Variants, func(v ItemSourceVariant[D]) bool { return v.ID == id })
		})
		if !found {
			t.Fatalf("equipment key %s has no record %s", join.key, id)
		}
	}
}

// requireCreationNames checks that a creation fixture names each gear family.
// The fixtures spell apostrophes in ASCII.
func requireCreationNames[D any](t *testing.T, families []ItemSourceFamily[D], build *SupportedBuild, names map[string]string) {
	t.Helper()
	for name, family := range names {
		if !slices.Contains(build.Equipment, name) {
			t.Fatalf("the %s creation fixture no longer names %s", build.Class, name)
		}
		if !slices.ContainsFunc(families, func(f ItemSourceFamily[D]) bool { return strings.ReplaceAll(f.Name, "’", "'") == family }) {
			t.Fatalf("creation equipment %s has no record %s", name, family)
		}
	}
}

// rejectInvalidGearFixtures mutates a decoded gear batch. Every mutation of the
// identity, requirement, price or payload shape must fail.
func rejectInvalidGearFixtures[D any](t *testing.T, load func() (ItemSourceBatch[D], error), decode func([]byte) (ItemSourceBatch[D], error), gear func(*D) *ItemGear) {
	t.Helper()
	type batch = ItemSourceBatch[D]
	v := func(b *batch) *ItemSourceVariant[D] { return &b.Families[0].Variants[0] }
	row := func(b *batch) *ItemGear { return gear(v(b).Details) }
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-other" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 255 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[1:] }},
		{"reordered families", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"adjacent family", func(b *batch) {
			extra := b.Families[0]
			extra.ID, extra.Name = prefix+"items/other", "Other"
			b.Families = append(b.Families, extra)
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = prefix + "items/unreviewed" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Unreviewed Name" }},
		{"magic kind", func(b *batch) { b.Families[0].Kind = "magic" }},
		{"category", func(b *batch) { b.Families[0].Category = "wondrous-item" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"any-simple-weapon"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{90, 90} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"magic presentation", func(b *batch) { b.Families[0].Presentations[0].ID = prefix + "magic-items/other" }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = b.Families[0].Variants[1:] }},
		{"variant id", func(b *batch) { v(b).ID = "" }},
		{"variant name", func(b *batch) { v(b).Name = "Unreviewed" }},
		{"invented rarity", func(b *batch) { v(b).Rarity = "common" }},
		{"attunement", func(b *batch) { v(b).Attunement.Mode = "required" }},
		{"invented charges", func(b *batch) { v(b).Charges.Maximum = 10 }},
		{"equip mode", func(b *batch) { v(b).Equip.Mode = "wear" }},
		{"use mode", func(b *batch) { v(b).Use.Mode = "drink" }},
		{"invented use restriction", func(b *batch) { v(b).Use.Restrictions = []string{"proficiency"} }},
		{"unlisted price", func(b *batch) { v(b).Cost.Kind = "unspecified" }},
		{"changed price", func(b *batch) { v(b).Cost.Amount++ }},
		{"changed currency", func(b *batch) { v(b).Cost.Currency = "ep" }},
		{"cost page", func(b *batch) { v(b).Cost.SourcePages = []int{90, 90} }},
		{"missing mechanic", func(b *batch) { v(b).RequiredMechanicKeys = v(b).RequiredMechanicKeys[1:] }},
		{"unknown mechanic", func(b *batch) { v(b).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b).Details = nil }},
		{"empty details", func(b *batch) { v(b).Details = new(D) }},
		{"unknown action", func(b *batch) { row(b).Actions = []string{"free-action"} }},
		{"unknown spellcasting class", func(b *batch) { row(b).SpellcastingFocus = []string{"fighter"} }},
		{"prose trait", func(b *batch) { row(b).Traits = append(row(b).Traits, "Do the thing.") }},
		{"invented trait", func(b *batch) { row(b).Traits = append(row(b).Traits, "grants-advantage") }},
		{"invented check", func(b *batch) {
			row(b).Checks = append(row(b).Checks, ItemGearCheck{"other", "ability-check", "strength", "athletics", 10, "fixed"})
		}},
		{"saving throw with a skill", func(b *batch) {
			row(b).Checks = []ItemGearCheck{{"other", "saving-throw", "strength", "athletics", 10, "fixed"}}
		}},
		{"zero DC", func(b *batch) {
			row(b).Checks = []ItemGearCheck{{"other", "saving-throw", "strength", "none", 0, "fixed"}}
		}},
		{"invented damage", func(b *batch) { row(b).Damage = append(row(b).Damage, ItemGearDamage{"other", 1, 4, 0, "fire"}) }},
		{"both damage forms", func(b *batch) { row(b).Damage = []ItemGearDamage{{"other", 1, 4, 1, "fire"}} }},
		{"zero denominator", func(b *batch) { row(b).Measures = []ItemGearMeasure{{"other", 1, 0, "feet"}} }},
		{"negative measure", func(b *batch) { row(b).Measures = []ItemGearMeasure{{"other", -1, 1, "feet"}} }},
		{"invented content", func(b *batch) { row(b).Contents = append(row(b).Contents, ItemGearContent{"backpack", 1, "item"}) }},
		{"zero content quantity", func(b *batch) { row(b).Contents = []ItemGearContent{{"backpack", 0, "item"}} }},
		{"listed weight without a value", func(b *batch) {
			row(b).Weight = ItemGearWeight{Listed: true, PoundsDenominator: 1, Condition: "none"}
		}},
		{"unlisted weight with a value", func(b *batch) {
			row(b).Weight = ItemGearWeight{PoundsNumerator: 1, PoundsDenominator: 1, Condition: "none"}
		}},
		{"zero weight denominator", func(b *batch) { row(b).Weight.PoundsDenominator = 0 }},
		{"unknown weight condition", func(b *batch) { row(b).Weight.Condition = "empty" }},
		{"first source page", func(b *batch) { row(b).SourcePages = []int{93, 100} }},
		{"overlong source pages", func(b *batch) {
			p := row(b)
			p.SourcePages = []int{p.SourcePages[0], p.SourcePages[0] + 3}
		}},
		{"short source pages", func(b *batch) { row(b).SourcePages = row(b).SourcePages[:1] }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := load()
			if err != nil {
				t.Fatal(err)
			}
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := decode(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

// rejectChangedGearValues changes, in turn, every value of every gear payload
// outside weight and source_pages, and drops and repeats the last member of
// every list there. Each decode must fail.
func rejectChangedGearValues(t *testing.T, source []byte, decode func([]byte) error) {
	t.Helper()
	var root map[string]any
	if err := json.Unmarshal(source, &root); err != nil {
		t.Fatal(err)
	}
	changed := 0
	try := func(what string, set func(any), old, replacement any) {
		t.Helper()
		set(replacement)
		data, err := json.Marshal(root)
		set(old)
		if err != nil {
			t.Fatal(err)
		}
		if err := decode(data); err == nil {
			t.Fatalf("changed value accepted: %s", what)
		}
		changed++
	}
	var walk func(path string, node any, set func(any))
	walk = func(path string, node any, set func(any)) {
		switch n := node.(type) {
		case map[string]any:
			for key, child := range n {
				if key != "weight" && key != "source_pages" {
					walk(path+"."+key, child, func(x any) { n[key] = x })
				}
			}
		case []any:
			if len(n) > 0 {
				try(path+" without its last member", set, n, n[:len(n)-1])
				try(path+" with a repeated member", set, n, append(slices.Clone(n), n[len(n)-1]))
			}
			for i, child := range n {
				walk(path, child, func(x any) { n[i] = x })
			}
		case string:
			try(path, set, n, n+"-x")
		case float64:
			try(path, set, n, n+1)
		default:
			t.Fatalf("unexpected payload value at %s: %T", path, node)
		}
	}
	for _, family := range root["families"].([]any) {
		for _, variant := range family.(map[string]any)["variants"].([]any) {
			gear := variant.(map[string]any)["details"].(map[string]any)["gear"].(map[string]any)
			walk("gear", gear, nil)
		}
	}
	if changed == 0 {
		t.Fatal("no payload values found")
	}
	data, err := json.Marshal(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := decode(data); err != nil {
		t.Fatalf("the restored document no longer decodes: %v", err)
	}
}

func gearFamilies(t *testing.T) (slugs map[string]bool, contents []ItemGearContent) {
	t.Helper()
	slugs = map[string]bool{"potion-of-healing": true}
	add := func(id string, gear *ItemGear) {
		slugs[strings.TrimPrefix(id, prefix+"items/")] = true
		contents = append(contents, gear.Contents...)
	}
	collectGear(t, LoadItems256, itemGear256, add)
	collectGear(t, LoadItems257, itemGear257, add)
	collectGear(t, LoadItems258, itemGear258, add)
	collectGear(t, LoadItems259, itemGear259, add)
	collectGear(t, LoadItems260, itemGear260, add)
	collectGear(t, LoadItems261, itemGear261, add)
	collectGear(t, LoadItems262, itemGear262, add)
	return slugs, contents
}

func collectGear[D any](t *testing.T, load func() (ItemSourceBatch[D], error), gear func(*D) *ItemGear, add func(string, *ItemGear)) {
	t.Helper()
	b, err := load()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range b.Families {
		for _, v := range f.Variants {
			add(f.ID, gear(v.Details))
		}
	}
}

// Every item a pack lists is a gear family. The record names it; it grants
// nothing.
func TestItemGearPackContentsNameCatalogueFamilies(t *testing.T) {
	slugs, contents := gearFamilies(t)
	if len(slugs) != 80 || len(contents) != 64 {
		t.Fatalf("expected 80 gear families and 64 pack lines, got %d and %d", len(slugs), len(contents))
	}
	for _, c := range contents {
		if !slugs[c.Item] {
			t.Fatalf("pack content %s has no gear record", c.Item)
		}
	}
}

// TestItemCatalogueMatchesPlanningInventory compares the whole item catalogue,
// the Healing contract and batches #225 to #262, with every item record of the
// planning inventory, in both directions. There are two reviewed differences:
// the inventory omits Heavy Crossbow, and the contract names its family Potion
// of Healing where the p. 236 heading is plural.
func TestItemCatalogueMatchesPlanningInventory(t *testing.T) {
	composed, err := ComposeItemSourceIdentities(itemIdentitiesThrough(t, 262)...)
	if err != nil {
		t.Fatal(err)
	}
	presented := map[string]string{}
	variants := 0
	for _, f := range composed {
		variants += len(f.VariantIDs)
		for _, p := range f.PresentationIDs {
			presented[p] = f.Name
		}
	}
	if len(composed) != 413 || variants != 513 || len(presented) != 415 {
		t.Fatalf("unexpected item catalogue: %d families, %d variants, %d presentations", len(composed), variants, len(presented))
	}
	counts, listed := map[string]int{}, map[string]bool{}
	for _, r := range itemInventory(t) {
		if r.Family == "spells" {
			continue
		}
		counts[r.Family]++
		listed[r.ID] = true
		name, ok := presented[r.ID]
		switch r.ID {
		case prefix + "magic-items/potions-of-healing":
			if name != "Potion of Healing" || r.Name != "Potions of Healing" {
				t.Fatalf("the Potions of Healing join changed: %q %q", name, r.Name)
			}
		default:
			if !ok || name != r.Name {
				t.Fatalf("inventory record %s has no matching catalogue record: %q", r.ID, name)
			}
		}
	}
	want := map[string]int{"magic-items": 258, "weapons": 37, "armor": 13, "tools": 25, "adventuring-gear": 81}
	if !maps.Equal(counts, want) {
		t.Fatalf("inventory item counts changed: %v", counts)
	}
	extra := []string{}
	for id := range presented {
		if !listed[id] {
			extra = append(extra, id)
		}
	}
	if !slices.Equal(extra, []string{heavyCrossbowID}) {
		t.Fatalf("catalogue records outside the inventory: %v", extra)
	}
}
