package srd521

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func loadItems237(t *testing.T) ItemSourceBatch[ItemDetails237] {
	t.Helper()
	b, err := LoadItems237()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 235-236, in planning inventory order for issue #237. Potion of
// Giant Strength prints "Rarity Varies": each row of its table is a variant,
// with the rarity printed in that row. Every potion lists no price.
func TestItems237SourceIdentities(t *testing.T) {
	type variant struct{ slug, name, rarity string }
	expected := []struct {
		name, slug string
		page       int
		variants   []variant
	}{
		{"Potion of Clairvoyance", "potion-of-clairvoyance", 235, []variant{{"standard", "Potion of Clairvoyance", "rare"}}},
		{"Potion of Climbing", "potion-of-climbing", 235, []variant{{"standard", "Potion of Climbing", "common"}}},
		{"Potion of Diminution", "potion-of-diminution", 235, []variant{{"standard", "Potion of Diminution", "rare"}}},
		{"Potion of Flying", "potion-of-flying", 235, []variant{{"standard", "Potion of Flying", "very-rare"}}},
		{"Potion of Gaseous Form", "potion-of-gaseous-form", 235, []variant{{"standard", "Potion of Gaseous Form", "rare"}}},
		{"Potion of Giant Strength", "potion-of-giant-strength", 235, []variant{{"hill", "Potion of Giant Strength (hill)", "uncommon"}, {"frost-or-stone", "Potion of Giant Strength (frost or stone)", "rare"}, {"fire", "Potion of Giant Strength (fire)", "rare"}, {"cloud", "Potion of Giant Strength (cloud)", "very-rare"}, {"storm", "Potion of Giant Strength (storm)", "legendary"}}},
		{"Potion of Growth", "potion-of-growth", 235, []variant{{"standard", "Potion of Growth", "uncommon"}}},
		{"Potion of Heroism", "potion-of-heroism", 236, []variant{{"standard", "Potion of Heroism", "rare"}}},
		{"Potion of Invisibility", "potion-of-invisibility", 236, []variant{{"standard", "Potion of Invisibility", "rare"}}},
		{"Potion of Invulnerability", "potion-of-invulnerability", 236, []variant{{"standard", "Potion of Invulnerability", "rare"}}},
		{"Potion of Longevity", "potion-of-longevity", 236, []variant{{"standard", "Potion of Longevity", "very-rare"}}},
		{"Potion of Mind Reading", "potion-of-mind-reading", 236, []variant{{"standard", "Potion of Mind Reading", "rare"}}},
	}
	b := loadItems237(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-13" || b.CatalogueIssue != 237 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #237")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 235–236") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	total := 0
	for i, want := range expected {
		t.Run(want.name, func(t *testing.T) {
			f := b.Families[i]
			id := "srd-5.2.1:content/items/" + want.slug
			if f.ID != id || f.Name != want.name || f.Kind != "magic" || f.Category != "potion" || len(f.Base.Options) != 0 || len(f.Base.Excluded) != 0 {
				t.Fatalf("source identity mismatch: %+v", f)
			}
			if len(f.Presentations) != 1 || f.Presentations[0].ID != "srd-5.2.1:content/magic-items/"+want.slug || !slices.Equal(f.Presentations[0].Pages, []int{want.page, want.page}) {
				t.Fatal("source presentation mismatch")
			}
			if len(f.Variants) != len(want.variants) {
				t.Fatalf("the source lists %d forms of %s, got %d variants", len(want.variants), want.name, len(f.Variants))
			}
			for j, wv := range want.variants {
				v := f.Variants[j]
				if v.ID != id+"/variant/"+wv.slug || v.Name != wv.name || v.Rarity != wv.rarity || v.Attunement.Mode != "none" || len(v.Attunement.Restrictions) != 0 {
					t.Fatalf("variant mismatch: %+v", v)
				}
				if v.Cost.Kind != "unspecified" || v.Cost.Amount != 0 || v.Cost.Currency != "none" || !slices.Equal(v.Cost.SourcePages, []int{want.page, want.page}) {
					t.Fatal("no price is listed; cost must stay unspecified")
				}
				if *v.Charges != (ItemCharges{Mode: "none", Recharge: "none"}) {
					t.Fatalf("no potion lists charges: %+v", *v.Charges)
				}
				if v.MechanicsAvailable() {
					t.Fatal("catalogue record became available")
				}
				total++
			}
		})
	}
	if total != 16 {
		t.Fatalf("expected 16 source-listed variants, got %d", total)
	}
	// Join every family to the decision inventory the issue cites, in its order.
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
		Batches []struct {
			Key     string `json:"key"`
			Pages   string `json:"pages"`
			Members []struct {
				Name string `json:"name"`
				Page int    `json:"source_page"`
			} `json:"members"`
		} `json:"batches"`
	}
	if err := json.Unmarshal(data, &inventory); err != nil {
		t.Fatal(err)
	}
	index := 0
	for _, r := range inventory.Records {
		if r.Issue != 237 {
			continue
		}
		if index >= len(b.Families) {
			t.Fatal("inventory lists more families than the batch")
		}
		f := b.Families[index]
		if r.Family != "magic-items" || r.Batch != b.Batch || r.Status != "unsupported" || r.ID != f.Presentations[0].ID || r.Name != f.Name || r.Page != f.Presentations[0].Pages[0] {
			t.Fatalf("inventory record %s disagrees with batch family %s", r.ID, f.Name)
		}
		index++
	}
	if index != 12 {
		t.Fatalf("inventory assigns %d families to issue #237", index)
	}
	// The planning batch names the same members and pages, in the same order.
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "235, 236" || len(batch.Members) != 12 {
			t.Fatalf("planning batch pages or member count changed: %q, %d", batch.Pages, len(batch.Members))
		}
		for i, m := range batch.Members {
			if m.Name != expected[i].name || m.Page != expected[i].page {
				t.Fatalf("planning batch member %d is %s p.%d", i, m.Name, m.Page)
			}
		}
	}
	if planned != 1 {
		t.Fatalf("expected one planning batch %s, got %d", b.Batch, planned)
	}
	b.Families[0].Variants[0].Details.Clairvoyance.Spell = "changed"
	if fresh := loadItems237(t); fresh.Families[0].Variants[0].Details.Clairvoyance.Spell != "clairvoyance" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all 16 variants, transcribed from pinned PDF pp.
// 235-236. The five Potion of Giant Strength variants are the five rows of its
// table; the "frost or stone" row is one variant listing two giants.
func TestItems237SourceValues(t *testing.T) {
	expected := []string{
		`{"clairvoyance":{"spell":"clairvoyance","concentration_required":false,"appearance":"yellowish-liquid-with-bobbing-eyeball-that-vanishes-when-opened"}}`,
		`{"climbing":{"speed_kind":"climb","speed_equals":"your-speed","duration_hours":1,"check_ability":"strength","check_skill":"athletics","check_mode":"advantage","check_subject":"climbing","appearance":"brown-silver-gray-layers-like-bands-of-stone-that-shaking-fails-to-mix"}}`,
		`{"diminution":{"spell":"enlarge-reduce","spell_effect":"reduce","duration_dice":1,"duration_die_sides":4,"duration_unit":"hours","concentration_required":false,"appearance":"red-contracts-to-tiny-bead-then-expands-to-color-clear-liquid-shaking-fails-to-interrupt"}}`,
		`{"flying":{"speed_kind":"fly","speed_equals":"your-speed","duration_hours":1,"can_hover":true,"falls_when_it_ends":"when-in-the-air","fall_exception":"other-means-of-staying-aloft","appearance":"clear-liquid-floats-at-top-with-cloudy-white-impurities-drifting"}}`,
		`{"gaseous_form":{"spell":"gaseous-form","duration_hours":1,"ends_early_action":"bonus","concentration_required":false,"appearance":"container-seems-to-hold-fog-that-moves-and-pours-like-water"}}`,
		`{"giant_strength":{"giants":["hill"],"strength_score":21,"duration_hours":1,"score_change":"strength-score-becomes-giant-score","no_effect_at_or_above_score":21,"appearance":"transparent-liquid-with-floating-sliver-of-light-resembling-giants-fingernail"}}`,
		`{"giant_strength":{"giants":["frost","stone"],"strength_score":23,"duration_hours":1,"score_change":"strength-score-becomes-giant-score","no_effect_at_or_above_score":23,"appearance":"transparent-liquid-with-floating-sliver-of-light-resembling-giants-fingernail"}}`,
		`{"giant_strength":{"giants":["fire"],"strength_score":25,"duration_hours":1,"score_change":"strength-score-becomes-giant-score","no_effect_at_or_above_score":25,"appearance":"transparent-liquid-with-floating-sliver-of-light-resembling-giants-fingernail"}}`,
		`{"giant_strength":{"giants":["cloud"],"strength_score":27,"duration_hours":1,"score_change":"strength-score-becomes-giant-score","no_effect_at_or_above_score":27,"appearance":"transparent-liquid-with-floating-sliver-of-light-resembling-giants-fingernail"}}`,
		`{"giant_strength":{"giants":["storm"],"strength_score":29,"duration_hours":1,"score_change":"strength-score-becomes-giant-score","no_effect_at_or_above_score":29,"appearance":"transparent-liquid-with-floating-sliver-of-light-resembling-giants-fingernail"}}`,
		`{"growth":{"spell":"enlarge-reduce","spell_effect":"enlarge","duration_minutes":10,"concentration_required":false,"appearance":"red-expands-from-tiny-bead-to-color-clear-liquid-then-contracts-shaking-fails-to-interrupt"}}`,
		`{"heroism":{"temporary_hit_points":10,"duration_hours":1,"spell":"bless","spell_duration":"same-as-temporary-hit-points","concentration_required":false,"appearance":"blue-liquid-bubbles-and-steams-as-if-boiling"}}`,
		`{"invisibility":{"condition":"invisible","duration_hours":1,"ends_early_on":["make-attack-roll","deal-damage","cast-spell"],"appearance":"container-looks-empty-but-feels-as-though-it-holds-liquid"}}`,
		`{"invulnerability":{"resistance":"all-damage","duration_minutes":1,"appearance":"syrupy-liquid-looks-like-liquefied-iron"}}`,
		`{"longevity":{"age_reduction_years":{"dice":1,"sides":6,"bonus":6},"minimum_age_years":13,"aging_chance_percent":10,"aging_chance_cumulative":true,"aging_chance_trigger":"each-subsequent-potion-of-longevity","aging_instead_years":{"dice":1,"sides":6,"bonus":6},"appearance":"amber-liquid-with-suspended-tiny-beating-heart-that-vanishes-when-opened"}}`,
		`{"mind_reading":{"spell":"detect-thoughts","save_dc":13,"duration_minutes":10,"concentration_required":false,"appearance":"dense-purple-liquid-with-floating-ovoid-pink-cloud"}}`,
	}
	// use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"consume:drink-one:item.consume spell.effect-clairvoyance spell.no-concentration",
		"consume:drink-one:item.consume speed.climb-equal-to-speed duration.1-hour check.strength-athletics-advantage-climbing",
		"consume:drink-one:item.consume spell.effect-enlarge-reduce-reduce duration.1d4-hours spell.no-concentration",
		"consume:drink-one:item.consume speed.fly-equal-to-speed duration.1-hour movement.hover movement.fall-when-ends-in-air",
		"consume:drink-one:item.consume spell.effect-gaseous-form duration.1-hour action.bonus-end-effect spell.no-concentration",
		"consume:drink-one:item.consume ability.strength-score-set-by-giant-type duration.1-hour effect.none-at-or-above-score",
		"consume:drink-one:item.consume ability.strength-score-set-by-giant-type duration.1-hour effect.none-at-or-above-score",
		"consume:drink-one:item.consume ability.strength-score-set-by-giant-type duration.1-hour effect.none-at-or-above-score",
		"consume:drink-one:item.consume ability.strength-score-set-by-giant-type duration.1-hour effect.none-at-or-above-score",
		"consume:drink-one:item.consume ability.strength-score-set-by-giant-type duration.1-hour effect.none-at-or-above-score",
		"consume:drink-one:item.consume spell.effect-enlarge-reduce-enlarge duration.10-minutes spell.no-concentration",
		"consume:drink-one:item.consume hit-points.temporary-10 duration.1-hour spell.effect-bless spell.no-concentration",
		"consume:drink-one:item.consume condition.invisible duration.1-hour effect.ends-early-on-attack-damage-or-spell",
		"consume:drink-one:item.consume defense.resistance-all-damage duration.1-minute",
		"consume:drink-one:item.consume age.reduce-1d6-plus-6-years age.minimum-13-years age.cumulative-10-percent-chance-to-age-instead",
		"consume:drink-one:item.consume spell.effect-detect-thoughts spell.detect-thoughts-save-dc-13 duration.10-minutes spell.no-concentration",
	}
	index := 0
	for _, f := range loadItems237(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails237
				d := json.NewDecoder(strings.NewReader(expected[index]))
				d.DisallowUnknownFields()
				if err := d.Decode(&want); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(v.Details, &want) {
					t.Fatalf("PDF p.%d values mismatch: %+v", f.Presentations[0].Pages[0], v.Details)
				}
				got := fmt.Sprintf("%s:%s:%s", v.Use.Mode, strings.Join(v.Use.Restrictions, " "), strings.Join(v.RequiredMechanicKeys, " "))
				if got != requirements[index] || v.Equip.Mode != "none" || len(v.Equip.Restrictions) != 0 {
					t.Fatalf("source requirements mismatch: %s", got)
				}
			})
			index++
		}
	}
	if index != len(expected) || index != len(requirements) {
		t.Fatal("fixture count disagrees with the batch")
	}
}

func TestItems237CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems237(t)
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
	if len(composed) != 97 || variants != 139 || composed[0].Name != "Potion of Healing" || composed[85].Name != "Potion of Clairvoyance" || composed[96].Name != "Potion of Mind Reading" {
		t.Fatalf("expected Potion of Healing, batches #225 to #231 and exactly the twelve families of #237, got %d families and %d variants", len(composed), variants)
	}
	if _, err := ComposeItemSourceIdentities(batch.SourceIdentities(), batch.SourceIdentities()); err == nil {
		t.Fatal("duplicate batch accepted")
	}
	clash := batch.SourceIdentities()[:1]
	clash[0].FamilyID, clash[0].Name, clash[0].PresentationIDs = "srd-5.2.1:content/items/other", "Other", []string{"srd-5.2.1:content/magic-items/other"}
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
			if strings.Contains(fixture, "items_237") {
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
func TestItems237RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items237JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 500 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items237JSON, &fresh); err != nil {
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
		if _, err := DecodeItems237(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems237RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails237]
	v := func(b *batch, family, variant int) *ItemSourceVariant[ItemDetails237] {
		return &b.Families[family].Variants[variant]
	}
	const (
		clairvoyance, climbing, diminution, flying, gaseous, giant, growth, heroism, invisibility, invuln, longevity, mind = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11
	)
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-12" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 236 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/potion-of-poison"
			b.Families[12].Name = "Potion of Poison"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/potion-of-speed" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Potion of Clairaudience" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "wondrous-item" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[climbing].Base.Options = []string{"flask"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{236, 236} }},
		{"heroism page", func(b *batch) { b.Families[heroism].Presentations[0].Pages = []int{235, 235} }},
		{"page range", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{235, 236} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"missing giant variant", func(b *batch) { b.Families[giant].Variants = b.Families[giant].Variants[:4] }},
		{"merged giant variants", func(b *batch) { b.Families[giant].Variants = b.Families[giant].Variants[:1] }},
		{"split frost and stone", func(b *batch) {
			f := &b.Families[giant]
			f.Variants = append(f.Variants[:1], append([]ItemSourceVariant[ItemDetails237]{f.Variants[1], f.Variants[1]}, f.Variants[2:]...)...)
		}},
		{"reordered giant variants", func(b *batch) {
			s := b.Families[giant].Variants
			s[0], s[1] = s[1], s[0]
		}},
		{"variant id", func(b *batch) { v(b, 0, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0, 0).Name = "Clairvoyance" }},
		{"giant name", func(b *batch) { v(b, giant, 1).Name = "Potion of Giant Strength (frost)" }},
		{"rarity", func(b *batch) { v(b, 0, 0).Rarity = "common" }},
		{"climbing rarity", func(b *batch) { v(b, climbing, 0).Rarity = "uncommon" }},
		{"giant rarity", func(b *batch) { v(b, giant, 4).Rarity = "very-rare" }},
		{"unknown rarity", func(b *batch) { v(b, 0, 0).Rarity = "mythic" }},
		{"invented attunement", func(b *batch) { v(b, 0, 0).Attunement.Mode = "required" }},
		{"attunement restriction", func(b *batch) { v(b, 0, 0).Attunement.Restrictions = []string{"spellcaster"} }},
		{"missing charges", func(b *batch) { v(b, 0, 0).Charges = nil }},
		{"invented charges", func(b *batch) { v(b, 0, 0).Charges.Maximum = 1 }},
		{"equip mode", func(b *batch) { v(b, 0, 0).Equip.Mode = "wear" }},
		{"equip restriction", func(b *batch) { v(b, 0, 0).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, 0, 0).Use.Mode = "magic-action" }},
		{"dropped use restriction", func(b *batch) { v(b, 0, 0).Use.Restrictions = []string{} }},
		{"listed price", func(b *batch) { v(b, 0, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0, 0).Cost.Amount = 500 }},
		{"cost currency", func(b *batch) { v(b, 0, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = v(b, 0, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0, 0).Details = &ItemDetails237{} }},
		{"mismatched payload", func(b *batch) { v(b, 0, 0).Details = v(b, 1, 0).Details }},
		{"second payload", func(b *batch) { v(b, 0, 0).Details.Climbing = v(b, 1, 0).Details.Climbing }},
		{"blank key", func(b *batch) { v(b, 0, 0).Details.Clairvoyance.Spell = " " }},
		{"prose value", func(b *batch) { v(b, 0, 0).Details.Clairvoyance.Spell = "Clairvoyance spell" }},
		{"negative value", func(b *batch) { v(b, climbing, 0).Details.Climbing.DurationHours = -1 }},
		{"clairvoyance concentration", func(b *batch) { v(b, clairvoyance, 0).Details.Clairvoyance.ConcentrationRequired = true }},
		{"climbing duration", func(b *batch) { v(b, climbing, 0).Details.Climbing.DurationHours = 2 }},
		{"diminution dice", func(b *batch) { v(b, diminution, 0).Details.Diminution.DurationDice = 2 }},
		{"diminution die", func(b *batch) { v(b, diminution, 0).Details.Diminution.DurationDieSides = 6 }},
		{"diminution concentration", func(b *batch) { v(b, diminution, 0).Details.Diminution.ConcentrationRequired = true }},
		{"flying hover", func(b *batch) { v(b, flying, 0).Details.Flying.CanHover = false }},
		{"flying duration", func(b *batch) { v(b, flying, 0).Details.Flying.DurationHours = 0 }},
		{"gaseous duration", func(b *batch) { v(b, gaseous, 0).Details.GaseousForm.DurationHours = 2 }},
		{"gaseous concentration", func(b *batch) { v(b, gaseous, 0).Details.GaseousForm.ConcentrationRequired = true }},
		{"giant score", func(b *batch) { v(b, giant, 0).Details.GiantStrength.StrengthScore = 22 }},
		{"giant threshold", func(b *batch) { v(b, giant, 2).Details.GiantStrength.NoEffectAtOrAboveScore = 23 }},
		{"giant duration", func(b *batch) { v(b, giant, 0).Details.GiantStrength.DurationHours = 2 }},
		{"giant frost only", func(b *batch) { v(b, giant, 1).Details.GiantStrength.Giants = []string{"frost"} }},
		{"giant no giants", func(b *batch) { v(b, giant, 0).Details.GiantStrength.Giants = []string{} }},
		{"giant row swapped", func(b *batch) { v(b, giant, 3).Details.GiantStrength = v(b, giant, 4).Details.GiantStrength }},
		{"growth duration", func(b *batch) { v(b, growth, 0).Details.Growth.DurationMinutes = 60 }},
		{"growth concentration", func(b *batch) { v(b, growth, 0).Details.Growth.ConcentrationRequired = true }},
		{"heroism hit points", func(b *batch) { v(b, heroism, 0).Details.Heroism.TemporaryHitPoints = 5 }},
		{"heroism duration", func(b *batch) { v(b, heroism, 0).Details.Heroism.DurationHours = 0 }},
		{"heroism concentration", func(b *batch) { v(b, heroism, 0).Details.Heroism.ConcentrationRequired = true }},
		{"invisibility duration", func(b *batch) { v(b, invisibility, 0).Details.Invisibility.DurationHours = 2 }},
		{"invisibility no ends", func(b *batch) { v(b, invisibility, 0).Details.Invisibility.EndsEarlyOn = []string{} }},
		{"invisibility duplicate end", func(b *batch) { v(b, invisibility, 0).Details.Invisibility.EndsEarlyOn[1] = "make-attack-roll" }},
		{"invulnerability duration", func(b *batch) { v(b, invuln, 0).Details.Invuln.DurationMinutes = 10 }},
		{"longevity reduction", func(b *batch) { v(b, longevity, 0).Details.Longevity.AgeReduction.Bonus = 0 }},
		{"longevity aging", func(b *batch) { v(b, longevity, 0).Details.Longevity.AgingInstead.Sides = 8 }},
		{"longevity minimum", func(b *batch) { v(b, longevity, 0).Details.Longevity.MinimumAgeYears = 18 }},
		{"longevity chance", func(b *batch) { v(b, longevity, 0).Details.Longevity.AgingChancePercent = 20 }},
		{"longevity not cumulative", func(b *batch) { v(b, longevity, 0).Details.Longevity.AgingChanceCumulative = false }},
		{"string clairvoyance spell", func(b *batch) { v(b, clairvoyance, 0).Details.Clairvoyance.Spell = "scrying" }},
		{"string climbing speed kind", func(b *batch) { v(b, climbing, 0).Details.Climbing.SpeedKind = "fly" }},
		{"string climbing speed equals", func(b *batch) { v(b, climbing, 0).Details.Climbing.SpeedEquals = "30-feet" }},
		{"string climbing skill", func(b *batch) { v(b, climbing, 0).Details.Climbing.CheckSkill = "acrobatics" }},
		{"string climbing ability", func(b *batch) { v(b, climbing, 0).Details.Climbing.CheckAbility = "dexterity" }},
		{"string climbing mode", func(b *batch) { v(b, climbing, 0).Details.Climbing.CheckMode = "proficiency" }},
		{"string diminution spell", func(b *batch) { v(b, diminution, 0).Details.Diminution.Spell = "haste" }},
		{"string diminution effect", func(b *batch) { v(b, diminution, 0).Details.Diminution.SpellEffect = "enlarge" }},
		{"string diminution unit", func(b *batch) { v(b, diminution, 0).Details.Diminution.DurationUnit = "minutes" }},
		{"string flying speed kind", func(b *batch) { v(b, flying, 0).Details.Flying.SpeedKind = "climb" }},
		{"string flying falls", func(b *batch) { v(b, flying, 0).Details.Flying.FallsWhenItEnds = "never" }},
		{"string flying exception", func(b *batch) { v(b, flying, 0).Details.Flying.FallException = "none" }},
		{"string gaseous spell", func(b *batch) { v(b, gaseous, 0).Details.GaseousForm.Spell = "fog-cloud" }},
		{"string gaseous ends early", func(b *batch) { v(b, gaseous, 0).Details.GaseousForm.EndsEarlyAction = "reaction" }},
		{"string giant score change", func(b *batch) { v(b, giant, 0).Details.GiantStrength.ScoreChange = "add-bonus" }},
		{"string growth spell", func(b *batch) { v(b, growth, 0).Details.Growth.Spell = "haste" }},
		{"string growth effect", func(b *batch) { v(b, growth, 0).Details.Growth.SpellEffect = "reduce" }},
		{"string heroism spell", func(b *batch) { v(b, heroism, 0).Details.Heroism.Spell = "haste" }},
		{"string heroism spell duration", func(b *batch) { v(b, heroism, 0).Details.Heroism.SpellDuration = "1-minute" }},
		{"string invisibility condition", func(b *batch) { v(b, invisibility, 0).Details.Invisibility.Condition = "poisoned" }},
		{"string invulnerability resistance", func(b *batch) { v(b, invuln, 0).Details.Invuln.Resistance = "fire" }},
		{"string longevity trigger", func(b *batch) { v(b, longevity, 0).Details.Longevity.AgingChanceTrigger = "each-potion" }},
		{"string mind spell", func(b *batch) { v(b, mind, 0).Details.MindReading.Spell = "clairvoyance" }},
		{"invisibility ends early", func(b *batch) { v(b, invisibility, 0).Details.Invisibility.EndsEarlyOn = []string{"deal-damage"} }},
		{"mind save", func(b *batch) { v(b, mind, 0).Details.MindReading.SaveDC = 15 }},
		{"mind duration", func(b *batch) { v(b, mind, 0).Details.MindReading.DurationMinutes = 1 }},
		{"mind concentration", func(b *batch) { v(b, mind, 0).Details.MindReading.ConcentrationRequired = true }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems237(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems237(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems237RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/potion-of-clairvoyance/variant/standard"`)
	if !bytes.Contains(items237JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["drink"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items237JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems237(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"clairvoyance": {`, `"clairvoyance_plus": {`},
		"unknown detail":     {`"spell": "clairvoyance"`, `"spell_all": "clairvoyance"`},
		"string for number":  {`"save_dc": 13`, `"save_dc": "13"`},
		"fractional number":  {`"save_dc": 13`, `"save_dc": 13.5`},
		"number for boolean": {`"can_hover": true`, `"can_hover": 1`},
		"null giant list":    {`"giants": [`, `"giants": null, "x": [`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items237JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items237JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems237(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items237JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems237(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor another batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items237JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems231(items237JSON); err == nil {
		t.Fatal("#231 decoder accepted the #237 batch")
	}
	if _, err := DecodeItems237(items231JSON); err == nil {
		t.Fatal("#237 decoder accepted the #231 batch")
	}
}
