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

func loadItems239(t *testing.T) ItemSourceBatch[ItemDetails239] {
	t.Helper()
	b, err := LoadItems239()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Identities transcribed from the category/rarity line under each heading on
// pinned PDF pp. 238-239, in planning inventory order for issue #239.
func TestItems239SourceIdentities(t *testing.T) {
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
		{"Ring of Free Action", "ring-of-free-action", 238, "ring", nil, "required", noCharges, []variant{{"standard", "Ring of Free Action", "rare"}}},
		{"Ring of Invisibility", "ring-of-invisibility", 238, "ring", nil, "required", noCharges, []variant{{"standard", "Ring of Invisibility", "legendary"}}},
		{"Ring of Jumping", "ring-of-jumping", 238, "ring", nil, "required", noCharges, []variant{{"standard", "Ring of Jumping", "uncommon"}}},
		{"Ring of Mind Shielding", "ring-of-mind-shielding", 238, "ring", nil, "required", noCharges, []variant{{"standard", "Ring of Mind Shielding", "uncommon"}}},
		{"Ring of Protection", "ring-of-protection", 238, "ring", nil, "required", noCharges, []variant{{"standard", "Ring of Protection", "rare"}}},
		{"Ring of Regeneration", "ring-of-regeneration", 238, "ring", nil, "required", noCharges, []variant{{"standard", "Ring of Regeneration", "very-rare"}}},
		{"Ring of Resistance", "ring-of-resistance", 238, "ring", nil, "none", noCharges, []variant{{"standard", "Ring of Resistance", "rare"}}},
		{"Ring of Shooting Stars", "ring-of-shooting-stars", 238, "ring", nil, "required", ItemCharges{Mode: "limited", Maximum: 6, PerUse: 1, Recharge: "regain-1d6-at-dawn"}, []variant{{"standard", "Ring of Shooting Stars", "very-rare"}}},
		{"Ring of Spell Storing", "ring-of-spell-storing", 239, "ring", nil, "required", noCharges, []variant{{"standard", "Ring of Spell Storing", "rare"}}},
		{"Ring of Spell Turning", "ring-of-spell-turning", 239, "ring", nil, "required", noCharges, []variant{{"standard", "Ring of Spell Turning", "legendary"}}},
		{"Ring of Swimming", "ring-of-swimming", 239, "ring", nil, "none", noCharges, []variant{{"standard", "Ring of Swimming", "uncommon"}}},
		{"Ring of Telekinesis", "ring-of-telekinesis", 239, "ring", nil, "required", noCharges, []variant{{"standard", "Ring of Telekinesis", "very-rare"}}},
	}
	b := loadItems239(t)
	if len(b.Families) != 12 || b.Batch != "data-magic-items-15" || b.CatalogueIssue != 239 || b.SchemaVersion != 1 || b.RulesVersion != "srd-5.2.1" {
		t.Fatal("batch must hold exactly the twelve families of issue #239")
	}
	if b.Source.SHA256 != artifactSHA256 || b.Source.Attribution != Attribution || b.Source.License != "CC-BY-4.0" || !strings.Contains(b.Source.Changes, "pp. 238–239") {
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
	if total != 12 {
		t.Fatalf("expected 12 source-listed variants, got %d", total)
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
		if r.Issue != 239 {
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
	if index != 12 || len(members) != 12 {
		t.Fatalf("inventory assigns %d families to issue #239", len(members))
	}
	// The planning batch names the same members and pages, in the same order.
	planned := 0
	for _, batch := range inventory.Batches {
		if batch.Key != b.Batch {
			continue
		}
		planned++
		if batch.Pages != "238, 239" || len(batch.Members) != 12 {
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
	b.Families[0].Variants[0].Details.FreeAction.Condition = "changed"
	if fresh := loadItems239(t); fresh.Families[0].Variants[0].Details.FreeAction.Condition != "while-wearing" {
		t.Fatal("load must return independent data")
	}
}

// Literal source values for all 12 variants, transcribed from pinned PDF pp.
// 238-239.
func TestItems239SourceValues(t *testing.T) {
	expected := []string{
		`{"free_action":{"condition":"while-wearing","difficult_terrain_costs_extra_movement":false,"magic_cannot":["reduce-any-speed","cause-paralyzed-condition","cause-restrained-condition"]}}`,
		`{"invisibility":{"condition":"while-wearing","action":"magic","grants_condition":"invisible","ends_when":["ring-removed","bonus-action-to-become-visible"]}}`,
		`{"jumping":{"condition":"while-wearing","spell":"jump","target":"yourself-only"}}`,
		`{"mind_shielding":{"condition":"while-wearing","immune_to_magic_that":["reads-thoughts","detects-lies","reveals-alignment","reveals-creature-type"],"telepathy_needs_permission":true,"imperceptible":{"action":"magic","ends_when":["magic-action-to-make-perceptible","ring-removed","wearer-dies"]},"soul":{"trigger":"die-while-wearing","blocked_if":"ring-already-houses-a-soul","choices":["remain-in-ring","depart-for-afterlife"],"telepathy_while":"soul-is-in-the-ring","telepathy_with":"any-creature-wearing-the-ring","wearer_can_prevent":false}}}`,
		`{"protection":{"condition":"while-wearing","armor_class_bonus":1,"saving_throw_bonus":1}}`,
		`{"regeneration":{"condition":"while-wearing","healing":{"dice":1,"sides":6,"type":"hit-points"},"interval_minutes":10,"minimum_hit_points":1,"regrowth":{"trigger":"lose-a-body-part","result":"part-regrows-with-full-functionality","delay":{"dice":1,"sides":6,"type":"days"},"delay_bonus_days":1,"requires_hit_points_whole_time":true}}}`,
		`{"resistance_ring":{"condition":"while-wearing","defense":"resistance","type_shown_by":"gemstone","type_chosen_by":"gm-chooses-or-rolls-randomly","table_die_sides":10,"damage_types":[{"roll":1,"damage_type":"acid","gemstone":"pearl"},{"roll":2,"damage_type":"cold","gemstone":"tourmaline"},{"roll":3,"damage_type":"fire","gemstone":"garnet"},{"roll":4,"damage_type":"force","gemstone":"sapphire"},{"roll":5,"damage_type":"lightning","gemstone":"citrine"},{"roll":6,"damage_type":"necrotic","gemstone":"jet"},{"roll":7,"damage_type":"poison","gemstone":"amethyst"},{"roll":8,"damage_type":"psychic","gemstone":"jade"},{"roll":9,"damage_type":"radiant","gemstone":"topaz"},{"roll":10,"damage_type":"thunder","gemstone":"spinel"}]}}`,
		`{"shooting_stars":{"free_spells":["dancing-lights","light"],"max_charges":6,"recharge":{"dice":1,"sides":6,"type":"charges"},"recharge_when":"daily-at-dawn","faerie_fire":{"charge_cost":1,"spell":"faerie-fire"},"lightning_spheres":{"charge_cost":2,"action":"magic","max_spheres":4,"sphere_diameter_feet":3,"placement_range_feet":120,"concentration_required":true,"max_duration_minutes":1,"dim_light_radius_feet":30,"move_action":"bonus","move_feet":30,"leash_feet":120,"discharge_distance_feet":5,"discharge_excludes":["you","behind-total-cover"],"sphere_disappears_on_discharge":true,"save_ability":"dexterity","save_dc":15,"success_damage":"half","damage_table":[{"spheres":1,"damage":{"dice":4,"sides":12,"type":"lightning"}},{"spheres":2,"damage":{"dice":5,"sides":4,"type":"lightning"}},{"spheres":3,"damage":{"dice":2,"sides":6,"type":"lightning"}},{"spheres":4,"damage":{"dice":2,"sides":4,"type":"lightning"}}]},"shooting_stars":{"minimum_charges":1,"maximum_charges":3,"action":"magic","motes_per_charge":1,"range_feet":60,"area_shape":"cube","area_feet":15,"save_ability":"dexterity","save_dc":15,"damage":{"dice":5,"sides":4,"type":"radiant"},"success_damage":"half"},"source_pages":[238,239]}}`,
		`{"spell_storing":{"holds":"spells-cast-into-it-until-attuned-wearer-uses-them","capacity_spell_levels":5,"found_levels":{"dice":1,"sides":6,"type":"stored-spell-levels"},"found_levels_penalty":1,"found_spells_chosen_by":"gm","storing":{"caster":"any-creature","minimum_spell_level":1,"maximum_spell_level":5,"method":"touch-ring-as-spell-is-cast","other_effect":"none-besides-being-stored","if_it_cannot_fit":"spell-expended-without-effect","space_used":"slot-level-used"},"casting":{"condition":"while-wearing","uses_from_original_caster":["slot-level","spell-save-dc","spell-attack-bonus","spellcasting-ability"],"otherwise":"treated-as-if-you-cast-it","frees_space":true}}}`,
		`{"spell_turning":{"condition":"while-wearing","save_mode":"advantage","negate_max_spell_level":7,"negate_result":"spell-has-no-effect-on-you","deflect":{"requires":["spell-targeted-only-you","spell-created-no-area-of-effect"],"action":"reaction","target":"original-caster","caster_save_uses":"casters-own-spell-save-dc"}}}`,
		`{"swimming":{"condition":"while-wearing","swim_speed_feet":40}}`,
		`{"telekinesis":{"condition":"while-wearing","spell":"telekinesis"}}`,
	}
	// equip mode : use mode : use restrictions : required mechanic keys.
	requirements := []string{
		"wear:none::item.attunement item.wear movement.ignore-difficult-terrain-cost defense.magic-cannot-reduce-speed defense.magic-cannot-paralyze-or-restrain",
		"wear:magic-action::item.attunement item.wear action.magic condition.invisible effect.ends-on-removal-or-bonus-action",
		"wear:cast-spell:jump-targeting-yourself:item.attunement item.wear spell.cast-jump spell.target-self-only",
		"wear:none::item.attunement item.wear defense.immunity-mind-reading-magic communication.telepathy-only-if-allowed item.imperceptible-toggle item.soul-vessel-on-death communication.soul-telepathy",
		"wear:none::item.attunement item.wear defense.armor-class-bonus-1 save.all-bonus-1",
		"wear:none::item.attunement item.wear healing.1d6-every-10-minutes healing.regrow-body-part healing.requires-1-hit-point",
		"wear:none::item.wear defense.resistance-one-damage-type table.damage-type-and-gemstone-1d10",
		"wear:cast-spell:expend-listed-charges:item.attunement item.wear item.charges spell.cast-dancing-lights spell.cast-light spell.cast-faerie-fire action.magic-lightning-spheres save.dexterity-dc-15 table.sphere-damage action.magic-shooting-stars damage.radiant-5d4-per-charge",
		"wear:cast-spell:spell-stored-in-ring:item.attunement item.wear item.store-spells item.capacity-5-spell-levels spell.cast-stored-spell spell.uses-original-caster-statistics",
		"wear:reaction:successful-save-against-qualifying-spell:item.attunement item.wear save.advantage-against-spells spell.negate-level-7-or-lower-on-successful-save action.reaction-deflect-spell",
		"wear:none::item.wear movement.swim-speed-40-feet",
		"wear:cast-spell:telekinesis:item.attunement item.wear spell.cast-telekinesis",
	}
	index := 0
	for _, f := range loadItems239(t).Families {
		for _, v := range f.Variants {
			t.Run(v.Name, func(t *testing.T) {
				var want ItemDetails239
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

func TestItems239CompositionWithEarlierBatches(t *testing.T) {
	batch := loadItems239(t)
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
	if len(composed) != 180 || variants != 241 || composed[0].Name != "Potion of Healing" || composed[168].Name != "Ring of Free Action" || composed[179].Name != "Ring of Telekinesis" {
		t.Fatalf("expected Potion of Healing, batches #225 to #238 and the twelve families of #239: %d families, %d variants", len(composed), variants)
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
			if strings.Contains(fixture, "items_239") {
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
func TestItems239RejectMissingFields(t *testing.T) {
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
	if err := json.Unmarshal(items239JSON, &root); err != nil {
		t.Fatal(err)
	}
	walk(root, nil)
	if len(paths) < 600 {
		t.Fatalf("expected the whole batch to be walked, got %d members", len(paths))
	}
	for _, path := range paths {
		var fresh any
		if err := json.Unmarshal(items239JSON, &fresh); err != nil {
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
		if _, err := DecodeItems239(data); err == nil {
			t.Fatalf("missing member accepted: %v", path)
		}
	}
}

func TestItems239RejectInvalidFixtures(t *testing.T) {
	type batch = ItemSourceBatch[ItemDetails239]
	// v addresses one variant of a family by index; see the source identity test.
	v := func(b *batch, family, variant int) *ItemSourceVariant[ItemDetails239] {
		return &b.Families[family].Variants[variant]
	}
	tests := []struct {
		name   string
		mutate func(*batch)
	}{
		{"schema", func(b *batch) { b.SchemaVersion = 2 }},
		{"rules version", func(b *batch) { b.RulesVersion = "legacy" }},
		{"batch key", func(b *batch) { b.Batch = "data-magic-items-14" }},
		{"issue", func(b *batch) { b.CatalogueIssue = 238 }},
		{"source hash", func(b *batch) { b.Source.SHA256 = "wrong" }},
		{"source attribution", func(b *batch) { b.Source.Attribution = "" }},
		{"adaptation notice", func(b *batch) { b.Source.Changes = " " }},
		{"source license", func(b *batch) { b.Source.License = "MIT" }},
		{"missing family", func(b *batch) { b.Families = b.Families[:11] }},
		{"reordered family", func(b *batch) { b.Families[0], b.Families[1] = b.Families[1], b.Families[0] }},
		{"duplicate family", func(b *batch) { b.Families[1] = b.Families[0] }},
		{"adjacent family", func(b *batch) {
			b.Families = append(b.Families, b.Families[0])
			b.Families[12].ID = "srd-5.2.1:content/items/ring-of-three-wishes"
			b.Families[12].Name = "Ring of Three Wishes"
		}},
		{"unknown family", func(b *batch) { b.Families[0].ID = "srd-5.2.1:content/items/unreviewed" }},
		{"family name", func(b *batch) { b.Families[0].Name = "Unreviewed Name" }},
		{"mundane kind", func(b *batch) { b.Families[0].Kind = "mundane" }},
		{"category", func(b *batch) { b.Families[0].Category = "armor" }},
		{"missing base", func(b *batch) { b.Families[0].Base = nil }},
		{"invented base", func(b *batch) { b.Families[0].Base.Options = []string{"any-armor"} }},
		{"source page", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{239, 239} }},
		{"page range", func(b *batch) { b.Families[0].Presentations[0].Pages = []int{238, 239} }},
		{"missing presentation", func(b *batch) { b.Families[0].Presentations = nil }},
		{"foreign presentation", func(b *batch) { b.Families[0].Presentations[0].ID = b.Families[1].Presentations[0].ID }},
		{"invented variant", func(b *batch) { b.Families[0].Variants = append(b.Families[0].Variants, b.Families[0].Variants[0]) }},
		{"missing variant", func(b *batch) { b.Families[0].Variants = nil }},
		{"variant id", func(b *batch) { v(b, 0, 0).ID = "" }},
		{"variant name", func(b *batch) { v(b, 0, 0).Name = "Unreviewed" }},
		{"ring-of-free-action rarity", func(b *batch) { v(b, 0, 0).Rarity = "common" }},
		{"ring-of-invisibility rarity", func(b *batch) { v(b, 1, 0).Rarity = "common" }},
		{"ring-of-jumping rarity", func(b *batch) { v(b, 2, 0).Rarity = "common" }},
		{"ring-of-mind-shielding rarity", func(b *batch) { v(b, 3, 0).Rarity = "common" }},
		{"ring-of-protection rarity", func(b *batch) { v(b, 4, 0).Rarity = "common" }},
		{"ring-of-regeneration rarity", func(b *batch) { v(b, 5, 0).Rarity = "common" }},
		{"ring-of-resistance rarity", func(b *batch) { v(b, 6, 0).Rarity = "common" }},
		{"ring-of-shooting-stars rarity", func(b *batch) { v(b, 7, 0).Rarity = "common" }},
		{"ring-of-spell-storing rarity", func(b *batch) { v(b, 8, 0).Rarity = "common" }},
		{"ring-of-spell-turning rarity", func(b *batch) { v(b, 9, 0).Rarity = "common" }},
		{"ring-of-swimming rarity", func(b *batch) { v(b, 10, 0).Rarity = "common" }},
		{"ring-of-telekinesis rarity", func(b *batch) { v(b, 11, 0).Rarity = "common" }},
		{"unknown rarity", func(b *batch) { v(b, 0, 0).Rarity = "mythic" }},
		{"ring-of-free-action attunement", func(b *batch) { v(b, 0, 0).Attunement.Mode = "none" }},
		{"ring-of-invisibility attunement", func(b *batch) { v(b, 1, 0).Attunement.Mode = "none" }},
		{"ring-of-jumping attunement", func(b *batch) { v(b, 2, 0).Attunement.Mode = "none" }},
		{"ring-of-mind-shielding attunement", func(b *batch) { v(b, 3, 0).Attunement.Mode = "none" }},
		{"ring-of-protection attunement", func(b *batch) { v(b, 4, 0).Attunement.Mode = "none" }},
		{"ring-of-regeneration attunement", func(b *batch) { v(b, 5, 0).Attunement.Mode = "none" }},
		{"ring-of-resistance attunement", func(b *batch) { v(b, 6, 0).Attunement.Mode = "required" }},
		{"ring-of-shooting-stars attunement", func(b *batch) { v(b, 7, 0).Attunement.Mode = "none" }},
		{"ring-of-spell-storing attunement", func(b *batch) { v(b, 8, 0).Attunement.Mode = "none" }},
		{"ring-of-spell-turning attunement", func(b *batch) { v(b, 9, 0).Attunement.Mode = "none" }},
		{"ring-of-swimming attunement", func(b *batch) { v(b, 10, 0).Attunement.Mode = "required" }},
		{"ring-of-telekinesis attunement", func(b *batch) { v(b, 11, 0).Attunement.Mode = "none" }},
		{"attunement restriction", func(b *batch) { v(b, 0, 0).Attunement.Restrictions = []string{"wizard"} }},
		{"missing charges", func(b *batch) { v(b, 0, 0).Charges = nil }},
		{"invented charges", func(b *batch) { v(b, 0, 0).Charges.Maximum = 10 }},
		{"ring-of-shooting-stars dropped charges", func(b *batch) { *v(b, 7, 0).Charges = ItemCharges{Mode: "none", Recharge: "none"} }},
		{"ring-of-shooting-stars charge maximum", func(b *batch) { v(b, 7, 0).Charges.Maximum++ }},
		{"ring-of-shooting-stars charge recharge", func(b *batch) { v(b, 7, 0).Charges.Recharge = "regain-all-at-dawn" }},
		{"equip mode", func(b *batch) { v(b, 0, 0).Equip.Mode = "hold" }},
		{"equip restriction", func(b *batch) { v(b, 0, 0).Equip.Restrictions = []string{"proficiency"} }},
		{"use mode", func(b *batch) { v(b, 0, 0).Use.Mode = "reaction" }},
		{"dropped use restriction", func(b *batch) { v(b, 2, 0).Use.Restrictions = []string{} }},
		{"listed price", func(b *batch) { v(b, 0, 0).Cost.Kind = "listed" }},
		{"inferred price", func(b *batch) { v(b, 0, 0).Cost.Amount = 500 }},
		{"cost currency", func(b *batch) { v(b, 0, 0).Cost.Currency = "gp" }},
		{"cost page", func(b *batch) { v(b, 0, 0).Cost.SourcePages = []int{99, 99} }},
		{"missing mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = v(b, 0, 0).RequiredMechanicKeys[:1] }},
		{"no mechanics", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys = []string{} }},
		{"unknown mechanic", func(b *batch) { v(b, 0, 0).RequiredMechanicKeys[0] = "item.grant-ac" }},
		{"missing details", func(b *batch) { v(b, 0, 0).Details = nil }},
		{"empty details", func(b *batch) { v(b, 0, 0).Details = &ItemDetails239{} }},
		{"mismatched payload", func(b *batch) { v(b, 0, 0).Details = v(b, 1, 0).Details }},
		{"second payload", func(b *batch) { v(b, 0, 0).Details.Invisibility = v(b, 1, 0).Details.Invisibility }},
		{"blank key", func(b *batch) { v(b, 0, 0).Details.FreeAction.Condition = " " }},
		{"prose value", func(b *batch) { v(b, 0, 0).Details.FreeAction.Condition = "Printed Prose" }},
		{"negative value", func(b *batch) { v(b, 4, 0).Details.Protection.ArmorClassBonus = -1 }},
		{"free action terrain", func(b *batch) { p := v(b, 0, 0).Details.FreeAction; p.DifficultTerrainCostsExtraMovement = true }},
		{"free action duplicate", func(b *batch) { p := v(b, 0, 0).Details.FreeAction; p.MagicCannot[1] = "reduce-any-speed" }},
		{"free action dropped", func(b *batch) { p := v(b, 0, 0).Details.FreeAction; p.MagicCannot = p.MagicCannot[:1] }},
		{"invisibility action", func(b *batch) { p := v(b, 1, 0).Details.Invisibility; p.Action = "bonus" }},
		{"invisibility condition", func(b *batch) { p := v(b, 1, 0).Details.Invisibility; p.GrantsCondition = "hidden" }},
		{"invisibility ends", func(b *batch) { p := v(b, 1, 0).Details.Invisibility; p.EndsWhen = p.EndsWhen[:1] }},
		{"jumping spell", func(b *batch) { p := v(b, 2, 0).Details.Jumping; p.Spell = "fly" }},
		{"jumping target", func(b *batch) { p := v(b, 2, 0).Details.Jumping; p.Target = "any-creature" }},
		{"mind shielding immunity", func(b *batch) { p := v(b, 3, 0).Details.MindShielding; p.ImmuneToMagicThat = p.ImmuneToMagicThat[:3] }},
		{"mind shielding telepathy", func(b *batch) { p := v(b, 3, 0).Details.MindShielding; p.TelepathyNeedsPermission = false }},
		{"mind shielding soul", func(b *batch) { p := v(b, 3, 0).Details.MindShielding; p.Soul.WearerCanPrevent = true }},
		{"mind shielding duplicate end", func(b *batch) {
			p := v(b, 3, 0).Details.MindShielding
			p.Imperceptible.EndsWhen[1] = "magic-action-to-make-perceptible"
		}},
		{"protection armor", func(b *batch) { p := v(b, 4, 0).Details.Protection; p.ArmorClassBonus = 2 }},
		{"protection saves", func(b *batch) { p := v(b, 4, 0).Details.Protection; p.SavingThrowBonus = 0 }},
		{"regeneration dice", func(b *batch) { p := v(b, 5, 0).Details.Regeneration; p.Healing.Sides = 4 }},
		{"regeneration interval", func(b *batch) { p := v(b, 5, 0).Details.Regeneration; p.IntervalMinutes = 1 }},
		{"regeneration minimum", func(b *batch) { p := v(b, 5, 0).Details.Regeneration; p.MinimumHitPoints = 0 }},
		{"regeneration delay", func(b *batch) { p := v(b, 5, 0).Details.Regeneration; p.Regrowth.DelayBonusDays = 0 }},
		{"resistance die", func(b *batch) { p := v(b, 6, 0).Details.ResistanceRing; p.TableDieSides = 8 }},
		{"resistance missing row", func(b *batch) { p := v(b, 6, 0).Details.ResistanceRing; p.DamageTypes = p.DamageTypes[:9] }},
		{"resistance duplicate type", func(b *batch) { p := v(b, 6, 0).Details.ResistanceRing; p.DamageTypes[1].DamageType = "acid" }},
		{"resistance duplicate gem", func(b *batch) { p := v(b, 6, 0).Details.ResistanceRing; p.DamageTypes[1].Gemstone = "pearl" }},
		{"resistance row overlap", func(b *batch) { p := v(b, 6, 0).Details.ResistanceRing; p.DamageTypes[4].Roll = 4 }},
		{"stars charges", func(b *batch) { p := v(b, 7, 0).Details.ShootingStars; p.MaxCharges = 5 }},
		{"stars sphere dc", func(b *batch) { p := v(b, 7, 0).Details.ShootingStars; p.LightningSpheres.SaveDC = 13 }},
		{"stars sphere cost", func(b *batch) { p := v(b, 7, 0).Details.ShootingStars; p.LightningSpheres.ChargeCost = 1 }},
		{"stars sphere row", func(b *batch) {
			p := v(b, 7, 0).Details.ShootingStars
			p.LightningSpheres.DamageTable = p.LightningSpheres.DamageTable[:3]
		}},
		{"stars damage", func(b *batch) { p := v(b, 7, 0).Details.ShootingStars; p.ShootingStars.Damage.Sides = 6 }},
		{"stars pages", func(b *batch) { p := v(b, 7, 0).Details.ShootingStars; p.SourcePages = []int{238} }},
		{"storing capacity", func(b *batch) { p := v(b, 8, 0).Details.SpellStoring; p.CapacitySpellLevels = 9 }},
		{"storing found dice", func(b *batch) { p := v(b, 8, 0).Details.SpellStoring; p.FoundLevels.Sides = 4 }},
		{"storing max level", func(b *batch) { p := v(b, 8, 0).Details.SpellStoring; p.Storing.MaximumSpellLevel = 9 }},
		{"storing frees space", func(b *batch) { p := v(b, 8, 0).Details.SpellStoring; p.Casting.FreesSpace = false }},
		{"turning level", func(b *batch) { p := v(b, 9, 0).Details.SpellTurning; p.NegateMaxSpellLevel = 9 }},
		{"turning mode", func(b *batch) { p := v(b, 9, 0).Details.SpellTurning; p.SaveMode = "normal" }},
		{"turning action", func(b *batch) { p := v(b, 9, 0).Details.SpellTurning; p.Deflect.Action = "bonus" }},
		{"swimming speed", func(b *batch) { p := v(b, 10, 0).Details.Swimming; p.SwimSpeedFeet = 60 }},
		{"telekinesis spell", func(b *batch) { p := v(b, 11, 0).Details.Telekinesis; p.Spell = "levitate" }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := loadItems239(t)
			tc.mutate(&b)
			data, err := json.Marshal(b)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := DecodeItems239(data); err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestItems239RejectUnsupportedJSON(t *testing.T) {
	anchor := []byte(`"id": "srd-5.2.1:content/items/ring-of-free-action/variant/standard"`)
	if !bytes.Contains(items239JSON, anchor) {
		t.Fatal("anchor variant missing")
	}
	for _, field := range []string{`"quantity":10,`, `"ac":99,`, `"actions":["attack"],`, `"resources":{"slots":99},`, `"playable":true,`, `"available":true,`, `"properties":{"ac_bonus":9},`, `"healing":{"dice":2,"sides":4,"bonus":2},`, `"id":"hidden-duplicate",`, `"ID":"hidden-alias",`} {
		t.Run(field, func(t *testing.T) {
			data := bytes.Replace(items239JSON, anchor, append([]byte(field), anchor...), 1)
			if _, err := DecodeItems239(data); err == nil {
				t.Fatal("unsupported runtime grant/property accepted")
			}
		})
	}
	for name, change := range map[string][2]string{
		"null details":       {`"details": {`, `"details": null, "x": {`},
		"null list":          {`"excluded": []`, `"excluded": null`},
		"unknown payload":    {`"free_action": {`, `"free_action_plus": {`},
		"unknown detail":     {`"maximum": 0`, `"maximum_all": 0`},
		"string for number":  {`"maximum": 0`, `"maximum": "0"`},
		"fractional number":  {`"maximum": 0`, `"maximum": 0.5`},
		"number for boolean": {`"telepathy_needs_permission": true`, `"telepathy_needs_permission": 1`},
	} {
		t.Run(name, func(t *testing.T) {
			if !bytes.Contains(items239JSON, []byte(change[0])) {
				t.Fatal("replacement anchor missing")
			}
			data := bytes.Replace(items239JSON, []byte(change[0]), []byte(change[1]), 1)
			if _, err := DecodeItems239(data); err == nil {
				t.Fatal("ambiguous or unsupported JSON accepted")
			}
		})
	}
	for _, data := range [][]byte{append(slices.Clone(items239JSON), []byte(` {}`)...), []byte(`null`), []byte(`[]`), []byte(`{"schema_version":1,"schema_version":1}`)} {
		if _, err := DecodeItems239(data); err == nil {
			t.Fatal("ambiguous JSON accepted")
		}
	}
	// Neither the Healing contract decoder nor an earlier batch decoder accepts this batch.
	if _, err := DecodeItemCatalogue(items239JSON); err == nil {
		t.Fatal("contract decoder accepted a source batch")
	}
	if _, err := DecodeItems238(items239JSON); err == nil {
		t.Fatal("#238 decoder accepted the #239 batch")
	}
	if _, err := DecodeItems239(items238JSON); err == nil {
		t.Fatal("#239 decoder accepted the #238 batch")
	}
}
