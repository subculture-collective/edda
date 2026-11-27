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

func loadItems246(t *testing.T) ItemSourceBatch[ItemDetails246] {
	t.Helper()
	b, err := LoadItems246()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF p. 253, in planning inventory order for issue #246.
func TestItems246SourceIdentities(t *testing.T) {
	type variant struct{ slug, name, rarity string }
	noCharges := ItemCharges{Mode: "none", Recharge: "none"}
	expected := []struct {
		name, slug string
		page       int
		category   string
		base       []string
		attunement string
		charges    ItemCharges
		variants   []variant
	}{
		{"Weapon of Warning", "weapon-of-warning", 253, "weapon", []string{"any-simple-weapon", "any-martial-weapon"}, "required", noCharges, []variant{{"standard", "Weapon of Warning", "uncommon"}}},
		{"Weapon, +1, +2, or +3", "weapon-1-2-or-3", 253, "weapon", []string{"any-simple-weapon", "any-martial-weapon"}, "none", noCharges, []variant{{"plus-1", "Weapon, +1", "uncommon"}, {"plus-2", "Weapon, +2", "rare"}, {"plus-3", "Weapon, +3", "very-rare"}}},
		{"Well of Many Worlds", "well-of-many-worlds", 253, "wondrous-item", nil, "none", noCharges, []variant{{"standard", "Well of Many Worlds", "legendary"}}},
		{"Wind Fan", "wind-fan", 253, "wondrous-item", nil, "none", noCharges, []variant{{"standard", "Wind Fan", "uncommon"}}},
		{"Winged Boots", "winged-boots", 253, "wondrous-item", nil, "required", ItemCharges{Mode: "limited", Maximum: 4, PerUse: 1, Recharge: "regain-1d4-at-dawn"}, []variant{{"standard", "Winged Boots", "uncommon"}}},
		{"Wings of Flying", "wings-of-flying", 253, "wondrous-item", nil, "required", noCharges, []variant{{"standard", "Wings of Flying", "rare"}}},
	}
	b := loadItems246(t)
	if len(b.Families) != 6 || b.Batch != "data-magic-items-22" || b.CatalogueIssue != 246 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the six families of issue #246")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "p. 253") {
		t.Fatal("pin, license, attribution or adaptation notice changed")
	}
	total := 0
	for i, want := range expected {
		t.Run(want.name, func(t *testing.T) {
			f := b.Families[i]
			id := "srd-5.2.1:content/items/" + want.slug
			if f.ID != id || f.Name != want.name || f.Kind != "magic" || f.Category != want.category || len(f.Base.Excluded) != 0 || !slices.Equal(f.Base.Options, append([]string{}, want.base...)) {
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
				if v.ID != id+"/variant/"+wv.slug || v.Name != wv.name || v.Rarity != wv.rarity || v.Attunement.Mode != want.attunement || len(v.Attunement.Restrictions) != 0 {
					t.Fatalf("variant mismatch: %+v", v)
				}
				if v.Cost.Kind != "unspecified" || v.Cost.Amount != 0 || v.Cost.Currency != "none" || !slices.Equal(v.Cost.SourcePages, []int{want.page, want.page}) {
					t.Fatal("no price is listed; cost must stay unspecified")
				}
				if *v.Charges != want.charges {
					t.Fatalf("charges mismatch: %+v", *v.Charges)
				}
				if v.MechanicsAvailable() {
					t.Fatal("catalogue record became available")
				}
				total++
			}
		})
	}
	if total != 8 {
		t.Fatalf("expected 8 source-listed variants, got %d", total)
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
	type member struct {
		name string
		page int
	}
	var members []member
	index := 0
	for _, r := range inventory.Records {
		if r.Issue != 246 {
			continue
		}
		if r.Family != "magic-items" || r.Batch != b.Batch {
			t.Fatalf("inventory record %s changed family or batch", r.ID)
		}
		members = append(members, member{r.Name, r.Page})
		if index >= len(b.Families) {
			t.Fatal("inventory lists more families than the batch")
		}
		f := b.Families[index]
		if r.Status != "unsupported" || r.ID != f.Presentations[0].ID || r.Name != f.Name || r.Page != f.Presentations[0].Pages[0] {
			t.Fatalf("inventory record %s disagrees with batch family %s", r.ID, f.Name)
		}
		index++
	}
	if index != 6 || len(members) != 6 {
		t.Fatalf("inventory assigns %d families to issue #246", len(members))
	}
	// The planning batch names the same members and pages, in the same order.
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "253" || len(batch.Members) != 6 {
			t.Fatalf("planning batch pages or member count changed: %q, %d", batch.Pages, len(batch.Members))
		}
		for i, m := range batch.Members {
			if m.Name != members[i].name || m.Page != members[i].page {
				t.Fatalf("planning batch member %d is %s p.%d", i, m.Name, m.Page)
			}
		}
	}
	if planned != 1 {
		t.Fatalf("expected one planning batch %s, got %d", b.Batch, planned)
	}
	b.Families[0].Variants[0].Details.Warning.Condition = "changed"
	if fresh := loadItems246(t); fresh.Families[0].Variants[0].Details.Warning.Condition != "within-reach-while-attuned" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all 8 variants, transcribed from pinned PDF
// p. 253.
func TestItems246SourceValues(t *testing.T) {
	expected := []string{
		`{"warning":{"condition":"within-reach-while-attuned","radius_feet":30,"subjects":"you-and-allies","alarm":{"wakes":"naturally-sleeping-subjects","when":"combat-begins","wakes_from_magical_sleep":false},"supernatural_readiness":{"roll":"initiative","mode":"advantage"}}}`,
		`{"weapon_bonus":{"attack_bonus":1,"damage_bonus":1,"bonus_rule":"determined-by-rarity"}}`,
		`{"weapon_bonus":{"attack_bonus":2,"damage_bonus":2,"bonus_rule":"determined-by-rarity"}}`,
		`{"weapon_bonus":{"attack_bonus":3,"damage_bonus":3,"bonus_rule":"determined-by-rarity"}}`,
		`{"well":{"material":"fine-black-cloth","folded_size":"handkerchief","unfolded_diameter_feet":6,"open":{"action":"magic","placement":"solid-surface"},"portal":{"direction":"two-way","diameter_feet":6,"destination":"another-world-or-plane","chosen_by":"gm"},"close":{"action":"magic","within_feet":5,"method":"fold-cloth-by-edges"},"reuse_after":{"dice":1,"sides":8,"type":"hours"}}}`,
		`{"wind_fan":{"condition":"while-holding","spell":"gust-of-wind","save_dc":13,"repeat_use":{"failure_percent_cumulative":20,"applies_to":"each-subsequent-use-before-next-dawn","on_failure":"tears-into-nonmagical-tatters"}}}`,
		`{"winged_boots":{"max_charges":4,"recharge_dice":{"dice":1,"sides":4,"type":"charges"},"action":"magic","charge_cost":1,"fly_speed_feet":30,"duration_hours":1,"on_expiry":{"condition":"flying-when-duration-expires","descent_feet_per_round":30}}}`,
		`{"wings":{"form":"cloak-becomes-wings-on-back","action":"magic","duration_hours":1,"early_end_action":"magic","fly_speed_feet":60,"aloft_when_wings_disappear":"fall","reuse_after":{"dice":1,"sides":12,"type":"hours"}}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"within-reach:none::item.attunement item.magic-weapon effect.wake-natural-sleepers-at-combat-start roll.initiative-advantage aura.allies-within-30-feet",
		"wield:weapon-attack::item.magic-weapon bonus.attack-and-damage-by-rarity",
		"wield:weapon-attack::item.magic-weapon bonus.attack-and-damage-by-rarity",
		"wield:weapon-attack::item.magic-weapon bonus.attack-and-damage-by-rarity",
		"none:magic-action:unfold-and-place-on-solid-surface:action.magic-unfold effect.two-way-planar-portal item.gm-chooses-destination action.magic-close-portal item.1d8-hour-reuse",
		"hold:cast-spell:cumulative-failure-chance-before-next-dawn:item.hold spell.cast-gust-of-wind spell.save-dc-13 item.cumulative-20-percent-failure item.destroyed-on-failure",
		"wear:magic-action:expend-one-charge:item.attunement item.wear item.charges action.magic movement.fly-speed-30 duration.1-hour movement.gentle-descent-30-feet-per-round",
		"wear:magic-action:1d12-hour-wait-after-wings-disappear:item.attunement item.wear action.magic movement.fly-speed-60 duration.1-hour action.magic-end-early item.1d12-hour-reuse",
	}
	index := 0
	for _, f := range loadItems246(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails246
				d := json.NewDecoder(strings.NewReader(expected[index]))
				d.DisallowUnknownFields()
				if err := d.Decode(&want); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(v.Details, &want) {
					t.Fatalf("PDF p.%d values mismatch: %+v", f.Presentations[0].Pages[0], v.Details)
				}
				got := fmt.Sprintf("%s:%s:%s:%s", v.Equip.Mode, v.Use.Mode, strings.Join(v.Use.Restrictions, " "), strings.Join(v.RequiredMechanicKeys, " "))
				if got != requirements[index] || len(v.Equip.Restrictions) != 0 {
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

func TestItems246CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems246(t)
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
	if len(composed) != 258 || variants != 334 || composed[0].Name != "Potion of Healing" || composed[252].Name != "Weapon of Warning" || composed[257].Name != "Wings of Flying" {
		t.Fatalf("expected Potion of Healing, batches #225 to #245 and the six families of #246: %d families, %d variants", len(composed), variants)
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
			if strings.Contains(fixture, "items_246") {
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
func TestItems246RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items246JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 300 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items246JSON, &fresh); err != nil {
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
		if _, err := DecodeItems246(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems246RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails246]
	// v addresses one variant of a family by index; see the source identity test.
	v := func(b *batch, family, variant int) *ItemSourceVariant[ItemDetails246] {
		return &b.Families[family].Variants[variant]
	}
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-21" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 245 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:5] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[6].ID = "srd-5.2.1:content/items/adamantine-armor"
			b.Families[6].Name = "Adamantine Armor"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/unreviewed" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Unreviewed Name" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "armor" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"any-armor"} }},
		{"weapon-of-warning widened base", func(b *batch) { b.Families[0].Base.Options = append(b.Families[0].Base.Options, "any-weapon") }},
		{"weapon-1-2-or-3 widened base", func(b *batch) { b.Families[1].Base.Options = append(b.Families[1].Base.Options, "any-weapon") }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{254, 254} }},
		{"page range", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{253, 254} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"variant id", func(b *batch) { v(b, 0, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0, 0).Name = "Unreviewed" }},
		{"weapon-of-warning rarity", func(b *batch) { v(b, 0, 0).Rarity = "common" }},
		{"weapon-of-warning attunement", func(b *batch) { v(b, 0, 0).Attunement.Mode = "none" }},
		{"weapon-1-2-or-3 plus-1 rarity", func(b *batch) { v(b, 1, 0).Rarity = "common" }},
		{"weapon-1-2-or-3 plus-1 attunement", func(b *batch) { v(b, 1, 0).Attunement.Mode = "required" }},
		{"weapon-1-2-or-3 plus-2 rarity", func(b *batch) { v(b, 1, 1).Rarity = "common" }},
		{"weapon-1-2-or-3 plus-2 attunement", func(b *batch) { v(b, 1, 1).Attunement.Mode = "required" }},
		{"weapon-1-2-or-3 plus-3 rarity", func(b *batch) { v(b, 1, 2).Rarity = "common" }},
		{"weapon-1-2-or-3 plus-3 attunement", func(b *batch) { v(b, 1, 2).Attunement.Mode = "required" }},
		{"well-of-many-worlds rarity", func(b *batch) { v(b, 2, 0).Rarity = "common" }},
		{"well-of-many-worlds attunement", func(b *batch) { v(b, 2, 0).Attunement.Mode = "required" }},
		{"wind-fan rarity", func(b *batch) { v(b, 3, 0).Rarity = "common" }},
		{"wind-fan attunement", func(b *batch) { v(b, 3, 0).Attunement.Mode = "required" }},
		{"winged-boots rarity", func(b *batch) { v(b, 4, 0).Rarity = "common" }},
		{"winged-boots attunement", func(b *batch) { v(b, 4, 0).Attunement.Mode = "none" }},
		{"wings-of-flying rarity", func(b *batch) { v(b, 5, 0).Rarity = "common" }},
		{"wings-of-flying attunement", func(b *batch) { v(b, 5, 0).Attunement.Mode = "none" }},
		{"unknown rarity", func(b *batch) { v(b, 0, 0).Rarity = "mythic" }},
		{"attunement restriction", func(b *batch) { v(b, 0, 0).Attunement.Restrictions = []string{"wizard"} }},
		{"missing charges", func(b *batch) { v(b, 0, 0).Charges = nil }},
		{"invented charges", func(b *batch) { v(b, 0, 0).Charges.Maximum = 10 }},
		{"winged-boots dropped charges", func(b *batch) { *v(b, 4, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"winged-boots charge maximum", func(b *batch) { v(b, 4, 0).Charges.Maximum++ }},
		{"winged-boots charge recharge", func(b *batch) { v(b, 4, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"equip mode", func(b *batch) { v(b, 0, 0).Equip.Mode = "unreviewed" }},
		{"equip restriction", func(b *batch) { v(b, 0, 0).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, 0, 0).Use.Mode = "unreviewed" }},
		{"dropped use restriction", func(b *batch) { v(b, 2, 0).Use.Restrictions = []string{} }},
		{"listed price", func(b *batch) { v(b, 0, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0, 0).Cost.Amount = 500 }},
		{"cost currency", func(b *batch) { v(b, 0, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = v(b, 0, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0, 0).Details = &ItemDetails246{} }},
		{"mismatched payload", func(b *batch) { v(b, 0, 0).Details = v(b, 1, 0).Details }},
		{"second payload", func(b *batch) { v(b, 0, 0).Details.WeaponBonus = v(b, 1, 0).Details.WeaponBonus }},
		{"blank key", func(b *batch) { v(b, 0, 0).Details.Warning.Condition = " " }},
		{"prose value", func(b *batch) { v(b, 0, 0).Details.Warning.Condition = "Printed Prose" }},
		{"negative value", func(b *batch) { v(b, 0, 0).Details.Warning.RadiusFeet = -1 }},
		{"warning radius", func(b *batch) { p := v(b, 0, 0).Details.Warning; p.RadiusFeet = 60 }},
		{"warning magical sleep", func(b *batch) { p := v(b, 0, 0).Details.Warning; p.Alarm.WakesFromMagicalSleep = true }},
		{"warning roll", func(b *batch) { p := v(b, 0, 0).Details.Warning; p.SupernaturalReadiness.Roll = "attack" }},
		{"weapon attack bonus", func(b *batch) { p := v(b, 1, 0).Details.WeaponBonus; p.AttackBonus = 3 }},
		{"weapon damage bonus", func(b *batch) { p := v(b, 1, 0).Details.WeaponBonus; p.DamageBonus = 0 }},
		{"well diameter", func(b *batch) { p := v(b, 2, 0).Details.Well; p.Portal.DiameterFeet = 10 }},
		{"well close range", func(b *batch) { p := v(b, 2, 0).Details.Well; p.Close.WithinFeet = 10 }},
		{"well reuse", func(b *batch) { p := v(b, 2, 0).Details.Well; p.ReuseAfter.Sides = 12 }},
		{"fan dc", func(b *batch) { p := v(b, 3, 0).Details.WindFan; p.SaveDC = 15 }},
		{"fan failure", func(b *batch) { p := v(b, 3, 0).Details.WindFan; p.RepeatUse.FailurePercentCumulative = 10 }},
		{"boots charges", func(b *batch) { p := v(b, 4, 0).Details.WingedBoots; p.MaxCharges = 3 }},
		{"boots speed", func(b *batch) { p := v(b, 4, 0).Details.WingedBoots; p.FlySpeedFeet = 60 }},
		{"boots descent", func(b *batch) { p := v(b, 4, 0).Details.WingedBoots; p.OnExpiry.DescentFeetPerRound = 60 }},
		{"wings speed", func(b *batch) { p := v(b, 5, 0).Details.Wings; p.FlySpeedFeet = 30 }},
		{"wings duration", func(b *batch) { p := v(b, 5, 0).Details.Wings; p.DurationHours = 8 }},
		{"wings reuse", func(b *batch) { p := v(b, 5, 0).Details.Wings; p.ReuseAfter.Sides = 8 }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems246(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems246(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems246RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/weapon-of-warning/variant/standard"`)
	if !bytes.Contains(items246JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items246JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems246(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"warning": {`, `"warning_plus": {`},
		"unknown detail":     {`"radius_feet": 30`, `"radius_feet_all": 30`},
		"string for number":  {`"radius_feet": 30`, `"radius_feet": "30"`},
		"fractional number":  {`"radius_feet": 30`, `"radius_feet": 30.5`},
		"number for boolean": {`"wakes_from_magical_sleep": false`, `"wakes_from_magical_sleep": 1`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items246JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items246JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems246(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items246JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems246(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items246JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems245(items246JSON); err == nil {
		t.Fatal("#245 decoder accepted the #246 batch")
	}
	if _, err := DecodeItems246(items245JSON); err == nil {
		t.Fatal("#246 decoder accepted the #245 batch")
	}
}
